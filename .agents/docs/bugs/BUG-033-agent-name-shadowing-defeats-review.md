---
owner: bugs
date: 2026-09-16
id: BUG-033
status: fixed
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
`skills/review/SKILL.md` Step 3 and `skills/preflight/SKILL.md` Step 2 now check, before any dispatch, for a project agent shadowing a harness review agent — a file under `.claude/agents/` or `.agents/agents/` named `reviewer.md`, `security.md`, `verifier.md` or `explorer.md`, or whose frontmatter `name:` is one of those four. `/review` stops before dispatching, names the file, says the review cannot run until it is renamed or removed, and writes nothing; `/preflight` stops before any dispatch and goes to Step 5 with NO-SHIP, the file named in the gate report, and Step 5's blocking list names the shadowing agent. Matrix rows `review-shadow`, `review-shadow-name` and `preflight-shadow` added to `.agents/docs/testing/skill-matrix.md`. The `review` skill's description claim (the review runs through the harness's own reviewer/security/verifier/explorer agents) is left as-is since the check now enforces it. Verified in the sandbox (this tree, Claude Code 2.1.272, clones of the demo repository): a committed `.agents/agents/security.md` before the fix caused `/review --security` to dispatch the project copy, whose transcript carried the planted project prompt; after the fix the same scenario made no dispatch and named the file. A `.agents/agents/sec-shadow.md` carrying `name: security` was caught by the frontmatter arm. `/preflight` against the committed shadow returned `VERDICT: NO-SHIP` in `.agents/docs/reviews/2026-09-16-preflight-a2d9a9c.md`, naming the file, dispatching no agent, committing/tagging/pushing nothing. A control clone with no shadow made one `security` dispatch and the harness agent's report carried a `## Not examined` section. Verifier `.agents/docs/testing/2026-09-16-verify-bug-033-034-review-gate.md` PASS, counts pass 54, fail 0, skip 1. Code review `.agents/docs/reviews/2026-09-16-review-bug-033-034-code.md` REQUEST CHANGES then APPROVE WITH NITS; docs review `…-docs.md` APPROVE WITH NITS twice. Whether Claude Code resolves a dispatch by file name or by frontmatter `name:` stays unmeasured beyond this; the check covers both.

## Prevention
Before dispatching `reviewer`, `security`, `verifier`, or `explorer` for an independent judgement, list `.claude/agents/` and `.agents/agents/` (and the diff's changed paths there) for a name clash with a harness review agent; a clash stops the dispatch and is blocking in `/preflight` Step 5, not just Step 2.

## Related
BUG-034 (same review, the gate that should have caught a missing/shadowed security row); D41, G26 (`.agents/docs/testing/2026-09-15-verify-g26-precedence.md`); styrman-crew `2026-09-15-security-worktree.md` finding 1, merged in `2026-09-15-review-worktree-merged.md`.
