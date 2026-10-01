---
owner: bugs
date: 2026-10-01
id: BUG-128
status: fixed
area: internal/harness/apply.go (install manifest)
---
# BUG-128: a write that failed midway through `install` left files no manifest listed (IN2)

## Symptom
Files written before the failure stayed on disk, unlisted, so `uninstall` never removed them. The same after the resync loop's removal error.

## Repro
`TestInstallMidFailureKeepsManifestComplete`, `TestInstallResyncFailureKeepsManifest` (an unwritable or unremovable directory).

## Root cause
`Install` in `internal/harness/apply.go` saved the manifest only at the end. Two early returns (write error, stale-removal error) skipped the save.

## Fix
`saveOnError` (`apply.go:37`) saves (old entries still on disk and not links) union (files seen so far) once; both returns go through it (`apply.go:113`, `apply.go:146`). Found complete after code review round 2.

## Prevention
In a function that writes several files and records them, every return after the first write goes through the one path that saves the record. Grep the function for `return` before calling the fix done.

## Related
Audit IN2 in `plans/2026-09-30-recon-audit-fix.md`; reviews `reviews/2026-09-30-review-4764c5f-worktree-code-r2.md`; plan `plans/2026-09-30-audit-fix.md` Phase 1.
