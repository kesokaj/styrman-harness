---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-078 path-shadow, round 4   (2026-09-18)

Range reviewed: 4cdc0b5..28e4747

## Status: PASS

Scope: complete — ran per .agents/docs/testing/guide.md §§1 and 3, the targeted
test command, live install/status/list checks, and govulncheck; guide.md already
existed and needed no update.

## Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0.
- `go test -race -count=1 -v ./...` → exit 0, all packages `ok`.
  counts: pass 215, fail 0, skip 0 at 28e4747 for `go test -race -count=1 -v ./...`.
  vs previous: `.agents/docs/testing/2026-09-18-verify-bug078-path-shadow-3.md`,
  pass 215→215, skip 0→0. Unchanged, matching the dispatch note that this range
  adds no test function (only an extra assertion block inside
  `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`).
- `go test -race -count=1 -v -run 'TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice|TestInstallAndStatusNameAPathShadow' ./cmd/sh-agent/` → exit 0:
  `--- PASS: TestInstallAndStatusNameAPathShadow (0.10s)`
  `--- PASS: TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice (0.08s)`
- `go run ./cmd/sh-agent install --dry-run --home .agents/scratch/verify-bug078-4-home`
  (fresh scratch home) → exit 0. Output has `[dry-run] the binary would be rebuilt
  from ...`, 259 `write` lines, `[dry-run] 260 written, 0 unchanged, 0 removed, 0
  kept — 260 files managed`, `[dry-run] would create .../.claude/settings.json
  with "attribution": {"commit": "", "pr": ""}`, and one `note: ...settings.json
  does not set "attribution": {"commit": ""}...` line (about Claude Code
  attribution config, not PATH). Grepped the full output for `PATH` and `on
  PATH`: zero matches. No line containing "is not on PATH" about the bin dir, no
  "so `sh-agent` runs", no "comes before".
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug078-4-home`
  (same home, now real install) → exit 0. Printed:
  `note: .agents/scratch/verify-bug078-4-home/.local/bin is not on PATH, so
  `sh-agent` runs /home/kesokaj/.local/bin/sh-agent instead; put
  .agents/scratch/verify-bug078-4-home/.local/bin ahead of it`
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug078-4-home` →
  exit 0. Printed `drift      none — 259 files match the repo` (259 not 260
  because settings.json already existed unchanged from the install step above).
  No line containing "shadowed" anywhere in the output.
- `go run ./cmd/sh-agent list` → exit 0, prints skill/agent tables (27 skills, 7
  agents), no parse errors.
- `govulncheck ./...` → exit 0, `No vulnerabilities found.` (go1.27.1 SDK ahead
  of /usr/bin/go, per the environment note in AGENTS.md § Setup commands).

## Root cause  (FAIL only)
N/A — Status is PASS.

## Coverage of changed code
Diff `4cdc0b5..28e4747` touches source only in `cmd/sh-agent/main.go`, changed
lines: `installFiles` now guards the `pathHint` call with `!dryRun` (a dry run
writes no binary, so no PATH claim is made about it).
- `cmd/sh-agent/main.go` (`installFiles`, dry-run guard) — covered by the added
  assertion block at `cmd/sh-agent/main_test.go` lines ~1478-1485, inside
  `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`: it runs
  `install --dry-run` on a fresh, off-PATH home and asserts the output contains
  no PATH-claim string ("install --dry-run made a PATH claim about a binary it
  did not write"). Confirmed passing above, and confirmed live (the dry-run
  command above produced zero PATH-related output).
- No other `.go` source lines changed in this range (the rest of the diff is
  docs: `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`,
  `README.md`, and new review/testing records — no test files besides
  `main_test.go` changed, matching the one new assertion block).
- Suppressions check: `git diff 4cdc0b5..28e4747 | grep -nE
  '(//nolint|//lint:ignore|@ts-ignore|@ts-expect-error|eslint-disable|t\.Skip|t\.SkipNow|go:build ignore|\.skip\(|\.only\()'`
  — the only hits are inside `.agents/docs/reviews/*.md` and `testing/*.md`
  prose quoting the absence of suppressions in earlier rounds; no match inside
  `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`,
  `internal/harness/pathguard.go`, or `internal/harness/pathguard_test.go`. No
  suppressions added.

## Not examined
- Sandbox CLI starts (guide §4) and real-home smoke checks (guide §5): this
  range touches no `context/`, `agents/`, or `skills/` file (only
  `cmd/sh-agent/main.go`, its test, `internal/harness/pathguard*.go`, and docs),
  so those tiers do not apply per the guide's own scoping ("every change to
  context/, agents/, or skills/").
- `report --dry-run` (guide §3b): out of scope, no `internal/harness/report.go`
  or report wiring touched in this range.
- Plain-copy install tier (guide §3a): out of scope, `IsCheckout`/rebuild logic
  untouched in this range.
- Full skill-matrix (guide "Full matrix"): not run; not claiming "complete" in
  that sense, only that the scoped commands above are complete for this range.

## Environment
Linux 6.12.107+deb13-amd64. Go: user-level go1.27.1 SDK ahead of
/usr/bin/go (per AGENTS.md § Setup commands note on Debian trixie).
`which govulncheck` → `/home/kesokaj/.go/bin/govulncheck`. Tree was clean and
frozen at 28e4747 for the whole run (`git status --short` empty at start; no
tree-mutating git command run).
