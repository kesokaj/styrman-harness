---
owner: main
date: 2026-10-06
---
# Handoff — 2026-10-06 — releases without a key: v2.4.0 and v2.4.1 shipped, session closed

Session ID: ad142c6e-ead7-4625-9e5c-9767b55c9268 (its transcript was deleted on the operator's word, because it held a pasted private key: it cannot be resumed)
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

The v2.4.0 tag was blocked by a missing release key. The key was then exposed (pasted into this session; it was also dave's GitHub login). After a rotation attempt (ADR-052), the operator chose to drop release signing and keep versioning (ADR-053, R9 approved). v2.4.0 shipped as the bridge, v2.4.1 followed, and the follow-ups were closed.

## Completed

- v2.4.0 (`b68f6b9`, unsigned tag, plus one bridge `SHA256SUMS.sig` by the old key) and v2.4.1 (`7ea9068`) released; `update` and `--from-zip` check `SHA256SUMS` only, and the gh host is pinned (ADR-053, ADR-053a).
- Every place on v2.4.1 (operator, 2026-10-06; einherjar by hand after BUG-161).
- BUG-160 (kit credential check, fails closed) and BUG-161 (bridge check in `/preflight` Step 6), both graduated to mechanisms.
- GitHub Actions off; the old key gone from GitHub, scratch, the paste cache and this session's transcript.

## Open / In Progress

- None for the agent.
- Requirements edits awaiting approval: none.

## Key Decisions

- No release signing (ADR-053): trust is the GitHub account plus HTTPS; versioning unchanged. GitHub Actions for releases rejected: it breaks "no CI", keeps the token need on binary-only places, and reopens the publish path just closed.
- `gh` token not rotated (operator, 2026-10-06: single-user box, risk accepted; do not raise it again).
- v2.4.1 stays a PATCH; v2.4.0 a MINOR.

## Corrections

- None logged (`~/.config/sh-agent/feedback.jsonl` does not exist on this machine).

## Agent notes (unverified)

- The operator's corrections this session were not logged by `feedback-log`: "I told you that I did not want to use a signing key" (2026-10-05), and frustration at multi-step key handling. Candidate retro: prefer the simplest path the operator states over a more secure one they did not ask for.
- Trail slips by the main agent, left as they are: a committed `changes/` entry appended to (`29cd652`), and a committed handoff edited (`7c67518`). Both directories are append-only or immutable by rule.
- This box's home credentials are rewritten on every login shell by the styrman broker (`/etc/profile.d/zz-styrman-credentials.sh`): never keep a key at a default path here (operator declined saving it to project memory).

## Environment State

- Branch: main, clean, in sync with origin; newest tag v2.4.1.
- Running processes: none.

## Artifacts to Read

- `.agents/docs/architecture/DECISIONS.md` ADR-053, ADR-053a — the unsigned release model.
- `.agents/docs/deploy/remote.md` § Releases — the procedure, no key.
- `.agents/docs/bugs/BUG-161-release-after-bridge-stranded-binary-places.md` — why `/preflight` now checks for a bridge.

## Next Steps

1. The next release owes the Antigravity run of the changed `preflight` skill (AGENTS.md § Release); the bridge check is prose with no test.
