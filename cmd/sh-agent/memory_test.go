package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// TestDreamVerb: the verb reaches the package with its flags, --project taking
// a value (takesValue), and refuses neither-nor with the usage line.
func TestDreamVerb(t *testing.T) {
	home := t.TempDir()
	if err := run([]string{"dream", "--home", home}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("bare dream: %v", err)
	}
	var out bytes.Buffer
	if err := run([]string{"dream", "--all", "--dry-run", "--home", home, "--vendor", "claude"}, &out); err != nil {
		t.Fatalf("--all --dry-run on an empty home: %v", err)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("a dry run on an empty home wrote %v", entries)
	}
	if err := run([]string{"dream", "--project", t.TempDir(), "--dry-run", "--home", home}, &out); err == nil || !strings.Contains(err.Error(), "not a git checkout") {
		t.Errorf("--project outside a checkout: %v", err)
	}
}

// TestUninstallDryRunLeavesTheTimer: a dry run names the timer and runs no
// systemctl; the real uninstall removes it, with no manifest left.
func TestUninstallDryRunLeavesTheTimer(t *testing.T) {
	home := t.TempDir()
	unit := filepath.Join(home, ".config", "systemd", "user", "sh-agent-dream.timer")
	os.MkdirAll(filepath.Dir(unit), 0o755)
	os.WriteFile(unit, []byte("x"), 0o644)
	bin := t.TempDir()
	calls := filepath.Join(t.TempDir(), "calls")
	os.WriteFile(filepath.Join(bin, "systemctl"), []byte("#!/bin/sh\necho \"$*\" >> \""+calls+"\"\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "launchctl"), []byte("#!/bin/sh\necho \"$*\" >> \""+calls+"\"\n"), 0o755)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	var out bytes.Buffer
	run([]string{"uninstall", "--dry-run", "--home", home}, &out)
	if !strings.Contains(out.String(), "would remove the dream timer") {
		t.Errorf("dry run output:\n%s", out.String())
	}
	if _, err := os.Stat(calls); !os.IsNotExist(err) {
		t.Error("a dry run called the service manager")
	}
	if _, err := os.Stat(unit); err != nil {
		t.Error("a dry run removed the unit")
	}
	out.Reset()
	run([]string{"uninstall", "--home", home}, &out)
	if runtime.GOOS == "linux" {
		if _, err := os.Stat(unit); !os.IsNotExist(err) {
			t.Errorf("uninstall left the timer (no manifest):\n%s", out.String())
		}
	}
}
