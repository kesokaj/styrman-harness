---
owner: main
date: 2026-10-09
id: BUG-165
status: fixed
area: internal/search (reconcile)
---
# BUG-165: `sh-agent search` runs an unhardened `git status` in the checkout it indexes

## Symptom
Round 3 of the security review of the metrics range (`reviews/2026-10-08-security-2c0f215-166aa39.md`, an info item outside that range) noted that `internal/search/reconcile.go:132` runs `git status` with no hardening. fsmonitor is on and clean filters are live, every time `sh-agent search` or `index` runs. It is the same class as BUG-164: a repository's own config or `.gitattributes` can start a program. It has been released since search shipped.

## Repro
- `TestGitViewRunsNoCleanFilter` (`internal/search/hardening_test.go`) builds the checkout: `.gitattributes` routes `*.md` through `filter.x.clean`, a marker script; one file's mtime is moved. The fixture first proves a plain `git status` writes the marker. Before the fix, `gitView` wrote it too: the test was red ("search ran the repository's clean filter").
- Security review of the fix (`reviews/2026-10-09-security-worktree.md`, finding 1, medium) found the same class in `internal/memory`: `git check-ignore` started the checkout's `core.fsmonitor` program. `TestAddRunsNoProjectHook` was red before the fix.

## Root cause
BUG-164's fix put the hardened git calls in `internal/trail` and made them the rule for project roots taken from the index or a transcript. `search` and `memory` run in the directory the operator invoked them in, so that review's range never covered them, and they kept their own git calls.

## Fix
- **The builder.** `internal/trail.Command` (exported) is the one builder for git in a project repository. It sets no index lock, no fsmonitor, no signature check, and `GIT_NO_LAZY_FETCH=1`.
- **Search's git calls.** `internal/search/reconcile.go` `git()` and `search.go` `guardProject` (`check-ignore`) are built on it.
- **No more `status` in search.** `gitView` no longer runs `git status`:
  - A tracked file keeps the index's blob when its size and mtime still match the index's stat data (`ls-files -s --debug`, read line by line; a quoted name is skipped and its file hashed) and it is not racily clean.
  - Otherwise its blob ID is computed from its own bytes; modified means that ID differs from the index's.
  - Untracked files come from `ls-files --others --exclude-standard`.
  - Every `.sh-agent-noindex` marker's existence is checked whatever the prefix.
- **Memory's git calls.** `Root`, `storeTracked` (was `storeTrackedArgs`) and `Add`'s `check-ignore` in `internal/memory/memory.go` go through `trail.Command`.
- **Cost.** Search per query is about 116 ms, the same as v2.5.0 on this checkout. A first version, which hashed every tracked file, took about 600 ms.
- **Tests.** Each was red before its fix:
  - `TestGitViewRunsNoCleanFilter`
  - `TestGitViewRunsNoFsmonitor`
  - `TestGitViewRacilyCleanEditIsHashed` (its mutation, the racy check removed, turned it red)
  - `TestGitViewDeletedMarkerOutsidePrefix`
  - `TestIndexStatsIgnoresForgedNames`
  - `TestAddRunsNoProjectHook`
- **Reviews.** `reviews/2026-10-09-review-worktree-code-docs.md` and `reviews/2026-10-09-security-worktree.md`, two rounds each.

## Prevention
As BUG-164: every git call into a project repository goes through `internal/trail.Command` and only reads objects and files (never `status`). A new git command there gets a marker test against a repository whose config or attributes would start a program. `internal/harness` is the exception: it runs git only in the harness's own checkout.

**Partial clones (security review finding 3, closed 2026-10-09):** search refuses a partial clone only when this machine's git is older than 2.44 (`trail.LazyFetchSafe`); a newer git honours `GIT_NO_LAZY_FETCH` and searches it. Pinned by `TestGitViewPartialCloneNeedsNewGit` (red before the fix) and `TestVersionAtLeast`.
