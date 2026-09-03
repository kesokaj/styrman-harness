package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDeriveMapsToolsAndTiers(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "reviewer.md")
	write(t, src, `---
name: reviewer
description: Reviews code.
tools: Read, Grep, Bash, Write
model: opus
skills: review-standards
maxTurns: 30
---
Body referencing ~/.claude/skills/review-standards.
`)
	dst := filepath.Join(dir, "out.md")
	if err := DeriveFile(src, dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	out := string(got)

	for _, want := range []string{
		"subagent: true",
		"mainAgent: false",
		"model: pro",
		"tools: [view_file, grep_search, run_command, write_to_file]",
		"skills: [~/.gemini/config/skills/review-standards]",
		"~/.gemini/config/skills/review-standards.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("derived agent missing %q\n---\n%s", want, out)
		}
	}
	// Claude-only keys must not leak: agy would not understand them.
	if strings.Contains(out, "maxTurns") {
		t.Error("Claude-only key maxTurns leaked into the derived agent")
	}
}

// TestDeriveRejectsUnmappedTool guards the documented Antigravity failure mode:
// an unrecognised tool name hangs the subagent at dispatch, so an unmapped tool
// must fail the build rather than pass through.
func TestDeriveRejectsUnmappedTool(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.md")
	write(t, src, "---\nname: a\ndescription: d\ntools: Read, WebFetch\n---\nbody\n")

	err := DeriveFile(src, filepath.Join(dir, "out.md"))
	if err == nil {
		t.Fatal("unmapped tool was silently accepted")
	}
	if !strings.Contains(err.Error(), "WebFetch") {
		t.Errorf("error should name the offending tool, got: %v", err)
	}
}

func TestDeriveRequiresFrontmatter(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.md")
	write(t, src, "no frontmatter here\n")
	if err := DeriveFile(src, filepath.Join(dir, "out.md")); err == nil {
		t.Error("a file without frontmatter was accepted")
	}
}

// TestMirrorPrunesExtras is why the installer mirrors rather than copies: a
// skill deleted from the repo has to disappear from the deployment too.
func TestMirrorPrunesExtras(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	write(t, filepath.Join(src, "keep", "SKILL.md"), "keep me")
	write(t, filepath.Join(dst, "keep", "SKILL.md"), "stale content")
	write(t, filepath.Join(dst, "retired", "SKILL.md"), "should be pruned")

	if err := mirrorDir(src, dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dst, "keep", "SKILL.md"))
	if err != nil || string(got) != "keep me" {
		t.Errorf("mirror did not overwrite: %q, %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(dst, "retired")); !os.IsNotExist(err) {
		t.Error("mirror left a directory the source no longer has")
	}
}

func TestCollectCommandsFindsNestedHooks(t *testing.T) {
	doc := map[string]any{
		"hooks": map[string]any{
			"PreToolUse": []any{
				map[string]any{"hooks": []any{
					map[string]any{"type": "command", "command": "$HOME/.agents/bin/harness hook pretooluse"},
				}},
			},
		},
		"statusLine": map[string]any{"command": "~/.agents/bin/harness hook statusline"},
	}
	got := collectCommands(doc)
	if len(got) != 2 {
		t.Fatalf("collectCommands found %d commands, want 2: %v", len(got), got)
	}
}

func TestExecutablePathExpandsHomeForms(t *testing.T) {
	home := "/home/u"
	cases := map[string]string{
		"~/.agents/bin/harness hook stop":     "/home/u/.agents/bin/harness",
		"$HOME/.agents/bin/harness hook stop": "/home/u/.agents/bin/harness",
		"/usr/bin/true":                       "/usr/bin/true",
		"jq -r .":                             "", // resolved via PATH, not checkable
	}
	for cmd, want := range cases {
		if got := executablePath(home, cmd); got != want {
			t.Errorf("executablePath(%q) = %q, want %q", cmd, got, want)
		}
	}
}

// TestIsRealHomeGuardsExternalSteps: agy has no --home and writes to the real
// home regardless, so any install step that shells out to it must be gated on
// this. A scratch-home install once reached out and modified the operator's
// actual Antigravity configuration.
func TestIsRealHomeGuardsExternalSteps(t *testing.T) {
	real, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home directory")
	}
	if !(Env{Home: real}).isRealHome() {
		t.Error("the real home was not recognised as real")
	}
	if (Env{Home: t.TempDir()}).isRealHome() {
		t.Error("a scratch directory was treated as the real home")
	}
	// Trailing separators and unclean paths must not defeat the check.
	if !(Env{Home: real + "/."}).isRealHome() {
		t.Error("an unclean spelling of the real home was not recognised")
	}
}
