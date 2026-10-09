---
owner: verifier
date: 2026-10-09
---
# Verification — metrics-store   (2026-10-09)

Range reviewed: 2c0f215..166aa39 plus the uncommitted working tree (review fixes; 20 modified, 5 untracked at start)

## Status: PASS

Scope: gofmt, vet, full `go test -race ./...` and guide § 3g tier ran; govulncheck not run (not installed, per dispatch); sandbox runs not re-run (cited). Guide: .agents/docs/testing/guide.md (exists, filled). Not `complete` only because of govulncheck and the cited-not-rerun sandbox evidence.

## Commands
- 1. `pwd` = /home/ostling-simon/projects/styrman-harness; change set present; HEAD 166aa39.
- 4. `gofmt -l cmd internal` -> exit 0, no output.
- 4. `go vet ./...` -> exit 0.
- 4/5. `go test -race -v -count=1 -timeout 30m ./...` (background, polled; output .agents/scratch/verify-metrics-store/go-test-1.txt) -> exit 0; counts: pass 610, fail 0, skip 3 at 166aa39+dirty for that command; vs previous: .agents/docs/testing/ (v2.5.0 release log, pass 559, skip 3) pass 559->610, skip 3->3 (explained: metrics/trail/dream tests added in range). Skips: TestAgyCorpusLeaves, TestRecall, TestHelperIndex (measurement tests, env-gated, unchanged). All 8 packages ok (cmd/sh-agent 715s, harness 557s). Covers size pins and TestHarnessMap (part of ./... in internal/harness and root; both ok).
- 4/5. § 3g tier `go test -race -timeout 30m ./internal/metrics ./internal/trail ./internal/dream ./internal/search ./cmd/sh-agent` (run with -v -count=1 for counts; go-test-2.txt) -> exit 0; counts: pass 274, fail 0, skip 3 at 166aa39+dirty; vs previous: none with this command. All 5 packages ok.
- govulncheck ./... -> not run (not installed; dispatch).
- 7. Suppressions: added-line grep of `git diff 2c0f215` for nolint, lint:ignore, t.Skip, go:build ignore, .skip(, .only( -> no hits.
- 8. Tree: porcelain after runs equals the start except this log. No source/test/config path moved.

## Coverage of changed code
Changed non-test .go files (git diff --name-only 2c0f215) and covering tests:
- internal/metrics/{agy,claude,events,ingest,numbers,report,show,store,verb}.go — agy_test, claude_test, events_test, numbers_test, report_test, show_test, store_test, verb_test; ingest.go via claude_test.go Ingest( calls.
- internal/trail/trail.go — internal/trail/trail_test.go (untracked, new).
- internal/dream/{collect,metrics,run}.go — collect_test, metrics_test, run_test; dream.go, prompt.go (small edits, 6 and 4 lines) — covered at package level by internal/dream tests, no dedicated file.
- internal/search/{export,dream,store}.go — export.go by protofields_test.go (ProtoFields) and dream_test.go (AgyTime/OpenCopy/KindOf users); dream.go, store.go (2-3 line edits) by dream_test/search_test at package level.
- internal/harness/report.go — report_test.go, harness_test.go.
- cmd/sh-agent/main.go — cmd/sh-agent/metrics_test.go.
No changed source file lacks a covering test file in its package. Line-level coverage was not measured (no `-cover` run).

## Not examined
- govulncheck (not installed).
- Sandbox runs on Claude Code and Antigravity: cited from .agents/docs/testing/2026-10-08-verify-metrics-sandbox.md, not re-run; that log records a FAIL row (Claude Code `retro` matrix check 0 of 2) which is an Open item, not part of this tree's tests. This PASS does not cover it.
- Line-level coverage (`go test -cover`); dream.go, prompt.go, search dream.go/store.go edits mapped at package level only.
- Env-gated measurements TestAgyCorpusLeaves, TestRecall (skipped).
- Docs/review content (reviews and bug records in the tree not read).

## Environment
Linux 6.18, go1.27.1 linux/amd64, one box; no Docker involved.
