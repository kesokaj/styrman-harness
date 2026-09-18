package harness

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// newRepo builds a minimal harness checkout: two skills, one agent, two
// context files. Tests never touch the real $HOME.
func newRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		p := filepath.Join(repo, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("context/AGENTS.md", "# AGENTS.md\n\nBe honest.\n")
	write("context/dev-setup.md", "## Environment & Stack\n\n- GCP only.\n")
	write("skills/alpha/SKILL.md", "---\nname: alpha\ndescription: |\n  Alpha does things.\n  Use it for alpha work.\nmetadata:\n  version: \"1\"\n---\nbody\n")
	write("skills/alpha/references/notes.md", "notes\n")
	write("skills/beta/SKILL.md", "---\nname: beta\ndescription: \"Beta: slash only\"\ndisable-model-invocation: true\n---\nbody\n")
	write("skills/VENDORED.md", "not a skill\n")
	write("commands/deploy.md", "---\ndescription: Deploy the thing\n---\nSteps\n")
	write("rules/go-style.md", "# Go style\n\nUse slog.\n")
	write("agents/reviewer.md", "---\nname: reviewer\ndescription: Reviews diffs\ntools: Read, Grep, Bash\nmodel: opus\nskills: review-standards\neffort: high\n---\nLoad ~/.claude/skills/review-standards first.\n")
	return repo
}

func TestParseFrontmatter(t *testing.T) {
	front := "name: x\ndescription: |\n  line one\n  line two\nquoted: \"a: b\"\nmetadata:\n  version: \"1.2\"\n  author: me\ntools: [Read, Grep]\nskills:\n  - alpha\n  - beta\nplain: first\n  second\nempty:\n"
	fm := Parse(front)
	want := map[string]string{
		"name":        "x",
		"description": "line one line two",
		"quoted":      "a: b",
		"tools":       "[Read, Grep]",
		"skills":      "[alpha, beta]",
		"plain":       "first second",
		"empty":       "",
	}
	for k, v := range want {
		if fm[k] != v {
			t.Errorf("%s = %q, want %q", k, fm[k], v)
		}
	}
	if _, ok := fm["metadata"]; ok {
		t.Errorf("nested map should be skipped, got %q", fm["metadata"])
	}
	if got := SplitList("[Read, Grep]"); strings.Join(got, ",") != "Read,Grep" {
		t.Errorf("SplitList brackets = %v", got)
	}
	if got := SplitList("Read Grep"); strings.Join(got, ",") != "Read,Grep" {
		t.Errorf("SplitList spaces = %v", got)
	}
}

func TestLintRejectsUnquotedComment(t *testing.T) {
	// A plain scalar with " #" is a YAML comment start: strict parsers
	// truncate it or, inside a flow sequence, fail the whole file. The
	// Antigravity CLI dropped the `bug` and `feature` skills that way.
	bad := []string{
		"name: x\nargument-hint: [feature description or #issue]\n",
		"name: x\ndescription: fixes #1 and #2\n",
		"name: x\ndescription: #comment only\n",
		"name: x\ndescription: first line\n  wraps with #issue\n",
		"name: x\ndescription: value\t#tab-comment\n",
		"name: x\ndescription: \"quoted\" # trailing\n",
	}
	for _, front := range bad {
		if err := Lint(front); err == nil {
			t.Errorf("Lint accepted %q", front)
		} else if !strings.Contains(err.Error(), "#") {
			t.Errorf("error should name the marker: %v", err)
		}
	}
	good := []string{
		"name: x\nargument-hint: \"[feature description or #issue]\"\n",
		"name: x\nargument-hint: '[#issue]'\n",
		"name: x\ndescription: |\n  a #b\n  c #d\nmetadata:\n  tag: \"#1\"\n",
		"name: x\ntools: [Read, Grep]\ncolor: \"#ff0000\"\n",
	}
	for _, front := range good {
		if err := Lint(front); err != nil {
			t.Errorf("Lint rejected %q: %v", front, err)
		}
	}
}

func TestLoadSkillsRejectsUnquotedComment(t *testing.T) {
	repo := newRepo(t)
	p := filepath.Join(repo, "skills", "gamma", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("---\nname: gamma\ndescription: Gamma\nargument-hint: [topic or #issue]\n---\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadSkills(DirSource(repo))
	if err == nil || !strings.Contains(err.Error(), "gamma") || !strings.Contains(err.Error(), "argument-hint") {
		t.Fatalf("LoadSkills should name the skill and the key, got %v", err)
	}
}

func TestSplitRejectsMissingFrontmatter(t *testing.T) {
	if _, _, err := Split("# no frontmatter\n"); err == nil {
		t.Fatal("expected an error")
	}
	if _, _, err := Split("---\nname: x\n"); err == nil {
		t.Fatal("expected an error for an unterminated block")
	}
}

func TestDeriveAgy(t *testing.T) {
	repo := newRepo(t)
	agents, err := LoadAgents(DirSource(repo))
	if err != nil {
		t.Fatal(err)
	}
	out, err := DeriveAgy(agents[0])
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{
		"name: reviewer\n",
		"description: \"Reviews diffs\"\n",
		"subagent: true\n",
		"mainAgent: false\n",
		"model: pro\n",
		"tools: [view_file, grep_search, run_command]\n",
		"~/.gemini/config/skills/review-standards",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("derived agent missing %q:\n%s", want, s)
		}
	}
	for _, unwanted := range []string{"effort:", "skills:", "~/.claude/"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("derived agent should not contain %q:\n%s", unwanted, s)
		}
	}

	bad := agents[0]
	bad.Front = Frontmatter{"tools": "Read, WebSearch"}
	if _, err := DeriveAgy(bad); err == nil {
		t.Fatal("unmapped tool should be an error")
	}
}

func TestInventoryAndRender(t *testing.T) {
	repo := newRepo(t)
	skills, err := LoadSkills(DirSource(repo))
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 2 || skills[0].Name != "alpha" || skills[1].Name != "beta" {
		t.Fatalf("skills = %+v", skills)
	}
	if !skills[1].SlashOnly || skills[0].SlashOnly {
		t.Errorf("SlashOnly wrong: %+v", skills)
	}
	inv, err := LoadInventory(DirSource(repo))
	if err != nil {
		t.Fatal(err)
	}
	if len(inv.Commands) != 1 || inv.Commands[0].Name != "deploy" || inv.Commands[0].Description != "Deploy the thing" {
		t.Errorf("commands = %+v", inv.Commands)
	}
	if len(inv.Rules) != 1 || inv.Rules[0].Description != "Go style" {
		t.Errorf("rules = %+v", inv.Rules)
	}
	ctx, err := RenderContext(DirSource(repo), inv)
	if err != nil {
		t.Fatal(err)
	}
	s := string(ctx)
	if !strings.HasPrefix(s, Marker+"\n# AGENTS.md") {
		t.Errorf("context should start with the marker and AGENTS.md:\n%s", s[:80])
	}
	for _, want := range []string{"## Environment & Stack", "`alpha`", "`/beta`", "`reviewer` (opus)", "`/deploy`", "`go-style`"} {
		if !strings.Contains(s, want) {
			t.Errorf("context missing %q", want)
		}
	}
	if strings.Contains(s, "Alpha does things") {
		t.Errorf("deployed context should carry names only, not descriptions")
	}
	full := RenderInventory(inv)
	for _, want := range []string{"| `alpha` | auto |", "| `beta` | `/beta` |", "| `reviewer` | opus |", "| `/deploy` |", "| `go-style` |", "Beta: slash only"} {
		if !strings.Contains(full, want) {
			t.Errorf("full inventory missing %q", want)
		}
	}
}

func TestInstallStatusResyncUninstall(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()
	self := []byte("#!/bin/sh\necho sh-agent\n")

	// A CLAUDE.md the operator wrote must survive an install.
	claudeMD := filepath.Join(home, ClaudeRoot, "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(claudeMD, []byte("my own notes\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := Plan(DirSource(repo), self, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Install(DirSource(repo), home, files, VendorAll, false, &out); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		".claude/AGENTS.md",
		".claude/skills/alpha/SKILL.md",
		".claude/skills/alpha/references/notes.md",
		".claude/skills/beta/SKILL.md",
		".claude/agents/reviewer.md",
		".gemini/config/rules/AGENTS.md",
		".gemini/config/skills/alpha/SKILL.md",
		".gemini/config/agents/reviewer/agent.md",
		".claude/commands/deploy.md",
		".gemini/config/workflows/deploy.md",
		".claude/rules/go-style.md",
		".gemini/config/rules/go-style.md",
		".local/bin/sh-agent",
		".config/sh-agent/state.json",
	} {
		if _, err := os.Stat(filepath.Join(home, rel)); err != nil {
			t.Errorf("expected %s after install: %v", rel, err)
		}
	}
	if got, _ := os.ReadFile(claudeMD); string(got) != "my own notes\n" {
		t.Errorf("foreign CLAUDE.md was overwritten: %q", got)
	}
	if !strings.Contains(out.String(), "keep    ~/.claude/CLAUDE.md") {
		t.Errorf("install should report the kept file:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "move only rules meant for every machine into a rule file under rules/") || !strings.Contains(out.String(), "never credentials or personal or host-specific notes") {
		t.Fatalf("keep line lacks the rules/ guidance or its exclusion of credentials and personal notes:\n%s", out.String())
	}
	if info, _ := os.Stat(filepath.Join(home, BinPath)); info.Mode().Perm()&0o100 == 0 {
		t.Errorf("binary should be executable, mode %v", info.Mode())
	}

	report, err := Status(home, files)
	if err != nil {
		t.Fatal(err)
	}
	if report.Drift() || report.Counts["foreign"] != 1 {
		t.Errorf("fresh install should be clean apart from the foreign CLAUDE.md: %+v", report.Counts)
	}
	// status and diff plan without the binary; the installed one is not stale.
	noSelf, err := Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if report, _ := Status(home, noSelf); report.Drift() {
		t.Errorf("installed binary must not count as stale: %+v", report.Counts)
	}

	// Edit a deployed file, delete another: both show up as drift.
	if err := os.WriteFile(filepath.Join(home, ".claude/skills/alpha/SKILL.md"), []byte("edited"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(home, ".gemini/config/agents/reviewer/agent.md")); err != nil {
		t.Fatal(err)
	}
	report, _ = Status(home, files)
	if report.Counts["differs"] != 1 || report.Counts["missing"] != 1 {
		t.Errorf("drift not detected: %+v", report.Counts)
	}

	// Remove a skill from the repo and reinstall: the deployed copies go away.
	if err := os.RemoveAll(filepath.Join(repo, "skills", "beta")); err != nil {
		t.Fatal(err)
	}
	files, err = Plan(DirSource(repo), self, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := Install(DirSource(repo), home, files, VendorAll, false, &out); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{".claude/skills/beta", ".gemini/config/skills/beta"} {
		if _, err := os.Stat(filepath.Join(home, rel)); err == nil {
			t.Errorf("stale %s should have been removed", rel)
		}
	}
	if !strings.Contains(out.String(), "2 removed") {
		t.Errorf("resync should report removals:\n%s", out.String())
	}
	report, _ = Status(home, files)
	if report.Drift() {
		t.Errorf("after resync there should be no drift: %+v", report.Counts)
	}

	// Dry-run changes nothing.
	if err := Uninstall(home, true, &out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude/AGENTS.md")); err != nil {
		t.Errorf("dry-run uninstall removed a file")
	}

	if err := Uninstall(home, false, &out); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{".claude/AGENTS.md", ".claude/skills", ".claude/agents", ".gemini/config/skills", ".gemini/config/rules", ".local/bin/sh-agent", ".config/sh-agent"} {
		if _, err := os.Lstat(filepath.Join(home, rel)); err == nil {
			t.Errorf("%s should be gone after uninstall", rel)
		}
	}
	for _, rel := range []string{".claude", ".gemini/config", ".claude/CLAUDE.md"} {
		if _, err := os.Lstat(filepath.Join(home, rel)); err != nil {
			t.Errorf("%s must survive uninstall: %v", rel, err)
		}
	}
	if got, _ := os.ReadFile(claudeMD); string(got) != "my own notes\n" {
		t.Errorf("foreign CLAUDE.md touched by uninstall: %q", got)
	}
}

func TestInstallWritesOwnClaudeStub(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()
	files, err := Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if err := Install(DirSource(repo), home, files, VendorAll, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(home, ClaudeRoot, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !IsOurs(got) || !strings.Contains(string(got), "\n@AGENTS.md\n") {
		t.Errorf("stub = %q", got)
	}
}

func TestVendorSelection(t *testing.T) {
	// NormalizeVendor checks
	for _, in := range []string{"", "all", "ALL", " all "} {
		if got, err := NormalizeVendor(in); err != nil || got != VendorAll {
			t.Errorf("NormalizeVendor(%q) = %q, err %v", in, got, err)
		}
	}
	for _, in := range []string{"antigravity", "Antigravity", "agy", "gemini"} {
		if got, err := NormalizeVendor(in); err != nil || got != VendorAntigravity {
			t.Errorf("NormalizeVendor(%q) = %q, err %v", in, got, err)
		}
	}
	for _, in := range []string{"claude", "Claude", "claudecode"} {
		if got, err := NormalizeVendor(in); err != nil || got != VendorClaude {
			t.Errorf("NormalizeVendor(%q) = %q, err %v", in, got, err)
		}
	}
	if _, err := NormalizeVendor("invalid-vendor"); err == nil {
		t.Error("NormalizeVendor should reject unknown vendor")
	}

	repo := newRepo(t)
	self := []byte("#!/bin/sh\necho sh-agent\n")

	// Plan for antigravity only
	agyPlan, err := Plan(DirSource(repo), self, VendorAntigravity)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range agyPlan {
		if strings.HasPrefix(f.Rel, ClaudeRoot) {
			t.Errorf("antigravity plan contains Claude file %s", f.Rel)
		}
	}
	foundAgyAgents := false
	for _, f := range agyPlan {
		if f.Rel == filepath.Join(AgyRoot, "rules", "AGENTS.md") {
			foundAgyAgents = true
		}
	}
	if !foundAgyAgents {
		t.Error("antigravity plan missing rules/AGENTS.md")
	}

	// Plan for claude only
	claudePlan, err := Plan(DirSource(repo), self, VendorClaude)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range claudePlan {
		if strings.HasPrefix(f.Rel, AgyRoot) {
			t.Errorf("claude plan contains Antigravity file %s", f.Rel)
		}
	}

	// Install antigravity-only to clean home
	home := t.TempDir()
	var out bytes.Buffer
	if err := Install(DirSource(repo), home, agyPlan, VendorAntigravity, false, &out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ClaudeRoot)); err == nil {
		t.Errorf(".claude should not exist after antigravity-only install")
	}
	if _, err := os.Stat(filepath.Join(home, AgyRoot, "rules", "AGENTS.md")); err != nil {
		t.Errorf("expected antigravity rules/AGENTS.md to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, BinPath)); err != nil {
		t.Errorf("expected binary to be installed: %v", err)
	}

	// Verify state records vendor
	st, found, err := LoadState(home)
	if err != nil || !found {
		t.Fatalf("LoadState failed: %v, found: %v", err, found)
	}
	if st.Vendor != VendorAntigravity {
		t.Errorf("state vendor = %q, want %q", st.Vendor, VendorAntigravity)
	}

	// Status check with antigravity plan has no drift
	noSelfAgy, err := Plan(DirSource(repo), nil, VendorAntigravity)
	if err != nil {
		t.Fatal(err)
	}
	report, err := Status(home, noSelfAgy)
	if err != nil {
		t.Fatal(err)
	}
	if report.Drift() {
		t.Errorf("antigravity status has drift: %+v", report.Counts)
	}

	// Switch from antigravity to all, then back to antigravity:
	// Verify resync removes Claude files when returning to antigravity
	allPlan, err := Plan(DirSource(repo), self, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if err := Install(DirSource(repo), home, allPlan, VendorAll, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ClaudeRoot, "AGENTS.md")); err != nil {
		t.Fatalf("expected .claude/AGENTS.md after installing all: %v", err)
	}

	// Now switch back to antigravity
	out.Reset()
	if err := Install(DirSource(repo), home, agyPlan, VendorAntigravity, false, &out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ClaudeRoot, "AGENTS.md")); err == nil {
		t.Errorf(".claude/AGENTS.md should have been removed by resync when switching back to antigravity")
	}
	st, _, err = LoadState(home)
	if err != nil {
		t.Fatal(err)
	}
	if st.Vendor != VendorAntigravity {
		t.Errorf("state vendor = %q, want %q", st.Vendor, VendorAntigravity)
	}
}

func TestAdopt(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()
	src := filepath.Join(home, ".agents", "skills", "gamma")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("---\nname: gamma\ndescription: g\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// A symlink the way `npx skills add -g` leaves one for Claude Code.
	link := filepath.Join(home, ClaudeRoot, "skills", "gamma")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(src, link); err != nil {
		t.Fatal(err)
	}

	dest, err := Adopt(repo, home, "gamma", false)
	if err != nil {
		t.Fatal(err)
	}
	if dest != filepath.Join(repo, "skills", "gamma") {
		t.Errorf("dest = %s", dest)
	}
	if _, err := os.Stat(filepath.Join(dest, "SKILL.md")); err != nil {
		t.Error(err)
	}
	if _, err := Adopt(repo, home, "gamma", false); err == nil {
		t.Error("second adopt without --force should fail")
	}
	if _, err := Adopt(repo, home, "gamma", true); err != nil {
		t.Errorf("adopt --force: %v", err)
	}
	if _, err := Adopt(repo, home, "nope", false); err == nil {
		t.Error("unknown skill should fail")
	}
}

// gitInit makes a directory its own git checkout, which is what IsCheckout
// reads and what separates a tree update can move forward from a plain copy.
func gitInit(t *testing.T, dir string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	cmd := exec.Command("git", "init", "--quiet", dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
}

func TestResolveSource(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	cwd := t.TempDir()
	t.Chdir(cwd)
	// The walk starts at the working directory, so a temp directory INSIDE a
	// harness checkout (a TMPDIR under one) would be found and every case
	// below would read that tree instead (security review 2026-09-17).
	if _, err := ResolveSource("", t.TempDir(), ""); err == nil {
		t.Skip("the temp directory happens to be inside a harness repo")
	}
	repo := newRepo(t)
	home := t.TempDir()

	// A development build with no tree found keeps the error: only a release
	// binary may install the files it carries (D14).
	for _, v := range []string{"", "b64ffe886729", "b64ffe886729+dirty", "(devel)", "v1.2.3-rc.1", "v1.2.3-0.20260917090000-b64ffe886729", "1.2.3", "v1.2"} {
		_, err := ResolveSource("", home, v)
		if err == nil || !strings.Contains(err.Error(), "harness repo not found") {
			t.Errorf("version %q with no tree: err = %v, want harness repo not found", v, err)
			continue
		}
		// The message says WHY the embedded fallback was refused — this build
		// has no harness in it — and names the release-binary way out. Saying
		// only "run from inside the checkout" sent the operator looking for a
		// repository they had been told they would not need (2026-09-18).
		if !strings.Contains(err.Error(), "no harness built in") || !strings.Contains(err.Error(), "release binary") {
			t.Errorf("version %q: err = %v, want it to name the missing harness and the release-binary path", v, err)
		}
		if v != "" && v != "unknown" && !strings.Contains(err.Error(), v) {
			t.Errorf("version %q: err = %v, want it to quote the version", v, err)
		}
	}

	src, err := ResolveSource("", home, "v1.2.3")
	if err != nil {
		t.Fatalf("release binary with no tree: %v", err)
	}
	if src.Dir != "" || src.Version != "v1.2.3" || src.FS == nil {
		t.Errorf("release binary with no tree: %+v, want the embedded source", src)
	}
	if _, err := fs.Stat(src.FS, "context/AGENTS.md"); err != nil {
		t.Errorf("the embedded source has no context/AGENTS.md: %v", err)
	}

	// The working-directory walk, for a development build: standing inside a
	// harness tree, or anywhere under it, finds that tree. Nothing pinned this
	// before, so deleting the walk kept the suite green (code review
	// 2026-09-18).
	nested := filepath.Join(repo, "skills", "alpha", "references")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{repo, nested} {
		t.Chdir(dir)
		if src, err := ResolveSource("", t.TempDir(), ""); err != nil || src.Dir != repo {
			t.Errorf("development build in %s: %+v, %v — want the tree walked up to", dir, src, err)
		}
	}
	t.Chdir(cwd)

	// A development build takes the path it is given.
	if src, err := ResolveSource(repo, home, ""); err != nil || src.Dir != repo {
		t.Errorf("explicit repo, development build: %+v, %v", src, err)
	}
	if err := SaveState(home, State{Repo: repo}); err != nil {
		t.Fatal(err)
	}
	// The recorded tree comes before the embedded files for a DEVELOPMENT
	// build, which has nothing of its own to offer.
	if src, err := ResolveSource("", home, ""); err != nil || src.Dir != repo {
		t.Errorf("repo from state, development build: %+v, %v", src, err)
	}
	// A RELEASE binary takes NO tree on this machine: not the recorded one,
	// whether a plain copy or a real checkout, and not the working directory.
	// What it installs is the published repository at the tagged commit, which
	// it carries; `update` replaces the whole binary to change that (ADR-034 —
	// the operator, after a box installed months-old content from a recorded
	// plain copy while reporting no drift: "the binary should only use the
	// published git repo in github always and only").
	if src, err := ResolveSource("", home, "v1.2.3"); err != nil || src.Dir != "" || src.Version != "v1.2.3" {
		t.Errorf("recorded plain copy, release binary: %+v, %v — want the embedded files", src, err)
	}
	checkout := newRepo(t)
	gitInit(t, checkout)
	if err := SaveState(home, State{Repo: checkout}); err != nil {
		t.Fatal(err)
	}
	if src, err := ResolveSource("", home, "v1.2.3"); err != nil || src.Dir != "" {
		t.Errorf("recorded checkout, release binary: %+v, %v — want the embedded files", src, err)
	}
	// Standing inside a tree does not change it either.
	t.Chdir(checkout)
	if src, err := ResolveSource("", home, "v1.2.3"); err != nil || src.Dir != "" {
		t.Errorf("inside a checkout, release binary: %+v, %v — want the embedded files", src, err)
	}
	// A named path is REFUSED, not ignored: doing something other than what
	// was typed is what D15 exists to prevent.
	// The refusal names WHICH input carried the path: an exported
	// $SH_AGENT_REPO makes every command on the machine fail, and "unset it"
	// is unactionable without knowing which one (code review 2026-09-18).
	if _, err := ResolveSource(checkout, home, "v1.2.3"); err == nil ||
		!strings.Contains(err.Error(), "release binary") || !strings.Contains(err.Error(), "--repo="+checkout) {
		t.Errorf("--repo on a release binary: err = %v, want a refusal naming --repo and the path", err)
	}
	t.Setenv("SH_AGENT_REPO", checkout)
	if _, err := ResolveSource("", home, "v1.2.3"); err == nil ||
		!strings.Contains(err.Error(), "release binary") || !strings.Contains(err.Error(), "$SH_AGENT_REPO="+checkout) {
		t.Errorf("$SH_AGENT_REPO on a release binary: err = %v, want a refusal naming the variable", err)
	}
	t.Setenv("SH_AGENT_REPO", "")
	os.Unsetenv("SH_AGENT_REPO")
	t.Chdir(cwd)
	if err := SaveState(home, State{Repo: repo}); err != nil {
		t.Fatal(err)
	}

	// A manifest written by an embedded install records no path: its repo
	// value is never read as one, even when a directory of that name holding
	// a harness sits in the working directory.
	if err := os.CopyFS(filepath.Join(cwd, "embedded:v1.0.0"), os.DirFS(repo)); err != nil {
		t.Fatal(err)
	}
	if err := SaveState(home, State{Repo: "embedded:v1.0.0"}); err != nil {
		t.Fatal(err)
	}
	if src, err := ResolveSource("", home, "v1.2.3"); err != nil || src.Dir != "" {
		t.Errorf("recorded embedded install: %+v, %v, want the embedded source", src, err)
	}
	if _, err := ResolveSource("", home, ""); err == nil {
		t.Error("a development build must not read a recorded embedded:<version> as a path")
	}
}

func TestBehindAndPull(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	base := t.TempDir()
	run := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	origin := filepath.Join(base, "origin.git")
	run(base, "init", "--bare", "-q", "-b", "main", origin)

	a := filepath.Join(base, "a")
	run(base, "clone", "-q", origin, a)
	if err := os.WriteFile(filepath.Join(a, "f"), []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(a, "add", "f")
	run(a, "commit", "-q", "-m", "one")
	run(a, "push", "-q", "-u", "origin", "main")

	b := filepath.Join(base, "b")
	run(base, "clone", "-q", origin, b)

	if n, err := Behind(b); err != nil || n != 0 {
		t.Fatalf("fresh clone behind = %d, %v", n, err)
	}

	if err := os.WriteFile(filepath.Join(a, "f"), []byte("2"), 0o644); err != nil {
		t.Fatal(err)
	}
	run(a, "commit", "-q", "-am", "two")
	run(a, "push", "-q")

	if n, err := Behind(b); err != nil || n != 1 {
		t.Fatalf("behind = %d, %v; want 1", n, err)
	}

	if err := os.WriteFile(filepath.Join(b, "f"), []byte("dirty"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Pull(b, &bytes.Buffer{}); err == nil {
		t.Fatal("pull over a dirty tree must be refused")
	}
	run(b, "checkout", "--", "f")
	if err := Pull(b, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if n, _ := Behind(b); n != 0 {
		t.Errorf("after pull behind = %d", n)
	}
}

// TestNestedCopyIsNotACheckout: a harness-shaped tree inside another git
// repository must not inherit that repository's branch, commit, or dirty
// state, and git-driven update must refuse it.
func TestNestedCopyIsNotACheckout(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := newRepo(t)
	gitIn(t, root, "init", "-q")
	gitIn(t, root, "add", ".")
	gitIn(t, root, "commit", "-q", "-m", "fixture")

	nested := filepath.Join(root, "mirror")
	for _, rel := range []string{"context", "skills"} {
		if err := os.MkdirAll(filepath.Join(nested, rel), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(nested, "context", "AGENTS.md"), []byte("# copy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !IsRepo(nested) {
		t.Fatal("fixture: nested copy should look like a harness")
	}

	if got := Inspect(nested); got != (RepoInfo{}) {
		t.Errorf("nested copy Inspect = %+v, want empty", got)
	}
	if _, err := Behind(nested); !errors.Is(err, errPlainCopy) {
		t.Errorf("nested copy Behind err = %v, want %v", err, errPlainCopy)
	}
	if err := Pull(nested, &bytes.Buffer{}); !errors.Is(err, errPlainCopy) {
		t.Errorf("nested copy Pull err = %v, want %v", err, errPlainCopy)
	}

	// An inherited GIT_DIR/GIT_WORK_TREE must not redirect the check to the
	// enclosing repository: the check looks for .git itself and never asks git.
	t.Setenv("GIT_DIR", filepath.Join(root, ".git"))
	t.Setenv("GIT_WORK_TREE", nested)
	if got := Inspect(nested); got != (RepoInfo{}) {
		t.Errorf("nested copy with GIT_DIR/GIT_WORK_TREE set: Inspect = %+v, want empty", got)
	}
	if _, err := Behind(nested); !errors.Is(err, errPlainCopy) {
		t.Errorf("nested copy with GIT_DIR set: Behind err = %v, want %v", err, errPlainCopy)
	}
	t.Setenv("GIT_DIR", "")
	t.Setenv("GIT_WORK_TREE", "")
	os.Unsetenv("GIT_DIR")
	os.Unsetenv("GIT_WORK_TREE")

	got := Inspect(root)
	if got.Branch == "" || got.Commit == "" || !got.Dirty {
		t.Errorf("repo root Inspect = %+v, want branch, commit, and dirty (the untracked mirror)", got)
	}
	if _, err := Behind(root); errors.Is(err, errPlainCopy) {
		t.Errorf("repo root Behind must not be refused as a plain copy: %v", err)
	}
}

func TestLastRetro(t *testing.T) {
	home := t.TempDir()
	if _, ok := LastRetro(home); ok {
		t.Fatal("no proposals yet")
	}
	dir := filepath.Join(home, ".config", "sh-agent", "proposals")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "2026-09-03-retro.md"), []byte("# retro\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	first, ok := LastRetro(home)
	if !ok {
		t.Fatal("proposal should be found")
	}
	// A project-level retro appends to the index without a proposal file;
	// that counts as a retro too, and its newer mtime wins.
	index := filepath.Join(home, ".config", "sh-agent", "retros.jsonl")
	if err := os.WriteFile(index, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	later := first.Add(48 * time.Hour)
	if err := os.Chtimes(index, later, later); err != nil {
		t.Fatal(err)
	}
	if got, ok := LastRetro(home); !ok || !got.Equal(later) {
		t.Fatalf("LastRetro = %v, %v; want %v", got, ok, later)
	}
}

func TestSlashNoteOnlyWhenSomeSkillIsSlashOnly(t *testing.T) {
	repo := newRepo(t)
	inv, err := LoadInventory(DirSource(repo))
	if err != nil {
		t.Fatal(err)
	}
	const marked = "Skills marked `/name` are typed by the user"
	if !strings.Contains(RenderInventoryCompact(inv), marked) || !strings.Contains(RenderInventory(inv), marked) {
		t.Errorf("fixture has a slash-only skill; the note should be printed")
	}
	for i := range inv.Skills {
		inv.Skills[i].SlashOnly = false
	}
	for _, out := range []string{RenderInventoryCompact(inv), RenderInventory(inv)} {
		if strings.Contains(out, marked) {
			t.Errorf("no slash-only skill; the note must not be printed:\n%s", out)
		}
		if !strings.Contains(out, "can also be typed as `/name`") {
			t.Errorf("expected the all-model-invocable sentence:\n%s", out)
		}
	}
}

// gitIn runs git in dir and fails the test on error.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, outb)
	}
	return strings.TrimSpace(string(outb))
}

// An explicit --repo or $SH_AGENT_REPO is used or refused, never replaced by
// the embedded files, even in a release binary (D15).
func TestResolveSourceRefusesNamedNonRepo(t *testing.T) {
	home := t.TempDir()
	bogus := filepath.Join(t.TempDir(), "not-a-checkout")
	for _, v := range []string{"", "v1.2.3"} {
		t.Setenv("SH_AGENT_REPO", "")
		if _, err := ResolveSource(bogus, home, v); err == nil {
			t.Errorf("version %q: an explicit path that is not a checkout must be an error, not a fallback", v)
		}
		t.Setenv("SH_AGENT_REPO", bogus)
		if _, err := ResolveSource("", home, v); err == nil {
			t.Errorf("version %q: $SH_AGENT_REPO pointing at a non-checkout must be an error, not a fallback", v)
		}
	}
	t.Setenv("SH_AGENT_REPO", "")
	repo := newRepo(t)
	if got, err := ResolveSource(repo, home, ""); err != nil || got.Dir != repo {
		t.Errorf("explicit checkout: %+v, %v", got, err)
	}
}

func TestInstallMarksDirtyTree(t *testing.T) {
	repo := newRepo(t)
	gitIn(t, repo, "init", "-q")
	gitIn(t, repo, "add", ".")
	gitIn(t, repo, "commit", "-q", "-m", "fixture")
	short := gitIn(t, repo, "rev-parse", "--short", "HEAD")
	home := t.TempDir()
	files, err := Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Install(DirSource(repo), home, files, VendorAll, false, &out); err != nil {
		t.Fatal(err)
	}
	st, _, err := LoadState(home)
	if err != nil {
		t.Fatal(err)
	}
	if st.Commit != short {
		t.Errorf("clean tree: commit = %q, want %q", st.Commit, short)
	}
	if err := os.WriteFile(filepath.Join(repo, "context", "AGENTS.md"), []byte("# AGENTS.md\n\nBe honest and brief.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	files, err = Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if err := Install(DirSource(repo), home, files, VendorAll, false, &out); err != nil {
		t.Fatal(err)
	}
	st, _, err = LoadState(home)
	if err != nil {
		t.Fatal(err)
	}
	if st.Commit != short+"+dirty" {
		t.Errorf("dirty tree: commit = %q, want %q", st.Commit, short+"+dirty")
	}
}

// TestReinstallBuildsAndInstalls runs the real rebuild path against this
// checkout: go build into the repo's scratch space, then that build installs
// into a scratch home and copies itself there, stamped with the revision.
func TestReinstallBuildsAndInstalls(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not on PATH")
	}
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	var out bytes.Buffer
	if err := Reinstall(repo, []string{"--home", home}, &out); err != nil {
		t.Fatalf("Reinstall: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "files managed") {
		t.Errorf("child install output missing:\n%s", out.String())
	}
	bin := filepath.Join(home, BinPath)
	ver, err := exec.Command(bin, "version").Output()
	if err != nil {
		t.Fatalf("installed binary does not run: %v", err)
	}
	if v := strings.TrimSpace(string(ver)); strings.Contains(v, "(devel)") || strings.Contains(v, "unknown") {
		t.Errorf("installed binary is unstamped: %q", v)
	}
	builds, _ := filepath.Glob(filepath.Join(repo, ".agents", "scratch", "sh-agent-build-*"))
	if len(builds) != 0 {
		t.Errorf("build directories left behind: %v", builds)
	}
}

// vendoredSkills reads skills/VENDORED.md and returns the skill directories
// listed there. Vendored text is never edited in place (AGENTS.md), so a
// prose check that walks it would report a defect nobody may fix; the fix
// for those is a re-fetch or a fork under another name.
func vendoredSkills(t *testing.T, repo string) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repo, "skills", "VENDORED.md"))
	if err != nil {
		t.Fatalf("reading VENDORED.md: %v", err)
	}
	out := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "|") {
			continue
		}
		cell := strings.TrimSpace(strings.Trim(strings.SplitN(strings.TrimSpace(line), "|", 3)[1], " `"))
		if cell == "" || strings.Contains(cell, " ") || strings.HasPrefix(cell, "-") {
			continue
		}
		if info, err := os.Stat(filepath.Join(repo, "skills", cell)); err == nil && info.IsDir() {
			out[cell] = true
		}
	}
	if len(out) == 0 {
		t.Fatal("VENDORED.md named no skill directory; the parser or the file changed shape")
	}
	return out
}

// TestNoTmpInHarnessProse guards the scratch rule at the source: the
// context, the agents, the skills — their SKILL.md AND their references,
// where the rule drifted unseen (BUG-019) — the commands, the rules, and the
// repo docs may mention /tmp only to forbid it. Four records (BUG-003,
// BUG-007, BUG-014, BUG-019) are the reason this is a test and not a
// reminder. Vendored skills are exempt: their text may not be edited here.
func TestNoTmpInHarnessProse(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	roots := []string{"context", "agents", "commands", "rules", "README.md", "AGENTS.md"}
	vendored := vendoredSkills(t, repo)
	dirs, _ := filepath.Glob(filepath.Join(repo, "skills", "*"))
	for _, d := range dirs {
		info, err := os.Stat(d)
		if err != nil || !info.IsDir() || vendored[filepath.Base(d)] {
			continue
		}
		rel, _ := filepath.Rel(repo, d)
		roots = append(roots, rel)
	}
	allowed := func(line string) bool {
		l := strings.ToLower(line)
		return strings.Contains(l, "never /tmp") || strings.Contains(l, "never `/tmp`") ||
			strings.Contains(l, "not /tmp") || strings.Contains(l, "not `/tmp`") || strings.Contains(l, "instead of `/tmp`")
	}
	for _, root := range roots {
		err := filepath.WalkDir(filepath.Join(repo, root), func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for i, line := range strings.Split(string(data), "\n") {
				if strings.Contains(line, "/tmp") && !allowed(line) {
					rel, _ := filepath.Rel(repo, path)
					t.Errorf("%s:%d mentions /tmp outside a 'never /tmp' clause: %s", rel, i+1, strings.TrimSpace(line))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// A source that resolves to the repo, to skills/, or to the destination must
// never become the deletion target: ".", ".." and a relative "skills/<name>"
// reached Base unchanged and named the deletion root itself (BUG-037).
func TestAdoptRefusesDestructiveSources(t *testing.T) {
	home := t.TempDir()
	skillDir := func(t *testing.T) string {
		t.Helper()
		d := filepath.Join(t.TempDir(), "delta")
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(d, "SKILL.md"), []byte("---\nname: delta\ndescription: d\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return d
	}

	t.Run("dot adopts the directory itself, never skills/", func(t *testing.T) {
		repo := newRepo(t)
		t.Chdir(skillDir(t))
		dest, err := Adopt(repo, home, ".", true)
		if err != nil {
			t.Fatalf("adopt . --force: %v", err)
		}
		if dest != filepath.Join(repo, "skills", "delta") {
			t.Errorf("dest = %s", dest)
		}
		if _, err := os.Stat(filepath.Join(repo, "skills", "alpha", "SKILL.md")); err != nil {
			t.Errorf("skills/ was removed: %v", err)
		}
	})

	t.Run("dotdot adopts the parent, never the repo", func(t *testing.T) {
		repo := newRepo(t)
		sub := filepath.Join(skillDir(t), "sub")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Chdir(sub)
		dest, err := Adopt(repo, home, "..", true)
		if err != nil {
			t.Fatalf("adopt .. --force: %v", err)
		}
		if dest != filepath.Join(repo, "skills", "delta") {
			t.Errorf("dest = %s", dest)
		}
		if _, err := os.Stat(filepath.Join(repo, "context", "AGENTS.md")); err != nil {
			t.Errorf("the repo was removed: %v", err)
		}
	})

	t.Run("a skill adopted onto itself", func(t *testing.T) {
		repo := newRepo(t)
		src := filepath.Join(repo, "skills", "alpha")
		if _, err := Adopt(repo, home, src, true); err == nil {
			t.Error("adopting a skill onto itself should fail")
		}
		if _, err := os.Stat(filepath.Join(src, "SKILL.md")); err != nil {
			t.Errorf("the source was removed: %v", err)
		}
	})

	t.Run("a relative path into the repo's own skills/", func(t *testing.T) {
		repo := newRepo(t)
		t.Chdir(repo)
		if _, err := Adopt(repo, home, filepath.Join("skills", "alpha"), true); err == nil {
			t.Error("adopting skills/alpha from the repo root should fail")
		}
		if _, err := os.Stat(filepath.Join(repo, "skills", "alpha", "SKILL.md")); err != nil {
			t.Errorf("the source was removed: %v", err)
		}
	})
}

func TestWithin(t *testing.T) {
	for _, c := range []struct {
		parent, child string
		want          bool
	}{
		{"/a/skills", "/a/skills", true},
		{"/a/skills", "/a/skills/", true},
		{"/a/skills", "/a/skills/delta", true},
		{"/a/skills/delta", "/a/skills", false},
		{"/a/skills", "/a/skills-old", false},
		{"/a/skills", "/a/skillsdelta", false},
		{"/a/skills", "/b/skills/delta", false},
		{"/a", "/a/skills/delta", true},
	} {
		if got := within(c.parent, c.child); got != c.want {
			t.Errorf("within(%q, %q) = %v, want %v", c.parent, c.child, got, c.want)
		}
	}
}

// TestWriteAtomicIsAtomicAndDurable covers the write every deployed file
// goes through: the bytes land whole under the asked-for mode, and no temp
// file is left beside them on either outcome. The Sync added on 2026-09-16
// is not directly observable from a test — a crash is what would show it —
// so what is asserted here is the contract it protects: after the call
// returns, the destination holds exactly the bytes and nothing else does.
func TestWriteAtomicIsAtomicAndDurable(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "nested", "context.md")
	data := []byte("# AGENTS.md\n\nwhole file\n")
	if err := writeAtomic(dest, data, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Errorf("content = %q, want %q", got, data)
	}
	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Errorf("mode = %v, want 0644", info.Mode().Perm())
	}
	entries, err := os.ReadDir(filepath.Dir(dest))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".sh-agent-") {
			t.Errorf("temp file %s left beside the destination", e.Name())
		}
	}
	// Overwriting an existing file keeps the same guarantee.
	if err := writeAtomic(dest, []byte("second\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(dest); string(got) != "second\n" {
		t.Errorf("after overwrite = %q", got)
	}
}

// TestRenderContextRefusesSymlinks is BUG-036's class, tested directly on
// RenderContext rather than only through Plan: a symlinked context directory
// or context file would let a tree outside the repo decide what every
// machine reads.
func TestRenderContextRefusesSymlinks(t *testing.T) {
	t.Run("context directory is a link", func(t *testing.T) {
		repo := newRepo(t)
		outside := t.TempDir()
		if err := os.WriteFile(filepath.Join(outside, "AGENTS.md"), []byte("# hostile\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(filepath.Join(repo, "context")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(repo, "context")); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		inv, _ := LoadInventory(DirSource(repo))
		if _, err := RenderContext(DirSource(repo), inv); err == nil {
			t.Fatal("a symlinked context directory must be refused")
		} else if !strings.Contains(err.Error(), "symlink") {
			t.Errorf("the error should name the symlink, got %v", err)
		}
	})

	t.Run("a file inside context is a link", func(t *testing.T) {
		repo := newRepo(t)
		outside := filepath.Join(t.TempDir(), "extra.md")
		if err := os.WriteFile(outside, []byte("# hostile\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(repo, "context", "extra.md")); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		inv, _ := LoadInventory(DirSource(repo))
		if _, err := RenderContext(DirSource(repo), inv); err == nil {
			t.Fatal("a symlinked context file must be refused")
		} else if !strings.Contains(err.Error(), "symlink") {
			t.Errorf("the error should name the symlink, got %v", err)
		}
	})
}

// TestAdoptFailedCopyKeepsTheExistingSkill: adopt --force used to delete the
// destination and then copy, so a copy that failed part-way left the operator
// with neither the old skill nor the new one. The copy is staged and renamed
// into place, so a refusal mid-walk changes nothing.
func TestAdoptFailedCopyKeepsTheExistingSkill(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()

	dest := filepath.Join(repo, "skills", "widget")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	original := "---\nname: widget\ndescription: the version already here\n---\nold body\n"
	if err := os.WriteFile(filepath.Join(dest, "SKILL.md"), []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	// A source whose copy fails part-way: copyTree refuses symlinks.
	src := filepath.Join(t.TempDir(), "widget")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("---\nname: widget\ndescription: replacement\n---\nnew body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/etc/passwd", filepath.Join(src, "linked.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if _, err := Adopt(repo, home, src, true); err == nil {
		t.Fatal("a source holding a symlink must be refused")
	}
	got, err := os.ReadFile(filepath.Join(dest, "SKILL.md"))
	if err != nil {
		t.Fatalf("the existing skill must survive a failed adopt: %v", err)
	}
	if string(got) != original {
		t.Errorf("the existing skill was modified: %q", got)
	}
	assertSkillsClean(t, repo, "after a failed adopt")

	// A rename that cannot complete must leave the ORIGINAL in place, never
	// delete both copies: dest is moved aside and put back (release review of
	// v1.0.0, finding 3). Forced here by making the destination's parent
	// read-only so the rename into it fails.
	if os.Geteuid() != 0 {
		if err := os.Remove(filepath.Join(src, "linked.md")); err != nil {
			t.Fatal(err)
		}
		skillsDir := filepath.Join(repo, "skills")
		if err := os.Chmod(skillsDir, 0o500); err != nil {
			t.Fatal(err)
		}
		_, err := Adopt(repo, home, src, true)
		if chmodErr := os.Chmod(skillsDir, 0o755); chmodErr != nil {
			t.Fatal(chmodErr)
		}
		if err == nil {
			t.Error("a rename into an unwritable skills/ should fail")
		}
		got, readErr := os.ReadFile(filepath.Join(dest, "SKILL.md"))
		if readErr != nil {
			t.Fatalf("the original skill must survive a failed rename: %v", readErr)
		}
		if string(got) != original {
			t.Errorf("the original skill changed: %q", got)
		}
		if err != nil && !strings.Contains(err.Error(), ".agents/scratch") {
			t.Errorf("the error should name where the prepared copy is, got %v", err)
		}
	}

	// The happy path still replaces the skill and leaves no staging behind.
	os.Remove(filepath.Join(src, "linked.md"))
	if _, err := Adopt(repo, home, src, true); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); !strings.Contains(string(got), "new body") {
		t.Errorf("the skill should have been replaced, got %q", got)
	}
	assertSkillsClean(t, repo, "after a successful adopt")

	// The staging directory lives in .agents/scratch and is removed either
	// way; if one survived there with a complete SKILL.md, LoadSkills would
	// not see it (it walks skills/ only), which is the point of the move.
	leftovers, _ := filepath.Glob(filepath.Join(repo, ".agents", "scratch", "adopt-*"))
	if len(leftovers) != 0 {
		t.Errorf("staging directories left in scratch: %v", leftovers)
	}
}

// assertSkillsClean fails if skills/ holds anything but real skill
// directories: a staging or partial copy there would be loaded as a second
// skill of the same name and installed over the real one (code review
// 2026-09-16).
func assertSkillsClean(t *testing.T, repo, when string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(repo, "skills"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if strings.Contains(e.Name(), "staging") || strings.Contains(e.Name(), "adopt-") || strings.HasPrefix(e.Name(), ".") {
			t.Errorf("%s: skills/ holds %q, which is not a skill", when, e.Name())
		}
	}
	skills, err := LoadSkills(DirSource(repo))
	if err != nil {
		t.Fatalf("%s: LoadSkills: %v", when, err)
	}
	seen := map[string]bool{}
	for _, s := range skills {
		if seen[s.Name] {
			t.Errorf("%s: two skills named %q are loaded", when, s.Name)
		}
		seen[s.Name] = true
	}
}

// TestLoadSkillsRefusesDuplicateNames: two directories whose frontmatter
// claims one name install to a single destination, so the walk order alone
// would decide which text every machine reads. That is what an interrupted
// `adopt` used to be able to leave behind (its staging now lives outside
// skills/), and what a hand-copied directory can still produce.
func TestLoadSkillsRefusesDuplicateNames(t *testing.T) {
	repo := newRepo(t)
	dup := filepath.Join(repo, "skills", "alpha-copy")
	if err := os.MkdirAll(dup, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: alpha\ndescription: a second directory claiming the same name\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(dup, "SKILL.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadSkills(DirSource(repo))
	if err == nil {
		t.Fatal("two skills named alpha must be refused")
	}
	for _, want := range []string{"alpha", "alpha-copy"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error should name %s, got %v", want, err)
		}
	}
}

// TestPlanRefusesSecretShapedFiles: a credential that slipped past
// .gitignore into a skill directory is invisible to review (git never shows
// it) and install is the one step that would copy it into both vendor
// directories on every machine (security review 2026-09-16).
func TestPlanRefusesSecretShapedFiles(t *testing.T) {
	for _, name := range []string{".env", ".env.local", "server.key", "cert.pem", "creds.tfvars", "application_default_credentials.json", "sa-key.json"} {
		t.Run(name, func(t *testing.T) {
			repo := newRepo(t)
			if err := os.WriteFile(filepath.Join(repo, "skills", "alpha", name), []byte("secret\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := Plan(DirSource(repo), nil, VendorAll); err == nil {
				t.Fatalf("%s must not be deployed", name)
			} else if !strings.Contains(err.Error(), name) {
				t.Errorf("the error should name the file, got %v", err)
			}
		})
	}
	t.Run("an ordinary reference file still plans", func(t *testing.T) {
		repo := newRepo(t)
		if err := os.WriteFile(filepath.Join(repo, "skills", "alpha", "keys-explained.md"), []byte("about keys\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		files, err := Plan(DirSource(repo), nil, VendorAll)
		if err != nil {
			t.Fatalf("a normal file must still plan: %v", err)
		}
		var found bool
		for _, f := range files {
			if strings.HasSuffix(f.Rel, "keys-explained.md") {
				found = true
			}
		}
		if !found {
			t.Error("the reference file should have been planned")
		}
	})
}

// TestInspectCountsUntrackedFiles: `status.showUntrackedFiles=no` used to
// make a tree holding a whole new untracked skill read as clean, which
// `update` would then deploy everywhere (security review 2026-09-16).
func TestInspectCountsUntrackedFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	repo := newRepo(t)
	env := append(os.Environ(), "GIT_CONFIG_GLOBAL="+filepath.Join(t.TempDir(), "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		cmd.Env = env
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "--quiet", "-b", "main")
	git("config", "user.email", "fixture@example.invalid")
	git("config", "user.name", "fixture")
	git("add", "-A")
	git("commit", "--quiet", "-m", "fixture")
	if Inspect(repo).Dirty {
		t.Fatal("a freshly committed tree is clean")
	}

	// The operator's own setting, which Inspect must not inherit.
	git("config", "status.showUntrackedFiles", "no")
	if err := os.MkdirAll(filepath.Join(repo, "skills", "smuggled"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "skills", "smuggled", "SKILL.md"), []byte("---\nname: smuggled\ndescription: never committed\n---\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !Inspect(repo).Dirty {
		t.Error("an untracked skill directory makes the checkout dirty, whatever status.showUntrackedFiles says")
	}
}

// TestInspectFailsClosedOnUnknownStatus: "git could not answer" must never
// read as "the tree is clean" to a caller deciding whether to deploy
// (release review of v1.0.0, finding 4).
func TestInspectFailsClosedOnUnknownStatus(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	repo := newRepo(t)
	// A .git that is a file, not a directory: IsCheckout is satisfied, and
	// every git command then fails rather than reporting a clean tree.
	if err := os.WriteFile(filepath.Join(repo, ".git"), []byte("not a git dir\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	info := Inspect(repo)
	if !info.StatusUnknown {
		t.Error("a git status that failed must be reported as unknown")
	}
	if err := Pull(repo, io.Discard); err == nil {
		t.Error("Pull must refuse a checkout whose status git cannot report")
	}
}

// TestReinstallNeverNeedsToDownloadAToolchain exercises Reinstall in the
// environment that broke it: a module whose go.mod asks for a toolchain that
// is not present, with no proxy to fetch it from. Under the default
// GOTOOLCHAIN=auto that build exits 1 with "toolchain not available"; the
// same build must succeed through Reinstall, which pins GOTOOLCHAIN=local.
// Asserted by running it, not by grepping the source (round-2 release review).
func TestReinstallNeverNeedsToDownloadAToolchain(t *testing.T) {
	for _, bin := range []string{"go", "git"} {
		if _, err := exec.LookPath(bin); err != nil {
			t.Skipf("%s is not on PATH", bin)
		}
	}
	repo := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		p := filepath.Join(repo, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// A toolchain that will never exist locally, so auto must try to fetch it.
	write("go.mod", "module fixture\n\ngo 1.24.4\n\ntoolchain go1.99.0\n")
	write("cmd/sh-agent/main.go", "package main\n\nfunc main() {}\n")
	// A real repository: IsCheckout needs .git, and `go build` stamps VCS
	// information, which fails on a .git that is not one.
	for _, args := range [][]string{
		{"init", "--quiet", "-b", "main"},
		{"config", "user.email", "fixture@example.invalid"},
		{"config", "user.name", "fixture"},
		{"add", "-A"},
		{"commit", "--quiet", "-m", "fixture"},
	} {
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+filepath.Join(t.TempDir(), "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	cache := t.TempDir()
	t.Setenv("GOMODCACHE", cache)
	t.Setenv("GOPROXY", "off")
	t.Setenv("GOFLAGS", "")

	// Must-fail: the same build under the default auto.
	t.Setenv("GOTOOLCHAIN", "auto")
	probe := exec.Command("go", "build", "-o", filepath.Join(t.TempDir(), "probe"), "./cmd/sh-agent")
	probe.Dir = repo
	if out, err := probe.CombinedOutput(); err == nil {
		t.Skipf("this environment resolved go1.99.0, so the failure cannot be staged: %s", out)
	}

	// Must-pass: through Reinstall, which forces the local toolchain. The
	// child install is the fixture's own no-op main, so success means the
	// BUILD succeeded, which is what this test is about.
	var out bytes.Buffer
	if err := Reinstall(repo, []string{"--home", t.TempDir()}, &out); err != nil {
		t.Fatalf("Reinstall must not need to download a toolchain: %v\n%s", err, out.String())
	}
}

// TestPlanDropsGroupAndOtherWrite: a skill file that is group-writable in
// the repo must not land group-writable in ~/.claude, where its contents are
// instructions every session reads (security review of v1.0.0).
func TestPlanDropsGroupAndOtherWrite(t *testing.T) {
	repo := newRepo(t)
	loose := filepath.Join(repo, "skills", "alpha", "loose.md")
	if err := os.WriteFile(loose, []byte("reference\n"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(loose, 0o666); err != nil {
		t.Fatal(err)
	}
	files, err := Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	var seen int
	for _, f := range files {
		if f.Mode&0o022 != 0 {
			t.Errorf("%s would be deployed mode %o — group or other writable", f.Rel, f.Mode)
		}
		if strings.HasSuffix(f.Rel, "loose.md") {
			seen++
			if f.Mode.Perm() != 0o644 {
				t.Errorf("loose.md mode = %o, want 644", f.Mode.Perm())
			}
		}
	}
	if seen == 0 {
		t.Error("the loose file should still be planned")
	}
}

// TestInstallCorrectsTheModeOfAnUnchangedFile: same bytes is not the same
// file. A deployed copy left group-writable stays that way on the
// "unchanged" branch, so the machine that has the problem is the one that
// never gets the fix (security review of v1.0.0, round 2, which measured it
// on a real install).
func TestInstallCorrectsTheModeOfAnUnchangedFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root's umask games make the mode assertions meaningless")
	}
	repo := newRepo(t)
	home := t.TempDir()
	files, err := Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := Install(DirSource(repo), home, files, VendorAll, false, &buf); err != nil {
		t.Fatal(err)
	}
	victim := filepath.Join(home, ClaudeRoot, "skills", "alpha", "SKILL.md")
	if err := os.Chmod(victim, 0o666); err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	if err := Install(DirSource(repo), home, files, VendorAll, false, &buf); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(victim)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o022 != 0 {
		t.Errorf("a second install left %s mode %o — group or other writable", victim, info.Mode().Perm())
	}
	// And a dry run changes nothing.
	if err := os.Chmod(victim, 0o666); err != nil {
		t.Fatal(err)
	}
	buf.Reset()
	if err := Install(DirSource(repo), home, files, VendorAll, true, &buf); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(victim); info.Mode().Perm() != 0o666 {
		t.Errorf("--dry-run changed a mode: %o", info.Mode().Perm())
	}
}

// TestAdoptRefusesSecretShapedFiles guards the entrance: Plan refuses a
// credential on the way out, but adopt copies into skills/, where it would
// be committed and pushed before any consumer's install ever refused it.
func TestAdoptRefusesSecretShapedFiles(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()
	src := filepath.Join(t.TempDir(), "widget")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "SKILL.md"), []byte("---\nname: widget\ndescription: has a key beside it\n---\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "id_ed25519"), []byte("PRIVATE\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := Adopt(repo, home, src, false)
	if err == nil {
		t.Fatal("adopt must refuse a skill carrying a private key")
	}
	if !strings.Contains(err.Error(), "id_ed25519") {
		t.Errorf("the error should name the file, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(repo, "skills", "widget")); !errors.Is(statErr, os.ErrNotExist) {
		t.Error("nothing should have been copied into skills/")
	}
}

// install names each file it changes, on every run: `add` for a file that
// did not exist, `update` for changed bytes (or a previous file it could not
// read), `mode` for a permission it corrects, `remove` for a file the source
// no longer has. Unchanged files print nothing; --dry-run prints the same
// lines and writes nothing (plan sh-agent-works-everywhere, Phase 1, D4, D6-D9).
func TestInstallNamesEveryChange(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()
	lines := func(out string, verb string) []string {
		var got []string
		for _, l := range strings.Split(out, "\n") {
			if strings.HasPrefix(l, verb+" ") {
				got = append(got, l)
			}
		}
		return got
	}
	plan := func() []File {
		t.Helper()
		files, err := Plan(DirSource(repo), []byte("#!/bin/sh\necho sh-agent\n"), VendorAll)
		if err != nil {
			t.Fatal(err)
		}
		return files
	}
	install := func(files []File, dryRun bool) string {
		t.Helper()
		var out bytes.Buffer
		if err := Install(DirSource(repo), home, files, VendorAll, dryRun, &out); err != nil {
			t.Fatalf("install: %v\n%s", err, out.String())
		}
		return out.String()
	}

	// 1. First install: one `add` per managed file, nothing else.
	files := plan()
	out := install(files, false)
	if adds := lines(out, "add"); len(adds) != len(files) {
		t.Errorf("first install: %d `add` lines for %d files:\n%s", len(adds), len(files), out)
	}
	for _, verb := range []string{"update", "remove", "mode"} {
		if got := lines(out, verb); len(got) != 0 {
			t.Errorf("first install printed %s lines: %v", verb, got)
		}
	}
	if !strings.Contains(out, strconv.Itoa(len(files))+" written, 0 unchanged") {
		t.Errorf("first install summary:\n%s", out)
	}

	// A second, identical install: nothing named at all.
	out = install(files, false)
	for _, verb := range []string{"add", "update", "remove", "mode"} {
		if got := lines(out, verb); len(got) != 0 {
			t.Errorf("identical install printed %s lines: %v", verb, got)
		}
	}

	// 2. One changed source file: `update` for each deployed copy of it, nothing else.
	alpha := filepath.Join(repo, "skills", "alpha", "SKILL.md")
	if err := os.WriteFile(alpha, []byte("---\nname: alpha\ndescription: changed\n---\nnew body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	files = plan()
	out = install(files, false)
	updates := lines(out, "update")
	if len(updates) == 0 {
		t.Errorf("a changed source file printed no `update` line:\n%s", out)
	}
	for _, l := range updates {
		if !strings.HasSuffix(l, "skills/alpha/SKILL.md") {
			t.Errorf("unexpected update line: %q", l)
		}
	}
	for _, verb := range []string{"add", "remove", "mode"} {
		if got := lines(out, verb); len(got) != 0 {
			t.Errorf("changed-file install printed %s lines: %v", verb, got)
		}
	}

	// 3. A removed source file: `remove … (no longer in the source)` on the real run.
	if err := os.RemoveAll(filepath.Join(repo, "skills", "beta")); err != nil {
		t.Fatal(err)
	}
	files = plan()
	out = install(files, false)
	removes := lines(out, "remove")
	if len(removes) == 0 {
		t.Errorf("a removed source file printed no `remove` line:\n%s", out)
	}
	for _, l := range removes {
		if !strings.Contains(l, "skills/beta/") || !strings.HasSuffix(l, "(no longer in the source)") {
			t.Errorf("unexpected remove line: %q", l)
		}
	}

	// 4. A mode-only difference: one `mode` line, still counted unchanged.
	if os.Geteuid() == 0 {
		t.Skip("root reads a 0000 file and its umask blurs the mode cases")
	}
	victim := filepath.Join(home, ClaudeRoot, "skills", "alpha", "SKILL.md")
	if err := os.Chmod(victim, 0o666); err != nil {
		t.Fatal(err)
	}
	out = install(files, false)
	modes := lines(out, "mode")
	if len(modes) != 1 || !strings.Contains(modes[0], "skills/alpha/SKILL.md (0666 -> 0644)") {
		t.Errorf("mode-only difference: want one `mode … (0666 -> 0644)` line, got %v:\n%s", modes, out)
	}
	if !strings.Contains(out, "0 written, ") {
		t.Errorf("a mode-only correction counted as written:\n%s", out)
	}
	if info, _ := os.Stat(victim); info.Mode().Perm() != 0o644 {
		t.Errorf("mode not corrected: %o", info.Mode().Perm())
	}

	// 5. A previous file that cannot be read is written over and the line says why (D6, D9).
	if err := os.Chmod(victim, 0o000); err != nil {
		t.Fatal(err)
	}
	out = install(files, false)
	updates = lines(out, "update")
	if len(updates) != 1 || !strings.Contains(updates[0], "skills/alpha/SKILL.md (previous file unreadable: ") {
		t.Errorf("unreadable previous file: want one `update … (previous file unreadable: …)` line, got %v:\n%s", updates, out)
	}
	if got, err := os.ReadFile(victim); err != nil || !strings.Contains(string(got), "new body") {
		t.Errorf("the unreadable file was not written over: %v %q", err, got)
	}

	// 6. and 7. --dry-run: the same lines for a changed source, a removed
	// source and a mode difference, and nothing written or chmodded.
	if err := os.WriteFile(alpha, []byte("---\nname: alpha\ndescription: changed again\n---\nthird body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(victim, 0o666); err != nil {
		t.Fatal(err)
	}
	gamma := filepath.Join(home, ClaudeRoot, "skills", "alpha", "references", "notes.md")
	if err := os.Chmod(gamma, 0o666); err != nil {
		t.Fatal(err)
	}
	files = plan()
	out = install(files, true)
	if got := lines(out, "update"); len(got) == 0 || !strings.HasSuffix(got[0], "skills/alpha/SKILL.md") {
		t.Errorf("--dry-run: want `update` lines for the changed skill, got %v:\n%s", got, out)
	}
	if got := lines(out, "mode"); len(got) != 1 || !strings.Contains(got[0], "references/notes.md (0666 -> 0644)") {
		t.Errorf("--dry-run: want one `mode` line for the untouched-but-chmodded file, got %v:\n%s", got, out)
	}
	if got, _ := os.ReadFile(victim); strings.Contains(string(got), "third body") {
		t.Errorf("--dry-run wrote a file")
	}
	if info, _ := os.Stat(gamma); info.Mode().Perm() != 0o666 {
		t.Errorf("--dry-run changed a mode: %o", info.Mode().Perm())
	}
}
