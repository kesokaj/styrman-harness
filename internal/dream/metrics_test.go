package dream

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kesokaj/styrman-harness/internal/metrics"
)

// claudeLine puts one Claude Code transcript with a Bash call in the
// fixture's home.
func claudeLine(t *testing.T, f fixture) {
	t.Helper()
	dir := filepath.Join(f.home, ".claude", "projects", "-p")
	os.MkdirAll(dir, 0o755)
	line := `{"type":"assistant","sessionId":"S9","cwd":"` + f.root + `","timestamp":"2026-09-30T19:00:00Z","message":{"id":"m","content":[{"type":"tool_use","id":"t","name":"Bash","input":{"command":"git status"}}]}}` + "\n"
	if err := os.WriteFile(filepath.Join(dir, "S9.jsonl"), []byte(line), 0o600); err != nil {
		t.Fatal(err)
	}
}

// A real run ingests every session's shape into the metrics store before it
// dreams (R3, ADR-055), and says so in dream.log.
func TestRunIngestsMetricsFirst(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	claudeLine(t, f)
	fakeCLI(t, "claude", canned, 0)
	var out bytes.Buffer
	if err := Run(runOpts(f, &out, "claude")); err != nil {
		t.Fatalf("run: %v\n%s", err, out.String())
	}
	s, err := metrics.Open(metrics.Path(f.home))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	sum, err := s.Show(metrics.ShowOptions{Since: now.AddDate(0, 0, -2), Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if sum.B.ToolCalls["claude"] != 1 {
		t.Fatalf("the dream did not ingest: tool calls %v", sum.B.ToolCalls)
	}
	if log := read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log")); !strings.Contains(log, "metrics: 1 sessions") {
		t.Errorf("dream.log lacks the ingest line:\n%s", log)
	}
}

// A dry run ingests nothing and writes no store; it says what it would do.
func TestDryRunDoesNotIngest(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	claudeLine(t, f)
	var out bytes.Buffer
	o := runOpts(f, &out, "claude")
	o.DryRun = true
	if err := Run(o); err != nil {
		t.Fatalf("dry run: %v\n%s", err, out.String())
	}
	if _, err := os.Stat(metrics.Path(f.home)); !os.IsNotExist(err) {
		t.Fatalf("a dry run created the metrics store: %v", err)
	}
	if !strings.Contains(out.String(), "[dry-run] would ingest metrics into "+metrics.Path(f.home)) {
		t.Errorf("dry run does not name the ingest:\n%s", out.String())
	}
}

// An ingest failure is logged and changes neither the exit status nor
// last_run.error: the metrics never stop a dream.
func TestIngestFailureDoesNotFailDream(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	fakeCLI(t, "claude", canned, 0)
	saved := ingestMetrics
	t.Cleanup(func() { ingestMetrics = saved })
	ingestMetrics = func(string) (string, error) { return "", errors.New("disk on fire") }
	var out bytes.Buffer
	if err := Run(runOpts(f, &out, "claude")); err != nil {
		t.Fatalf("run failed because ingest failed: %v", err)
	}
	st, err := readState(f.home)
	if err != nil {
		t.Fatal(err)
	}
	if st.LastRun.Error != "" {
		t.Errorf("last_run.error = %q, want empty", st.LastRun.Error)
	}
	if log := read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log")); !strings.Contains(log, "metrics ingest failed: disk on fire") {
		t.Errorf("dream.log lacks the failure:\n%s", log)
	}
}
