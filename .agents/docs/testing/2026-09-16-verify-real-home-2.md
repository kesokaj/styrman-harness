---
owner: main
date: 2026-09-16
---
# Verification — the harness reinstalled into this box's real home after the security session   (2026-09-16)

Range reviewed: the harness at dbde0aa, installed into `/home/kesokaj`; the smoke checks ran from that install.

## Status: PASS for both CLIs

Scope: the three smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI, after this session's five fixes (BUG-033 to BUG-038) changed `/review`, `/preflight`, four agent definitions, the context's Git paragraph and `internal/harness`. The runner was `.agents/scratch/realhome-smoke.sh`, as on 2026-09-15; outputs in `.agents/scratch/realhome-smoke/` and `.agents/scratch/realhome2/`.

## Install and status

| step | command | result |
|---|---|---|
| before | `~/.local/bin/sh-agent status` | exit 1: installed from a81ef4b, `binary a81ef4be3cb8 — older than the repo (dbde0aa)` |
| install at dbde0aa, clean checkout | `go run ./cmd/sh-agent install` | exit 0: `15 written, 245 unchanged, 0 removed, 0 kept — 260 files managed` |
| status | `~/.local/bin/sh-agent status` | exit 0: installed from dbde0aa, `binary dbde0aa7f769` with no warning, `drift none — 259 files match the repo` |

`~/.claude/skills/review/SKILL.md` and `~/.gemini/config/skills/preflight/SKILL.md` carry the BUG-033 check.

## Smoke checks

Each ran with the real `HOME` in a fresh clone of a sandbox demo. Antigravity's settings were widened for the checks and restored on exit; `cmp` against the backup reports the file identical.

**Check 1 — context, no tools.**
- Claude Code: PASS (6 s). `# AGENTS.md — global agent context`; 27 skills, 7 subagents; step 2's reviewer sentence verbatim.
- Antigravity: PASS (8 s). The same three, verbatim.

**Check 2 — skill auto-load, change nothing.**
- Claude Code: PASS (85 s). Reviewed `greet.go` against `golang` and `review-standards`, wrote nothing, and named the rule it applied; clone clean.
- Antigravity: PASS (66 s). Rule-by-rule table citing `go-review.md`; no change; clone clean.

**Check 3 — agent dispatch.**
- Claude Code: PASS (59 s). `verifier` in the clone: `PASS`, 13 tests, log `testing/2026-09-16-verify-greet-comment.md`.
- Antigravity: PASS (133 s). `explorer` wrote `plans/2026-09-16-package-map.md`; `verifier` reported `Status: PASS`, log `testing/2026-09-16-verify-uncommitted.md`.

`grep -c "Failed to parse"` prints `0` on the newest Antigravity log, `cli-20260916_071831.log`, which this run created.

## Not verified here

- The Mac.
- The full matrix from the real home; this session's targeted rows ran in the sandbox home (`testing/2026-09-16-verify-bug-033-034-review-gate.md`).
- Interactive prompts on either CLI.

## Environment

- Linux 6.12.107+deb13-amd64, Debian 13; Claude Code 2.1.272, Antigravity CLI 1.2.3, go1.24.4.
- Harness dbde0aa installed; the closing handoff commit changes trail files only.
