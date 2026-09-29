---
owner: bugs
date: 2026-09-29
id: BUG-106
status: fixed
area: testing (behavioural run kits)
---
# BUG-106: headless fixture clones placed inside the harness checkout wrote to its own trail

## Symptom
During Phase 3 of `.agents/docs/plans/2026-09-29-search-index.md`, 27 headless
behavioural runs (`claude -p --agents <json>`) used fixture clones at
`.agents/scratch/bkit/runs/<run>/notebox` — inside the harness checkout. All
three `changelog` agent runs wrote their change entry into the harness's own
trail, `.agents/docs/changes/2026-09-29-1848-store-add-close.md`, describing a
`Close()` method in `internal/search/store.go` that does not exist. Found by
the main agent reading the run results and `git status` showing the stray
untracked file.

## Repro
Run `claude -p --agents <json>` with a `changelog` dispatch inside a fixture
clone placed under `.agents/scratch/` of this checkout (any ancestor with
`.agents/docs`). The agent resolves the docs root upward/absolute and writes
into the harness's `.agents/docs/changes/`, not the clone's.

## Root cause
The main agent applied the scratch rule ("generated files go to
`.agents/scratch/`) to fixture clones as well as ordinary temp files, placing
a writable harness trail and the harness `AGENTS.md` above the fixture in the
directory tree. Claude Code loads ancestor `AGENTS.md` files, so every run saw
the harness project's rules too. The main agent had already noticed the
ancestor-`AGENTS.md` risk while the runs were in flight and let them continue
instead of relocating first.

## Fix
No code changed. Remaining and future behavioural/sandbox runs are relocated
outside the checkout (`~/sandbox/search/`, matching the 2026-09-16 simulation
and the Antigravity runs of the same phase). The stray untracked file was
deleted by name. Damage verified: no commit/tag/push occurred; re-applying the
main agent's edit script to a clean `git archive HEAD` and `cmp`-comparing
confirmed the 16 edited context/agent/skill files hold only the intended
edits. The three changelog-agent results are void; the other 24 runs' results
for search behaviour stand.

## Prevention
A fixture clone for a headless run must live outside every git checkout the
agent works in (never under `.agents/scratch/` of this repo), and the run
script asserts it before starting the CLI: `git -C "$(dirname <clone>)"
rev-parse` must fail. This is a specific instance of `AGENTS.md` § Testing's
existing sentence that a must-fail control "runs in a project that does not
contain the harness source" — that clause did not name fixture placement for
ordinary (must-pass) runs too, and the main agent did not extend it on its
own while the risk was visible.

## Related
BUG-092 (probe truncation wrote a fake change entry into this repo via a
different mechanism — shared symptom, different cause); BUG-087/090/091
(behavioural/sandbox run setup must be proven before its output counts).
Candidate for `AGENTS.md` § Testing: the must-fail sentence should be
broadened to cover every fixture clone, not only the must-fail control.
