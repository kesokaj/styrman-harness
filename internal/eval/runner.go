package eval

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Runner orchestrates the execution of tasks in isolated workspaces.
type Runner struct {
	RepoRoot    string
	WorkDirRoot string
	Config      Config
	Storage     *Storage
}

// NewRunner creates an execution runner.
func NewRunner(repoRoot string, cfg Config) *Runner {
	return &Runner{
		RepoRoot: repoRoot,
		Config:   cfg,
		Storage:  NewStorage(repoRoot, cfg),
	}
}

// RunOptions configures an evaluation run.
type RunOptions struct {
	TaskID   string
	Vendor   string // "claude" | "agy" | "mock"
	Baseline bool   // Run without harness active
	Timeout  time.Duration
}

// Run executes a single task.
func (r *Runner) Run(ctx context.Context, opts RunOptions) (EvalRunResult, error) {
	task, ok := GetTask(opts.TaskID)
	if !ok {
		return EvalRunResult{}, fmt.Errorf("task %q not found", opts.TaskID)
	}

	runID := randomID(8)
	mode := "harnessed"
	if opts.Baseline {
		mode = "baseline"
	}

	baseDir := r.WorkDirRoot
	if baseDir == "" {
		baseDir = filepath.Join(os.TempDir(), "eval-runs")
	}
	workspaceDir := filepath.Join(baseDir, runID, task.ID)
	if err := os.MkdirAll(workspaceDir, 0o755); err != nil {
		return EvalRunResult{}, err
	}

	// 1. Initialize Git Workspace & Seed Files
	if err := r.initWorkspace(workspaceDir, task, opts.Baseline); err != nil {
		return EvalRunResult{}, fmt.Errorf("failed to init workspace: %w", err)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 10 * time.Minute
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Adapt prompt for baseline mode (strip harness-specific requirements)
	prompt := task.Prompt
	if opts.Baseline {
		prompt = strings.ReplaceAll(prompt, "Satisfy the project Definition of Done and commit your work cleanly.", "Commit your work cleanly.")
		prompt = strings.ReplaceAll(prompt, "Follow the Resume Protocol: inspect TODO.md and recent trail artifacts.\n", "")
		prompt = strings.ReplaceAll(prompt, "Satisfy the Definition of Done (run `harness dispatch` or record required artifacts) and commit cleanly.", "Commit cleanly.")
	}

	// 2. Launch Agent
	start := time.Now()
	telemetry, runErr := r.executeAgent(runCtx, workspaceDir, opts.Vendor, prompt, opts.Baseline)
	duration := time.Since(start)

	// 3. Grade the Outcome
	gradeResult, gradeErr := Grade(runCtx, task, workspaceDir, r.Config)
	if gradeErr != nil {
		gradeResult.Error = gradeErr.Error()
	}

	// 4. Populate Telemetry & Metadata
	gradeResult.RunID = runID
	gradeResult.Vendor = opts.Vendor
	gradeResult.Mode = mode
	gradeResult.Duration = duration
	gradeResult.DurationSec = duration.Seconds()
	gradeResult.Turns = telemetry.Turns
	gradeResult.ToolCalls = telemetry.ToolCalls
	gradeResult.TotalTokens = telemetry.TotalTokens
	gradeResult.InputTokens = telemetry.InputTokens
	gradeResult.OutputTokens = telemetry.OutputTokens
	gradeResult.CacheTokens = telemetry.CacheTokens
	gradeResult.CostUSD = telemetry.CostUSD
	gradeResult.WorkspaceDir = workspaceDir

	if runErr != nil && gradeResult.Error == "" {
		gradeResult.Error = runErr.Error()
	}

	// 5. Save to Ledger & GCS
	_ = r.Storage.SaveResult(ctx, gradeResult)

	return gradeResult, nil
}

func (r *Runner) initWorkspace(dir string, task Task, baseline bool) error {
	gitInit := exec.Command("git", "init")
	gitInit.Dir = dir
	if err := gitInit.Run(); err != nil {
		return err
	}

	// Set local author to satisfy git rules
	_ = exec.Command("git", "-C", dir, "config", "user.name", "eval-agent").Run()
	_ = exec.Command("git", "-C", dir, "config", "user.email", "eval@agent.local").Run()

	for relPath, content := range task.InitialFiles {
		fullPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			return err
		}
	}

	if baseline {
		vanillaAgentsMD := "# Project Rules\n\nBuild the requested application cleanly and directly using the standard Go library.\nWrite unit tests and commit your work cleanly to git when done.\nDo not create any .agents directories or extra harness tracking files.\n"
		_ = os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte(vanillaAgentsMD), 0o644)
	}

	// Initial commit
	_ = exec.Command("git", "-C", dir, "add", ".").Run()
	_ = exec.Command("git", "-C", dir, "commit", "-m", "initial repository setup").Run()
	return nil
}

func (r *Runner) executeAgent(ctx context.Context, dir, vendor, prompt string, baseline bool) (TelemetryData, error) {
	switch vendor {
	case "mock":
		return r.executeMock(dir, prompt)
	case "claude":
		return r.executeClaude(ctx, dir, prompt, baseline)
	case "agy":
		return r.executeAgy(ctx, dir, prompt, baseline)
	default:
		return TelemetryData{}, fmt.Errorf("unknown vendor: %s (supported: claude, agy, mock)", vendor)
	}
}

func (r *Runner) executeClaude(ctx context.Context, dir, prompt string, baseline bool) (TelemetryData, error) {
	enrichPath()
	args := []string{"-p", prompt, "--output-format", "json", "--dangerously-skip-permissions"}
	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = dir
	cmd.WaitDelay = 5 * time.Second

	out, err := cmd.CombinedOutput()
	telem := ParseEnvelopeJSON(out)

	// Also check for transcript logs
	if tFile := FindTranscriptInWorkspace(dir); tFile != "" {
		if fileTelem, parseErr := ParseTranscriptJSONL(tFile); parseErr == nil {
			if fileTelem.Turns > telem.Turns {
				telem.Turns = fileTelem.Turns
			}
			if len(fileTelem.ToolCalls) > 0 {
				telem.ToolCalls = fileTelem.ToolCalls
			}
		}
	}

	return telem, err
}

func (r *Runner) executeAgy(ctx context.Context, dir, prompt string, baseline bool) (TelemetryData, error) {
	enrichPath()
	args := []string{"-p", prompt, "--add-dir", dir, "--output-format", "json", "--dangerously-skip-permissions"}
	cmd := exec.CommandContext(ctx, "agy", args...)
	cmd.Dir = dir
	cmd.WaitDelay = 5 * time.Second

	out, err := cmd.CombinedOutput()
	telem := ParseEnvelopeJSON(out)

	if tFile := FindTranscriptInWorkspace(dir); tFile != "" {
		if fileTelem, parseErr := ParseTranscriptJSONL(tFile); parseErr == nil {
			telem.Turns = fileTelem.Turns
			telem.ToolCalls = fileTelem.ToolCalls
			if fileTelem.TotalTokens > 0 {
				telem.TotalTokens = fileTelem.TotalTokens
				telem.InputTokens = fileTelem.InputTokens
				telem.OutputTokens = fileTelem.OutputTokens
				telem.CacheTokens = fileTelem.CacheTokens
			}
		}
	}

	return telem, err
}

func enrichPath() {
	home, _ := os.UserHomeDir()
	extra := []string{
		filepath.Join(home, ".local", "bin"),
		filepath.Join(home, ".agents", "bin"),
		filepath.Join(home, ".go", "bin"),
		filepath.Join(home, "go", "bin"),
		"/usr/local/go/bin",
		"/opt/homebrew/bin",
		os.Getenv("PATH"),
	}
	os.Setenv("PATH", strings.Join(extra, string(os.PathListSeparator)))
}

// executeMock simulates an agent for quick verification testing without external LLM calls.
func (r *Runner) executeMock(dir, prompt string) (TelemetryData, error) {
	// Seed a valid working CLI solution
	mainGo := `package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type LogEntry struct {
	Status int ` + "`json:\"status\"`" + `
	Path   string ` + "`json:\"path\"`" + `
}

func main() {
	var reader io.Reader = os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err == nil {
			defer f.Close()
			reader = f
		}
	}

	counts := make(map[int]int)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err == nil && entry.Status > 0 {
			counts[entry.Status]++
		}
	}

	for status, count := range counts {
		fmt.Printf("Status %d: %d\n", status, count)
	}
}
`
	_ = os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0o644)

	mainTest := `package main

import "testing"

func TestDummy(t *testing.T) {
	// Passes
}
`
	_ = os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(mainTest), 0o644)

	// Commit work
	_ = exec.Command("git", "-C", dir, "add", ".").Run()
	_ = exec.Command("git", "-C", dir, "commit", "-m", "feat: implement log aggregator").Run()

	return TelemetryData{
		Turns: 3,
		ToolCalls: map[string]int{
			"write_to_file": 2,
			"run_command":   2,
		},
		TotalTokens:  1500,
		InputTokens:  1200,
		OutputTokens: 300,
		CostUSD:      0.0045,
	}, nil
}

func randomID(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
