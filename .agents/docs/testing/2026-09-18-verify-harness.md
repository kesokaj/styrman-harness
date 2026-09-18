---
owner: verifier
date: 2026-09-18
---
# Verification — harness (internal/harness, cmd/sh-agent)   (2026-09-18)

Range reviewed: whole-file scope, working tree — every non-test `.go` file in
`internal/harness/` and `cmd/sh-agent/` (not a diff)

## Status: PASS (with two coverage gaps, ranked below)

Scope: `gofmt`/`vet`/`go test -race` + coverage run over both packages; every
exported function and every `main.go` verb-dispatch case mapped to its test(s)
by name. Guide `.agents/docs/testing/guide.md` cited (tiers 2-5, harness
self-check/install/report/sandbox) were NOT run — out of scope per dispatch,
which asked only for the Go-check tier and coverage mapping.

## Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0. counts: pass 112, fail 0, skip 0
  at `d737282`+dirty (untracked review/testing files only, no source changes)
  for `go test -race -v -count=1 ./...`.
  Packages: `github.com/kesokaj/styrman-harness` ok, `.../cmd/sh-agent` ok
  (4.228s), `.../internal/harness` ok (2.816s).
  vs previous: no earlier `testing/*.md` log in this repo carries a
  `counts:` line for this exact command (checked `.agents/docs/testing/`);
  vs previous: none.
- `go test -race -count=1 -cover ./...` → coverage: root pkg "no statements";
  `cmd/sh-agent` 77.1%; `internal/harness` 86.9%.
- `go test -count=1 -coverprofile=.agents/scratch/verify-harness/cover.out ./...`
  then `go tool cover -func` → `.agents/scratch/verify-harness/cover-func.txt`.

## Root cause
N/A — no failures. Two functions genuinely uncovered (0.0%), one verb alias
untested; see below.

## Exported-function / verb-dispatch map

cmd/sh-agent/main.go
- `main` — none (0.0%, entry point; not itself unit-tested, expected — its
  body is `run` + `os.Exit`, and `run` is exercised extensively)
- `(exitCode) Error` — none (0.0%; the `error` interface method on the exit-code
  type, never called because tests read `err.Error()` via the wrapped errors
  `run` returns, not via an `exitCode` value)
- verb `version` — internal/harness/harness_test.go:TestReinstallBuildsAndInstalls
  (execs the installed binary with `version`)
- verb `uninstall` — cmd/sh-agent/main_test.go (line 161, unnamed at call site;
  containing test uses `run([]string{"uninstall", ...})`)
- verb `report` — cmd/sh-agent/main_test.go (lines 1031, 1166, 1304, 1366)
- verb `install` — cmd/sh-agent/main_test.go, many tests (e.g. line 137, 305,
  349, 1208, 1221, 1238, 1243, 1396, 1451, 1478)
- verb `sync` (alias of `install`) — **none**. Grepped every `*_test.go` in
  the repo for `"sync"`: zero hits. The alias shares `install`'s code path
  (`case "install", "sync":`) so it is behaviorally covered by every install
  test, but the dispatch string `"sync"` itself is never invoked by any test.
- verb `status` — cmd/sh-agent/main_test.go, many (147, 172, 353, 662, 1258, ...)
- verb `diff` — cmd/sh-agent/main_test.go:1276
- verb `update` — cmd/sh-agent/main_test.go (1127, 1136, 1297, 1329, 1532-1533)
- verb `list` — cmd/sh-agent/main_test.go:1280
- verb `render` — cmd/sh-agent/main_test.go:1284
- verb `adopt` — cmd/sh-agent/main_test.go:1288; internal/harness/harness_test.go
- helper `resolveVendor` — 100% (harness_test.go / main_test.go vendor-flag cases)
- helper `shedKeptClone` — 50% covered (main_test.go; one branch untested,
  not further isolated in this pass — see Not examined)
- helper `samePath`, `shouldRebuild`, `binaryOlder`, `reorder`,
  `abandonTemp`, `attributionHint`, `attributionOff`, `releaseFromInfo`,
  `versionFromInfo` — all 100%, exercised via main_test.go install/update/report
  paths and `TestReorder` (guide's own cited test for the `takesValue`/reorder
  path).

internal/harness/adopt.go
- `Adopt` — harness_test.go (TestAdopt*)

internal/harness/apply.go
- `Install` — harness_test.go, pathguard_test.go, report_test.go, source_test.go
- `(Report) Drift` — harness_test.go
- `Status` — harness_test.go
- `Diff` — **none (0.0%)**. Zero references to `Diff(` in any `_test.go` in
  the repo; only exercised indirectly at the CLI layer via
  `cmd/sh-agent/main_test.go:1276` (`run([]string{"diff", ...})`), which
  covers the `diff` verb's plumbing but not the `harness.Diff` function's own
  branches as a unit — `go tool cover -func` still reports it 0%, meaning the
  CLI path taken by that one `diff` test does not reach `Diff`'s statements
  (worth a look: possibly `diff` verb calls something else, or that test's
  home has no drift so `Diff` short-circuits before its body). **FAIL-worthy
  coverage gap** — flagged, not patched.
- `Uninstall` — harness_test.go
- `pruneEmptyDirs` — 90% (harness_test.go)

internal/harness/derive.go — `DeriveAgy` — harness_test.go

internal/harness/frontmatter.go — `Split`, `Parse`, `Lint`, `SplitList`,
`isQuoted`, `isBlockIndicator`, `yamlQuote` — all 100%, harness_test.go
(and `Split` also referenced from repo-root `embed_test.go`, release_test.go,
report_test.go)

internal/harness/inventory.go
- `identLike`, `validName` — 100%, harness_test.go
- `checkManagedTree` — 90%, harness_test.go
- `LoadInventory`, `LoadSkills`, `LoadAgents` — harness_test.go (LoadSkills/
  LoadAgents also pathguard_test.go, source_test.go)
- `LoadDocs` — 92.9%, no direct `LoadDocs(` call in any test; covered
  indirectly as `LoadInventory`'s callee (inventory.go:138,141) exercised by
  harness_test.go's LoadInventory tests.
- `RenderInventory`, `RenderInventoryCompact`, `slashNote` — 100%,
  harness_test.go

internal/harness/pathguard.go
- `PathShadows`, `RunningElsewhere` — pathguard_test.go
- `executableAt` — 100%, pathguard_test.go
- `NoteRunningElsewhere` — 83.3%, no direct call in any test; covered
  indirectly via release.go:178 (inside `UpdateFromRelease`, tested in
  release_test.go) and cmd/sh-agent/main.go:509.

internal/harness/plan.go
- `NormalizeVendor` — harness_test.go
- `Plan` — harness_test.go, pathguard_test.go, report_test.go, source_test.go,
  repo-root plan_test.go
- `IsOurs`, `secretShaped` — 100%, harness_test.go

internal/harness/release.go
- `AssetName`, `UpdateFromRelease` — release_test.go
- `LatestRelease` — 100%, no direct call in tests; covered indirectly as
  `UpdateFromRelease`'s callee (release.go:94), exercised by release_test.go.
- `checksumFor`, `globOrNil`, `compareRelease` — 100%, release_test.go
- `firstLineOfError` — 80%, release_test.go
- `firstStderrLine` — 100%, release_test.go

internal/harness/render.go — `RenderContext` — harness_test.go

internal/harness/repo.go
- `ResolveSource`, `IsRepo`, `Pull` — harness_test.go
- `IsCheckout`, `Behind`, `Reinstall` — cmd/sh-agent/main_test.go and
  harness_test.go
- `Inspect` — cmd/sh-agent/main_test.go, harness_test.go
- `gitOut`, `gitOutC` — 100%, harness_test.go

internal/harness/report.go
- `Collect`, `ResolveFeedback`, `RemoveKeptClone`, `SendReport` — report_test.go
- `BuildRecord` — 90%, report_test.go, source_test.go
- `labelAt`, `needleIn`, `blockEnd`, `marker`, `redactLines`, `decoded`,
  `redactDocument`, `throwawayDir`, `scrubCredentials` — 100%, report_test.go
- `verifyFeedback` — 80%, report_test.go

internal/harness/source.go
- `DirSource` — harness_test.go, pathguard_test.go, report_test.go,
  source_test.go, repo-root plan_test.go
- `EmbeddedSource` — 100%, no direct call in tests; covered indirectly via
  repo.go:47 inside `ResolveSource` (tested in harness_test.go) and used
  directly by cmd/sh-agent/main.go:300.
- `IsReleaseVersion` — 100%, same indirect path (repo.go:34, main.go:299/916).
- `path`, `readRegular` — 100%, source_test.go

internal/harness/state.go
- `LoadState`, `SaveState` — cmd/sh-agent/main_test.go, harness_test.go,
  pathguard_test.go, source_test.go
- `statePath` — 100%, harness_test.go
- `RemoveState` — 80%, no direct call in tests; covered indirectly via
  apply.go:273 inside `Uninstall` (tested in harness_test.go).
- `LastRetro` — harness_test.go

## Coverage of changed code (0.0% functions, cross-checked by name)

`go tool cover -func` 0.0% rows, confirmed with `grep -P '\t0\.0%$'`:
1. `cmd/sh-agent/main.go:59 main` — entry point; not independently testable,
   not flagged as a real gap.
2. `cmd/sh-agent/main.go:73 (exitCode) Error` — trivial `Sprintf`, unreached
   because no test constructs a bare `exitCode` and calls `.Error()` on it
   directly (they inspect the wrapped error from `run` instead). Low risk,
   ranked lowest.
3. `internal/harness/apply.go:206 Diff` — **real gap, ranked highest**. `Diff`
   is the function behind the `diff` verb, which every place running
   `sh-agent diff` (an operator debugging drift) hits directly. It has a CLI
   test (`main_test.go:1276`) but that test does not reach `Diff`'s own
   statements per the coverage tool, meaning drift-formatting output is
   unverified. This ranks above the other two because `diff` is a documented,
   user-facing verb (guide.md doesn't cover it either — tier 3 only checks
   `status`'s "drift none").

Ranked by how likely a real machine hits it:
1. `Diff` (apply.go:206) — hit on every operator `sh-agent diff` call.
2. verb-string `"sync"` never invoked by name in tests — alias path shares
   `install`'s tested logic, so behavior is covered even though the literal
   dispatch string is not; medium-low likelihood a real place ever types
   `sync` instead of `install` (no doc references `sync` as a public verb).
3. `(exitCode) Error` — internal plumbing, extremely unlikely to matter.
4. `main` — never independently hit; expected.

## Not examined
- Testing-guide tiers 2-5 (harness self-check install/status against a
  scratch home, plain-copy install, `report --dry-run`, sandbox CLI starts,
  real-home smoke checks) — out of scope for this dispatch, which asked only
  for the Go-check tier + coverage mapping over the named files.
- Why `Diff`'s statements read 0% despite `main_test.go:1276` calling the
  `diff` verb — not root-caused beyond "the CLI test's fixture likely has no
  drift, so `Diff` never executes its formatting body"; not confirmed by
  tracing `cmd/sh-agent/main.go`'s `diff` case to `harness.Diff`'s call site.
- `shedKeptClone`'s untested branch (50% coverage) — not isolated to a
  specific `if` inside the function in this pass.
- Root package `github.com/kesokaj/styrman-harness` (repo-root `embed_test.go`,
  `plan_test.go`) — outside the named scope (`internal/harness`, `cmd/sh-agent`
  only), noted only because it appeared in `go test ./...` output.

## Environment
Linux 6.12.107+deb13-amd64, go1.27.1 linux/amd64 (user-level SDK per this
repo's `AGENTS.md`). No Docker/Colima involved — no Dockerfile in scope.
