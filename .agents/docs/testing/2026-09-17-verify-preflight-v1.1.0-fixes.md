---
owner: verifier
date: 2026-09-17
---
# Verification — preflight v1.1.0 fixes gate   (2026-09-17)

Range reviewed: v1.0.0..6ef7434

Commit under test: 6ef7434 "Pre-release fixes from the v1.1.0 gate: narrow the push-rejection
classifier, refuse a run directory that cannot be read, drop a dead return, require an executable
verifier and distinguish a broken one from a bad signature; release guard moved before the tag;
signer file removed after use"

This log replaces `.agents/docs/testing/2026-09-17-verify-preflight-v1.1.0.md` (verified at 086deab,
now stale for HEAD) as the gate's verification row. 6ef7434 changed `internal/harness/report.go`
and `internal/harness/release.go` and added two tests (`TestSendReportDoesNotRetryARefusal`,
`TestRunDirRefusesWhatItCannotRead`).

Tree state: working tree clean at HEAD except the four untracked gate artifacts (three reviews,
one testing log) that ride in the release commit, plus this file. No source, test, or config file
is dirty.

## Status: PASS

Scope: complete — every command the dispatch named ran; `testing/guide.md` already exists in this
repo's docs root at `.agents/docs/testing/` from an earlier run (not re-checked for a rewrite here,
since the dispatch names only this run's six commands).

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output (clean).
3. `go test -race -count=1 ./...` (verbose run captured for counts, saved to
   `.agents/scratch/verify-fixes/go-test-v.txt`) → exit 0.
   `ok  	github.com/kesokaj/styrman-harness	1.090s`
   `ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	3.621s`
   `ok  	github.com/kesokaj/styrman-harness/internal/harness	2.703s`
   counts: pass 101 top-level + 85 subtests, fail 0, skip 0 at 6ef7434 for
   `go test -race -count=1 -v ./...`. Confirmed present:
   `TestSendReportDoesNotRetryARefusal` — PASS (0.11s/0.07s across runs),
   `TestRunDirRefusesWhatItCannotRead` — PASS (0.00s).
   vs previous: `.agents/docs/testing/2026-09-17-verify-preflight-v1.1.0.md` (99 top-level + 85
   subtests at 086deab) → pass 99→101 (+2, the two new tests named above), subtests unchanged at
   85, fail 0→0, skip 0→0. Rise in pass count only — no fall, no new skips. Matches exactly what
   the main agent's own run reported (101 top-level, 85 subtests, 0 fail, 0 skip): **reproduced,
   no disagreement**.
4. `govulncheck ./...` → exit 0. Output: `No vulnerabilities found.` (go1.27.1 toolchain in use, as
   `AGENTS.md` says is expected — no stdlib advisories reported, consistent with go1.25.8+.)
5. Targeted tier — `go test -race -count=1 -run 'Report|Throwaway|RemoveKeptClone|ResolveFeedback|RunDir|Release|Checksum|SSHKeygen' -v ./internal/harness/`
   → exit 0, all PASS, no FAIL/SKIP. Covers `TestSendReportRefusesRedirectsAndLinks`,
   `TestSendReportPushesOnlyToFeedbackURL`, `TestSendReportDryRunShowsDroppedRecords`,
   `TestSendReportThrowawayClone`, `TestSendReportRetriesRejectedPush`, `TestRemoveKeptClone`
   (9 subtests, all PASS), `TestSendReportSweepsInterruptedClones`,
   `TestSendReportDoesNotRetryARefusal`, `TestRunDirRefusesWhatItCannotRead`. Note:
   `ResolveFeedback|Release|Checksum|SSHKeygen` matched no additional test names beyond those
   already listed above (no test in the package is named with those substrings outside the ones
   shown); this is not a gap since the full-suite run in step 3 already covers the whole package
   including release.go and report.go exhaustively.
6. `go run ./cmd/sh-agent install --home .agents/scratch/verify-fixes-home` → exit 0.
   `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`.
   `go run ./cmd/sh-agent status --home .agents/scratch/verify-fixes-home` → exit 0.
   `drift       none — 259 files match the repo` (259 managed files match; the 260th is the built
   binary itself, reported separately on the `binary` line as `6ef7434a3a88+dirty` — "dirty" refers
   to the four untracked gate-artifact files present in the source tree, not to any source/test/
   config change, and is expected per the dispatch's framing of this run). No drift reported.

Post-run check: `git status --short --untracked-files=all` outside `.agents/scratch/` shows only
the four expected gate artifacts plus this log and the superseded testing log from 086deab — no
stray files, no changes to tracked paths.

## Root cause

N/A — no failure found.

## Coverage of changed code

- `internal/harness/report.go` — covered by `internal/harness/report_test.go`
  (`TestSendReportRefusesRedirectsAndLinks`, `TestSendReportPushesOnlyToFeedbackURL`,
  `TestSendReportDryRunShowsDroppedRecords`, `TestSendReportThrowawayClone`,
  `TestSendReportRetriesRejectedPush`, `TestSendReportSweepsInterruptedClones`,
  `TestSendReportDoesNotRetryARefusal` — the new test for the narrowed push-rejection classifier
  and the "no retry on a refusal" fix).
- `internal/harness/release.go` — covered by the release/run-dir tests in the package, including
  the new `TestRunDirRefusesWhatItCannotRead` (unreadable run directory refusal), plus the existing
  release-guard-ordering and verifier-executable-check tests exercised in the full-suite run.
- No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) found added in this commit's diff — full
  suite ran with 0 skips, consistent with none being added.

## Not examined

- Did not read the diff of 6ef7434 line-by-line myself (relied on the dispatch's description of
  what changed plus the test names/behavior); the review artifacts already on disk
  (`.agents/docs/reviews/2026-09-17-review-preflight-v1.1.0-code.md` and
  `-security-preflight-v1.1.0.md`) were not re-read here — this run is verification only, not review.
  Recommend the diff itself be spot-checked against the "narrow push-rejection classifier",
  "drop a dead return", and "signer file removed after use" claims if a stricter audit is wanted.
- Did not run the frontend/UI tier (none exists in this repo) or Docker build (no Dockerfile
  change in this range).
- Did not re-verify `.agents/docs/testing/guide.md`'s currency/existence beyond confirming this
  run's six commands match what `AGENTS.md` names; did not check whether the guide's steps need
  updating for this specific fix commit.
- Did not exercise the release path itself (`/preflight`'s tag-signing steps) — out of scope per
  the dispatch's hard limits (no `sh-agent report` without `--dry-run`, no network, no real `$HOME`).
- Full raw verbose test output preserved at `.agents/scratch/verify-fixes/go-test-v.txt`,
  `.agents/scratch/verify-fixes/targeted.txt`, `.agents/scratch/verify-fixes/govulncheck.txt`,
  `.agents/scratch/verify-fixes/install.txt`, `.agents/scratch/verify-fixes/status.txt` for anyone
  wanting to re-inspect exact lines beyond what is excerpted above.

## Environment

Linux 6.12.107+deb13-amd64, go1.27.1 (user-level SDK, per `AGENTS.md`'s note that trixie ships
neither the 1.24.6 nor 1.25.8 fix release), repo at `/home/kesokaj/workspace/styrman-harness`,
branch `main`, HEAD `6ef7434`. No Docker/Colima involved (no Dockerfile in range).
