package trail

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestScopeChangesWithTheTree is the property dates never had: an artifact
// stops covering work once that work moves on.
func TestScopeChangesWithTheTree(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	write(t, filepath.Join(proj.Root, "a.go"), "package main\n")
	first := proj.Scope()
	if first == "" || !strings.HasPrefix(first, "changeset:") {
		t.Fatalf("dirty tree produced scope %q", first)
	}

	write(t, filepath.Join(proj.Root, "a.go"), "package main\n\nvar x = 1\n")
	if second := proj.Scope(); second == first {
		t.Error("scope did not change when the working tree did")
	}
}

// TestScopeSurvivesTheCommitThatLandsIt is the regression for the worst bug the
// scope model shipped with: keying uncommitted and committed work differently
// meant that committing re-owed every row for work the gate had just accepted,
// and it self-perpetuated for the rest of the day.
func TestScopeSurvivesTheCommitThatLandsIt(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	write(t, filepath.Join(proj.Root, "feature.go"), strings.Repeat("package main\n", 80))
	before := proj.Scope()

	// Satisfy every owed row against that scope.
	for _, row := range []Row{RowChangelog, RowReviewer, RowVerifier} {
		write(t, filepath.Join(proj.Docs, row.Class(), string(row)+"-report.md"),
			"---\nowner: "+string(row)+"\ndate: "+Today()+"\nscope: "+before+"\n---\n")
	}
	if due := proj.Owed(DefaultConfig(), Today()).Due; len(due) != 0 {
		t.Fatalf("rows still due before the commit: %+v", due)
	}

	commit(t, proj.Root, "land the feature")

	if after := proj.Scope(); after != before {
		t.Errorf("scope changed across the commit: %q -> %q", before, after)
	}
	if due := proj.Owed(DefaultConfig(), Today()).Due; len(due) != 0 {
		t.Errorf("committing re-owed rows the gate had already accepted: %+v", due)
	}
}

// TestScopeChangesAcrossDifferentCommits: distinct work must still get distinct
// scopes, or the fix above would satisfy everything forever.
func TestScopeChangesAcrossDifferentCommits(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "a.go"), "package main\n")
	commit(t, proj.Root, "one")
	scope := proj.Scope()

	write(t, filepath.Join(proj.Root, "b.go"), "package main\n")
	commit(t, proj.Root, "two")
	if next := proj.Scope(); next == scope {
		t.Error("two different commits produced the same scope")
	}
}

// TestSecondChangeSameDayOwesRows is the regression this whole mechanism
// exists for. Under date matching, a second change on the same day was
// satisfied by the morning's artifact — it hid owed rows five times on the day
// the gap was found, including for changes to the gate itself.
func TestSecondChangeSameDayOwesRows(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	// First change, fully recorded for its own scope.
	write(t, filepath.Join(proj.Root, "one.go"), strings.Repeat("package main\n", 80))
	firstScope := proj.Scope()
	write(t, filepath.Join(proj.Docs, "changes", "first.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: "+firstScope+"\n---\n")

	if rep := proj.Owed(DefaultConfig(), Today()); hasRow(rep.Due, RowChangelog) {
		t.Fatal("changelog still due for the scope its own artifact covers")
	}

	// A second, different change on the same day.
	write(t, filepath.Join(proj.Root, "two.go"), strings.Repeat("package main\n", 80))
	rep := proj.Owed(DefaultConfig(), Today())
	if !hasRow(rep.Due, RowChangelog) {
		t.Error("changelog not due for a new change set — the date gap is back")
	}
	if !hasRow(rep.Due, RowReviewer) {
		t.Error("reviewer not due for a second 80-line change on the same day")
	}
}

// TestArtifactWithoutScopeCoversNothing: artifacts predating the mechanism, or
// written by an agent that ignored the instruction, must not silently satisfy
// a row they may have nothing to do with.
func TestArtifactWithoutScopeCoversNothing(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "a.go"), "package main\n")
	write(t, filepath.Join(proj.Docs, "changes", "legacy.md"),
		"---\nowner: changelog\ndate: "+Today()+"\n---\n")

	if !hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("a scope-less artifact satisfied a row")
	}
}

// TestCoversRequiresAnExactMatch: the range form was removed as unreachable, so
// coverage is exact-match only. This asserts the whole of that contract.
func TestCoversRequiresAnExactMatch(t *testing.T) {
	proj := newRepo(t)
	cases := []struct {
		name           string
		artifact, want string
		covers         bool
	}{
		{"identical", "changeset:abc123", "changeset:abc123", true},
		{"different digest", "changeset:abc123", "changeset:def456", false},
		{"artifact has none", "", "changeset:abc123", false},
		{"nothing to cover", "changeset:abc123", "", false},
		{"a range covers nothing now", "aaa..bbb", "changeset:abc123", false},
	}
	for _, tc := range cases {
		got := proj.Covers(Artifact{Owner: RowReviewer, Scope: tc.artifact}, tc.want)
		if got != tc.covers {
			t.Errorf("%s: Covers(%q, %q) = %v, want %v", tc.name, tc.artifact, tc.want, got, tc.covers)
		}
	}
}

// TestScopeIsIndependentOfCallerSetup: classification used to be a field callers
// assigned by hand, and the gates did not — so on a project declaring
// sourceGlobs the gates hashed a different file set than the dispatcher, the
// blocking gate demanded a scope nothing could write, and it blocked forever.
func TestScopeIsIndependentOfCallerSetup(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, ".agents", "harness.json"),
		`{"sourceGlobs":["docs/**/*.md"]}`)
	write(t, filepath.Join(proj.Root, "docs", "spec.md"), "# spec\n")

	// A freshly resolved project must classify the same way as one whose caller
	// remembered to load config.
	resolved, ok := Resolve(proj.Root)
	if !ok {
		t.Fatal("Resolve failed")
	}
	manual := resolved
	cfg, err := LoadConfig(proj.Root)
	if err != nil {
		t.Fatal(err)
	}
	manual.SourceGlobs = cfg.SourceGlobs

	if resolved.Scope() != manual.Scope() {
		t.Errorf("scope depends on caller setup: %q vs %q", resolved.Scope(), manual.Scope())
	}
	if resolved.Scope() == "" {
		t.Error("a declared source glob was not honoured by a resolved project")
	}
}
