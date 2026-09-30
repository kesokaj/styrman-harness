package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		// Project review 2026-09-30, security finding 3: token shapes with no
		// keyword on the line — Stripe, PyPI (every token opens with the same
		// macaroon location), a Slack webhook, a JWT, an OpenAI key, an npm
		// token, a Hugging Face token, and a URL whose userinfo is a bare
		// token with no password.
		"sk_live_FAKEVALUE", "rk_live_FAKEVALUE", "whsec_FAKEVALUE",
		"pypi-AgEIcHlwaS5vcmcFAKEVALUE",
		"https://hooks.slack.com/services/T0/B0/FAKEVALUE",
		"eyJhbGciOiJIUzI1NiJ9.eyJFAKEVALUEx.FAKEVALUE",
		"sk-FAKEVALUEFAKEVALUEFAKEVALUEFAKEVALUE",
		"npm_FAKEVALUEFAKEVALUEFAKEVALUEFAKEVALUE",
		"hf_FAKEVALUEFAKEVALUEFAKEVALUEFAKEVALUE",
		"https://FAKEVALUEFAKEVALUEFAKEVALUE@github.com/o/r.git",
		// Security check of the fix, finding 2: the same tokens after a
		// literal `\n` (a JSONL line that failed to decode keeps its
		// escapes) and after an underscore, and a test.pypi.org token.
		`log:\neyJhbGciOiJIUzI1NiJ9.eyJFAKEVALUEx.FAKEVALUE`, `key_sk-FAKEVALUEFAKEVALUEFAKEVALUEFAKEVALUE`,
		`x\tnpm_FAKEVALUEFAKEVALUEFAKEVALUEFAKEVALUE`, `_hf_FAKEVALUEFAKEVALUEFAKEVALUEFAKEVALUE`,
		"pypi-AgENdGVzdC5weXBpLm9yZwFAKEVALUE",
		// Round 2 of the check: test-mode keys, service and admin keys, a
		// workflow webhook, a percent-escaped userinfo token.
		"sk_test_FAKEVALUE", "rk_test_FAKEVALUE", "sk-svcacct-FAKEVALUE", "sk-admin-FAKEVALUE",
		"https://hooks.slack.com/workflows/T0/A0/FAKEVALUE", "https://FAKEVALUE%2FFAKEVALUE%3DFAKE@host/r",
		// Round 3, one rule per kind: any Slack webhook path in any case, a
		// userinfo token of sixteen characters, a padded one, an sk-None key.
		"https://HOOKS.SLACK.COM/services/T0/B0/FAKEVALUE", "https://hooks.slack.com/triggers/T0/FAKEVALUE",
		"https://FAKEVALUE1234567@host/r", "https://FAKEVALUEFAKEVAL==@host/r", "sk-None-FAKEVALUE",
		// The security check of that rule: a webhook host with a port, a
		// trailing dot or a gov domain, an encoded path, an empty password.
		"https://hooks.slack.com:443/services/T0/B0/FAKEVALUE", "https://hooks.slack.com./services/FAKEVALUE",
		"https://hooks.slack-gov.com/services/FAKEVALUE", "https://hooks.slack.com%2Fservices%2FFAKEVALUE",
		"https://FAKEVALUEFAKEVALUE:@host/r",
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
	// Accepted over-redaction of the url-token rule, pinned so a change to
	// it is deliberate: a username of sixteen or more characters, and a
	// pathless URL whose line holds an unrelated `@` later on.
	for _, over := range []string{"ssh://first.lastname.long@bitbucket.org/x", `["https://app.example.com","ops@example.com"]`} {
		if _, recs := redactLines("f", []byte(over+"\n"), false); len(recs) != 1 {
			t.Errorf("%q: over-redaction is the pinned behaviour, got %+v", over, recs)
		}
	}
	for _, benign := range []string{
		"the task-list grew", "the token expired", "passwords are hard", "max_tokens 5", "ASIAN food",
		// The shapes above, too short or too plain to be a token.
		"npm_config_registry=https://r.example", "ssh://git@github.com/o/r.git", "sk-1234 is a short id",
		"the eyJ prefix marks base64 JSON", "hf_hub is the client library", "a whsec is a webhook secret",
		"https://deploy-bot@host/r", "see slack.com/apps for hooks",
	} {
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
	rec, notes := BuildRecord(t.TempDir(), Source{}, "dave", now, nil)
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"ts", "host", "os", "arch", "vendor", "files", "harness_commit", "harness_version", "installed_at", "drift", "skills", "agents", "go", "redactions"} {
		if _, ok := m[k]; !ok {
			t.Errorf("record lacks %q: %s", k, data)
		}
	}
	for _, k := range []string{"vendor", "files", "harness_commit", "harness_version", "installed_at", "drift", "skills", "agents"} {
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
	files, err := Plan(DirSource(repo), nil, VendorAll)
	if err != nil {
		t.Fatal(err)
	}
	if err := Install(DirSource(repo), home, files, VendorAll, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	rec, _ = BuildRecord(home, DirSource(repo), "dave", now, nil)
	if rec.HarnessCommit != nil || rec.HarnessVersion != nil || rec.Drift == nil || *rec.Drift != 0 || rec.Skills == nil || *rec.Skills != 2 || rec.Agents == nil || *rec.Agents != 1 || rec.Files == nil || rec.Vendor == nil {
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

	// Dry run: clones nothing, writes nothing, and says where it would clone.
	var out bytes.Buffer
	dry := run
	dry.DryRun = true
	if err := SendReport(dry, &out); err != nil {
		t.Fatal(err)
	}
	if entries, _ := os.ReadDir(configDir(home)); len(entries) != 1 || entries[0].Name() != "feedback.jsonl" {
		t.Errorf("dry run touched %s: %v", configDir(home), entries)
	}
	for _, want := range []string{"[dry-run] would clone", "and remove it after the push", "feedback.jsonl:1 (secret assignment)", "[dry-run] would write"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("dry run output lacks %q:\n%s", want, out.String())
		}
	}

	// First run: clones the empty repository, names the branch main, pushes,
	// and keeps nothing behind.
	out.Reset()
	if err := SendReport(run, &out); err != nil {
		t.Fatalf("first run: %v\n%s", err, out.String())
	}
	if got := git(t, origin, "rev-list", "--count", "main"); got != "1" {
		t.Fatalf("origin main has %s commits after the first run\n%s", got, out.String())
	}
	assertKeepsNothing(t, home)
	runRel := "reports/dave/" + now.Format("2006-01-02-150405")
	sent := git(t, origin, "show", "main:"+runRel+"/feedback.jsonl")
	if strings.Contains(sent, "hunter2") {
		t.Errorf("sent feedback.jsonl = %q", sent)
	}
	var rec Record
	if err := json.Unmarshal([]byte(git(t, origin, "show", "main:"+runRel+"/report.json")), &rec); err != nil || len(rec.Redactions) != 1 {
		t.Errorf("report.json: %v %+v", err, rec)
	}

	// Second run, same second: a second directory beside the first, both on
	// the remote, still nothing kept locally.
	out.Reset()
	if err := SendReport(run, &out); err != nil {
		t.Fatalf("second run: %v\n%s", err, out.String())
	}
	if got := git(t, origin, "rev-list", "--count", "main"); got != "2" {
		t.Errorf("origin main has %s commits after the second run", got)
	}
	if files := git(t, origin, "show", "--name-only", "--format=", "main"); !strings.Contains(files, runRel+"-2/") {
		t.Errorf("second run directory not on the remote: %q", files)
	}
	assertKeepsNothing(t, home)
}

// assertKeepsNothing: after a run without --feedback-repo, ~/.config/sh-agent
// holds no checkout of the feedback repository and no leftover clone.
func assertKeepsNothing(t *testing.T, home string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(configDir(home), "feedback")); !os.IsNotExist(err) {
		t.Errorf("a kept feedback checkout was left behind: stat err = %v", err)
	}
	left, _ := filepath.Glob(filepath.Join(configDir(home), "report-*"))
	if len(left) != 0 {
		t.Errorf("left behind: %v", left)
	}
}

// TestSendReportNamedCheckout: --feedback-repo keeps today's behaviour exactly
// — the checkout is the operator's, so it is written into and never removed,
// a stray file in it is left alone and named, and a rejected push leaves the
// commit there for them to rebase (plan § Report throwaway clone).
func TestSendReportNamedCheckout(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	origin := filepath.Join(base, "origin.git")
	git(t, base, "init", "--bare", "-q", "-b", "main", origin)
	checkout := filepath.Join(base, "feedback")
	git(t, base, "clone", "-q", origin, checkout)

	home := filepath.Join(base, "home")
	writeConfig(t, home, "feedback.jsonl", `{"user_quote":"a correction"}`+"\n")
	now := time.Date(2026, 9, 17, 10, 0, 0, 0, time.Local)
	run := ReportRun{Home: home, Feedback: checkout, URL: origin, Host: "dave", Now: now}

	var out bytes.Buffer
	if err := SendReport(run, &out); err != nil {
		t.Fatalf("first run: %v\n%s", err, out.String())
	}
	first := filepath.Join(checkout, "reports", "dave", now.Format("2006-01-02-150405"))
	if _, err := os.Stat(filepath.Join(first, "report.json")); err != nil {
		t.Errorf("the named checkout was not written into: %v", err)
	}
	if _, err := os.Stat(checkout); err != nil {
		t.Fatalf("the named checkout was removed: %v", err)
	}

	// A stray file is left alone and named, and a dry run changes nothing.
	if err := os.WriteFile(filepath.Join(checkout, "stray.txt"), []byte("mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := git(t, checkout, "status", "--porcelain", "--untracked-files=all")
	dry := run
	dry.DryRun = true
	out.Reset()
	if err := SendReport(dry, &out); err != nil {
		t.Fatal(err)
	}
	if after := git(t, checkout, "status", "--porcelain", "--untracked-files=all"); after != before {
		t.Errorf("dry run changed the checkout: %q -> %q", before, after)
	}
	out.Reset()
	if err := SendReport(run, &out); err != nil {
		t.Fatalf("second run: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "left alone") || !strings.Contains(git(t, checkout, "status", "--porcelain"), "?? stray.txt") {
		t.Errorf("stray file staged or not named:\n%s", out.String())
	}
	if files := git(t, checkout, "show", "--name-only", "--format=", "HEAD"); strings.Contains(files, "stray") {
		t.Errorf("stray file committed: %s", files)
	}

	// Behind the remote: the push is rejected and the commit stays local.
	other := filepath.Join(base, "other")
	git(t, base, "clone", "-q", origin, other)
	behind := run
	behind.Feedback = other
	behind.Host = "erin"
	out.Reset()
	if err := SendReport(run, &out); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	err := SendReport(behind, &out)
	if err == nil || !strings.Contains(err.Error(), "push rejected") || !strings.Contains(err.Error(), "pull --rebase") {
		t.Fatalf("behind checkout: err = %v", err)
	}
	if got := git(t, other, "rev-list", "--count", "@{upstream}..HEAD"); got != "1" {
		t.Errorf("rejected commit not kept locally: %s ahead", got)
	}
}

func TestResolveFeedbackFailures(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	missing := filepath.Join(t.TempDir(), "nope.git")

	plain := t.TempDir()
	if _, err := ResolveFeedback(plain, missing); err == nil || !strings.Contains(err.Error(), "never cloned into") {
		t.Errorf("explicit non-checkout: err = %v", err)
	}

	harnessRepo := newRepo(t)
	git(t, harnessRepo, "init", "-q")
	if _, err := ResolveFeedback(harnessRepo, missing); err == nil || !strings.Contains(err.Error(), "not the feedback repository") {
		t.Errorf("--feedback-repo naming a harness checkout: err = %v", err)
	}

	// Any checkout whose origin is not the feedback URL is refused: its remote
	// would receive the report (security review 2026-09-17, finding 3).
	feedbackOrigin := filepath.Join(base, "feedback.git")
	otherOrigin := filepath.Join(base, "other.git")
	git(t, base, "init", "--bare", "-q", "-b", "main", feedbackOrigin)
	git(t, base, "init", "--bare", "-q", "-b", "main", otherOrigin)
	other := filepath.Join(base, "other")
	git(t, base, "clone", "-q", otherOrigin, other)
	if _, err := ResolveFeedback(other, feedbackOrigin); err == nil || !strings.Contains(err.Error(), "not the feedback repository") {
		t.Errorf("explicit checkout of another repository: err = %v", err)
	}
	// The refusal names the expected URL, never the checkout's origin, which
	// can carry a token (round 2, finding 4).
	git(t, other, "remote", "set-url", "origin", "https://x-access-token:ghs_FAKETOKEN@example.test/a/b.git")
	if _, err := ResolveFeedback(other, feedbackOrigin); err == nil || strings.Contains(err.Error(), "FAKETOKEN") {
		t.Errorf("origin with a token: err = %v", err)
	}
	right := filepath.Join(base, "right")
	git(t, base, "clone", "-q", feedbackOrigin, right)
	if _, err := ResolveFeedback(right, feedbackOrigin); err != nil {
		t.Errorf("checkout of the feedback repository refused: %v", err)
	}

	// An unreachable URL on the throwaway path leaves nothing behind.
	home := filepath.Join(t.TempDir(), "home")
	err := SendReport(ReportRun{Home: home, URL: missing, Host: "dave", Now: time.Now()}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "could not reach "+missing) || !strings.Contains(err.Error(), "--feedback-repo") {
		t.Errorf("unreachable URL: err = %v", err)
	}
	entries, _ := os.ReadDir(configDir(home))
	if len(entries) != 0 {
		t.Errorf("a failed clone left %v behind", entries)
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

// seedFeedback builds a bare feedback repository holding two hosts' reports,
// with server-side filtering enabled — a file:// remote refuses
// --filter=blob:none without it, and the clone then silently fetches every
// blob, which would make the isolation check below vacuous (measured
// 2026-09-17).
func seedFeedback(t *testing.T, base string) string {
	t.Helper()
	origin := filepath.Join(base, "feedback.git")
	git(t, base, "init", "--bare", "-q", "-b", "main", origin)
	git(t, origin, "config", "uploadpack.allowFilter", "true")
	seed := filepath.Join(base, "seed")
	git(t, base, "clone", "-q", origin, seed)
	for host, body := range map[string]string{"dave": `{"host":"dave"}`, "otherbox": "OTHERBOX-SECRET-CANARY"} {
		dir := filepath.Join(seed, "reports", host, "2026-09-01-000000")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "report.json"), []byte(body+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, seed, "add", "-A")
	git(t, seed, "commit", "-qm", "seed")
	git(t, seed, "push", "-q", origin, "HEAD:main")
	return "file://" + origin
}

// localBlobHolds reports whether any object the clone holds LOCALLY contains
// needle. Lazy fetching is disabled throughout: a plain `git cat-file` on a
// promised object would fetch it and answer yes about a blob the clone never
// had (measured 2026-09-17).
func localBlobHolds(t *testing.T, dir, needle string) bool {
	t.Helper()
	list := exec.Command("git", "-C", dir, "cat-file", "--batch-all-objects", "--batch-check")
	list.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1")
	out, err := list.Output()
	if err != nil {
		t.Fatalf("cat-file --batch-all-objects: %v", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) < 2 || f[1] != "blob" {
			continue
		}
		show := exec.Command("git", "-C", dir, "cat-file", "blob", f[0])
		show.Env = append(os.Environ(), "GIT_NO_LAZY_FETCH=1")
		if data, err := show.Output(); err == nil && strings.Contains(string(data), needle) {
			return true
		}
	}
	return false
}

func configDir(home string) string { return filepath.Join(home, ".config", "sh-agent") }

// TestSendReportThrowawayClone: without --feedback-repo the run clones blob-less
// and sparse, pushes, and keeps nothing — no kept checkout, no report-<pid>
// directory — and the clone never holds another place's report (plan D8).
func TestSendReportThrowawayClone(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	url := seedFeedback(t, base)
	home := filepath.Join(base, "home")
	writeConfig(t, home, "feedback.jsonl", `{"user_quote":"a correction"}`+"\n")
	now := time.Date(2026, 9, 17, 12, 0, 0, 0, time.Local)

	var inspected string
	var sawCanary, sawOtherHostFile bool
	testHooks.beforeRemove = func(dir string) {
		inspected = dir
		sawCanary = localBlobHolds(t, dir, "OTHERBOX-SECRET-CANARY")
		_, err := os.Stat(filepath.Join(dir, "reports", "otherbox"))
		sawOtherHostFile = err == nil
	}
	defer func() { testHooks.beforeRemove = nil }()

	var out bytes.Buffer
	if err := SendReport(ReportRun{Home: home, URL: url, Host: "dave", Now: now}, &out); err != nil {
		t.Fatalf("report: %v\n%s", err, out.String())
	}
	if inspected == "" {
		t.Fatal("no throwaway clone was made")
	}
	if sawCanary {
		t.Error("the clone held another place's report content")
	}
	if sawOtherHostFile {
		t.Error("the clone checked out another place's reports directory")
	}
	if _, err := os.Stat(inspected); !os.IsNotExist(err) {
		t.Errorf("the throwaway clone was not removed: stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(configDir(home), "feedback")); !os.IsNotExist(err) {
		t.Error("a kept feedback checkout was left behind")
	}
	left, _ := filepath.Glob(filepath.Join(configDir(home), "report-*"))
	if len(left) != 0 {
		t.Errorf("left behind: %v", left)
	}

	// The run reached the feedback repository.
	origin := strings.TrimPrefix(url, "file://")
	files := git(t, origin, "show", "--name-only", "--format=", "main")
	want := "reports/dave/" + now.Format("2006-01-02-150405") + "/report.json"
	if !strings.Contains(files, want) {
		t.Errorf("the pushed commit holds %q, want %s", files, want)
	}
	if got := git(t, origin, "rev-list", "--count", "main"); got != "2" {
		t.Errorf("feedback repository has %s commits, want 2", got)
	}
}

// TestSendReportRetriesRejectedPush: a place that pushed between this run's
// clone and its push makes the first push a non-fast-forward. The run re-clones
// and re-applies once; a second rejection says so and keeps nothing.
func TestSendReportRetriesRejectedPush(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	url := seedFeedback(t, base)
	origin := strings.TrimPrefix(url, "file://")
	home := filepath.Join(base, "home")
	writeConfig(t, home, "feedback.jsonl", `{"user_quote":"a correction"}`+"\n")
	now := time.Date(2026, 9, 17, 13, 0, 0, 0, time.Local)

	// Another place pushes first, once, just before this run's own push.
	wedge := func(n int) func(string, int) {
		return func(_ string, attempt int) {
			if attempt > n {
				return
			}
			other := filepath.Join(base, fmt.Sprintf("other-%d-%d", n, attempt))
			git(t, base, "clone", "-q", origin, other)
			dir := filepath.Join(other, "reports", "elsewhere", fmt.Sprintf("2026-09-17-14%02d%02d", n, attempt))
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "report.json"), []byte("{}\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			git(t, other, "add", "-A")
			git(t, other, "commit", "-qm", "another place")
			git(t, other, "push", "-q", origin, "HEAD:main")
		}
	}

	testHooks.beforePush = wedge(1)
	var out bytes.Buffer
	err := SendReport(ReportRun{Home: home, URL: url, Host: "dave", Now: now}, &out)
	testHooks.beforePush = nil
	if err != nil {
		t.Fatalf("one rejection should be retried: %v\n%s", err, out.String())
	}
	files := git(t, origin, "show", "--name-only", "--format=", "main")
	if !strings.Contains(files, "reports/dave/"+now.Format("2006-01-02-150405")) {
		t.Errorf("the retry did not re-apply the run: %q", files)
	}
	if !strings.Contains(out.String(), "push rejected") {
		t.Errorf("the retry was silent:\n%s", out.String())
	}

	// Rejected on both attempts: named, and nothing kept.
	testHooks.beforePush = wedge(2)
	out.Reset()
	err = SendReport(ReportRun{Home: home, URL: url, Host: "dave", Now: now.Add(time.Hour)}, &out)
	testHooks.beforePush = nil
	if err == nil || !strings.Contains(err.Error(), "push rejected twice") || !strings.Contains(err.Error(), "nothing is kept locally") {
		t.Fatalf("two rejections: err = %v\n%s", err, out.String())
	}
	left, _ := filepath.Glob(filepath.Join(configDir(home), "*"))
	for _, p := range left {
		if strings.Contains(filepath.Base(p), "report-") || filepath.Base(p) == "feedback" {
			t.Errorf("left behind after two rejections: %s", p)
		}
	}
}

// TestRemoveKeptClone: `update` removes the clone earlier versions kept, but
// only when it is that clone, clean, and not ahead of its upstream (plan D9).
func TestRemoveKeptClone(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	url := seedFeedback(t, base)
	origin := strings.TrimPrefix(url, "file://")

	newHome := func(t *testing.T) string {
		t.Helper()
		home := filepath.Join(t.TempDir(), "home")
		if err := os.MkdirAll(configDir(home), 0o755); err != nil {
			t.Fatal(err)
		}
		git(t, base, "clone", "-q", url, filepath.Join(configDir(home), "feedback"))
		return home
	}

	t.Run("clean clone is removed", func(t *testing.T) {
		home := newHome(t)
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(configDir(home), "feedback")); !os.IsNotExist(err) {
			t.Errorf("not removed: %v", err)
		}
		if !strings.Contains(out.String(), "removed the kept feedback clone") {
			t.Errorf("silent removal:\n%s", out.String())
		}
	})

	t.Run("dirty clone is kept and named", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		if err := os.WriteFile(filepath.Join(kept, "mine.txt"), []byte("mine\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("a dirty clone must be kept: %v", err)
		}
		if !strings.Contains(out.String(), "uncommitted changes") {
			t.Errorf("reason not named:\n%s", out.String())
		}
	})

	t.Run("clone with unpushed commits is kept and named", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		if err := os.WriteFile(filepath.Join(kept, "local.txt"), []byte("local\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(t, kept, "add", "-A")
		git(t, kept, "commit", "-qm", "local work")
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("a clone with unpushed work must be kept: %v", err)
		}
		if !strings.Contains(out.String(), "none of its remote-tracking refs hold") {
			t.Errorf("reason not named:\n%s", out.String())
		}
		// The message says what the count was measured against, and how to
		// settle it, because that basis goes stale (2026-09-18).
		// The settling command must cover the same refs as the count, or it
		// tells the operator "nothing unique" about a commit on another
		// branch (code review 2026-09-18).
		if !strings.Contains(out.String(), "may well be pushed already") || !strings.Contains(out.String(), "--all --not --remotes FETCH_HEAD") {
			t.Errorf("the message does not say what it measured or how to check:\n%s", out.String())
		}
		// And it must not paste the directory into that command: a path with
		// a space or a metacharacter would ride along (security review).
		if strings.Contains(out.String(), "git -C") {
			t.Errorf("the message interpolates a path into a command to paste:\n%s", out.String())
		}
	})

	t.Run("another repository is kept and named", func(t *testing.T) {
		home := filepath.Join(t.TempDir(), "home")
		if err := os.MkdirAll(configDir(home), 0o755); err != nil {
			t.Fatal(err)
		}
		other := filepath.Join(base, "other.git")
		git(t, base, "init", "--bare", "-q", "-b", "main", other)
		git(t, base, "clone", "-q", other, filepath.Join(configDir(home), "feedback"))
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(configDir(home), "feedback")); err != nil {
			t.Errorf("another repository must be kept: %v", err)
		}
		if !strings.Contains(out.String(), "not the feedback repository") {
			t.Errorf("reason not named:\n%s", out.String())
		}
	})

	// A branch with no upstream holds commits no remote has, and the old
	// guard (`@{upstream}..HEAD`) could not even count them: its git call
	// fails, and a failed count was read as "nothing to lose" (security review
	// 2026-09-17, finding 1).
	t.Run("unpushed work on a branch with no upstream is kept", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		git(t, kept, "checkout", "-q", "-b", "local-only")
		if err := os.WriteFile(filepath.Join(kept, "mine.md"), []byte("mine\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(t, kept, "add", "-A")
		git(t, kept, "commit", "-qm", "work on a branch with no upstream")
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("a clone with unpushed work on an untracked branch must be kept: %v", err)
		}
		if !strings.Contains(out.String(), "none of its remote-tracking refs hold") {
			t.Errorf("reason not named:\n%s", out.String())
		}
		// The message says what the count was measured against, and how to
		// settle it, because that basis goes stale (2026-09-18).
		// The settling command must cover the same refs as the count, or it
		// tells the operator "nothing unique" about a commit on another
		// branch (code review 2026-09-18).
		if !strings.Contains(out.String(), "may well be pushed already") || !strings.Contains(out.String(), "--all --not --remotes FETCH_HEAD") {
			t.Errorf("the message does not say what it measured or how to check:\n%s", out.String())
		}
		// And it must not paste the directory into that command: a path with
		// a space or a metacharacter would ride along (security review).
		if strings.Contains(out.String(), "git -C") {
			t.Errorf("the message interpolates a path into a command to paste:\n%s", out.String())
		}
	})

	// The real shape this hit on the operator's own box: `report` pushed with
	// `-u <url>`, so the branch tracks a URL, `@{upstream}` does not resolve,
	// and no remote-tracking ref moves — every commit counts as unheld even
	// after it reached the remote. The clone is kept either way; what matters
	// is that the message does not claim the work is unpushed (2026-09-18).
	t.Run("a stale tracking ref keeps the clone without claiming lost work", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		plain := strings.TrimPrefix(url, "file://")
		git(t, kept, "config", "branch.main.remote", plain)
		if err := os.WriteFile(filepath.Join(kept, "run.json"), []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(t, kept, "add", "-A")
		git(t, kept, "commit", "-qm", "a run that was pushed")
		git(t, kept, "push", "-q", plain, "HEAD:main")
		// Key the guard on what RemoveKeptClone counts — commits no
		// remote-tracking ref holds — rather than on @{upstream}, which is a
		// symptom of the same setup and not what the code reads (code review
		// 2026-09-18).
		if n, err := gitOut(kept, "rev-list", "--count", "--all", "--not", "--remotes"); err != nil || n == "0" {
			t.Skipf("the stale-ref shape did not reproduce here (count %q, err %v)", n, err)
		}
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("the clone must be kept when the count cannot be trusted: %v", err)
		}
		if strings.Contains(out.String(), "unpushed") {
			t.Errorf("the message claims unpushed work it cannot know about:\n%s", out.String())
		}
		if !strings.Contains(out.String(), "may well be pushed already") {
			t.Errorf("the message does not say the commits may be pushed:\n%s", out.String())
		}
	})

	// A count that cannot be taken is not "nothing to lose": with the refs
	// unreadable, git exits non-zero and the clone is kept.
	t.Run("an uncountable clone is kept", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		packed := filepath.Join(kept, ".git", "packed-refs")
		if err := os.Chmod(packed, 0o000); err != nil {
			t.Skip("cannot make the refs unreadable here")
		}
		t.Cleanup(func() { _ = os.Chmod(packed, 0o644) })
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("a clone whose refs cannot be read must be kept: %v", err)
		}
		if !strings.Contains(out.String(), "kept ") {
			t.Errorf("reason not named:\n%s", out.String())
		}
	})

	// A stash holds work that is in no commit and no working tree.
	t.Run("stashed work is kept", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		if err := os.WriteFile(filepath.Join(kept, "reports", "dave", "2026-09-01-000000", "report.json"), []byte("changed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(t, kept, "stash", "push", "-q")
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("a clone with a stash must be kept: %v", err)
		}
		if !strings.Contains(out.String(), "stash") {
			t.Errorf("reason not named:\n%s", out.String())
		}
	})

	// A link at that path is not the clone: removing it would be removing
	// whatever it points at from the operator's view (BUG-036's rule).
	t.Run("a symlink there is kept", func(t *testing.T) {
		home := filepath.Join(t.TempDir(), "home")
		if err := os.MkdirAll(configDir(home), 0o755); err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(base, "link-target")
		git(t, base, "clone", "-q", url, target)
		if err := os.Symlink(target, filepath.Join(configDir(home), "feedback")); err != nil {
			t.Skip("symlinks unavailable here")
		}
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Lstat(filepath.Join(configDir(home), "feedback")); err != nil {
			t.Errorf("the link was removed: %v", err)
		}
		if _, err := os.Stat(filepath.Join(target, ".git")); err != nil {
			t.Errorf("the link's target was removed: %v", err)
		}
		if !strings.Contains(out.String(), "symlink") {
			t.Errorf("reason not named:\n%s", out.String())
		}
	})

	// A relative --home is what this repo's own rules prescribe for testing
	// (`--home .agents/scratch/x`), so the path is resolved before git sees it.
	t.Run("a relative home still resolves", func(t *testing.T) {
		home := newHome(t)
		wd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		rel, err := filepath.Rel(wd, home)
		if err != nil {
			t.Skip("no relative path from the working directory")
		}
		var out bytes.Buffer
		if err := RemoveKeptClone(rel, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(configDir(home), "feedback")); !os.IsNotExist(err) {
			t.Errorf("a clean clone under a relative home was not removed: %v", err)
		}
	})

	// git's answers can be redirected by its environment, and every guard but
	// the two Lstat checks asks git — on a path that DELETES. SendReport
	// refuses these variables for the same reason (security review
	// 2026-09-18).
	t.Run("a redirect variable keeps the clone", func(t *testing.T) {
		home := newHome(t)
		kept := filepath.Join(configDir(home), "feedback")
		t.Setenv("GIT_DIR", filepath.Join(t.TempDir(), "elsewhere.git"))
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(kept); err != nil {
			t.Errorf("a clone must be kept while git can be redirected: %v", err)
		}
		if !strings.Contains(out.String(), "GIT_DIR is set") {
			t.Errorf("reason not named:\n%s", out.String())
		}
	})

	t.Run("nothing there says nothing", func(t *testing.T) {
		home := filepath.Join(t.TempDir(), "home")
		var out bytes.Buffer
		if err := RemoveKeptClone(home, url, &out); err != nil {
			t.Fatal(err)
		}
		if out.String() != "" {
			t.Errorf("said something with no clone: %q", out.String())
		}
	})
	_ = origin
}

// TestSendReportSweepsInterruptedClones: a run killed before its deferred
// removal (SIGINT, a crash) leaves report-<pid> behind, holding a clone of the
// feedback repository — the very thing this phase removes from a machine. The
// next run sweeps such a directory when its process is gone, and leaves one
// whose process is still alive (security review 2026-09-17).
func TestSendReportSweepsInterruptedClones(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	url := seedFeedback(t, base)
	home := filepath.Join(base, "home")
	writeConfig(t, home, "feedback.jsonl", `{"user_quote":"a correction"}`+"\n")

	// A pid that has exited, and one that is alive (this test process).
	done := exec.Command("true")
	if err := done.Run(); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(configDir(home), fmt.Sprintf("report-%d", done.Process.Pid))
	// pid 1 is alive and is never this process, whose own report-<pid> path
	// the run under test clones into.
	alive := filepath.Join(configDir(home), "report-1")
	for _, dir := range []string{stale, alive} {
		if err := os.MkdirAll(filepath.Join(dir, "reports"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	var out bytes.Buffer
	if err := SendReport(ReportRun{Home: home, URL: url, Host: "dave", Now: time.Now()}, &out); err != nil {
		t.Fatalf("report: %v\n%s", err, out.String())
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Errorf("an interrupted run's clone was not swept: %v", err)
	}
	if !strings.Contains(out.String(), "removed") {
		t.Errorf("the sweep was silent:\n%s", out.String())
	}
	if _, err := os.Stat(alive); err != nil {
		t.Errorf("a live run's clone was swept: %v", err)
	}
}

// TestCloneFailureDoesNotPrintCredentials: git's stderr can carry the URL a
// rewrite produced, credentials included; verifyFeedback deliberately never
// prints an origin, and this path must not either (security review
// 2026-09-17, finding 6).
func TestCloneFailureDoesNotPrintCredentials(t *testing.T) {
	gitEnv(t)
	home := filepath.Join(t.TempDir(), "home")
	url := "https://x-access-token:ghs_FAKETOKEN@example.invalid/kesokaj/nope.git"
	err := SendReport(ReportRun{Home: home, URL: url, Host: "dave", Now: time.Now()}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("cloning an unreachable URL must fail")
	}
	if strings.Contains(err.Error(), "ghs_FAKETOKEN") {
		t.Errorf("the error carries the credential: %v", err)
	}
	if !strings.Contains(err.Error(), "could not reach") {
		t.Errorf("err = %v, want it to name the unreachable repository", err)
	}
}

// TestSendReportDoesNotRetryARefusal: only a non-fast-forward is the race the
// retry exists for. A pre-receive hook, a protected branch or a read-only
// credential produces `! [remote rejected]`, which will not go away by
// re-running — telling the operator to re-run, or retrying it once, is a wrong
// diagnosis of a permanent refusal (code review 2026-09-17).
func TestSendReportDoesNotRetryARefusal(t *testing.T) {
	gitEnv(t)
	base := t.TempDir()
	url := seedFeedback(t, base)
	origin := strings.TrimPrefix(url, "file://")
	hook := filepath.Join(origin, "hooks", "pre-receive")
	if err := os.MkdirAll(filepath.Dir(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hook, []byte("#!/bin/sh\necho 'policy: no' >&2\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(base, "home")
	writeConfig(t, home, "feedback.jsonl", `{"user_quote":"a correction"}`+"\n")

	clones := 0
	testHooks.beforePush = func(string, int) { clones++ }
	defer func() { testHooks.beforePush = nil }()

	var out bytes.Buffer
	err := SendReport(ReportRun{Home: home, URL: url, Host: "dave", Now: time.Now()}, &out)
	if err == nil {
		t.Fatal("a refused push must be an error")
	}
	if strings.Contains(err.Error(), "twice") || strings.Contains(out.String(), "re-cloning") {
		t.Errorf("a permanent refusal was retried as a race: %v\n%s", err, out.String())
	}
	if !strings.Contains(err.Error(), "push failed") {
		t.Errorf("err = %v, want it to say the push failed", err)
	}
	if clones != 1 {
		t.Errorf("clone attempts = %d, want 1", clones)
	}
	assertKeepsNothing(t, home)
}

// TestRunDirRefusesWhatItCannotRead: a name whose state cannot be read is an
// error, not "taken" — treating it as taken looped forever (code review
// 2026-09-17).
func TestRunDirRefusesWhatItCannotRead(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can read a mode-000 directory")
	}
	feedback := t.TempDir()
	reports := filepath.Join(feedback, "reports")
	if err := os.MkdirAll(filepath.Join(reports, "dave"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(reports, 0o000); err != nil {
		t.Skip("cannot make the directory unreadable here")
	}
	t.Cleanup(func() { _ = os.Chmod(reports, 0o755) })

	done := make(chan struct{})
	var err error
	go func() {
		_, err = runDir(feedback, "dave", time.Now(), false)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("runDir did not return: it is looping over a name it cannot read")
	}
	if err == nil {
		t.Error("an unreadable reports/<host> must be an error")
	}
}

func TestCollectNeverSendsSearchDB(t *testing.T) {
	home := t.TempDir()
	dir := filepath.Join(home, ".config", "sh-agent")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"feedback.jsonl", "search.db", "search.db-wal"} {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("{}\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	c, err := Collect(home)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range c.Files {
		if strings.HasPrefix(f.Name, "search.db") {
			t.Errorf("report would send %s", f.Name)
		}
	}
	named := strings.Join(c.Skipped, "\n")
	if !strings.Contains(named, "search.db (the global search index, a cache; never sent)") || !strings.Contains(named, "search.db-wal") {
		t.Errorf("skipped list does not name the index:\n%s", named)
	}
}

// TestPushErrorScrubsCredentials: a failed push is reported without the
// credentials a URL rewrite can put into git's message, as the clone failure
// already is, and only a non-fast-forward is classed as the retryable race
// (project review 2026-09-30, security finding 5).
func TestPushErrorScrubsCredentials(t *testing.T) {
	msg := "fatal: unable to access 'https://user:FAKEPW@host/x.git/': could not resolve host"
	err := pushError(errors.New(msg))
	if errors.Is(err, errPushRejected) {
		t.Errorf("an unreachable host is not a rejected push: %v", err)
	}
	if strings.Contains(err.Error(), "FAKEPW") || !strings.Contains(err.Error(), "[redacted]@") {
		t.Errorf("credentials reached the message: %v", err)
	}
	rej := pushError(errors.New("To https://user:FAKEPW@host/x.git\n ! [rejected] HEAD -> main (non-fast-forward)"))
	if !errors.Is(rej, errPushRejected) {
		t.Errorf("a non-fast-forward is the retryable race: %v", rej)
	}
	if strings.Contains(rej.Error(), "FAKEPW") {
		t.Errorf("credentials reached the rejection: %v", rej)
	}
}
