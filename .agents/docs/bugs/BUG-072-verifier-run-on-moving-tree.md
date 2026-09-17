---
owner: bugs
date: 2026-09-17
id: BUG-072
status: fixed
area: process/verifier-dispatch
---
# BUG-072: verifier dispatched against a tree the main agent kept editing

## Symptom
`verifier` was dispatched at commit 4d45155 (Phase 4). While it was still running,
the main agent edited `internal/harness/release.go` and `release_test.go` to fix
that phase's security-review findings. The verifier reported PASS (96 pass, 0 fail,
0 skip), but the result could not be attributed to any commit: not 4d45155, whose
content had moved under it, and not the fixed state, which did not exist when the
run started. The log (`.agents/docs/testing/2026-09-17-release-binary-phase-4.md`)
had to be marked inconclusive and replaced by a second, clean run against a frozen
tree at f15d533 (99 top-level tests, 85 subtests, 0 fail, 0 skip;
`.agents/docs/testing/2026-09-17-release-binary-phase-4-after-fixes.md`), which is
what the Phase 4 Checkpoint now rests on.

## Repro
Dispatch `verifier`; before it returns, edit a file within its scope; compare the
returned pass/fail counts against the commit the dispatch named — no commit's tree
matches what was actually compiled and tested.

## Root cause
The verifier reads the LIVE working tree by design (BUG-005 removed worktree
isolation for exactly this reason: to see uncommitted work). The main agent treated
"the dispatch is asynchronous" as licence to keep editing the same files while it
ran, destroying the correspondence between the artifact and any single commit.

## Fix
First log caveated with the defect and a pointer to the clean re-run; the plan's
Phase 4 Checkpoint row names both runs and states which one the phase rests on.
Not a code fix — the verifier and the security fixes were both correct; the first
verification run's evidence was simply unusable.

## Prevention
While a `verifier` dispatch is in flight, do not edit any file in its scope — keep
the tree frozen until the artifact lands; note anything found meanwhile and apply
it after. When a fix cannot wait, stop the run (TaskStop) rather than let it verify
a moving tree, then dispatch again after committing.

## Related
BUG-005 (verifier reads the live tree by design — that's the mechanism this bug
exploited by accident), BUG-013 (read a run's status before citing it as
complete). `.agents/docs/testing/2026-09-17-release-binary-phase-4.md`,
`.agents/docs/testing/2026-09-17-release-binary-phase-4-after-fixes.md`,
`.agents/docs/plans/2026-09-17-release-binary.md` Phase 4 Checkpoint.
