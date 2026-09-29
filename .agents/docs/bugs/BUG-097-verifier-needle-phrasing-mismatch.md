---
owner: bugs
date: 2026-09-29
id: BUG-097
status: fixed
area: process (verifier dispatch — check needle vs. text under test)
---
# BUG-097: verifier grep needle didn't match the sentence it was checking

## Symptom
Round-2 verification (`.agents/docs/testing/2026-09-29-verify-dream-skill-phase-2-r2.md`, Status FAIL) ran `grep -c "distill and dream"` against the installed `skills/feedback-log/SKILL.md` and got 0. The file's actual sentence read "`/distill` and `/dream` records there too" (slash-command form), while the needle assumed the bare-noun form used in `skills/handoff/SKILL.md` and `context/AGENTS.md` ("retro, distill and dream").

## Repro
`grep -c "distill and dream" .agents/scratch/verify-dream-p2r2-home/.claude/skills/feedback-log/SKILL.md` → 0, vs. 1 for the same grep on `skills/handoff/SKILL.md`. See verifier log line 27 and Root cause section (line 34).

## Root cause
The main agent wrote both the target sentence in `skills/feedback-log/SKILL.md` and the acceptance-check grep it handed to the verifier without running the grep once, bare, on the working tree in the must-pass state first — so the phrasing mismatch (slash-command form vs. bare-noun form) across three files restating the same claim went undetected until the verifier ran it.

## Fix
`skills/feedback-log/SKILL.md`'s sentence aligned to the other two restatements ("Distill and dream records go there too."), confirmed by a bare `grep -c` returning 1, and round 3 (`.agents/docs/testing/2026-09-29-verify-dream-skill-phase-2-r3.md`) passed.

## Prevention
Before dispatching a verifier check built on a grep, the main agent runs that grep bare on the working tree in the must-pass state first — this is the existing global rule ("a check is validated once in the must-pass and once in the must-fail state before its verdict is cited") applied to the dispatch-authoring step, which was skipped here. This is a test failing for a reason other than the change under test and cost one extra verifier round.

## Related
`.agents/docs/testing/2026-09-29-verify-dream-skill-phase-2-r2.md`, `-r3.md`; `context/AGENTS.md` § Hard rules, "Silence from your own instrumentation" bullet (the must-pass/must-fail validation rule this extends to dispatch authorship).
