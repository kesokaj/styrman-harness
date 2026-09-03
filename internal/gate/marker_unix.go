//go:build unix

package gate

import (
	"os"
	"syscall"
)

// ownedByUs reports whether the path is owned by the current user. A gate must
// not trust state it does not own.
func ownedByUs(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	return int(st.Uid) == os.Getuid()
}
