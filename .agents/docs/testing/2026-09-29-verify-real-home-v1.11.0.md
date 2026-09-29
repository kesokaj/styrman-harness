---
owner: main
date: 2026-09-29
---
# Verification — real home on the Linux box after v1.11.0   (2026-09-29)

Range reviewed: b936ef8 (main HEAD, the checkout place installs the working tree; it is v1.11.0's tree plus the two handoff commits).

## Status: PASS for both CLIs

## Update

- The operator had already run `sh-agent update` at 09:01 (the manifest said `from commit b936ef8`). A second run at 09:26 reported `up to date with upstream (b936ef8)`, rebuilt from the checkout, `0 written, 258 unchanged, 0 removed`.
- `hash -r; type sh-agent` → `~/.local/bin/sh-agent`; `sh-agent version` → `b936ef8041c3`; `sh-agent status` → `binary b936ef8041c3`, `drift none — 257 files match the repo`.
- Both rendered contexts (`~/.claude/AGENTS.md`, `~/.gemini/config/rules/AGENTS.md`) end with the GANNET canary line.
- `update` kept `~/.config/sh-agent/feedback` (an older full clone with two commits its remote-tracking refs do not hold) and printed the by-hand check for it; unchanged since earlier sessions.

## Runs (real `$HOME`, cwd a fresh demo Go repo with one uncommitted comment edit)

The fixture was built under `.agents/scratch/smoke/demo` (since deleted): a `greet.go` with a test, `AGENTS.md` Setup commands, `testing/guide.md`, one commit, then a doc-comment edit left uncommitted. Each CLI got the prompt in print mode; Claude Code with `--allowedTools` naming the tools each check needs, Antigravity (`agy` 1.2.12, `-p`) with the real home's settings; no `.err` file holds an "auto-denied" line.

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools | PASS — `# AGENTS.md — global agent context`; 27 skills, 7 subagents; DoD step 3 = docs; `/grilling`; GANNET. 1 turn, no tool | PASS — the same five answers (`/grilling (or /grill-me)`); the newest `conversations/*.db` holds no `<truncated` line and 5 GANNET hits (BUG-085) |
| Knowledge skill auto-load (`golang` on `greet.go`, no edits) | PASS — loaded only `golang`, style findings, no edits, no agents | PASS — loaded `golang` and `review-standards`, style findings, no edits |
| Agent dispatch (`verifier` on the uncommitted edit) | PASS — `verifier` ran in the clone, `## Status: PASS`, log `.agents/docs/testing/2026-09-29-verify-greet.md`, the only file written | PASS — `verifier` ran in the clone, `## Status: PASS`, `Scope: complete`, same log path (it overwrote Claude Code's log of the same name; both read PASS before and after) |

## Not verified here

- The Mac and the devbox (operator).
- `explorer` on Antigravity (the 2026-09-03 procedure lists it as an extra; not run).
- The zip-place dry run with the real v1.11.0 zip on this box: step 3b of the release did the same check with the same binary and zip (`VERIFIED`), so it was not repeated.
