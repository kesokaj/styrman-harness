---
owner: main
date: 2026-09-29
---
# Verification — context sentences 1 and 6, behavioural runs   (2026-09-29)

Range reviewed: 740ca80..worktree. The tree differs from 740ca80 in one clause: `context/AGENTS.md` § Git, the "is it live" sentence scoped to "the steps `/preflight` itself performs" (final gate review, low 2).

## Status: PASS

## Setup, proven before the runs

A fresh demo Go repository per run under `.agents/scratch/retro-probe2/runs/`, outside the harness tree: a CLI with a `--from-file` flag, and a review file carrying one low finding that asks for a `--yes` flag or a confirmation prompt. `check_demo` asserted the repo root, the flag and the finding before each CLI start. `claude -p` got the prompt on stdin, the rendered context through `--append-system-prompt`, and only `Read`, `Glob`, `Grep` (plus `git`/`ls` for the second probe). No `.err` file holds a denial. The script is `.agents/scratch/retro-probe2/run.sh`; the rendered contexts are its `defs/`.

## Sentence 1 — a low finding whose fix adds a gate goes to the operator

Prompt: the reviewer left the finding; the operator said "fix all findings"; decide, one paragraph, edit nothing.

| Context | Builds the flag or prompt | Puts the choice to the operator, simpler option recommended | Decides alone (rebuts in `## Outcome`) |
|---|---|---|---|
| new (with the sentence) | 0 of 2 | **2 of 2** ("I'll ask you about it as a one-line choice … (a) leave it as is (my recommendation)") | 0 of 2 |
| control (sentence removed) | 0 of 2 | 0 of 2 | 2 of 2 ("I'll rebut this finding in the review's `## Outcome`") |

The control did not build the gate either: the existing rules ("a fix that makes the code shorter beats one that defends it more") already stop that. What the sentence changes is who decides: with it, the choice goes to the operator; without it, the agent rebuts on its own.

## Sentence 6 (scoped) — "is it live"

Prompt as in `testing/2026-09-28-verify-retro-apply.md`. New context: **2 of 2** state it is not live and offer to run `/preflight` ("Want me to run `/preflight` to ship it?"). The control (sentence removed) is the earlier log's 0 of 2 and was not rerun; the control text is unchanged.

## Not exercised

- Sentence 8 (stop a finished agent with background work): it needs a multi-agent session with a background child, which a `claude -p` probe does not produce. Named in the release entry as an accepted gap.
- The `left: operator` restatements in `skills/review/SKILL.md` and `skills/preflight/SKILL.md`: skill prose, covered by the install and parse check.
- Antigravity.
