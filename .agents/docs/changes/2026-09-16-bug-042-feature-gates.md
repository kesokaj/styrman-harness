---
owner: changelog
date: 2026-09-16
---
# BUG-042: feature skill's cold-read cap, scope rule and decision rows all had gaps the simulation exposed     (2026-09-16)

## Summary

Three related defects in the feature skill: (a) Claude Code's `/feature` ran a fourth cold-read round after rounds 1–3 all ended "design needs revision", though the skill's three-round cap is meant to gate the work — root cause was "rounds without convergence", which a near-miss third round reads as license for one more; (b) Antigravity's plan grew past the brief ("take the simplest reading") with unrequested DELETE shapes and GET-by-id, but the cold reader never asked whether the request wanted any of it — Step 2 had no rule forcing one shape per interface ambiguity; (c) after security-driven code changes, plan rows D2 ("stdlib only"), D13 ("ID via crypto/rand") and D14 ("0755/0644") stood unchanged while the code shipped an external dependency and different permissions. Fixed: cap now counts reports, not rounds — "the third report that is not both tokens ready ends the check — there is no fourth round"; Scope In holds only what the request names, and an ambiguity is resolved to ONE shape recorded as a D-row with alternatives under Rejected; the cold-read brief now asks three things with (2) scope conformance; a review fix that reverses a D-row must write `supersedes D<n>` in the same change, with that row form added to the Decisions template. Reviewer agent description updated to say "answer the questions the dispatch asks, in its numbering" (no longer hard-coded "two").

## Type

bug fix

## Files changed

- `skills/feature/SKILL.md` — Step 2: ambiguity resolved to ONE shape, Rejected alternatives recorded; Step 3 brief: three questions asked (1 Gaps, 2 Scope, 3 Readiness); gate counts reports not subjective convergence; Step 6 Decisions template: new row form for review fix that supersedes a D-row.
- `agents/reviewer.md` — § Scope description updated: answer the questions the dispatch asks, not hard-coded "two questions".

## Agents dispatched

reviewer (`reviews/2026-09-16-review-bug-g15-feature-gates.md` APPROVE WITH NITS; two findings fixed: `agents/reviewer.md` Scope description updated to match feature skill's new three-question brief, Decisions template added the `supersedes D<n>` row form), verifier (shared log: rendered context 93 lines, 27 skills 7 agents), bugs, changelog.

## Verification

Verifier log `.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`: rendered context output unchanged (93 lines). Skill text reads as internally consistent: the new Step 3 brief (three questions), the cap (third report ends check), and the Step 2 scope rule (one shape per ambiguity) all agree; no other living doc restates the old "two questions" / no-scope-check wording.

Not verified by re-run (see plan Phase 3 "Not verified here"): a `/feature` build at the fixed tree would cost an hour per CLI; the gap is visible from reading the corrected skill text and reading the control-run divergence (Claude Code's four rounds vs. the skill's three-round cap after the fix).

## Docs

`skills/feature/SKILL.md` (Steps 2, 3, 6) and `agents/reviewer.md` (§ Scope) updated to restate the three-question brief and scope rule consistently; staleness sweep found no other living doc carrying the old two-question wording (only the dated plan `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`, out of the living-docs set).

## Requirements

n/a
