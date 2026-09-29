---
owner: bugs
date: 2026-09-29
id: BUG-107
status: fixed
area: testing (behavioural run kits)
---
# BUG-107: `sh-agent search` inside headless fixture clones wrote the real home's global search index

## Symptom
During Phase 3 of `.agents/docs/plans/2026-09-29-search-index.md`, headless Claude Code
behavioural runs (which must run from the operator's real home, since a scratch `HOME`
cannot authenticate) called `sh-agent search` inside fixture clones. `sh-agent search`
registers every checkout it runs in into the global index at `~/.config/sh-agent/search.db`
of the home it runs under — here the operator's real home. After the runs, the real
`~/.config/sh-agent/` held `search.db`, `search.db-wal`, `search.db-shm`, registering 17
test clones (`~/sandbox/search/cruns/`, `.agents/scratch/bkit/runs/` in this checkout) plus
this checkout itself. Found by the main agent while designing the Phase 5 kit, listing
`~/.config/sh-agent/` and reading the `projects` table read-only.

## Repro
Run a headless behavioural round of an agent whose dispatch or tools can call
`sh-agent search` (or any verb that registers a checkout), from the real home, against
fixture clones on disk; check `~/.config/sh-agent/search.db` afterward.

## Root cause
A verb with a side effect on harness state (registration in the global search index) is
invoked by the agent under test, not by the kit itself, so the kit's `--home` discipline —
applied only to the kit's own `sh-agent` calls — never reaches it. `sh-agent` has no
environment variable for the home, so there is no way to redirect a call the kit does not
construct.

## Fix
The three cache files (`search.db`, `search.db-wal`, `search.db-shm`) were deleted by name
after being read; `feedback.jsonl`, `retros.jsonl`, `proposals/`, `state.json` were listed
before and after and left untouched. The Phase 5 kit now puts a wrapper named `sh-agent`
first on PATH that execs the real build with `--home <scratch home>` appended, so any call
the agent under test makes lands in the scratch home, not the real one.

## Prevention
A behavioural run in which the agent under test may itself invoke `sh-agent` puts a wrapper
first on PATH that pins `--home` to a scratch home (not only the kit's own calls); the run
script then asserts afterward that the real `~/.config/sh-agent/` gained no file.

## Related
BUG-106 (same run batch, fixture clones placed inside the checkout). Three testing-kit
isolation bugs in one batch (BUG-106, BUG-107) plus the general behavioural-run-setup rule
already graduated for BUG-087/090/091 (`AGENTS.md` § Testing, "A behavioural or sandbox run
proves its setup before its output counts"); this PATH-wrapper rule extends that paragraph
and belongs in the same home when it graduates.
