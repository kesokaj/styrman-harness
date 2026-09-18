---
owner: explorer
date: 2026-09-18
---
# Recon — sh-agent verb x place-kind audit (2026-09-18)

## Question
For every `sh-agent` verb, on every place kind (checkout / release binary / plain copy), where is that cell handled: exercised, refused, or silent — and which cells have a covering test.

## Summary
Dispatch splits at `cmd/sh-agent/main.go:104` (uninstall/report/version bypass `ResolveSource`) and `:124` (everything else, which needs a `Source`). `ResolveSource` (`internal/harness/source.go` uses `repo.go`'s `IsRepo`) resolves a checkout and a plain copy to the *same* kind of `Source` (`Dir` set, `Version` ""); the checkout/plain-copy distinction is made downstream per-verb via `harness.IsCheckout(repo)` (`repo.go:127`). A release binary is a distinct `Source` (`Dir==""`, embedded FS) chosen only when `IsReleaseVersion(version)` (`source.go:34`). No verb falls through silently for the three place kinds on the main dispatch paths I traced — every cell I found is either exercised or explicitly refused with a message. `adopt` and `update` are the two verbs with a real per-place branch; `install`, `status`, `report` also branch; `list`, `render`, `diff`, `uninstall`, `version` do not branch on place kind at all (uniform code path). Gaps are mostly in *test coverage*, not in silent behavior — see § Risks/untested cells.

## Verb x place-kind matrix (exercised / refused / silent)

| verb | checkout (Dir set, IsCheckout=true) | plain copy (Dir set, IsCheckout=false) | release binary (Dir=="") |
|---|---|---|---|
| install | exercised: rebuild path, `main.go:158-166` -> `harness.Reinstall` (repo.go:201) -> re-exec installs | exercised: no rebuild, installs running binary as-is, `main.go:168-174`, then `installFiles` `main.go:264` | exercised: embedded files + running binary installed, `main.go:130-138` |
| update | exercised: `update()` `main.go:461`, `Behind`+`Pull`+`Reinstall`; dirty tree **refused** `main.go:484-490` | refused: `Behind(repo)` returns `errPlainCopy` (repo.go:118,151-154), surfaced as error text "not its own git checkout; a plain copy cannot be updated with git" via `update()` `main.go:462-464` | exercised: `harness.UpdateFromRelease` `main.go:191-219` / `release.go:87` |
| status | exercised: git state line `main.go:353` | exercised (informational, not refused as a command): `main.go:350-351` prints "plain copy, not its own git checkout; update is refused" (Inspect returns empty RepoInfo, repo.go:135-138) | exercised: `main.go:337-339` embedded-source line, no git state |
| report | exercised: `ResolveSource` success -> normal report, `main.go:304-308` | exercised: same path as checkout (Dir set either way; report does not call IsCheckout) | exercised: bypasses `ResolveSource` entirely, uses `EmbeddedSource` directly and warns if `--repo` given, `main.go:298-303` (deliberate, per comment: avoids the ResolveSource refusal `report` would otherwise inherit) |
| adopt | exercised: `harness.Adopt(repo,...)` `main.go:252`, no IsCheckout check inside `adopt.go` | exercised: same as checkout — `Adopt` only needs a readable tree with `skills/`, not a git checkout | refused: `main.go:249-251` "adopt needs a harness checkout; this install is a release binary" (repo=="") |
| list | exercised: `harness.LoadInventory(src)` `main.go:226`, no place branch | exercised: same code path | exercised: same code path, reads embedded FS |
| render | exercised: `main.go:234-243`, no place branch | exercised: same | exercised: same, embedded FS |
| diff | exercised: `main.go:180-189`, `harness.Plan`+`harness.Diff`, no place branch | exercised: same | exercised: same |
| uninstall | exercised: `main.go:108-109` -> `harness.Uninstall` (apply.go:241) — dispatched *before* `ResolveSource`, reads only the manifest under `home`, never touches `src`/place kind at all | exercised: same, place-kind irrelevant | exercised: same, place-kind irrelevant |
| version | exercised: `main.go:105-107`, `buildVersion()` only, no `Source` involved | exercised: same | exercised: same |
| sync (alias of install) | same as install — `main.go:125` `case "install", "sync"` | same as install | same as install |

No `silent` cell found on the traced dispatch paths for these ten verbs across the three place kinds — every branch either runs code or returns/prints a refusal. (I did not exhaustively trace every helper under `internal/harness/apply.go`'s per-file switch or `report.go`'s ~900 lines; see § Not examined.)

## First install into empty home vs update over existing install

- First install: `installFiles` (`main.go:264`) always runs regardless of ordering; `harness.Install` (`apply.go:19`) writes files + manifest fresh (no prior `state.json`).
- `update` (checkout): `main.go:220-223` -> `update()` `main.go:461` -> `Reinstall` (`repo.go:201-236`, builds then re-execs `install --repo <repo> ...`). Not `Install` directly — always goes through a rebuilt binary's own `install`.
- `update` (release): `main.go:191-219` -> `harness.UpdateFromRelease` (`release.go:87`), which downloads+verifies+`writeAtomic`s the binary then execs `<dest> install ...` (`release.go:172-178`). Distinct code path from `Reinstall`; both eventually call the same `installFiles`/`Install` via a subprocess, not directly.
- `update` (plain copy): refused before any install step (`errPlainCopy`, see matrix above) — never reaches `Reinstall` or `UpdateFromRelease`.
- BUG-079 "note" (process not the installed file) is printed by `harness.NoteRunningElsewhere` (`pathguard.go:112`), called from both `update()` (`main.go:509`, checkout path) and inside `UpdateFromRelease` itself (`release.go:180`) — the plan's Interfaces section says only `UpdateFromRelease` can test the positive case since only it knows "replaced".

## Second table — test coverage per cell

| verb | checkout | plain copy | release binary |
|---|---|---|---|
| install | `TestInstallStatusResyncUninstall`, `TestReinstallBuildsAndInstalls`, `TestInstallMarksDirtyTree`, `TestDryRunNeverBuilds` (main_test.go:209) | `TestPlainCopyNeverRebuilds` (main_test.go:295) | `TestReleaseBinaryWithoutCheckout` (main_test.go:1201) |
| update | `TestUpdateNeverInstallsDirtyTree`, `TestUpdateRefusesWhenStatusIsUnknown`, `TestCheckoutUpdateNamesTheBinaryItRanFrom`, `TestBehindAndPull` (harness_test.go:673, unit level for `Behind`/`Pull`) | **none at CLI level** — `errPlainCopy` only unit-tested via `Behind`/`Pull` directly (`harness_test.go:762-766`, on a *nested* copy, not `update` dispatched through `main.go`) | `TestUpdateFromReleaseReplacesAndInstalls`, `TestUpdateFromReleaseUpToDateAndCheck`, `TestUpdateFromReleaseRefusals`, `TestReleaseBinaryUpdateShedsKeptClone`, `TestReleaseBinaryUpdateCheckExitsTwo` |
| status | implicit inside `TestInstallStatusResyncUninstall`, `TestInstallAndStatusNameAPathShadow` | **none found** — no test asserts the "plain copy, not its own git checkout; update is refused" status line | **none found** at CLI level for the embedded-source status line (`main.go:339`); `TestEmbeddedInstallStateAndRecord` (source_test.go:217) covers install/state, not the `status` command text |
| report | `TestSendReport*` family (report_test.go, many), `TestReportWithoutCheckout` (main_test.go:1027) | not distinguished from checkout (same code path) | `TestReportOnReleaseBinaryIgnoresRepoFlags` (main_test.go:1156) |
| adopt | `TestAdopt`, `TestAdoptRefusesDestructiveSources`, `TestAdoptFailedCopyKeepsTheExistingSkill`, `TestAdoptRefusesSecretShapedFiles` (harness_test.go) — fixtures appear to be git checkouts (`git init` at harness_test.go:536 pattern) | **not found** — no test adopts from a plain-copy repo tree specifically | `case repo == ""` refusal at `main.go:250` — **no test found** asserting the exact "adopt needs a harness checkout; this install is a release binary" message |
| list/render/diff | `TestInventoryAndRender`, diff covered inside `TestInstallStatusResyncUninstall` (uses `harness.Diff`) | not distinguished (same code); no dedicated plain-copy test | `TestEmbeddedInstallStateAndRecord`-adjacent, but no dedicated `list`/`render`/`diff` test against `EmbeddedSource` found by name |
| uninstall | `TestInstallStatusResyncUninstall`, `TestRunInstallStatusUninstall` | place-irrelevant (manifest-only) — covered by the same tests | place-irrelevant — covered by the same tests |
| version | `TestVersionFromInfo`, `TestReleaseFromInfo` | n/a | n/a |
| sync | not found by name (only tested as `install`'s alias, presumably untested as `sync` specifically — grep shows no `"sync"` in *_test.go dispatch) | same | same |

## Flag behavior (deliverable 4) — per-place differences not in the table

- `--repo`: refused outright on a release binary before any command-specific logic other than `report` (`source.go:34-47`, `ResolveSource`); `report` deliberately bypasses this refusal for a release binary and instead just prints a note (`main.go:301-303`).
- `$SH_AGENT_REPO`: same refusal as `--repo` on a release binary (`source.go:36-46`); on checkout/plain-copy dev builds it's tried before repo-walking (`source.go:56-64`).
- `--home`: uniform across place kinds — just changes the destination; no place-kind branch found.
- `--vendor`: uniform; resolved via `resolveVendor` (`main.go:322-331`) using `LoadState`, independent of place kind.
- `--check`: only meaningful for `update`; behaves differently by place — checkout: `Behind()==0` early-return before the dirty check (`main.go:466-476`, never reaches `shedKeptClone`); release: reported via `ErrNewerRelease` -> `exitCode(2)` (`main.go:212-214`, `release.go:127-130`); plain copy: `update` is refused before `--check` is even consulted (refusal happens inside `update()`'s `Behind` call, `main.go:462-464`, ahead of the `checkOnly` branches).
- `--dry-run`: on `install`, release binary prints a distinct dry-run line (`main.go:133-137`) vs. checkout/plain-copy's `[dry-run] the binary would be rebuilt...` (`main.go:156-157`) or `[dry-run] plain copy...` (`main.go:168-169`) — three different dry-run message shapes for the three place kinds on the same verb.

## Silent / untested cells ranked by real-machine likelihood

1. **HIGH** — `update` on a plain copy, exercised only through `main.go` dispatch: no CLI-level test drives `run(["update", ...])` against a plain-copy fixture to confirm the `errPlainCopy` message actually surfaces through `update()`'s wrapping (`main.go:462-464`, `fmt.Errorf` there is untested — only the lower-level `Behind`/`Pull` are). A real machine hits this whenever someone runs `sh-agent update` inside an unpacked tarball or `git archive` export.
2. **MEDIUM** — `status` on a plain copy: the "update is refused" status line (`main.go:351`) has no test; likely hit whenever a plain-copy install is inspected.
3. **MEDIUM** — `adopt` refusal message on a release binary (`main.go:250`) has no test pinning its exact wording — plan's Failure modes rule requires refusal messages to name cause+alternative; this one does not name an alternative (just states the fact), worth flagging to reviewer.
4. **LOW-MEDIUM** — `list`/`render`/`diff`/`status` against `EmbeddedSource` specifically (not just install) — likely fine since code path is shared, but no test names it directly for these verbs.
5. **LOW** — `sync` as a literal invoked verb string is untested by name (always tested as `install`), though it's the same `case` arm so risk is low.


## Exhaustive grep pass (`IsCheckout`/`src.Dir ==`/`IsReleaseVersion` across internal/harness/*.go, non-test)

Confirms the matrix above is complete for the traced dispatch paths, and surfaces two more place-kind branches not reachable from `main.go`'s verb dispatch directly but load-bearing for `report` and manifest state:

- `internal/harness/apply.go:108-114` (`Install`) — manifest's recorded `repo`/`commit`: embedded -> `EmbeddedPrefix+src.Version`, commit `""`; tree on disk (checkout or plain copy alike) -> `src.Dir`, commit from `Inspect` (empty for plain copy, "no commit of another repo is recorded"). This is what `main.go`'s `knownCheckout` (`main.go:546-570`) later reads back to refuse an install/update from an unmatched checkout.
- `internal/harness/report.go:345-357` (`buildRecord`-type function) — three-way branch for the report payload: `src.FS==nil` (no harness at all) -> notes null fields; `src.Dir==""` (release) -> `harness_commit` null, `harness_version` set; `default` (checkout or plain copy, same branch) -> `Inspect(src.Dir).Commit`, with a note "harness is a plain copy, not a git checkout: harness_commit is null" when `Commit==""`. Exercised, no refusal.
- `internal/harness/report.go:391` (`ResolveFeedback`, for `--feedback-repo`) — refuses a non-checkout: `"--feedback-repo %s is not a git checkout; it is never cloned into"`. This is about the *feedback* repo, not the harness source, so it's orthogonal to the verb x place-kind matrix but is another `IsCheckout` refusal site worth knowing about.
- `internal/harness/report.go:539` (kept-clone sweep, `RemoveKeptClone`) — refuses to touch anything that is `not a git checkout` before deleting it; same rationale, not part of the harness-source matrix.
- `internal/harness/inventory.go:91` (`checkManagedTree`) — `src.Dir == ""` short-circuits (embedded files hold no symlink to guard against); uniform for checkout and plain copy otherwise (both get the same symlink-guard treatment via `src.path`).

No additional silent cells found in this pass; every branch is either an exercised code path or a refusal with a message.

## Not examined (ran out of budget / out of declared scope for this pass)

- `internal/harness/apply.go`'s full per-file `switch` (add/update/mode/remove/keep lines) — Phase 1 already covered this per the plan; I did not re-verify it here.
- `internal/harness/report.go` in full (900 lines) beyond confirming the release-binary bypass at `main.go:298-303` and test names — did not trace every internal branch for place-kind sensitivity beyond what's in the table.
- `internal/harness/pathguard.go` beyond `NoteRunningElsewhere`'s call sites — did not check whether `PathShadows`/`RunningElsewhere` behave differently by place kind (unlikely, since they operate on PATH/binary identity, not source kind).
- Did not run any `go run ./cmd/sh-agent ...` commands against scratch homes to confirm behavior live (relied on static reading + existing test names); the plan's Phase 3 is where that integration matrix is meant to be built.
- Did not check `internal/harness/state.go`, `internal/harness/plan.go`, `internal/harness/derive.go`, `internal/harness/inventory.go`, `internal/harness/frontmatter.go`, or `internal/harness/render.go` for place-kind-specific branches — skimmed only `Source`/`repo.go`/`release.go`/`apply.go`(partial)/`adopt.go`(partial); a full pass would re-grep these for `IsCheckout`/`src.Dir ==`/`IsReleaseVersion` to be exhaustive.
- Did not confirm whether any test exercises `update --check` on a plain copy (likely untested given item 1 above covers the base case too).

## Open questions

- The exhaustive grep pass above (done) found no further silent cells, but I did not trace `internal/harness/state.go`, `plan.go`, `derive.go`, `frontmatter.go`, or `render.go` line-by-line — only grepped them for the three markers, none of which hit. If `/review --all` (Phase 2's next step) wants deeper coverage of those files, it should read them directly.

recon complete
