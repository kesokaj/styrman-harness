---
owner: bugs
date: 2026-09-15
id: BUG-017
status: fixed
area: tooling/workflows
---
# BUG-017: Workflow resumeFromRunId re-ran finished proposer agents live instead of replaying them from cache

## Symptom
A multi-agent research Workflow run (wf_921298bb-241) had 37 verification/critic agents fail with "You've reached your Fable limit" after the operator switched the session model mid-run with `/model`. The main agent pinned `model: 'opus'` on the failed calls and resumed with `resumeFromRunId`. The journal then showed new "started" entries for `propose:elephant-goldfish`, `propose:goldfish`, `propose:beads` — in a different order than the first run (first run: goldfish, beads, elephant-goldfish). The three proposer agents, whose prompts and options were unchanged, ran live again instead of replaying from cache, which would have produced a different candidate set than the one already read.

A second failure in the same task: the initial workflow script launch failed immediately with `ReferenceError: goldfish is not defined`, from `REPO.elephant-goldfish.path` being parsed as subtraction on an unbracketed hyphenated key.

## Repro
Not reproduced against the Workflow runtime directly. Observed: run wf_921298bb-241, journal start order per proposer differed between the original run and the `resumeFromRunId` resume after the model-pin fix.

## Root cause
The Workflow tool replays only "the longest unchanged prefix of agent() calls" on resume; the first edited/new call and everything after it runs live. The script's pipeline() stage-2 calls (one proposer per repo) start in completion order of stage 1, which is non-deterministic across runs, so the prefix broke at the first proposer even though its prompt/options were identical. This matches the tool's documented behavior and the observed order change; the runtime itself was not inspected, so this is the likely cause, not a proven one.

## Fix
Stopped the resumed run (task wruz8wwpf). Wrote the cached reader briefs, harness map, and 60 candidates to `.agents/scratch/research/` (briefs/, batches/, harness-map.json) and launched a new verification-only workflow (wf_1e892864-c31) that reads them by path with the model pinned on every agent; it completed with 49 agents, 0 errors. The hyphenated-key script bug was fixed in place with bracket access (`REPO['elephant-goldfish'].path`).

## Prevention
Before resuming a Workflow run whose earlier stage used concurrent/pipelined agent() calls, assume `resumeFromRunId` cannot guarantee a stable prefix — completion-order-dependent stages break the replay even with unchanged prompts. To reuse finished agents after a failure, write their outputs to scratch files (`.agents/scratch/`) and start a new run that reads them by path instead of resuming. On any long-running workflow, pin `model` on every agent explicitly, since the session model can change mid-run and pinning only the failed calls does not protect calls not yet started. Always access hyphenated object keys with bracket notation, never dot notation.

## Related
BUG-013 (read a run's exit status/behavior before citing it as expected).
