package search

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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

// noIndexMarker names the file that takes its directory out of the index.
const noIndexMarker = ".sh-agent-noindex"

// gitView lists the checkout's current files: the index entries of
// `git ls-files -s`, overridden by what `git status` reports modified or
// untracked (their blob IDs computed from their bytes), minus deleted paths,
// symlinks, submodules, secret-shaped names and every path under a
// directory that holds a noIndexMarker file. A file's content is read
// again when its blob is indexed: a file edited and restored between the two
// reads is indexed with the edit until `index --rebuild` (a known limit; the
// operator's choice at the security review's cap, plan D35).
// prefix limits it to one directory.
func gitView(root, prefix string) ([]entry, error) {
	under := func(p string) bool { return prefix == "" || strings.HasPrefix(p, prefix) }
	// Every read goes through an os.Root: a path that leaves the checkout
	// through `..` or a symlinked directory is refused by the standard
	// library, whatever git lists (security check of the v1.13.0 gate fixes).
	fsRoot, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	defer fsRoot.Close()
	out, err := git(root, nil, "ls-files", "-s", "-z")
	if err != nil {
		return nil, err
	}
	var order []string
	view := map[string]entry{}
	// Markers are collected from the whole checkout, not only under prefix:
	// a marker above the prefix covers it (a checkout marked at its root has
	// no docs trail for the global index either).
	markers := map[string]bool{}
	for _, rec := range strings.Split(string(out), "\x00") {
		meta, p, ok := strings.Cut(rec, "\t")
		if ok && path.Base(p) == noIndexMarker {
			markers[p] = true
		}
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
	// --no-renames: a rename is two records (D and A), never a pair the
	// parser has to keep in step — a work-tree rename was once read as an
	// entry with three characters cut off, which pointed outside the
	// checkout (release gate security review, v1.13.0).
	out, err = git(root, nil, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--no-renames")
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
		if path.Base(p) == noIndexMarker {
			if _, err := fsRoot.Lstat(filepath.FromSlash(p)); y == 'D' || err != nil {
				delete(markers, p) // a marker deleted from the work tree no longer counts
			} else {
				markers[p] = true
			}
		}
		if !under(p) {
			continue
		}
		fi, err := fsRoot.Lstat(filepath.FromSlash(p))
		if y == 'D' || err != nil || !fi.Mode().IsRegular() || fi.Size() > maxFileBytes {
			delete(view, p)
			continue
		}
		// The blob ID comes from the file's own bytes: `git hash-object
		// --stdin-paths` strips a trailing CR and unquotes a leading '"'
		// in a name, binding an ID to another file (security review,
		// finding 3).
		data, err := readAll(fsRoot, filepath.FromSlash(p))
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
	// A directory holding a noIndexMarker is data the project carries, not its
	// own text: nothing below it is indexed. The marker counts in every state
	// git reports it in (committed, staged, not yet added): the first search
	// after the marker is written must already leave the directory out.
	var marked []string
	for p := range markers {
		dir := path.Dir(p) + "/"
		if dir == "./" {
			dir = "" // a marker at the root: nothing of the checkout is indexed
		}
		marked = append(marked, dir)
	}
	noIndex := func(p string) bool {
		for _, dir := range marked {
			if strings.HasPrefix(p, dir) {
				return true
			}
		}
		return false
	}
	entries := make([]entry, 0, len(view))
	for _, p := range order {
		if harness.SecretShaped(path.Base(p)) {
			continue // never indexed, even untracked (security review, finding 6)
		}
		if noIndex(p) {
			continue
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
	// Listed by name, never through a pattern built from the config path: a
	// `[` in it would be read as a character class (the golang rule).
	if entries, err := os.ReadDir(filepath.Join(dir, "proposals")); err == nil {
		var props []string
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				props = append(props, "proposals/"+e.Name())
			}
		}
		sort.Strings(props)
		names = append(names, props...)
	}
	// Plain reads: ~/.config/sh-agent is the operator's own state, and a
	// dotfiles arrangement may symlink its proposals directory.
	var out []entry
	for _, n := range names {
		fi, err := os.Stat(filepath.Join(dir, n))
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

// readAll reads one file through the root, so the read cannot leave it. On
// a Go before 1.25.12, os.Root.Open carries advisory GO-2026-4970 (an escape
// through a path ending in "/"); every call here follows an Lstat of the
// same path that reported a regular file, which such a path never is.
func readAll(fsRoot *os.Root, rel string) ([]byte, error) {
	f, err := fsRoot.Open(rel)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// readIndexable returns a file's text, or ok=false when it is too large,
// binary, gone, or outside the root.
func readIndexable(fsRoot *os.Root, rel string) (string, bool) {
	fi, err := fsRoot.Lstat(rel)
	if err != nil || !fi.Mode().IsRegular() || fi.Size() > maxFileBytes {
		return "", false
	}
	data, err := readAll(fsRoot, rel)
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
	// it appears at in the view (D14). Reads go through an os.Root.
	fsRoot, err := os.OpenRoot(src.dir)
	if err != nil {
		return 0, err
	}
	defer fsRoot.Close()
	var in *inserter
	for _, e := range src.view {
		if _, known := blobs[e.oid]; known {
			continue
		}
		text, ok := readIndexable(fsRoot, filepath.FromSlash(e.path))
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
