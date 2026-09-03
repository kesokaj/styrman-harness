package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/trail"
)

// Grade evaluates the outcome of an agent's work in a workspace.
func Grade(ctx context.Context, task Task, workspaceDir string, cfg Config) (EvalRunResult, error) {
	res := EvalRunResult{
		TaskID:    task.ID,
		TaskName:  task.Name,
		Category:  task.Category,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// 1. Check Git Cleanliness
	clean, gitErr := isGitClean(workspaceDir)
	res.CleanGit = clean && gitErr == nil

	// 2. Check for Prohibited Files (secrets)
	if hasForbiddenFiles(workspaceDir) {
		res.Verdict = "FAIL"
		res.Detail = "prohibited secret or credential file found in workspace"
		return res, nil
	}

	// 3. Check Trail Freshness (if opted in)
	if proj, ok := trail.Resolve(workspaceDir); ok && proj.Opted() {
		rep := proj.CheckFreshness()
		res.TrailFresh = !rep.Stale
		res.DoDCompliant = proj.HasArtifact(trail.RowChangelog, proj.Scope())
	} else {
		res.TrailFresh = true
		res.DoDCompliant = true
	}

	// 4. Run Task Verification
	if task.Verify != nil {
		pass, detail, err := task.Verify(workspaceDir)
		if err != nil {
			res.Error = err.Error()
		}
		res.Detail = detail
		if pass {
			res.BuildSuccess = true
			res.TestPassRate = 1.0
		} else {
			res.TestPassRate = 0.0
		}
	} else {
		// Default build check
		buildCmd := exec.CommandContext(ctx, "go", "build", "./...")
		buildCmd.Dir = workspaceDir
		if out, err := buildCmd.CombinedOutput(); err == nil {
			res.BuildSuccess = true
			res.TestPassRate = 1.0
		} else {
			res.BuildSuccess = false
			res.Detail = string(out)
		}
	}

	// 5. LLM-as-a-Judge (Vertex AI / Gemini)
	judgeScore, judgeFeedback := runLLMJudge(ctx, task, workspaceDir, cfg)
	res.JudgeScore = judgeScore
	res.JudgeFeedback = judgeFeedback

	// 6. Synthesize Overall Score
	deterministicWeight := 0.70
	judgeWeight := 0.30

	detScore := 0.0
	if res.BuildSuccess {
		detScore += 40.0
	}
	detScore += res.TestPassRate * 40.0
	if res.CleanGit {
		detScore += 10.0
	}
	if res.TrailFresh {
		detScore += 10.0
	}

	if judgeScore > 0 {
		res.OverallScore = (detScore * deterministicWeight) + (judgeScore * judgeWeight)
	} else {
		res.OverallScore = detScore
	}

	switch {
	case res.OverallScore >= 80 && res.BuildSuccess && res.TestPassRate >= 1.0:
		res.Verdict = "PASS"
	case res.BuildSuccess || res.OverallScore >= 40:
		res.Verdict = "PARTIAL"
	default:
		res.Verdict = "FAIL"
	}

	return res, nil
}

func isGitClean(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(out)) == 0, nil
}

func hasForbiddenFiles(dir string) bool {
	forbidden := []string{".env", ".env.local", "id_rsa", "id_ed25519", "credentials.json"}
	for _, f := range forbidden {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			return true
		}
	}
	return false
}

// runLLMJudge queries Vertex AI (Gemini) to evaluate code quality against the task rubric.
func runLLMJudge(ctx context.Context, task Task, workspaceDir string, cfg Config) (float64, string) {
	if cfg.ProjectID == "" {
		return 0, "LLM judge skipped: GCP project not configured"
	}

	diffCmd := exec.CommandContext(ctx, "git", "diff", "HEAD~1")
	diffCmd.Dir = workspaceDir
	diffOut, err := diffCmd.CombinedOutput()
	diff := string(diffOut)
	if err != nil || strings.TrimSpace(diff) == "" {
		// Try uncommitted diff
		diffCmd = exec.CommandContext(ctx, "git", "diff")
		diffCmd.Dir = workspaceDir
		diffOut, _ = diffCmd.CombinedOutput()
		diff = string(diffOut)
	}

	// Fetch ADC access token
	tokenCmd := exec.CommandContext(ctx, "gcloud", "auth", "application-default", "print-access-token")
	tokenOut, err := tokenCmd.Output()
	if err != nil {
		// Try gcloud print-access-token
		tokenCmd = exec.CommandContext(ctx, "gcloud", "auth", "print-access-token")
		tokenOut, err = tokenCmd.Output()
		if err != nil {
			return 0, "LLM judge skipped: unable to obtain GCP access token"
		}
	}
	token := strings.TrimSpace(string(tokenOut))

	prompt := fmt.Sprintf(
		"You are an expert software engineer evaluating an AI agent's coding task.\n"+
			"Task: %s\nPrompt: %s\nRubric:\n%s\n\n"+
			"The agent submitted the following git diff:\n```\n%s\n```\n\n"+
			"Evaluate the implementation against the rubric. "+
			"Respond ONLY with a valid JSON object in this format:\n"+
			"{\"score\": <number 0 to 100>, \"feedback\": \"<concise explanation>\"}",
		task.Name, task.Prompt, task.Rubric, truncateDiff(diff, 6000),
	)

	reqBody, _ := json.Marshal(map[string]any{
		"contents": []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]string{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":      0.1,
			"responseMimeType": "application/json",
		},
	})

	url := fmt.Sprintf(
		"https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:generateContent",
		cfg.Region, cfg.ProjectID, cfg.Region, cfg.VertexModel,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Sprintf("request error: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Sprintf("call error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Sprintf("vertex ai status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return 0, fmt.Sprintf("decode error: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return 0, "empty response from judge"
	}

	rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text
	var judgeResult struct {
		Score    float64 `json:"score"`
		Feedback string  `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(rawJSON), &judgeResult); err != nil {
		return 0, fmt.Sprintf("invalid judge JSON: %s", rawJSON)
	}

	return judgeResult.Score, judgeResult.Feedback
}

func truncateDiff(diff string, maxChars int) string {
	if len(diff) <= maxChars {
		return diff
	}
	return diff[:maxChars] + "\n... [diff truncated]"
}
