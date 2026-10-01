---
owner: verifier
date: 2026-10-01
---
# Verification — gate-lows-r2   (2026-10-01)

Range reviewed: 8b6c6de..working tree

## Status: PASS

Scope: complete for dispatched steps 1-7; guide .agents/docs/testing/guide.md exists. Previous log: .agents/docs/testing/2026-10-01-verify-gate-lows.md. Outputs in .agents/scratch/verify-lows-r2/.

## Commands
- 1. `gofmt -l cmd internal` -> exit 0, printed nothing
- 2. `go vet ./...` -> exit 0
- 3. `go test -race -v -count=1 ./...` -> exit 0; all six packages ok (root 1.2s, cmd/sh-agent 15.5s, dream 6.3s, harness 3.8s, memory 1.2s, search 7.4s); counts: pass 525 (subtests included), fail 0, skip 3 at 8b6c6de+dirty; top-level `^--- PASS` 348; vs previous: .agents/docs/testing/2026-10-01-verify-gate-lows.md, pass 524->525, skip 3->3, top-level 347->348 (explained: the new TestLintPermissionKeysByKind); skips TestAgyCorpusLeaves, TestRecall, TestHelperIndex unchanged.
- 4. harness named tests -> ok: PASS TestAgentDescriptionRules, TestLintDescriptionLimit, TestLintRefusesAllowedTools, TestLintPermissionKeysByKind, TestRenderedContextFitsAntigravityBlock (15137 bytes, budget 22000), TestRenderedContextUnderTarget (hand-written 14252 bytes, target 15000), TestSkillLinksResolve. memory: PASS TestStoreRunsNoProjectHook, TestSymlinkedIndexRefused.
- 5. Must-fail controls (each from a cp copy in scratch, compiled, restored; `git diff --stat` identical to start after each):
  - (a) `unreadLine` returns "", false unconditionally -> TestLintPermissionKeysByKind FAIL at harness_test.go:2428: "skill with double-quoted key loaded", "single-quoted key loaded", "spaced colon loaded".
  - (b) "hooks" dropped from skill permissionKey (inventory.go:206) -> FAIL: "skill with hooks loaded".
  - (c) "permissionMode" dropped from agent check (inventory.go:260) -> FAIL at harness_test.go:2449: agent with "permissionMode: bypassPermissions" loaded.
  - (d) Lstat refusal removed in memory.go readLines (lines 248-250) -> TestSymlinkedIndexRefused FAIL, three assertions fire: memory_test.go:395 "Show followed a symlinked index", :398 "Add followed a symlinked index", :403 "the index is no longer the link: Add wrote a regular MEMORY.md from the outside file" (the re-aimed third is among them).
  - Restored; `git diff --stat` equals the start; tests pass in the unmutated state (steps 3-4).
- 6. `go run ./cmd/sh-agent install --home .agents/scratch/verify-lows-home-r2` -> exit 0, 266 written; `status` -> exit 0: "repo has 34 skills, 8 agents", "drift none — 265 files match the repo", binary line without warning; `list` -> exit 0; installed dirs hold 34 skills, 8 agents. All pass the stricter lint.
- 7. `git status --short --untracked-files=all` at end equals start plus this log only.

## Coverage of changed code
Not mapped file by file this run (dispatch scoped to the named tests); changed Go files frontmatter.go, inventory.go (harness_test.go: TestLintPermissionKeysByKind), memory.go (memory_test.go: TestSymlinkedIndexRefused), dream/collect.go, harness/report.go have tests in their packages, all green.

## Not examined
Per-line coverage (`go test -cover`) not run; suppressions scan of the diff not run; docs/prose changes not read; Antigravity/Claude behavioural runs not done; govulncheck not run.

## Environment
Linux 6.12 (Debian 13), Go as on PATH, no Docker involved (no Dockerfile change).
