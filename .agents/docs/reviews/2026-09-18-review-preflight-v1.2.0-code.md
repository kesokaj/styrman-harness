---
owner: reviewer
date: 2026-09-18
---
# Review — preflight v1.2.0 code   (2026-09-18)

Range reviewed: v1.1.0..bba384d

## Verdict: REQUEST CHANGES

Two defects, both in the release path this gate is about to walk, both fixable in a
few lines. (1) `.agents/docs/deploy/remote.md` § Releases now carries prose saying
"upload the assets one call at a time, not all of them on `gh release create`" while
the copy-paste block 20 lines above it still passes all five assets to
`gh release create` — the operator following the guide's executable step runs exactly
the call the same guide, and `skills/preflight/SKILL.md` item 5, were changed in this
range to forbid, and a single HTTP 500 rolls v1.2.0's release back the way it did
twice for v1.1.0. (2) `RemoveKeptClone`'s new message suggests a settling command
narrower than the count it explains: the count is over ALL local refs, the command
checks HEAD only, so a commit held on any other local ref is counted but not listed
— and `deploy/local.md`, also changed in this range, tells the operator an empty list
means "the clone holds nothing unique and you can remove it by hand". Measured, both.
The source guard's decision paths are untouched and correct; the test additions are
sound and cannot pass for the wrong reason; the six context rules and the README
`go install` block hold up.

## Findings

- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:126` (the block's asset list) vs
  `:131` (prose) — the release procedure contradicts itself. Step 4 of the code block is
  still `gh release create vX.Y.Z … --notes-file … sh-agent-linux-amd64 sh-agent-darwin-arm64
  sh-agent-linux-arm64 SHA256SUMS SHA256SUMS.sig`; the paragraph added in this range says
  create with `--notes-file` alone and then `gh release upload … --clobber` per file,
  and `skills/preflight/SKILL.md:87` (item 5) now says the same. `/preflight` item 5
  instructs "build them from the tagged tree exactly as that guide says", and a
  copy-paste block is what gets executed. Failure: one 500 on one asset and the whole
  v1.2.0 release rolls back, tag pushed, no release — the recorded failure recurs
  under the very fix meant to prevent it. Fix: replace the block's step 4 with
  `gh release create vX.Y.Z --repo … --verify-tag --title … --notes-file ../release-notes-vX.Y.Z.md`
  followed by a `for asset in …; do gh release upload vX.Y.Z --repo … "$asset" --clobber; done`
  (or five explicit upload lines), so block and prose say one thing. This is criterion 3
  (a changed rule with a restatement left standing) on a path that runs in minutes.

- **[medium / confidence: high]** `internal/harness/report.go:569` — the suggested
  settling command does not cover what the count measured. `ahead` comes from
  `rev-list --count --all --not --remotes` (line 557: every local ref), but the message
  offers `git -C <dir> fetch && git -C <dir> log --oneline HEAD --not FETCH_HEAD`, which
  examines HEAD only. Measured in `.agents/scratch/rev-fetchhead`: after `push -u <url>`
  plus one commit on a second local branch, the count is `2` and the suggested command
  prints **nothing**. `.agents/docs/deploy/local.md:77` (changed in this range) then
  licenses the wrong conclusion in as many words — "an empty list means the clone holds
  nothing unique and you can remove it by hand". Failure: the operator deletes the only
  copy of a commit on a non-HEAD ref, on the guide's own instruction, which is precisely
  the loss this guard exists to prevent (D9, and the security finding it was hardened for).
  Fix: make the command mirror the count —
  `git -C <dir> fetch && git -C <dir> log --oneline --all --not --remotes FETCH_HEAD`
  — and update the matching sentence in `local.md`. Low likelihood (report never creates
  a second branch in the throwaway clone) but the function's own doc comment says the
  clone "holds work this tool did not put there", so a stray local ref is in scope by design.

- **[low / confidence: high]** `internal/harness/report_test.go:936` — the new subtest's
  `t.Skip` is keyed on `@{upstream}` failing to resolve, but the guard never consults
  `@{upstream}`; it counts against `refs/remotes/*`. The condition proven is therefore a
  proxy for the one under test (the new context rule at `context/AGENTS.md:57` asks a test
  to prove *the* condition holds). It runs, not skips, on this machine (`go test -run
  TestRemoveKeptClone -v`: `a_stale_tracking_ref_keeps_the_clone_without_claiming_lost_work`
  PASS, 0.05s), and a git that resolved `@{upstream}` would silently drop the subtest while
  the two updated subtests still pin the new strings — so it cannot hide a failure, only
  coverage. Fix, if taken: assert `git merge-base --is-ancestor HEAD refs/remotes/origin/main`
  fails, i.e. the tracking ref really is stale, instead of the `@{upstream}` probe.

- **[low / confidence: medium]** `skills/preflight/SKILL.md:91` — the new paragraph ("A
  finding that does not block is still a decision … fixed before the tag") is indented as a
  continuation of Step 6's numbered list and sits *after* item 7. A reader following Step 6
  in order has already tagged and pushed at items 3-4 before reaching an instruction whose
  deadline is "before the tag". Fix: move it above item 1 of Step 6, or into Step 5 where the
  verdict is written. Also a mild overlap with Step 4 (`:65`), which routes non-blocking docs
  gaps to the plan's Follow-ups rather than to the release entry — both can be satisfied, so
  it is ambiguity, not contradiction, but "named in the release entry" now reads as also
  covering every docs gap.

## Answers to the dispatch's questions

1. **`report.go` guard and message.** The decision paths are untouched: the diff changes
   only the `fmt.Sprintf` argument inside the existing `if ahead != "0"` branch; the
   `err != nil` → `leave(...)` path (line 559-561), the symlink, checkout, `verifyFeedback`,
   `status --ignored`, and `stash list` gates are byte-identical, and every failure still
   keeps the clone. The claim is **true**, measured: `git push -u <url> HEAD:main` sets
   `branch.main.remote` to the URL and `branch.main.merge` to `refs/heads/main`,
   `@{upstream}` then fails ("not stored as a remote-tracking branch"), no
   `refs/remotes/*` ref moves, and `rev-list --count --all --not --remotes` reports `1`
   for a commit that is on the remote. The suggested command is **right for HEAD** — plain
   `git fetch` does fetch from the URL in `branch.main.remote` and writes `FETCH_HEAD`
   (exit 0), and `log --oneline HEAD --not FETCH_HEAD` is empty — but too narrow for the
   count (finding 2).
2. **`report_test.go`.** The two updated assertions trace to the source change and are
   the right shape (they name the new basis, not the old word). The new subtest **cannot
   pass for the wrong reason**: `"may well be pushed already"` appears only on the
   `ahead != "0"` branch, so any earlier `leave(...)` fails the third assertion, and the
   `git()` helper (`report_test.go:325`) `t.Fatalf`s on a failed command, so the
   `push -q plain HEAD:main` that establishes "the commit really did reach the remote"
   cannot silently no-op. If a future git updated the tracking ref, `ahead` would be `0`,
   the clone would be removed, and `os.Stat(kept)` would fail loudly. The `t.Skip` guard
   does **not** hide a real failure — it names its condition and a reason, and it is
   inactive here — but it is keyed on the wrong predicate (finding 3). The reported
   "three subtests fail against the old message" is consistent with the diff: two updated
   assertions plus the new subtest's `"may well be pushed already"` check.
3. **`context/AGENTS.md`, six rules.** None contradicts an existing rule; none restates one
   already in the file (I grepped `destructive`, `dry-run`, `mutation`, `embed`, and the
   frozen-tree nouns across `context/`, `skills/`, `agents/`, `rules/`, `commands/` — the
   nearest neighbours are the `rm -r`/`git ls-files` hard rule and the "silence from your own
   instrumentation" bullet the mutation text extends, and the extension is coherent with it).
   The guard rule (`err != nil` keeps the data) matches what `RemoveKeptClone` already does,
   so it graduates a rule the code honours. One ambiguity worth a word, not a finding: the
   frozen-tree rule in step 1 says "the main agent edits no file in its scope" — "its scope"
   is left to judgement, and a session writing `changes/` or a handoff while the verifier runs
   will have to decide whether that counts (it should not; the rule would be sharper saying
   "no source, test or config file"). None of the six is broad enough to block ordinary work.
   Note for the record, not a finding: the rules are not yet in the operator's installed
   `~/.claude/AGENTS.md`, which is expected — places pick them up with `sh-agent update`.
4. **`skills/preflight/SKILL.md`.** The upload-one-at-a-time addition is clear, carries its
   evidence, and is consistent with the rest of item 5 (verify-then-publish, never publish
   an unverified artifact) — but its restatement in the deploy guide's code block was left
   standing (finding 1). The non-blocking-finding paragraph is clear in substance; its
   placement is the problem (finding 4).
5. **`README.md` § Install without cloning.** It works as written and changes **no**
   persistent config. `GIT_CONFIG_COUNT`/`KEY_0`/`VALUE_0` is process-scoped (git ≥ 2.31);
   `url.git@github.com:.insteadOf` parses as section `url`, subsection `git@github.com:`,
   key `insteadOf` (git splits on the first and last dot, and the subsection holds no dot);
   the `go` command inherits and forwards the environment to git, and `GOPRIVATE` covers
   both `GONOPROXY` and the checksum database, so the direct SSH fetch is what happens.
   The module path matches `go.mod` (`github.com/kesokaj/styrman-harness`). The
   `git config --global` alternative is inside a comment and clearly labelled as the
   opt-in. One observation, not a finding: the example stays pinned at `@v1.1.0` after
   v1.2.0 ships; the "for v1.1.0 or newer" wording carries it, and the repo tracks no
   version file by rule, so nothing must change — but a reader copying it gets v1.1.0.

## Criteria

1. **Dead code — pass.** Nothing orphaned: the old message string has no other referent, and
   both grep'd test assertions were updated rather than duplicated. No commented-out code,
   no leftover debug output in the +8 source lines.
2. **Flow tracing — fail (finding 2).** `RemoveKeptClone` has one caller path (`update`,
   non-`--check`); the changed message is output-only, so no caller, type or error path
   shifts, and the `error` return is unchanged (`leave` still returns nil, by design). The
   hop that breaks is documentation-to-operator: the count's basis and the remedy's basis
   differ, and `local.md` turns that gap into an instruction. Every new assertion can fail
   (verified: the dispatch reports three subtests failing against the old message, and the
   new subtest's third assertion pins the branch).
3. **Standards conformance — fail (findings 1, 2).** Go side is clean: `gofmt -l cmd internal`
   silent, `go vet ./...` and `go test -race ./internal/harness` pass, error wrapping and the
   stdlib-only rule intact. The failure is the repo's own describing-text rule and the
   restatement rule — a changed procedure whose executable restatement in `remote.md` was
   left standing, and a changed message whose restatement in `local.md` over-promises.
4. **Directory layout — pass.** Test beside its source (`report_test.go` next to `report.go`),
   docs under `.agents/docs/`, context and skill changes in `context/` and `skills/<name>/`,
   scratch fixture under `.agents/scratch/`. Nothing generated committed.
5. **Simplicity (DIY-except-security) — pass.** No new dependency; `go.mod` untouched and
   still `toolchain`-free per ADR-029. No crypto or auth is implemented here — the release
   path consumes `ssh-keygen -Y` and `sha256sum` as designed, which is not hand-rolling. The
   message is long for a single line but it buys an operator decision, not abstraction.
6. **Decoupling — pass.** The change is confined to one function's output string and one test
   file; no new seam, no cross-package reach-in, no file that two engineers would now collide in.

## Not examined

- The non-source docs in the range read only as restatements of the changed rules:
  `BUG-073/074/075`, `INDEX.md`, the two `changes/` entries, `retros/2026-09-17-retro-4.md`,
  `handoffs/2026-09-17-v1.1.0-released.md`, `plans/2026-09-17-release-binary.md` and
  `testing/2026-09-17-release-binary-e2e.md` were not read line by line; the untracked
  `.agents/docs/testing/2026-09-18-verify-preflight-v1.2.0.md` is outside the range and
  was not read at all.
- The full suite was run only for `./internal/harness` plus `go vet ./...`; `TestHarnessMap`,
  the installer against a scratch home, and `govulncheck` were not run here (the verifier's rows).
- No network command was run, no `sh-agent report`, no `sh-agent update`, no git command that
  changes the working tree. The `gh release create`/`upload` behaviour in finding 1 is taken
  from the recorded v1.1.0 measurement in `remote.md`, not re-measured.
- `ROADMAP.md`'s one-line change was read and is consistent with the released state; not a
  criterion beyond that.

## Outcome

Read against the diff by the main agent, 2026-09-18. Both mediums fixed before the tag, as the skill's own SHIP paragraph requires; the verification row taken at bba384d is therefore superseded by a re-run at the commit below.

1. **[medium] `remote.md`'s copy-paste block contradicted the prose above it.** Correct, and the sharper half of the finding is that the BLOCK is what gets executed: `/preflight` item 5 says to publish "exactly as that guide says", so v1.2.0 would have used the all-in-one `gh release create` and could have rolled back on a single 500 — under the very fix meant to prevent it. Fixed: the block now creates the release with `--notes-file` alone, loops one `gh release upload … --clobber` per asset with a retry note, and ends with the `--json assets` listing. This was my error in the retro that added the prose: I wrote the rule beside the commands instead of into them.
2. **[medium] the settling command was narrower than the count it explained.** Correct, and measured here before and after in `.agents/scratch/fetchhead-check/`: a clone with one commit on a second local branch gives count `1`; `git log --oneline HEAD --not FETCH_HEAD` prints NOTHING, while `git log --oneline --all --not --remotes FETCH_HEAD` prints `ad6e33c only on side`. With `local.md` saying an empty list means the clone holds nothing unique, the old pair was an instruction to hand-delete the only copy of that commit. Fixed in the message and in `local.md`'s sentence, and the test now asserts the command's refs match the count's.
3. **[low] the new subtest keyed its skip on `@{upstream}`** rather than on what the guard reads. Fixed: it now skips only when `rev-list --count --all --not --remotes` is `0` or errors — the exact quantity the guard uses — so the shape is confirmed, not inferred from a symptom.
4. **[low] the non-blocking-finding paragraph sat after the items that tag and push**, while its deadline is "before the tag". Fixed: moved above item 1, reworded to say "fixed NOW, before the tag".

Security's one low finding for the same range (`.agents/docs/reviews/2026-09-18-security-preflight-v1.2.0.md`) is fixed by the same restructure: the message no longer interpolates the clone's path at all — `leave()` already names the directory — so nothing from a path rides into text the operator pastes, and the test asserts `git -C` never appears.

Must-fail, both well-formed (the first attempt left a `%s` without an argument and `go test`'s vet step rejected the build, so nothing ran — recorded here because a mutation that does not compile proves nothing): narrowing the command's refs fails two subtests; pasting a path into it fails two. Suite after the fixes: `gofmt`, `go vet`, `go test -race -count=1 ./...` clean.
