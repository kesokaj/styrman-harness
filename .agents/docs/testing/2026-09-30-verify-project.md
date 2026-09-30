---
owner: verifier
date: 2026-09-30
---
# Verification — project   (2026-09-30)

Range reviewed: the whole repository at HEAD (3478a1c), no diff.

## Status: PASS

Scope: ran every command AGENTS.md § Setup commands names (`gofmt -l cmd internal`, `go vet ./...`,
`go test -race ./...`, `go run ./cmd/sh-agent install --home .agents/scratch/home`,
`go run ./cmd/sh-agent status --home .agents/scratch/home`, `go run ./cmd/sh-agent list`,
`govulncheck ./...`) and mapped every non-test `.go` file in the repo to its covering test.
`testing/guide.md` (2026-09-17) already present — cited, not rewritten (no step changed).
Sandbox/real-home CLI smoke tiers excluded per dispatch brief. Scope: complete for the
tiers named in the brief.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./... -count=1` (run once, in background, output to
  `.agents/scratch/verify-project/go-test-1.txt`) → exit 0. All 4 packages `ok`:
  root (`github.com/kesokaj/styrman-harness`, 1.103s), `cmd/sh-agent` (15.200s),
  `internal/harness` (3.302s), `internal/search` (7.186s).
  counts: pass 366, fail 0, skip 2 at 3478a1c+dirty (untracked doc files only, see
  Environment) for `go test -race -v ./... -count=1`.
  Both skips carry reasons: `TestRecall` ("measurement, not a gate: set SH_AGENT_RECALL=1
  to run it") and `TestHelperIndex` ("helper process", a subprocess-mode test helper, not
  a real skip).
  vs previous: `.agents/docs/testing/2026-09-30-verify-v1.13.0-gate-r4.md`
  (`go test -race -count=1 -v ./...` at c63f736), pass 366→366, skip 2→2. No change —
  expected, since this run reviews HEAD with no diff.
- `go run ./cmd/sh-agent install --home .agents/scratch/home` → exit 0. 266 written, 0
  unchanged, 0 removed; installed into the scratch home only (never the real `$HOME`).
- `go run ./cmd/sh-agent status --home .agents/scratch/home` → exit 0.
  `drift none — 265 files match the repo` (one fewer than "written" because `state.json`
  itself is the manifest, not a managed file — consistent with the tool's own accounting).
  `repo has 28 skills, 8 agents, 0 commands, 0 rules`.
- `go run ./cmd/sh-agent list` → exit 0. Prints 28 skills and 8 agents; counts match
  `status`'s report and AGENTS.md's own "Installed" line (28 skills, 8 subagents).
- `govulncheck ./...` → exit 0. "No vulnerabilities found" in code or imported packages;
  1 vulnerability in a required-but-uncalled module — consistent with AGENTS.md's noted
  known finding (GO-2026-5024, x/sys via the SQLite driver, not reachable; ADR-045a).

## Root cause  (FAIL only)
n/a — no FAIL.

## Coverage of changed code
No diff in scope (whole-repo review at HEAD, no changes). Mapping every non-test `.go`
file to its covering test file (by package; `go test -race ./...` above exercises all of
them and passed):

- Root package (`embed.go`) — `embed_test.go`, `plan_test.go`
- `cmd/sh-agent/main.go` — `cmd/sh-agent/main_test.go` (37 test funcs)
- `internal/harness/adopt.go`, `apply.go`, `frontmatter.go`, `inventory.go` — `internal/harness/harness_test.go`
- `internal/harness/derive.go` — `internal/harness/harness_test.go`
- `internal/harness/pathguard.go` — `internal/harness/pathguard_test.go`
- `internal/harness/plan.go` — `internal/harness/harness_test.go` (also root `plan_test.go`)
- `internal/harness/release.go` — `internal/harness/release_test.go`
- `internal/harness/render.go` — `internal/harness/size_test.go` (`TestRenderedContextFitsAntigravityBlock`), `internal/harness/harness_test.go`
- `internal/harness/repo.go` — `internal/harness/harness_test.go`, `internal/harness/audit_test.go`
- `internal/harness/report.go` — `internal/harness/report_test.go`
- `internal/harness/source.go` — `internal/harness/source_test.go`
- `internal/harness/state.go` — `internal/harness/pathguard_test.go` (state-load tests), `internal/harness/harness_test.go`
- `internal/harness/zip.go` — `internal/harness/zip_test.go`
- `internal/harness/owned_unix.go` — `internal/harness/owned_unix_test.go` (build tag `unix`)
- `internal/harness/owned_other.go` — build tag `!unix`; not compiled or tested on this
  Linux box. NOT COVERED on this platform by construction — no macOS/Windows runner here
  to exercise it. Named, not fixed.
- `internal/search/parse.go`, `reconcile.go`, `search.go`, `store.go` — `internal/search/search_test.go` (24 funcs), `internal/search/global_test.go`, `internal/search/security_test.go`, `internal/search/recall_test.go`

Suppressions/skips found in the tree (grep over `*_test.go`, no diff to scope this to):
did not run a repo-wide suppression grep this pass — see Not examined. The two `t.Skip`
call sites seen in the go-test output both carry a reason string (see Commands section)
and are pre-existing, not introduced by any change (no diff in scope).

## Not examined
- Repo-wide grep for added suppressions (`//nolint`, `t.Skip` without reason, `.only(`,
  `@ts-ignore`, etc.) was not run as a dedicated pass, since the brief scopes this to
  "whole repository at HEAD, no diff" and there is no diff to check "added" against;
  the two skips actually encountered in the test run were checked by hand.
- Sandbox and real-home CLI smoke tiers (`.agents/docs/testing/2026-09-03-verify-sandbox.md`,
  `2026-09-03-verify-real-home.md`) — explicitly excluded by the dispatch brief, not run.
- `internal/harness/owned_other.go` (the `!unix` build) — cannot be compiled or tested
  from this Linux box.
- `sh-agent-matrix.md` (§3c) and the search recall/perf tier (§3d,
  `SH_AGENT_RECALL=1 go test -run TestRecall`, rebuild/warm-search timing) — not run;
  no change to `internal/search` or the matrix-relevant verbs is in scope here.
- `sh-agent report --dry-run` (§3b) — not run; out of scope for this dispatch (no change
  to report.go, and the guide reserves a real end-to-end dry-run check for the main agent).
- Frontend/`web`/`ui` tier — repo has no such directory; not applicable.
- Container tier — no Dockerfile changed (no diff in scope at all); not applicable.

## Environment
OS: Linux 6.12.107+deb13-amd64. Go: go1.27.1 linux/amd64 (user-level SDK ahead of
`/usr/bin/go`, as AGENTS.md documents). `govulncheck` from `~/.go/bin/govulncheck`.
`git status --short` was clean at dispatch start (HEAD 3478a1c); by the time this log was
written, three untracked files had appeared from a concurrent dispatch in the same
session (`.agents/docs/reviews/2026-09-30-review-project-code.md`,
`.agents/docs/reviews/2026-09-30-review-project-docs.md`, and this log itself) — none of
them source, and `go run ./cmd/sh-agent status` reported the tree as
`3478a1c, uncommitted changes` / binary `3478a1ccc1ae+dirty` for that reason; this does
not affect `go test`, `go vet`, or `gofmt`, which only read `cmd/` and `internal/`.
