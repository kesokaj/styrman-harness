---
owner: bugs
date: 2026-09-30
id: BUG-116
status: fixed
area: internal/dream (agy invocation), plan parity-memory-dream D7
---
# BUG-116: the agy tool-free probe used a toy prompt, so D7a declared an invocation that failed on the real prompt

## Symptom
D7 probe ran `agy -p '<one-line JSON request>' --mode plan --disable-slash-commands` from an empty untrusted cwd and returned the JSON; D7a recorded the invocation as working. The first real run with the 14 KB dream prompt called a `command` tool, was auto-denied and ended with no answer (`fail invalid output`). stderr: `--mode plan has no effect while slash command expansion is disabled`.

## Repro
Run `sh-agent dream` on a real project under the D7a invocation. Evidence: `.agents/docs/testing/2026-09-30-verify-parity-memory-dream-phase-4.md`.

## Root cause
A probe with a trivial payload does not exercise the model's urge to act on a real one (the transcripts hold commands and paths it wants to run). Also `--disable-slash-commands` switches `--mode plan` off, so the two flags together removed the restriction.

## Fix
D7b: a warning at the top of the prompt (both CLIs) and `--mode plan` without `--disable-slash-commands`. Later a private HOME that grants nothing (BUG-119).

## Prevention
A probe that decides an invocation runs with the real payload's size and shape, never a toy one; read the flags' own stderr warnings before recording a probe as a pass.

## Related
BUG-119 (the private HOME), BUG-115 (sandbox settings lesson), BUG-113 (skills/dream).
