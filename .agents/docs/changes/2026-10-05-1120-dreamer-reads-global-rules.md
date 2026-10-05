---
owner: main
date: 2026-10-05
---
# Dreamer and dream read the global rule files

## Summary

`agents/dreamer.md:32` and `skills/dream/SKILL.md:30,39` name the global rule files (`~/.claude/rules/*.md`, `~/.gemini/config/rules/*.md`) among what already exists, so a memory whose rule lives there is left alone or marked graduated rather than proposed again. Left for the operator by `reviews/2026-10-05-review-worktree-docs-2.md`; approved in their own turn: "You fix both 1 and 2. Approved."

## Type

task

## Files changed

- `agents/dreamer.md` — What you read § 3.
- `skills/dream/SKILL.md` — the read list and the `graduated` verdict.
- `.agents/docs/specs/sh-agent/memory-loop.md:23` — the dreamer's Reads row (by `documentation`).

## Agents dispatched

reviewer (docs lens, `reviews/2026-10-05-review-worktree-docs-3.md`, APPROVE WITH NITS; Outcome: 1 fixed), documentation (spec row), verifier (PASS). Entry written by main from the diff.

## Verification

- `testing/2026-10-05-verify-dreamer-global-rules.md` (main, PASS): new definition 2 of 2, control 2 of 2 — the probe does not discriminate on Claude Code, which loads `~/.claude/rules/*.md` into subagents.
- `testing/2026-10-05-verify-dreamer-global-rules-go.md` (verifier, PASS): § 1 553 pass, 0 fail, 3 skip; § 2; § 3 drift none; derived Antigravity `dreamer` carries the text, one `## Records` each.

## Notes

Antigravity runs of the changed `dreamer` and `dream` are owed before the next tag (project release rule).
