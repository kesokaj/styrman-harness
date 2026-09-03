package gate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kesokaj/styrman-harness/internal/hookio"
	"github.com/kesokaj/styrman-harness/internal/trail"
)

func TestIsSecretPath(t *testing.T) {
	cases := map[string]bool{
		".env":                    true,
		".env.production":         true,
		"config/prod.tfvars":      true,
		"keys/id_rsa":             true,
		"sa/service-account.json": true,
		"certs/server.pem":        true,
		"main.go":                 false,
		"environment.ts":          false,
		"README.md":               false,
	}
	for path, want := range cases {
		if got := isSecretPath(path); got != want {
			t.Errorf("isSecretPath(%q) = %v, want %v", path, got, want)
		}
	}
}

// newSecretRepo builds a repo containing a credential-shaped file, plus an
// ordinary file and a template that must NOT be treated as one.
func newSecretRepo(t *testing.T) string {
	t.Helper()
	root := newContractRepo(t)
	for name, body := range map[string]string{
		".env":         "SECRET=1",
		".env.example": "SECRET=",
		"main.go":      "package main\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "secrets"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestSecretDecisionBypasses covers every spelling a security review used to
// defeat the argument-parsing version of this gate. Parsing shell and git
// syntax was the mistake; the gate now asks git what it actually holds.
func TestSecretDecisionBypasses(t *testing.T) {
	root := newSecretRepo(t)
	deny := []string{
		"git add -A",
		"git -C . add -A",
		"git -c user.name=x add -A",
		"git add .env",
		"git add ./",
		"git add .",
		"git add -u",
		"git add .e*",
		"git add secrets/",
		"git add .en'v'",
		"git commit -am wip",
		"echo hi && git add -A",
		"bash -c 'git add .env'",
	}
	for _, cmd := range deny {
		t.Run("deny/"+cmd, func(t *testing.T) {
			got := SecretDecision(hookio.Payload{CWD: root, ToolInput: map[string]any{"command": cmd}})
			if got.Verb != "deny" {
				t.Errorf("%q was allowed: %+v", cmd, got)
			}
		})
	}
}

// TestSecretDecisionDoesNotFailClosed: a gate that denies ordinary work is a
// gate that gets turned off.
func TestSecretDecisionDoesNotFailClosed(t *testing.T) {
	root := newSecretRepo(t)
	allow := []string{
		"git add main.go", // names an ordinary file
		"git commit -m \"handle credentials.json\"", // prose, not a path
		"git log --oneline",                         // not a staging command
		"cat .env",                                  // reads, does not stage
		"git status",
	}
	for _, cmd := range allow {
		t.Run("allow/"+cmd, func(t *testing.T) {
			got := SecretDecision(hookio.Payload{CWD: root, ToolInput: map[string]any{"command": cmd}})
			if got.Verb != "" {
				t.Errorf("%q was denied: %+v", cmd, got)
			}
		})
	}
}

// TestSecretDecisionCatchesStagedSecret: once a credential is in the index, any
// later commit ships it — and that commit need not look like a staging command.
func TestSecretDecisionCatchesStagedSecret(t *testing.T) {
	root := newSecretRepo(t)
	cmd := exec.Command("git", "add", "-f", ".env")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("staging fixture: %v: %s", err, out)
	}

	got := SecretDecision(hookio.Payload{CWD: root, ToolInput: map[string]any{"command": "git commit -m wip"}})
	if got.Verb != "deny" {
		t.Fatalf("a commit with .env already staged was allowed: %+v", got)
	}
	if !strings.Contains(got.Reason, ".env") {
		t.Errorf("denial should name the staged file: %s", got.Reason)
	}
}

// TestSecretDecisionReadsAntigravityPayload: the gate is wired on both vendors,
// and Antigravity spells its tool arguments differently.
func TestSecretDecisionReadsAntigravityPayload(t *testing.T) {
	root := newSecretRepo(t)
	var p hookio.Payload
	p.WorkspacePaths = []string{root}
	p.ToolCall.Args = map[string]any{"Command": "git add -A"}

	if got := SecretDecision(p); got.Verb != "deny" {
		t.Errorf("Antigravity payload shape was not read: %+v", got)
	}
}

// TestExampleFilesAreNotSecrets: .env.example exists to be committed.
func TestExampleFilesAreNotSecrets(t *testing.T) {
	for _, name := range []string{".env.example", ".env.sample", "prod.tfvars.template"} {
		if isSecretPath(name) {
			t.Errorf("%s was treated as a credential", name)
		}
	}
}

// TestSubagentStopReleasesAfterCap is the guard against the failure this whole
// redesign started from: a hook that blocks unconditionally costs eight forced
// continuations before the vendor gives up. The second guard must release the
// agent even when stop_hook_active never arrives.
func TestSubagentStopReleasesAfterCap(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	proj := newContractRepo(t)

	p := hookio.Payload{
		AgentType: "reviewer",
		AgentID:   "agent-under-test",
		CWD:       proj,
	}
	first := SubagentStopDecision(p)
	if first.Verb != "block" {
		t.Fatalf("first stop without an artifact should hold once, got %+v", first)
	}
	second := SubagentStopDecision(p)
	if second.Verb != "" {
		t.Errorf("second stop must release rather than loop, got %+v", second)
	}
}

func TestSubagentStopHonoursStopHookActive(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	p := hookio.Payload{
		AgentType:      "reviewer",
		AgentID:        "a",
		CWD:            newContractRepo(t),
		StopHookActive: true,
	}
	if got := SubagentStopDecision(p); got.Verb != "" {
		t.Errorf("stop_hook_active must release immediately, got %+v", got)
	}
}

func TestSubagentStopIgnoresNonOwners(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	p := hookio.Payload{AgentType: "general-purpose", AgentID: "a", CWD: newContractRepo(t)}
	if got := SubagentStopDecision(p); got.Verb != "" {
		t.Errorf("an agent owning no row must never be held, got %+v", got)
	}
}

// newContractRepo returns a git repo that has adopted the trail contract.
func newContractRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents", "docs", "changes"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestBlockMessageNamesTheScopeThatClearsIt is the regression for the failure
// that ended five rounds of the same bug: coverage was decided with one scope
// function and the message printed another, so following the instruction
// exactly left the gate blocking. The gate and the dispatcher now share one
// computation, and this asserts the property that matters — do what it says,
// and it stops saying it.
func TestBlockMessageNamesTheScopeThatClearsIt(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	root := newContractRepo(t)
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	run("config", "user.email", "t@e.com")
	run("config", "user.name", "t")
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", "-A")
	run("commit", "-q", "-m", "unrecorded work")

	p := hookio.Payload{CWD: root}
	first := StopDecision(p)
	if first.Verb != "block" {
		t.Fatalf("an unrecorded commit did not block: %+v", first)
	}

	// Extract the scope the message demands and write exactly that entry.
	idx := strings.Index(first.Reason, "changeset:")
	if idx < 0 {
		t.Fatalf("the block message names no scope: %s", first.Reason)
	}
	scope := strings.Fields(strings.TrimPrefix(first.Reason[idx:], ""))[0]
	scope = strings.Trim(scope, "`—,. ")

	entry := filepath.Join(root, ".agents", "docs", "changes", "recorded.md")
	body := "---\nowner: changelog\ndate: " + trail.Today() + "\nscope: " + scope + "\n---\n"
	if err := os.WriteFile(entry, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	if second := StopDecision(p); second.Verb != "" {
		t.Errorf("following the block message exactly did not clear it: %+v", second)
	}
}
