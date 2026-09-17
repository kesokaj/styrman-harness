---
owner: verifier
date: 2026-09-17
---
# Verification — preflight v1.1.0   (2026-09-17)

Range reviewed: v1.0.0..086deab

## Status: PASS

Scope: complete — every tier the guide names for this release (go checks, harness
self-check, scratch-home install/drift, release-shaped embedded install/drift and
file-set comparison) ran; govulncheck ran clean. `testing/guide.md` exists (cited,
not rewritten — no steps changed). Sandbox CLI starts (guide §4) and real-home smoke
checks (guide §5) were not in scope of this dispatch and are listed under Not examined.

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 -v ./...` → exit 0. `ok github.com/kesokaj/styrman-harness/internal/harness 2.671s`, `PASS`.
   counts: pass 184 (99 top-level `--- PASS:` + 85 subtest `    --- PASS:`), fail 0,
   skip 0 at 086deab for `go test -race -count=1 -v ./...`.
   vs previous: `.agents/docs/testing/2026-09-17-release-binary-phase-4-after-fixes.md`,
   pass 184→184 (99+85 both runs), skip 0→0. Exact reproduction, no regression.
4. `govulncheck ./...` → exit 0. "No vulnerabilities found." (go1.27.1 toolchain, so no
   stdlib advisories expected per AGENTS.md notes).
5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-v110-home` → exit 0,
   "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
   `go run ./cmd/sh-agent status --home .agents/scratch/verify-v110-home` → exit 0,
   "drift none — 259 files match the repo" (259 vs 260 written is the manifest's own
   state.json, expected — not a discrepancy).
6. `go run ./cmd/sh-agent list` → exit 0. `go run ./cmd/sh-agent render` → exit 0
   (93 lines of output; both parse every skill and agent without error).
7. Release-shaped embedded install:
   - `CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=v0.0.0" -o
     .agents/scratch/verify-v110-bin/sh-agent ./cmd/sh-agent` → exit 0.
   - Run from cwd `/`, `SH_AGENT_REPO` unset:
     `install --home .agents/scratch/verify-v110-embed` → exit 0, "release binary
     v0.0.0: installing the harness built into it and the running binary", "260
     written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
   - `status --home .agents/scratch/verify-v110-embed` → exit 0, `source embedded
     v0.0.0 (release binary; update downloads the newest release)`, "drift none —
     259 files match the repo". Matches the expected source line and drift-none from
     the dispatch exactly.
   - File-set comparison: `find .claude .gemini -type f | sort` on both the step-5
     checkout install and the step-7 embedded install → 260 files each, `diff` exit 0
     (byte-identical path lists).

## Root cause (FAIL only)

N/A — Status is PASS.

## Coverage of changed code

Not separately mapped file-by-file: this is a release-gate verification over 38
commits / 148 files (v1.0.0..HEAD), not a single-change review. The full `go test
-race` suite (184 sub/top-level tests, 0 fail, 0 skip) covers `cmd/sh-agent` and
`internal/harness`, which is the entirety of this repo's Go code; `context/`,
`agents/`, `skills/`, `commands/`, `rules/` are exercised indirectly via `install`,
`status`, `list`, and `render` succeeding against all 27 skills and 7 agents. No
added `//nolint`, `t.Skip`, `.skip(`/`.only(` suppressions were observed in the test
run (0 skip, 0 fail) — a diff-level suppression grep over the 148-file range was not
run in this dispatch (see Not examined).

## Environment note (not a failure)

`sh-agent status` reported `binary 086deab42962+dirty` and `repo (main @ 086deab,
uncommitted changes)` on both the checkout and embedded installs. Checked
`git status --short --untracked-files=all` at that point: two untracked doc files —
this verifier's own report-in-progress
(`.agents/docs/testing/2026-09-17-verify-preflight-v1.1.0.md`) and a concurrently
written `.agents/docs/reviews/2026-09-17-security-preflight-v1.1.0.md` from a
parallel security dispatch. HEAD 086deab itself was reported clean by the dispatcher
before this run started. This is an artifact of running verification tooling
concurrently in the live tree, not a defect in the release code; the actual release
commit and tree were clean at freeze time. No source, config, or test file was dirty.

## Not examined

- Guide §4 (sandbox start of each CLI, grep for parse errors) and §5 (real-home smoke
  checks) — out of scope for this dispatch's tool-call budget; not run.
- The `report` tier (§3b) — not exercised; dispatch's hard limits reserve
  `report --dry-run` and this release's diff was not specifically flagged as
  touching `internal/harness/report.go` by the dispatch, so it was skipped to stay
  in budget.
- Plain-copy / read-only mirror install tier (§3a) — not run; would only be required
  if this range touched `IsCheckout`/rebuild-decision logic, which was not confirmed
  either way in this dispatch.
- A per-file suppression grep (`//nolint`, `t.Skip`, `.only(`, etc.) across the full
  38-commit / 148-file range was not performed.
- `sh-agent update` and `sh-agent report` (non-dry-run) were correctly never run, per
  hard limits.

## Environment

OS: Linux 6.12.107+deb13-amd64. Go: go1.27.1 linux/amd64 (user-level SDK ahead of
/usr/bin/go per AGENTS.md). No Docker/Colima involved (no Dockerfile in this repo).
Network: not exercised; embedded install ran with `SH_AGENT_REPO` unset and cwd `/`,
no `gh` invocation triggered.
