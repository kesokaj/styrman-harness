package dispatch

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kesokaj/styrman-harness/internal/trail"
)

// stubCLI writes a script that behaves like `claude -p --agent X --output-format
// json`: it writes the artifact the row expects and prints the result envelope.
//
// This is what makes the dispatcher's CLI contract testable in CI. The real
// binaries need credentials and cost money per run, so without a stand-in the
// contract was only ever exercised by hand — and every seam that was only
// checked by hand in this project turned out to have a bug in it.
func stubCLI(t *testing.T, behaviour string) cliTool {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "stub-cli")

	script := `#!/bin/bash
# Args mirror the real CLIs: -p <prompt> --agent <name> --output-format json
prompt=""; agent=""
while [ $# -gt 0 ]; do
  case "$1" in
    -p) prompt="$2"; shift 2 ;;
    --agent) agent="$2"; shift 2 ;;
    *) shift ;;
  esac
done
` + behaviour + `
`
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return cliTool{
		name: path,
		args: func(agent, prompt string) []string {
			return []string{"-p", prompt, "--agent", agent, "--output-format", "json"}
		},
	}
}

// writesArtifact makes the stub extract docs root, owner, date and scope from
// the prompt and write a conforming artifact — exactly what a real agent does.
const writesArtifact = `
docs=$(printf '%s' "$prompt" | grep -oE '[^ ]+/(changes|reviews|testing|plans)/' | head -1)
day=$(printf '%s' "$prompt" | grep -oE 'date: [0-9-]+' | head -1 | cut -d' ' -f2)
scope=$(printf '%s' "$prompt" | grep -oE 'changeset:[a-f0-9]+' | head -1)
mkdir -p "$docs"
printf -- '---\nowner: %s\ndate: %s\nscope: %s\n---\n\nstub\n' "$agent" "$day" "$scope" > "$docs/stub-$agent.md"
printf '{"result":"wrote it","total_cost_usd":0.0123,"is_error":false}\n'
`

// writesNothing makes the stub claim success without writing anything.
const writesNothing = `printf '{"result":"all done, trust me","total_cost_usd":0.5,"is_error":false}\n'`

func scratchProject(t *testing.T) trail.Project {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"}, {"config", "user.email", "t@e.com"}, {"config", "user.name", "t"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents", "docs", "changes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	proj, ok := trail.Resolve(root)
	if !ok {
		t.Fatal("Resolve failed")
	}
	return proj
}

func TestDispatchAgentRowWritesAndDetectsArtifact(t *testing.T) {
	proj := scratchProject(t)
	cli := stubCLI(t, writesArtifact)
	scope := proj.Scope()

	res := runAgentRow(proj, trail.DefaultConfig(), cli, trail.RowChangelog, trail.Today(), scope)

	if !res.OK {
		t.Fatalf("row not satisfied: %+v", res)
	}
	if res.CostUSD != 0.0123 {
		t.Errorf("cost not parsed from the envelope: got %v", res.CostUSD)
	}
	if res.Artifact == "" {
		t.Error("no artifact path reported")
	}
	if !proj.HasArtifact(trail.RowChangelog, scope) {
		t.Error("the artifact the stub wrote does not satisfy the row")
	}
}

// TestDispatchRejectsUnwrittenArtifact: the row is satisfied by the artifact
// existing, never by the agent saying it wrote one.
func TestDispatchRejectsUnwrittenArtifact(t *testing.T) {
	proj := scratchProject(t)
	cli := stubCLI(t, writesNothing)

	res := runAgentRow(proj, trail.DefaultConfig(), cli, trail.RowChangelog, trail.Today(), proj.Scope())

	if res.OK {
		t.Fatal("a row was satisfied by an agent that wrote nothing")
	}
	if !strings.Contains(res.Detail, "without writing") {
		t.Errorf("detail should name the problem, got %q", res.Detail)
	}
	if res.CostUSD != 0.5 {
		t.Errorf("cost should still be recorded for a failed row: got %v", res.CostUSD)
	}
}

// TestDispatchTimeoutIsEnforced: a hung agent must not hang the dispatcher.
func TestDispatchTimeoutIsEnforced(t *testing.T) {
	proj := scratchProject(t)
	cli := stubCLI(t, "sleep 30")
	cfg := trail.DefaultConfig()
	cfg.DispatchTimeout = 1

	res := runAgentRow(proj, cfg, cli, trail.RowChangelog, trail.Today(), proj.Scope())

	if res.OK {
		t.Fatal("a timed-out row reported success")
	}
	if !strings.Contains(res.Detail, "timed out") {
		t.Errorf("detail should say it timed out, got %q", res.Detail)
	}
}

// TestVerifierRowRunsWithoutAModel is the redesign's central claim, asserted:
// a passing suite costs nothing because no CLI is invoked at all.
func TestVerifierRowRunsWithoutAModel(t *testing.T) {
	proj := scratchProject(t)
	if err := os.WriteFile(filepath.Join(proj.Root, ".agents", "harness.json"),
		[]byte(`{"verify":["true"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// A CLI that fails loudly if it is ever called.
	cli := stubCLI(t, `echo "the verifier row invoked a model" >&2; exit 3`)

	res := runVerifier(proj, trail.DefaultConfig(), cli, trail.Today(), proj.Scope())

	if !res.OK {
		t.Fatalf("verification of a passing suite failed: %+v", res)
	}
	if res.CostUSD != 0 {
		t.Errorf("a passing verification cost %v; it must invoke no model", res.CostUSD)
	}
}

// TestVerifierEscalatesFailureWithOutput: the failing path must hand the agent
// the captured output, since the agent is forbidden from re-running the suite.
func TestVerifierEscalatesFailureWithOutput(t *testing.T) {
	proj := scratchProject(t)
	if err := os.WriteFile(filepath.Join(proj.Root, ".agents", "harness.json"),
		[]byte(`{"verify":["echo distinctive-failure-marker; exit 1"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	captured := filepath.Join(t.TempDir(), "prompt.txt")
	cli := stubCLI(t, `printf '%s' "$prompt" > `+captured+`
`+writesArtifact)

	runVerifier(proj, trail.DefaultConfig(), cli, trail.Today(), proj.Scope())

	got, err := os.ReadFile(captured)
	if err != nil {
		t.Fatal("the verifier agent was never invoked on failure")
	}
	if !strings.Contains(string(got), "distinctive-failure-marker") {
		t.Error("the failure output was not passed to the verifier agent")
	}
}

func TestCLIRunParsesEnvelope(t *testing.T) {
	cli := stubCLI(t, `printf '{"result":"ok","total_cost_usd":1.5}\n'`)
	out, err := cli.run(context.Background(), t.TempDir(), "changelog", "prompt")
	if err != nil {
		t.Fatal(err)
	}
	if out.Result != "ok" || out.CostUSD != 1.5 {
		t.Errorf("got %+v", out)
	}
}
