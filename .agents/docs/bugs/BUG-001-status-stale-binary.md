---
owner: bugs
date: 2026-09-03
id: BUG-001
status: fixed
area: installer/status
---
# BUG-001: `status` reported the installed binary as stale right after `install`

## Symptom
`sh-agent status --home <scratch>` immediately after a clean install printed `stale ~/.local/bin/sh-agent` and exited 1.

## Repro
`go run ./cmd/sh-agent install --home /tmp/x && go run ./cmd/sh-agent status --home /tmp/x` on the first e2e run of the rebuild.

## Root cause
`status` builds its plan with `selfExe == nil` (the binary's bytes depend on the toolchain, not the repo), so the binary was absent from the plan while the manifest listed it; the stale check flagged the difference.

## Fix
`internal/harness/apply.go`: `Status` skips `BinPath` in the stale loop. Regression test in `TestInstallStatusResyncUninstall`.

## Prevention
A path that install writes but status cannot compare must be excluded from drift explicitly, and the exclusion needs a test. Before adding a new managed path, decide how status compares it.

## Related
`.agents/docs/testing/2026-09-03-verify-rebuild.md`
