---
owner: bugs
date: 2026-09-28
id: BUG-091
status: fixed
area: process / sandbox runs of the CLIs
---
# BUG-091: first E2E 1 attempt invalid — headless agy auto-denied, control run couldn't isolate context

## Symptom
Two agy runs from sandbox home `.agents/scratch/tsn` produced no output: headless agy auto-denied the
first shell command ("a tool required the \"command\" permission that headless mode cannot prompt for,
so it was auto-denied"). Separately, the must-fail control — agy from the real home with the old v1.6.0
context — answered correctly by citing the NEW rule.

## Repro
Run headless `agy` against a sandbox home built by copying only the credential, no
`~/.gemini/antigravity-cli/settings.json`; the first tool call needing the "command" permission is
auto-denied. Run the must-fail control inside a fresh clone of styrman-harness itself; the clone's
working tree carries `context/AGENTS.md` with the new rule, so the run can read it regardless of which
context text was injected.

## Root cause
1. The sandbox home was built by copying only the credential file, omitting the real home's
   `~/.gemini/antigravity-cli/settings.json` (which holds `toolPermission: always-proceed`,
   `artifactReviewPolicy`, `trustedWorkspaces`, and no secret) — headless mode has no prompt path, so a
   tool call needing permission is silently auto-denied and the run produces nothing.
2. The control for a context change was run inside a clone of styrman-harness, whose tree already
   contains the new rule in `context/AGENTS.md`; any run there can read the new rule directly from the
   repo regardless of which context was injected, so the control cannot isolate the context change.

## Fix
Copied the real home's agy `settings.json` (no credential in it) into the sandbox home beside the
copied token, and reran the E2E on the harness clone plus a clean demo project holding no harness
source. New context picked the right handoff 2/2 on each; old-context control 2/2 as well, verified via
`git log`: 9 calls against 7. See `.agents/docs/testing/2026-09-28-verify-time-sortable-names-e2e1.md`.

## Prevention
A sandbox home for headless agy runs gets the real home's
`~/.gemini/antigravity-cli/settings.json` (no credential in it) copied in beside the copied token; a run
that uses tools is checked for the "auto-denied" stderr line before its output is read. A control for a
context change runs in a project that does not contain the harness source, never in a clone of
styrman-harness.

## Related
BUG-087 (fixture state must be proven before a behavioural run), BUG-090 (recurrence of BUG-087's rule).
This is a third instance of "prove the run setup before trusting its output" — BUG-087, BUG-090, and
BUG-091 now share that shape; candidate for graduating a standing check into
`.agents/docs/testing/2026-09-03-verify-sandbox.md` or the behavioural-run procedure in `AGENTS.md` §
Testing.
