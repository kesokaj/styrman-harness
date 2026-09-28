---
owner: main
date: 2026-09-28
---
# Verification — time-sortable names, Phase 2: behavioural runs   (2026-09-28)

Range: `e5a49d4` (Phase 1), the new definitions against the definitions at `e5a49d4~1`. Plan: `plans/2026-09-28-time-sortable-trail-names.md`.

## Status: PASS

Every run with the new text wrote timed names, and every must-fail control wrote date-only names.

## Runs

| Probe | New text | Must-fail control (text at `e5a49d4~1`) |
|---|---|---|
| `changelog` agent (Claude Code, `--agents` as `changelog-probe`; tools Read, Grep, Glob, Edit, Write plus Agent, no Bash; dispatch "Stamp: <`date +%F-%H%M`>…", no format named) | **3 of 3** timed: `2026-09-28-1431-add-healthz-endpoint.md`, `-1432-healthz-endpoint.md`, `-1433-healthz-endpoint.md` | 0 of 3 timed: `2026-09-28-healthz-endpoint.md` ×2, `2026-09-28-add-healthz-endpoint.md`. The same stamp in the dispatch did not produce a time, so the definition carries the behaviour. |
| `handoff` skill (copied as project skill `handoff-probe`, `name:` changed, invoked `/handoff-probe`) | **3 of 3** timed, for the handoff and for the change entry its Step 3 wrote (for example `2026-09-28-1435-healthz-endpoint.md`, `2026-09-28-1435-healthz-followup.md`) | 0 of 3 timed (`2026-09-28-healthz.md` and similar) |
| Context rule (`claude -p --append-system-prompt` with the rendered new context; "write a short handoff note…", no skill, no format named) | **3 of 3** timed: `2026-09-28-1442-demo-setup.md` | 0 of 3 timed, without the appended context: `2026-09-28-demo-setup.md` |
| Sandbox start, Claude Code (`HOME=.agents/scratch/tsn`, new harness installed, drift none) | 1 of 1: named `$(date +%F-%H%M)-<slug>.md`, "the timed name is newer", GANNET | — |
| Sandbox start, Antigravity (same home) | 1 of 1: the same three answers; the newest `cli-*.log` has 0 parse-error lines | — |
| E2E 2 must fire: `/init-project` on a clone of this repo at `e5a49d4` (sandbox home) | 1 of 1: "160 handoff and change files have date-only names: 55 in `handoffs/` and 105 in `changes/`…", with the warning text as specified. 0 renames in `git status`. The expected count from the plan's own command was 55 + 105. | — |
| E2E 2 must stay silent: `/init-project` on a project whose `changes/` and `handoffs/` hold only timed names | 1 of 1: no warning ("Both of those have timed names, so no old-style filenames need a warning") | — |

**Attribution.** In every changelog run the file was written by the subagent (`parent_tool_use_id` set on its `Write`); the main agent only ran `git status`. The tool order was read from the stream-json output of each run (`.agents/scratch/tsn-probe/runs/*.jsonl`).

## What went wrong on the way

The first changelog batch used a fixture whose `cmd/server/main.go` did not contain the `/healthz` endpoint the dispatch claimed. The agent correctly refused to write an entry. The batch was stopped, and the fixture now commits the endpoint, with a guard that aborts when the line is missing. Every run above is from the corrected fixture. This is BUG-087's class again, and a bug record follows at Done.

## Not exercised

- Antigravity for the changelog, handoff and context probes (those ran on Claude Code only); Antigravity got the sandbox start.
- The real home, which still has v1.6.0 installed.
- E2E 1, the agy cold resume, which is Phase 3.

The two `/init-project` runs could not commit, because the sandbox home has no git identity; that is expected for a scratch `HOME`, and it did not affect the warning.

## Environment

- Claude Code 2.1.283 and agy 1.2.8.
- The sandbox home `.agents/scratch/tsn`. Credential copies were `chmod 600` and deleted by name after the runs, and are copied again only for E2E 1.
- Probe scripts: `.agents/scratch/tsn-probe/run.sh`, with results in `results.tsv`.
