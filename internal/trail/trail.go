// Package trail answers when a project's trail records changed, from git:
// the nightly dreamer's watermark and the metrics store's events both need
// it, and neither imports the other (plan 2026-10-08-metrics-store D31). It
// imports no SQLite driver.
package trail

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Dates is, for every tracked file under docs, the time of its newest
// commit (paths relative to root, slash-separated), and the set of paths git
// reports as modified or untracked there. Both empty when git is absent or
// root is no checkout: every file then keeps its mtime.
func Dates(root, docs string) (dates map[string]time.Time, modified map[string]bool) {
	dates, modified = map[string]time.Time{}, map[string]bool{}
	out, err := exec.Command("git", "-C", root, "log", "--format=%ct", "--name-only", "--no-renames", "--", docs).Output()
	if err != nil {
		return dates, modified
	}
	var t time.Time
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		if secs, err := strconv.ParseInt(line, 10, 64); err == nil {
			t = time.Unix(secs, 0)
			continue
		}
		if _, seen := dates[line]; !seen && !t.IsZero() {
			dates[line] = t
		}
	}
	out, err = exec.Command("git", StatusArgs(root, docs)...).Output()
	if err != nil {
		return dates, modified
	}
	for _, line := range strings.Split(string(out), "\n") {
		if len(line) > 3 {
			modified[strings.TrimSpace(line[3:])] = true
		}
	}
	return dates, modified
}

// StatusArgs is the status command run in a project: the global options
// come first (git rejects them after the verb), it holds no index lock, and
// it runs no fsmonitor hook the project may have configured (security review
// rounds 2 and 3, 2026-10-01).
func StatusArgs(root, docs string) []string {
	return []string{"--no-optional-locks", "-c", "core.fsmonitor=false", "-C", root, "status", "--porcelain", "--untracked-files=all", "--", docs}
}
