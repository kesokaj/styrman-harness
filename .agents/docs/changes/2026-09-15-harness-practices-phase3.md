---
owner: main
date: 2026-09-15
---
# Harness practices Phase 3: review loop and review standards (G10, G11, G15, G16 skill half)     (2026-09-15)

## Summary

Phase 3 of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`: edits I-12a–f, I-13a–b, I-14a–c and I-15a–f, applied verbatim from the plan's old/new blocks. A `--fix` review round keeps a ledger line per finding and re-dispatches with only the fixed lines. A rebuttal needs one checkable reason. The review standards gain checks for debug leftovers, partial writes and races, missing generated counterparts, changed assertions, unexplained suppressions and stale describing text. A suppression names its linter and its reason.

## Type

task

## Files changed

- `skills/review/SKILL.md` — I-12a–d, I-12f and I-13a.
  - I-13a: a rebuttal gives one checkable reason — not a defect because X, with X at file:line; a trade-off recorded at file:line, in an ADR or in an `AGENTS.md` section; or the user's words, quoted. It is never an intent the next reader cannot verify, and never a command that exited non-zero. A finding outside the scope is `left`, not `rebutted`.
  - I-12a–d and I-12f: at most three rounds, unless the user picks one more at the cap. Each round writes a `## Round N` section with one line per finding so far: `fixed at <file:line>`, `rebutted: <reason>`, `in flight: <who>` or `left: <owner>`. The re-dispatch carries each brief, then `---` and that section's fixed lines only, asking the agent to verify each fix and look for defects the fixes introduced; rebuttal reasons stay out of the dispatch. `--docs` findings in the `documentation` agent's files go to it, and the main agent fixes the rest (skills, agents, rules, the context, other `AGENTS.md` sections). The loop stops when a round reports nothing beyond the ledger's `rebutted` and `left` lines and no line is still `in flight`, after the third round, or when a round's findings land on the mechanism the previous round fixed or a fix introduced that round's new defect (sweep F78).
- `context/AGENTS.md` — I-12e: Definition of Done step 2's docs findings go to `documentation` in step 3, except findings in skills, agents, rules, the context or other `AGENTS.md` sections, which the main agent fixes. The rendered context stays at 93 lines.
- `skills/preflight/SKILL.md` — I-13b: a `rebutted` finding in the gate report carries a reason as `/review` Step 4 defines it.
- `skills/review-standards/SKILL.md` — I-14a–c: checklist items for leftover debug output and temporary scaffolding; a write that can fail halfway, a check-then-act race, a fallback that masks the error it replaces, and empty, NaN, Unicode and timezone inputs; a source change whose generated counterpart is absent from the diff; an assertion changed in an existing test that does not trace to a behaviour change, and an added suppression, skip, `.only` or lint-config exclusion with no named rule and reason; comments, doc-comments, test names and spec rows that no longer describe the changed behaviour.
- `skills/review-standards/references/go-review.md` — I-15a: `t.Skip` only behind a named prerequisite with a message; `//nolint:<linter> // <reason>` and `//lint:ignore <check> <reason>`, never bare.
- `skills/review-standards/references/typescript-review.md` — I-15b: no `it.only` or `describe.only`, and `.skip` and `eslint-disable` carry a reason. I-15e: `@ts-expect-error` needs a reason rather than a linked issue (sweep F95).
- `skills/golang/SKILL.md` — I-15c: a vet, lint or test failure is fixed in the source; a `//nolint` or a skip is for a demonstrated false positive and names the linter and the reason.
- `skills/golang/references/golang-ecosystem.md` — I-15d: `nolintlint` joins the starting linter set. I-15f: a project that runs golangci-lint keeps `.golangci.yml` as the source of truth, and `go vet` and `staticcheck` are the baseline every Go project passes (sweep F49).
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — boxes 3.1, 3.2 and the Phase 3 checkpoint ticked with evidence.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-phase3.md` — the Phase 3 verifier log.

## Agents dispatched

verifier. This entry was written by the main agent from the phase diff.

Review: no `reviewer` ran over the Phase 3 diff. The plan's Done list runs `reviewer` and a `--docs` review once over `c192622..HEAD`.

## Verification

The verifier log reads Status: PASS. The main agent also ran, each bare with its exit status read: `go test -race ./...` exit 0; `gofmt -l .` printed nothing; `go vet ./...` exit 0; a scratch install into `.agents/scratch/phase3-checkpoint-home` exit 0, then `status` with drift none; `render` 93 lines. The checks of boxes 3.1 and 3.2 met their expects; over `c192622` every 3.1 pattern printed 0, and "without a linked issue" went from 1 to 0.

Not run: the sandbox start of each CLI (testing guide §4), deferred to the plan's Phase 5 checkpoint (D36).

## Docs

The changed living docs are the skill and context texts above. `README.md`, the testing guide and the deploy guides are unchanged in this phase.

## Requirements

n/a
