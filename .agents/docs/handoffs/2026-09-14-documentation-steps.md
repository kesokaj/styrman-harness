---
owner: main
date: 2026-09-14
---
# Handoff — 2026-09-14 — the ten documentation steps shipped on both CLIs

Session ID: session_01CtMYM9AGp33edxBQGURNaB
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean once this handoff's commit is pushed

## Context

The operator's `/btw` review produced `.agents/docs/plans/2026-09-14-documentation-steps.md` ("so the agent has total knowledge of the entire app it is building"); the operator asked for it as a feature and for the residue of the cold reads to be resolved by recommendation. Supersedes `2026-09-14-bug-012-fixed.md`.

## Completed

- Phase 1 (d87df30): step 3 names the files the docs step owns and says a new command, flag, endpoint, or deploy step is not done until its file says so; step 2 runs `/review --docs` on interface changes; every agent carries a `## Records` section; the changelog template carries `## Docs`. `.agents/docs/changes/2026-09-14-documentation-steps-phase-1.md`.
- Phase 2 (this commit): the verifier writes `testing/<app>/guide.md` on first use (it wrote this repo's `.agents/docs/testing/sh-agent/guide.md`); review and gate records get an `## Outcome` section from the task that acts on them; retro proposals carry a status from day one and `/retro --harness` lists the open ones and pending `REQUIREMENTS.md` edits first (dry-run verified against the real queue: four queued 2026-09-04 proposals, one pending Known-limits sentence); the handoff template names pending requirements edits; deploy guides carry Operator steps; the skill matrix has a `documentation` row; the six queue files are backfilled with statuses. `.agents/docs/changes/2026-09-14-documentation-steps-phase-2.md`; ADR-018.
- Cold read: two rounds of three readers, residue resolved by recommendation and recorded in the plan. Reviewer APPROVE WITH NITS, three fixed. Verifier PASS.
- The documentation matrix row: PASS on both CLIs for the behaviour it proves (`testing/2026-09-14-verify-docs-matrix-row.md`); Claude Code's run exhausted its 25-minute budget before its own commit.

## Open / In Progress

Nothing queued in this repo; `TODO.md` and `ROADMAP.md` are empty. Operator-only, unchanged: the Mac (`sh-agent update`, three smoke checks; Claude Code there never checked); Go ≥ 1.24.6 on both machines; the downstream mirrors in styrman-crew and argolis-devbox (`deploy/remote.md`, Operator steps); improvement-review item 20 after the first Mac retro; the full skill matrix.
- Requirements edits awaiting approval: in `REQUIREMENTS.md` Known limits, replace "The Mac build is confirmed at the operator's next `sh-agent update` there (see the rebuild plan)." with "Confirmed on the Mac for Antigravity on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); Claude Code on the Mac is unverified."
- Still queued in the harness queue (`/retro --harness` lists them): 2026-09-04 P1 (retro before handoff), P2 (a behaviour change updates the text that describes it), P3 (check generated files are text; `git ls-files` before `rm -r`), P4 (every new test can fail). Decide apply or reject at the next harness retro.

## Key Decisions

- Items 6 and 7 kept at the operator's word although the check had recommended striking them; item 3 applies per task and in `/preflight`.
- The demo fixtures gained deploy guides; the matrix row's task is a setting the fixture lacks (`SHUTDOWN_TIMEOUT`), because a row must prove a change, not find it already present.
- A headless task through the full Definition of Done can exceed 25 minutes on Claude Code; accepted, recorded, not worked around.

## Corrections

- "DoD is when todo and roadmap is done plus the documentation is up to date and all the instructions are there. And ofc everything is pushed to git." — feedback log 2026-09-14 (restatement, `AGENTS.md § Definition of Done`); this feature is its implementation.

## Agent notes (unverified)

- Claude Code's matrix run wrote two bug records in the demo about its own headless permissions (a denied live-process test); fixture noise, left in the clone.
- The Antigravity run also edited the demo's `AGENTS.md`; harmless.

## Environment State

- Real home on this box: the operator asked for `sh-agent update` after this commit; it is run right after the push. Sandbox home at the working tree; clones `/home/kesokaj/sandbox/runs/{claude,agy}-docs2`; evidence `.agents/scratch/e2e-docs2/`, `e2e-docs3/`. No running processes.

## Artifacts to Read

- `.agents/docs/plans/2026-09-14-documentation-steps.md` — every box ticked, the residue sections
- `.agents/docs/testing/2026-09-14-verify-docs-matrix-row.md` — both CLI runs of the new row
- `.agents/docs/testing/sh-agent/guide.md` — this repo's standing verification procedure
- `.agents/docs/architecture/DECISIONS.md` — ADR-018

## Next Steps

1. Operator: the Mac; the mirrors; the four queued 2026-09-04 proposals at the next `/retro --harness`.
2. First real task in a project under the new rules: watch that README, the deploy guide, and the changes entry's Docs line move with the code.
