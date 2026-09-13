---
owner: changelog
date: 2026-09-13
---
# Done means empty TODO and ROADMAP     (2026-09-13)

## Summary
Operator clarified Definition of Done: at task completion, both `ROADMAP.md` and `TODO.md` must be empty (headers and comments only). Finished work is recorded in `changes/`; operator-only actions go to handoff Next Steps; scheduled work goes to plan files; neither file is a parking lot for ideas or unscheduled directions.

## Type
task

## Files changed
- `context/AGENTS.md` — Definition of Done step 4 updated (at Done both files are EMPTY, finished lines deleted, neither is a parking lot)
- `skills/handoff/SKILL.md` — Step 3 item 3 clarified (ROADMAP and TODO empty when task done, items go to handoff Next Steps)
- `skills/init-project/assets/TODO.md` — template header comment added (empty when work is done, ideas go to handoff)
- `skills/init-project/assets/ROADMAP.md` — template header comment added (both sections empty when done, operator actions/ideas go to handoff)
- `ROADMAP.md` — emptied to headers only (Now, Next, Later); last reviewed set to 2026-09-13
- `TODO.md` — emptied to headers only

## Agents dispatched
changelog (this entry)

## Verification
`go run ./cmd/sh-agent render` shows updated Definition of Done step 4 sentence. Sandbox e2e being run by main agent — see the handoff.

## Requirements
n/a
