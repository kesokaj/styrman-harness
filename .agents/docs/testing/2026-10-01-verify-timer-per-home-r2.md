---
owner: verifier
date: 2026-10-01
---
# Verification — timer-per-home r2   (2026-10-01)

Range reviewed: 0971e2a..working tree

## Status: PASS

Scope: complete for the dispatch's steps 1-7. Guide: .agents/docs/testing/guide.md exists (cited). Not run: darwin/launchd on a real macOS host (renderers and path logic are unit-tested only).

## Commands
- 1. `gofmt -l cmd internal` -> exit 0, printed nothing.
- 2. `go vet ./...` -> exit 0.
- 3. `go test -race -count=1 -v ./...` (log .agents/scratch/verify-timer-home-r2-logs/go-test-1.txt) -> exit 0; counts: pass 520, fail 0, skip 3 at 0971e2a+dirty for go test -race -count=1 -v ./...; vs previous: .agents/docs/testing/2026-10-01-verify-timer-per-home.md, pass 520->520, skip 3->3 (TestAgyCorpusLeaves, TestRecall, TestHelperIndex; conditional, not in diff). Packages: root ok, cmd/sh-agent ok 14.9s, internal/dream ok 6.3s, internal/harness ok 3.7s, internal/memory ok 1.2s, internal/search ok 7.5s.
- 4. targeted run: PASS TestRenderSystemdUnits, TestRenderLaunchdPlist, TestRemoveTimerLinux, TestTimerMacOS, TestInstallTimerRefusesUnsafePaths, TestUnitNamesPerHome; `ok`.
- 5a. Must-fail: first attempt (only line 204 `unit + ".timer"` -> `unitName + ".timer"`) did NOT compile (`declared and not used: unit`, timer.go:196), so no verdict from it. Second attempt added `; _ = unit` on line 196: compiles (vet 0); `--- FAIL: TestUnitNamesPerHome` ("timer_test.go:250: remove side keeps the bare name"); TestRemoveTimerLinux still PASS (does not catch it alone). Restored, cmp identical.
- 5b. Must-fail: dropped ` --home "%s"` and `home` arg in renderSystemd (timer.go:70-71): compiles (vet 0); `--- FAIL: TestRenderSystemdUnits` and `--- FAIL: TestUnitNamesPerHome`. Restored, cmp identical; `git diff --stat` = 9 files, 132 insertions, 23 deletions, same as start.
- 6. `go run ./cmd/sh-agent install --home .agents/scratch/verify-timer-home-r2` -> exit 0 (266 written); `status --home` -> exit 0: `binary 0971e2aaf99f+dirty` (no warning), `drift none — 265 files match the repo`.
- 7. `git status --short --untracked-files=all` at end: identical to start plus this log. No reviewer report appeared; no source/test/config path moved.

## Coverage of changed code
- internal/dream/timer.go, dream.go — covered by internal/dream/timer_test.go (TestUnitNamesPerHome, TestRenderSystemdUnits, TestRenderLaunchdPlist, TestTimerMacOS, TestRemoveTimerLinux, TestInstallTimerRefusesUnsafePaths) and run_test.go; both mutations were caught.
- cmd/sh-agent/memory_test.go, internal/dream/run_test.go, timer_test.go are tests.
- Docs-only changes under .agents/docs need no test.
- Suppressions: one added `t.Skip("no symlinks here")` in the symlink test (reason on its line; it did not skip here, skip count unchanged at 3). No nolint, .only, ts-ignore.

## Not examined
Real macOS launchd behaviour; real systemd --user enable/disable; the symlink case was only judged through the passing TestUnitNamesPerHome, not mutated.

## Environment
Linux 6.12 (Debian 13), go from PATH (version not recorded), no Docker needed (no Dockerfile changed).
