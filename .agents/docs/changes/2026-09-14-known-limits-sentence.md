---
owner: changelog
date: 2026-09-14
---
# Known limits sentence: verified facts, not promises     (2026-09-14)

## Summary

REQUIREMENTS.md Known limits, first bullet, updated to state what was verified and what was not, instead of a promise about future verification. The Mac was verified for Antigravity on 2026-09-04; Claude Code on the Mac remains unverified. Applied per the hard rule: proposal (2026-09-13) → operator approval (2026-09-14) → ADR-019 → REQUIREMENTS.md and ROADMAP in one commit. Completes the trail's visibility of stale Known-limits items.

## Type

task

## Files changed

- REQUIREMENTS.md — Known limits, first bullet: "The Mac build is confirmed at the operator's next `sh-agent update` there (see the rebuild plan)." → "Confirmed on the Mac for Antigravity on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); Claude Code on the Mac is unverified."

## Agents dispatched

changelog (this agent)

## Verification

Old sentence confirmed absent by:
```
grep -n "The Mac build is confirmed at the operator's next" REQUIREMENTS.md
# (no output)
```

New sentence confirmed present by:
```
grep -n "Confirmed on the Mac for Antigravity on 2026-09-04" REQUIREMENTS.md
# 35:Confirmed on the Mac for Antigravity on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); Claude Code on the Mac is unverified.
```

Only REQUIREMENTS.md line 35 changed; ADR-019 appended to DECISIONS.md; this entry written; handoff remains unchanged.

## Docs

REQUIREMENTS.md changed (one sentence in Known limits section); no README, spec, testing guide, or deploy guide involved. The sentence itself documents the verified scope.

## Requirements

R6 holds for the Mac (operator-only action on Mac: `sh-agent update`, smoke checks on both CLIs). The changed sentence satisfies the non-goal of hard-coding unverified paths; Known-limits now follow ADR-019 (verified facts named explicitly, gaps named, no promises).
