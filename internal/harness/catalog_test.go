package harness

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// The Google skills catalog (plan 2026-10-01-google-skills-catalog): 141 skills
// published by Google, vendored verbatim under skills/google-catalog/catalog/
// and reached through one router skill, because installed as global skills
// they push the harness's own skills out of both CLIs' listings.
//
// CATALOG_DIR points the catalog tests at another directory and
// CATALOG_WRITE=1 makes the two derived files be written instead of compared:
// that is how scripts/refresh-catalog.sh builds INDEX.md and CHECKSUMS, so the
// index is derived by the installer's own frontmatter parser and nothing is
// parsed twice.

// catalogHarnessFiles are written by the harness, not vendored: CHECKSUMS does
// not cover them.
var catalogHarnessFiles = map[string]bool{"INDEX.md": true, "MANIFEST.md": true, "CHECKSUMS": true, ".sh-agent-noindex": true}

const catalogIndexHeader = "# Google skills catalog — index\n\nOne line per entry: the directory name under `catalog/`, then the start of its description. Grep this file for the product or task words, then read `catalog/<name>/SKILL.md`.\n\n"

func catalogDir(t *testing.T) string {
	t.Helper()
	if d := os.Getenv("CATALOG_DIR"); d != "" {
		return d
	}
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(repo, "skills", "google-catalog", "catalog")
}

// catalogEntries are the directories holding a SKILL.md, sorted. An empty
// catalog fails the caller: a test over zero files proves nothing.
func catalogEntries(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	des, _ := os.ReadDir(dir)
	for _, de := range des {
		if !de.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, de.Name(), "SKILL.md")); err == nil {
			out = append(out, de.Name())
		} else {
			t.Errorf("catalog entry %s: no SKILL.md", de.Name())
		}
	}
	if len(out) == 0 {
		t.Fatalf("the catalog is empty: no entry with a SKILL.md under %s", dir)
	}
	sort.Strings(out)
	return out
}

// catalogSums maps every vendored file (slash path relative to the catalog)
// to its sha256.
func catalogSums(t *testing.T, dir string) map[string]string {
	t.Helper()
	sums := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)
		if catalogHarnessFiles[rel] {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		sums[rel] = hex.EncodeToString(sum[:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return sums
}

// TestCatalogMatchesChecksums pins the vendored tree to what the refresh
// fetched, in both directions: a changed or added file has no matching line,
// and a line with no file is a deletion.
func TestCatalogMatchesChecksums(t *testing.T) {
	dir := catalogDir(t)
	catalogEntries(t, dir)
	got := catalogSums(t, dir)
	file := filepath.Join(dir, "CHECKSUMS")
	if os.Getenv("CATALOG_WRITE") == "1" {
		paths := make([]string, 0, len(got))
		for p := range got {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		var b strings.Builder
		for _, p := range paths {
			fmt.Fprintf(&b, "%s  %s\n", got[p], p)
		}
		if err := os.WriteFile(file, []byte(b.String()), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("no CHECKSUMS: %v", err)
	}
	want := map[string]string{}
	for _, line := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
		sum, path, ok := strings.Cut(line, "  ")
		if !ok || len(sum) != 64 {
			t.Fatalf("CHECKSUMS line is not `<sha256>  <path>`: %q", line)
		}
		want[path] = sum
	}
	for p, sum := range got {
		switch w, ok := want[p]; {
		case !ok:
			t.Errorf("%s is in the catalog and not in CHECKSUMS", p)
		case w != sum:
			t.Errorf("%s differs from its CHECKSUMS line: the vendored file was edited", p)
		}
	}
	for p := range want {
		if _, ok := got[p]; !ok {
			t.Errorf("CHECKSUMS names %s, which is not in the catalog", p)
		}
	}
}

// catalogIndexLine is one INDEX.md line: the entry's directory name and its
// description as Parse reads it, cut at 220 runes.
func catalogIndexLine(name, description string) string {
	r := []rune(description)
	if len(r) > 220 {
		r = r[:220]
	}
	return "- `" + name + "` — " + strings.TrimSpace(string(r))
}

// TestCatalogIndexMatchesEntries re-derives the index from every entry's
// frontmatter through the installer's own parser: a missing line, an extra
// line and a stale description all differ from it.
func TestCatalogIndexMatchesEntries(t *testing.T) {
	dir := catalogDir(t)
	entries := catalogEntries(t, dir)
	var lines []string
	for _, name := range entries {
		data, err := os.ReadFile(filepath.Join(dir, name, "SKILL.md"))
		if err != nil {
			t.Fatal(err)
		}
		front, _, err := Split(string(data))
		if err != nil {
			t.Fatalf("catalog entry %s: SKILL.md has no frontmatter: %v", name, err)
		}
		fm := Parse(front)
		if fm["name"] == "" || fm["description"] == "" {
			t.Fatalf("catalog entry %s: SKILL.md frontmatter lacks name or description", name)
		}
		lines = append(lines, catalogIndexLine(name, fm["description"]))
	}
	want := catalogIndexHeader + strings.Join(lines, "\n") + "\n"
	file := filepath.Join(dir, "INDEX.md")
	if os.Getenv("CATALOG_WRITE") == "1" {
		if err := os.WriteFile(file, []byte(want), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("no INDEX.md: %v", err)
	}
	if string(data) == want {
		return
	}
	have := map[string]bool{}
	for _, l := range strings.Split(string(data), "\n") {
		have[l] = true
	}
	wanted := map[string]bool{}
	for _, l := range strings.Split(want, "\n") {
		wanted[l] = true
		if !have[l] {
			t.Errorf("INDEX.md lacks the line %q", l)
		}
	}
	for l := range have {
		if !wanted[l] {
			t.Errorf("INDEX.md has a line no entry produces: %q", l)
		}
	}
	if !t.Failed() {
		t.Error("INDEX.md holds the right lines in the wrong order or count")
	}
}

// TestCatalogFilesAreCarryable: what the harness cannot carry byte-identical or
// would mis-handle is refused at the source — an executable bit (the embed has
// none, BUG-137), a symlink (the installer refuses it), and a Go file or module
// file (the Go tool would compile the one and the embed would drop the
// directory for the other).
func TestCatalogFilesAreCarryable(t *testing.T) {
	dir := catalogDir(t)
	catalogEntries(t, dir)
	root := dir
	if os.Getenv("CATALOG_DIR") == "" {
		root = filepath.Dir(dir) // the whole router skill, its scripts included
	}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			t.Errorf("%s is a symlink", rel)
		case d.IsDir():
		case info.Mode()&0o111 != 0:
			t.Errorf("%s has an executable bit", rel)
		case d.Name() == "go.mod" || d.Name() == "go.sum" || strings.HasSuffix(d.Name(), ".go"):
			t.Errorf("%s is a Go file inside the module", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestCatalogInstallsClean: the install plan over the real checkout has no
// refusal — a secret-shaped name or a symlink anywhere under skills/ fails it,
// through the installer's own list rather than a second copy in the script.
func TestCatalogInstallsClean(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Plan(DirSource(repo), nil, VendorAll); err != nil {
		t.Fatalf("the checkout does not install: %v", err)
	}
}

// TestRouterDescription: the router is the one catalog line both CLIs list, so
// its description is short (Claude Code's listing budget, TestSkillListingBudget)
// and says what it must not be used for.
func TestRouterDescription(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repo, "skills", "google-catalog", "SKILL.md"))
	if err != nil {
		t.Fatalf("no skills/google-catalog/SKILL.md: %v", err)
	}
	front, _, err := Split(string(data))
	if err != nil {
		t.Fatal(err)
	}
	desc := Parse(front)["description"]
	if n := len(desc); n < 240 || n > 300 {
		t.Errorf("the router's description is %d bytes, want 240 to 300", n)
	}
	if !strings.Contains(desc, "Do NOT use") {
		t.Errorf("the router's description names no skill that wins over it: %q", desc)
	}
}

// TestRefreshCatalogScript runs the refresh against two local fixture
// repositories: the happy path, idempotence, and one fixture per refusal
// class, each of which must exit 1, name the entry, and leave the destination
// byte-identical. It never reaches the network.
func TestRefreshCatalogScript(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(repo, "skills", "google-catalog", "scripts", "refresh-catalog.sh")
	if _, err := os.Stat(script); err != nil {
		t.Fatalf("no refresh-catalog.sh: %v", err)
	}
	if testing.Short() {
		t.Skip("runs git and a nested go test")
	}
	for _, tool := range []string{"bash", "git", "go"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not on PATH", tool)
		}
	}
	tmp := t.TempDir()
	git := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.invalid")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	put := func(dir string, files map[string]string) {
		t.Helper()
		for rel, body := range files {
			p := filepath.Join(dir, rel)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	commit := func(dir string) string {
		t.Helper()
		git(dir, "add", "-A")
		git(dir, "commit", "-qm", "c", "--allow-empty")
		return git(dir, "rev-parse", "HEAD")
	}
	skill := func(name, description string) string {
		return "---\nname: " + name + "\n" + description + "---\n\n# " + name + "\n"
	}
	agents, gskills := filepath.Join(tmp, "agents-cli"), filepath.Join(tmp, "google-skills")
	for _, d := range []string{agents, gskills} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		git(d, "init", "-q")
	}
	put(agents, map[string]string{
		"LICENSE":               "Apache-2.0 (fixture)\n",
		"skills/alpha/SKILL.md": skill("alpha", "description: plain one-line description of alpha\n"),
		"extensions/langchain/template/.agents/skills/beta/SKILL.md": skill("beta", "description: \"a quoted description: with a colon\"\n"),
	})
	pinA := commit(agents)
	put(gskills, map[string]string{
		"LICENSE":                              "Apache-2.0 (fixture)\n",
		"skills/cloud/gamma/SKILL.md":          skill("gamma", "description: >-\n  a block scalar description\n  folded over two lines\n"),
		"skills/cloud/gamma/scripts/run.sh":    "#!/bin/sh\necho hi\n",
		"skills/cloud/delta/SKILL.md":          skill("delta", "metadata:\n  version: \"1.0\"\ndescription: delta with a metadata map\n"),
		"skills/developers/epsilon/SKILL.md":   skill("epsilon", "description: epsilon\n"),
		"skills/ads/not-asked-for/SKILL.md":    skill("not-asked-for", "description: outside the four source sets\n"),
		"skills/cloud/gamma/references/ref.md": "reference text\n",
	})
	if err := os.Chmod(filepath.Join(gskills, "skills/cloud/gamma/scripts/run.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	pinG := commit(gskills)

	dest, work := filepath.Join(tmp, "dest"), filepath.Join(tmp, "work")
	run := func(args ...string) (string, int) {
		t.Helper()
		cmd := exec.Command("bash", append([]string{script, "--dest", dest, "--work", work}, args...)...)
		cmd.Dir = repo
		cmd.Env = append(os.Environ(), "CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
		var out bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &out
		err := cmd.Run()
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return out.String(), ee.ExitCode()
		}
		if err != nil {
			t.Fatal(err)
		}
		return out.String(), 0
	}
	snapshot := func() string {
		t.Helper()
		var b strings.Builder
		filepath.WalkDir(dest, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			data, _ := os.ReadFile(path)
			sum := sha256.Sum256(data)
			info, _ := d.Info()
			rel, _ := filepath.Rel(dest, path)
			fmt.Fprintf(&b, "%s %o %x\n", rel, info.Mode().Perm(), sum[:6])
			return nil
		})
		return b.String()
	}

	out, code := run(pinA, pinG)
	if code != 0 {
		t.Fatalf("happy path: exit %d\n%s", code, out)
	}
	for _, name := range []string{"alpha", "beta", "gamma", "delta", "epsilon"} {
		if _, err := os.Stat(filepath.Join(dest, name, "SKILL.md")); err != nil {
			t.Errorf("entry %s not vendored: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dest, "not-asked-for")); err == nil {
		t.Error("an entry outside the four source sets was vendored")
	}
	index, _ := os.ReadFile(filepath.Join(dest, "INDEX.md"))
	for _, want := range []string{
		"- `gamma` — a block scalar description folded over two lines\n",
		"- `beta` — a quoted description: with a colon\n",
		"- `delta` — delta with a metadata map\n",
	} {
		if !strings.Contains(string(index), want) {
			t.Errorf("INDEX.md lacks %q:\n%s", want, index)
		}
	}
	if n := strings.Count(string(index), "\n- `"); n != 5 {
		t.Errorf("INDEX.md has %d entry lines, want 5", n)
	}
	sums, _ := os.ReadFile(filepath.Join(dest, "CHECKSUMS"))
	for _, want := range []string{"  gamma/scripts/run.sh\n", "  LICENSE-agents-cli\n", "  LICENSE-google-skills\n"} {
		if !strings.Contains(string(sums), want) {
			t.Errorf("CHECKSUMS lacks %q", want)
		}
	}
	if strings.Contains(string(sums), "INDEX.md") || strings.Contains(string(sums), "MANIFEST.md") {
		t.Error("CHECKSUMS covers a harness-written file")
	}
	manifest, _ := os.ReadFile(filepath.Join(dest, "MANIFEST.md"))
	for _, want := range []string{pinA[:7], pinG[:7], "## Excluded", "## Notes"} {
		if !strings.Contains(string(manifest), want) {
			t.Errorf("MANIFEST.md lacks %q", want)
		}
	}
	if info, err := os.Stat(filepath.Join(dest, "gamma", "scripts", "run.sh")); err != nil || info.Mode()&0o111 != 0 {
		t.Errorf("the executable bit was not cleared: %v %v", info, err)
	}
	if _, err := os.Stat(filepath.Join(dest, ".sh-agent-noindex")); err != nil {
		t.Errorf("the search marker was not written: %v", err)
	}
	if !strings.Contains(out, "next: git add") {
		t.Errorf("the script does not say what to run next:\n%s", out)
	}
	good := snapshot()
	if out, code := run(pinA, pinG); code != 0 || snapshot() != good {
		t.Fatalf("a second run at the same pins changed the catalog (exit %d)\n%s", code, out)
	}

	// Relative --dest and --work resolve against the caller's directory, not
	// against the directory the nested go test runs in (BUG-146).
	relCmd := exec.Command("bash", script, "--dest", "rel-dest", "--work", "rel-work", pinA, pinG)
	relCmd.Dir = tmp
	relCmd.Env = append(os.Environ(), "CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
	if out, err := relCmd.CombinedOutput(); err != nil {
		t.Errorf("relative --dest and --work: %v\n%s", err, out)
	} else if _, err := os.Stat(filepath.Join(tmp, "rel-dest", "gamma", "SKILL.md")); err != nil {
		t.Errorf("relative --dest: the catalog did not land under the caller's directory: %v", err)
	}

	// Refusals: each leaves the destination as it was and names the entry.
	refuse := func(label, names string, files map[string]string, mutate func(), exclude string) {
		t.Helper()
		put(gskills, files)
		if mutate != nil {
			mutate()
		}
		pin := commit(gskills)
		out, code := run(pinA, pin)
		if code != 1 || !strings.Contains(out, names) {
			t.Errorf("%s: want exit 1 naming %q, got exit %d\n%s", label, names, code, out)
		}
		if snapshot() != good {
			t.Errorf("%s: the refused run changed the destination", label)
		}
		if exclude != "" {
			if out, code := run("--exclude", exclude, pinA, pin); code != 0 {
				t.Errorf("%s: --exclude %s did not clear the refusal (exit %d)\n%s", label, exclude, code, out)
			}
			if _, err := os.Stat(filepath.Join(dest, exclude)); err == nil && exclude != "alpha" {
				t.Errorf("%s: the excluded entry %s was vendored", label, exclude)
			}
			// back to the good state for the next fixture
			if out, code := run(pinA, pinG); code != 0 || snapshot() != good {
				t.Fatalf("%s: could not restore the good catalog (exit %d)\n%s", label, code, out)
			}
		}
		git(gskills, "reset", "-q", "--hard", pinG)
		git(gskills, "clean", "-qfd")
	}
	if out, code := run(pinA, "0000000000000000000000000000000000000000"); code != 1 || !strings.Contains(out, "does not resolve") || snapshot() != good {
		t.Errorf("unresolvable pin: want exit 1 and an untouched destination, got exit %d\n%s", code, out)
	}
	refuse("no frontmatter", "bad", map[string]string{"skills/cloud/bad/SKILL.md": "no frontmatter here\n"}, nil, "bad")
	refuse("collision", "alpha", map[string]string{"skills/cloud/alpha/SKILL.md": skill("alpha", "description: a second alpha\n")}, nil, "")
	refuse("symlink", "gamma", nil, func() {
		if err := os.Symlink("SKILL.md", filepath.Join(gskills, "skills/cloud/gamma/link.md")); err != nil {
			t.Skip("no symlinks here")
		}
	}, "gamma")
	refuse("go.mod", "delta", map[string]string{"skills/cloud/delta/go.mod": "module x\n"}, nil, "")
	refuse("go file", "epsilon", map[string]string{"skills/developers/epsilon/main.go": "package main\n"}, nil, "")
}
