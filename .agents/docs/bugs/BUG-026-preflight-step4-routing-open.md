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

Phase 7 reopening: the D46/ADR-023 fix (`skills/preflight/SKILL.md:65`, `:77`) did not hold. The
Phase 7 reviews (`.agents/docs/reviews/2026-09-15-review-641eccd-worktree-code.md` and
`…-worktree-docs.md`) found three defects in the same mechanism:
1. `:65` carries a rebutted docs finding forward only from the newest prior gate report's
   `## Outcome`, and nothing writes it into the next report's Outcome. So one gate report later it
   is raised again and can block. A same-day re-run also overwrites the report that held the
   rebuttal, because the report path is `$(date +%F)-preflight-<version>.md`.
2. `:77` excludes a rebutted finding, but not a REQUEST CHANGES verdict that rests on it. The docs
   reviewer is never shown the prior rebuttals.
3. `:81` still says the SHIP path "has no findings left and needs none". Step 4 now records
   non-blocking gaps as `left` on the SHIP path.

This is the third round of review findings on this one mechanism (round 1, round 2, Phase 7). Under
Definition of Done step 2, patching stops and the operator is asked what the mechanism is for.

## Repro
Not reproduced as a running scenario; found by the Done review reading `skills/preflight/SKILL.md`
Step 4 and Step 6 against each other.

## Root cause
This release-gate routing was argued across three cold-read rounds (I-21a–c cut, I-21d–g resolved by
recommendation, plan D36) and patched again in review. Each patch moved the defect to a different line
instead of settling what the mechanism (Step 4 docs freshness) is for.

## Fix
Operator settled the mechanism's purpose on 2026-09-15 ("Do your recommendation"; plan row D46, ADR-023).
`skills/preflight/SKILL.md:65` now records each Step 4 docs finding with severity and confidence like a
Step 2 finding; a finding the newest prior `reviews/*-preflight-*` report's `## Outcome` rebutted, the
same statement at the same file:line, is cited as rebutted there and not raised again; a gap that blocks
is fixed after NO-SHIP as ordinary work and `/preflight` runs again; every other gap is recorded as
`left` with an owner. `skills/preflight/SKILL.md:77` now blocks Step 5 only on a Step 4 REQUEST CHANGES
verdict or a high-severity docs finding, not counting a finding a prior gate report rebutted. The
sentence "the gate never edits the tree after it has run" is gone, resolving R2-1 and R2-2 in
`.agents/docs/reviews/2026-09-15-review-c192622-22c106c-merged.md`.

Phase 7: not fixed. No patch was attempted for the three defects above; the question of what this
mechanism is for is with the operator.

## Prevention
When one mechanism draws findings in two consecutive rounds — a cold read or a review — stop patching
and settle its purpose with the operator before a third change.

## Related
BUG-025 shares the Done review; recurs with BUG-004/BUG-003/BUG-007-style "patch without checking
recorded intent" shape — see BUG-028's Related note on where this class of rule should live. The
stop-and-ask this bug's prevention rule calls for was settled by the operator in one decision (D46,
ADR-023) rather than a third patch — the prevention rule held for round 2, but the Phase 7 reviews
found the same mechanism drawing findings a third time; this is the third round, and the prevention
rule applies again: stop and ask, do not patch.
