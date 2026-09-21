---
owner: main
date: 2026-09-21
---
# Verification — real-home smoke checks after the context consolidation, both CLIs   (2026-09-21)

Range reviewed: harness `94a6596` installed in the real `$HOME` by `sh-agent update` on the operator's word ("Run it on this machine"; 17 written, 243 unchanged; `hash -r`, `sh-agent version` → `sh-agent 94a65966d102`, `type sh-agent` → hashed `~/.local/bin/sh-agent`; `status`: binary `94a65966d102`, drift none; rendered context 20,205 bytes, identical on both vendors, last line `Canary: this context ends with the word GANNET.`). Phase 4 of `plans/2026-09-21-context-consolidation.md`; the e2e golden path.

## Status: PASS

Procedure as in `2026-09-17-verify-real-home-smoke-4.md`: fresh demo clones `.agents/scratch/realhome-smoke/smoke5-c` (Claude Code) and `smoke5-a` (Antigravity), each CLI in its own subshell; one uncommitted doc-comment edit for check 3; no commit in either clone.

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools: first heading, skill count and pruned names, subagents, DoD step 3, the grilling command, the canary word | PASS — `# AGENTS.md — global agent context`, 27 skills, `grill-me`/`skill-repair` absent, 7 agents, step 3 quoted with the deploy guides, `/grilling`, GANNET | PASS — same six; "no truncation is reported"; conversation store: one `<RULE[~/.gemini/config/rules/AGENTS.md]>` block, the canary line inside it, **zero `<truncated` markers** — the first real-home Antigravity session without a cut since v1.2.0 (BUG-085) |
| Knowledge skill auto-load (`golang` on `greet.go`, no edits) | PASS — `golang` loaded, two findings, no edits | PASS — `golang` and `review-standards` loaded, no edits |
| Agent dispatch (`verifier` on the uncommitted edit) | PASS — log `.agents/docs/testing/2026-09-21-greet-comment-verify.md`, `Status: PASS`; the verifier also wrote the missing `guide.md` | PASS — `invoke_subagent` of `verifier` in the store; log `2026-09-21-verify-greet-smoke.md`, `## Status: PASS` |

Antigravity checks 2 and 3 first produced no output: headless mode auto-denied `read_file` and `command`, because the real `~/.gemini/antigravity-cli/settings.json` today holds only `trustedWorkspaces` (no `permissions.allow`; the 2026-09-16 log records six narrow entries then). They were rerun with the per-invocation `--dangerously-skip-permissions` flag inside the throwaway clone, as the 2026-09-16 run did; the operator's settings file was not modified. Parse errors: `Failed to parse skill` count 0 in the newest `cli-*.log`.

## Not exercised

The Mac (Antigravity-only place): the operator runs `sh-agent update` there and asks for the canary word. The full skill matrix. The devbox.
