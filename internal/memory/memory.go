// Package memory is a project's memory store, .agents/memory/ at the
// repository root: one file per memory plus the MEMORY.md index. Both CLIs
// write and read it through `sh-agent memory`, so a memory is taken the same
// way on Claude Code and on Antigravity (R10, plan
// 2026-09-30-parity-memory-dream).
package memory

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

// The index caps match what Claude Code loads of its own MEMORY.md.
const (
	maxLines = 200
	maxBytes = 25 * 1024
	lockAge  = 30 * time.Second
)

// lockWait is how long Add retries a live lock before refusing.
var lockWait = 5 * time.Second

var (
	slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	types  = map[string]bool{"user": true, "feedback": true, "project": true, "reference": true}
)

// Entry is one memory's frontmatter as the caller gives it.
type Entry struct {
	Type        string // user, feedback, project or reference
	Name        string // the slug: file name and index title
	Description string // one line
}

// Refusal is a store the command will not touch, or input it will not take;
// the caller exits 2 with its message.
type Refusal struct{ Msg string }

func (r *Refusal) Error() string { return r.Msg }

func refuse(format string, a ...any) error { return &Refusal{Msg: fmt.Sprintf(format, a...)} }

// IsRefusal reports whether err is a Refusal.
func IsRefusal(err error) bool {
	var r *Refusal
	return errors.As(err, &r)
}

// Root is the main worktree's top level for the checkout that holds dir, so
// every worktree of a repository shares one store (D26).
func Root(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if err != nil {
		return "", refuse("not a git checkout: %s", dir)
	}
	common := strings.TrimSpace(string(out))
	if common == "" {
		return "", refuse("not a git checkout: %s", dir)
	}
	return filepath.Dir(common), nil
}

// trusted returns the store path of root after the checks every command
// makes: the store is the project's own directory, never a link a commit
// could carry and never content git tracks (memory-loop spec, ADR-043).
func trusted(root string) (string, error) {
	agents := filepath.Join(root, ".agents")
	store := filepath.Join(agents, "memory")
	for _, p := range []string{agents, store} {
		if info, err := os.Lstat(p); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", refuse("%s is a symlink: the memory store must be the project's own directory", p)
		}
	}
	out, err := exec.Command("git", "-C", root, "ls-files", "--", ".agents/memory").Output()
	if err != nil {
		return "", refuse("cannot ask git what it tracks under .agents/memory: %v", err)
	}
	if strings.TrimSpace(string(out)) != "" {
		return "", refuse("git tracks files under .agents/memory: the store is the operator's own, never the repository's")
	}
	return store, nil
}

// Add writes or replaces the memory e.Name with body, and its index line.
func Add(dir string, e Entry, body io.Reader, now time.Time) (string, error) {
	switch {
	case !types[e.Type]:
		return "", refuse("--type must be user, feedback, project or reference, not %q", e.Type)
	case !slugRe.MatchString(e.Name):
		return "", refuse("--name must match %s, not %q", slugRe, e.Name)
	case strings.TrimSpace(e.Description) == "":
		return "", refuse("--description is required")
	case strings.ContainsAny(e.Description, "\r\n"):
		return "", refuse("--description must be one line")
	}
	raw, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return "", refuse("the memory's body comes on stdin and is empty")
	}
	root, err := Root(dir)
	if err != nil {
		return "", err
	}
	store, err := trusted(root)
	if err != nil {
		return "", err
	}
	rel := ".agents/memory/" + e.Name + ".md"
	if exec.Command("git", "-C", root, "check-ignore", "-q", rel).Run() != nil {
		return "", refuse("%s is not ignored by git: run /init-project first", rel)
	}

	// The index is read under the lock, so two writers never lose a line.
	if err := os.MkdirAll(store, 0o700); err != nil {
		return "", err
	}
	unlock, err := lock(filepath.Join(store, ".lock"))
	if err != nil {
		return "", err
	}
	defer unlock()

	index := filepath.Join(store, "MEMORY.md")
	lines, err := readLines(index)
	if err != nil {
		return "", err
	}
	target := "(" + e.Name + ".md)"
	line := fmt.Sprintf("- [%s](%s.md) — %s", e.Name, e.Name, e.Description)
	replaced := false
	for i, l := range lines {
		if strings.Contains(l, target) {
			lines[i], replaced = line, true
		}
	}
	if !replaced {
		if len(lines)+1 > maxLines {
			return "", refuse("MEMORY.md already holds %d lines, the most a session loads: run /dream to consolidate", len(lines))
		}
		lines = append(lines, line)
	}

	desc, _ := json.Marshal(e.Description)
	file := fmt.Sprintf("---\nname: %s\ndescription: %s\nmetadata:\n  node_type: memory\n  type: %s\n  modified: %s\n---\n\n%s\n",
		e.Name, desc, e.Type, now.UTC().Format(time.RFC3339), text)
	path := filepath.Join(store, e.Name+".md")
	if err := harness.WriteAtomic(path, []byte(file), 0o600); err != nil {
		return "", err
	}
	if err := harness.WriteAtomic(index, []byte(strings.Join(lines, "\n")+"\n"), 0o600); err != nil {
		return "", err
	}
	return path, nil
}

// Show returns the index, capped at what a session loads, plus the rollout
// hint when Claude Code's native memory may still be on here (D20, D27):
// claudeInstalled is true unless the install manifest names Antigravity alone.
func Show(dir string, claudeInstalled bool) (string, error) {
	root, err := Root(dir)
	if err != nil {
		return "", err
	}
	store, err := trusted(root)
	if err != nil {
		return "", err
	}
	lines, err := readLines(filepath.Join(store, "MEMORY.md"))
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for i, l := range lines {
		if i == maxLines || b.Len()+len(l)+1 > maxBytes {
			fmt.Fprintf(&b, "(… %d more lines in .agents/memory/MEMORY.md)\n", len(lines)-i)
			break
		}
		b.WriteString(l + "\n")
	}
	if nativeOn(root, claudeInstalled) {
		b.WriteString("Claude Code's native memory is still on here: run /init-project to switch it off (R10)\n")
	}
	return b.String(), nil
}

// nativeOn reports whether this project's settings leave Claude Code's own
// auto memory on: a settings file without "autoMemoryEnabled": false, or no
// file at all on a machine where Claude Code is installed.
func nativeOn(root string, claudeInstalled bool) bool {
	data, err := os.ReadFile(filepath.Join(root, ".claude", "settings.local.json"))
	if err != nil {
		return claudeInstalled
	}
	var s map[string]any
	if json.Unmarshal(data, &s) != nil {
		return true
	}
	v, ok := s["autoMemoryEnabled"].(bool)
	return !ok || v
}

// readLines returns the non-empty lines of path; a missing file is empty.
func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		if l := strings.TrimRight(sc.Text(), "\r"); l != "" {
			lines = append(lines, l)
		}
	}
	return lines, sc.Err()
}

// lock takes the store's lock file, holding this process's pid. A lock whose
// pid is gone, or older than lockAge, is stale and replaced; a live one is
// retried for lockWait, then refused.
func lock(path string) (func(), error) {
	deadline := time.Now().Add(lockWait)
	for {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			fmt.Fprint(f, os.Getpid())
			f.Close()
			return func() { os.Remove(path) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if stale(path) {
			os.Remove(path)
			continue
		}
		if time.Now().After(deadline) {
			return nil, refuse("the memory store is busy: another sh-agent memory holds %s", path)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func stale(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if time.Since(info.ModTime()) > lockAge {
		return true
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	return err != nil || !harness.Alive(pid)
}
