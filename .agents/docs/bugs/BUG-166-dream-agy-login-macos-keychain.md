---
owner: main
date: 2026-10-09
id: BUG-166
status: fixed
area: internal/dream (agyHome, the private agy HOME)
---
# BUG-166: `sh-agent dream` cannot log in to Antigravity on macOS

## Symptom
On the operator's Mac (v2.6.1, release binary, agy logged in and working), `sh-agent dream --all` exited 2 before any project:
`sh-agent dream: no Antigravity login at /Users/spgo/.gemini/antigravity-cli/antigravity-oauth-token: run agy once and log in`. `sh-agent status` showed the same in its `dream` line. This reached the operator in a release.

## Repro
- On macOS with agy logged in: `sh-agent dream --all`.
- Operator probes, 2026-10-09:
  1. `ls ~/.gemini/antigravity-cli/`: there is no token file. `security dump-keychain` lists Antigravity and gemini entries.
  2. agy headless with `HOME` set to a fresh private directory (the dream job's shape): the log said `silent auth failed`, then `triggering interactive OAuth`. agy printed `Authentication required. Please visit the URL to log in:`, a URL, and `Waiting for authentication (timeout 60s)...`, then waited. After a login completed there, agy wrote an `antigravity-oauth-token` file into that HOME.
  3. The same, with `$HOME/Library/Keychains` linked into the private HOME: `silent auth succeeded`, no token file written, and the tool call was still refused (`soft-denying tool confirmation "RunCommand"`).
- `TestAgyHomeMacKeychain` (`internal/dream/run_test.go`) was red before the fix with the operator's own error text.

## Root cause
`agyHome` (`internal/dream/run.go`) was written and verified on Linux only. There agy keeps its login in `~/.gemini/antigravity-cli/antigravity-oauth-token`. On macOS agy keeps it in the login Keychain and finds that Keychain through `$HOME/Library/Keychains`. The private HOME the dream job gives agy (so the operator's tool grants do not apply, D28) had neither, so the file check refused the run. Dropping the check alone would not have fixed it: agy would have hung on a browser login. The macOS tests covered only the timer, against a fake launchd (`testing/2026-10-01-verify-timer-per-home.md`). No real macOS dream run was owed before a release.

## Fix
- **The login.** `agyHome(conf, goos)` in `internal/dream/run.go`. On darwin, when `$HOME/Library/Keychains` is a directory, it is linked (never copied) into the private HOME, and the token file is copied only when it exists. The run is refused only when neither exists, and the error names both places. Linux is unchanged. `os.RemoveAll` removes the link without following it; the test checks that the real folder survives cleanup.
- **The hint.** `agyAuthHint` adds `agy is not logged in inside the dream job's private HOME (...on macOS the login Keychain must be unlocked...)` to a project's failure line when agy's output asks for a login. It is added on both failure paths: a non-zero exit, and no answer. Before, a locked Keychain (no user logged in, the nightly timer) read as `invalid output`.
- **Tests.**
  - `TestAgyHomeMacKeychain` and `TestAgyAuthHint` were red before the fix.
  - `TestRunAgyAuthFailureNamed` goes through `Run`. Removing the hint from either failure path turned its matching case red.
- **Docs.** `specs/sh-agent/dream.md` § the agy command, `deploy/local.md` (the dream job's login), and `testing/guide.md` § 3e (a real macOS dream run owed before a release that changes the login path).
- **Not yet shown.** A real macOS `dream --all` with the fixed binary. It is owed before the release that carries this fix (operator step). The nightly launchd run against a locked Keychain is also unchecked (deferred, operator 2026-10-09).

## Prevention
How agy finds its login differs per OS: a file on Linux, the Keychain through `$HOME/Library/Keychains` on macOS. Any change to `agyHome` or to the private HOME owes one real `sh-agent dream --all` on a macOS place before the release (`testing/guide.md` § 3e). A fake CLI cannot show where the real one looks.
