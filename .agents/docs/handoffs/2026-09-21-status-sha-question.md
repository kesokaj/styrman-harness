---
owner: main
date: 2026-09-21
---
# Handoff — 2026-09-21 — why status shows a sha; nothing changed; the consolidation Feature is still next

Session ID: 5209d64b-84a1-4fe8-a4e2-434185fe9d52
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 3e92bef (this handoff's commit follows it), clean and pushed. Tags `v1.0.0` … `v1.4.0`, all published.

## Context

A one-question session. The operator asked why `sh-agent status` on this box shows a sha (`binary 3e92befe39e4`) rather than a version. Answered from `cmd/sh-agent/main.go` (`status`, `releaseFromInfo`, `versionFromInfo`) and `delivery.md` § The model; the operator said "leave it as is". No file changed; no plan is active; `ROADMAP.md` "Now" and `TODO.md` are empty.

## Completed

- Answer only. A checkout place's binary is built by `go build` from the checkout, and `status` renders a checkout place's binary as its VCS revision (twelve characters, `+dirty` when modified); it renders a `vX.Y.Z` only for a binary-only place with a release binary. This box's checkout is one commit past `v1.4.0` (`git describe` → `v1.4.0-1-g3e92bef`), and even at a clean tag the checkout branch of `status` would still print the sha. Both machines-on-release show `v1.3.0`/`v1.4.0` because they are binary-only places. Nothing is broken: `installed from commit 3e92bef`, `binary 3e92befe39e4`, drift none.

## Open / In Progress

- Everything carried from the 2026-09-18 handoff stands unchanged: the Mac and devbox `sh-agent update` to v1.4.0 (not confirmed this session), the consolidation Feature first next session, the carried Tasks (`skills/preflight/SKILL.md:77` first), `close-map-gaps` waiting on another place's proposal (ADR-035).
- Requirements edits awaiting approval: none.

## Key Decisions

- "Leave it as is" — the operator's: a checkout place keeps showing the sha; no `v1.4.0-1-g<sha>` rendering. Recommended and accepted because the sha is the honest identity for a place that follows `main`.

## Corrections

- None this session (no feedback log entry).

## Agent notes (unverified)

- This box's real home is now at 3e92bef (`installed 2026-09-21T08:37:37+02:00 from commit 3e92bef`), so the hourly timer has moved it past the b06268d the previous handoff recorded.

## Environment State

- `main` at 3e92bef + this handoff, clean, pushed. No source change since v1.4.0's release commit ce39eac beyond the previous handoff.
- Real home on this box: 3e92bef, drift none, 27 skills, 7 agents. Mac and devbox: last known v1.3.0 (update to v1.4.0 not confirmed).
- `ROADMAP.md` "Now" empty; `TODO.md` empty.

## Artifacts to Read

- `.agents/docs/handoffs/2026-09-18-v1.4.0-released.md` — the previous handoff; its Open and Next Steps are still the work.
- `.agents/docs/specs/sh-agent/delivery.md` § The model — why a checkout place is identified by commit.
- `cmd/sh-agent/main.go` — `status`'s binary line, `releaseFromInfo`, `versionFromInfo`.

## Next Steps

1. **Next session:** the consolidation Feature (scope confirmation first; the operator decides what stays) — unchanged from the 2026-09-18 handoff.
2. **Operator:** `sh-agent update` on the Mac and devbox if not yet done; `sh-agent version` should print `sh-agent v1.4.0`.
3. Then the carried Tasks, `:77` first; `close-map-gaps` when ADR-035's other proposal arrives.
