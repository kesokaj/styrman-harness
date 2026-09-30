---
owner: verifier
date: 2026-09-30
---
# Verification — parity-memory-dream phases 4-5 (2026-09-30)

Range reviewed: 1bf5f4f..worktree (uncommitted, HEAD 1bf5f4f)

## Status: PASS

Scope: guide §1, §3, §3b run; §3c rows NOT present (Phase 6 owns; listed below); mutation and audit done. Not "complete": no sandbox CLI starts, no timer against a real manager (by instruction). guide.md exists (no change).

## Commands
- `gofmt -l cmd internal` -> prints nothing
- `go vet ./...` -> exit 0
- `go test -race -count=1 -v ./...` -> exit 0 (.agents/scratch/verify-p45/go-test-1.txt); counts: pass 443, fail 0, skip 2 at 1bf5f4f+dirty for go test -race -count=1 -v ./...; vs previous: .agents/docs/testing/2026-09-30-verify-parity-memory-dream-phase-3.md, pass 426->443, skip 2->2 (TestRecall, TestHelperIndex, env-gated, unchanged). 6 packages ok.
  New tests, all PASS: internal/dream run_test.go TestRunWritesProposalsDraftsAndState, TestRunAgyGetsThePromptAsArgument, TestRunInvalidOutputKeepsRawAndWatermark, TestRunCLIExitFails, TestRunNoCLIFailsBeforeProjects, TestRunLock, TestRunPrunesOldFailures, TestStatusLine; timer_test.go TestInstallTimerLinux, TestInstallTimerRefusals, TestRemoveTimerLinux, TestTimerMacOS; harness TestReportCollect (extended); cmd TestHarnessMap, TestDreamVerb, TestUninstallDryRunLeavesTheTimer, TestMemoryRefusalExitsTwo.
- § 3: `go run ./cmd/sh-agent install --home .agents/scratch/verify-p45/home` -> 268 written; `status` -> "drift none — 267 files match the repo"; line `dream      never run; timer not installed`.
- § 3b: scratch home holding proposals/dream/harness/2026-09-30-h.md, proposals/dream/projA/p.md, proposals/dream/loose.md, proposals/2026-09-30-retro.md, drafts/d.md. `report --dry-run --home ...` would write only: state.json, proposals/2026-09-30-retro.md, proposals/dream/harness/2026-09-30-h.md, report.json. projA file, loose.md, drafts NOT listed. PASS. (Symlink guard: unit-tested only, not hand-run.)
- `go test -cover ./internal/dream ./internal/harness` -> dream 88.7%, harness 87.7%. Per function run.go/timer.go: every function has some coverage, none untested; lowest: parse 63.6%, unique 75.0%, keepFailure 75.0%, renderLaunchd 80.0%, realRun 82.9%.
- Mutation (review-standards § 2): inserted `st.Projects[p.Root] = t0; saveState(o.Home, st)` before `logf(... "fail invalid output")` in run.go; compiled; `go test ./internal/dream -run TestRunInvalidOutputKeepsRawAndWatermark` -> FAIL "run_test.go:166: watermark advanced on a failed project". Restored from copy: sha256 4871e74a...5ed421 identical before/after; test passes again.
- Suppressions added (nolint/t.Skip/ts-ignore) in diff and new files: none.

## Coverage of changed code
- internal/dream/run.go, timer.go -> run_test.go, timer_test.go
- internal/dream/dream.go, prompt.go -> dream tests (TestStatusLine, dry-run tests)
- internal/harness/report.go -> report_test.go TestReportCollect
- cmd/sh-agent/main.go -> main_test.go, memory_test.go
- docs/skills (README, AGENTS.md, local.md, dream.md spec, retro SKILL.md): no code; rendered-context size test passed in go test.
- Gap (not FAIL): the prompt.go warning text is not asserted by name in a test I could identify (not checked).

## § 3c rows missing from testing/sh-agent-matrix.md
No row for `memory add|show`, `dream` (--dry-run, --all, --project), `dream --install-timer`/`--remove-timer` (or whatever the flags are named), nor the timer line of `uninstall`. `grep dream|memory` on the file: no hits. Phase 6 owns adding.

## Audit of .agents/docs/testing/2026-09-30-verify-parity-memory-dream-phase-4.md
Supported by files: E3 dream.log lines (ok 2 / ok 2 / 4 proposals; second run unchanged x2); 53 s (state t0 20:38:49 -> finish 18:39:42Z); proposals proj-a/proj-b and 2 drafts; proj-a content (queued/project/rule, BUG-001/002 and both memories evidence); `proj-c skip gone` and proj-c dir absent; FAKEVALUE: leak.md exists in proj-a memory, dry-a2.txt shows memory/leak.md with `[redacted by sh-agent report]`, grep FAKEVALUE = 0 in dry-a.txt, dry-a2.txt and all proposals (not vacuous); agy attempt 1: dream-failed/2026-09-30-proj-d.txt stderr matches the quoted warning and denial, stdout empty, no proj-d watermark before the later ok; agy runs 2-4: three ok lines, proj-d.md/-2/-3 and three drafts; Claude `--project proj-a` ok 1 proposals (18:46:40Z).
Claims NOT fully supported or imprecise:
1. The log says "drafts/2026-09-30-proj-a-tag-release/SKILL.md" and "drafts/..." as if beside the home's own `drafts/`; the files are at `<home>/.config/sh-agent/proposals/dream/drafts/...` (the proposal's own text uses the relative `drafts/...`). Path wording only.
2. "Second run, nothing changed ... proposal files 4 -> 4" counts md plus draft files; the dir now holds 3 md + 3 drafts because the later Claude `--project proj-a` run added proj-a-2.md and a draft. Consistent only if read as a point-in-time count; not re-verifiable now.
3. "Claude Code with the same warning ... exit 0": exit code is not recorded anywhere on disk (only the log line). The "both CLIs" attribution for runs 2-4 is unverifiable: dream.log in p2-agy shows three agy-home runs; the Claude run is in p4home. Fine, but the warning wording in the prompt is not checked by me.
4. "sh-agent status ... on a home with nothing installed stops before its retro and dream lines" and the `status` dream line on an installed home: second confirmed here (scratch home); the first not re-run.
5. Timing 53 s is derived from state t0, not from a recorded wall-clock; `run1.txt`/`run1.err` in ~/sandbox/p4 are empty (0 bytes), so they support nothing.
No claim contradicted by the files.

## Not examined
Sandbox CLI starts (forbidden); systemctl/launchctl against a real manager; a hand-run symlink case for report Collect; prompt.go warning text against a test; real `report` push; docs prose (README, dream.md spec, local.md, AGENTS.md rule 1, retro source 9) read for accuracy; search.db content of the p4 home; flag names for timer verbs.

## Environment
Linux 6.12 (Debian 13), go as on PATH, real HOME untouched, scratch under .agents/scratch/verify-p45. Tree at end: same 18 paths as at start plus this log. No Docker/Colima needed.
