package harness

import (
	"io/fs"
	"path/filepath"
	"testing"
)

// The Antigravity CLI builds ONE <user_rules> block from every file in
// ~/.gemini/config/rules/ and cuts it at this many bytes, marking the cut
// with "<truncated N bytes>" (BUG-085, measured 2026-09-21; the vendor page
// says 12,000 characters per file). A second file costs the first one room.
const antigravityRulesBlockBytes = 24000

// renderedContextBudgetBytes is the ceiling this repo pins itself to: an 8%
// margin under the observed cap (ADR-035). It is never raised by an
// implementer; a rule that does not fit is paid for by a deletion or a
// merge, or lives in a skill (REQUIREMENTS § Known limits).
const renderedContextBudgetBytes = 22000

// TestRenderedContextFitsAntigravityBlock renders THIS checkout the way
// `sh-agent install` does and measures the whole block Antigravity would
// build: the rendered context plus every rules/<name>.md, which deploys
// beside it. It reads the repository at ../.. like TestHarnessMap, so it
// fails — never skips — when run from anywhere else.
func TestRenderedContextFitsAntigravityBlock(t *testing.T) {
	src := DirSource(filepath.Join("..", ".."))
	if _, err := fs.Stat(src.FS, "context/AGENTS.md"); err != nil {
		t.Fatalf("this test reads the checkout at ../..: %v", err)
	}
	inv, err := LoadInventory(src)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := RenderContext(src, inv)
	if err != nil {
		t.Fatal(err)
	}
	rules := 0
	for _, r := range inv.Rules {
		rules += len(r.Raw)
	}
	total := len(ctx) + rules
	if total > renderedContextBudgetBytes {
		t.Fatalf("Antigravity rules block is %d bytes (context %d, rules %d), budget %d (the CLI cuts the block at %d): cut %d bytes",
			total, len(ctx), rules, renderedContextBudgetBytes, antigravityRulesBlockBytes, total-renderedContextBudgetBytes)
	}
	t.Logf("Antigravity rules block is %d bytes (context %d, rules %d), budget %d", total, len(ctx), rules, renderedContextBudgetBytes)
}
