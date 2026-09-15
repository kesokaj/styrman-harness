---
owner: main
date: 2026-09-15
---
# Harness practices Phase 7, part 3: the release gate's docs review is advisory     (2026-09-15)

## Summary

The release gate's docs routing drew review findings in four rounds (BUG-026). Under Definition of Done step 2 the operator was asked what the mechanism is for and approved the recommended option C ("Everything is approved. Go go go."): `/preflight` Step 4 still reviews the docs over `<base>..HEAD` and records each gap with severity, confidence and an owner, but no docs gap blocks the release; the Definition of Done's task-time `/review --docs` is where a docs finding blocks (plan D49, ADR-024, which supersedes ADR-023). A coverage-rule version tried before it was reviewed, drew four findings, and was never committed.

## Type

task

## Files changed

- `skills/preflight/SKILL.md` — :65 Step 4 runs the `--docs` lens over `<base>..HEAD` and records each docs gap, which is advisory, in the gate report with severity, confidence and an owner who fixes it as ordinary work; the rebuttal carry-forward from the newest prior gate report is gone. :77 "Any blocking finding means NO-SHIP."; the blocking list has no docs clause and says "Step 4's docs gaps never block". :81 Step 6 no longer says the SHIP path has no findings left.
- `.agents/docs/architecture/DECISIONS.md` — ADR-024, the release gate's docs review is advisory; ADR-023's heading marked "Superseded by ADR-024". ADR-024 was first appended by `changelog` for the coverage version and rewritten by the main agent for option C before any commit.
- `.agents/docs/bugs/BUG-026-preflight-step4-routing-open.md` and `INDEX.md` — fixed, with the round 4 findings and the option C resolution.
- `.agents/docs/reviews/2026-09-15-review-c514936-preflight-docs.md` — the docs-lens review of the uncommitted coverage version, with its Outcome; `2026-09-15-review-c514936-preflight-advisory-docs.md` — the docs-lens review of option C.
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — D49, task 7.7, and a Follow-ups line for a same-day gate report overwrite.

## Agents dispatched

reviewer (docs lens, Opus, twice: the coverage version and option C), bugs (BUG-026 reopened and closed), changelog (the first ADR-024 draft). The code change is a few lines in one file, below the reviewer threshold: the main agent read the diff and searched the living docs for wording of a blocking docs gate or the carry-forward and found none. This entry was written by the main agent.

## Verification

The option C docs-lens review reads APPROVE WITH NITS: four low findings, two fixed in Step 4's wording (a named place for the owner's work, and where a docs finding still blocks), one rebutted with ADR-024's trade-off, and this entry's placeholder. `grep -cF` in `skills/preflight/SKILL.md`: "docs gaps never block" `1`, "not already covered" `1` (Step 2 only), "not raised again" `0`, "high-severity docs finding" `0`, "no findings left" `0`; `^## ADR-024: The release gate` counts `1` in `DECISIONS.md`.

## Follow-ups

- A second `/preflight` run on the same day for the same version overwrites that day's gate report, losing a NO-SHIP report's Outcome (plan Follow-ups, found in 7.7).

## Requirements

n/a
