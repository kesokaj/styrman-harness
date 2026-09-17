package styrman

import (
	"bytes"
	"io/fs"
	"os"
	"testing"
)

// TestFilesMatchTree: every regular file under the five harness directories is
// embedded with the same bytes — `.gitkeep` files included, which `embed`
// skips without `all:` (plan D2).
func TestFilesMatchTree(t *testing.T) {
	disk := os.DirFS(".")
	n := 0
	for _, root := range []string{"context", "skills", "agents", "commands", "rules"} {
		err := fs.WalkDir(disk, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			want, err := fs.ReadFile(disk, path)
			if err != nil {
				return err
			}
			got, err := fs.ReadFile(Files, path)
			if err != nil {
				t.Errorf("%s not embedded: %v", path, err)
				return nil
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%s embedded with different bytes", path)
			}
			n++
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if n == 0 {
		t.Fatal("no files walked; the test ran in the wrong directory")
	}
}
