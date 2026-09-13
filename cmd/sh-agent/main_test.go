package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

func TestReorder(t *testing.T) {
	cases := map[string][]string{
		// BUG-002: the flag package stops at the first positional argument.
		"adopt my-skill --home /x":          {"--home", "/x", "adopt", "my-skill"},
		"--vendor agy install":              {"--vendor", "agy", "install"},
		"adopt a --force":                   {"--force", "adopt", "a"},
		"adopt -- --looks-like-a-flag":      {"adopt", "--looks-like-a-flag"},
		"--repo r --home h adopt x --force": {"--repo", "r", "--home", "h", "--force", "adopt", "x"},
	}
	for in, want := range cases {
		if got := reorder(strings.Fields(in)); !reflect.DeepEqual(got, want) {
			t.Errorf("reorder(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestResolveVendor(t *testing.T) {
	home := t.TempDir()
	if v, err := resolveVendor("", home); err != nil || v != harness.VendorAll {
		t.Errorf("no flag, no manifest: %q, %v; want all", v, err)
	}
	if v, err := resolveVendor("agy", home); err != nil || v != harness.VendorAntigravity {
		t.Errorf("flag agy: %q, %v", v, err)
	}
	if _, err := resolveVendor("vscode", home); err == nil {
		t.Error("unknown vendor should fail")
	}
	if err := harness.SaveState(home, harness.State{Vendor: harness.VendorClaude}); err != nil {
		t.Fatal(err)
	}
	if v, err := resolveVendor("", home); err != nil || v != harness.VendorClaude {
		t.Errorf("manifest vendor: %q, %v; want claude", v, err)
	}
	if v, err := resolveVendor("all", home); err != nil || v != harness.VendorAll {
		t.Errorf("flag beats manifest: %q, %v", v, err)
	}
}

func TestShouldRebuild(t *testing.T) {
	cases := []struct{ goOnPath, child, want bool }{
		{true, false, true},
		{true, true, false},
		{false, false, false},
		{false, true, false},
	}
	for _, c := range cases {
		if got := shouldRebuild(c.goOnPath, c.child); got != c.want {
			t.Errorf("shouldRebuild(%v, %v) = %v", c.goOnPath, c.child, got)
		}
	}
}

func TestBinaryOlder(t *testing.T) {
	cases := []struct {
		rev, commit string
		want        bool
	}{
		{"b64ffe886729", "b64ffe8", false},
		{"b64ffe8", "b64ffe886729", false},
		{"b64ffe886729+dirty", "b64ffe8", false},
		{"50ac771abcde", "b64ffe8", true},
		{"(devel)", "b64ffe8", true},
		{"unknown", "b64ffe8", true},
		{"", "b64ffe8", false},
		{"b64ffe886729", "", false},
	}
	for _, c := range cases {
		if got := binaryOlder(c.rev, c.commit); got != c.want {
			t.Errorf("binaryOlder(%q, %q) = %v", c.rev, c.commit, got)
		}
	}
}

// TestRunInstallStatusUninstall drives the CLI end to end against this
// checkout and a scratch home. ChildEnv is set so install takes the plain
// path: under `go test` the running binary is the test binary, and a rebuild
// here would be a real `go build` of the repo for nothing.
func TestRunInstallStatusUninstall(t *testing.T) {
	t.Setenv(harness.ChildEnv, "1")
	home := t.TempDir()
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"install", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("install: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "files managed") {
		t.Errorf("install output:\n%s", out.String())
	}
	if _, err := os.Stat(filepath.Join(home, harness.BinPath)); err != nil {
		t.Errorf("binary not installed: %v", err)
	}
	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("status: %v\n%s", err, out.String())
	}
	for _, want := range []string{"binary     ", "drift      none"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("status missing %q:\n%s", want, out.String())
		}
	}
	out.Reset()
	if err := run([]string{"uninstall", "--home", home}, &out); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	// The vendor roots themselves are kept on purpose (pruneEmptyDirs stops
	// there); everything the manifest listed under them is gone.
	for _, rel := range []string{filepath.Join(harness.ClaudeRoot, "AGENTS.md"), filepath.Join(harness.ClaudeRoot, "skills"), filepath.Join(harness.AgyRoot, "agents"), harness.BinPath} {
		if _, err := os.Stat(filepath.Join(home, rel)); !os.IsNotExist(err) {
			t.Errorf("~/%s should be gone after uninstall, stat err = %v", rel, err)
		}
	}
	out.Reset()
	if err := run([]string{"status", "--repo", repo, "--home", home}, &out); err == nil {
		t.Error("status after uninstall should report exit 1")
	}
}
