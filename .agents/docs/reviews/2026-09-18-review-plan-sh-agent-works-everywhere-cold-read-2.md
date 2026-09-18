---
owner: reviewer
date: 2026-09-18
---
# Review — plan-sh-agent-works-everywhere-cold-read-2   (2026-09-18)

Range reviewed: plan file, round 2 — `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` and `internal/harness/apply.go`, `skills/review/SKILL.md`

## (1) Gaps

- **BLOCKING (new)** `apply.go:49-55`: the mode-correction check (`os.Stat` + `os.Chmod`) is wrapped entirely inside `if !dryRun { ... }`. D7 requires the `mode  ~/<rel> (0644 -> 0755)` line on `--dry-run` too ("`--dry-run` prints exactly the same lines and writes nothing", Interfaces line 25), which means the `os.Stat` mode probe must run unconditionally (only the `os.Chmod` stays gated). The plan never states this restructuring, and Phase 1's six named unit-test cases (line 47) do not include "mode line appears under `--dry-run` without chmod-ing" as one of them — a naive additive implementation (print new text but keep the existing `if !dryRun` gate around the whole block) silently leaves `--dry-run` mute on mode-only drift, breaking the plan's own "exactly the same lines" claim, and no listed test would catch it. Needs a decision/explicit test before Phase 1 starts.
- **Resolved from round 1**: Phase 2's `/review --project --all` mis-scoping is gone — the plan now reads `/review internal/harness cmd/sh-agent --all` (paths scope, both dirs read whole), which matches `skills/review/SKILL.md:25`'s "paths (any number)" scope class. No longer a gap.
- **Resolved from round 1**: the EACCES/unreadable-file ambiguity and the mode-only-branch silence ambiguity are now settled by D6 and D7, each with an explicit output format and an explicit unit-test case in Phase 1 task 1.
- **Minor, non-blocking**: the real-run (non-dry) removal path (`apply.go:87-90`) currently prints nothing at all — only the dry-run branch (`apply.go:84`) has a `remove` line. Phase 1 correctly implies this needs adding ("per-file lines on every run"), so it's ordinary implementation work, not a plan ambiguity — but the existing dry-run text carries a `(no longer in the repo)` suffix that the Interfaces line's `remove` example (line 25) doesn't show. Verification's "Both hold whatever the output lines are called" (line 40) makes exact wording non-binding, so this is not blocking.
- **Carried over from round 1, still present**: Phase 1's checkpoint (line 49) still asserts a literal "260 `add` lines" with no caveat that the count is illustrative/re-derived at checkpoint time. Low severity, unchanged from round 1.

## (2) Scope vs Why

- D5 now explicitly authorizes the Phase 2 audit ("the operator... 'All other functions of sh-agent need to work'... cold read round 1, scope note") — the round-1 concern about the audit lacking a D-row is resolved.
- No new extra endpoints/config found beyond round 1's read: "No new flags, no new verbs, no config keys" still holds against the diff described.

## (3) Readiness

1. **BLOCKING** — Does the `--dry-run` mode-line require restructuring the `os.Stat` probe outside the `if !dryRun` guard (apply.go:49), and is that covered by an explicit test? Not stated; not in Phase 1's six named cases.
2. **NON-BLOCKING** — Exact wording of the real-run `remove` line (with or without the `(no longer in the repo)` suffix already in the dry-run branch) — Verification line 40 says wording is not load-bearing; a reasonable implementer keeps the existing suffix.
3. **NON-BLOCKING** — The "260 add lines" checkpoint number will drift if `context/`/`skills/`/`agents/` change before the checkpoint runs; a reasonable implementer re-derives it at run time rather than treating it as a hard assertion.

## Verdict

`design ready` is close but not clean: one new BLOCKING gap (the `--dry-run` mode-probe restructuring, apply.go:49) surfaced by reading the code the plan's own D7 depends on. `implementation not ready` until that is decided and a test case for it is added to Phase 1 task 1's list.

## Not examined

Did not read `apply.go` beyond line 119 (the full file was covered, but not its tests `apply_test.go`), `main.go`, `plan.go`, `repo.go`, `source.go`, `release.go`, `report.go`, `adopt.go`, `state.go`, `render.go`, `inventory.go`, `derive.go`, `pathguard.go` or any of their tests; `README.md`, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/guide.md`; Phase 2's recon brief (not yet written); Phase 3-5's real-machine claims. Given the 25-tool-call budget, this pass re-checked the two files the round-1 report's BLOCKING findings pointed at (`apply.go`, `skills/review/SKILL.md`) plus the full plan text, to verify D5/D6/D7 actually close round 1's gaps and to look for anything D7's own target code raises that D7 doesn't settle.
