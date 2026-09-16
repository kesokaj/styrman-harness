---
owner: changelog
date: 2026-09-16
---
# Linkshelf simulation: both CLIs build one project under the harness  (2026-09-16)

## Summary
The harness was verified end-to-end by running both CLIs through `/init-project` and `/feature` from empty repositories in the sandbox, then comparing their results against the harness rules. Both CLIs closed the Definition of Done trail (plan, decisions, reviews, testing guides, changes, handoff) with no human reply required. 19 gaps found between harness and runs: 7 fixed under the bug skill (BUG-039..043), 6 filed as Follow-ups, 6 informational. Phase 3 re-runs and Phase 4 install confirmed fixes and closed all owed sandbox rows.

## Type
verification and process

## Files changed
- `.agents/docs/testing/2026-09-16-linkshelf-simulation.md` — verification log, 103 lines; § Status, Runs, Init/Feature checklists, Gaps, Re-runs
- `.agents/docs/testing/2026-09-16-verify-real-home-3.md` — install and smoke-check log, 55 lines; § Install/status, Smoke checks, Re-check
- `agents/verifier.md` § What to run (line 54) and § Records (line 97) — unconditional duty: "whenever guide.md does not exist … however many earlier logs `testing/` holds" (not "first time"); `-count=1` added to test counts
- `agents/documentation.md` § Your lane (line 35) — "whenever it runs … and finds the guide absent" (not "first time")
- `context/AGENTS.md` step 3 — "written by the verifier whenever it is absent" (not "created on first use"); § Git rule now names `settings.json` and `attribution` setting
- `skills/init-project/SKILL.md` Step 4, Step 5 template — "whenever it is absent" (not "first time"); Test block added; Configure row deleted when uncalled

## Agents dispatched
verifier ×2, reviewer ×1, bugs ×1, documentation ×1

## Verification
`.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § Status: "PASS — both CLIs built the service and closed the trail without a human reply; 19 gaps found (G1–G19), 7 fixed under the bug skill (BUG-039 to BUG-043), 6 filed as Follow-ups, 6 informational". `.agents/docs/testing/2026-09-16-verify-real-home-3.md` § Status: "PASS for both CLIs, with one FAIL from a verifier that was right". Install at a634c3d: `drift none`, `note:` printed by install output as specified in BUG-039. Reinstall from working tree (this commit): `drift none`, `binary a634c3d+dirty`. Every owed re-run row executed per testing log § Re-runs.

## Docs
Affected: agents/verifier.md, agents/documentation.md, context/AGENTS.md, skills/init-project/SKILL.md (paths, line ranges as listed above). Review record `reviews/2026-09-16-review-bug-043-guide-wording.md` found and fixed all restatements of pre-fix wording.

## Requirements
n/a
