---
owner: main
date: 2026-09-29
---
# Verification — real home on the Linux box with `sh-agent search`, both CLIs   (2026-09-29)

Range reviewed: f41e5c0 installed in the real home by `sh-agent update` from this checkout, on the operator's word ("Don't forget to update local harness"). Plan `plans/2026-09-29-search-release-readiness.md`, Phase 4.

## Status: PASS on both CLIs

## Install

| Check | Result |
|---|---|
| `sh-agent update` (second run, after the context change) | `5 written, 261 unchanged, 0 removed, 0 kept — 266 files managed`; the first run at 2716d08 wrote 33 |
| `hash -r`, `type sh-agent` | `/home/kesokaj/.local/bin/sh-agent` |
| `sh-agent status` | `repo … (main @ f41e5c0, clean)`, `installed … from commit f41e5c0`, `binary f41e5c0d1627`, no warning |
| Both rendered contexts | end with `GANNET.`; each holds "Search first" and `sh-agent search --open` (what is open) once |
| `sh-agent search --refs ADR-040` in this checkout, installed binary | `ADR-040 … Superseded by ADR-041`, `referenced by 25:` |

## Smoke checks

Fresh clones of a fixture (the `tasklog` sample project: Go CLI, tests passing, a trail with five bug records, plans, reviews and handoffs) under `~/sandbox/search/smoke/runs/`, outside every checkout (BUG-106). Real `$HOME`, installed skills and agents, the installed `sh-agent`.

| Check | Claude Code 2.1.284 | Antigravity 1.2.13 |
|---|---|---|
| Context questions, no tools | PASS — "AGENTS.md — global agent context", 28 skills and 8 subagents, GANNET, `sh-agent search <words>` before Grep, `sh-agent search --open` last in the resume order | PASS — the same four answers; the conversation store holds the `<RULE[…AGENTS.md]>` block with the canary line and no `<truncated` line |
| Skill load (review `store.go`, no edits) | PASS — `Skill: golang` in the transcript; no edit, no command | PASS — `golang` and `review-standards` named; tree clean |
| Agent dispatch (`verifier` on an uncommitted comment) | PASS — returned; `testing/2026-09-29-2216-smoke.md`, `## Status: PASS`, 72 tests | PASS — returned; `testing/2026-09-29-verify-smoke.md`, `## Status: PASS` |
| Trail question ("what does the trail record about `store.go`? Use the harness's way") | PASS — ran `sh-agent search --refs internal/tasklog/store.go`, 9 records, answered with `file:line` | PASS — search output in the conversation store; answered with `file:line` |
| Logs | — | no "auto-denied" line in any run's stderr; no "invalid allow entry" in `cli.log` |

The Antigravity runs used the real settings (`toolPermission: always-proceed`) in a directory outside its trusted workspaces; every run completed.

## After the runs

The two clones that ran `sh-agent search` had registered in the real global search index. They were deleted, and `sh-agent index --global` dropped their registrations; the index then held this checkout alone.

## Not verified here

The Mac and the devbox (after the release, by the operator's `sh-agent update` there).
