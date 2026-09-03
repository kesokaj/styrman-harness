package trail

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// newRepo builds a throwaway git repo that has opted into the trail contract.
func newRepo(t *testing.T) Project {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "test"},
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
	proj, ok := Resolve(root)
	if !ok {
		t.Fatal("Resolve failed on a fresh repo")
	}
	return proj
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestIsSourcePath(t *testing.T) {
	proj := Project{Root: "/r", Docs: "/r/.agents/docs"}
	cases := map[string]bool{
		"main.go":                   true,
		"web/app.tsx":               true,
		"Dockerfile":                true,
		"ops/k8s/deploy.yaml":       true,
		"go.mod":                    true,
		"scripts/build.sh":          true,
		"docker-compose.yml":        true,
		"charts/x/values-prod.yaml": true,
		"README.md":                 false,
		".agents/docs/changes/x.md": false,
		"settings.json":             false,
		".vscode/settings.json":     false,
		"assets/logo.png":           false,
	}
	for path, want := range cases {
		if got := proj.IsSourcePath(path); got != want {
			t.Errorf("IsSourcePath(%q) = %v, want %v", path, got, want)
		}
	}
}

// TestArtifactMatchedByFrontmatterNotFilename is the regression test for the
// failure that motivated frontmatter matching: an artifact whose filename does
// not follow the convention still counts, and one that merely looks right but
// declares nothing does not.
func TestArtifactMatchedByFrontmatterNotFilename(t *testing.T) {
	proj := newRepo(t)
	plans := filepath.Join(proj.Docs, "plans")

	// Off-convention filename, correct frontmatter: must count.
	write(t, filepath.Join(plans, "harness-recon-notes.md"),
		"---\nowner: explorer\ndate: 2026-09-03\nscope: worktree:abc123\n---\n\n# Recon\n")
	// Convention-shaped filename, no frontmatter: must not count.
	write(t, filepath.Join(plans, "2026-09-03-recon-decoy.md"), "# Recon\n")

	if !proj.HasArtifact(RowExplorer, "worktree:abc123") {
		t.Error("artifact with correct frontmatter but off-convention filename was not found")
	}
	if proj.HasArtifact(RowExplorer, "worktree:different") {
		t.Error("artifact matched a scope it does not declare")
	}
	if !proj.HasArtifactDated(RowExplorer, "2026-09-03") {
		t.Error("dated lookup did not find the artifact")
	}
	if got := len(proj.Artifacts("plans")); got != 1 {
		t.Errorf("Artifacts returned %d entries; the decoy without frontmatter should be ignored", got)
	}
}

func TestOwedProportionality(t *testing.T) {
	proj := newRepo(t)
	cfg := Config{ReviewLines: 50, ReviewFiles: 5, DispatchTimeout: 60}
	day := "2026-09-03"

	// A one-line change: changelog and verifier are due, reviewer is not.
	write(t, filepath.Join(proj.Root, "main.go"), "package main\n")
	rep := proj.Owed(cfg, day)
	if !hasRow(rep.Due, RowChangelog) {
		t.Error("changelog is a resumability row and must always be due")
	}
	if !hasRow(rep.Due, RowVerifier) {
		t.Error("verifier must be due when source changed")
	}
	if hasRow(rep.Due, RowReviewer) {
		t.Error("reviewer fired on a change below both thresholds")
	}
	if !hasRow(rep.Skipped, RowReviewer) {
		t.Error("a skipped row must be reported with its reason, not silently dropped")
	}

	// Security surface touched: the security row becomes due.
	write(t, filepath.Join(proj.Root, "auth", "login.go"), "package auth\n")
	rep = proj.Owed(cfg, day)
	if !hasRow(rep.Due, RowSecurity) {
		t.Error("security row did not fire for a path touching auth")
	}
}

func TestOwedSatisfiedByExistingArtifact(t *testing.T) {
	proj := newRepo(t)
	day := Today()
	write(t, filepath.Join(proj.Root, "main.go"), "package main\n")
	// The artifact must declare the scope of the very change it covers.
	write(t, filepath.Join(proj.Docs, "changes", "anything.md"),
		"---\nowner: changelog\ndate: "+day+"\nscope: "+proj.Scope()+"\n---\n")

	rep := proj.Owed(DefaultConfig(), day)
	if hasRow(rep.Due, RowChangelog) {
		t.Error("changelog was still due despite a matching artifact on disk")
	}
}

func TestOwedSilentOutsideContract(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	proj, ok := Resolve(root)
	if !ok {
		t.Fatal("Resolve failed")
	}
	write(t, filepath.Join(root, "main.go"), "package main\n")

	if rep := proj.Owed(DefaultConfig(), "2026-09-03"); len(rep.Due) != 0 {
		t.Errorf("gate fired in a repo that never adopted the contract: %+v", rep.Due)
	}
}

func TestLoadConfigRejectsMalformed(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, ".agents", "harness.json"), "{not json")
	if _, err := LoadConfig(root); err == nil {
		t.Error("malformed harness.json was silently ignored")
	}
}

func TestLoadConfigDefaultsWhenAbsent(t *testing.T) {
	cfg, err := LoadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("missing config should not be an error: %v", err)
	}
	if !reflect.DeepEqual(cfg, DefaultConfig()) {
		t.Errorf("got %+v, want defaults %+v", cfg, DefaultConfig())
	}
}

func hasRow(list []Due, row Row) bool {
	for _, d := range list {
		if d.Row == row {
			return true
		}
	}
	return false
}

// TestSourceGlobsMakeProseSource: a project whose deliverable is markdown must
// be able to say so. Without this, rewriting every agent definition and the
// constitution itself reports that nothing is owed.
func TestSourceGlobsMakeProseSource(t *testing.T) {
	proj := newRepo(t)
	proj.SourceGlobs = []string{"shared/AGENTS.md", "claude/agents/**"}

	for _, want := range []string{"shared/AGENTS.md", "claude/agents/reviewer.md"} {
		if !proj.IsSourcePath(want) {
			t.Errorf("%s should be source once declared in sourceGlobs", want)
		}
	}
	for _, notWant := range []string{"README.md", "docs/notes.md", ".agents/docs/changes/x.md"} {
		if proj.IsSourcePath(notWant) {
			t.Errorf("%s should not be source", notWant)
		}
	}
}

// TestChangedLinesCountsUntracked: git's numstat says nothing about files it
// has never seen, so an all-new change set would otherwise score zero and skip
// the reviewer row — the case least likely to have been read by anyone.
func TestChangedLinesCountsUntracked(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "new.go"), strings.Repeat("package main\n", 120))

	if got := proj.ChangedLines(); got < 120 {
		t.Errorf("ChangedLines() = %d, want at least 120 for an untracked file", got)
	}
}

// TestMatchGlobSupportsDoublestar: init-project documents patterns like
// `**/*.go` and `claude/agents/**`, and Go's path.Match implements neither —
// its `*` never crosses a separator, so `**/*.go` would match exactly one
// directory level. The matcher has to mean what the documentation says.
func TestMatchGlobSupportsDoublestar(t *testing.T) {
	cases := []struct {
		glob, path string
		want       bool
	}{
		{"**/*.go", "main.go", true},
		{"**/*.go", "internal/trail/git.go", true},
		{"**/*.go", "a/b/c/d/e.go", true},
		{"**/*.go", "web/app.ts", false},
		{"claude/agents/**", "claude/agents/reviewer.md", true},
		{"claude/agents/**", "claude/agents/sub/x.md", true},
		{"claude/agents/**", "claude/settings.json", false},
		{"shared/AGENTS.md", "shared/AGENTS.md", true},
		{"shared/AGENTS.md", "shared/skills/AGENTS.md", false},
		{"web/**/*.tsx", "web/app.tsx", true},
		{"web/**/*.tsx", "web/routes/deep/page.tsx", true},
		{"web/**/*.tsx", "api/app.tsx", false},
		{"*.md", "README.md", true},
		{"*.md", "docs/README.md", false},
	}
	for _, tc := range cases {
		if got := matchGlob(tc.glob, tc.path); got != tc.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", tc.glob, tc.path, got, tc.want)
		}
	}
}

// TestMatchGlobHandlesUnicodeAndClasses: an earlier version quoted one BYTE of
// a multi-byte rune, producing a pattern that could never match, and turned
// character classes into literal brackets. Both failed silently — a gate that
// matches nothing has switched itself off without saying so.
func TestMatchGlobHandlesUnicodeAndClasses(t *testing.T) {
	cases := []struct {
		glob, path string
		want       bool
	}{
		{"src/år/*.go", "src/år/main.go", true},
		{"src/år/*.go", "src/ar/main.go", false},
		{"**/naïve/*.ts", "web/naïve/x.ts", true},
		{"[a-z]*.go", "main.go", true},
		{"[a-z]*.go", "Main.go", false},
		{"[!a-z]*.go", "Main.go", true},
		{"[!a-z]*.go", "main.go", false},
		{"cmd/?arness/*.go", "cmd/harness/main.go", true},
		{"cmd/?arness/*.go", "cmd/hharness/main.go", false},
		{"a**b", "axxb", true},
		{"a**b", "ax/xb", true},
	}
	for _, tc := range cases {
		if got := matchGlob(tc.glob, tc.path); got != tc.want {
			t.Errorf("matchGlob(%q, %q) = %v, want %v", tc.glob, tc.path, got, tc.want)
		}
	}
}

// TestCompileGlobRejectsUnterminatedClass asserts on compileGlob directly. An
// earlier test called matchGlob("[unterminated", …) and passed for the wrong
// reason: the pattern compiled fine as a literal, so it proved nothing.
func TestCompileGlobRejectsUnterminatedClass(t *testing.T) {
	if _, err := compileGlob("src/[abc/*.go"); err == nil {
		t.Error("an unterminated character class compiled without error")
	}
	if matchGlob("src/[abc/*.go", "src/[abc/main.go") {
		t.Error("a malformed glob matched instead of being rejected")
	}
}

// TestOwedSeesInitialCommit: `git diff-tree HEAD` prints nothing for a root
// commit unless --root is passed, so without it a project's very first commit
// owed no rows at all.
func TestOwedSeesInitialCommit(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "main.go"), strings.Repeat("package main\n", 80))
	commit(t, proj.Root, "initial")

	rep := proj.Owed(DefaultConfig(), Today())
	if len(rep.Changed) == 0 {
		t.Fatal("the initial commit reported no changed source files")
	}
	if !hasRow(rep.Due, RowChangelog) {
		t.Error("changelog row did not fire for the initial commit")
	}
}

// TestOwedMeasuresCommittedLines: a clean tree has no working diff, so the
// proportionality rule has to read the commit itself or every committed change
// scores zero and skips the reviewer row however large it is.
func TestOwedMeasuresCommittedLines(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")
	write(t, filepath.Join(proj.Root, "big.go"), strings.Repeat("package main\n", 200))
	commit(t, proj.Root, "a large committed change")

	rep := proj.Owed(DefaultConfig(), Today())
	if rep.Lines < 200 {
		t.Errorf("Lines = %d for a 200-line commit against a clean tree", rep.Lines)
	}
	if !hasRow(rep.Due, RowReviewer) {
		t.Error("reviewer row did not fire for a 200-line committed change")
	}
}

// TestInitialCommitLinesCounted covers --root on the numstat path, which the
// earlier tests exercised only via --name-only.
func TestInitialCommitLinesCounted(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "main.go"), strings.Repeat("package main\n", 90))
	commit(t, proj.Root, "initial")

	if got := proj.HeadLinesToday(); got < 90 {
		t.Errorf("HeadLinesToday() = %d on a root commit of 90 lines", got)
	}
}

// TestLineCountingIsStableAcrossCommit: both counters must filter paths the
// same way, or a change scores differently before and after it is committed and
// the reviewer row fires on one side of the commit but not the other.
func TestLineCountingIsStableAcrossCommit(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	// A source file and a non-source file of equal size.
	write(t, filepath.Join(proj.Root, "app.go"), strings.Repeat("package main\n", 60))
	write(t, filepath.Join(proj.Root, "notes.md"), strings.Repeat("prose\n", 60))

	before := proj.ChangedLines()
	commit(t, proj.Root, "add both")
	after := proj.HeadLinesToday()

	if before != after {
		t.Errorf("line count changed across the commit: %d before, %d after", before, after)
	}
	if before < 60 || before >= 120 {
		t.Errorf("got %d lines; expected the .go file counted and the .md file ignored", before)
	}
}

// TestMergeCommitIsVisible: diff-tree prints nothing for a merge commit, with
// or without --root, so a merge as HEAD reported no changed files and owed no
// rows. It is diffed against its first parent instead.
func TestMergeCommitIsVisible(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "base.go"), "package main\n")
	commit(t, proj.Root, "base")

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = proj.Root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	run("checkout", "-q", "-b", "feature")
	write(t, filepath.Join(proj.Root, "feature.go"), strings.Repeat("package main\n", 70))
	commit(t, proj.Root, "feature work")
	run("checkout", "-q", "-")
	run("merge", "-q", "--no-ff", "-m", "merge feature", "feature")

	if !proj.headIsMerge() {
		t.Fatal("fixture did not produce a merge commit")
	}
	changed := proj.HeadSourceToday()
	if len(changed) == 0 {
		t.Error("a merge commit reported no changed source files")
	}
	if got := proj.HeadLinesToday(); got < 70 {
		t.Errorf("HeadLinesToday() = %d across a merge of a 70-line branch", got)
	}
}
