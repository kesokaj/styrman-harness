---
owner: reviewer
date: 2026-09-03
scope: changeset:f6f6c28eee97
---

# Review — narrative row exception (final state)   (2026-09-03)

Range reviewed: the working-tree diff present at review start — `git diff` over
`internal/trail/{check.go,owed.go,scope.go,trail.go,scope_test.go}`, `REQUIREMENTS.md`,
`ROADMAP.md`, `TODO.md`, `.agents/docs/architecture/DECISIONS.md`,
`.agents/docs/testing/2026-09-03-verify-local.md` (10 files, 280 insertions / 57
deletions). That tree was committed mid-review; the same content is now
`e50533b..73bb93c` (`24323ab` for the code and docs, `73bb93c` for the changelog entry).
Verified independently: `gofmt -l .` clean, `go vet ./...` clean,
`go test -race -count=1 ./internal/trail` PASS, `harness doctor` 10/10.

## Verdict: APPROVE WITH NITS

Nothing here is functionally broken. The narrative exception is correctly bounded:
`commitBoundary` returns early on a dirty tree, so `hasInProgressEntry` can never reach
the one blocking gate, and `TestChangelogExceptionLapsesAtTheCommitBoundary` pins that.
The freshness fix is right — `CheckFreshness` asking coverage directly instead of
inheriting the exception removes a report that printed OK for a tree nothing covered.
The worst items are process and record-keeping, in a repo whose product *is* the process:
the whole feature landed under the commit message `asd`, and the verification artifact
that satisfies this change set records only cached test output, so it cannot be
distinguished from a hand-retargeted scope line.

## Findings

- **[medium / confidence: high]** `24323ab` — the entire feature (code, tests,
  `REQUIREMENTS.md`, ADR-009, 416 insertions) is committed under the message `asd`. The
  constitution requires descriptive commit messages; the git history is the only record
  that survives when the docs trail is read selectively. History is already built on top,
  so do not rewrite it: add a `git notes` entry, or accept that
  `.agents/docs/changes/2026-09-03-adr-009-narrative-row-exception.md` is now the sole
  narrative for that commit and say so in the next entry. Do not repeat it — a placeholder
  message on a commit that changes gate semantics is the failure mode ADR-009 itself was
  written to prevent.

- **[medium / confidence: medium]** `.agents/harness.json:19` +
  `.agents/docs/testing/2026-09-03-verify-local.md` — the verify command is
  `go test -race ./...`, which accepts Go's build cache. The artifact for scope
  `changeset:f6f6c28eee97` shows `(cached)` for *every* package including
  `internal/trail`, whose sources changed in this very diff; a first run at this code
  state cannot be cached. Either dispatch ran twice at one state (benign) or the file's
  scope line was retargeted by hand over an older run's body (not benign — that is exactly
  what R7 exists to prevent, and the diff for this file is the scope line and nothing
  else). Fix at the root rather than by trust: use `go test -race -count=1 ./...` in
  `verify`, so a verification log is always evidence of a run against the code it claims,
  and consider stamping the wall-clock time of each command into the artifact.

- **[medium / confidence: high]** `.agents/docs/architecture/DECISIONS.md` (ADR-009,
  Decision paragraph) — the ADR states the rule as "while the working tree is dirty, the
  changelog row is satisfied by any entry from today that declares some scope". The
  implementation is deliberately stricter: `hasInProgressEntry` also requires the entry to
  be *uncommitted* and to carry the `changeset:` prefix, and
  `TestNarrativeExceptionRequiresCurrentTaskEntry` asserts that the looser rule is wrong.
  The ADR now misdescribes the shipped mechanism, and it is committed and append-only, so
  amend by appending a correction paragraph (or an ADR-010 superseding note) — never by
  editing the text in place. Same wording gap in `REQUIREMENTS.md` R7 ("today's entry
  stands"), which is the operator-approved text and should be left alone until it next
  needs approval.

- **[low / confidence: high]** `internal/trail/trail.go:248` (`hasInProgressEntry`) — the
  exception matches *any* uncommitted, today-dated, `changeset:`-scoped entry in
  `changes/`, not one tied to the current task. Two unrelated tasks in one uncommitted
  session therefore share a single entry: task A's in-progress entry silences the
  changelog notice for task B until a commit lands. The blocking gate bounds the damage
  (after the commit the entry is tracked and the exception lapses), and the harness cannot
  compute "same task" from a diff — so this is a documented limit rather than a bug. Say
  so explicitly in the function comment, or add it to the **Known limits** section of
  `REQUIREMENTS.md` at the next approved edit, so the next reader does not read it as an
  oversight.

- **[low / confidence: high]** `internal/trail/trail.go:52` — `Row.Resumability()` has no
  non-test caller anywhere in the tree (only a comment in `owed.go:131` mentions the
  concept), and the new `Row.NarrativeRow()` has an identical body (`r == RowChangelog`).
  Two exported predicates with the same implementation and different meanings will drift
  the first time a second resumability row exists. Delete `Resumability()` — the contract
  it documents lives in `REQUIREMENTS.md` R5 — or, if it is kept as documentation, state
  in its comment that nothing calls it.

- **[low / confidence: high]** `internal/trail/owed.go:117` — `dirty := len(p.DirtySource()) > 0`
  re-shells `git status --porcelain -uall` after `rep.Changed = p.DirtySource()` already
  ran it eleven lines earlier. `Owed` runs on every prompt and every stop hook, where R2
  budgets in milliseconds. Capture the first result (`dirtyFiles := p.DirtySource()`;
  `dirty := len(dirtyFiles) > 0`) before the `HeadSourceToday` fallback overwrites
  `rep.Changed`.

- **[low / confidence: medium]** `internal/trail/trail.go:266` (`uncommittedArtifacts`) —
  the porcelain parser strips surrounding quotes but does not undo git's C-style escaping,
  so an artifact filename with a non-ASCII or quoted character will not match
  `filepath.Base(a.Path)` and the exception will silently not apply to it. `DirtySource`
  has the same shape, so this is consistent rather than novel. If it is worth closing, add
  `-c core.quotePath=false` (already the pattern used for `core.fsmonitor` in `git()`) or
  parse `-z` output.

- **[low / confidence: medium]** process — the session's own hook reported
  `reviewer: 263 changed lines ≥ threshold 50` as owed, and the change was committed
  (`24323ab`, `73bb93c`) before an artifact covering `changeset:f6f6c28eee97` existed. The
  prior review on disk covers `changeset:318242dcea4a`, an earlier state. This report
  closes the row after the fact; the DoD says the row is owed *before* the commit. Worth
  noting only because the repo's product is that contract.

## Criteria

1. **Dead code — fail (low).** `Row.Resumability()` is unreferenced outside tests and
   duplicates `NarrativeRow()`. Nothing else is orphaned: `LoadConfig` dropped from
   `CheckFreshness` left no unused import (vet clean), the removed `TODO.md` items are
   genuinely settled by this change, and no commented-out block was left behind.
2. **Flow tracing — pass.** The three `Owed` callers (`gate.TrailNotice`,
   `gate.commitBoundary`, `dispatch`) all still typecheck against the unchanged `Report`;
   `commitBoundary` short-circuits on a dirty tree, so the new exception is unreachable
   from the blocking gate — the property the comment in `owed.go` claims, verified by
   reading `trailgate.go:84` and by `TestChangelogExceptionLapsesAtTheCommitBoundary`.
   `CheckFreshness` no longer consumes `Report`, and its remaining consumer
   (`cmd/harness/main.go:111`) is unaffected. `Scope()` remains the single scope
   computation shared by gate, dispatcher and freshness. No stub is called.
3. **Standards conformance — pass with nits.** Idiomatic Go; comments explain *why*, which
   is this repo's house style; `gofmt`/`vet` clean; tests are table-free but explicit and
   each documents the regression it pins. Nits: the double `DirtySource()` call and the
   duplicate predicate above. `REQUIREMENTS.md` was edited, which agents may not do
   unilaterally — but ADR-009 records the operator's explicit "Approve both changes", so
   the rule is satisfied and the trace is on disk where a later reader can find it. The
   `asd` commit message violates the descriptive-commit rule.
4. **Directory layout — pass.** Code stayed in `internal/trail/`, tests beside their
   subject in `scope_test.go`, the ADR appended to `architecture/DECISIONS.md`, the
   changelog entry in `changes/`, the verification log in `testing/`. Nothing new outside
   the documented layout in `AGENTS.md`.
5. **Simplicity, DIY-except-security — pass.** No new dependency; stdlib and `git`
   subprocesses only, consistent with the rest of the package. The exception is two
   predicates and one helper — about the smallest shape that expresses the rule. No
   security surface is touched: no auth, crypto, secrets, or input parsing changed, and
   nothing is hand-rolled that should be a vetted library.
6. **Decoupling — pass with a caution.** `CheckFreshness` now computes coverage itself
   (`scope != "" && p.HasArtifact(...)`) rather than reading it off `Owed`. That is the
   right call — inheriting the exception was the bug — but it re-creates, in mild form,
   the "two functions that must agree will eventually not" hazard that the previous round
   removed. Consider a single `func (p Project) CoversCurrent(row Row) bool` used by both,
   with the exception applied only in `Owed`; then the shared half cannot drift.
