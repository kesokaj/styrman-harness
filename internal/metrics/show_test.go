package metrics

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestShowBothCLIs: the Claude and Antigravity fixtures in one store give
// per-vendor numbers, and what Antigravity cannot record is null, not 0.
func TestShowBothCLIs(t *testing.T) {
	chome, _ := claudeFixture(t)
	ahome, _ := agyFixture(t)
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	ingestT(t, s, Options{Home: chome, HomeGiven: true, Vendor: "claude"})
	ingestT(t, s, Options{Home: ahome, HomeGiven: true, Vendor: "agy"})

	sum, err := s.Show(ShowOptions{Since: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Now: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := sum.WriteJSON(&buf); err != nil {
		t.Fatal(err)
	}
	var j map[string]any
	if err := json.Unmarshal(buf.Bytes(), &j); err != nil {
		t.Fatal(err)
	}
	get := func(path string) any {
		var v any = j
		for _, k := range strings.Split(path, ".") {
			m, ok := v.(map[string]any)
			if !ok {
				return "<missing>"
			}
			v = m[k]
		}
		return v
	}
	want := map[string]any{
		"sessions.claude":                       1.0,
		"sessions.agy":                          1.0,
		"subagents.claude.verifier":             1.0,
		"subagents.agy.explorer":                1.0,
		"B.tool_calls.claude":                   5.0, // 4 main + 1 subagent
		"B.tool_calls.agy":                      4.0,
		"B.tool_calls_by_tool.claude.Bash":      3.0,
		"B.tool_calls_by_tool.agy.Bash":         1.0,
		"B.cost_usd.claude":                     0.5,
		"B.cost_usd.agy":                        nil,
		"B.compactions.claude":                  1.0,
		"B.compactions.agy":                     nil,
		"B.peak_context.agy":                    42000.0,
		"C.skills.claude.feature":               1.0,
		"C.skills.agy.golang":                   1.0,
		"C.forbidden.commit chained to test":    1.0,
		"C.first_reads.handoff":                 1.0,
		"C.first_reads.plan":                    1.0,
		"D.denials.claude":                      1.0,
		"D.denials.agy":                         nil,
		"D.interrupts.claude":                   1.0,
		"D.subagent_tool_calls.claude.verifier": 1.0,
	}
	for k, w := range want {
		if g := get(k); g != w {
			t.Errorf("%s = %v, want %v", k, g, w)
		}
	}
	for _, v := range []string{"claude", "agy"} {
		if sum.B.ActiveMin[v] > sum.B.WallMin[v]+1e-9 {
			t.Errorf("%s active %.3f min > wall %.3f min: both must count main sessions only", v, sum.B.ActiveMin[v], sum.B.WallMin[v])
		}
	}
	if sum.C.Forbidden["git add ."] != 1 { // a key with a dot: read directly
		t.Errorf("forbidden git add . = %d, want 1", sum.C.Forbidden["git add ."])
	}
	if r := sum.D.ToolErrorRate["agy"]; r == nil || *r != 0.25 {
		t.Errorf("agy tool error rate = %v, want 0.25", r)
	}
	if !strings.Contains(strings.Join(sum.Notes, "\n"), "agy: cost is not exposed") {
		t.Errorf("notes do not name agy's missing cost: %q", sum.Notes)
	}
	var text bytes.Buffer
	sum.WriteText(&text)
	for _, w := range []string{"A. outcome", "B. efficiency", "C. process adherence", "D. reliability", "n/a"} {
		if !strings.Contains(text.String(), w) {
			t.Errorf("text output lacks %q", w)
		}
	}
}

// Labelled sandbox sessions are out of the default view, in with their label
// or --all-labels, and give pass^k.
func TestShowLabels(t *testing.T) {
	chome, _ := claudeFixture(t)
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	ingestT(t, s, Options{Source: chome, Label: "arm=lean", Outcome: "pass", Vendor: "claude"})
	o := ShowOptions{Since: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Now: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)}
	if sum, _ := s.Show(o); sum.Sessions["claude"] != 0 {
		t.Errorf("default view shows labelled sessions: %v", sum.Sessions)
	}
	o.Label = "arm=lean"
	sum, _ := s.Show(o)
	if sum.Sessions["claude"] != 1 || sum.A.PassK["arm=lean"] != (PassK{K: 1, Passed: 1, PassK: 1}) {
		t.Errorf("label view: sessions %v pass_k %v", sum.Sessions, sum.A.PassK)
	}
}

// Wall and active time count main sessions only: a subagent's steps run
// inside its parent's wall time and would be counted twice (spec B2).
func TestShowActiveTimeMainOnly(t *testing.T) {
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	base := time.Date(2026, 10, 8, 10, 0, 0, 0, time.UTC).UnixMilli()
	min := int64(60000)
	for _, q := range []struct {
		sql  string
		args []any
	}{
		{`INSERT INTO sessions(id, vendor, started_at, ended_at) VALUES ('M', 'agy', ?, ?)`, []any{base, base + 2*min}},
		{`INSERT INTO sessions(id, vendor, parent_id, agent_type, started_at, ended_at) VALUES ('S', 'agy', 'M', 'explorer', ?, ?)`, []any{base, base + 4*min}},
		{`INSERT INTO steps(session_id, ext_id, ts, kind) VALUES ('M', 'm1', ?, 'assistant'), ('M', 'm2', ?, 'assistant')`, []any{base, base + 2*min}},
		{`INSERT INTO steps(session_id, ext_id, ts, kind) VALUES ('S', 's1', ?, 'assistant'), ('S', 's2', ?, 'assistant')`, []any{base, base + 4*min}},
	} {
		if _, err := s.db.Exec(q.sql, q.args...); err != nil {
			t.Fatal(err)
		}
	}
	sum, err := s.Show(ShowOptions{Since: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Now: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if sum.B.WallMin["agy"] != 2 || sum.B.ActiveMin["agy"] != 2 {
		t.Fatalf("wall %.1f active %.1f, want 2 and 2", sum.B.WallMin["agy"], sum.B.ActiveMin["agy"])
	}
}

// A resumed or forked Claude session replays earlier history into a new
// transcript: a tool call or a message seen under one id is counted once,
// whichever sessions repeat it (release verification of v2.6.0).
func TestShowCountsReplayedHistoryOnce(t *testing.T) {
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	base := time.Date(2026, 10, 8, 10, 0, 0, 0, time.UTC).UnixMilli()
	for _, q := range []string{
		`INSERT INTO sessions(id, vendor, started_at, ended_at) VALUES ('A', 'claude', ?, ?), ('B', 'claude', ?, ?)`,
	} {
		if _, err := s.db.Exec(q, base, base+60000, base+120000, base+180000); err != nil {
			t.Fatal(err)
		}
	}
	for _, sid := range []string{"A", "B"} { // B replays A's history
		s.db.Exec(`INSERT INTO steps(session_id, ext_id, ts, kind, tool) VALUES (?, 't1', ?, 'tool', 'Bash')`, sid, base+1000)
		s.db.Exec(`INSERT INTO steps(session_id, ext_id, ts, kind, tok_in, tok_out) VALUES (?, 'msg:m1', ?, 'assistant', 100, 10)`, sid, base+1000)
	}
	s.db.Exec(`INSERT INTO steps(session_id, ext_id, ts, kind, tool) VALUES ('B', 't2', ?, 'tool', 'Read')`, base+130000)
	sum, err := s.Show(ShowOptions{Since: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Now: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if sum.B.ToolCalls["claude"] != 2 || sum.B.Tokens["claude"].In != 100 {
		t.Fatalf("tool calls %d, tokens in %d; want 2 and 100 (replayed history counted once)", sum.B.ToolCalls["claude"], sum.B.Tokens["claude"].In)
	}
}

// Replayed history is deduplicated within one label only: a sandbox
// transcript that reuses a real session's ids never hides the real session's
// steps under --all-labels (security review 6118921..c478d86, finding 3).
func TestShowReplayDedupeStaysWithinALabel(t *testing.T) {
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	base := time.Date(2026, 10, 8, 10, 0, 0, 0, time.UTC).UnixMilli()
	// The sandbox session starts first, so without the label in the key it
	// would claim the shared ids and the real session's steps would vanish.
	if _, err := s.db.Exec(`INSERT INTO sessions(id, vendor, label, started_at, ended_at) VALUES ('sbx', 'claude', 'arm=x', ?, ?)`, base, base+60000); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(`INSERT INTO sessions(id, vendor, started_at, ended_at) VALUES ('real', 'claude', ?, ?)`, base+120000, base+180000); err != nil {
		t.Fatal(err)
	}
	for _, sid := range []string{"sbx", "real"} {
		s.db.Exec(`INSERT INTO steps(session_id, ext_id, ts, kind, tool) VALUES (?, 't1', ?, 'tool', 'Bash')`, sid, base+1000)
	}
	sum, err := s.Show(ShowOptions{AllLabels: true, Since: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), Now: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if sum.B.ToolCalls["claude"] != 2 {
		t.Fatalf("tool calls %d; want 2 (one per label)", sum.B.ToolCalls["claude"])
	}
}

// A session belongs to the window it starts in, [since, until): two adjacent
// windows never count the same session, even one starting at the boundary
// or running across it (release review of v2.6.0, round 7).
func TestShowAdjacentWindowsShareNoSession(t *testing.T) {
	s := openT(t, filepath.Join(t.TempDir(), "metrics.db"))
	day := func(d, h int) int64 { return time.Date(2026, 10, d, h, 0, 0, 0, time.UTC).UnixMilli() }
	s.db.Exec(`INSERT INTO sessions(id, vendor, started_at, ended_at) VALUES
		('across', 'claude', ?, ?), ('boundary', 'claude', ?, ?), ('early', 'claude', ?, ?)`,
		day(4, 23), day(5, 1), day(5, 0), day(5, 2), day(3, 10), day(3, 11))
	split := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	before, _ := s.Show(ShowOptions{Since: split.AddDate(0, 0, -3), Now: split})
	after, _ := s.Show(ShowOptions{Since: split, Now: split.AddDate(0, 0, 3)})
	if before.Sessions["claude"] != 2 || after.Sessions["claude"] != 1 {
		t.Fatalf("before %d, after %d; want 2 (early, across) and 1 (boundary)", before.Sessions["claude"], after.Sessions["claude"])
	}
}
