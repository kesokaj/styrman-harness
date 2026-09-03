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
		"---\nowner: explorer\ndate: 2026-09-03\n---\n\n# Recon\n")
	// Convention-shaped filename, no frontmatter: must not count.
	write(t, filepath.Join(plans, "2026-09-03-recon-decoy.md"), "# Recon\n")

	if !proj.HasArtifact(RowExplorer, "2026-09-03") {
		t.Error("artifact with correct frontmatter but off-convention filename was not found")
	}
	if proj.HasArtifact(RowExplorer, "2026-09-04") {
		t.Error("artifact matched a date it does not declare")
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
	day := "2026-09-03"
	write(t, filepath.Join(proj.Root, "main.go"), "package main\n")
	write(t, filepath.Join(proj.Docs, "changes", "anything.md"),
		"---\nowner: changelog\ndate: "+day+"\n---\n")

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
