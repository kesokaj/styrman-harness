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
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The Google skills catalog (plan 2026-10-01-google-skills-catalog): the skills
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

// catalogWriteMode reports whether the derived files are to be written. That
// happens only on the refresh script's build directory: CATALOG_WRITE without
// CATALOG_DIR would regenerate the checked-in CHECKSUMS from whatever is on
// disk and turn the integrity test into its own answer (security review).
func catalogWriteMode(t *testing.T) bool {
	t.Helper()
	if os.Getenv("CATALOG_WRITE") != "1" {
		return false
	}
	if os.Getenv("CATALOG_DIR") == "" {
		t.Fatal("CATALOG_WRITE=1 needs CATALOG_DIR: the checked-in INDEX.md and CHECKSUMS are written only by the refresh script, into its build directory")
	}
	return true
}

func catalogDir(t *testing.T) string {
	t.Helper()
	if d := os.Getenv("CATALOG_DIR"); d != "" {
		if os.Getenv("CATALOG_WRITE") != "1" {
			// Compare mode on another tree would pass for a catalog that is
			// not the checked-in one (security review T1).
			t.Fatal("CATALOG_DIR is only for the refresh script's write mode (CATALOG_WRITE=1)")
		}
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
	if catalogWriteMode(t) {
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
	if catalogWriteMode(t) {
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
	// The search marker travels with the catalog: without it `sh-agent search`
	// indexes the entries and the checkout's own files drop in the hits.
	if _, err := os.Stat(filepath.Join(dir, ".sh-agent-noindex")); err != nil {
		t.Errorf("the catalog has no .sh-agent-noindex marker: %v", err)
	}
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

// catalogExcluded reads the names under the manifest's "## Excluded".
func catalogExcluded(t *testing.T, dir string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "MANIFEST.md"))
	if err != nil {
		t.Fatalf("no MANIFEST.md: %v", err)
	}
	var out []string
	in := false
	for _, line := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(line, "## "):
			in = line == "## Excluded"
		case in && strings.HasPrefix(line, "- `"):
			if name, _, ok := strings.Cut(line[3:], "`"); ok {
				out = append(out, name)
			}
		}
	}
	return out
}

// TestCatalogExcludedEntriesAreAbsent: an entry the security review excluded
// stays out — the manifest's list is what the refresh applies, and a refresh
// that brought one back fails here.
func TestCatalogExcludedEntriesAreAbsent(t *testing.T) {
	dir := catalogDir(t)
	catalogEntries(t, dir)
	for _, name := range catalogExcluded(t, dir) {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			t.Errorf("%s is listed under ## Excluded in MANIFEST.md and is in the catalog", name)
		}
	}
}

// TestCatalogWriteModeNeedsCatalogDir: CATALOG_WRITE=1 alone must not rewrite
// the checked-in CHECKSUMS.
func TestCatalogWriteModeNeedsCatalogDir(t *testing.T) {
	if testing.Short() {
		t.Skip("runs a nested go test")
	}
	if os.Getenv("CATALOG_WRITE") == "1" {
		t.Skip("inside the write mode itself")
	}
	file := filepath.Join(catalogDir(t), "CHECKSUMS")
	before, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.WriteFile(file, before, 0o644) }) // a broken guard must not leave the file rewritten
	cmd := exec.Command("go", "test", ".", "-run", "^TestCatalogMatchesChecksums$", "-count=1")
	cmd.Env = fixtureEnv("CATALOG_WRITE=1", "CATALOG_DIR=")
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "needs CATALOG_DIR") {
		t.Errorf("CATALOG_WRITE=1 without CATALOG_DIR was accepted (err %v):\n%s", err, out)
	}
	if after, _ := os.ReadFile(file); !bytes.Equal(before, after) {
		t.Error("the checked-in CHECKSUMS was rewritten")
	}
	// The other half of the pair: CATALOG_DIR without the write mode would
	// compare some other tree and pass for a catalog that is not this one.
	cmd = exec.Command("go", "test", ".", "-run", "^TestCatalogMatchesChecksums$", "-count=1")
	cmd.Env = fixtureEnv("CATALOG_WRITE=", "CATALOG_DIR="+t.TempDir())
	out, err = cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "only for the refresh script's write mode") {
		t.Errorf("CATALOG_DIR without CATALOG_WRITE=1 was accepted (err %v):\n%s", err, out)
	}
}

// TestCatalogNotesNeverAllowQuiet: the router does not follow an entry's
// `--quiet` (step 4), and a note never relaxes the router. A note that names
// the flag therefore says to leave it out, in one of three fixed phrases;
// "only on read-only commands" or "after approval" is an exception the router
// does not have (security reviews of the v2.1.0 gate: found three times as
// instances, pinned here as a kind).
func TestCatalogNotesNeverAllowQuiet(t *testing.T) {
	// The helper first, on lines that must and must not pass: a carve-out
	// beside an allowed phrase once passed (code review, R3-1).
	for line, bad := range map[string]bool{
		"- `x` — No `--quiet` on any command.":                                      false,
		"- `x` — bind on the bucket, drop `--quiet`, and confirm every IAM change.": false,
		"- `x` — run it without `--quiet`; ask before `services enable`.":           false,
		"- `x` — the operator only said so; drop `--quiet`; ask after the dry run.": false,
		"- `x` — No `--quiet` except after approval.":                               true,
		"- `x` — drop `--quiet` unless the command is read-only.":                   true,
		"- `x` — `--quiet` only on read-only commands.":                             true,
		"- `x` — No `--quiet` once the operator approved the plan.":                 true,
		"- `x` — it runs `gcloud compute ssh --quiet`.":                             true,
		"- `x` — drop `--quiet` if the command changes state.":                      true,
		"- `x` — No `--quiet` when the operator has not approved.":                  true,
		"- `x` — drop `--quiet`, but keep it on read commands.":                     true,
		"- `x` — pass `--project` on each command.":                                 false,
	} {
		if got := quietNoteProblem(line); (got != "") != bad {
			t.Errorf("quietNoteProblem(%q) = %q, want a problem: %v", line, got, bad)
		}
	}
	if os.Getenv("CATALOG_WRITE") == "1" {
		t.Skip("the refresh script's build directory has no notes yet")
	}
	data, err := os.ReadFile(filepath.Join(catalogDir(t), "MANIFEST.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, notes, ok := strings.Cut(string(data), "\n## Notes\n")
	if !ok {
		t.Fatal("MANIFEST.md has no ## Notes section")
	}
	seen := 0
	for _, line := range strings.Split(notes, "\n") {
		if !strings.HasPrefix(line, "- `") || !strings.Contains(line, "--quiet") {
			continue
		}
		seen++
		if problem := quietNoteProblem(line); problem != "" {
			t.Errorf("%s:\n%s", problem, line)
		}
	}
	if seen == 0 {
		t.Error("no note names --quiet: the test checks nothing (was the section renamed?)")
	}
}

// fixtureEnv is the environment for every subprocess the catalog tests
// start, the fixture's git calls and each run of the script: the caller's,
// without any GIT_* variable, plus the test's own additions. An exported
// GIT_DIR or GIT_INDEX_FILE (git sets them inside a hook) would aim the
// fixture's add, commit and reset --hard at another repository (the v2.1.0
// gate's L10).
func fixtureEnv(extra ...string) []string {
	var env []string
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") {
			env = append(env, kv)
		}
	}
	return append(env, extra...)
}

// quietNoteProblem judges one note: every clause that names `--quiet` must
// say to leave it out, in one of three fixed phrases, and must not carve an
// exception beside the phrase. A clause ends at a semicolon or at a full
// stop followed by a space. It returns "" or the reason.
func quietNoteProblem(line string) string {
	for _, clause := range quietClauseEnd.Split(line, -1) {
		if !strings.Contains(clause, "--quiet") {
			continue
		}
		lower := strings.ToLower(clause)
		if !strings.Contains(lower, "no `--quiet`") && !strings.Contains(lower, "drop `--quiet`") && !strings.Contains(lower, "without `--quiet`") {
			return "a note names --quiet without saying to leave it out (no `--quiet`, drop `--quiet` or without `--quiet`)"
		}
		if word := quietCarveOut.FindString(lower); word != "" {
			return "a note carves an exception (\"" + word + "\") in the clause that names --quiet; the router has none"
		}
	}
	return ""
}

var (
	quietClauseEnd = regexp.MustCompile(`;|\. `)
	quietCarveOut  = regexp.MustCompile(`\b(except|unless|only|once|after|if|when|but)\b`)
)

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
		cmd.Env = fixtureEnv("GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.invalid")
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
		cmd.Env = fixtureEnv("CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
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
	relCmd.Env = fixtureEnv("CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
	if out, err := relCmd.CombinedOutput(); err != nil {
		t.Errorf("relative --dest and --work: %v\n%s", err, out)
	} else if _, err := os.Stat(filepath.Join(tmp, "rel-dest", "gamma", "SKILL.md")); err != nil {
		t.Errorf("relative --dest: the catalog did not land under the caller's directory: %v", err)
	}

	// An entry the old CHECKSUMS names but whose directory is already gone
	// (a half-done refresh, a hand deletion) must not stop the swap half way:
	// the loop's last test once ended the pipeline with status 1, silently,
	// after the other entries were deleted (code review 2026-10-01).
	if err := os.RemoveAll(filepath.Join(dest, "gamma")); err != nil {
		t.Fatal(err)
	}
	if out, code := run(pinA, pinG); code != 0 || snapshot() != good {
		t.Fatalf("a refresh over a catalog with a missing entry directory did not restore it (exit %d)\n%s", code, out)
	}

	// The swap removes what the old CHECKSUMS names and nothing else: a
	// directory someone put beside the entries survives a refresh.
	put(dest, map[string]string{"bystander/keep.txt": "kept\n"})
	if out, code := run(pinA, pinG); code != 0 {
		t.Fatalf("a refresh over a catalog with a bystander directory failed (exit %d)\n%s", code, out)
	}
	if _, err := os.Stat(filepath.Join(dest, "bystander", "keep.txt")); err != nil {
		t.Errorf("the swap removed a directory the old CHECKSUMS does not name: %v", err)
	}
	if err := os.RemoveAll(filepath.Join(dest, "bystander")); err != nil {
		t.Fatal(err)
	}
	if snapshot() != good {
		t.Fatal("the catalog is not the good one after the bystander fixture")
	}

	// The same loop's second way out: a CHECKSUMS that names no entry at all
	// (an empty one) once made `grep /` end the pipeline the same way.
	bare := filepath.Join(tmp, "bare-dest")
	put(bare, map[string]string{"CHECKSUMS": ""})
	bareCmd := exec.Command("bash", script, "--dest", bare, "--work", work, pinA, pinG)
	bareCmd.Dir = repo
	bareCmd.Env = fixtureEnv("CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
	if out, err := bareCmd.CombinedOutput(); err != nil {
		t.Errorf("a refresh over a destination whose CHECKSUMS names no entry failed: %v\n%s", err, out)
	} else if _, err := os.Stat(filepath.Join(bare, "gamma", "SKILL.md")); err != nil {
		t.Errorf("a refresh over a destination whose CHECKSUMS names no entry vendored nothing: %v", err)
	}

	// The destination guard: an empty --dest, and a directory that is not a
	// catalog, are refused before anything is deleted (security review S1).
	guardDir := filepath.Join(tmp, "not-a-catalog")
	if err := os.MkdirAll(filepath.Join(guardDir, "precious"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"", guardDir} {
		cmd := exec.Command("bash", script, "--dest", d, "--work", work, pinA, pinG)
		cmd.Dir = guardDir
		cmd.Env = fixtureEnv("CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("--dest %q was accepted:\n%s", d, out)
		}
		if _, err := os.Stat(filepath.Join(guardDir, "precious")); err != nil {
			t.Fatalf("--dest %q: a directory that is not a catalog entry was deleted", d)
		}
	}
	// The places the guard names, each refused by its own message before the
	// not-a-catalog check is reached: /, the home directory, the checkout,
	// and an ancestor of either. HOME is a fixture directory here.
	fakeHome := filepath.Join(tmp, "home", "user")
	put(fakeHome, map[string]string{"keep.txt": "kept\n"})
	repoLink := filepath.Join(tmp, "link-to-the-checkout")
	if err := os.Symlink(repo, repoLink); err != nil {
		repoLink = repo // no symlinks here: the case repeats the plain one
	}
	for _, c := range []struct{ dest, want string }{
		{"/", "is not a catalog directory"},
		{fakeHome, "is not a catalog directory"},
		{fakeHome + "/", "is not a catalog directory"},
		{repo, "is not a catalog directory"},
		{filepath.Dir(repo), "contains the checkout"},
		{filepath.Dir(fakeHome), "contains the home directory"},
		// Another spelling of the same place is the same place: a path
		// through `..` and a symlink (security review of the v2.1.0 gate).
		{filepath.Join(repo, "skills") + "/..", "is not a catalog directory"},
		{repoLink, "is not a catalog directory"},
		// `..` below a directory that does not exist cannot be resolved, and
		// `mkdir -p` would later make it real: the home directory, the
		// checkout or / reached that way (code review of 7dd0645..f62513b).
		{fakeHome + "/no-such-dir/..", "below a directory that does not exist"},
		{repo + "/no-such-dir/..", "below a directory that does not exist"},
		{"/no-such-dir-for-this-test/..", "below a directory that does not exist"},
		{fakeHome + "/no-such-dir/./catalog", "below a directory that does not exist"},
		// A leading run of slashes is the same place (the v2.1.0 gate's L2).
		{"/" + fakeHome, "is not a catalog directory"},
		{"//" + fakeHome, "is not a catalog directory"},
	} {
		cmd := exec.Command("bash", script, "--dest", c.dest, "--work", work, pinA, pinG)
		cmd.Dir = tmp
		// The sources name nothing: if the guard and the not-a-catalog check
		// both regressed, the run stops at the clone and writes nowhere
		// (code review round 3).
		nowhere := filepath.Join(tmp, "no-such-repository")
		cmd.Env = fixtureEnv("HOME="+fakeHome, "CATALOG_AGENTS_CLI_REPO="+nowhere, "CATALOG_GOOGLE_SKILLS_REPO="+nowhere, "CATALOG_DIR=", "CATALOG_WRITE=")
		out, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(out), c.want) {
			t.Errorf("--dest %s: want a refusal saying %q, got err %v\n%s", c.dest, c.want, err, out)
		}
	}
	if _, err := os.Stat(filepath.Join(fakeHome, "keep.txt")); err != nil {
		t.Fatalf("the guard fixture's home directory was touched: %v", err)
	}
	// Two ways the guard itself could lose its footing (code review of
	// 7dd0645..b4e3d79, round 2): a parent directory it cannot enter, which
	// once resolved to a path under /, and no HOME to compare against.
	guardRun := func(dest string, env []string) (string, error) {
		cmd := exec.Command("bash", script, "--dest", dest, "--work", work, pinA, pinG)
		cmd.Dir = tmp
		nowhere := filepath.Join(tmp, "no-such-repository")
		cmd.Env = append(env, "CATALOG_AGENTS_CLI_REPO="+nowhere, "CATALOG_GOOGLE_SKILLS_REPO="+nowhere, "CATALOG_DIR=", "CATALOG_WRITE=")
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	if os.Geteuid() != 0 {
		locked := filepath.Join(tmp, "locked")
		if err := os.Mkdir(locked, 0o000); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { os.Chmod(locked, 0o755) })
		if out, err := guardRun(filepath.Join(locked, "x"), fixtureEnv("HOME="+fakeHome)); err == nil || !strings.Contains(out, "cannot resolve") {
			t.Errorf("--dest below a directory that cannot be entered: want a refusal saying \"cannot resolve\", got err %v\n%s", err, out)
		}
	}
	// --work is cloned into and its build directory removed, so it is a
	// directory of this script or a new one (the v2.1.0 gate's L7): never /,
	// the home directory or the checkout, never nested with --dest, and when
	// it holds files, only with the marker the script writes before its
	// first clone, or the origins.txt of an earlier completed run.
	workRun := func(destArg, workArg string, pins ...string) (string, error) {
		cmd := exec.Command("bash", append([]string{script, "--dest", destArg, "--work", workArg}, pins...)...)
		cmd.Dir = tmp
		cmd.Env = fixtureEnv("HOME="+fakeHome, "CATALOG_AGENTS_CLI_REPO="+agents, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	// The three forbidden places are tried with a COPY of the script inside a
	// stand-in checkout, and with sources that cannot be cloned: when this
	// guard did not exist yet, the same cases run against the real checkout
	// with working sources cloned the fixtures into it (BUG-150).
	fakeRepo := filepath.Join(tmp, "fake-checkout")
	scriptCopy := filepath.Join(fakeRepo, "skills", "google-catalog", "scripts", "refresh-catalog.sh")
	scriptBytes, err := os.ReadFile(script)
	if err != nil {
		t.Fatal(err)
	}
	put(fakeRepo, map[string]string{"skills/google-catalog/scripts/refresh-catalog.sh": string(scriptBytes), "README.md": "a stand-in checkout\n"})
	places := []string{fakeHome, fakeRepo}
	if os.Geteuid() != 0 {
		places = append(places, "/") // as root a regressed guard could write a marker there
	}
	for _, w := range places {
		cmd := exec.Command("bash", scriptCopy, "--dest", filepath.Join(tmp, "dest-w"), "--work", w, pinA, pinG)
		cmd.Dir = tmp
		nowhere := filepath.Join(tmp, "no-such-repository")
		cmd.Env = fixtureEnv("HOME="+fakeHome, "CATALOG_AGENTS_CLI_REPO="+nowhere, "CATALOG_GOOGLE_SKILLS_REPO="+nowhere, "CATALOG_DIR=", "CATALOG_WRITE=")
		if out, err := cmd.CombinedOutput(); err == nil || !strings.Contains(string(out), "must not be /, the home directory or the checkout") {
			t.Errorf("--work %s: want a refusal, got err %v\n%s", w, err, out)
		}
	}
	for _, stray := range []string{filepath.Join(fakeRepo, "agents-cli"), filepath.Join(fakeRepo, ".refresh-catalog-work"), filepath.Join(fakeHome, "agents-cli"), filepath.Join(fakeHome, ".refresh-catalog-work")} {
		if _, err := os.Stat(stray); err == nil {
			t.Errorf("a refused --work was written into: %s", stray)
		}
	}
	nest := filepath.Join(tmp, "nest")
	for _, c := range [][2]string{{nest, filepath.Join(nest, "work")}, {filepath.Join(nest, "dest"), nest}} {
		if out, err := workRun(c[0], c[1], pinA, pinG); err == nil || !strings.Contains(out, "must not contain each other") {
			t.Errorf("--dest %s --work %s: want a refusal, got err %v\n%s", c[0], c[1], err, out)
		}
	}
	foreign := filepath.Join(tmp, "someone-elses-directory")
	put(foreign, map[string]string{"agents-cli/notes.txt": "the operator's own clone\n", "build/keep.txt": "kept\n"})
	if out, err := workRun(filepath.Join(tmp, "dest-f"), foreign, pinA, pinG); err == nil || !strings.Contains(out, "is not a work directory of this script") || !strings.Contains(out, "name an empty or new directory, or delete this one") {
		t.Errorf("a non-empty --work without the marker: want a refusal, got err %v\n%s", err, out)
	}
	for _, f := range []string{"agents-cli/notes.txt", "build/keep.txt"} {
		if _, err := os.Stat(filepath.Join(foreign, f)); err != nil {
			t.Errorf("the refused --work lost %s: %v", f, err)
		}
	}
	// A first run that fails after the checks (a pin that does not resolve)
	// must leave a work directory the corrected run accepts.
	freshWork, freshDest := filepath.Join(tmp, "fresh-work"), filepath.Join(tmp, "fresh-dest")
	if out, err := workRun(freshDest, freshWork, pinA, "0000000000000000000000000000000000000000"); err == nil || !strings.Contains(out, "does not resolve") {
		t.Fatalf("a first run with a bad pin: want the pin refusal, got err %v\n%s", err, out)
	}
	if out, err := workRun(freshDest, freshWork, pinA, pinG); err != nil {
		t.Errorf("the corrected run on the same --work was refused: %v\n%s", err, out)
	} else if _, err := os.Stat(filepath.Join(freshDest, "gamma", "SKILL.md")); err != nil {
		t.Errorf("the corrected run vendored nothing: %v", err)
	}
	// The source overrides take an https URL or an absolute path, nothing
	// else: git has transports that run a command (`ext::`).
	for _, bad := range []string{"ext::sh -c id", "relative/path", "http://example.invalid/x", "file:///x", "https://user:token@example.invalid/x", "http://user:token@example.invalid/x", "https://example.invalid/x?access_token=abc"} {
		cmd := exec.Command("bash", script, "--dest", filepath.Join(tmp, "dest-o"), "--work", filepath.Join(tmp, "work-o"), pinA, pinG)
		cmd.Dir = tmp
		cmd.Env = fixtureEnv("HOME="+fakeHome, "CATALOG_AGENTS_CLI_REPO="+bad, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
		out, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(out), "CATALOG_AGENTS_CLI_REPO must be an https:// URL or an absolute path") {
			t.Errorf("override %q: want a refusal naming the variable, got err %v\n%s", bad, err, out)
		}
		// No refused value is echoed: any of them may hold a credential
		// (an http:// URL with user information is refused for its scheme).
		if strings.Contains(string(out), bad) || strings.Contains(string(out), "token") || strings.Contains(string(out), "example.invalid") {
			t.Errorf("override %q: the refusal echoes the value:\n%s", bad, out)
		}
		if strings.HasPrefix(bad, "https://") && !strings.Contains(string(out), "without credentials") {
			t.Errorf("override with credentials: the refusal does not say why:\n%s", out)
		}
	}
	if _, err := os.Stat(filepath.Join(tmp, "work-o")); err == nil {
		t.Error("a refused override still created the work directory")
	}
	// A source that passes the check and then cannot be cloned is named by
	// its role, not by its value: a token can sit in a URL's path.
	{
		cmd := exec.Command("bash", script, "--dest", filepath.Join(tmp, "dest-c"), "--work", filepath.Join(tmp, "work-c"), pinA, pinG)
		cmd.Dir = tmp
		secretPath := filepath.Join(tmp, "no-such-SECRETPATH-repository")
		cmd.Env = fixtureEnv("HOME="+fakeHome, "CATALOG_AGENTS_CLI_REPO="+secretPath, "CATALOG_GOOGLE_SKILLS_REPO="+gskills, "CATALOG_DIR=", "CATALOG_WRITE=")
		out, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(out), "cannot clone the agents-cli source") || strings.Contains(string(out), "SECRETPATH") {
			t.Errorf("a source that cannot be cloned: want it named by role and its value kept out of the output, git's own message included, got err %v\n%s", err, out)
		}
	}
	var noHome []string
	for _, kv := range fixtureEnv() {
		if !strings.HasPrefix(kv, "HOME=") {
			noHome = append(noHome, kv)
		}
	}
	if out, err := guardRun(filepath.Join(tmp, "dest-without-home"), noHome); err == nil || !strings.Contains(out, "HOME is not set") {
		t.Errorf("no HOME: want a refusal saying \"HOME is not set\", got err %v\n%s", err, out)
	}
	if out, code := run(pinA, "main"); code == 0 || !strings.Contains(out, "commit id") {
		t.Errorf("a pin that is not a commit id was accepted (exit %d)\n%s", code, out)
	}
	// An entry named under the manifest's ## Excluded stays out without the flag (S5).
	mf := filepath.Join(dest, "MANIFEST.md")
	orig, _ := os.ReadFile(mf)
	if err := os.WriteFile(mf, bytes.Replace(orig, []byte("## Excluded\n"), []byte("## Excluded\n\n- `epsilon` — excluded by the fixture\n"), 1), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, code := run(pinA, pinG); code != 0 {
		t.Errorf("a refresh with a manifest exclusion failed (exit %d)\n%s", code, out)
	} else if _, err := os.Stat(filepath.Join(dest, "epsilon")); err == nil {
		t.Error("an entry listed under ## Excluded was vendored")
	}
	os.WriteFile(mf, orig, 0o644)
	if out, code := run(pinA, pinG); code != 0 || snapshot() != good {
		t.Fatalf("could not restore the good catalog after the exclusion fixture (exit %d)\n%s", code, out)
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
	// The hint is printed shell-quoted: the name comes from upstream and the
	// operator pastes the hint into a shell (code review round 3).
	refuse("name with a space", `(--exclude bad\ name)`, map[string]string{"skills/cloud/bad name/SKILL.md": skill("bad-name", "description: a directory name with a space\n")}, nil, "bad name")
	refuse("uppercase name", `(--exclude Loud)`, map[string]string{"skills/cloud/Loud/SKILL.md": skill("loud", "description: a directory name with a capital\n")}, nil, "Loud")
	refuse("shell characters in a name", `(--exclude a\$\(id\))`, map[string]string{"skills/cloud/a$(id)/SKILL.md": skill("a-id", "description: a directory name a shell would expand\n")}, nil, "a$(id)")
	// The refusal prints an upstream name only shell-quoted: a raw escape
	// sequence in a directory name must not reach the operator's terminal
	// (security review of the v2.1.0 gate).
	put(gskills, map[string]string{"skills/cloud/a\x1b[31mb/SKILL.md": skill("a-b", "description: a directory name with a terminal escape\n")})
	escPin := commit(gskills)
	if out, code := run(pinA, escPin); code != 1 || strings.Contains(out, "\x1b") || !strings.Contains(out, "not a plain skill name") {
		t.Errorf("a name with an escape sequence: want exit 1 and no raw escape in the message, got exit %d\n%q", code, out)
	}
	git(gskills, "reset", "-q", "--hard", pinG)
	git(gskills, "clean", "-qfd")
	// The same for a path inside an entry: the dotfile, symlink and Go-file
	// refusals print the path shell-quoted. One fixture per message: each is
	// its own line in the script (the gate's verifier removed the quoting
	// from one and nothing failed).
	for _, c := range []struct {
		label, want string
		make        func()
	}{
		{"dotfile", "holds a dotfile", func() {
			put(gskills, map[string]string{"skills/cloud/delta/.x\x1b[31m": "a dotfile whose name holds an escape\n"})
		}},
		{"Go file", "holds a Go file", func() {
			put(gskills, map[string]string{"skills/cloud/delta/a\x1b[31m.go": "package main\n"})
		}},
		{"symlink", "holds a symlink", func() {
			if err := os.Symlink("SKILL.md", filepath.Join(gskills, "skills/cloud/delta/l\x1b[31mnk")); err != nil {
				t.Skip("no symlinks here")
			}
		}},
	} {
		c.make()
		pin := commit(gskills)
		if out, code := run(pinA, pin); code != 1 || strings.Contains(out, "\x1b") || !strings.Contains(out, c.want) {
			t.Errorf("a nested %s path with an escape sequence: want exit 1, %q and no raw escape in the message, got exit %d\n%q", c.label, c.want, code, out)
		}
		git(gskills, "reset", "-q", "--hard", pinG)
		git(gskills, "clean", "-qfd")
	}
	// An excluded name is one name: excluding "delta epsilon" must not also
	// exclude the entries delta and epsilon.
	put(gskills, map[string]string{"skills/cloud/delta epsilon/SKILL.md": skill("delta-epsilon", "description: two entry names with a space between\n")})
	wordsPin := commit(gskills)
	if out, code := run("--exclude", "delta epsilon", pinA, wordsPin); code != 0 {
		t.Errorf("--exclude \"delta epsilon\" failed (exit %d)\n%s", code, out)
	}
	for _, name := range []string{"delta", "epsilon"} {
		if _, err := os.Stat(filepath.Join(dest, name, "SKILL.md")); err != nil {
			t.Errorf("excluding \"delta epsilon\" also excluded the entry %s", name)
		}
	}
	git(gskills, "reset", "-q", "--hard", pinG)
	git(gskills, "clean", "-qfd")
	if out, code := run(pinA, pinG); code != 0 || snapshot() != good {
		t.Fatalf("could not restore the good catalog after the two-word exclusion (exit %d)\n%s", code, out)
	}
	refuse("dotfile", "delta", map[string]string{"skills/cloud/delta/.gitattributes": "* export-ignore\n"}, nil, "")
	refuse("symlinked entry", "zeta", nil, func() {
		if err := os.Symlink("gamma", filepath.Join(gskills, "skills/cloud/zeta")); err != nil {
			t.Skip("no symlinks here")
		}
	}, "")
	refuse("symlinked source set", "leaves the clone through a symlink", nil, func() {
		// skills/developers becomes a link to a directory elsewhere in the clone.
		if err := os.Rename(filepath.Join(gskills, "skills/developers"), filepath.Join(gskills, "elsewhere")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("../elsewhere", filepath.Join(gskills, "skills/developers")); err != nil {
			t.Skip("no symlinks here")
		}
	}, "")
	refuse("symlinked LICENSE", "LICENSE", nil, func() {
		os.Remove(filepath.Join(gskills, "LICENSE"))
		if err := os.Symlink("skills/cloud/gamma/SKILL.md", filepath.Join(gskills, "LICENSE")); err != nil {
			t.Skip("no symlinks here")
		}
	}, "")
}
