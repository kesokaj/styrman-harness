---
owner: main
date: 2026-09-29
---
# Behavioural run — `/review --prompts` brief on the reviewer agent   (2026-09-29)

Range reviewed: working tree over 6b1b486 (the `--prompts` lens row of `skills/review/SKILL.md`, the vendored reference under `skills/review-standards/references/prompt-audit/`).

## Status: PASS — 3 of 3 runs on the first brief, 3 of 3 on the corrected brief; the control did not fail (see round 2)

## What was under test

The `--prompts` brief as written in the lens table of `skills/review/SKILL.md`, handed verbatim to the `reviewer` agent (its working-tree definition, `agents/reviewer.md`, unchanged by this task) through `claude -p --agents <json>`, per `AGENTS.md` § Testing. The dispatch text was assembled the way `/review` Step 3 assembles it: the paths scope sentence, the date, the slug `AGENTS.md-prompts`, the brief, the artifact sentence, a 30-call budget. One substitution: the brief's reference path `~/.claude/skills/review-standards/references/prompt-audit/prompt-audit.md` was replaced by the working-tree path of the same file, because the real home has no copy until the operator runs `sh-agent update`; the file was made readable with `--add-dir skills/review-standards`. The kit is `.agents/scratch/prompts-probe/` (`agents.json` built from `agents/reviewer.md`, `prompt.txt` built from the skill file's lens row, `run.sh`).

## Fixture (`.agents/scratch/prompts-probe/fixture`, asserted by `grep`, `test` and `diff` before any run)

A three-file Go service with these seeds in its instruction files:

| Seed | Where | Expected verdict | Since |
|---|---|---|---|
| `The entry point is cmd/nothing/main.go` (no such file) | `AGENTS.md:9` | Group 2 volatile specific, high, `rewrite` | commit 1fc87a5 |
| `go test ./...` in `AGENTS.md:8` against `make test` in `.agents/rules/tests.md:3` | both | Group 2 contradiction, high, `rewrite` or `flag` | 1fc87a5 |
| `CRITICAL: You MUST ALWAYS run gofmt before every edit.` | `AGENTS.md:13` | Group 1a pressure language; the Claude-only part `flag` under the two-model rule | 1fc87a5 |
| `Never commit .env files (BUG-001)` with a matching `bugs/INDEX.md` row | `AGENTS.md:14` | keep-list item 5, not in the diff | 1fc87a5 |
| `.agents/agents/reviewer-go.md` and `reviewer-ts.md`, same tools, near-duplicate prompts differing in the file suffix | both | Group 4 sub-agent roster check: one finding with a roster edit | e4e5c8a (round 2) |

## Round 1 — first brief, runs 1 to 3 (fresh clone each, real home for credentials, prompt on stdin, `--permission-mode acceptEdits`, `--allowedTools` = the reviewer's tools plus read-only Bash prefixes)

| Run | Report written | Dead path | Conflict | Pressure line | BUG-001 kept | `## Proposed diff` | Criteria `n/a — prompts lens` | Tree after run |
|---|---|---|---|---|---|---|---|---|
| 1 | yes, Write tool | found, high, rewrite | found, high, rewrite | found, medium; emphasis part flagged as Claude-specific, the timing defect rewritten | kept, cited keep-list 5 | yes, 4 hunks | 6 of 6 | only `.agents/docs/reviews/` untracked |
| 2 | yes, Write tool | found | found | found | kept | yes | 6 of 6 | only `.agents/docs/reviews/` |
| 3 | yes, Write tool | found | found | found | kept | yes | 6 of 6 | only `.agents/docs/reviews/` |

Tool calls, from the transcripts under `~/.claude/projects/-home-kesokaj-workspace-styrman-harness--agents-scratch-prompts-probe-runs-<n>-clone/`, executed and refused counted separately:

| Run | Calls | Read of `prompt-audit.md` | Executed, non-zero exit | Refused by the permission gate |
|---|---|---|---|---|
| 1 | 9 | 1 | 1 (`git blame AGENTS.md .agents/rules/tests.md`, the reviewer's own bad revision syntax; the rest of the command ran) | 0 |
| 2 | 6 | 1 | 0 | 1 (a `cat > … <<'EOF'` heredoc write of the report; the reviewer then used the Write tool) |
| 3 | 7 | 1 | 1 (a compound `cat … ; git log --a…`, exit 2) | 1 (the same heredoc write, then the Write tool) |

## Round 2 — corrected brief (the Group 4 clause now names the reference's roster-check exception; the roster seed added), runs 4 to 6, plus one control with the round-1 wording

| Run | Brief | Dead path | Conflict | Pressure line | BUG-001 kept | Roster pair found (`reviewer-go`, `reviewer-ts`) | `## Proposed diff` | Criteria n/a | Tree | Matrix check (`review-prompts` row, all seven greps) |
|---|---|---|---|---|---|---|---|---|---|---|
| 4 | corrected | found | found | found | kept | yes, high, one `reviewer` agent proposed | yes | 6 of 6 | only `reviews/` | PASS |
| 5 | corrected | found | found | found | kept | yes | yes | 6 of 6 | only `reviews/` | PASS |
| 6 | corrected | found | found | found | kept | yes | yes | 6 of 6 | only `reviews/` | PASS |
| ctl | round-1 wording (no roster clause) | found | found | found | kept | **yes** — "Group 4 = roster check only (1); no request-building code" | yes | 6 of 6 | only `reviews/` | PASS |

| Run | Calls | Read of `prompt-audit.md` | Executed, non-zero exit | Refused by the permission gate |
|---|---|---|---|---|
| 4 | 10 | 1 | 0 | 2 |
| 5 | 12 | 1 | 0 | 1 |
| 6 | 11 | 1 | 0 | 1 |
| ctl | 10 | 1 | 0 | 1 |

Reading of round 2: the corrected brief finds the roster pair in 3 of 3. The control also found it in 1 of 1, because the reviewer read the reference and applied its Group 4 exception even where the brief's paraphrase said Group 4 was not applicable. So the review finding's wording defect was real (the brief contradicted the file it points at) and is fixed, but its predicted effect, a silently disabled roster check, did not reproduce in the one control run. The refused calls are all the same shape as round 1: a Bash heredoc write of the report, redone with the Write tool. Cost per run about $0.47 to $0.55, 2 to 3 main-agent turns.

## Not exercised

- Antigravity: no run. The reviewer's agy flavour gets the same brief with the `~/.gemini/config/skills/…` path (`derive.go` rewrites nothing in skills; the brief spells both paths). The sandbox home this repo's procedure names no longer exists on this box; the real-home smoke checks after `sh-agent update` are the operator's, per `AGENTS.md` § Testing.
- The shipped reference path in the real home (`~/.claude/skills/review-standards/references/prompt-audit/prompt-audit.md`): substituted by the working-tree path, see above. The install tier (`testing/2026-09-29-verify-prompts-lens.md` and `-r2.md`) shows the file lands there on install.
- `/review --prompts` typed as a command through the `review` skill's Step 1 parsing: the runs dispatched the reviewer directly with the composed text, so the skill's parse of the `--prompts` token was not exercised; it is one token added to an existing flag list.
- A control without any brief was not run: without it the reviewer does a code review, and the seeds cannot be found by accident.
- `--fix` on the prompts lens (the main agent applying hunks): not run.
- More than one control run for the round-2 clause: one run, so "did not reproduce" is 0 of 1, not a rate.
