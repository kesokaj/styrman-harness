---
owner: bugs
date: 2026-09-21
id: BUG-087
status: fixed
area: process / behavioural runs of agents/*.md
---
# BUG-087: a behavioural-run fixture did not hold the condition its dispatch claimed, twice

## Symptom
Behavioural runs of the changed `agents/bugs.md` (Phase 2, same plan) used demo-service clones whose dispatch said "the commit at HEAD carries the line agentId: a1b2c3". The fixture's HEAD commit had no such line: in 2 of 3 runs the main agent checked `git log -1 --format=%B`, found nothing, and asked the user instead of dispatching; the third dispatched and the record said "not reproduced at HEAD". A second fixture that did have the line in the commit body was also wrong for the rule under test: the new bullet only bites when the recommended permanent home is `context/AGENTS.md`, and in a project repo the bugs agent never proposes the global context, so that arrangement could not exercise the rule at all.

## Repro
`git log -1 --format=%B | grep -c 'agentId: a1b2c3'` on the first fixture's HEAD returns 0, not the 1 the dispatch assumed.

## Root cause
The arrangement was built from the dispatch text, not from what the rule measures. Global hard rule "a test that sets up an unusual condition proves the condition HOLDS before the test is trusted" and BUG-081's "the test builds the ARRANGEMENT the finding measured" both apply, and neither was checked before running.

## Fix
Replaced both fixtures with three clones of this harness repo itself (where graduated rules actually live in the context), subject a real failure of the day (BUG-086); 3 of 3 runs then named the existing sentence to widen.

## Prevention
Before a behavioural run, prove the fixture's claimed state with the same command the agent will use (e.g. `git log -1 --format=%B | grep -c <line>` → 1), and choose the repository kind in which the changed rule can actually fire — a rule about the global context runs in a clone of the harness, not of a project.

## Related
BUG-071, BUG-081, BUG-083 — a fourth instance of "fixture didn't hold the condition it claimed"; candidate for a standing pre-run check alongside BUG-081's arrangement rule in `context/AGENTS.md` § Hard rules. Found in `.agents/docs/plans/2026-09-21-context-consolidation.md` Phase 2.
