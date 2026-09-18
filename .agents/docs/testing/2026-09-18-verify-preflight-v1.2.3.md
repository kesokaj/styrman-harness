---
owner: verifier
date: 2026-09-18
---
# Verification — preflight v1.2.3   (2026-09-18)

Range reviewed: v1.2.2..30f7488

## Status: PASS

Scope: complete — guide tiers 1, 2, 3 (scratch install/status/perms), 3c (verb x place-kind
matrix), plus the release-gate commands the dispatch named (cross-compile, govulncheck,
embed test, tracked-tree git status). Tiers 3a/3b/4/5 not triggered by this range (no change
to `IsCheckout`/rebuild logic, `report.go`, `context/`/`agents/`/`skills/`, or a real-home
install) and are not run. `testing/guide.md` already existed; not rewritten (no step change).

## Commands

- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 -v ./...` → exit 0. counts: pass 250, fail 0, skip 0 at 30f7488 for
  `go test -race -count=1 -v ./...`; vs previous:
  `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-2-3.md`, pass
  250→250, skip 0→0 (matches the dispatch's stated expectation that the last commit added no
  test). Log: `.agents/scratch/verify-v123-run/go-test-1.txt`.
  (Counting note: `grep -c '^--- PASS'` under-counts because indented subtest lines are
  anchored; recounted with `grep -c -- '--- PASS:'` — both PASS and FAIL/SKIP confirmed this
  way.)
- `GOOS=windows go build ./...` → exit 0 (WIN_OK).
- `GOOS=darwin GOARCH=arm64 go build ./...` → exit 0 (DARWIN_OK).
- `GOOS=linux GOARCH=arm64 go build ./...` → exit 0 (LINUXARM_OK).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v123` (fresh home) → 260 `add`
  lines (one per managed file), `260 written, 0 unchanged, 0 removed, 0 kept — 260 files
  managed`, PATH note, settings-created note. Log:
  `.agents/scratch/verify-v123-run/install-1.txt`.
- Same install, second run → `0 written, 260 unchanged, 0 removed, 0 kept — 260 files
  managed` (no `add`/`update` lines). Log:
  `.agents/scratch/verify-v123-run/install-2.txt`.
- `stat -c %a .agents/scratch/verify-v123/.config/sh-agent` → `700`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v123` → `drift none — 259 files
  match the repo` (259, not 260: settings.json is created but not manifest-diffed the same
  way — consistent with the install-2 count of 260 managed vs the matrix doc's own
  "drift none — 259 files match the repo" wording at other commits, so not a regression).
- `go run ./cmd/sh-agent list` → prints 27 skills, 7 agents tables, no parse errors.
- `govulncheck ./...` → `go version go1.27.1 linux/amd64`; "No vulnerabilities found." (matches
  AGENTS.md's go1.25.8+ expectation).
- `go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...` → PASS, ok (all 3 packages,
  2 report "no tests to run").
- `git status --porcelain --ignored --untracked-files=all -- context skills agents commands
  rules` → empty (no output).

### § 3c — verb x place-kind matrix

Removed `.agents/scratch/matrix/{checkout,dev,home-checkout,home-plain,home-release,nogh,
plain,rel,widget,run.log}` (kept `run.sh`), then ran `bash .agents/scratch/matrix/run.sh` →
exit 0, `run.log` 485 lines, 53 `[exit ...]` cells (matches
`.agents/docs/testing/sh-agent-matrix.md`'s stated cell count).

13 cells exit 1, all documented `refused` cells in `sh-agent-matrix.md`: checkout/plain/
release `status (after drift)`, checkout/plain/release `status (after uninstall)`, plain/
release `update --check`, plain/release `update`, release `adopt`, release
`install --repo`, release `$SH_AGENT_REPO`. All 40 remaining cells exit 0.

Spot-checked quoted refusal text against the doc's Matrix table, verbatim match:
- `release install --repo`: `sh-agent: this is the release binary v9.9.9: it installs the
  harness built into it and never a tree on disk, so --repo=<path> is refused — unset it, or
  build sh-agent from that tree and run that build's install`
- `release $SH_AGENT_REPO`: same message, naming `$SH_AGENT_REPO=<path>` — matches.

One cosmetic deviation from the doc, not a code regression: `checkout version` printed
`sh-agent 30f7488928d0+dirty` rather than the doc's example `sh-agent (devel)`. Root cause:
this verification session's own untracked files (this log and a pre-existing untracked
`.agents/docs/reviews/2026-09-18-review-v1.2.2-30f7488-docs.md`, confirmed via `git status
--short` in the outer repo) make the checkout's VCS build-info stamp `+dirty` when
`go build ./cmd/sh-agent` runs from this tree — the same "uncommitted changes" the `status`
verb also reports. This is expected verifier-session noise (the tree itself, at 30f7488, is
otherwise clean per the dispatch's own git-status snapshot), not a functional difference; the
doc's exact version string was captured at a moment with no untracked files present and is
stale in that one respect. No behavioral cell disagrees with the spec.

## Root cause

N/A — PASS, no failures found.

## Coverage of changed code

Non-test functions added/changed in v1.2.2..30f7488 (`git diff --stat`, 80 files changed,
Phase 2 round 2/3 + audit changes): `settingsUnwritable`, `diffOne`, `ownedByCaller` (both
build-tag variants), `RunningElsewhere`, `NoteRunningElsewhere`, `ConfigDir`,
`TightenConfigDir` in `internal/harness/{apply,owned_unix,owned_other,state}.go`. Each has a
directly-named test added in the same range: `TestStatusWarnsOnOversizedSettings`,
`TestDiffVerbShowsDrift`/`TestDiffUsesTheHomeNotTmp`, `TestOwnedByCaller`,
`TestRunningElsewhere`/`TestNoteRunningElsewhereOnlyForTheRealHome`/
`TestNoteRunningElsewhereThroughASymlinkedHome`, `TestInstallMakesTheConfigDirPrivate`/
`TestThrowawayParentIsPrivate`. All pass in the `go test -race -count=1 -v ./...` run above
(250/250). No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) found added in this range —
not separately grepped line-by-line in this pass; see Not examined.

## Not examined

- Did not grep the full range diff line-by-line for added `//nolint`/`t.Skip`/`.only(`
  suppressions with reasons (budget); the test-count comparison (250→250, 0 skip both) is
  the indirect signal that none were added as skips.
- Guide tiers 3a (plain-copy install), 3b (`report` collection/redaction tests), 4 (sandbox
  CLI start), 5 (real-home smoke checks) not run — none triggered by this range per the
  guide's own trigger conditions, and 5 is out of scope for a `--home`-only, no-real-home
  dispatch.
- Did not independently verify every one of the 250 individual test bodies' assertions
  beyond reading PASS/FAIL/SKIP counts and a tail of the log.

## Environment

Linux 6.12.107+deb13-amd64, go1.27.1 (user-level SDK ahead of `/usr/bin/go`, per AGENTS.md),
govulncheck from `~/.go/bin/govulncheck`. No Docker/Colima involved (no Dockerfile in range).
