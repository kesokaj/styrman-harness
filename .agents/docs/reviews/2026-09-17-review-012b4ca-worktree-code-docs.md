---
owner: reviewer
date: 2026-09-17
---
# Review — worktree changes over 012b4ca (code + docs lenses)   (2026-09-17)

Range reviewed: working tree vs `012b4ca`, limited to the files named by the dispatch:
`internal/harness/report.go`, `internal/harness/report_test.go`, `context/AGENTS.md`,
`README.md`, `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`,
`.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/sh-agent/guide.md`,
`.agents/docs/plans/2026-09-17-delivery-rework.md`.

Blobs read (pinned via `git hash-object`):
- `internal/harness/report.go` = `9bb378e5005e99c0fea1eb8cfb79fff9ac0f11e7`
- `internal/harness/report_test.go` = `06b295af2041085d8f1e65389fcb5c0f280d503d`
- `context/AGENTS.md` = `e1aae1c9899b9901c5d4fdf1b1fd952178557a51`
- `README.md` = `2e282ef77a484bdb33360730dd7f2b32682bf2cb`
- `.agents/docs/deploy/local.md` = `820a9a070b641246be03fc2c86c7568847c84ce8`
- `.agents/docs/deploy/remote.md` = `85629dd209398a795a297b017ec2985e7adc329b`
- `.agents/docs/specs/sh-agent/delivery.md` = `c4d0cf16643a5516a12692efd1a8c5c18527f66c`
- `.agents/docs/testing/sh-agent/guide.md` = `92fc106a33d353cd93e08e2159b9fa12433a76ff`
- `.agents/docs/plans/2026-09-17-delivery-rework.md` = `f4526ced8595d7d09bad0549aa80399edade530c`

Base commit for the diff: `012b4ca93838373bd1823672eee3a5a8fd584de2`.

## Verdict: APPROVE WITH NITS

This is D13's fix for the second security review: it strikes the per-value JSON walk
(D9's earlier design) in favour of whole-line/whole-document redaction, and hardens the
push to name the feedback URL literally rather than the `origin` remote name (closing a
`pushurl`/second-`url` redirect). `gofmt`, `go vet`, and the full `internal/harness` test
suite pass; the new tests (`TestSendReportPushesOnlyToFeedbackURL`, the widened-needle
cases, the structural-gap cases) are concrete and each one I probed does fail when the
mechanism it names is reverted — with one exception below. Docs (README, local.md,
remote.md, delivery.md, the plan's D9–D13 rows, the testing guide) were cross-checked
against the code and against each other; they agree, including a genuinely verified claim
in the testing guide (the `-run`/`-list` regex was run and does list exactly the seven
report-related tests). Nothing here is broken enough to block; one test's name overstates
what it covers, which the project's own hard rule treats as a description-vs-behaviour
defect.

## Findings

- **[low / confidence: high]** `internal/harness/report_test.go:110` — the case named
  `"needle behind a JSON unicode escape"` is `{"user_quote":"AKIAFAKEVALUE"}`, which
  contains no escape at all: `AKIA` sits verbatim in the raw line and is caught by the
  first `needleIn(body)` pass in `redactLines` (report.go:145), never reaching the new
  `decoded(body)` fallback (report.go:146-148) that the test's name and the surrounding
  comment ("the ways a needle was present... two from round 2") claim to add coverage for.
  Verified: I copied `internal/harness` into `.agents/scratch/mutant-test/`, replaced
  `decoded()`'s body with `return ""` (disabling the escape-unwrap fallback entirely), and
  `TestReportRedactStructuralGaps` still passed in full — this specific case does not
  detect a broken or absent `decoded()`. (A different, pre-existing case in
  `TestReportRedactShapes` — the JSON-nested `"config was {\"password\": ...}"` string —
  does depend on `decoded()` and does fail when it's disabled, so the function overall
  isn't untested; just this one named case is.) Per the project's own hard rule ("Claims
  of test COVERAGE are the sharpest case: name what is not covered"), fix by using an
  actual escape, e.g. `{"user_quote":"AKIAFAKEVALUE"}` or a token split by `\/`.

- **[low / confidence: low]** `internal/harness/report.go:586` — `commitAndPush` now
  pushes `HEAD:main` to the literal feedback URL regardless of the local branch's name.
  For the default path this is always correct (the code forces an unborn HEAD to
  `refs/heads/main` before the first commit, report.go:574-578, and every subsequent run
  reuses that same checkout). For a pre-existing checkout supplied via `--feedback-repo`,
  `verifyFeedback` checks only that the checkout is the top of its own tree and that
  `origin` matches the URL — it never checks the local branch name. A `--feedback-repo`
  checkout left on a branch other than `main` would still push its `HEAD` into remote
  `main` without error, silently disconnected from whatever the local branch name tracks.
  No test constructs this case (the `--feedback-repo` paths in `TestSendReport` and
  `TestResolveFeedbackFailures` all end up on `main` via a real `git clone`). Low
  real-world risk for a single-operator tool where the documented flow always leaves the
  checkout on `main`, so not blocking.

## Criteria

1. **Dead code** — pass. `walkJSON`/`redactJSON`/`Redaction.Path` were fully removed with
   no leftover references (`grep` confirms); no new dead code introduced.
2. **Flow tracing** — pass, with the one test-naming gap above. `Collect` →
   `redactLines`/`redactDocument` → `SendReport` → `commitAndPush` all agree on the new
   `([]byte, []Redaction)` signature (the dropped `error` return is legitimate: nothing in
   the new code path can fail, so keeping `error` would have been the dead code). The
   `verifyFeedback`/`commitAndPush` split correctly moves push-target trust to the literal
   URL argument, and `gitRedirects` env-var refusal is unchanged and still checked first
   in `SendReport`.
3. **Standards conformance** — pass. Stdlib only (no go.mod/go.sum diff), `gofmt -l` clean,
   `go vet ./...` clean. Grepped for stale restatements of "JSON key path" / redaction
   shape lists outside the reviewed doc set (skills, other docs) — none found beyond the
   files already in scope. The `context/AGENTS.md` asking-format change has no other
   restatement in this repo needing a matching edit (grepped `SHORT LIST`/`numbered table`
   across the tree; the other hits are unrelated generic usages).
4. **Directory layout** — pass, no new files, no misplacement.
5. **Simplicity / DIY-except-security** — pass. The rewrite is a net simplification
   (removes `walkJSON`, `redactJSON`, `sort`, `strconv`, `bytes` imports) with no crypto/
   auth surface touched.
6. **Decoupling** — pass. `marker`/`decoded`/`redactLines`/`redactDocument` are small,
   independently testable functions; `commitAndPush` taking `url` as an explicit parameter
   (rather than reaching for a package-level constant or the remote name) keeps the trust
   decision at the caller (`SendReport`, which already ran `verifyFeedback`).

Docs lens: `README.md`, `local.md`, `remote.md`, `delivery.md`, and the plan's D9–D13 rows
all describe the current whole-line/whole-document redaction, the literal-URL push, and
the no-key-path `report.json` schema consistently with the code and with each other. The
testing guide's claim that `-run 'TestReport|TestSendReport|TestResolveFeedback'` covers
every report-related test was independently verified with `go test -list` (exact set
match, 7 tests). Historical checklist/decision rows (D9, D11, plan line ~120) still
describing the old `-u origin HEAD`/JSON-walk behaviour are append-only history entries
correctly left standing, not live descriptions — they read as historical, not current.

## Not examined

- `.agents/docs/bugs/INDEX.md`, `.agents/docs/reviews/2026-09-17-security-6af54f7-015cf2d.md`,
  and the four new `BUG-053..056` files, plus the two new `2026-09-17-review-015cf2d-012b4ca-*`
  review files: all shown as changed/untracked in `git status` but not in the dispatch's
  named file list, so not diffed or checked for internal consistency here.
- No live end-to-end run against the real `styrman-harness-feedback` GitHub repository —
  only the local-bare-repository test suite (already passing) was exercised; the dispatch
  scope is a diff review, not the e2e tier the project's own testing guide assigns to the
  main agent.
- `cmd/sh-agent` and other `internal/harness/*.go` files outside the named diff were not
  read beyond grepping for stale references to removed symbols.

## Outcome

- Finding 1 (the "unicode escape" case carries no escape, so `decoded()` was untested there): fixed — a real defect. The fixture had been written through a tool that turned `A` into `A`; it now holds the six bytes `A` (`grep -n` shows them), and the test header names the case. Must-fail seen with `decoded()` stubbed to return "" in a scratch copy (`.agents/scratch/r2-decoded-mustfail.log`, the case fails); must-pass in the tree. `TestReportRedactShapes`' escaped-quote `password` case also fails without `decoded()`, as seen when this change was first written. Bug record dispatched.
- Finding 2 (`HEAD:main` from a non-`main` branch of a `--feedback-repo` checkout): rebutted — the feedback repository has one branch by design (D4), and the report is meant to land on `main` whatever the local branch is called; a branch that does not descend from `main` is rejected as non-fast-forward and the commit stays local, the same path as a behind checkout.
