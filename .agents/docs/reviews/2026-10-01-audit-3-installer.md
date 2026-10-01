---
owner: reviewer
date: 2026-10-01
---
# Review — audit-3-installer   (2026-10-01)

Range reviewed: no diff; scope at HEAD 2ba4a50: `cmd/sh-agent/*.go`, `internal/harness/*.go` and their tests, `go.mod`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/sh-agent-matrix.md`, `README.md` § Daily commands. Checklist: recon section `## IN` of `.agents/docs/plans/2026-09-30-recon-audit-fix.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No defect above low severity. Every IN and S item the plan kept is fixed or partly fixed; nothing is open at medium or higher. The worst items are low: a contradictory `dream` flag combination that silently drops one flag, a manifest written by a failed FIRST install that carries no source or commit, and two exported functions (`Lint`, `Parse`) left with only test callers after S8. Measured: `gofmt -l` clean, `go vet ./...` clean, all tests pass; coverage cmd/sh-agent 85.7%, internal/harness 89.0% (plan floor 84%); `govulncheck ./...`: 0 reachable (1 in a required module, not called); staticcheck not on PATH (not run); no system temp dir use in non-test source (grep empty).

## Findings

| ID | was | now | evidence |
|---|---|---|---|
| IN1 | med: over-redaction in `Redact` | fixed | `Redact` probed through a scratch program: `git checkout -b feature/x`, `## Secrets: how it works`, `the token: field is set`, `passwords: are hard` survive; `password=correct`, `--password letmein`, `API_KEY=sk-live-...` still go. `?sig=<16+>` is still redacted by the Teams/Azure SAS needle on purpose and pinned as accepted over-redaction at `report_test.go:146-148`. `token=abcdefg` stays (accepted under-redaction, documented at `report.go:~70`). |
| IN2 | med: manifest saved only at the end | fixed | `apply.go:36-56` `saveOnError` saves (written so far) ∪ (old entries still on disk, not links) on write and stale-remove errors; `TestInstallMidFailureKeepsManifestComplete` (`harness_test.go:1946`) makes an unwritable directory, asserts the new file is listed and uninstall removes it. Skips as root (fails closed to a skip, not a false pass). |
| IN3 | low: symlinked destination followed | fixed | `apply.go:59-65` Lstat plus `underLink` before any read/chmod; Status (`:~183`), Diff and Uninstall use the same pair. `underLink` also covers linked parent dirs. |
| IN4 | low: `update --dry-run` skipped the dirty refusal | fixed | `main.go:update` checks `Dirty`/`StatusUnknown` above the `checkOnly` return; probe `update --dry-run --check --repo .` printed the refusal, exit 1. README row (`README.md:185`) says so. |
| IN5 | low: one FlagSet for every verb | fixed | `verbs` table (`main.go:96-165`), one FlagSet per verb; probes: `list --check`, `list --force`, `status --dry-run`, `version --dry-run`, `index --kind bug`, `search --dry-run`, `adopt --all` all exit 2 naming the flag; `diff foo`, `install -- --dry-run` exit 2 "takes no arguments". `reorder` reads value flags from the FlagSet; the hand-kept map is gone. |
| IN6 | low: origin string parsed by hand at six sites | fixed | `ParseOrigin`/`State.Origin`/`Source.OriginString`/`ZipOrigin` in `repo.go:122-187`; the only remaining uses of the prefix constants are in `repo.go`. |
| IN7 | low: `report` re-implemented source resolution | fixed | `harness.ResolveSourceLenient` (`repo.go:106`), called from `main.go:report`. |
| IN8 | low: no lock across installs | struck (D33) | not re-opened. |
| IN9 | low: always-nil error returns | fixed | `pathHint`, `AttributionHint`, `AttributionSet` return nothing; `go vet` clean. |
| IN10 | low: system temp dir | fixed | grep over non-test source empty; `VerifyZipSums` uses `os.MkdirTemp(dir, ...)` beside the zip, `UnpackZip` with `parent==""` unpacks beside the zip; `TestNoSystemTempDirInSource` plus the behavioural `TestVerifyZipSumsWritesNothingUnderTempDir` (TMPDIR watched). See finding 3 for what the dry run still writes. |
| IN11 | low: dream timer `--dry-run` missing from usage | fixed | `usage` and the doc comment both show `[--dry-run]`. |
| IN12 | low: `--install-timer --remove-timer` exit 1 | fixed | probe: exit 2, message on stdout. |
| IN13 | low: refusals to `os.Stderr` | fixed | the only `os.Stderr` left in non-test source is `main()`'s generic error line (`main.go:84`). |
| IN14 | low: `Inspect` three times | fixed | `update()` inspects once, and again only after a pull moved the tree (`main.go` "inspected again"); `status` once. |
| S1 | move `attribution*` to harness | fixed | `internal/harness/settings.go` (213 lines); main only calls `AttributionSet`/`AttributionHint`. |
| S2 | split `run` | fixed | `run` is 88 lines (was 324); `main.go` 1142 lines (was 1255 at the Phase 6 commit); 15 table entries = 14 verbs + `sync`, the same set as the pre-audit usage (D9, D30). |
| S3 | collapse verbs | struck (D30) | verbs unchanged; verified against the usage text. |
| S4 | delete `adopt` | struck (D30) | `adopt.go` still 177 lines, still a verb. |
| S5 | cut incident-narrative comments | partly | comments citing `BUG-`/review/audit still number 22 in `report.go`, 14 in `main.go`, 5 in `repo.go` (one grep over all three words; the audit counted by another method: 27/19/7). Several are now rationale; many still read as a review log (e.g. `main.go:158-161` "security review 2026-10-01, round 2"; `apply.go` WriteAtomic's "Round 2 of the 2026-09-30 review"). Not a defect; it is the one S item that was a judgement call. |
| S6 | merge PATH walkers | fixed | `pathCandidates` (`pathguard.go:48`) is the one walk for `PathShadows` and `RunningElsewhere`. `internal/dream/timer.go:154` has its own `SplitList` loop; it is outside this slice. |
| S7 | one `readCapped` | fixed | `ReadCapped` (`release.go:359`) used by release, zip and settings; fails closed past the cap (FIFO case tested at `phase8_test.go:213`). `zip.go`'s `copyCapped` is a stream copy with a different job. |
| S8 | `Lint` as a flag on `Parse` | partly | `ParseStrict` does lint + parse in one pass, but `Lint` (`frontmatter.go:120`) and `Parse` (`:70`) stay exported with no production caller (see finding 1). |

Counts: fixed 17 (IN1-7, IN9-14, S1, S2, S6, S7), partly 2 (S5, S8), open 0, struck 3 (IN8, S3, S4).

New findings:

1. **[low / confidence: high]** `internal/harness/frontmatter.go:70,120` — `Parse` and `Lint` are exported and called only from tests (`phase8_test.go:165`, `harness_test.go:88-101,2333-2353`); production uses `ParseStrict`. S8 added a function rather than removing one. Delete `Parse`/`Lint` (or unexport them as test helpers) and point the tests at `ParseStrict`.
2. **[low / confidence: medium]** `cmd/sh-agent/main.go:runDream` — a timer flag together with a run flag is accepted and half ignored: `dream --all --install-timer` and `dream --project DIR --remove-timer` (probed with `--dry-run`, exit 0) act on the timer and drop `--all`/`--project`, while `--all --project` and `--install-timer --remove-timer` are refused with exit 2. Add a case refusing `(installTimer||removeTimer) && (dreamAll||dreamProject != "")` with the same exit 2.
3. **[low / confidence: medium]** `cmd/sh-agent/main.go:updateFromZip`, `release.go:VerifyZipSums` — `update --from-zip --dry-run` is not read-only: it writes a signers directory beside the zip and unpacks the tree (into `~/.config/sh-agent/` when that exists, beside the zip when not, so it resolves a different path than the real run, which always unpacks under the config dir). Both are removed on exit and README/delivery.md say so (`README.md:176`), but the README's "without writing a file" and the test name `TestUpdateFromZipDryRunWritesNothing` overstate it, and a zip in a read-only folder fails the dry run. Fix by wording ("writes and removes a temporary directory") or by planning from the `zip.Reader` without unpacking.
4. **[low / confidence: medium]** `internal/harness/apply.go:36-56` — `saveOnError` on a FIRST install (no manifest yet) saves the zero `State` with only `Files` set: empty `Repo`, `Commit`, `InstalledAt`. `status` then prints `installed  <empty> from commit  (...)`, and `knownCheckout` treats the empty repo as "no checkout on record". Cosmetic and uninstall still works, but fill `Repo`/`Commit`/`InstalledAt`/`Vendor` from the same values the success path computes.
5. **[low / confidence: medium]** tests that can miss a plausible regression: `TestNoSystemTempDirInSource` (`phase8_test.go:16`) greps for `MkdirTemp("", ` and `CreateTemp("", ` only, so `os.MkdirTemp(os.TempDir(), ...)` or `os.TempDir()` passes (the behavioural TMPDIR test covers `VerifyZipSums` and the unpack, not every site); and `TestUsageMatchesDocComment` (`search_test.go:82`) compares usage with the doc comment but nothing compares the `verbs` table with either, so a verb added to or dropped from the table (D9/D30: verbs unchanged) is caught by no test (`TestHarnessMap` checks only that mapped verbs appear in usage). Add `os.TempDir` to the needles and one test that the `verbs` keys equal the usage verbs plus `sync`.

Not findings, checked: `run` refuses `--from-zip` on every other verb by a raw scan (`hasFlag`) while `update` decides from the parsed FlagSet; probed `uninstall --from-zip x`, `install --from-zip x` (refused, exit 1), `update --home -- --from-zip z` (decided from the FlagSet, covered by `TestUpdateFromZipDecidedByParsedFlags`). Usage errors that are not flag/argument errors (`adopt` with no name, `memory` with no subverb, `search` with no words, an unknown command, `--from-zip` without a file) exit 1 while README line 199 and D41 promise exit 2 only for an undefined flag or a stray argument, so no contradiction. `TestWriteAtomicClosesBeforeRename` is a source-text pin but follows the § 2 exception (code lines only, comment says what cannot be observed). `TestVerbRefusesAnotherVerbsFlag` and `TestInstallMidFailureKeepsManifestComplete` would fail if the single FlagSet or the end-only manifest save came back (the second only for the write-error path; the stale-remove path is covered by `TestInstallResyncFailureKeepsManifest`, name read, body not).

## Criteria
1. Dead code — fail (low): `Parse` and `Lint` have no production caller (finding 1); no other exported harness function is test-only except `SwapReleaseSigner`, which says so.
2. Flow tracing — pass: verb table to `run`, `update` dirty check above the `--check` return, `saveOnError` on both error paths, `ReadCapped` callers, origin parsing in one place; two low gaps (findings 3, 4).
3. Standards conformance — pass with a nit: gofmt/vet clean; guards keep data on error (stale removal and uninstall skip links and stop on error); describing text mostly current (finding 3 names one overstatement; S5 comment narrative partly remains).
4. Directory layout — pass: settings code in `internal/harness/settings.go`, tests beside code, nothing new outside `cmd/` and `internal/`.
5. Simplicity (DIY-except-security) — pass: stdlib only in both packages, `go.mod` unchanged; signature verification uses `ssh-keygen`, nothing hand-rolled.
6. Decoupling — pass: `main.go` is a table plus one function per verb; harness holds the attribution, origin and source logic; `internal/search`, `dream`, `memory` imported only through their entry points.

## Not examined
`report.go` beyond `Redact`'s needles and the probes above (the clone/push path and `TestReportRedact*` bodies not read); `inventory.go`, `plan.go`, `derive.go`, `render.go`, `frontmatter.go` bodies beyond the parse/lint section; `release.go` download/verify internals beyond `UpdateFromRelease`, `VerifyZipSums`, `ReadCapped`; `zip.go` entry checks (`checkZipTrees`, `zipNameSafe`) read only at the unpack seam; `adopt.go`, `pathguard.go` `RunningElsewhere` semantics; most of the 4,000+ lines of test code (names and two to four tests read; no mutation run); `delivery.md` and `README.md` read only on the lines the code claims touched (no restatement-grep script run); staticcheck unavailable. Probes ran a scratch build against `.agents/scratch/audit3-home`; no real home was touched.
