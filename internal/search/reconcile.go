package search

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

// maxFileBytes is the largest file indexed; larger files, and files with a
// NUL byte in their first 8 KiB, are skipped.
const (
	maxFileBytes = 1 << 20
	binarySniff  = 8 << 10
	gcAge        = 14 * 24 * 60 * 60 // seconds an unreferenced blob is kept
	touchEvery   = 24 * 60 * 60      // seconds between last_seen refreshes
)

// entry is one path of the current view.
type entry struct {
	path, oid, state string
}

func git(dir string, stdin []byte, args ...string) ([]byte, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	// Git's messages in its C locale: dropReason matches "not a git
	// repository", which a translated catalogue would reword.
	cmd.Env = append(os.Environ(), "LC_ALL=C", "LANGUAGE=")
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return out, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// checkoutRoot returns the top of the git work tree holding dir.
func checkoutRoot(dir string) (string, error) {
	out, err := git(dir, nil, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", errors.New("no work tree")
	}
	return root, nil
}

// docsRootOf is the trail's root: .agents/docs when it exists, else docs,
// else .agents/docs.
func docsRootOf(root string) string {
	for _, d := range []string{".agents/docs", "docs"} {
		if fi, err := os.Stat(filepath.Join(root, d)); err == nil && fi.IsDir() {
			return d
		}
	}
	return ".agents/docs"
}

// gitView lists the checkout's current files: the index entries of
// `git ls-files -s`, overridden by what `git status` reports modified or
// untracked (their blob IDs computed from their bytes), minus deleted paths,
// symlinks, submodules and secret-shaped names. A file's content is read
// again when its blob is indexed: a file edited and restored between the two
// reads is indexed with the edit until `index --rebuild` (a known limit; the
// operator's choice at the security review's cap, plan D35).
// prefix limits it to one directory.
func gitView(root, prefix string) ([]entry, error) {
	under := func(p string) bool { return prefix == "" || strings.HasPrefix(p, prefix) }
	out, err := git(root, nil, "ls-files", "-s", "-z")
	if err != nil {
		return nil, err
	}
	var order []string
	view := map[string]entry{}
	for _, rec := range strings.Split(string(out), "\x00") {
		meta, p, ok := strings.Cut(rec, "\t")
		if !ok || !under(p) {
			continue
		}
		f := strings.Fields(meta) // mode oid stage
		if len(f) != 3 || f[0] == "120000" || f[0] == "160000" {
			continue
		}
		if _, dup := view[p]; dup { // an unmerged path: status decides below
			continue
		}
		view[p] = entry{path: p, oid: f[1], state: "tracked"}
		order = append(order, p)
	}
	out, err = git(root, nil, "status", "--porcelain=v1", "-z", "--untracked-files=all")
	if err != nil {
		return nil, err
	}
	recs := strings.Split(string(out), "\x00")
	for i := 0; i < len(recs); i++ {
		rec := recs[i]
		if len(rec) < 4 {
			continue
		}
		x, y, p := rec[0], rec[1], rec[3:]
		if x == 'R' || x == 'C' {
			i++ // the original path follows; the new one is p
		}
		if !under(p) {
			continue
		}
		full := filepath.Join(root, p)
		fi, err := os.Lstat(full)
		if y == 'D' || err != nil || !fi.Mode().IsRegular() || fi.Size() > maxFileBytes {
			delete(view, p)
			continue
		}
		// The blob ID comes from the file's own bytes: `git hash-object
		// --stdin-paths` strips a trailing CR and unquotes a leading '"'
		// in a name, binding an ID to another file (security review,
		// finding 3).
		data, err := os.ReadFile(full)
		if err != nil {
			delete(view, p)
			continue
		}
		state := "modified"
		if x == '?' {
			state = "untracked"
		}
		if _, known := view[p]; !known {
			order = append(order, p)
		}
		view[p] = entry{path: p, oid: blobID(data), state: state}
	}
	entries := make([]entry, 0, len(view))
	for _, p := range order {
		if harness.SecretShaped(path.Base(p)) {
			continue // never indexed, even untracked (security review, finding 6)
		}
		if e, ok := view[p]; ok {
			entries = append(entries, e)
			delete(view, p) // an unmerged path appears once
		}
	}
	return entries, nil
}

// blobID is the git blob ID of data under SHA-1 ("blob <len>\x00" + data).
// It only has to name content consistently inside this index; a checkout in
// SHA-256 object format gets SHA-1 IDs for its modified files, which never
// collide with its index's IDs.
func blobID(data []byte) string {
	h := sha1.New()
	fmt.Fprintf(h, "blob %d\x00", len(data))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// configView lists the global index's own sources in ~/.config/sh-agent/,
// keyed by their SHA-256 since git does not hold them.
func configView(dir string) ([]entry, error) {
	names := []string{"feedback.jsonl", "retros.jsonl"}
	props, _ := filepath.Glob(filepath.Join(dir, "proposals", "*.md"))
	sort.Strings(props)
	for _, p := range props {
		names = append(names, filepath.ToSlash(strings.TrimPrefix(p, dir+string(filepath.Separator))))
	}
	var out []entry
	for _, n := range names {
		fi, err := os.Lstat(filepath.Join(dir, n))
		if err != nil || !fi.Mode().IsRegular() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, n))
		if err != nil {
			continue
		}
		sum := sha256.Sum256(data)
		out = append(out, entry{path: n, oid: "sha256:" + hex.EncodeToString(sum[:]), state: "tracked"})
	}
	return out, nil
}

// readIndexable returns a file's text, or ok=false when it is too large,
// binary, or gone.
func readIndexable(full string) (string, bool) {
	fi, err := os.Lstat(full)
	if err != nil || !fi.Mode().IsRegular() || fi.Size() > maxFileBytes {
		return "", false
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return "", false
	}
	sniff := data
	if len(sniff) > binarySniff {
		sniff = sniff[:binarySniff]
	}
	if bytes.IndexByte(sniff, 0) >= 0 {
		return "", false
	}
	return string(data), true
}

// source is what apply needs to know about one project.
type source struct {
	project  string
	dir      string // where the view's paths are read from
	docsRoot string // "" for ~config
	view     []entry
}

// apply brings one project's rows up to its current view and returns how
// many paths changed. When nothing changed and last_seen is fresh it takes
// no write lock at all.
func (s *store) apply(src source) (int, error) {
	now := nowFunc().Unix()
	changed, err := s.diffCount(src)
	if err != nil {
		return 0, err
	}
	if changed == 0 && !s.touchDue(src.project, now) {
		return 0, nil
	}
	if _, err := s.db.Exec(`BEGIN IMMEDIATE`); err != nil {
		if isBusy(err) {
			return 0, errBusy
		}
		return 0, err
	}
	n, err := s.applyLocked(src, now)
	if err != nil {
		s.db.Exec(`ROLLBACK`)
		if isBusy(err) {
			return 0, errBusy
		}
		return 0, err
	}
	if _, err := s.db.Exec(`COMMIT`); err != nil {
		s.db.Exec(`ROLLBACK`)
		if isBusy(err) {
			return 0, errBusy
		}
		return 0, err
	}
	if n > 0 {
		s.db.Exec(`PRAGMA wal_checkpoint(TRUNCATE)`)
	}
	return n, nil
}

// current reads the project's files rows and the blob table's skip marks.
func (s *store) current(project string) (map[string]string, map[string]bool, error) {
	files := map[string]string{}
	rows, err := s.db.Query(`SELECT path, oid FROM files WHERE project = ?`, project)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var p, o string
		if err := rows.Scan(&p, &o); err != nil {
			rows.Close()
			return nil, nil, err
		}
		files[p] = o
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	blobs := map[string]bool{} // oid -> skipped
	rows, err = s.db.Query(`SELECT oid, skipped FROM blobs`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var o string
		var sk int
		if err := rows.Scan(&o, &sk); err != nil {
			return nil, nil, err
		}
		blobs[o] = sk != 0
	}
	return files, blobs, rows.Err()
}

// wanted is the view minus paths whose blob is known to be skipped.
func wanted(view []entry, blobs map[string]bool) map[string]entry {
	w := map[string]entry{}
	for _, e := range view {
		if skipped, known := blobs[e.oid]; known && skipped {
			continue
		}
		w[e.path] = e
	}
	return w
}

func (s *store) diffCount(src source) (int, error) {
	files, blobs, err := s.current(src.project)
	if err != nil {
		return 0, err
	}
	w := wanted(src.view, blobs)
	n := 0
	for p, e := range w {
		if files[p] != e.oid {
			n++
		}
	}
	for p := range files {
		if _, ok := w[p]; !ok {
			n++
		}
	}
	return n, nil
}

func (s *store) touchDue(project string, now int64) bool {
	var v int64
	err := s.db.QueryRow(`SELECT CAST(value AS INTEGER) FROM meta WHERE key = ?`, "touched:"+project).Scan(&v)
	return err != nil || now-v >= touchEvery
}

func (s *store) applyLocked(src source, now int64) (int, error) {
	files, blobs, err := s.current(src.project)
	if err != nil {
		return 0, err
	}
	// Index every blob the database has not seen, parsed by the first path
	// it appears at in the view (D14).
	var in *inserter
	for _, e := range src.view {
		if _, known := blobs[e.oid]; known {
			continue
		}
		text, ok := readIndexable(filepath.Join(src.dir, filepath.FromSlash(e.path)))
		if _, err := s.db.Exec(`INSERT INTO blobs(oid, last_seen, skipped) VALUES (?, ?, ?)`, e.oid, now, boolInt(!ok)); err != nil {
			return 0, err
		}
		blobs[e.oid] = !ok
		if !ok {
			continue
		}
		if in == nil {
			if in, err = s.newInserter(); err != nil {
				return 0, err
			}
			defer in.close()
		}
		if err := in.content(e, text, src.docsRoot); err != nil {
			return 0, err
		}
	}
	w := wanted(src.view, blobs)
	n := 0
	for p, oid := range files {
		if e, ok := w[p]; !ok || e.oid != oid {
			if _, err := s.db.Exec(`DELETE FROM files WHERE project = ? AND path = ?`, src.project, p); err != nil {
				return 0, err
			}
			n++
		}
	}
	for p, e := range w {
		if old, ok := files[p]; ok && old == e.oid {
			if _, err := s.db.Exec(`UPDATE files SET state = ? WHERE project = ? AND path = ? AND state != ?`, e.state, src.project, p, e.state); err != nil {
				return 0, err
			}
			continue
		}
		if _, ok := files[p]; !ok {
			n++
		}
		if _, err := s.db.Exec(`INSERT INTO files(project, path, oid, state) VALUES (?, ?, ?, ?)`, src.project, p, e.oid, e.state); err != nil {
			return 0, err
		}
	}
	if _, err := s.db.Exec(`UPDATE blobs SET last_seen = ? WHERE oid IN (SELECT oid FROM files WHERE project = ?)`, now, src.project); err != nil {
		return 0, err
	}
	if _, err := s.db.Exec(`INSERT INTO meta(key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, "touched:"+src.project, now); err != nil {
		return 0, err
	}
	return n, s.gc(now)
}

// inserter holds the three insert statements, prepared once per reconcile:
// a full build runs some fifty thousand inserts.
type inserter struct{ chunk, record, ref *sql.Stmt }

func (s *store) newInserter() (*inserter, error) {
	var in inserter
	var err error
	if in.chunk, err = s.db.Prepare(`INSERT INTO chunks(oid, line, heading, body) VALUES (?, ?, ?, ?)`); err != nil {
		return nil, err
	}
	if in.record, err = s.db.Prepare(`INSERT INTO records(oid, id, line, title, status, date) VALUES (?, ?, ?, ?, ?, ?)`); err != nil {
		in.close()
		return nil, err
	}
	if in.ref, err = s.db.Prepare(`INSERT INTO refs(oid, line, target) VALUES (?, ?, ?)`); err != nil {
		in.close()
		return nil, err
	}
	return &in, nil
}

func (in *inserter) close() {
	for _, st := range []*sql.Stmt{in.chunk, in.record, in.ref} {
		if st != nil {
			st.Close()
		}
	}
}

func (in *inserter) content(e entry, text, docsRoot string) error {
	for _, c := range chunkText(e.path, text) {
		if _, err := in.chunk.Exec(e.oid, c.line, c.heading, c.body); err != nil {
			return err
		}
	}
	rel, _ := docsRelative(e.path, docsRoot)
	for _, r := range parseRecords(rel, text) {
		if _, err := in.record.Exec(e.oid, r.id, r.line, r.title, r.status, r.date); err != nil {
			return err
		}
	}
	for _, r := range parseRefs(text) {
		if _, err := in.ref.Exec(e.oid, r.line, r.target); err != nil {
			return err
		}
	}
	return nil
}

// gc drops blobs no project's view references that were last seen more
// than gcAge ago, with everything derived from them; the chunks trigger
// keeps chunks_fts in step.
func (s *store) gc(now int64) error {
	cut := now - gcAge
	old := `SELECT oid FROM blobs WHERE last_seen < ? AND oid NOT IN (SELECT oid FROM files)`
	for _, t := range []string{"chunks", "records", "refs"} {
		if _, err := s.db.Exec(`DELETE FROM `+t+` WHERE oid IN (`+old+`)`, cut); err != nil {
			return err
		}
	}
	_, err := s.db.Exec(`DELETE FROM blobs WHERE oid IN (`+old+`)`, cut)
	return err
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
