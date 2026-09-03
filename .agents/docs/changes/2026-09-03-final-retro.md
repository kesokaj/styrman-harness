---
owner: main
date: 2026-09-03
---
# Final retro of the verification day: two proposals applied     (date: 2026-09-03)

## Summary

`/retro` over the session's three feedback entries (one correction, two gaps — all the same shape: "done" said at a narrower scope than the reader meant). Proposal file: `~/.config/sh-agent/proposals/2026-09-03-retro-final.md`. Applied:

- `skills/handoff/SKILL.md`: a cold-reader check closes Step 6 — a handoff is finished when a fresh reader (or the `explorer` agent) has resumed from it and the gaps are fixed, not when it is saved. Today's first cold read found four gaps in a minutes-old handoff.
- `AGENTS.md` (repo) Testing: tested means tested from where it runs — after any install into a real home, the three smoke checks in `testing/2026-09-03-verify-real-home.md` run from that home on each CLI.

Cost count per ADR-009: two dispatches in this repo today (`reviewer` on Sonnet, one real finding; `Explore` for the doc audit, nine findings). The ADR-009 no-tuning clock (ends 2026-09-17) is unchanged: these come from verification, not use.

## Type

task

## Files changed

`skills/handoff/SKILL.md`, `AGENTS.md`

## Agents dispatched

None.

## Verification

`go test` ok; scratch install 256 files, no drift; real home reinstalled at this commit.

## Requirements

R7 (proposals stay in `~/.config/sh-agent/`).
