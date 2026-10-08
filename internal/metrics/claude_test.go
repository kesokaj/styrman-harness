package metrics

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const canary = "SECRETCANARY123"

// jl encodes one transcript line.
func jl(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b) + "\n"
}

type obj = map[string]any

// claudeFixture writes a home with one project: a main session that runs a
// Bash command, a Skill, a Read of a handoff, gets a denial and an error, is
// interrupted, compacts, and ends with a cost line; and one verifier
// subagent. The canary sits in every field that carries content.
func claudeFixture(t *testing.T) (home, root string) {
	t.Helper()
	home = t.TempDir()
	root = filepath.Join(t.TempDir(), "proj")
	os.MkdirAll(filepath.Join(root, ".agents", "docs", "handoffs"), 0o755)
	dir := filepath.Join(home, ".claude", "projects", "-proj")
	os.MkdirAll(filepath.Join(dir, "S1", "subagents"), 0o755)

	base := obj{"sessionId": "S1", "cwd": root}
	line := func(extra obj) obj {
		m := obj{}
		for k, v := range base {
			m[k] = v
		}
		for k, v := range extra {
			m[k] = v
		}
		return m
	}
	usage := func(in, out int) obj {
		return obj{"input_tokens": in, "output_tokens": out, "cache_read_input_tokens": 10, "cache_creation_input_tokens": 5}
	}
	var main strings.Builder
	for _, l := range []obj{
		line(obj{"type": "user", "uuid": "u1", "timestamp": "2026-10-08T10:00:00Z", "message": obj{"role": "user", "content": "please do " + canary}}),
		// one streamed message over two lines: usage is the max, one row
		line(obj{"type": "assistant", "uuid": "a1", "timestamp": "2026-10-08T10:00:01Z", "message": obj{"id": "m1", "model": "claude-opus-5-5",
			"usage": usage(3, 10), "content": []obj{{"type": "text", "text": "thinking about " + canary}}}}),
		line(obj{"type": "assistant", "uuid": "a2", "timestamp": "2026-10-08T10:00:02Z", "message": obj{"id": "m1", "model": "claude-opus-5-5",
			"usage": usage(3, 40), "content": []obj{
				{"type": "tool_use", "id": "t1", "name": "Bash", "input": obj{"command": "git commit -m '" + canary + "' && go test ./... # " + canary}},
				{"type": "tool_use", "id": "t2", "name": "Skill", "input": obj{"skill": "feature", "args": canary}},
				{"type": "tool_use", "id": "t3", "name": "Read", "input": obj{"file_path": filepath.Join(root, ".agents", "docs", "handoffs", canary+".md")}},
				{"type": "tool_use", "id": "t4", "name": "Bash", "input": obj{"command": "rm -rf /" + canary}},
			}}}),
		line(obj{"type": "user", "uuid": "u2", "timestamp": "2026-10-08T10:00:03Z", "message": obj{"role": "user", "content": []obj{
			{"type": "tool_result", "tool_use_id": "t1", "is_error": true, "content": "Exit code 1\n" + canary},
			{"type": "tool_result", "tool_use_id": "t2", "is_error": false, "content": canary},
			{"type": "tool_result", "tool_use_id": "t3", "is_error": false, "content": []obj{{"type": "text", "text": canary}}},
			{"type": "tool_result", "tool_use_id": "t4", "is_error": true, "content": "Permission for this action was denied by the Claude Code auto mode classifier. " + canary},
		}}}),
		line(obj{"type": "user", "uuid": "u3", "timestamp": "2026-10-08T10:00:04Z", "message": obj{"role": "user", "content": []obj{{"type": "text", "text": "[Request interrupted by user]"}}}}),
		line(obj{"type": "user", "uuid": "u4", "isMeta": true, "timestamp": "2026-10-08T10:00:05Z", "message": obj{"role": "user", "content": canary}}),
		line(obj{"type": "system", "subtype": "compact_boundary", "uuid": "c1", "timestamp": "2026-10-08T10:00:06Z"}),
		{"type": "cost-state", "sessionId": "S1", "totalCostUSD": 0.5, "totalDuration": 6000},
	} {
		main.WriteString(jl(t, l))
	}
	main.WriteString("not json " + canary + "\n")
	os.WriteFile(filepath.Join(dir, "S1.jsonl"), []byte(main.String()), 0o600)

	sub := jl(t, obj{"type": "assistant", "sessionId": "S1", "agentId": "A1", "cwd": root, "uuid": "x1", "timestamp": "2026-10-08T10:00:02Z",
		"message": obj{"id": "m9", "model": "claude-sonnet-5-5", "usage": usage(1, 2), "content": []obj{
			{"type": "tool_use", "id": "t9", "name": "Bash", "input": obj{"command": "go vet ./..."}}}}})
	os.WriteFile(filepath.Join(dir, "S1", "subagents", "agent-A1.jsonl"), []byte(sub), 0o600)
	os.WriteFile(filepath.Join(dir, "S1", "subagents", "agent-A1.meta.json"), []byte(`{"agentType":"verifier","description":"`+canary+`"}`), 0o600)
	return home, root
}

func ingestT(t *testing.T, s *Store, o Options) Result {
	t.Helper()
	r, err := Ingest(s, o, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestClaudeIngestShape(t *testing.T) {
	home, root := claudeFixture(t)
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	r := ingestT(t, s, Options{Home: home, HomeGiven: true, Vendor: "claude"})
	if r.Sessions != 2 || len(r.Failed) != 0 || r.Skipped != 1 {
		t.Fatalf("result %+v, want 2 sessions, 0 failed, 1 skipped", r)
	}

	var cost float64
	var dur int64
	var proj, model string
	s.db.QueryRow(`SELECT cost_usd, duration_ms, project_root, model FROM sessions WHERE id = 'S1'`).Scan(&cost, &dur, &proj, &model)
	if cost != 0.5 || dur != 6000 || proj != root || model != "claude-opus-5-5" {
		t.Fatalf("S1 = %v %v %q %q", cost, dur, proj, model)
	}
	var parent, agent string
	s.db.QueryRow(`SELECT parent_id, agent_type FROM sessions WHERE id = 'A1'`).Scan(&parent, &agent)
	if parent != "S1" || agent != "verifier" {
		t.Fatalf("A1 parent %q type %q", parent, agent)
	}

	q := func(query string) int {
		var n int
		if err := s.db.QueryRow(query).Scan(&n); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
		return n
	}
	checks := map[string]int{
		`SELECT count(*) FROM steps WHERE session_id='S1' AND kind='assistant'`:                          1,
		`SELECT tok_out FROM steps WHERE ext_id='msg:m1'`:                                                40,
		`SELECT count(*) FROM steps WHERE session_id='S1' AND kind='tool'`:                               4,
		`SELECT count(*) FROM steps WHERE kind='tool' AND is_error=1`:                                    2,
		`SELECT count(*) FROM steps WHERE kind='tool' AND denied=1`:                                      1,
		`SELECT count(*) FROM steps WHERE kind='interrupt'`:                                              1,
		`SELECT count(*) FROM steps WHERE kind='user'`:                                                   1,
		`SELECT count(*) FROM steps WHERE kind='compaction'`:                                             1,
		`SELECT count(*) FROM steps WHERE skill='feature'`:                                               1,
		`SELECT count(*) FROM steps WHERE target_kind='handoff'`:                                         1,
		`SELECT count(*) FROM step_verbs WHERE verb IN ('git commit -m', 'go test', 'rm -rf', 'go vet')`: 4,
		`SELECT count(*) FROM steps WHERE session_id='A1' AND kind='tool'`:                               1,
	}
	for query, want := range checks {
		if got := q(query); got != want {
			t.Errorf("%s = %d, want %d", query, got, want)
		}
	}

	// Re-ingesting adds nothing.
	if r2 := ingestT(t, s, Options{Home: home, HomeGiven: true, Vendor: "claude"}); r2.Sessions+r2.Steps+r2.Events != 0 {
		t.Fatalf("second ingest added %+v", r2)
	}
}

// No content reaches the store: the canary sits in a prompt, a reply, a
// command, a Skill's args, a path, every tool result and a subagent's
// description, and is absent from every byte the store wrote.
func TestClaudeIngestStoresNoContent(t *testing.T) {
	home, _ := claudeFixture(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "metrics.db")
	s := openT(t, p)
	ingestT(t, s, Options{Home: home, HomeGiven: true, Vendor: "claude"})
	exp := filepath.Join(dir, "export.jsonl")
	if err := s.Export(exp); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{p, p + "-wal", exp} {
		b, err := os.ReadFile(f)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if bytes.Contains(b, []byte(canary)) {
			t.Errorf("%s holds the canary", filepath.Base(f))
		}
	}
}

// A replaced (shorter) transcript is read again from the start and adds no
// duplicate rows; a growing one adds only its new lines; a partial last line
// waits.
func TestClaudeIngestIncremental(t *testing.T) {
	home, root := claudeFixture(t)
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	o := Options{Home: home, HomeGiven: true, Vendor: "claude"}
	ingestT(t, s, o)
	f := filepath.Join(home, ".claude", "projects", "-proj", "S1.jsonl")

	next := jl(t, obj{"type": "assistant", "sessionId": "S1", "cwd": root, "uuid": "a9", "timestamp": "2026-10-08T10:01:00Z",
		"message": obj{"id": "m2", "content": []obj{{"type": "tool_use", "id": "t5", "name": "Grep", "input": obj{}}}}})
	fh, _ := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0)
	fh.WriteString(next[:len(next)/2])
	fh.Close()
	if r := ingestT(t, s, o); r.Steps != 0 {
		t.Fatalf("a partial line was ingested: %+v", r)
	}
	fh, _ = os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0)
	fh.WriteString(next[len(next)/2:])
	fh.Close()
	if r := ingestT(t, s, o); r.Steps != 2 {
		t.Fatalf("growth added %d steps, want 2 (message + tool)", r.Steps)
	}

	os.WriteFile(f, []byte(next), 0o600) // replaced by a shorter file
	if r := ingestT(t, s, o); r.Steps != 0 || r.Sessions != 0 {
		t.Fatalf("re-read of a replaced file added %+v", r)
	}
}

// --home and --source win over CLAUDE_CONFIG_DIR; the variable applies only
// when neither flag is given (D9, R6).
func TestClaudeRootRule(t *testing.T) {
	decoy := t.TempDir()
	env := func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return decoy
		}
		return ""
	}
	cases := []struct {
		o    Options
		want string
	}{
		{Options{Home: "/h", HomeGiven: true, Getenv: env}, "/h/.claude/projects"},
		{Options{Home: "/h", Source: "/s", Getenv: env}, "/s/.claude/projects"},
		{Options{Home: "/h", Getenv: env}, filepath.Join(decoy, "projects")},
		{Options{Home: "/h", Getenv: func(string) string { return "" }}, "/h/.claude/projects"},
	}
	for _, c := range cases {
		if got := c.o.claudeRoot(); got != c.want {
			t.Errorf("claudeRoot(%+v) = %q, want %q", c.o, got, c.want)
		}
	}
}

func TestSourceLabelAndOutcome(t *testing.T) {
	home, _ := claudeFixture(t)
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	ingestT(t, s, Options{Home: "/nonexistent", Source: home, Label: "arm=lean", Outcome: "pass", Vendor: "claude"})
	var n int
	s.db.QueryRow(`SELECT count(*) FROM sessions WHERE label='arm=lean' AND outcome='pass'`).Scan(&n)
	if n != 2 {
		t.Fatalf("labelled sessions = %d, want 2", n)
	}
	for _, o := range []Options{{Source: "x"}, {Label: "nokv"}, {Outcome: "maybe", Source: "x", Label: "a=b"}, {Outcome: "pass"}, {Label: "a=b"}, {Vendor: "gpt"}} {
		if o.Validate() == nil {
			t.Errorf("Validate(%+v) accepted", o)
		}
	}
}

// The kit ingests a run first and learns its outcome later: a second ingest
// with --outcome stamps the label's sessions that have none, even though no
// transcript changed, and leaves an earlier run's outcome alone.
func TestOutcomeStampedAfterTheRun(t *testing.T) {
	run1, _ := claudeFixture(t)
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	ingestT(t, s, Options{Source: run1, Label: "arm=x", Outcome: "fail", Vendor: "claude"})

	run2, _ := claudeFixture(t)
	// a second run of the same arm: other session ids
	for _, f := range []string{filepath.Join(run2, ".claude", "projects", "-proj", "S1.jsonl"),
		filepath.Join(run2, ".claude", "projects", "-proj", "S1", "subagents", "agent-A1.jsonl")} {
		b, _ := os.ReadFile(f)
		os.WriteFile(f, bytes.ReplaceAll(bytes.ReplaceAll(b, []byte(`"S1"`), []byte(`"R2"`)), []byte(`"A1"`), []byte(`"B2"`)), 0o600)
	}
	ingestT(t, s, Options{Source: run2, Label: "arm=x", Vendor: "claude"})
	ingestT(t, s, Options{Source: run2, Label: "arm=x", Outcome: "pass", Vendor: "claude"})
	got := map[string]string{}
	rows, _ := s.db.Query(`SELECT id, COALESCE(outcome, '') FROM sessions WHERE parent_id IS NULL`)
	for rows.Next() {
		var id, o string
		rows.Scan(&id, &o)
		got[id] = o
	}
	rows.Close()
	if got["S1"] != "fail" || got["R2"] != "pass" {
		t.Fatalf("outcomes %v, want S1 fail (kept) and R2 pass (stamped late)", got)
	}
}
