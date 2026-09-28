---
owner: reviewer
date: 2026-09-28
---
# Review — update-from-zip-code-r2   (2026-09-28)

Range reviewed: 1cb2608..worktree (`git diff 1cb2608` plus the untracked `internal/harness/zip.go`, `internal/harness/zip_test.go`, `cmd/sh-agent/zip_test.go`)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The round-1 findings are fixed inside `UnpackZip` and the flag handling, but BUG-093 is still open one layer up. `updateFromZip` hands the unpacked tree to `installFiles`, and `harness.Plan` refuses zip content there: a skill with no description, a secret-shaped file, an unquoted `#`, an unmapped agent tool. By that point `UnpackZip` has already created `~/.config/sh-agent`. I reproduced it: a zip whose `skills/bash/SKILL.md` has no `description:` exits 1 on a fresh home and leaves `home/.config/sh-agent` behind.

Second, `TestFromZipEmptyValueRefused` has no isolation. If the regression it pins comes back (an empty value read as "no flag"), the red run falls through to a normal `update` from the package directory. `ResolveSource` then finds the operator's real checkout, and `update` runs `git fetch` with the operator's credentials. On a clean tree that is behind, it also runs `git pull --ff-only` and `Reinstall`. That is BUG-094's class: a red test with real side effects.

gofmt, vet and `go test -race ./...` are clean (run 2026-09-28).

## Round-1 findings — status

| # | Round-1 finding | Fixed? | Test that fails without the fix |
|---|---|---|---|
| 1 | Refusals after `MkdirAll(parent)` leave `~/.config/sh-agent` | Fixed inside `UnpackZip`: `checkZipTrees` validates kinds, duplicates and sizes before `mkdirAllTracked`, and the defer removes the created chain. Still open at the CLI level (new finding 1). | `TestUnpackZipLateRefusalCreatesNoParent` (symlink, duplicate; fresh parent) fails under the round-1 code. Nothing fails if the chain-cleanup loop is removed (finding 4). |
| 2 | Empty `--from-zip` value falls through to `update` | Fixed: `fs.Visit` detects the flag and the empty value is refused (`main.go:119-132`). | `TestFromZipEmptyValueRefused`. Its red path is unsafe (finding 2). |
| 3 | `uninstall`/`report --from-zip` run | Fixed: the refusal block sits before the first `switch`. | `TestFromZipRefusedOnOtherCommands`, 8 commands. It asserts the home stays empty and nothing is printed, with `PATH` emptied, which is correct isolation. |
| 4 | `diff` compares a zip home to the embedded files | Fixed (`main.go` diff case, `ZipPrefix` check). | `TestDiffAfterZip` (release version; would print drift containing ZIPMARKER). |
| 5 | Interrupted runs leave `zip-*` | Fixed for later real runs: `sweepStaleZipDirs` removes entries older than 1 h. A dry run never sweeps, by design. | `TestUpdateFromZipSweepsStaleDirs`. It uses `Chtimes` for the old entry, keeps the fresh one, and asserts that a dry run does not sweep. |
| 6 | `update --dry-run` on a release binary downloads and replaces | Fixed: `dry` sets `*check`, and `ErrNewerRelease` under dry prints and returns nil. | `TestReleaseUpdateDryRunChangesNothing`, with a fake `gh` alone on `PATH`. The checkout path has no test (finding 3). |

## Findings

- **[medium / confidence: high] BUG-093 recurrence.** `cmd/sh-agent/main.go:1064` (`UnpackZip` creates the parent chain) vs `:1075` (`installFiles` → `harness.Plan` refuses the zip's content).
  - `UnpackZip` validates only the archive's shape. Plan's content refusals come after the parent exists: a SKILL.md without a description, a secret-shaped name (`plan.go:120`), an unquoted `#` in frontmatter, an unmapped agent tool. On success `UnpackZip` returns with `~/.config/sh-agent` created. The caller's `defer os.RemoveAll(dir)` removes only the `zip-*` tree.
  - Reproduced with a scratch binary built from the worktree, input `.agents/scratch/rv2-zip/nodesc.zip`. Output: `sh-agent: skill bash: SKILL.md has no description`, exit 1. `find home1` → `home1/.config/sh-agent`.
  - The comment at `main.go:1051-1052`, "UnpackZip creates the parent only after the zip's names pass, so a refused zip writes nothing", is false for this case. It is also stale on "names": validation now covers kinds, duplicates and sizes too. Under the describing-text rule that is part of this finding.
  - No CLI test refuses a zip at the Plan stage on a fresh home.
  - Fix: have `updateFromZip` own the chain. Call a tracked mkdir before `UnpackZip`, so `UnpackZip` always receives an existing parent. When `installFiles` returns an error, remove the `zip-*` tree and then `os.Remove` the created chain, deepest first. `os.Remove` keeps anything a partial install wrote.
  - Add a fresh-home CLI test with a description-less skill that asserts `os.ReadDir(home)` is empty.

- **[medium / confidence: high] Red-path isolation (BUG-094 class).** `cmd/sh-agent/zip_test.go:272-284` (`TestFromZipEmptyValueRefused`).
  - The test sets `SH_AGENT_REPO=""` but neither `PATH` nor the working directory. `version` is the dev build.
  - Under the exact regression it pins (round-1 code: `if *fromZip != ""`), `update --from-zip= --home <tmp>` falls through to the normal `update`. `ResolveSource` walks the cwd (`cmd/sh-agent`) up to the real checkout (`repo.go:65-72`), and `knownCheckout` passes because the scratch home has no manifest. `update` then calls `harness.Behind` → `git fetch --quiet` in the operator's checkout, with their credentials. If the tree is clean and behind, it also runs `Pull` and `Reinstall`.
  - I did not run the mutation, deliberately, because it would fetch. The call chain was traced by reading.
  - The same exposure exists, less likely, if the whole `zipSet` block regressed:
    - `TestUpdateFromZipRefusals` (`:108-126`): the `--repo ../..` row names the real repository root.
    - `TestUpdateFromZip`, `TestStatusAfterZip` and `TestDiffAfterZip` set `version = "v1.2.3"`. A fall-through there reaches `UpdateFromRelease`, which runs the real `gh`.
  - The sibling `TestFromZipRefusedOnOtherCommands` shows the right pattern and cites BUG-094.
  - Fix: in every `--from-zip` CLI test, `t.Setenv("PATH", t.TempDir())` (or a fake-`gh` dir) and `t.Chdir(t.TempDir())`. The zip path needs neither `PATH` nor the cwd: `selfExecutable` is `os.Executable`. Confirm the tests stay green after the change.

- **[low / confidence: high] Checkout-path `--dry-run` untested.** `cmd/sh-agent/main.go:232-235, 271-276`.
  - D11 makes `update --dry-run` read-only "on the checkout and release paths". Before this range, the checkout path ignored `--dry-run` entirely: it pulled and reinstalled.
  - No test runs `update --dry-run` against a checkout. `grep` in `main_test.go` finds `--check` and plain `update` only. Removing `*check = true` would leave the checkout path pulling on a dry run with every test green.
  - Fix: extend the existing local-remote checkout test near `main_test.go:1129` with `update --dry-run` on a behind checkout. Assert HEAD is unchanged, the `[dry-run] update would pull` line appears, and exit is 0.

- **[low / confidence: high] Chain cleanup has no test.** `internal/harness/zip.go:84-95` (the created-chain cleanup) and `zip_test.go:336` (`TestUnpackZipLateRefusalCreatesNoParent`).
  - After D11, both cases in that test are refused in `checkZipTrees`, before `mkdirAllTracked` runs. The test proves validation precedes creation, but it never reaches the cleanup loop.
  - The post-creation refusals are the final `IsRepo` (a zip with `context/AGENTS.md` and no `skills/`) and I/O errors. The loop works: I reproduced `noskills.zip` on a fresh home, exit 1 and the home left empty. But deleting the loop fails no test.
  - Fix: add a `noskills` case with a fresh parent (`home/.config/sh-agent`). Or move the `skills/` presence check into validation, so the loop guards only I/O errors.

## Items the dispatch asked about

- **`UnpackZip` validate-then-write, and cleanup of the parent chain.**
  - The split is sound. `readCapped` → hash → parse → name pass → root → `checkZipTrees` (kinds, duplicates, sizes counted by decompressing into `io.Discard`) → `mkdirAllTracked` → write.
  - The named return `err` is what the defer reads. Every error after `mkdirAllTracked` returns through the named result, so the shadowed `err` inside the loop does not bypass it.
  - The chain is removed deepest first with `os.Remove`, which never deletes a non-empty directory.
  - The write phase re-checks sizes and uses `O_EXCL`, as defence in depth.
  - Gaps: findings 1 and 4.
- **`fs.Visit` block before the first `switch`.** Correct.
  - Only home resolution (`os.UserHomeDir`, `filepath.Abs`) runs before it, and neither writes.
  - It covers `--from-zip=`, `--from-zip ""` and both dash forms, since `reorder` lists both.
  - Ordering: other command first, then empty value, then `--check`/`--repo`.
- **`update --dry-run` sets `--check`.**
  - Release path: `shedKeptClone` is skipped (`!*check`). `UpdateFromRelease` receives `check=true`, and `ErrNewerRelease` becomes the message plus exit 0. Tested.
  - Checkout path: `update()` returns before `shedKeptClone`/`Pull`/`Reinstall` under `checkOnly`. `Behind` still runs `git fetch`, which updates remote-tracking refs, the same as `--check` always has. This is not a finding, but "changes nothing" in the comment at `main.go:229-231` holds only for the working tree and the home.
  - Checkout path untested: finding 3.
- **`diff` zip branch.** Correct and tested.
  - `src` is resolved before the zip check, so a DEV build on a zip home with no checkout errors with `ResolveSource`'s zip note rather than "not compared". That is acceptable and states the way out.
- **`sweepStaleZipDirs`.** Correct.
  - It globs only `<home>/.config/sh-agent/zip-*`. `Lstat` plus `IsDir` means a symlink named `zip-x` is not followed.
  - The age comes from mtime, which the unpack never back-dates because `unpackFile` sets no `Chtimes`.
  - It runs only on a real run, before `UnpackZip`. A refused zip may therefore still sweep hour-old debris; that is harmless and removes only this command's own leftovers.
- **`secretShaped`'s `*oauth-token*`.**
  - Consistent with the context's Git rule (`context/AGENTS.md:45`) and preflight.
  - No tracked harness file matches: `git ls-files | grep -i oauth` is empty.
  - `TestSecretShapedCLITokenFiles` fails without the pattern for `antigravity-oauth-token`.
  - The pattern is broad (a future doc named `*oauth-token*.md` would be refused at install and adopt). No such file exists, so this is not a finding.
- **`.gitignore` line.** `git ls-files -ci --exclude-standard` is empty, so no tracked file becomes ignored. Pass.
- **Network or credentials from tests.**
  - On the green path, no test reaches either: the full suite passed, and `TestReleaseUpdateDryRunChangesNothing` and `TestFromZipRefusedOnOtherCommands` isolate `PATH`. No package uses `net/http`.
  - On the red path, `TestFromZipEmptyValueRefused` can, and so can the zip tests that set a release version (finding 2).
  - Whether the main agent's "empty-value" mutation ran a real `git fetch` depends on how it was mutated. Deleting the `case *fromZip == ""` line routes `""` into `UnpackZip` and fails safely. Reverting to the round-1 `!= ""` gate reaches `git fetch`. The dispatch does not say which one was run.

## Criteria

1. **Dead code: pass.** No unused functions, commented-out blocks or debug output.
   - `short()` has two callers (`status`, `knownCheckout`). `repo.go` inlines the same truncation, which is duplication, not dead code.
2. **Flow tracing: fail.** Findings 1 (a refusal after creation, across the `UnpackZip` → `installFiles` hop) and 2 (a red path reaching the network).
   - Errors are wrapped with the zip path.
   - Resources are closed on every path in `zip.go`.
3. **Standards conformance: fail (minor part).**
   - gofmt, vet and `go test -race ./...` are clean.
   - Stdlib only.
   - The stale comment at `main.go:1051-1052` is part of finding 1.
   - Missing tests for changed behaviour: findings 3 and 4.
4. **Directory layout: pass.** Tests sit beside `zip.go` and `main.go`.
5. **Simplicity (DIY-except-security): pass.**
   - `archive/zip`; `crypto/sha256` is used as a checksum, as designed, which is not hand-rolling.
   - `mkdirAllTracked` is small and needed for exact cleanup.
6. **Decoupling: pass.**
   - `zip:sha256:` parsing is repeated in three places (`status`, `knownCheckout`, `ResolveSource`), as in round 1, with no failure mode.

## Not examined

- `internal/harness/zip_test.go` lines 1-130 and 176-335 were skimmed by test name and parent setup only, not read line by line.
- `readCapped`'s body: it is outside `zip.go` and I did not open it. Its behaviour is inferred from `TestUnpackZipRefusesOversizedArchive`'s name.
- `UpdateFromRelease`'s internals: the release dry-run was judged from its call site and its test.
- The docs hunks (README, delivery.md, deploy guides, REQUIREMENTS, DECISIONS/ADR-038, skill-matrix, VENDORED.md): that is the docs review's scope.
- The phase-1/phase-2 verify logs, and the e2e log beyond its Not exercised and D11 rerun sections.
- No mutations of my own. Finding 2's fall-through was traced by reading, not run, because running it would fetch from the real remote.
- macOS.

Scratch used: `.agents/scratch/rv2-zip/` (scratch binary, `nodesc.zip`, `noskills.zip`, homes `home1`, `home2`).
