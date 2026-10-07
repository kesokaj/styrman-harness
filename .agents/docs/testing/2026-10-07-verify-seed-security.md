---
owner: main
date: 2026-10-07
---
# Verify — the `security` agent at effort high vs max, seeded weaknesses

Range reviewed: 6c9eaf9 plus the uncommitted lean-process-3 tree (`testing/seed-security/run.sh`; `agents/security.md` at `effort: high` since ADR-054a).

## Status: PASS (a measurement, plan lean-process-3 D4)

## Method
- **The fixture.** `testing/seed-security/run.sh` builds a Go service whose one commit plants six flaws:
  1. shell injection through `sh -c`
  2. path traversal
  3. a hard-coded admin password
  4. an unauthenticated `/admin/purge`
  5. a non-constant-time token compare
  6. `math/rand` session IDs
- **The run.** Claude Code headless (the ab-kit isolation, with `kit-preamble.sh`) is asked only to dispatch `security` over that commit. The sandbox copy of `security.md` carries the effort under test.
- **The score** is the six flaws matched by keyword in the artifact.
- **The scorer** was validated before any run: a synthetic all-six artifact scored 6/6 and a no-findings one 0/6. The `high-2` artifact was read by hand: all six are real findings with file:line, plus four unplanted ones, including the token check failing open when `SVC_TOKEN` is empty.

## Results

| Run | Effort | Found | Time | Cost |
|---|---|---|---|---|
| high-1 | high | 6/6 | 98 s | $0.49 |
| high-2 | high | 6/6 | 76 s | $0.44 |
| high-3 | high | 6/6 | 103 s | $0.46 |
| high-4 (revised kit) | high | 6/6 | 83 s | $0.48 |
| high-5 (fixture bound to 127.0.0.1:0) | high | 6/6 | 126 s | $0.48 |
| high-6 (with `kit_leftovers`) | high | 6/6 | 109 s | $0.54 |
| max-1 | max | 6/6 | 343 s | $1.16 |

Every run printed `kit: real config unchanged` and `kit: no credentials left in KIT_HOME`.

At `high`, `security` found all six planted flaws in 6 of 6 runs, the same as `max` in 1 of 1, at about 26% of the time and 40% of the cost. This is a sample of five runs on one fixture, not a proof that `high` never misses what `max` would catch.
