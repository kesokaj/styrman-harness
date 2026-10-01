---
owner: main
date: 2026-10-01
---
# 2026-10-01-1717 — The v2.1.0 gate's fourth pass: two tests for a coverage gap; the remaining lows recorded, not fixed

After ff4c181: the verifier's log `testing/2026-10-01-verify-v2.1.0-gate-r4.md` (FAIL on one coverage gap, everything else green: 364 top-level tests with -race, three must-fail controls, recall 0.633, no vulnerabilities, three targets compile, zip 2,572 entries), code review round 3 (APPROVE WITH NITS, two lows) and the security read of b4e3d79..ff4c181 (four lows).

- **Tests only** (`internal/harness/catalog_test.go`): one escape-sequence fixture per quoted message of the refresh script (dotfile, Go file, symlink). The verifier had removed the quoting from the symlink message and nothing failed; now each of the two that lacked a fixture fails when its message prints the path raw (checked on a saved copy of the script, restored).
- **Recorded, not fixed** (Outcomes in `reviews/2026-10-01-security-b4e3d79-ff4c181.md` and `reviews/2026-10-01-review-7dd0645-f62513b-code.md` § Outcome (round 3)): one note that lets a skip-prompt flag through after a dry run; `--dest` spelled with a leading `//`; the timer cleanup's missing check for an unlisted account; a note test that a carve-out beside an allowed phrase passes; one note asking for what the entry's script hard-codes. No source changed for these: three fix rounds each moved the tree, and none of the five is medium or high.
- **Both CLIs on ff4c181's router and notes**: catalog rows 9 of 9 on Claude Code and 9 of 9 on Antigravity (no entry listed as a skill; 14 of 14 routed runs read the manifest note first; a deploy request stays with the harness).

Verification before the commit (main agent): `go test ./internal/harness -run TestRefreshCatalogScript` ok; the two new fixtures seen failing with their fix removed. The gate's verifier runs over this commit next (`-r5`). Review: a test-only change of one file, read by the main agent against the diff; a security read of that one path is dispatched with the verifier.

Agents dispatched: verifier ×1, reviewer ×1 (round 3), security ×1.
