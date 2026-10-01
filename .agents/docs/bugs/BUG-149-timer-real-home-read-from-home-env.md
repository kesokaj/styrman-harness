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
`accountHome()` in `internal/dream/timer.go` reads the account's home from the account database by uid (`os/user`, `user.LookupId`) and has no fallback: an account the database does not list has no real home, and then every home gets the suffixed name. The first version of this fix used `user.Current()` with a fallback to `os.UserHomeDir`, and the code review of 7dd0645..f62513b showed that `Current` itself answers with `HOME` when the lookup fails — on Linux a build without cgo reads only `/etc/passwd`, so LDAP, SSSD and systemd-homed accounts were unprotected — and that the limit the first version named (macOS without cgo) was wrong: macOS asks the system resolver either way. Tests in `internal/dream/timer_test.go`: `TestUnitForIgnoresHomeOverride` (no skip) and `TestAccountHomeNeverFallsBackToHOME` (a lookup that fails, and an entry with an empty home). Same change: `renderSystemd` escapes `$` as `$$` in `ExecStart` (`TestRenderSystemdEscapesDollar`).

## Prevention
Never decide "the operator's own" from an environment variable a sandbox sets (`HOME`, `XDG_*`, `PATH`): ask the account database or the service itself. A fix for an isolation bug is tested in the isolation setup the project really uses (here `HOME` overridden AND `--home` passed), not only through a seam that replaces the lookup.

## Related
BUG-144 (this is a defect in its own fix); R6; BUG-107 and BUG-123 (kit guards, same family); `reviews/2026-10-01-security-v2.0.0-7dd0645-code-text.md`; `testing/2026-10-01-verify-v2.1.0-gate.md`.
