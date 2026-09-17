package harness

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const pem = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAsecretbody\nmoresecretbody\n-----END RSA PRIVATE KEY-----"

func TestReportRedactJSONL(t *testing.T) {
	pemJSON, _ := json.Marshal("here is my key " + pem)
	quote := `{"n":12345678901234567890,"user_quote":` + string(pemJSON) + `}`
	akia := `{"project":"/p","user_quote":"use AKIAIOSFODNN7EXAMPLE for the demo","tags":["ok"]}`
	malformed := `{"user_quote": "token=abc123 broken`
	clean := `{"user_quote":"<b>no secrets</b> & more","n":1.50}`
	in := strings.Join([]string{quote, akia, malformed, clean, "not json at all", ""}, "\n")

	out, recs := redactLines("feedback.jsonl", []byte(in), true)
	lines := strings.Split(string(out), "\n")
	if len(lines) != 6 || lines[5] != "" {
		t.Fatalf("line structure changed: %q", out)
	}
	for _, leak := range []string{"secretbody", "IOSFODNN7EXAMPLE", "abc123"} {
		if strings.Contains(string(out), leak) {
			t.Errorf("output still carries %q:\n%s", leak, out)
		}
	}
	// A matching line goes whole, and in a JSONL file its marker is itself a
	// JSON object, so every redacted line parses — the malformed one included.
	for i, want := range []string{marker("[redacted by sh-agent report: 1 lines]", true), marker(redactedValue, true), marker(redactedValue, true)} {
		if lines[i] != want || !json.Valid([]byte(lines[i])) {
			t.Errorf("line %d = %q, want %q", i+1, lines[i], want)
		}
	}
	if lines[3] != clean || lines[4] != "not json at all" {
		t.Errorf("clean lines not byte-identical: %q %q", lines[3], lines[4])
	}
	want := []Redaction{
		{File: "feedback.jsonl", Line: 1, Pattern: "-----BEGIN", Lines: 1},
		{File: "feedback.jsonl", Line: 2, Pattern: "AKIA"},
		{File: "feedback.jsonl", Line: 3, Pattern: "secret assignment"},
	}
	if len(recs) != len(want) {
		t.Fatalf("redactions = %+v, want %+v", recs, want)
	}
	for i := range want {
		if recs[i] != want[i] {
			t.Errorf("redaction %d = %+v, want %+v", i, recs[i], want[i])
		}
	}
}

// TestReportRedactShapes: the common credential forms, in any case, are
// caught inside a JSONL value and on a text line; ordinary prose that only
// resembles them is not (security review 2026-09-17, finding 1).
func TestReportRedactShapes(t *testing.T) {
	secrets := []string{
		"export GITHUB_TOKEN=FAKEVALUE", "DB_PASSWORD=FAKEVALUE", "password: FAKEVALUE",
		`config was {"password": "FAKEVALUE"}`, "Authorization: Bearer FAKEVALUE",
		"github_pat_FAKEVALUE", "ghp_FAKEVALUE", "sk-ant-FAKEVALUE", "aws_secret_access_key = FAKEVALUE",
		"ASIAFAKEVALUE", "postgres://app:FAKEVALUE@db/x", "token = FAKEVALUE", "xoxb-FAKEVALUE",
		"GOCSPX-FAKEVALUE", "api_key=FAKEVALUE", "API-KEY: FAKEVALUE",
		// Round 2, finding 8: an empty URL user, and a non-Bearer scheme.
		"redis://:FAKEVALUE@localhost:6379/0", "Authorization: Basic FAKEVALUE",
		// Round 3, finding 3: the header name quoted, as code and JSON write it.
		`headers: {"Authorization": "Basic FAKEVALUE"}`, `requests.get(u, headers={'Authorization': 'Basic FAKEVALUE'})`,
	}
	for _, sec := range secrets {
		quoted, _ := json.Marshal(map[string]string{"user_quote": "they said " + sec})
		for name, in := range map[string]struct {
			data  string
			jsonl bool
		}{"jsonl": {string(quoted) + "\n", true}, "text": {"note: " + sec + "\n", false}} {
			out, recs := redactLines("f", []byte(in.data), in.jsonl)
			if len(recs) != 1 || strings.Contains(string(out), "FAKEVALUE") {
				t.Errorf("%s %q: out=%q recs=%+v", name, sec, out, recs)
			}
		}
	}
	for _, benign := range []string{"the task-list grew", "the token expired", "passwords are hard", "max_tokens 5", "ASIAN food"} {
		// "ASIAN" contains the ASIA prefix: accepted over-redaction, pinned
		// here so a change to it is deliberate.
		_, recs := redactLines("f", []byte(benign+"\n"), false)
		if want := benign == "ASIAN food"; (len(recs) == 1) != want {
			t.Errorf("%q: redactions %+v, want redacted=%v", benign, recs, want)
		}
	}
}

// TestReportRedactStructuralGaps: the ways a needle was present, content still
// left, and the run reported success — four from security review 2026-09-17
// finding 4, two from round 2 (a key split across array values, a secret-named
// key beside another hit, and a needle behind a JSON escape, which only the
// decoded text shows) — plus a key that named its own secret in report.json
// (finding 2).
func TestReportRedactStructuralGaps(t *testing.T) {
	cases := []struct {
		name, in string
		jsonl    bool
		leak     string
	}{
		{"duplicate key shadows the secret", `{"user_quote":"AKIAFAKEVALUE","user_quote":"clean"}` + "\n", true, "FAKEVALUE"},
		{"two keys concatenated on one line", "-----BEGIN A KEY-----\nFAKEBODY1\n-----END A KEY----------BEGIN B KEY-----\nFAKEBODY2\n-----END B KEY-----\nafter\n", false, "FAKEBODY"},
		{"dashes rewritten", "–––––BEGIN RSA PRIVATE KEY–––––\nFAKEBODY\n–––––END RSA PRIVATE KEY–––––\nafter\n", false, "FAKEBODY"},
		{"JSONL header line that is valid JSON alone", `"-----BEGIN RSA PRIVATE KEY-----"` + "\nFAKEBODY\n-----END RSA PRIVATE KEY-----\n" + `{"ok":"after"}` + "\n", true, "FAKEBODY"},
		{"key split across array values on one line", `{"tool_output":["-----BEGIN OPENSSH PRIVATE KEY-----","FAKEBODY1","FAKEBODY2","-----END OPENSSH PRIVATE KEY-----"]}` + "\n" + `{"ok":"after"}` + "\n", true, "FAKEBODY"},
		{"needle behind a JSON unicode escape", `{"user_quote":"\u0041KIAFAKEVALUE"}` + "\n" + `{"ok":"after"}` + "\n", true, "FAKEVALUE"},
		// Round 3: a header line that only mentions a footer (finding 2b); a
		// number beyond float64, and two records on one line, beside a needle
		// only the decoded text shows (finding 4).
		{"header line that mentions a footer", "pasted -----BEGIN OPENSSH PRIVATE KEY----- up to -----END below\nFAKEBODY1\nFAKEBODY2\n-----END OPENSSH PRIVATE KEY-----\nafter\n", false, "FAKEBODY"},
		{"number beyond float64 beside an escaped needle", `{"n":1e400,"q":"\u0041KIAFAKEVALUE"}` + "\n" + `{"ok":"after"}` + "\n", true, "FAKEVALUE"},
		{"two records on one line", `{"a":1}{"q":"\u0041KIAFAKEVALUE"}` + "\n" + `{"ok":"after"}` + "\n", true, "FAKEVALUE"},
		{"secret-named key beside another hit", `{"user_quote":"use ghp_FAKEPREFIX","config":{"api_key":"FAKEVALUE"}}` + "\n" + `{"ok":"after"}` + "\n", true, "FAKEVALUE"},
	}
	for _, c := range cases {
		out, recs := redactLines("f", []byte(c.in), c.jsonl)
		if recs == nil || strings.Contains(string(out), c.leak) {
			t.Errorf("%s: out=%q recs=%+v", c.name, out, recs)
		}
		if strings.Contains(c.in, "after") && !strings.Contains(string(out), "after") {
			t.Errorf("%s: redacted past the block: %q", c.name, out)
		}
	}
	for _, doc := range []string{`{"repo":"ghp_FAKEVALUE","repo":"/r"}`, `{"repo":"ya29.FAKE","db":{"password":"FAKEVALUE"}}`} {
		out, recs := redactDocument("state.json", []byte(doc))
		if recs == nil || strings.Contains(string(out), "FAKE") {
			t.Errorf("state.json %s: out=%q recs=%+v", doc, out, recs)
		}
	}

	_, recs := redactLines("retros.jsonl", []byte(`{"tool_output":{"https://api.example.test/v1?token=FAKEVALUE":"404"}}`), true)
	rec, _ := json.Marshal(recs)
	if len(recs) != 1 || strings.Contains(string(rec), "FAKEVALUE") {
		t.Errorf("secret in a key reached the record: %s", rec)
	}
}

func TestReportRedactCleanIsByteIdentical(t *testing.T) {
	jsonl := "{\"a\":\"x\",  \"b\":[1,2]}\n\n{\"c\":null}"
	md := "# Proposal\n\nNothing to hide.\n"
	state := "{\n    \"repo\": \"/r\",\n\"files\": []\n}\n"
	for name, tc := range map[string]struct {
		data   string
		redact func(string, []byte) ([]byte, []Redaction)
	}{
		"feedback.jsonl": {jsonl, func(n string, d []byte) ([]byte, []Redaction) { return redactLines(n, d, true) }},
		"proposals/x.md": {md, func(n string, d []byte) ([]byte, []Redaction) { return redactLines(n, d, false) }},
		"state.json":     {state, redactDocument},
	} {
		out, recs := tc.redact(name, []byte(tc.data))
		if recs != nil || !bytes.Equal(out, []byte(tc.data)) {
			t.Errorf("%s: out=%q recs=%v, want the input unchanged", name, out, recs)
		}
	}
}

func TestReportRedactText(t *testing.T) {
	in := "# Proposal\n\nbefore\n" + pem + "\nafter password=hunter2\nlast\n"
	out, recs := redactLines("proposals/p.md", []byte(in), false)
	want := "# Proposal\n\nbefore\n[redacted by sh-agent report: 4 lines]\n" + redactedValue + "\nlast\n"
	if string(out) != want {
		t.Errorf("got:\n%s\nwant:\n%s", out, want)
	}
	if len(recs) != 2 || recs[0].Line != 4 || recs[0].Lines != 4 || recs[0].Unterminated || recs[1].Pattern != "secret assignment" || recs[1].Line != 8 {
		t.Errorf("redactions = %+v", recs)
	}

	unterminated := "keep\n-----BEGIN OPENSSH PRIVATE KEY-----\nbody\nbody"
	out, recs = redactLines("proposals/q.md", []byte(unterminated), false)
	if string(out) != "keep\n[redacted by sh-agent report: 3 lines]" || len(recs) != 1 || !recs[0].Unterminated || recs[0].Lines != 3 {
		t.Errorf("unterminated: out=%q recs=%+v", out, recs)
	}
}

func TestReportRedactDocument(t *testing.T) {
	// Any needle anywhere replaces the whole document with a JSON marker.
	out, recs := redactDocument("state.json", []byte(`{"repo":"/r","nested":{"list":["x","ya29.token"]}}`))
	if string(out) != marker(redactedValue, true)+"\n" || !json.Valid(out) {
		t.Fatalf("out=%s", out)
	}
	if len(recs) != 1 || recs[0].Pattern != "ya29." || recs[0].Line != 0 {
		t.Errorf("redactions = %+v", recs)
	}
	// A needle only the decoded document shows (round 3, finding 6).
	if out, recs := redactDocument("state.json", []byte(`{"x":"\u0041KIAFAKE"}`)); len(recs) != 1 || strings.Contains(string(out), "FAKE") {
		t.Errorf("escaped needle in state.json: out=%s recs=%+v", out, recs)
	}
	// A needle in a KEY takes its line like any other hit.
	out, recs = redactLines("feedback.jsonl", []byte(`{"token=abc":"v"}`), true)
	if string(out) != marker(redactedValue, true) || len(recs) != 1 {
		t.Errorf("key hit: out=%s recs=%+v", out, recs)
	}
}

func writeConfig(t *testing.T, home, rel, content string) string {
	t.Helper()
	p := filepath.Join(home, ".config", "sh-agent", filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestReportCollect(t *testing.T) {
	home := t.TempDir()
	c, err := Collect(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Files) != 0 || strings.Join(c.Skipped, ",") != "feedback.jsonl (absent),retros.jsonl (absent),state.json (absent),proposals/ (absent)" {
		t.Errorf("empty home: %+v", c)
	}

	original := `{"user_quote":"AKIAIOSFODNN7EXAMPLE"}` + "\n"
	fb := writeConfig(t, home, "feedback.jsonl", original)
	writeConfig(t, home, "proposals/a.md", "fine\n")
	writeConfig(t, home, "proposals/notes.txt", "not a proposal\n")
	outside := filepath.Join(t.TempDir(), "secret.md")
	if err := os.WriteFile(outside, []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(home, ".config", "sh-agent", "proposals", "link.md")); err != nil {
		t.Fatal(err)
	}
	c, err = Collect(home)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, f := range c.Files {
		names = append(names, f.Name)
	}
	if strings.Join(names, ",") != "feedback.jsonl,proposals/a.md" {
		t.Errorf("collected %v", names)
	}
	if !strings.Contains(strings.Join(c.Skipped, ","), "proposals/link.md (not a regular file, not sent)") {
		t.Errorf("symlink not skipped and named: %v", c.Skipped)
	}
	if len(c.Redactions) != 1 {
		t.Errorf("redactions = %+v", c.Redactions)
	}
	if data, _ := os.ReadFile(fb); string(data) != original {
		t.Errorf("original modified: %q", data)
	}
}

func TestReportRecord(t *testing.T) {
	now := time.Date(2026, 9, 17, 9, 0, 0, 0, time.UTC)
	rec, notes := BuildRecord(t.TempDir(), "", "dave", now, nil)
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"ts", "host", "os", "arch", "vendor", "files", "harness_commit", "installed_at", "drift", "skills", "agents", "go", "redactions"} {
		if _, ok := m[k]; !ok {
			t.Errorf("record lacks %q: %s", k, data)
		}
	}
	for _, k := range []string{"vendor", "files", "harness_commit", "installed_at", "drift", "skills", "agents"} {
		if m[k] != nil {
			t.Errorf("%s = %v with no checkout and no install, want null", k, m[k])
		}
	}
	if r, ok := m["redactions"].([]any); !ok || len(r) != 0 {
		t.Errorf("redactions = %v, want []", m["redactions"])
	}
	if m["ts"] != "2026-09-17T09:00:00Z" || len(notes) != 2 {
		t.Errorf("ts=%v notes=%v", m["ts"], notes)
	}

	// With a plain-copy harness and an install, only harness_commit is null.
	repo, home := newRepo(t), t.TempDir()
	files, err := Plan(repo, nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if err := Install(repo, home, files, VendorAll, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	rec, _ = BuildRecord(home, repo, "dave", now, nil)
	if rec.HarnessCommit != nil || rec.Drift == nil || *rec.Drift != 0 || rec.Skills == nil || *rec.Skills != 2 || rec.Agents == nil || *rec.Agents != 1 || rec.Files == nil || rec.Vendor == nil {
		t.Errorf("installed record = %+v", rec)
	}
}

// gitEnv isolates git from the operator's own config for the tests that
// drive SendReport, which runs git in the process environment on purpose.
func gitEnv(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	cfg := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(cfg, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{
		"GIT_CONFIG_GLOBAL": cfg, "GIT_CONFIG_NOSYSTEM": "1",
		"GIT_AUTHOR_NAME": "t", "GIT_AUTHOR_EMAIL": "t@t", "GIT_COMMITTER_NAME": "t", "GIT_COMMITTER_EMAIL": "t@t",
	} {
		t.Setenv(k, v)
	}
	// Git reads an EMPTY GIT_WORK_TREE as set, so these are unset, not blanked.
	for _, k := range []string{"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE"} {
		t.Setenv(k, "") // registers the restore
		os.Unsetenv(k)
	}
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := gitOut(dir, args...)
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return out
}

func TestSendReport(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	origin := filepath.Join(base, "origin.git")
	// An EMPTY repository, as GitHub creates one: no commits, HEAD on main.
	git(t, base, "init", "--bare", "-q", "-b", "main", origin)

	home := filepath.Join(base, "home")
	writeConfig(t, home, "feedback.jsonl", `{"user_quote":"password=hunter2"}`+"\n")
	now := time.Date(2026, 9, 17, 9, 0, 0, 0, time.Local)
	run := ReportRun{Home: home, URL: origin, Host: "dave", Now: now}
	feedback := filepath.Join(home, ".config", "sh-agent", "feedback")

	// Dry run with the default path absent: clones nothing, writes nothing.
	var out bytes.Buffer
	dry := run
	dry.DryRun = true
	if err := SendReport(dry, &out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(feedback); !os.IsNotExist(err) || !strings.Contains(out.String(), "[dry-run] would clone") || !strings.Contains(out.String(), "feedback.jsonl:1 (secret assignment)") {
		t.Fatalf("dry run: lstat err=%v\n%s", err, out.String())
	}

	// First run: clones the empty repository, names the branch main, and
	// points main at the URL it pushed to, so a plain pull fetches from there.
	out.Reset()
	if err := SendReport(run, &out); err != nil {
		t.Fatalf("first run: %v\n%s", err, out.String())
	}
	if got := git(t, origin, "rev-list", "--count", "main"); got != "1" {
		t.Fatalf("origin main has %s commits after the first run\n%s", got, out.String())
	}
	if got := git(t, feedback, "config", "branch.main.remote"); got != origin {
		t.Errorf("branch.main.remote = %q, want %q", got, origin)
	}
	first := filepath.Join(feedback, "reports", "dave", now.Format("2006-01-02-150405"))
	sent, err := os.ReadFile(filepath.Join(first, "feedback.jsonl"))
	if err != nil || strings.Contains(string(sent), "hunter2") {
		t.Errorf("sent feedback.jsonl = %q, %v", sent, err)
	}
	var rec Record
	if data, err := os.ReadFile(filepath.Join(first, "report.json")); err != nil || json.Unmarshal(data, &rec) != nil || len(rec.Redactions) != 1 {
		t.Errorf("report.json: %v %+v", err, rec)
	}

	// A stray file in the checkout is left alone and named. Dry run against a
	// real checkout leaves git status exactly as it was.
	stray := filepath.Join(feedback, "stray.txt")
	if err := os.WriteFile(stray, []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := git(t, feedback, "status", "--porcelain", "--untracked-files=all")
	out.Reset()
	if err := SendReport(dry, &out); err != nil {
		t.Fatal(err)
	}
	if after := git(t, feedback, "status", "--porcelain", "--untracked-files=all"); after != before {
		t.Errorf("dry run changed the checkout: %q -> %q", before, after)
	}

	// Second run, same second: a second directory, the first untouched, the
	// stray not staged, a plain push.
	out.Reset()
	if err := SendReport(run, &out); err != nil {
		t.Fatalf("second run: %v\n%s", err, out.String())
	}
	if _, err := os.Stat(first + "-2"); err != nil {
		t.Errorf("second run directory: %v", err)
	}
	if got := git(t, origin, "rev-list", "--count", "main"); got != "2" {
		t.Errorf("origin main has %s commits after the second run", got)
	}
	if !strings.Contains(out.String(), "left alone") || !strings.Contains(git(t, feedback, "status", "--porcelain"), "?? stray.txt") {
		t.Errorf("stray file staged or not named:\n%s", out.String())
	}
	if files := git(t, feedback, "show", "--name-only", "--format=", "HEAD"); strings.Contains(files, "stray") {
		t.Errorf("stray file committed: %s", files)
	}

	// A behind checkout: the push is rejected, the commit stays local.
	other := filepath.Join(base, "other")
	git(t, base, "clone", "-q", origin, other)
	if err := SendReport(run, &out); err != nil {
		t.Fatal(err)
	}
	other2 := run
	other2.Feedback = other
	// A different host, or the commit is byte-identical to the one already
	// pushed — same tree, same second — and the push has nothing to reject.
	other2.Host = "erin"
	out.Reset()
	err = SendReport(other2, &out)
	if err == nil || !strings.Contains(err.Error(), "push rejected") || !strings.Contains(err.Error(), "pull --rebase") {
		t.Fatalf("behind checkout: err = %v", err)
	}
	if got := git(t, other, "rev-list", "--count", "@{upstream}..HEAD"); got != "1" {
		t.Errorf("rejected commit not kept locally: %s ahead", got)
	}
}

func TestResolveFeedbackFailures(t *testing.T) {
	gitEnv(t)
	home := t.TempDir()
	missing := filepath.Join(t.TempDir(), "nope.git")
	_, _, err := ResolveFeedback("", home, missing, false)
	if err == nil || !strings.Contains(err.Error(), "could not reach "+missing) || !strings.Contains(err.Error(), "--feedback-repo") {
		t.Errorf("failed clone: err = %v", err)
	}
	entries, _ := os.ReadDir(filepath.Join(home, ".config", "sh-agent"))
	if len(entries) != 0 {
		t.Errorf("failed clone left %v behind", entries)
	}

	plain := t.TempDir()
	if _, _, err := ResolveFeedback(plain, home, missing, false); err == nil || !strings.Contains(err.Error(), "never cloned into") {
		t.Errorf("explicit non-checkout: err = %v", err)
	}

	harnessRepo := newRepo(t)
	git(t, harnessRepo, "init", "-q")
	if _, _, err := ResolveFeedback(harnessRepo, home, missing, false); err == nil || !strings.Contains(err.Error(), "not the feedback repository") {
		t.Errorf("--feedback-repo naming a harness checkout: err = %v", err)
	}

	// Any checkout whose origin is not the feedback URL is refused, on the
	// explicit path and on the default one (security review 2026-09-17,
	// finding 3): its remote would receive the report.
	base := t.TempDir()
	feedbackOrigin := filepath.Join(base, "feedback.git")
	otherOrigin := filepath.Join(base, "other.git")
	git(t, base, "init", "--bare", "-q", "-b", "main", feedbackOrigin)
	git(t, base, "init", "--bare", "-q", "-b", "main", otherOrigin)
	other := filepath.Join(base, "other")
	git(t, base, "clone", "-q", otherOrigin, other)
	if _, _, err := ResolveFeedback(other, home, feedbackOrigin, false); err == nil || !strings.Contains(err.Error(), "not the feedback repository") {
		t.Errorf("explicit checkout of another repository: err = %v", err)
	}
	wrongHome := t.TempDir()
	git(t, base, "clone", "-q", otherOrigin, filepath.Join(wrongHome, ".config", "sh-agent", "feedback"))
	if _, _, err := ResolveFeedback("", wrongHome, feedbackOrigin, true); err == nil || !strings.Contains(err.Error(), "not the feedback repository") {
		t.Errorf("default path holding another repository, dry run: err = %v", err)
	}
	// The refusal names the expected URL, never the checkout's origin, which
	// can carry a token (round 2, finding 4).
	git(t, other, "remote", "set-url", "origin", "https://x-access-token:ghs_FAKETOKEN@example.test/a/b.git")
	if _, _, err := ResolveFeedback(other, home, feedbackOrigin, true); err == nil || strings.Contains(err.Error(), "FAKETOKEN") {
		t.Errorf("origin with a token: err = %v", err)
	}
	right := filepath.Join(base, "right")
	git(t, base, "clone", "-q", feedbackOrigin, right)
	if _, _, err := ResolveFeedback(right, home, feedbackOrigin, false); err != nil {
		t.Errorf("checkout of the feedback repository refused: %v", err)
	}
}

func TestSendReportRefusesRedirectsAndLinks(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	origin := filepath.Join(base, "origin.git")
	git(t, base, "init", "--bare", "-q", "-b", "main", origin)
	checkout := filepath.Join(base, "feedback")
	git(t, base, "clone", "-q", origin, checkout)
	run := ReportRun{Home: t.TempDir(), Feedback: checkout, URL: origin, Host: "dave", Now: time.Now()}

	// GIT_DIR would send every git call below to another repository.
	elsewhere := filepath.Join(base, "elsewhere")
	git(t, base, "init", "-q", elsewhere)
	t.Setenv("GIT_DIR", filepath.Join(elsewhere, ".git"))
	if err := SendReport(run, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "GIT_DIR is set") {
		t.Errorf("GIT_DIR set: err = %v", err)
	}
	os.Unsetenv("GIT_DIR")

	// A reports/<host> link pulled from the remote must not carry writes out.
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(checkout, "reports"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(checkout, "reports", "dave")); err != nil {
		t.Fatal(err)
	}
	if err := SendReport(run, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("reports/dave link: err = %v", err)
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 {
		t.Errorf("wrote through the link: %v", entries)
	}
}

// TestSendReportPushesOnlyToFeedbackURL: a pushurl, or an earlier url listed
// before the one verifyFeedback reads, would have sent the report elsewhere
// while it reported success (security review round 2, finding 1).
func TestSendReportPushesOnlyToFeedbackURL(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	origin := filepath.Join(base, "origin.git")
	elsewhere := filepath.Join(base, "elsewhere.git")
	git(t, base, "init", "--bare", "-q", "-b", "main", origin)
	git(t, base, "init", "--bare", "-q", "-b", "main", elsewhere)
	for name, setup := range map[string][][]string{
		"pushurl":    {{"config", "remote.origin.pushurl", elsewhere}},
		"second url": {{"config", "remote.origin.url", elsewhere}, {"config", "--add", "remote.origin.url", origin}},
	} {
		checkout := filepath.Join(base, strings.ReplaceAll(name, " ", "-"))
		git(t, base, "clone", "-q", origin, checkout)
		for _, args := range setup {
			git(t, checkout, args...)
		}
		run := ReportRun{Home: t.TempDir(), Feedback: checkout, URL: origin, Host: name, Now: time.Now()}
		var out bytes.Buffer
		if err := SendReport(run, &out); err != nil {
			t.Fatalf("%s: %v\n%s", name, err, out.String())
		}
		if _, err := gitOut(elsewhere, "rev-parse", "--verify", "--quiet", "main"); err == nil {
			t.Errorf("%s: the report reached the other repository", name)
		}
	}
	if got := git(t, origin, "rev-list", "--count", "main"); got != "2" {
		t.Errorf("feedback repository has %s commits, want 2", got)
	}
}

// TestSendReportDryRunShowsDroppedRecords: an unterminated key block redacts to
// the end of its file, so the dry-run list says so, where the operator reads
// before sending (security review round 3, finding 7).
func TestSendReportDryRunShowsDroppedRecords(t *testing.T) {
	home := t.TempDir()
	writeConfig(t, home, "feedback.jsonl", `{"q":"the header starts with -----BEGIN"}`+"\n"+`{"q":"a later correction"}`+"\n")
	var out bytes.Buffer
	run := ReportRun{Home: home, URL: filepath.Join(t.TempDir(), "absent.git"), Host: "dave", Now: time.Now(), DryRun: true}
	if err := SendReport(run, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "feedback.jsonl:1 (-----BEGIN) 2 lines, unterminated") {
		t.Errorf("dry run does not show the dropped records:\n%s", out.String())
	}
}
