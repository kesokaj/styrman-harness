---
owner: main
date: 2026-10-08
---
# Handoff — 2026-10-08 — v2.5.0 follow-ups closed: the preamble's leftover check, hard seeds, the broker note

Session ID: 728a4838-8b3c-4819-a62d-2b01e563a77f
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

The operator said "go" to the v2.5.0 handoff's follow-ups 2.1–2.3. Plan: `plans/2026-10-08-kit-leftovers-and-hard-seeds.md`, closed.

## Completed

- **`bcc3c4c`, the leftover check in `kit-preamble.sh`.** Every kit gets it:
  - it runs first in the exit trap and matches by cwd or executable under the clone or `KIT_HOME`, or by `HOME=$KIT_HOME`
  - roots are refused at load when they are `/`, the home, the checkout, or a parent of either
  - a new exit line, `kit: no processes left by the run`
  - pinned by `TestKitPreambleKillsLeftoverProcesses`, with mutations seen failing
- **`seed-security` hard mode,** three multi-step, fail-open flaws: `high` 3/3 in 3/3 runs at about 26% of the time of `max`, which also scored 3/3.
- **The `remote.md` broker note.**
- **Gates:** a reviewer and a security review, 2 rounds each by resume; the verifier PASSed (560/0/3).

## Open / In Progress

- None for the agent.
- Requirements edits awaiting approval: none.

## Key Decisions

- The match is by cwd, executable or HOME, not by new listeners (D1), because a listener diff would also catch the operator's own servers.
- The hard flaws are scored apart from the six (D2/D3).

## Corrections

- None from the operator.

## Agent notes (unverified)

- **A mutation went wrong.** The run with the checkout-refusal guard removed made the trap treat this checkout as `KIT_HOME`, around 08:50 UTC on 2026-10-08:
  - it wrote `bin/sh-agent`, since deleted
  - it began a credential scan of the checkout, which was stopped
  - it may have sent SIGTERM to same-uid processes working inside the checkout

  Lesson: a mutation that removes a safety guard must run against a copy of the guarded file in a fixture tree, never the real one. This is a candidate rule for `review-standards` § 2, not yet proposed.
- **The size pins pass, but the context has only 64 bytes of headroom.**

## Environment State

- Branch: main, clean, pushed; newest tag v2.5.0. Nothing is released since; the preamble changes are test kits, not shipped content.
- Real home: a checkout place at `9c87ddb`. `sh-agent update` brings it to `bcc3c4c` (only the test files differ).

## Artifacts to Read

- `changes/2026-10-08-0901-kit-leftovers-preamble-hard-seeds.md`
- `testing/2026-10-08-verify-seed-security-hard.md`

## Next Steps

1. **Operator:** update the other places to v2.5.0 (carried from the last handoff).
2. **Operator choice:** add the "guard-removal mutations run on a fixture copy" rule to `review-standards` § 2. Recommended.

The next task should start in a fresh session from this handoff.
