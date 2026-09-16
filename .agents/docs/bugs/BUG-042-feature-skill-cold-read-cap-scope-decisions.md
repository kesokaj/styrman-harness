---
owner: bugs
date: 2026-09-16
id: BUG-042
status: fixed
area: skills/feature (cold-read cap, scope, decisions)
---
# BUG-042: feature skill's cold-read cap, scope rule and decision rows all had gaps the simulation exposed

## Symptom
Three related defects in one run: (a) Claude Code's `/feature` ran a fourth cold-read round after
rounds 1-3 all ended "design needs revision" (commit 8c882a3 "cold-read ready in 4 rounds"), where the
skill's three-round cap is meant to be a gate that ends the work with a handoff question, headless;
(b) Antigravity's plan grew past the brief under "take the simplest reading" — three DELETE shapes, an
unrequested GET-by-id, flags plus environment variables with precedence rules — and the cold reader
never asked whether the request wanted any of it; (c) after security-driven changes, plan rows D2
("stdlib only"), D13 ("ID via crypto/rand") and D14 ("0755/0644") stood unchanged while the shipped
code had an external dependency and 0700/0600 permissions — describing text that outlived the
behaviour it described.

## Repro
The two sandbox repositories' plan files and cold-read reports
(`.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § Feature checklist row "Plan and cold
read", G8/G15).

## Root cause
The cap sentence counted "rounds without convergence", which a near-miss third round (still not both
tokens ready) reads as license for one more rather than a hard stop. Step 2 had no rule forcing an
interface ambiguity to resolve to one shape before planning, so "simplest reading" was read
permissively as "build every candidate". The cold-reader brief asked only for gaps and readiness, not
for scope conformance, so scope creep was invisible to the gate meant to catch it. Step 6 had no clause
saying a review fix that reverses an earlier decision must update that decision's row in the same
change.

## Fix
`skills/feature/SKILL.md`: the cap now counts reports, not rounds — "the third report that is not both
tokens ready ends the check — there is no fourth round"; Scope In holds only what the request names,
and an ambiguity is resolved to ONE shape recorded as a D-row with alternatives under Rejected; the
cold-read brief now asks three things, the second being scope conformance; a review fix that reverses a
D-row must write `supersedes D<n>` in the same change, with the row form added to the Decisions
template. `agents/reviewer.md` § Scope no longer hard-codes "two questions".

## Prevention
A gate must count something a print-mode (headless) run cannot argue its way around — reports, not
subjective "convergence". A design ambiguity is resolved by one recorded, named choice, never left open
for the agent to build every candidate under "simplest reading". A decision row that a later fix
reverses is updated in the same change that reverses it, not left standing as stale description.

## Related
Review: `.agents/docs/reviews/2026-09-16-review-bug-g15-feature-gates.md`. Testing:
`.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`. Sandbox re-run not owed here — the
G15 cap and G7 scope rule are verified by reading the fixed skill text, not by a second `/feature`
build, per `.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § Not verified here. Three bugs
from this session (BUG-039, BUG-040, BUG-041, BUG-042) each fixed a restatement gap across multiple
agent/skill files — recommend a standing reviewer check that greps every restatement of a changed rule
(echoes BUG-003/007/019/027); this rule is actionable (a grep-for-restatements check) and generalizable
beyond this session.
