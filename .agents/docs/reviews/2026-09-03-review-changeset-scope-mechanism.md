---
owner: reviewer
date: 2026-09-03
scope: changeset:09e0aa231bad
---

# Review — changeset scope mechanism   (2026-09-03)

What this covers is the `scope:` key in the frontmatter above, not a line in this body.
Diff established with `git diff` + `git status --porcelain -uall` on the working tree of
`styrman-harness`: 30 modified files (+360/−161) plus 6 untracked
(`internal/trail/scope.go`, `internal/trail/scope_test.go`,
`internal/dispatch/integration_test.go`, `internal/hookio/hookio_test.go`, and two
`.agents/docs/` artifacts). `gofmt -l .` and `go vet ./...` are clean; the verification
log at `.agents/docs/testing/2026-09-03-verify-local.md` carries this same scope and a
PASS verdict.

## Verdict: REQUEST CHANGES

The worst problem is that the change ships a scope format its own governing documents do
not describe. `REQUIREMENTS.md` R7 — the immutable file — was rewritten in this diff to
state that "a scope is `worktree:<digest>` for uncommitted work … or `commit:<sha>` for a
clean tree", and ADR-007 records the same as the accepted decision. The code emits
exactly one form, `changeset:<digest>`, and `internal/trail/scope.go:24` explicitly names
the two-form design as the bug it was written to remove. So R7 and ADR-007 now document a
rejected design as the contract, and a downstream reader (or the next agent) implementing
against them would reintroduce the regression `TestScopeSurvivesTheCommitThatLandsIt`
guards. Second, and independently: the one gate that actually *blocks* — `StopDecision`
via `CoversHead` — was never migrated, so the same-day coverage hole this whole change set
exists to close is still open at the boundary where it matters most. Everything else is
good work: the `WaitDelay` bound, the merge-commit diff, the stub-CLI dispatcher tests and
the `hookio` refactor are all real improvements with real tests.

## Findings

- **[high / confidence: high]** `REQUIREMENTS.md:73-74`, `.agents/docs/architecture/DECISIONS.md:117-118` —
  Both describe the scope format as `worktree:<digest>` / `commit:<sha>`; `Scope()`
  (`internal/trail/scope.go:27-41`) only ever returns `changeset:<digest>`, and its own
  doc comment says the two-form split "re-owed every row the moment the work was
  committed". The requirements file is the project's contract and is now false about the
  shipped binary; ADR-007's Decision paragraph likewise records the design that was
  rejected, including a rationale ("a digest over the tracked diff plus the content of
  untracked source files") that does not match `digestOf`, which hashes path + content for
  *every* changed source file, tracked or not. Fix: correct R7 and ADR-007 to
  `changeset:<digest>` — one form spanning both sides of the commit — and state the
  commit-survival property as the reason. Note also that `REQUIREMENTS.md` is marked
  "Immutable without the operator's explicit approval"; if that approval was not given for
  this edit, the edit itself is the violation and the file should be reverted before the
  correction is proposed.
- **[high / confidence: high]** `internal/trail/head.go:38-46`, `internal/gate/trailgate.go:85` and `:109` —
  `CoversHead` still answers "does the trail cover HEAD?" by comparing *dates*
  (`newest >= head`), and it is the sole input to the only blocking gate in the harness
  (R4). Commit A at 09:00 with its changelog, then commit B at 15:00 with none: `newest`
  and `head` are both today, so the gate reports covered and lets the stop through. That
  is precisely the failure `TestSecondChangeSameDayOwesRows` was written to prevent, left
  intact at the boundary that blocks. `internal/trail/check.go:61` inherits the same
  weakness through `HasArtifactDated`, and the block message still tells the agent to
  "pass today's date" rather than the scope. Fix: on a clean tree, decide coverage with
  `p.HasArtifact(RowChangelog, p.Scope())` and fall back to the date comparison only when
  `Scope()` is empty (HEAD not today, nothing to cover); keep `HasArtifactDated` for the
  genuinely date-shaped "was the trail touched today" note, and put `proj.Scope()` in the
  block message the way `subagent.go:79` now does.
- **[medium / confidence: high]** `internal/trail/scope.go:80-97` —
  The range branch of `Covers` is unreachable in production. It delegates to `sameCommit`,
  which returns false unless the *target* scope carries a `commit:` prefix, and no code
  path produces such a scope any more: `Owed`, `dispatch` and `gate` all pass
  `p.Scope()`, which is always `changeset:…`. Consequences: (a) `sameCommit` and the
  `strings.Cut` branch are dead code by the criterion that dead code must go; (b) the
  contract `shared/skills/preflight/SKILL.md:47` now advertises — a reviewer writing
  `scope: <base>..<head>` satisfies the rows for the commits inside the range — cannot be
  honoured by `harness dispatch`, so every release will re-dispatch the reviewer row it
  just paid for; (c) `TestRangeScopeResolvesRevisions` passes only because it fabricates
  `"commit:" + head[:12]` by hand, so the test asserts a path the system cannot reach and
  gives false confidence. Fix one of two ways: either teach `Covers` to resolve a range
  head against the current clean-tree change set (compare `rev-parse <head>` to HEAD when
  the target is a `changeset:` scope produced from `HeadSourceToday`), or delete the range
  branch, `sameCommit` and its test, and state in the preflight skill that range coverage
  is a judgement the gate model makes from frontmatter rather than something the binary
  computes.
- **[low / confidence: medium]** `internal/trail/scope_test.go:135-138` —
  The assertion is guarded by `&& spelling != "main" && spelling != "master"`, which
  silently skips two of the five spellings under test because the fixture's default branch
  name is unknown. A test that cannot fail for a case it enumerates is documentation, not
  coverage. Fix: resolve the branch once with `git symbolic-ref --short HEAD` in the
  fixture and assert on that single, known spelling.
- **[low / confidence: high]** `.agents/docs/changes/2026-09-03-adr-006-008-implementation.md:4`, `.agents/docs/reviews/2026-09-03-review-scope-based-dod-coverage.md:4` —
  Both carry `scope: worktree:9366002f2d81`, which under the rule this change set
  introduces covers nothing at all, and the changelog entry's body (lines 17-18) repeats
  the `worktree:`/`commit:` claim contradicted above. By the change's own contract the
  changelog row for `changeset:09e0aa231bad` is still unsatisfied — the DoD hook agrees.
  Fix: rewrite the changelog entry for the current scope (and correct its description of
  the format) before committing; the superseded review can stay on disk as history since a
  scope-less-to-this-change artifact is inert by design.
- **[low / confidence: medium]** `internal/trail/scope.go:52-62` —
  `digestOf` swallows read errors (`if data, err := os.ReadFile(…); err == nil`), so a file
  that is transiently unreadable hashes identically to an empty one. The realistic effect
  is small but it points the wrong way: an IO error can make a stale artifact appear to
  cover the current tree. Fix: on error, write a distinguishing marker byte into the hash
  rather than nothing, so an unreadable file can never collide with an empty file.
- **[low / confidence: medium]** `internal/dispatch/dispatch.go:228-232` —
  The prompt now says "YAML frontmatter containing exactly these three keys". The
  harness's own verification log writes a fourth (`verdict: PASS`,
  `internal/dispatch/verify.go:146-150`), and preflight reports want a verdict key too. An
  agent taking "exactly" literally would strip useful metadata. Fix: "containing at least
  these three keys" — the emphasis wanted is on copying the scope verbatim, which the next
  sentence already carries.

## Criteria

1. **Dead code — fail.** `sameCommit` and the range branch of `Covers`
   (`internal/trail/scope.go:80-97`) cannot be reached with any scope the binary produces;
   removed code elsewhere (the old `readFrontmatter` signature, the inline `diff-tree`
   invocations, the duplicated `EmitAgyContext` body) was cleaned up properly, and no
   orphaned files or commented-out blocks were left behind.
2. **Flow tracing — fail.** The `day → scope` threading through `Run` → `runRows` →
   `runAgentRow`/`runVerifier` → `rowPrompt`/`writeVerificationLog` is complete and
   type-consistent, and `HasArtifactDated` correctly preserves the one date-shaped caller —
   but the trail's *blocking* consumer (`CoversHead` → `StopDecision`/`AgyStopDecision`)
   was left on date semantics, so the new mechanism stops short of the gate it most needed
   to reach, and the range-scope path leads into a function nothing can satisfy.
3. **Standards conformance — partial.** `gofmt -l .` and `go vet ./...` are clean, errors
   are handled at every new call site, comments explain *why* (the `WaitDelay` note and the
   `headDiff` git-quirk note are exemplary), and the new tests are table-driven and
   named for the regression they pin. Against project rules, the failure is documentary:
   `REQUIREMENTS.md` R7 and ADR-007 do not describe the code that ships, which is the
   high finding above.
4. **Directory layout — pass.** `scope.go`/`scope_test.go` sit beside the rest of
   `internal/trail`; `integration_test.go` is in the package it exercises; `hookio_test.go`
   fills the one package that had no test; every artifact lands under `.agents/docs/<class>/`
   as the contract requires.
5. **Simplicity, DIY-except-security — pass.** No new dependencies; `crypto/sha256`,
   `encoding/hex` and `sort` are stdlib, and the digest is a content fingerprint rather than
   a security primitive, so the DIY inversion does not apply. Nothing in the diff touches
   auth, crypto, password hashing, sessions or JWT/OIDC. The stub-CLI harness is a shell
   script written by the test itself into `t.TempDir()`, which is simpler than a fake
   binary package and is the right call.
6. **Decoupling — pass, with one note.** `trail` owns the scope, `dispatch` and `gate`
   consume it, and neither reaches into the other; the `hookio` split into
   `AgyContextJSON` + a one-line `EmitAgyContext` is exactly the seam that made the
   Antigravity shape testable. The note: `dispatch` computes the scope once for the prompt
   while `gate.SubagentStopDecision` recomputes it independently at stop time, so any
   source edit in between silently invalidates every artifact the dispatch produced. The
   preflight skill documents the mitigation ("make no further edits between running
   `harness dispatch` and its rows landing"), and `.agents/**` being excluded from
   `IsSourcePath` keeps the artifacts themselves from moving the scope — but the
   dependency on a frozen tree is real and belongs in the ADR's Consequences.
