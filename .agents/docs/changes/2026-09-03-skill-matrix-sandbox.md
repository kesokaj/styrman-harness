---
owner: main
date: 2026-09-03
---
# Every skill and agent exercised on both CLIs     (date: 2026-09-03)

## Summary

Operator asked whether the sandbox had covered all skills and the deploy guides: it had not. A mock scenario per skill (31 scenarios, the six workflow agents through the skills that dispatch them, `explorer` directly, and a context question on the deploy guides) now runs headless on Claude Code and Antigravity in a fresh clone each, with a shell check per scenario and the CLI's own transcript read for skills loaded and agents dispatched. 62 of 62 pass; two Claude runs first failed only because headless mode cannot reach `~/.config/sh-agent/` and passed with the directory allowed. No harness defect surfaced in this pass; one ambiguity recorded (`/preflight` base ref when the newest review says `..worktree`).

The kit lives in `.agents/scratch/skilltest/` (gitignored, reproducible: `matrix.txt`, `prompts/`, `checks/`, `run-one.sh`, `summarize.py`). Lesson recorded in the log: never edit the driver script while runs are in flight — bash reads it incrementally and seven result lines were lost (outputs intact, checks re-run by hand).

## Type

task

## Files changed

- `.agents/docs/testing/2026-09-03-verify-sandbox-skills.md` — new
- `.agents/docs/plans/2026-09-03-rebuild.md` — section 12 ticked

## Agents dispatched

None in this repo; inside the demo clones the scenarios dispatched `verifier`, `reviewer`, `security`, `documentation`, `bugs`, `changelog`, `explorer` as their skills prescribe.

## Verification

See the log. Executed-command audit over every transcript: no cloud, container, or registry command ran.

## Requirements

R2 (both vendors read and act on every file).
