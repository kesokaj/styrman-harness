---
owner: verifier
date: 2026-10-01
---
# Verification - audit-fix-phase-4-r2 (2026-10-01)

Range reviewed: 1430d53..working tree

## Status: PASS

## Commands
- gofmt -l cmd internal -> exit 0, prints nothing
- go vet ./... -> exit 0
- go test -race -v -count=1 ./... -> exit 0; counts: pass 488, fail 0, skip 3 at 1430d53+dirty for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-10-01-verify-audit-fix-phase-4.md, pass 488->488
- sh-agent list / status -> 29 skills, 8 agents; render | wc -c -> 14882 (pins: 22,000 and 15,000; TestRenderedContext* 2 pass)
- TestHarnessMap 1 pass; go test . -run TestEmbedded -count=1 -> ok (2 pass in full run)
- install --home .agents/scratch/verify-p4r2 --vendor all -> exit 0; status drift none (269 files match); scratch home removed
- derived reviewer agent.md: names config/skills/review-standards/SKILL.md (1 hit), no model/effort/maxTurns/skills line; security agent.md: model: pro
- Coverage tests PASS: TestAgentDescriptionRules, TestLintDescriptionLimit, TestDeriveAgyRealReviewer, TestDeriveAgy (2)

Scope: guide 1, 2, 3 ran. Guide 4 and behavioural runs are in testing/2026-10-01-verify-agents-behavioural.md. Changed Go hunks (derive.go, frontmatter.go, inventory.go) are covered by harness_test.go tests above; per-hunk cover profile not run.
## Not examined
-cover per hunk, govulncheck, any real CLI, real home.
## Environment
Linux 6.12, Go as on PATH, no Docker needed. Tree unchanged during run (status diff empty).
