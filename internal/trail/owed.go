package trail

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config tunes the proportionality rule. A project overrides it by committing
// .agents/harness.json; the zero value is never used, defaults come from
// DefaultConfig.
type Config struct {
	// ReviewLines is the added+removed line count at or above which the
	// reviewer row fires. Below it, an in-process /code-review is enough and a
	// full subagent is not worth its cost.
	ReviewLines int `json:"reviewLines"`
	// ReviewFiles is the changed-file count at or above which the reviewer row
	// fires regardless of line count — breadth is its own risk.
	ReviewFiles int `json:"reviewFiles"`
	// DispatchTimeout caps a single dispatched row, in seconds.
	DispatchTimeout int `json:"dispatchTimeoutSeconds"`
	// SourceGlobs declare paths that count as source even when the default
	// classification would not treat them as such — a project whose product is
	// prose needs this, or a rewrite of its entire deliverable owes nothing.
	SourceGlobs []string `json:"sourceGlobs"`
}

// DefaultConfig is deliberately permissive: the accountability rows exist to
// catch what a solo operator's own reading misses, and below these thresholds
// the operator's own reading is cheaper and just as good.
func DefaultConfig() Config {
	return Config{ReviewLines: 50, ReviewFiles: 5, DispatchTimeout: 900}
}

// LoadConfig reads .agents/harness.json from the project root, falling back to
// DefaultConfig. A malformed file is an error rather than a silent default —
// a gate that quietly ignores its own configuration is worse than one that
// refuses to start.
func LoadConfig(root string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(filepath.Join(root, ".agents", "harness.json"))
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), fmt.Errorf(".agents/harness.json: %w", err)
	}
	if cfg.ReviewLines <= 0 {
		cfg.ReviewLines = DefaultConfig().ReviewLines
	}
	if cfg.ReviewFiles <= 0 {
		cfg.ReviewFiles = DefaultConfig().ReviewFiles
	}
	if cfg.DispatchTimeout <= 0 {
		cfg.DispatchTimeout = DefaultConfig().DispatchTimeout
	}
	return cfg, nil
}

// Due is one owed row plus why it is owed, so a gate can explain itself
// instead of merely refusing.
type Due struct {
	Row    Row
	Reason string
}

// Report is the full picture of what a project owes right now.
type Report struct {
	Changed []string // changed source paths
	Lines   int      // added+removed lines
	Due     []Due    // rows owed and not yet satisfied
	Skipped []Due    // rows deliberately not fired, and why
}

// Owed computes the rows due for the current working state.
//
// Rows split by purpose. The changelog is a resumability row: it is what lets
// a zero-context agent pick the work up again, it is cheap, and it is always
// due. The reviewer, security and verifier rows are accountability rows —
// re-derivable from the diff, expensive to produce, and worth their cost only
// in proportion to the change. Firing all of them on a three-line edit is how
// a harness ends up burning model time for no benefit.
func (p Project) Owed(cfg Config, day string) Report {
	var rep Report
	if !p.Opted() {
		return rep
	}
	p.SourceGlobs = cfg.SourceGlobs

	rep.Changed = p.DirtySource()
	if len(rep.Changed) == 0 {
		rep.Changed = p.HeadSourceToday()
	}
	if len(rep.Changed) == 0 {
		return rep
	}
	rep.Lines = p.ChangedLines()

	add := func(row Row, reason string) {
		if p.HasArtifact(row, day) {
			return
		}
		rep.Due = append(rep.Due, Due{Row: row, Reason: reason})
	}
	skip := func(row Row, reason string) {
		rep.Skipped = append(rep.Skipped, Due{Row: row, Reason: reason})
	}

	// Resumability — always due.
	add(RowChangelog, fmt.Sprintf("%d source file(s) changed", len(rep.Changed)))

	// Verification is cheap: it runs the project's own test commands as a
	// script and only escalates to a model to explain a failure.
	add(RowVerifier, "source changed; tests must run before any completion claim")

	// Accountability — proportional.
	switch {
	case rep.Lines >= cfg.ReviewLines:
		add(RowReviewer, fmt.Sprintf("%d changed lines ≥ threshold %d", rep.Lines, cfg.ReviewLines))
	case len(rep.Changed) >= cfg.ReviewFiles:
		add(RowReviewer, fmt.Sprintf("%d changed files ≥ threshold %d", len(rep.Changed), cfg.ReviewFiles))
	default:
		skip(RowReviewer, fmt.Sprintf("%d lines across %d file(s) is below the review threshold (%d lines / %d files) — read it yourself or run /code-review",
			rep.Lines, len(rep.Changed), cfg.ReviewLines, cfg.ReviewFiles))
	}

	if p.OwesSecurity(rep.Changed) {
		add(RowSecurity, "a changed path touches auth, crypto, secrets, input parsing, or dependencies")
	} else {
		skip(RowSecurity, "no changed path touches the security surface")
	}

	return rep
}

// Rows returns just the owed rows, in dispatch order: cheap and mandatory
// first so a failure surfaces before the expensive rows have run.
func (r Report) Rows() []Row {
	order := []Row{RowVerifier, RowChangelog, RowSecurity, RowReviewer}
	var out []Row
	for _, want := range order {
		for _, due := range r.Due {
			if due.Row == want {
				out = append(out, want)
			}
		}
	}
	return out
}
