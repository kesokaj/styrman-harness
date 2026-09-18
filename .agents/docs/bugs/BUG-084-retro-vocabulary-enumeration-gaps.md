---
owner: bugs
date: 2026-09-18
id: BUG-084
status: fixed
area: skills/feedback-log, agents/verifier
---
# BUG-084: retro-applied edits left a vocabulary enumerated in one place stale in another

## Symptom
The v1.3.0 gate's code review (`reviews/2026-09-18-review-v1.2.3-3b5b71a-code.md`, findings 1-3) found:
- 1358bcc added `redirection` to `skills/feedback-log/SKILL.md`'s trigger list and Signal Types table but not the JSON schema line an agent copies (`:48`, still `correction|override|restatement|format|gap`).
- The same commit rewrote `agents/verifier.md`'s early-log rule to make `Status: IN PROGRESS` load-bearing while the file's own log template (`:71`) still enumerated only `PASS | FAIL | BLOCKED`.
- The rule said "after the FIRST command you run", but the agent's Pre-flight runs `pwd`/`git status` before any verification command, so it was unsatisfiable as written.
The docs-lens review of the same edits (`reviews/2026-09-18-review-worktree-retro-docs.md`, finding 2) had already flagged the schema vocabulary gap; the fix that landed added only the table row, missing the schema line — a fix that took more than one attempt.

## Repro
`grep -n 'correction|override' skills/feedback-log/SKILL.md` at 1358bcc shows the schema line without the new value; `grep -n '^## Status' agents/verifier.md` shows the template without `IN PROGRESS`.

## Root cause
A value was added to a vocabulary (a signal name) in prose without grepping the file for every other enumeration of that same vocabulary (schema line, table). Separately, a rule named an ordinal step ("the first command") without checking the agent's own documented step order (Pre-flight runs first).

## Fix
94bec82: `skills/feedback-log/SKILL.md:48` gains `|redirection`; `agents/verifier.md:71` template becomes `IN PROGRESS | PASS | FAIL | BLOCKED`; `agents/verifier.md:18` and `context/AGENTS.md` § Work say "the FIRST verification command … once the docs root is resolved".

## Prevention
When a value joins a vocabulary (a signal, a status, a category), grep the file and its consumers for every enumeration of that vocabulary — schema lines, templates, tables — and change them in the same edit. A rule that names an ordinal step ("the first command") is read against the agent's own step order before it is written.

## Related
BUG-083, BUG-058. Third instance of "enumeration left stale" pattern with BUG-003/BUG-007/BUG-027; candidate for a standing reviewer/docs-lens check per BUG-027's note.
Graduated into `context/AGENTS.md` § Hard rules, "A behaviour change includes the text that DESCRIBES it" bullet, the new sentence beginning "Two shapes hide best: a value joining a VOCABULARY …".
