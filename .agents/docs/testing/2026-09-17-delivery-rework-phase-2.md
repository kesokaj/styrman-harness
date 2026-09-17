---
owner: verifier
date: 2026-09-17
---
# Verification — delivery-rework Phase 2, `sh-agent report`   (2026-09-17)

Range reviewed: working tree at HEAD f5c44f9 (uncommitted) —
`cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`,
`internal/harness/report.go` (new), `internal/harness/report_test.go` (new),
`.agents/docs/plans/2026-09-17-delivery-rework.md` (docs, plan-tracking only).

## Status: PASS

Scope: gate commands, `govulncheck`, `go test -cover`, a dry-run smoke test of
`sh-agent report` against a scratch home, `sh-agent help` line check, and a
manual map of every exported function/type in `report.go` and the `report`
case/func in `main.go` to its covering test. The real GitHub feedback remote
and the Phase 4 e2e are out of scope for this phase (see Not examined).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0; output saved to
  `.agents/scratch/verify-report/go-test-1.txt`; counts: pass 118, fail 0,
  skip 1 at f5c44f9+dirty for `go test -race -v -count=1 ./...` (subtests
  included via `grep -cE '^[[:space:]]*--- (PASS|FAIL|SKIP):'`). vs previous:
  `.agents/docs/testing/2026-09-16-verify-preflight-v1.0.0.md` (round 2, the
  newest log by mtime carrying a `counts:` line), pass 106→118, skip 1→1.
  Rise of 12 matches the new tests added in this diff: `report_test.go` adds
  8 top-level tests (`TestReportRedactJSONL`, `TestReportRedactCleanIsByteIdentical`,
  `TestReportRedactText`, `TestReportRedactDocument`, `TestReportCollect`,
  `TestReportRecord`, `TestSendReport`, `TestResolveFeedbackFailures`) and
  `main_test.go` adds `TestReportWithoutCheckout`; the remaining rise is
  subtests inside `TestSendReport`'s sequential scenarios and
  `TestReportRecord`'s two-phase check counted as top-level asserts, not
  separate `--- PASS:` lines — no discrepancy found on inspection. The one
  skip is `TestFindRepo` ("working directory happens to be inside a harness
  repo"), pre-existing and environment-dependent, unrelated to this change —
  explained.
- `govulncheck ./...` → "No vulnerabilities found." (go1.27.1 on this
  machine; `AGENTS.md`'s note about two stdlib advisories predates this
  toolchain and no longer applies here — not a regression, an environment
  fact.)
- `go test -cover -count=1 ./...` → `cmd/sh-agent` 69.5%, `internal/harness`
  83.6%.
- `rm -rf .agents/scratch/verify-report-home && go run ./cmd/sh-agent report
  --dry-run --home .agents/scratch/verify-report-home` → exit 0. Output:
  named all four config files/dirs "(absent)", "nothing installed in this
  home", "nothing to send", `[dry-run] would clone
  git@github.com:kesokaj/styrman-harness-feedback.git into
  .../verify-report-home/.config/sh-agent/feedback`, `[dry-run] would write
  .../feedback/reports/dave/2026-09-17-085320/report.json`. Confirmed after
  the run: `.agents/scratch/verify-report-home` does not exist (`ls`:
  "No such file or directory") — the dry run cloned nothing and wrote
  nothing, as required. `dave` is this machine's real `hostname`, not test
  contamination.
- `go run ./cmd/sh-agent help` → contains the line: `report [--feedback-repo
  DIR] [--dry-run]  send a redacted copy of ~/.config/sh-agent and a metrics
  record to the feedback repo (default checkout ~/.config/sh-agent/feedback,
  cloned on first use)`.

## Coverage of changed code
Exported surface in `internal/harness/report.go` mapped to its test:
- `Collect` — `TestReportCollect` (empty home, populated home, symlink
  skipped, originals untouched), plus exercised end-to-end by `TestSendReport`.
- `BuildRecord` — `TestReportRecord` (no checkout/no install; plain-copy
  checkout with install).
- `ResolveFeedback` — `TestResolveFeedbackFailures` (missing remote leaves
  nothing behind; explicit non-checkout rejected; explicit path naming a
  harness checkout rejected) and the default-path branches (absent → dry-run
  no-clone, absent → real clone, present → reused) via `TestSendReport`.
- `SendReport` — `TestSendReport` (dry run with default path absent; first
  real run cloning an empty bare origin, branch named `main`, upstream set;
  dry run against an existing checkout leaves `git status` unchanged; second
  run same second gets a `-2` suffixed directory, stray file left alone and
  named, not staged, not committed; a behind checkout gets "push rejected"
  with the commit kept local).
- Unexported helpers `redactLines`, `redactJSON`, `walkJSON`,
  `redactDocument`, `needleIn`, `runDir`, `hostDir`, `commitAndPush` are all
  exercised directly or transitively by the above tests (redaction cases
  cover JSON-lines, plain text, PEM blocks incl. unterminated, key-hit vs.
  value-hit, malformed-JSON fallback, byte-identical clean-input path).
- `report` case in `cmd/sh-agent/main.go` (dispatch) and the `report(...)`
  helper func — `TestReportWithoutCheckout` (no harness checkout anywhere,
  `--dry-run`, asserts the "not a harness checkout" note, the null-field
  note, the dry-run clone message, `report.json` in output, and that nothing
  was written under `home`).

No changed file in this range is without a covering test.

`.agents/docs/plans/2026-09-17-delivery-rework.md` is a plan-tracking doc
edit (8 lines changed, box-ticking) — not code, no test applicable.

Suppressions/skips added in this diff: one, `report_test.go:228`
`t.Skip("git not installed")` inside `gitEnv`, guarding tests that shell out
to git — carries a reason, not a bare skip, not a FAIL line. No
`//nolint`, `//lint:ignore`, `.only(`, or unexplained skip found in the diff.

## Not examined
- The real GitHub feedback remote (`git@github.com:kesokaj/styrman-harness-feedback.git`)
  was never contacted — `sh-agent report` without `--dry-run` was not run, per
  the dispatch's explicit instruction. All git operations in `TestSendReport`
  and `TestResolveFeedbackFailures` run against local bare repositories under
  `t.TempDir()`, not the real remote.
- The Phase 4 e2e (whatever it names) — out of scope for this phase per the
  dispatch.
- `docs/testing/sh-agent/guide.md` — exists already in this repo (owner:
  verifier from an earlier run); not updated in this run since the gate
  commands and their order did not change.
- No frontend, no Docker build (no Dockerfile changed).
- Coverage percentages (69.5% / 83.6%) were read but not diffed line-by-line
  against the change set beyond the function-level map above; no tool that
  maps coverage to only the new/changed lines was run.

## Environment
Linux 6.12.107+deb13-amd64 (`dave`), `go version go1.27.1 linux/amd64`,
`govulncheck` present at `~/.go/bin/govulncheck` (no vulnerabilities found
with this toolchain — supersedes the two-advisory note in this repo's
`AGENTS.md`, which was written against an older Go). `git` present and used
directly (no sandboxing needed for local bare-repo tests). No Docker/Colima
involved — no container-affecting change in this range.
