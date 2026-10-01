---
owner: bugs
date: 2026-10-01
id: BUG-134
status: fixed
area: dream transcripts
---
# BUG-134: harness-generated user content reached the dreamer's prompt

## Symptom
Slash-command output (`<local-command-stdout>`), `<command-message>`, `<task-notification>` results (474 on this machine), `<system-reminder>` text and `isMeta` lines (529) were read as the user's prose and sent to the dreamer.

## Repro
`TestSessionTextSkipsCommandOutput`, `TestSessionTextSkipsMetaLines`.

## Root cause
`contentText` (`internal/dream/collect.go:433`) and `sessionText` (line 390) took every user and assistant line as prose, with no list of the harness's own line shapes. Found by the audit (SD10) and widened by security rounds 2 and 3, which named the further prefixes.

## Fix
String content starting with those four prefixes is dropped, and lines with `isMeta: true` are skipped. In the working tree, committed with Phase 1.

## Prevention
A transcript reader lists the harness's own line shapes (meta flag, injected tag prefixes) and drops them; a new tag the harness starts writing is added to that list, with a test line per prefix.

## Related
Audit ID SD10; `reviews/2026-09-30-security-4764c5f-worktree-r2.md`, `reviews/2026-10-01-security-4764c5f-worktree-r3.md`; BUG-124 (same reader family).
