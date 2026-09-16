---
owner: main
date: 2026-09-16
---
# Verification — Antigravity's real-home smoke checks, completed   (2026-09-16)

Range reviewed: the installed tree at `8cf30e0` (the v1.0.0 release plus its handoff)

## Status: PASS — all three checks, both CLIs, from the real `$HOME`

This closes the one item every earlier log of the day left open. Claude Code passed all three checks earlier; Antigravity had passed the context check while its other two were blocked on a tool-permission grant this session's classifier refused three times. The operator approved the permission, the flag took effect (`Print mode: --dangerously-skip-permissions set, auto-approving all tool permissions` in `cli-20260916_213200.log`), and both checks then ran to completion. The earlier six-minute timeouts were this session's own `--print-timeout`, not a block: check 2 needed about four minutes and check 3 about nine.

## Runs

| Check | Result |
|---|---|
| Context questions, no tools | PASS (recorded earlier today) — the two new context rules quoted verbatim, 27 skills and 7 subagents, zero `Failed to parse` lines |
| Knowledge skill auto-load (`golang` on `internal/greet/greet.go`, no edits) | **PASS**, 21:38:28→21:42:35, exit 0. Named `golang` and reached it through `review-standards`, the same route as 2026-09-03. Three real findings: `fmt.Sprintf` for a single concatenation, two exported functions with identical signatures returning zero, and both uncovered and uncalled. It wrote nothing — `git status` on the demo showed only the edits this session had staged for it |
| Agent dispatch (`verifier` on an uncommitted change) | **PASS**, 21:42:35→21:51:57, exit 0 |

## What the dispatch check proved beyond itself

The Antigravity run printed `Waiting for the subagent to complete...` and did not end its turn until the verifier returned — the rule this session added for BUG-031, observed working on the CLI whose own failure mode was the opposite one.

The dispatched `verifier` then exercised three more of today's changes without being asked to:

- **`<app>` resolved, not left literal.** It wrote `.agents/docs/testing/server/guide.md`, headed `# Verification guide — server`. The demo's binary is `cmd/server` and its repository directory is `demo`; the G20 fix says `<app>` is the single `cmd/<app>` name. It took `server`. That is the exact gap G20 was filed for, closed and now demonstrated on the CLI, in a project it had never seen.
- **Scratch in the right place.** Its working files went to `.agents/scratch/verify-greet-count/`, which is BUG-014's rule.
- **The right verdict.** `## Status: FAIL`, root cause "adds `Count()` without a covering test", `Scope: complete`, and a `Range reviewed:` line naming the working tree against `3ff76cb`. A coverage gap reported as FAIL and never patched is what the agent definition demands; it changed no source.

## Not examined

The full skill matrix (`testing/skill-matrix.md`) still has not run from a real home on either CLI. The Mac. Nothing here touched the harness repository itself: the checks ran in the throwaway service under `.agents/scratch/realhome-smoke/demo`, which is disposable.

## Environment

Linux (Debian 13, amd64), real `$HOME`, harness installed from `8cf30e0`, `drift none`. Antigravity's `~/.gemini/antigravity-cli/settings.json` was NOT modified for this run — the per-invocation `--dangerously-skip-permissions` flag was used instead, so the operator's allow list still holds its six narrow entries.
