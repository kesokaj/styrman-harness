package harness

import (
	"os"
	"path/filepath"
)

// PathShadows returns the executables named like binPath that a shell finds on
// PATH BEFORE binPath itself, in PATH order, without repeats.
//
// An empty result means an unqualified `sh-agent` runs the binary the install
// just wrote. A non-empty one means it does not — which is how a downloaded,
// signature-verified, checksum-verified release binary was installed to
// /usr/local/bin and then ignored for two more PATH entries, while every
// command reported success and `status` reported drift none, because drift is
// measured against whatever source the binary that actually ran resolved
// (BUG-078).
//
// resolve reports the real path of an executable file at a location, and false
// when there is none; it is injected so the tests need no real files.
func PathShadows(pathEnv, binPath string, resolve func(string) (string, bool)) []string {
	name := filepath.Base(binPath)
	binDir := filepath.Clean(filepath.Dir(binPath))
	target, haveTarget := resolve(binPath)

	var shadows []string
	seen := map[string]bool{}
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			dir = "." // POSIX: an empty PATH entry means the working directory.
		}
		if filepath.Clean(dir) == binDir {
			// PATH reaches the installed binary here; nothing later can win.
			break
		}
		cand := filepath.Join(dir, name)
		if seen[cand] {
			continue
		}
		got, ok := resolve(cand)
		if !ok {
			continue
		}
		// A symlink or bind that resolves to the installed binary runs the same
		// file, so it is not a shadow. This is how BUG-078 was settled on the
		// operator's box, and the warning must go quiet once it is.
		if haveTarget && got == target {
			continue
		}
		seen[cand] = true
		shadows = append(shadows, cand)
	}
	return shadows
}

// ResolveExecutable reports the real path of an executable file at p. A
// directory, a missing file, a dangling symlink, or a file with no execute bit
// is not one.
func ResolveExecutable(p string) (string, bool) {
	real, err := filepath.EvalSymlinks(p)
	if err != nil {
		return "", false
	}
	info, err := os.Stat(real)
	if err != nil || info.IsDir() || info.Mode().Perm()&0o111 == 0 {
		return "", false
	}
	return real, true
}
