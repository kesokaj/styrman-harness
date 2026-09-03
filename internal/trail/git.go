package trail

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// git runs a git command with repo-local config that could execute a program
// pinned off. core.fsmonitor runs an arbitrary binary on `git status`, and the
// hook payload chooses the directory we run in — so a gate must never let a
// surveyed repo's own config execute.
func git(dir string, args ...string) (string, error) {
	full := append([]string{"-c", "core.fsmonitor=false", "-c", "diff.external=", "-C", dir}, args...)
	out, err := exec.Command("git", full...).Output()
	return strings.TrimSpace(string(out)), err
}

func gitToplevel(dir string) (string, error) {
	return git(dir, "rev-parse", "--show-toplevel")
}

// DirtySource lists source files modified, staged, or untracked right now. A
// session that only creates new files owes the same rows, so untracked files
// count (-uall lists them individually inside untracked directories).
func (p Project) DirtySource() []string {
	out, err := git(p.Root, "status", "--porcelain", "-uall")
	if err != nil {
		return nil
	}
	var files []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		path := line[3:]
		// Renames arrive as "old -> new"; the new path is what changed.
		if _, after, found := strings.Cut(path, " -> "); found {
			path = after
		}
		path = strings.Trim(path, `"`)
		if p.IsSourcePath(path) {
			files = append(files, path)
		}
	}
	return files
}

// HeadSourceToday lists source files in HEAD when HEAD was committed today.
// Committed-but-unrecorded still owes the trail.
func (p Project) HeadSourceToday() []string {
	day, err := git(p.Root, "log", "-1", "--format=%cd", "--date=short")
	if err != nil || day != Today() {
		return nil
	}
	// diff-tree names files for merge commits too, which `git show` suppresses.
	out, err := git(p.Root, "diff-tree", "-r", "--no-commit-id", "--name-only", "HEAD")
	if err != nil {
		return nil
	}
	var files []string
	for _, path := range strings.Split(out, "\n") {
		if path != "" && p.IsSourcePath(path) {
			files = append(files, path)
		}
	}
	return files
}

// ChangedLines counts added+removed lines across the working tree and index,
// the size input to the proportionality rule.
//
// Untracked files are counted in full. git's numstat says nothing about a file
// it has never seen, so a change made entirely of new files would otherwise
// score zero and skip the reviewer row — which is exactly backwards, since all-
// new code is the case least likely to have been read by anyone.
func (p Project) ChangedLines() int {
	total := 0
	for _, args := range [][]string{
		{"diff", "--numstat"},
		{"diff", "--numstat", "--cached"},
	} {
		out, err := git(p.Root, args...)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(out, "\n") {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			total += atoiSafe(fields[0]) + atoiSafe(fields[1])
		}
	}
	return total + p.untrackedLines()
}

// untrackedLines counts the lines of untracked source files.
func (p Project) untrackedLines() int {
	out, err := git(p.Root, "status", "--porcelain", "-uall")
	if err != nil {
		return 0
	}
	total := 0
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "??") || len(line) < 4 {
			continue
		}
		rel := strings.Trim(line[3:], `"`)
		if !p.IsSourcePath(rel) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(p.Root, rel))
		if err != nil {
			continue
		}
		total += bytes.Count(data, []byte{'\n'})
	}
	return total
}

// atoiSafe treats git's "-" (binary file) and any malformed count as zero.
func atoiSafe(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}
