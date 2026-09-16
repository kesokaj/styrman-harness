---
owner: bugs
date: 2026-09-16
id: BUG-034
status: fixed
area: skills/preflight (release gate)
---
# BUG-034: the release gate's blocking list fails open on an incomplete security row

## Symptom
Security review of the a81ef4b mirror refresh (`styrman-crew/.agents/docs/reviews/2026-09-15-security-worktree.md`, finding 2) found `skills/preflight/SKILL.md:77` blocks on "a security finding of high severity", but not on a security row that Step 2 (`:51`) says is owed and that is missing, cut off at its agent's turn limit, or written by a shadowing agent (BUG-033). Preflight then deploys without handing back to the user.

## Repro
Not reproduced by exploit; read from the skill text. A release diff touching auth/IAM/secrets/input-parsing/dependencies triggers the Step 2 security dispatch requirement; if that dispatch never completes (turn limit, agent not invoked, or produces only a finding-free partial artifact — a normal outcome per the agents' early-write rule), Step 5's blocking list has no clause for "row owed and missing" and the merged report can still read VERDICT: SHIP.

## Root cause
The blocking list was written to check the content of a security report (its findings' severity) but never the report's existence or completeness against what Step 2 required. This plan (`.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`) deliberately narrowed the gate to "any blocking finding" and cut a broader "no release path unexamined" rule at the cold-read cap (D35); this record is the security-specific angle on that decision, not a re-litigation of it.

## Fix
`skills/preflight/SKILL.md` Step 5's blocking list gains "a security report that Step 2 owes and that is missing, has no `## Not examined` section, or names a path from the release diff under `## Not examined`" — mirroring the existing verifier-BLOCKED clause. The operator chose "Security row only", a narrow reversal of D35 scoped to that row; D52 in `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` records it, and D35's row now points at D52. No end-to-end run: staging a missing security report would mean sabotaging an agent, so this is confirmed by the text and by round-1 review (which found the clause missing its "names a changed path" half and fixed it) rather than by a live SHIP/NO-SHIP demonstration; BUG-033's `/preflight` scenario exercises the adjacent shadowing path. Verifier `.agents/docs/testing/2026-09-16-verify-bug-033-034-review-gate.md` PASS, counts pass 54, fail 0, skip 1. Code review `.agents/docs/reviews/2026-09-16-review-bug-033-034-code.md` REQUEST CHANGES then APPROVE WITH NITS (round 1 also found the check matching file names only and `/preflight` not saying stop before dispatching, both fixed; round 2's causal-qualifier nit dropped). Docs review `…-docs.md` APPROVE WITH NITS twice (D52 added, then D35 annotated).

## Prevention
Before trusting a preflight VERDICT: SHIP, check that every dispatch Step 2 required actually produced a complete artifact — a report a dispatcher decided not to write is not distinguishable, at the gate, from one that found nothing.

## Related
BUG-033 (the shadowing mechanism that can produce a clean-but-untrustworthy security row); D35 (cold-read cap decision that cut the broader coverage rule); styrman-crew `2026-09-15-security-worktree.md` finding 2, merged in `2026-09-15-review-worktree-merged.md`.
