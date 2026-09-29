---
owner: main
date: 2026-09-29
---
# Matrix row `dream`: the store inside the clone     (date: 2026-09-29)

## Summary
The `dream` row of `testing/skill-matrix.md` still seeded its store outside the clone at `<run>/mem` and grepped for that path in the store line — a value `/dream` refuses since the v1.12.0 gate's exact-path rule (run 9 of `testing/2026-09-29-memory-loop-phase-2-behavioural.md`). The row now seeds the clone's `.agents/memory/` and greps for it. Found by the operator's question after the release ("Have you fixed everything in the gate?"); the second gate security round had named rows 31 and 32 as unchanged.

## Type
fix (docs only, two cells of one row)

## Agents dispatched
None: a docs-only change under ten lines (context § Definition of Done).
