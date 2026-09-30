package dream

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kesokaj/styrman-harness/internal/search"
)

var (
	now   = time.Date(2026, 9, 30, 20, 0, 0, 0, time.UTC)
	since = now.Add(-24 * time.Hour)
	old   = now.Add(-72 * time.Hour)
	fresh = now.Add(-1 * time.Hour)
)

// fixture is a scaffolded project and a home, with every input kind present
// once as changed and once as unchanged.
type fixture struct{ home, root string }

func newFixture(t *testing.T) fixture {
	t.Helper()
	f := fixture{home: t.TempDir(), root: t.TempDir()}
	run(t, f.root, "git", "init", "-q")
	write(t, filepath.Join(f.root, ".gitignore"), ".agents/memory\n.agents/scratch/*\n", old)
	write(t, filepath.Join(f.root, ".agents", "memory", "MEMORY.md"), "- [deploy-day](deploy-day.md) — fridays\n", old)
	write(t, filepath.Join(f.root, ".agents", "memory", "deploy-day.md"), "---\nname: deploy-day\n---\n\nWe deploy on Fridays.\n", old)
	write(t, filepath.Join(f.root, ".agents", "docs", "bugs", "INDEX.md"), "| BUG-001 | fresh bug row |\n", fresh)
	write(t, filepath.Join(f.root, ".agents", "docs", "changes", "2026-09-01-0000-old.md"), "an old change entry\n", old)
	write(t, filepath.Join(f.home, ".claude", "skills", "golang", "SKILL.md"), "---\nname: golang\ndescription: \"Go expert\"\n---\n", old)
	fb := func(ts time.Time, project, text string) string {
		b, _ := json.Marshal(map[string]string{"ts": ts.Format(time.RFC3339), "project": project, "original": text})
		return string(b) + "\n"
	}
	write(t, filepath.Join(f.home, ".config", "sh-agent", "feedback.jsonl"),
		fb(fresh, f.root, "fresh feedback for this project")+fb(old, f.root, "old feedback")+fb(fresh, "/elsewhere", "another project's feedback"), fresh)
	tl := func(ts time.Time, cwd, role, text string) string {
		b, _ := json.Marshal(map[string]any{"type": role, "cwd": cwd, "timestamp": ts.Format(time.RFC3339Nano),
			"message": map[string]any{"role": role, "content": []map[string]string{{"type": "text", "text": text}, {"type": "tool_use", "text": "not prose"}}}})
		return string(b) + "\n"
	}
	write(t, filepath.Join(f.home, ".claude", "projects", "-x", "s1.jsonl"),
		tl(fresh, f.root, "user", "please remember the deploy day")+
			tl(fresh, filepath.Join(f.root, "sub"), "assistant", "noted from a subdirectory")+
			tl(old, f.root, "user", "an old line before the watermark")+
			tl(fresh, "/elsewhere", "user", "another project's line")+
			`{"type":"summary","summary":"not a message"}`+"\n", fresh)
	return f
}

func write(t *testing.T, path, body string, mtime time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%s %v: %v\n%s", name, args, err, out)
	}
}

func project(f fixture) search.Project {
	return search.Project{Root: f.root, Name: filepath.Base(f.root), DocsRoot: ".agents/docs"}
}

func labels(c Collected) string {
	var b strings.Builder
	for _, in := range c.Inputs {
		b.WriteString(in.Kind + ":" + in.Label + "\n")
	}
	return b.String()
}

func TestCollectInputs(t *testing.T) {
	f := newFixture(t)
	c, err := collect(f.home, "claude", project(f), since, maxPrompt)
	if err != nil {
		t.Fatal(err)
	}
	if !c.Changed {
		t.Fatal("changed = false with a fresh bug row, feedback and transcript")
	}
	text := c.Prompt
	for _, want := range []string{
		"We deploy on Fridays.",           // memory, sent whole
		"fresh bug row",                   // trail after the watermark
		"fresh feedback for this project", // feedback for this root
		"Go expert",                       // installed skill
		"please remember the deploy day",  // transcript, at the root
		"noted from a subdirectory",       // transcript, under the root
	} {
		if !strings.Contains(text, want) {
			t.Errorf("prompt lacks %q\ninputs:\n%s", want, labels(c))
		}
	}
	for _, not := range []string{"an old change entry", "old feedback", "another project", "an old line before", "not prose", "not a message"} {
		if strings.Contains(text, not) {
			t.Errorf("prompt carries %q", not)
		}
	}
}

func TestCollectUnchanged(t *testing.T) {
	f := newFixture(t)
	later := now.Add(time.Hour)
	c, err := collect(f.home, "claude", project(f), later, maxPrompt)
	if err != nil {
		t.Fatal(err)
	}
	if c.Changed {
		t.Errorf("changed = true with nothing after the watermark:\n%s", labels(c))
	}
}

// Input 4 is global: a skill updated after the watermark never makes a
// project changed on its own (D-row in the plan's Interfaces).
func TestCollectSkillAloneIsNoChange(t *testing.T) {
	f := newFixture(t)
	later := now.Add(time.Hour)
	write(t, filepath.Join(f.home, ".claude", "skills", "golang", "SKILL.md"), "---\nname: golang\ndescription: \"Go expert, updated\"\n---\n", now.Add(2*time.Hour))
	c, err := collect(f.home, "claude", project(f), later, maxPrompt)
	if err != nil {
		t.Fatal(err)
	}
	if c.Changed {
		t.Error("a skill update alone made the project changed")
	}
}

func TestCollectRedactsAndHidesHome(t *testing.T) {
	f := newFixture(t)
	write(t, filepath.Join(f.root, ".agents", "memory", "leak.md"), "token = FAKEVALUE\nstored under "+f.home+"/x\n", fresh)
	c, err := collect(f.home, "claude", project(f), since, maxPrompt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(c.Prompt, "FAKEVALUE") {
		t.Error("a secret reached the prompt")
	}
	if strings.Contains(c.Prompt, f.home) || !strings.Contains(c.Prompt, "~/x") {
		t.Error("the home path was not rewritten to ~")
	}
}

func TestCollectCapDropsTranscriptThenTrail(t *testing.T) {
	f := newFixture(t)
	full, err := collect(f.home, "claude", project(f), since, maxPrompt)
	if err != nil {
		t.Fatal(err)
	}
	// Room for everything but the transcript.
	small := len(full.Prompt) - 40
	c, err := collect(f.home, "claude", project(f), since, small)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(c.Prompt, "please remember the deploy day") || !strings.Contains(c.Prompt, "fresh bug row") {
		t.Errorf("the transcript should go first and the trail stay:\n%s", labels(c))
	}
	if len(c.Dropped) == 0 || !strings.Contains(strings.Join(c.Dropped, " "), "transcript") {
		t.Errorf("dropped = %v, want the transcript named", c.Dropped)
	}
	// Room for less than memory alone: the project fails.
	if _, err := collect(f.home, "claude", project(f), since, 10); err == nil || !strings.Contains(err.Error(), "memory over cap") {
		t.Errorf("err = %v, want memory over cap", err)
	}
}

func TestDryRunWritesNothing(t *testing.T) {
	f := newFixture(t)
	var out bytes.Buffer
	o := Options{Home: f.home, Vendor: "claude", All: true, DryRun: true, Out: &out, Now: func() time.Time { return now },
		Projects: func(string) ([]search.Project, error) { return []search.Project{project(f)}, nil }}
	before := tree(t, f.home)
	if err := Run(o); err != nil {
		t.Fatal(err)
	}
	if after := tree(t, f.home); after != before {
		t.Errorf("dry run wrote:\nbefore %s\nafter  %s", before, after)
	}
	if !strings.Contains(out.String(), f.root) || !strings.Contains(out.String(), "claude -p") {
		t.Errorf("dry run output lacks the project or the command:\n%s", out.String())
	}
}

func tree(t *testing.T, dir string) string {
	t.Helper()
	var b strings.Builder
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err == nil {
			b.WriteString(p + " ")
		}
		return nil
	})
	return b.String()
}

func TestTargets(t *testing.T) {
	f := newFixture(t)
	sub := filepath.Join(f.root, "sub")
	os.MkdirAll(sub, 0o755)
	listed := []search.Project{{Root: "/p/one"}, {Root: "/p/two"}}
	o := Options{Home: f.home, Projects: func(string) ([]search.Project, error) { return listed, nil }}

	o.Project = sub
	got, err := o.targets()
	if err != nil || len(got) != 1 || !sameDir(got[0].Root, f.root) || got[0].DocsRoot != ".agents/docs" {
		t.Errorf("--project from a subdirectory: %+v, %v", got, err)
	}
	o.All = true
	if _, err := o.targets(); err == nil || !strings.Contains(err.Error(), "not both") {
		t.Errorf("--all with --project: %v", err)
	}
	o.Project = ""
	if got, err := o.targets(); err != nil || len(got) != 2 {
		t.Errorf("--all: %+v, %v", got, err)
	}
	o.All = false
	if _, err := o.targets(); err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("neither: %v", err)
	}
	o.Project = t.TempDir() // not a checkout
	if _, err := o.targets(); err == nil {
		t.Error("--project outside a checkout was accepted")
	}
}

func sameDir(a, b string) bool {
	ra, _ := filepath.EvalSymlinks(a)
	rb, _ := filepath.EvalSymlinks(b)
	return ra == rb
}

func TestSkipReason(t *testing.T) {
	f := newFixture(t)
	if why := skipReason(project(f)); why != "no .agents/docs" && why != "" {
		t.Fatalf("fixture: %q", why)
	}
	os.MkdirAll(filepath.Join(f.root, ".agents", "docs"), 0o755)
	if why := skipReason(project(f)); why != "" {
		t.Errorf("a scaffolded checkout skipped: %q", why)
	}
	if why := skipReason(search.Project{Root: filepath.Join(f.root, "gone")}); why != "gone" {
		t.Errorf("missing root: %q", why)
	}
	plain := t.TempDir()
	if why := skipReason(search.Project{Root: plain}); why != "not a git checkout" {
		t.Errorf("plain directory: %q", why)
	}
	bare := t.TempDir()
	run(t, bare, "git", "init", "-q")
	if why := skipReason(search.Project{Root: bare, DocsRoot: "docs"}); why != "no docs" {
		t.Errorf("no docs root: %q", why)
	}
}

// TestDryRunLines drives every line the dry run prints: a skipped root, an
// unchanged project, a project over the cap, and --project's prompt.
func TestDryRunLines(t *testing.T) {
	f := newFixture(t)
	calm := newFixture(t)
	later := now.Add(time.Hour)
	write(t, statePath(f.home), `{"projects":{"`+calm.root+`":"`+later.Format(time.RFC3339)+`"}}`, old)
	var out bytes.Buffer
	o := Options{Home: f.home, Vendor: "antigravity", All: true, DryRun: true, Out: &out, Now: func() time.Time { return now },
		Projects: func(string) ([]search.Project, error) {
			return []search.Project{project(f), project(calm), {Root: filepath.Join(f.root, "gone")}}, nil
		}}
	if err := Run(o); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{f.root + ": memory", "would run: agy -p", calm.root + ": unchanged since", "gone: skip (gone)"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("dry run lacks %q:\n%s", want, out.String())
		}
	}

	out.Reset()
	o.All, o.Project = false, f.root
	if err := Run(o); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "[dry-run] prompt:") || !strings.Contains(out.String(), "We deploy on Fridays.") {
		t.Errorf("--project dry run did not print the prompt:\n%s", out.String())
	}
}

func TestDryRunOverCapAndDropped(t *testing.T) {
	f := newFixture(t)
	big := strings.Repeat("a memory line that fills the store. ", 8000) // ~290 KB
	write(t, filepath.Join(f.root, ".agents", "memory", "huge.md"), big, fresh)
	var out bytes.Buffer
	o := Options{Home: f.home, Vendor: "claude", All: true, DryRun: true, Out: &out, Now: func() time.Time { return now },
		Projects: func(string) ([]search.Project, error) { return []search.Project{project(f)}, nil }}
	if err := Run(o); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "would fail (memory over cap)") {
		t.Errorf("over-cap memory not reported:\n%s", out.String())
	}

	// A store that refuses to be read is left out and named, and the run
	// goes on with the rest.
	g := newFixture(t)
	store := filepath.Join(g.root, ".agents", "memory")
	os.RemoveAll(store)
	if err := os.Symlink(t.TempDir(), store); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	o.Projects = func(string) ([]search.Project, error) { return []search.Project{project(g)}, nil }
	if err := Run(o); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "left out: memory:") || !strings.Contains(out.String(), "symlink") {
		t.Errorf("a refused store is not named:\n%s", out.String())
	}
}

func TestReadStateFallback(t *testing.T) {
	home := t.TempDir()
	if st := readState(home); st.Projects == nil || len(st.Projects) != 0 {
		t.Errorf("missing file: %+v", st)
	}
	write(t, statePath(home), "{not json", old)
	if st := readState(home); st.Projects == nil || len(st.Projects) != 0 {
		t.Errorf("corrupt file: %+v", st)
	}
	write(t, statePath(home), `{"projects":{"/p":"2026-09-29T00:00:00Z"}}`, old)
	o := Options{Now: func() time.Time { return now }}
	st := readState(home)
	if got := o.since(st, "/p"); !got.Equal(time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("saved watermark: %v", got)
	}
	if got := o.since(st, "/new"); !got.Equal(now.Add(-firstLook)) {
		t.Errorf("first look: %v", got)
	}
}

func TestCommand(t *testing.T) {
	if c := strings.Join(command("antigravity"), " "); !strings.HasPrefix(c, "agy -p") {
		t.Errorf("antigravity: %s", c)
	}
	for _, v := range []string{"claude", "all"} {
		if c := command(v); c[0] != "claude" || c[2] != "--tools" || c[3] != "" {
			t.Errorf("%s: %q", v, c)
		}
	}
}
