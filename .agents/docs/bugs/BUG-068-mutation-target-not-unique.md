---
owner: bugs
date: 2026-09-17
id: BUG-068
status: fixed
area: process / mutation testing (must-fail checks)
---
# BUG-068: a must-fail mutation anchored on non-unique text hit the wrong site and its restore didn't undo the damage

## Symptom
To prove BUG-067's new guard fails without it, the main agent scripted a text replacement of `strings.CutPrefix(st.Repo, harness.EmbeddedPrefix)`, but that construct appears twice in `main.go` (once in the `knownCheckout` guard, once in `status`). The script replaced the FIRST occurrence, in `status`, not the guard under test. The result was a build failure (`undefined: release`), not a failing test, and the restore step was itself a text replacement matching only the intended site — so the damaged `status` block survived the restore and the next run still failed to build.

## Repro
Observed during the release-binary Phase 2 checkpoint, on the mutation for the last entry of `.agents/docs/changes/2026-09-17-release-binary-phase-2.md` § Verification ("`knownCheckout`'s new guard removed → its test fails"). The two sites are `cmd/sh-agent/main.go:328` (`status`, which renders `from release vX.Y.Z`) and `main.go:483` (the `knownCheckout` guard); the script hit 328. Found by reading the build error (`undefined: release`) rather than a failing test, then repaired by hand, after which the mutation was redone against unique text and did fail its test as intended.

## Root cause
A mutation anchored on text that is not unique in the file, combined with a restore mechanism (a second text replacement) that cannot undo damage at a site it doesn't know it hit.

## Fix
The mutation was redone against unique text, the file restored, `gofmt`/`go vet`/`go test -race ./...` run clean, and `git diff` read to confirm only the intended lines remained.

## Prevention
A must-fail mutation asserts its target is unique before writing (count occurrences, fail if not exactly one). Undo it by restoring the saved copy of the file, never by a reverse text replacement. After restoring, read `git diff` to confirm nothing else changed.

## Related
Verification log: `.agents/docs/testing/2026-09-17-release-binary-phase-2.md`. Changes: `.agents/docs/changes/2026-09-17-release-binary-phase-2.md`.
