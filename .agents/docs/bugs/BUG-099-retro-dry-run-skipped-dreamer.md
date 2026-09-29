---
owner: bugs
date: 2026-09-29
id: BUG-099
status: fixed
area: skills/retro
---
# BUG-099: `/retro --dry-run` skipped the dreamer dispatch it was supposed to include

## Symptom
Two behavioural probes (t1, t2) refused to dispatch the dreamer agent under `--dry-run`, one verbatim: "Dreamer agent not run: it always writes an artifact file, and a dry run writes nothing."

## Repro
Run `/retro --dry-run` against the working-tree `skills/retro/SKILL.md` before the fix; the dreamer dispatch step is skipped in both probes.

## Root cause
The skill's usage line said `--dry-run` means "nothing written"; an agent reading that literally extended it to cover the dreamer's own artifact, treating "nothing written" as "no writes of any kind" rather than "no change to source/config".

## Fix
`skills/retro/SKILL.md` source 8 now states the dry run still dispatches the dreamer, because an agent's artifact is trail, not a change (`context/AGENTS.md` § Work). Reruns: 2 of 2. See `.agents/docs/testing/2026-09-29-memory-loop-phase-2-behavioural.md`.

## Prevention
A skill that gains a dispatch under a read-only mode states in the same sentence that the dispatched agent's artifact is exempt from that mode — do not rely on the reader inferring the trail/change distinction from elsewhere.

## Related
BUG-098, BUG-100 (same task); candidate for `context/AGENTS.md` § Definition of Done or the `bug`/`feature` skill's read-only-mode language, once a third instance recurs.
