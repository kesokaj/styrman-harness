---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-079 review round 1 fixes (0f2da25)   (2026-09-18)

Range reviewed: c868ed4..0f2da25

## Status: PASS

Scope: complete — gofmt/vet/test suite, focused BUG-079 tests, install/status against a scratch
home, govulncheck, and a live check of `RunningElsewhere` through a checkout update in a scratch
clone, both with and without the running-elsewhere file on PATH. `testing/guide.md` exists and was
followed; no dirty-tree guard needed since the Go source under `cmd/`/`internal/` was frozen for
this run (concurrent `.agents/docs/` and `README.md` edits by `documentation` showed up only as
`+dirty` in `sh-agent status`'s repo line, touching no test).

## Commands
- `gofmt -l cmd internal && go vet ./...` → both clean, no output.
- `go test -race -count=1 -v ./...` → exit 0; counts: pass 230, fail 0, skip 0 at 0f2da25 for
  `go test -race -count=1 -v ./...`;
  vs previous: `.agents/docs/testing/2026-09-18-verify-bug079-stale-running-binary.md`, pass
  229→230, skip 0→0 (matches the expected +1 subtest, `no installed binary: silent`, on
  `RunningElsewhere`). Package results:
  `ok  	github.com/kesokaj/styrman-harness	1.087s`
  `ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	4.223s`
  `ok  	github.com/kesokaj/styrman-harness/internal/harness	2.798s`
- `go test -race -count=1 -v -run 'TestRunningElsewhere|TestUpdateNamesTheBinaryItRanFrom|TestCheckoutUpdateNamesTheBinaryItRanFrom' ./...`
  → exit 0. `TestUpdateNamesTheBinaryItRanFrom` (4 subtests incl. `no installed binary: silent`),
  `TestCheckoutUpdateNamesTheBinaryItRanFrom`, `TestRunningElsewhere` (8 subtests incl.
  `installed file missing: nothing to be behind, silent`) all PASS.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug079-2` → exit 0; `260 written, 0
  unchanged, 0 removed, 0 kept — 260 files managed`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug079-2` → exit 0; `drift none — 259
  files match the repo`, `binary 0f2da25fab2f+dirty` (dirty = the concurrent docs edit noted
  above), repo commit shows `0f2da25`.
- `govulncheck ./...` → exit 0; `No vulnerabilities found.`

## Live check — RunningElsewhere through a checkout update
Built a second binary at `.agents/scratch/verify-bug079-2/other/sh-agent` (`go build`), copied it
(not linked) to `.agents/scratch/verify-bug079-2/home/.local/bin/sh-agent` as the "installed"
file, cloned this repo locally with no network (`git clone -q "$(pwd)" .../clone`, clean at
0f2da25).

With `other/` on PATH:
```
PATH=".../other:$PATH" .../other/sh-agent update --home .../home --repo .../clone
```
→ exit 0; note printed verbatim:
```
note: this update ran from /home/kesokaj/workspace/styrman-harness/.agents/scratch/verify-bug079-2/other/sh-agent, not from the installed binary /home/kesokaj/workspace/styrman-harness/.agents/scratch/verify-bug079-2/home/.local/bin/sh-agent; update replaces only the installed one, so the file it ran from stays as it is. A shell that has already resolved `sh-agent` to that path keeps running it until `hash -r` (bash) or `rehash` (zsh); point it at the installed binary (`ln -sfn`) or remove it
```
Matches the spec'd wording exactly.

With `other/` NOT on PATH (plain `.../other/sh-agent update ...`, invoked by path, not by name):
→ exit 0; `0 written, 260 unchanged, 0 removed, 0 kept`; only the pre-existing "not on PATH" note
about the home's own bin dir printed. No "note: this update ran from" line — matches expectation
(no PATH entry resolves `sh-agent` to `other/sh-agent`, so it is nobody's `sh-agent` and
`RunningElsewhere` correctly stays silent).

## Root cause
N/A — no failures found in this range.

## Coverage of changed code
Diff `c868ed4..0f2da25` touches only:
- `cmd/sh-agent/main.go` — covered by `cmd/sh-agent/main_test.go`
  (`TestUpdateNamesTheBinaryItRanFrom`, `TestCheckoutUpdateNamesTheBinaryItRanFrom`) and by the
  live checkout-update legs above.
- `internal/harness/pathguard.go` (`RunningElsewhere`, new function) — covered by
  `internal/harness/pathguard_test.go` (`TestRunningElsewhere`, 8 subtests including the new
  "installed file missing: nothing to be behind, silent" case) and by the live check.

No production `.go` file in this range is uncovered. No added `//nolint`, `//lint:ignore`,
`t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, or `.only(` in the diff (checked via
`git diff c868ed4..0f2da25 -- '*.go'` grepped for those markers — none found).

## Not examined
- No `.md`/docs changes in this range were reviewed for content (task scoped this run to the
  frozen Go source; `documentation` agent's concurrent edits were not read).
- Did not exercise the real `sh-agent report` path (excluded per hard limits) or any network-based
  release-shaped update (offline by design, substituted with the checkout-update live check as
  instructed).
- `--check` mode's silence on the note was exercised only via the unit test
  (`TestUpdateNamesTheBinaryItRanFrom/--check_says_nothing`), not via a second live run.

## Environment
Linux 6.12.107+deb13-amd64. `go` per project's user-level SDK ahead of `/usr/bin/go` (see
project `AGENTS.md`). `govulncheck` at `/home/kesokaj/.go/bin/govulncheck`, reported no
vulnerabilities. All work under `.agents/scratch/verify-bug079-2/`; no network; no git command
altered this working tree (only a local `git clone` into scratch, and a `git build`/`go run`
against the frozen source).
