---
owner: bugs
date: 2026-09-18
id: BUG-076
status: fixed
area: internal/harness/repo.go
---
# BUG-076: `ResolveSource`'s not-found error hid the real precondition (no harness built in)

## Symptom
Operator, following README § Install without cloning on their Mac, ran `sh-agent install -vendor antigravity` and got:
`sh-agent: harness repo not found: run from inside the checkout, pass --repo, or set SH_AGENT_REPO`
They asked "Do I need the repo on my machine? it is supposed to get it from git?? I should just download the binary." — the message contradicted the feature and README, sending them to look for a checkout they had been told they would not need.

## Repro
`go build -o <scratch>/x ./cmd/sh-agent` from a checkout, then run `<scratch>/x install --home <scratch>/home` from a directory outside any harness tree with `SH_AGENT_REPO` unset. Old message as above. New message (verified): `harness repo not found and this build is e61db1f7acd6+dirty, a development build rather than a vX.Y.Z release, so it has no harness built in: run from inside the checkout, pass --repo, set SH_AGENT_REPO, or install from a release binary (README § Install without cloning)`.

## Root cause
The code was correct: the binary on the operator's PATH (`sh-agent 2615454169fb`) was a development build from a stale checkout, and D14 (2026-09-17 release-binary plan) refuses the embedded-harness fallback for any non-`vX.Y.Z` build on purpose, so a stale embedded copy is never installed silently. The defect was the error message: it named only the ways out (checkout, `--repo`, `$SH_AGENT_REPO`) and never said WHY it failed — that this particular build carries no embedded harness at all. It could not distinguish "wrong directory" from "this binary has nothing in it", which are opposite problems needing opposite fixes.

## Fix
`internal/harness/repo.go` `ResolveSource`'s not-found error now names which case applies (version stamp present but not a release, or no stamp at all), states the build "has no harness built in", quotes the version, and adds the release-binary way out. `TestResolveSource` (internal/harness/harness_test.go) asserts, for each of eight non-release version strings, that the message contains "no harness built in", mentions "release binary", and quotes the version; the old message fails all eight. Uncommitted at time of writing, committed in the same task as this record.

## Prevention
When a fallback or guard refuses because a precondition failed, the error names the precondition that failed, not only the alternatives the reader might try. "X not found: try A, B or C" is right only when X's absence is the whole story; when the code knows something more specific — this build cannot do it at all, this file is the wrong kind, this mode forbids it — say that first, since the reader's next action depends on which it is.

## Related
Same family as BUG-074 (a message stating a conclusion its measurement did not support). D14 in `.agents/docs/plans/2026-09-17-release-binary.md`.
