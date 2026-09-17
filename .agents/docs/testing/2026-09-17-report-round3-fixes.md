---
owner: verifier
date: 2026-09-17
---
# Verification — sh-agent report round 3 fixes

Range reviewed: 012b4ca..worktree (uncommitted changes on top of 012b4ca; internal/harness/report.go, internal/harness/report_test.go modified, plus docs/bugs BUG-053..060, reviews, plan, spec, guide, README, REQUIREMENTS, context/AGENTS.md)

## Status: PASS

Scope: complete for the checks named in the dispatch — build/vet/format, full race test suite, test-list cross-check, govulncheck, install/status against a scratch home, dry-run no-write check, must-fail mutant log review, and behaviour-to-test coverage mapping for the five round-3 fix areas. Not exercised: the real feedback-repository remote (push path), full skill/agent install matrix.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 -v ./...` → exit 0. counts: pass 69, fail 0, skip 1 at 012b4ca+dirty for `go test -race -count=1 -v ./...`. The one skip is `TestFindRepo` ("working directory happens to be inside a harness repo"), an environment-dependent pre-existing skip unrelated to the report change. vs previous: `.agents/docs/testing/2026-09-17-report-round2-fixes.md` recorded `counts: not printed` for a bare (non-verbose) run, so no comparable prior counts line exists; treating as vs previous: none.
- `go test -list 'TestReport|TestSendReport|TestResolveFeedback' ./internal/harness/` → exit 0, lists 13 test names (TestReportRedactJSONL, TestReportRedactShapes, TestReportRedactStructuralGaps, TestReportRedactCleanIsByteIdentical, TestReportRedactText, TestReportRedactDocument, TestReportCollect, TestReportRecord, TestSendReport, TestResolveFeedbackFailures, TestSendReportRefusesRedirectsAndLinks, TestSendReportPushesOnlyToFeedbackURL, TestSendReportDryRunShowsDroppedRecords). `grep -c '^func Test\(Report\|SendReport\|ResolveFeedback\)' internal/harness/report_test.go` → 13. Match: 13 = 13.
- `govulncheck ./...` (govulncheck@v1.8.0, go1.27.1 toolchain) → "No vulnerabilities found."
- `go run ./cmd/sh-agent install --home .agents/scratch/r3-home` → exit 0, "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
- `go run ./cmd/sh-agent status --home .agents/scratch/r3-home` → `binary 012b4ca93838+dirty` (no warning), `drift none — 259 files match the repo` (259 vs 260 written accounted for by the generated settings.json excluded from drift comparison by design).
- Dry-run no-write check: file count under `.agents/scratch/r3-home` was 262 before `go run ./cmd/sh-agent report --dry-run --home .agents/scratch/r3-home` and 262 after. Command exit 0; output listed "would clone", "would write .../state.json", "would write .../report.json" — no files actually created.

## Must-fail mutant evidence (cited from .agents/scratch/r3-mutants/, not re-created)
- `2b-label.log` → exit non-zero, `--- FAIL: TestReportRedactStructuralGaps`. Covers the fix that a block closes only at a footer naming its own opening label, not any footer mention.
- `3-authz.log` → exit non-zero, `--- FAIL: TestReportRedactShapes`. Covers the quoted `"Authorization"` key shape.
- `4-decoded.log` → exit non-zero, `--- FAIL: TestReportRedactStructuralGaps`. Covers `decoded()` reading every JSON value on a line, numbers kept as written, including a number beyond float64 and two records concatenated on one line.
- `6-document.log` → exit non-zero, `--- FAIL: TestReportRedactDocument`. Covers `redactDocument`'s decoded branch (an escaped needle only visible after unescaping).
- `7-dryrun.log` → exit non-zero, `--- FAIL: TestSendReportDryRunShowsDroppedRecords`. Covers the dry-run list reporting an unterminated, dropped block.

One mutant per fix area, each caught by the test named in the dispatch; no fixture text reproduced here.

## Coverage of changed code (behaviour → test)
- Block closes only at a footer naming its own header label (not any footer mention) — `TestReportRedactStructuralGaps`, case `"header line that mentions a footer"` (report_test.go:117).
- Quoted `"authorization"` key shape (as code/JSON write it) — `TestReportRedactShapes`, entries `headers: {"Authorization": "Basic FAKEVALUE"}` and the Python-dict-style header (report_test.go:72).
- `decoded()` reads every value on a line, numbers kept, including a number beyond float64 and two JSON records on one line — `TestReportRedactStructuralGaps`, cases `"number beyond float64 beside an escaped needle"` and `"two records on one line"` (report_test.go:118-119).
- `redactDocument`'s decoded branch (needle only visible after unescaping) — `TestReportRedactDocument`, the `AKIAFAKE` case (report_test.go:192).
- Dry-run list prints per-file redaction lines including an unterminated block — `TestSendReportDryRunShowsDroppedRecords` (report_test.go:559-570), asserts the exact "feedback.jsonl:1 (-----BEGIN) 2 lines, unterminated" line.

No behaviour named in the dispatch is without a covering test.

## Not examined
- The real feedback-repository remote (git@github.com:kesokaj/styrman-harness-feedback.git) — not pushed to; dry-run and existing SendReport tests use local temp repos/URLs only.
- Full skill/agent install matrix (`.agents/docs/testing/skill-matrix.md`) — not run; this verification scoped to the report change only, per dispatch.
- CLI smoke starts (Claude Code / Antigravity parsing the installed context) from the scratch home — not run; dispatch scoped to install/status/report-dry-run only.

## Environment
Linux 6.12.107+deb13-amd64, go1.27.1 (user-level SDK per AGENTS.md), govulncheck@v1.8.0. Working tree dirty on top of 012b4ca (uncommitted report.go/report_test.go and docs changes under verification). Scratch home used: `.agents/scratch/r3-home`; mutant logs read from `.agents/scratch/r3-mutants/` (pre-existing, not regenerated).
