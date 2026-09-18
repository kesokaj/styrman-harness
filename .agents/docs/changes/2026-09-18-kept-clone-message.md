---
owner: main
date: 2026-09-18
---
# The kept-clone message says what its count was measured against     (2026-09-18)

## Summary

`sh-agent update` on the operator's own box refused to remove the old kept feedback clone at `~/.config/sh-agent/feedback` and reported `has 2 unpushed commit(s)`. The refusal was right — the guard fails closed by design (BUG-069) — but the reason was wrong in a way that cost time: those two commits were on the remote already. `report` pushes with `-u <url>`, which leaves the branch tracking a URL rather than a remote, so `@{upstream}` does not resolve and no remote-tracking ref ever moves; the count (`rev-list --all --not --remotes`) then measures against a ref frozen at the commit the clone was made from. The message now says what it measured, that the commits may well be pushed, and the two commands that settle it. The behaviour is unchanged: a count it cannot trust still keeps the clone.

Measured on the real clone: local `HEAD` `e973bac`, remote `main` `09586ff`, remote-tracking `origin/main` stale at `71f749e`; after a fetch, `git log --oneline HEAD --not FETCH_HEAD` was empty — the clone held nothing unique.

## Type

bugfix

## Files changed

- `internal/harness/report.go` — `RemoveKeptClone`'s ahead-count message, with the reason the basis goes stale in a comment beside it
- `internal/harness/report_test.go` — the two assertions on the old wording updated, plus a new subtest that reproduces the production shape (branch tracking a URL, a pushed commit, no moved tracking ref) and asserts the clone is kept WITHOUT the message claiming unpushed work
- `.agents/docs/deploy/local.md` — § Report back's row explains the message and gives the command that settles it

## Agents dispatched

None. The diff is 1 changed message, 1 new subtest and 1 doc row, below the review threshold (context § Definition of Done step 2), so the main agent read it: the only behavioural change is the string, the guard's decision paths are untouched, and the new subtest fails against the old message together with the two updated assertions (seen: 3 failures, then restored).

## Verification

`gofmt`, `go vet`, `go test -race -count=1 ./...` clean. `TestRemoveKeptClone` is eleven subtests, all passing, and the new one did not SKIP — its `@{upstream}` guard confirms the shape reproduced rather than silently passing. Must-fail: with the old message restored, three subtests fail. Re-run of `sh-agent update` against the real home afterwards shows the new wording on the real clone.

## Docs

`.agents/docs/deploy/local.md` § Report back, above. No other living doc quoted the old message.

## Requirements

None touched.
