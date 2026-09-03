package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/trail"
)

// ResumeContext builds the Resume Protocol orientation block for a directory:
// at most a handful of short lines telling an agent where the work left off.
// It returns "" outside a git repo or in a project that never adopted the
// trail contract — a gate has nothing to say about a foreign clone.
func ResumeContext(dir string) string {
	proj, ok := trail.Resolve(dir)
	if !ok || !proj.Opted() {
		return ""
	}

	var lines []string
	lines = append(lines, "Resume Protocol: AGENTS.md -> REQUIREMENTS.md -> ROADMAP.md -> TODO.md (open items) -> newest handoff -> DECISIONS.md (index below) -> 3 newest changes/")

	if h := newestDated(filepath.Join(proj.Docs, "handoffs")); h != "" {
		rel, _ := filepath.Rel(proj.Root, h)
		if age, ok := ageDays(filepath.Base(h)); ok {
			lines = append(lines, fmt.Sprintf("Latest handoff: %s (%dd old)", safeText(rel), age))
		} else {
			lines = append(lines, "Latest handoff: "+safeText(rel))
		}
	} else {
		lines = append(lines, "no handoffs recorded")
	}

	if line := adrIndex(proj); line != "" {
		lines = append(lines, line)
	}
	if line := changesIndex(proj); line != "" {
		lines = append(lines, line)
	}
	if line := missingRootFiles(proj); line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// adrIndex summarises DECISIONS.md rather than quoting it. The file is
// append-only and grows without bound, so the Resume Protocol reads this line
// plus the newest few ADRs, never the whole thing.
func adrIndex(proj trail.Project) string {
	path := filepath.Join(proj.Docs, "architecture", "DECISIONS.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var titles []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "## ADR") {
			titles = append(titles, strings.TrimPrefix(line, "## "))
		}
	}
	if len(titles) == 0 {
		return ""
	}
	newest := titles
	if len(newest) > 3 {
		newest = newest[len(newest)-3:]
	}
	for i, t := range newest {
		newest[i] = safeText(truncate(t, 80))
	}
	rel, _ := filepath.Rel(proj.Root, path)
	return fmt.Sprintf("ADRs: %d in %s — newest: %s", len(titles), rel, strings.Join(newest, " | "))
}

// changesIndex names the three newest changelog entries. Most sessions never
// need to open one — the index line answers the question by itself.
func changesIndex(proj trail.Project) string {
	arts := proj.Artifacts("changes")
	if len(arts) == 0 {
		return ""
	}
	sort.Slice(arts, func(i, j int) bool { return arts[i].Date > arts[j].Date })
	if len(arts) > 3 {
		arts = arts[:3]
	}
	var names []string
	for _, a := range arts {
		names = append(names, safeText(strings.TrimSuffix(filepath.Base(a.Path), ".md")))
	}
	return "Changes (newest 3): " + strings.Join(names, " | ")
}

// missingRootFiles names the Resume Protocol's expected files that are absent
// or empty, so the agent reports the gap instead of guessing intent.
func missingRootFiles(proj trail.Project) string {
	var missing []string
	for _, name := range []string{"AGENTS.md", "REQUIREMENTS.md", "ROADMAP.md", "TODO.md"} {
		fi, err := os.Stat(filepath.Join(proj.Root, name))
		if err != nil || fi.Size() == 0 {
			missing = append(missing, name)
		}
	}
	if len(missing) == 0 {
		return ""
	}
	return "Missing or empty: " + strings.Join(missing, " ")
}

// newestDated returns the newest date-prefixed file in dir, chosen by filename
// rather than mtime — mtime lies after a clone, checkout, or rsync, and the
// trail's filenames are date-prefixed by contract.
func newestDated(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && datePrefix(e.Name()) != "" {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	return filepath.Join(dir, names[len(names)-1])
}

// datePrefix returns the leading YYYY-MM-DD of a filename, or "".
func datePrefix(name string) string {
	if len(name) < 10 {
		return ""
	}
	head := name[:10]
	for i, r := range head {
		if i == 4 || i == 7 {
			if r != '-' {
				return ""
			}
			continue
		}
		if r < '0' || r > '9' {
			return ""
		}
	}
	return head
}

func ageDays(name string) (int, bool) {
	day := datePrefix(name)
	if day == "" {
		return 0, false
	}
	t, err := time.Parse("2006-01-02", day)
	if err != nil {
		return 0, false
	}
	return int(time.Since(t).Hours() / 24), true
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// safeText reduces repo-sourced text to a bounded, inert character set before
// it is embedded in a line of model-facing context.
//
// Filenames and ADR titles are repo content — data, not instructions. Character
// filtering alone is not a data/instruction boundary and this does not pretend
// to be one; what it does is remove the punctuation that markup and fenced
// blocks are built from, and cap the length, so a crafted filename cannot
// restructure the orientation block or crowd out the rest of it. The real
// boundary is that this block only ever states facts the agent can re-derive.
func safeText(s string) string {
	const maxLen = 120
	var b strings.Builder
	for _, r := range s {
		if b.Len() >= maxLen {
			b.WriteString("…")
			break
		}
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case strings.ContainsRune(" ._/:()-", r):
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}
