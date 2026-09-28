---
owner: main
date: 2026-09-28
---
# Time-sortable trail names: E2E 1 and plan closed     (date: 2026-09-28)

## Summary
E2E 1 passed. A cold agy resume on a fresh clone of `9a17e0a` named the true newest handoff 2 of 2, although an older date-only handoff sorts last by name and every mtime is equal. The same held 2 of 2 on a clean demo project.

The old-context control on the demo was also right 2 of 2. It broke the tie with `git log` (9 calls against 7), so in a two-file tie the rule saves calls rather than changing the answer. The first attempt was invalid, for a missing agy settings file in the sandbox home and a control contaminated by the harness's own tree (BUG-091). The plan is closed and ROADMAP "Now" is empty.

## Type
feature

## Files changed
- `.agents/docs/testing/2026-09-28-verify-time-sortable-names-e2e1.md` — E2E 1.
- `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` — the last boxes ticked, with evidence.
- `ROADMAP.md` — "Now" emptied.
- `.agents/docs/bugs/BUG-091-e2e1-sandbox-run-invalid.md` and `INDEX.md` — the record.

## Agents dispatched
- `bugs`: BUG-091.

## Verification
- **E2E runs:** 6 agy runs (listed in the testing log). The condition was proven before them: name-sort last is the older file, there is one mtime, and `git log` names the newest.
- **Credentials:** `chmod 600` and deleted by name after each batch.
- **Clones:** unmodified.
- **Source:** none changed since the verifier's PASS on `9a17e0a`'s tree (`testing/2026-09-28-verify-time-sortable-names-fixes.md`).

## Docs
No docs changed.

## Requirements
n/a
