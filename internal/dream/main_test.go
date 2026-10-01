package dream

import (
	"os"
	"strings"
	"testing"
)

// TestMain clears CLAUDE_CONFIG_DIR for the package's tests. The collector
// honours it, as the CLI does, so a caller who has it set (a sandbox run, an
// operator who moved Claude Code's directory) made four collect tests read
// transcripts from somewhere other than their fixture home. A test that wants
// the variable sets it itself. Found by the first run of the harness-audit
// skill, 2026-10-01.
func TestMain(m *testing.M) {
	os.Unsetenv("CLAUDE_CONFIG_DIR")
	// And every GIT_* variable, as the other packages' TestMain does: the
	// collector runs git in fixture checkouts.
	for _, kv := range os.Environ() {
		if k, _, _ := strings.Cut(kv, "="); strings.HasPrefix(k, "GIT_") {
			os.Unsetenv(k)
		}
	}
	os.Exit(m.Run())
}
