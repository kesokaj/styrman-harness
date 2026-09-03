---
owner: reviewer
date: 2026-09-03
scope: changeset:318242dcea4a
---

# Review — changelog narrative exception   (2026-09-03)

Range reviewed: working-tree diff against `HEAD` (`e50533b`) — `git diff` over
`internal/trail/{owed.go,trail.go,scope_test.go}`, `REQUIREMENTS.md`, `TODO.md`,
`.agents/docs/architecture/DECISIONS.md`, `.agents/docs/testing/2026-09-03-verify-local.md`.
155 insertions / 34 deletions across 7 files (3 source files, 142 source lines).

## Verdict: REQUEST CHANGES

The narrative exception is broader than the ADR that justifies it, and it is
already firing wrongly on this repo's own trail. `hasScopedArtifactDated`
accepts *any* changelog artifact dated today whose scope is merely non-empty —
including entries committed hours earlier for unrelated tasks, and including the
obsolete `worktree:` scope form the `changeset:` mechanism replaced. Verified
live: `go run ./cmd/harness dispatch --dry-run` on this 142-line, 3-file change
set reports only `reviewer` owed and **no changelog row**, because
`.agents/docs/changes/2026-09-03-adr-006-008-implementation.md` (`scope:
worktree:9366002f2d81`, a landed and committed task) satisfies it. The doc
comment at `trail.go:236-241` claims that guard specifically prevents "entries
written before the mechanism existed" from satisfying a row; the entry now
satisfying it is exactly such an entry. Secondary to that, `REQUIREMENTS.md` —
declared immutable without operator approval — was edited by this change, and a
`TODO.md` item recording a known unreviewed commit was deleted without the
review being produced. Everything else is clean: `gofmt`, `go vet ./...` and
`go test ./...` all pass, the tests added are well-targeted, and the commit
boundary genuinely does still catch the original bug.

## Findings

- **[high / confidence: medium]** `REQUIREMENTS.md:80-85` — the file's own header
  says "Immutable without the operator's explicit approval. Agents never edit
  this file", and the global constitution repeats it. This diff edits R7. The
  text itself is accurate to the new behaviour, so this is a process violation
  rather than a wrong statement: nothing in the trail (no handoff, no changelog,
  no note in ADR-009) records the operator approving it. Fix: get and record
  explicit approval, or revert the R7 paragraph and keep the rationale in
  ADR-009 only. Confidence is medium because approval may have been given in a
  turn that left no artifact.

- **[medium / confidence: high]** `internal/trail/owed.go:118-124`,
  `internal/trail/trail.go:242-249` — the exception's guard is `a.Scope != ""`,
  which is a date match wearing a scope check. Any today-dated changelog entry
  satisfies the row for any in-progress work, including an entry for a
  completely different task that has already landed, and including a scope in a
  form this binary no longer produces. Demonstrated above on the current tree:
  the changelog row for this change set is reported as satisfied by
  `worktree:9366002f2d81`. Concrete fix: restrict the exception to an entry that
  plausibly belongs to the *current* task — require the entry to be uncommitted
  relative to `HEAD` (an entry written for work in progress is by definition not
  yet committed; `git status --porcelain -uall -- <docs>/changes` already gives
  this), and additionally require `strings.HasPrefix(a.Scope, "changeset:")` so
  a pre-mechanism scope form cannot satisfy anything. That preserves every case
  ADR-009 argues for — the entry you just wrote survives your next keystroke —
  while closing the "yesterday's task covers today's" hole.

- **[medium / confidence: high]** `internal/trail/check.go:60-75` (unchanged file,
  behaviour changed by this diff) — `CheckFreshness` derives `changelogDue` from
  `Owed()`, so it inherits the exception. Its own doc comment states it
  deliberately answers a *wider* question than the Stop gate and reports on the
  working tree; it now prints `OK: a changelog entry covers the current change
  set` for a dirty tree that no entry covers. `/audit` and a mid-work
  `/preflight` will therefore report FRESH on an uncovered tree. Either pass a
  flag so the freshness path evaluates rows without the narrative exception, or
  state the divergence in the comment. Impact is bounded — at the commit
  boundary the tree is clean and the exception lapses — but the comment and the
  behaviour now disagree.

- **[medium / confidence: high]** `TODO.md:3-17` — the deleted item "Reviewer row
  owed for the delta after review round 5" is not resolved by anything in this
  diff. No artifact in `.agents/docs/reviews/` declares `scope:
  changeset:713c8dc5e016`, the change set that entry names (it is the scope on
  `changes/2026-09-03-stop-gate-scope-unification.md` and on the previous
  verification log). Deleting the item closes the gap on paper only. Fix:
  restore the item, or dispatch the reviewer against that commit and cite the
  artifact.

- **[low / confidence: high]** `internal/trail/trail.go:52` and `:65` —
  `Resumability()` and `NarrativeRow()` are byte-identical predicates
  (`r == RowChangelog`), and `Resumability()` has no callers anywhere in the
  repo (`grep -rn Resumability` finds the definition and one unrelated comment).
  This diff added a second copy of an already-dead predicate. Delete
  `Resumability()`, or define it as the single source and have `NarrativeRow()`
  express its own meaning — the two will diverge the day a handoff row becomes a
  resumability row.

- **[low / confidence: high]** `internal/trail/owed.go:118` — `dirty :=
  len(p.DirtySource()) > 0` shells out to `git status --porcelain -uall` a
  second time; `owed.go:96` already made that exact call, and `p.Scope()` at
  `:105` makes a third. `Owed()` runs on every hook invocation, and R2 justifies
  the whole Go rewrite on hook latency (~2.6 ms vs ~22.8 ms). Capture the slice
  once: `dirtyFiles := p.DirtySource()` at :96, then `dirty := len(dirtyFiles) >
  0`.

- **[low / confidence: medium]** `internal/trail/scope_test.go:186-215` — the new
  tests pin the same-task case and the commit boundary, but nothing pins the
  guard the implementation comment claims: that a foreign or obsolete scope must
  not satisfy the exception. Had such a test existed it would have failed on
  this repo's own trail. Add a case writing `changes/other.md` with `date:
  today` and `scope: worktree:deadbeef` and assert the changelog is still due.

## Criteria

1. **Dead code — fail.** `Row.Resumability()` is unreferenced and now duplicated
   verbatim by `Row.NarrativeRow()`. No commented-out blocks or orphaned files
   otherwise; the `TODO.md` deletions are intentional (one of them prematurely).
2. **Flow tracing — fail.** The new predicate flows correctly through the one
   call site in `Owed()`, and all three consumers (`gate/trailgate.go:31,94`,
   `dispatch/dispatch.go:63`, `trail/check.go:64`) share that computation as
   designed — but two of them inherit a relaxation they were not audited for:
   `CheckFreshness` now reports OK on an uncovered dirty tree. The Stop gate is
   correctly unaffected, since `dirty` is false on a clean tree.
3. **Standards conformance — pass.** `gofmt -l` clean, `go vet ./...` clean, `go
   test ./...` green. Comments explain *why*, exported identifiers documented,
   the helper is unexported and method-scoped, no stutter, error-free paths only.
   The one conformance concern is the `REQUIREMENTS.md` edit under Findings.
4. **Directory layout — pass.** Changes land in `internal/trail/` beside the code
   they extend, the ADR is appended (not rewritten) to
   `.agents/docs/architecture/DECISIONS.md`, and the tests sit in the existing
   `scope_test.go` next to the mechanism they cover.
5. **Simplicity (DIY-except-security) — pass with a caveat.** No new
   dependencies, stdlib only, no security surface touched (no auth, crypto,
   secrets, or session handling in the diff — the `security` row is correctly
   reported as not due). The caveat is that the simplest thing was reached for
   where a slightly less simple predicate was needed: "any scoped artifact from
   today" is the coarse date rule R7 exists to replace, reintroduced behind a
   `dirty` check.
6. **Decoupling — pass.** `NarrativeRow()` is a property of `Row`,
   `hasScopedArtifactDated` a private method on `Project`, and the exception is
   applied at exactly one place in `add()`. The seam is clean enough that fixing
   Finding 2 touches one function body.
