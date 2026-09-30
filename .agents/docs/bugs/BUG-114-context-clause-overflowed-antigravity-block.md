---
owner: bugs
date: 2026-09-30
id: BUG-114
status: fixed
area: context/AGENTS.md (size budget)
---
# BUG-114: the context clause for BUG-113 overflowed the Antigravity rules block

## Symptom
Reviewer finding 1 (`.agents/docs/reviews/2026-09-30-1721-review-agy-only-init.md`): the clause "unless the harness is installed for Antigravity alone" at `context/AGENTS.md:38` made the rendered Antigravity block 22050 bytes; `TestRenderedContextFitsAntigravityBlock` (budget 22000, BUG-085) failed. HEAD was 21996.

## Repro
`go test ./internal/harness -run TestRenderedContext -v` on the tree with the clause.

## Root cause
The main agent ran only the new test (`go test -run TestInitProjectVendorCheckMatchesManifest .`) before dispatching review and verify, not `go test ./...`, while `context/AGENTS.md` had 4 bytes of headroom. The size test existed (BUG-085); it was not run before the edit was handed on.

## Fix
`context/AGENTS.md:38` now reads "with `.claude/*` linked into it by `/init-project` (Claude Code only)" and the same line drops "in a project, ". 21998 bytes, test passes. The first verifier run was stopped and re-dispatched after the fix.

## Prevention
An edit to `context/AGENTS.md` runs `go test ./internal/harness -run TestRenderedContext -v` and reads the byte count before any dispatch; with a few bytes of headroom, every added word is paid for by a cut in the same edit. The size rule already existed as a test (BUG-085); it was not run, so it was not followed.

## Related
BUG-085 (the size budget), BUG-113 (the change that added the clause). `context/AGENTS.md` is named by 53 records.
