package harness

import (
	"os"
	"path/filepath"
)

// PathShadows walks PATH the way a shell does and reports two facts about the
// binary at binPath: the executables of the same name that PATH finds BEFORE
// it (in PATH order, without repeats), and whether its directory is on PATH at
// all. The callers print from these; they never re-derive "on PATH" by another
// rule, because two definitions of it disagreed on a `.local/bin//` entry and
// hid a real shadow (review of 2026-09-18).
//
// An empty shadows result with onPath true means an unqualified `sh-agent`
// runs the binary the install just wrote. Anything else means it does not —
// which is how a downloaded, signature-verified, checksum-verified release
// binary was installed to /usr/local/bin and then ignored for two more PATH
// entries, while every command reported success and `status` reported drift
// none, because drift is measured against whatever source the binary that
// actually ran resolved (BUG-078).
func PathShadows(pathEnv, binPath string) (shadows []string, onPath bool) {
	name := filepath.Base(binPath)
	binDir := filepath.Clean(filepath.Dir(binPath))
	target, haveTarget := executableAt(binPath)

	seen := map[string]bool{}
	for _, dir := range filepath.SplitList(pathEnv) {
		if filepath.Clean(dir) == binDir {
			// PATH reaches the installed binary here; nothing later can win.
			return shadows, true
		}
		cand := filepath.Join(dir, name)
		if dir == "" {
			// POSIX: an empty PATH entry is the working directory. Join would
			// print it as a bare name that reads like a command, not a file.
			cand = "." + string(filepath.Separator) + name
		}
		if seen[cand] {
			continue
		}
		info, ok := executableAt(cand)
		if !ok {
			continue
		}
		// A symlink or hard link to the installed binary runs the same file,
		// so it is not a shadow. This is how BUG-078 was settled on the
		// operator's box, and the warning must go quiet once it is.
		if haveTarget && os.SameFile(info, target) {
			continue
		}
		seen[cand] = true
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
