---
owner: main
date: 2026-10-01
---
# Verification — real-home install of the v2.0.0 gated head and the guide § 5 smoke on both CLIs (2026-10-01)

Range reviewed: v1.16.2..9c6c3ed — two installs into the real home: 0b1b5b8 (the first gated head) and 9c6c3ed (the final one, after the gate fixes), each with the six smoke checks; installed on the Linux box (`go run ./cmd/sh-agent install` per `deploy/local.md` § Build and install, the operator's ship request being the deploy approval)

## Status: PASS — 6 of 6 smoke checks on both CLIs at 0b1b5b8, and 6 of 6 again at 9c6c3ed

Install: `3 written, 263 unchanged, 0 removed, 0 kept — 266 files managed`; `sh-agent status`: `repo has 34 skills, 8 agents, 0 commands, 0 rules`, `drift none — 265 files match the repo`; `sh-agent version` → `sh-agent 0b1b5b8453da` (after `hash -r`, the bare name).

Smoke (`.agents/docs/testing/2026-09-03-verify-real-home.md` procedure, run by `.agents/scratch/audit-fix/kit-realhome.sh <cli> <check>` from a fresh clone of the demo service outside every checkout; the CLI runs with the REAL home on purpose, guide § 5):

| CLI | check | result | seconds |
|---|---|---|---|
| claude | context | PASS (denied 0) | 6 |
| claude | agent | PASS (denied 0) | 44 |
| agy | context | PASS (denied 0) | 11 |
| agy | skill | PASS (denied 0) | 57 |
| claude | skill | PASS (denied 0) | 119 |
| agy | agent | PASS (denied 0) | 84 |

Checks: `context` — the canary word GANNET and the skill count 34 answered from the global context with no tool; `skill` — the golang skill applied to `internal/greet/greet.go` with no edit; `agent` — the verifier agent dispatched over a one-comment change, its log under `.agents/docs/testing/` with a `## Status` line. The Antigravity runs used the real `~/.gemini/config` (rules, skills, derived agents) with no `auto-denied` line.

## The final head, 9c6c3ed

Install: `15 written, 251 unchanged, 0 removed, 0 kept — 266 files managed`; `drift none — 265 files match the repo`; `sh-agent version` → `sh-agent 9c6c3ed2e23f+dirty` (the `+dirty` is the gate's untracked records under `.agents/docs/`, which the release commit takes in; the binary is rebuilt by the deploy step after it). Smoke, same procedure:

| CLI | check | result | seconds |
|---|---|---|---|
| `context` | claude | PASS (denied 0) | 7 |
| `agent` | claude | PASS (denied 0) | 43 |
| `skill` | claude | PASS (denied 0) | 55 |
| `context` | agy | PASS (denied 0) | 12 |
| `skill` | agy | PASS (denied 0) | 63 |
| `agent` | agy | PASS (denied 0) | 97 |

## What was not exercised

The Mac (no access this session); `sh-agent update` at another place (the standing operator step of `deploy/local.md` § Operator steps, not pending work).
