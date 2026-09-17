---
owner: bugs
date: 2026-09-17
id: BUG-070
status: fixed
area: cmd/sh-agent/main.go (update)
---
# BUG-070: `sh-agent update --check` deleted files despite its read-only contract

## Symptom
`update --check`, whose whole contract is to report and change nothing,
deleted files. The kept-clone cleanup (`harness.RemoveKeptClone`) was called
at the top of `update`, above both `checkOnly` returns, so a read-only
invocation removed `~/.config/sh-agent/feedback`.

## Repro
Security review (`.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md`,
finding 2): `RemoveKeptClone` is called at `cmd/sh-agent/main.go:415`, before
`Behind`, before the `checkOnly` early returns, and before the dirty-tree
refusal — so `sh-agent update --check` mutates the filesystem.

## Root cause
A side effect (the kept-clone cleanup, added in the same phase) was appended
to a function without checking where that function's read-only mode returns;
it landed above the `checkOnly` branch instead of below it.

## Fix
Commit `dab2565`. The `RemoveKeptClone` call moved below every `checkOnly`
return, beside the work that actually updates.
`TestUpdateKeptCloneCleanupOnlyWhenUpdating` asserts both halves: `--check`
leaves the clone and says nothing, the same update without `--check` removes
it and reports so. With the call moved back above the returns the test
fails, naming the removal.

## Prevention
When adding a side effect to a command that has a read-only or dry-run mode,
find that mode's early returns first and place the effect below them — and
the test for the effect asserts the read-only mode does NOT perform it, not
only that the normal mode does.

## Related
BUG-069 (same review, same phase, same class of destructive-side-effect
defect — recommend the shared rule "a destructive side effect is placed and
tested against every early-return mode of its caller" for
`context/AGENTS.md` § Hard rules if a third instance appears). Review:
`.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md`. Changes:
`.agents/docs/changes/2026-09-17-release-binary-phase-3.md`.
