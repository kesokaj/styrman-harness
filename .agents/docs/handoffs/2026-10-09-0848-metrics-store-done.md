---
owner: main
date: 2026-10-09
---
# Handoff — 2026-10-09 — metrics store done; operator items queued

Session ID: 5dd0d90e-7c79-4c04-8f88-4cd1feee99e8
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

The operator asked to use the local SQLite for trajectory data and agent/sandbox metrics from every session, for suggested evaluation metrics, for a measure of how much the agent used the SQLite, for visual reports, for feedback into the learning loop, for sandbox tests, and for it to work on both Claude Code and Antigravity. All of it is built, reviewed, verified and pushed. The plan is closed.

## Completed

- Plan `.agents/docs/plans/2026-10-08-metrics-store.md` (D1–D39), ADR-055 and ADR-055a, and the REQUIREMENTS R3/R4/R7 edits plus the R3 addendum (approved).
- `internal/metrics`:
  - store, migrations, `export` without the salt;
  - Claude and Antigravity ingesters, which store shape only;
  - `Verb`, the shape normaliser;
  - events joined from git, the trail and feedback;
  - `show`, `report` (HTML with SVG) and `Numbers`.
- `internal/trail` holds the hardened git reads.
- The `sh-agent metrics ingest|show|report|export` verbs.
- Dream ingests first and carries a numbers-only metrics block.
- `/retro` gained Step 1 source 11.
- The A/B kit moved onto `metrics ingest --source`.
- Changes entries: `changes/2026-10-08-1100-…phase0` through `changes/2026-10-09-0847-metrics-store-done.md`.
- Verification: `testing/2026-10-09-verify-metrics-store.md` PASS (610 pass, 0 fail). Sandbox runs on both CLIs: `testing/2026-10-08-verify-metrics-sandbox.md`.
- BUG-164 fixed: git in project roots ran repository-defined programs; the dreamer had done so since v2.0.0.

## Open / In Progress

- Claude Code `retro` matrix-row check: 0 of 2. Both runs scoped the `-race` rule to the project under Step 2. Antigravity passes, 1 of 1. A control at 2c0f215 has not been run (plan `Open:` row).
- BUG-165 (open): `sh-agent search` runs an unhardened `git status`.
- Requirements edit awaiting approval, in R3's addendum:
  - replace "a summary of the metrics store (`sh-agent metrics show --json` for each project's window, numbers only)"
  - with "a numbers-only summary of the metrics store (its metric keys and values for each project's window and the window before, never a tool, skill, agent, label or path)".

  The code sends `Numbers()`, which is narrower than `show --json` (code review finding 6).
- The interactive map page (claude.ai artifact "Harness Delivery Map") has lagged the README map since 2026-09-17: it still shows bundles and signed tags, and has no zip place, binary place or dreamer. It needs a rebuild from the map (AGENTS.md rule 5).
- `govulncheck` is not installed on this box, so it was not run for this range.

## Key Decisions

- SQLite is the store of record. Rows hold shape only, plus three identifiers. Nothing goes through `report`; feedback travels through dream and `/retro` proposals (D1, D22, D27).
- Antigravity parity: tool names and command lines come from `steps.metadata`, never `step_payload` (D28). Missing metrics are null and noted, never 0 (D29).
- No `git status` in project roots. Modified is a raw-blob comparison, and partial clones are not read (BUG-164, D39 — the operator's "Fix both in this feature").

## Corrections

- None logged. The feedback log has no entry from this session.

## Agent notes (unverified)

- The operator's mid-session "this needs to work for both claude and agy" was taken as a requirement (D28, D29), not logged as a correction.
- The Claude `retro` 0-of-2 looks like Step 2 scope judgement on Sonnet, not this change. That is unproven.

## Environment State

- Branch: main, clean, pushed at 9187d89 (before this handoff). Newest tag v2.5.0.
- Running processes: none. Sandbox homes and clones are kept under `~/sandbox/metrics-sbx/` (runs `claude`, `claude-1`, `agy`) as evidence. The kit scripts are in `.agents/scratch/kits/` (gitignored).

## Artifacts to Read

- `.agents/docs/specs/sh-agent/metrics.md` — the surface, schema, metric definitions and failure modes.
- `.agents/docs/reviews/2026-10-08-security-2c0f215-166aa39.md` and `…-code-docs.md` — the findings and their Outcome.
- `.agents/docs/bugs/BUG-164-…`, `BUG-165-…` — the git hardening rule and the open search case.

## Next Steps

1. Operator: the first real ingest, soon, because transcripts expire. After `sh-agent update` or `install`, run `sh-agent metrics ingest`, then `sh-agent metrics report`.
2. Operator: approve or reject the R3 wording edit above (one line).
3. Agent, next task: BUG-165, which moves search's git calls onto `internal/trail.command`.
4. Before the next release:
   - the Claude `retro` control run at 2c0f215;
   - the release rule's Antigravity runs of the changed skill (`retro`) — the 2026-10-08 sandbox run is one, cite it or re-run it;
   - `govulncheck`.
5. Operator's call: rebuild the "Harness Delivery Map" artifact page from the README map.
