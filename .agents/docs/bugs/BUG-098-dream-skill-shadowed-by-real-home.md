---
owner: bugs
date: 2026-09-29
id: BUG-098
status: fixed
area: skills/dream (behavioural probes)
---
# BUG-098: nine `/dream` probe runs measured the installed skill, not the working-tree edit

## Symptom
While a `/dream` edit was under probe through a project symlink, the operator installed the real home from commit 4426e58 at 14:03. Nine subsequent headless `/dream` runs loaded `~/.claude/skills/dream/SKILL.md` instead of the working-tree file; their first lines kept the installed skill's format and showed no change.

## Repro
Run a headless `claude -p` probe of a working-tree skill via a project symlink without `--setting-sources project,local`, with a real-home install of the same skill name present. The real home wins.

## Root cause
Claude Code's setting-source resolution can prefer the installed (user-level) copy of a skill over the project/working-tree copy when both exist under the same name and the probe does not pin sources; the operator's own `sh-agent install` mid-probe created exactly that collision.

## Fix
Probe kits now pass `--setting-sources project,local`. Three reruns then showed 3 of 3 runs reading the edited file. See `.agents/docs/testing/2026-09-29-memory-loop-phase-1-behavioural.md` § split-store, plan D15.

## Prevention
Already written: `AGENTS.md` § Testing — "a headless probe of a working-tree skill passes `--setting-sources project,local`, or first shows the real home has no installed copy of that name."

## Related
Plan `.agents/docs/plans/2026-09-29-dreamer-memory-loop.md` D15; testing log above.
