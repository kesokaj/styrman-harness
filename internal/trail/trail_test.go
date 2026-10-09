package trail

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// signedRepo is a repository whose own config makes `git log` run a program:
// log.showSignature plus gpg.program, and one commit signed by that program.
// The program writes marker when git asks it to verify.
func signedRepo(t *testing.T) (root, marker string) {
	t.Helper()
	root = t.TempDir()
	marker = filepath.Join(t.TempDir(), "gpg-ran")
	prog := filepath.Join(t.TempDir(), "fake-gpg")
	// The program drains stdin before answering: git writes the payload to
	// it, and a program that exits first can fail the signing under load
	// ("gpg failed to sign", security review round 2).
	script := "#!/bin/sh\nfor a in \"$@\"; do [ \"$a\" = \"--verify\" ] && { cat >/dev/null; touch '" + marker + "'; exit 0; }; done\n" +
		"cat >/dev/null\necho '[GNUPG:] SIG_CREATED D 1 8 00 0 0' >&2\nprintf -- '-----BEGIN PGP SIGNATURE-----\\n\\nx\\n-----END PGP SIGNATURE-----\\n'\n"
	if err := os.WriteFile(prog, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.invalid")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("init", "-q")
	os.MkdirAll(filepath.Join(root, ".agents", "docs", "handoffs"), 0o755)
	os.WriteFile(filepath.Join(root, ".agents", "docs", "handoffs", "h.md"), []byte("h\n"), 0o644)
	git("add", ".agents")
	git("-c", "gpg.program="+prog, "-c", "user.signingkey=x", "commit", "-q", "-S", "-m", "signed")
	git("config", "gpg.program", prog)
	git("config", "log.showSignature", "true")
	// The fixture must trip the marker with the very format the code asks
	// for, minus the hardening (silence is not evidence).
	if out, err := exec.Command("git", "-C", root, "log", "--format=%ct", "--name-only", "--no-renames", "--", ".agents/docs").CombinedOutput(); err != nil {
		t.Fatalf("plain git log: %v\n%s", err, out)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("the fixture does not run gpg.program on a plain git log, so it proves nothing: %v", err)
	}
	os.Remove(marker)
	return root, marker
}

// Dates runs git log in project roots taken from transcripts: a repository's
// own config must not make it start a program (security review 2026-10-08).
func TestDatesRunsNoSignatureProgram(t *testing.T) {
	root, marker := signedRepo(t)
	dates, _ := Dates(root, ".agents/docs")
	if len(dates) == 0 {
		t.Fatal("Dates read nothing; the test would pass vacuously")
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("git log ran the repository's gpg.program")
	}
}

// Log is the one way the metrics store and the dreamer run git log in a
// project root (internal/metrics commitsIn uses it too).
func TestLogRunsNoSignatureProgram(t *testing.T) {
	root, marker := signedRepo(t)
	out, err := Log(root, "--format=%H %ct", "--all")
	if err != nil || len(out) == 0 {
		t.Fatalf("Log: %v, %q", err, out)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("git log ran the repository's gpg.program")
	}
}

// filterRepo is a repository whose .gitattributes routes every .md file
// through a clean filter its own config defines; the filter writes marker.
// One tracked record is then edited, one left as committed, one added
// untracked.
func filterRepo(t *testing.T) (root, marker string) {
	t.Helper()
	root = t.TempDir()
	marker = filepath.Join(t.TempDir(), "filter-ran")
	prog := filepath.Join(t.TempDir(), "fake-clean")
	if err := os.WriteFile(prog, []byte("#!/bin/sh\ntouch '"+marker+"'\ncat\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.invalid")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	docs := filepath.Join(root, ".agents", "docs", "handoffs")
	os.MkdirAll(docs, 0o755)
	os.WriteFile(filepath.Join(docs, "kept.md"), []byte("kept\n"), 0o644)
	os.WriteFile(filepath.Join(docs, "edited.md"), []byte("before\n"), 0o644)
	git("init", "-q")
	git("add", ".agents")
	git("commit", "-q", "-m", "records")
	os.WriteFile(filepath.Join(root, ".gitattributes"), []byte("*.md filter=x\n"), 0o644)
	git("config", "filter.x.clean", prog)
	os.WriteFile(filepath.Join(docs, "edited.md"), []byte("after, longer\n"), 0o644)
	os.WriteFile(filepath.Join(docs, "new.md"), []byte("new\n"), 0o644)
	// kept.md's content and size are unchanged but its mtime moves: status
	// cannot settle it from the index's stat data and must hash it, through
	// the filter, every time (a size change alone is decided without one).
	later := time.Now().Add(time.Hour)
	os.Chtimes(filepath.Join(docs, "kept.md"), later, later)
	// The fixture must trip the marker with a plain status (silence is not evidence).
	if out, err := exec.Command("git", "-C", root, "status", "--porcelain", "--", ".agents/docs").CombinedOutput(); err != nil {
		t.Fatalf("plain git status: %v\n%s", err, out)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("the fixture does not run the clean filter on a plain git status, so it proves nothing: %v", err)
	}
	os.Remove(marker)
	return root, marker
}

// Dates runs no program the repository's config defines (no clean filter,
// security review round 2), and still tells an edited record from one as
// committed.
func TestDatesRunsNoCleanFilter(t *testing.T) {
	root, marker := filterRepo(t)
	dates, modified := Dates(root, ".agents/docs")
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("Dates ran the repository's clean filter")
	}
	kept, edited := ".agents/docs/handoffs/kept.md", ".agents/docs/handoffs/edited.md"
	if _, ok := dates[kept]; !ok || modified[kept] {
		t.Errorf("the record as committed: dated %v, modified %v; want dated and not modified", ok, modified[kept])
	}
	if !modified[edited] {
		t.Error("the edited record is not reported modified")
	}
	if _, ok := dates[".agents/docs/handoffs/new.md"]; ok {
		t.Error("an untracked record got a commit date")
	}
}

// Every git call into a project root carries GIT_NO_LAZY_FETCH=1: a partial
// clone would otherwise fetch through its own core.sshCommand.
func TestCommandsNeverLazyFetch(t *testing.T) {
	cmd := Command("/r", "log")
	found := false
	for _, e := range cmd.Env {
		found = found || e == "GIT_NO_LAZY_FETCH=1"
	}
	if !found {
		t.Fatal("a hardened git command can lazy-fetch")
	}
	if got := strings.Join(cmd.Args[1:], " "); got != "--no-optional-locks -c core.fsmonitor=false -c log.showSignature=false -C /r log" {
		t.Errorf("args %q", got)
	}
}

// A partial clone is not read at all: a git older than 2.44 ignores
// GIT_NO_LAZY_FETCH and would fetch through the repository's own
// core.sshCommand (security review round 3).
func TestPartialCloneIsNotRead(t *testing.T) {
	root, _ := filterRepo(t)
	if out, err := exec.Command("git", "-C", root, "config", "extensions.partialClone", "origin").CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if dates, _ := Dates(root, ".agents/docs"); len(dates) != 0 {
		t.Fatalf("a partial clone was read: %v", dates)
	}
	if out, err := Log(root, "--oneline"); err == nil {
		t.Fatalf("Log ran in a partial clone: %q", out)
	}
}

// Toplevel answers a parent of the directory or the directory itself: a
// repository's core.worktree cannot point a session at another project.
func TestToplevelIsAParent(t *testing.T) {
	root, _ := filterRepo(t)
	sub := filepath.Join(root, ".agents", "docs")
	if got := Toplevel(sub); got != root {
		t.Errorf("Toplevel(sub) = %q, want %q", got, root)
	}
	elsewhere := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "config", "core.worktree", elsewhere).CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if got := Toplevel(sub); got != filepath.Clean(sub) {
		t.Errorf("with core.worktree elsewhere, Toplevel(sub) = %q, want the directory itself", got)
	}
}
