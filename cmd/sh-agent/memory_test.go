package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMemoryVerb runs `memory add` then `memory show` the way a session on
// either CLI calls them: flags after the subcommand, a description with
// spaces, the body on stdin.
func TestMemoryVerb(t *testing.T) {
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".agents/memory\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	home := t.TempDir()
	stdin = strings.NewReader("We deploy on Fridays.\n")
	defer func() { stdin = os.Stdin }()

	var out bytes.Buffer
	args := []string{"memory", "add", "--type", "project", "--name", "deploy-day", "--description", "we deploy on Fridays", "--home", home}
	if err := run(args, &out); err != nil {
		t.Fatalf("memory add: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "deploy-day.md") {
		t.Errorf("add printed %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "memory", "deploy-day.md")); err != nil {
		t.Fatal(err)
	}

	out.Reset()
	if err := run([]string{"memory", "show", "--home", home, "--vendor", "agy"}, &out); err != nil {
		t.Fatalf("memory show: %v", err)
	}
	if got := out.String(); got != "- [deploy-day](deploy-day.md) — we deploy on Fridays\n" {
		t.Errorf("show = %q", got)
	}
}

// TestMemoryRefusalExitsTwo: a refusal is exit 2, not an error the session
// would read as a crash.
func TestMemoryRefusalExitsTwo(t *testing.T) {
	t.Chdir(t.TempDir()) // not a git checkout
	stdin = strings.NewReader("fact\n")
	defer func() { stdin = os.Stdin }()
	err := run([]string{"memory", "add", "--type", "user", "--name", "x", "--description", "d", "--home", t.TempDir()}, &bytes.Buffer{})
	var code exitCode
	if !errors.As(err, &code) || code != 2 {
		t.Fatalf("err = %v, want exit 2", err)
	}
}
