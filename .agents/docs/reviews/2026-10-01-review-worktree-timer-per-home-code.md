---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-timer-per-home-code   (2026-10-01)

Range reviewed: 96a7e48..working tree (internal/dream, cmd/sh-agent)

## Verdict: APPROVE WITH NITS

The naming change is sound: install, remove, status and the unit/plist paths all derive from one `UnitFor(home)`, the real home keeps the bare names (a v2.0.0 timer is still found and removed), and `--home` is made absolute in `main.go:236` before it reaches `dream`. `go test -race ./internal/dream ./cmd/sh-agent` passes. The worst problem is a test gap: the security fix on the REMOVE side (the original finding's `systemctl disable --now` on the real unit) has no test that fails if `RemoveTimer` goes back to the bare name. A second, residual issue: a scratch macOS install still loads a job that runs `dream --all` against the real home.

## Findings
- **[medium / confidence: high]** `internal/dream/timer_test.go:135-155`, `cmd/sh-agent/memory_test.go:106-136` — no test fails if `RemoveTimer` (timer.go:186-194) reverts to the bare `unitName` for a non-real home. `TestRemoveTimerLinux` and `TestTimerMacOS` use `timerHome`, which maps the home to the bare name; `TestUnitNamesPerHome` checks only the Linux INSTALL calls. `TestUninstallDryRunLeavesTheTimer` runs a real uninstall on a scratch home with a recording `systemctl`, but never reads `calls` afterwards, so it would pass with `disable --now sh-agent-dream.timer`, the exact call the security finding named. Same for the launchd label: no test reads the plist label for a non-real home (`renderLaunchd(plistName, ...)` would pass). Fix: in `TestUnitNamesPerHome` also call `RemoveTimer(other, "linux")` and assert `--user disable --now <suffixed>.timer`; assert the darwin plist for `other` contains `<string>com.styrman.sh-agent-dream-XXXXXXXX</string>`. By the project rule a fix counts only when a test fails without it. (I reasoned this from the tests' code; I did not run a mutation, since I edit no source.)
- **[medium / confidence: medium]** `internal/dream/timer.go:55-60,135-138` — the rest of the original security finding remains. On macOS, `dream --install-timer --home <scratch>` still runs `launchctl bootstrap gui/<uid> <scratch plist>`, and the plist's ProgramArguments are `<scratch>/.local/bin/sh-agent dream --all` with no `--home`, so the nightly job dreams over the REAL home's projects and sends their content to the model. The new suffix stops the name collision but arms a second, separate real job instead. On Linux `systemctl --user enable <name>` of a unit in a scratch directory normally fails (unit not on the manager's search path), so it is mostly harmless there, but the same no-`--home` ExecStart is written. Fix: either write the unit/plist with `--home <home>` when home is not the real one, or skip the manager calls and say so for a non-real home (the remediation the finding gave). If the operator accepts it as residual, say so in the changelog; the current comment on `UnitFor` claims the scratch install "never touches the operator's own timer", which is true of the name only. I could not run launchd here.
- **[low / confidence: medium]** `internal/dream/timer.go:36-38` — the real-home test is a lexical `filepath.Clean` compare. If `--home` is a symlink to (or alias of) `os.UserHomeDir()` the home gets a suffixed name: install registers a second timer that also runs `dream --all` nightly beside any bare-name one, and uninstall/remove from the other spelling does not find the first. Not an orphaning of a v2.0.0 timer under the normal spelling. Fix: compare with `filepath.EvalSymlinks` on both sides, or `os.SameFile` as `harness/pathguard.go:113-121` does; fall back to Clean when the path does not exist.
- **[low / confidence: low]** `internal/dream/timer.go:36` — if `os.UserHomeDir()` errors, every home, including the real one, gets a suffixed name. Unreachable from the CLI (`main.go:225-229` fails first on the default home), but a `--home` run with `$HOME` unset would suffix the real home's name. Only note if the case is expected to occur.

## Answers to the brief's points
- Relative `home`: `main.go:236` runs `filepath.Abs` before any dream verb, so the hash is of an absolute path; callers of `dream` other than `main` pass absolute homes. No collision with the real name. A scratch home cannot map to the bare name unless it equals the real home.
- Symlinked real home: gets a suffixed name (low finding above); the bare-name timer is not removed by that spelling.
- Install/remove agreement: `TimerInstalled`, `InstallTimer`, `RemoveTimer` and the paths all go through `UnitFor(home)` (`dream.go:219-226`, `timer.go:44-47,125,186`); they agree. The `systemctl` unit argument and the file names use the same `unit` value.
- `userHomeDir` seam: acceptable. It is a package variable restored by `t.Cleanup`, which the `golang` checklist's "no package-level mutable state" rule disfavours, but it is the smallest seam, and tests that set it cannot use `t.Parallel()` (none in the diff do). Not a finding.
- v2.0.0 upgrade: a real-home timer under the bare name is still found by `TimerInstalled` and removed by `RemoveTimer` (real home -> `unitName`/`plistName`). Confirmed by `TestUnitNamesPerHome`'s first assertion and `TestRemoveTimerLinux`.
- Eight-hex suffix: `sh-agent-dream-<8 hex>` uses only characters valid in a systemd unit name; `com.styrman.sh-agent-dream-<8 hex>` is a valid launchd label. 32 bits give a collision risk of no concern for a handful of scratch homes (a collision would only make two scratch homes share a name).

## Criteria
1. Dead code — pass. No unused symbols; `unitName`/`plistName` remain used.
2. Flow tracing — pass with the two findings: names agree end to end; Remove has no test pinning the per-home name; the macOS service still runs against the real home.
3. Standards conformance — pass. Errors in the diff are unchanged; doc comments on `UnitFor` and `RemoveTimer` match, except the "never touches the operator's own timer" claim noted above. I did not run the restatement grep; docs naming `sh-agent-dream.timer` (plan 2026-09-30, deploy docs) were not checked for a per-home note (see Not examined).
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass. `sha256` over a path for a name suffix is an identifier, not crypto use.
6. Decoupling — pass. `cmd/sh-agent` test imports the exported `dream.UnitFor`, not an unexported helper.

## Not examined
Docs and specs that describe the unit names (`.agents/docs/specs`, `deploy/`, README) for a per-home-name note: out of scope (diff only). No mutation was run; no macOS run. The `gofmt`/`go vet` gate was not run by me (only `go test -race` on the two packages).

## Round 2

Range reviewed: 96a7e48..working tree, round 2

`go test -race ./internal/dream ./cmd/sh-agent` passes. Re-read the whole diff.

### Outcome
- Finding 1 (remove side and darwin label untested) — **fixed.** `TestUnitNamesPerHome` now asserts `RemoveTimer(other, "linux")` calls `disable --now <suffixed>.timer`, that the darwin plist for the scratch home carries the suffixed label and a `<label>.plist` path, and that `bootout` names that plist. Reverting `RemoveTimer`'s unit to the bare name fails the Linux assertion; reverting the label fails the plist assertion. Read from the test code; no mutation was run.
- Finding 2 (scratch job dreams over the real home) — **fixed.** Both renderers now emit `dream --all --home <home>` (systemd with `%` escaped for the home; launchd through the XML escaper); the home sits inside `bin`, so the existing quote/backslash/newline refusal covers it. `TestRenderSystemdUnits`, `TestRenderLaunchdPlist`, the `%` escape test and `TestUnitNamesPerHome` pin it.
- Finding 3 (symlinked spelling of the real home) — **fixed.** `UnitFor` falls back to `os.SameFile` on the two stats after the Clean compare; the new test symlinks the real home and expects the bare name. The unreadable-`$HOME` low note (finding 4) is left, as unreachable from the CLI.

### New observation (low / confidence: medium)
`internal/dream/timer.go:98` — the real-home unit text changes (it now carries `--home "<home>"`). A v2.0.0 install holding the old text, on re-running `dream --install-timer`, hits the existing "exists with other text; kept ... run --remove-timer first" refusal (exit 2). That is the designed behaviour and the old unit keeps working (it defaults to the real home), so no data is lost; the changelog or delivery note should say that an upgraded home runs `--remove-timer` then `--install-timer` to pick up `--home`. Not blocking.

## Verdict (round 2): APPROVE

## Status: COMPLETE
