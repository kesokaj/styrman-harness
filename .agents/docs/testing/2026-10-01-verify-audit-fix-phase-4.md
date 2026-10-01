---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-4   (2026-10-01)

Range reviewed: 1430d53..working tree (dirty); Go hunks in derive.go, frontmatter.go, inventory.go, harness_test.go

## Status: FAIL (one dispatch criterion unmet: reviewer derived file has no `model: pro`; everything else passes)

Scope: guide § 1, § 2, § 3 ran (`.agents/docs/testing/guide.md` exists). § 4 CLI starts and behavioural runs are recorded in `testing/2026-10-01-verify-agents-behavioural.md`, not here. No real CLI run, real home untouched.

## Commands
- `gofmt -l cmd internal` -> exit 0, printed nothing
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0, 6 packages ok; counts: pass 488, fail 0, skip 3 at 1430d53+dirty for go test -race -v -count=1 ./...; vs previous: testing/2026-10-01-verify-audit-fix-phase-3-r2.md, pass 485->488, skip 3->3 (TestAgyCorpusLeaves, TestRecall, TestHelperIndex, unchanged). Output .agents/scratch/verify-audit-fix-phase-4/go-test-1.txt
- `go run ./cmd/sh-agent list` -> exit 0; Skills (29), Subagents (8)
- `go run ./cmd/sh-agent render | wc -c` -> 14882 bytes
- `go test ./internal/harness -run 'TestRenderedContextFitsAntigravityBlock|TestRenderedContextUnderTarget' -v` -> both PASS
- `go test ./cmd/sh-agent -run TestHarnessMap` -> ok
- `go test . -run TestEmbedded` -> ok
- `sh-agent install --home .agents/scratch/verify-p4 --vendor all` -> exit 0, 270 written; `status` -> exit 0, "drift none — 269 files match the repo", 29 skills, 8 agents
- Derived `.gemini/config/agents/reviewer/agent.md`: names `config/skills/review-standards/SKILL.md` (1 hit); no `effort:`, `maxTurns:`, `skills:` key; **no `model:` line at all** (dispatch expected `model: pro`).

## Root cause
The `model: pro` expectation does not match the design: `agents/reviewer.md` has `model: sonnet`, and `modelMap` in derive.go (`opus->pro`, `haiku->flash`) omits unlisted aliases, which Antigravity reads as inherit. `security` (opus) derives `model: pro`, `changelog` (haiku) `model: flash`. Likely the dispatch's expectation is wrong, not the code; the caller decides. TestDeriveAgyRealReviewer passes (its assertions on model not confirmed here).

## Coverage of changed code
- derive.go (DeriveAgy hunk) — TestDeriveAgy (extended, PASS) and TestDeriveAgyRealReviewer (PASS)
- frontmatter.go / inventory.go (description limit and lint) — TestLintDescriptionLimit PASS, TestAgentDescriptionRules PASS
- No hunk without a covering test found. Added nolint/t.Skip/.only in Go diff: none.

## Not examined
§ 4 CLI starts and behavioural runs (see the behavioural log); govulncheck; markdown content of agents/skills beyond install and size checks. Tree check: status at end vs start differed only by other agents' paths: `.agents/docs/reviews/2026-10-01-review-1430d53-worktree-code.md`, `...-worktree-prompts-docs.md`, and this log; no source, test or config path changed.

## Environment
Linux 6.12 (Debian 13), Go from PATH (version not recorded), no Docker needed. Scratch home removed after the run.
