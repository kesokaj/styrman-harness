---
owner: changelog
date: 2026-09-13
---
# Real home updated to a35b567 and smoke-checked   (2026-09-13)

## Summary
On the operator's instruction, `sh-agent update` was run in the real home on the Linux box. The 2026-09-03 binary rebuilt itself from the current source (commit a35b567), installing 73 new/changed files out of 260 total managed. The resulting binary `a35b567582fc` reported no drift and no warnings. Verification: the three smoke checks from `.agents/docs/testing/2026-09-03-verify-real-home.md` ran from the real home in fresh demo clones on both Claude Code 2.1.270 and Antigravity 1.2.2, all PASS (context questions including step 4 of Definition of Done; `golang` auto-load; `verifier` dispatch writing logs in clones). No source files were changed. No vendor settings were edited.

## Type
task

## Files changed
None — no source or configuration files changed. Harness binary rebuilt and reinstalled.

## Agents dispatched
changelog (this entry only)

## Verification
Verification log: `.agents/docs/testing/2026-09-13-verify-real-home.md` — `sh-agent update` step, six smoke check runs across two CLIs, all PASS, zero parse errors in Antigravity log. Handoff: `.agents/docs/handoffs/2026-09-13-real-home-updated.md`.

## Requirements
R6 (operator's action, run on their instruction).
