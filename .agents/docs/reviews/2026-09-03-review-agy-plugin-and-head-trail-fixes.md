---
owner: reviewer
date: 2026-09-03
---

# Review — agy plugin ownership and HEAD-based trail measurement (2026-09-03)

Scope: the uncommitted working-tree diff on top of `e677a97` — `internal/install/install.go`,
`internal/trail/git.go`, `internal/trail/owed.go`, `internal/trail/trail_test.go`.
Earlier commits were not reviewed. All claims below were reproduced on this machine
(git 2.47.3, `agy` present at `/usr/local/bin/agy`).

## Verdict: REQUEST CHANGES

The `AgyOwned` skip is guarded on the wrong predicate. `Run` skips the mirror whenever
`agy` is on PATH, but `registerAgyPlugin` skips the `agy` call whenever `--home` is a
scratch directory. On a machine with `agy` installed — the operator's machine — those two
guards disagree and `harness install --home /tmp/scratch`, the procedure `AGENTS.md`
mandates for every development install, deploys the Antigravity plugin by *neither*
mechanism and then fails its own doctor. Reproduced verbatim below. The two trail fixes are
correct as far as they go (both new tests are genuine regression tests: each fails when its
fix is reverted, and passes when the other fix alone is reverted), but the `--root` fix
leaves merge commits still invisible while the comment above it claims the opposite, and
`Owed`'s new "same source" comment overstates what the code does — `HeadLinesToday` filters
by `IsSourcePath` and `ChangedLines` does not.

## Findings

- **[high / confidence: high]** `internal/install/install.go:89-93` — the plugin is deployed
  by neither mechanism on a scratch-home install. `haveAgy := agyAvailable()` skips the
  `AgyOwned` mirror, but `registerAgyPlugin` returns early at line 267 (`!env.isRealHome()`)
  without running `agy`. Reproduced:

  ```
  $ go run ./cmd/harness install --home /tmp/.../home2
  skipping agy plugin install: --home is a scratch directory and agy would write to the real home
  deployed; running doctor
  DRIFT ~/.gemini/config/plugins/harness   8 file(s) differ or missing, e.g. agents/changelog.md, ...
  harness: 1 of 10 checks failed — run `harness install`
  $ ls /tmp/.../home2/.gemini/config/plugins/
  ls: cannot access ...: No such file or directory
  ```

  Re-running never fixes it. This breaks the workflow `AGENTS.md` requires ("The installer
  must never be run against the operator's real `$HOME` during development. Use
  `--home /tmp/<scratch>` for every test") and makes `install` exit non-zero on every dev
  run. Fix: derive one predicate and use it for both decisions, so the skip and the
  registration can never disagree —

  ```go
  agyInstalls := agyAvailable() && env.isRealHome()   // computed once, before the target loop
  ...
  if op.AgyOwned && agyInstalls { continue }
  ```

  and have `registerAgyPlugin` consult the same value. Related nit: `haveAgy` is recomputed
  inside the target loop, so `exec.LookPath` runs once per target.

- **[medium / confidence: medium]** `internal/install/install.go:91` + `doctor.go:118`
  (`dirDiff`) — skipping the mirror also skips `pruneExtras`, and nothing replaces it.
  `mirrorDir` was the only thing deleting destination entries the repo no longer has, and
  `dirDiff` walks **src** only, so a file in the deployed plugin that is absent from the
  repo is neither removed nor reported. Delete an agent from
  `antigravity/plugins/harness/agents/` and the stale copy lives in
  `~/.gemini/config/plugins/harness/agents/` indefinitely, invisible to doctor. That is a
  direct hit on R3 ("Deployment overwrites… drift is loud"), which exists because this
  harness once ran a configuration nobody had shipped. Confidence is medium only because
  `agy plugin install` may prune on its own; that cannot be verified without writing to the
  real home. Fix: either have doctor's `AgyOwned` check also walk dst and report extras, or
  keep the mirror and hand `agy` the repo path afterwards (mirroring the repo into the
  destination is harmless — the truncation bug was `agy` copying dst onto dst, not the
  mirror).

- **[medium / confidence: high]** `internal/trail/git.go:61` — the comment "diff-tree names
  files for merge commits too, which `git show` suppresses" is false, and `--root` does not
  change that. Verified on git 2.47.3: for a merge HEAD, `git diff-tree -r [--root]
  --no-commit-id --name-only HEAD` prints **nothing** at all, with or without `--root`, and
  `--numstat` likewise. So merge-commit behaviour is indeed unchanged by this diff (as
  asked), but it is unchanged at *broken*: merging today's feature branch leaves a HEAD that
  owes no rows, exactly the hole `--root` was added to close for root commits.
  `--diff-merges=first-parent` produces the expected file list (verified) and would close
  it. At minimum the comment must stop asserting behaviour git does not have.

- **[medium / confidence: high]** `internal/trail/owed.go:94-102` vs `git.go:88-94` and
  `git.go:120-126` — the two line sources do not measure the same thing, so the new comment
  ("The line count has to come from the same place") is not what the code delivers.
  `HeadLinesToday` filters numstat rows through `IsSourcePath(fields[2])`; `ChangedLines`
  filters nothing (only its `untrackedLines` helper does). A working tree holding a 3-line
  `.go` edit plus a 900-line non-source file scores 903 and fires the reviewer row; commit
  it and the identical change scores 3 and skips the row. Fix: add the same
  `!p.IsSourcePath(fields[2])` guard to `ChangedLines`' numstat loop — it makes both paths
  mean "source lines" and matches what `untrackedLines` already does.

- **[medium / confidence: high]** `internal/trail/owed.go:97-102` — a dirty tree masks
  today's commit entirely. Any single dirty source file makes `Owed` take the working-tree
  branch, so one line edited after a 500-line commit made today reports `Lines = 1` and
  drops the reviewer row that the commit alone owed. Answering the two questions in the
  dispatch directly: **no**, there is no double-counting today, because the branches are
  mutually exclusive; and there would be none even if they were summed, since
  `git diff [--cached] --numstat` measures HEAD→worktree while `HeadLinesToday` measures
  parent→HEAD — disjoint spans of the same day's work. Summing (and unioning `Changed`) is
  the correct shape; the either/or is a hole.

- **[low / confidence: high]** `internal/trail/git.go:83` — the `--root` in the **numstat**
  command is untested. Removing it while leaving everything else intact leaves the entire
  suite green (`go test ./internal/...` passes), so a project's first commit could silently
  score zero lines again. `TestOwedSeesInitialCommit` already commits 80 lines; asserting
  `rep.Lines >= 80` there covers it in one line.

- **[low / confidence: high]** `internal/trail/git.go:98-101` — `headIsToday` re-implements
  `HeadDate` from `head.go` (identical `git log -1 --format=%cd --date=short`). Make it
  `return p.HeadDate() == Today()` so there is one definition of HEAD's date.

- **[low / confidence: high]** `internal/install/install.go:263-265` — in `--dry-run` with
  `agy` present, the plugin mirror op is skipped (never printed) *and* `registerAgyPlugin`
  prints `agy plugin install` before checking `isRealHome`, so a dry run against a scratch
  home advertises a step the real run will not take, and shows no step for the one it will.
  Move the `isRealHome` check above the `dryRun` branch.

- **[informational / confidence: high]** The rename concern raised in the dispatch does not
  bite. `git diff-tree` is plumbing and does not honour `diff.renames` (verified: with
  `diff.renames=true` set in repo config, a rename still prints two plain 3-field rows).
  The `pkg/{b.go => c.go}` layout appears only with an explicit `-M`, which this code does
  not pass. `fields[2]` is therefore a real path today — but leave a comment saying so,
  because adding `-M` later would silently corrupt the filter.

## Criteria

1. **Dead code** — pass. No orphans introduced; `agyAvailable` has two callers,
   `HeadLinesToday` and `headIsToday` are both reached from `Owed`, no commented-out blocks,
   the replaced `path` variable in `registerAgyPlugin` is gone rather than kept.
2. **Flow tracing** — fail. Two hops are inconsistent: `Run`'s `AgyOwned` skip vs
   `registerAgyPlugin`'s `isRealHome` guard (finding 1), and `Owed`'s two line sources
   applying different source filters (finding 4). Error propagation itself is fine —
   `agy plugin install` failure is wrapped with `%w` and fails the install, and the git
   helpers' "error means zero/nil" convention is consistent with the surrounding fail-open
   design required by R4.
3. **Standards conformance** — pass with nits. gofmt-clean, `go vet ./...` clean,
   `go test ./...` green. Comments explain *why* (per the constitution), and the new ones do
   so well — except where they assert behaviour git does not have (findings 3 and 4).
4. **Directory layout** — pass. Changes stay in `internal/install` and `internal/trail`;
   tests sit beside the code they cover in `internal/trail/trail_test.go`; no new files.
5. **Simplicity (DIY-except-security)** — pass. No new dependencies, stdlib and `git`
   plumbing only. `AgyOwned` is a single bool on an existing struct rather than a new
   abstraction. Nothing here touches auth, crypto, hashing, sessions or JWT/OIDC; the
   `sha256` use in doctor is content-identity, not security, and is unchanged.
6. **Decoupling** — pass with a nit. `Op` stays the single vendor-neutral description of a
   deployment step and doctor keeps consuming it unchanged, which is what preserves drift
   detection over the agy-owned path. The nit is duplication rather than coupling:
   `headIsToday` vs `HeadDate` (finding 7).

## Verification of the new tests

Confirmed by reverting each fix independently in a scratch copy:

| Reverted | `TestOwedSeesInitialCommit` | `TestOwedMeasuresCommittedLines` |
|---|---|---|
| `--root` removed from `--name-only` | FAIL ("the initial commit reported no changed source files") | PASS |
| `Owed` restored to `rep.Lines = p.ChangedLines()` only | PASS | FAIL ("Lines = 0 for a 200-line commit against a clean tree") |
| `--root` removed from `--numstat` | PASS | PASS — **uncovered**, see finding 6 |

Both tests are genuine regression tests for their own fix and independent of each other.
