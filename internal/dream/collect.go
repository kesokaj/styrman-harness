package dream

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/harness"
	"github.com/kesokaj/styrman-harness/internal/memory"
	"github.com/kesokaj/styrman-harness/internal/search"
)

// The kinds, in prompt order. Transcripts are dropped first when the prompt
// is over the cap, then the trail, then feedback, then skills; memory is
// always sent whole (D21).
const (
	kindMemory     = "memory"
	kindFeedback   = "feedback"
	kindTrail      = "trail"
	kindSkill      = "skill"
	kindTranscript = "transcript"
)

var promptOrder = []string{kindMemory, kindFeedback, kindTrail, kindSkill, kindTranscript}

// transcriptLineMax is the longest transcript line read; a longer one ends
// that file's read, which is named as left out. A test lowers it.
var transcriptLineMax = 16 * 1024 * 1024

// trailDirs are the records a dream reads under a project's docs root.
var trailDirs = []string{"bugs", "reviews", "retros", "changes", "handoffs"}

// collect gathers what project p would send: every memory file, and the
// feedback, trail and transcripts after since, with the installed skills.
// Changed is true when anything but the skills is newer than since.
func collect(home, vendor string, p search.Project, since time.Time, limit int) (Collected, error) {
	c := Collected{Root: p.Root, Name: p.Name}
	add := func(kind, label, text string, t time.Time) {
		if strings.TrimSpace(text) == "" {
			return
		}
		c.Inputs = append(c.Inputs, Input{Kind: kind, Label: label, Text: text, Time: t})
		if kind != kindSkill && t.After(since) {
			c.Changed = true
		}
	}

	if store, err := memory.Store(p.Root); err != nil {
		c.Dropped = append(c.Dropped, "memory: "+err.Error())
	} else {
		for _, f := range mdFiles(store) {
			body, t := readFile(f)
			add(kindMemory, "memory/"+filepath.Base(f), body, t)
		}
	}

	for _, e := range feedback(home, p.Root, since) {
		add(kindFeedback, "feedback.jsonl", e.text, e.ts)
	}

	docs := p.DocsRoot
	if docs == "" {
		docs = ".agents/docs"
	}
	dates, modified := trailDates(p.Root, docs)
	for _, dir := range trailDirs {
		// A symlinked docs root or record directory is never followed: a
		// repository could commit one pointing anywhere (security review
		// 2026-09-30).
		if linked(p.Root, docs, dir) {
			c.Dropped = append(c.Dropped, filepath.Join(docs, dir)+": a symlink on the way, not followed")
			continue
		}
		for _, f := range mdFiles(filepath.Join(p.Root, docs, dir)) {
			body, t := readFile(f)
			// A tracked, unmodified record is dated by its newest commit:
			// its mtime is the clone's or the checkout's, and by mtime a
			// whole trail looked new and evicted the transcripts (audit
			// 2026-09-30, SD4). An untracked or modified file keeps its mtime.
			if rel, err := filepath.Rel(p.Root, f); err == nil {
				if gt, ok := dates[filepath.ToSlash(rel)]; ok && !modified[filepath.ToSlash(rel)] {
					t = gt
				}
			}
			if t.After(since) {
				add(kindTrail, filepath.Join(dir, filepath.Base(f)), body, t)
			}
		}
	}

	for _, s := range skills(home, vendor) {
		add(kindSkill, "skill "+s.name, s.name+": "+s.description, time.Time{})
	}

	trs, trErrs := claudeTranscript(home, p.Root, since)
	for _, tr := range trs {
		add(kindTranscript, "claude "+tr.label, tr.text, tr.ts)
	}
	c.Dropped = append(c.Dropped, trErrs...)
	// An Antigravity read error fails the project, so its watermark stays and
	// the conversations are read again next time, never skipped for good
	// (code review 2026-09-30).
	agy, err := search.AgyConversationText(home, p.Root, since)
	if err != nil {
		return c, fmt.Errorf("reading Antigravity's conversations: %w", err)
	}
	for _, a := range agy {
		add(kindTranscript, "antigravity "+a.ID, a.Text, a.Modified)
	}

	for i := range c.Inputs {
		c.Inputs[i].Text = scrub(c.Inputs[i].Text, home)
	}
	for i := range c.Dropped {
		c.Dropped[i] = scrub(c.Dropped[i], home)
	}
	if err := c.fit(limit); err != nil {
		return c, err
	}
	c.Prompt = prompt(c)
	return c, nil
}

// linked reports whether any directory from root down to docs/dir is a
// symlink: .agents, .agents/docs and the record directory alike (code
// review round 2).
func linked(root, docs, dir string) bool {
	p := root
	for _, part := range strings.Split(filepath.ToSlash(filepath.Join(docs, dir)), "/") {
		p = filepath.Join(p, part)
		if info, err := os.Lstat(p); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

// scrub is what every input passes before it leaves sh-agent: report's
// redaction, and the home path written as ~.
func scrub(text, home string) string {
	text = harness.Redact(text)
	if home != "" && home != "/" {
		text = strings.ReplaceAll(text, strings.TrimSuffix(home, "/"), "~")
	}
	return text
}

// fit drops inputs until the prompt is under limit: transcripts oldest
// first, then trail files oldest first, then feedback, then skills (D21). A
// project whose memory alone is over the limit fails.
func (c *Collected) fit(limit int) error {
	for _, kind := range []string{kindTranscript, kindTrail, kindFeedback, kindSkill} {
		for len(prompt(*c)) > limit {
			i := oldest(c.Inputs, kind)
			if i < 0 {
				break
			}
			c.Dropped = append(c.Dropped, fmt.Sprintf("%s %s (over the %d KB cap)", c.Inputs[i].Kind, c.Inputs[i].Label, limit/1024))
			c.Inputs = append(c.Inputs[:i], c.Inputs[i+1:]...)
		}
	}
	if len(prompt(*c)) > limit {
		return errors.New("memory over cap")
	}
	return nil
}

func oldest(in []Input, kind string) int {
	best := -1
	for i, x := range in {
		if x.Kind == kind && (best < 0 || x.Time.Before(in[best].Time)) {
			best = i
		}
	}
	return best
}

// mdFiles lists the .md files directly in dir, by name, never through a
// pattern built from a path (the golang skill's glob rule).
func mdFiles(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if e.Type().IsRegular() && strings.HasSuffix(e.Name(), ".md") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(out)
	return out
}

// trailDates is, for every tracked file under docs, the time of its newest
// commit (paths relative to root, slash-separated), and the set of paths git
// reports as modified or untracked there. Both empty when git is absent or
// root is no checkout: every file then keeps its mtime.
func trailDates(root, docs string) (dates map[string]time.Time, modified map[string]bool) {
	dates, modified = map[string]time.Time{}, map[string]bool{}
	out, err := exec.Command("git", "-C", root, "log", "--format=%ct", "--name-only", "--no-renames", "--", docs).Output()
	if err != nil {
		return dates, modified
	}
	var t time.Time
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		if secs, err := strconv.ParseInt(line, 10, 64); err == nil {
			t = time.Unix(secs, 0)
			continue
		}
		if _, seen := dates[line]; !seen && !t.IsZero() {
			dates[line] = t
		}
	}
	out, err = exec.Command("git", trailStatusArgs(root, docs)...).Output()
	if err != nil {
		return dates, modified
	}
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) > 3 {
			modified[strings.TrimSpace(line[3:])] = true
		}
	}
	return dates, modified
}

// trailStatusArgs is the status command the nightly job runs in a project:
// the global options come first (git rejects them after the verb), it holds
// no index lock, and it runs no fsmonitor hook the project may have
// configured (security review rounds 2 and 3, 2026-10-01).
func trailStatusArgs(root, docs string) []string {
	return []string{"--no-optional-locks", "-c", "core.fsmonitor=false", "-C", root, "status", "--porcelain", "--untracked-files=all", "--", docs}
}

func readFile(path string) (string, time.Time) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", time.Time{}
	}
	body, err := os.ReadFile(path)
	if err != nil {
		return "", time.Time{}
	}
	return string(body), info.ModTime()
}

type entry struct {
	ts   time.Time
	text string
}

// feedback is the feedback log's lines for root after since, each as it is
// in the file.
func feedback(home, root string, since time.Time) []entry {
	f, err := os.Open(filepath.Join(home, ".config", "sh-agent", "feedback.jsonl"))
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []entry
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for sc.Scan() {
		var e struct {
			TS      string `json:"ts"`
			Project string `json:"project"`
		}
		if json.Unmarshal(sc.Bytes(), &e) != nil || e.Project != root {
			continue
		}
		ts, err := time.Parse(time.RFC3339, e.TS)
		if err != nil || !ts.After(since) {
			continue
		}
		out = append(out, entry{ts: ts, text: sc.Text()})
	}
	return out
}

type skill struct{ name, description string }

// skills are the installed skills' names and descriptions, from the skill
// directory of the vendor the manifest names.
func skills(home, vendor string) []skill {
	dir := filepath.Join(home, ".claude", "skills")
	if vendor == harness.VendorAntigravity {
		dir = filepath.Join(home, ".gemini", "config", "skills")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []skill
	for _, e := range entries {
		body, _ := readFile(filepath.Join(dir, e.Name(), "SKILL.md"))
		name, desc := frontField(body, "name"), frontField(body, "description")
		if name == "" {
			name = e.Name()
		}
		if desc != "" {
			out = append(out, skill{name, desc})
		}
	}
	return out
}

// frontField is a single-line frontmatter value, unquoted.
func frontField(body, key string) string {
	if !strings.HasPrefix(body, "---\n") {
		return ""
	}
	end := strings.Index(body[4:], "\n---")
	if end < 0 {
		return ""
	}
	for _, l := range strings.Split(body[4:4+end], "\n") {
		if v, ok := strings.CutPrefix(l, key+":"); ok {
			v = strings.TrimSpace(v)
			if uq, err := unquote(v); err == nil {
				return uq
			}
			return strings.Trim(v, `"'`)
		}
	}
	return ""
}

func unquote(v string) (string, error) {
	var s string
	err := json.Unmarshal([]byte(v), &s)
	return s, err
}

type transcript struct {
	label, text string
	ts          time.Time
}

// claudeTranscript is Claude Code's session text for root after since: the
// user and assistant text parts of lines whose cwd is root or under it.
// CLAUDE_CONFIG_DIR moves Claude Code's directory, as it does for the CLI.
func claudeTranscript(home, root string, since time.Time) ([]transcript, []string) {
	base := filepath.Join(home, ".claude")
	if d := os.Getenv("CLAUDE_CONFIG_DIR"); d != "" {
		base = d
	}
	projects := filepath.Join(base, "projects")
	dirs, err := os.ReadDir(projects)
	if err != nil {
		return nil, nil
	}
	root = filepath.Clean(root)
	var out []transcript
	var errs []string
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		files, _ := os.ReadDir(filepath.Join(projects, d.Name()))
		for _, f := range files {
			path := filepath.Join(projects, d.Name(), f.Name())
			info, err := f.Info()
			if err != nil || !info.Mode().IsRegular() || !strings.HasSuffix(f.Name(), ".jsonl") || !info.ModTime().After(since) {
				continue
			}
			t, err := sessionText(path, root, since)
			if err != nil {
				errs = append(errs, "claude transcript "+f.Name()+": "+err.Error())
			}
			if t.text != "" {
				t.label = strings.TrimSuffix(f.Name(), ".jsonl")
				out = append(out, t)
			}
		}
	}
	return out, errs
}

func sessionText(path, root string, since time.Time) (transcript, error) {
	f, err := os.Open(path)
	if err != nil {
		return transcript{}, err
	}
	defer f.Close()
	var b strings.Builder
	var last time.Time
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), transcriptLineMax)
	for sc.Scan() {
		var l struct {
			Type      string `json:"type"`
			Cwd       string `json:"cwd"`
			Timestamp string `json:"timestamp"`
			IsMeta    bool   `json:"isMeta"`
			Message   struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		// A meta line is the harness talking (a caveat, an injected
		// reminder), never the user or the model (security review 2026-10-01).
		if json.Unmarshal(sc.Bytes(), &l) != nil || (l.Type != "user" && l.Type != "assistant") || l.IsMeta {
			continue
		}
		cwd := filepath.Clean(l.Cwd)
		if cwd != root && !strings.HasPrefix(cwd, root+string(filepath.Separator)) {
			continue
		}
		ts, err := time.Parse(time.RFC3339Nano, l.Timestamp)
		if err != nil || !ts.After(since) {
			continue
		}
		if text := contentText(l.Message.Content); text != "" {
			fmt.Fprintf(&b, "%s: %s\n", l.Type, text)
			last = ts
		}
	}
	return transcript{text: b.String(), ts: last}, sc.Err()
}

// contentText is a message's text: a plain string, or the text parts of a
// parts array (tool calls and results are not prose).
func contentText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		s = strings.TrimSpace(s)
		// Claude Code writes a slash command's output and its command line
		// as user-typed string content; that is command output, not prose
		// (audit 2026-09-30, SD10).
		// Also a background task's or subagent's result and a system
		// reminder, which the harness writes into user turns (security
		// review of the fix).
		for _, p := range []string{"<local-command", "<command-", "<task-notification", "<system-reminder"} {
			if strings.HasPrefix(s, p) {
				return ""
			}
		}
		return s
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &parts) != nil {
		return ""
	}
	var texts []string
	for _, p := range parts {
		if p.Type == "text" && strings.TrimSpace(p.Text) != "" {
			texts = append(texts, strings.TrimSpace(p.Text))
		}
	}
	return strings.Join(texts, "\n")
}
