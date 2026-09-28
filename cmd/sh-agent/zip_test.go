package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

// zipOfTree writes tree's files into a zip under one top folder, the way
// GitHub's "Download ZIP" lays them out, and returns the zip's path.
func zipOfTree(t *testing.T, tree string) string {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	err := filepath.WalkDir(tree, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(tree, p)
		w, err := zw.Create("styrman-harness-main/" + filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		_, err = w.Write(b)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	zp := filepath.Join(t.TempDir(), "styrman-harness-main.zip")
	if err := os.WriteFile(zp, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return zp
}

// markedZip is this checkout's harness with a marker line in its context,
// so an install can be told apart from the embedded files.
func markedZip(t *testing.T) string {
	t.Helper()
	tree := plainCopy(t)
	ctx := filepath.Join(tree, "context", "AGENTS.md")
	b, err := os.ReadFile(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ctx, append([]byte("<!-- ZIPMARKER -->\n"), b...), 0o644); err != nil {
		t.Fatal(err)
	}
	return zipOfTree(t, tree)
}

// isolate cuts a test's reach: no git or gh on PATH and a scratch working
// directory, so a regression that falls through to a real update or report
// cannot fetch, pull or push (BUG-094).
func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("SH_AGENT_REPO", "")
	t.Setenv("PATH", t.TempDir())
	t.Chdir(t.TempDir())
}

func noZipDirsUnder(t *testing.T, home string) {
	t.Helper()
	left, _ := filepath.Glob(filepath.Join(home, ".config", "sh-agent", "zip-*"))
	if len(left) != 0 {
		t.Errorf("unpack directories left behind: %v", left)
	}
}

func TestUpdateFromZip(t *testing.T) {
	isolate(t)
	defer func(v string) { version = v }(version)
	version = "v1.2.3" // a release binary: its source would otherwise be the embedded files
	zp := markedZip(t)
	home := t.TempDir()

	var out bytes.Buffer
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &out); err != nil {
		t.Fatalf("update --from-zip: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "UNVERIFIED") || !strings.Contains(out.String(), "nothing in the zip replaces") || !strings.Contains(out.String(), "sha256") {
		t.Errorf("output does not say what was installed and that the binary stays:\n%s", out.String())
	}
	ctx, err := os.ReadFile(filepath.Join(home, harness.ClaudeRoot, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ctx), "ZIPMARKER") {
		t.Errorf("installed context is not the zip's (no ZIPMARKER): the release binary's embedded files won")
	}
	st, found, err := harness.LoadState(home)
	if err != nil || !found {
		t.Fatalf("manifest: found=%v err=%v", found, err)
	}
	if !strings.HasPrefix(st.Repo, harness.ZipPrefix+"sha256:") || st.Commit != "" {
		t.Errorf("manifest repo=%q commit=%q, want zip:sha256:… and no commit", st.Repo, st.Commit)
	}
	noZipDirsUnder(t, home)
}

func TestUpdateFromZipRefusals(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	for _, args := range [][]string{
		{"update", "--from-zip", zp, "--check"},
		{"update", "--from-zip", zp, "--repo", "../.."},
		{"install", "--from-zip", zp},
	} {
		home := t.TempDir()
		var out bytes.Buffer
		err := run(append(args, "--home", home), &out)
		if err == nil || !strings.Contains(err.Error(), "--from-zip") {
			t.Errorf("%v: err = %v, want a refusal naming --from-zip", args, err)
		}
		if entries, _ := os.ReadDir(home); len(entries) != 0 {
			t.Errorf("%v: refused, yet wrote %d entries under the home", args, len(entries))
		}
	}
}

func TestUpdateFromZipDryRunWritesNothing(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	home := t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home, "--dry-run"}, &out); err != nil {
		t.Fatalf("dry run: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "[dry-run] zip") {
		t.Errorf("dry run does not say it is one:\n%s", out.String())
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		names := []string{}
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("dry run on a fresh home wrote %v", names)
	}
}

func TestStatusAfterZip(t *testing.T) {
	isolate(t)
	defer func(v string) { version = v }(version)
	version = "v1.2.3"
	zp := markedZip(t)
	home := t.TempDir()
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"status", "--home", home}, &out); err != nil {
		t.Fatalf("status after a zip install: %v\n%s", err, out.String())
	}
	for _, want := range []string{"from zip sha256:", "(unverified)", "installed from a zip", "not compared", "manifest files present"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("status lacks %q:\n%s", want, out.String())
		}
	}
	if strings.Contains(out.String(), "differ") {
		t.Errorf("status compared a zip install against the embedded files:\n%s", out.String())
	}

	st, _, _ := harness.LoadState(home)
	gone := st.Files[0]
	if err := os.Remove(filepath.Join(home, gone)); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	err := run([]string{"status", "--home", home}, &out)
	var code exitCode
	if !errors.As(err, &code) || code != 1 {
		t.Errorf("status with a missing manifest file: err = %v, want exit 1", err)
	}
	if !strings.Contains(out.String(), "missing  ~/"+gone) {
		t.Errorf("status does not list the missing %s:\n%s", gone, out.String())
	}
}

// A manifest written before the verdict was recorded ("zip:sha256:<hex>", no
// ";…" part) still reads: status names the zip and claims no verdict.
func TestStatusLegacyZipManifest(t *testing.T) {
	isolate(t)
	defer func(v string) { version = v }(version)
	version = "v1.2.3"
	home := t.TempDir()
	if err := harness.SaveState(home, harness.State{Repo: harness.ZipPrefix + "sha256:0123456789abcdef0123"}); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	_ = run([]string{"status", "--home", home}, &out)
	if !strings.Contains(out.String(), "from zip sha256:0123456789ab ") {
		t.Errorf("status does not name the legacy zip install:\n%s", out.String())
	}
	if strings.Contains(out.String(), "verified") {
		t.Errorf("status claims a verdict the legacy manifest never recorded:\n%s", out.String())
	}
}

func TestKnownCheckoutZip(t *testing.T) {
	isolate(t)
	home := t.TempDir()
	if err := harness.SaveState(home, harness.State{Repo: harness.ZipPrefix + "sha256:0123456789abcdef0123"}); err != nil {
		t.Fatal(err)
	}
	err := knownCheckout("", home, "/some/checkout")
	if err == nil || !strings.Contains(err.Error(), "installed from a zip (sha256 0123456789ab)") || !strings.Contains(err.Error(), "--repo /some/checkout") {
		t.Errorf("err = %v, want the zip install named and the --repo way out", err)
	}
}

// A refused zip writes nothing under the home, not even ~/.config/sh-agent.
func TestUpdateFromZipRefusedZipWritesNothing(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	zr, err := zip.OpenReader(zp)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range zr.File {
		w, err := zw.Create(f.Name)
		if err != nil {
			t.Fatal(err)
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		b := new(bytes.Buffer)
		b.ReadFrom(rc)
		rc.Close()
		w.Write(b.Bytes())
	}
	zr.Close()
	w, err := zw.Create("styrman-harness-main/x/../../escape")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte("x"))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	slip := filepath.Join(t.TempDir(), "slip.zip")
	if err := os.WriteFile(slip, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	err = run([]string{"update", "--from-zip", slip, "--unverified", "--home", home}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "would escape") {
		t.Fatalf("err = %v, want the escaping entry refused", err)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("a refused zip left %d entries under the home", len(entries))
	}
}

// --from-zip belongs to update: every other command refuses it before doing
// anything — uninstall must not uninstall, report must not report.
func TestFromZipRefusedOnOtherCommands(t *testing.T) {
	// No git and no gh on PATH: if the refusal ever regresses, report cannot
	// reach the network. A red run of this test before the refusal existed
	// pushed a real report (BUG-094).
	isolate(t)
	for _, cmd := range []string{"install", "uninstall", "report", "status", "diff", "list", "render", "version"} {
		home := t.TempDir()
		var out bytes.Buffer
		err := run([]string{cmd, "--from-zip", "x.zip", "--home", home}, &out)
		if err == nil || !strings.Contains(err.Error(), "--from-zip belongs to update") {
			t.Errorf("%s --from-zip: err = %v, want the refusal", cmd, err)
		}
		if entries, _ := os.ReadDir(home); len(entries) != 0 {
			t.Errorf("%s --from-zip wrote %d entries before refusing", cmd, len(entries))
		}
		if out.Len() != 0 {
			t.Errorf("%s --from-zip printed before refusing:\n%s", cmd, out.String())
		}
	}
}

// An empty value is refused, never read as "no --from-zip": that would fall
// through to a normal update, which on a release binary downloads and
// replaces the binary.
func TestFromZipEmptyValueRefused(t *testing.T) {
	isolate(t)
	for _, args := range [][]string{{"update", "--from-zip="}, {"update", "--from-zip", ""}} {
		home := t.TempDir()
		err := run(append(args, "--home", home), &bytes.Buffer{})
		if err == nil || !strings.Contains(err.Error(), "--from-zip needs a file") {
			t.Errorf("%q: err = %v, want the empty-value refusal", args, err)
		}
		if entries, _ := os.ReadDir(home); len(entries) != 0 {
			t.Errorf("%q wrote %d entries", args, len(entries))
		}
	}
}

func TestDiffAfterZip(t *testing.T) {
	isolate(t)
	defer func(v string) { version = v }(version)
	version = "v1.2.3"
	zp := markedZip(t)
	home := t.TempDir()
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"diff", "--home", home}, &out); err != nil {
		t.Fatalf("diff after a zip install: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "not compared") || strings.Contains(out.String(), "ZIPMARKER") {
		t.Errorf("diff compared a zip install against the embedded files:\n%s", out.String())
	}
}

// A zip-* directory an interrupted run left behind is swept once it is an
// hour old; a fresh one, which may belong to a run in progress, is kept; a
// dry run sweeps nothing.
func TestUpdateFromZipSweepsStaleDirs(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	home := t.TempDir()
	cfg := filepath.Join(home, ".config", "sh-agent")
	old, fresh := filepath.Join(cfg, "zip-old"), filepath.Join(cfg, "zip-fresh")
	// An old directory that is not an unpack directory (the feedback clone
	// older versions kept, say) is never the sweep's to remove.
	notOurs := filepath.Join(cfg, "feedback-old")
	for _, d := range []string{old, fresh, notOurs} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	past := time.Now().Add(-2 * time.Hour)
	for _, d := range []string{old, notOurs} {
		if err := os.Chtimes(d, past, past); err != nil {
			t.Fatal(err)
		}
	}
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home, "--dry-run"}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); err != nil {
		t.Errorf("a dry run swept %s", old)
	}
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); err == nil {
		t.Errorf("the stale %s was not swept", old)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Errorf("the fresh %s was swept: %v", fresh, err)
	}
	if _, err := os.Stat(notOurs); err != nil {
		t.Errorf("the sweep removed %s, which is not an unpack directory: %v", notOurs, err)
	}
}

// update --dry-run on a release binary changes nothing: it never downloads,
// never replaces the binary, never sheds the kept clone. It says what a real
// update would do.
func TestReleaseUpdateDryRunChangesNothing(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("no sh")
	}
	t.Setenv("SH_AGENT_REPO", "")
	t.Chdir(t.TempDir())
	defer func(v string) { version = v }(version)
	version = "v1.0.0"
	bin := t.TempDir()
	log := filepath.Join(t.TempDir(), "gh.log")
	script := "#!/bin/sh\n" + "echo \"$@\" >> " + log + "\n" + "printf 'v9.9.9\\n'\n"
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	home := t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"update", "--dry-run", "--home", home}, &out); err != nil {
		t.Fatalf("update --dry-run: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "[dry-run]") || !strings.Contains(out.String(), "v9.9.9") {
		t.Errorf("dry run does not say what update would do:\n%s", out.String())
	}
	calls, _ := os.ReadFile(log)
	if strings.Contains(string(calls), "download") {
		t.Errorf("update --dry-run downloaded:\n%s", calls)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("update --dry-run wrote %d entries under the home", len(entries))
	}
}

// A zip whose content the install refuses (here a skill with no description)
// leaves a fresh home as it was: the directory chain updateFromZip created
// is removed on any failure, not only an unpack refusal (BUG-093).
func TestUpdateFromZipRefusedContentWritesNothing(t *testing.T) {
	isolate(t)
	tree := plainCopy(t)
	bad := filepath.Join(tree, "skills", "zz-broken", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(bad), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bad, []byte("---\nname: zz-broken\n---\nno description\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	zp := zipOfTree(t, tree)
	home := t.TempDir()
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "has no description") {
		t.Fatalf("err = %v, want the install's own refusal of a skill with no description", err)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("a refused install left %d entries under a fresh home", len(entries))
	}
}

// The sweep matches names, not a glob built from the home: a home whose path
// holds glob metacharacters sweeps its own stale directory and nobody else's.
func TestUpdateFromZipSweepWithGlobCharsInHome(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	base := t.TempDir()
	home, other := filepath.Join(base, "h[1]"), filepath.Join(base, "h1")
	past := time.Now().Add(-2 * time.Hour)
	var mine, theirs string
	for _, h := range []string{home, other} {
		d := filepath.Join(h, ".config", "sh-agent", "zip-old")
		if err := os.MkdirAll(d, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(d, past, past); err != nil {
			t.Fatal(err)
		}
		if h == home {
			mine = d
		} else {
			theirs = d
		}
	}
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(mine); err == nil {
		t.Errorf("the home's own stale %s was not swept", mine)
	}
	if _, err := os.Stat(theirs); err != nil {
		t.Errorf("another home's %s was swept", theirs)
	}
}

// A dry run on a fresh home unpacks under the system temp directory and
// leaves nothing there either.
func TestUpdateFromZipDryRunCleansTempDir(t *testing.T) {
	isolate(t)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	zp := markedZip(t)
	before, _ := os.ReadDir(tmp)
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", t.TempDir(), "--dry-run"}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadDir(tmp)
	if len(after) != len(before) {
		t.Errorf("the dry run left %d new entries in TMPDIR", len(after)-len(before))
	}
}

// update --dry-run on a checkout that is behind reports and changes nothing:
// no pull, no rebuild, no install.
func TestCheckoutUpdateDryRunChangesNothing(t *testing.T) {
	t.Setenv("SH_AGENT_REPO", "")
	repo, marker := gitFixture(t)
	base := filepath.Dir(repo)
	other := filepath.Join(base, "other")
	skill := filepath.Join(repo, "skills", "alpha", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skill), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skill, []byte("---\nname: alpha\ndescription: Alpha.\n---\nbody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"git", "-C", repo, "add", "-A"},
		{"git", "-C", repo, "commit", "--quiet", "-m", "harness-shaped"},
		{"git", "-C", repo, "push", "--quiet"},
		{"git", "clone", "--quiet", filepath.Join(base, "origin.git"), other},
		{"git", "-C", other, "-c", "user.email=f@example.invalid", "-c", "user.name=f", "commit", "--quiet", "--allow-empty", "-m", "ahead"},
		{"git", "-C", other, "push", "--quiet", "origin", "HEAD:main"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = append(gitCleanEnv(), "GIT_CONFIG_GLOBAL="+filepath.Join(base, "gitconfig"), "GIT_CONFIG_NOSYSTEM=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out)
		}
	}
	head := func() string {
		out, _ := exec.Command("git", "-C", repo, "rev-parse", "HEAD").Output()
		return strings.TrimSpace(string(out))
	}
	before := head()
	home := t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"update", "--dry-run", "--repo", repo, "--home", home}, &out); err != nil {
		t.Fatalf("update --dry-run: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "behind upstream by 1") || !strings.Contains(out.String(), "[dry-run]") {
		t.Errorf("the dry run does not say what update would do:\n%s", out.String())
	}
	if head() != before {
		t.Errorf("update --dry-run pulled: HEAD moved from %s", before)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Errorf("update --dry-run ran an install")
	}
}

// The sweep runs only after a successful install: a run refused at any point
// (the unpack here; the install below) changes nothing, not even a stale leftover.
func TestUpdateFromZipRefusedRunSweepsNothing(t *testing.T) {
	isolate(t)
	home := t.TempDir()
	stale := filepath.Join(home, ".config", "sh-agent", "zip-old")
	if err := os.MkdirAll(stale, 0o700); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-3 * time.Hour)
	if err := os.Chtimes(stale, past, past); err != nil {
		t.Fatal(err)
	}
	notZip := filepath.Join(t.TempDir(), "x.zip")
	if err := os.WriteFile(notZip, []byte("not a zip"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"update", "--from-zip", notZip, "--unverified", "--home", home}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "not a readable zip") {
		t.Fatalf("err = %v, want the not-a-zip refusal", err)
	}
	if _, err := os.Stat(stale); err != nil {
		t.Errorf("a refused run swept %s", stale)
	}
}

// A value that looks like a flag is refused: `update --from-zip --dry-run`
// must not read --dry-run as the file and run a real update.
func TestFromZipDashValueRefused(t *testing.T) {
	isolate(t)
	home := t.TempDir()
	err := run([]string{"update", "--from-zip", "--dry-run", "--home", home}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "--from-zip needs a file") {
		t.Errorf("err = %v, want the file-name refusal", err)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("wrote %d entries", len(entries))
	}
}

// A zip the INSTALL refuses (it unpacks fine) sweeps nothing either: the sweep
// sits after installFiles, not between the unpack and the install.
func TestUpdateFromZipInstallRefusalSweepsNothing(t *testing.T) {
	isolate(t)
	tree := plainCopy(t)
	bad := filepath.Join(tree, "skills", "zz-broken", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(bad), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bad, []byte("---\nname: zz-broken\n---\nno description\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	zp := zipOfTree(t, tree)
	home := t.TempDir()
	stale := filepath.Join(home, ".config", "sh-agent", "zip-old")
	if err := os.MkdirAll(stale, 0o700); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-3 * time.Hour)
	if err := os.Chtimes(stale, past, past); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", home}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "has no description") {
		t.Fatalf("err = %v, want the install's own refusal of a skill with no description", err)
	}
	if _, err := os.Stat(stale); err != nil {
		t.Errorf("an install refusal swept %s", stale)
	}
}

// A zip with signed checksums beside it that do not verify — signed by a key
// that is not the release signer — is refused, and the home stays as it was.
func TestUpdateFromZipForeignSignatureRefused(t *testing.T) {
	keygen := "/usr/bin/ssh-keygen"
	if _, err := os.Stat(keygen); err != nil {
		t.Skip("no /usr/bin/ssh-keygen")
	}
	isolate(t)
	zp := markedZip(t)
	dir := filepath.Dir(zp)
	data, err := os.ReadFile(zp)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	sums := "# styrman-harness v9.9.9\n" + hex.EncodeToString(sum[:]) + "  " + filepath.Base(zp) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "SHA256SUMS"), []byte(sums), 0o644); err != nil {
		t.Fatal(err)
	}
	key := filepath.Join(t.TempDir(), "key")
	for _, args := range [][]string{
		{"-q", "-t", "ed25519", "-N", "", "-C", "not-the-signer", "-f", key},
		{"-Y", "sign", "-q", "-f", key, "-n", "sh-agent-release", filepath.Join(dir, "SHA256SUMS")},
	} {
		if out, err := exec.Command(keygen, args...).CombinedOutput(); err != nil {
			t.Fatalf("ssh-keygen %v: %v\n%s", args, err, out)
		}
	}
	// --unverified allows a MISSING signature, never a failed one (R9, ADR-040).
	for _, extra := range [][]string{nil, {"--unverified"}} {
		home := t.TempDir()
		var out bytes.Buffer
		err = run(append([]string{"update", "--from-zip", zp, "--home", home}, extra...), &out)
		if err == nil || !strings.Contains(err.Error(), "does not verify") {
			t.Fatalf("%v: err = %v, want the signature refusal\n%s", extra, err, out.String())
		}
		if strings.Contains(strings.ToLower(err.Error()), "unverified") {
			t.Errorf("%v: a failed-check refusal names --unverified (ADR-040): %v", extra, err)
		}
		if entries, _ := os.ReadDir(home); len(entries) != 0 {
			t.Errorf("%v: a refused zip left %d entries", extra, len(entries))
		}
	}
}

// A zip with neither SHA256SUMS nor its .sig beside it is refused unless
// --unverified is passed — on a real run and a dry run alike, before anything
// is written (R9, ADR-040).
func TestUpdateFromZipUnsignedRefusedWithoutFlag(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	if _, err := os.Stat(filepath.Join(filepath.Dir(zp), "SHA256SUMS")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("setup: the zip's folder must hold no SHA256SUMS (stat err = %v)", err)
	}
	for _, extra := range [][]string{nil, {"--dry-run"}} {
		home := t.TempDir()
		var out bytes.Buffer
		err := run(append([]string{"update", "--from-zip", zp, "--home", home}, extra...), &out)
		if err == nil || !strings.Contains(err.Error(), "no signed SHA256SUMS") || !strings.Contains(err.Error(), "--unverified") {
			t.Fatalf("%v: err = %v, want the unsigned refusal naming --unverified\n%s", extra, err, out.String())
		}
		if !strings.Contains(err.Error(), "an agent never adds it on its own") {
			t.Errorf("%v: the refusal does not tell an agent the flag is the operator's call: %v", extra, err)
		}
		if entries, _ := os.ReadDir(home); len(entries) != 0 {
			t.Errorf("%v: a refused zip left %d entries", extra, len(entries))
		}
	}
}

// `--unverified false` is not a way to say no: the stray word is refused, and
// nothing is installed.
func TestFromZipStrayArgumentRefused(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	home := t.TempDir()
	err := run([]string{"update", "--from-zip", zp, "--unverified", "false", "--home", home}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "takes no other arguments") {
		t.Errorf("err = %v, want the stray argument refused", err)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("a refused run left %d entries", len(entries))
	}
}

// --unverified means nothing outside update --from-zip, and is refused there
// rather than ignored.
func TestUnverifiedOnlyWithFromZip(t *testing.T) {
	isolate(t)
	for _, args := range [][]string{{"update", "--unverified"}, {"install", "--unverified"}, {"status", "--unverified"}} {
		home := t.TempDir()
		err := run(append(args, "--home", home), &bytes.Buffer{})
		if err == nil || !strings.Contains(err.Error(), "--unverified belongs to update --from-zip") {
			t.Errorf("%v: err = %v, want the flag refused", args, err)
		}
		if entries, _ := os.ReadDir(home); len(entries) != 0 {
			t.Errorf("%v: a refused run left %d entries", args, len(entries))
		}
	}
}

// A zip with nothing beside it installs with --unverified, marked UNVERIFIED.
func TestUpdateFromZipUnsignedSaysUnverified(t *testing.T) {
	isolate(t)
	zp := markedZip(t)
	var out bytes.Buffer
	if err := run([]string{"update", "--from-zip", zp, "--unverified", "--home", t.TempDir()}, &out); err != nil {
		t.Fatalf("%v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "UNVERIFIED") || strings.Contains(out.String(), "VERIFIED against") {
		t.Errorf("an unsigned zip is not marked UNVERIFIED:\n%s", out.String())
	}
}

// A release zip with its signed SHA256SUMS beside it installs and says
// VERIFIED; the same zip with one byte changed is refused, home untouched.
func TestUpdateFromZipVerifiedAndTampered(t *testing.T) {
	keygen := "/usr/bin/ssh-keygen"
	if _, err := os.Stat(keygen); err != nil {
		t.Skip("no /usr/bin/ssh-keygen")
	}
	isolate(t)
	src := markedZip(t)
	dir := t.TempDir()
	zp := filepath.Join(dir, "styrman-harness-v1.2.3.zip")
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(zp, data, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	if err := os.WriteFile(filepath.Join(dir, "SHA256SUMS"), []byte("# styrman-harness v1.2.3\n"+hex.EncodeToString(sum[:])+"  styrman-harness-v1.2.3.zip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	key := filepath.Join(t.TempDir(), "key")
	for _, args := range [][]string{
		{"-q", "-t", "ed25519", "-N", "", "-C", "test@signer", "-f", key},
		{"-Y", "sign", "-q", "-f", key, "-n", harness.ReleaseNamespace, filepath.Join(dir, "SHA256SUMS")},
	} {
		if out, err := exec.Command(keygen, args...).CombinedOutput(); err != nil {
			t.Fatalf("ssh-keygen %v: %v\n%s", args, err, out)
		}
	}
	pub, err := os.ReadFile(key + ".pub")
	if err != nil {
		t.Fatal(err)
	}
	f := strings.Fields(string(pub))
	defer harness.SwapReleaseSigner("test@signer " + f[0] + " " + f[1])()

	var out bytes.Buffer
	vhome := t.TempDir()
	if err := run([]string{"update", "--from-zip", zp, "--home", vhome}, &out); err != nil {
		t.Fatalf("verified zip: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "VERIFIED against the signed SHA256SUMS of release v1.2.3") {
		t.Errorf("a verified zip does not say so:\n%s", out.String())
	}
	// The verdict outlives the output: status reads it from the manifest. It
	// needs a release version to resolve a source, as in TestStatusAfterZip.
	defer func(v string) { version = v }(version)
	version = "v1.2.3"
	var st bytes.Buffer
	if err := run([]string{"status", "--home", vhome}, &st); err != nil {
		t.Fatalf("status after a verified zip: %v\n%s", err, st.String())
	}
	if !strings.Contains(st.String(), "(verified v1.2.3)") {
		t.Errorf("status does not record the verified install:\n%s", st.String())
	}

	data[len(data)/2] ^= 0xff
	if err := os.WriteFile(zp, data, 0o644); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	err = run([]string{"update", "--from-zip", zp, "--home", home}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Errorf("tampered zip: err = %v, want checksum mismatch", err)
	}
	if strings.Contains(fmt.Sprint(err), "move SHA256SUMS") {
		t.Errorf("the refusal suggests removing the checksums: %v", err)
	}
	if strings.Contains(strings.ToLower(fmt.Sprint(err)), "unverified") {
		t.Errorf("a failed-check refusal names --unverified (ADR-040): %v", err)
	}
	if entries, _ := os.ReadDir(home); len(entries) != 0 {
		t.Errorf("a tampered zip left %d entries", len(entries))
	}
}
