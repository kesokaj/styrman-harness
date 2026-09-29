---
owner: verifier
date: 2026-09-29
---
# Verification — search-index-final   (2026-09-29)

Range reviewed: c04f2d7..HEAD (HEAD=aae16af)

## Status: PASS

Scope: complete — ran AGENTS.md § Setup Commands and testing/guide.md sections 1, 2, 3 and 3d over the range: gofmt, vet, race test suite, go.mod version check, govulncheck, three release cross-builds, a `/usr/bin/go` build, the context size test, the install/status pair with context grep, index rebuild/warm-search timing, TestRecall, and coverage mapping of every changed non-test .go file. `testing/guide.md` exists (cited above, no changes needed to it).

## Anomaly noted (not part of this verification's scope)
While this run was in progress, an untracked file appeared in the working tree: `.agents/docs/reviews/2026-09-29-review-search-index-code-r2.md` (created 2026-09-29 20:04:21, mid-run — a placeholder `Status: IN PROGRESS` reviewer artifact, apparently a concurrent `reviewer` dispatch over `565bb1e..aae16af` writing its own file). This is another agent's own trail artifact, not source/test/config, and it does not touch any file this verification exercises, so it did not affect the results above. It is flagged because the dispatch stated "the working tree is clean and stays frozen while you run," and a file appearing mid-run — even a harmless one — is the condition that rule exists to catch; the caller should confirm whether a reviewer was dispatched in parallel.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race -v -count=1 ./...` → all packages `ok` (styrman-harness, cmd/sh-agent, internal/harness, internal/search), no FAIL line. Output saved to `.agents/scratch/verify-search-index-final/go-test-1.txt`.
  counts: pass 215, fail 0, skip 2 at aae16af for `go test -race -v -count=1 ./...`
  vs previous: `.agents/docs/testing/2026-09-29-verify-search-index-phase-2.md`, pass 204→215, skip 2→2 (no rise in skip, no fall in pass). +11 pass matches the round-2/round-3 security tests added in this range (`TestSkipsSecretShapedFiles`, `TestGlobalDirIsPrivate`, `TestGlobalDropsCheckoutWithoutAgentsDocs`, `TestHintAndNotesArePrintable`, `TestSkipsMoreSecretNames`, `TestRegistrationKeptWhenCheckErrors`, `TestRebuiltNoteIsPrintable`, plus others in security_test.go). The 2 skips are unchanged: `TestRecall` (gated by `SH_AGENT_RECALL=1`, run separately below, carries a reason) and `TestHelperIndex` (a subprocess-helper pattern, carries a reason "helper process").
- `grep '^go ' go.mod` → `go 1.24.4` (matches required, no toolchain directive per ADR-029)
- `govulncheck ./...` and `govulncheck -show verbose ./...` → "No vulnerabilities found" for code; 1 vulnerability in a required module not called by our code, identified via `-show verbose` as `GO-2026-5024` — matches the expected, known advisory named in the dispatch. No other advisories present.
- Cross-builds, `CGO_ENABLED=0 GOTOOLCHAIN=local go build ./cmd/sh-agent`:
  - `GOOS=linux GOARCH=amd64` → exit 0, `.agents/scratch/verify-search-index-final/builds/sh-agent-linux-amd64` (13,738,527 bytes)
  - `GOOS=darwin GOARCH=arm64` → exit 0, `.../sh-agent-darwin-arm64` (13,457,090 bytes)
  - `GOOS=linux GOARCH=arm64` → exit 0, `.../sh-agent-linux-arm64` (13,101,460 bytes)
- `/usr/bin/go version` → `go version go1.24.4 linux/amd64`
- `GOTOOLCHAIN=local /usr/bin/go build -o .agents/scratch/verify-search-index-final/builds/sh-agent-usrbingo ./cmd/sh-agent` → exit 0
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -v` → PASS. Log line: "Antigravity rules block is 21988 bytes (context 21988, rules 0), budget 22000" — under the 22,000-byte budget (ADR-035).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-final` → "266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed", exit 0 (a PATH-priority note and a settings.json creation note, neither an error)
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-final` → `binary  aae16af6f28d+dirty` (no warning line); `drift  none — 265 files match the repo`
- `grep -rn "Search first" .agents/scratch/verify-final/.claude/AGENTS.md` and `.../.gemini/config/rules/AGENTS.md` → both found, identical line: "**Search first:** `sh-agent search <words>` ranks files and trail records; ..."
- `sh-agent index --rebuild --home .agents/scratch/verify-final-home`, native build, 3 runs: 2157 ms, 2329 ms, 2313 ms — all under the 10 s budget (3 of 3).
- Warm `sh-agent search harness --home .agents/scratch/verify-final-home` (one untimed warm-up run, then 5 timed runs): 39 ms, 38 ms, 39 ms, 38 ms, 39 ms — all under the 100 ms budget (5 of 5, end to end including process start).
- `SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v -count=1` → PASS (3.73s–3.93s across two runs). Both arms' means (all-41-question mean): `grep mean all 0.379`, `index mean all 0.677` — the index arm scores well above the grep baseline; per-category means also reported (words, drift, refs, open) for both arms in the log.
- `go test -cover ./cmd/sh-agent/... ./internal/harness/... ./internal/search/...` → cmd/sh-agent 83.3%, internal/harness 87.2%, internal/search 84.1%, all `ok`.

## Root cause  (FAIL only)
Not applicable — no failures.

## Coverage of changed code
Every changed non-test `.go` file in the range maps to a covering test file in the same package:
- `cmd/sh-agent/main.go` (new `search`/`index` flags and dispatch) — covered by `cmd/sh-agent/search_test.go` (new, 112 lines)
- `internal/harness/plan.go`, `internal/harness/state.go` — covered by `internal/harness/harness_test.go` (existing suite, package coverage 87.2%)
- `internal/harness/report.go` — covered by `internal/harness/report_test.go` (new assertions, 26 lines added)
- `internal/search/parse.go`, `internal/search/reconcile.go`, `internal/search/search.go`, `internal/search/store.go` — covered by `internal/search/search_test.go`, `internal/search/global_test.go`, `internal/search/security_test.go`, `internal/search/recall_test.go` (package coverage 84.1%)

No changed non-test `.go` file was found without a covering test.

### Suppressions added in the range
`git diff c04f2d7..HEAD -- '*.go'` shows three added `t.Skip(...)` calls, all with a reason on the line:
- `internal/search/recall_test.go` (TestRecall): `t.Skip("measurement, not a gate: set SH_AGENT_RECALL=1 to run it")`
- a helper-process test: `t.Skip("helper process")`
- a root-privilege-dependent test: `t.Skip("root reads a mode-000 directory; the check cannot fail here")`

No added `.only(`, no bare/unreasoned skip, no `//nolint`, `//lint:ignore`, `@ts-ignore`, `@ts-expect-error`, or `eslint-disable` lines found in the diff.

## Not examined
- `.agents/docs/testing/skill-matrix.md` full matrix — out of scope for this dispatch (not requested).
- Behavioural/sandbox CLI runs (Claude Code / Antigravity actual starts) — not requested by this dispatch; the dispatch asked for install/status + context grep only, which was done.
- `go test -cover` was run only for the three touched packages (cmd/sh-agent, internal/harness, internal/search), not `./...`; the untouched root package (`cmd/sh-agent` vs top-level `styrman-harness` module test) was covered implicitly by the full `-race` run above but not re-measured for coverage percentage.
- Did not investigate who/what produced the concurrent `.agents/docs/reviews/2026-09-29-review-search-index-code-r2.md` file noted above — flagged for the caller, not chased down (out of scope for a verifier, and doing so would require reading process state outside this dispatch's remit).

## Environment
Linux 6.12.107+deb13-amd64. Default `go` on PATH: go1.27.1 linux/amd64 (user-level SDK ahead of `/usr/bin/go`, per AGENTS.md § Setup Commands). `/usr/bin/go`: go1.24.4 linux/amd64 (matches go.mod's `go 1.24.4`, no toolchain directive, ADR-029). All scratch artifacts under `.agents/scratch/verify-search-index-final/`, `.agents/scratch/verify-final/`, and `.agents/scratch/verify-final-home/` (gitignored, never `/tmp`, never the real home).
