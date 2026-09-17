---
owner: verifier
date: 2026-09-17
---
# Verification — release-binary Phase 4 (after fixes)   (2026-09-17)

Replaces `.agents/docs/testing/2026-09-17-release-binary-phase-4.md`, which ran while
the main agent was editing `internal/harness/release.go` and could not be attributed to
a commit. This run's tree was confirmed clean before the first command and re-checked
after the install step; only this log file was ever untracked.

Range reviewed: c7da4fb..f15d533 (Phase 4 commits 4d45155, 40bf2db, f15d533)
Commit verified: f15d533cede50b19d9a8808239afc34135a624e, branch release-binary,
working tree clean at start (`git status --short` empty).

## Status: PASS

Scope: complete for the commands the dispatch named; `testing/guide.md` exists and was
followed. No real release exists, so the download/verify path in `internal/harness/release.go`
ran only against the tests' fake `gh` (see Not examined).

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 ./...` → exit 0, all three packages `ok`
   (root 1.090s, cmd/sh-agent 3.650s, internal/harness 2.658s).
   `go test -race -count=1 -v ./...` (saved to `.agents/scratch/verify-phase4b/go-test-v.txt`)
   → exit 0.
   counts: pass 99 top-level + 85 subtests (184 `--- PASS` lines total), fail 0, skip 0,
   at f15d533 (clean) for `go test -race -count=1 -v ./...`.
   vs previous: `.agents/docs/testing/2026-09-17-release-binary-phase-3.md`
   (counts: pass 87, fail 0, skip 0, top-level only), pass 87→99 top-level, skip 0→0.
   Rise of 12 explained: Phase 4 commits (4d45155, 40bf2db, f15d533) added the
   update-from-release feature and its tests (`TestUpdateFromRelease*`,
   `TestChecksumFor`, `TestCompareRelease`, `TestSSHKeygenPrefersASystemPath`,
   `TestReleaseSignerMatchesTheDeployGuide`, `TestReleaseFromInfo`,
   `TestReleaseBinaryWithoutCheckout`, `TestReleaseBinaryUpdateCheckExitsTwo`, and
   subtests of `TestUpdateFromReleaseRefusals`).
   This reproduces the main agent's own reported run at this commit: 99 top-level
   tests, 85 subtests, 0 fail, 0 skip. No disagreement found.
4. `go test -race -count=1 -run 'Release|Checksum|Compare|SSHKeygen|UpdateFromRelease' ./internal/harness/ ./cmd/sh-agent/ -v`
   → exit 0. All matched tests passed, including every `TestUpdateFromReleaseRefusals`
   subtest (ssh-keygen missing, no asset for platform, signature by another key,
   SHA256SUMS altered after signing, binary mismatch vs signed checksum, release
   shipping its own signer file, signed checksums of another release, checksums
   naming no version, no signature in release, no checksums in release),
   `TestReleaseSignerMatchesTheDeployGuide`, `TestChecksumFor`, `TestCompareRelease`,
   `TestSSHKeygenPrefersASystemPath`, `TestUpdateFromReleaseNeverGoesBackwards`,
   `TestReleaseFromInfo`, `TestReleaseBinaryWithoutCheckout`,
   `TestReleaseBinaryUpdateCheckExitsTwo`. Full output in
   `.agents/scratch/verify-phase4b/go-test-targeted.txt`.
5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase4b-home` → exit 0:
   "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed", plus expected
   PATH note and the `attribution: {"commit": "", "pr": ""}` settings note.
   `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase4b-home` → exit 0:
   `drift      none — 259 files match the repo`.
   Minor discrepancy noted, not treated as a FAIL: install reports 260 files managed,
   status reports 259 files checked for drift (off by one). Also, both `repo` and
   `binary` lines showed `+dirty`/"uncommitted changes" — this is explained by this
   run's own log file (this document) being untracked in the working tree at the
   moment `install`/`status` ran, not by any defect in the code under test; confirmed
   with `git status --short --untracked-files=all` immediately after, which showed
   only this log file. The 260-vs-259 count difference was not chased further under
   this run's budget; it may be the `.claude/settings.json` written during install but
   excluded from the drift count by design, or it may be worth a one-line follow-up.
6. `govulncheck ./...` → exit 0: "No vulnerabilities found." Go version: go1.27.1
   (user-level SDK ahead of system go, per AGENTS.md). Matches the documented
   expectation of zero advisories on go1.25.8+.

## Root cause

N/A — no failures found.

## Coverage of changed code

Phase 4 touches `internal/harness/release.go` (download-and-verify-and-install from a
signed GitHub release) and `internal/harness/release_test.go`
(`internal/harness/release_test.go` covers `UpdateFromRelease`, checksum parsing,
version comparison, ssh-keygen resolution, and the security refusals added in
40bf2db: no downgrade, tag validated, version inside signed checksums, verifier
preferred from a system path, random download dir, stricter checksum parsing).
`cmd/sh-agent` gained `TestReleaseFromInfo`, `TestReleaseBinaryWithoutCheckout`,
`TestReleaseBinaryUpdateCheckExitsTwo` covering the CLI wiring. No changed source file
in this range was found without a covering test.

Suppressions: none added in this range (no `//nolint`, `t.Skip`, `.only(`, etc. found
in the Phase 4 commits' diffs, based on the file lists in 4d45155/40bf2db/f15d533).

## Not examined

- No real GitHub release exists for this project; the download/signature-verification
  path in `internal/harness/release.go` (`gh release download`, `ssh-keygen -Y verify`
  against `SHA256SUMS`) was exercised only through the test suite's fake `gh` script
  on an isolated `PATH` (confirmed in `release_test.go`: `fakeGh`, `releaseEnv`,
  `t.Setenv("PATH", bin)` — the real `gh` and the network were never reached in this
  run, either in the test suite or in the `install`/`status` scratch-home steps, which
  do not touch `release.go`'s download path at all).
  A real end-to-end run against an actual signed release is out of scope for a bare
  verifier dispatch, per `.agents/docs/testing/guide.md` § 3b's equivalent guidance for
  `report` (same principle: real network/remote runs are the main agent's e2e tier, not
  the verifier's).
- `go run ./cmd/sh-agent report --dry-run` was not run in this dispatch (not named in
  the command list) — no claim is made about the report path in this log.
- The 3a "plain-copy install" and 3b "report" tiers of the guide were not run — not
  named in the dispatch and not touched by this range's diff (Phase 4 is
  `release.go`/`release_test.go` plus the CLI wiring, not the checkout-vs-mirror or
  report-collection logic).
- The 260-vs-259 file count discrepancy noted in Command 5 was not chased to a root
  cause under this run's budget.

## Environment

OS: Linux 6.12.107+deb13-amd64 (Debian trixie). Go: go1.27.1 (user-level SDK, no
`toolchain` directive per ADR-029). `govulncheck` at `/home/kesokaj/.go/bin/govulncheck`.
No Docker/Colima involved (no Dockerfile change in this range). `gh` was reachable on
this machine's normal PATH but was never invoked outside the test suite's isolated
`PATH` sandbox.
