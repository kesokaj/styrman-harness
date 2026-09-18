---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-079 stale running binary (round 3)

Range reviewed: c868ed4..a078848

## Status: PASS

Scope: Go checks (gofmt/vet/test), targeted BUG-079 test subset, scratch-home
install+status drift check, govulncheck, and a live checkout-update rehearsal
(update then update --check) with the running binary reachable by name via
PATH. `testing/guide.md` § 5's real-home post-update check (`hash -r`,
`sh-agent version`, `sh-agent status` shadowed-line check) is explicitly out
of scope — it requires a real `$HOME`, not exercised here.

## Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0.
- `go test -race -count=1 -v ./...` → exit 0; counts: pass 226, fail 0, skip 0
  at a078848 for `go test -race -count=1 -v ./...`.
  vs previous: `.agents/docs/testing/2026-09-18-verify-bug079-stale-running-binary-2.md`,
  pass 230→226, skip 0→0. Fall in pass count explained: this range renamed
  `TestUpdateNamesTheBinaryItRanFrom` (4 subtests) to
  `TestUpToDateUpdateSaysNothingAboutTheRunningBinary` (no subtests, one
  top-level case) — a net −4, matching 230−4=226 exactly; confirmed no other
  test was removed (`git diff c868ed4..a078848` adds tests only, no `t.Skip`
  or removed `Test` funcs found).
- `go test -race -count=1 -v -run 'TestUpdateFromReleaseReplacesAndInstalls|TestUpdateFromReleaseUpToDateAndCheck|TestRunningElsewhere|TestUpToDateUpdateSaysNothingAboutTheRunningBinary|TestCheckoutUpdateNamesTheBinaryItRanFrom' ./...`
  → exit 0. All five names found and PASS: `TestUpToDateUpdateSaysNothingAboutTheRunningBinary`,
  `TestCheckoutUpdateNamesTheBinaryItRanFrom` (cmd/sh-agent), `TestRunningElsewhere`
  (8 subtests, all PASS), `TestUpdateFromReleaseReplacesAndInstalls`,
  `TestUpdateFromReleaseUpToDateAndCheck` (internal/harness).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug079-3/home`
  → "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed";
  expected PATH note (scratch home not on PATH) and settings-attribution note
  printed, both unrelated to BUG-079.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug079-3/home`
  → `drift      none — 259 files match the repo`. (`binary` line reported
  `a0788481afbb+dirty` because this repo's working tree carries one untracked
  file — this very log, written before the run — the only untracked path
  present; no other drift.)
- `govulncheck ./...` (go1.27.1) → "No vulnerabilities found."

## Live checkout-update rehearsal (running binary reachable by name)

Setup: built `.agents/scratch/verify-bug079-3/other/sh-agent` from this repo;
planted a placeholder file at `.agents/scratch/verify-bug079-3/home2/.local/bin/sh-agent`
(a different, non-matching file, to be replaced); cloned this repo (at
a078848) to `.agents/scratch/verify-bug079-3/clone`; ran the "other" binary
with it prepended to PATH.

`update` (full):
```
up to date with upstream (a078848)
building sh-agent from .../verify-bug079-3/clone
260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
note: .../home2/.local/bin is not on PATH, so `sh-agent` runs .../other/sh-agent instead; put .../home2/.local/bin ahead of it
settings: .../home2/.claude/settings.json created with "attribution": {"commit": "", "pr": ""} — Claude Code's commit and PR trailers are off
note: this update ran from .../other/sh-agent, not from the installed binary .../home2/.local/bin/sh-agent; update replaces only the installed one, so the file it ran from stays as it is. A shell that has already resolved `sh-agent` to that path keeps running it until `hash -r` (bash) or `rehash` (zsh); point it at the installed binary (`ln -sfn`) or remove it — when that path is root-owned, prefer removing it
```
Ends exactly with "when that path is root-owned, prefer removing it" as expected.

`update --check` (same PATH/home/repo): output is exactly
```
up to date with upstream (a078848)
```
No BUG-079 note printed, as expected (`--check` takes no action, replaces
nothing, so `NoteRunningElsewhere` is never reached on that path — matches
`main.go`'s update function calling it only after a successful `Reinstall`).

## Root cause
N/A — no failures found in this range.

## Coverage of changed code
Non-test functions changed in this range (`git diff c868ed4..a078848 -- '*.go'`,
test files excluded):
- `internal/harness/pathguard.go`: `RunningElsewhere` (new) — covered by
  `internal/harness/pathguard_test.go` `TestRunningElsewhere` (8 subtests, all
  cases: self, link-to-installed, system-copy-by-name, PATH-order,
  unreachable-by-name, link-reported-at-PATH-entry, installed-file-missing,
  exe-missing — run above, all PASS).
- `internal/harness/pathguard.go`: `NoteRunningElsewhere` (new) — covered by
  `internal/harness/release_test.go` `TestUpdateFromReleaseReplacesAndInstalls`
  (asserts the note after a replace) and `TestUpdateFromReleaseUpToDateAndCheck`
  (asserts no note when up to date / --check), and by
  `cmd/sh-agent/main_test.go` `TestCheckoutUpdateNamesTheBinaryItRanFrom` and
  `TestUpToDateUpdateSaysNothingAboutTheRunningBinary` for the checkout `update`
  path — all run above, all PASS. Also exercised live above (both branches:
  note present after replace, note absent after --check).
- `cmd/sh-agent/main.go` `update()`: new call to `harness.NoteRunningElsewhere`
  after a successful `Reinstall` — covered by the same two `cmd/sh-agent`
  tests above, plus the live rehearsal.
- `internal/harness/release.go` `UpdateFromRelease`: new call to
  `NoteRunningElsewhere` after a successful install — covered by the two
  `internal/harness/release_test.go` tests above.

No suppressions (`//nolint`, `//lint:ignore`, `@ts-ignore`, `@ts-expect-error`,
`eslint-disable`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, `.only(`)
added anywhere in this range's `*.go` diff — checked with
`git diff c868ed4..a078848 -- '*.go'` and a grep for those markers; none found.

## Not examined
- `.agents/docs/testing/guide.md` § 5's post-update check added this range
  (`hash -r`/`rehash`, `sh-agent version`, `sh-agent status` for a `shadowed`
  line) — explicitly out of scope per dispatch: it is a real-home check, and
  no real `$HOME` was touched in this run.
- Frontend/UI: none in this range (no `web/`, `ui/`, or `package.json` touched).
- Docs-only files changed in this range (plans, reviews, README, ROADMAP,
  deploy guides, `delivery.md`) were not proofread beyond the grep for
  suppression markers; their content review is `/review --docs`' job, not
  this verifier dispatch's.

## Environment
- OS: Linux 6.12.107+deb13-amd64 (Debian trixie).
- Go: go1.27.1 (`go.mod` carries no toolchain directive, per ADR-029).
- `govulncheck`: 1.1.4-class binary at `~/.go/bin/govulncheck`, clean run.
- No Docker/Colima involved (no Dockerfile in this range).
- All work confined to `.agents/scratch/verify-bug079-3/`; no git command
  altered the working tree; `sh-agent report` never invoked; no network
  access used.
