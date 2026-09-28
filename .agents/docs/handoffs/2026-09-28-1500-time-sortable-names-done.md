---
owner: main
date: 2026-09-28
---
# Handoff — 2026-09-28 — time-sortable trail names done

Session ID: f538b90a-5aa7-44d1-8311-a160b5bb0167
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main`, starting from 9a17e0a. This handoff's commit closes the plan.

## Context

The operator asked for handoffs and change entries to be named so that name order is time order, for the rule to live in the global context for every project, and for `/init-project` to warn on an active project with old names. It is done. This handoff supersedes `2026-09-28-1450-time-sortable-names.md`, which had E2E 1 open.

## Completed

- **Plan:** `plans/2026-09-28-time-sortable-trail-names.md`, closed with every box ticked.
- **Change entries:** `changes/2026-09-28-1428-…`, `-1450-…` and `-1500-…`.
- **Rule:**
  - `handoffs/` and `changes/` files are named `$(date +%F-%H%M)-<slug>.md`;
  - a date-only name is older than every timed name of its day;
  - it is in the context and in the handoff skill; the changelog agent gets the stamp from its dispatcher;
  - `/init-project` warns in update mode.
- **Evidence:** behavioural 3 of 3 with every control 0 of 3; both CLIs in the sandbox; E2E 1 and 2 passed.
- **Bugs:** BUG-089, BUG-090, BUG-091.

## Open / In Progress

- Nothing is in progress, and `ROADMAP.md` "Now" and `TODO.md` are empty.
- Requirements edits awaiting approval: none.

## Key Decisions

- Plan D1–D5.

## Corrections

- None since `2026-09-28-okf-scope-fix.md`.

## Agent notes (unverified)

- BUG-087, BUG-090 and BUG-091 share one rule: prove the fixture or the sandbox before trusting a run. That is three occurrences. The `bugs` agent suggests it graduate into `AGENTS.md` § Testing or the sandbox procedure. It was not applied, because it is the operator's call.

## Environment State

- `main`, clean after this commit is pushed. The scratch sandbox home is `.agents/scratch/tsn`, with no credentials in it.

## Artifacts to Read

- The plan; `testing/2026-09-28-verify-time-sortable-names-phase-2.md` and `-e2e1.md`.

## Next Steps

1. **Operator:** release. The context changes for every machine, so the proposal is a MINOR, `/preflight v1.7.0`, from your own turn.
2. **Operator, after updating any machine:** the smoke checks, including the canary. The Mac is on v1.6.0 and the devbox has not been updated.
3. **Only if wanted:** graduate the prove-the-fixture rule (BUG-087/090/091) into `AGENTS.md` § Testing.
4. **`/retro`** when convenient: one correction today, and the retro is due by the 14-day rule after 2026-09-23.
