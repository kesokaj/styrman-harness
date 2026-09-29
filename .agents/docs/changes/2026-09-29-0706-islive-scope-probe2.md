---
owner: main
date: 2026-09-29
---
# "Is it live" sentence scoped; sentences 1 and 6 probed     (date: 2026-09-29)

## Summary
The final v1.11.0 gate review (`reviews/2026-09-29-review-v1.10.0-740ca80-code.md`, APPROVE WITH NITS) left two nits, both closed here.

- **The "is it live" sentence** in `context/AGENTS.md` § Git now says "never with the steps `/preflight` itself performs listed as the operator's work", so it no longer clashes with steps only the operator can take (a place's `sh-agent update`, a root step without its standing-permission record).
- **Behavioural runs** for context sentence 1 (a low finding whose fix adds a gate goes to the operator): new 2 of 2 put the choice to the operator; the control's 2 of 2 rebutted alone. Sentence 6 as scoped: 2 of 2 offer `/preflight`. Log: `testing/2026-09-29-verify-retro-probe2.md`.

Sentence 8 (stop a finished agent with background work) has no behavioural run; it needs a multi-agent session. That gap is accepted and named in the release entry.

## Type
fix

## Files changed
`context/AGENTS.md` (one clause); the log above.

## Agents dispatched
None for the review: the diff is one clause, read by the main agent. `verifier` over the result.

## Verification
- The context size test passes at 21,444 bytes.
- The probes above.
