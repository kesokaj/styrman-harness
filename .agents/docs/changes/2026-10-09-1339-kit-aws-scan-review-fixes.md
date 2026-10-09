---
owner: main
date: 2026-10-09
---
# v2.6.1 gate: the kit's AWS check counts, never captures

Fixes from the v2.6.1 gate's reviews of f6d9854: `reviews/2026-10-09-security-v2.6.0-f6d9854.md` (security, one medium) and `reviews/2026-10-09-review-v2.6.0-f6d9854-code-docs.md` (reviewer, approve with nits). Both end in `## Outcome`.

- **`testing/kit-preamble.sh`, `kit_credential_leftovers`:**
  - `kit_cred_text` loses its `AKIA` alternative and keeps v2.6.0's `grep -q` check.
  - AWS keys are checked apart. The matches go through a pipe into `grep -cvxF <example>`, and only the count is assigned, so a key never reaches a variable that `bash -x` would trace, nor a here-string's temp file.
  - Both use `command grep`, so a shell function named `grep` cannot blind the scan.
  - The function ends in `return 0`, so its own status never trips the EXIT trap's `set -e`.
- **`TestKitPreambleRefusesLeftoverCredentials`** gains four cases:
  - a key found under `bash -x` stays out of the output;
  - a `grep` function that finds nothing still FAILs the run;
  - one file holding only the example passes, with the scan's status 0;
  - a GitHub token holding the example still FAILs.

  The paths-only check now covers the AWS key. Against f6d9854's preamble the first three fail; the fourth is a guard and passes on both. All pass after the fix.
- **Verification:** `go test -run TestKitPreambleRefusesLeftoverCredentials ./internal/harness` exit 0. The gate's security round 2 and verifier run after this commit.
