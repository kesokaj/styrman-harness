package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/eval"
)

const usage = `eval-harness — evaluation suite for styrman-harness

Usage:
  eval-harness list
        List all available evaluation tasks.
  eval-harness run [--task <id>] [--vendor claude|agy|mock] [--baseline] [--json]
        Execute an evaluation task and calculate tokens, tools, turns, time, and rating.
  eval-harness report [--json]
        Display aggregated evaluation history from the ledger.
  eval-harness compare [--task <id>] [--json]
        Compare newest harnessed run vs baseline run.
  eval-harness doctor
        Verify local agent CLIs, Go toolchain, and GCP configuration.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "eval:", err)
		os.Exit(1)
	}

	cfg := eval.LoadConfig(repoRoot)

	switch os.Args[1] {
	case "list":
		runList()
	case "run":
		runEval(repoRoot, cfg, os.Args[2:])
	case "report":
		runReport(repoRoot, cfg, os.Args[2:])
	case "compare":
		runCompare(repoRoot, cfg, os.Args[2:])
	case "doctor":
		runDoctor(cfg)
	case "-h", "--help", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "eval: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runList() {
	tasks := eval.AllTasks()
	fmt.Printf("%-18s %-12s %s\n", "TASK ID", "CATEGORY", "NAME")
	fmt.Println(strings.Repeat("-", 70))
	for _, t := range tasks {
		fmt.Printf("%-18s %-12s %s\n", t.ID, t.Category, t.Name)
	}
}

func runEval(repoRoot string, cfg eval.Config, args []string) {
	fset := flag.NewFlagSet("run", flag.ContinueOnError)
	taskID := fset.String("task", "gf-go-cli", "task ID to run (default: gf-go-cli)")
	vendor := fset.String("vendor", "mock", "agent vendor: claude | agy | mock")
	baseline := fset.Bool("baseline", false, "run baseline unharnessed agent")
	asJSON := fset.Bool("json", false, "output result as JSON")
	timeoutMin := fset.Int("timeout", 10, "timeout in minutes")

	if err := fset.Parse(args); err != nil {
		os.Exit(2)
	}

	runner := eval.NewRunner(repoRoot, cfg)
	ctx := context.Background()

	opts := eval.RunOptions{
		TaskID:   *taskID,
		Vendor:   *vendor,
		Baseline: *baseline,
		Timeout:  time.Duration(*timeoutMin) * time.Minute,
	}

	fmt.Printf("Starting eval run for %q (vendor: %s, mode: %s)...\n",
		*taskID, *vendor, modeName(*baseline))

	res, err := runner.Run(ctx, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval error: %v\n", err)
		os.Exit(1)
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(res)
		return
	}

	printResultSummary(res)
}

func runReport(repoRoot string, cfg eval.Config, args []string) {
	fset := flag.NewFlagSet("report", flag.ContinueOnError)
	asJSON := fset.Bool("json", false, "output report as JSON")
	if err := fset.Parse(args); err != nil {
		os.Exit(2)
	}

	storage := eval.NewStorage(repoRoot, cfg)
	results, err := storage.ReadLedger()
	if err != nil {
		fmt.Fprintln(os.Stderr, "report error:", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No evaluation records found in ledger.")
		return
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
		return
	}

	fmt.Printf("%-10s %-16s %-8s %-10s %-7s %-8s %-7s %s\n",
		"RUN ID", "TASK", "VENDOR", "MODE", "VERDICT", "SCORE", "TOKENS", "DURATION")
	fmt.Println(strings.Repeat("-", 80))
	for _, r := range results {
		fmt.Printf("%-10s %-16s %-8s %-10s %-7s %-8.1f %-7d %s\n",
			r.RunID, r.TaskID, r.Vendor, r.Mode, r.Verdict, r.OverallScore, r.TotalTokens, r.Duration.Round(time.Second))
	}
}

func runDoctor(cfg eval.Config) {
	fmt.Println("eval-harness doctor:")
	fmt.Println(strings.Repeat("-", 40))

	// 1. Go
	if _, err := exec.LookPath("go"); err == nil {
		out, _ := exec.Command("go", "version").Output()
		fmt.Printf("ok    Go Toolchain: %s", string(out))
	} else {
		fmt.Println("FAIL  Go Toolchain not found")
	}

	// 2. Git
	if _, err := exec.LookPath("git"); err == nil {
		fmt.Println("ok    Git available")
	} else {
		fmt.Println("FAIL  Git not found")
	}

	// 3. Agent CLIs
	if _, err := exec.LookPath("claude"); err == nil {
		fmt.Println("ok    Claude Code CLI (claude) found")
	} else {
		fmt.Println("warn  Claude Code CLI (claude) not on PATH")
	}

	if _, err := exec.LookPath("agy"); err == nil {
		fmt.Println("ok    Antigravity CLI (agy) found")
	} else {
		fmt.Println("warn  Antigravity CLI (agy) not on PATH")
	}

	// 4. GCP Config
	if cfg.ProjectID != "" {
		fmt.Printf("ok    GCP Project: %s\n", cfg.ProjectID)
	} else {
		fmt.Println("warn  GCP Project not set in .agents/eval-config.json")
	}

	if cfg.GCSBucket != "" {
		fmt.Printf("ok    GCS Bucket: gs://%s\n", cfg.GCSBucket)
	} else {
		fmt.Println("info  GCS Bucket not configured (local ledger only)")
	}

	if _, err := exec.Command("gcloud", "auth", "application-default", "print-access-token").Output(); err == nil {
		fmt.Println("ok    Application Default Credentials (ADC) valid")
	} else {
		fmt.Println("warn  ADC not configured; LLM judge and GCS sync require 'gcloud auth application-default login'")
	}
}

func modeName(baseline bool) string {
	if baseline {
		return "baseline"
	}
	return "harnessed"
}

func printResultSummary(r eval.EvalRunResult) {
	fmt.Println("")
	fmt.Println("Evaluation Result:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Run ID:        %s\n", r.RunID)
	fmt.Printf("Task:          %s (%s)\n", r.TaskName, r.TaskID)
	fmt.Printf("Verdict:       %s (Score: %.1f/100)\n", r.Verdict, r.OverallScore)
	fmt.Printf("Duration:      %s\n", r.Duration.Round(time.Millisecond))
	fmt.Printf("Turns:         %d\n", r.Turns)
	fmt.Printf("Tool Calls:    %d\n", len(r.ToolCalls))
	for tool, count := range r.ToolCalls {
		fmt.Printf("  - %s: %d\n", tool, count)
	}
	fmt.Printf("Total Tokens:  %d (Input: %d, Output: %d, Cache: %d)\n",
		r.TotalTokens, r.InputTokens, r.OutputTokens, r.CacheTokens)
	if r.CostUSD > 0 {
		fmt.Printf("Estimated Cost: $%.4f\n", r.CostUSD)
	}
	fmt.Println("")
	fmt.Println("Breakdown:")
	fmt.Printf("  Build Success:  %v\n", r.BuildSuccess)
	fmt.Printf("  Test Pass Rate: %.1f%%\n", r.TestPassRate*100)
	fmt.Printf("  Clean Git:      %v\n", r.CleanGit)
	fmt.Printf("  Trail Fresh:    %v\n", r.TrailFresh)
	if r.JudgeScore > 0 {
		fmt.Printf("  Judge Score:    %.1f/100\n", r.JudgeScore)
		if r.JudgeFeedback != "" {
			fmt.Printf("  Judge Feedback: %s\n", r.JudgeFeedback)
		}
	}
	if r.Detail != "" {
		fmt.Printf("Detail: %s\n", r.Detail)
	}
	if r.Error != "" {
		fmt.Printf("Error:  %s\n", r.Error)
	}
}

func runCompare(repoRoot string, cfg eval.Config, args []string) {
	fset := flag.NewFlagSet("compare", flag.ContinueOnError)
	taskID := fset.String("task", "gf-go-cli", "task ID to compare")
	asJSON := fset.Bool("json", false, "output comparison as JSON")
	if err := fset.Parse(args); err != nil {
		os.Exit(2)
	}

	storage := eval.NewStorage(repoRoot, cfg)
	results, err := storage.ReadLedger()
	if err != nil {
		fmt.Fprintln(os.Stderr, "compare error:", err)
		os.Exit(1)
	}

	var newestHarnessed, newestBaseline *eval.EvalRunResult
	for i := len(results) - 1; i >= 0; i-- {
		r := results[i]
		if r.TaskID != *taskID {
			continue
		}
		if r.Mode == "harnessed" && newestHarnessed == nil {
			cp := r
			newestHarnessed = &cp
		}
		if r.Mode == "baseline" && newestBaseline == nil {
			cp := r
			newestBaseline = &cp
		}
		if newestHarnessed != nil && newestBaseline != nil {
			break
		}
	}

	if newestHarnessed == nil || newestBaseline == nil {
		fmt.Fprintf(os.Stderr, "Need both a harnessed and baseline run for task %q to compare.\n", *taskID)
		os.Exit(1)
	}

	h := newestHarnessed
	b := newestBaseline

	tokenDeltaPct := 0.0
	if b.TotalTokens > 0 {
		tokenDeltaPct = float64(h.TotalTokens-b.TotalTokens) / float64(b.TotalTokens) * 100.0
	}
	costDelta := h.CostUSD - b.CostUSD

	if *asJSON {
		summary := eval.ComparisonSummary{
			TaskID:          *taskID,
			HarnessedScore:  h.OverallScore,
			BaselineScore:   b.OverallScore,
			ScoreDelta:      h.OverallScore - b.OverallScore,
			HarnessedTokens: h.TotalTokens,
			BaselineTokens:  b.TotalTokens,
			TokenDeltaPct:   tokenDeltaPct,
			HarnessedCost:   h.CostUSD,
			BaselineCost:    b.CostUSD,
			CostDeltaUSD:    costDelta,
			HarnessedTime:   h.DurationSec,
			BaselineTime:    b.DurationSec,
			HarnessedTurns:  h.Turns,
			BaselineTurns:   b.Turns,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(summary)
		return
	}

	fmt.Printf("A/B Baseline Comparison for Task: %s\n", *taskID)
	fmt.Println(strings.Repeat("=", 68))
	fmt.Printf("%-20s %-16s %-16s %s\n", "METRIC", "HARNESSED", "BASELINE", "DELTA")
	fmt.Println(strings.Repeat("-", 68))
	fmt.Printf("%-20s %-16.1f %-16.1f %+.1f\n", "Score", h.OverallScore, b.OverallScore, h.OverallScore-b.OverallScore)
	fmt.Printf("%-20s %-16s %-16s -\n", "Verdict", h.Verdict, b.Verdict)
	fmt.Printf("%-20s %-16s %-16s %+.1fs\n", "Duration", h.Duration.Round(time.Second), b.Duration.Round(time.Second), h.DurationSec-b.DurationSec)
	fmt.Printf("%-20s %-16d %-16d %+d\n", "Turns", h.Turns, b.Turns, h.Turns-b.Turns)
	fmt.Printf("%-20s %-16d %-16d %+d (%+.1f%%)\n", "Total Tokens", h.TotalTokens, b.TotalTokens, h.TotalTokens-b.TotalTokens, tokenDeltaPct)
	fmt.Printf("%-20s %-16d %-16d %+d\n", "Input Tokens", h.InputTokens, b.InputTokens, h.InputTokens-b.InputTokens)
	fmt.Printf("%-20s %-16d %-16d %+d\n", "Output Tokens", h.OutputTokens, b.OutputTokens, h.OutputTokens-b.OutputTokens)
	fmt.Printf("%-20s %-16d %-16d %+d\n", "Cache Tokens", h.CacheTokens, b.CacheTokens, h.CacheTokens-b.CacheTokens)
	fmt.Printf("%-20s $%-15.4f $%-15.4f +$%.4f\n", "Cost (USD)", h.CostUSD, b.CostUSD, costDelta)
	fmt.Printf("%-20s %-16d %-16d %+d\n", "Total Tools", len(h.ToolCalls), len(b.ToolCalls), len(h.ToolCalls)-len(b.ToolCalls))
	fmt.Printf("%-20s %-16v %-16v -\n", "Clean Git", h.CleanGit, b.CleanGit)
	fmt.Printf("%-20s %-16v %-16v -\n", "Build Success", h.BuildSuccess, b.BuildSuccess)
	fmt.Printf("%-20s %-16s %-16s -\n", "Test Pass Rate", fmt.Sprintf("%.1f%%", h.TestPassRate*100), fmt.Sprintf("%.1f%%", b.TestPassRate*100))
}
