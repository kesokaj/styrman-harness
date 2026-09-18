package harness

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PathShadows walks PATH the way a shell does and reports two facts about the
// binary at binPath: the executables of the same name that PATH finds BEFORE
// it (in PATH order, without repeats), and whether PATH reaches it at all.
//
// "Reaches it" means one thing only: a PATH entry holds an executable that IS
// the installed file, by os.SameFile — through its own directory, a symlinked
// directory, a symlink or a hard link to it. There is no lexical comparison
// of directory names anywhere: two rounds of review found a directory form
// (`.local/bin//`, then a symlinked directory) on which a lexical "on PATH"
// and the file identity disagreed, and each time the disagreement hid a real
// shadow. One identity, one walk.
//
// An empty shadows result with onPath true means an unqualified `sh-agent`
// runs the binary the install just wrote. Anything else means it does not —
// which is how a downloaded, signature-verified, checksum-verified release
// binary was installed to /usr/local/bin and then ignored for two more PATH
// entries, while every command reported success and `status` reported drift
// none, because drift is measured against whatever source the binary that
// actually ran resolved (BUG-078). When binPath does not exist, nothing can
// reach it: onPath is false and every executable of that name is a shadow.
func PathShadows(pathEnv, binPath string) (shadows []string, onPath bool) {
	name := filepath.Base(binPath)
	target, haveTarget := executableAt(binPath)

	seen := map[string]bool{}
	for _, dir := range filepath.SplitList(pathEnv) {
		cand := filepath.Join(dir, name)
		if dir == "" {
			// POSIX: an empty PATH entry is the working directory. Join would
			// print it as a bare name that reads like a command, not a file.
			cand = "." + string(filepath.Separator) + name
		}
		if seen[cand] {
			continue
		}
		seen[cand] = true
		info, ok := executableAt(cand)
		if !ok {
			continue
		}
		if haveTarget && os.SameFile(info, target) {
			// The shell would run the installed file here; nothing later can win.
			return shadows, true
		}
		shadows = append(shadows, cand)
	}
	return shadows, false
}

// executableAt reports the file a shell would execute at p, following
// symlinks. A directory, a missing file, a dangling symlink, or a file with no
// execute bit is not one.
func executableAt(p string) (os.FileInfo, bool) {
	info, err := os.Stat(p)
	if err != nil || info.IsDir() || info.Mode().Perm()&0o111 == 0 {
		return nil, false
	}
	return info, true
}

// RunningElsewhere reports where a shell finds the binary at exe by the name
// of binPath, when exe is NOT the file at binPath. `update` replaces the file
// at binPath and nothing else, so a process started from any other path has
// just updated a file it is not running — and a shell that resolved the name
// to that other path keeps running it whatever PATH says, because the shell
// caches the lookup (BUG-079: /usr/local/bin/sh-agent ran `update`,
// ~/.local/bin got v1.2.2, bash kept running v1.2.1 with the right PATH
// order). A process no PATH entry reaches by name — a test binary, a `go run`
// build — is nobody's `sh-agent` and is not reported.
func RunningElsewhere(pathEnv, binPath, exe string) (string, bool) {
	self, err := os.Stat(exe)
	if err != nil {
		return "", false
	}
	// No installed file, nothing to be behind: an `update` that installed
	// nothing into a home with no binary must not name a file that does not
	// exist and call the running one stale (security review 2026-09-18).
	inst, err := os.Stat(binPath)
	if err != nil || os.SameFile(self, inst) {
		return "", false
	}
	name := filepath.Base(binPath)
	for _, dir := range filepath.SplitList(pathEnv) {
		cand := filepath.Join(dir, name)
		if dir == "" {
			cand = "." + string(filepath.Separator) + name // as PathShadows prints it
		}
		if info, err := os.Stat(cand); err == nil && os.SameFile(info, self) {
			return cand, true
		}
	}
	return "", false
}

// NoteRunningElsewhere says, after an update wrote the installed binary, when
// the process that ran the update is not that binary and a shell can reach the
// process's file by name: that file stays as it is, and a shell that already
// resolved `sh-agent` to it keeps running it whatever PATH says (BUG-079).
// Called only after a binary was actually written — an up-to-date or
// older-release run installs nothing, so there is nothing to be behind, and on
// the older-release branch the running file is the NEWER one (security review
// 2026-09-18, round 2).
func NoteRunningElsewhere(home string, out io.Writer) {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	installed := filepath.Join(home, BinPath)
	if from, ok := RunningElsewhere(os.Getenv("PATH"), installed, exe); ok {
		fmt.Fprintf(out, "note: this update ran from %s, not from the installed binary %s; update replaces only the installed one, so the file it ran from stays as it is. A shell that has already resolved `sh-agent` to that path keeps running it until `hash -r` (bash) or `rehash` (zsh); point it at the installed binary (`ln -sfn`) or remove it — when that path is root-owned, prefer removing it\n", from, installed)
	}
}
