---
owner: bugs
date: 2026-09-13
id: BUG-007
status: fixed
area: docs/context
---
# BUG-007: repo docs told tests to use /tmp while the deployed hard rule forbids it

## Symptom
`AGENTS.md` (Setup commands, Rules, Testing) and `README.md` ("Working on the
harness") told every test to use `--home /tmp/sh-agent-test` / `/tmp/x`, while
the global hard rule this same repo deploys says "Scratch is
`.agents/scratch/` … never `/tmp`". Found by the 2026-09-13 research pass
(`.agents/docs/testing/2026-09-13-verify-research-pass.md`, finding 3).

## Repro
`grep -rn '/tmp' AGENTS.md README.md` at commit b64ffe8.

## Root cause
The hard rule was added to `context/AGENTS.md` on 2026-09-03; BUG-003 fixed
the same contradiction inside `context/dev-setup.md` at that time. But the
repo-level docs (`AGENTS.md`, `README.md`, `.agents/docs/deploy/local.md`)
were never grepped for the old `/tmp` wording — BUG-003's prevention rule
("grep every context file for the old wording") was scoped to files under
`context/` and not applied to the repo's own docs, which deploy that same
global context.

## Fix
`AGENTS.md` lines 13/34/46, `README.md` line 98, and
`.agents/docs/deploy/local.md` now say `.agents/scratch/<name>` (uncommitted
working tree, Phase 4 of
`.agents/docs/plans/2026-09-13-learning-loop-project-level.md`).

## Prevention
Widen BUG-003's rule: when a global rule changes, `git grep -n '<old
wording>'` across the WHOLE repo, not only `context/` — the repo's own
`AGENTS.md`/`README.md` count as context because the global context is
deployed from here.

## Related
BUG-003 (same contradiction, narrower scope).
