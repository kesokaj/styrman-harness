---
owner: reviewer
date: 2026-09-03
scope: changeset:06cc436fce2d
---

# Review — scope-based coverage, gates, and merge/diff handling   (2026-09-03)

Change set reviewed: the uncommitted working tree at `6873ce7` — 33 modified and 9
untracked files, 896 lines by the dispatcher's count. Go: `internal/trail/scope.go`
(new), `trail.go`, `head.go`, `check.go`, `owed.go`, `git.go`, `internal/gate/{trailgate,
subagent}.go`, `internal/dispatch/{dispatch,verify,cli}.go`, `internal/hookio/hookio.go`,
plus three new test files. Prose: `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`,
`shared/AGENTS.md` (+ generated copies), the five agent definitions (+ derived),
`shared/skills/{preflight,VENDORED}`.

`go vet ./...`, `go test ./...` and `gofmt -l .` are all clean, and the generated
copies (`claude/CLAUDE.md`, `antigravity/GEMINI.md`, the derived agents) match their
sources — R1 holds. The problems below are behavioural, and two of them are live in
this repo right now.

## Verdict: REQUEST CHANGES

The scope mechanism is a genuine improvement over dates and its tests are good, but it
is wired in three ways that break the gate it was built to fix. The worst: the scope a
gate computes depends on whether that particular caller happened to load
`.agents/harness.json`, and three of the four gates do not. In this repo, right now,
`harness check` says FRESH against `changeset:06cc436fce2d` while `harness hook stop`
blocks against `changeset:a086d675cc96` — two different answers to the same question,
about the same tree, from the same binary. Second: `CoversHead()` no longer answers the
question its name and its callers ask; a single uncommitted source file makes it report
a fully recorded HEAD as uncovered, so the one blocking gate now blocks at every turn
end with a statement that is false — exactly the false-positive cost R4 exists to avoid.

## Findings

- **[high / confidence: high]** `internal/gate/trailgate.go:81,107` and
  `internal/gate/subagent.go:58` — `trail.Resolve()` is called without
  `trail.LoadConfig()`, so `proj.SourceGlobs` is empty when `proj.Scope()` runs.
  `Scope()` digests `DirtySource()`, which filters through `IsSourcePath`, which
  consults `SourceGlobs` — so the three stop gates hash a *different file set* than
  `Owed()` and `CheckFreshness()`, which do load the config. Any project that sets
  `sourceGlobs` (this one does, for exactly the reason R5 gives) gets two incompatible
  scopes. Reproduced live, unmodified tree, both binaries agreeing:
  `harness check` → `OK: newest changelog entry covers HEAD … RESULT: FRESH`;
  `harness hook stop` → `{"decision":"block", … scope: changeset:a086d675cc96}`.
  The consequence is not cosmetic: the changelog entry `harness dispatch` writes
  (`changeset:06cc436fce2d`) can never satisfy the blocking gate, and the entry the
  gate demands can never satisfy the dispatcher. The gate blocks permanently and its
  remediation text is wrong. Fix: make the classification part of the project rather
  than a field callers must remember to set — have `trail.Resolve` load
  `.agents/harness.json` and populate `SourceGlobs` (or have `Scope()`/`IsSourcePath`
  load it), and delete the two hand-assignments in `Owed` and `CheckFreshness`. A
  regression test should assert that `Scope()` on a repo with `sourceGlobs` is the
  same value whether or not the caller pre-loaded config.

- **[high / confidence: high]** `internal/trail/head.go:47-51` — `CoversHead()` asks
  `p.Scope()`, which prefers the *working tree* (`DirtySource()`) over HEAD. So it no
  longer reports whether the changelog covers HEAD; it reports whether the changelog
  covers whatever is on disk. Reproduced in a scratch repo: commit source plus a
  changelog carrying that scope → `OK: newest changelog entry covers HEAD`; then add
  one new `.go` file, touching nothing else → `STALE: no changelog entry covers the
  current change set (HEAD 2026-09-03, newest entry 2026-09-03)`, and `hook stop`
  returns `decision: block` with the text "HEAD was committed 2026-09-03 and no
  changelog entry covers it", which is false. This contradicts `StopDecision`'s own
  doc comment ("Uncommitted work in progress is reported by TrailNotice instead …
  the DoD is a commit contract, not a turn contract") and R4 ("blocks at exactly one
  boundary — a commit that the changelog does not cover"), so it is a REQUIREMENTS
  violation and blocking. Fix: give `Scope()` an explicit HEAD-scoped sibling — e.g.
  `HeadScope()` built from `HeadSourceToday()` only — and have `CoversHead` use that,
  leaving the worktree-preferring `Scope()` to `Owed`/`TrailNotice` where it belongs.
  Note the two are equal exactly when the tree is clean, which is why the new tests
  (`TestCoversHeadIsNotFooledByADateMatch` commits before asserting) do not catch it.

- **[medium / confidence: high]** `internal/trail/scope.go:46` with
  `internal/dispatch/dispatch.go:88,201` — the scope is recomputed from disk by every
  caller, but the dispatcher captures it once and then checks satisfaction against
  that captured value (`proj.HasArtifact(row, scope)`). Any source edit between the
  dispatcher computing the scope and the agents finishing silently invalidates the
  whole run: every artifact is written with a scope no gate will ever ask for again,
  and `harness dispatch` still reports the rows green. Nothing detects it — the
  preflight skill's "if the tree moves, re-dispatch" is discipline, not a check. Cheap
  fix: after the rows land, recompute `p.Scope()` and, if it differs from the captured
  one, report the run as invalidated rather than successful.

- **[medium / confidence: medium]** `shared/skills/preflight/SKILL.md:44-48` — the
  skill now reads `scope:` from frontmatter (good, that was the point), but it still
  specifies range and commit forms: "Tell the agent to write `scope: <base>..<head>`",
  and "an exact `commit:<sha>` match, or a range `<base>..<head>` whose head is HEAD".
  `Project.Covers` is exact-string-match only and `scope.go` documents the range form
  as deliberately removed; R7 says "A scope is `changeset:<digest>` — one form". So a
  preflight report written per this skill carries a scope no gate can match, and a
  reviewer report written per `claude/agents/reviewer.md` (verbatim copy of the
  dispatch prompt's `changeset:…`) is not recognised by the preflight range logic —
  the "two conventions side by side" item that `TODO.md` deletes as resolved. Pick
  one: either state in the skill that release-range coverage is a human judgement made
  by reading `changeset:` frontmatter, or re-introduce a range form in `Covers` with
  tests. Do not leave the two documents disagreeing.

- **[high / confidence: medium]** `REQUIREMENTS.md` — modified in this change set
  (R7 rewritten, R8 extended, a new "Known limits" section added). Its own header says
  "Immutable without the operator's explicit approval. Agents never edit this file",
  and the constitution repeats it. The content is accurate and matches the code, so if
  the operator approved this in-session there is no issue — but the approval is not on
  the record anywhere in the trail, and this is the one file where that has to be
  explicit. Confidence is medium only because I cannot see the session that produced
  it. Confirm the approval in the changelog entry, or revert.

- **[low / confidence: high]** `internal/trail/check.go:70` — `case !covered && newest
  == ""` is unreachable: `case newest == ""` two lines above already matched it. Dead
  on arrival; delete it.

- **[low / confidence: high]** `internal/gate/trailgate.go:96,135` — both stop gates
  tell the agent to hand-write an entry with a specific scope, and that string is
  obtainable *only* by triggering the block. Neither `TrailNotice` nor `harness
  dispatch --dry-run` prints it (verified: the dry-run output names rows and reasons,
  never the scope). Add the scope to the notice header and the dry-run output.

- **[low / confidence: medium]** `internal/trail/scope.go:76-79` — `digestOf` ignores
  `os.ReadFile` errors, so a listed-but-deleted file, an unreadable file, and an empty
  file all contribute the same bytes. Two different change sets can collide (delete
  `a.go` vs truncate `a.go`). Hash a distinguishing marker on the error path.

- **[low / confidence: high]** `internal/trail/scope.go:95` — `Covers` never uses its
  receiver; it is a pure function of `(Artifact, string)`. Fine as an API-symmetry
  choice, but worth a comment saying so, since a method that ignores `p` reads like an
  oversight. Related: `runAgentRow(proj, cfg, cli, row, day, scope string, extra
  ...string)` now has two adjacent, interchangeable string parameters that the compiler
  cannot tell apart — one transposition and every artifact is written with a date in
  its `scope:` key. A two-field `rowContext{day, scope string}` removes the class.

## Criteria

1. **Dead code** — fail (minor). One unreachable switch case in `check.go:70`.
   Otherwise clean: `HasArtifactDated` is genuinely used by the freshness check,
   `NewestChangeDate` still has callers, the removed `worktree:`/`commit:` scope forms
   and the old `readFrontmatter` signature left nothing orphaned, and the deleted
   `TODO.md`/`ROADMAP.md` sections point at real destinations rather than vanishing.
2. **Flow tracing** — fail. Callers of the changed signatures are all updated and the
   types agree (`runRows`/`runAgentRow`/`runVerifier`/`writeVerificationLog` threading
   `scope` is consistent, and `EmitAgyContext`'s split into `AgyContextJSON` is correct
   and now tested). The breaks are semantic: `Scope()` depends on caller-set state that
   three gates never set (finding 1), and `CoversHead()`'s meaning changed under
   callers that still ask the old question (finding 2). The merge-commit handling in
   `git.go` (`headDiff`/`headIsMerge`) is a real fix, correctly scoped to first-parent,
   and its test proves it.
3. **Standards conformance** — pass with the R4/R7 caveats above. `gofmt`, `go vet`
   and `go test` clean; comments explain *why* throughout and are unusually good; no
   `interface{}`-style looseness; errors are handled rather than swallowed except at
   the digest read noted above. R1 (generated copies) verified by direct comparison.
4. **Directory layout** — pass. `scope.go` beside the model it extends,
   `scope_test.go`/`hookio_test.go` beside their packages, `integration_test.go` in
   `internal/dispatch`; no new top-level files, no stray scripts, R2 intact.
5. **Simplicity (DIY-except-security)** — pass. Stdlib only (`crypto/sha256`,
   `os/exec`, `bufio`); no new dependency in `go.mod`. Nothing in this change set
   touches auth, crypto-as-security, sessions or tokens — the sha256 here is a content
   fingerprint, not a security primitive, which is the correct use of DIY hashing. The
   `cmd.WaitDelay = 5s` fix in `cli.go` is the right minimal answer to the described
   thirty-second timeout, and the bash stub CLI is a good, cheap way to test the
   dispatch contract without credentials.
6. **Decoupling** — fail (mild). `trail.Project.SourceGlobs` is an exported field that
   every caller must remember to populate before calling anything that classifies
   paths; that implicit protocol is precisely what produced finding 1, and it is the
   one seam here that two engineers could not work in parallel without colliding.
   Folding the config into `Resolve` makes `Project` self-contained and removes the
   footgun. The gate/dispatch/trail boundaries are otherwise clean.
