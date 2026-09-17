---
owner: bugs
date: 2026-09-17
id: BUG-069
status: fixed
area: internal/harness/report.go (RemoveKeptClone)
---
# BUG-069: `RemoveKeptClone`'s unpushed-commit guard fails open, deleting the operator's only copy

## Symptom
`RemoveKeptClone`'s own comment says "a commit that never reached the remote
is the operator's only copy", but the guard was
`if ahead, err := gitOut(dir, "rev-list", "--count", "@{upstream}..HEAD"); err == nil && ahead != "0"`.
When `@{upstream}` cannot be resolved (no tracking branch, detached HEAD, a
renamed branch), `gitOut` errors, the condition is false, and
`os.RemoveAll(dir)` runs on a checkout holding unpushed commits.

## Repro
Security review probe (`.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md`,
finding 1): a scratch home with an unpushed commit on a branch with no
upstream. `git -C … rev-list --count '@{upstream}..HEAD'` prints
`fatal: no upstream configured for branch 'main'`; `RemoveKeptClone` then
prints "removed the kept feedback clone" and the directory and its commit are
gone.

## Root cause
The guard was written as "keep it when I can prove there is something to
lose" instead of "remove it only when I can prove there is not". An error
from the safety check was treated the same as "nothing to protect" rather
than as "cannot tell, so protect it".

## Fix
Commit `dab2565`. Count over all local refs against all remote-tracking refs
(`rev-list --count --all --not --remotes`) and keep the clone whenever the
count cannot be taken at all; also refuse a symlink at that path, treat a
stash (`git stash list`) and ignored files (`status --ignored`) as work, and
resolve a relative `--home` with `filepath.Abs`. `TestRemoveKeptClone`'s
subtests (unpushed work on a branch with no upstream, an uncountable clone
via unreadable `.git/packed-refs`, stashed work, a symlink, a relative home)
were each seen failing before the fix.

## Prevention
A guard standing in front of a destructive action is written so that an
ERROR from the check keeps the data — `if err != nil` leaves, never falls
through to the delete — and its test suite includes the case where the check
itself cannot run, not only the cases where it runs and answers.

## Related
BUG-070 (same review, same phase, second destructive-side-effect defect);
BUG-036 (symlink refusal pattern reused here). Review:
`.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md`. Changes:
`.agents/docs/changes/2026-09-17-release-binary-phase-3.md`.
