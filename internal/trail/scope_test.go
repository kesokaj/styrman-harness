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

// TestSecondChangeSameDayOwesRows is the regression the scope mechanism exists
// for. Under date matching, a second change on the same day was satisfied by the
// morning's artifact.
//
// The changelog carries a stated exception while work is in progress — a task
// narrative does not stop being true because a line moved — so this asserts the
// boundary precisely: accountability rows re-owe immediately, the changelog
// re-owes at the commit. The original bug was "a second substantive COMMIT owed
// nothing", and that stays fixed.
func TestSecondChangeSameDayOwesRows(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	// First change, fully recorded for its own scope.
	write(t, filepath.Join(proj.Root, "one.go"), strings.Repeat("package main\n", 80))
	first := proj.Scope()
	for _, row := range []Row{RowChangelog, RowReviewer, RowVerifier} {
		write(t, filepath.Join(proj.Docs, row.Class(), string(row)+".md"),
			"---\nowner: "+string(row)+"\ndate: "+Today()+"\nscope: "+first+"\n---\n")
	}
	if due := proj.Owed(DefaultConfig(), Today()).Due; len(due) != 0 {
		t.Fatalf("rows due for the scope their own artifacts cover: %+v", due)
	}

	// A second, different change on the same day, still uncommitted.
	write(t, filepath.Join(proj.Root, "two.go"), strings.Repeat("package main\n", 80))
	rep := proj.Owed(DefaultConfig(), Today())
	if !hasRow(rep.Due, RowReviewer) {
		t.Error("reviewer not re-owed for a second 80-line change — the date gap is back")
	}
	if !hasRow(rep.Due, RowVerifier) {
		t.Error("verifier not re-owed for changed code")
	}
	if hasRow(rep.Due, RowChangelog) {
		t.Error("changelog re-owed mid-work; the narrative exception should hold until the commit")
	}

	// Commit it. The exception lapses and the changelog must cover HEAD.
	commit(t, proj.Root, "second change, unrecorded")
	if !hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("a second commit on the same day owed no changelog — the original bug is back")
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

// TestChangelogSurvivesAnEditWhileWorkIsInProgress: a task narrative does not
// stop being true because a line moved. Holding the changelog to a byte-exact
// scope made it unsatisfiable on a working tree — the entry was stale as soon
// as it was written, and the agent's anti-duplication rule declined to write
// another.
func TestChangelogSurvivesAnEditWhileWorkIsInProgress(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	write(t, filepath.Join(proj.Root, "work.go"), "package main\n")
	write(t, filepath.Join(proj.Docs, "changes", "entry.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: "+proj.Scope()+"\n---\n")
	if hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Fatal("changelog due for the scope its own entry covers")
	}

	// Keep working. The narrative has not changed.
	write(t, filepath.Join(proj.Root, "work.go"), "package main\n\nvar x = 1\n")
	if hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("changelog re-owed by an edit while work was still in progress")
	}
	// The accountability rows must NOT get that exception: a review of the old
	// code is not a review of the new code.
	write(t, filepath.Join(proj.Root, "work.go"), "package main\n"+strings.Repeat("var y = 1\n", 80))
	write(t, filepath.Join(proj.Docs, "reviews", "r.md"),
		"---\nowner: reviewer\ndate: "+Today()+"\nscope: changeset:stale000000\n---\n")
	if !hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowReviewer) {
		t.Error("reviewer accepted a stale scope — accountability rows must re-run on code change")
	}
}

// TestChangelogExceptionLapsesAtTheCommitBoundary: the tree is clean at a
// commit, the scope is HEAD's own and stable, and the entry must cover it
// exactly. This is where the blocking gate fires, so the exception must not
// reach it.
func TestChangelogExceptionLapsesAtTheCommitBoundary(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	// An entry dated today, but covering nothing that is now committed.
	write(t, filepath.Join(proj.Docs, "changes", "old.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: changeset:stale000000\n---\n")
	write(t, filepath.Join(proj.Root, "landed.go"), "package main\n")
	commit(t, proj.Root, "unrecorded commit")

	if !hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("a today-dated entry satisfied the changelog row at a commit boundary")
	}
}

// TestNarrativeExceptionRequiresCurrentTaskEntry: the first version of the
// exception asked only whether a scope was non-empty, which made it a date match
// in disguise — any of today's entries satisfied any of today's work, including
// one carrying a scope form the binary no longer produces.
func TestNarrativeExceptionRequiresCurrentTaskEntry(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	// A finished task's entry: today-dated, scoped, and COMMITTED.
	write(t, filepath.Join(proj.Docs, "changes", "finished.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: changeset:aaaaaaaaaaaa\n---\n")
	commit(t, proj.Root, "a finished task")

	// New work begins.
	write(t, filepath.Join(proj.Root, "new.go"), "package main\n")
	if !hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("a committed entry for a finished task satisfied new in-progress work")
	}

	// An entry carrying a superseded scope form must not satisfy anything either.
	write(t, filepath.Join(proj.Docs, "changes", "oldform.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: worktree:9366002f2d81\n---\n")
	if !hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("an entry with a pre-mechanism scope form satisfied the row")
	}

	// An uncommitted, current-form entry is what the exception is for.
	write(t, filepath.Join(proj.Docs, "changes", "inprogress.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: "+proj.Scope()+"\n---\n")
	if hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("the in-progress entry did not satisfy the row it was written for")
	}
	// ...and it keeps satisfying it as the work continues.
	write(t, filepath.Join(proj.Root, "new.go"), "package main\n\nvar x = 1\n")
	if hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Error("the in-progress entry stopped covering the task after an edit")
	}
}

// TestFreshnessDoesNotInheritTheException: the freshness report says whether the
// trail covers the code as it stands. Reading the exception off Owed() made it
// print OK for a dirty tree no entry covered.
func TestFreshnessDoesNotInheritTheException(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "seed.go"), "package main\n")
	commit(t, proj.Root, "seed")

	write(t, filepath.Join(proj.Root, "work.go"), "package main\n")
	write(t, filepath.Join(proj.Docs, "changes", "entry.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: "+proj.Scope()+"\n---\n")

	// The dispatcher is satisfied; the tree then moves.
	write(t, filepath.Join(proj.Root, "work.go"), "package main\n\nvar x = 1\n")
	if hasRow(proj.Owed(DefaultConfig(), Today()).Due, RowChangelog) {
		t.Fatal("fixture wrong: the exception should still hold for the dispatcher")
	}
	if f := proj.CheckFreshness(); !f.Stale {
		t.Errorf("freshness reported OK for a tree no entry covers:\n%s", strings.Join(f.Lines, "\n"))
	}
}
