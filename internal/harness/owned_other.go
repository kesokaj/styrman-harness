//go:build !unix

package harness

import "os"

// ownedByCaller has no ownership to read here; owned is the full remedy.
func ownedByCaller(os.FileInfo) bool { return true }
