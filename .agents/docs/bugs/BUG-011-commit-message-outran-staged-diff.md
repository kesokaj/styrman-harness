---
owner: bugs
date: 2026-09-13
id: BUG-011
status: fixed
area: git/process
---
# BUG-011: commit message described 60 files, commit contained one

## Symptom
Commit 8a8b276, whose message describes the whole "improvement review applied"
change (one TODO/ROADMAP owner, handoff trimmed, agents without private memory,
plain-copy checkout test, status reads installed binary, gcloud skills merged,
docker/frontend rewritten, golang references cut, ADR-016), actually contains
exactly one change: deletion of `skills/gcloud-auth-verification/SKILL.md`
(49 lines). Pushed HEAD briefly had the skill removed without the rules
merged into `gcloud-setup`, until e5833a3 (the rest) followed minutes later.

## Repro
`git show --stat 8a8b276` — one file, 49 deletions, against a message
describing ~60 files.

## Root cause
Staging was done from a saved path list via `xargs -a list git add --`. The
list still named the already-`git rm`'d file. `git add` aborted on the stale
pathspec, but the command line continued with `; git commit`, chaining the
commit after the add without reading the add's exit status — the commit ran
on whatever happened to be staged (just the earlier `git rm`), not on the
full change. Same family as BUG-033/041 in styrman-crew and the global hard
rule on verification chains; that rule had only been applied to test runs,
not to the staging step itself.

## Fix
Remaining paths (list minus the deleted path) staged by name, `git diff
--cached --stat` read (60 files, nothing left unstaged), commit issued as
its own call: e5833a3. History was not rewritten (global git rule forbids
it); 8a8b276's message is misleading and e5833a3's message says so.

## Prevention
Staging is a verification step like a test run: run `git add` bare, read its
exit status and `git diff --cached --stat`, then commit in a separate call —
never `git add … ; git commit` or `… && git commit`. A path list saved
earlier is stale after any delete or rename in the same session; regenerate
it from `git status --porcelain` immediately before staging, don't reuse it.

## Related
BUG-010 (a stale list applied half-way, same session); global context Hard
rules bullet on verification chains leading to a commit. Three records now
share the "chain into a commit" mechanism (BUG-010, BUG-011, and the global
rule's own styrman-crew precedent BUG-033/041) — worth a project-level
`.agents/rules/` entry making "stage, read, commit" its own three-step
sequence, not just the test-command case the global rule states.
