---
owner: bugs
date: 2026-09-13
id: BUG-009
status: fixed
area: cmd/tests
---
# BUG-009: `TestRunInstallStatusUninstall` asserted uninstall removes `~/.claude` itself

## Symptom
New test in `cmd/sh-agent/main_test.go` failed on its first run: `~/.claude should be gone after uninstall, stat err = <nil>`.

## Repro
`go test ./cmd/sh-agent/...` before the fix; `TestRunInstallStatusUninstall` failed at the assertion checking `~/.claude` no longer exists after `uninstall`.

## Root cause
The assertion contradicted documented behaviour. `pruneEmptyDirs` in `internal/harness/apply.go` deliberately stops at the vendor roots ("uninstall never removes `~/.claude` or `~/.gemini/config` themselves"). The test author (the main agent) asserted the root directory itself was gone without reading that stop list first.

## Fix
Assert on the manifest's paths under the roots (`~/.claude/AGENTS.md`, `~/.claude/skills`, `~/.gemini/config/agents`, the binary) rather than the roots themselves. `cmd/sh-agent/main_test.go`, second run green.

## Prevention
Before asserting that something is removed, read the code's documented stop/keep list first (here `pruneEmptyDirs`'s comment) and assert on what it actually says it removes, not on an assumption. A test that fails on its first run for a reason other than the change under test is itself a recorded failure, per the global context's "every failure gets a number" rule.

## Related
BUG-006 (same session, installer rebuild work); global context rule on recording every dispatched failure.
