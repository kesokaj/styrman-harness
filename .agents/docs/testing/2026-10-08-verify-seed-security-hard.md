---
owner: main
date: 2026-10-08
---
# Verify — the `security` agent on harder, multi-step seeded flaws, high vs max

Range reviewed: 280c1df plus the uncommitted plan-2026-10-08 tree (`testing/seed-security/run.sh` hard mode, `kit-preamble.sh`'s leftover check).

## Status: PASS (a measurement, plan 2026-10-08 D3)

## Method
- **The hard commit.** `seed-security/run.sh <label> <effort> hard` builds the skeleton commit, then ONE commit of three flaws, and dispatches `security` over that commit only:
  1. a session cookie whose HMAC check is skipped when the cookie has no signature part (`strings.Cut` found=false)
  2. `strings.HasPrefix(filepath.Clean(p), "/srv/files")` with no separator, so `/srv/files-evil` passes
  3. a token middleware that serves the request when the token store errors (fail-open)
- **The flaws were proven exploitable before use,** by a Go test in scratch: an unsigned "admin" cookie is accepted, the prefix bypass passes, and a store error calls the handler.
- **The scorer** was validated before any run: a synthetic all-three artifact scored 3/3 and a no-findings one 0/3. The kit's own fixture, extracted from `run.sh`, compiles and vets clean.
- **One `high` artifact (`hard-high-2`) was read by hand.** All three are real findings with file:line, plus five unplanted ones: an unescaped token in the store URL, an empty HMAC key accepted, a hand-rolled session format, an unauthenticated `/raw` that follows symlinks, and no client timeout.

## Results

| Run | Effort | Set | Score | Time | Cost |
|---|---|---|---|---|---|
| hard-high-1 | high | hard | 3/3 | 144 s | $0.52 |
| hard-high-2 | high | hard | 3/3 | 157 s | $0.53 |
| hard-high-3 | high | hard | 3/3 | 155 s | $0.56 |
| hard-max-1 | max | hard | 3/3 | 590 s | $1.75 |
| six-high-7 | high | six (default mode unchanged) | 6/6 | 146 s | $0.50 |

Every run printed `kit: no processes left by the run`, `kit: real config unchanged` and `kit: no credentials left in KIT_HOME`. These are the first runs through the preamble's new leftover check.

On multi-step, fail-open flaws, `high` found all three in 3 of 3 runs, the same as `max` in 1 of 1, at about 26% of the time and 31% of the cost. That answers the v2.5.0 security review's low 3 for this fixture. It remains a sample, not a proof.
