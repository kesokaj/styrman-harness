package harness

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SweepStaleDirs removes the directories directly under parent whose names
// start with prefix and were last modified more than age ago: the leftovers
// of an interrupted update or zip install. Names come from ReadDir and match
// by prefix, never from a glob built from the home path, which reads `[`,
// `*` and `?` in the home as pattern and reaches a sibling directory
// (security review 2026-09-30, finding 1). A file, a symlink and a directory
// younger than age are left alone: the young one may belong to a run still
// in flight.
func SweepStaleDirs(parent, prefix string, age time.Duration, out io.Writer) {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), prefix) || e.Type()&fs.ModeSymlink != 0 || !e.IsDir() {
			continue
		}
		d := filepath.Join(parent, e.Name())
		info, err := os.Lstat(d)
		if err != nil || time.Since(info.ModTime()) <= age {
			continue
		}
		if err := os.RemoveAll(d); err != nil {
			fmt.Fprintf(out, "warning: could not remove the leftover %s: %v\n", d, err)
			continue
		}
		fmt.Fprintf(out, "removed %s, left by an interrupted run\n", d)
	}
}
