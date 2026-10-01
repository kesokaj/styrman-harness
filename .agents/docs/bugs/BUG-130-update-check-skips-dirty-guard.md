---
owner: bugs
date: 2026-10-01
id: BUG-130
status: fixed
area: cmd/sh-agent update
---
# BUG-130: `update --check` / `--dry-run` on a dirty checkout reported success where the real run refuses (IN4)

## Symptom
On a dirty checkout, `--check` and `--dry-run` said "up to date" or "behind, would pull"; the real `update` refuses.

## Repro
`TestUpdateNeverInstallsDirtyTree`, subtests "--check on a dirty tree refuses" and "dirty and behind with --check".

## Root cause
`update()` in `cmd/sh-agent/main.go` returned on `checkOnly` (`main.go:719`) before the dirty check (`main.go:705`). The read-only mode skipped a guard the real run applies.

## Fix
The dirty check moved above both check-only returns. `--check`/`--dry-run` now exit 1 with the refusal on a dirty tree, 2 when clean and behind.

## Prevention
A test of every read-only mode runs it in the state the real run refuses and asserts the same refusal.

## Related
Audit IN4 in `plans/2026-09-30-recon-audit-fix.md`; plan `plans/2026-09-30-audit-fix.md` Phase 1. Recurs the global hard rule "a read-only mode resolves every guard the real run does" and the class of BUG-086 and BUG-122: the rule existed and was not followed.
