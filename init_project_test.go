package styrman_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kesokaj/styrman-harness/internal/harness"
)

// TestInitProjectVendorCheckMatchesManifest: /init-project skips every
// `.claude/` write on a machine where sh-agent is installed for Antigravity
// alone, and /dream then takes the store without the setting it would have
// read there; both tell that machine from the install manifest with one grep.
// The grep's needle and the manifest SaveState writes must agree: the needle
// matches the antigravity manifest and no other, or an Antigravity-only
// project gets a `.claude/` directory again, /dream finds no store there, or a
// Claude Code project loses its links.
func TestInitProjectVendorCheckMatchesManifest(t *testing.T) {
	const needle = `"vendor": "antigravity"`
	for _, path := range []string{"skills/init-project/SKILL.md", "skills/dream/SKILL.md"} {
		skill, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if want := "grep -q '" + needle + "' ~/.config/sh-agent/state.json"; !strings.Contains(string(skill), want) {
			t.Errorf("%s lacks the vendor check %s", path, want)
		}
	}
	for _, vendor := range []string{harness.VendorAntigravity, harness.VendorAll, harness.VendorClaude} {
		home := t.TempDir()
		if err := harness.SaveState(home, harness.State{Vendor: vendor}); err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(filepath.Join(home, ".config", "sh-agent", "state.json"))
		if err != nil {
			t.Fatal(err)
		}
		if got, want := strings.Contains(string(data), needle), vendor == harness.VendorAntigravity; got != want {
			t.Errorf("vendor %s: manifest matches %s = %v, want %v:\n%s", vendor, needle, got, want, data)
		}
	}
}
