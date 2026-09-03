---
owner: reviewer
date: 2026-09-03
---

# Review — glob matcher and doc sync   (2026-09-03)

Scope: uncommitted diff at 6e6aef6 (diff sha1 5dbbb22a939e) — 13 modified files
(the 12 named in the dispatch plus `TODO.md`, which was already modified in the
tree; it is prose-only and not reviewed further). Commit 6e6aef6 itself is out of
scope, having been reviewed in `2026-09-03-review-harness-go-binary.md`.

## Verdict: REQUEST CHANGES

The index arithmetic in `compileGlob` is correct — I verified every arm of the
switch against a standalone harness, and `**/`, `/**`, bare `**`, `*` and `?` all
advance exactly the right number of bytes. The problems are elsewhere. The worst
one is a documentation claim the binary cannot honour: `claude/agents/verifier.md`
now instructs the verifier to "read the dispatch prompt for the exact command(s)
that failed and their output — that is your starting point, not a blind rerun",
but `rowPrompt` in `internal/dispatch/dispatch.go:221` has no `RowVerifier` case,
so the escalation prompt is the generic default string, and the command output
`runVerifier` captured in its `log` builder is discarded on the failure path
without ever being written or passed on. The agent is told not to blind-rerun and
is then handed nothing else to work from. Behind that: `compileGlob` silently
mangles any non-ASCII glob and silently demotes `path.Match`'s character classes
to literals, and both failures are invisible — a `sourceGlob` that matches nothing
disables the Definition-of-Done rows for those paths without a word, which is
exactly the failure mode R2's "twelve files drift independently and silently"
rationale exists to prevent. `gofmt -l`, `go vet ./...` and `go test ./...` all
pass.

## Findings

- **[high / confidence: high]** `claude/agents/verifier.md:11` (and the mirrored
  `antigravity/plugins/harness/agents/verifier.md`, plus the "Pre-flight" step 2
  and the `## Commands dispatch ran` template section) — the agent is told its
  starting point is "the dispatch prompt for the exact command(s) that failed and
  their output". No such content exists. `internal/dispatch/dispatch.go:226-235`
  switches on `RowChangelog`, `RowReviewer`, `RowSecurity` and falls through to
  `"Complete your Definition-of-Done row for the current changes."` for the
  verifier; `internal/dispatch/verify.go` builds the per-command log in a
  `strings.Builder` and, on `!passed`, calls `runAgentRow` without it and lets the
  builder go out of scope. Fix: add a `case trail.RowVerifier` to `rowPrompt` (or
  a prompt-extension parameter) that embeds the command list and the tail of the
  failing output, then the doc is true. Until then the doc describes behaviour the
  code does not have, and a doc that lies about its input is worse than a silent
  agent because it forbids the fallback the agent would otherwise take.

- **[medium / confidence: high]** `internal/trail/source.go:149` — `string(glob[i])`
  converts a *byte* to a rune, so every non-ASCII byte in a glob is reinterpreted
  as a Latin-1 code point. Verified: `src/år/*.go` compiles to
  `^src/Ã¥r/[^/]*\.go$` and does not match `src/år/x.go`; `dôc/**` never matches
  `dôc/a.md`. The old `path.Match` handled these correctly. The consequence is the
  one the whole matcher exists to avoid: a declared source path is silently not
  source, so the DoD rows silently do not fire. Fix: iterate runes
  (`for i, r := range glob` with an explicit skip counter, or build a
  `[]rune` first) and `regexp.QuoteMeta(string(r))`.

- **[medium / confidence: high]** `shared/skills/preflight/SKILL.md:50` — "Run
  `harness dispatch` ... over the full release diff `<base>..<head>`". `harness
  dispatch` accepts no range: `dispatch.Run` parses only `--dry-run`, `--json`,
  `--vendor`, `--only`, and `Project.Owed` computes its change set from
  `DirtySource()` falling back to `HeadSourceToday()`. After Step 1 has committed
  the dirty work, Step 2 therefore evaluates *today's commits*, not the release
  range — a release spanning several days gets its reviewer and security rows
  decided from a fraction of the diff, and the Step 2 security bullet's "a changed
  path in the release diff" is false for the same reason. Fix: either add a
  `--range <base>..<head>` flag to dispatch and have `Owed` use it, or change the
  skill to say plainly that dispatch covers today's work only and that the release
  range is still the `reviewer`/`security` agents' own responsibility via their
  `Scope:` lines.

- **[medium / confidence: high]** `shared/skills/audit/SKILL.md:28` — "run
  `harness dispatch` for the verification row" alongside three subagents already
  dispatched in parallel. A bare `harness dispatch` runs *every* owed row; since
  the parallel `reviewer` and `security` subagents have not yet written their
  dated artifacts, `HasArtifact` is false for both and dispatch will dispatch them
  a second time, concurrently with the ones the skill just started. The flag that
  makes the sentence true already exists. Fix: `harness dispatch --only verifier`.

- **[medium / confidence: high]** `shared/skills/preflight/SKILL.md:56` — the
  verification row bullet claims `harness dispatch` runs "dual-arch ARM64+AMD64
  image builds when a Dockerfile changed". `verifyCommands`
  (`internal/dispatch/verify.go:88`) has no Dockerfile detection and no `buildx`
  path: it returns the repo-declared `verify` array, else `go vet ./...` +
  `go test -race ./...` for a `go.mod` and `npm test --if-present` for a
  `package.json`. Nothing else. Fix: drop the claim, or say the builds happen only
  when the project lists them in `verify`.

- **[medium / confidence: high]** `shared/skills/init-project/SKILL.md:212-214` —
  three inaccuracies in the field documentation the dispatch asked me to verify
  against the code. (a) Auto-detection is described as "`go test ./... -race` for a
  Go module"; the code also runs `go vet ./...`, and the flag order differs.
  (b) "the project's own `npm test`/`vitest` for a `web/`or `ui/` with
  `package.json`" — `verifyCommands` stats `filepath.Join(root, "package.json")`
  only, so a `web/package.json` is not detected at all, and the command is always
  `npm test --if-present`; `vitest` is never invoked. (c) "Omit the field to fall
  back to auto-detection" is safe advice only for Go/npm repos: in a project with
  neither manifest, `verifyCommands` returns empty, `runVerifier` returns
  `OK == false`, and `runRows` then skips *every* remaining row with "verification
  failed" — omitting `verify` silently disables the whole dispatch. That case
  deserves an explicit warning in the same bullet. (Also: `dispatchTimeoutSeconds`
  is shown as `300` while `DefaultConfig` is `900`, directly under a bullet that
  says "Defaults (shown above) apply"; and the field also caps each individual
  `verify` command, not just "a subagent".)

- **[medium / confidence: medium]** `internal/trail/source.go:128` — character
  classes silently regressed. `path.Match`, which this replaces, supports
  `[a-z]`, `[^abc]`, `[!abc]` and `\` escapes. `compileGlob` `QuoteMeta`s every
  non-`*`/`?` byte, so `[a-z]*.go` compiles to `^\[a-z\][^/]*\.go$` and matches
  only a path literally named `[a-z]…`. Nothing warns; a project that had such a
  glob loses its gates. Either implement the class syntax (translating `[!` to
  `[^`) or reject a glob containing an unescaped `[` at `LoadConfig` time so the
  drift is loud, per R3.

- **[medium / confidence: high]** `internal/trail/source.go:118-121` and
  `internal/trail/trail_test.go:242` — the compile-error branch is unreachable and
  the test that claims to cover it is vacuous. Every byte outside `*`/`?` is
  `QuoteMeta`'d, so the generated pattern is always a valid regexp and
  `regexp.Compile` cannot fail; `if err != nil { return false }` is dead code (and
  the question of whether a compile error is cached or retried is moot for the
  same reason — it is not cached, but it also cannot happen). Worse,
  `TestMatchGlobRejectsBadPattern` passes for the wrong reason: `[unterminated`
  compiles cleanly to `^\[unterminated$`, which simply does not match
  `anything.go`. The test asserts nothing about malformed input and would keep
  passing if the error handling were deleted. Delete the unreachable branch and
  the test, or make `compileGlob` genuinely validate (see the character-class
  finding) and then the test has something to assert.

- **[low / confidence: high]** `internal/trail/trail_test.go:216` — coverage gaps
  in the new test. Two of `compileGlob`'s five switch arms are never exercised:
  bare `**` not adjacent to a separator (`a**b` → `^a.*b$`) and `?`
  (`?.go` → `^[^/]\.go$`). Also untested: the middle-doublestar collapse
  (`a/**/b` correctly matches both `a/x/y/b` and `a/b`), and that
  `claude/agents/**` matches the bare prefix `claude/agents`. The table covers the
  five shapes `init-project` documents and nothing riskier. Add the two missing
  arms at minimum — for a function that decides which gates fire, an untested arm
  is an ungated arm.

- **[low / confidence: medium]** `internal/trail/source.go:137` — `/**` compiles
  to `(?:/.*)?`, an *optional* group, so the separator the pattern demands can
  vanish: `docs/**.md` matches the file `docs.md`, and `x/**` matches a file
  literally named `x`. Harmless for the trailing-`/**` shapes this repo uses (git
  never reports directories), but it is a genuine over-match for any pattern with
  a suffix after the `**`. If a prefix is meant to be required, emit `(?:/.*)?`
  only when `/**` ends the glob and `/(?:.*/)?` otherwise.

- **[low / confidence: high]** `claude/agents/changelog.md:25` — the new guidance
  to write `date: DATE-UNKNOWN` into frontmatter collides with
  `Freshness.checkDates` (`internal/trail/check.go`), which compares dates as
  strings: `"DATE-UNKNOWN" > "2026-09-03"`, so `harness check` will report the
  entry as a **future-dated artifact** and mark the trail STALE. Failing loudly is
  arguably right, but the reason it prints is wrong and will send a reader looking
  for a clock problem. Either special-case `DATE-UNKNOWN` in `checkDates` with its
  own message, or keep the marker out of the `date:` field and in the heading only.

- **[low / confidence: high]** `shared/skills/preflight/SKILL.md:60` — "Its
  failing checks:" enumerates three but `CheckFreshness` has four: it also fails
  on uncommitted source files with no changelog entry dated today. Add it, or drop
  the closed-list phrasing.

- **[low / confidence: medium]** `shared/skills/init-project/SKILL.md:214` — the
  `sourceGlobs` bullet ("which changed paths count as 'source' for triggering
  rows") reads as a whitelist, but `IsSourcePath` consults the globs *additively*:
  a non-matching path can still be source via `sourceExts`/`sourceNames`, and
  `.agents/` and the docs root are excluded *before* the globs are consulted, so
  they cannot be opted back in. The example
  `["**/*.go", "web/**/*.ts", "web/**/*.tsx"]` is a no-op — all three extensions
  are already in `sourceExts` — which teaches the field's purpose backwards. Say
  it is an *addition* for projects whose deliverable the default classification
  misses, and use this repo's own prose globs as the example.

- **[low / confidence: high]** `internal/trail/source.go:124` — the bare type
  assertion `re.(*regexp.Regexp)` panics if anything else ever lands in
  `globCache`. Only one store site exists today so this is safe, but a
  comma-ok assertion (or a typed wrapper struct around the `sync.Map`) costs one
  line and removes the class of bug. Note also that the `Load`/`LoadOrStore`
  sequence compiles the same glob redundantly under concurrent first use; benign
  here, and `LoadOrStore`'s return is used correctly so all callers share one
  instance.

## Criteria

1. **Dead code — fail.** `compileGlob`'s error return can never be non-nil, so
   `matchGlob`'s `if err != nil { return false }` is unreachable, and
   `TestMatchGlobRejectsBadPattern` exercises nothing. Both should go, or the
   validation should become real. Nothing else in the diff is orphaned; the
   `path.Match` call and the `CutSuffix` special case it replaced were both
   removed cleanly, and `strings.CutSuffix`'s removal left no unused import.

2. **Flow tracing — fail.** The one seam this diff creates —
   `runVerifier` failure → `runAgentRow` → `rowPrompt` → the verifier agent — is
   documented as carrying the failing commands and their output and does not carry
   them; the captured log is dropped. `IsSourcePath` → `matchGlob` → `compileGlob`
   traces correctly otherwise, `p.SourceGlobs` is populated identically in both
   `Owed` and `CheckFreshness` so the dispatcher and the freshness report see the
   same file set, and no error is swallowed on that path except the unreachable
   compile error. Secondary: an escalated verifier writes an artifact with today's
   date, which makes `HasArtifact(RowVerifier, today)` true and marks the row
   satisfied for the rest of the day even though it recorded FAIL — pre-existing in
   6e6aef6 and already logged in `TODO.md`, but this diff's hardcoded
   `## Status: FAIL` template makes it certain rather than incidental.

3. **Standards conformance — pass with one exception.** `gofmt -l .` is clean,
   `go vet ./...` and `go test ./...` pass. The code is idiomatic Go, comments
   explain *why* per `AGENTS.md`, and the doc comments on `matchGlob` and
   `compileGlob` are genuinely good — they state the constraint (`path.Match` has
   no `**`) rather than restating the code. The exception is the `string(byte)`
   conversion, which `go vet`'s `stringintconv` does not catch for `byte` and which
   the `golang` skill's rune-vs-byte guidance rules out. R1 is satisfied: I
   re-derived `antigravity/plugins/harness/agents/*.md` with
   `go run ./cmd/harness install --dry-run --home <scratch>` and the working tree
   did not move — the three generated agent files match their `claude/agents/`
   sources. R2 (one binary), R3, R6 (ten-line replies), R7 (frontmatter, not
   filenames) are all respected by the doc edits.

4. **Directory layout — pass.** `matchGlob`/`compileGlob` sit in
   `internal/trail/source.go` beside the classifier that is their only caller;
   tests are in `internal/trail/trail_test.go` next to the code. The doc edits stay
   in `shared/skills/` and `claude/agents/` with generated copies under
   `antigravity/plugins/harness/agents/`, exactly as `AGENTS.md` prescribes.

5. **Simplicity (DIY-except-security) — pass.** Hand-rolling a `**` translator
   rather than taking `bmatcuk/doublestar` is the right call for ~25 lines against
   R2's one-artifact rule and the constitution's stdlib-first default, and no
   dependency was added (`go.mod` untouched). No security-sensitive code is
   hand-rolled here — the diff touches no auth, crypto, hashing, session or
   JWT/OIDC surface. The regexp translation is RE2, so no catastrophic
   backtracking is reachable from a hostile `harness.json`, and `QuoteMeta` means
   no pattern injection. The `sync.Map` cache is the one piece of unneeded
   complexity: the pattern set is at most a handful of globs per process and the
   process is short-lived, so a plain map (or no cache) would do — but it is cheap
   and correctly used, so this is an observation, not a finding.

6. **Decoupling — pass.** `compileGlob` is a pure function of a string, testable
   in isolation, with no reach into `Project`; `matchGlob` adds only memoisation.
   The `globCache` package-level global is the one shared-state seam and is keyed
   solely on the glob, so it stays correct across `Project` values. The doc files
   are independently editable per agent and per skill; the one coupling worth
   naming is editorial rather than structural — the same claim about what
   `harness dispatch` runs is now restated in `preflight`, `audit`,
   `init-project` and `verifier.md`, so the four inaccuracies above are one fact
   stated wrongly in four places. A single "what `harness dispatch` does" section
   that the others link to would keep them from drifting apart again.
