---
owner: verifier
date: 2026-10-09
---
# Verification — project   (2026-10-09)

Range reviewed: the whole repository at HEAD (c267990), no diff.

## Status: PASS

Scope: AGENTS.md § Setup commands ran in full except govulncheck (not on PATH); guide .agents/docs/testing/guide.md § 1 and § 3e corpus measurement ran; § 3d/3g real-data checks, § 2/§ 4 sandbox starts not run (project audit, as dispatched). Not `complete`. Coverage gaps (below) are listed separately and do not change the command verdict.

## Commands
- 1. `pwd`, `git status --short`: clean at start; HEAD c267990; go1.27.2 linux/amd64; 16 CPUs.
- 4. `gofmt -l cmd internal` -> exit 0, printed nothing.
- 4. `go vet ./...` -> exit 0, printed nothing.
- 4. `go test -race -count=1 -timeout 30m -json -coverprofile=... ./...` -> exit 0 (real 12m16s, run in one go); all 8 packages ok (cmd/sh-agent 731s, internal/harness 577s, search 128s, dream 49s, metrics 35s, memory 22s, trail 18s, root 2s). counts: pass 649, fail 0, skip 3 at c267990 for go test -race ./... (counted from -json, subtests included; -v was via -json, not a text file); vs previous: none compared (no earlier log read). Statement coverage total 87.5%. Output: .agents/scratch/verify-project/test.json, cover.out, func.txt.
- 4. `go run ./cmd/sh-agent install --home .agents/scratch/home` -> exit 0; "1790 written, 0 unchanged, 0 removed, 0 kept"; settings attribution note printed.
- 4. `go run ./cmd/sh-agent list` -> exit 0, 63 lines.
- 4. `govulncheck ./...` -> NOT RUN: not on PATH, not in ~/go/bin or ~/.local/bin. Reported, not installed.
- 8. `git status --short --untracked-files=all` after: only .agents/docs paths other agents wrote (plans/2026-10-09-recon-project.md, reviews/2026-10-09-review-project-{code-docs,prompts,simplify}.md, reviews/2026-10-09-security-project.md) plus this log; no source/test/config path changed. Tree did not move.
- 3e. `SH_AGENT_CORPUS=<copy of ~/.gemini/antigravity-cli/conversations/*.db in scratch> go test -count=1 -v -run 'TestAgyCorpusLeaves$' ./internal/search` -> exit 0: "23 databases, 197 prose leaves kept, 0 with a tool-argument key".

## The three skips (only these fired of ~50 t.Skip sites)
- internal/search/dream_test.go:397 TestAgyCorpusLeaves: measurement, needs SH_AGENT_CORPUS. Reason true by design (it has no assertion, see below).
- internal/search/recall_test.go:250 TestRecall: measurement, needs SH_AGENT_RECALL=1.
- internal/search/search_test.go:529 TestHelperIndex: helper process, not a test.
Other Skip reasons (root umask, git/go/bash not on PATH, symlinks) did not fire here; box has go, git, sh, agy. The root/mode-000 skips are true only when run as root; this run was not root.

## Verb/flag end-to-end map (cmd/sh-agent tests drive run(); `main` itself is 0%)
- install/sync, status, diff, update (+ --check, --dry-run, --from-zip), uninstall, list, render, adopt (+ --force), report (+ --feedback-repo, --dry-run), search (--kind, --limit, --global, --refs, --open), index (--rebuild), memory add/show, dream (--all/--project, --install-timer, --remove-timer, --dry-run), metrics ingest/show/report/export (--source, --label, --all-labels, --json), --repo/--home/--vendor: all have test hits in main_test.go, zip_test.go, search_test.go, memory_test.go, metrics_test.go (e.g. TestSyncIsInstall, TestDiffVerbShowsDrift, TestDiffAfterZip, TestDreamVerb, TestDreamTimerDryRun, TestAdoptOn*, TestReport*). Mapping is by grep of verb/flag strings in tests plus per-function coverage of run*; I did not read each test to prove it asserts the outcome.
- GAP cmd/sh-agent/main.go:309 `runVersion` 0.0% -> verb `sh-agent version` is not driven through run(); only `buildVersion`/TestVersionFromInfo (main_test.go:126) covers the version string logic. A change that breaks the `version` dispatch or its output passes.
- GAP cmd/sh-agent/main.go:84 `main` 0.0% (exit-code wiring, `Error` at :98 0%): no test runs the real binary's exit status for a failing verb.
- Weaker: runRender 75%, runDiff 80%, runList 80%, run 82%, runDream 84% (error legs unexercised).

## Coverage of changed code (per source file, statement %, own-package tests, -race run)
- cmd/sh-agent/main.go 87% - main_test.go, zip_test.go, search_test.go, memory_test.go, metrics_test.go
- internal/dream: collect.go 93, dream.go 93, metrics.go 83, prompt.go 100, run.go 85, state.go 83, timer.go 85 - collect_test.go, run_test.go, metrics_test.go, timer_test.go, main_test.go
- internal/harness: adopt.go 76, apply.go 84, derive.go 100, frontmatter.go 99, inventory.go 93, owned_unix.go 100, pathguard.go 92, plan.go 94, release.go 91, render.go 92, repo.go 93, report.go 88, settings.go 81, source.go 100, state.go 85, sweep.go 79, zip.go 85; owned_other.go not compiled on linux (build-tagged, not measured)
- internal/memory/memory.go 91 - memory_test.go
- internal/metrics: agy 87, claude 86, events 86, ingest 82, numbers 93, open_unix 100, report 92, show 89, store 75, verb 94; open_other.go not compiled here - *_test.go per file
- internal/search: dream.go 82, export.go 66, parse.go 98, reconcile.go 85, search.go 84, store.go 74 - dream_test, search_test, hardening_test, security_test, global_test, protofields_test
- internal/trail/trail.go 87 - trail_test.go
- embed.go (root): no statements; embed_test.go, plan_test.go, init_project_test.go cover its embedded content.
No source file has 0% coverage. Note: the profile is per-package (no -coverpkg), so a function exercised only by another package's tests shows 0% below though it may be covered there; I did not trace cross-package callers.

### Functions at 0% in their own package (potential gaps, file:line)
- cmd/sh-agent/main.go:309 runVersion (real gap, see above); :84 main; :98 Error; :1089 shedKeptClone 50%
- internal/harness/plan.go:208 SecretShaped; internal/harness/report.go:289 RedactJSONL (redaction helpers; if only used cross-package, covered incidentally at best; unverified)
- internal/dream/timer.go:276 OtherSystemTimerFiles, :296 TimerInstalledHere (likely driven from cmd TestTimerOfAnotherSystemIsNamed, memory_test.go:219; unverified); dream.go:191 Error
- internal/metrics/ingest.go:34 Line; report.go:34 ReportPath; claude.go:402 resultText 42.9%
- internal/search/export.go:15-38 PrepareDBFiles, DSN, IsBusy, OpenCopy, AgyTime, Printable, KindOf: 0% in search; exports for dream/metrics (incidental unless metrics/dream tests assert them)
- internal/search: reconcile.go:426 apply 58%, :608 newInserter 55%; search.go:378 dropProject 50%; store.go:177 Error 0%, :202 ensureSchema 57%
- internal/trail/trail.go:193 LazyFetchSafe 0%

## Tests that cannot fail
- internal/search/dream_test.go:394 TestAgyCorpusLeaves: no assertion, only t.Logf (by design, documented "measurement, not a gate"). Any source change to proseLeaves passes it. Counted as passing/skipped in counts but guards nothing.
- internal/search/recall_test.go:~246 TestRecall: no assertion by static scan, skipped without SH_AGENT_RECALL; same property.
- TestHelperIndex: a helper process, not a test.
- A heuristic scan of all 77+ test functions in cmd and all packages for bodies with no t.Error/Fatal/Run/assert-like helper found only the two above. Not found by this scan (not checked): self-comparisons, fixtures that already hold the expected output, helper-hidden missing assertions.

## Slowest tests (go test -race, elapsed)
1. internal/harness TestRefreshCatalogScript 309s (runs the refresh script, nested go test/git; skip-gated in catalog_test.go:709-722)
2. cmd/sh-agent TestAttributionSet 185s
3. cmd/sh-agent TestAttributionNote 147s (subtest no_settings_file alone 20s)
4. internal/harness TestKitPreambleRefusesLeftoverCredentials 96s
5. internal/harness TestKitPreambleAssertsCommit 48s
6. internal/harness TestKitPreambleKillsLeftoverProcesses 32s
7. cmd/sh-agent TestRunInstallStatusUninstall 29s
8. cmd/sh-agent TestStatusWarnsOnOversizedSettings 28s
9. internal/harness TestReinstallBuildsAndInstalls 28s
10. cmd/sh-agent TestPlainCopyNeverRebuilds 20s
Not checked: whether each of these builds the binary or installs into a scratch home more than once (I did not read their bodies); the names (Attribution*, Reinstall*, PlainCopy*) suggest `go build` per case. The two Attribution tests (332s combined, about 45% of cmd/sh-agent's 731s) are the first thing to read. Durations are inflated by -race and parallel package runs on 16 CPUs (this box, not the one-CPU box AGENTS.md cites; 12m16s wall here anyway).

## TestAgyCorpusLeaves measurement
Guide/AGENTS.md figure: 971 prose leaves, 0 tool-argument keys (2026-09-30). Today on this box: 23 conversation databases (copied to scratch), 197 leaves, 0 with a key. The "0 keys" invariant holds; the 971 does NOT reproduce, because the corpus differs (not the 2026-09-30 set; the reader code was not diffed). 971 is not re-derivable here.

## Not examined
- govulncheck (not installed here): new advisories unchecked.
- Guide § 3d/3g real-data checks, § 2/§ 4 sandbox starts, behavioural agent runs, § 5: not run.
- Per-test body read for behavioural vs incidental coverage: not done for individual tests; judgement above is from per-function coverage and grep, so "incidental-only" files cannot be named beyond the 0%-function list.
- Whether slow tests rebuild repeatedly; cross-package coverage (-coverpkg) of export.go, SecretShaped, RedactJSONL, Timer helpers.
- No comparison with an earlier log's counts.

## Environment
Linux 6.18.28, go1.27.2 linux/amd64, 16 CPUs, agy at /usr/local/bin/agy, no Docker needed. Scratch: .agents/scratch/verify-project/, .agents/scratch/home (gitignored).
