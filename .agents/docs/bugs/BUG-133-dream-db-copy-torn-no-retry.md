---
owner: bugs
date: 2026-10-01
id: BUG-133
status: fixed
area: dream readers
---
# BUG-133: a conversation database copied mid-checkpoint failed the project (a torn global-index copy failed the whole run)

## Symptom
A copy of a conversation database taken while its writer was mid-checkpoint opened as NOTADB/CORRUPT and failed the project; a torn global-index copy failed the whole run.

## Repro
`TestCopyDBRetriesOnce` (fixture: the copy's 16-byte header zeroed; its first subtest proves the fixture yields NOTADB on the first query, since `sql.Open` is lazy) and `TestCopyDBNeverOpensSource`.

## Root cause
`internal/search/dream.go` copied db, -wal and -shm once, with no retry, and opened lazily, so a torn copy surfaced only at the first query, far from the copy.

## Fix
`openCopy(dir, name, tmp)` (`internal/search/dream.go:252`) copies, opens, forces one read, and on `isUnreadable` removes the three copies, waits `copyRetryWait` (a package var, zero in tests), then copies and opens once more. The live database is never opened. In the working tree, committed with Phase 1.

## Prevention
A read of a live SQLite database is a copy, retried once. A test of an open failure forces a read, since `sql.Open` does not touch the file; the fixture proves its condition in a subtest before the retry test is trusted.

## Related
Audit ID SD7 (`plans/2026-09-30-recon-audit-fix.md`); plan `plans/2026-09-30-audit-fix.md` Phase 1; the AGENTS.md rule that a test proves its unusual condition holds.
