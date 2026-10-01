---
owner: bugs
date: 2026-10-01
id: BUG-149
status: fixed
area: internal/dream (timer, UnitFor), uninstall
---
# BUG-149: the timer's "real home" was read from `$HOME`, which a sandbox run overrides

## Symptom
BUG-144's fix gave every home but the real one a suffixed unit name and decided which home is real with `os.UserHomeDir`, that is `$HOME`. The testing guide's sandbox kits run with `HOME=<sandbox>` and pass `--home <sandbox>`, so inside a sandbox the sandbox WAS the real home: `sh-agent dream --remove-timer` or `sh-agent uninstall` there used the bare name `sh-agent-dream` and would have disabled the operator's real nightly timer through the shared user service manager, while `status` in the real home still read the unit file as installed. The exact case BUG-144 existed for.

## Repro
Found by reading, in the security review of the v2.1.0 gate (`reviews/2026-10-01-security-v2.0.0-7dd0645-code-text.md`, finding 1); no real `systemctl` call was made. The new test reproduces it without a service manager: with `HOME` set to a temporary directory, `UnitFor(<that directory>)` returned the bare names (red: "a HOME override got the real home's unit names: sh-agent-dream, com.styrman.sh-agent-dream").

## Root cause
The yardstick for "the operator's own" was an environment variable the caller controls; the service manager belongs to the account, not to `HOME`. BUG-144's tests replaced the lookup through a seam (`userHomeDir`), so no test ran the real lookup under an overridden `HOME`: the fix was verified in an isolation setup the project does not use.

## Fix
`accountHome()` in `internal/dream/timer.go` reads the account's home from the account database (`os/user`, `user.Current().HomeDir`) and falls back to `os.UserHomeDir` only when that gives nothing. `TestUnitForIgnoresHomeOverride` (`internal/dream/timer_test.go`) judges by the platform's own account lookup and skips where that lookup itself falls back to `HOME`. Same change: `renderSystemd` escapes `$` as `$$` in `ExecStart` (`TestRenderSystemdEscapesDollar`). Known limit, written in `specs/sh-agent/dream.md`: a build without cgo on macOS cannot read the account database for a normal user, so there the real home is `HOME` again (the release binaries are built without cgo).

## Prevention
Never decide "the operator's own" from an environment variable a sandbox sets (`HOME`, `XDG_*`, `PATH`): ask the account database or the service itself. A fix for an isolation bug is tested in the isolation setup the project really uses (here `HOME` overridden AND `--home` passed), not only through a seam that replaces the lookup.

## Related
BUG-144 (this is a defect in its own fix); R6; BUG-107 and BUG-123 (kit guards, same family); `reviews/2026-10-01-security-v2.0.0-7dd0645-code-text.md`; `testing/2026-10-01-verify-v2.1.0-gate.md`.
