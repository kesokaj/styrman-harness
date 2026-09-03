package trail

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
)

// scopePrefix is the form this binary produces. Anything else is from a
// superseded mechanism and must not satisfy a row.
const scopePrefix = "changeset:"

// Scope identifies *what* a change set is, so an artifact can record what it
// covered rather than merely when it was written.
//
// Dates were the previous answer and they were too coarse by a wide margin: any
// artifact dated today satisfied its row, so a second substantive commit on the
// same day owed nothing. That gap hid owed rows five separate times on the day
// it was found, including for changes to the gate logic itself.
//
// The scope is a digest over the changed source files — their paths and their
// contents. That choice is what makes it survive the commit that lands the work:
// committing moves files from untracked to tracked without altering a byte on
// disk, so the digest is identical either side of it. An earlier version keyed
// the two states differently (`worktree:` before, `commit:` after) and re-owed
// every row the moment the work was committed — a gate rejecting work it had
// just accepted, which is precisely the false-positive cost R4 exists to avoid.
func (p Project) Scope() string {
	files := p.DirtySource()
	if len(files) == 0 {
		// A clean tree: the change set is whatever HEAD landed. The same files,
		// with the same contents, therefore the same digest.
		if !p.headIsToday() {
			return ""
		}
		files = p.HeadSourceToday()
	}
	if len(files) == 0 {
		return ""
	}
	return scopePrefix + p.digestOf(files)
}

// digestOf fingerprints a set of source files by path and content.
//
// Content, not diff. `git diff` would have been the obvious input and was the
// original one, but it cannot survive a commit — and it also drags in anything
// else the tree happens to be carrying, so an artifact written during a dispatch
// could move the scope out from under the artifacts written beside it.
func (p Project) digestOf(files []string) string {
	sorted := append([]string(nil), files...)
	sort.Strings(sorted)

	h := sha256.New()
	for _, rel := range sorted {
		h.Write([]byte(rel))
		h.Write([]byte{0})
		if data, err := os.ReadFile(filepath.Join(p.Root, rel)); err == nil {
			h.Write(data)
		}
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))[:12]
}

// Covers reports whether an artifact's recorded scope covers the given one.
//
// Exact match, and nothing else. An artifact with no scope covers nothing —
// artifacts written before scopes existed, or by an agent that ignored the
// instruction, must not silently satisfy a row they may have nothing to do with.
//
// There is deliberately no `<base>..<head>` range form. One was written, and it
// was unreachable: it could only match a scope carrying a `commit:` prefix, and
// no code path produces one. Rather than leave a branch that looks like a
// feature and never fires, range coverage is a judgement the release gate makes
// by reading frontmatter — not something this binary computes.
func (p Project) Covers(a Artifact, scope string) bool {
	return a.Scope != "" && scope != "" && a.Scope == scope
}
