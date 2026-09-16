---
owner: bugs
date: 2026-09-16
id: BUG-046
status: open
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
Not fixed. Candidate: `update` reinstalls only after a successful `Pull`,
or refuses a dirty checkout the way `Pull` does and says so;
`remote.md` § The resync loop then matches the code; add a test for
"not behind, dirty → no install".

## Prevention
Before relying on `sh-agent update` (via a timer or manually) on a
checkout used for active development, confirm it refuses a dirty tree on
every path, not just the "behind upstream" path — R6 ("development never
touches the real home") is defeated by any timer that installs from a
dirty working tree. Until fixed, a handoff on this box states the timer
exists so a session does not lose hours attributing a mystery write to
its own agents.

## Related
Investigation cost four rounds before the timer was found (agent
transcripts, nested sandbox CLI runs, hooks, profiles) — the pattern
matches BUG-039/041 style "silent instrumentation" surprises; worth a
standing rule that any process able to write the real home is named in
`deploy/remote.md § Downstream mirrors` or an equivalent inventory.
