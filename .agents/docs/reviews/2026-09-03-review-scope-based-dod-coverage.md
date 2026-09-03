---
owner: reviewer
date: 2026-09-03
scope: worktree:9366002f2d81
---

# Review — scope-based DoD coverage   (2026-09-03)

What this covers is the `scope:` key in the frontmatter, not a line in this body:
the uncommitted working tree at `worktree:9366002f2d81` — 26 files, +349/−154, on
`main` (4 commits ahead of `origin/main`). Establishing command: `git diff` plus the
four untracked files (`internal/trail/scope.go`, `internal/trail/scope_test.go`,
`internal/dispatch/integration_test.go`, `internal/hookio/hookio_test.go`).

## Verdict: REQUEST CHANGES

The scope mechanism does not survive the commit boundary it exists to gate. I
reproduced it against this tree's own code in a scratch copy: dispatch a change
fully (every artifact declares `worktree:<digest>`), commit it, and `Owed` reports
`changelog` and `verifier` due again — because the scope has become
`commit:<sha>` and nothing bridges the two forms. The dispatch → commit path is the
only path the constitution sanctions, so on today's HEAD every fully-satisfied row
comes due a second time immediately after the commit, and a second dispatch would be
invalidated by the next commit in the same way. Underneath that, `worktreeDigest`'s
`git diff HEAD` call always fails at exit 128 — the repo's own `git` wrapper passes
`-c diff.external=`, which git tries to execute as a program — so the diff never
contributes to the digest, the error is swallowed, and both the code comment and
ADR-007 describe a computation that does not run. The rest of the change is good
work: the stub-CLI integration tests, the `WaitDelay` fix, the merge-commit diff
fix, and the `hookio` split are all clean, tested, and well-motivated.

## Findings

- **[high / confidence: high]** `internal/trail/scope.go:27` — `Scope()` returns
  `worktree:<digest>` before a commit and `commit:<sha>` after it, and `Covers`
  relates the two forms only through a `<base>..<head>` range that `harness dispatch`
  never writes. Reproduced: seed a repo, add `feature.go`, write `changelog` and
  `reviewer` artifacts declaring `proj.Scope()` (rows go quiet), `git commit`, then
  `Owed` returns `[changelog, verifier]` again — scope moved
  `worktree:19806b7baba7` → `commit:39b0966fbe1e`. The gate now re-owes every row for
  work it just accepted, which is exactly the false-positive cost R4 warns about, and
  it self-perpetuates until `headIsToday()` decays at midnight. Fix: make the commit
  form reachable from the pre-commit form — e.g. when the tree is clean and HEAD is
  today, also accept an artifact whose scope equals the worktree digest recomputed
  over HEAD's own diff; or have `harness dispatch` record the digest→sha mapping it
  alone knows. Whichever, add the missing regression test: "a fully dispatched change
  set stays satisfied across the commit that lands it."
- **[medium / confidence: high]** `internal/trail/scope.go:52` — `git(p.Root,
  "diff", "HEAD")` can never succeed. The wrapper at `internal/trail/git.go:16`
  injects `-c diff.external=`, and git runs the empty string as an external differ:
  `error: cannot run : No such file or directory` / `fatal: external diff died`,
  exit 128, whenever there is any textual diff to render (verified on git 2.47.3;
  `--numstat` and `--name-only` are unaffected, which is why no existing caller
  noticed). The `err == nil` guard swallows it, so the digest is built purely from
  the sorted paths and contents of `DirtySource()`. It still changes when source
  changes — the tests pass for the right reason by accident — but the intent is dead
  code, and the comment at line 42 and ADR-007's "a digest over the tracked diff plus
  the content of untracked source files" both describe something that never executes.
  Fix: `git diff --no-ext-diff HEAD` (the flag must follow the subcommand), or delete
  the call and say plainly that the digest is content-based. Note that repairing it
  has a side effect worth deciding deliberately: the diff includes `.agents/docs/**`,
  so the verifier rewriting a *tracked* `testing/<day>-verify-*.md` mid-dispatch would
  then move the scope out from under the artifacts already written against it.
- **[medium / confidence: high]** `internal/gate/subagent.go:71-75` — the subagent
  stop gate still instructs a held agent to write frontmatter with "`owner: %s` and
  `date: %s`", with no `scope:`. Under the new `Covers` rule an artifact with no scope
  covers nothing, so the gate's own remediation produces an artifact that cannot
  satisfy the row it is holding the agent for. `proj` is in hand at that call site;
  include `proj.Scope()` in the message. (`hasFreshArtifact` matching on `a.Date` is
  fine — freshness within a run genuinely is a date question.)
- **[medium / confidence: high]** `shared/AGENTS.md:27` — the constitution, which R1
  names as the single source and from which `claude/CLAUDE.md` and
  `antigravity/GEMINI.md` are generated, still reads "Every artifact carries
  frontmatter `owner:` and `date:`". R7 in this same diff now mandates three keys and
  all five agent definitions were updated. CI cannot catch this (it is semantic, not
  byte drift — I ran `harness derive` over all five agents and confirmed the agy
  copies are in sync), so the main agent writing a handoff or an inline docs-only
  changelog will keep producing scope-less artifacts that cover nothing.
- **[high / confidence: medium]** `REQUIREMENTS.md` — modified (+43/−6: R7 rewritten,
  a "Known limits" section added). The file's own header says "Immutable without the
  operator's explicit approval. Agents never edit this file," and `AGENTS.md` and the
  global constitution repeat it. The content is accurate and improves the document, so
  the severity is procedural, not technical: confidence is medium only because I
  cannot see from the tree whether the operator approved it. If they did, ignore this;
  if not, it needs their sign-off before the commit.
- **[medium / confidence: medium]** `internal/trail/scope.go:83` — range coverage
  compares `strings.TrimPrefix(scope, "commit:")` against the literal head text, so it
  matches only when the head is spelled exactly as `Scope()` spells it
  (`git rev-parse --short=12`). `shared/skills/preflight/SKILL.md:42` tells the agent
  to write `scope: <base>..<head>`, where `<head>` in practice is `HEAD`, a tag, or a
  7-char abbreviation — none of which will ever match. `scope_test.go:90` passes only
  because the fixture hand-picks matching literals. Resolve both ends through
  `git rev-parse --short=12` before comparing, and say so in the skill.
- **[low / confidence: high]** `internal/dispatch/dispatch.go:174,228` — `runAgentRow`
  and `rowPrompt` now take `(… day, scope string)`, two adjacent untyped strings whose
  transposition compiles silently and would write the date into the scope key. A named
  `type Scope string` in `trail`, or one struct carrying both, removes the hazard for
  free.
- **[low / confidence: medium]** `TODO.md:3` / `ROADMAP.md:5` — "Nothing open. Nothing
  planned. Nothing in flight" is asserted over a tree with 26 uncommitted files, four
  unpushed commits, and (per this review) an unresolved coverage bug. The old TODO's
  "Decide the coverage model — the largest open design question" is deleted as settled;
  ADR-007 settles the design but not the commit boundary. Re-open it, or reword.

## Criteria

1. **Dead code** — fail. The `git diff HEAD` call in `worktreeDigest` executes but can
   never contribute (finding 2); everything else is clean. No orphaned files, no
   commented-out blocks, `EmitAgyContext`/`AgyContextJSON`/`HasArtifactDated` each have
   live callers, and `Artifact.Date` is still used by `head.go`, `check.go` and
   `gate/resume.go`.
2. **Flow tracing** — fail. `HasArtifact`'s meaning changed from date to scope and
   every in-tree caller was updated correctly (`owed.go:112`, `dispatch.go:201`,
   `check.go:61` → `HasArtifactDated`), but the flow breaks at two seams the compiler
   cannot see: the worktree→commit transition (finding 1) and the stop gate's
   remediation text (finding 3). `go vet ./...` and `go test ./...` both pass.
3. **Standards conformance** — mostly pass. Idiomatic Go, errors returned not
   panicked, comments explain *why* throughout and are unusually good. Two documents
   now contradict the code or each other: the constitution (finding 4) and ADR-007's
   description of the digest (finding 2).
4. **Directory layout** — pass. `scope.go` beside the trail logic it extends,
   `integration_test.go` in `internal/dispatch`, `hookio_test.go` in `internal/hookio`;
   the layout in `AGENTS.md` is respected and no new top-level file appears.
5. **Simplicity, DIY-except-security** — pass. Stdlib only, no new dependency
   (`go.mod` untouched), `crypto/sha256` for a non-security fingerprint is the right
   tool. Nothing security-sensitive is hand-rolled: no auth, crypto, hashing of
   secrets, sessions or tokens is introduced. The bash stub CLI is the cheapest way to
   test a contract that otherwise costs money per run — good judgement, and ADR-008
   records it.
6. **Decoupling** — pass with a nit. Scope computation lives entirely in `trail` and
   `dispatch` consumes it as data, so the two seams could be worked in parallel. The
   stringly-typed `(day, scope)` pair is the only place where that boundary is thinner
   than it looks (finding 7).

## What is good, and should not be lost in the rewrite

- `internal/dispatch/integration_test.go` — the stub CLI is the highest-value addition
  here. It covers argument shape, envelope parsing, cost capture, timeout enforcement,
  artifact detection, the "agent claimed success but wrote nothing" path, and the
  zero-cost verifier claim, with no credentials and no spend.
- `internal/dispatch/cli.go:71` — `WaitDelay` is a correct diagnosis of a real bug
  (`Output()` waiting on pipes inherited by a grandchild), with the measurement in the
  comment.
- `internal/trail/git.go:97` — `headDiff` turns a documented known-limitation about
  merge commits into a fix, with `TestMergeCommitIsVisible` building a real merge to
  prove it.
- `shared/skills/VENDORED.md` — the licensing gap is resolved honestly: attribution
  taken from each skill's own frontmatter, the fetch origin marked "not recorded"
  rather than guessed, and the Apache-2.0 `LICENSE`-file obligation named as a
  condition on future publication.
