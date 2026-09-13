---
owner: main
date: 2026-09-13
---
# Handoff — 2026-09-13 — Linux real home updated to a35b567 and smoke-checked; the Mac remains

Session ID: session_01CtMYM9AGp33edxBQGURNaB
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean, pushed (this handoff's commit is the last of the day)

## Context

Closes the first operator item of `2026-09-13-improvement-items-applied.md` (which this file supersedes): on the operator's instruction the agent ran `sh-agent update` in the real home on this box and the three smoke checks from it on both CLIs.

## Completed

- `sh-agent update` rebuilt and installed the binary (`a35b567582fc`, no warning, no drift); the 2026-09-03 binary is gone from this box. Log: `.agents/docs/testing/2026-09-13-verify-real-home.md`; entry: `.agents/docs/changes/2026-09-13-real-home-updated.md`.
- Smoke checks from the real home, both CLIs, all PASS: context (27 skills, step 4 quoted), `golang` auto-load, `verifier` dispatch writing its log in the clone. No settings were edited.

## Open / In Progress

Nothing queued in this repo. Operator-only:

- The Mac: `sh-agent update`, then the same three checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`); Claude Code there has never been checked.
- Optional `REQUIREMENTS.md` Known-limits sentence about the Mac (wording in `2026-09-13-improvement-items-applied.md`).
- Upgrade Go to ≥ 1.24.6 on both machines (govulncheck advisories, toolchain-fixed).
- Refresh the downstream mirrors in styrman-crew and argolis-devbox from those projects (`.agents/docs/deploy/remote.md`, "Downstream mirrors").
- Item 20 of the improvement review (per-machine queue and index): decide after the first `/retro` on the Mac. Full skill matrix when convenient.

## Key Decisions

- The real Antigravity settings (`toolPermission: always-proceed`, `/home/kesokaj` trusted) made the temporary allow-list of the 2026-09-03 procedure unnecessary; the procedure file is unchanged, this log says which path was used.

## Corrections

- None in this part of the session.

## Environment State

- Real home on this box: harness a35b567, both CLIs verified from it. Sandbox home at the same tree. Clones under `/home/kesokaj/sandbox/runs/realhome-*` hold the smoke-check logs; running processes: none.

## Artifacts to Read

- `.agents/docs/testing/2026-09-13-verify-real-home.md` — the update and the six runs
- `.agents/docs/handoffs/2026-09-13-improvement-items-applied.md` — everything else from today (superseded only on the Linux update item)

## Next Steps

1. Operator, on the Mac: `sh-agent update`; `sh-agent status` shows a stamped `binary` line; the three checks on each CLI.
2. Refresh the two downstream mirrors from their own projects.
3. First real `/retro` in a project; then `/retro --harness` here.
