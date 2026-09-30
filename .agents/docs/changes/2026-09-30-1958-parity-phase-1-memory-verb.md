---
owner: main
date: 2026-09-30
---
# Phase 1 of the parity plan: `sh-agent memory add|show`     (2026-09-30)

## Summary
The first `sh-agent` feature built under R10 (plan `plans/2026-09-30-parity-memory-dream.md`, Phase 1). It adds `sh-agent memory add --type T --name SLUG --description TEXT`, with the body on stdin, and `sh-agent memory show`. Both work in the checkout holding the working directory and use the main worktree's store (D26).
- `add` writes `<root>/.agents/memory/<slug>.md` (0600, in a 0700 store) in Claude Code's memory shape, without `originSessionId`. It sets that slug's line in `MEMORY.md`, found by its link target, so native lines titled otherwise survive.
- The index is read and written under a PID lock. A dead or 30-second-old lock is replaced; a live one is retried for 5 s, then refused.
- Refusals exit 2 and write nothing:
  - not a checkout;
  - `.agents` or the store is a symlink;
  - git tracks the store;
  - the path is not ignored (`git check-ignore`);
  - the index would pass 200 lines;
  - a bad slug, type or description, or an empty body.
- `show` prints the index, capped at 200 lines and 25 KB. It adds the rollout hint when Claude Code's native memory is still on in the project (D20, D27).

The skills and context do not call it yet; that is Phase 2. `harness.alive` is exported as `harness.Alive` for the lock.

## Type
feature (phase 1 of 6)

## Files changed
- `internal/memory/memory.go` and `memory_test.go` (new)
- `cmd/sh-agent/main.go`: the verb, its flags in `takesValue` (BUG-002's shape: the test went red on it first), the usage and the doc comment
- `cmd/sh-agent/memory_test.go` (new)
- `internal/harness/report.go` and `zip_test.go` (`Alive`)
- `README.md` verbs table
- the plan's Phase 1 boxes

## Agents dispatched
verifier (`testing/2026-09-30-verify-parity-memory-dream-phase-1.md`, PASS). The reviewer runs once over the whole feature at Done.

## Verification
`go test -race -count=1 ./...` exit 0. Two must-fail mutations, each compiled with a named test failing and each restored byte-identical: mine on the tracked-store check, the verifier's on the check-ignore check. A fixture run outside this repo (`~/sandbox/p1fixture`). Not exercised: either CLI calling the verb (Phase 2).

## Requirements
R10 (ADR-047); no requirement text changed in this phase.
