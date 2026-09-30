---
owner: verifier
date: 2026-09-30
---
# Verification — parity-memory-dream-phase-1   (2026-09-30)

Range reviewed: HEAD cb799ff + uncommitted tree (git diff HEAD; untracked internal/memory/memory.go, internal/memory/memory_test.go, cmd/sh-agent/memory_test.go)

## Status: PASS

Scope: complete per dispatch (gofmt, vet, race tests, scratch-home install/drift, coverage map, one mutation); guide exists (.agents/docs/testing/guide.md). Gaps below are named, not FAIL.

## Commands
- `gofmt -l cmd internal` → exit 0, printed nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0; output .agents/scratch/verify-p1/go-test-1.txt; counts: pass 404, fail 0, skip 2 at cb799ff+dirty for go test -race -v -count=1 ./...; vs previous: none with this exact command (earlier logs carry no comparable counts line). Skips: TestRecall, TestHelperIndex (internal/search, neither touched by this change set).
  - ok: root, cmd/sh-agent, internal/harness, internal/memory, internal/search
  - PASS TestMemoryVerb, PASS TestMemoryRefusalExitsTwo (cmd/sh-agent)
  - internal/memory: PASS TestAddWritesFileAndIndex, TestAddReplacesInPlace, TestAddRefusals (+11 subtests: index_full, empty_body, not_a_checkout, agents_is_a_symlink, store_is_a_symlink, store_not_ignored, bad_slug, bad_type, newline_description, empty_description, git_tracks_the_store), TestWorktreeSharesTheMainStore, TestStaleLockIsReplaced, TestLiveLockRefusesAfterRetry, TestShow, TestShowCaps, TestShowRefusesTrackedStore, TestShowRolloutHint (+6 subtests)
- Scratch-home install: `go build -o .agents/scratch/verify-p1/sh-agent ./cmd/sh-agent` exit 0; `sh-agent install --home .agents/scratch/verify-p1/home` exit 0; `sh-agent status --home ...` exit 0: "drift none — 265 files match the repo". Real home: only this session's own transcript changed; no install output there.
- Mutation (check-ignore refusal, memory.go): saved copy .agents/scratch/verify-p1/memory.go.orig; changed `check-ignore ... .Run() != nil` to `.Run() == nil && false` (1 match); `go build ./internal/memory` succeeded; `go test ./internal/memory -run TestAddRefusals` FAILED on TestAddRefusals/store_not_ignored; restored from copy, `cmp` identical, package test ok again.

## Coverage of changed code (internal/memory/memory.go refusals)
- not a checkout (Root) — TestAddRefusals/not_a_checkout (Add path only)
- symlinked .agents / symlinked store — TestAddRefusals/agents_is_a_symlink, store_is_a_symlink (Add path only)
- tracked store — TestAddRefusals/git_tracks_the_store, TestShowRefusesTrackedStore
- not ignored — TestAddRefusals/store_not_ignored (mutation-verified)
- index full — TestAddRefusals/index_full
- bad slug / type / description (empty, newline) / empty body — matching subtests
- live lock refusal — TestLiveLockRefusesAfterRetry; stale lock — TestStaleLockIsReplaced
- cmd/sh-agent/main.go memory verb — TestMemoryVerb, TestMemoryRefusalExitsTwo; internal/harness/report.go and zip_test.go changes covered by the passing internal/harness suite (diff not examined line by line)
- Gaps (branches with no test, all shared code reached via Add): Show on a non-checkout and Show with symlinked .agents/store (only Add tested; trusted() is shared); trusted() "cannot ask git" (ls-files error) branch; MkdirAll/WriteAtomic/readLines I/O error returns; nativeOn unparseable-JSON branch (verify: TestShowRolloutHint cases list does not name it).
- Suppressions added: none (grep of diff and new files for nolint, lint:ignore, t.Skip, go:build ignore found nothing).

## Not examined
Line-by-line diff of report.go, main.go, README.md and plan; the Gaps above; plain-copy install tier (3a) not triggered; Docker not applicable.

## Environment
Linux 6.12 (Debian), go from PATH (version not recorded), no Docker used.
