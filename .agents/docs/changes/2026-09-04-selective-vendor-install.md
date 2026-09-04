---
owner: main
date: 2026-09-04
---
# Selective vendor installation support via `--vendor`     (date: 2026-09-04)

## Summary

Added `--vendor <all|antigravity|claude>` flag to `sh-agent` (`install`, `status`, `diff`, `update`) with `agy` alias for `antigravity`:
- `Plan` takes a target vendor and filters planned files accordingly: `antigravity` produces only `~/.gemini/config/...` and `~/.local/bin/sh-agent`; `claude` produces only `~/.claude/...` and `~/.local/bin/sh-agent`; `all` produces both.
- `State` saves the active vendor (`st.Vendor`) so daily commands (`status`, `diff`, `update`) automatically operate on the installed vendor without needing to re-pass the flag.
- Default remains `all` when unconfigured to preserve backwards compatibility.
- Switching vendors (e.g. from `all` to `antigravity`) automatically prunes the unselected vendor's files via existing resync logic.
- Documented in `ADR-012`, `deploy/local.md`, and `README.md`.

## Type

feature

## Files changed

`cmd/sh-agent/main.go`, `internal/harness/plan.go`, `internal/harness/state.go`, `internal/harness/apply.go`, `internal/harness/harness_test.go`, `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/deploy/local.md`, `README.md`

## Verification

- Unit & race tests pass: `go test -race ./...` (including new `TestVendorSelection`).
- Scratch home install (`/tmp/sh-agent-test-agy`): 128 files installed under `~/.gemini/config/` and `~/.local/bin/`, `.claude` does not exist, `status` reports drift none.
- Scratch home vendor switch: switched to `all` (256 files), switched back to `antigravity` (128 removed, 128 unchanged), `uninstall` removes all managed files cleanly.
