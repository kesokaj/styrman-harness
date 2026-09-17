---
owner: bugs
date: 2026-09-17
id: BUG-073
status: fixed
area: process / dated trail records
---
# BUG-073: new retro record overwrote an existing same-day retro instead of taking the next free suffix

## Symptom
The main agent wrote a new retro to `.agents/docs/retros/2026-09-17-retro-3.md`, a path that already held a third retro from an earlier session (committed in 8af1875). The write replaced 7 existing lines with 71 new ones. `git status` showed the file as MODIFIED, not new — that is what exposed it.

## Repro
With `.agents/docs/retros/<date>-retro.md`, `-retro-2.md`, and `-retro-3.md` already on disk, write a new record for that date by choosing the next suffix from the session's own memory of how many retros had run, rather than from `ls` on the directory.

## Root cause
The suffix was derived from in-context knowledge of retro count, not from the directory. The session had read `.agents/docs/retros/` earlier and seen only two files; a third was committed later the same day by another session, going stale in this session's context.

## Fix
The clobbered record was restored from git (content verified against 8af1875); the new record rewritten as `2026-09-17-retro-4.md` with its title marked "fourth"; all references corrected — the changes entry renamed to `2026-09-17-retro-4-applied.md`, the harness queue file's record path, and the last line of `~/.config/sh-agent/retros.jsonl`. Not lost only because the clobbered version was in git.

## Prevention
Before writing any DATED trail artifact whose name may already be taken (retro, handoff, changes entry, testing log, review), `ls` the target directory in the same turn as the write and pick the next free suffix from what is there, never from what this session remembers writing — an in-context directory listing goes stale the moment another session commits. When a write meant to create a file instead shows as MODIFIED in `git status`, that is read before staging, per BUG-011.

## Related
BUG-011 (read git status/diff --cached before trusting a command did what it was meant to). Three records now share the "read git status before trusting an intended effect" theme (BUG-011, BUG-013, BUG-073); worth considering for a standing pre-write directory-listing rule in `context/AGENTS.md` § Hard rules if it recurs again.
