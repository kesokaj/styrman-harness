---
owner: verifier
date: 2026-09-28
---
# Verification — update-from-zip Phase 2   (2026-09-28)

Range reviewed: 1cb2608..worktree

## Status: PASS

Scope: uncommitted Phase 2 (cmd/sh-agent/main.go --from-zip flag/updateFromZip/status zip lines/knownCheckout + cmd/sh-agent/zip_test.go), per .agents/docs/plans/2026-09-28-update-from-zip.md. No commit at this checkpoint (plan D2).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race ./...` (bare) → exit 0; ok for module root, cmd/sh-agent, internal/harness (full pkg output, no per-test names since not -v); log at .agents/scratch/verify-updatefromzip-p2/go-test-full.txt

- `go test -race -count=1 -v -run 'TestUpdateFromZip|TestStatusAfterZip|TestKnownCheckoutZip|TestHarnessMap' ./...` → exit 0; counts: pass 6, fail 0, skip 0 at 1cb2608+dirty for the above command; tests: TestHarnessMap PASS, TestUpdateFromZip PASS, TestUpdateFromZipRefusals PASS, TestUpdateFromZipDryRunWritesNothing PASS, TestStatusAfterZip PASS, TestKnownCheckoutZip PASS; vs previous: none (first run of these tests)

- Mutation (a): cmd/sh-agent/main.go:1047 `return installFiles(src, home, vendor, dryRun, out)` → `return installFiles(src, home, vendor, false, out)` (target unique: 1 occurrence). Saved copy .agents/scratch/verify-updatefromzip-p2/main.go.orig (sha256 5bb1216083b2497657647aed073998c5c6ca50dadf349c8d9861c0fc15d65083). Mutation compiled (go build ./... exit 0). `go test -race -count=1 -v -run TestUpdateFromZipDryRunWritesNothing ./cmd/sh-agent/...` → FAIL as expected: "dry run on a fresh home wrote [.claude .config .gemini .local]". File restored; sha256 verified identical to saved copy.
- Mutation (b): cmd/sh-agent/main.go:372 `fromZip = found && fromZip` → `fromZip = false` (target unique: 1 occurrence). Mutation compiled (go build ./... exit 0). `go test -race -count=1 -v -run TestStatusAfterZip ./cmd/sh-agent/...` → FAIL as expected (exit 1, status output shows no zip-install lines when fromZip forced false). File restored; sha256 verified identical to saved copy.

- `go run ./cmd/sh-agent --help` → exit 0; output shows `update --from-zip FILE [--dry-run]  zip place: unpack a repository zip you downloaded, install its harness with the running binary (the binary is not replaced), delete the unpacked copy`
- `git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules` → exit 0, no output (nothing untracked or ignored in the deployed-file source trees)

## Root cause
None. All commands, named tests, and both must-fail mutations behaved as expected; no failures to explain.

## Coverage of changed code
- cmd/sh-agent/main.go (`--from-zip` flag, `updateFromZip`, `status` zip lines, `knownCheckout`) — covered by cmd/sh-agent/zip_test.go: TestUpdateFromZip, TestUpdateFromZipRefusals, TestUpdateFromZipDryRunWritesNothing, TestStatusAfterZip, TestKnownCheckoutZip (all PASS); mutation testing on `updateFromZip`'s dry-run guard and `status`'s fromZip detection confirms these tests actually exercise the behavior (both must-fail mutations produced the expected named failure).
- cmd/sh-agent/zip_test.go — new test file itself, exercised by the -run above.
- No added `//nolint`, `//lint:ignore`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, or `.only(` found in the reviewed diff (not separately diffed here beyond the named-test and full-suite runs; see Not examined for what a full diff-suppression grep would additionally need).

## Not examined
- Phase 1 files (`internal/harness/apply.go`, `internal/harness/repo.go`, `internal/harness/source.go`, `internal/harness/zip.go`, `internal/harness/zip_test.go`) and `ROADMAP.md` — covered by `.agents/docs/testing/2026-09-28-verify-update-from-zip-phase-1.md`, not re-verified in this Phase 2 run (dispatch scoped to Phase 2 files only).
- A dedicated diff-based grep for added suppression/skip directives (`//nolint`, `t.Skip`, `.only(`, etc.) across the full `git diff 1cb2608..worktree` was not run as a separate step; the full-suite (`go test -race ./...`) and named-test runs showed no skipped or failing tests, which is partial but not equivalent evidence.
- `govulncheck ./...` (listed in AGENTS.md Setup Commands) was not run — not part of the assigned command list for this dispatch.
- `go run ./cmd/sh-agent install --home .agents/scratch/x` scratch-home install/drift check (guide.md § 3) was not run — not part of the assigned command list; no change to `context/`, `agents/`, or `skills/` is present in this diff to trigger it, and step 5's git-status check confirms those trees are untouched.
- govulncheck / toolchain version were not inspected in this run.

## Environment
Linux 6.12.107+deb13-amd64, repo at 1cb2608 + uncommitted Phase 1-2 changes.
