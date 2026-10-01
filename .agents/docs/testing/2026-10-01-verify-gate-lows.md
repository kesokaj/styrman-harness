---
owner: verifier
date: 2026-10-01
---
# Verification — gate-lows   (2026-10-01)

Range reviewed: 8b6c6de..working tree

## Status: PASS

Scope: complete for the dispatched steps 1-8; guide .agents/docs/testing/guide.md exists. Output files in .agents/scratch/verify-gate-lows/.

## Commands
- 1. `gofmt -l cmd internal` -> exit 0, printed nothing
- 2. `go vet ./...` -> exit 0
- 3. `go test -race -count=1 ./...` -> exit 0; all six packages ok (root 1.1s, cmd/sh-agent 15.5s, dream 6.8s, harness 3.7s, memory 1.2s, search 7.7s)
- 3. `go test -race -v -count=1 ./...` -> exit 0; counts: pass 524 (subtests included), fail 0, skip 3 at 8b6c6de+dirty; top-level `^--- PASS` = 347 (dispatch expected 346: one more than expected, no fail; the baseline 342 was not re-measured here); vs previous: none (skips TestAgyCorpusLeaves, TestRecall, TestHelperIndex are pre-existing, not in the diff)
- 4. new tests, `-v`: `--- PASS: TestStoreRunsNoProjectHook (0.01s)`, `--- PASS: TestSymlinkedIndexRefused (0.01s)`, `--- PASS: TestFeedbackRedactsEscapedSecret (0.02s)`, `--- PASS: TestLintRefusesAllowedTools (0.01s)`. No SKIP in the hook test, so its control fired.
- 5. must-fail controls (each compiled; no "build failed"; each restored from cp copy, cmp identical, `git diff HEAD --stat` identical to start):
  - (a) dropped `"-c", "core.fsmonitor=false",` -> `--- FAIL: TestStoreTrackedArgs`, `--- FAIL: TestStoreRunsNoProjectHook`
  - (b) removed Lstat symlink block in readLines -> `--- FAIL: TestSymlinkedIndexRefused`
  - (c) `feedback` text -> `sc.Text()` -> `--- FAIL: TestFeedbackRedactsEscapedSecret`
  - (d) removed allowed-tools refusal block -> `--- FAIL: TestLintRefusesAllowedTools`
- 6. `go test ./internal/harness -run 'TestRenderedContext|TestSkillLinksResolve' -v` -> PASS x3: Antigravity block 15137 bytes (budget 22000); hand-written context 14252 bytes (target 15000); TestSkillLinksResolve PASS
- 7. `go run ./cmd/sh-agent install --home .agents/scratch/verify-lows-home` -> exit 0, 266 files written. `status`: `binary 8b6c6de8dceb+dirty` (no warning), `repo has 34 skills, 8 agents`, `drift none — 265 files match the repo`. `list --home ... >/dev/null && echo ok` -> ok (every skill passes the new lint).
- 7. coverage/suppressions: each changed Go source has a covering test in the same package (memory.go -> memory_test.go, collect.go -> collect_test.go, inventory.go -> harness_test.go); report.go change covered by the internal/harness package tests (no dedicated new test named in the dispatch; not separately mutated). Added-suppression scan of the diff not run line by line (see Not examined).
- 8. `git status --short --untracked-files=all` end vs start: source/test/config paths unchanged (the 12 modified files plus the prior probe log). New paths, other agents' artifacts, listed as seen: M .agents/docs/bugs/BUG-143-skill-allowed-tools-and-curl-sh.md, M .agents/docs/bugs/INDEX.md, ?? .agents/docs/bugs/BUG-145-test-needle-doubled-escape.md, ?? .agents/docs/reviews/2026-10-01-review-worktree-gate-lows-docs.md, ?? .agents/docs/reviews/2026-10-01-security-worktree-gate-lows.md, ?? this log.

## Coverage of changed code
- internal/memory/memory.go — TestStoreRunsNoProjectHook, TestStoreTrackedArgs, TestSymlinkedIndexRefused (mutation-verified)
- internal/dream/collect.go — TestFeedbackRedactsEscapedSecret (mutation-verified)
- internal/harness/inventory.go — TestLintRefusesAllowedTools (mutation-verified)
- internal/harness/report.go — package tests only; no mutation control run (not requested)
- Docs/skills (README, data-gcp, docker, feature, makefile-patterns) — covered only by the link/lint/size pins in step 6 and `list`.

## Not examined
- report.go change not mutation-tested; added-suppression grep (nolint/skip/only) over the diff not run; no behavioural agent runs, no real-home install, govulncheck not run.

## Environment
Linux 6.12.107 (Debian), Go via PATH; no Docker involved. HEAD 8b6c6de + dirty.
