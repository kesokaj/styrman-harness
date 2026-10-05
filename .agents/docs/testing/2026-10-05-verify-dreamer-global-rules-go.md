---
owner: verifier
date: 2026-10-05
---
# Verification — dreamer-global-rules-go   (2026-10-05)

Range reviewed: 507244a..working tree (modified: agents/dreamer.md, skills/dream/SKILL.md, .agents/docs/specs/sh-agent/memory-loop.md; untracked docs-root: reviews/2026-10-05-review-worktree-docs-3.md, testing/2026-10-05-verify-dreamer-global-rules.md)

## Status: PASS

Scope: guide § 1, § 2, § 3 and the § 3 agent check ran; guide § 4 behavioural rows NOT run by me (main agent's log testing/2026-10-05-verify-dreamer-global-rules.md). No Go source changed; no verification command edited. Guide: .agents/docs/testing/guide.md (exists, not rewritten). Not "complete" because § 4 is out of my scope.

## Commands
- 1. `gofmt -l cmd internal` -> printed nothing
- 1. `go vet ./...` -> exit 0
- 4/5. `go test -race -timeout 30m -v -count=1 ./...` (1 CPU, run in background, polled) -> exit 0; all 6 packages ok (cmd/sh-agent 641 s, harness 98.9 s, search 16.5, dream 7.4, memory 1.8, root 1.6); counts: pass 553, fail 0, skip 3 at 507244a+dirty for go test -race -timeout 30m -v -count=1 ./...; vs previous: .agents/docs/testing/2026-10-05-verify-critical-thinking-r2-go.md, pass 553->553, skip 3->3. Output: .agents/scratch/verify-dreamer-global-rules-go/go-test-1.txt
- 2. `go run ./cmd/sh-agent list` -> exit 0; 36 skills, 8 subagents, 1 rule
- 2. `go run ./cmd/sh-agent render | wc -l` -> 87 (under 100)
- 2. `go test ./internal/harness/... -run 'TestRenderedContextFitsAntigravityBlock|TestRenderedContextUnderTarget' -v -count=1` -> PASS both; Antigravity block 16907 bytes (budget 22000); hand-written context 14700 bytes (target 15000)
- 3. `go run ./cmd/sh-agent install --home .agents/scratch/verify-dr-home` -> exit 0
- 3. `go run ./cmd/sh-agent status --home .agents/scratch/verify-dr-home` -> "drift none — 1789 files match the repo"; repo has 36 skills, 8 agents, 1 rules
- 3 agent check: agents/dreamer.md has 1 `## Records`; derived .gemini/config/agents/dreamer/agent.md has 1 `## Records`. Derived file carries the new text "and the global rule files (`~/.claude/rules/*.md`, or `~/.gemini/config/rules/*.md` beside that `AGENTS.md`)" at line 31 (source line 32).

## Coverage of changed code
- agents/dreamer.md, skills/dream/SKILL.md: markdown, covered by the embed/lint/size/listing-budget tests in `go test ./...` (all pass) and the scratch install drift check.
- memory-loop.md: docs only, no test applies.
- No Go source changed; no suppressions or skips added (diff is markdown only).

## Not examined
- Guide § 4 (sandbox start of each CLI, behavioural rows): main agent's log, not run by me.
- Guide § 5 real-home checks; govulncheck; 3a-3f subsections (no code in their scope changed).
- Content of the two untracked docs-root files was not reviewed.

## Environment
Linux 6.18.28, nproc 1, scratch home .agents/scratch/verify-dr-home. Tree check (step 8): the only status delta from step 1 is this log; no source/test/config path moved during the run.
