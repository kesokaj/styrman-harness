---
owner: bugs
date: 2026-10-01
id: BUG-144
status: fixed
area: internal/dream (timer), uninstall
---
# BUG-144: `--home <scratch>` timer commands drove the operator's real service-manager unit

## Symptom
The v2.0.0 gate's security review of the Go code (`reviews/2026-10-01-security-v1.16.2-0b1b5b8-code.md`, low/high; row 6 of `reviews/2026-10-01-preflight-v2.0.0.md`): `sh-agent dream --install-timer`, `--remove-timer` and `uninstall` with `--home <scratch>` wrote unit files under the scratch home but named the unit `sh-agent-dream.timer` (launchd label `com.styrman.sh-agent-dream`) and drove `systemctl --user` / `launchctl` with that name. A scratch `--remove-timer` or `uninstall` disabled the operator's real timer on Linux; a scratch `--install-timer` on macOS loaded a nightly job over the real home.

## Repro
`sh-agent dream --install-timer --home .agents/scratch/x` on a box with the real timer installed, then `sh-agent dream --remove-timer --home .agents/scratch/x`; `systemctl --user list-timers` no longer lists `sh-agent-dream.timer`.

## Root cause
The unit name and label were package constants, chosen before `--home` existed as a development seam. Every path under the home was parameterised; the one thing outside the home, the manager's namespace (shared by every home of the user), was not. R6 ("development never touches the real home") held for files and failed for the service manager. The parity plan's review (`reviews/2026-09-30-cold-read-parity-memory-dream-r1.md`, R6 finding) led to tests with a fake manager, which hid the collision: the fake accepted any name.

## Fix
`dream.UnitFor(home)` in `internal/dream/timer.go`: the real home (`os.UserHomeDir`, cleaned) keeps the bare names; any other home gets `-<first 8 hex of sha256(clean path)>` appended to both. `InstallTimer`, `RemoveTimer`, `TimerInstalled` and the path helpers derive the name from the same home. `TestUnitNamesPerHome` fails without the fix (bare names for a temporary home; measured red, then green); existing tests pin the real home to their temp dir through the `userHomeDir` seam. Operator decision 2026-10-01: "Go with your recommendations on all three" on the gate's one-line choice. Working tree, committed with this record.

## Prevention
Anything a `--home`-scoped command touches OUTSIDE the home (a service manager's unit namespace, a PATH entry, a global cache) is named or keyed by the home as well, and its test asserts the name the manager receives, not only the files under the home.

## Related
R6; BUG-107 and BUG-123 (the kits' real-config guards: same family, guarding files and config, not the manager); BUG-143 (same gate); BUG-149 (a defect in this fix: the real home was read from `$HOME`).
