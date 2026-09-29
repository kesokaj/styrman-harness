---
owner: verifier
date: 2026-09-29
---
# Verification — search-index-final-2   (2026-09-29)

Range reviewed: aae16af..HEAD (HEAD=4c9a327)

## Status: PASS

Scope: complete — re-verified HEAD after aae16af..HEAD changed internal/search/reconcile.go (git runs with LC_ALL=C), internal/search/search.go (comment only), internal/search/security_test.go (new TestDropReasonIgnoresLocale), plus docs (plan, review, efficiency testing log, the previous final log). Ran gofmt, vet, the full race suite, the new locale test in isolation, go.mod version check, three release cross-builds, and the context-size test. `.agents/docs/testing/guide.md` exists; not rewritten (no procedure change).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race -v -count=1 ./...` → all packages `ok` (styrman-harness, cmd/sh-agent, internal/harness, internal/search), no FAIL line, exit 0. Output saved to `.agents/scratch/verify-search-index-final-2/go-test-1.txt`.
  counts: pass 359, fail 0, skip 2 at 4c9a327 for `go test -race -v -count=1 ./...`
  vs previous: `.agents/docs/testing/2026-09-29-verify-search-index-final.md` (as stated in its text) claimed "pass 215, fail 0, skip 2" — but its own saved artifact `.agents/scratch/verify-search-index-final/go-test-1.txt` (still on disk, unmodified) actually contains 358 `--- PASS:` lines and 2 `--- SKIP:` lines: the previous log's stated count (215) does not match its own saved output (358). Flagging this as a discrepancy in that prior record, not a defect in this run. Comparing this run's raw output against that prior run's raw output (the ground truth artifacts, both greppable): 358→359 pass, 2→2 skip. The +1 pass is exactly `TestDropReasonIgnoresLocale`, the one test added in this range (`git diff aae16af..HEAD -- internal/search/security_test.go`); no other test count changed. No fall in pass, no rise in skip.
  The 2 skips are unchanged from the prior run: `TestRecall` (gated by `SH_AGENT_RECALL=1`, carries a reason) and `TestHelperIndex` (subprocess-helper pattern, carries a reason).
- `go test -count=1 -v ./internal/search -run TestDropReasonIgnoresLocale` → RAN (not skipped), `--- PASS: TestDropReasonIgnoresLocale (0.00s)`, exit 0. Confirms the new test exercises the reconcile.go change: it installs a fake `git` that prints the German "not a git repository" message unless `LC_ALL=C`, sets the caller's `LC_ALL=de_DE.UTF-8`, and asserts `dropReason` still returns "no longer a git checkout" — i.e. the code's own `cmd.Env = append(os.Environ(), "LC_ALL=C", "LANGUAGE=")` line is what makes the match succeed.
- `grep '^go ' go.mod` → `go 1.24.4` (unchanged, no toolchain directive per ADR-029)
- Cross-builds, `CGO_ENABLED=0 GOTOOLCHAIN=local go build -o .agents/scratch/verify-search-index-final-2/builds/sh-agent-<os>-<arch> ./cmd/sh-agent`:
  - `GOOS=linux GOARCH=amd64` → exit 0, 13,738,423 bytes
  - `GOOS=darwin GOARCH=arm64` → exit 0, 13,457,090 bytes
  - `GOOS=linux GOARCH=arm64` → exit 0, 13,100,956 bytes
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -v` → PASS. Log line: "Antigravity rules block is 21988 bytes (context 21988, rules 0), budget 22000" — same byte count as the previous final log (context/rules under this range did not change), under the 22,000-byte budget (ADR-035).

## Root cause  (FAIL only)
Not applicable — no failures.

## Coverage of changed code
Changed non-test `.go` files in aae16af..HEAD:
- `internal/search/reconcile.go` — added `cmd.Env = append(os.Environ(), "LC_ALL=C", "LANGUAGE=")` to the `git` helper — covered directly by `internal/search/security_test.go:TestDropReasonIgnoresLocale` (new in this range; forces a foreign `LC_ALL` on the caller and a fake `git` binary that only emits the C-locale string when it sees `LC_ALL=C`, then asserts `dropReason` still matches it).
- `internal/search/search.go` — comment-only change (clarifying `printable`'s tab/newline handling); no code path changed, no test required.
- `internal/search/security_test.go` — the test file itself (new test, not source under test).
No added suppressions found in the diff: `git diff aae16af..HEAD` was checked for `//nolint`, `//lint:ignore`, `@ts-ignore`, `@ts-expect-error`, `eslint-disable`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, `.only(` — none present.

## Not examined
- `govulncheck ./...` — not run this pass (no dependency or toolchain change in this range; the previous final log already recorded the one known advisory `GO-2026-5024` on go1.24.4, unaffected by this range's files).
- `go run ./cmd/sh-agent install --home ...` / `status` drift check, `sh-agent index --rebuild` timing, warm `search` timing, `SH_AGENT_RECALL=1 TestRecall`, and `go test -cover` percentages — not re-run; nothing in aae16af..HEAD touches the installer, index/search runtime behavior, or recall scoring (only `dropReason`'s locale handling and a comment changed), and the prior final log already exercised all of these against the same `internal/search` package at aae16af with the git helper behavior otherwise identical. This run only re-verifies what actually changed plus the full non-behavioral suites (gofmt/vet/race/build/size).
- `/usr/bin/go` build with the system Go — not re-run; `go.mod`'s `go 1.24.4` line is unchanged and the prior log already confirmed it builds there.
- No behavioral/Antigravity CLI run was requested or performed for this range (no `context/`, `agents/`, or `skills/` files changed).

## Environment
OS: Linux 6.12.107+deb13-amd64. `go version`: go1.24.4 linux/amd64 (module `go.mod` line matches, no toolchain pin, ADR-029). Docker/Colima not used (no Dockerfile change in range).
