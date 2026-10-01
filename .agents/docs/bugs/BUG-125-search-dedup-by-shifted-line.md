---
owner: bugs
date: 2026-10-01
id: BUG-125
status: fixed
area: search
---
# BUG-125: `search: 2 hits ..., 1 with every word` printed the same chunk twice (SD2)

## Symptom
The same chunk appeared twice in one result list.

## Repro
A file whose matched term is on line 5 of a chunk that starts at line 1, queried with two words.

## Root cause
`find()` in `internal/search/search.go` keyed its dedup set by the hit's line AFTER `matchOffset` shifted it. The any-word round compared the chunk's unshifted line, so the keys never matched.

## Fix
`hit.id` (the chunk id) is the dedup key (`internal/search/search.go`, `find`, line 528 on). `TestNoDuplicateHitAcrossRounds`.

## Prevention
A dedup key names the entity (the chunk), never a derived display value (the printed line).

## Related
Audit SD2 in `plans/2026-09-30-recon-audit-fix.md`; plan `plans/2026-09-30-audit-fix.md` Phase 1; verify log `testing/2026-09-30-verify-audit-fix-phase-1.md`.
