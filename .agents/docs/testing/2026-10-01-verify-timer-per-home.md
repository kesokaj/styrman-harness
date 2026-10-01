---
owner: verifier
date: 2026-10-01
---
# Verification — timer-per-home   (2026-10-01)

Range reviewed: 96a7e48..working tree (actual HEAD is 0971e2a, a records-only commit on top of 96a7e48; the six named paths are the modified ones)

## Status: PASS

Scope: complete for the dispatch's seven steps. Guide .agents/docs/testing/guide.md exists (modified, uncommitted; not rewritten by this run). Raw output in .agents/scratch/verify-timer/.

## Commands
- 1. `gofmt -l cmd internal` -> exit 0, printed nothing.
- 2. `go vet ./...` -> exit 0.
- 4. `go test -race -count=1 -v ./...` -> exit 0; counts: pass 520, fail 0, skip 3 at 0971e2a+dirty for go test -race -count=1 -v ./...; vs previous: none compared (the gate logs ran other command shapes; not read). Skips: TestAgyCorpusLeaves, TestRecall, TestHelperIndex (pre-existing conditional skips, not in the diff). Packages: root ok 1.1s, cmd/sh-agent ok 15.8s, internal/dream ok 7.4s, internal/harness ok 3.8s, internal/memory ok 1.2s, internal/search ok 7.8s.
- 4b. `go test ./internal/dream -run 'TestUnitNamesPerHome|TestInstallTimerLinux|TestRemoveTimerLinux|TestTimerMacOS|TestStatusLine' -v` -> ok; PASS for all five.
- 4c. `go test ./cmd/sh-agent -run TestUninstallDryRunLeavesTheTimer -v` -> ok; PASS.
- 5. Must-fail control: copied timer.go to scratch, changed UnitFor line 37 to `if true {` (returns the bare names for every home), `go test ./internal/dream -run TestUnitNamesPerHome` -> FAIL: `timer_test.go:225: other home: sh-agent-dream com.styrman.sh-agent-dream`. File restored (cmp identical to the copy), `git diff --stat` identical to the start, test passes again (ok).
- 6. `go run ./cmd/sh-agent install --home .agents/scratch/verify-timer-home` -> exit 0 (266 written). `go run ./cmd/sh-agent status --home ...` -> exit 0: `binary 0971e2aaf99f+dirty` (no warning), `drift none - 265 files match the repo`, `dream never run; timer not installed`.
- 7. Suppressions: grep of added lines in the Go diff for nolint, lint:ignore, t.Skip, go:build ignore, .only( -> none.

## Coverage of changed code
- internal/dream/timer.go (UnitFor, linuxServicePath, plist/unit use) - covered by internal/dream/timer_test.go. TestUnitNamesPerHome asserts (a) suffix `sh-agent-dream-<8 hex>` and plist `com.styrman.<unit>` for a non-real home, and the bare names for the real home; (b) determinism; (c) two TempDirs give distinct names; (d) InstallTimer's systemctl calls are exactly `--user daemon-reload` and `--user enable --now <suffixed>.timer`, and the suffixed .service file exists. Fails without the fix (step 5).
- internal/dream/dream.go - exercised by internal/dream/run_test.go and the dream package run; I did not read its diff line by line.
- cmd/sh-agent uninstall path - TestUninstallDryRunLeavesTheTimer (cmd/sh-agent/memory_test.go) passes.
- Not shown by a test I ran: the macOS launchctl label for a non-real home beyond the string assertion in TestUnitNamesPerHome (TestTimerMacOS passes against the fake manager).

## Not examined
- Real systemd/launchd behaviour (fake managers only); no real-home install; no Antigravity run.
- Whether timer.go's `--remove-timer`/`uninstall` paths against a real shared user manager leave the operator's timer alone (asserted only through fakes).
- Diffs of dream.go and run_test.go were not read; guide.md diff not read.
- Docs/other agents' artifacts.

## Tree check
Source, test and config paths at the end equal those at the start (the six modified files). Paths seen new by name during the run, not mine: modified .agents/docs/bugs/INDEX.md, new .agents/docs/bugs/BUG-144-dream-timer-unit-not-scoped-by-home.md (another agent's records). Scratch home .agents/scratch/verify-timer-home is gitignored.

## Environment
Linux 6.12.107+deb13-amd64, Go from PATH (version not recorded), no Docker/Colima use.
