---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p8-code   (2026-10-01)

Range reviewed (round 1): working tree (uncommitted) vs HEAD 7c2e30a: cmd/sh-agent/main.go, main_test.go, search_test.go; internal/harness/{apply,frontmatter,inventory,pathguard,release,repo,report,zip}.go; README.md; .agents/docs/testing/sh-agent-matrix.md; untracked internal/harness/{settings.go,settings_test.go,phase8_test.go,repo_origin_test.go}

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No correctness or security defect found. `go vet`, `gofmt -l` and `go test -race -count=1 ./...` are clean on this tree (all six packages ok). The refactor keeps the verbs' outputs; the worst problem is one user-visible exit-code change (a flag error now exits 2, was 1) that collides with `update`'s documented exit 2, and two new tests that cannot fail for the behaviour they name.

## Findings
- **[medium / high]** `cmd/sh-agent/main.go` (`run`, the `c.fs.Parse` error branch) — a flag the verb does not define now returns `exitCode(2)` with nothing on stderr; the old single FlagSet returned the parse error, so `main` printed `sh-agent: flag provided but not defined` and exited 1. Measured with a built binary: `install --bogus` rc=2 (was 1), `update --bogus` rc=2. `update --check` documents exit 2 as "behind", so a script cannot tell a typo from "update available". Plan D9/D30 says verbs unchanged; the matrix row and README say "usage error" but not the code change. Fix: return the parse error (exit 1), or state exit 2 in the README/spec/plan and accept the collision on `update` — operator's choice.
- **[medium / high]** `cmd/sh-agent/zip_test.go:470-485` `TestUpdateFromZipDryRunCleansTempDir` (not in the diff, made stale by it) — comment and check still say the dry run unpacks under the system temp directory. Since IN10 the unpack is beside the zip, so the test watches a directory nothing writes to. Mutation (scratch copy, compiled, vet clean): replacing the deferred `os.RemoveAll(dir)` in `updateFromZip` with a no-op fails `TestUpdateFromZip` and `TestUpdateFromZipRefusedContentWritesNothing` (real runs) but NOT this test — a dry run on a fresh home that leaves `zip-*` beside the zip is unpinned at the verb level. Fix: compare the zip's folder entries before/after, drop the TMPDIR claim.
- **[low / high]** `internal/harness/phase8_test.go` `TestVerifyZipSumsWritesNothingUnderTempDir` — its TMPDIR assertion cannot fail: the old code removed its temp dir by defer too. Mutation (scratch copy, compiled): reverting to `os.MkdirTemp("", "sh-agent-verify-*")` leaves the test green; only the source-text `TestNoSystemTempDirInSource` fails. The "folder unchanged" half is real (mutation removing the defer fails it). Drop the TMPDIR half or accept that the source scan is the pin.
- **[low / high]** `internal/harness/phase8_test.go:14` `TestNoSystemTempDirInSource` — scans `internal/harness` and `cmd/sh-agent` only; `internal/dream`, `internal/search`, `internal/memory` are not read, so a `MkdirTemp("", ` added there passes. `os.TempDir()` is no needle either. Fix: walk `../` (all of `internal` plus `cmd`) and add `os.TempDir()`.
- **[low / medium]** `internal/harness/release.go` `VerifyZipSums` (signers scratch `MkdirTemp(dir, ...)`) and `zip.go` `UnpackZip` with `parent == ""` — both now need the zip's own folder to be writable; before, a signed zip in a read-only folder (mounted media) verified and a no-config `--dry-run` worked via the system temp dir. Now `update --from-zip` fails with "cannot write the signers scratch beside X" (verify) or a MkdirTemp error (dry run). Intended by IN10 and removed on every exit path (I checked each: defer in both, and the unpack dir by `updateFromZip`'s defer), but it is a behaviour change on the fallback path for exactly the machines that cannot take another way in. Operator decision; fallback to `~/.config/sh-agent` (which the real run creates) would keep scratch out of /tmp.
- **[low / high]** stale describing text after the moves: `cmd/sh-agent/main_test.go:410` (`attributionSet`), `:625` and `internal/harness/settings_test.go:50` (`readSettings`) name functions that are now `AttributionSet` / `ReadSettings` in `settings.go`.

## Checks asked for
1. Verb table: `--repo/--home/--vendor` are registered on every FlagSet in `run` before the verb's own; built-binary run of `version list status diff render` with all three globals accepted. `reorder` takes value flags from `fs.VisitAll` minus `IsBoolFlag` (Int/String flags take a value, so every old hand-listed flag is covered: repo home vendor feedback-repo from-zip limit kind refs type name description project; `--flag=value` is one token, not in the map, passes through). `help`, `-h`, `--help`, no args print usage; `version` ok; `sync` shares install's entry; `adopt x y` still errors rc=1, `adopt` positional works with flags either side. A verb flag before the verb (`--home H list`) is "unknown command", as before. Per-verb narrowing (IN5) is documented in README and the matrix; the exit code is the finding above. `hasFlag` ahead of Parse keeps the two `--from-zip` messages.
2. Origin: every former hand-parse site (ResolveSource, status, knownCheckout, runDiff, updateFromZip via ZipOrigin, apply.Install via OriginString) goes through `ParseOrigin`/`Origin()`. Round trips: checkout and plain copy (path, `Dir`), release (`embedded:vX`), zip (`ZipOrigin` with/without verified tag, and pre-verdict `zip:sha256:<sum>`). `ParseOrigin` requires `zip:sha256:` where the old `diff` check used the bare `zip:` prefix; only `ZipOrigin` writes that prefix, so no real manifest differs.
3. `ResolveSourceLenient`: for a release binary it skips only the `--repo`/`$SH_AGENT_REPO` refusal (same as the old inline code, same note text); for a development build it calls `ResolveSource` unchanged and turns its error into the note. It never reads a tree a release binary's `install` would refuse. Equivalent to the old `report`.
4. `settings.go`: moved logic matches the old text line for line (diffed with comments stripped): `Lstat` kept, `O_EXCL` create, `WriteAtomic` for edits, hint texts unchanged; `ReadSettings` now uses `ReadCapped` after the same Stat size gate; `AttributionHint` no longer returns an error (no caller used it).
5. Signers scratch and unpack: removed on every exit path (defer RemoveAll; UnpackZip removes its own dir on error; updateFromZip defers RemoveAll(dir)). Mutations 2 and 4 below confirm the real tests fail without them. Read-only caveat above.
6. `ReadCapped`/`copyCapped`: over-limit is an error at every site (release assets, SHA256SUMS, zip archive, per-entry cap in both `entrySize` and `unpackFile`; `unpackFile` still closes the writer and keeps the first error). `ReadCapped` is a Stat pre-check then `LimitReader(max)`, so a file that grows past max after the Stat is truncated silently — unchanged from the old `readCapped`, and every caller's next step (checksum, parse) fails on a truncated body; not a finding.
7. `ParseStrict` = `entries()` once, `lintEntries` then `parseEntries`; `Lint` is the thin wrapper over the same `lintEntries`; call sites in `LoadSkills/Agents/Docs` unchanged in effect; `TestParseStrict` compares the error text with `Lint`.
8. `pathCandidates` equals both old walks: same empty-entry `./name`, same relative entries, same dedupe; `RunningElsewhere` gains the dedupe, which cannot change its first-match result.
9. Mutations run by me in a scratch copy under `.agents/scratch/p8review/repo` (the working tree was never edited; all four compiled and `go vet` passed before the run): (a) `MkdirTemp("", ...)` back in `VerifyZipSums` -> only `TestNoSystemTempDirInSource` fails (finding 3); (b) `defer os.RemoveAll(signersDir)` removed -> `TestVerifyZipSumsWritesNothingUnderTempDir` fails (leftover folder); (c) `info.Size() > max` weakened -> `TestReadCappedRefusesOversize` fails; (d) `parent = filepath.Dir(zipPath)` -> `""` -> `TestUnpackZipEmptyParentUnpacksBesideTheZip` fails; (e) `RemoveAll(dir)` in `updateFromZip` removed -> two real-run zip tests fail, the dry-run test does not (finding 2). The fork's IN5 claim is consistent with the old code (one FlagSet; `list --check --force` ran); `TestVerbRefusesAnotherVerbsFlag` asserts exit 2 and the "flag provided but not defined" line, so it fails on the old dispatcher.
10. `TestNoSystemTempDirInSource` covers `internal/harness` and `cmd/sh-agent` only (finding 4); a `MkdirTemp("", ` in a new file of either would be caught (gofmt guarantees the space).

## Criteria
1. Dead code — pass. `short()` removed with its last caller; `ReadSettings` is exported with no caller outside its package (harmless inside `internal/`).
2. Flow tracing — pass with findings: the exit-code change (finding 1) and the untested dry-run cleanup (finding 2); errors propagate; resources closed on every path.
3. Standards conformance — pass with nits: stale names in three comments and one stale test comment (finding 6, 2); README updated for the `index` narrowing.
4. Directory layout — pass: settings code in `internal/harness/settings.go`, tests beside it.
5. Simplicity (DIY-except-security) — pass: stdlib only, no new dependency, no hand-rolled security primitive (signature verification still goes through `ssh-keygen`).
6. Decoupling — pass: the verb table and `ctx` replace one 600-line switch; settings moved out of `cmd`; one origin parser.

## Not examined
`main_test.go` and `search_test.go` were read only where they touch reorder, the per-verb flags and the settings moves (not every adapted assertion); `report.go`, `frontmatter.go`, `pathguard.go`, `apply.go` and `inventory.go` comment-trim hunks were read in the diff but not each comment checked against code; Antigravity runs, the real-home smoke and `govulncheck` were not run; the two untracked verifier/review artifacts of this phase (`testing/2026-10-01-verify-audit-fix-phase-8.md`, `reviews/*p8-docs`, `security-worktree-p8`) were not read. Scratch artifacts of this review: `.agents/scratch/p8review/`.

## Outcome
- [medium] undefined flag exits 2, message on the writer, overlapping `update --check`'s exit 2 — rebutted in part, fixed in part: exit 2 for usage errors stays (D41: the `dream` verbs set the precedent, IN12; the message names the flag so a human is never misled), and README § Daily commands documents the overlap; the message goes to the verb's writer as before (not a regression, as the verifier measured)
- [medium] `TestUpdateFromZipDryRunCleansTempDir` watches TMPDIR only — fixed at `cmd/sh-agent/zip_test.go` (the zip's own folder is watched too); must-fail: `rerr := error(nil)` in place of the `RemoveAll`, the test failed on `zip-…` left beside the zip, restored. `TestVerifyZipSumsWritesNothingUnderTempDir`'s TMPDIR half stays as a belt; its zip-folder half and the source scan are the pins
- [low] `TestNoSystemTempDirInSource` scans two packages — fixed (walks every package under `cmd/` and `internal/`)
- [low] a zip on read-only media no longer verifies — documented at `README.md` § Install from a zip and `delivery.md` § The model (the zip's directory is the operator's own writable folder)
- [low] stale function-name comments — fixed (`takesValue` in three test comments and the guide)

---

# Round 2   (2026-10-01)

Range reviewed: working tree — `git status --short --untracked-files=all` (Phase 8 of the audit-fix plan); `git diff HEAD -- <files>`; untracked files whole; round 2 read the round-1 fixes in `cmd/sh-agent/main.go` (`reorder`, the verb table, the stray-argument refusal, the `--from-zip` pre-check), `internal/harness/{release.go,settings.go,apply.go}`, `phase8_test.go`, and the two new tests in `main_test.go`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

All five round-1 fixes hold; `gofmt -l`, `go vet` and `go test -race -count=1 ./...` are clean (six packages ok). One defect the `--` fix introduced, low severity: `reorder` keeps the `--` marker where it stood, but the flag package stops at the first non-flag word, so a `--` that follows a positional word is never consumed and becomes a positional itself. Measured with a built binary: `search foo -- bar` searches for the three words `foo -- bar` (the old code dropped the marker and searched `foo bar`), `adopt foo --` and `memory show --` are now usage errors (rc 1; both worked before), and `list foo --` reports the `--` inside the stray list. No test pins a positional before the marker.

## Findings (round 2)
- **[low / high]** `cmd/sh-agent/main.go:935-963` (`reorder`) — with `positional` words ahead of `--`, `append(flags, positional...)` puts those words first and `Parse` stops there, so the later `--` stays in `Args()`. Built-binary runs: `search --home H foo -- bar` -> hits for "foo -- bar" (was "foo bar"); `adopt --home H foo --` -> usage error (NArg 2); `memory --home H show --` -> usage error. Where `--` leads (the only shape the new tests cover: `reorder(fs, ["--home","h","--","--from-zip","z"])`) it is consumed correctly and `Args()` is exactly the words after it (`search -- --open`, `adopt -- nosuchskill`, `memory -- show` all measured fine). Fix: emit the marker once, after the flags and before every positional, and drop the original: `append(append(flags, "--"), positional...)` with `--` removed from `positional`; the flag package then consumes it and all positionals stay positional whatever their shape. Add a test with a positional before the marker (`search foo -- --open` -> words `foo --open`).

## Round-2 checks asked for
1. `--` kept: `adopt -- name`, `memory -- show`, `search -- words`, `search -- --open` (a flag-shaped word stays positional), `search --limit 1 -- --open x` all work with the marker leading. Broken only when a positional precedes it (finding above). A flag taking a value consumes `--` as its value (`--kind --`), as the flag package itself does; unchanged.
2. Stray-argument refusal: `memory` and `search`/`adopt` are `positional: true`, so `memory add`/`memory show` still parse their subverb; `dream x`, `dream -- x`, `index extra`, `version --home H extra`, `list foo --` exit 2 with `<verb> takes no arguments: [...]`. `help`, `-h`, `--help` and no args return before the table; `version` and `list --` (a lone leading marker) are accepted. `update -- --from-zip nowhere.zip`: `hasFlag` stops at `--`, `Parse` leaves the words, the refusal fires before any update logic runs (no download path reached). Note `index` shares `runSearch` but is not positional, so `index words` is refused where it was ignored; intended.
3. `ReadCapped`: a file of exactly `max` bytes reads `max` bytes through `LimitReader(max+1)` and passes (`len > max` is false); only a stream of more than `max` fails. Callers (release asset, SHA256SUMS, zip archive, settings) already Stat-gate or use this Stat check first, so a regular file's behaviour is unchanged. `max` is a constant at every call site, so `max+1` cannot overflow. (A FIFO still blocks at `os.Open` with no writer; unchanged, and `ReadSettings` rejects non-regular files before it.)
4. `underLink` in `AttributionSet`: it walks `Dir(rel)` upward from `.claude` and never examines `home` itself, so a symlinked `$HOME` with a real `.claude` under it is not refused (right: a linked home is the operator's choice); only a linked `home/.claude` is refused, with the hand-edit hint, including a dangling one that `MkdirAll`/`O_EXCL` would otherwise have followed. Consistent with every other installed path in `apply.go`.
5. Re-read declarations touched: `ctx`/`verb`/`verbs` fields all used (`positional` read in `run`); `takesValue` remains only as a local in `reorder` (the guide no longer names a function of that name); `TestNoSystemTempDirInSource` now walks every package under `cmd/` and `internal/`; the zip dry-run test watches the zip's folder (the author's must-fail protocol was not repeated by me).

## Criteria (round 2)
1. Dead code — pass.
2. Flow tracing — pass with the finding above (marker handling).
3. Standards conformance — pass; README/guide restatements of exit 2 and the writable zip folder read and consistent.
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass; no new dependency.
6. Decoupling — pass.

## Not examined (round 2)
The author's must-fail mutations were not repeated; `search.md`, `delivery.md`, the README wording and the verifier/security artifacts were not read line by line; Antigravity runs, the real-home smoke and `govulncheck` were not run. Scratch: `.agents/scratch/p8r2/` (a built binary and a scratch home).

## Outcome (round 2)
- [low] a positional before `--` left the marker in `Args()` — fixed at `cmd/sh-agent/main.go` `reorder` (one marker re-emitted after the flags, ahead of every positional, only when the input had one); `TestDoubleDashKeepsFollowingArgumentsPositional` gains the `foo -- bar` case, the two reorder expectations updated; measured with the built binary: `search foo -- bar` searches "foo bar", `list --` exits 0, `version -- x` exits 2; `adopt foo --` now reaches adopt (exit 1: no such skill), as before the Phase 8 change
- `ReadCapped` at exactly `max` bytes, `underLink` on a linked `$HOME`, the stray-argument refusal per verb — confirmed as the reviewer measured; no change
- the local `takesValue` map inside `reorder` renamed `valueFlags`, so the retired name appears nowhere (the round-2 verifier's one grep deviation)
