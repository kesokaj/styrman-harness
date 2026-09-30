---
owner: verifier
date: 2026-09-30
---
# Verification — redaction-rules   (2026-09-30)

Range reviewed: working tree — changed files: internal/harness/report.go internal/harness/report_test.go; read the diff with git diff HEAD -- internal/harness/report.go internal/harness/report_test.go.

## Status: PASS

Scope: setup commands from AGENTS.md § Setup commands, mutation protocol on the round-3 redaction needles, then the report --dry-run tier from testing/guide.md § 3b. testing/guide.md already existed; not rewritten.

## Diff summary
`internal/harness/report.go`: prefix-needle list drops `hooks.slack.com/services/` and `hooks.slack.com/workflows/`, adds `sk-None-`; `url token` keyword regex loosens from `[A-Za-z0-9_.%-]{20,}` to `[^/\s:@]{16,}`; new `slack webhook` keyword rule `(?i)hooks\.slack\.com/` covers any Slack webhook path in any case (this is now the *only* Slack-webhook detector — the two old prefix needles were removed, not kept alongside it).
`internal/harness/report_test.go`: `TestReportRedactShapes` gains 5 new must-redact rows (mixed-case `services` webhook, `triggers` webhook, a 16-char userinfo token, a `==`-padded userinfo token, `sk-None-`) and 2 new must-NOT-redact rows (`https://deploy-bot@host/r`, `see slack.com/apps for hooks`).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0. counts: pass 372, fail 0, skip 2 at c7ffb05+dirty for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-09-30-verify-v1.13.0-gate-r4.md, pass 366→372, skip 2→2 (explained: HEAD c7ffb05 is 4 commits ahead of that log's c63f736 — `git log --oneline c63f736..c7ffb05` shows e59638c "Review fixes: ... redaction needles ..." plus 3 doc/handoff commits landed in between and added tests unrelated to the still-uncommitted diff under review here; the working-tree diff itself adds rows inside the existing loop-based `TestReportRedactShapes`/`TestReportRedactStructuralGaps`, no new `t.Run` subtests, so it does not change the PASS-line count on its own). All packages `ok`. Output saved: .agents/scratch/verify-redaction/go-test-1.txt.
- `go run ./cmd/sh-agent install --home .agents/scratch/home` → exit 0. "1 written, 265 unchanged, 0 removed, 0 kept — 266 files managed"; binary rebuilt and written under `.agents/scratch/home/.local/bin/sh-agent` (confirmed by `ls -la`, timestamp matches the run; the pre-existing `~/.local/bin/sh-agent` in the real home was untouched, different timestamp).
- `go run ./cmd/sh-agent status --home .agents/scratch/home` → exit 0. `binary c7ffb05fd9bc+dirty` with no warning; `drift none — 265 files match the repo`.
- `govulncheck ./...` → exit 0. "No vulnerabilities found... 0 vulnerabilities in packages you import and 1 vulnerability in modules you require" — matches AGENTS.md's documented GO-2026-5024 (not reachable); no stdlib os/os-exec advisories reported (box's go1.27.1 SDK is past the go1.25.8/1.24.6 threshold named in AGENTS.md). No new findings. Output saved: .agents/scratch/verify-redaction/govulncheck.txt.

### Mutation protocol (review-standards § 2) on the round-3 needles
Backed up `internal/harness/report.go` to `.agents/scratch/verify-redaction/report.go.orig`, then hand-edited three spots in place: removed the `slack webhook` keyword rule, reverted the `url token` regex to the pre-round-3 `[A-Za-z0-9_.%-]{20,}` class, and removed `"sk-None-"` from the prefix-needle slice (leaving the two round-2 Slack prefix needles removed, since they were already deleted by the diff under review, not restored — this mutation is a strict subset removal from the current working tree, not a full round-2 reconstruction).
- `go build ./...` after the mutation → exit 0 (compiles).
- `go test ./internal/harness/... -run TestReportRedactShapes -v` → exit 1 (FAIL), as expected. Actual failures logged, by secret string:
  - `https://hooks.slack.com/services/T0/B0/FAKEVALUE` (pre-existing round-2 row)
  - `https://hooks.slack.com/workflows/T0/A0/FAKEVALUE` (pre-existing round-2 row)
  - `https://HOOKS.SLACK.COM/services/T0/B0/FAKEVALUE` (round-3 row)
  - `https://hooks.slack.com/triggers/T0/FAKEVALUE` (round-3 row)
  - `https://FAKEVALUE1234567@host/r` (round-3 row, 16-char userinfo token)
  - `https://FAKEVALUEFAKEVAL==@host/r` (round-3 row, padded userinfo token)
  - `sk-None-FAKEVALUE` (round-3 row)
  - Note on the dispatch's expectation: the dispatch text said this mutation "must fail on the five round-3 rows" — it does, but 2 pre-existing round-2 rows (`services`, `workflows`) also failed. This is not a discrepancy in the code: those two rows are covered *only* by the new generic `slack webhook` rule in the current diff (their dedicated prefix needles were deleted by this same diff, not kept as a fallback), so removing that one generic rule removes Slack-webhook coverage entirely. Confirmed this is expected, not a bug, by re-reading the diff: `hooks.slack.com/services/` and `hooks.slack.com/workflows/` prefix needles are gone from the working tree, full stop.
- Restored `internal/harness/report.go` from the backup: `git diff --stat -- internal/harness/report.go` after restore shows `1 file changed, 11 insertions(+), 8 deletions(-)`, identical to the pre-mutation `git diff HEAD` shape (also 11+/8-, confirmed by direct comparison — restore is exact).
- `go build ./...` after restore → exit 0.
Verdict: mutation protocol satisfied for the round-3 rows named in the dispatch (5/5 fail with the rules removed), plus it surfaced that Slack-webhook coverage in this diff now depends entirely on one generic rule rather than on independent prefix needles — noted below, not a FAIL, since the generic rule's own test rows exist and pass in the unmutated tree.

### Report dry-run tier (testing/guide.md § 3b)
- `go run ./cmd/sh-agent report --dry-run --home .agents/scratch/home` → exit 0. Output: "report: skipped feedback.jsonl (absent)", "report: skipped retros.jsonl (absent)", "report: skipped proposals/ (absent)", then three `[dry-run] would clone ...` / `would write ...` lines naming a scratch path under `.agents/scratch/home/.config/sh-agent/report-*`. No network call made (dry-run only per the guide's warning), no write to the real home.

## Root cause
N/A — no FAIL in scope. See the mutation-protocol note above for a non-blocking observation (Slack-webhook detection is now single-point via one generic regex rather than defense-in-depth via separate prefixes), reported as-is, not fixed.

## Coverage of changed code
- `internal/harness/report.go` (prefix/keyword needle changes) — covered by `internal/harness/report_test.go::TestReportRedactShapes`, updated in the same diff with 5 new must-redact rows and 2 new must-NOT-redact rows exercising exactly the new/changed rules. Mutation protocol confirms the coverage is real (rules fail to redact when removed).
- No added `//nolint`, `//lint:ignore`, `t.Skip`, `.only(`, or other suppression found in the diff (`git diff HEAD -- internal/harness/report.go internal/harness/report_test.go` contains none).
- The 2 pre-existing skips in the full suite are outside `internal/harness/report.go`/`report_test.go` and unrelated to this diff (skip count unchanged vs. previous log, 2→2); not traced further as out of scope for this diff.

## Not examined
- Sandbox and real-home CLI smoke tiers (explicitly excluded by the dispatch).
- The identity of the 2 skipped tests in the full suite (not traced to a name — unrelated files, unchanged skip count vs. previous log, judged out of scope).
- `go run ./cmd/sh-agent report` in its real (non-dry-run, network-pushing) mode — correctly out of scope per the guide's own warning never to run it without `--dry-run` in a dev/verification session.
- `/review --docs` over this range was not run (not requested by the dispatch; this diff touches only `internal/harness/report.go`/`report_test.go`, no interface/command/deploy-step change apparent).

## Environment
OS: Linux 6.12.107+deb13-amd64. Go: user-level go1.27.1 SDK per AGENTS.md. Working tree dirty (the diff under review) at base commit c7ffb05. No Docker/Colima involved (no Dockerfile in this change).
