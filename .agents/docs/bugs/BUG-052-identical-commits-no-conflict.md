---
owner: bugs
date: 2026-09-17
id: BUG-052
status: fixed
area: internal/harness (report_test.go)
---
# BUG-052: rejected-push case of `TestSendReport` returned nil error, no rejection

## Symptom
The "behind checkout" case of `TestSendReport` expected `SendReport` to fail
with a "push rejected" error; it returned `err == nil` instead.

## Repro
Two checkouts of one bare origin repository; the second checkout was behind
the first. Both ran `SendReport` with the same fixed clock (`now`) and the
same `Host`, writing byte-identical report directories/files.

## Root cause
With identical tree contents, author/committer, commit message, and a
second-resolution timestamp, git produced the same commit hash (confirmed
with a debug test printing both logs: `dbe0cd1` on both sides). The "behind"
checkout's push was therefore a fast-forward no-op that pushed nothing new,
not a rejected non-fast-forward push — there was no conflict to reject.

## Fix
`internal/harness/report_test.go` `TestSendReport`: the behind checkout
(`other2`) uses a different `Host` (`"erin"` vs `"dave"`) so its commit
content, and therefore hash, differs from the one already pushed.

## Prevention
A test asserting a git conflict (rejected push, merge conflict, divergent
history) must make the competing commits differ in content, not just in
timing or checkout path; a fixed clock combined with otherwise-identical
inputs produces byte-identical, same-hash commits, which git treats as
already-applied rather than conflicting.

## Related
BUG-051 (same test, same first run).
