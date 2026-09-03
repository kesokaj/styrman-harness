package eval

import "time"

// TaskCategory classifies evaluation scenarios.
type TaskCategory string

const (
	CategoryGreenfield TaskCategory = "greenfield"
	CategoryBrownfield TaskCategory = "brownfield"
)

// Task defines an evaluation challenge for an agent.
type Task struct {
	ID           string                                          `json:"id"`
	Name         string                                          `json:"name"`
	Category     TaskCategory                                    `json:"category"`
	Description  string                                          `json:"description"`
	Prompt       string                                          `json:"prompt"`
	InitialFiles map[string]string                               `json:"initialFiles,omitempty"`
	Verify       func(workspaceDir string) (bool, string, error) `json:"-"`
	Rubric       string                                          `json:"rubric"`
}

// EvalRunResult captures the comprehensive telemetry and grading for one task execution.
type EvalRunResult struct {
	RunID        string         `json:"runId"`
	TaskID       string         `json:"taskId"`
	TaskName     string         `json:"taskName"`
	Category     TaskCategory   `json:"category"`
	Vendor       string         `json:"vendor"` // "claude" | "agy" | "mock"
	Mode         string         `json:"mode"`   // "harnessed" | "baseline"
	Timestamp    string         `json:"timestamp"`
	Duration     time.Duration  `json:"duration"`
	DurationSec  float64        `json:"durationSec"`
	Turns        int            `json:"turns"`
	ToolCalls    map[string]int `json:"toolCalls"`
	TotalTokens  int            `json:"totalTokens"`
	InputTokens  int            `json:"inputTokens"`
	OutputTokens int            `json:"outputTokens"`
	CacheTokens  int            `json:"cacheTokens"`
	CostUSD      float64        `json:"costUsd"`

	// Deterministic outcomes
	BuildSuccess bool    `json:"buildSuccess"`
	TestPassRate float64 `json:"testPassRate"` // 0.0 - 1.0
	DoDCompliant bool    `json:"dodCompliant"`
	TrailFresh   bool    `json:"trailFresh"`
	CleanGit     bool    `json:"cleanGit"`

	// Model-based outcomes (LLM-as-a-judge)
	JudgeScore    float64 `json:"judgeScore"` // 0.0 - 100.0
	JudgeFeedback string  `json:"judgeFeedback,omitempty"`

	// Synthesized result
	OverallScore float64 `json:"overallScore"` // 0.0 - 100.0
	Verdict      string  `json:"verdict"`      // "PASS" | "FAIL" | "PARTIAL"
	Detail       string  `json:"detail,omitempty"`
	Error        string  `json:"error,omitempty"`
	WorkspaceDir string  `json:"workspaceDir,omitempty"`
}

// ComparisonSummary compares a Harnessed run against a Baseline run.
type ComparisonSummary struct {
	TaskID          string  `json:"taskId"`
	HarnessedScore  float64 `json:"harnessedScore"`
	BaselineScore   float64 `json:"baselineScore"`
	ScoreDelta      float64 `json:"scoreDelta"`
	HarnessedTokens int     `json:"harnessedTokens"`
	BaselineTokens  int     `json:"baselineTokens"`
	TokenDeltaPct   float64 `json:"tokenDeltaPct"`
	HarnessedCost   float64 `json:"harnessedCost"`
	BaselineCost    float64 `json:"baselineCost"`
	CostDeltaUSD    float64 `json:"costDeltaUsd"`
	HarnessedTime   float64 `json:"harnessedTimeSec"`
	BaselineTime    float64 `json:"baselineTimeSec"`
	HarnessedTurns  int     `json:"harnessedTurns"`
	BaselineTurns   int     `json:"baselineTurns"`
}
