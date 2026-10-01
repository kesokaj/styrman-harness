---
owner: bugs
date: 2026-10-01
id: BUG-131
status: fixed
area: dream collection
---
# BUG-131: after a clone every trail record looked new and the size cap dropped the transcripts

## Symptom
After a clone or checkout every trail record looked changed since the watermark, and the size cap then dropped the transcripts (the freshest input) to keep old records.

## Repro
`TestTrailDateFromGit`: a record committed before the watermark, then given a fresh mtime; it must not count as new.

## Root cause
`internal/dream/collect.go` dated trail files by `ModTime()`. A clone or checkout resets every mtime to the time of the write, so the filesystem date says nothing about when the record changed.

## Fix
`trailDates(root, docs)` (`internal/dream/collect.go:207`) runs one `git log --format=%ct --name-only -- docs` per project. A tracked, unmodified record is dated by its newest commit; an untracked or modified one keeps its mtime (`git status --porcelain --untracked-files=all` decides). In the working tree, committed with Phase 1.

## Prevention
"Changed since" for a file git tracks comes from git, never from the filesystem's mtime, which a clone resets.

## Related
Audit ID SD4 (`plans/2026-09-30-recon-audit-fix.md`); plan `plans/2026-09-30-audit-fix.md` Phase 1; verify log `testing/2026-09-30-verify-audit-fix-phase-1.md`. See BUG-135 (the status call this fix added).
