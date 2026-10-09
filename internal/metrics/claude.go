package metrics

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/search"
	"github.com/kesokaj/styrman-harness/internal/trail"
)

// deniedPrefixes open a tool result Claude Code writes when a permission
// check refused the call (probe 2026-10-08; spec § Invariants). Only the
// flag is stored, never the text.
var deniedPrefixes = []string{
	"Permission for this action was denied",
	"Permission for this command was denied",
	"The user doesn't want to proceed with this tool use",
}

const interruptPrefix = "[Request interrupted by user"

// pathTools are the tools whose file_path is classified into a trail kind;
// the path itself is never stored.
var pathTools = map[string]bool{"Read": true, "Edit": true, "Write": true, "NotebookEdit": true}

// claudeLine is the part of a transcript line the store reads. Content stays
// raw until a block's type says what may be taken from it.
type claudeLine struct {
	Type      string          `json:"type"`
	Subtype   string          `json:"subtype"`
	UUID      string          `json:"uuid"`
	SessionID string          `json:"sessionId"`
	AgentID   string          `json:"agentId"`
	Cwd       string          `json:"cwd"`
	Timestamp string          `json:"timestamp"`
	IsMeta    bool            `json:"isMeta"`
	Message   json.RawMessage `json:"message"`
	CostUSD   *float64        `json:"totalCostUSD"`
	Duration  *int64          `json:"totalDuration"`
}

type claudeMessage struct {
	ID      string          `json:"id"`
	Model   string          `json:"model"`
	Content json.RawMessage `json:"content"`
	Usage   *struct {
		In         int64 `json:"input_tokens"`
		Out        int64 `json:"output_tokens"`
		CacheRead  int64 `json:"cache_read_input_tokens"`
		CacheWrite int64 `json:"cache_creation_input_tokens"`
	} `json:"usage"`
}

type claudeBlock struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Input     json.RawMessage `json:"input"`
	ToolUseID string          `json:"tool_use_id"`
	IsError   bool            `json:"is_error"`
	Content   json.RawMessage `json:"content"`
	Text      string          `json:"text"`
}

// toolInput is the only part of a tool's input the store reads: a Bash
// command (normalised, hashed, never stored), a Skill's name and a file
// path (classified, never stored). D22.
type toolInput struct {
	Command  string `json:"command"`
	Skill    string `json:"skill"`
	FilePath string `json:"file_path"`
	Notebook string `json:"notebook_path"`
}

// ingestClaude walks <root>/<project>/<session>.jsonl and
// <root>/<project>/<session>/subagents/agent-*.jsonl.
func (s *Store) ingestClaude(o Options, r *Result, log io.Writer) {
	root := o.claudeRoot()
	projects, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return
	}
	if err != nil {
		r.Failed = append(r.Failed, root)
		fmt.Fprintln(log, printable(fmt.Sprintf("metrics: %s: %v", root, err)))
		return
	}
	for _, p := range projects {
		if !p.IsDir() {
			continue
		}
		dir := filepath.Join(root, p.Name())
		files, _ := filepath.Glob(filepath.Join(dir, "*.jsonl"))
		subs, _ := filepath.Glob(filepath.Join(dir, "*", "subagents", "agent-*.jsonl"))
		for _, f := range append(files, subs...) {
			if err := s.ingestClaudeFile(f, o, r); err != nil {
				r.Failed = append(r.Failed, f)
				fmt.Fprintln(log, printable(fmt.Sprintf("metrics: %s: %v", f, err)))
			}
		}
	}
}

// subagentMeta is <agent>.meta.json beside a subagent transcript.
type subagentMeta struct {
	AgentType string `json:"agentType"`
}

// ingestClaudeFile reads path from its stored offset in one transaction.
// Only whole lines are consumed: a line still being written is read next
// time. A file shorter than its offset was replaced and is read from 0; the
// primary keys absorb what was already stored.
func (s *Store) ingestClaudeFile(path string, o Options, r *Result) error {
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("not a regular file")
	}
	var size, mtime, offset int64
	err = s.db.QueryRow(`SELECT size, mtime, offset FROM ingest_state WHERE source = ?`, path).Scan(&size, &mtime, &offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if size == fi.Size() && mtime == fi.ModTime().UnixMilli() && offset == fi.Size() {
		return nil
	}
	if fi.Size() < offset {
		offset = 0
	}
	f, ofi, err := openRegular(path)
	if err != nil {
		return err
	}
	fi = ofi // sizes and times from the file actually opened
	defer f.Close()
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	sub := strings.Contains(filepath.ToSlash(path), "/subagents/")
	agentType := ""
	if sub {
		var m subagentMeta
		if b, ok := readSmallRegular(strings.TrimSuffix(path, ".jsonl")+".meta.json", 64<<10); ok {
			json.Unmarshal(b, &m)
		}
		agentType = identifier(m.AgentType)
	}

	// The salt is read before Begin: the store has one connection, which the
	// transaction then holds.
	salt := s.salt()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	w := writer{tx: tx, salt: salt, label: o.Label, outcome: o.Outcome}

	br := bufio.NewReaderSize(f, 1<<20)
	consumed := offset
	for {
		line, n, tooLong, err := readLine(br)
		if err == io.EOF {
			break // a partial last line waits for the next ingest
		}
		if err != nil {
			return err
		}
		consumed += n
		if tooLong {
			r.Skipped++ // never held in memory whole
			continue
		}
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		if !w.claudeLine(line, sub, agentType) {
			r.Skipped++
		}
		if w.err != nil {
			return w.err
		}
	}
	if _, err := tx.Exec(`INSERT INTO ingest_state(source, size, mtime, offset) VALUES (?, ?, ?, ?)
		ON CONFLICT(source) DO UPDATE SET size = excluded.size, mtime = excluded.mtime, offset = excluded.offset`,
		path, fi.Size(), fi.ModTime().UnixMilli(), consumed); err != nil {
		return err
	}
	return tx.Commit()
}

// writer turns parsed lines into rows inside one transaction; the first
// error stops it.
type writer struct {
	tx             *sql.Tx
	salt           string
	label, outcome string
	err            error
	project        map[string]string // project root -> docs root
}

func (w *writer) exec(q string, args ...any) {
	if w.err != nil {
		return
	}
	_, w.err = w.tx.Exec(q, args...)
}

// claudeLine stores one line's shape and reports whether the line was
// understood.
func (w *writer) claudeLine(raw []byte, sub bool, agentType string) bool {
	var l claudeLine
	if err := json.Unmarshal(raw, &l); err != nil || l.Type == "" {
		return false
	}
	sid, parent := l.SessionID, ""
	if sub {
		if l.AgentID == "" {
			return false
		}
		sid, parent = l.AgentID, l.SessionID
	}
	if sid == "" {
		return true // a line with no session (a queue operation): nothing to store
	}
	ts := parseTS(l.Timestamp)
	w.session(sid, parent, agentType, l.Cwd, ts)

	switch l.Type {
	case "cost-state":
		if l.CostUSD != nil {
			w.exec(`UPDATE sessions SET cost_usd = ? WHERE id = ?`, *l.CostUSD, sid)
		}
		if l.Duration != nil {
			w.exec(`UPDATE sessions SET duration_ms = ? WHERE id = ?`, *l.Duration, sid)
		}
	case "system":
		if l.Subtype == "compact_boundary" && l.UUID != "" {
			w.step(sid, l.UUID, ts, "compaction")
		}
	case "assistant":
		var m claudeMessage
		if json.Unmarshal(l.Message, &m) != nil {
			return false
		}
		w.assistant(sid, ts, m, l.Cwd)
	case "user":
		if l.IsMeta {
			return true
		}
		var m claudeMessage
		if json.Unmarshal(l.Message, &m) != nil {
			return false
		}
		w.user(sid, l.UUID, ts, m)
	}
	return true
}

// session creates the session row on first sight and widens its window.
func (w *writer) session(id, parent, agentType, cwd string, ts int64) {
	var tsArg any
	if ts > 0 {
		tsArg = ts
	}
	w.exec(`INSERT INTO sessions(id, vendor, project_root, parent_id, agent_type, label, outcome, started_at, ended_at)
		VALUES (?, 'claude', NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			project_root = COALESCE(sessions.project_root, excluded.project_root),
			parent_id = COALESCE(sessions.parent_id, excluded.parent_id),
			agent_type = COALESCE(sessions.agent_type, excluded.agent_type),
			label = COALESCE(sessions.label, excluded.label),
			outcome = COALESCE(sessions.outcome, excluded.outcome),
			started_at = MIN(COALESCE(sessions.started_at, excluded.started_at), COALESCE(excluded.started_at, sessions.started_at)),
			ended_at = MAX(COALESCE(sessions.ended_at, excluded.ended_at), COALESCE(excluded.ended_at, sessions.ended_at))`,
		id, cwd, parent, agentType, w.label, w.outcome, tsArg, tsArg)
}

func (w *writer) step(sid, ext string, ts int64, kind string) {
	w.exec(`INSERT INTO steps(session_id, ext_id, ts, kind) VALUES (?, ?, ?, ?)
		ON CONFLICT(session_id, ext_id) DO NOTHING`, sid, ext, ts, kind)
}

// assistant stores one assistant line: the message row (usage kept as the
// maximum over the lines a streamed message repeats) and a row per tool call.
func (w *writer) assistant(sid string, ts int64, m claudeMessage, cwd string) {
	if m.ID != "" {
		var in, out, cr, cw int64
		if m.Usage != nil {
			in, out, cr, cw = m.Usage.In, m.Usage.Out, m.Usage.CacheRead, m.Usage.CacheWrite
		}
		w.exec(`INSERT INTO steps(session_id, ext_id, ts, kind, tok_in, tok_out, cache_read, cache_write)
			VALUES (?, ?, ?, 'assistant', ?, ?, ?, ?)
			ON CONFLICT(session_id, ext_id) DO UPDATE SET
				tok_in = MAX(steps.tok_in, excluded.tok_in), tok_out = MAX(steps.tok_out, excluded.tok_out),
				cache_read = MAX(steps.cache_read, excluded.cache_read), cache_write = MAX(steps.cache_write, excluded.cache_write)`,
			sid, "msg:"+m.ID, ts, in, out, cr, cw)
		if m.Model != "" && identifier(m.Model) != "" {
			w.exec(`UPDATE sessions SET model = COALESCE(model, ?) WHERE id = ?`, m.Model, sid)
		}
	}
	var blocks []claudeBlock
	if json.Unmarshal(m.Content, &blocks) != nil {
		return
	}
	for _, b := range blocks {
		if b.Type != "tool_use" || b.ID == "" {
			continue
		}
		var in toolInput
		json.Unmarshal(b.Input, &in)
		var skill, kind, hash any
		if b.Name == "Skill" && in.Skill != "" {
			if id := identifier(in.Skill); id != "" {
				skill = id
			}
		}
		if pathTools[b.Name] {
			p := in.FilePath
			if p == "" {
				p = in.Notebook
			}
			kind = w.targetKind(cwd, p)
		}
		if b.Name == "Bash" && in.Command != "" {
			hash = w.hash(in.Command)
		}
		w.exec(`INSERT INTO steps(session_id, ext_id, ts, kind, tool, skill, target_kind, cmd_hash)
			VALUES (?, ?, ?, 'tool', ?, ?, ?, ?)
			ON CONFLICT(session_id, ext_id) DO UPDATE SET ts = excluded.ts, kind = 'tool', tool = excluded.tool,
				skill = excluded.skill, target_kind = excluded.target_kind, cmd_hash = excluded.cmd_hash`,
			sid, b.ID, ts, toolName(b.Name), skill, kind, hash)
		if b.Name == "Bash" {
			for n, v := range Verb(in.Command) {
				w.exec(`INSERT INTO step_verbs(session_id, ext_id, n, verb) VALUES (?, ?, ?, ?)
					ON CONFLICT DO NOTHING`, sid, b.ID, n, v)
			}
		}
	}
}

// user stores a typed turn, an interrupt, or the error and denial flags of
// tool results. No text is stored.
func (w *writer) user(sid, uuid string, ts int64, m claudeMessage) {
	var text string
	if json.Unmarshal(m.Content, &text) == nil {
		w.userTurn(sid, uuid, ts, text)
		return
	}
	var blocks []claudeBlock
	if json.Unmarshal(m.Content, &blocks) != nil {
		return
	}
	turn := false
	for _, b := range blocks {
		switch b.Type {
		case "tool_result":
			if b.ToolUseID == "" {
				continue
			}
			denied := 0
			if b.IsError && hasDeniedPrefix(resultText(b.Content)) {
				denied = 1
			}
			w.exec(`INSERT INTO steps(session_id, ext_id, ts, kind, is_error, denied) VALUES (?, ?, ?, 'tool', ?, ?)
				ON CONFLICT(session_id, ext_id) DO UPDATE SET is_error = excluded.is_error, denied = excluded.denied`,
				sid, b.ToolUseID, ts, boolInt(b.IsError), denied)
		case "text":
			if !turn {
				w.userTurn(sid, uuid, ts, b.Text)
				turn = true
			}
		}
	}
}

func (w *writer) userTurn(sid, uuid string, ts int64, text string) {
	if uuid == "" || strings.HasPrefix(text, "<") {
		return // a command echo or a system block, not a typed turn
	}
	kind := "user"
	if strings.HasPrefix(text, interruptPrefix) {
		kind = "interrupt"
	}
	w.step(sid, uuid, ts, kind)
}

func resultText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var parts []struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &parts) == nil && len(parts) > 0 {
		return parts[0].Text
	}
	return ""
}

func hasDeniedPrefix(s string) bool {
	for _, p := range deniedPrefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// targetKind is the trail kind of an absolute path inside the session's
// project, "outside" for one elsewhere; the path is discarded.
func (w *writer) targetKind(root, p string) any {
	if p == "" || root == "" {
		return nil
	}
	rel, err := filepath.Rel(root, p)
	if err != nil || rel == ".." || strings.HasPrefix(rel, "../") {
		return "outside"
	}
	return search.KindOf(filepath.ToSlash(rel), w.docsRoot(root))
}

func (w *writer) docsRoot(root string) string {
	if w.project == nil {
		w.project = map[string]string{}
	}
	if d, ok := w.project[root]; ok {
		return d
	}
	d := trail.DocsRoot(root)
	w.project[root] = d
	return d
}

func (w *writer) hash(cmd string) string {
	h := sha256.Sum256([]byte(w.salt + "\x00" + cmd))
	return hex.EncodeToString(h[:])[:16]
}

func (s *Store) salt() string {
	var v string
	s.db.QueryRow(`SELECT value FROM meta WHERE key = 'salt'`).Scan(&v)
	return v
}

func parseTS(s string) int64 {
	if s == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return 0
	}
	return t.UnixMilli()
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// nameShape is what an identifier taken from a transcript must look like to
// be stored: a skill, a model, an agent type. Anything else (escape bytes,
// prose a prompt-injected session put there) is dropped, so it reaches
// neither the store nor the operator's terminal (security review of v2.6.0,
// part 1).
var nameShape = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:@/-]{0,63}$`)

// identifier is s when it is name-shaped, else "".
func identifier(s string) string {
	if nameShape.MatchString(s) {
		return s
	}
	return ""
}

// readSmallRegular reads path when it is a regular file of at most max
// bytes: a .meta.json under --source could be a FIFO or a link to /dev/zero.
func readSmallRegular(path string, max int64) ([]byte, bool) {
	if fi, err := os.Lstat(path); err != nil || !fi.Mode().IsRegular() {
		return nil, false
	}
	f, fi, err := openRegular(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	if fi.Size() > max {
		return nil, false
	}
	b, err := io.ReadAll(io.LimitReader(f, max))
	return b, err == nil
}

// readLine reads one newline-terminated line, holding at most lineMax bytes
// of it: a longer line is consumed to its newline and reported tooLong, with
// no bytes. n counts every byte consumed. A line with no newline yet is
// io.EOF, read again next time.
func readLine(br *bufio.Reader) (line []byte, n int64, tooLong bool, err error) {
	var buf []byte
	for {
		chunk, err := br.ReadSlice('\n')
		n += int64(len(chunk))
		if !tooLong {
			if len(buf)+len(chunk) > lineMax {
				tooLong, buf = true, nil
			} else {
				buf = append(buf, chunk...)
			}
		}
		switch {
		case err == nil:
			return buf, n, tooLong, nil
		case errors.Is(err, bufio.ErrBufferFull):
			continue
		case err == io.EOF && tooLong:
			// An unterminated tail already over the cap is passed over, not
			// re-read on every ingest; whatever is appended later starts a
			// new (malformed, skipped) line.
			return nil, n, true, nil
		default:
			return nil, n, tooLong, err
		}
	}
}

// lineMax is the longest transcript line read; a longer one is skipped and
// counted. A variable so a test can lower it.
var lineMax = 16 << 20

// toolName is a tool's name as stored: name-shaped, or "unnamed".
func toolName(s string) string {
	if id := identifier(s); id != "" {
		return id
	}
	return "unnamed"
}
