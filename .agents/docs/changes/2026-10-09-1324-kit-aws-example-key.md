---
owner: main
date: 2026-10-09
---
# The kit's credential scan sets aside AWS's published example key

The operator chose "fix then preflight", taking the v2.6.0 handoff's one Next Step.

- **The problem.** In the v2.6.0 gate the Antigravity sandbox kit FAILed twice on `AKIAIOSFODNN7EXAMPLE`, AWS's documented example access key, which the agent read in `internal/harness/report_test.go`. Each time the run had to be cleared by hand: delete by name, re-scan, planted-key control.
- **`testing/kit-preamble.sh`:** `kit_credential_leftovers` collects a file's matches with `grep -ao` (in memory, never printed) and reports the file when any match is not exactly `kit_cred_example`. Every other credential kind and name check is unchanged.
- **`TestKitPreambleRefusesLeftoverCredentials`** gains three cases:
  - a home holding only the example key, in a transcript and a binary database, must pass. This failed before the change;
  - a real-shaped AWS key alone must FAIL;
  - a real-shaped AWS key beside the example in one binary file must FAIL. A mutation that ignores any file containing the example made this case fail, so the check is exact.
- **`testing/guide.md`:** the credential paragraph names the exception.
- **Verification:** `go test -run TestKitPreamble ./internal/harness` exit 0. Diff read by me (21 lines, below the reviewer threshold); the release gate's reviewer, security and verifier cover it.
