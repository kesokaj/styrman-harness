---
owner: verifier
date: 2026-09-17
---
# Verification — release-binary Phase 4   (2026-09-17)

Range reviewed: working tree at HEAD 4d45155 (branch release-binary, clean at dispatch time)

Guide followed: `.agents/docs/testing/guide.md` (existing, unchanged).

## Status: PASS

Scope: gofmt/vet/full race suite/-v counts, Phase 4 test tier, scratch-home install+status
drift check, govulncheck — all ran. See "Not examined" for what this run could not cover
(no real signed release, no `sh-agent update` against a live GitHub release, no `report`).

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 -v ./...` → exit implied 0 (all three packages report `ok`):
   ```
   ok  	github.com/kesokaj/styrman-harness	1.085s
   ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	3.817s
   ok  	github.com/kesokaj/styrman-harness/internal/harness	2.624s
   ```
   Full output saved to `.agents/scratch/verify-phase4/go-test-full.txt`.
   counts: pass 96, fail 0, skip 0 at 4d45155 for `go test -race -count=1 -v ./...`
   vs previous: `.agents/docs/testing/2026-09-17-release-binary-phase-3.md`, pass 87→96, skip 0→0.
   Rise in pass count only (no fall, no new skips) — consistent with Phase 4 adding new
   Release/UpdateFromRelease/ChecksumFor tests; not a FAIL condition.
4. Phase 4 tier: `go test -race -count=1 -run 'Release|ChecksumFor|UpdateFromRelease' -v ./internal/harness/ ./cmd/sh-agent/` → exit 0.
   ```
   --- PASS: TestUpdateFromReleaseReplacesAndInstalls (0.02s)
   --- PASS: TestUpdateFromReleaseUpToDateAndCheck (0.01s)
   --- PASS: TestUpdateFromReleaseRefusals (0.15s)
   --- PASS: TestReleaseSignerMatchesTheDeployGuide (0.00s)
   --- PASS: TestChecksumFor (0.00s)
   ok  	github.com/kesokaj/styrman-harness/internal/harness	1.188s
   --- PASS: TestReleaseFromInfo (0.00s)
   --- PASS: TestReleaseBinaryWithoutCheckout (0.20s)
   --- PASS: TestReleaseBinaryUpdateCheckExitsTwo (0.00s)
   ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	1.225s
   ```
   8 PASS, 0 FAIL, 0 SKIP. Output saved to `.agents/scratch/verify-phase4/go-test-phase4tier.txt`.
5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase4-home` → exit 0.
   "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; settings.json created
   with attribution off, as expected.
   `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase4-home` → exit 0.
   `drift      none — 259 files match the repo` (one file, settings.json, is the
   installer-created local file not tracked by drift count — expected, matches guide §3
   "must report drift none"). Status also reported `repo ... uncommitted changes` /
   `commit 4d45155+dirty` — this reflects an untracked file already present in the
   working tree before this run (`.agents/docs/reviews/2026-09-17-security-dab2565-4d45155.md`,
   not created by this verification and not part of the release-binary code), plus this
   log file itself; `git status --short` confirms no tracked-file modifications. Not a
   defect in Phase 4 code.
6. `govulncheck ./...` → exit 0. "No vulnerabilities found." (matches AGENTS.md expectation for go1.27.1). Output saved to `.agents/scratch/verify-phase4/govulncheck.txt`.

## Root cause  (FAIL only)
N/A — no failures found; all commands passed.

## Coverage of changed code
Phase 4 adds the release-download/verify/install path (`internal/harness` release.go /
update-from-release logic, `cmd/sh-agent` release wiring). All of it is covered by the
Phase 4 tier tests listed above (TestUpdateFromRelease*, TestReleaseSignerMatchesTheDeployGuide,
TestChecksumFor, TestReleaseFromInfo, TestReleaseBinary*), which ran and passed as part of
both the full suite and the targeted tier.

## Not examined
- No real signed GitHub release exists yet: the update path is only exercised against the
  tests' fake `gh` / fake release server, never a live `gh release download` or a real
  SSH-signed SHA256SUMS from GitHub.
- `sh-agent update` was NOT run in this session outside a scratch home, and never against
  the network, per the dispatch's hard limit.
- `sh-agent report` was not run at all in this session (dispatch only required it be run
  with `--dry-run` if run; it was not needed for Phase 4 and was skipped).
- No plain-copy (`.git`-less mirror) install tier (guide §3a) was run — Phase 4 does not
  touch `IsCheckout`/rebuild-decision logic, so it was judged out of scope.
- Sandbox CLI starts (guide §4) and real-home smoke checks (guide §5) were not run — no
  change to `context/`, `agents/`, or `skills/` in this range.

## Environment
Linux 6.12.107+deb13-amd64, go1.27.1 linux/amd64 (`.agents/docs`-recorded user-level SDK,
ahead of any system Go). Colima/Docker not applicable (no Dockerfile change).


## Caveat added by the main agent, 2026-09-17

This run is INCONCLUSIVE as evidence for any one commit. It was dispatched against the live
working tree at 4d45155, and while it ran the main agent was editing
`internal/harness/release.go` and `release_test.go` to fix the Phase 4 security review's
findings. The test binaries this run compiled may therefore have held a partially edited
tree rather than 4d45155's content, and its counts (96 pass) predate the fixes either way.

Nothing here is retracted — every command it ran did exit 0 — but the phase's verification
rests on the clean re-run at `.agents/docs/testing/2026-09-17-release-binary-phase-4-after-fixes.md`
instead. The process defect (dispatching a verifier against a tree the main agent then edits)
is recorded as a bug.
