---
owner: reviewer
date: 2026-09-03
scope: changeset:af2dcaa6c4d3
---

# Review — scope coverage across the gates, merge diffs, CLI stub (2026-09-03)

Change set reviewed: the uncommitted working tree at `6873ce7` — 33 modified files plus
10 untracked (`git diff` + `git status --porcelain -uall`), 1002 changed source lines by
the harness's own count. Substance: `internal/trail/scope.go` (new), scope-based coverage
in `trail`/`gate`/`dispatch`, merge-commit-aware `headDiff`, `cmd.WaitDelay`, the CLI stub
integration tests, `hookio` tests, and the doc/agent/skill updates that follow ADR-007/008.

## Verdict: REQUEST CHANGES

The blocking gate and the thing that is supposed to satisfy it disagree about which scope
counts. `CoversHead()` decides with `HeadScope()`, but every message the gate prints — and
`Owed()`, which is what `harness dispatch` writes artifacts against — uses `Scope()`, which
prefers the working tree. Whenever source is dirty at a stop, the gate blocks demanding
HEAD's digest while instructing the agent to write the worktree's digest. I reproduced it
on a scratch repo: following the block message verbatim produced a byte-identical block on
the next stop, and `harness dispatch` — the message's other recommendation — cannot clear
it either. That is the harness's one blocking gate telling an agent to do something that
provably does not work, which is precisely the false-positive cost R4 exists to prevent.
Everything else here is good work: the content-digest design is sound, the regression tests
are the right ones, the merge-commit fix closes a real hole, and the stub finally makes the
CLI contract testable without spend.

## Findings

- **[high / confidence: high]** `internal/gate/trailgate.go:93`, `:135`, and
  `internal/trail/check.go:50` — coverage is decided by `proj.CoversHead()` (which uses
  `HeadScope()`, `internal/trail/head.go:52`) but every remediation string interpolates
  `proj.Scope()`, which returns the *worktree* digest whenever `DirtySource()` is non-empty.
  Reproduced: repo with one commit lacking a changelog, plus one uncommitted `.go` file →
  block names `changeset:61dd68990bf0`; the gate actually requires `changeset:cf742af6275a`;
  writing the entry exactly as instructed yields the same block again. `harness check` prints
  the same wrong scope ("an entry must declare `scope: …`"), so `/preflight` and `/audit`
  inherit it. `Owed()` (`internal/trail/owed.go:106`) also satisfies rows against `Scope()`,
  so the dispatcher cannot satisfy the blocking gate on a dirty tree either.
  Fix: interpolate `proj.HeadScope()` in all three messages, and either have `CoversHead()`
  accept a changelog artifact carrying `HeadScope()` **or** `Scope()`, or have `Owed()` add
  the HEAD scope as a second satisfier for `RowChangelog`. Whichever is chosen, one test
  should assert "do what the block message says → the block clears".

- **[medium / confidence: high]** `claude/agents/reviewer.md:15`, `claude/agents/verifier.md:25`,
  `shared/AGENTS.md:21` (and the generated `claude/CLAUDE.md:21`,
  `antigravity/GEMINI.md:21`, `antigravity/plugins/harness/agents/{reviewer,verifier}.md`) —
  the old body-line convention survived in the very files the change set was meant to clean.
  reviewer.md:15 still says "state it on the report's `Scope:` line — `/preflight` matches
  release ranges against that line, so a report without one covers nothing", while the report
  template ten lines below now says the opposite ("not a line in the body") and
  `shared/skills/preflight/SKILL.md` has moved to `Range reviewed:`. The constitution still
  describes a `/code-review` report saved "with only `Source:` and `Scope:` lines prepended",
  which under the new rule carries no frontmatter scope and therefore satisfies nothing —
  a DoD row documented in a form that cannot be satisfied. ADR-007's stated goal was one
  mechanism; three files still ship two. Fix: rewrite those lines to `Range reviewed:` /
  frontmatter `scope:` and regenerate (`harness install --dry-run`).

- **[medium / confidence: high]** `.agents/docs/testing/2026-09-03-verify-local.md` — a tracked
  verification artifact was hand-edited: `scope: changeset:af2dcaa6c4d3` added, and a captured
  output line changed from `? … internal/hookio [no test files]` to `ok … internal/hookio
  (cached)`. The body of that same file asserts it was "Produced by `harness dispatch`, which
  ran the project's declared commands directly. No model was involved." The verifier row for
  this change set is thus satisfied by a retro-fitted record rather than a run. The *claim* is
  true — I ran `gofmt -l .` (clean), `go vet ./...` (clean) and `go test ./...` (all packages
  ok) independently — but the artifact is no longer an honest record, and the new constitution
  text ("this applies to anything you write into the trail by hand") is what makes that a
  contract violation rather than a tidy-up. Fix: revert the edit and let `harness dispatch`
  write a fresh log for this scope.

- **[low / confidence: high]** `internal/trail/check.go:49` — `case !covered && newest == ""`
  is unreachable; the preceding `case newest == ""` catches every state it could match. Dead
  branch introduced by this change; delete it.

- **[low / confidence: high]** `TODO.md` — now reads "Nothing open", and the item "Push. Three
  commits are local-only; nothing has been published" was deleted while still true: `main` is
  ahead of `origin/main` by 4. A resumability document that states a false state costs the next
  zero-context agent more than an untidy list does. Either push, or keep the item.

- **[low / confidence: medium]** `internal/dispatch/dispatch.go:92` and `:108` — the
  invalidation condition is evaluated twice against two separate `proj.Scope()` calls, so a
  tree that moves between them prints INVALIDATED but returns nil, or returns the error without
  the explanation. Compute `after := proj.Scope()` once before the ledger write and reuse it.

- **[low / confidence: medium]** `claude/agents/documentation.md` — the only agent of the six
  that did not receive the `scope:` paragraph. Defensible (REQUIREMENTS's Known limits say the
  documentation row is not computed, and that agent is told never to stamp frontmatter into
  README/spec files), but it now means a hand-written documentation artifact covers nothing.
  Worth one explicit sentence in the agent file so the omission reads as a decision.

## Criteria

1. **Dead code — fail.** One unreachable `switch` case (`check.go:49`), and three stale
   `Scope:`-body-line instructions that are worse than dead: they actively contradict the new
   mechanism. No orphaned files; every new/changed export (`AgyContextJSON`, `HasArtifactDated`,
   `headIsMerge`, `HeadScope`, `Covers`) has a live caller; `Config.SourceGlobs` is still read
   by `LoadConfig`→`Resolve`.
2. **Flow tracing — fail.** The scope threaded through `Run → runRows → runAgentRow/runVerifier
   → rowPrompt` and into the frontmatter is coherent, and `artifactPath`/`HasArtifact` were both
   moved to `Covers` so nothing calls the old date matcher by accident. The break is the
   `Scope()` / `HeadScope()` seam between `Owed`+gate messages and `CoversHead` (finding 1).
   Moving `SourceGlobs` loading into `Resolve` is the right fix for the classification split and
   is tested (`TestScopeIsIndependentOfCallerSetup`); the only surviving hazard is a `Project`
   built literally instead of via `Resolve`, and `internal/` has none outside tests.
3. **Standards conformance — pass.** `gofmt -l .` clean, `go vet ./...` clean, `go test ./...`
   green including the new `-race`-safe dispatch tests. Comments explain *why* throughout, which
   is the house style here. Generated artifacts are not drifted: I re-derived all six Antigravity
   agents and diffed them, and both `CLAUDE.md`/`GEMINI.md` are byte-identical to
   `shared/AGENTS.md` — R1 holds. Content of `REQUIREMENTS.md` was edited (R7 rewritten, "Known
   limits" added); that file is operator-owned and needs explicit approval — flagged for the
   operator to confirm rather than scored as a violation, since approval is not visible from here.
4. **Directory layout — pass.** `scope.go`/`scope_test.go` sit with the rest of `trail`;
   `integration_test.go` sits with the dispatcher it exercises; the stub script is generated into
   `t.TempDir()` rather than committed. Nothing landed outside the layout in `AGENTS.md`.
5. **Simplicity, DIY-except-security — pass.** `crypto/sha256` + `sort` from stdlib, `go.mod`
   untouched, no new dependency. No hand-rolled auth/crypto/session handling: the digest is an
   identity key, not a security control, so the 12-hex-char truncation is fine (collisions there
   cost a missed row, not a bypass — and only for an adversary who is already editing the tree).
   `cmd.WaitDelay` is the minimal correct fix for the pipe-inheritance hang, correctly explained.
   The stub CLI is a shell script executed from tests inside the microVM boundary R8 describes.
6. **Decoupling — mixed.** `trail` exposes coverage as `Scope`/`HeadScope`/`Covers` and both
   `gate` and `dispatch` consume it without reaching into internals, which is the right shape.
   But the boundary does not say which of the two scopes is authoritative for which caller, and
   finding 1 is exactly the bug that ambiguity produces — a doc comment on `Scope` naming its one
   legitimate caller set (the dispatcher and the informational notice) and on `HeadScope` naming
   the commit-boundary gate would keep two engineers off each other's toes here.
