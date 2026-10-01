---
owner: bugs
date: 2026-10-01
id: BUG-132
status: fixed
area: dream answer
---
# BUG-132: one extra JSON key in the model's answer failed the run every night

## Symptom
A model answer with one extra key (`"notes"`) was `fail invalid output`; the watermark stayed, so the same prompt was retried every night.

## Repro
`TestAnswerToleratesExtraKeys`.

## Root cause
`parse()` in `internal/dream/run.go` (line 339) used `DisallowUnknownFields()` on free-form model output. A model adding a harmless key turned into a permanent failure, and a failure holds the watermark.

## Fix
Unknown keys are ignored; the fields actually used are validated as before. In the working tree, committed with Phase 1.

## Prevention
A decoder of a model's free-form JSON ignores unknown keys and validates the fields it uses; strict decoding of model output is a nightly failure loop.

## Related
Audit ID SD5 (`plans/2026-09-30-recon-audit-fix.md`); plan `plans/2026-09-30-audit-fix.md` Phase 1; `testing/2026-09-30-verify-audit-fix-phase-1.md`.
