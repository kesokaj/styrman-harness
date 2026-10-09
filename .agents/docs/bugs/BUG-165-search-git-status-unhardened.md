---
owner: main
date: 2026-10-09
id: BUG-165
status: open
area: internal/search (reconcile)
---
# BUG-165: `sh-agent search` runs an unhardened `git status` in the checkout it indexes

## Symptom
Round 3 of the security review of the metrics range (`reviews/2026-10-08-security-2c0f215-166aa39.md`, an info item outside that range) noted that `internal/search/reconcile.go:132` runs `git status` with no hardening. fsmonitor is on and clean filters are live, every time `sh-agent search` or `index` runs. It is the same class as BUG-164: a repository's own config or `.gitattributes` can start a program. It has been released since search shipped.

## Repro
Not reproduced yet. The expected repro is a checkout whose `.gitattributes` routes a tracked file through `filter.x.clean`, with that file's mtime moved, then `sh-agent search x` run in it. The filter would run.

## Root cause
BUG-164's fix put the hardened git calls in `internal/trail` and made them the rule for project roots taken from the index or a transcript. `search` runs in the directory the operator invoked it in, so it was outside that review's range and was never moved onto the helper.

## Fix
Open. The likely shape: route search's git calls through `internal/trail.command`, and replace `status` with the `ls-tree` and `hash-object --no-filters` comparison `trail.Dates` uses, under a marker test like `TestDatesRunsNoCleanFilter`. Search also needs untracked files, which `ls-files --others --exclude-standard` lists without running filters.

## Prevention
As BUG-164: every git call into a repository goes through the hardened builder and only reads objects and files. Graduates when search uses it.
