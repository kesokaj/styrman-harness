package harness

import (
	"io/fs"
	"path/filepath"
	"strings"
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

// renderedContextTargetBytes is the harness's own ceiling for the hand-written
// context (the render before its generated `## Installed` inventory), under the Antigravity block pin above: the audit of 2026-09-30
// found the context 94 bytes from that pin, so no rule could be added and no
// rules/ file deployed (plan audit-fix, C1, D2b). Adherence tracks the number
// of distinct instructions, not bytes, so the cut moved whole rules into the
// skill or agent that acts on them. Goal 14,500; this test fails above 15,000.
const renderedContextTargetBytes = 15000

func TestRenderedContextUnderTarget(t *testing.T) {
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
	// The hand-written part only: the generated `## Installed` inventory grows with
	// every skill dropped in, and a file drop must never fail this pin (R5); the
	// whole render, inventory included, stays under the Antigravity block in the
	// test above.
	written := string(ctx)
	if i := strings.Index(written, "\n## Installed"); i >= 0 {
		written = written[:i]
	}
	if len(written) > renderedContextTargetBytes {
		t.Fatalf("hand-written context is %d bytes, target %d (goal 14,500): move a whole rule into the skill or agent that acts on it, never shorten a sentence", len(written), renderedContextTargetBytes)
	}
	t.Logf("hand-written context is %d bytes (render %d with the inventory), target %d", len(written), len(ctx), renderedContextTargetBytes)
}
