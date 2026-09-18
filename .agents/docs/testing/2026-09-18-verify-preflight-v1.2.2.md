---
owner: verifier
date: 2026-09-18
---
# Verification — preflight v1.2.2   (2026-09-18)

Range reviewed: v1.2.1..f91257c

## Status: PASS

Scope: complete for this range — only `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`,
`internal/harness/pathguard.go`, `internal/harness/pathguard_test.go` are Go changes; no
file under `context/`, `agents/`, or `skills/` changed, so guide.md §§2, 3a, 3b, 4, 5 and
the skill-matrix do not apply to this range. Ran guide.md §1 (Go checks), the live
install/status/list checks, govulncheck, the embed-contamination guard test, and the
untracked-files guard the dispatch named. testing/guide.md already existed and needed no
update.

## Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0.
- `go test -race -count=1 ./...` → exit 0, all three packages `ok`.
- `go test -race -count=1 -v ./...` (saved to `.agents/scratch/verify-v122/go-test-v.txt`)
  → exit 0.
  counts: pass 215, fail 0, skip 0 at f91257c+dirty for `go test -race -count=1 -v ./...`
  (`grep -cE '^\s*--- (PASS|FAIL|SKIP)'`; an unindented-only grep undercounts subtests —
  first pass wrongly showed 108 before the regex was corrected to allow leading
  whitespace).
  vs previous: `.agents/docs/testing/2026-09-18-verify-bug078-path-shadow-4.md`,
  pass 215→215, skip 0→0. Unchanged, matching the dispatch note that the range adds no
  test.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v122-home` → exit 0.
  Printed: `note: .agents/scratch/verify-v122-home/.local/bin is not on PATH, so
  `sh-agent` runs /home/kesokaj/.local/bin/sh-agent instead; put
  .agents/scratch/verify-v122-home/.local/bin ahead of it`. No `warning:` line about PATH.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v122-home` → exit 0.
  Printed (in full): `repo       .../styrman-harness (main @ f91257c, uncommitted
  changes)` / `installed  ... from commit f91257c+dirty (vendor: all, 260 files, ...)` /
  `binary     f91257c62e0f+dirty` / `retro      none yet ...` / `repo has   27 skills, 7
  agents, 0 commands, 0 rules` / `claude     ~/.claude  (...)` / `agy
  ~/.gemini/config  (...)` / `drift      none — 259 files match the repo`. No `shadowed`
  line; ends `drift none`. ("uncommitted changes" / "+dirty" reflects this verify log
  file itself, newly written under a tracked docs directory before this snapshot; no
  source file is modified.)
- `go run ./cmd/sh-agent list` → exit 0; printed 27 skills, 7 agents, tables well-formed
  (spot-checked head and tail).
- `govulncheck ./...` (go1.27.1) → "No vulnerabilities found."
- `go test -run TestEmbedHoldsOnlyTrackedFiles ./...` → exit 0 in all 3 packages
  (2 report "[no tests to run]", the harness package runs it and passes).
- `git status --porcelain --ignored --untracked-files=all -- context skills agents
  commands rules` → empty output (embed-contamination guard clean).

## Coverage of changed code
- `internal/harness/pathguard.go`:
  - `PathShadows` — covered by `internal/harness/pathguard_test.go:TestPathShadows`
    (plus edge cases at lines 350, 358).
  - `executableAt` — covered by `TestExecutableAt`.
- `cmd/sh-agent/main.go`:
  - `pathHint` (rewritten to call `harness.PathShadows`, three-way switch) — covered by
    `TestInstallAndStatusNameAPathShadow` and
    `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`.
  - `status`'s new `shadowed` line — covered by the same two tests (shadow-present case
    in the first, off-PATH absence case in the second).
  - the `!dryRun` gate around `pathHint` in `install` — covered by
    `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`'s final block (`install
    --dry-run` on a fresh home, asserted to print no PATH line).
- No added `//nolint`, `//lint:ignore`, `t.Skip`, `.only(`, or similar suppression found
  in the range's diff (checked by eye against the 4-file Go diff above; the range also
  touches only docs/markdown files otherwise).

## Not examined
- Tiers 2 ("harness self-check"), 3a (plain-copy install), 3b (`report` dry-run), 4
  (sandbox CLI starts), 5 (real-home smoke checks), and the full skill-matrix: none apply
  — this range touches no file under `context/`, `agents/`, or `skills/`, and no file
  under `IsCheckout`/rebuild-decision logic (`cmd/sh-agent` rebuild path, `internal/
  harness/repo.go`) beyond what `TestDryRunNeverBuilds`/`TestPlainCopyNeverRebuilds`
  already exercise unchanged, and no file touching `report.go` or the `reorder`
  `takesValue` map.
- Non-Go changed files (bug record, changes entries, plans, reviews, deploy docs,
  README, ROADMAP, handoff, testing logs from the round-1..4 BUG-078 work) were not
  reviewed for content — out of scope for a Go-verification pass; `/review --docs` is
  the mechanism named in `AGENTS.md` for docs changes, not this dispatch.
- Container/Docker tier: no Dockerfile in this range, not run.

## Environment
Linux 6.12.107+deb13-amd64; `go version` → go1.27.1 linux/amd64 (user-level SDK per
AGENTS.md § Setup commands, ahead of `/usr/bin/go`); no toolchain directive in `go.mod`.
`govulncheck` clean (0 advisories), consistent with go1.25.8+ per AGENTS.md.
