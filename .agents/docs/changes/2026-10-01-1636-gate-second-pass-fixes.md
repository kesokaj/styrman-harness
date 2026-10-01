---
owner: main
date: 2026-10-01
---
# 2026-10-01-1636 — The v2.1.0 gate's second pass: the timer's account lookup fails closed, the script's guard resolves through missing directories, notes never relax the router

Fixes after f62513b, from three reports on that commit: the code review of 7dd0645..f62513b (REQUEST CHANGES, one medium), the security read of the catalog's manifest and that delta (six lows), and the verifier's log (PASS).

- **Dream timer, correcting the 1615 entry**: that entry said the real home comes from the account database and named macOS without cgo as the limit. Both halves were off. `user.Current()` answers with `HOME` when the lookup fails, which on Linux without cgo (the release binaries) is every account outside `/etc/passwd`; macOS asks the system resolver either way. Now `accountHome` uses `user.LookupId` by uid with no fallback: an account the database does not list has no real home and every home there gets the suffixed unit name. `TestAccountHomeNeverFallsBackToHOME`; `TestUnitForIgnoresHomeOverride` no longer skips. `specs/sh-agent/dream.md`, `deploy/local.md` and BUG-149 say the same.
- **Refresh script**: `--dest` is resolved through its nearest existing ancestor, and a `.` or `..` below a directory that does not exist is refused (it could name the home directory, the checkout or `/` once `mkdir -p` ran); an unset `HOME` no longer stops the script; a refusal prints an upstream name shell-quoted, never raw.
- **Catalog notes**: a note narrows an entry and never relaxes the router's steps 3 to 5 (router step 2, manifest header). Six notes reworded: no `--quiet` in three, `--project` per command in place of `gcloud config set project`, no plugin or MCP server added by the agent in two, and the scaffold's generated guidance reaches an existing project's `AGENTS.md` only as a diff the operator has seen.
- **Docs**: `deploy/local.md` names the files to delete by hand for a timer that v2.0.0 or earlier wrote under the bare name from a home that is not the account's.

Left with an owner (gate report and handoff): a notice for such a legacy timer in `--remove-timer`, `uninstall` and `status`; the items the 1615 entry lists; the check that the claude.ai artifact named in the v2.0.0 handoff is not shared by link.

Verification before the commit (main agent): `gofmt -l cmd internal` empty, `go vet ./...`, `go test -race -count=1 ./...` all six packages ok; red first: the unlisted account, the four missing-directory destinations, the escape sequence in a name. Both CLIs on f62513b: catalog rows 9 of 9 on Claude Code and on Antigravity (no-leak by two identical self-listings of 37 names); the six harness skills changed since v2.0.0 ran on Antigravity, 6 of 6 (`bash`, `cloud-run-deploy`, `data-gcp`, `review`, `feature`, `project-memory`). The router rows run again on this commit's text. Not covered: a real service manager call; an account outside `/etc/passwd` (the test simulates the failed lookup through a seam).

Agents dispatched: verifier ×1 (PASS over v2.0.0..f62513b, `testing/2026-10-01-verify-v2.1.0-gate-r2.md`), reviewer ×1 (`reviews/2026-10-01-review-7dd0645-f62513b-code.md`), security ×1 (`reviews/2026-10-01-security-v2.1.0-catalog-harness-files.md`).
