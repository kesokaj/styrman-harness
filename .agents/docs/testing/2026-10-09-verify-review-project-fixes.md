---
owner: verifier
date: 2026-10-09
---
# Verification — review-project-fixes   (2026-10-09)

Range reviewed: commits dbd776d..6cec23b (git log dbd776d..6cec23b); the diff is git diff dbd776d...6cec23b.

## Status: PASS

Scope: guide .agents/docs/testing/guide.md (exists, not rewritten); setup commands, deploy/remote.md verification command, coverage map and two must-fail controls ran; corpus measurement and CLI sandbox starts did not (see Not examined), so not `complete`.

## Commands
- 1. pwd /home/ostling-simon/projects/styrman-harness; `git status --short --untracked-files=all` clean at 6cec23b.
- 4. `gofmt -l cmd internal` -> exit 0, no output.
- 4. `go vet ./...` -> exit 0, no output.
- 4. `go test -count=1 -run TestDeliverySpecCurrent ./cmd/sh-agent` (quoted from deploy/remote.md Releases as added by the range) -> exit 0: ok 0.012s.
- 4. `go test -race -v -count=1 -timeout 30m ./...` (foreground-equivalent, output .agents/scratch/verify-rpf/go-test-1.txt) -> exit 0; ok: root 2.1s, cmd/sh-agent 716.6s, dream 46.2s, harness 562.1s, memory 18.9s, metrics 28.8s, search 111.2s, trail 15.1s. counts: pass 682, fail 0, skip 5 at 6cec23b for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-10-09-verify-project.md, pass 649->682, skip 3->5; explained: the two new skips are the helper processes TestHelperMain and TestHelperCall added by the range (skip when not re-executed); the other skips are TestAgyCorpusLeaves, TestRecall, TestHelperIndex (measurements/helpers).
- 4. `go run ./cmd/sh-agent install --home .agents/scratch/home` -> exit 0 ("9 written, 1781 unchanged, 0 removed, 0 kept - 1790 files managed"; PATH note only).
- 4. `go run ./cmd/sh-agent list` -> exit 0.
- 4. `"$(go env GOPATH)/bin/govulncheck" ./...` -> exit 0: "No vulnerabilities found... 0 in packages you import and 1 in modules you require" (the accepted GO-2026-5024, not reached); no new advisory.
- 7. Named tests, all PASS in the run above: TestLintPermissionKeysByKind, TestCallReturnsAtTheDeadlineDespiteAGrandchild, TestCallKeepsTheAnswerWhenAHelperOutlivesTheCLI, TestCallKillsTheGroupOnInterrupt, TestUsageErrorsExit2, TestVersionVerb, TestMainExitStatus, TestDeliverySpecCurrent, TestInstallPassArgsForwardsAnExplicitVendor, TestReportRedactShapes, TestDocsRoot, TestSkillsReadsBlockScalarDescription, TestFromZipEmptyValueRefused, TestFromZipDashValueRefused, TestUpdateFromZipRefusals, TestFromZipRefusedOnOtherCommands, TestSearchVerb, TestSkillListingBudget. TestHelperMain/TestHelperCall are SKIP by design when run directly; their parents ran them as child processes and passed.
- 7. Must-fail control 1 (copy under .agents/scratch/verify-rpf/copy, own git init; must-pass first: ok): removed "mcpServers" from the permissionKey call in internal/harness/inventory.go:261 -> TestLintPermissionKeysByKind FAIL ("agent with mcpServers ... loaded", both rows); restored from saved copy (cmp identical), re-run ok.
- 7. Must-fail control 2: broke the `credentials key` needle in internal/harness/report.go:167 (`credentials?` -> `credentialsXX?`) -> TestReportRedactShapes FAIL; restored (cmp identical), re-run ok. Copy deleted.
- 7. Suppressions: added `t.Skip` lines are the helper-process skips (usage_test.go:63, run_test.go:795, search_test.go:529) each with a reason on its line; no nolint/ts-ignore/.only seen in the Go diff (grep limited to the Go range, not to skills/docs).
- 8. `git status --short --untracked-files=all` after: only the expected log path (this file); no source, test or config path changed. .agents/scratch/home and verify-rpf are gitignored.

## Coverage of changed code (cmd/internal non-test .go, by reading the tests)
- cmd/sh-agent/main.go (usage exits, version, installPassArgs) - usage_test.go, zip_test.go, search_test.go, delivery_test.go
- internal/dream/run.go (process group kill, WaitDelay, lingering helper) - run_test.go three Call tests
- internal/dream/collect.go, dream.go (block scalar description, docsOf) - collect_test.go TestSkillsReadsBlockScalarDescription; docsOf covered indirectly by existing dream tests (no dedicated test for the docsOf default)
- internal/harness/inventory.go (mcpServers, loadDocs/unexport) - TestLintPermissionKeysByKind (mutation-proven)
- internal/harness/report.go - TestReportRedactShapes (mutation-proven on one needle)
- internal/harness/frontmatter.go, release.go, memory.go, search/export.go - comment/unexport/removal only; compile-checked, existing tests
- internal/search/search.go (usageError, IsUsage) - TestSearchVerb and TestUsageErrorsExit2 at the CLI level; no search-package unit test of IsUsage named
- internal/search/reconcile.go, metrics/claude.go, metrics/events.go (trail.DocsRoot) - TestDocsRoot for the resolver; the call sites are covered only through existing search/metrics tests (not read line by line)
- internal/metrics/events.go readSmall (io.ReadAll over LimitReader replaces a single Read): no test of a file over 1 MiB found; behavioural change on that path is NOT COVERED by a dedicated test (weak finding; the 1 MiB cap is untested)
- internal/trail/trail.go - trail_test.go TestDocsRoot

Coverage gap list (separate from command status): metrics readSmall 1 MiB cap untested; docsOf default and search IsUsage have only indirect coverage.

## Not examined
- Guide 3e corpus measurement (TestAgyCorpusLeaves): not run; the range does not touch the Antigravity reader (internal/search/dream.go unchanged) and ~/.gemini/antigravity/conversations is absent on this box.
- No CLI sandbox starts (not owed per dispatch); no -cover profile run; no line-level coverage tool used, coverage mapped by reading.
- Skill-file changes covered only via TestSkillListingBudget (PASS) and the scratch-home install (exit 0).
- No instruction-like text found in anything read.

## Environment
Linux 6.18.28 amd64, go1.27.2, no Docker involved. Previous run's go test took ~12 min; this one ~12 min.
