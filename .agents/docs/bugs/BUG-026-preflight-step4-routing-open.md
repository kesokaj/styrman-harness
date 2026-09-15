---
owner: bugs
date: 2026-09-15
id: BUG-026
status: open
area: skills/preflight (Step 4 docs freshness)
---
# BUG-026: Step 4 docs-freshness routing patched twice, each patch introducing new defects

## Symptom
Edit I-21d made Step 4 fix skills, agents, rules, the context and `AGENTS.md` sections after the gate,
and those fixes rode into the tagged commit unverified (merged review F2). Round 1 fixed this at
`skills/preflight/SKILL.md:65`, `:77`, `:85`, `:90`: every Step 4 docs finding blocks, is fixed after
NO-SHIP, and `/preflight` runs again. That fix introduced two new defects:
- R2-1: a correctly rebutted docs finding blocks every release — Step 4 has no coverage rule and no
  rebuttal path (`:65`, `:77`, `:81`).
- R2-2: "the gate never edits the tree after it has run" contradicts Step 6's changelog entry and
  version bump.

## Repro
Not reproduced as a running scenario; found by the Done review reading `skills/preflight/SKILL.md`
Step 4 and Step 6 against each other.

## Root cause
This release-gate routing was argued across three cold-read rounds (I-21a–c cut, I-21d–g resolved by
recommendation, plan D36) and patched again in review. Each patch moved the defect to a different line
instead of settling what the mechanism (Step 4 docs freshness) is for.

## Fix
None yet. The review loop stopped at round 2 under Definition of Done step 2 (two consecutive rounds of
findings on the same mechanism); the question goes to the operator in the handoff.

## Prevention
When one mechanism draws findings in two consecutive rounds — a cold read or a review — stop patching
and settle its purpose with the operator before a third change.

## Related
BUG-025 shares the Done review; recurs with BUG-004/BUG-003/BUG-007-style "patch without checking
recorded intent" shape — see BUG-028's Related note on where this class of rule should live.
