package metrics

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Options are one ingest's inputs (spec § Verbs).
type Options struct {
	Home string
	// HomeGiven is true when the operator passed --home: the transcript-root
	// rule then ignores CLAUDE_CONFIG_DIR (D9).
	HomeGiven bool
	// Vendor is "claude", "agy" or "all"; empty means all.
	Vendor string
	// Source is a sandbox home; Label (KEY=VALUE) is required with it, and
	// Outcome ("pass", "fail" or "") is stamped with it (D11).
	Source, Label, Outcome string
	Getenv                 func(string) string
}

// Result counts one ingest.
type Result struct {
	Sessions, Steps, Events, Skipped int
	Failed                           []string
}

// Line is the summary ingest prints.
func (r Result) Line() string {
	return fmt.Sprintf("metrics: %d sessions, %d steps, %d events added; %d lines skipped (unknown shape); %d sources failed",
		r.Sessions, r.Steps, r.Events, r.Skipped, len(r.Failed))
}

// Validate is the usage check: exit 2 when it fails.
func (o Options) Validate() error {
	switch o.Vendor {
	case "", "all", "claude", "agy", "antigravity":
	default:
		return fmt.Errorf("--vendor %q: want claude, agy or all", o.Vendor)
	}
	if o.Source != "" && o.Label == "" {
		return fmt.Errorf("--source needs --label KEY=VALUE")
	}
	if o.Label != "" {
		k, v, ok := strings.Cut(o.Label, "=")
		if !ok || k == "" || v == "" {
			return fmt.Errorf("--label %q: want KEY=VALUE", o.Label)
		}
	}
	if o.Outcome != "" && o.Outcome != "pass" && o.Outcome != "fail" {
		return fmt.Errorf("--outcome %q: want pass or fail", o.Outcome)
	}
	if (o.Outcome != "" || o.Label != "") && o.Source == "" {
		return fmt.Errorf("--label and --outcome are for a --source run")
	}
	return nil
}

// claudeRoot is where Claude Code's transcripts are read from: --source, then
// --home, then CLAUDE_CONFIG_DIR only when neither flag was given, then the
// home (D9). The environment variable never redirects a scratch-home run to
// the operator's real transcripts (R6).
func (o Options) claudeRoot() string {
	switch {
	case o.Source != "":
		return filepath.Join(o.Source, ".claude", "projects")
	case o.HomeGiven:
		return filepath.Join(o.Home, ".claude", "projects")
	}
	getenv := o.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	if dir := getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "projects")
	}
	return filepath.Join(o.Home, ".claude", "projects")
}

func (o Options) agyRoot() string {
	if o.Source != "" {
		return filepath.Join(o.Source, ".gemini", "antigravity-cli")
	}
	return filepath.Join(o.Home, ".gemini", "antigravity-cli")
}

func (o Options) wants(vendor string) bool {
	v := o.Vendor
	if v == "antigravity" {
		v = "agy"
	}
	return v == "" || v == "all" || v == vendor
}

// Ingest reads every transcript the options name into s. A source that fails
// is named in Result.Failed and the rest are still read; the error is only for
// a store that cannot be written at all.
func Ingest(s *Store, o Options, log io.Writer) (Result, error) {
	var r Result
	before, err := s.counts()
	if err != nil {
		return r, err
	}
	if o.wants("claude") {
		s.ingestClaude(o, &r, log)
	}
	if o.wants("agy") {
		s.ingestAgy(o, &r, log)
	}
	// A kit learns its run's outcome after the run (accept.py): the label's
	// sessions with no outcome yet get this one, though no transcript changed.
	// An earlier run, stamped after its own acceptance, keeps its outcome.
	if o.Outcome != "" {
		if _, err := s.db.Exec(`UPDATE sessions SET outcome = ? WHERE label = ? AND outcome IS NULL`, o.Outcome, o.Label); err != nil {
			r.Failed = append(r.Failed, "outcome")
			fmt.Fprintf(log, "metrics: outcome: %v\n", err)
		}
	}
	// Events come from the real home's projects and feedback log, also for a
	// --source run: its sessions' project roots are where git and the trail are.
	notes, err := s.joinEvents(o.Home)
	if err != nil {
		r.Failed = append(r.Failed, "events")
		fmt.Fprintf(log, "metrics: events: %v\n", err)
	}
	for _, n := range notes {
		fmt.Fprintln(log, "metrics: note: "+n)
	}
	after, err := s.counts()
	if err != nil {
		return r, err
	}
	r.Sessions, r.Steps, r.Events = after[0]-before[0], after[1]-before[1], after[2]-before[2]
	return r, nil
}

func (s *Store) counts() ([3]int, error) {
	var c [3]int
	for i, t := range []string{"sessions", "steps", "events"} {
		if err := s.db.QueryRow(`SELECT count(*) FROM ` + t).Scan(&c[i]); err != nil {
			return c, err
		}
	}
	return c, nil
}
