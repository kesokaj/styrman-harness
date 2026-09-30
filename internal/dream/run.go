package dream

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/harness"
	"github.com/kesokaj/styrman-harness/internal/search"
)

const (
	cliTimeout = 10 * time.Minute
	keepFailed = 30 * 24 * time.Hour
)

var skillName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// response is the one JSON object the CLI is asked for.
type response struct {
	Proposals []proposal `json:"proposals"`
	Skills    []draft    `json:"skills"`
}

type proposal struct {
	Title, Kind, Target, File, Text, Evidence string
}

type draft struct {
	Name, Description, Body string
}

// realRun dreams over every project and writes what the CLI proposes.
func (o Options) realRun(projects []search.Project, st state) error {
	conf, err := harness.ConfigDir(o.Home)
	if err != nil {
		return &ExitError{2, "sh-agent dream: " + err.Error()}
	}
	unlock, held, err := lock(filepath.Join(conf, "dream.lock"))
	if err != nil {
		return &ExitError{2, "sh-agent dream: " + err.Error()}
	}
	if held != "" {
		fmt.Fprintf(o.Out, "sh-agent dream: already running (%s)\n", held)
		return nil
	}
	defer unlock()

	t0 := o.Now()
	logf := logger(filepath.Join(conf, "dream.log"), o.Now)
	prune(filepath.Join(conf, "dream-failed"), t0)
	finish := func(n, proposals int, runErr string) {
		st.LastRun = lastRun{At: t0, Projects: n, Proposals: proposals, Error: runErr}
		saveState(o.Home, st)
		line := fmt.Sprintf("- run %d projects, %d proposals", n, proposals)
		if runErr != "" {
			line += ", error " + runErr
		}
		logf(line)
	}

	argv := command(o.Vendor)
	if _, err := exec.LookPath(argv[0]); err != nil {
		finish(0, 0, fmt.Sprintf("no %s on PATH", argv[0]))
		return &ExitError{2, fmt.Sprintf("sh-agent dream: no %s on PATH: the dream job runs the CLI the manifest names (--vendor)", argv[0])}
	}
	work := filepath.Join(conf, "dream-work")
	if err := os.MkdirAll(work, 0o700); err != nil {
		finish(0, 0, err.Error())
		return &ExitError{2, "sh-agent dream: " + err.Error()}
	}

	failed, total := 0, 0
	for _, p := range projects {
		if why := skipReason(p); why != "" {
			logf(p.Root + " skip " + why)
			continue
		}
		c, err := collect(o.Home, o.Vendor, p, o.since(st, p.Root), maxPrompt)
		if err != nil {
			failed++
			logf(p.Root + " fail " + err.Error())
			continue
		}
		if !c.Changed {
			logf(p.Root + " unchanged")
			continue
		}
		stdout, stderr, err := call(argv, c.Prompt, work)
		if err != nil {
			failed++
			o.keepFailure(conf, c, t0, stdout, stderr)
			logf(p.Root + " fail " + err.Error())
			continue
		}
		resp, ok := parse(stdout)
		if !ok {
			failed++
			o.keepFailure(conf, c, t0, stdout, stderr)
			logf(p.Root + " fail invalid output")
			continue
		}
		if err := o.write(conf, c, resp, argv[0], t0); err != nil {
			failed++
			logf(p.Root + " fail " + err.Error())
			continue
		}
		st.Projects[p.Root] = t0
		saveState(o.Home, st)
		total += len(resp.Proposals)
		logf(fmt.Sprintf("%s ok %d proposals", p.Root, len(resp.Proposals)))
	}
	finish(len(projects), total, "")
	if failed > 0 {
		return &ExitError{1, fmt.Sprintf("sh-agent dream: %d project(s) failed; see %s", failed, filepath.Join(conf, "dream.log"))}
	}
	return nil
}

// call runs the CLI in work: claude takes the prompt on stdin, agy as the
// -p argument. The child keeps sh-agent's own environment (the plan's
// Interfaces), so it authenticates as the operator's CLI does.
func call(argv []string, prompt, work string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cliTimeout)
	defer cancel()
	args := append([]string(nil), argv[1:]...)
	var stdin string
	if argv[0] == "agy" {
		for i, a := range args {
			if a == "<prompt>" {
				args[i] = prompt
			}
		}
	} else {
		stdin = prompt
	}
	cmd := exec.CommandContext(ctx, argv[0], args...)
	cmd.Dir = work
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), errb.String(), fmt.Errorf("%s timed out after %s", argv[0], cliTimeout)
	}
	if err != nil {
		return out.String(), errb.String(), fmt.Errorf("%s: %v", argv[0], err)
	}
	return out.String(), errb.String(), nil
}

// parse takes the first {...} in out that decodes as the response.
func parse(out string) (response, bool) {
	for i := strings.IndexByte(out, '{'); i >= 0; {
		var r response
		dec := json.NewDecoder(strings.NewReader(out[i:]))
		dec.DisallowUnknownFields()
		if dec.Decode(&r) == nil && (r.Proposals != nil || r.Skills != nil) {
			return r, true
		}
		next := strings.IndexByte(out[i+1:], '{')
		if next < 0 {
			break
		}
		i += next + 1
	}
	return response{}, false
}

// write renders the proposals: project targets to the project's file,
// harness targets to proposals/dream/harness/ (D23), every valid draft skill
// to its own directory; the model's text is scrubbed again first.
func (o Options) write(conf string, c Collected, r response, cli string, t0 time.Time) error {
	day := t0.Format("2006-01-02")
	dir := filepath.Join(conf, "proposals", "dream")
	var projectP, harnessP []proposal
	for _, p := range r.Proposals {
		if p.Target == "harness" {
			harnessP = append(harnessP, p)
		} else {
			projectP = append(projectP, p)
		}
	}
	var drafts []draft
	var dropped []string
	for _, d := range r.Skills {
		if skillName.MatchString(d.Name) {
			drafts = append(drafts, d)
		} else {
			dropped = append(dropped, d.Name)
		}
	}
	for _, d := range drafts {
		path, err := unique(filepath.Join(dir, "drafts"), day+"-"+c.Name+"-"+d.Name, "")
		if err != nil {
			return err
		}
		desc, _ := json.Marshal(o.clean(d.Description))
		body := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s\n", d.Name, desc, strings.TrimSpace(o.clean(d.Body)))
		if err := os.MkdirAll(path, 0o700); err != nil {
			return err
		}
		if err := harness.WriteAtomic(filepath.Join(path, "SKILL.md"), []byte(body), 0o600); err != nil {
			return err
		}
	}
	if len(projectP) > 0 || len(drafts) > 0 || len(dropped) > 0 {
		var extra []string
		for _, d := range drafts {
			extra = append(extra, "Draft skill: drafts/"+day+"-"+c.Name+"-"+d.Name+"/SKILL.md")
		}
		for _, n := range dropped {
			extra = append(extra, fmt.Sprintf("Draft skill %q dropped: not a kebab-case name", n))
		}
		if err := o.render(dir, day, c, cli, projectP, extra); err != nil {
			return err
		}
	}
	if len(harnessP) > 0 {
		if err := o.render(filepath.Join(dir, "harness"), day, c, cli, harnessP, nil); err != nil {
			return err
		}
	}
	return nil
}

func (o Options) render(dir, day string, c Collected, cli string, ps []proposal, extra []string) error {
	path, err := unique(dir, day+"-"+c.Name, ".md")
	if err != nil {
		return err
	}
	counts := map[string]int{}
	for _, in := range c.Inputs {
		counts[in.Kind]++
	}
	var inputs []string
	for _, k := range promptOrder {
		inputs = append(inputs, fmt.Sprintf("%s %d", k, counts[k]))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "---\nowner: sh-agent dream\ndate: %s\nproject: %s\ncli: %s\ninputs: %s\n", day, c.Root, cli, strings.Join(inputs, ", "))
	if len(c.Dropped) > 0 {
		fmt.Fprintf(&b, "left out: %q\n", strings.Join(c.Dropped, "; "))
	}
	fmt.Fprintf(&b, "---\n# Dream proposals — %s   (%s)\n", c.Name, day)
	for i, p := range ps {
		target := "project"
		if p.Target == "harness" {
			target = "harness"
		}
		fmt.Fprintf(&b, "\n### %d. %s\n**Status:** queued\n**Target:** %s   **Kind:** %s   **File:** %s\n\n%s\n\nEvidence: %s\n",
			i+1, o.clean(p.Title), target, o.clean(p.Kind), o.clean(p.File), strings.TrimSpace(o.clean(p.Text)), o.clean(p.Evidence))
	}
	for _, e := range extra {
		b.WriteString("\n" + e + "\n")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	return harness.WriteAtomic(path, []byte(b.String()), 0o600)
}

func (o Options) clean(s string) string { return scrub(s, o.Home) }

// unique is dir/base+ext, or with -2, -3 … when that exists.
func unique(dir, base, ext string) (string, error) {
	for n := 1; n < 1000; n++ {
		name := base
		if n > 1 {
			name += "-" + strconv.Itoa(n)
		}
		p := filepath.Join(dir, name+ext)
		if _, err := os.Lstat(p); errors.Is(err, os.ErrNotExist) {
			return p, nil
		}
	}
	return "", fmt.Errorf("no free name for %s in %s", base, dir)
}

// keepFailure keeps the CLI's output of a failed project, scrubbed, for the
// operator: never in the CLI's working directory (D27).
func (o Options) keepFailure(conf string, c Collected, t0 time.Time, stdout, stderr string) {
	dir := filepath.Join(conf, "dream-failed")
	if os.MkdirAll(dir, 0o700) != nil {
		return
	}
	path, err := unique(dir, t0.Format("2006-01-02")+"-"+c.Name, ".txt")
	if err != nil {
		return
	}
	body := "stdout:\n" + o.clean(stdout) + "\nstderr:\n" + o.clean(stderr) + "\n"
	harness.WriteAtomic(path, []byte(body), 0o600)
}

// prune removes kept failures older than keepFailed.
func prune(dir string, t0 time.Time) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if info, err := e.Info(); err == nil && info.Mode().IsRegular() && t0.Sub(info.ModTime()) > keepFailed {
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}

// logger appends one timestamped line per call.
func logger(path string, now func() time.Time) func(string) {
	return func(line string) {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			return
		}
		fmt.Fprintf(f, "%s %s\n", now().UTC().Format(time.RFC3339), line)
		f.Close()
	}
}

func saveState(home string, st state) {
	if data, err := json.MarshalIndent(st, "", "  "); err == nil {
		harness.WriteAtomic(statePath(home), append(data, '\n'), 0o600)
	}
}

// lock takes path holding this pid. A lock whose pid is gone is stale and
// replaced; a live one is returned in held, and nothing is taken.
func lock(path string) (unlock func(), held string, err error) {
	for tries := 0; tries < 2; tries++ {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			fmt.Fprint(f, os.Getpid())
			f.Close()
			return func() { os.Remove(path) }, "", nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, "", err
		}
		data, _ := os.ReadFile(path)
		pid, perr := strconv.Atoi(strings.TrimSpace(string(data)))
		if perr == nil && harness.Alive(pid) {
			return nil, "pid " + strconv.Itoa(pid), nil
		}
		os.Remove(path)
	}
	return nil, "", errors.New("could not take " + path)
}
