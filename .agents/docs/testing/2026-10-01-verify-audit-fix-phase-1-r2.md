---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-1-r2   (2026-10-01)

Range reviewed: 4764c5f..working tree (dirty); source under cmd/sh-agent, internal/{dream,harness,search} (17 .go files, +1151/-99)

## Status: PASS

Scope: guide § 1, § 3 (scratch-home install + drift), § 3e, govulncheck, named-test coverage; at file level, not per hunk. Guide: .agents/docs/testing/guide.md (exists). No real CLI run, real home untouched.

## Commands
- `gofmt -l cmd internal` -> exit 0, no output
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0; counts: pass 484, fail 0, skip 3 at 4764c5f+dirty for go test -race -v -count=1 ./... (log .agents/scratch/verify-audit-fix-phase-1-r2/go-test-1.txt); vs previous: .agents/docs/testing/2026-09-30-verify-audit-fix-phase-1.md, pass 478→484, skip 3→3
- § 3 `go run ./cmd/sh-agent install --home .agents/scratch/verify-p1r2` -> exit 0, "268 written ... 268 files managed"; `status --home ...` -> exit 0, "drift none — 267 files match the repo", binary 4764c5f82dc0+dirty
- § 3e `go test -race -count=1 ./internal/memory ./internal/dream ./internal/search ./cmd/sh-agent` -> exit 0, all four ok
- `govulncheck ./...` -> exit 0, "Your code is affected by 0 vulnerabilities"

## Coverage of changed code
Named tests, all PASS in the -v run: TestUninstallAndLinkedDirs, TestUninstallSkipsLink, TestUninstallSkipsFilesUnderLink, TestInstallResyncFailureKeepsManifest, TestTrailStatusArgs, TestSessionTextSkipsMetaLines, TestReportRedactShapes, TestPromptEscapesFence, TestTrailDateFromGit.
Every changed source file has its _test.go changed in the same diff: main.go/main_test.go+memory_test.go; collect.go/collect_test.go; dream.go, prompt.go, run.go/run_test.go; apply.go/harness_test.go; report.go/report_test.go; search/dream.go/dream_test.go; search.go/search_test.go. Hunk-by-hunk mapping not performed.
Suppressions added: three t.Skip, all with a reason (two "root ignores directory modes", one "measurement, not a gate: set SH_AGENT_CORPUS=..."), same three skips as before. No .only, nolint or ignore directives.

## Not examined
Per-hunk coverage mapping; mutation checks; real-CLI tiers (not asked); § 3d recall measurement (TestRecall skipped, env-gated).
Tree check: source/test/config paths unchanged during the run. Seen, written by other agents under .agents/docs: modified INDEX.md, deploy/remote.md, plans/2026-09-30-audit-fix.md, specs delivery/dream/search.md, testing/guide.md, README.md; new BUG-124..130 records and this log.

## Environment
Linux 6.12 (Debian 13), go1.27.1 linux/amd64, no Docker used.
