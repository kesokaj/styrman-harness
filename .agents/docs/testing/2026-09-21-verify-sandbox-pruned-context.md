---
owner: main
date: 2026-09-21
---
# Verification — sandbox start of both CLIs on the pruned context   (2026-09-21)

Range reviewed: working tree over f236e63 (Phase 1 of `plans/2026-09-21-context-consolidation.md`: pruned `context/AGENTS.md` and `context/dev-setup.md`, canary last, size test). Guide § 4 and the three smoke checks of § 5, run in a sandbox rather than the real home.

## Status: PASS

## Setup

The 2026-09-03 sandbox home no longer exists on this box. A new one: `.agents/scratch/sandbox-home` (gitignored) with the operator's own credential files copied from the real home (`~/.claude/.credentials.json`, `~/.claude.json`, `~/.gemini/antigravity-cli/{antigravity-oauth-token,settings.json}`; the Antigravity copy given the matrix's headless `permissions.allow` list), then `sh-agent install --home <that>`: 259 written, 1 unchanged, `status` drift none, rendered context 20,012 bytes ending `Canary: this context ends with the word GANNET.`. A first install attempt was cut off because its output was piped into `head -1` (SIGPIPE after the first line) — the installer's fault in nothing; never pipe it. Demo service: fresh clones of `.agents/scratch/realhome-smoke/demo` (4f1c122, a three-file Go HTTP service) into `.agents/scratch/sandbox-runs/{claude,agy}-demo`. Versions: Claude Code 2.1.278, Antigravity CLI 1.2.7.

## Runs (`HOME=<sandbox>`, headless)

| # | Check | Claude Code | Antigravity |
|---|---|---|---|
| 1 | Context questions, no tools: first heading; skill and subagent counts; the Bug bullet quoted; the "Silence from your own instrumentation" sentence and whether BUG-082 is mentioned; the canary word; which skill for a failing test / a new endpoint | PASS — `# AGENTS.md — global agent context`; 27 / 7; `- **Bug** (anything broken): load the `bug` skill.`; the new first sentence, BUG-082 absent; GANNET; bug / feature | PASS — same six answers; "the loaded rule does not report any truncation"; conversation store `89b56840-*.db`: `<RULE[…/sandbox-home/.gemini/config/rules/AGENTS.md]>` … `Canary: this context ends with the word GANNET.` … `</RULE[…]>`, no `<truncated` line |
| 2 | Knowledge skill auto-load (`golang` on `internal/greet/greet.go`, no edits) | PASS — loaded `golang`, five findings, no edits (`git status` clean) | PASS — loaded `golang` and `review-standards`, flagged the smoke comment as scaffolding, no edits |
| 3 | Agent dispatch (`verifier` on an uncommitted comment edit) | PASS — `Agent` tool use with `subagent_type: verifier` in the transcript; log `.agents/docs/testing/2026-09-21-verify-greet-comment.md`, `## Status: PASS`; the verifier also wrote `testing/guide.md` because the demo had none; no commit | PASS — `invoke_subagent` of `verifier` in the store; log `2026-09-21-verify-greet.md`, `## Status: PASS`; no commit |

Parse errors: none in the sandbox's Antigravity `cli-*.log` (`Failed to parse skill` count 0); none in the Claude transcripts.

## Caveat

Checks 2 and 3 for Antigravity ran in the `claude-demo` clone, not `agy-demo`: the two CLI runs were issued as parallel shell calls and the shell's working directory is shared, so the second call's `cd` was overtaken by the first's. Both CLIs therefore edited and verified the same clone (the smoke comment was appended twice; Antigravity's skill check flagged the first copy as scaffolding). The checks still prove what they claim — which skill loaded, which agent ran, where its log went — but a future run gives each CLI its own shell call or `cd`s inside a subshell.

## Not exercised

The real home (Phase 4); the full skill matrix; the behavioural runs of `agents/bugs.md` (Phase 2); govulncheck.
