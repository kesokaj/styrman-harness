// Package dispatch runs the Definition-of-Done rows a working tree owes.
//
// It is deliberately not an agent. Putting a model in charge of deciding which
// agents to run adds a layer that can stall, summarises twice, and moves a
// determination that belongs in a gate into model judgement. This computes the
// owed rows from the diff, runs them in parallel with hard timeouts, records
// what each one cost, and hands back one compact table.
package dispatch

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kesokaj/styrman-harness/internal/trail"
)

// Result is one row's outcome.
type Result struct {
	Row      trail.Row `json:"row"`
	OK       bool      `json:"ok"`
	Detail   string    `json:"detail"`
	CostUSD  float64   `json:"costUsd"`
	Duration string    `json:"duration"`
	Artifact string    `json:"artifact,omitempty"`
}

// Run is the `harness dispatch` subcommand.
func Run(args []string) error {
	fset := flag.NewFlagSet("dispatch", flag.ContinueOnError)
	dryRun := fset.Bool("dry-run", false, "print the owed rows and exit")
	asJSON := fset.Bool("json", false, "emit results as JSON")
	vendor := fset.String("vendor", "", "cli to dispatch with (claude|agy); default: whichever is on PATH")
	only := fset.String("only", "", "run only these owed rows (comma-separated)")
	if err := fset.Parse(args); err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	proj, ok := trail.Resolve(wd)
	if !ok {
		return fmt.Errorf("not a git repository: %s", wd)
	}
	if !proj.Opted() {
		return fmt.Errorf("no trail contract here: %s/changes does not exist", proj.DocsRel())
	}
	cfg, err := trail.LoadConfig(proj.Root)
	if err != nil {
		return err
	}

	day := trail.Today()
	report := proj.Owed(cfg, day)
	if len(report.Due) == 0 {
		fmt.Println("nothing owed: the trail covers the current working tree")
		return nil
	}

	rows, err := selectRows(report.Rows(), *only)
	if err != nil {
		return err
	}
	if *dryRun {
		for _, due := range report.Due {
			fmt.Printf("owed  %-12s %s\n", due.Row, due.Reason)
		}
		for _, s := range report.Skipped {
			fmt.Printf("skip  %-12s %s\n", s.Row, s.Reason)
		}
		return nil
	}

	cli, err := resolveCLI(*vendor)
	if err != nil {
		return err
	}

	results := runRows(proj, cfg, cli, rows, day)
	if err := appendLedger(proj, results, day); err != nil {
		fmt.Fprintln(os.Stderr, "warning: ledger:", err)
	}
	return present(results, *asJSON)
}

// selectRows narrows the owed rows to those named in --only. Naming a row that
// is not owed is an error rather than a silent no-op: it usually means the
// caller expected work the contract does not actually require.
func selectRows(owed []trail.Row, only string) ([]trail.Row, error) {
	if strings.TrimSpace(only) == "" {
		return owed, nil
	}
	isOwed := make(map[trail.Row]bool, len(owed))
	for _, r := range owed {
		isOwed[r] = true
	}
	wanted := make(map[trail.Row]bool)
	for _, name := range strings.Split(only, ",") {
		row := trail.Row(strings.TrimSpace(name))
		if row == "" {
			continue
		}
		if !isOwed[row] {
			return nil, fmt.Errorf("row %q is not owed right now; `--dry-run` lists what is", row)
		}
		wanted[row] = true
	}

	var out []trail.Row
	for _, r := range owed {
		if wanted[r] {
			out = append(out, r)
		}
	}
	return out, nil
}

// runRows executes the owed rows concurrently. The verifier goes first and
// alone: it is the cheapest row and the one most likely to make the rest moot,
// since there is no point reviewing code whose tests do not pass.
func runRows(proj trail.Project, cfg trail.Config, cli cliTool, rows []trail.Row, day string) []Result {
	var results []Result
	remaining := rows

	if len(rows) > 0 && rows[0] == trail.RowVerifier {
		remaining = rows[1:]
		res := runVerifier(proj, cfg, cli, day)
		results = append(results, res)
		if !res.OK {
			// Tests failed. Reviewing and recording work that does not build
			// wastes the rows and buries the actual problem.
			for _, row := range remaining {
				results = append(results, Result{
					Row: row, OK: false, Duration: "0s",
					Detail: "skipped — verification failed; fix the failure first",
				})
			}
			return results
		}
	}

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	for _, row := range remaining {
		wg.Add(1)
		go func(row trail.Row) {
			defer wg.Done()
			res := runAgentRow(proj, cfg, cli, row, day)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(row)
	}
	wg.Wait()

	sort.Slice(results, func(i, j int) bool { return results[i].Row < results[j].Row })
	return results
}

// runAgentRow dispatches one row to a subagent headlessly.
func runAgentRow(proj trail.Project, cfg trail.Config, cli cliTool, row trail.Row, day string) Result {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.DispatchTimeout)*time.Second)
	defer cancel()

	prompt := rowPrompt(proj, row, day)
	out, err := cli.run(ctx, proj.Root, string(row), prompt)
	res := Result{Row: row, Duration: time.Since(start).Round(time.Second).String()}

	switch {
	case ctx.Err() == context.DeadlineExceeded:
		res.Detail = fmt.Sprintf("timed out after %ds", cfg.DispatchTimeout)
		return res
	case err != nil:
		res.Detail = strings.TrimSpace(err.Error())
		return res
	}
	res.CostUSD = out.CostUSD

	// The row is satisfied by its artifact existing, not by what the agent
	// said about it. Chat is not the record.
	if proj.HasArtifact(row, day) {
		res.OK = true
		res.Artifact = artifactPath(proj, row, day)
		res.Detail = "artifact written"
		return res
	}
	res.Detail = "agent finished without writing its artifact"
	return res
}

func artifactPath(proj trail.Project, row trail.Row, day string) string {
	for _, a := range proj.Artifacts(row.Class()) {
		if a.Owner == row && a.Date == day {
			rel, err := filepath.Rel(proj.Root, a.Path)
			if err != nil {
				return a.Path
			}
			return rel
		}
	}
	return ""
}

// rowPrompt is the dispatch prompt for one row. Each ends with the same
// instruction to keep the reply short: the detail belongs in the artifact, and
// repeating it into the caller's context is how a harness fills its own window
// with things it has already written to disk.
func rowPrompt(proj trail.Project, row trail.Row, day string) string {
	base := fmt.Sprintf(
		"Today is %s. Write your artifact to %s/%s/ with YAML frontmatter containing `owner: %s` and `date: %s`.",
		day, proj.DocsRel(), row.Class(), row, day)
	var task string
	switch row {
	case trail.RowChangelog:
		task = "Record the current working-tree changes as a changelog entry."
	case trail.RowReviewer:
		task = "Review the current working-tree diff against the project's review standards."
	case trail.RowSecurity:
		task = "Security-review the current working-tree diff and its blast radius."
	default:
		task = "Complete your Definition-of-Done row for the current changes."
	}
	return task + " " + base +
		" Reply with at most ten lines: a verdict, the artifact path, and up to three findings. Nothing else."
}

// present prints the dispatch table.
func present(results []Result, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}
	var total float64
	failed := 0
	fmt.Printf("%-12s %-6s %-8s %-9s %s\n", "ROW", "OK", "TIME", "COST", "DETAIL")
	for _, r := range results {
		mark := "yes"
		if !r.OK {
			mark, failed = "NO", failed+1
		}
		total += r.CostUSD
		detail := r.Detail
		if r.Artifact != "" {
			detail = r.Artifact
		}
		fmt.Printf("%-12s %-6s %-8s $%-8.4f %s\n", r.Row, mark, r.Duration, r.CostUSD, detail)
	}
	fmt.Printf("\n%d row(s), $%.4f total\n", len(results), total)
	if failed > 0 {
		return fmt.Errorf("%d row(s) failed", failed)
	}
	return nil
}

// appendLedger records what this dispatch cost. Without it, "the subagents run
// for no benefit" stays a feeling; with it, it is a number that can be argued
// with.
func appendLedger(proj trail.Project, results []Result, day string) error {
	dir := filepath.Join(proj.Docs, "testing")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, "dispatch-ledger.jsonl"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	entry, err := json.Marshal(map[string]any{
		"date":    day,
		"ts":      time.Now().UTC().Format(time.RFC3339),
		"results": results,
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, string(entry))
	return err
}

// exists reports whether a path is present.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
