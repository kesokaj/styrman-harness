---
owner: main
date: 2026-09-30
---
# Phase 2 of the parity plan: memory on both CLIs through `sh-agent memory`     (2026-09-30)

## Summary
Both CLIs now take and load memory the same way (R10).
- The new `project-memory` skill says when to save (the four types) and how (`sh-agent memory add`). `feedback-log` rule 7 calls the verb on every CLI.
- The context's resume step runs `sh-agent memory show` on every CLI.
- `/init-project` Step 5b writes `"autoMemoryEnabled": false` beside `autoMemoryDirectory`, switching Claude Code's own writer off.
- `/dream` Step 1 always uses `<root>/.agents/memory` (the Antigravity-only check kept for the by-hand move hint). `agents/dreamer.md`, `/retro`, `/distill`, the memory spec and README say the same.

Verified on three sandbox homes (Claude Code only, Antigravity only, both):
- A memory saved on one CLI was recalled by a fresh session on the same CLI and, where both are installed, on the other.
- A control with the pre-change harness showed Claude Code's native writer (`originSessionId`); with the key off, only the verb wrote.
- The changed `dreamer` read the store at `<root>/.agents/memory` in 5 of 5 runs.

The gate verifier's first pass failed on a leftover line (`/dream`: "because that is all Claude Code loads at start") and on my log's unexplained "0 auto-denied" claim. Both were fixed, and the re-run passed.

## Type
feature (phase 2 of 6)

## Files changed
- `skills/project-memory/SKILL.md` (new)
- `skills/{feedback-log,init-project,dream,retro,distill}/SKILL.md`
- `agents/dreamer.md`
- `context/AGENTS.md` (resume step; "auto memory" → "memory"): 21906 of 22000 bytes
- `.agents/docs/specs/sh-agent/memory-loop.md`
- `README.md`
- `.agents/docs/testing/guide.md` § 4 (three-setup kit)
- the plan (a § 3b report-tier task added for Phase 1's `takesValue` change)

## Agents dispatched
verifier (`testing/2026-09-30-verify-parity-memory-dream-phase-2-gate.md`: FAIL, then PASS). The behavioural log `testing/2026-09-30-verify-parity-memory-dream-phase-2.md` is the main agent's own, audited file by file by that verifier.

## Verification
`go test -race -count=1 ./...` 406 pass. Scratch installs for the three vendors show no drift. E2E 1 and E2E 2 as above. `dreamer` 4 of 5 on the matrix check: r3 cited the bug rows by line, not by id.

## Requirements
R10 (ADR-047), the Non-goals memory-keys exception (ADR-047a).
