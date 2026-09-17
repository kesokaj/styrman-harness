package harness

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// FeedbackURL is where `sh-agent report` clones the feedback repository on
// first use, and where it pushes: the push names this URL rather than origin,
// so origin's own pushurl or extra url values cannot redirect it. Git's URL
// rewrites (url.*.insteadOf, pushInsteadOf) and a remote named after this URL
// still apply, as they do to every push on the machine: they are the
// operator's own configuration, trusted rather than refused (D14). SSH, the same transport and
// credential the operator already uses for this repository; change it by
// editing this constant (plan D4).
const FeedbackURL = "git@github.com:kesokaj/styrman-harness-feedback.git"

const redactedValue = "[redacted by sh-agent report]"

// Needles LOCATE a credential; they never define how much is removed. AKIA is
// a prefix and a PEM body carries no needle at all, so replacing the matched
// span would ship the key while reporting it redacted. What goes is always a
// whole line, a whole BEGIN-to-END block, or a whole document (D9, D13).
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
	{"authorization", regexp.MustCompile(`(?i)authorization["']?\s*:\s*\S`)},
	// The user may be empty: redis://:pass@host.
	{"url credentials", regexp.MustCompile(`://[^/\s:@]*:[^/\s@]+@`)},
}

// A key block opens at a PEM header, or at a header whose dashes a chat UI
// rewrote, and closes at the matching footer. A footer followed on the same
// line by another header keeps the block open (two keys concatenated).
var (
	blockOpen  = regexp.MustCompile(`-----BEGIN|BEGIN[A-Z ]*PRIVATE KEY`)
	blockClose = regexp.MustCompile(`-----END[^-]*(-----)?|END[A-Z ]*PRIVATE KEY`)
	blockLabel = regexp.MustCompile(`^BEGIN ([A-Z0-9][A-Z0-9 ]*[A-Z0-9])`)
)

// labelAt returns the label of the header that starts at s ("RSA PRIVATE
// KEY"), or "" when the header names none.
func labelAt(s string) string {
	s = s[strings.Index(s, "BEGIN"):]
	if m := blockLabel.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return ""
}

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

// Redaction records where content was removed and which needle fired — never
// the text. Line is 1-based and absent when a whole document went.
type Redaction struct {
	File         string `json:"file"`
	Line         int    `json:"line,omitempty"`
	Pattern      string `json:"pattern"`
	Lines        int    `json:"lines,omitempty"`
	Unterminated bool   `json:"unterminated,omitempty"`
}

// blockEnd returns the index of the line that closes the key block opened on
// line i, or the last non-empty line and false when nothing closes it.
func blockEnd(lines []string, i int) (int, bool) {
	open := blockOpen.FindStringIndex(lines[i])
	j, off, opened, label := i, open[1], i, labelAt(lines[i][open[0]:])
	for j < len(lines) {
		loc := blockClose.FindStringIndex(lines[j][off:])
		if loc == nil {
			j, off = j+1, 0
			continue
		}
		after := off + loc[1]
		// On the line its header opened, a block closes only at a footer naming
		// that header's label: a line that merely mentions a footer ("... up to
		// -----END below") would otherwise close it above the key body (D14).
		if j == opened && label != "" && !strings.Contains(lines[j][off+loc[0]:after], label) {
			off = after
			continue
		}
		if next := blockOpen.FindStringIndex(lines[j][after:]); next != nil {
			opened, label = j, labelAt(lines[j][after+next[0]:])
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

// marker is what replaces removed content. In a JSONL file it is a JSON
// object, so every line of the sent file still parses on its own.
func marker(text string, jsonl bool) string {
	if !jsonl {
		return text
	}
	b, _ := json.Marshal(map[string]string{"redacted": text})
	return string(b)
}

// redactLines redacts line-oriented content: a key block goes whole, and any
// other line carrying a needle goes whole. Nothing on a matching line
// survives, so no structure inside it — JSON values, keys, arrays — can carry
// a secret past the rule (D13, which retired the per-value JSON walk).
func redactLines(name string, data []byte, jsonl bool) ([]byte, []Redaction) {
	lines := strings.SplitAfter(string(data), "\n")
	var out strings.Builder
	var recs []Redaction
	for i := 0; i < len(lines); i++ {
		raw := lines[i]
		if raw == "" {
			continue
		}
		body := strings.TrimSuffix(raw, "\n")
		if blockOpen.MatchString(body) {
			end, closed := blockEnd(lines, i)
			rec := Redaction{File: name, Line: i + 1, Pattern: "-----BEGIN", Lines: end - i + 1, Unterminated: !closed}
			last := lines[end]
			out.WriteString(marker(fmt.Sprintf("[redacted by sh-agent report: %d lines]", rec.Lines), jsonl))
			out.WriteString(last[len(strings.TrimSuffix(last, "\n")):])
			recs = append(recs, rec)
			i = end
			continue
		}
		n, ok := needleIn(body)
		if !ok && jsonl {
			n, ok = needleIn(decoded(body))
		}
		if ok {
			out.WriteString(marker(redactedValue, jsonl) + raw[len(body):])
			recs = append(recs, Redaction{File: name, Line: i + 1, Pattern: n})
			continue
		}
		out.WriteString(raw)
	}
	return []byte(out.String()), recs
}

// decoded is a JSON line's text with its escapes undone, for matching only: a
// secret behind \" or \u0041 does not match the raw line. fmt prints each key
// followed by a colon, so a secret-named key still reads as an assignment.
// Every value on the line is decoded, numbers kept as written, up to the first
// value that does not decode; that value and everything after it on the line
// are not seen here (leading text or a BOM, NaN, a raw control character, a
// truncated record), nor is a value shadowed by a duplicate key. A documented
// limit, not an exhaustive list (D14).
func decoded(line string) string {
	dec := json.NewDecoder(strings.NewReader(line))
	dec.UseNumber()
	var b strings.Builder
	for {
		var v any
		if dec.Decode(&v) != nil {
			return b.String()
		}
		fmt.Fprintln(&b, v)
	}
}

// redactDocument redacts state.json, one JSON document: any needle anywhere in
// it replaces the whole document with a JSON marker. The metrics record reads
// the install manifest itself, so nothing report.json needs is lost.
func redactDocument(name string, data []byte) ([]byte, []Redaction) {
	n, ok := needleIn(string(data))
	if !ok {
		n, ok = needleIn(decoded(string(data)))
	}
	if !ok {
		return data, nil
	}
	return []byte(marker(redactedValue, true) + "\n"), []Redaction{{File: name, Pattern: n}}
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
	add := func(name string, redact func(string, []byte) ([]byte, []Redaction)) error {
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
		red, recs := redact(name, data)
		c.Files = append(c.Files, Collected{Name: name, Data: red})
		c.Redactions = append(c.Redactions, recs...)
		return nil
	}
	jsonl := func(name string, data []byte) ([]byte, []Redaction) { return redactLines(name, data, true) }
	text := func(name string, data []byte) ([]byte, []Redaction) { return redactLines(name, data, false) }

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
// operator's corrections into that repository (security review 2026-09-17,
// finding 3). A harness checkout fails it too. Where the commit is pushed is
// commitAndPush's, which names url itself.
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
		// The origin itself is never printed: it can carry a token.
		return fmt.Errorf("%s is not the feedback repository: its origin is not %q", dir, url)
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
			fmt.Fprintf(out, " (%s)", x.Pattern)
			if x.Lines > 0 {
				fmt.Fprintf(out, " %d lines", x.Lines)
			}
			if x.Unterminated {
				fmt.Fprint(out, ", unterminated: everything after it in this file was dropped")
			}
			fmt.Fprintln(out)
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
	if err := commitAndPush(dir, r.URL, paths, fmt.Sprintf("report: %s %s", host, filepath.Base(runRel)), out); err != nil {
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
func commitAndPush(dir, url string, paths []string, msg string, out io.Writer) error {
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
	// To the URL, not to origin, whose pushurl and extra urls would otherwise
	// decide the destination; the operator's URL rewrites still apply (see
	// FeedbackURL). -u makes a later plain `git pull --rebase` fetch from the
	// same place.
	if _, err := gitOut(dir, "push", "--quiet", "-u", url, "HEAD:main"); err != nil {
		if strings.Contains(err.Error(), "rejected") {
			return fmt.Errorf("push rejected; run `git -C %s pull --rebase` and re-run — the commit stays local", dir)
		}
		return fmt.Errorf("push failed, the commit stays local: %w", err)
	}
	return nil
}
