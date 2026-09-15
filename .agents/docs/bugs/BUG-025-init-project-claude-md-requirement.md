---
owner: bugs
date: 2026-09-15
id: BUG-025
status: fixed
area: skills/init-project
---
# BUG-025: init-project's "keep content, add `@AGENTS.md`" edit left CLAUDE.md not a one-line import

## Symptom
Edit I-23b added "on no, keep the content and add `@AGENTS.md` as its first line" to
`skills/init-project/SKILL.md:44`. That leaves a project `CLAUDE.md` that is not the one-line import
`REQUIREMENTS.md:11` (R2) and `.agents/docs/architecture/DECISIONS.md:41` require. Found by the Done
review, `.agents/docs/reviews/2026-09-15-review-c192622-22c106c-merged.md` finding F1, both lenses.

## Repro
Not reproduced as a running scenario; found by review reading `SKILL.md:44` against `REQUIREMENTS.md:11`.

## Root cause
The plan's edits went through three cold reads and a rehearsal, but no step compared an edit against
`REQUIREMENTS.md`, which changes only with the operator's approval of the exact edit — the check that
would have caught the contradiction was never run.

## Fix
`skills/init-project/SKILL.md:44` now leaves the file untouched when the user declines to fold it. Round
2 review found that line's new wording still inaccurate (R2-3); left open with the operator's R2
question in the handoff. Operator kept R2 on 2026-09-15 ("Yes. We should keep the AGENTS.md style.";
plan row D45). R2-3 is now fixed at `skills/init-project/SKILL.md:44`: on "no", the summary says
Antigravity never reads the kept file's content, and Claude Code imports `AGENTS.md` only if the file
has the line `@AGENTS.md`.

## Prevention
A plan's cold-reader check and rehearsal compare every edit against `REQUIREMENTS.md`; an edit that
changes a requirement's invariant stops at the operator before it lands.

## Related
BUG-026 (same Done review, same mechanism: a fix applied without checking recorded decisions/requirements).
