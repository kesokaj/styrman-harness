package search

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Findings of reviews/2026-09-29-security-search-index.md, one test each.

func TestPathWithURISpecialCharacters(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "q?mark#frag%20x")
	if err := os.Rename(newRepo(t), root); err != nil {
		t.Fatal(err)
	}
	out := mustSearch(t, Options{Dir: root, Home: t.TempDir(), Words: []string{"demo"}})
	if !strings.Contains(out, "README.md") {
		t.Errorf("search in a checkout with ?, # and %% in its path:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents", "scratch", "search.db")); err != nil {
		t.Errorf("the index is not where it belongs: %v", err)
	}
	entries, _ := os.ReadDir(parent)
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("files appeared beside the checkout: %v", names)
	}
}

func TestRegistersOnlyHarnessProjects(t *testing.T) {
	home := t.TempDir()
	plain := newRepo(t)
	writeFile(t, plain, "docs/guide.md", "# a pelican guide\n")
	commitAll(t, plain, "docs")
	mustSearch(t, Options{Dir: plain, Home: home, Words: []string{"pelican"}})
	harnessed := newRepo(t)
	writeFile(t, harnessed, ".agents/docs/plans/p.md", "# the pelican plan\n")
	commitAll(t, harnessed, "trail")
	mustSearch(t, Options{Dir: harnessed, Home: home, Words: []string{"pelican"}})
	out := mustSearch(t, Options{Dir: t.TempDir(), Home: home, Global: true, Words: []string{"pelican"}})
	if strings.Contains(out, "guide.md") || !strings.Contains(out, "plans/p.md") {
		t.Errorf("only a checkout with .agents/docs registers:\n%s", out)
	}
}

func TestHashesOddNamesItself(t *testing.T) {
	root := newRepo(t)
	home := t.TempDir()
	writeFile(t, root, `"quoted.md`, "a lemming in quotes\n")
	writeFile(t, root, "cr\r.md", "a lemming with a carriage return\n")
	writeFile(t, root, ".agents/docs/plans/\"odd\r.md", "# lemming plan\n")
	mustSearch(t, Options{Dir: root, Home: home, Words: []string{"demo"}}) // registers
	out := mustSearch(t, Options{Dir: root, Home: home, Words: []string{"lemming"}})
	if !strings.Contains(out, "quoted.md") || !strings.Contains(out, "carriage") {
		t.Errorf("odd file names not indexed with their own content:\n%s", out)
	}
	if _, err := runSearch(t, Options{Dir: t.TempDir(), Home: home, Global: true, Words: []string{"lemming"}}); err != nil {
		t.Errorf("an odd name under a registered docs root broke --global: %v", err)
	}
}

func TestOutputStripsControlCharacters(t *testing.T) {
	root := newRepo(t)
	writeFile(t, root, ".agents/docs/bugs/BUG-001-x.md", "---\nid: BUG-001\nstatus: open\x1b[2J\ndate: 2026-09-01\n---\n# BUG-001: \x1b]0;pwned\x07title\n\nthe \x1b[31mbadger\x1b[0m escaped\n")
	commitAll(t, root, "add")
	home := t.TempDir()
	for _, o := range []Options{
		{Dir: root, Home: home, Words: []string{"badger"}},
		{Dir: root, Home: home, Refs: "BUG-001"},
		{Dir: root, Home: home, Open: true},
	} {
		out := mustSearch(t, o)
		if strings.ContainsAny(out, "\x1b\x07") {
			t.Errorf("%+v printed a control character:\n%q", o, out)
		}
	}
}

func TestGlobalDirIsPrivate(t *testing.T) {
	root := newRepo(t)
	writeFile(t, root, ".agents/docs/plans/p.md", "# plan\n")
	commitAll(t, root, "trail")
	home := t.TempDir()
	mustSearch(t, Options{Dir: root, Home: home, Words: []string{"plan"}})
	fi, err := os.Stat(filepath.Join(home, ".config", "sh-agent"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o700 {
		t.Errorf("~/.config/sh-agent created %04o, want 0700", fi.Mode().Perm())
	}
}

func TestSkipsSecretShapedFiles(t *testing.T) {
	root := newRepo(t)
	home := t.TempDir()
	writeFile(t, root, ".env", "TOKEN=wolverine-secret\n")
	writeFile(t, root, "keys/server.pem", "wolverine private key\n")
	writeFile(t, root, "notes.md", "a wolverine note\n")
	out := mustSearch(t, Options{Dir: root, Home: home, Words: []string{"wolverine"}})
	if strings.Contains(out, ".env") || strings.Contains(out, "server.pem") || !strings.Contains(out, "notes.md") {
		t.Errorf("secret-shaped files indexed:\n%s", out)
	}
}

// Round 2 of the security review.

func TestGlobalDropsCheckoutWithoutAgentsDocs(t *testing.T) {
	home := t.TempDir()
	root := newRepo(t)
	writeFile(t, root, ".agents/docs/plans/p.md", "# plan\n")
	writeFile(t, root, "docs/evil.md", "an ocelot planted in docs\n")
	commitAll(t, root, "both")
	mustSearch(t, Options{Dir: root, Home: home, Words: []string{"plan"}})
	gitT(t, root, "rm", "-q", "-r", ".agents/docs")
	commitAll(t, root, "drop trail")
	out := mustSearch(t, Options{Dir: t.TempDir(), Home: home, Global: true, Words: []string{"ocelot"}})
	if strings.Contains(out, "evil.md") || !strings.Contains(out, "(no longer has .agents/docs)") {
		t.Errorf("a registration without .agents/docs still answered:\n%s", out)
	}
}

func TestHintAndNotesArePrintable(t *testing.T) {
	root := newRepo(t)
	home := t.TempDir()
	writeFile(t, root, "a\x1b[31mb/z.go", "package x\n")
	commitAll(t, root, "odd dir")
	out := mustSearch(t, Options{Dir: root, Home: home, Refs: "z.go"})
	if !strings.Contains(out, "hint: did you mean") || strings.ContainsRune(out, 0x1b) {
		t.Errorf("hint line:\n%q", out)
	}
	parent := t.TempDir()
	gone := filepath.Join(parent, "gone\x1b[2J")
	if err := os.Rename(newRepo(t), gone); err != nil {
		t.Fatal(err)
	}
	writeFile(t, gone, ".agents/docs/plans/p.md", "# plan\n")
	commitAll(t, gone, "trail")
	mustSearch(t, Options{Dir: gone, Home: home, Words: []string{"plan"}})
	os.RemoveAll(gone)
	out = mustSearch(t, Options{Dir: t.TempDir(), Home: home, Global: true, Words: []string{"plan"}})
	if !strings.Contains(out, "note: dropped") || strings.ContainsRune(out, 0x1b) {
		t.Errorf("note line:\n%q", out)
	}
}

func TestSkipsMoreSecretNames(t *testing.T) {
	root := newRepo(t)
	writeFile(t, root, ".git-credentials", "https://u:caracal@example.invalid\n")
	writeFile(t, root, ".aws/credentials", "caracal key\n")
	writeFile(t, root, "config/.ENV", "CARACAL=1\n")
	writeFile(t, root, "notes.md", "a caracal note\n")
	out := mustSearch(t, Options{Dir: root, Home: t.TempDir(), Words: []string{"caracal"}})
	for _, bad := range []string{".git-credentials", ".aws/credentials", ".ENV"} {
		if strings.Contains(out, bad) {
			t.Errorf("%s indexed:\n%s", bad, out)
		}
	}
	if !strings.Contains(out, "notes.md") {
		t.Errorf("control file missing:\n%s", out)
	}
}

// Round 3 of the security review.

func TestRegistrationKeptWhenCheckErrors(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads a mode-000 directory; the check cannot fail here")
	}
	home := t.TempDir()
	root := newRepo(t)
	writeFile(t, root, ".agents/docs/plans/p.md", "# the tamarin plan\n")
	commitAll(t, root, "trail")
	mustSearch(t, Options{Dir: root, Home: home, Words: []string{"tamarin"}})
	agents := filepath.Join(root, ".agents")
	if err := os.Chmod(agents, 0); err != nil {
		t.Fatal(err)
	}
	// The condition holds: the check cannot read .agents/docs.
	if _, err := os.Stat(filepath.Join(agents, "docs")); err == nil || os.IsNotExist(err) {
		t.Fatalf("fixture: want a permission error on .agents/docs, got %v", err)
	}
	out := mustSearch(t, Options{Dir: t.TempDir(), Home: home, Global: true, Words: []string{"tamarin"}})
	os.Chmod(agents, 0o755)
	if strings.Contains(out, "dropped") || !strings.Contains(out, "note: skipped") {
		t.Errorf("a registration whose check errored was dropped, or skipped silently:\n%s", out)
	}
	out = mustSearch(t, Options{Dir: t.TempDir(), Home: home, Global: true, Words: []string{"tamarin"}})
	if !strings.Contains(out, "plans/p.md") {
		t.Errorf("the registration did not survive the error:\n%s", out)
	}
}

func TestRebuiltNoteIsPrintable(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "esc\x1b[2Jrepo")
	if err := os.Rename(newRepo(t), root); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, ".agents/scratch/search.db", "not a database, and long enough to be read as a header by sqlite")
	out := mustSearch(t, Options{Dir: root, Home: t.TempDir(), Words: []string{"demo"}})
	if !strings.Contains(out, "rebuilt") || strings.ContainsRune(out, 0x1b) {
		t.Errorf("rebuilt note:\n%q", out)
	}
}

// Review of the security fixes (reviews/2026-09-29-review-search-index-code-r2.md): git's
// "not a git repository" is matched in git's C-locale text, whatever the caller's locale.
func TestDropReasonIgnoresLocale(t *testing.T) {
	bin := t.TempDir()
	fake := "#!/bin/sh\nif [ \"$LC_ALL\" = C ]; then echo 'fatal: not a git repository (or any of the parent directories): .git' >&2; else echo 'fatal: kein Git-Repository (oder irgendeines der Elternverzeichnisse): .git' >&2; fi\nexit 128\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LC_ALL", "de_DE.UTF-8")
	reason, err := dropReason(t.TempDir())
	if err != nil || reason != "no longer a git checkout" {
		t.Errorf("under a German locale: reason %q, err %v; want the checkout dropped", reason, err)
	}
}
