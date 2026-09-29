---
owner: main
date: 2026-09-29
---
# Handoff — 2026-09-29 — final: Linux box on v1.11.0, scratch and old feedback clone removed

Session ID: f538b90a-5aa7-44d1-8311-a160b5bb0167
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 0d7dc32; this handoff's commit follows it.

## Context

The closing steps of the session that released v1.9.0, v1.10.0 and v1.11.0. The full account is in `handoffs/2026-09-29-0711-v1.11.0-released.md` and `handoffs/2026-09-29-0717-session-closed.md`; this handoff records only what happened after them, all of it operator-requested housekeeping on this machine.

## Completed

- **Scratch cleared:** `.agents/scratch/` holds only the tracked `.gitkeep` (3.5G removed; a Go module cache inside needed `chmod -R u+w` first).
- **Real home updated and smoked** (operator's request): `sh-agent update` current, drift none, GANNET on both vendors; the three smoke checks PASS on Claude Code and Antigravity. Log: `testing/2026-09-29-verify-real-home-v1.11.0.md` (0d7dc32).
- **Old feedback clone removed:** `~/.config/sh-agent/feedback` held no commit GitHub lacks (checked after a fresh fetch) and no local changes; deleted. `sh-agent update` no longer prints the "kept" warning.
- No commit since the previous handoff changed a path outside `.agents/docs/`.

## Open / In Progress

- Nothing. `ROADMAP.md` "Now" and `TODO.md` are empty; no plan has an unticked box.
- Requirements edits awaiting approval: none.

## Key Decisions

- None new. Earlier ones are in the two handoffs named above.

## Corrections

- None since the previous handoff.

## Agent notes (unverified)

- Carried from the previous handoff: BUG-095's prevention rule is a candidate for `skills/review-standards/SKILL.md` at the next retro; retro proposal 7 is `queued` in `~/.config/sh-agent/proposals/2026-09-23-retro.md`.
- The Antigravity verifier wrote its smoke log under the same name as Claude Code's and overwrote it; both read PASS. A future smoke could give each CLI its own fixture clone.

## Environment State

- Branch: `main`, clean, pushed. Tags v1.9.0, v1.10.0, v1.11.0 pushed and verified.
- Real home on this box: installed from `0d7dc32`, drift none. `~/.config/sh-agent/` holds `feedback.jsonl`, `proposals/`, `retros.jsonl`, `state.json` only.
- No running processes; no scratch.

## Artifacts to Read

- `handoffs/2026-09-29-0711-v1.11.0-released.md` — what v1.11.0 ships, what was struck, retro statuses.
- `reviews/2026-09-29-preflight-v1.11.0.md` — the gate report.
- `testing/2026-09-29-verify-real-home-v1.11.0.md` — this box's smoke after the update.

## Next Steps

1. **Operator, other machines:** `sh-agent update` on the Mac (on v1.7.0) and the devbox, then `hash -r`, `sh-agent status`, and the three smoke checks. Not tracked here; the operator said they need not report back.
2. Nothing else is queued.
