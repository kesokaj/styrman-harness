---
owner: bugs
date: 2026-09-29
id: BUG-100
status: fixed
area: agents/dreamer
---
# BUG-100: dreamer sampled `bugs/INDEX.md` and missed the newest rows

## Symptom
On its first E2E run over this checkout's real trail, the dreamer listed BUG-084 through BUG-091 and never opened BUG-096, producing zero rule candidates where a real candidate existed.

## Repro
Dispatch `agents/dreamer.md` (pre-fix) over `bugs/INDEX.md` with a `Since` date bound; it reads a subset of rows ("recent") rather than every row on or after that date.

## Root cause
The dreamer's definition told it to read the index for candidates without stating that "every row on or after `Since`" is required; left open, the agent sampled a plausible-looking recent window and stopped there.

## Fix
The definition now reads every index row dated on or after `Since` and opens each record. The rerun produced the expected candidate. See `.agents/docs/retros/2026-09-29-dreamer.md` (first run, missed) and `-2.md` (rerun, found).

## Prevention
An agent told to read an index states the bound explicitly (a date, or "all") and says "every row", never "recent rows" — an unbounded instruction over a growing table is read as a sample.

## Related
BUG-098, BUG-099 (same task, same session); `.agents/docs/retros/2026-09-29-dreamer.md`, `-dreamer-2.md`.
