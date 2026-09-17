---
owner: bugs
date: 2026-09-17
id: BUG-066
status: fixed
area: embed.go / release build
---
# BUG-066: `go:embed all:` compiles untracked or gitignored files into the release binary

## Symptom
`//go:embed all:context all:skills all:agents all:commands all:rules` reads the WORKING TREE, so any untracked or gitignored file sitting in those five directories at build time is compiled into a released binary and installed on every machine, though no loader reads it and review never saw it.

## Repro
`printf 'TOKEN=canary\n' > rules/.env` (ignored by `.gitignore:12`), rebuild: the embedded set grew to hold `rules/.env`, while `TestFilesMatchTree` and `TestEmbeddedPlanEqualsCheckout` both stayed green, because such a file is carried but never planned or installed. (Reported as finding 1 of `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`, which measured the `go list` EmbedFiles behaviour on a probe module under `.agents/scratch/sec-review/`; the `rules/.env` canary above was then run in THIS repository by the main agent, and the file removed again.)

## Root cause
The only embed test compared disk ⊆ embed, never embed ⊆ git. `go:embed` has no concept of gitignore, so a secret-shaped or stray file on disk in the embedded directories ships silently — it was never checked against what the release is actually supposed to contain.

## Fix
`TestEmbedHoldsOnlyTrackedFiles` in `embed_test.go` pins the embedded set to `git ls-files` in both directions (commit 73f2c36). The release-time half (refusing a release build when anything is untracked or ignored under the five directories) is a Phase 5 task box, decision D23, not yet built.

## Prevention
A test that a generated or embedded artifact matches a source must assert BOTH directions — extra content is as much a defect as missing content. An artifact built from the working tree, not from git, must be pinned to what git tracks, not just checked for completeness.

## Related
BUG-067 (same review, finding 2). Review: `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`. Changes: `.agents/docs/changes/2026-09-17-release-binary-phase-2.md`.
