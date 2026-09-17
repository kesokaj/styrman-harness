---
owner: verifier
date: 2026-09-17
---
# Verification — release-binary Phase 2 checkpoint (2026-09-17)

Range reviewed: working tree at HEAD 21ba8dc (branch release-binary, clean at dispatch time; 3
untracked docs artifacts appeared during this run from concurrently dispatched agents —
`.agents/docs/reviews/2026-09-17-review-f7ee817-21ba8dc-code.md`,
`.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`, and this log itself — no tracked
file was modified by this run)

Followed `.agents/docs/testing/guide.md` § 1 and § 3, plus the plan's own Phase 2 task-4
integration tier (`.agents/docs/plans/2026-09-17-release-binary.md`).

## Status: PASS

Scope: complete for the five command groups named in the dispatch (Go checks, race test with
counts, checkout-source scratch-home install/status, release-shaped embedded-source
install/status/update, file-set comparison). guide.md exists and needed no update. Not run:
guide.md § 2 (harness self-check: `list`, `render`), § 3a (plain-copy/mirror tier — this change
does not touch `IsCheckout`/rebuild logic), § 3b (`report` — untouched by Phase 2), § 4/§ 5
(sandbox/real-home CLI smoke checks — out of scope for this diff-free working-tree check),
skill-matrix. Phase 2's own remaining checkpoint items (branch commit/push) are not this
verifier's job.

## Commands

- `gofmt -l cmd internal *.go` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race -count=1 -v ./...` → exit 0; output saved to
  `.agents/scratch/verify-phase2/go-test.txt`; counts: pass 81, fail 0, skip 0 at 21ba8dc for
  `go test -race -count=1 -v ./...`; package summary: `ok github.com/kesokaj/styrman-harness
  1.086s`, `ok .../cmd/sh-agent 3.555s`, `ok .../internal/harness 1.816s`. vs previous: none
  found in `testing/` with a counts line for this exact command (earlier logs in this repo predate
  the counts-line convention or cover different commands) — writing this as the first baseline.
- `go test -race -count=1 ./...` (non-verbose, per dispatch item 3) → exit 0, all three packages
  `ok`.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase2-home` → exit 0; "260
  written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase2-home` → exit 0; `repo
  /home/kesokaj/workspace/styrman-harness (release-binary @ 21ba8dc, uncommitted changes)`;
  `drift none — 259 files match the repo`. Note: "uncommitted changes" reflects the untracked
  docs artifacts named above (present when this status ran), not a modified tracked file — `git
  status --short` at every point in this run showed only untracked additions, never a modified
  tracked path.
- `CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=v0.0.0" -o
  .agents/scratch/verify-phase2-bin/sh-agent ./cmd/sh-agent` → exit 0.
- From `/`, `SH_AGENT_REPO` unset, using the built binary:
  - `install --home .agents/scratch/verify-phase2-embed` → exit 0; first line "release binary
    v0.0.0: installing the harness built into it and the running binary"; "260 written, 0
    unchanged, 0 removed, 0 kept — 260 files managed".
  - `status --home .agents/scratch/verify-phase2-embed` → exit 0; `source     embedded v0.0.0
    (release binary; update downloads the newest release)`; `installed ... from release v0.0.0`;
    `drift none — 259 files match the repo`.
  - `update --home .agents/scratch/verify-phase2-embed` → exit 1 (expected refusal): "this install
    is the release binary v0.0.0 and has no checkout; download a newer release by hand (README §
    Install without cloning)".
- File-set comparison: `find .../.claude .../.gemini -type f` (paths relativized), sorted, diffed
  between `verify-phase2-home` (checkout source) and `verify-phase2-embed` (embedded/release
  source) → empty diff, exit 0. The two installs produce an identical file set under both vendor
  directories.

## Root cause

N/A — no failures.

## Coverage of changed code

Phase 2's changed surfaces (per the plan's task 4 evidence: `internal/harness/source.go`
`ResolveSource`, `EmbeddedSource`, `IsReleaseVersion`, `EmbeddedPrefix`; `Install` taking a
`Source`; `BuildRecord`'s `harness_version`; `main.go`'s `version`, per-command embedded-source
behaviour, `releaseFromInfo`) are exercised by `TestResolveSource`,
`TestResolveSourceRefusesNamedNonRepo`, `TestEmbeddedInstallStateAndRecord`,
`TestReleaseBinaryWithoutCheckout`, `TestReleaseFromInfo` (all present and passing in this run's
`go test` output) and independently reproduced end-to-end here via the release-shaped build. No
changed file in this working tree lacks a corresponding test — `git diff` against `main` was not
re-examined here since the plan's own task evidence already lists the tests and mutation
must-fails for each; this run's job was to reproduce the integration tier fresh, which it did.

Suppressions added: none observed in the files touched by Phase 2 (no new `//nolint`,
`t.Skip`, `@ts-ignore`, etc. encountered while reading test output or running commands); a full
diff-line grep for suppressions was not performed in this run (see Not examined).

## Not examined

- guide.md § 2 (`sh-agent list`, `sh-agent render | wc -l`) — not run this pass.
- guide.md § 3a (plain-copy/read-only mirror tier) — Phase 2 does not touch `IsCheckout`/rebuild
  logic; the plan's own task-4 evidence already covers the checkout-vs-embedded comparison this
  tier would otherwise add.
- guide.md § 3b (`report` dry-run and report-specific tests) — `report.go` is untouched until
  Phase 3.
- guide.md § 4/§ 5 (sandbox CLI starts, real-home smoke checks) — no change to `context/`,
  `agents/`, or `skills/` content in this range; out of scope for a Phase 2 code checkpoint.
- A line-by-line `git diff main...release-binary` suppression grep (added `t.Skip`, `//nolint`,
  etc.) was not performed; only the working tree's current test output and the plan's own stated
  mutation-tested evidence were read.
- `govulncheck ./...` (named in `AGENTS.md` § Setup commands) was not run in this pass.
- Whether `.agents/scratch/verify-phase2-home` vs `-embed` differ under directories other than
  `.claude`/`.gemini` (e.g. `.config/sh-agent/state.json` content) was not diffed byte-for-byte
  beyond the file-name/path comparison above; `state.json`'s `"repo"` field difference
  (`/home/.../styrman-harness` vs `embedded:v0.0.0`) is expected per D20/D17 and was seen directly
  in each `status` output, not diffed as files.

## Environment

Linux 6.12.107+deb13-amd64, go1.27.1 (`go version go1.27.1 linux/amd64`), branch release-binary,
HEAD 21ba8dc, working tree clean of tracked-file changes throughout. No Docker/Colima involved
(not applicable to this tier).
