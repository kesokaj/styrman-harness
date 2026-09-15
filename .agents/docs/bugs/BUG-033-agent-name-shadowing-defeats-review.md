---
owner: bugs
date: 2026-09-16
id: BUG-033
status: open
area: context/agents (precedence)
---
# BUG-033: agent-name shadowing defeats the independent review rows

## Symptom
Security review of the a81ef4b mirror refresh (`styrman-crew/.agents/docs/reviews/2026-09-15-security-worktree.md`, finding 1) found `context/AGENTS.md:37` states precedence between a project and harness agent of the same name "is not guaranteed" and that a project agent/skill "never reuses a harness name" — a convention, not a check. The G26 run (`.agents/docs/testing/2026-09-15-verify-g26-precedence.md`, D41) measured Claude Code 2.1.272 loading a project `reviewer` over the harness one.

## Repro
Not reproduced by exploit; measured mechanism: add `.claude/agents/security.md` (or `.agents/agents/security.md` behind the `/init-project` link) naming itself `security`/`reviewer`/`verifier`/`explorer`; it wins the name clash per D41's measurement. `/preflight` Step 2 and `/review` Step 3 dispatch by name with no clash check first, so the shadowing agent can write its own clean report and the release ships.

## Root cause
Precedence rules were written as prose guidance for authors ("never reuses a harness name") instead of as a check any dispatcher runs before trusting an independent agent's report. `/preflight` and `/review` assume the name they dispatch resolves to the harness definition.

## Fix
Not fixed. Remediation named by the finding: `/preflight` (Step 2) and `/review` (when it dispatches a review agent) refuse to run while the project layer holds an agent file named like a harness review agent (`reviewer`, `security`, `verifier`, `explorer`), naming the file; correct the `review` skill description's claim of "the harness's own reviewer, security, verifier, and explorer agents". Not reachable in styrman-crew today (`.agents/agents/` holds only `.gitkeep`).

## Prevention
Before dispatching `reviewer`, `security`, `verifier`, or `explorer` for an independent judgement, list `.claude/agents/` and `.agents/agents/` (and the diff's changed paths there) for a name clash with a harness review agent; a clash stops the dispatch and is blocking in `/preflight` Step 5, not just Step 2.

## Related
BUG-034 (same review, the gate that should have caught a missing/shadowed security row); D41, G26 (`.agents/docs/testing/2026-09-15-verify-g26-precedence.md`); styrman-crew `2026-09-15-security-worktree.md` finding 1, merged in `2026-09-15-review-worktree-merged.md`.
