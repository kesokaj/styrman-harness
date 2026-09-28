---
owner: main
date: 2026-09-28
---
# Retro rules: gate review fixes     (date: 2026-09-28)

## Summary
The v1.11.0 gate code review returned REQUEST CHANGES on the rules just applied. All of its findings are fixed.

**`agents/verifier.md`:**
- A long suite is redirected to a file with its own status line, never `tee`.
- A background job's end is checked with `kill -0`, not `jobs`.
- Shared calls are only for output read inline.
- `&&` is never used between the commands.

**`skills/feature/SKILL.md`:** the cold-read early exit is consistent with the gate list, the box text and the tick rule, and is triggered by `implementation ready`.

**The "put to the operator" outcome** is recorded `left: operator` in `context/AGENTS.md` and `skills/review/SKILL.md`.

**The verify log** now counts executed shared calls, not attempted ones.

## Type
fix

## Files changed
- `agents/verifier.md`
- `context/AGENTS.md`
- `skills/feature/SKILL.md`
- `skills/review/SKILL.md`
- `testing/2026-09-28-verify-retro-apply.md`
- `reviews/2026-09-28-review-v1.10.0-d3cc2b7-code.md` (with `## Outcome`)

## Agents dispatched
- **Gate:** `reviewer` (opus: code, docs), `verifier`.
- **After these fixes:** `reviewer` and `verifier`.

## Verification
- `go test -race ./...` is green.
- The context is 21,517 of 22,000 bytes (`TestRenderedContextFitsAntigravityBlock` passes).
