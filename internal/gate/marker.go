package gate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// markerDir returns a private per-user scratch directory for gate state,
// creating it if needed. It returns ok=false rather than falling back to a
// shared path: a marker another local user can forge or drop would let them
// steer a gate's decision.
func markerDir(name string) (string, bool) {
	dir := filepath.Join(os.TempDir(), name+"-"+strconv.Itoa(os.Getuid()))

	fi, err := os.Lstat(dir)
	switch {
	case err == nil && fi.Mode()&os.ModeSymlink != 0:
		return "", false // a symlink could redirect writes anywhere
	case err == nil && !fi.IsDir():
		return "", false
	case os.IsNotExist(err):
		if err := os.Mkdir(dir, 0o700); err != nil {
			return "", false
		}
	case err != nil:
		return "", false
	}
	if !ownedByUs(dir) {
		return "", false
	}
	return dir, true
}

// safeName reduces untrusted payload text (conversation and agent ids) to a
// character set that cannot escape the marker directory.
func safeName(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9',
			r == '.', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
		if b.Len() >= 128 {
			break
		}
	}
	return b.String()
}

// touch creates or truncates a marker file, reporting whether it worked.
func touch(path string) bool {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return false
	}
	return f.Close() == nil
}

// bump increments a counter stored in a marker file and returns the new value.
// A missing or unreadable file counts as zero, so the caller sees 1.
func bump(path string) int {
	n := 0
	if data, err := os.ReadFile(path); err == nil {
		n, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}
	n++
	_ = os.WriteFile(path, []byte(strconv.Itoa(n)), 0o600)
	return n
}
