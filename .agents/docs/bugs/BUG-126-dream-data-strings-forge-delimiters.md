---
owner: bugs
date: 2026-10-01
id: BUG-126
status: fixed
area: dream prompt/proposals
---
# BUG-126: data strings could forge a proposal section or a prompt data block (SD3, SD8)

## Symptom
A proposal's free text, written unfenced into `proposals/dream/<date>-<project>.md`, could open a `### N.` / `**Status:** queued` section that `/retro` parses. An input or label holding `>>>` or `<<<data` (or a newline in a label) could close or open a data block in the prompt.

## Repro
`TestTextFieldCannotForgeASection`, `TestPromptEscapesFence`.

## Root cause
Strings from data were written next to the delimiters another parser reads, unescaped: the proposal text in `internal/dream/run.go` render, and the text, label and "Left out" line in `internal/dream/prompt.go`. It took review rounds 2 and 3 to reach labels and the summary line after the first fix covered only the inputs.

## Fix
Text is indented four spaces (`indent`, `run.go:458`), fence runs of either direction are spaced (`fenceRun`, `spaced`, `prompt.go:34`), labels are folded to one line (`oneLine`, `prompt.go:38`) in the prompt and the dry run.

## Prevention
Every string that comes from data and lands beside a delimiter (heading marker, fence, label) is escaped for that delimiter, labels and summary lines included. The test forges each delimiter from each position.

## Related
Audit SD3, SD8 in `plans/2026-09-30-recon-audit-fix.md`; reviews `reviews/2026-09-30-review-4764c5f-worktree-code{,-r2}.md`, `reviews/2026-10-01-review-4764c5f-worktree-code-r3.md`; plan `plans/2026-09-30-audit-fix.md` Phase 1.
