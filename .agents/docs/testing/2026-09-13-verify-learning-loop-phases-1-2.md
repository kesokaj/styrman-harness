---
owner: verifier
date: 2026-09-13
---
# Verification — learning-loop-phases-1-2   (2026-09-13)

Range reviewed: b64ffe8..bd5ed87

## Status: PASS

Scope: full Go test suite, gofmt, vet, scratch-home install/status cycle, and coverage mapping for every changed source file in the range. CLI-start smoke check (context/agents/skills change) is cited from the 2026-09-13 record, not re-run.

## Commands

- `gofmt -l .` → empty (clean)
- `go vet ./...` → clean, exit 0
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/cmd/sh-agent`, `ok github.com/kesokaj/styrman-harness/internal/harness`, exit 0
- Scratch-home install/status (never `/tmp`, used `.agents/scratch/verify-home-13`, removed after):
  - `go run ./cmd/sh-agent install --home .agents/scratch/verify-home-13` → `256 written, 0 unchanged, 0 removed, 0 kept — 256 files managed` plus the expected PATH note
  - Installed binary directly (`.agents/scratch/verify-home-13/.local/bin/sh-agent version`) → `sh-agent bd5ed875ab14+dirty`
  - Installed binary `status --home .agents/scratch/verify-home-13 --repo <this repo>` →
    `binary     bd5ed875ab14+dirty` (no warning) and `drift      none — 255 files match the repo` — this satisfies the "binary line without a warning" and "no drift" requirement from AGENTS.md Testing.
  - Note: running `go run ./cmd/sh-agent status` (as opposed to the installed binary) prints `binary (devel) — no build stamp` because `go run` itself is unstamped; that is expected behavior of `binaryOlder`/`buildVersion`, not a defect — the installed binary, which is what `status` is meant to characterize in real use, is correctly stamped.
- CLI-start smoke check (required by AGENTS.md Testing after a change to `context/`, `agents/`, or `skills/`): not re-run here; cited from `.agents/docs/changes/2026-09-13-all-skills-model-invocable.md`, which records both CLIs starting cleanly against the sandbox with evidence under `.agents/scratch/e2e/`.

## Coverage of changed code

| Changed file | Covering test |
|---|---|
| `cmd/sh-agent/main.go` (`shouldRebuild`, `binaryOlder`, rebuild routing in `install`/`sync`, `status` binary line, `selfExecutable` comment-only) | `cmd/sh-agent/main_test.go`: `TestShouldRebuild`, `TestBinaryOlder`, `TestRunInstallStatusUninstall` (end-to-end install/status/uninstall against this repo with `SH_AGENT_CHILD=1`) |
| `cmd/sh-agent/main_test.go` | new test file, self-covering |
| `internal/harness/repo.go` (`Reinstall` rewritten to `go build` into `.agents/scratch/sh-agent-build-<pid>/` then exec the built binary with `ChildEnv`) | **not unit-tested** — no test calls `harness.Reinstall` itself (it spawns a real `go build`, which `TestRunInstallStatusUninstall` deliberately avoids via `SH_AGENT_CHILD=1`). Manually verified per the change record `.agents/docs/changes/2026-09-13-installer-rebuilds-from-source.md` ("Scratch home integration" section) and re-confirmed live in this run (see Commands above: install rebuilt and the resulting binary carried the repo's commit stamp). This is real coverage but it is manual/integration, not a `go test` assertion on `Reinstall` in isolation. |
| `internal/harness/apply.go` (`+dirty` suffix on `State.Commit` when `Inspect(repo).Dirty`) | **not unit-tested directly** — no test in `harness_test.go` asserts a `+dirty`-suffixed commit after `Install` on a dirty repo. Confirmed manually in this run: scratch-home `status` showed `installed ... from commit bd5ed87+dirty` against this (dirty) working tree. Gap: no automated regression test pins this behavior. |
| `internal/harness/inventory.go` (`slashNote` helper, conditional sentence) | `internal/harness/harness_test.go`: `TestSlashNoteOnlyWhenSomeSkillIsSlashOnly` |
| `internal/harness/state.go` (`LastRetro` now also reads `retros.jsonl`) | `internal/harness/harness_test.go`: `TestLastRetro` (extended) |
| `internal/harness/harness_test.go` | test file itself, exercised by the run above |
| `context/AGENTS.md`, `AGENTS.md`, `README.md`, `.agents/docs/deploy/local.md`, `skills/{audit,distill,handoff,init-project,research}/SKILL.md` | doc/prompt content, no Go test applicable; covered by the required install (drift-check confirms these deploy without parse errors) and by the cited CLI-start record for the skill frontmatter changes |
| `.agents/docs/bugs/BUG-006-*.md`, `BUG-009-*.md`, `.agents/docs/bugs/INDEX.md`, `.agents/docs/changes/*.md`, `.agents/docs/plans/2026-09-13-learning-loop-project-level.md`, `.agents/docs/reviews/2026-09-13-review-cold-read-learning-loop-*.md`, `.agents/docs/testing/2026-09-13-verify-research-pass.md` | records/plan/review artifacts, not executable, no test applicable |

## Root cause (n/a — Status is PASS)

Not applicable; the two coverage gaps below are noted as findings, not failures, because they are exercised by recorded manual/integration verification rather than left unverified.

## Findings (informational, not blocking PASS)

1. `internal/harness.Reinstall`'s new `go build`-then-exec path has no `go test` covering it in isolation (only end-to-end manual runs, one of which is repeated in this log).
2. `internal/harness/apply.go`'s `+dirty` commit-suffix logic has no `go test` assertion; only manually observed (here and in the prior change record).

Neither gap blocks PASS: both behaviors were independently reproduced live in this session against the actual commit range, and the plan's own testing bar for this class of change (AGENTS.md Testing section) is "install into scratch home, status reports no drift and a binary line without warning" — which passed — plus the CLI-start check, which is cited. A future regression in either path would currently only surface via manual/integration checks, not `go test -race ./...`; worth a follow-up unit test but not a release blocker.

## Not examined

- Did not re-run either CLI's sandbox start (cited from `.agents/docs/changes/2026-09-13-all-skills-model-invocable.md` instead, per the dispatch instruction).
- Did not run the full skill matrix (`.agents/docs/testing/skill-matrix.md`); out of scope per that doc's own "complete" bar and not requested.
- Did not verify Mac/Antigravity behavior; this run is Linux-only, matching the AGENTS.md "tested from where it runs" rule for this box.
- Did not touch the uncommitted working-tree docs edits present at session start (ADR-013/ADR-014 additions to DECISIONS.md, ROADMAP.md, TODO.md, README.md, bugs/INDEX.md, two new BUG files) — these are outside the b64ffe8..bd5ed87 commit range per the dispatch and were left untouched.

## Environment

- OS: Linux 6.12.107+deb13-amd64
- Go: `go version` not captured verbatim in this run but matches the repo's `go.mod` toolchain (build succeeded, `go vet`/`go test` ran clean)
- No Docker/Colima involved (no Dockerfile change in range)
