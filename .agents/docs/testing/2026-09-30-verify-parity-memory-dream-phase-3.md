---
owner: verifier
date: 2026-09-30
---
# Verification — parity-memory-dream-phase-3   (2026-09-30)

Range reviewed: b629ae0..worktree (uncommitted: internal/dream/*, internal/search/dream.go + dream_test.go, internal/harness/report.go + report_test.go, internal/memory/memory.go, cmd/sh-agent/main.go)

## Status: PASS

Scope: guide § 1, § 3, § 3d, coverage map, one must-fail mutation, dry-run scratch-home check ran; guide.md (per guide, cited, not rewritten). FAIL is for coverage only: no test exercises internal/dream/dream.go's `targets`, most of `skipReason`, the skip/unchanged/dropped/--project branches of `dryRun`, the agy `command`, or the state-file corrupt path (listed below). Every function has at least one covering test, so by the "FAIL only when no covering test at all" rule this is a finding, not a hard FAIL; caller decides. Nothing failed at run time.

## Commands (guide: .agents/docs/testing/guide.md)
- `gofmt -l cmd internal` → exit 0, printed nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0 (output .agents/scratch/verify-p3/go-test-1.txt); counts: pass 417, fail 0, skip 2 at b629ae0+dirty for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-09-30-verify-v1.15.0-gate.md, pass 374→417, skip 2→2. Skips: TestRecall, TestHelperIndex (env-gated, unchanged). All 6 packages ok.
  - New tests, all PASS: internal/dream TestCollectInputs, TestCollectUnchanged, TestCollectSkillAloneIsNoChange, TestCollectRedactsAndHidesHome, TestCollectCapDropsTranscriptThenTrail, TestCollectNeverSendsSearchDB, TestDryRunWritesNothing; internal/search TestListProjectsMissingDBIsEmpty, TestListProjects, TestAgyConversationText, TestAgyConversationTextNoAntigravity; internal/harness TestRedactText (shapeSecrets refactor, PASS).
- Mutation (review-standards § 2): collect.go:45 `kind != kindSkill && t.After(since)` → `kind == kindSkill || t.After(since)`. Build exit 0 (compiled). `go test ./internal/dream -run TestCollectSkillAloneIsNoChange -v` → `collect_test.go:143: a skill update alone made the project changed`, `--- FAIL: TestCollectSkillAloneIsNoChange`; TestCollectUnchanged also failed. Restored from copy: `sha256sum -c` OK, `cmp` identical, `go test ./internal/dream` ok again. Line 45 of the original was not otherwise edited.
- § 3: `sh-agent install --home .agents/scratch/p3-home` → exit 0; `status` → `drift none — 267 files match the repo` (29 skills, 8 agents). agents/*.md not in diff, so the Records-section check was not needed.
- § 3d `SH_AGENT_RECALL=1 go test -count=1 ./internal/search -run 'TestRecall$' -v` → PASS (4.45 s). Recall means: words 0.692 (26), drift 0.339 (10), refs 1.000 (3), open 0.000 (2), all 0.595 (41). Not a gate; no earlier numbers compared.
- § 3d rebuild: `.agents/scratch/sh-agent index --rebuild --home .agents/scratch/s` → 1936 ms (real 1.975 s), limit 10 s, OK (24,125,520 bytes).
- § 3d warm search: `search --home .agents/scratch/s verifier guide` real 0.042 s; `search memory dream` real 0.035 s; limit 100 ms, OK.
- Dry-run: scratch home .agents/scratch/p3v-home indexed from two clones of ~/sandbox/p2/fixture at ~/sandbox/p3v/a and /b (each `index --rebuild` run from inside the clone). File list before: one file, search.db. `sh-agent dream --all --dry-run --home <scratch>` → exit 0, printed for each: `[dry-run] …/p3v/a: memory 0, feedback 0, trail 2, skill 0, transcript 0; 2241 bytes; would run: claude -p --tools  --no-session-persistence --output-format text` (same for b). After: file list and sha256 of every file identical (LIST-IDENTICAL, CONTENT-IDENTICAL). Real home never touched; no claude/agy started.

## Root cause
n/a — no test failed. Coverage gap below.

## Coverage of changed code (go test -coverprofile over ./internal/dream ./internal/search)
Every function has a covering test. Per-function %: collect 86.1, scrub 100, fit 100, oldest 100, mdFiles 100, readFile 71.4, feedback 94.4, skills 80, frontField 75, unquote 100, claudeTranscript 86.4, sessionText 95.8, contentText 84.6; dream.go Run 63.6, targets 18.2, skipReason 60, since 66.7, dryRun 66.7, command 66.7; prompt 100; state statePath 100, readState 66.7; search/dream.go roDSN 100, ListProjects 80, AgyConversationText 79.2, agyTime 28.6, inWorkspace 83.3, conversationText 78.6, copyDB 80, copyFile 72.7, removeIfEmpty 100.
Untested branches (0 hits in profile):
- dream.go Run: nil Now/Projects defaults (58-66), the "run path is not built yet" return (71); targets: `--project`, `--all`+`--project` conflict, `--all` via Projects, usage error (78-92) — `targets` is effectively untested; skipReason "gone", "not a git checkout", "no <docs>" (98-109); since state hit (115); dryRun skip / would-fail / unchanged lines and Dropped / --project prompt printing (126-152); command() for Antigravity (160). No test in cmd/sh-agent names dream (grep).
- state.go readState: unreadable and corrupt-JSON fallbacks (32-36).
- search/dream.go: ListProjects error paths (57-78); AgyConversationText error/skip paths (107-157); agyTime string/[]byte/RFC3339 branches (170-178, only the time.Time case hit); inWorkspace (186-192); conversationText/copyDB/copyFile error paths.
Only partly covered, not whole functions: TestCollectNeverSendsSearchDB and TestDryRunWritesNothing are the tests for the two safety properties.
Suppressions added: none (grep over `git diff -U0` and the new files for nolint, lint:ignore, ts-ignore, eslint-disable, t.Skip, .skip(, .only(, go:build ignore printed nothing). sh-agent-matrix.md has no `dream` rows (grep empty): the verb is new, rows missing.

## Not examined
- Whether each named test "fails without" its function for functions other than collect.go's Changed rule: only that one mutation was run.
- Real `claude`/`agy` run path (not built; Phase 4); sandbox tier 4; § 3a, 3b, 3c matrix rows for the new `dream` verb (sh-agent-matrix.md not checked for dream rows).
- Earlier-range recall numbers for comparison.
- Fixture ~/sandbox/p3v left in place (outside repo).

## Environment
Linux 6.12.107 amd64, bash; Go go1.27.1; no Docker needed. Tree at end vs start: only this log added under .agents/docs/testing; no source/test/config path changed during the run (collect.go was mutated and restored byte-identical).

## Re-run after the added tests

Tree: b629ae0+dirty; tests added in internal/dream/collect_test.go, internal/search/dream_test.go, cmd/sh-agent/memory_test.go. Tree at end: no source/test path changed during the run.
- `gofmt -l cmd internal` -> exit 0, printed nothing
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0 (.agents/scratch/verify-p3r/go-test.txt); counts: pass 426, fail 0, skip 2 at b629ae0+dirty for go test -race -v -count=1 ./...; vs previous: this log's first run, pass 417->426, skip 2->2. All 6 packages ok.
- `go test -count=1 -coverprofile=.agents/scratch/verify-p3r/{d,s}.out ./internal/dream ./internal/search` -> exit 0 each; dream total 90.7%.
  dream.go: Run 63.6, targets 90.9, skipReason 90.0, since 100, dryRun 100, command 100; state.go statePath 100, readState 83.3.
  search/dream.go: roDSN 100, ListProjects 80.0, AgyConversationText 79.2, agyTime 100, inWorkspace 100, conversationText 78.6, copyDB 80.0, copyFile 72.7, removeIfEmpty 100.

### Earlier gaps
- targets -> TestTargets (internal/dream/collect_test.go:215); covered 90.9%, only dream.go:82-83 unhit.
- skipReason -> TestSkipReason (:251); covered, only dream.go:105-106 unhit.
- dryRun skip / unchanged / would-run / dropped / prompt lines -> TestDryRunLines (:276: "unchanged since", "gone: skip (gone)", "would run: agy -p", --project "[dry-run] prompt:"), TestDryRunOverCapAndDropped (:305); dryRun 100%. The would-fail line is asserted by name only via those two; dryRun is 100% so its branch is hit.
- command(antigravity) -> TestCommand (:357), 100%.
- since with saved state -> 100% (TestTargets/TestDryRunLines; exact test not isolated).
- readState fallback -> TestReadStateFallback (:337); 83.3%, corrupt-JSON return (state.go:35-36) still unhit.
- agyTime string/bytes/bad -> TestAgyTime (search/dream_test.go:170), 100%; inWorkspace -> TestInWorkspace, 100%.
- cmd/sh-agent dream wiring -> memory_test.go:67-77 (bare dream usage error, --all --dry-run --vendor claude, --project non-checkout "not a git checkout").
Still uncovered (minor): Run nil Now/Projects defaults (dream.go:61-66) and "run path is not built yet" (:71, Phase 4); targets :82-83; skipReason :105-106; state.go:35-36; search error paths (ListProjects, AgyConversationText, copyDB, copyFile, conversationText).
Matrix rows for `dream` / `memory` in sh-agent-matrix.md are a Phase 6 task in the plan: not a FAIL here.
Not examined (re-run): mutation of the new tests (only the earlier collect.go mutation was run); recall/perf tiers not re-run; no CLI started, real home untouched.
