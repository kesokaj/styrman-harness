---
owner: bugs
date: 2026-09-29
id: BUG-102
status: fixed
area: skills/dream (Step 1, settings read)
---
# BUG-102: `/dream` Step 1 read the committed settings file with no trust gate

## Symptom
Security review finding (b): `/dream` Step 1 read the committed `.claude/settings.json` with no trust gate, so a cloned repository could point the dream skill and the dreamer agent at any directory of the attacker's choosing.

## Repro
Not reproduced against a live run; found by review reading Step 1's settings resolution order against the trust boundary between tracked and untracked/user settings files.

## Root cause
Step 1 resolved the memory-directory setting from whichever settings file had it set, including the tracked (committed, repo-controlled) `.claude/settings.json`, without restricting the source to files the operator's own machine controls.

## Fix
Step 1 now reads only the untracked local settings file or the user settings file; a resolved value outside `~/.config/sh-agent/memory/` and `~/.claude/projects/` is refused.

## Prevention
Same rule as BUG-101: a step that derives a filesystem destination from a project-supplied setting checks the source file's tracked status before trusting it, and constrains the resolved path to a known-owned directory tree regardless of source.

## Related
BUG-101, BUG-103 (same review, same task); `.agents/docs/reviews/2026-09-29-security-memory-loop.md` finding (b); candidate for a standing `skills/dream/SKILL.md` / `skills/init-project/SKILL.md` rule on settings trust, given two occurrences in this review alone.
