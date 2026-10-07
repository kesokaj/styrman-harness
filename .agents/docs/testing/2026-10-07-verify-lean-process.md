---
owner: verifier
date: 2026-10-07
---
# Verification — lean-process   (2026-10-07)

Range reviewed: working tree (uncommitted) on 4793d5a

## Status: PASS

Scope: guide §1 commands and scratch-home install ran on the uncommitted tree; Antigravity/behavioural runs (guide §2/§4) NOT run, so not complete.

## Commands
- 1. pwd /home/ostling-simon/projects/styrman-harness; git status: 22 entries (17 modified, 5 untracked), matches the dispatch
- 4. `gofmt -l cmd internal` -> exit 0, printed nothing
- 4. `go vet ./...` -> exit 0
- 4. `go test -race -timeout 30m -v -count=1 ./...` -> exit 0 (all 6 packages ok: root 9.5s, cmd/sh-agent 694s, internal/dream 20s, internal/harness 406s, internal/memory 12s, internal/search 81s); counts: pass 559, fail 0, skip 3 at 4793d5a+dirty for go test -race -timeout 30m -v -count=1 ./...; vs previous: none (no earlier log records counts in this form). Output: .agents/scratch/verify-lean-process/go-test-1.txt
- 4. `go test ./internal/harness -run TestRenderedContext -v -count=1` -> exit 0 (verification-only narrowing, for the pin numbers): TestRenderedContextFitsAntigravityBlock 17128 bytes (context 15672 + rules 1456) vs budget 22000; TestRenderedContextUnderTarget hand-written context 14719 bytes vs target 15000 (goal 14,500: 219 over the goal, 281 under the ceiling)
- 6. `go run ./cmd/sh-agent install --home .agents/scratch/home-verify-lean` -> exit 0, 1790 written
- 6. `go run ./cmd/sh-agent list` -> exit 0, 63 lines, no parse error
- 6. installed .agents/scratch/home-verify-lean/.claude/AGENTS.md contains "reached a commit, a deploy or the user" (1 match); installed .claude/skills/review/SKILL.md contains "<scope>-code-docs" (1 match)

## Coverage of changed code
No Go source changed. Changed files are prose, covered only by generic tests in `go test ./...`: lint of frontmatter/descriptions/permission keys (TestAgentDescriptionRules, TestLintDescriptionLimit, TestLintRefusesAllowedTools, TestLintPermissionKeysByKind and kin), TestSkillListingBudget, the two size pins (context/AGENTS.md), and the install itself. The init-project asset templates are covered by install/scaffold tests only as copied content. NOT covered: the semantics of the new process wording (no test asserts it beyond the two grep checks above), behavioural runs of changed skills/agents on Claude Code or Antigravity (guide §2/§4), the real-home smoke (§5).

## Suppressions
Diff is prose only; no suppression markers examined in Go/TS (none expected). Not grepped line by line.

## Not examined
Antigravity sandbox start and behavioural runs per changed skill/agent; real-home smoke; govulncheck; `/review --docs`; trail files under .agents/docs (read as data only by tests that scan them, not reviewed). Previous-count comparison.

## Tree check
Step 8: status after the runs equals status at start, plus this log only. No source/test/config path moved.

## Environment
Linux 6.18.28, Go via local toolchain (version not recorded), no Docker used.
