---
owner: bugs
date: 2026-09-13
id: BUG-008
status: fixed
area: trail/process
---
# BUG-008: five retro commits changed context/agents/skills with no trail

## Symptom
Five commits (830f1fb on 2026-09-04; 67616cb and 9f457a1 on 2026-09-08;
50ac771 on 2026-09-10; b64ffe8 on 2026-09-11) changed `context/`, `agents/`
and `skills/` with no `changes/` entry, no verification log, no sandbox start
of either CLI (required by this repo's Testing section after such a change),
and no handoff. `ROADMAP.md` "Now" still said "No tuning before two weeks of
use (ADR-009)" while three retros had already tuned. Found by the 2026-09-13
research pass (finding 4).

## Repro
`git log --stat d524c65..b64ffe8` vs
`ls .agents/docs/changes/ .agents/docs/testing/ .agents/docs/handoffs/`.

## Root cause
The retro skill's Step 5 ("apply the changes … run `sh-agent install` …
commit … push") described applying proposals as a plain commit, not as a
task under the Definition of Done. The sessions ran from another project
(styrman-crew), where the harness repo was touched as a side effect rather
than being the project under work, so nothing in the skill flagged that the
harness repo's own Testing rule (sandbox start of each CLI, changes entry)
still applied there.

## Fix
Catch-up entry
`.agents/docs/changes/2026-09-13-catch-up-2026-09-04-to-09-11.md`; the
research pass verified the resulting context on both CLIs from the sandbox;
plan decision D8 makes applying retro proposals a task under the Definition
of Done (changes entry, sandbox start, commit, push, handoff), landing in
`skills/retro/SKILL.md` in Phase 3 of
`.agents/docs/plans/2026-09-13-learning-loop-project-level.md`.

## Prevention
A change to another repo made from inside a project session is a task IN
THAT REPO — its own AGENTS.md rules (here: sandbox start of each CLI, changes
entry) and the Definition of Done apply there before the push; the retro
skill must say so explicitly.

## Related
BUG-007; ADR-009 (the two-week clock these commits ran ahead of).
