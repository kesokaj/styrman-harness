package eval

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTasksRegistry(t *testing.T) {
	tasks := AllTasks()
	if len(tasks) == 0 {
		t.Fatal("expected default tasks to be registered")
	}

	task, ok := GetTask("gf-go-cli")
	if !ok {
		t.Fatal("expected gf-go-cli task to exist")
	}
	if task.Category != CategoryGreenfield {
		t.Fatalf("expected category greenfield, got %s", task.Category)
	}
}

func TestTelemetryParseEnvelope(t *testing.T) {
	raw := []byte(`{
		"total_cost_usd": 0.0125,
		"total_tokens": 4500,
		"input_tokens": 3800,
		"output_tokens": 700,
		"turns": 4,
		"tool_calls": {
			"run_command": 3,
			"write_to_file": 2
		}
	}`)

	telem := ParseEnvelopeJSON(raw)
	if telem.TotalTokens != 4500 {
		t.Errorf("expected 4500 tokens, got %d", telem.TotalTokens)
	}
	if telem.Turns != 4 {
		t.Errorf("expected 4 turns, got %d", telem.Turns)
	}
	if telem.CostUSD != 0.0125 {
		t.Errorf("expected cost 0.0125, got %f", telem.CostUSD)
	}
	if telem.ToolCalls["run_command"] != 3 {
		t.Errorf("expected 3 run_command calls, got %d", telem.ToolCalls["run_command"])
	}
}

func TestMockRunnerExecution(t *testing.T) {
	tempDir := t.TempDir()
	cfg := Config{
		LocalLedger: filepath.Join(tempDir, "ledger.jsonl"),
	}

	runner := NewRunner(tempDir, cfg)
	runner.WorkDirRoot = filepath.Join(tempDir, "workspaces")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	opts := RunOptions{
		TaskID:  "gf-go-cli",
		Vendor:  "mock",
		Timeout: 30 * time.Second,
	}

	res, err := runner.Run(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected runner error: %v", err)
	}

	if res.Verdict != "PASS" {
		t.Fatalf("expected PASS verdict for mock runner, got %s (detail: %s)", res.Verdict, res.Detail)
	}
	if !res.BuildSuccess {
		t.Fatal("expected build success")
	}
	if res.TestPassRate < 1.0 {
		t.Fatalf("expected 1.0 test pass rate, got %f", res.TestPassRate)
	}
	if res.Turns <= 0 {
		t.Errorf("expected positive turns, got %d", res.Turns)
	}

	// Verify ledger was written
	storage := NewStorage(tempDir, cfg)
	records, err := storage.ReadLedger()
	if err != nil {
		t.Fatalf("failed to read ledger: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 ledger record, got %d", len(records))
	}
	if records[0].RunID != res.RunID {
		t.Errorf("expected ledger run ID %s, got %s", res.RunID, records[0].RunID)
	}
}

func TestConfigDefaults(t *testing.T) {
	tempDir := t.TempDir()
	cfg := LoadConfig(tempDir)
	if cfg.Region == "" {
		t.Error("expected non-empty default region")
	}
	if cfg.VertexModel == "" {
		t.Error("expected non-empty default vertex model")
	}

	// Write custom eval-config.json
	agentsDir := filepath.Join(tempDir, ".agents")
	_ = os.MkdirAll(agentsDir, 0o755)
	_ = os.WriteFile(filepath.Join(agentsDir, "eval-config.json"), []byte(`{
		"projectId": "my-test-proj",
		"region": "us-central1"
	}`), 0o644)

	customCfg := LoadConfig(tempDir)
	if customCfg.ProjectID != "my-test-proj" {
		t.Errorf("expected projectId my-test-proj, got %s", customCfg.ProjectID)
	}
	if customCfg.Region != "us-central1" {
		t.Errorf("expected region us-central1, got %s", customCfg.Region)
	}
}
