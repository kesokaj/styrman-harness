package dream

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	for _, dir := range trailDirs {
		for _, f := range mdFiles(filepath.Join(p.Root, docs, dir)) {
			if body, t := readFile(f); t.After(since) {
				add(kindTrail, filepath.Join(dir, filepath.Base(f)), body, t)
			}
		}
	}

	for _, s := range skills(home, vendor) {
		add(kindSkill, "skill "+s.name, s.name+": "+s.description, time.Time{})
	}

	for _, tr := range claudeTranscript(home, p.Root, since) {
		add(kindTranscript, "claude "+tr.label, tr.text, tr.ts)
	}
	agy, err := search.AgyConversationText(home, p.Root, since)
	if err != nil {
		c.Dropped = append(c.Dropped, "antigravity transcripts: "+err.Error())
	}
	for _, a := range agy {
		add(kindTranscript, "antigravity "+a.ID, a.Text, a.Modified)
	}

	for i := range c.Inputs {
		c.Inputs[i].Text = scrub(c.Inputs[i].Text, home)
	}
	if err := c.fit(limit); err != nil {
		return c, err
	}
	c.Prompt = prompt(c)
	return c, nil
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
func claudeTranscript(home, root string, since time.Time) []transcript {
	base := filepath.Join(home, ".claude")
	if d := os.Getenv("CLAUDE_CONFIG_DIR"); d != "" {
		base = d
	}
	projects := filepath.Join(base, "projects")
	dirs, err := os.ReadDir(projects)
	if err != nil {
		return nil
	}
	root = filepath.Clean(root)
	var out []transcript
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
			if t := sessionText(path, root, since); t.text != "" {
				t.label = strings.TrimSuffix(f.Name(), ".jsonl")
				out = append(out, t)
			}
		}
	}
	return out
}

func sessionText(path, root string, since time.Time) transcript {
	f, err := os.Open(path)
	if err != nil {
		return transcript{}
	}
	defer f.Close()
	var b strings.Builder
	var last time.Time
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for sc.Scan() {
		var l struct {
			Type      string `json:"type"`
			Cwd       string `json:"cwd"`
			Timestamp string `json:"timestamp"`
			Message   struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(sc.Bytes(), &l) != nil || (l.Type != "user" && l.Type != "assistant") {
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
	return transcript{text: b.String(), ts: last}
}

// contentText is a message's text: a plain string, or the text parts of a
// parts array (tool calls and results are not prose).
func contentText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return strings.TrimSpace(s)
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
