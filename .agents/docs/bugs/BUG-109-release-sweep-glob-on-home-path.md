---
owner: bugs
date: 2026-09-30
id: BUG-109
status: fixed
area: internal/harness (release.go update sweep)
---
# BUG-109: `update` swept leftover release directories with a glob built from the home path

## Symptom
`internal/harness/release.go:144-148` swept leftover `release-*` download
directories with `filepath.Glob` on a pattern built from the home path; a home
path containing `[`, `*` or `?` matched entries in a sibling directory, and
every match was `os.RemoveAll`'d with no age check. Found by the project-wide
security review, finding 2.

## Repro
Home path `.../u[1]/.config/sh-agent`; the `[1]` in the path is read by
`filepath.Glob` as a character class, so the pattern matches `u1`'s entries
too and removes them. A second concurrent update could also delete the
first's in-flight directory since there was no age check.

## Root cause
`filepath.Glob` was used to enumerate directories under a path the operator
controls (the home), treating an operator-supplied path segment as pattern
syntax instead of a literal prefix to match by name. The two sibling sweeps
(`report.go` `sweepInterruptedClones`, `main.go` `sweepStaleZipDirs`) had
already been fixed for this exact class of defect; this third call site was
not found when that fix landed.

## Fix
`internal/harness/sweep.go` `SweepStaleDirs(parent, prefix, age, out)`:
`os.ReadDir` + prefix match by `Name()`, skips symlinks and non-directories,
and only removes entries older than `age` (one hour). Used from
`release.go` and `main.go` (replacing `sweepStaleZipDirs`). Tests:
`TestSweepStaleDirsStaysInParent`, `TestUpdateFromReleaseSweepStaysInHome`
(`internal/harness/release_test.go`).

## Prevention
A path built from the home directory or any operator-supplied path is never
passed to `filepath.Glob`; list the parent with `ReadDir` and match by name.
When a class of defect is fixed at one call site, `git grep -n
'filepath.Glob' | grep -v _test` before closing the fix, to find the other
call sites of the same class.

## Related
BUG-093 (refused command still creates state); review
`.agents/docs/reviews/2026-09-30-review-project-merged.md` finding 2, Round 1
and Round 2 ("age limit untested").
