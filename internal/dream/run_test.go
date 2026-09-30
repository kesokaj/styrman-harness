package dream

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kesokaj/styrman-harness/internal/search"
)

// fakeCLI puts an executable named name first on PATH. It records its
// arguments, working directory and stdin under rec, prints out and exits
// with code.
func fakeCLI(t *testing.T, name, out string, code int) (rec string) {
	t.Helper()
	bin := t.TempDir()
	rec = t.TempDir()
	if err := os.WriteFile(filepath.Join(rec, "out"), []byte(out), 0o600); err != nil {
		t.Fatal(err)
	}
	script := "#!/bin/sh\n" +
		"printf '%s\\n' \"$@\" > \"" + rec + "/args\"\n" +
		"pwd > \"" + rec + "/cwd\"\n" +
		"cat > \"" + rec + "/stdin\"\n" +
		"cat \"" + rec + "/out\"\n" +
		"exit " + strconv.Itoa(code) + "\n"
	if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return rec
}

const canned = `Here is the answer.
{"proposals":[
 {"title":"Vet before tag","kind":"rule","target":"project","file":"AGENTS.md","text":"Run go vet before every tag.","evidence":"memory/deploy-day.md"},
 {"title":"Release skill","kind":"skill","target":"harness","file":"skills/release/SKILL.md","text":"A release skill.","evidence":"trail changes/x.md"}],
 "skills":[{"name":"release-steps","description":"cut a release","body":"# Release\n1. vet"},{"name":"../escape","description":"bad","body":"x"}]}
trailing words`

func runOpts(f fixture, out *bytes.Buffer, vendor string) Options {
	return Options{Home: f.home, Vendor: vendor, All: true, Out: out, Now: func() time.Time { return now },
		Projects: func(string) ([]search.Project, error) { return []search.Project{project(f)}, nil }}
}

func scaffold(t *testing.T, f fixture) {
	t.Helper()
	os.MkdirAll(filepath.Join(f.root, ".agents", "docs"), 0o755)
}

func TestRunWritesProposalsDraftsAndState(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	write(t, filepath.Join(f.root, ".agents", "memory", "leak.md"), "token = FAKEVALUE\n", fresh)
	rec := fakeCLI(t, "claude", canned, 0)
	var out bytes.Buffer
	if err := Run(runOpts(f, &out, "claude")); err != nil {
		t.Fatalf("run: %v\n%s", err, out.String())
	}
	name := project(f).Name
	dir := filepath.Join(f.home, ".config", "sh-agent", "proposals", "dream")
	proj := read(t, filepath.Join(dir, "2026-09-30-"+name+".md"))
	for _, want := range []string{"owner: sh-agent dream", "project: " + f.root, "cli: claude", "### 1. Vet before tag", "**Status:** queued", "**Target:** project", "Run go vet before every tag."} {
		if !strings.Contains(proj, want) {
			t.Errorf("project file lacks %q:\n%s", want, proj)
		}
	}
	if strings.Contains(proj, "Release skill") {
		t.Error("a harness proposal landed in the project file (D23)")
	}
	harn := read(t, filepath.Join(dir, "harness", "2026-09-30-"+name+".md"))
	if !strings.Contains(harn, "### 1. Release skill") || !strings.Contains(harn, "**Target:** harness") {
		t.Errorf("harness file:\n%s", harn)
	}
	draft := read(t, filepath.Join(dir, "drafts", "2026-09-30-"+name+"-release-steps", "SKILL.md"))
	if !strings.HasPrefix(draft, "---\nname: release-steps\ndescription: \"cut a release\"\n---\n") {
		t.Errorf("draft:\n%s", draft)
	}
	if !strings.Contains(proj, "../escape") || !strings.Contains(proj, "dropped") {
		t.Error("the invalid draft name is not listed as dropped")
	}
	if strings.Contains(read(t, filepath.Join(rec, "stdin")), "FAKEVALUE") {
		t.Error("a secret reached the CLI")
	}
	args := read(t, filepath.Join(rec, "args"))
	if !strings.Contains(args, "--tools\n\n") || !strings.Contains(args, "--no-session-persistence") {
		t.Errorf("claude called without its tool-free flags:\n%s", args)
	}
	if cwd := strings.TrimSpace(read(t, filepath.Join(rec, "cwd"))); !sameDir(cwd, filepath.Join(f.home, ".config", "sh-agent", "dream-work")) {
		t.Errorf("cwd = %s, want dream-work", cwd)
	}
	st := readState(f.home)
	if !st.Projects[f.root].Equal(now) || st.LastRun.Projects != 1 || st.LastRun.Proposals != 2 || st.LastRun.Error != "" {
		t.Errorf("state = %+v", st)
	}
	log := read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log"))
	if !strings.Contains(log, f.root+" ok 2 proposals") || !strings.Contains(log, "- run 1 projects, 2 proposals") {
		t.Errorf("log:\n%s", log)
	}
	if _, err := os.Stat(filepath.Join(f.home, ".config", "sh-agent", "dream.lock")); !os.IsNotExist(err) {
		t.Error("lock left behind")
	}
	// A second run on the same day with no change calls nothing.
	os.Remove(filepath.Join(rec, "stdin"))
	if err := Run(runOpts(f, &out, "claude")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(rec, "stdin")); !os.IsNotExist(err) {
		t.Error("an unchanged project called the CLI")
	}
	if !strings.Contains(read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log")), f.root+" unchanged") {
		t.Error("unchanged not logged")
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRunAgyGetsThePromptAsArgument(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	rec := fakeCLI(t, "agy", `{"proposals":[],"skills":[]}`, 0)
	var out bytes.Buffer
	if err := Run(runOpts(f, &out, "antigravity")); err != nil {
		t.Fatal(err)
	}
	args := strings.Split(read(t, filepath.Join(rec, "args")), "\n")
	if args[0] != "-p" || !strings.HasPrefix(args[1], "IMPORTANT: every tool call") || !strings.Contains(strings.Join(args, "\n"), "You are the dreamer") || !strings.Contains(strings.Join(args, " "), "--mode plan --output-format text") {
		t.Errorf("agy args: %q", args[:3])
	}
	// Zero proposals: no file, watermark advanced, logged.
	if _, err := os.Stat(filepath.Join(f.home, ".config", "sh-agent", "proposals")); !os.IsNotExist(err) {
		t.Error("zero proposals wrote a file")
	}
	if st := readState(f.home); !st.Projects[f.root].Equal(now) {
		t.Errorf("watermark not advanced: %+v", st)
	}
	if !strings.Contains(read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log")), "ok 0 proposals") {
		t.Error("zero proposals not logged")
	}
}

func TestRunInvalidOutputKeepsRawAndWatermark(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	fakeCLI(t, "claude", "no json here\ntoken = FAKEVALUE\n", 0)
	var out bytes.Buffer
	err := Run(runOpts(f, &out, "claude"))
	var ee *ExitError
	if !errors.As(err, &ee) || ee.Code != 1 {
		t.Fatalf("err = %v, want exit 1", err)
	}
	if _, ok := readState(f.home).Projects[f.root]; ok {
		t.Error("watermark advanced on a failed project")
	}
	failed := filepath.Join(f.home, ".config", "sh-agent", "dream-failed")
	entries, _ := os.ReadDir(failed)
	if len(entries) == 0 {
		t.Fatal("raw output not kept")
	}
	raw := read(t, filepath.Join(failed, entries[0].Name()))
	if !strings.Contains(raw, "no json here") || strings.Contains(raw, "FAKEVALUE") {
		t.Errorf("kept raw output unredacted or missing:\n%s", raw)
	}
	if !strings.Contains(read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log")), "fail invalid output") {
		t.Error("failure not logged")
	}
}

func TestRunCLIExitFails(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	fakeCLI(t, "claude", `{"proposals":[],"skills":[]}`, 3)
	err := Run(runOpts(f, &bytes.Buffer{}, "claude"))
	var ee *ExitError
	if !errors.As(err, &ee) || ee.Code != 1 {
		t.Fatalf("err = %v, want exit 1", err)
	}
	if _, ok := readState(f.home).Projects[f.root]; ok {
		t.Error("watermark advanced on a failing CLI")
	}
}

func TestRunNoCLIFailsBeforeProjects(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	t.Setenv("PATH", t.TempDir())
	err := Run(runOpts(f, &bytes.Buffer{}, "claude"))
	var ee *ExitError
	if !errors.As(err, &ee) || ee.Code != 2 {
		t.Fatalf("err = %v, want exit 2", err)
	}
	st := readState(f.home)
	if !strings.Contains(st.LastRun.Error, "claude") {
		t.Errorf("last_run.error = %q", st.LastRun.Error)
	}
	if !strings.Contains(read(t, filepath.Join(f.home, ".config", "sh-agent", "dream.log")), "error") {
		t.Error("pre-project failure not logged")
	}
}

func TestRunLock(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	rec := fakeCLI(t, "claude", `{"proposals":[],"skills":[]}`, 0)
	lock := filepath.Join(f.home, ".config", "sh-agent", "dream.lock")
	write(t, lock, strconv.Itoa(os.Getpid()), fresh)
	var out bytes.Buffer
	if err := Run(runOpts(f, &out, "claude")); err != nil || !strings.Contains(out.String(), "already running") {
		t.Fatalf("live lock: %v, %q", err, out.String())
	}
	if _, err := os.Stat(filepath.Join(rec, "stdin")); !os.IsNotExist(err) {
		t.Error("ran under a live lock")
	}
	write(t, lock, "0", fresh) // pid 0 is never alive: stale
	if err := Run(runOpts(f, &out, "claude")); err != nil {
		t.Fatalf("stale lock: %v", err)
	}
	if _, err := os.Stat(filepath.Join(rec, "stdin")); err != nil {
		t.Error("a stale lock blocked the run")
	}
}

func TestRunPrunesOldFailures(t *testing.T) {
	f := newFixture(t)
	scaffold(t, f)
	fakeCLI(t, "claude", `{"proposals":[],"skills":[]}`, 0)
	failed := filepath.Join(f.home, ".config", "sh-agent", "dream-failed")
	write(t, filepath.Join(failed, "2026-08-01-x.txt"), "old", now.Add(-40*24*time.Hour))
	write(t, filepath.Join(failed, "2026-09-29-y.txt"), "recent", now.Add(-24*time.Hour))
	if err := Run(runOpts(f, &bytes.Buffer{}, "claude")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(failed, "2026-08-01-x.txt")); !os.IsNotExist(err) {
		t.Error("a failure older than 30 days was kept")
	}
	if _, err := os.Stat(filepath.Join(failed, "2026-09-29-y.txt")); err != nil {
		t.Error("a recent failure was pruned")
	}
}

func TestStatusLine(t *testing.T) {
	home := t.TempDir()
	if got := StatusLine(home); got != "never run; timer not installed" {
		t.Errorf("fresh: %q", got)
	}
	st := state{Projects: map[string]time.Time{}, LastRun: lastRun{At: now, Projects: 2, Proposals: 3}}
	b, _ := json.Marshal(st)
	write(t, statePath(home), string(b), fresh)
	if got := StatusLine(home); got != "last run 2026-09-30: 2 projects, 3 proposals; timer not installed" {
		t.Errorf("ok: %q", got)
	}
	st.LastRun.Error = "no claude on PATH"
	b, _ = json.Marshal(st)
	write(t, statePath(home), string(b), fresh)
	write(t, filepath.Join(home, ".config", "systemd", "user", "sh-agent-dream.timer"), "x", fresh)
	if got := StatusLine(home); got != "last run 2026-09-30 failed: no claude on PATH; timer installed" {
		t.Errorf("failed: %q", got)
	}
}
