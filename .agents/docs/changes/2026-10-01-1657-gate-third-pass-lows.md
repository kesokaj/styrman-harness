---
owner: main
date: 2026-10-01
---
# 2026-10-01-1657 — The v2.1.0 gate's third pass: lows from the review and security reads of b4e3d79

After b4e3d79: the verifier's log (PASS, one coverage gap), code review round 2 (APPROVE WITH NITS, three lows) and the security read of f62513b..b4e3d79 (four lows). No medium or high.

- **Refresh script**: stops when a directory on `--dest`'s path cannot be entered (the unresolved path once became one under `/`) and when `HOME` is not set (the home-directory guard was skipped); the symlink, Go-file and dotfile refusals print the nested path shell-quoted. Tests for each in `TestRefreshCatalogScript`, red first.
- **Catalog notes**: `TestCatalogNotesNeverAllowQuiet` pins the kind: a note that names `--quiet` says to leave it out. Seven notes reworded (five named the flag with a scoped or after-approval exception; `google-cloud-storage-basics` no longer endorses installing the MCP Toolbox binary; `google-cloud-recipe-onboarding` shows the billing link command and runs it on the operator's go).
- **Docs**: `deploy/local.md` § Operator steps gives the by-hand cleanup of a timer that v2.0.0 or earlier wrote under the bare name, as two cases (a scratch home; the real home of an account the database does not list).

Left with an owner (gate report and handoff): one choice covering three findings — `--remove-timer`, `uninstall` and `status` say what they found and print the cleanup command for a legacy bare-named timer; `printf %q` on macOS's bash 3.2 not run.

Verification before the commit (main agent): `gofmt -l cmd internal` empty, `go vet ./...`, `go test -race -count=1 ./...` all six packages ok; red first: the unreadable parent, the unset `HOME`, five notes under the new test, the nested path with an escape sequence. Router rows on b4e3d79: 9 of 9 on Claude Code and on Antigravity; they run again on this commit's notes.

Agents dispatched: verifier ×1 (`testing/2026-10-01-verify-v2.1.0-gate-r3.md`, PASS), reviewer ×1 (round 2 of `reviews/2026-10-01-review-7dd0645-f62513b-code.md`), security ×1 (`reviews/2026-10-01-security-f62513b-b4e3d79.md`).
