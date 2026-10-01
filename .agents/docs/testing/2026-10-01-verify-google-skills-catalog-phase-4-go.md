---
owner: verifier
date: 2026-10-01
---
# Verification — google-skills-catalog-phase-4-go   (2026-10-01)

Range reviewed: bfd25d3..working tree (Phase 4, Go and script changes)

## Status: PASS

Scope: steps 1-9 of the dispatch; guide at .agents/docs/testing/guide.md exists.

## Commands
- 1. `gofmt -l cmd internal` → exit 0, printed nothing.
- 2. `go vet ./...` → exit 0.
- 3. `go test -race -count=1 ./...` → exit 0; all 6 packages ok (root, cmd/sh-agent, dream, harness, memory, search). `go test -count=1 -v ./...` (output .agents/scratch/verify-p4go/v.txt) → exit 0; counts: pass 357 (top-level `^--- PASS`; 534 incl. subtests), fail 0, skip 3 at bfd25d3+dirty for go test -count=1 -v ./...; vs previous: none run (dispatch expects 357, matches). Skips are env-gated, pre-existing: TestAgyCorpusLeaves, TestRecall, TestHelperIndex.
- 4. TestNoIndexMarkerSkipsTheDirectory PASS (internal/search). harness: PASS TestCatalogMatchesChecksums, TestCatalogIndexMatchesEntries, TestCatalogFilesAreCarryable, TestCatalogInstallsClean, TestRouterDescription, TestRefreshCatalogScript, TestSkillListingBudget; no SKIP.
- 5. `go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles -v .` → PASS.
- 6. Must-fail controls (each from cp copy, restored, cmp identical): (a) noIndex returns false -> TestNoIndexMarkerSkipsTheDirectory FAIL (compiled, leaked docs shown); (b) marker removed from catalogHarnessFiles -> TestCatalogMatchesChecksums FAIL ".sh-agent-noindex is in the catalog and not in CHECKSUMS"; (c) marker line deleted from script -> TestRefreshCatalogScript FAIL "the search marker was not written". `git diff --stat` after equals start value.
- 7. `SH_AGENT_RECALL=1 go test ./internal/search -run TestRecall -v` -> "index mean all 0.633 over 41". `sh-agent search bigquery reservation slots | grep -c skills/google-catalog/catalog/` -> 0. Control `search google-catalog router` -> 10 hits, outside catalog (reviews files).
- 8. `sh-agent install --home .agents/scratch/verify-catalog-p4-home` exit 0; status: 35 skills, 8 agents, drift none (1837 files match); marker count in deployed catalog dir = 1.
- 9. `git status --short | grep -v security-catalog-` equals start (the start already listed the Phase 4 log; this log is new and untracked, plus nothing else changed).

## Coverage of changed code
- internal/search/reconcile.go: search_test.go TestNoIndexMarkerSkipsTheDirectory
- refresh-catalog.sh and catalog marker: catalog_test.go TestRefreshCatalogScript, TestCatalogMatchesChecksums
- No added suppressions or skips reviewed beyond the diff (not scanned in detail).

## Not examined
Added-suppression grep over the diff not run (step 7 suppression scan); docs diffs not read; Antigravity runs not done.

## Environment
Linux 6.12 amd64, Go per local SDK, no Docker used.
