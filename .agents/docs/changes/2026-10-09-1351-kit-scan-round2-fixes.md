---
owner: main
date: 2026-10-09
---
# v2.6.1 gate: the kit's credential scan fails closed on read errors and ignores shell functions

From the security review's round 2 (`reviews/2026-10-09-security-v2.6.0-f6d9854.md`, two lows). At the review cap the operator chose "Fix both, round 3".

- **`testing/kit-preamble.sh`, `kit_credential_leftovers`:**
  - The AWS check keeps grep's exit status: a read error (exit 2) adds a line that is not the example, so the file is reported, never passed. The previous commit's `|| true` had made that case pass, which v2.6.0 did not allow.
  - The function first runs `unset -f find head od tr grep`, so no shell function stands in for the scan's tools. The trap calls it in `$( )`, so the unset stays in that subshell.
- **`TestKitPreambleRefusesLeftoverCredentials`** gains three cases, each failing against eaa029a and passing now:
  - a `grep` on PATH that exits 2 for `-aoE`;
  - a `find` function that lists nothing;
  - an `od` function that says every file is ELF.
- **Verification:** `go test -run 'TestKitPreambleRefusesLeftoverCredentials|TestScriptsHaveNoEarlyExitPipe' ./...` exit 0. Security round 3 over these lines and the gate's verifier follow this commit.
