package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// FeedbackURL is where `sh-agent report` clones the feedback repository on
// first use, and the only origin it pushes to. SSH, the same transport and
// credential the operator already uses for this repository; change it by
// editing this constant (plan D4).
const FeedbackURL = "git@github.com:kesokaj/styrman-harness-feedback.git"

const redactedValue = "[redacted by sh-agent report]"

// Needles LOCATE a credential; they never define how much is removed. AKIA is
// a prefix and a PEM body carries no needle at all, so replacing the matched
// span would ship the key while reporting it redacted. What goes is always a
// whole JSON string value, a whole line, or a whole BEGIN-to-END block (D9).
//
// prefixNeedles are vendor token shapes, matched case-sensitively: they are
// distinctive as written. Nothing shorter than a vendor's own prefix is
// listed — a bare `sk-` would match every "task-".
var prefixNeedles = []string{
	"-----BEGIN", "PRIVATE KEY", "AKIA", "ASIA", "AIza", "ya29.", "GOCSPX-",
	"ghp_", "gho_", "ghs_", "ghu_", "github_pat_", "glpat-",
	"sk-ant-", "sk-proj-", "xoxb-", "xoxp-", "xoxa-", "xoxs-",
}

// keywordNeedles are assignment and header shapes, matched without regard to
// case, because the common forms are GITHUB_TOKEN=, DB_PASSWORD: and
// Authorization: Bearer (security review 2026-09-17, finding 1). Pattern is
// the name recorded in report.json, never the matched text.
var keywordNeedles = []struct {
	pattern string
	re      *regexp.Regexp
}{
	{"secret assignment", regexp.MustCompile(`(?i)(password|passwd|secret|token|api[_-]?key)[a-z0-9_-]*["']?\s*[:=]`)},
	{"bearer", regexp.MustCompile(`(?i)bearer\s+\S`)},
	{"url credentials", regexp.MustCompile(`://[^/\s:@]+:[^/\s@]+@`)},
}

// A key block opens at a PEM header, or at a header whose dashes a chat UI
// rewrote, and closes at the matching footer. A footer followed on the same
// line by another header keeps the block open (two keys concatenated).
var (
	blockOpen  = regexp.MustCompile(`-----BEGIN|BEGIN[A-Z ]*PRIVATE KEY`)
	blockClose = regexp.MustCompile(`-----END[^-]*(-----)?|END[A-Z ]*PRIVATE KEY`)
)

func needleIn(s string) (string, bool) {
	for _, n := range prefixNeedles {
		if strings.Contains(s, n) {
			return n, true
		}
	}
	for _, k := range keywordNeedles {
		if k.re.MatchString(s) {
			return k.pattern, true
		}
	}
	return "", false
}

// Redaction records where a value was removed and which needle fired — never
// the text. Line is 1-based and absent for a whole-document walk; Path is the
// JSON key path when the value sat inside JSON, and absent when a key itself
// matched, since that key is the secret.
type Redaction struct {
	File         string `json:"file"`
	Line         int    `json:"line,omitempty"`
	Path         string `json:"path,omitempty"`
	Pattern      string `json:"pattern"`
	Lines        int    `json:"lines,omitempty"`
	Unterminated bool   `json:"unterminated,omitempty"`
}

// blockEnd returns the index of the line that closes the key block opened on
// line i, or the last non-empty line and false when nothing closes it.
func blockEnd(lines []string, i int) (int, bool) {
	j, off := i, blockOpen.FindStringIndex(lines[i])[1]
	for j < len(lines) {
		loc := blockClose.FindStringIndex(lines[j][off:])
		if loc == nil {
			j, off = j+1, 0
			continue
		}
		after := off + loc[1]
		if next := blockOpen.FindStringIndex(lines[j][after:]); next != nil {
			off = after + next[1]
			continue
		}
		return j, true
	}
	last := len(lines) - 1
	if lines[last] == "" {
		last--
	}
	return last, false
}

// redactLines redacts line-oriented content. A key block is found first. With
// jsonLines, a line that is valid JSON has its string values walked, unless a
// block runs past it; a line the walk leaves carrying a needle (a shadowed
// duplicate key) and every other line fall through to the text rule, which
// never crashes on malformed input.
func redactLines(name string, data []byte, jsonLines bool) ([]byte, []Redaction, error) {
	lines := strings.SplitAfter(string(data), "\n")
	var out strings.Builder
	var recs []Redaction
	for i := 0; i < len(lines); i++ {
		raw := lines[i]
		if raw == "" {
			continue
		}
		body := strings.TrimSuffix(raw, "\n")
		nl := raw[len(body):]
		valid := jsonLines && json.Valid([]byte(body))

		if blockOpen.MatchString(body) {
			end, closed := blockEnd(lines, i)
			// A key wholly inside one JSON line is a string value: the walk
			// replaces it and the line stays JSON.
			if !(valid && closed && end == i) {
				rec := Redaction{File: name, Line: i + 1, Pattern: "-----BEGIN", Lines: end - i + 1, Unterminated: !closed}
				last := lines[end]
				fmt.Fprintf(&out, "[redacted by sh-agent report: %d lines]%s", rec.Lines, last[len(strings.TrimSuffix(last, "\n")):])
				recs = append(recs, rec)
				i = end
				continue
			}
		}

		if valid {
			red, lineRecs, err := redactJSON(name, i+1, []byte(body), false)
			if err != nil {
				return nil, nil, err
			}
			if lineRecs != nil {
				out.Write(red)
				out.WriteString(nl)
				recs = append(recs, lineRecs...)
				continue
			}
		}

		if n, ok := needleIn(body); ok {
			out.WriteString(redactedValue + nl)
			recs = append(recs, Redaction{File: name, Line: i + 1, Pattern: n})
			continue
		}
		out.WriteString(raw)
	}
	return []byte(out.String()), recs, nil
}

// redactJSON walks one JSON value. It returns nil records when nothing fired,
// so the caller ships the original bytes untouched — or, when the raw bytes
// still carry a needle, falls back to its text rule. A needle in an object KEY
// cannot be replaced without colliding keys, so the whole value becomes the
// marker as a JSON string — still valid JSON, still one line — and the record
// carries no path, which would be the key.
func redactJSON(name string, line int, data []byte, indent bool) ([]byte, []Redaction, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, nil, fmt.Errorf("%s: decoding JSON that validated: %w", name, err)
	}
	var recs []Redaction
	keyNeedle := ""
	v = walkJSON(v, "", func(path, needle string, isKey bool) {
		if isKey && keyNeedle == "" {
			keyNeedle = needle
		}
		recs = append(recs, Redaction{File: name, Line: line, Path: path, Pattern: needle})
	})
	if recs == nil {
		return nil, nil, nil
	}
	if keyNeedle != "" {
		v = redactedValue
		recs = []Redaction{{File: name, Line: line, Pattern: keyNeedle}}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if indent {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(v); err != nil {
		return nil, nil, fmt.Errorf("%s: re-encoding redacted JSON: %w", name, err)
	}
	outBytes := bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
	// A redacted line that is no longer JSON is a bug in this file, and it
	// fails the run rather than shipping something the reader cannot parse.
	if !json.Valid(outBytes) {
		return nil, nil, fmt.Errorf("%s line %d: redacted output is not valid JSON", name, line)
	}
	return outBytes, recs, nil
}

func walkJSON(v any, path string, fire func(path, needle string, isKey bool)) any {
	switch t := v.(type) {
	case string:
		if n, ok := needleIn(t); ok {
			fire(path, n, false)
			return redactedValue
		}
	case []any:
		for i := range t {
			t[i] = walkJSON(t[i], path+"["+strconv.Itoa(i)+"]", fire)
		}
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			p := k
			if path != "" {
				p = path + "." + k
			}
			if n, ok := needleIn(k); ok {
				fire(p, n, true)
			}
			t[k] = walkJSON(t[k], p, fire)
		}
	}
	return v
}

// redactDocument redacts a single JSON document, or treats it as text when it
// is not JSON, or when the walk found nothing but the raw bytes still carry a
// needle (a duplicate key shadowing a secret).
func redactDocument(name string, data []byte) ([]byte, []Redaction, error) {
	if !json.Valid(data) {
		return redactLines(name, data, false)
	}
	red, recs, err := redactJSON(name, 0, data, true)
	if err != nil {
		return nil, nil, err
	}
	if recs == nil {
		return redactLines(name, data, false)
	}
	return append(red, '\n'), recs, nil
}

// Collected is one file as it will be sent: already redacted.
type Collected struct {
	Name string // relative to ~/.config/sh-agent, slash-separated
	Data []byte
}

// Collection is everything `report` gathered from one home.
type Collection struct {
	Files      []Collected
	Skipped    []string // name and reason, for the output
	Redactions []Redaction
}

// Collect reads the learning-loop files under ~/.config/sh-agent/ and
// redacts each copy. The originals are never modified. A link or other
// non-regular file is skipped and named, never followed (BUG-036).
func Collect(home string) (Collection, error) {
	dir := filepath.Join(home, ".config", "sh-agent")
	var c Collection
	add := func(name string, redact func(string, []byte) ([]byte, []Redaction, error)) error {
		path := filepath.Join(dir, filepath.FromSlash(name))
		info, err := os.Lstat(path)
		switch {
		case errors.Is(err, os.ErrNotExist):
			c.Skipped = append(c.Skipped, name+" (absent)")
			return nil
		case err != nil:
			return fmt.Errorf("reading %s: %w", name, err)
		case !info.Mode().IsRegular():
			c.Skipped = append(c.Skipped, name+" (not a regular file, not sent)")
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}
		red, recs, err := redact(name, data)
		if err != nil {
			return err
		}
		c.Files = append(c.Files, Collected{Name: name, Data: red})
		c.Redactions = append(c.Redactions, recs...)
		return nil
	}
	jsonl := func(name string, data []byte) ([]byte, []Redaction, error) { return redactLines(name, data, true) }
	text := func(name string, data []byte) ([]byte, []Redaction, error) { return redactLines(name, data, false) }

	for _, name := range []string{"feedback.jsonl", "retros.jsonl"} {
		if err := add(name, jsonl); err != nil {
			return Collection{}, err
		}
	}
	if err := add("state.json", redactDocument); err != nil {
		return Collection{}, err
	}

	propDir := filepath.Join(dir, "proposals")
	info, err := os.Lstat(propDir)
	switch {
	case errors.Is(err, os.ErrNotExist):
		c.Skipped = append(c.Skipped, "proposals/ (absent)")
	case err != nil:
		return Collection{}, fmt.Errorf("reading proposals/: %w", err)
	case !info.IsDir():
		c.Skipped = append(c.Skipped, "proposals/ (not a directory, not sent)")
	default:
		entries, err := os.ReadDir(propDir)
		if err != nil {
			return Collection{}, fmt.Errorf("reading proposals/: %w", err)
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				if err := add("proposals/"+e.Name(), text); err != nil {
					return Collection{}, err
				}
			}
		}
	}
	return c, nil
}

// Record is report.json: one machine, one run. Every pointer field is null
// when its source is missing — no checkout, no install manifest — because a
// machine without them still reports that it is alive.
type Record struct {
	TS            string      `json:"ts"`
	Host          string      `json:"host"`
	OS            string      `json:"os"`
	Arch          string      `json:"arch"`
	Vendor        *string     `json:"vendor"`
	Files         *int        `json:"files"`
	HarnessCommit *string     `json:"harness_commit"`
	InstalledAt   *string     `json:"installed_at"`
	Drift         *int        `json:"drift"`
	Skills        *int        `json:"skills"`
	Agents        *int        `json:"agents"`
	Go            string      `json:"go"`
	Redactions    []Redaction `json:"redactions"`
}

// BuildRecord assembles the metrics record. repo is "" when no harness
// checkout was found. Nothing here is fatal: every source that cannot be read
// leaves its fields null and adds a note naming why.
func BuildRecord(home, repo, host string, now time.Time, redactions []Redaction) (Record, []string) {
	rec := Record{
		TS: now.Format(time.RFC3339), Host: host, OS: runtime.GOOS, Arch: runtime.GOARCH,
		// The RUNNING binary's toolchain, which built it — not `go` on PATH,
		// which can be a different one entirely (see Reinstall).
		Go:         runtime.Version(),
		Redactions: redactions,
	}
	if rec.Redactions == nil {
		rec.Redactions = []Redaction{}
	}
	var notes []string

	st, found, err := LoadState(home)
	vendor := VendorAll
	switch {
	case err != nil:
		notes = append(notes, fmt.Sprintf("install manifest unreadable (%v): vendor, files and installed_at are null", err))
	case !found:
		notes = append(notes, "nothing installed in this home: vendor, files and installed_at are null")
	default:
		vendor = st.Vendor
		files := len(st.Files)
		rec.Vendor, rec.Files, rec.InstalledAt = &vendor, &files, &st.InstalledAt
	}

	if repo == "" {
		notes = append(notes, "no harness checkout: harness_commit, drift, skills and agents are null")
		return rec, notes
	}
	if info := Inspect(repo); info.Commit != "" {
		rec.HarnessCommit = &info.Commit
	} else {
		notes = append(notes, "harness is a plain copy, not a git checkout: harness_commit is null")
	}
	if inv, err := LoadInventory(repo); err != nil {
		notes = append(notes, fmt.Sprintf("inventory unreadable (%v): skills and agents are null", err))
	} else {
		skills, agents := len(inv.Skills), len(inv.Agents)
		rec.Skills, rec.Agents = &skills, &agents
	}
	files, err := Plan(repo, nil, vendor)
	if err == nil {
		var r Report
		if r, err = Status(home, files); err == nil {
			drift := r.Counts["missing"] + r.Counts["differs"] + r.Counts["stale"]
			rec.Drift = &drift
		}
	}
	if err != nil {
		notes = append(notes, fmt.Sprintf("drift not computable (%v): drift is null", err))
	}
	return rec, notes
}

// ResolveFeedback returns the feedback checkout. A path given explicitly must
// already be a git checkout and is never cloned into. The default path,
// ~/.config/sh-agent/feedback, is cloned from url when absent — never on a
// dry run — through a temporary sibling, so a failed clone leaves nothing.
// Whichever it is, a checkout that exists must be the feedback repository
// itself (verifyFeedback) before anything is written into it.
func ResolveFeedback(explicit, home, url string, dryRun bool) (dir string, present bool, err error) {
	if explicit != "" {
		abs, err := filepath.Abs(explicit)
		if err != nil {
			return "", false, err
		}
		if !IsCheckout(abs) {
			return "", false, fmt.Errorf("--feedback-repo %s is not a git checkout; it is never cloned into", explicit)
		}
		return abs, true, verifyFeedback(abs, url)
	}
	dir = filepath.Join(home, ".config", "sh-agent", "feedback")
	if _, err := os.Lstat(dir); err == nil {
		if !IsCheckout(dir) {
			return "", false, fmt.Errorf("%s exists but is not a git checkout; move it aside or pass --feedback-repo", dir)
		}
		return dir, true, verifyFeedback(dir, url)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", false, err
	}
	if dryRun {
		return dir, false, nil
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return "", false, err
	}
	tmp := fmt.Sprintf("%s.clone-%d", dir, os.Getpid())
	_ = os.RemoveAll(tmp)
	if _, err := gitOut(filepath.Dir(dir), "clone", "--quiet", url, tmp); err != nil {
		_ = os.RemoveAll(tmp)
		return "", false, fmt.Errorf("could not reach %s: %v; clone it by hand or pass --feedback-repo", url, err)
	}
	if err := os.Rename(tmp, dir); err != nil {
		_ = os.RemoveAll(tmp)
		return "", false, err
	}
	return dir, true, nil
}

// verifyFeedback refuses any checkout but the feedback repository: dir must be
// the top of its own checkout and its origin exactly url, compared as
// configured rather than after insteadOf rewrites. Without it a typo in
// --feedback-repo, or a directory nested in another repository, commits the
// operator's corrections into that repository and pushes them to its remote
// (security review 2026-09-17, finding 3). A harness checkout fails it too.
func verifyFeedback(dir, url string) error {
	top, err := gitOut(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("%s: %w", dir, err)
	}
	want, err1 := filepath.EvalSymlinks(dir)
	got, err2 := filepath.EvalSymlinks(top)
	if err1 != nil || err2 != nil || want != got {
		return fmt.Errorf("%s is not the top of its own git checkout (git reports %s)", dir, top)
	}
	if origin, _ := gitOut(dir, "config", "--get", "remote.origin.url"); origin != url {
		return fmt.Errorf("%s is not the feedback repository: its origin is %q, not %q", dir, origin, url)
	}
	return nil
}

// gitRedirects are the variables that point git at a repository other than
// the directory it runs in.
var gitRedirects = []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_OBJECT_DIRECTORY", "GIT_COMMON_DIR"}

// ReportRun is one invocation of `sh-agent report`.
type ReportRun struct {
	Home, Repo string // Repo is "" when no harness checkout was found
	Feedback   string // --feedback-repo, or ""
	URL        string // cloned into the default path when it is absent
	Host       string
	Now        time.Time
	DryRun     bool
}

// SendReport collects, redacts, writes one run directory into the feedback
// checkout, and commits and pushes only the paths it wrote. It runs no git
// command in the harness repository beyond Inspect.
func SendReport(r ReportRun, out io.Writer) error {
	// With any of these set, every git call below would commit and push in
	// another repository whatever the checkout checks say (finding 3).
	for _, k := range gitRedirects {
		if _, set := os.LookupEnv(k); set {
			return fmt.Errorf("%s is set; report commits only in the feedback checkout — unset it and re-run", k)
		}
	}
	col, err := Collect(r.Home)
	if err != nil {
		return err
	}
	host := hostDir(r.Host)
	rec, notes := BuildRecord(r.Home, r.Repo, host, r.Now, col.Redactions)
	for _, s := range col.Skipped {
		fmt.Fprintf(out, "report: skipped %s\n", s)
	}
	for _, n := range notes {
		fmt.Fprintf(out, "report: %s\n", n)
	}
	if len(col.Files) == 0 {
		fmt.Fprintf(out, "report: nothing to send from %s; the metrics record still goes\n", filepath.Join(r.Home, ".config", "sh-agent"))
	}
	if n := len(col.Redactions); n > 0 {
		fmt.Fprintf(out, "report: redacted %d value(s) — see report.json\n", n)
		for _, x := range col.Redactions {
			fmt.Fprintf(out, "  %s", x.File)
			if x.Line > 0 {
				fmt.Fprintf(out, ":%d", x.Line)
			}
			if x.Path != "" {
				fmt.Fprintf(out, " %s", x.Path)
			}
			fmt.Fprintf(out, " (%s)\n", x.Pattern)
		}
	}

	recData, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	files := append(col.Files, Collected{Name: "report.json", Data: append(recData, '\n')})

	dir, present, err := ResolveFeedback(r.Feedback, r.Home, r.URL, r.DryRun)
	if err != nil {
		return err
	}
	runRel, err := runDir(dir, host, r.Now, !r.DryRun)
	if err != nil {
		return err
	}
	if r.DryRun {
		if !present {
			fmt.Fprintf(out, "[dry-run] would clone %s into %s\n", r.URL, dir)
		}
		for _, f := range files {
			fmt.Fprintf(out, "[dry-run] would write %s\n", filepath.Join(dir, runRel, filepath.FromSlash(f.Name)))
		}
		return nil
	}

	if dirty, err := gitOut(dir, "status", "--porcelain", "--untracked-files=all"); err != nil {
		return fmt.Errorf("reading the feedback checkout's state: %w", err)
	} else if dirty != "" {
		fmt.Fprintf(out, "report: left alone in %s, not staged:\n%s\n", dir, dirty)
	}

	paths := make([]string, 0, len(files))
	for _, f := range files {
		rel := filepath.Join(runRel, filepath.FromSlash(f.Name))
		if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(rel)), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, rel), f.Data, 0o644); err != nil {
			return err
		}
		paths = append(paths, rel)
	}
	if err := commitAndPush(dir, paths, fmt.Sprintf("report: %s %s", host, filepath.Base(runRel)), out); err != nil {
		return err
	}
	sha, _ := gitOut(dir, "rev-parse", "--short", "HEAD")
	fmt.Fprintf(out, "report: sent %s (%s)\n", filepath.ToSlash(runRel), sha)
	return nil
}

// runDir picks reports/<host>/<date>-<HHMMSS>, suffixed -2, -3… when that
// directory exists, so no run ever overwrites another. With create it makes
// the directory, which is what reserves the name.
func runDir(feedback, host string, now time.Time, create bool) (string, error) {
	base := filepath.Join("reports", host, now.Format("2006-01-02-150405"))
	// The checkout's content comes from a shared remote: a reports or
	// reports/<host> link pulled from it would carry every write below out of
	// the checkout, so a link is refused, never followed (BUG-036's rule).
	for _, rel := range []string{"reports", filepath.Join("reports", host)} {
		if info, err := os.Lstat(filepath.Join(feedback, rel)); err == nil && !info.IsDir() {
			return "", fmt.Errorf("%s in the feedback checkout is not a directory; refusing to write through it", rel)
		}
	}
	if create {
		if err := os.MkdirAll(filepath.Join(feedback, "reports", host), 0o755); err != nil {
			return "", err
		}
	}
	for n := 1; ; n++ {
		rel := base
		if n > 1 {
			rel = fmt.Sprintf("%s-%d", base, n)
		}
		abs := filepath.Join(feedback, rel)
		if !create {
			if _, err := os.Lstat(abs); errors.Is(err, os.ErrNotExist) {
				return rel, nil
			}
			continue
		}
		err := os.Mkdir(abs, 0o755)
		if err == nil {
			return rel, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return "", err
		}
	}
}

// hostDir makes a hostname safe as one path element.
func hostDir(host string) string {
	host = strings.NewReplacer("/", "_", `\`, "_").Replace(host)
	if host == "" || host == "." || host == ".." {
		return "unknown-host"
	}
	return host
}

// commitAndPush stages the named paths only, shows the staged stat, commits
// those paths only — so anything the operator already staged there stays
// out — and pushes. Each is its own git call, as the harness's Git rule asks.
func commitAndPush(dir string, paths []string, msg string, out io.Writer) error {
	add := append([]string{"add", "--"}, paths...)
	if _, err := gitOut(dir, add...); err != nil {
		return err
	}
	stat, err := gitOut(dir, append([]string{"diff", "--cached", "--stat", "--"}, paths...)...)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, stat)
	// A freshly created feedback repository has no commits, so the clone's
	// branch is unborn and named after the local git default; it is named
	// main before the first commit.
	if _, err := gitOut(dir, "rev-parse", "--verify", "--quiet", "HEAD"); err != nil {
		if _, err := gitOut(dir, "symbolic-ref", "HEAD", "refs/heads/main"); err != nil {
			return err
		}
	}
	commit := append([]string{"commit", "--quiet", "-m", msg, "--"}, paths...)
	if _, err := gitOut(dir, commit...); err != nil {
		return err
	}
	// Always to origin, which verifyFeedback checked, never to whatever
	// upstream the branch happens to track; -u sets main's upstream the first
	// time and is a no-op after.
	if _, err := gitOut(dir, "push", "--quiet", "-u", "origin", "HEAD"); err != nil {
		if strings.Contains(err.Error(), "rejected") {
			return fmt.Errorf("push rejected; run `git -C %s pull --rebase` and re-run — the commit stays local", dir)
		}
		return fmt.Errorf("push failed, the commit stays local: %w", err)
	}
	return nil
}
