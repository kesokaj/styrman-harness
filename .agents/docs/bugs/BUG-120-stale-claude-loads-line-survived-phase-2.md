---
owner: bugs
date: 2026-09-30
id: BUG-120
status: fixed
area: skills/dream, docs (claim reversal)
---
# BUG-120: a stale "Claude Code loads" line survived Phase 2

## Symptom
Verifier FAIL (`.agents/docs/testing/2026-09-30-verify-parity-memory-dream-phase-2-gate.md`): `skills/dream/SKILL.md` said the index cap exists "because that is all Claude Code loads at start" after native loading had been switched off.

## Repro
`grep -n "loads" skills/dream/SKILL.md` before b629ae0.

## Root cause
Before declaring the claim reversed I grepped the setting's name (autoMemoryDirectory) and "auto memory", not the claim's own words ("loads", "Claude Code"), so the one restatement using neither survived.

## Fix
Fixed in b629ae0.

## Prevention
Before calling a claim reversed, grep its key words across every living file, not only the setting's name (the global context's describing-text rule, applied).

## Related
BUG-113 and BUG-096, BUG-102 (earlier records on skills/dream).
