package harness

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type zipEntry struct {
	name string
	body string
	mode fs.FileMode // 0 means a regular 0o644 file
}

// writeZip builds an archive at path from entries, in order.
func writeZip(t *testing.T, path string, entries []zipEntry) {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, e := range entries {
		h := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
		mode := e.mode
		if mode == 0 {
			mode = 0o644
		}
		h.SetMode(mode)
		w, err := zw.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(e.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

// harnessEntries is newRepo's tree as zip entries under prefix.
func harnessEntries(t *testing.T, prefix string) []zipEntry {
	t.Helper()
	repo := newRepo(t)
	var entries []zipEntry
	err := filepath.WalkDir(repo, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(repo, p)
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		entries = append(entries, zipEntry{name: prefix + filepath.ToSlash(rel), body: string(data)})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return entries
}

// noZipDirs asserts UnpackZip left no zip-* directory under parent.
func noZipDirs(t *testing.T, parent string) {
	t.Helper()
	left, _ := filepath.Glob(filepath.Join(parent, "zip-*"))
	if len(left) != 0 {
		t.Errorf("temporary directories left under %s: %v", parent, left)
	}
}

func TestUnpackZipGitHubTopFolder(t *testing.T) {
	dir := t.TempDir()
	zp := filepath.Join(dir, "styrman-harness-main.zip")
	entries := harnessEntries(t, "styrman-harness-main/")
	entries = append(entries,
		zipEntry{name: "styrman-harness-main/README.md", body: "readme\n"},
		zipEntry{name: "styrman-harness-main/cmd/sh-agent/main.go", body: "package main\n"},
		// This repository's own tracked symlink: GitHub's zip keeps it as a
		// symlink entry, and it sits outside the five trees (D6).
		zipEntry{name: "styrman-harness-main/.agents/skills", body: "../skills", mode: fs.ModeSymlink | 0o777},
	)
	writeZip(t, zp, entries)
	parent := t.TempDir()

	tree, _, err := UnpackZip(zp, parent)
	if err != nil {
		t.Fatalf("UnpackZip: %v", err)
	}
	if !IsRepo(tree) {
		t.Fatalf("%s is not harness-shaped after unpacking", tree)
	}
	if !strings.HasPrefix(tree, parent+string(filepath.Separator)) {
		t.Errorf("tree %s is not under the parent %s", tree, parent)
	}
	for _, rel := range []string{"skills/alpha/references/notes.md", "agents/reviewer.md", "commands/deploy.md", "rules/go-style.md", "context/dev-setup.md"} {
		if _, err := os.Stat(filepath.Join(tree, rel)); err != nil {
			t.Errorf("%s not unpacked: %v", rel, err)
		}
	}
	// Outside the five trees: name-checked, never written.
	for _, rel := range []string{"README.md", "cmd", ".agents"} {
		if _, err := os.Lstat(filepath.Join(tree, rel)); err == nil {
			t.Errorf("%s was written, but it is outside the five trees", rel)
		}
	}
}

func TestUnpackZipRootLayout(t *testing.T) {
	dir := t.TempDir()
	zp := filepath.Join(dir, "flat.zip")
	writeZip(t, zp, harnessEntries(t, ""))
	tree, _, err := UnpackZip(zp, t.TempDir())
	if err != nil {
		t.Fatalf("UnpackZip: %v", err)
	}
	if !IsRepo(tree) {
		t.Fatalf("%s is not harness-shaped", tree)
	}
}

func TestUnpackZipRefuses(t *testing.T) {
	base := func(t *testing.T) []zipEntry { return harnessEntries(t, "top/") }
	cases := []struct {
		name  string
		extra zipEntry
		want  string
	}{
		{"parent element", zipEntry{name: "top/../../escape", body: "x"}, "escape"},
		{"absolute name", zipEntry{name: "/etc/escape", body: "x"}, "escape"},
		{"backslash name", zipEntry{name: `top\..\escape`, body: "x"}, "escape"},
		{"symlink inside a tree", zipEntry{name: "top/skills/alpha/link", body: "/etc/passwd", mode: fs.ModeSymlink | 0o777}, "symlink"},
		{"device-like mode inside a tree", zipEntry{name: "top/skills/alpha/pipe", body: "", mode: fs.ModeNamedPipe | 0o644}, "not a regular file"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			zp := filepath.Join(dir, "bad.zip")
			writeZip(t, zp, append(base(t), c.extra))
			// Prove the fixture holds the condition before trusting the refusal.
			zr, err := zip.OpenReader(zp)
			if err != nil {
				t.Fatal(err)
			}
			last := zr.File[len(zr.File)-1]
			if last.Name != c.extra.name || (c.extra.mode != 0 && last.Mode().Type() != c.extra.mode.Type()) {
				t.Fatalf("fixture: last entry %q mode %v, want %q mode %v", last.Name, last.Mode(), c.extra.name, c.extra.mode)
			}
			zr.Close()

			parent := t.TempDir()
			tree, _, err := UnpackZip(zp, parent)
			if err == nil {
				t.Fatalf("UnpackZip accepted %s (tree %s)", c.name, tree)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("error %q does not name %q", err, c.want)
			}
			noZipDirs(t, parent)
			if _, err := os.Lstat(filepath.Join(filepath.Dir(parent), "escape")); err == nil {
				t.Errorf("an escaping entry was written beside the parent")
			}
		})
	}
}

func TestUnpackZipLimits(t *testing.T) {
	saved := [3]int64{int64(maxZipEntries), maxZipFileBytes, maxZipTotalBytes}
	t.Cleanup(func() {
		maxZipEntries, maxZipFileBytes, maxZipTotalBytes = int(saved[0]), saved[1], saved[2]
	})

	t.Run("too many entries", func(t *testing.T) {
		entries := harnessEntries(t, "top/")
		maxZipEntries = len(entries) + 2
		for i := 0; i < 3; i++ {
			entries = append(entries, zipEntry{name: "top/docs/f" + string(rune('a'+i)), body: "x"})
		}
		zp := filepath.Join(t.TempDir(), "many.zip")
		writeZip(t, zp, entries)
		parent := t.TempDir()
		if _, _, err := UnpackZip(zp, parent); err == nil || !strings.Contains(err.Error(), "entries") {
			t.Errorf("err = %v, want an entry-count refusal", err)
		}
		noZipDirs(t, parent)
	})
	maxZipEntries = int(saved[0])

	t.Run("one file too large", func(t *testing.T) {
		maxZipFileBytes = 1024
		entries := append(harnessEntries(t, "top/"), zipEntry{name: "top/skills/alpha/big.md", body: strings.Repeat("a", 2048)})
		zp := filepath.Join(t.TempDir(), "big.zip")
		writeZip(t, zp, entries)
		parent := t.TempDir()
		if _, _, err := UnpackZip(zp, parent); err == nil || !strings.Contains(err.Error(), "big.md") {
			t.Errorf("err = %v, want a per-file size refusal naming big.md", err)
		}
		noZipDirs(t, parent)
	})
	maxZipFileBytes = saved[1]

	t.Run("total too large", func(t *testing.T) {
		entries := harnessEntries(t, "top/")
		var total int64
		for _, e := range entries {
			total += int64(len(e.body))
		}
		maxZipTotalBytes = total + 100
		entries = append(entries, zipEntry{name: "top/skills/alpha/more.md", body: strings.Repeat("b", 200)})
		zp := filepath.Join(t.TempDir(), "total.zip")
		writeZip(t, zp, entries)
		parent := t.TempDir()
		if _, _, err := UnpackZip(zp, parent); err == nil || !strings.Contains(err.Error(), "total") {
			t.Errorf("err = %v, want a total-size refusal", err)
		}
		noZipDirs(t, parent)
	})
}

func TestUnpackZipNotHarness(t *testing.T) {
	zp := filepath.Join(t.TempDir(), "other.zip")
	writeZip(t, zp, []zipEntry{{name: "proj/README.md", body: "x"}, {name: "proj/main.go", body: "package main"}})
	parent := t.TempDir()
	if _, _, err := UnpackZip(zp, parent); err == nil || !strings.Contains(err.Error(), "not a harness zip") {
		t.Errorf("err = %v, want not a harness zip", err)
	}
	noZipDirs(t, parent)

	notZip := filepath.Join(t.TempDir(), "plain.zip")
	if err := os.WriteFile(notZip, []byte("not a zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := UnpackZip(notZip, parent); err == nil || !strings.Contains(err.Error(), "not a readable zip") {
		t.Errorf("err = %v, want not a readable zip", err)
	}
	noZipDirs(t, parent)
}

func TestInstallRecordsZipOrigin(t *testing.T) {
	repo := newRepo(t)
	home := t.TempDir()
	src := DirSource(repo)
	src.Origin = ZipPrefix + "sha256:abc123"
	files, err := Plan(src, []byte("#!/bin/sh\n"), VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if err := Install(src, home, files, VendorAll, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	st, found, err := LoadState(home)
	if err != nil || !found {
		t.Fatalf("manifest: found=%v err=%v", found, err)
	}
	if st.Repo != src.Origin || st.Commit != "" {
		t.Errorf("manifest repo=%q commit=%q, want %q and no commit — never the unpack directory %s", st.Repo, st.Commit, src.Origin, repo)
	}
}

func TestResolveSourceSkipsZipManifest(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	cwd := t.TempDir()
	t.Chdir(cwd)
	if _, err := ResolveSource("", t.TempDir(), ""); err == nil {
		t.Skip("the temp directory happens to be inside a harness repo")
	}
	origin := ZipPrefix + "sha256:0123456789abcdef"
	// A harness-shaped directory whose relative name IS the manifest value:
	// a walk that treated the value as a path would find it here.
	decoy := filepath.Join(cwd, origin)
	if err := os.MkdirAll(filepath.Join(decoy, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(decoy, "context"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(decoy, "context", "AGENTS.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !IsRepo(origin) {
		t.Fatalf("fixture: %s must look like a harness tree relative to %s", origin, cwd)
	}
	home := t.TempDir()
	if err := SaveState(home, State{Repo: origin}); err != nil {
		t.Fatal(err)
	}
	src, err := ResolveSource("", home, "")
	if err == nil {
		t.Fatalf("ResolveSource returned %+v from a zip manifest value", src)
	}
	if !strings.Contains(err.Error(), "installed from a zip (sha256 0123456789ab)") || !strings.Contains(err.Error(), "update --from-zip") {
		t.Errorf("err = %v, want it to name the zip install and the --from-zip way forward first", err)
	}
	if !strings.Contains(err.Error(), "harness repo not found") {
		t.Errorf("err = %v, want the existing not-found message kept", err)
	}
}

func TestUnpackZipRefusesDuplicateName(t *testing.T) {
	entries := harnessEntries(t, "top/")
	entries = append(entries, zipEntry{name: "top/skills/alpha/SKILL.md", body: "---\nname: alpha\ndescription: second copy\n---\n"})
	zp := filepath.Join(t.TempDir(), "dup.zip")
	writeZip(t, zp, entries)
	zr, err := zip.OpenReader(zp)
	if err != nil {
		t.Fatal(err)
	}
	seen := 0
	for _, f := range zr.File {
		if f.Name == "top/skills/alpha/SKILL.md" {
			seen++
		}
	}
	zr.Close()
	if seen != 2 {
		t.Fatalf("fixture: %d entries named top/skills/alpha/SKILL.md, want 2", seen)
	}
	parent := t.TempDir()
	if _, _, err := UnpackZip(zp, parent); err == nil || !strings.Contains(err.Error(), "SKILL.md") {
		t.Errorf("err = %v, want the duplicate entry refused, never overwritten", err)
	}
	noZipDirs(t, parent)
}

// A refusal found after the names pass — a symlink, a duplicate, a size —
// must still leave no parent directory behind (BUG-093's second half).
func TestUnpackZipLateRefusalCreatesNoParent(t *testing.T) {
	cases := map[string]zipEntry{
		"symlink":   {name: "top/skills/alpha/link", body: "/etc/passwd", mode: fs.ModeSymlink | 0o777},
		"duplicate": {name: "top/skills/alpha/SKILL.md", body: "---\nname: alpha\ndescription: again\n---\n"},
	}
	for name, extra := range cases {
		t.Run(name, func(t *testing.T) {
			zp := filepath.Join(t.TempDir(), "late.zip")
			writeZip(t, zp, append(harnessEntries(t, "top/"), extra))
			home := t.TempDir()
			parent := filepath.Join(home, ".config", "sh-agent")
			if _, _, err := UnpackZip(zp, parent); err == nil {
				t.Fatal("UnpackZip accepted the zip")
			}
			if entries, _ := os.ReadDir(home); len(entries) != 0 {
				t.Errorf("a late refusal left %d entries under the home (the parent chain)", len(entries))
			}
		})
	}
}

func TestUnpackZipRefusesDeepOrLongNames(t *testing.T) {
	deep := "top/" + strings.Repeat("d/", 40) + "f"
	long := "top/" + strings.Repeat("n", 1100)
	for label, name := range map[string]string{"deep": deep, "long": long} {
		t.Run(label, func(t *testing.T) {
			zp := filepath.Join(t.TempDir(), "names.zip")
			writeZip(t, zp, append(harnessEntries(t, "top/"), zipEntry{name: name, body: "x"}))
			parent := t.TempDir()
			if _, _, err := UnpackZip(zp, parent); err == nil || !strings.Contains(err.Error(), "too ") {
				t.Errorf("err = %v, want a too-deep or too-long refusal", err)
			}
		})
	}
}

// The sum names the very bytes that were unpacked: one read, hashed and
// parsed from the same buffer.
func TestUnpackZipSumIsOfTheBytesRead(t *testing.T) {
	zp := filepath.Join(t.TempDir(), "sum.zip")
	writeZip(t, zp, harnessEntries(t, "top/"))
	data, err := os.ReadFile(zp)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("%x", sha256.Sum256(data))
	_, sum, err := UnpackZip(zp, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if sum != want {
		t.Errorf("sum = %q, want %q", sum, want)
	}
}

func TestUnpackZipRefusesOversizedArchive(t *testing.T) {
	saved := maxZipArchiveBytes
	t.Cleanup(func() { maxZipArchiveBytes = saved })
	zp := filepath.Join(t.TempDir(), "whole.zip")
	writeZip(t, zp, harnessEntries(t, "top/"))
	maxZipArchiveBytes = 100
	if _, _, err := UnpackZip(zp, t.TempDir()); err == nil || !strings.Contains(err.Error(), "archive") {
		t.Errorf("err = %v, want the archive-size refusal", err)
	}
}

// The Antigravity CLI's token file is secret-shaped: a skill tree carrying
// one is refused like a key file.
func TestSecretShapedCLITokenFiles(t *testing.T) {
	for _, name := range []string{"antigravity-oauth-token", ".credentials.json"} {
		if !secretShaped(name) {
			t.Errorf("secretShaped(%q) = false, want true", name)
		}
	}
}

// The report sweep matches names, never a glob built from the home: a home
// holding glob metacharacters must not reach another home's clone.
func TestSweepInterruptedClonesGlobCharsInHome(t *testing.T) {
	base := t.TempDir()
	home, other := filepath.Join(base, "h[1]"), filepath.Join(base, "h1")
	// A pid that is not running: the sweep removes such a directory.
	mine, theirs := filepath.Join(home, ".config", "sh-agent", "report-999999"), filepath.Join(other, ".config", "sh-agent", "report-999999")
	for _, d := range []string{mine, theirs} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	sweepInterruptedClones(home, &bytes.Buffer{})
	if _, err := os.Stat(theirs); err != nil {
		t.Errorf("another home's %s was swept", theirs)
	}
	if _, err := os.Stat(mine); err == nil {
		t.Errorf("the home's own stale %s was not swept", mine)
	}
}
