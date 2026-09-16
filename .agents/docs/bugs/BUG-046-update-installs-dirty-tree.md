---
owner: bugs
date: 2026-09-16
id: BUG-046
status: fixed
area: cmd/sh-agent (update)
---
# BUG-046: `sh-agent update` reinstalled an uncommitted, dirty checkout

## Symptom
At 17:04:46 today the operator's user-level systemd timer
`sh-agent-update.timer` (hourly, `RandomizedDelaySec=5min`, unit files
created 16:55:34 today in `~/.config/systemd/user/`, not mentioned
anywhere in this repository's trail) ran `sh-agent update --repo
~/workspace/styrman-harness`; the checkout was up to date with upstream
(8b95f22) but dirty with this session's uncommitted edits, and `update`
reinstalled it anyway: the real home's `~/.claude/AGENTS.md`,
`~/.claude/skills/preflight/SKILL.md`,
`~/.gemini/config/rules/AGENTS.md`, and `~/.config/sh-agent/state.json`
(`commit: 8b95f22+dirty`, `installedAt 17:04:48`) took mid-task text
(journal: "up to date with upstream (8b95f22)", "building sh-agent from
…", "10 written, 250 unchanged").

## Repro
Have a dirty, up-to-date checkout and run
`sh-agent update --repo <checkout>`; it reinstalls from the working tree
as it stands instead of refusing.

## Root cause
`cmd/sh-agent/main.go` `update()` calls `harness.Reinstall`
unconditionally: when the checkout is not behind upstream it skips
`Pull` (which refuses a dirty tree) and installs the working tree as it
stands. `.agents/docs/deploy/remote.md` describes `update` as "fetch; if
behind: pull --ff-only, rebuild, reinstall", which is not what the code
does — the "not behind" path has no dirty-tree guard.

## Fix
Fixed 2026-09-16 (Phase 1, clear the board): `cmd/sh-agent/main.go`
`update()` refuses a dirty checkout on every path, not just "behind
upstream" (`cmd/sh-agent/main.go:350`, "A dirty checkout is never
installed, on either path"); `remote.md` § The resync loop and `README.md`
updated to match; `ADR-027` records the resync loop's shape; `REQUIREMENTS.md`
R4 references the fix. See BUG-049 for a gap the same day in the dirty
check this fix relies on.

## Prevention
`update` refuses a dirty checkout on every path, not just "behind
upstream" — graduated into `cmd/sh-agent/main.go`, `ADR-027`, and
`REQUIREMENTS.md` R4. Before relying on `sh-agent update` on a checkout used
for active development, confirm the dirty check itself is pinned (see
BUG-049): a guard that trusts git's configurable default can still pass a
dirty-looking-clean tree.

## Related
Investigation cost four rounds before the timer was found (agent
transcripts, nested sandbox CLI runs, hooks, profiles) — the pattern
matches BUG-039/041 style "silent instrumentation" surprises; worth a
standing rule that any process able to write the real home is named in
`deploy/remote.md § Downstream mirrors` or an equivalent inventory.
