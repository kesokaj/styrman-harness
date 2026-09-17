---
owner: bugs
date: 2026-09-17
id: BUG-065
status: fixed
area: skills/feedback
---
# BUG-065: `/feedback` loaded the harness skill instead of Claude Code's built-in

## Symptom
Operator typed `/feedback` in Claude Code intending the CLI's built-in feedback
command (send a queued bug report to Anthropic). The harness's own skill named
`feedback` (`skills/feedback/SKILL.md`, the passive correction logger) loaded
instead; the report was not sent.

## Repro
Type `/feedback` in Claude Code with the harness installed. Not independently
re-reproduced here beyond the reported symptom.

## Root cause
A harness skill was named `feedback`, identical to a built-in slash command of
one of the target CLIs (Claude Code). Skill names are not checked against
either vendor's built-in command namespace before install, so the harness's
own skill wins the name and shadows the built-in.

## Fix
Skill renamed to `feedback-log`: `git mv skills/feedback skills/feedback-log`,
frontmatter `name: feedback-log`, and its six mentions updated in
`context/AGENTS.md` § Learning loop, `skills/distill/SKILL.md` (two),
`skills/handoff/SKILL.md`, `skills/retro/SKILL.md` (two), and the README map's
session label. `go test -count=1 ./...` passes; fresh scratch install shows
`feedback-log` and no `feedback` for both vendors, drift none; installing over
an older scratch home holding `feedback` removes it and leaves
`feedback-log`. Not yet verified: that `/feedback` now reaches Claude Code's
built-in interactively, and that `feedback-log` still auto-loads by
description — both need the operator's `sh-agent update` plus an interactive
session.

## Prevention
A skill, agent, or command in this harness never takes the name of a built-in
slash command of either target CLI; check the CLI's built-in command list
before adding or renaming one. This is the same class of clash one level up
from `AGENTS.md` (project) § Rules for working here, which already forbids a
project skill reusing a harness name — that rule should be read as covering
vendor built-ins too, or a matching check added alongside it.

## Related
`AGENTS.md` (project) § Rules for working here (project-vs-harness name
clash rule, same mechanism).
