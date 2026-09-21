---
owner: main
date: 2026-09-21
---
# Context consolidation, Phase 1 — under the cap: 27,163 → 20,012 bytes, pinned by a test     (2026-09-21)

## Summary

Second phase of `.agents/docs/plans/2026-09-21-context-consolidation.md` (BUG-085). The always-on context is pruned from 27,163 to 20,012 bytes, under Antigravity's 24,000-byte block with the 22,000-byte budget ADR-035 set, and a test now fails the build when it grows past that. Every rule stayed; what went was bug-number stories, a layout tree the `init-project` skill owns, and copies of what skills already say. Three procedures moved into skills. Both CLIs, started from a fresh sandbox home with this harness, answered every context question including the canary, and Antigravity's own conversation store shows no truncation marker for the first time since v1.2.0.

## Type

feature phase (context, skills, one test, docs)

## What changed

- `context/AGENTS.md` (74 → 66 lines, 23,324 → 16,569 bytes) and `context/dev-setup.md` (2,947 → 2,551 bytes): 44 actions from `plans/2026-09-21-context-pruning-review.md` — 34 lines shortened, 7 deleted, 3 moved — classified by three independent readers (Workflow `wf_ebd17196-757`), merged by majority, approved by the operator ("go all"). Discrete instructions 214 → 188 by the evidence reader's count.
- `skills/review-standards/SKILL.md` § 2: the must-fail-mutation protocol (BUG-068/082), the generated-artifact both-directions rule (BUG-066) and, from the docs review, the self-hosting-suite test-design rule. `skills/golang/SKILL.md` § Error Handling: the `err == nil && somethingToLose` guard example (BUG-069). `skills/bug/SKILL.md` Step 3 and `skills/feature/SKILL.md` Step 4 cite the protocol's new home.
- `internal/harness/size_test.go`: `TestRenderedContextFitsAntigravityBlock` renders the real checkout plus every `rules/*.md` and fails above 22,000 bytes (red first at 27,163; green at 20,012).
- `.agents/docs/bugs/INDEX.md`: Graduated rows BUG-066/068/069/082 point at the rules' new homes (`bugs`).
- Plan: Phase 1 boxes and checkpoint ticked; D16 (the saved-diff-file allowance dropped on purpose).

## Verification

- `testing/2026-09-21-context-consolidation-phase-1.md` (`verifier`) — PASS, 16 commands: gofmt, vet, `go test -race -count=1 ./...`; size test 20,012 / 22,000; `sh-agent list` 27 skills, 7 agents; scratch install and status drift none; canary last, both vendors identical; BUG-057 escape present once.
- `testing/2026-09-21-verify-sandbox-pruned-context.md` (main agent) — PASS: a logged-in sandbox home in scratch, both CLIs headless from a demo clone: six context questions right on both (canary GANNET, the new Bug bullet, no BUG-082 in the instrumentation bullet), `golang` loaded with no edits, `verifier` dispatched with a PASS log. Caveat recorded: Antigravity's checks 2 and 3 ran in the Claude clone because parallel shell calls share a working directory.
- Not exercised: the real home (Phase 4), the behavioural runs (Phase 2), the must-fail mutation (Phase 4).

## Review

`reviews/2026-09-21-review-worktree-merged.md`: code lens APPROVE, no findings; docs lens two findings — the self-hosting-suite rule fixed into review-standards § 2, the saved-diff allowance rebutted (never enforced; the rule is now stricter, D16).

## Agents dispatched

`reviewer` ×2 (code and docs lenses), `bugs` (index rows), `verifier`; the classification ran as a Workflow of three general subagents, read-only.

## Not done

The sandbox home keeps copied credential files until Phase 2's behavioural runs are over; they are deleted then and the handoff says so.
