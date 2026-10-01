---
owner: bugs
date: 2026-10-01
id: BUG-135
status: fixed
area: dream verb and collection wiring
---
# BUG-135: dream verb inconsistencies, and a git option placed after the verb that silently emptied the modified-file set

## Symptom
(a) `dream --install-timer --remove-timer` exited 1 while every other dream usage error exits 2. (b) dream and memory refusals went to stderr, invisible to the writer tests see. (c) `--dry-run` of the timer flags was missing from the usage and synopsis lines. (d) `update` ran `Inspect` three times. (e) an agy stdout line over 16 MiB was logged as "invalid output". (f) `removeCopies` used the wall clock. (g) the first fix for the nightly `git status` put `--no-optional-locks` after the verb; git rejects that (exit 129), so the modified-file set was silently empty, and the project's `core.fsmonitor` hook ran at 03:00.

## Repro
`TestDreamVerb` (new cases), `TestAnswerLineTooLong`, `TestRemoveCopiesUsesInjectedClock`, `TestTrailStatusArgs` (pins the argument order and runs the command).

## Root cause
Small inconsistencies in `cmd/sh-agent/main.go` and `internal/dream/run.go`. For (g): a git global option placed after the verb, a caller that swallowed the exit status, and no test that ran the command. The fix for one finding caused the next round's defect.

## Fix
Exit 2 and the message on `out`; `[--dry-run]` in both usage lines; one `Inspect`; `answer` returns the scanner error ("output line too long or unreadable"); `removeCopies(conf, now)` (`internal/dream/run.go:322`); `trailStatusArgs` (`internal/dream/collect.go:242`) = `--no-optional-locks -c core.fsmonitor=false -C root status ...`. In the working tree, committed with Phase 1.

## Prevention
A command line built in code gets a test that runs it (an exec of the real binary against an empty directory), since git rejects a global option after the verb with an exit status the caller may swallow. Every usage refusal of one verb exits with the same code and writes to the caller's writer.

## Related
Audit IDs IN11-IN14; `reviews/2026-09-30-security-4764c5f-worktree-r2.md` finding 1, `reviews/2026-10-01-security-4764c5f-worktree-r3.md` finding 2; BUG-131 (the status call).
