---
owner: reviewer
date: 2026-09-28
---
# Review — update-from-zip-code-r3   (2026-09-28)

Range reviewed: 1cb2608..worktree (`git diff 1cb2608` plus the untracked `internal/harness/zip.go`, `internal/harness/zip_test.go`, `cmd/sh-agent/zip_test.go`)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing in this round is broken. All four r2 findings are fixed, and each fix has a test that fails without it. I ran two mutations myself in a scratch copy, and both compiled and failed a named test:
- The created-chain cleanup, rewritten to the shadowed `created, err :=` form, fails `TestUpdateFromZipRefusedContentWritesNothing` ("a refused install left 1 entries under a fresh home").
- The checkout dry run with `*check = true` removed fails `TestCheckoutUpdateDryRunChangesNothing` on both assertions: HEAD moved, and an install ran.

The remaining items are describing-text nits:
- The comment r2 named at `main.go:1052-1053` is still there, and it is now false.
- The usage text and the package doc comment do not mention `update --from-zip` or `update --dry-run`.
- D13 says "every" zip CLI test gets a scratch cwd, which is one test short.

`gofmt -l cmd internal` printed nothing, `go vet ./...` passed and `go test -race -count=1 ./...` passed (run 2026-09-28, this worktree).

## Round-2 findings — status

| # | r2 finding | Fixed? | Test that fails without the fix |
|---|---|---|---|
| 1 | BUG-093 at the CLI: `installFiles`/`Plan` refuses the zip after `~/.config/sh-agent` was created | **Fixed.** `updateFromZip` creates the chain with `MkdirAllTracked` before `UnpackZip` and removes it in a `defer` keyed on the named result `err` (`main.go:1057-1075`). `UnpackZip` now creates only its own `zip-*` (`zip.go:79-90`). | `TestUpdateFromZipRefusedContentWritesNothing` (`cmd/sh-agent/zip_test.go:395`), a description-less skill on a fresh home. I checked it by mutation (shadowed `err`): it compiled, and the test failed at `:411`. The stale comment this finding named is **not** fixed (finding 1 below). |
| 2 | Red-path isolation (BUG-094 class) | **Fixed.** `isolate()` (`zip_test.go:72-77`) empties `SH_AGENT_REPO` and `PATH` and moves the cwd to a scratch dir. It is called by every `--from-zip` test except `TestFromZipRefusedOnOtherCommands`, which already had `PATH` emptied. `plainCopy` resolves the checkout once at package load (`main_test.go:257-259`), so a moved cwd still finds it. The `--repo ../..` row now resolves under the scratch cwd, not the real repository. | This is an isolation property, not a behaviour, so no single test pins it. I checked it by reading: a regression to `!= ""` now reaches `ResolveSource` from a scratch cwd (dev build: error) or `UpdateFromRelease` with no `gh` on `PATH` (error). Neither can fetch or push. |
| 3 | Checkout-path `update --dry-run` untested | **Fixed.** | `TestCheckoutUpdateDryRunChangesNothing` (`zip_test.go:468`) runs against a local bare origin. I checked it by mutation (deleted `*check = true`): it compiled, and the test failed at `:508` (HEAD moved) and `:511` (install ran). |
| 4 | Chain cleanup had no test | **Fixed by moving the cleanup.** `UnpackZip` no longer owns the chain. Its own `zip-*` is removed by its `defer` on a late refusal. The chain is covered by row 1's test. `TestUnpackZipLateRefusalCreatesNoParent` now pins that `UnpackZip` never creates the parent: it would fail if a `MkdirAll(parent)` came back before a late refusal. | See row 1. The dispatch reports the main agent's mutations for the TMPDIR cleanup (`TestUpdateFromZipDryRunCleansTempDir`) and the `zip-` prefix (`TestUpdateFromZipSweepsStaleDirs`, the `feedback-old` case). I did not repeat those two. |

## Findings

- **[low / confidence: high]** `cmd/sh-agent/main.go:1052-1053` — stale comment (describing-text rule).
  - It says "UnpackZip creates the parent only after the zip's names pass, so a refused zip writes nothing." Since D13, `UnpackZip` never creates the parent. `updateFromZip` creates it a few lines below, and a refused zip does write: the chain, removed again by the `defer`, and the stale-`zip-*` sweep, which runs before `UnpackZip`.
  - r2 finding 1 named this comment. It survived the fix.
  - Fix: replace the first sentence with what the code does, or delete it. The dry-run sentence that follows is still true.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:42` and the package doc comment `main.go:8` — the command's own help does not describe the new behaviour.
  - The usage line `update [--check]` does not mention that `--dry-run` is now accepted on the checkout and release paths and prints `[dry-run] …` with exit 0.
  - The package doc lists `update [--check]` and has no `--from-zip` line. The usage block at `:43` does have one.
  - README `:168`/`:177` does document both, so this is a nit, not a missing doc.
  - Fix: add `[--dry-run]` to both `update` lines, and a `--from-zip` line to the package doc.

- **[low / confidence: high]** `.agents/docs/plans/2026-09-28-update-from-zip.md` D13 — the record overstates the fix.
  - D13 says "Every zip CLI test runs with an empty `PATH` and a scratch working directory". `TestFromZipRefusedOnOtherCommands` (`zip_test.go:257`) empties `PATH` but keeps the package directory as cwd.
  - No concrete failure follows. Under a regression, `install`/`uninstall` write only into the temp home, and `report` has no `git`/`gh`.
  - Fix: call `isolate(t)` there too, or narrow D13's sentence.

## Items the dispatch asked about

- **`updateFromZip` created-chain cleanup (the named `err`, the `defer`).** Correct.
  - `vendor, err := resolveVendor(...)` and `dir, sum, err := harness.UnpackZip(...)` sit in the function's outermost scope, so `:=` reuses the named result. It does not shadow it.
  - The chain assignment uses `=` inside the `else` block, and the closure reads the named `err`.
  - `return installFiles(...)` assigns the result before the defers run.
  - The defers run in LIFO order: `os.RemoveAll(dir)` removes `zip-*` first, then `RemoveCreated` can empty the chain.
  - If `installFiles` fails after `Install` wrote the manifest (a `pathHint`/`attributionHint` error), `os.Remove` keeps the non-empty `~/.config/sh-agent`. No written data is removed.
  - The dry-run branch registers no chain defer and creates nothing. Its `err` shadow at `if _, err := os.Stat(parent)` is harmless.
- **`MkdirAllTracked` / `RemoveCreated`** (`zip.go:180-215`). Correct.
  - It walks up to the first existing entry with `Lstat`. A symlinked ancestor counts as existing and is not recreated.
  - It creates outermost first, with 0700, which matches README `:176` ("tightened to 0700 by every install").
  - `ErrExist` from a concurrent creator is skipped rather than recorded, so it is never removed.
  - A mid-chain `Mkdir` error removes what this call made.
  - `RemoveCreated` uses `os.Remove`, deepest first, so it never deletes a non-empty directory.
  - Untested branches: the `ErrExist` skip and the mid-chain error cleanup. Neither is reachable from `updateFromZip` except under a race or ENOSPC, since every directory it creates is 0700 and owned by the caller. It is not a finding.
- **`sweepStaleZipDirs`** (`main.go:1096-1114`). Correct.
  - It reads the directory with `os.ReadDir`, never a glob built from the home (`TestUpdateFromZipSweepWithGlobCharsInHome`).
  - It matches the `zip-` prefix and skips symlinks and non-directories.
  - The age comes from `Lstat` mtime, and `RemoveAll` on a path that became a symlink meanwhile removes the link, not its target.
  - A run in progress could only be swept if it stayed more than an hour between unpack and delete.
  - It runs only on a real run, after the chain exists and before `UnpackZip`. A refused zip may still sweep hour-old leftovers, which is harmless and only touches this command's own debris.
- **Network or operator credentials from tests (BUG-094).**
  - Every `--from-zip` CLI test has `PATH` emptied. All but one also have a scratch cwd (finding 3).
  - `TestReleaseUpdateDryRunChangesNothing` puts only a fake `gh` on `PATH` and moves the cwd. Under the mutation its fake returns garbage, so the signed-checksum verification fails and nothing is replaced.
  - `TestCheckoutUpdateDryRunChangesNothing` keeps the real `git`, which it needs, but passes `--repo` to a fixture whose `origin` is a local bare repository. A red run fetches and pulls from that file path only. I ran that red run, and it touched only the fixture.
  - The internal tests (`internal/harness/zip_test.go`) do not shell out. `TestResolveSourceSkipsZipManifest` moves its cwd to a temp dir.
  - I found no test that can reach the network or the operator's credentials.

## Criteria

1. **Dead code: pass.**
   - `MkdirAllTracked`/`RemoveCreated` are called by `updateFromZip`.
   - `short()` has two callers.
   - No commented-out blocks or debug output.
2. **Flow tracing: pass.**
   - The chain is created, then the stale dirs are swept, then `UnpackZip` writes only `zip-*`, then `installFiles`/`Plan`/`Install` with `Origin` runs, then both defers.
   - Every error path returns through the named result.
   - The dry-run flag reaches both `update` paths, and both are tested.
   - Resources are closed on every path in `zip.go`.
3. **Standards conformance: pass, with nits.**
   - gofmt, vet and race tests are clean.
   - The describing-text rule is not met in three places: findings 1–3.
   - Restatements of the limits, the sweep age and the `zip-*` location in the living docs: I grepped README, AGENTS, REQUIREMENTS, specs, deploy, DECISIONS, context and skills. None is stale beyond the finding-1 comment.
4. **Directory layout: pass.** `zip.go` is in `internal/harness`, with tests beside `zip.go` and `main.go`.
5. **Simplicity (DIY-except-security): pass.**
   - Stdlib only.
   - `crypto/sha256` is used as a checksum, as designed, which is not hand-rolling.
   - `MkdirAllTracked` is the minimum needed for exact cleanup.
6. **Decoupling: pass.**
   - The chain ownership sits in one caller.
   - `UnpackZip`'s contract ("parent exists; I remove only my own dir") is narrower than in r2.
   - The `zip:sha256:` parse is still repeated in `status`, `knownCheckout` and `ResolveSource`, with no failure mode.

## Not examined

- `internal/harness/zip_test.go` lines 1–243 and 357–410: read by test name and grep only (no `exec`, `Setenv` or `Chdir` outside `TestResolveSourceSkipsZipManifest`), not line by line.
- `harness.Install`, `attributionSet`, `pathHint` and `UpdateFromRelease` internals: judged from their call sites and the passing suite.
- I did not repeat the main agent's mutations for the TMPDIR cleanup and the sweep's `zip-` prefix.
- The docs hunks (README, delivery.md, deploy guides, REQUIREMENTS, DECISIONS/ADR-038, skill-matrix, VENDORED.md, ROADMAP, AGENTS.md, .gitignore): these belong to the docs review, and I only grepped them for restatements.
- The verify logs and the cold-read and security reports.
- macOS.

Scratch used: `.agents/scratch/rv3/` (a copy of the tree without `.git`, and `main.go.orig`). The mutations were made only there, and the copy was restored after each run.

## Outcome

See `changes/2026-09-28-1752-update-from-zip.md` and plan D11/D13. After this round:
- **Fixed:** the stale sweep now runs only after a successful install; `--from-zip` refuses a value that starts with `-`; a failed `RemoveAll` of the unpacked copy prints a warning; the stale comment, usage and package doc are updated; `TestFromZipRefusedOnOtherCommands` uses `isolate`; ADR-038's Rejected entry now says deferred; ADR-034 is marked superseded in part.
- **Left:** the unsigned-zip medium, open by decision and put to the operator; the pre-existing `report.go:411` glob, outside this diff and a follow-up; the matrix curl grant, accepted under option 1 (the allowlist is stated not to be a boundary).
