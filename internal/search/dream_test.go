package search

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestListProjectsMissingDBIsEmpty(t *testing.T) {
	home := t.TempDir()
	got, err := ListProjects(home)
	if err != nil || len(got) != 0 {
		t.Fatalf("got %v, %v; want none, nil", got, err)
	}
	if _, err := os.Stat(globalDB(home)); !os.IsNotExist(err) {
		t.Error("ListProjects created the global db")
	}
}

func TestListProjects(t *testing.T) {
	home := t.TempDir()
	path := globalDB(home)
	os.MkdirAll(filepath.Dir(path), 0o700)
	db, err := sql.Open("sqlite", dsn(path))
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`CREATE TABLE projects(root TEXT PRIMARY KEY, name TEXT, docs_root TEXT, last_seen INTEGER)`,
		`INSERT INTO projects VALUES('/p/one','one','.agents/docs',1),('/p/two','two','.agents/docs',2)`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()
	os.Remove(path + "-wal")
	os.Remove(path + "-shm")
	before := listDir(t, filepath.Dir(path))
	got, err := ListProjects(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Root != "/p/one" || got[1].Name != "two" || got[1].DocsRoot != ".agents/docs" {
		t.Errorf("got %+v", got)
	}
	// A read-only open of a WAL database creates -wal and -shm beside it;
	// ListProjects reads a copy, so a dry run leaves the directory as it was.
	if after := listDir(t, filepath.Dir(path)); after != before {
		t.Errorf("ListProjects changed %s:\nbefore %s\nafter  %s", filepath.Dir(path), before, after)
	}
}

func listDir(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return strings.Join(names, " ")
}

// agyFixture writes the two databases Antigravity keeps, in the shapes the
// D8 probe found: summaries map conversations to workspaces, and each
// conversation's steps hold protobuf blobs with the text as UTF-8 runs.
func agyFixture(t *testing.T, home string, convs map[string]struct {
	workspace string
	modified  string
	payloads  [][]byte
}) {
	t.Helper()
	dir := filepath.Join(home, ".gemini", "antigravity-cli")
	os.MkdirAll(filepath.Join(dir, "conversations"), 0o755)
	sum, err := sql.Open("sqlite", dsn(filepath.Join(dir, "conversation_summaries.db")))
	if err != nil {
		t.Fatal(err)
	}
	defer sum.Close()
	if _, err := sum.Exec("CREATE TABLE conversation_summaries (`conversation_id` text, `last_modified_time` datetime NOT NULL, `workspace_uris` text NOT NULL)"); err != nil {
		t.Fatal(err)
	}
	for id, c := range convs {
		if _, err := sum.Exec(`INSERT INTO conversation_summaries VALUES(?,?,?)`, id, c.modified, `["file://`+c.workspace+`"]`); err != nil {
			t.Fatal(err)
		}
		db, err := sql.Open("sqlite", dsn(filepath.Join(dir, "conversations", id+".db")))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec("CREATE TABLE steps (`idx` integer, `step_type` integer, `step_payload` blob)"); err != nil {
			t.Fatal(err)
		}
		for i, p := range c.payloads {
			// Step types as the D8 probe found them: 14 the user's input, 15
			// the model's reply, 132 a tool's output.
			typ := 15
			if i == 0 {
				typ = 14
			}
			if bytes.Contains(p, []byte("TOOL OUTPUT")) {
				typ = 132
			}
			if _, err := db.Exec(`INSERT INTO steps VALUES(?,?,?)`, i, typ, p); err != nil {
				t.Fatal(err)
			}
		}
		db.Close()
	}
}

func blob(parts ...string) []byte {
	var b []byte
	for _, p := range parts {
		b = append(b, 0x0a, 0x8f, 0x01)
		b = append(b, p...)
		b = append(b, 0x12, 0x00)
	}
	return b
}

func TestAgyConversationText(t *testing.T) {
	home := t.TempDir()
	type conv = struct {
		workspace string
		modified  string
		payloads  [][]byte
	}
	agyFixture(t, home, map[string]conv{
		"new-here": {"/p/one", "2026-09-30 16:14:07.567153506+00:00", [][]byte{
			blob("b$7863135a-0a35-4648-9cef-f196af3f2197", "Remember that we deploy on Fridays please."),
			blob("Remember that we deploy on Fridays please.", "Saved: deploys happen on Fridays, noted for later."),
			blob("TOOL OUTPUT: the contents of a private file read by a tool"),
		}},
		"old-here":  {"/p/one", "2026-09-01 10:00:00.000000000+00:00", [][]byte{blob("this is an old conversation text")}},
		"new-other": {"/p/two", "2026-09-30 16:00:00.000000000+00:00", [][]byte{blob("another project entirely, not ours")}},
		"new-under": {"/p/one/sub", "2026-09-30 15:00:00.000000000+00:00", [][]byte{blob("started in a subdirectory of the root")}},
	})
	since := time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)
	got, err := AgyConversationText(home, "/p/one", since)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, c := range got {
		ids[c.ID] = c.Text
	}
	if len(ids) != 2 || ids["new-here"] == "" || ids["new-under"] == "" {
		t.Fatalf("conversations %v, want new-here and new-under only", ids)
	}
	text := ids["new-here"]
	if !strings.Contains(text, "Remember that we deploy on Fridays please.") || !strings.Contains(text, "Saved: deploys happen on Fridays") {
		t.Errorf("text %q lacks the prompt or the reply", text)
	}
	if strings.Count(text, "Remember that we deploy") != 1 {
		t.Errorf("text %q repeats a run: dedup failed", text)
	}
	if strings.Contains(text, "private file read by a tool") {
		t.Errorf("text %q carries a tool's output (security review 2026-09-30)", text)
	}
	if strings.Contains(text, "7863135a") {
		t.Errorf("text %q keeps an id run", text)
	}
	// home had no ~/.config/sh-agent before: the run leaves none behind.
	if _, err := os.Stat(filepath.Join(home, ".config", "sh-agent")); !os.IsNotExist(err) {
		t.Errorf("~/.config/sh-agent left behind: %v", err)
	}
}

func TestAgyConversationTextNoAntigravity(t *testing.T) {
	got, err := AgyConversationText(t.TempDir(), "/p/one", time.Time{})
	if err != nil || len(got) != 0 {
		t.Fatalf("got %v, %v; want none, nil", got, err)
	}
}

func TestAgyTime(t *testing.T) {
	want := time.Date(2026, 9, 30, 16, 14, 7, 567153506, time.UTC)
	for name, v := range map[string]any{
		"time":    want,
		"string":  "2026-09-30 16:14:07.567153506+00:00",
		"bytes":   []byte("2026-09-30 16:14:07.567153506+00:00"),
		"rfc3339": "2026-09-30T16:14:07.567153506Z",
	} {
		if got, ok := agyTime(v); !ok || !got.Equal(want) {
			t.Errorf("%s: %v, %v", name, got, ok)
		}
	}
	for _, v := range []any{"yesterday", 42, nil} {
		if _, ok := agyTime(v); ok {
			t.Errorf("%v parsed", v)
		}
	}
}

func TestInWorkspace(t *testing.T) {
	for uris, want := range map[string]bool{
		`["file:///p/one"]`:             true,
		`["file:///p/one/sub"]`:         true,
		`["file:///p/one-other"]`:       false,
		`["https:///p/one"]`:            false,
		`not json`:                      false,
		`["file:///x","file:///p/one"]`: true,
	} {
		if got := inWorkspace(uris, "/p/one"); got != want {
			t.Errorf("%s: %v, want %v", uris, got, want)
		}
	}
}
