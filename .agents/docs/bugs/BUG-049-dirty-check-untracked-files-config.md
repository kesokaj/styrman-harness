---
owner: bugs
date: 2026-09-16
id: BUG-049
status: fixed
area: internal/harness (dirty check)
---
# BUG-049: the BUG-046 dirty-checkout guard trusted git's configurable default

## Symptom
The dirty-checkout guard added for BUG-046 the same day used a bare `git
status --porcelain`, which inherits the operator's
`status.showUntrackedFiles`. The security review built a probe repository
and proved that under `status.showUntrackedFiles=no` a checkout holding a
wholly untracked new skill directory reports clean, so `update` would deploy
that never-committed skill to every vendor directory — the exact outcome
BUG-046's guard was written to stop.

## Repro
`git -C <probe> config status.showUntrackedFiles no`, add an untracked
`skills/foo/` directory, run the pre-fix `Inspect`; it reports the tree
clean.

## Root cause
The guard was built on a git command whose output depends on operator
configuration outside the repo. The fix written for BUG-046 the same day
reproduced its own failure mode in a narrower case: a config key the guard
never pinned.

## Fix
`internal/harness/repo.go:94` passes `--untracked-files=all` explicitly;
`TestInspectCountsUntrackedFiles` sets `status.showUntrackedFiles=no` in a
fixture repo and asserts the tree still reads dirty.

## Prevention
A guard never trusts a git command's default when a config key can change
it — pass the flag that pins the behaviour, and test with that key set
against you.

## Related
BUG-046 (the guard this narrows, same day); the "validate a check in the
state it must fail as well as the state it must pass" family (BUG-029,
BUG-024, BUG-031). Third occurrence of a same-day fix reproducing its own
bug's shape (see BUG-026); worth a standing reviewer check on any fix landed
without a second pass. Graduated to `TestInspectCountsUntrackedFiles`.
