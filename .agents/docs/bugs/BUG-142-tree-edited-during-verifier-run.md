---
owner: bugs
date: 2026-10-01
id: BUG-142
status: open
area: process/verifier-dispatch
---
# BUG-142: working tree edited while the verifier ran; its log ended BLOCKED

## Symptom
The Phase 9 verifier dispatch of plan `2026-09-30-audit-fix` (log `.agents/docs/testing/2026-10-01-verify-audit-fix-phase-9-r2.md`) ended `## Status: BLOCKED`. The verifier saw the tree change under it: `skills/distill/SKILL.md`, `skills/retro/SKILL.md` and `agents/changelog.md` were edited while its `go test -race ./...` and the render-size pins ran. A PASS would have belonged to no commit. A third run (`-r3`) had to be dispatched on a frozen tree.

## Repro
Dispatch `verifier` over the working tree, then edit any file the tests read before the log is on disk.

## Root cause
The main agent applied docs-lens review fixes to three skill and agent files while the verifier was in flight, treating "docs-only" edits as outside the freeze. The freeze rule (global context, Definition of Done step 1: "the tree stays frozen while it runs") covers every file: here the size pins and the agent-description tests read `skills/`, `agents/` and `context/`, not only Go source. The same slip happened once before in this plan (Phase 7's r1, recovered the same way). This is the second occurrence.

## Fix
None in code. `-r3` was dispatched after all edits, with the edit queue emptied first.

## Prevention
Before dispatching `verifier`, empty your own edit queue: make every fix still intended, read `git status --short`, then dispatch. While the dispatch is in flight do only reads (`Read`, `grep`, `git log`); a fix that cannot wait stops the run (`TaskStop`), is committed, and the verifier is re-dispatched. A "docs-only" edit is not exempt: in this repository the test suite reads `skills/`, `agents/` and `context/`.
The rule already exists (`context/AGENTS.md` § Definition of Done, step 1), so it was not followed.

## Related
- BUG-072 (earlier frozen-tree rule, graduated into the Definition of Done sentence). This record is the repeat.
- Phase 7 r1 of the same plan (same slip, same recovery).
- Graduation candidate: a line in the `review` skill's `--fix` step (Step 5): the verifier is dispatched after the last fix of the round, never alongside it. Actionable (the `git status` read before dispatch is a specific check) and generalizable (any repository whose tests read docs).
