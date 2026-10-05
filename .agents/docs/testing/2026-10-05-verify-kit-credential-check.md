---
owner: verifier
date: 2026-10-05
---
# Verification — kit-credential-check   (2026-10-05)

Range reviewed: a372e95..ee28fe3

## Status: PASS

Scope: steps 1-6 as dispatched all ran; guide `.agents/docs/testing/guide.md` exists (changed in range, read-only for me). Not complete: behavioural run for the context/rule change not done (see Not examined).

## Commands
- 1. orient: pwd /home/ostling-simon/projects/styrman-harness; git status clean at start; HEAD ee28fe3
- 1. `gofmt -l cmd internal` -> exit 0, printed nothing
- 2. `go vet ./...` -> exit 0
- 3. `go test -race -count=1 -timeout 30m ./...` -> exit 0 (output .agents/scratch/verify-ee28/go-test-1.txt, run without -v per dispatch, so per-test counts: not printed; package ok lines): root 2.0s, cmd/sh-agent 673.1s, internal/dream 8.0s, internal/harness 132.2s, internal/memory 1.9s, internal/search 19.8s. vs previous: not compared (no -v counts).
- 4. `go test -count=1 -run 'TestKitPreamble|TestRenderedContext|TestSkillListingBudget' -v ./internal/harness/` -> exit 0; counts: pass 10, fail 0, skip 0 at ee28fe3 (5 top-level + 5 subtests of TestKitPreambleRefusesLeftoverCredentials); vs previous: none. Includes TestRenderedContextFitsAntigravityBlock, TestRenderedContextUnderTarget, TestSkillListingBudget, TestKitPreambleAssertsCommit.
- 5. `bash -n .agents/docs/testing/kit-preamble.sh` -> exit 0
- 6. `go run ./cmd/sh-agent install --home .agents/scratch/verify-ee28-home` -> exit 0, "1790 written, 0 unchanged, 0 removed, 0 kept"
- 6. `go run ./cmd/sh-agent status --home .agents/scratch/verify-ee28-home` -> exit 0; `binary ee28fe33055a+dirty` (no warning), `drift none — 1789 files match the repo`, 36 skills, 8 agents, 1 rule. "+dirty" comes from untracked/modified trail files listed below, not source.
- 7. suppressions added in range: internal/harness/catalog_test.go t.Skip at diff lines ~14, 18, 28, each with a reason string on the same line (no bash/git; missing tool; temp dir inside a checkout). None ran in step 4 (skip 0). No nolint, .only, ts-ignore.

## Coverage of changed code
- .agents/docs/testing/kit-preamble.sh — covered by internal/harness/catalog_test.go TestKitPreambleRefusesLeftoverCredentials (5 subtests) and TestKitPreambleAssertsCommit; syntax checked by bash -n
- internal/harness/catalog_test.go — is the test
- context/AGENTS.md, rules/critical-thinking.md — covered structurally by size pins TestRenderedContextFitsAntigravityBlock / TestRenderedContextUnderTarget, install lint in cmd/sh-agent and internal/harness (full suite green); wording itself has no unit test
- .agents/docs/testing/guide.md, changes entry — docs, no test

## Not examined
- Behavioural run (guide § 4) owed for the context/AGENTS.md and rules/critical-thinking.md wording change: main agent's, not run.
- No must-fail control for the preamble check was run by me (the test's subtests are the refusals; I did not mutate a copy).
- Antigravity/Claude sandbox starts, real-home smoke (§ 5), govulncheck: not run.
- Step 3 had no -v, so pass/skip counts for the full suite are not captured.

## Tree movement (step 8)
Changed/new paths since start, all trail artifacts of other agents or me, no source/test/config moved: M .agents/docs/bugs/BUG-160-credentials-left-in-scratch-homes.md, M .agents/docs/bugs/INDEX.md, ?? .agents/docs/reviews/2026-10-05-review-kit-credential-check.md, ?? this log. Source, tests and config unchanged.

## Environment
Linux 6.18.28 amd64, go1.27.1, no Docker/Colima use. Background job pid 403015 ended (kill -0 fails).
