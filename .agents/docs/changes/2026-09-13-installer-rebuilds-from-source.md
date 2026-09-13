---
owner: changelog
date: 2026-09-13
---

# Binary rebuild: from source, versioned, dirty marker   (2026-09-13)

## Summary

Phase 1 of learning-loop plan: `sh-agent install`/`sync` now always rebuild the binary from the repo with `go build` (VCS-stamped) and let that build execute the install, unless already a child build (`SH_AGENT_CHILD=1`) or `go` is unavailable (prints warning, falls back to self-copy). Fixes BUG-006 (stale installed binary). `status` prints a `binary` line showing the running binary's revision and warns when older than the repo's HEAD or unstamped. `state.json` records `<commit>+dirty` when the tree had uncommitted changes at install. `LastRetro` now reads `~/.config/sh-agent/retros.jsonl` mtime. Inventory sentence "Skills marked `/name`..." prints only when at least one skill is slash-only. New unit and integration tests cover the rebuild decision, binary age comparison, and end-to-end install/status/uninstall.

## Type

feature

## Files changed

- cmd/sh-agent/main.go — `install`/`sync` route through `shouldRebuild()` and `harness.Reinstall()`; `status` prints binary line + older-than-repo or no-stamp warning; `binaryOlder()` logic
- cmd/sh-agent/main_test.go — NEW: TestReorder (BUG-002 site), TestResolveVendor, TestShouldRebuild, TestBinaryOlder, TestRunInstallStatusUninstall (end-to-end)
- internal/harness/repo.go — `Reinstall(repo, args, out)` builds with `go build -o <repo>/.agents/scratch/sh-agent-build-<pid>/` and runs the built binary with `SH_AGENT_CHILD=1`; `ChildEnv` constant
- internal/harness/apply.go — `Install` records `Commit: commit + "+dirty"` when `Inspect(repo).Dirty`
- internal/harness/state.go — `LastRetro` reads both `~/.config/sh-agent/proposals/*.md` and `~/.config/sh-agent/retros.jsonl` mtime
- internal/harness/inventory.go — `slashNote()` conditional: prints only when `any(s.SlashOnly)` true
- internal/harness/harness_test.go — `TestLastRetro` extended to verify retros.jsonl read; `TestSlashNoteOnlyWhenSomeSkillIsSlashOnly` added

## Agents dispatched

reviewer (×9, three rounds of three cold readers; reports in `reviews/2026-09-13-review-cold-read-learning-loop-{1,2,3,r2-1,r2-2,r2-3,r3-1,r3-2,r3-3}.md`), changelog (you)

## Verification

- `gofmt -l .` → empty
- `go vet ./...` → clean
- `go test -race -count=1 ./...` → ok: cmd/sh-agent, internal/harness
- Scratch home integration: `go run ./cmd/sh-agent install --home <scratch>` rebuilds with `go build` (output shows build), installed binary `version` prints repo's short hash; `status` prints `binary <hash>+dirty` and `drift none`; `go run ./cmd/sh-agent status` prints `binary (devel) — no build stamp; run sh-agent update`; with `SH_AGENT_CHILD=1` set, skips rebuild; with `go` off PATH, prints warning and self-copies; `--dry-run` passes through; `state.json` commit line is `b64ffe8+dirty`; no orphaned `sh-agent-build-*` directories remain
- E2E: sandbox home `.agents/scratch/home-test`, both CLIs start, no parse errors

## Requirements

R4 (one binary, idempotent install) — `install` always routes through a single rebuild path or self-copy, idempotent on second run.
R6 (development only against `--home` scratch homes) — tests use `t.TempDir()` and pass `--home` explicitly.
