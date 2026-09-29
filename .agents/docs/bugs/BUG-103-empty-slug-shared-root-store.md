---
owner: bugs
date: 2026-09-29
id: BUG-103
status: fixed
area: skills/init-project / skills/dream (slug derivation)
---
# BUG-103: an empty slug from a failed `git rev-parse` made the store the shared root, and migration ran even when the settings write was refused

## Symptom
Security review finding (c): when `git rev-parse` failed, the derived slug was empty, and the empty slug made the effective memory store the shared root directory rather than a per-project one. Separately, a refused settings write still went on to migrate memory into a store Claude Code would not load.

## Repro
Not reproduced against a live run; found by review reading the slug-derivation and migration order against the failure path of `git rev-parse` and against a refused settings write.

## Root cause
Two missing checks in sequence: (1) the slug derivation did not stop on an empty result, letting an unset slug collapse the per-project path to the shared parent directory; (2) the migration step ran unconditionally instead of only after the settings write it depends on had actually succeeded, so a refused write and a completed migration could disagree about where memory lived.

## Fix
Slug derivation now stops on an empty slug. Migration now runs only when the settings write is confirmed in place.

## Prevention
A step that derives a path from a project file or a command's output checks the result (non-empty, well-formed) before any move; a step that depends on a prior write's success checks that write succeeded before acting on its result, never assuming order implies success.

## Related
BUG-101, BUG-102 (same review, same task, same class of missing-check-before-move); `.agents/docs/reviews/2026-09-29-security-memory-loop.md` finding (c). Three related records (BUG-101/102/103) share "check the source/result before trusting it for a move" — candidate for a `context/AGENTS.md` § Hard rules bullet or a `skills/init-project`/`skills/dream` standing rule; actionable (a specific tracked-file / empty-result check) and generalizable beyond this feature.
