---
owner: main
date: 2026-09-17
---
# BUG-065: the `feedback` skill renamed to `feedback-log`     (2026-09-17)

## Summary
Typing `/feedback` in Claude Code loaded the harness's passive correction logger instead of Claude Code's built-in feedback command, so a queued bug report could not be sent. The operator chose to rename the skill: `skills/feedback` is now `skills/feedback-log`, and its six mentions are updated. The skill loads by its description, not by being typed, so nothing the operator types changes.

## Type
bugfix

## Files changed
- `skills/feedback/SKILL.md` → `skills/feedback-log/SKILL.md` — frontmatter `name: feedback-log`
- `context/AGENTS.md` — § Learning loop names `feedback-log`
- `skills/distill/SKILL.md` (two), `skills/handoff/SKILL.md`, `skills/retro/SKILL.md` (two) — the skill's name
- `README.md` — the map's session label
- `.agents/docs/bugs/BUG-065-feedback-skill-shadows-builtin-slash-command.md`, `INDEX.md`

## Agents dispatched
reviewer (APPROVE), bugs (BUG-065). This entry was written by the main agent from the diff.

## Verification
- `go test -count=1 ./...` passes.
- A fresh scratch install (`.agents/scratch/rename-home`): `feedback-log` present and `feedback` absent for both vendors; drift none.
- Upgrade: installing over an older scratch home (`.agents/scratch/retro2-home`) that held `feedback` removed it and left `feedback-log`; drift none.
- Review `.agents/docs/reviews/2026-09-17-review-4c63232-feedback-log-rename.md`.

## Docs
Context, three skills and the README label, above.

## Requirements
n/a

## Not exercised
- `/feedback` reaching Claude Code's built-in in an interactive session, and `feedback-log` still auto-loading on a correction — both need `sh-agent update` and an interactive session.
- Whether other harness names (`bug`, `review`) clash with a built-in command of either CLI: no local list of built-ins to check against.
