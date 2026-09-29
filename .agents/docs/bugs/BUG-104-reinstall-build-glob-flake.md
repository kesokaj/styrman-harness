---
owner: bugs
date: 2026-09-29
id: BUG-104
status: fixed
area: internal/harness (tests)
---
# BUG-104: `TestReinstallBuildsAndInstalls` flaked under a concurrent `go test`

## Symptom
During the v1.12.0 gate (2026-09-29) the release reviewer's `go test ./...`
ran while the gate verifier's suite also ran; `TestReinstallBuildsAndInstalls`
failed once with "build directories left behind" and passed when run alone.

## Repro
Run two `go test ./internal/harness/...` invocations concurrently on the
same checkout; one asserting the glob may fail depending on scheduling.
Not independently re-reproduced beyond the gate's observed failure — the
mutation evidence below confirms the mechanism instead.

## Root cause
`internal/harness/harness_test.go` asserted that
`filepath.Glob(<repo>/.agents/scratch/sh-agent-build-*)` was empty after
`Reinstall`. `Reinstall` (`internal/harness/repo.go:218`) names its build
directory `sh-agent-build-<pid>` under the shared repository scratch
directory, so the glob also matches another process's in-flight build
directory when two `go test` runs (or a test binary and a real install)
share the same checkout concurrently. The test's own cleanup was correct;
the assertion's scope was too wide for a directory other processes write to.

## Fix
`harness_test.go:958` now stats only its own path,
`sh-agent-build-<os.Getpid()>`, instead of globbing `sh-agent-build-*`.
Not yet released — lands in v1.12.1.

## Prevention
A test that asserts a shared directory is clean names the exact path its
own run created, never a glob over a directory other processes also write
to.

## Related
Mutation evidence: with the test's `defer os.RemoveAll(dir)` commented out
(sole occurrence, restored after, `git diff` clean), the mutant compiled
and the named test failed with "build directory left behind:
…/sh-agent-build-504536"; restored, it passes.
