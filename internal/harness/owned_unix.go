//go:build unix

package harness

import (
	"os"
	"syscall"
)

// ownedByCaller reports whether the file or directory belongs to the user
// this process runs as; no unix stat behind the FileInfo reads as owned.
func ownedByCaller(info os.FileInfo) bool {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return true
	}
	return int(st.Uid) == os.Geteuid()
}
