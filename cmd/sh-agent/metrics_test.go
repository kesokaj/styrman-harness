package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMetricsVerb: ingest reads <home>/.claude/projects under --home even with
// CLAUDE_CONFIG_DIR pointing elsewhere (D9, R6), prints its summary line, and
// export writes the store and refuses to overwrite.
func TestMetricsVerb(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".claude", "projects", "-p")
	os.MkdirAll(dir, 0o755)
	line := `{"type":"assistant","sessionId":"S1","cwd":"/p","timestamp":"2026-10-08T10:00:00Z","message":{"id":"m1","content":[{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"git status"}}]}}` + "\n"
	os.WriteFile(filepath.Join(dir, "S1.jsonl"), []byte(line), 0o600)

	decoy := filepath.Join(t.TempDir(), "decoy")
	os.MkdirAll(filepath.Join(decoy, "projects", "-q"), 0o755)
	os.WriteFile(filepath.Join(decoy, "projects", "-q", "D.jsonl"), []byte(strings.ReplaceAll(line, "S1", "DECOY")), 0o600)
	t.Setenv("CLAUDE_CONFIG_DIR", decoy)

	var out bytes.Buffer
	if err := run([]string{"metrics", "ingest", "--home", home, "--vendor", "claude"}, &out); err != nil {
		t.Fatalf("ingest: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "metrics: 1 sessions, 2 steps") {
		t.Fatalf("ingest printed %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(home, ".config", "sh-agent", "metrics.db")); err != nil {
		t.Fatal(err)
	}

	exp := filepath.Join(t.TempDir(), "m.jsonl")
	out.Reset()
	if err := run([]string{"metrics", "export", "--home", home, exp}, &out); err != nil {
		t.Fatalf("export: %v\n%s", err, out.String())
	}
	b, _ := os.ReadFile(exp)
	if !bytes.Contains(b, []byte(`"S1"`)) || bytes.Contains(b, []byte("DECOY")) {
		t.Fatalf("export holds %s", b)
	}
	if err := run([]string{"metrics", "export", "--home", home, exp}, &out); err == nil {
		t.Fatal("export overwrote an existing file")
	}
}

func TestMetricsUsageErrorsExit2(t *testing.T) {
	home := t.TempDir()
	for _, args := range [][]string{
		{"metrics"},
		{"metrics", "frobnicate"},
		{"metrics", "ingest", "extra"},
		{"metrics", "ingest", "--source", home},
		{"metrics", "ingest", "--label", "a=b"},
		{"metrics", "ingest", "--source", home, "--label", "a=b", "--outcome", "maybe"},
		{"metrics", "ingest", "--vendor", "gpt"},
		{"metrics", "export"},
	} {
		err := run(append(args, "--home", home), &bytes.Buffer{})
		var code exitCode
		if !errors.As(err, &code) || code != 2 {
			t.Errorf("%q: err = %v, want exit 2", args, err)
		}
	}
	if _, err := os.Stat(filepath.Join(home, ".config", "sh-agent", "metrics.db")); !os.IsNotExist(err) {
		t.Errorf("a usage error created the store: %v", err)
	}
}
