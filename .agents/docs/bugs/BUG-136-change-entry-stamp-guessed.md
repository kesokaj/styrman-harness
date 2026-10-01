---
owner: bugs
date: 2026-10-01
id: BUG-136
status: fixed
area: trail naming
---
# BUG-136: change entry named with a guessed time stamp

## Symptom
The Phase 0 change entry was first written as `changes/2026-09-30-2330-audit-fix-phase-0.md` at 23:02. Caught and renamed before staging.

## Repro
Compare the file's stamp with its commit time.

## Root cause
The main agent typed the stamp from an estimate during a long turn instead of taking it from `date +%F-%H%M`. Same class as BUG-073.

## Fix
Renamed to `2302` before the commit. Phase 5 of the plan (WS2/TR3) adds "the stamp is `date +%F-%H%M` at write time" to the handoff skill and the changelog agent.

## Prevention
A time-stamped file name is produced by running `date +%F-%H%M` in the same tool call that writes the file, never typed.

## Related
Recurrence of BUG-073. Audit ID TR3; plan `plans/2026-09-30-audit-fix.md`.
