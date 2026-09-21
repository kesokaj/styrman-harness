---
owner: bugs
date: 2026-09-21
id: BUG-086
status: fixed
area: process / shell usage
---
# BUG-086: installer piped into `head` was killed by SIGPIPE, leaving files with no manifest

## Symptom
`go run ./cmd/sh-agent install --home .agents/scratch/sandbox-home 2>&1 | head -1` printed only "building sh-agent from …", then `sh-agent status --home .agents/scratch/sandbox-home` reported `installed nothing`; only `.claude/AGENTS.md` had been written, no manifest.

## Repro
Pipe any writer with side effects into a filter that exits early: `head -1` closes the pipe after one line and the producer receives SIGPIPE at its next write, while the pipeline's exit status is the filter's (0), hiding the kill.

## Root cause
Operator-side command shape, not the installer: `apply.go` writes the manifest last, so an interrupted install leaves files on disk but no manifest — `status` correctly reads that as nothing installed, which is what made the failure visible.

## Fix
Re-ran the install bare (no pipe); 259 files written, drift none.

## Prevention
Never pipe a command that has side effects (install, update, report, a git write) into a filter that can exit before consuming its output (`head`, `grep -q`, `grep -m1`, `sed q`); run it bare, redirect its output to a file under `.agents/scratch/` if it must be shortened, and read that file. Recommend widening the existing `context/AGENTS.md` § Hard rules sentence "A verification command is never piped into a filter inside a `&&` chain that leads to a commit…" with one clause covering a side-effecting command piped into an early-exiting filter, rather than a new bullet — the rendered context is size-pinned at 22,000 bytes (ADR-035, `TestRenderedContextFitsAntigravityBlock`).

## Related
BUG-011 (unread exit status in an `add … ; commit` chain — same family, different mechanism); BUG-011's index row still sits in the main table although its rule already lives in `context/AGENTS.md`. Found in plan `.agents/docs/plans/2026-09-21-context-consolidation.md` Phase 1, `testing/2026-09-21-verify-sandbox-pruned-context.md` § Setup.
