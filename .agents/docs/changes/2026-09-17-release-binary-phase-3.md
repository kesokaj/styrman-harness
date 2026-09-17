---
owner: main
date: 2026-09-17
---
# `sh-agent report` keeps nothing — Phase 3 of the release-binary plan     (2026-09-17)

## Summary

`sh-agent report` no longer keeps a clone of the private feedback repository on every machine, which is what let each place read every other place's reports ("But why can I see dave reports?"). Without `--feedback-repo`, a run now clones `~/.config/sh-agent/report-<pid>` with `--filter=blob:none --no-checkout --depth 1`, limits it with `git sparse-checkout set --no-cone reports/<host>/`, writes its run directory, commits, pushes to the feedback URL, and removes that directory on every path out — a failed push, a refusal and an early return included. No other machine's report content is fetched to the disk; this machine's own earlier runs do come along, which is what lets the `-2` suffix pick a free directory name. A server that does not support the filter makes git warn and send everything: slower, still correct. Because a throwaway clone has nowhere to keep a rejected commit, a non-fast-forward push is retried once from a fresh clone; a second rejection says `push rejected twice; re-run sh-agent report — nothing is kept locally` and leaves nothing behind. `--feedback-repo DIR` is unchanged in every respect: that checkout is the operator's, so it is written into and never removed, a stray file in it is named and left alone, and a rejected push leaves the commit there with today's `pull --rebase` message. `sh-agent update` removes the clone earlier versions kept at `~/.config/sh-agent/feedback`, but only when it is that clone, clean, and not ahead of its upstream; anything else is left in place with the reason named.

## Type

task

## Files changed

- `internal/harness/report.go` — `ResolveFeedback` reduced to the `--feedback-repo` resolver; new `throwawayDir`, `cloneThrowaway`, `sweepInterruptedClones`, `alive`, `RemoveKeptClone`, `writeAndPush`, `scrubCredentials`, `errPushRejected` and the `testHooks` pair; `SendReport` split into the dry-run, named-checkout and throwaway paths, the last a two-attempt loop with removal deferred inside each attempt; `commitAndPush` returns the rejection sentinel so each caller words it
- `internal/harness/report_test.go` — `TestSendReport` rewritten for the default path (reads what was sent from the origin, `assertKeepsNothing` after each run); new `TestSendReportNamedCheckout`, `TestSendReportThrowawayClone`, `TestSendReportRetriesRejectedPush`, `TestRemoveKeptClone` (nine subtests), `TestSendReportSweepsInterruptedClones`, `TestCloneFailureDoesNotPrintCredentials`; `TestResolveFeedbackFailures` rewritten for the new signature plus the unreachable-URL case; helpers `seedFeedback`, `localBlobHolds`, `configDir`
- `cmd/sh-agent/main.go` — `update` calls `harness.RemoveKeptClone` below every `--check` return, reported and never fatal; usage text and the command comment describe the throwaway clone
- `internal/harness/repo.go` — `gitOut` split so `gitOutC` can run one call with `LC_ALL=C`: the push whose rejection message the retry has to recognise
- `cmd/sh-agent/main_test.go` — `TestUpdateKeptCloneCleanupOnlyWhenUpdating`
- `.agents/docs/deploy/local.md` — § Report back rewritten: a new "Nothing is kept here" bullet, both rejection paths, and what `update` does with an old kept clone
- `.agents/docs/plans/2026-09-17-release-binary.md` — Phase 3's three task boxes ticked with evidence; the partial-clone assumption marked settled; D25

## Agents dispatched

- `verifier` — PASS, `.agents/docs/testing/2026-09-17-release-binary-phase-3.md` (87 pass, 0 fail, 0 skip at 83f3cf6)
- `security` — 2 medium, 6 low, `.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md`; five fixed with tests, three rebutted in its `## Outcome`. The two mediums were real defects in this phase's own code: `RemoveKeptClone`'s unpushed-commit guard failed OPEN when the count could not be taken (a branch with no upstream), and the cleanup call sat above `update --check`'s returns, so the read-only mode deleted the clone. Both now have a test that fails without the fix, and the function additionally refuses a symlink at that path, treats a stash and ignored files as work, and resolves a relative `--home`. Also fixed: an interrupted run's `report-<pid>` is swept by the next run (a deferred removal cannot cover a signal), git's clone error no longer repeats a URL a rewrite can fill with a token, and the one push whose message the retry must recognise runs with `LC_ALL=C`.

## Verification

`gofmt`, `go vet` and `go test -race -count=1 ./...` clean at each commit. The git mechanics were settled before any code was written (`.agents/scratch/p3-probe/probe.sh`), which caught two things that would have made the isolation check vacuous: a `file://` remote ignores `--filter` unless it sets `uploadpack.allowFilter`, and `git cat-file` on a promised object lazily FETCHES it, so a naive check answers yes about a blob the clone never held. The tests therefore seed `uploadpack.allowFilter` and read with `GIT_NO_LAZY_FETCH=1`. `TestSendReportThrowawayClone` inspects the clone through `testHooks.beforeRemove`: the other host's directory is never checked out and its canary is in no local blob, the run lands on the remote, and afterwards there is no kept checkout and no `report-*`. Must-fail: every new and rewritten test failed to BUILD against the old `ResolveFeedback`, so the guards were mutated one at a time instead (`.agents/scratch/p3-mustfail/mutate.py`, which asserts each target text is unique and restores from a saved copy, per BUG-068) — removal not deferred, `--filter` dropped, `sparse-checkout` dropped, the retry removed, and the dirty check dropped each failed exactly the test that covers it, and the file was restored. One real run against GitHub, authorised by the operator (D25) and logged in `.agents/docs/testing/2026-09-17-release-binary-phase-3-real-report.md`, settled the last assumption: the blob-less shallow clone over SSH succeeded, the push landed as `09586ff`, and nothing was left in `~/.config/sh-agent/`. The security review then found two real defects in this phase's own code (see § Agents dispatched); each was reproduced, fixed, and covered by a test seen failing without the fix, and the suite ended at 90 top-level tests and 73 subtests, none failing or skipped. Not exercised: what the real clone contained (it is removed by design, so only the tests can show that), the retry against the real remote, `RemoveKeptClone` against this box's own kept clone (which is `sh-agent update`'s job in the operator's home), and a pid collision between two runs in separate pid namespaces (rebutted in the review's `## Outcome`).

## Docs

`.agents/docs/deploy/local.md` § Report back is current with this change. `README.md` § How it works (the map's report line), `.agents/docs/deploy/remote.md` § The four places and `.agents/docs/specs/sh-agent/delivery.md` still describe a kept clone; they land with R8, R9 and the map in Phase 5's single commit to `main` (D11, R9).

## Requirements

R9's "Report back" bullet changes to the text the operator approved in the plan's Interfaces, applied in Phase 5. R4 holds: stdlib only, git and ssh shelled out as before, no dependency added. Nothing here reaches a machine until that commit lands and a place updates.
