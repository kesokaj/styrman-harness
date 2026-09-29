---
owner: bugs
date: 2026-09-29
id: BUG-105
status: open
area: skills/preflight (release gate verdict)
---
# BUG-105: gate report's VERDICT: SHIP contradicted its own "What was not checked" section

## Symptom
`.agents/docs/reviews/2026-09-29-preflight-v1.12.1.md` wrote `VERDICT: SHIP`, with the verdict paragraph asserting `AGENTS.md` § Release's rule (a gate naming a missing Antigravity run of a changed skill is NO-SHIP) "is met" by a Claude Code run of the changed skill (`init-project`, linked-store row). The same report's "What was not checked" section correctly listed that `/init-project`'s changed text had never run on Antigravity. v1.12.1 was tagged on that verdict.

## Repro
Read `.agents/docs/reviews/2026-09-29-preflight-v1.12.1.md`: the verdict paragraph and the "What was not checked" list name the same missing run and reach opposite conclusions.

## Root cause
The verdict paragraph was composed from the rows already satisfied (a Claude Code run existed) rather than against the rule's exact condition (a run ON Antigravity); "What was not checked" was compiled afterward and never read back against the verdict already written.

## Fix
Not fixed as a rule. The missing run was performed post-release (PASS, 1 of 1, `.agents/docs/testing/2026-09-29-verify-real-home-v1.12.1.md`) and an `## Addendum` was appended to the gate report noting the verdict was overstated. The release itself needed no correction, but no file yet stops the next gate report from writing a verdict before checking it against its own "What was not checked" list.

## Prevention
Before writing a VERDICT line, read every item in "What was not checked" against every blocking condition (including project-added ones like `AGENTS.md` § Release's Antigravity-run rule); an item a condition names forces NO-SHIP. Write the verdict last, after that section is final, never before it. Candidate home: `skills/preflight/SKILL.md` Step 5.

## Related
BUG-083 (claiming "works" from an unproven run), BUG-085. Three related process/verdict-accuracy bugs (BUG-083, BUG-085, BUG-105) now share the theme "a gate's stated conclusion must be checked against its own evidence list before it is written" — recommend this graduate into `skills/preflight/SKILL.md` Step 5 as an explicit ordering rule, actionable (a re-read-before-write check) and generalizable beyond this release.
