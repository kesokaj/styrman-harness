---
owner: main
date: 2026-10-09
---
# Verification — real-home install at c478d86, then 5d8ab96, and the guide § 5 smoke on both CLIs (v2.6.0 gate)

Range reviewed: v2.5.0..5d8ab96

**The install.** The real home on the Linux box was already a checkout place, at 9c87ddb. It was installed twice from this checkout with `go run ./cmd/sh-agent install`:
- first at 56b4aaa;
- then at c478d86, after the `retro` template fix.

The authority is the project's Release rule (the real-home smoke on both CLIs before the tag, because `skills/retro` changed) and the operator's release instruction of 2026-10-09 ("try to fix all in the next release").

After the second install:
- `sh-agent status`: installed from commit `c478d86+dirty`, binary `c478d86d8106+dirty`, drift none (1789 files match).
- `type -a`: one `sh-agent`.

"+dirty" is the uncommitted trail records under `.agents/docs`, which install nothing.

## Status: PASS

6 of 6 rows passed, each against its own criterion, at c478d86 and again at 5d8ab96 (below). The checks are the three of `testing/2026-09-03-verify-real-home.md`, run by the scratch `realhome/smoke.sh`, in fresh clones under `~/sandbox/realhome`.

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context: GANNET, the Antigravity poll rule present, DoD step 2 quoted | 1 of 1, denied 0 | 12 |
| Claude Code | skill: Go review names `golang`, changes nothing | 1 of 1, denied 0 | 26 |
| Claude Code | agent: `verifier` over a comment edit, log written | 1 of 1, denied 0 | 64 |
| Antigravity | context ("GANNET", 36 harness skills, "Yes", DoD step 2) | 1 of 1, denied 0 | 145 |
| Antigravity | skill | 1 of 1, denied 0 | 150 |
| Antigravity | agent: `verifier` log `2026-10-09-verify-greet-comment.md` | 1 of 1, denied 0 | 174 |

**Claude Code's skill and agent rows failed twice before they passed.** Each time the output was "Failed to authenticate: OAuth session expired": the operator's Claude login had expired. That happened once at 56b4aaa and once at c478d86. After the login was renewed, both rows were re-run at c478d86 and passed. The failed attempts' outputs are kept under `~/sandbox/old/`.

**The Antigravity store check** reads a copy of conversation `6c06f0d3…`, deleted after:
- `<RULE[~/.gemini/config/rules/AGENTS.md]>` and `<RULE[…/critical-thinking.md]>` each appear as often as the demo project's own `AGENTS.md` block: 2 each, one per stored copy of the prompt. So there is one block of each per prompt.
- GANNET is present.
- There is no `<truncated` line.

## Run at the release head 5d8ab96

**The reinstall.** `go run ./cmd/sh-agent install` at 5d8ab96 wrote 3 files.
- `sh-agent status`: installed from commit 5d8ab96 (no "+dirty"), binary `5d8ab969bbf9`, drift none (1789 files match).

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context | 1 of 1, denied 0 | 12 |
| Claude Code | skill | 1 of 1, denied 0 | 36 |
| Claude Code | agent | 1 of 1, denied 0 | 75 |
| Antigravity | context ("GANNET", 38, "yes", DoD step 2) | 1 of 1, denied 0 | 99 |
| Antigravity | skill | 1 of 1, denied 0 | 161 |
| Antigravity | agent | 1 of 1, denied 0 | 242 |

**The store check** reads a copy of conversation `3eb1407a…`, deleted after. The global `rules/AGENTS.md`, `critical-thinking.md` and project `AGENTS.md` blocks appear 2 each, one per stored prompt copy. GANNET is present, with no `<truncated`.

## Not examined

- The Mac.
- The nightly dream timer (not installed on this box).
- An interactive session on either CLI (the operator's).
- A real-home `sh-agent metrics ingest`: the store already holds the operator's earlier ingest. This gate did not run one.
