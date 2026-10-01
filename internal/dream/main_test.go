package dream

import (
	"os"
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
	os.Exit(m.Run())
}
