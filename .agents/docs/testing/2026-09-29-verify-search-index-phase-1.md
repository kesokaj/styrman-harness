---
owner: verifier
date: 2026-09-29
---
# Verification — search-index-phase-1   (2026-09-29)

Range reviewed: c04f2d7..7c24503 (HEAD)

## Status: PASS

Scope: complete — AGENTS.md § Setup commands, install/status pair, three release cross-builds, go.mod pin check, govulncheck, and the plan's Phase 1 behavioural checks (index --rebuild timing, warm search timing repeated, search --refs ADR-040/041, scratch cleanliness), plus coverage-of-changed-code mapping and suppression scan. `testing/guide.md` exists (not modified; no new procedure needed).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0. Output saved at `.agents/scratch/verify-search-index-phase-1/go-test-1.txt`.
  counts: pass 198, fail 0, skip 2 at 7c24503 for `go test -race -v ./...`
  Skips (both carry a reason, confirmed by grepping the surrounding lines): `TestRecall` — "measurement, not a gate: set SH_AGENT_RECALL=1 to run it" (internal/search/recall_test.go); `TestHelperIndex` — "helper process" (internal/search/search_test.go, a subprocess-helper harness, not a real test).
  vs previous: none — no earlier log in `testing/` has a counts line for the exact command `go test -race -v ./...` over the whole module.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-p1` → exit 0. 266 written, 0 unchanged, 0 removed, 0 kept. Output at `.agents/scratch/verify-search-index-phase-1/install.txt`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-p1` → `binary     7c24503ea3a1+dirty` (no warning line); `drift      none — 265 files match the repo`. The repo line's `+dirty` / "uncommitted changes" reflects this run's own log file (`.agents/docs/testing/2026-09-29-verify-search-index-phase-1.md`), which was untracked at the time of that status call — confirmed via `git status --short --untracked-files=all` showing only that one path both before and after. No code drift found.
- Cross-builds, `CGO_ENABLED=0 GOTOOLCHAIN=local go build ./cmd/sh-agent`, output under `.agents/scratch/verify-search-index-phase-1/`:
  - `GOOS=linux GOARCH=amd64` → exit 0, binary 13714185 bytes.
  - `GOOS=darwin GOARCH=arm64` → exit 0, binary 13423218 bytes.
  - `GOOS=linux GOARCH=arm64` → exit 0, binary 13092206 bytes.
- `grep '^go ' go.mod` → `go 1.24.4` (unchanged, per ADR-045's pin promise).
- `govulncheck ./...` → installed at `~/.go/bin/govulncheck`, go1.27.1 on PATH. Exit 0: "No vulnerabilities found... 0 vulnerabilities in packages you import and 1 vulnerability in modules you require, but your code doesn't appear to call [it]" — consistent with AGENTS.md's note that an older-than-go1.25.8 stdlib carries two known advisories; here on go1.27.1 the direct/imported count is 0 and only an unreachable required-module advisory remains, no NEW advisory. Output at `.agents/scratch/verify-search-index-phase-1/govulncheck.txt`.
- Plan Phase 1 behavioural checks, binary built at `.agents/scratch/verify-search-index-phase-1/sh-agent` (native build), run from the repo root (the tool's per-project index lives at `<repo>/.agents/scratch/search.db`; `--home` is not consumed by `index`/`search` for a per-project run — it only matters for `--global`, so the dispatch's `--home .agents/scratch/verify-p1-home` was not applicable here and was not used for these two commands):
  - `index --rebuild` → exit 0, `files 1054, chunks 10377, records 915, refs 39773, changed 1054, ms 2108`. Timed with shell `time`: real 0m2.135s — under the 10 s bound.
  - Warm `search "delivery model" --limit 3`, timed 5 of 5 runs with `date +%s%N` deltas (wall clock, end to end including process startup): 32 ms, 32 ms, 30 ms, 25 ms, 28 ms — all 5 of 5 under the 100 ms bound.
  - `search --refs ADR-040` → `ADR-040 ... Superseded by ADR-041 ... 2026-09-28`, and the reference list includes an ADR-041 entry ("adr ADR-041 Accepted") whose body quotes "supersedes «ADR-040»" — matches the plan's expected output shape (D23 grouping). Full output saved at `.agents/scratch/verify-search-index-phase-1/refs-adr040.txt`.
  - `git status --short --untracked-files=all` after all of the above → only `.agents/docs/testing/2026-09-29-verify-search-index-phase-1.md` (this log, untracked as expected of a verifier run); no `.agents/scratch/*` path appears (confirmed gitignored — `search.db`, `search.db-wal`, `search.db-shm` and all cross-build binaries stayed untracked).
- `go test -cover ./internal/search/... ./cmd/sh-agent/... ./internal/harness/...` → all three `ok`: internal/search 72.5%, cmd/sh-agent 83.3%, internal/harness 87.2% (informational, not a gate line by itself — see Coverage of changed code below).

## Root cause
None — no failing command, no FAIL condition found.

## Coverage of changed code
Non-test Go files changed in the range (`git diff --stat c04f2d7..HEAD -- '*.go'`):
- `internal/search/parse.go` (new, 281 lines) — covered by `internal/search/search_test.go` (chunking/heading tests, e.g. `TestChunkMarkdownHeadingsAndWindows`) and `internal/search/recall_test.go`.
- `internal/search/reconcile.go` (new, 415 lines) — covered by `internal/search/search_test.go` (`TestConcurrentReconcile`, busy/refusal tests, `TestRefuseOutsideCheckout`, `TestRefuseNotIgnored`, `TestRefuseSymlinkedScratch`, `TestBusyRefusesNotStale`) and the rebuild/reconcile path exercised live above (`index --rebuild`).
- `internal/search/search.go` (new, 777 lines) — covered by `internal/search/search_test.go` (`TestSearchRanksAndCaps`, `TestSearchOperatorCharactersAreLiteral`, `TestRefs`, `TestRefsIDsAndPaths`, `TestOpen`, `TestArgumentErrors`, `TestRecordsBugADRPlan`) and the live `search`/`search --refs` runs above.
- `internal/search/store.go` (new, 210 lines) — covered by the same `internal/search/search_test.go` suite (schema/store exercised through every test above; store.go has no standalone `store_test.go` but every test in the package opens a store).
- `cmd/sh-agent/main.go` (+30 lines, new `search`/`index` subcommands) — covered by `cmd/sh-agent/search_test.go` (new, 112 lines).
- `internal/harness/report.go` (+7 lines, skips `search.db`) — covered by `internal/harness/report_test.go` (+26 lines; plan names `TestCollectNeverSendsSearchDB`, confirmed present and passing in the `go test -race -v` run).
- `internal/harness/state.go` (comment-only change, `ConfigDir` doc comment names `search.db`) — no behavior change; not a FAIL, no test needed.

No changed non-test Go function found without a covering test.

Suppression scan over the diff (`git diff c04f2d7..HEAD | grep -nE 'nolint|lint:ignore|ts-ignore|ts-expect-error|eslint-disable|t\.Skip|SkipNow|go:build ignore|\.skip\(|\.only\('`): two added `t.Skip` calls, both in `internal/search`, both carry an inline reason string ("measurement, not a gate: set SH_AGENT_RECALL=1 to run it" for `TestRecall`; "helper process" for `TestHelperIndex`, a subprocess-mode helper, not a disabled real test). No added `.only(`, no unreasoned suppression.

## Not examined
- Global index behaviour (`--global`, cross-project registry) — plan line 223 marks this an open (unticked) Phase-2-adjacent item, not part of Phase 1's done checklist; not exercised here.
- Context entry / GANNET render check, consumer D13 lines, and behavioural CLI runs (`claude -p --append-system-prompt`) — these are separate unticked plan items (lines 227-230), outside the scope this dispatch named (Phase 1's install/status/build/index/search checks); not run.
- Vector-fusion gate probe and the larger-project eff-kit task (plan lines 238-244) — explicitly future/gated phases, not run.
- `govulncheck -show verbose` detail on the one required-but-unreached module advisory — not requested by the dispatch; only the summary was captured.
- Full `-race` detector stress beyond the single `go test -race -v ./...` run (e.g. repeated runs to shake out flaky concurrency in `TestConcurrentReconcile`) — ran once, passed once; not repeated.

## Environment
Linux 6.12.107+deb13-amd64. `go version go1.27.1 linux/amd64` (ahead of `/usr/bin/go`, per `.agents/docs/deploy/local.md` § Operator steps). `govulncheck` present at `~/.go/bin/govulncheck`. No Docker/Colima involved (no Dockerfile changed in this range).
