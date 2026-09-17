package styrman_test

import (
	"bytes"
	"testing"

	styrman "github.com/kesokaj/styrman-harness"
	"github.com/kesokaj/styrman-harness/internal/harness"
)

// TestEmbeddedPlanEqualsCheckout: a release binary installs exactly what a
// checkout at the same commit installs — the same paths, bytes and modes, for
// every vendor. It fails when a harness directory is not embedded, when embed
// drops a file, or when a file's mode depends on the disk (embed records every
// file as 0444, so an executable file in the tree would deploy without its
// execute bit from the binary).
func TestEmbeddedPlanEqualsCheckout(t *testing.T) {
	for _, vendor := range []string{harness.VendorAll, harness.VendorClaude, harness.VendorAntigravity} {
		fromTree, err := harness.Plan(harness.DirSource("."), nil, vendor)
		if err != nil {
			t.Fatalf("%s: Plan from the tree: %v", vendor, err)
		}
		fromEmbed, err := harness.Plan(harness.Source{FS: styrman.Files}, nil, vendor)
		if err != nil {
			t.Fatalf("%s: Plan from the embedded files: %v", vendor, err)
		}
		if len(fromTree) == 0 {
			t.Fatalf("%s: the tree plans no files; the test ran in the wrong directory", vendor)
		}
		embedded := map[string]harness.File{}
		for _, f := range fromEmbed {
			embedded[f.Rel] = f
		}
		for _, want := range fromTree {
			got, ok := embedded[want.Rel]
			switch {
			case !ok:
				t.Errorf("%s: ~/%s is planned from the tree, not from the embedded files", vendor, want.Rel)
			case !bytes.Equal(got.Data, want.Data):
				t.Errorf("%s: ~/%s differs in bytes", vendor, want.Rel)
			case got.Mode != want.Mode || got.KeepForeign != want.KeepForeign:
				t.Errorf("%s: ~/%s mode %v keep %v from the embedded files, %v %v from the tree", vendor, want.Rel, got.Mode, got.KeepForeign, want.Mode, want.KeepForeign)
			}
			delete(embedded, want.Rel)
		}
		for rel := range embedded {
			t.Errorf("%s: ~/%s is planned from the embedded files, not from the tree", vendor, rel)
		}
	}
}
