---
owner: bugs
date: 2026-09-29
id: BUG-101
status: fixed
area: skills/init-project (Step 5b, memory migration)
---
# BUG-101: memory migration trusted an `autoMemoryDirectory` value that could come from a tracked file

## Symptom
Security review (`.agents/docs/reviews/2026-09-29-security-memory-loop.md`) finding (a): `/init-project` Step 5b used an existing `autoMemoryDirectory` value as the migration target even when that value came from a possibly tracked (committed) settings file, so a committed settings file in a cloned repo could redirect the operator's memory migration to an attacker-chosen directory.

## Repro
Not reproduced against a live run; found by review reading Step 5b's logic against the settings-file trust boundary (tracked vs. untracked/local/user).

## Root cause
Step 5b read the value of `autoMemoryDirectory` without first checking which file it came from; a value from a file the repo controls (tracked) was treated as equally trustworthy as a value only the operator's own machine controls (untracked local or user-level).

## Fix
Files now move only into `~/.config/sh-agent/memory/<slug>/`; a value found in a tracked file is reported to the operator and nothing moves on its authority.

## Prevention
A step that derives a filesystem destination from a project-supplied setting checks whether that setting's source file is tracked before trusting its value for a write or move; a tracked source is reported, never acted on.

## Related
BUG-102, BUG-103 (same review, same task); `.agents/docs/reviews/2026-09-29-security-memory-loop.md` finding (a); spec `.agents/docs/specs/sh-agent/memory-loop.md`.
