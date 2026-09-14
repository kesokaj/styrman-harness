---
owner: main
date: 2026-09-14
---
# Handoff — 2026-09-14 — Known-limits sentence approved and applied; nothing pending in the requirements

Session ID: session_01CtMYM9AGp33edxBQGURNaB
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean once this handoff's commit is pushed

## Context

The operator: "approve the Known-limits sentence". The sentence had been proposed in `2026-09-13-improvement-items-applied.md`, re-raised by `/retro --harness`, and carried in every handoff since. Supersedes `2026-09-14-retro-harness.md`.

## Completed

- `REQUIREMENTS.md` Known limits, first bullet: "The Mac build is confirmed at the operator's next `sh-agent update` there (see the rebuild plan)." → "Confirmed on the Mac for Antigravity on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); Claude Code on the Mac is unverified." Per the hard rule: proposal → approval → ADR-019 → `REQUIREMENTS.md` and `ROADMAP.md` (unchanged, empty by design) in one commit. Entry `.agents/docs/changes/2026-09-14-known-limits-sentence.md`.

## Open / In Progress

Nothing queued in this repo; `TODO.md` and `ROADMAP.md` empty.
- Requirements edits awaiting approval: none.
Operator-only, unchanged: the Mac (`sh-agent update`, the three smoke checks on each CLI; Claude Code there is the unverified half the new sentence names); Go ≥ 1.24.6 on both machines; the downstream mirrors; improvement-review item 20; the full skill matrix.

## Key Decisions

- The sentence records what was verified and what was not, instead of a promise about the next update.

## Corrections

- None.

## Environment State

- Real home on this box updated after the push (`sh-agent status`: stamped `binary`, no drift). No running processes.

## Artifacts to Read

- `.agents/docs/handoffs/2026-09-14-retro-harness.md` and `2026-09-14-documentation-steps.md` — the rest of today's state and the operator items

## Next Steps

1. Operator, on the Mac: `sh-agent update`, then the three smoke checks on each CLI — that closes the "unverified" half.
2. The mirrors, the Go upgrade, item 20 — as before.
