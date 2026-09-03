package trail

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func commit(t *testing.T, root, message string) {
	t.Helper()
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", message}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
}

func joined(f Freshness) string { return strings.Join(f.Lines, "\n") }

func TestCheckFreshnessQuietOutsideContract(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	proj, _ := Resolve(root)

	f := proj.CheckFreshness()
	if f.Stale {
		t.Errorf("a repo that never adopted the contract was reported stale:\n%s", joined(f))
	}
	if !strings.Contains(joined(f), "contract not adopted") {
		t.Errorf("expected an explanation of why nothing was checked:\n%s", joined(f))
	}
}

func TestCheckFreshnessFlagsUncoveredCommit(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "main.go"), "package main\n")
	commit(t, proj.Root, "add source with no changelog")

	f := proj.CheckFreshness()
	if !f.Stale {
		t.Errorf("a commit with no changelog artifact was reported fresh:\n%s", joined(f))
	}
	if !strings.Contains(joined(f), "no changelog entry exists") {
		t.Errorf("report did not name the problem:\n%s", joined(f))
	}
}

func TestCheckFreshnessAcceptsCoveredCommit(t *testing.T) {
	proj := newRepo(t)
	today := Today()
	write(t, filepath.Join(proj.Root, "main.go"), "package main\n")
	write(t, filepath.Join(proj.Root, "ROADMAP.md"), "# Roadmap\n")
	write(t, filepath.Join(proj.Root, "TODO.md"), "# Todo\n")
	// The entry must declare the scope of the change it covers — a bare date no
	// longer satisfies the blocking gate.
	write(t, filepath.Join(proj.Docs, "changes", "whatever-name.md"),
		"---\nowner: changelog\ndate: "+today+"\nscope: "+proj.Scope()+"\n---\n\nrecorded\n")
	commit(t, proj.Root, "source plus its changelog")

	f := proj.CheckFreshness()
	if f.Stale {
		t.Errorf("a commit covered by a changelog entry was reported stale:\n%s", joined(f))
	}
}

// TestCheckFreshnessRejectsFutureDate guards the failure mode where a date came
// from model memory rather than the clock: it silently breaks every comparison
// that decides whether the trail covers the code.
func TestCheckFreshnessRejectsFutureDate(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Docs, "changes", "ahead.md"),
		"---\nowner: changelog\ndate: 2999-01-01\n---\n")
	commit(t, proj.Root, "future-dated entry")

	f := proj.CheckFreshness()
	if !f.Stale {
		t.Errorf("a future-dated artifact was accepted:\n%s", joined(f))
	}
	if !strings.Contains(joined(f), "future-dated") {
		t.Errorf("report did not name the future date:\n%s", joined(f))
	}
}

// TestCheckFreshnessEnforcesAppendOnly: the changelog is a record, and a
// rewritten record is a stale trail regardless of what it now says.
func TestCheckFreshnessEnforcesAppendOnly(t *testing.T) {
	proj := newRepo(t)
	entry := filepath.Join(proj.Docs, "changes", "recorded.md")
	write(t, entry, "---\nowner: changelog\ndate: "+Today()+"\n---\n\noriginal\n")
	commit(t, proj.Root, "record something")

	if err := os.WriteFile(entry, []byte("---\nowner: changelog\ndate: "+Today()+"\n---\n\nrewritten\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	f := proj.CheckFreshness()
	if !f.Stale {
		t.Errorf("a rewritten changelog entry was accepted:\n%s", joined(f))
	}
	if !strings.Contains(joined(f), "append-only") {
		t.Errorf("report did not name the append-only violation:\n%s", joined(f))
	}
}

// TestFreshnessIsNotFooledByADateMatch: two commits in one day, only the first
// recorded. A date comparison called that covered.
func TestFreshnessIsNotFooledByADateMatch(t *testing.T) {
	proj := newRepo(t)
	write(t, filepath.Join(proj.Root, "one.go"), "package main\n")
	write(t, filepath.Join(proj.Docs, "changes", "first.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: "+proj.Scope()+"\n---\n")
	commit(t, proj.Root, "first, recorded")

	if f := proj.CheckFreshness(); f.Stale {
		t.Fatalf("a recorded commit was reported stale:\n%s", joined(f))
	}

	write(t, filepath.Join(proj.Root, "two.go"), "package main\n")
	commit(t, proj.Root, "second, unrecorded")

	if f := proj.CheckFreshness(); !f.Stale {
		t.Errorf("an unrecorded second commit passed:\n%s", joined(f))
	}
}

func TestCheckFreshnessDocOnlyCommit(t *testing.T) {
	proj := newRepo(t)
	// Seed a recorded change so a baseline changelog entry exists.
	write(t, filepath.Join(proj.Root, "main.go"), "package main\n")
	write(t, filepath.Join(proj.Docs, "changes", "first.md"),
		"---\nowner: changelog\ndate: "+Today()+"\nscope: "+proj.Scope()+"\n---\n")
	commit(t, proj.Root, "first commit, recorded")

	if f := proj.CheckFreshness(); f.Stale {
		t.Fatalf("first commit reported stale:\n%s", joined(f))
	}

	// Now commit only a markdown documentation file outside sourceGlobs.
	write(t, filepath.Join(proj.Root, "notes.md"), "# Some notes\n")
	commit(t, proj.Root, "add doc file only")

	// CheckFreshness should not report STALE on an empty scope when HEAD has no source files.
	f := proj.CheckFreshness()
	if f.Stale {
		t.Errorf("a doc-only commit with no source changes was reported stale:\n%s", joined(f))
	}
	if !strings.Contains(joined(f), "no source changes in current state") {
		t.Errorf("expected explanation that no source changes exist:\n%s", joined(f))
	}
}
