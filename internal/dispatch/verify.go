package dispatch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/trail"
)

// runVerifier satisfies the verification row.
//
// Verification is the one Definition-of-Done row that needs no judgement: it
// runs the project's own commands and reports whether they passed. Spending a
// model to relay an exit code is the clearest case of a subagent working for
// no benefit. A model is invoked only when something fails, where explaining
// the failure is genuinely useful work.
func runVerifier(proj trail.Project, cfg trail.Config, cli cliTool, day string) Result {
	start := time.Now()
	cmds := verifyCommands(proj.Root)
	res := Result{Row: trail.RowVerifier}

	if len(cmds) == 0 {
		res.Duration = time.Since(start).Round(time.Second).String()
		res.Detail = "no test command detected — declare `verify` in .agents/harness.json"
		return res
	}

	// Print the commands before running them. Verifying a repository means
	// executing that repository's code — its test files at minimum, and
	// whatever `verify` in .agents/harness.json declares. That is inherent to
	// verification, not a flaw to be papered over; the microVM is the boundary
	// that makes it acceptable. What is not acceptable is doing it silently, so
	// the operator always sees exactly what ran.
	fmt.Printf("verifier: running %d command(s) from %s\n", len(cmds), verifySource(proj.Root))
	for _, cmd := range cmds {
		fmt.Printf("  $ %s\n", cmd)
	}

	var log strings.Builder
	passed := true
	for _, cmd := range cmds {
		ctx, cancel := context.WithTimeout(context.Background(),
			time.Duration(cfg.DispatchTimeout)*time.Second)
		out, err := runShell(ctx, proj.Root, cmd)
		cancel()

		status := "PASS"
		if err != nil {
			status, passed = "FAIL", false
		}
		fmt.Fprintf(&log, "\n## `%s` — %s\n\n```\n%s\n```\n", cmd, status, tail(out, 4000))
		if err != nil {
			break // the first failure is the one worth reporting
		}
	}
	res.Duration = time.Since(start).Round(time.Second).String()

	if !passed {
		// Hand the failure to a model to explain and record — the only path
		// where the verification row costs anything. The captured output goes
		// with it: the verifier agent is forbidden from blind-rerunning the
		// suite, so without the log it would have nothing to work from.
		res.Detail = "tests failed — escalating to the verifier agent to explain"
		esc := runAgentRow(proj, cfg, cli, trail.RowVerifier, day, log.String())
		esc.Duration = time.Since(start).Round(time.Second).String()
		if esc.OK {
			esc.OK = false // the artifact exists, but verification still failed
			esc.Detail = "FAIL — see " + esc.Artifact
		}
		return esc
	}

	path, err := writeVerificationLog(proj, day, cmds, log.String())
	if err != nil {
		res.Detail = "tests passed but the log could not be written: " + err.Error()
		return res
	}
	res.OK = true
	res.Artifact = path
	res.Detail = "all commands passed"
	return res
}

// verifyCommands returns the project's verification commands: whatever
// .agents/harness.json declares, otherwise what the repo's own manifests imply.
func verifyCommands(root string) []string {
	if declared := declaredVerify(root); len(declared) > 0 {
		return declared
	}
	var cmds []string
	if exists(filepath.Join(root, "go.mod")) {
		cmds = append(cmds, "go vet ./...", "go test -race ./...")
	}
	if exists(filepath.Join(root, "package.json")) {
		cmds = append(cmds, "npm test --if-present")
	}
	return cmds
}

// verifySource names where the commands came from, so a repo-declared command
// is never mistaken for one the harness chose.
func verifySource(root string) string {
	if len(declaredVerify(root)) > 0 {
		return ".agents/harness.json (repo-declared)"
	}
	return "the repo's own manifests (auto-detected)"
}

// declaredVerify reads the optional `verify` array from .agents/harness.json.
func declaredVerify(root string) []string {
	type harnessFile struct {
		Verify []string `json:"verify"`
	}
	data, err := os.ReadFile(filepath.Join(root, ".agents", "harness.json"))
	if err != nil {
		return nil
	}
	var hf harnessFile
	if err := json.Unmarshal(data, &hf); err != nil {
		return nil
	}
	return hf.Verify
}

func runShell(ctx context.Context, dir, cmd string) (string, error) {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	c.Dir = dir
	out, err := c.CombinedOutput()
	return string(out), err
}

// writeVerificationLog records a passing run. The script owns this artifact:
// it did the work, so it writes the record.
func writeVerificationLog(proj trail.Project, day string, cmds []string, body string) (string, error) {
	dir := filepath.Join(proj.Docs, "testing")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, day+"-verify-local.md")
	content := fmt.Sprintf(`---
owner: verifier
date: %s
verdict: PASS
---

# Verification — local tier   (%s)

Produced by `+"`harness dispatch`"+`, which ran the project's declared commands
directly. No model was involved: every command below exited zero.

Commands: %s
%s`, day, day, "`"+strings.Join(cmds, "`, `")+"`", body)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}
	rel, err := filepath.Rel(proj.Root, path)
	if err != nil {
		return path, nil
	}
	return rel, nil
}

// tail keeps the last n bytes of command output; a failing test suite's useful
// part is at the end.
func tail(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return "...\n" + s[len(s)-n:]
}
