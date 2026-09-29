---
owner: verifier
date: 2026-09-29
---
# Verification — search-index-phase-2   (2026-09-29)

Range reviewed: bfcedcc..HEAD (HEAD = 7fd56d5)

## Status: PASS

Scope: gofmt/vet/race suite, repeat run of the six TestGlobal* tests, and the Phase 2 hand integration across two scratch git checkouts (index + `search --global`, `--global --open`, `--global --refs`). Complete.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race -v -count=1 ./...` → exit 0. Output at `.agents/scratch/verify-search-index-phase-2/go-test-1.txt`.
  counts: pass 204, fail 0, skip 2 at 7fd56d5 for `go test -race -v -count=1 ./...`
  vs previous: `.agents/docs/testing/2026-09-29-verify-search-index-phase-1.md`, pass 198→204, skip 2→2 (no rise, no fall in fail). The +6 pass matches the 6 new `TestGlobal*` subtests added by this range's only source file, `internal/search/global_test.go`.
  Skips: `TestRecall` (gated by `SH_AGENT_RECALL=1`, unset here — pre-existing, unrelated to this range), `TestHelperIndex` (a `-run TestHelperMain`-only helper-process test, prints "helper process" and skips when run directly — pre-existing, unrelated to this range).
  Note: an earlier accidental duplicate `go test -race` invocation (started via a shell `&` before switching to `run_in_background`) produced a stale `(cached)` result for one package; that run's output was discarded and the command rerun with `-count=1` exactly as above before any counts were recorded — see the process-management steps in this session, no test artifact from the bad run was kept or cited.
- `go test -race -count=3 -v ./internal/search -run TestGlobal` → exit 0. Output at `.agents/scratch/verify-search-index-phase-2/go-test-global-repeat.txt`.
  18 of 18 pass (6 subtests × 3 repeats), 0 fail: `TestGlobalRegistersProject`, `TestGlobalRunsOutsideACheckout`, `TestGlobalIndexesFeedbackRetrosProposals`, `TestGlobalDropsVanishedProject`, `TestGlobalNameCollision`, `TestGlobalRebuildDropsRegistry` — all stable across 3 repeats under `-race`.

### Phase 2 hand integration (all under `.agents/scratch/`, never `/tmp`, never the real home)
- Built the binary: `go build -o .agents/scratch/verify-p2-bin/sh-agent ./cmd/sh-agent`.
- Two fresh git repos: `.agents/scratch/verify-p2-repo-a` and `.agents/scratch/verify-p2-repo-b`, each `git init`, each with `.gitignore` = `.agents/scratch/*`, committed.
- Put `.agents/scratch/verify-p2-repo-b/.agents/docs/bugs/BUG-001-zabranthium.md` (frontmatter `id: BUG-001`, `status: open`, `date: 2026-09-29`, unique word `zabranthium`) — committed in repo-b only.
- **Before running the binary**, grepped both repos for `zabranthium`: found in repo-b (2 lines), `NOTFOUND-A` in repo-a — fixture proven before any tool ran.
- `sh-agent index --home .agents/scratch/verify-p2-home` from inside repo-b → exit 0: `files 2, chunks 3, records 1, refs 3, changed 2, bytes 98384, ms 118`.
- `sh-agent search --global zabranthium --home .agents/scratch/verify-p2-home` from inside repo-a → exit 0:
  ```
  search: 1 hits for "zabranthium" (index current: 1 files, reconciled 2 changed in 148 ms)
  verify-p2-repo-b:.agents/docs/bugs/BUG-001-zabranthium.md:6  bug BUG-001 open  2026-09-29
      # BUG-001: the »zabranthium« bug Unique marker word: »zabranthium«
  ```
  Names the second repository (`verify-p2-repo-b`) and `bug BUG-001 open` — matches the required assertion exactly.
- `sh-agent search --global --open --home .agents/scratch/verify-p2-home` from repo-a → exit 0:
  ```
  search: 1 open records (index current: 1 files, reconciled 0 changed in 14 ms)
  verify-p2-repo-b:.agents/docs/bugs/BUG-001-zabranthium.md:1  bug BUG-001 open 2026-09-29 BUG-001: the zabranthium bug
  ```
- `sh-agent search --global --refs BUG-001 --home .agents/scratch/verify-p2-home` from repo-a → exit 0:
  ```
  BUG-001  verify-p2-repo-b:.agents/docs/bugs/BUG-001-zabranthium.md:1  open  2026-09-29  BUG-001: the zabranthium bug
  referenced by 1:
  verify-p2-repo-b:.agents/docs/bugs/BUG-001-zabranthium.md:2  bug BUG-001 open  2026-09-29  ×2
      id: »BUG-001«
  ```
  Both print correctly, both name the correct repo and record.

## Root cause
N/A — no failures found.

## Coverage of changed code
- `internal/search/global_test.go` (+152 lines, the range's only changed file) — this file IS the test; it exercises `internal/search`'s global-index code paths (registration, `~config` sources, run-outside-a-checkout, vanished-project drop, name collision, rebuild-drops-registry) via the `TestGlobal*` family, all present and passing (see counts above, run 3× for stability).
- No non-test source file changed in this range (`internal/search/global_test.go` plus two docs files: `.agents/docs/plans/2026-09-29-search-index.md`, `.agents/docs/specs/sh-agent/search.md`); no coverage gap to report.
- Suppressions added in range: none. `git diff bfcedcc..HEAD -- internal/search/global_test.go` grepped for `nolint`, `lint:ignore`, `ts-ignore`, `ts-expect-error`, `eslint-disable`, `t.Skip`, `SkipNow`, `go:build ignore`, `.skip(`, `.only(` — no matches.
- Phase 2 plan items cross-checked against the test names: registration → `TestGlobalRegistersProject`; `~config` sources / feedback-retros-proposals → `TestGlobalIndexesFeedbackRetrosProposals`; vanished project dropped → `TestGlobalDropsVanishedProject`; name collision → `TestGlobalNameCollision`; `--global` outside a checkout → `TestGlobalRunsOutsideACheckout`; rebuild dropping the registry → `TestGlobalRebuildDropsRegistry`. All six named behaviours have a passing test; matches the plan's own "done 2026-09-29" note at line 222 of the plan.

## Not examined
- Mutations G1–G5 named in the plan (line 222, "each compiled and failed a named test") were not re-run by this verifier; the plan's checkpoint note claims they were run when the code landed, but this session did not re-verify that mutation claim (no mutation harness was re-executed here).
- `govulncheck ./...` (listed in project `AGENTS.md` § Setup commands) was not run in this session — out of the scope of Phase 2's own plan checkpoint (line 224: "full test command exit 0, verifier log PASS"), not requested in the dispatch.
- `sh-agent status --home <scratch>` drift check after install was not run (the dispatch asked for `index`/`search`, not `install`/`status`).
- Phase 3–5 items (context/agent behavioural runs, Antigravity runs, real-home smoke, recall test, efficiency runs) are out of scope for this Phase 2 verification and were not touched.
- The `TestRecall` skip (gated by `SH_AGENT_RECALL=1`) was not exercised — Phase 4 work, not this range.
- No `TestHarnessMap`-style cross-check of `search --global` documentation vs behaviour beyond reading the printed output above; did not separately open `specs/sh-agent/search.md` to diff prose against behaviour (only confirmed it was touched in this range).

## Environment
Linux (host), go1.27.1 linux/amd64, repo at /home/kesokaj/workspace/styrman-harness, working tree clean and frozen throughout the run (git status --short empty before/after; only `.agents/scratch/` scratch dirs added, gitignored).
