// Package trail implements the Definition-of-Done trail contract: where a
// project's docs root lives, which changed files owe which DoD rows, and
// whether the artifact for a row already exists.
//
// Artifacts are matched by their frontmatter (owner + date), never by
// filename. The previous shell gates globbed on `<date>-<infix>*.md`, which
// silently failed whenever a dispatch prompt named an output path that did
// not happen to match the convention.
package trail

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Row is one Definition-of-Done row: an artifact class with an owning agent.
type Row string

// The documentation row is deliberately absent. "Behaviour or interfaces
// changed" is a judgement the dispatcher cannot compute from a diff, so that
// row stays model-triggered rather than being declared here and never fired.
const (
	RowChangelog Row = "changelog"
	RowReviewer  Row = "reviewer"
	RowSecurity  Row = "security"
	RowVerifier  Row = "verifier"
	RowExplorer  Row = "explorer"
)

// Class is the docs-root subdirectory a row's artifact lands in.
func (r Row) Class() string {
	switch r {
	case RowChangelog:
		return "changes"
	case RowReviewer, RowSecurity:
		return "reviews"
	case RowVerifier:
		return "testing"
	case RowExplorer:
		return "plans"
	default:
		return ""
	}
}

// Resumability reports whether this row exists so a zero-context agent can
// resume the work. Those rows are mandatory. The remaining rows exist for
// accountability — they are re-derivable from the diff and fire proportionally.
func (r Row) Resumability() bool { return r == RowChangelog }

// Project is a resolved harness project: a git work tree whose docs root has
// opted into the trail contract.
type Project struct {
	Root string // absolute path to the git work tree
	Docs string // absolute path to the docs root

	// SourceGlobs are project-declared paths that count as source even when the
	// default classification would not treat them as such — see IsSourcePath.
	SourceGlobs []string
}

// DocsRel is the docs root relative to the project root, for display and for
// path filtering.
func (p Project) DocsRel() string {
	rel, err := filepath.Rel(p.Root, p.Docs)
	if err != nil {
		return p.Docs
	}
	return rel
}

// Opted reports whether the project adopted the trail contract. A missing
// changes/ directory is the opt-in marker's absence: gates must stay silent
// in a foreign clone rather than nag about a contract it never accepted.
func (p Project) Opted() bool {
	fi, err := os.Stat(filepath.Join(p.Docs, "changes"))
	return err == nil && fi.IsDir()
}

// Resolve locates the project containing dir. It returns ok=false outside a
// git work tree.
func Resolve(dir string) (Project, bool) {
	root, err := gitToplevel(dir)
	if err != nil || root == "" {
		return Project{}, false
	}
	return Project{Root: root, Docs: docsRoot(root)}, true
}

// docsRoot prefers .agents/docs, falls back to a plain docs/ directory, and
// otherwise names .agents/docs as the path that would be created.
func docsRoot(root string) string {
	if fi, err := os.Stat(filepath.Join(root, ".agents", "docs")); err == nil && fi.IsDir() {
		return filepath.Join(root, ".agents", "docs")
	}
	if fi, err := os.Stat(filepath.Join(root, "docs")); err == nil && fi.IsDir() {
		return filepath.Join(root, "docs")
	}
	return filepath.Join(root, ".agents", "docs")
}

// Artifact is one on-disk DoD artifact, identified by its frontmatter.
type Artifact struct {
	Path  string
	Owner Row
	Date  string // YYYY-MM-DD
}

// HasArtifact reports whether an artifact owned by row and dated day exists
// under the project's docs root.
func (p Project) HasArtifact(row Row, day string) bool {
	for _, a := range p.Artifacts(row.Class()) {
		if a.Owner == row && a.Date == day {
			return true
		}
	}
	return false
}

// Artifacts lists every artifact in one docs-root class whose frontmatter
// declares an owner and a date. Files without that frontmatter are ignored:
// the trail records what claims ownership, not whatever happens to be filed.
func (p Project) Artifacts(class string) []Artifact {
	if class == "" {
		return nil
	}
	entries, err := os.ReadDir(filepath.Join(p.Docs, class))
	if err != nil {
		return nil
	}
	var out []Artifact
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(p.Docs, class, e.Name())
		owner, date, ok := readFrontmatter(path)
		if !ok {
			continue
		}
		out = append(out, Artifact{Path: path, Owner: Row(owner), Date: date})
	}
	return out
}

// readFrontmatter extracts `owner:` and `date:` from a leading YAML
// frontmatter block. It reads only the frontmatter, never the whole file.
func readFrontmatter(path string) (owner, date string, ok bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", false
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	if !sc.Scan() || strings.TrimSpace(sc.Text()) != "---" {
		return "", "", false
	}
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "---" {
			break
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		switch strings.TrimSpace(key) {
		case "owner":
			owner = value
		case "date":
			date = value
		}
	}
	if owner == "" || date == "" {
		return "", "", false
	}
	return owner, date, true
}

// Today is the current date in the trail's YYYY-MM-DD form. Every date in the
// contract comes from the system clock, never from a model.
func Today() string { return time.Now().Format("2006-01-02") }
