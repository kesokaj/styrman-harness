---
owner: bugs
date: 2026-09-15
id: BUG-028
status: fixed
area: process (review)
---
# BUG-028: docs-lens finding applied without checking a recorded decision that already set the threshold

## Symptom
Round 1 of the Done review applied the docs lens's F5, adding a bug-fix clause to preflight Step 2's
release-diff reviewer at `skills/preflight/SKILL.md:50`. ADR-022's Consequences
(`.agents/docs/architecture/DECISIONS.md:179`) and plan row D8 keep the 50-line/5-file threshold for
the release gate's own review. Round 2's code lens caught the contradiction.

## Repro
Not reproduced as a running scenario; found by round 2 review reading `skills/preflight/SKILL.md:50`
against `DECISIONS.md:179`.

## Root cause
The finding (F5) was applied without checking the recorded decisions for the rule it touched — ADR-022
explicitly amends the bug-fix threshold only, and explicitly keeps the release-gate threshold as-is.

## Fix
Reverted: `skills/preflight/SKILL.md:50` matches commit `22c106c` again; F5 is rebutted in the merged
record.

## Prevention
Before applying a review finding, grep `DECISIONS.md` and the active plan's decision rows for the rule
it changes. A recorded decision makes the finding a rebuttal, not a fix.

## Related
BUG-025 and BUG-026 (same Done review) share the shape "a fix/finding applied without checking recorded
intent first (REQUIREMENTS.md, prior review rounds, ADR)". Three records in one review session sharing
this shape recommends the check — "before applying a finding or an edit, grep DECISIONS.md /
REQUIREMENTS.md for the rule it touches" — graduate into `context/AGENTS.md`'s review section or the
`review-standards` skill as a standing step, not left to each reviewer's judgement.
