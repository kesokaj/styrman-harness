---
owner: reviewer
date: 2026-09-28
---
# Review — update-from-zip-code   (2026-09-28)

Range reviewed: 1cb2608..worktree (`git diff 1cb2608` plus the untracked `internal/harness/zip.go`, `internal/harness/zip_test.go`, `cmd/sh-agent/zip_test.go`)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem is that BUG-093 is only half fixed. `UnpackZip` creates `~/.config/sh-agent` once the NAMES pass, and five refusals come after that point: a symlink, a non-regular mode, a duplicate name, the size limits, and the final `IsRepo`. A zip refused for any of them still leaves an empty `~/.config/sh-agent` on a fresh home. I reproduced this with the CLI for the symlink and duplicate-name cases. The code comments claim the opposite, and no test can see it: the internal refusal tests pass a parent that already exists. Second, an EMPTY `--from-zip` value (`--from-zip "$ZIP"` with the variable unset, or `--from-zip=`) falls through to the normal `update`. On a release binary that downloads and replaces the binary, which is the opposite of what the flag promises, and it does so even with `--dry-run`. I reproduced this in a scratch home. Third, the "--from-zip belongs to update" refusal sits after the `uninstall`/`report`/`version` switch, so `uninstall --from-zip x` uninstalls. gofmt, vet and `go test -race ./...` are clean.

## Findings

- **[medium / confidence: high]** `internal/harness/zip.go:57-63` (parent created) vs `:90-93`, `:95-98`, `:100-101`, `:104-105` (later refusals). This is a recurrence of BUG-093 ("a command that can refuse its input creates nothing before the refusal point").
  - The code calls `os.MkdirAll(parent)` after the name pass. It then refuses on symlink, non-regular mode, duplicate name (the `O_EXCL` failure in `unpackFile`), per-file or total size, and the final `IsRepo`. The deferred cleanup removes only `tmp`, never the parent it created.
  - Reproduced with the scratch binary `.agents/scratch/rv-zip/sh-agent` (`-X main.version=v0.0.1`) and a fresh home for each case. `update --from-zip sym.zip` exits 1 and leaves `home/.config/sh-agent`; `dup.zip` leaves the same.
  - Two comments overclaim and are stale under the describing-text rule: `zip.go:57-58` ("a refused zip leaves nothing behind, not even the directory") and the `updateFromZip` comment in `cmd/sh-agent/main.go` ("so a refused zip writes nothing").
  - Tests that pass for the wrong reason: every refusal case in `TestUnpackZipRefuses`, `TestUnpackZipLimits`, `TestUnpackZipRefusesDuplicateName` and `TestUnpackZipNotHarness` passes `parent := t.TempDir()`, which already exists, and checks only `zip-*`. They cannot observe parent creation. The only fresh-home CLI test, `TestUpdateFromZipRefusedZipWritesNothing`, covers just the name refusal.
  - Fix, two parts:
    - Move the mode and duplicate-name checks into the first pass over `zr.File`. Both are known from the headers.
    - For the checks that need content (size, `IsRepo`), record whether `parent` existed before `MkdirAll` and remove it on error when this call created it.
  - Add fresh-home CLI tests for a symlink, a duplicate name and an oversize file. Each asserts the home is empty.

- **[medium / confidence: high]** `cmd/sh-agent/main.go:130` (`if *fromZip != ""`). An empty value is treated as "flag absent".
  - `sh-agent update --from-zip "$ZIP"` with `$ZIP` unset, or `--from-zip=`, runs the ordinary `update`. On a release binary that is `UpdateFromRelease`: it downloads the newest release and replaces `~/.local/bin/sh-agent`. That contradicts the stated contract "nothing is fetched and the binary itself is not replaced" (R9 zip-place sub-bullet; README § Install from a zip).
  - Reproduced: `update --from-zip= --home <scratch> --dry-run` printed `downloading sh-agent-linux-amd64 v1.7.0` / `verified and installed …` and wrote `<scratch>/.local/bin/sh-agent`. The real `~/.local/bin/sh-agent` was untouched (mtime 16:03, before the run).
  - Fix: detect the flag was set (`fs.Visit`, name `from-zip`) and refuse an empty value before anything else. Add a test that `update --from-zip=` returns an error and writes nothing.

- **[medium / confidence: high]** `cmd/sh-agent/main.go:113-122` vs `:130-133`. The `cmd != "update"` refusal runs after the `version`/`uninstall`/`report` switch has already returned.
  - `sh-agent uninstall --from-zip x.zip` uninstalls. Reproduced with `--dry-run`: exit 0, "would remove 258 files and the manifest". `report --from-zip x` sends a report.
  - The comment at `:127` says "Refusals come before anything is read". `TestUpdateFromZipRefusals` covers only `install`, which reaches the check.
  - Fix: move the `*fromZip != "" && cmd != "update"` refusal above the first `switch cmd`. Extend the refusal table with `uninstall` and `report`, and assert the manifest is still there after `uninstall`.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:203-212` (`diff`). `status` now says a zip-installed home's drift is "not compared" because the zip is not kept. `diff` has no zip branch: with a release binary it diffs the installed zip harness against the binary's embedded files and prints every difference between the zip's `main` and the release as drift. No test covers it. Fix: the same `ZipPrefix` check as `status`, print "not compared", return 0.

- **[low / confidence: high]** `cmd/sh-agent/main.go` `updateFromZip`: interrupted runs leave the unpacked tree behind.
  - The unpack directory is removed only by `defer`. A run killed mid-way, e.g. by SIGPIPE on `update --from-zip f.zip --dry-run | head -2`, leaves the whole unpacked tree in `~/.config/sh-agent/zip-*`. Reproduced: the dry run left `zip-1845427001/…` under an existing config dir.
  - Nothing ever sweeps stale `zip-*` directories. That breaks "a dry run writes nothing" for a realistic invocation.
  - Impact is limited: `report` collects named files only (`report.go:262-287`), so the leftover is not sent.
  - Fix: remove stale `zip-*` under the parent at the start of a non-dry run, or accept it and say so in the docs.

- **[low / confidence: high] pre-existing, outside the range, reached through finding 2.** `cmd/sh-agent/main.go:214-236`: `update` on a release binary ignores `--dry-run`. `UpdateFromRelease` receives `*check` only, so `update --dry-run` downloads and replaces the binary. Not introduced by this diff. Recorded because finding 2 makes it reachable from a command line that asked for a dry run.

## Items the dispatch asked about

- **Tests that pass for the wrong reason:** the internal refusal tests' pre-existing parent (finding 1). The others hold:
  - `TestUpdateFromZip` proves the zip, not the embedded files, was installed (`ZIPMARKER` under a release version).
  - `TestUpdateFromZipDryRunWritesNothing` would fail if the D10 `parent = ""` branch were removed, because `MkdirAll` would create `.config`.
  - `TestUnpackZipRefuses` proves each fixture holds its condition before trusting the refusal.
- **Behaviour with no test that fails without it:**
  - Refusals after the name pass on a fresh home (finding 1).
  - The empty value (finding 2).
  - Non-update commands other than `install` (finding 3).
  - `diff` on a zip home (finding 4).
  - The dry-run branch that unpacks under an EXISTING `~/.config/sh-agent`: its cleanup is covered only indirectly, by the shared `defer` that `TestUpdateFromZip` exercises.
- **Read-only mode with a side effect:**
  - A dry run with an existing `~/.config/sh-agent` unpacks the full tree into the home and deletes it afterwards. This is by design (plan D10, README says so) and is cleaned on normal exit. It is not cleaned on a kill (finding 5).
  - `--from-zip=` combined with `--dry-run` downloads and replaces the binary (findings 2 and 6).
- **Error path leaving state:** finding 1. When `Install` fails after unpacking, the tree is removed by the `defer`.
- **Flag handling:**
  - `reorder`'s value table has both `-from-zip` and `--from-zip`.
  - `--from-zip --dry-run` with the file forgotten takes `--dry-run` as the path and fails at `open`, which is safe.
  - The `--check`/`--repo` refusals come before the zip is read (`TestUpdateFromZipRefusals` would fail otherwise).
  - The non-update refusal is misplaced (finding 3), and an empty value is not refused (finding 2).

## Criteria

1. **Dead code: pass.** No unused functions, commented-out blocks or leftover debug output.
2. **Flow tracing: fail.**
   - Findings 1–4.
   - `Source.Origin` goes through `Install` correctly (manifest records `zip:sha256:…`, no commit, verified by `TestInstallRecordsZipOrigin`).
   - `ResolveSource` never treats the marker as a path (decoy test).
   - `knownCheckout` names the zip.
   - `status` reads the manifest first and presence-checks with exit 1 on missing.
   - Every file and reader opened in `zip.go` is closed on every path.
3. **Standards conformance: fail (minor part).**
   - gofmt, vet and `go test -race ./...` are clean (run 2026-09-28).
   - Stdlib only.
   - Errors are wrapped with the zip path.
   - The overclaiming comments in finding 1 and at `main.go:127` break the describing-text rule.
   - The "three ways in" restatement is updated in `AGENTS.md` rule 1, `README.md`, `delivery.md`, `deploy/{local,remote}.md` and `REQUIREMENTS.md` R8/R9 (text marked operator-approved, ADR-038).
4. **Directory layout: pass.** `zip.go` and its test are in `internal/harness`, the CLI test is beside `main.go`.
5. **Simplicity (DIY-except-security): pass.** `archive/zip`, plus `crypto/sha256` used as a checksum, which is stdlib as designed and not hand-rolled. The limits are variables only so tests can lower them.
6. **Decoupling: pass.** The `zip:sha256:` parsing and the 12-character truncation appear in three places (`status`, `knownCheckout`, `ResolveSource`). That is a small duplication with no failure mode today.

## Not examined

- ADR-038's full text in `.agents/docs/architecture/DECISIONS.md`: only its effect on R8/R9 and the delivery spec was read.
- The BUG-093 record body.
- The plan except § Interfaces, D10 and the phase checkboxes.
- The cold-read reviews and the verify logs.
- README prose beyond the diff hunks.
- `TestHarnessMap` against the new map lines (it passed as part of the suite; its content was not read).
- `internal/harness/release*.go`: finding 6 comes from observed behaviour and the call site only.
- No mutation runs of my own; the verifier log's mutation claims were not re-checked.
- macOS behaviour (a case-insensitive filesystem with two names differing only in case is refused by `O_EXCL`, by reasoning, not run).

Scratch used: `.agents/scratch/rv-zip/` (the scratch binary, crafted zips, scratch homes).
