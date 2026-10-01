---
owner: reviewer
date: 2026-10-01
---
# Review — audit-3-search-memory-dream   (2026-10-01)

Range reviewed: no diff; scope at HEAD 2ba4a50: `internal/search/*.go`, `internal/memory/*.go`, `internal/dream/*.go` and their tests, the `search`/`index`/`memory`/`dream` verbs in `cmd/sh-agent/main.go`, `.agents/docs/specs/sh-agent/{search,memory-loop}.md` (plus `dream.md`, `delivery.md:77` for SD-spec), `skills/{dream,project-memory,distill,retro}/SKILL.md`, `agents/dreamer.md`, `.agents/docs/testing/guide.md` § 3d and § 3e, `.agents/docs/bugs/BUG-12[4-9]-*.md` and `BUG-13[0-6]-*.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing in this slice is broken. All twelve SD items are fixed or partly fixed in the code, each Phase 1 fix has its named test present, and the suites are green. The worst remaining problems are low severity. A nightly run in which every project fails records `last_run.error` as empty, so `sh-agent status` reports it like a clean run. The recall number in search.md (0.653) no longer reproduces: 0.584 measured, because two "open" questions in `recall.tsv` point at records that have since closed. I did not re-run the must-fail mutations (read-only review). Whether the named tests fail without their fix rests on the plan's recorded runs and on my reading of the fixtures.

## Findings

| ID | was | now | evidence |
|---|---|---|---|
| SD1 | high | fixed | `walk` (internal/search/dream.go:~300) tries `walk(field)` nested first and `textField` only when that fails. `TestProseLeaves` (dream_test.go:245) holds both glued shapes (`pb(10, pb(1, …"command"…))`, `pb(11, pb(4, …))`) and asserts `id_rsa` and `hunter2` are absent. Corpus re-measured by me over copies of the 97 Antigravity conversation databases: 980 prose leaves, 0 with a tool-argument key (2026-09-30: 971 / 0). |
| SD2 | med | fixed | `have` keyed by chunk id (search.go:535-561). `TestNoDuplicateHitAcrossRounds` exists (search_test.go:633) and asserts one line and the header counts. |
| SD3 | med | fixed | `render` indents the text four spaces (`indent`) and folds title, kind, file and evidence through `line` (run.go). `TestTextFieldCannotForgeASection` exists (run_test.go:493) and asserts one `\n**Status:**` and no `\n### 9.`. |
| SD4 | low | fixed | `trailDates` dates a tracked, unmodified record by `git log` and keeps mtime for modified or untracked files (collect.go). `TestTrailDateFromGit` and `TestTrailStatusArgs` exist. The `--no-optional-locks -c core.fsmonitor=false` options come before the verb. |
| SD5 | low | fixed | `parse` decodes leniently and accepts the first `{` that yields proposals or skills. `TestAnswerToleratesExtraKeys` exists (run_test.go:513). |
| SD6 | low | partly | `reconcileGlobal` still runs `gitView` and `apply` for every registered checkout per `--global` query (search.go:~300), so behaviour is unchanged. search.md:46 now says the cost grows per checkout and is not the 15 ms figure. Measured 57 ms global against 34-37 ms project, with one checkout registered. The plan's Phase 1 list has no SD6 line, only the SD-spec doc fix. |
| SD7 | low | fixed | `openCopy` copies, forces a read, and on `isUnreadable` removes the three copies, waits `copyRetryWait` and retries once; it serves all three callers. `TestCopyDBRetriesOnce` (fixture proves NOTADB first) and `TestCopyDBNeverOpensSource` exist. |
| SD8 | low | fixed | `fenceRun` (`>{3,}` and `<{3,}`) escapes text, labels and dropped lines (prompt.go). `TestPromptEscapesFence` exists (run_test.go:527) and covers label, dropped-label, 5-run and `<<<data` shapes. |
| SD9 | low | partly | The corpus check is now a procedure (guide § 3e) and a test (`TestAgyCorpusLeaves`, env-gated, a measurement and not a gate; K=0 is read by hand). `prose()` still guards JSON only, and the corpus needle list is five JSON keys, so a non-JSON argument string with 3 or more spaces would pass unmeasured. |
| SD10 | low | fixed | `contentText` drops string content starting `<local-command`, `<command-`, `<task-notification` or `<system-reminder`, and `sessionText` drops `isMeta` lines. `TestSessionTextSkipsCommandOutput` and `TestSessionTextSkipsMetaLines` exist. I counted the operator's own transcripts for this project: 0 of 7,548 non-meta user lines carry harness text as array parts, so the string-only prefix check matches real data. |
| SD11 | low | fixed | `answer` returns "output line too long or unreadable" on `sc.Err()`, and `removeCopies(conf, o.Now())` uses the injected clock. `TestAnswerLineTooLong` and `TestRemoveCopiesUsesInjectedClock` exist. |
| SD-spec | med | fixed | delivery.md:77 now says "built, not yet released"; search.md:46 states the `--global` cost; dream.md:21,29,34 state git dating, lenient decode and the dropped-draft file; memory-loop.md:12 states `add` refuses the 201st line while `show` caps at 200 lines and 25 KB, and lists `metadata.node_type`; the parity plan's claude command (line 137) carries `--strict-mcp-config --setting-sources local`, which is the code's `command()` flag set. Exit statuses I ran by hand match dream.md: no args exits 2 with the usage line on stdout; `--all --dry-run` exits 0 and writes only `search.db` under the scratch `~/.config/sh-agent`; `memory show` outside a checkout exits 2. |

Phase 1 bug records in the scope: BUG-124, 125, 126, 131, 132, 133, 134, 135 (search/dream/memory) are all `fixed` in file and INDEX. BUG-127..130 (installer) and BUG-136 (trail) are outside this slice; I read only their status.

Test-first and must-fail evidence: every test named above exists and the three packages pass (`go test -cover`: search 84.2%, memory 89.1%, dream 89.6%; 2026-09-30: 83.9 / 89.1 / 89.1). I did not run mutations, so "fails without the fix" is the plan's recorded claim; the fixtures I read (TestProseLeaves, TestTextFieldCannotForgeASection, TestCopyDBRetriesOnce, TestPromptEscapesFence) assert on the defect's own shape and can fail.

New findings:

- **[low / confidence: high]** `internal/dream/run.go:~143` — `finish(len(projects), total, "")` runs after the loop, so a run where every project failed (exit 1) writes `last_run.error = ""`, and `StatusLine` (dream.go) prints `last run <date>: N projects, M proposals` with no failure. A nightly failure is visible only in `dream.log`. Fix: pass `fmt.Sprintf("%d project(s) failed", failed)` as `runErr` when `failed > 0`, and add the status assertion to a failed-project test.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/search.md:89` and `internal/search/testdata/recall.tsv:33,41` — `SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v` (guide § 3d) prints `index mean all 0.584 over 41`, not the spec's 0.653. The gap is largely `q33` and `q41`: their expected files are BUG-083/105 and the 2026-09-29 plan, which are no longer open, so the `open` mean is 0.000 even though `search --open` lists three real open records (BUG-139, BUG-140, the audit-fix plan). Fix: re-point those two rows at stable fixtures, or date the 0.653 as a 2026-09-29 measurement and drop the claim that it holds now.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/search.md:11` — "under two seconds for this repository's 1,080 files"; `index --rebuild` measured 2123 ms for 1,344 files (27 MB). The guide's "under 10 s" gate holds. Fix: refresh the figure and file count, or say "about two seconds".
- **[low / confidence: medium]** `cmd/sh-agent/main.go:592-593` — a `memory` usage error returns a plain error and exits 1 (run: `sh-agent memory bogus` exits 1), while every `memory` refusal and every `dream` usage error exits 2. No spec row covers memory usage errors, so this is a gap, not a contradiction. Fix: return `exitCode(2)` after printing the usage line, as `runDream` does.

Checked, no finding: the search reader (`os.OpenRoot`, Lstat, symlink and submodule modes skipped, tests in security_test.go); dream's reads of the trail (symlink chain refused by `linked`, `IsRegular`, Lstat in `readFile`); the memory store (`trusted` refuses a symlinked `.agents` or `.agents/memory` and tracked content); every lock, state, log, copy and failure file is under `~/.config/sh-agent/` (`ConfigDir(home)`) or the CLI's `dream-work/`; the three SQLite reads are copies opened with `mode=ro` (the live index is read from `dream-projects-*`, the live Antigravity files never opened; `ListProjects` and a dry run leave no copy behind); the prompt: every input is fenced and fence runs are escaped in text, labels and the dropped list; `agyHome` copies only two login files from the real `$HOME` into a 0700 temp directory under conf and removes it. Not a finding: `memory add` bounds lines (200) but not description length, so the index can pass `show`'s 25 KB cap; memory-loop.md:12 states the split.

## Criteria
1. Dead code — pass. No unreferenced helpers found in the three packages; `ListProjects`, `AgyConversationText`, `Store` and `ExitError` all have callers.
2. Flow tracing — pass with the `last_run.error` note above. `dream` exit codes (0, 1, 2) and the `memory` refusal path (exit 2) trace from the packages through `runDream` and `memoryCmd`; the lock, state and watermark order matches dream.md (`Advance` only after a successful write).
3. Standards conformance — pass with nits. The fix-site comments carry audit and review ids; the spec statements I checked match the code (SD-spec). Two measured claims in search.md are stale (recall 0.653, 1,080 files) and `memory`'s usage exit differs from dream's.
4. Directory layout — pass. The SQLite driver import is confined to `internal/search`, and dream and memory reach it only through exported helpers.
5. Simplicity (DIY-except-security) — pass. Stdlib plus the one allowed driver; flock via `syscall`, no hand-rolled crypto or auth; the dream prompt uses redaction from `harness.Redact`.
6. Decoupling — pass. dream consumes `search.Project`, `search.ListProjects`, `search.AgyConversationText` and `memory.Store`/`Root`; the injected `Projects` and `Now` seams keep tests off the real home.

## Not examined
- Must-fail mutations: none re-run (read-only review); the fail-first claims are from the plan and my reading of fixtures.
- Read only in part: `internal/search/{search,reconcile,store,parse}.go` (the dedupe, `reconcileGlobal`, symlink and mode checks only); `internal/dream/{collect,run}_test.go`, `timer_test.go`, `internal/memory/memory_test.go`, `global_test.go`, `search_test.go` (named tests only).
- Not run: a real `claude -p` or `agy` dream call (a real run sends content to a provider); `--global` cost with many registered checkouts (measured with one); `systemd-analyze verify` of the timer units (guide § 3e).
- Skills and agent: `skills/{dream,project-memory,distill,retro}/SKILL.md` and `agents/dreamer.md` were read by grep for exit, cap and Status restatements (retro.md:48 parsing of `**Status:**` against the four-space indent, memory 200 line/25 KB, `proposals/dream`), not end to end; no restatement contradicted the code.
- BUG-127..130 and BUG-136 content (installer and trail): status checked only.
- The 2026-09-30 corpus's 971 / 0 figure: I measured 97 databases / 980 / 0 on this machine; I did not re-check which databases changed.
