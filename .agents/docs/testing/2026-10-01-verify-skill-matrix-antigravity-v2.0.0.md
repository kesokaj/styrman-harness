---
owner: main
date: 2026-10-01
---
# Verification — the Antigravity rows the release rule still owed before v2.0.0 (2026-10-01)

Range reviewed: v1.16.2..9c6c3ed (the rows below ran from the sandbox home `~/sandbox/audit-fix/sbAgy` installed from 0b1b5b8, `sh-agent install --home … --vendor agy`; the four skills changed by the gate fixes were re-run from a re-install at 9c6c3ed — second table)

## Status: PASS — 14 of 14 rows, plus 4 of 4 re-run rows (bash, docker, golang, gcloud-auth) at 9c6c3ed

Why: `AGENTS.md` § Release makes a gate NO-SHIP until one run of each changed skill and agent exists on Antigravity. The Phase 9 matrix (`testing/2026-10-01-verify-skill-matrix-antigravity.md`, 24 of 24) left these without an Antigravity run in the range: skills `bug`, `review` (code lens), `feature`, `research`, `retro`, `feedback-log`, `preflight`, `dream`, `project-memory`; agents `bugs`, `changelog`, `documentation`, `dreamer`, `security` (`reviewer`, `explorer` and `verifier` ran in the Phase 9 rows and the real-home smoke).

Procedure: `.agents/scratch/audit-fix/kit-agy-ext.sh <row>` (not committed, ADR-010): `testing/kit-preamble.sh` sourced (the real `~/.config/sh-agent` fingerprinted before and after; 0 `kit: FAIL` lines in the batch log), a fresh clone of the demo service per row outside every checkout, `settings.json` written per run with `command(*)`, `read_file(*)`, `write_file(*)` and the clone as the only trusted workspace (BUG-141), `agy -p "The project root is the current directory, <clone>; ignore any other checkout. <prompt>" --mode accept-edits --print-timeout 24m`, the row's check cell from `testing/skill-matrix.md` (fixtures: the `TrimLeft` regression committed for `bug`; one uncommitted edit for `review`; a seeded three-file store for `dream` and the dreamer; two seeded feedback entries in the sandbox config for `retro`; two sessions for `project-memory`), the token copied in for the batch and deleted by name after it (`ls` of the path: no such file); `xargs -P 2`. A first attempt of this batch refused to start because the real-home install of 0b1b5b8 ran at the same time and changed the fingerprinted directory; the batch was re-run after the smoke finished, with nothing else touching that directory.

| row | kind | result | seconds |
|---|---|---|---|
| `feedback` | skill | PASS (denied 0) | 120 |
| `review` | skill | PASS (denied 0) | 162 |
| `bug` | skill | PASS (denied 0) | 570 |
| `research` | skill | PASS (denied 0) | 81 |
| `dream` | skill | PASS (denied 0) | 172 |
| `project-memory` | skill | PASS (denied 0) | 67 |
| `feature` | skill | PASS (denied 0) | 806 |
| `preflight` | skill | PASS (denied 0) | 308 |
| `ag-bugs` | agent dispatch | PASS (denied 0) | 133 |
| `retro` | skill | PASS (denied 0) | 320 |
| `ag-changelog` | agent dispatch | PASS (denied 0) | 79 |
| `ag-dreamer` | agent dispatch | PASS (denied 0) | 100 |
| `ag-security` | agent dispatch | PASS (denied 0) | 104 |
| `ag-documentation` | agent dispatch | PASS (denied 0) | 209 |

Each PASS is the row's check cell: `bug` — the fix compiles and the trailing-space test passes, a `BUG-0*` record exists, the output names the id; `review` — a `review-worktree-code` report; `feature` — a plan file with unticked boxes, `expect:` clauses, Checkpoint and Cold-reader boxes, `## Decisions`, and no code; `research` — ldflags discussed and the "continue" question asked; `retro` — a proposal file under the sandbox's `~/.config/sh-agent/proposals/`; `feedback` — the correction logged to the sandbox's `feedback.jsonl`; `preflight` — a verdict and a preflight report; `dream` — the store named and the broken and merge verdicts; `project-memory` — a verb-written memory file (no native `originSessionId`), a `MEMORY.md` slug line, and the second session answering Friday; agent rows — the artifact the agent owns exists (a BUG record and a new INDEX row; a `changes/` entry; README changed with `ROADMAP.md` untouched; a `retros/*dreamer*` artifact; a `reviews/*security*` findings file).

## Re-run of the four skills the gate fixes changed (9c6c3ed)

The gate fixes changed `bash`, `docker`, `golang` and `gcloud-setup` (89f3308 → 9c6c3ed), so their matrix rows ran again from a sandbox re-installed at 9c6c3ed (`kit-agy-matrix.sh`, same grants, token copied in and deleted by name after; `invalid_allow=6` is `cli.log`'s cumulative count from Phase 9's round 4, none added):

| row | result | seconds |
|---|---|---|
| `bash` | PASS (denied 0) | 127 |
| `golang` | PASS (denied 0) | 65 |
| `gcloud-auth` | PASS (denied 0) | 19 |
| `docker` | PASS (denied 0) | 283 |

## What was not exercised

The Mac; `init-memory*` rows (unchanged since their 2026-09-29/30 runs); `review --prompts`, `review-project`, the shadow rows (unchanged text since their last runs).
