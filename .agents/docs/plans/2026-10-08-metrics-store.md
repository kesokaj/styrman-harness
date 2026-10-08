---
owner: main
date: 2026-10-08
---
# Plan — trajectory and agent/sandbox metrics in a local SQLite store   (2026-10-08)

## Why
The operator (2026-10-08): "I want to use the local sqllite more to gather trajectory data and agent/sandbox metrics from all sessions. Give me suggestions on metrics for agent evaulations." Today the SQLite holds nothing about sessions: `sh-agent dream` reads transcript text and then drops it (`internal/dream/collect.go:435-466`). Tool calls, tokens and timings are never read. The only measurement so far is the one-off A/B kit (`testing/ab-kit/metrics.py`, n=1). Claude transcripts on this box cover only 2026-10-05..10-08, so every day without an ingest loses data for good.

Measured baseline (2026-10-08, Claude transcripts, deduplicated by `tool_use.id`, verbs counted at command position):
- 1,619 tool calls.
- 38 of those (2.3 %) were SQLite-backed verbs. 27 were the reviewer's `search --refs --kind bug`; the main thread made 11 of its 612 Bash calls this way.
- grep/rg-style lookups numbered about 848.
- `memory add|show`: 0.
- `search.db`: 29 MB.

## Scope
**In:**
- `~/.config/sh-agent/metrics.db`, a SQLite **store of record**.
- `sh-agent metrics ingest|show|export`.
- Ingest of Claude Code main and subagent transcripts and Antigravity conversation databases as **shape only**.
- Events joined from git, the trail and `feedback.jsonl`.
- Sandbox kit homes ingested under a label.
- `dream --all` ingests before it dreams.
- The metric catalogue below, computed by `show`.

**Out:**
- Arguments, prose or file paths in any row, except the two path columns D12 names.
- Sending anything through `report`.
- A live dashboard or server: the visual report is a static file (D26).
- Any hook.
- Any change to `context/`. The context has 64 bytes of headroom.

**Follow-ups noted, not built:**
- Redacted arguments, which would enable path-level checks (C-gaps).
- Aggregates in `report` (rejected for now, D27).
- A dated `export` run nightly.

**Consumers:**
1. The operator, running agent evaluations: `show`, `--json`.
2. `/retro` and the nightly dreamer, which get evidence (regressions, correction rate, rule violations) from `show --json` and turn it into proposals. Those proposals reach the harness through the existing loop (Phase 7, D27).
3. The A/B and seed kits: `ingest --source --label` replaces `ab-kit/metrics.py` (Phase 5).
4. `dream --all`, which ingests nightly.
5. The operator reading `sh-agent metrics report`, a self-contained HTML file with charts (Phase 6, D26).
6. Not the agents themselves: no agent queries metrics mid-task, and nothing needs it.

### Metric catalogue (what `show` computes)
- **A. Outcome**
  - A1 DoD closure: a commit plus a handoff in the session window.
  - A2 Verifier status (`## Status:` heading of `testing/*-verify-*.md` in the window).
  - A3 Review findings by severity (`**[high|medium|low / confidence: …]**` lines in `reviews/*.md` in the window).
  - A4 Escaped defects: BUG records dated after a session commit.
  - A5 Corrections: `feedback.jsonl` entries in the window.
  - A6 Sandbox pass rate per label, with pass^k (all k runs of a label pass).
- **B. Efficiency**
  - B1 Tokens: in, out, cache read and cache write; cache-hit ratio. Cost is taken from Claude's `cost-state.totalCostUSD`, with no price table.
  - B2 Wall time and active time (gaps under 5 min), turns, tool calls.
  - B3 Peak context and compactions.
  - B4 Subagents by type, with tool calls and SendMessage resumes.
- **C. Process adherence**
  - C1 search : grep ratio.
  - C2 Forbidden shapes (`git add -A`/`git add .`, `push --force`, `--no-verify`, a `git commit` in the same Bash call as `go test`).
  - C3 The first three Reads by trail kind.
  - C4 Skills loaded.
- **D. Reliability**
  - D1 Tool error rate, permission denials and interrupts.
  - D2 Retry loops: the same `cmd_hash` twice in a row.
  - D3 Subagent tool calls.

## Surfaces touched
- New: `internal/metrics/` (`store.go`, `verb.go`, `claude.go`, `agy.go`, `events.go`, `show.go`, tests, `testdata/`).
- `cmd/sh-agent` (verb wiring) and `internal/dream/run.go` (ingest first).
- `REQUIREMENTS.md` R3, R4 and R7 (Phase 0, after approval).
- `AGENTS.md` (Layout and Rules lines naming the one driver importer).
- `README.md` (verb table, state table rows 217-218, the map's dream→configDir label).
- `.agents/docs/specs/sh-agent/metrics.md` (new) and `DECISIONS.md` (ADR-055).
- `.agents/docs/testing/guide.md` § 3e, `testing/sh-agent-matrix.md`, `testing/ab-kit/`.

Code to reuse (package seams: D8; import direction `dream → metrics → search`, never back):
- Open rules: `internal/search/store.go:144-171` (`prepareFiles`, 0600/0700, symlink refusal), exported as `search.PrepareDBFiles`; WAL + busy timeout copied as the DSN string.
- Antigravity copy-then-read: `internal/search/dream.go` `copyDB`/`openCopy`, exported as `search.CopyDB`. A new `search.ProtoFields(b []byte) []Field{Path []int; Varint uint64; Bytes []byte}` sits beside `walk` (`:316-403`); `walk` is unchanged.
- Path kind: `internal/search/parse.go:270` `kindOf`, exported as `search.KindOf`.
- Trail record times: `trailDates` moves from `internal/dream/collect.go:207` to `search.TrailDates` (dream calls it there), so `metrics` never imports `dream`.
- The Claude walk is NOT reused: `claudeTranscript` (`collect.go:355-433`) filters by one root and skips `subagents/`. `metrics` owns its own all-projects walker (D9).

## Interfaces
The spec holds every shape: `.agents/docs/specs/sh-agent/metrics.md` (verbs, flags, the transcript-root rule, `Verb`, the shared tool vocabulary, `show --json` keys, the report, metric definitions, schema, invariants). The D-rows below record why each shape was chosen.

## Failure modes
The spec's `## Failure modes` table.

## Verification
Tiers (each on Claude Code AND Antigravity, D29): unit (`go test -race`, with fixtures for both vendors); an integration run of the binary against a scratch home holding a copy of real transcripts; one Antigravity session run, then `metrics ingest` + `show` (Phase 3).

E2E criteria:
- **Golden path.** `~/.claude/projects` is copied into `.agents/scratch/metrics-home/.claude/projects`; it may hold secrets, and the copy is deleted when the check ends. Then `env -u CLAUDE_CONFIG_DIR sh-agent metrics ingest --home .agents/scratch/metrics-home`. `show --json --all-labels` must report `B.tool_calls` and `B.tool_calls_by_tool`, main and subagent sessions included, equal to an independent count over the same copy: unique `tool_use.id` per tool name. The verb count is not compared, because that would need a second normaliser; `Verb` has its own table tests. A second ingest adds 0 rows.
- **Edge case (canary).**
  - `SECRETCANARY123` is placed in every field the ingester reads that must not be stored: a positional argument, a flag value, a Read path, a tool_result body, user text, assistant text, and a Skill `args`.
  - After ingest, the canary is absent from the bytes of `metrics.db`, `-wal` and the export.
  - Must-fail control: with `Verb` patched to keep positionals (on a copy, review-standards § 2), the test goes red.
- **Migration.** A v1 store with rows opens under a schema list one entry longer and keeps every row. A store whose `meta.schema` is above the binary's is refused, and its bytes are unchanged.

## Phases
- [x] Cold-reader check — one report (`reviews/2026-10-08-review-metrics-store-plan-cold-read.md`), every item answered by a D-row, a rebuttal or the user (D7–D24)
### Phase 0 — Requirements, ADR, spec
- [x] Spike: Claude `message.usage`, `tool_result.is_error`, `cost-state` and `subagents/*.meta.json` exist; Antigravity `gen_metadata` holds the model name and a usage-shaped message — verify: read-only probe, AI — expect: field list recorded in D2, D3 — evidence: probe 2026-10-08 (cost-state: one line per session end or resume, cumulative; 0 lines while a session runs)
- [x] Propose the exact R3/R4/R7 text to the operator — verify: operator's approval in their own turn, user — expect: approved text quoted in a D-row — evidence: D21
- [x] ADR-055 (`.agents/docs/architecture/DECISIONS.md`) + REQUIREMENTS edit + ROADMAP in one commit — verify: `git show --stat HEAD`, AI — expect: REQUIREMENTS.md, DECISIONS.md and ROADMAP.md in one commit — evidence: 4df050f
- [x] Spec `.agents/docs/specs/sh-agent/metrics.md` (agreed by the operator, D30) — docs: specs/sh-agent/metrics.md. The R10 rows in `memory-loop.md` § Known differences move to the Phase 3 spike, which produces them.
- [x] Checkpoint — expect: requirements, ADR and spec committed and pushed — evidence: `go test ./internal/harness` exit 0 (REQUIREMENTS text), commits 4df050f + Phase 0 close, 2026-10-08
### Phase 1 — Store
- [x] Exports in `internal/search`: `PrepareDBFiles`, `DSN`, `IsBusy`, `OpenCopy`, `KindOf`, `ProtoFields` (`internal/search/export.go`); `TrailDates` moves in Phase 4 to `internal/trail` (D31); the AGENTS.md and package-comment driver lines updated in the same change (restatement rule) — verify: `go test -race ./internal/search ./internal/dream` — expect: existing tests unchanged and green; a `ProtoFields` table test over a hand-built protobuf
- [x] `internal/metrics/store.go`: Open, migrations, salt, refusal of a newer schema — verify: `go test -race ./internal/metrics -run 'Store|Migrat'`, AI — expect: symlink refusal, newer-schema refusal (bytes unchanged) and v1 row survival pass, and fail with the migration loop removed
- [x] `export` (`Store.Export`; the CLI verb is wired with `ingest` in Phase 2) — docs: spec — verify: `go test -run Export` — expect: JSONL round-trips every table; refuses an existing file
- [x] Checkpoint — expect: `go test -race -timeout 30m ./...` exit 0 — evidence: `go test -race -timeout 30m ./...` exit 0 (all 7 packages ok), gofmt and vet clean; migration-loop mutation on a copy turned 3 tests red; 2026-10-08
### Phase 2 — Claude ingester (first, because transcripts expire)
- [x] `verb.go` `Verb()` table tests: chains, pipes, env prefixes, `go run ./cmd/sh-agent search --refs x`, quoted `;` inside a string, `git add .` (D22) — verify: `go test -run Verb` — expect: every case passes
- [x] `claude.go`: all-projects walker incl. `subagents/`, `.meta.json` agent type, usage (max per message id), the tool row plus result UPSERT, denied/interrupt/compaction markers (D14), `cost-state` (last line wins), offset, the transcript-root rule (D9) — verify: `go test -run Claude` over `testdata/claude/`, with `CLAUDE_CONFIG_DIR` set to a decoy dir in the test — expect: exact row counts; re-ingest adds 0; a replaced file adds no duplicates; the decoy is never read; canary absent (with the must-fail control)
- [x] `cmd/sh-agent` `metrics ingest` (Claude) and `metrics export` — docs: README verb table, usage text in `cmd/sh-agent/main.go`, spec — verify: `env -u CLAUDE_CONFIG_DIR go run ./cmd/sh-agent metrics ingest --home .agents/scratch/metrics-home` — expect: the summary line, exit 0 — evidence: a copy of the real transcripts ingested into a scratch home: 85 sessions, 3667 steps, 1837 tool calls equal to an independent Python count (0 differences); second ingest 0 rows; every stored program and flag audited as name-shaped (D33); copy deleted, 2026-10-08
- [x] `report` Skipped lines and the never-sent test; the uninstall-keeps test (D17) — docs: README state table — verify: `go test ./internal/harness -run 'Report|Uninstall'` — expect: pass, and fail with the Skipped entry removed
- [ ] Operator runs the first real ingest — verify: user, `sh-agent metrics ingest` after install — expect: a summary line with N>0 (handoff Next Steps if not done)
- [x] Checkpoint — expect: full test exit 0 — evidence: `go test -race -timeout 30m ./...` exit 0 (7 packages ok), gofmt and vet clean; mutation controls on copies: positionals kept → canary test red; heredoc stripping removed → 2 tests red; uninstall deleting the store → test red; 2026-10-08
### Phase 3 — Antigravity ingester
- [x] Spike (agy, read-only, on copies) — the known-difference rows are in memory-loop.md: the tool-name and command-line field numbers in `steps.metadata` for each tool type in the local DBs; how a subagent appears (a `has_subtrajectory` step, `parent_references`, or a separate conversation); whether a compaction or interrupt leaves a step; `steps.status` values for an error and a denial — verify: AI — expect: a D-row per item, with field numbers and a fixture, or a named R10 known difference, written into `specs/sh-agent/memory-loop.md` § Known differences — docs: memory-loop.md
- [x] Shared tool vocabulary (kept as `agyTools` in `agy.go`, the one place a vendor name is translated; Claude names are the vocabulary): the vendor tool name maps to a common name (`view_file`→Read, `run_command`→Bash, `grep_search`→Grep, …), and unknown names are kept as `agy:<name>`. A skill load is a Skill call (Claude) or a read of `skills/<name>/SKILL.md` (agy; only `<name>` is stored, D22) — verify: `go test -run Tools` — expect: both vendors' fixtures give the same common names
- [x] Check the D3 token mapping (D3a: paths fixed by the second spike; the store's generation and token rows match a direct count of the copied databases, 196 generations; whether input includes cached tokens stays unconfirmed, recorded as best-effort in memory-loop.md): sum fields 2, 3 and 5 for one conversation, and compare with `conversation_summaries` or the agy UI — verify: AI + user — expect: a match recorded as a D-row, or a D-row that sets agy tokens to NULL
- [x] `agy.go`: `CopyDB`, then read `conversation_summaries`; `steps` type and status per step (type 132: the tool name, and a `run_command` command line normalised by `Verb` in memory, read from `steps.metadata` via `search.ProtoFields`; `step_payload` never read, D28); `gen_metadata` model and tokens per D3 — docs: spec, README (`--vendor agy`) — verify: `go test -run Agy` over a synthetic DB in `testdata/agy/`, with the source directory chmod 0o500 as in `internal/search/dream_test.go:360-377` — expect: rows match; no `-wal`/`-shm` appears beside the source — evidence: `TestAgyIngestShape` passes; mutation opening the live file in place → red (both conversations fail in the 0o500 directory); real copy: 190/11/194/196 tool/user/reply/generation steps equal to direct counts, second ingest 0
- [-] One real Antigravity session, then `metrics ingest --vendor agy` + `show` on a scratch home holding a copy — verify: AI — expect: the session appears with its step count — dropped here: moved to § Sandbox tests, which runs a fresh session per CLI with `show` and `report` once they exist (2026-10-08)
- [x] Checkpoint — expect: full test exit 0 — evidence: `go test -race -timeout 30m ./...` exit 0 (7 packages ok; dream and memory cached, unchanged), gofmt and vet clean, 2026-10-08
### Phase 4 — Events and show
- [x] `events.go` (with `TrailDates` moved to `internal/trail`, D31): git commits per root and window (D19), `search.TrailDates` for handoff/testing/reviews/bugs, verifier `## Status:`, review severity lines, `feedback.jsonl`, bug↔commit link (D19), ambiguous overlap — verify: `go test -run Events` with a git fixture repo built in `t.TempDir()` — expect: the right event counts; overlap marked ambiguous
- [x] `show.go` with A–D, the label filter and `--json` — docs: README verb table, spec `## Output` — verify: `go test -run Show` golden output — expect: matches the fixture's hand-computed values
- [x] E2E golden path and canary edge case (Verification) — verify: AI, scratch home, copy deleted after — expect: counts equal; canary absent; second ingest 0 — evidence: copies of both CLIs' real data, 100 sessions; `show --json` `B.tool_calls_by_tool` equal to independent counts (claude 1899, agy 190); second ingest 0; 91 events joined; canary absence is pinned by `TestClaudeIngestStoresNoContent` and `TestAgyIngestShape` with must-fail controls; the E2E found active > wall on agy (subagent time counted), fixed under `TestShowActiveTimeMainOnly`; 2026-10-08
- [x] Checkpoint — expect: full test exit 0 — evidence: `go test -race -timeout 30m ./...` exit 0 (8 packages; `internal/trail` has no test file of its own: covered by dream's `StatusArgs` test and `TestJoinEvents`), gofmt and vet clean, 2026-10-08
### Phase 5 — Wiring (ONE commit: code, docs and map, rule 3; D20)
- [x] `dream` realRun ingests inside the lock, dry run does not (D16) — docs: README dream row, spec dream.md — verify: `go test ./internal/dream -run Ingest` — expect: ingest runs before collect; dry run writes no `metrics.db`; an injected failure leaves the exit status and `last_run.error` unchanged
- [x] README map label `proposals/dream/` → `proposals/dream/, metrics.db`, plus a state-table row — docs: README.md — verify: `go test ./internal/harness -run TestHarnessMap` — expect: pass
- [x] `ab-kit/run.sh` calls `metrics ingest --home "$AB_OUT" --source <home> --label arm=<arm>` after the run and prints the `--outcome` command for after `accept.py` (D24, D36); `metrics.py` removed — docs: testing/guide.md § ab-kit — verify: dry run of the kit to the ingest line — expect: rows under the label — evidence: `bash -n` ok; the exact ingest command is exercised in § Sandbox tests (a real kit run costs tokens and was not run); the late `--outcome` stamp is pinned by `TestOutcomeStampedAfterTheRun`
- [x] (AGENTS.md driver-importer lines done in Phase 1) testing guide § 3g (new, not § 3e) and the AGENTS.md pointer; sh-agent-matrix rows are written from the § Sandbox tests evidence — docs: AGENTS.md, testing/guide.md
- [x] Checkpoint — expect: full test exit 0; `sh-agent status` binary line with no warning from a scratch-home install — evidence: `go test -race -timeout 30m ./...` exit 0, gofmt and vet clean, 2026-10-08; the scratch-home `status` check runs in § Sandbox tests (one install for the whole feature)
### Phase 6 — Visual report
- [x] `report.go` in `internal/metrics`: `sh-agent metrics report [--home H] [--since] [--label|--all-labels] [--out FILE]`. It writes one self-contained HTML file: `html/template`, inline SVG drawn in Go, no script and no external fetch, light and dark through `prefers-color-scheme`. Default path: `~/.config/sh-agent/reports/<YYYY-MM-DD-HHMM>-metrics.html`, mode 0600, printed on stdout. An existing FILE is refused. Charts:
  - sessions, cost and tokens per day (line)
  - tool calls by tool (bar)
  - search:grep ratio over time (line)
  - error rate and denials per week (line)
  - subagents by type (bar)
  - pass^k by label (bar)
  - A1–A5 outcome table
  - docs: README verb table, spec `## Report`
  - verify: `go test -run Report` golden HTML over the fixture store; the output holds no `<script` and no `http`, and every number in it matches `show --json`
  - expect: pass; the golden test fails when a chart's data is altered
- [x] Render the E2E store's report and open it in a real browser (styrman-browser): both themes, plus a 400 px width — verify: AI — expect: screenshots in `.agents/scratch/`, no console errors, the charts are labelled — evidence: report from copies of both CLIs' real data in styrman-browser at 1100 px light and dark and at 400 px; screenshots `.agents/scratch/report-{light,dark,400}.png`; no console messages. The first render showed three defects, each fixed under `TestReportDrawsNoMissingValues` (red, then green): legend swatches sanitised to `ZgotmplZ`, Antigravity drawn as 0 cost and as a 0 ratio on days with no grep, and colliding end labels. Small ticks at 400 px are recorded as a known limit; 2026-10-08
- [x] Checkpoint — expect: full test exit 0 — evidence: `go test -race -timeout 30m ./...` exit 0, gofmt and vet clean, 2026-10-08
### Phase 7 — Feedback into the learning loop (D27)
- [x] Propose the R3 addendum: dream also reads "a summary of the metrics store (`metrics show --json` for the project's window)" — verify: operator's approval in their own turn, user — expect: approved text quoted in a D-row and committed with an ADR-055 amendment row in DECISIONS (append-only: ADR-055a)
- [x] Dream's prompt gains a `metrics` section: that project's `show --json` since its watermark, plus the previous window for comparison. Only numbers go in, and the existing redaction and 256 KB cap apply, with metrics dropped before transcripts. The answer schema accepts a proposal whose `evidence` names a metric — docs: spec dream.md, README dream row — verify: `go test ./internal/dream -run Metrics` — expect: the section appears with numbers only; a dry run prints it; a cap test drops it before transcripts — evidence: `TestPromptCarriesMetricsNumbersOnly` and `TestMetricsDroppedFirstAndNeverChanges` (both red before the code), `TestNumbersOnly` (its must-fail control on a copy, leaking tool names into keys, went red), `TestSnapshotReadOnly`; the dry run lists the `metrics` input through the shared `collect`; 2026-10-08
- [x] `skills/retro/SKILL.md` (Step 1 source 11; the written change — its behavioural runs are in § Sandbox tests): the orientation step runs `sh-agent metrics show --json --project .` and names regressions against the previous window as candidate items — docs: the skill itself — verify: a behavioural run on both CLIs (testing guide § 4) — expect: the retro artifact cites at least one metric
- [x] Checkpoint — expect: full test exit 0 — evidence: `go test -race -timeout 30m ./...` exit 0 twice (before and after the no-store refusal), gofmt and vet clean; the later `sh-agent` second-level-verb fix passed `go test -race ./internal/metrics` and is covered by the Definition of Done verifier run; 2026-10-08
### Sandbox tests (D25)
- [x] Testing guide § 2: a scratch-home install, with `sh-agent status` showing a clean binary line — verify: AI — expect: exit 0, no warning — evidence: `drift none — 1789 files match the repo`, binary line without a warning, `testing/2026-10-08-verify-metrics-sandbox.md`
- [x] § 4 behavioural runs, kit-preamble first: on Claude Code and on Antigravity, one session in a sandbox home that runs tools, then `metrics ingest --home <sandbox>` + `show` + `report`, plus one `/retro` run per CLI after Phase 7 — verify: AI, reported as n of N per CLI in `testing/` — expect: every run's session appears with its tool count (Claude) or step count (agy); the report renders — evidence: `testing/2026-10-08-verify-metrics-sandbox.md`: every metrics check 2 of 2 (Claude Code) and 1 of 1 (Antigravity), both retros read the metrics (source 11); the skill-matrix `retro` row check is 0 of 2 on Claude (project-scope decision under Step 2; a control at 2c0f215 is Open) and 1 of 1 on Antigravity
- [ ] Real-home smoke (§ 5) after the operator's install — verify: user — expect: `sh-agent metrics ingest` exit 0 (handoff Next Steps if not done)
### Done
- [ ] Definition of Done: `reviewer` (`--code --docs`) and `security` over the whole diff (a transcript parser feeding a persistent store; input parsing), then `verifier`, then the changes entry, plan closed, commit, push, `/handoff` — verify: artifacts under reviews/ and testing/ — expect: Outcome sections closed, verifier PASS

## Risks
Assumptions:
- Antigravity's usage-shaped message (fields 2, 3 and 5 under `gen_metadata`) is input, output and cached tokens. Checked by: the first task of Phase 3. If it cannot be confirmed, the token columns stay NULL for agy (R10 known difference).
- Claude Code keeps writing `cost-state` at session end. Checked by: the Phase 2 fixture comes from a current transcript. A running or crashed session has no cost line, so `cost_usd` is NULL and `show` notes it.

## Decisions
- D1 — 2026-10-08 — SQLite is the store of record at `~/.config/sh-agent/metrics.db`, separate from `search.db`; collection is a new verb plus `dream --all`; rows are shape only; nothing goes through `report` — by: operator, plan-mode answers 2026-10-08 ("SQLite as store of record", "New verb + nightly dream", "Shape only", "Local only for now").
- D2 — 2026-10-08 — Claude cost and duration come from the last `cost-state` line of a session (`totalCostUSD`, `totalDuration`; cumulative, written at session end or resume). No price table. Subagent cost is not separate (it is in the parent's total) — by: implementer (spike 2026-10-08: 1–4 lines per ended main session, 0 in a running one, none in subagent files).
- D3 — 2026-10-08 — Antigravity tokens are read from `gen_metadata` by field number (path ending 2 = input, 3 = output, 5 = cached; 4 under the parent = context window 256000), pinned by a fixture and reported as best-effort — by: implementer (spike; no schema published; confirmed or dropped in Phase 3).
- D3a — 2026-10-08 — Refines D3 with the exact paths from the second spike (one 18-generation conversation, read-only copy). One `gen_metadata` row is one model call. Its fields:
  - 1→4→2 input tokens, 1→4→3 output tokens, 1→4→5 cached tokens (absent when 0);
  - 1→19 the model name;
  - 1→9→10→1 the context used and 1→9→10→4 the window (256000), which gives B3 peak context on Antigravity.

  Each generation is stored as a `steps` row of kind `generation` (`ext_id` = `<conversation>:gen:<idx>`). Its time is that of the i-th model reply (step type 15), because the spike's conversation had 18 of each; with no such reply the time is null. Whether input includes the cached tokens is unconfirmed; the spec's known-difference row keeps "best-effort" — by: implementer (spike 2026-10-08).
- D4 — 2026-10-08 — New package `internal/metrics` as the second driver importer — by: implementer.
- D5 — 2026-10-08 — Forward-only migrations; a store is never deleted on a version bump (the opposite of `search/store.go:112-136`), because it is the only copy — by: implementer, from D1.
- D6 — 2026-10-08 — `cmd_hash` is salted with a per-store salt kept in `meta`. It allows equality detection only. It does not stop brute force by someone who can read the file, which is acceptable for a local 0600 file — by: implementer.
- D7 — 2026-10-08 — `ingest` reuses the global `--vendor` (`cmd/sh-agent/main.go` `c.fs.StringVar(&c.vendorFlag, "vendor", …)`). Empty means all — by: cold read round 1 (G2).
- D8 — 2026-10-08 — Seams: `metrics` imports `search` (exported `PrepareDBFiles`, `CopyDB`, `KindOf`, `TrailDates`, `ProtoFields`); `dream` imports `metrics`; `search` imports neither — by: cold read round 1 (G3).
- D9 — 2026-10-08 — Transcript root: `--source` > `--home` > `CLAUDE_CONFIG_DIR` (only with neither flag) > the real home. This closes the R6 leak through the environment variable — by: cold read round 1 (G7).
- D10 — 2026-10-08 — Step identity is `(session_id, ext_id)` per the Interfaces rules: max usage per message id, tool_result folded into its tool row, times in unix ms — by: cold read round 1 (G5).
- D11 — 2026-10-08 — A6 pass^k: a run is one `--source` ingest, with `sessions.label` + `sessions.outcome` from `--outcome`; pass^k for a label is 1 when all k runs pass — by: cold read round 1 (G1).
- D12 — 2026-10-08 — Path columns are limited to `sessions.project_root` and `ingest_state.source`; `sessions.source` is dropped — by: cold read round 1 (G1).
- D13 — 2026-10-08 — `events` payloads are fixed per kind (Interfaces) and hold no prose — by: cold read round 1 (G1).
- D14 — 2026-10-08 — Denials, interrupts and compactions are detected by fixed markers. Only a flag or a row kind is stored — by: cold read round 1 (G1).
- D15 — 2026-10-08 — "metrics record" keeps R8's meaning; this is "the metrics store" — by: cold read round 1 (G11).
- D16 — 2026-10-08 — Dream ingests inside the lock, store-wide, never in a dry run, and never fails the dream — by: cold read round 1 (G10).
- D17 — 2026-10-08 — `report` names `metrics.db*` as Skipped; `uninstall` keeps it; both pinned by tests — by: cold read round 1 (G11).
- D18 — 2026-10-08 — (superseded by D28) Antigravity step 132's payload is never decoded (the 2026-09-30 security decision at `internal/search/dream.go:203-206` stands). Tool steps are counted from the `step_type`/`status` columns, so agy has no tool names: C1, C2, C4 and per-tool B2 are Claude-only (R10 known difference) — by: cold read round 1 (G4).
- D19 — 2026-10-08 — A1: the main session's window is `[started_at, ended_at + 30 min]`; closed = at least one commit AND a handoff file whose `TrailDates` time is in the window. A4: a bug record is linked to a session when its text contains a ≥7-hex prefix of one of that session's commits, and only `BUG-NNN` + the SHA are stored. The project root is the first line's `cwd`. A2/A3 are counted for main sessions only — by: cold read round 1 (G6).
- D20 — 2026-10-08 — Rule 3: the README map label lands in the Phase 5 commit together with the dream wiring that creates the path. ADR-055 (Phase 0) says so — by: cold read round 1 (G11).
- D21 — 2026-10-08 — The R3/R4/R7 edits are approved as drafted in `.agents/scratch/2026-10-08-metrics-requirements-edit.md` — by: operator, "go all approved" (2026-10-08).
- D22 — 2026-10-08 — Q1: under "shape only", exactly three identifiers are stored: the Skill tool's skill name, `subagent_type`, and a bare `.` after `git add` — by: operator, "go all approved" (2026-10-08).
- D23 — 2026-10-08 — Q2: `metrics export` is kept — by: operator, "go all approved" (2026-10-08).
- D24 — 2026-10-08 — Q3: `ab-kit/metrics.py` is replaced by `metrics ingest --source --label --outcome`, and `accept.py` is kept — by: operator, "go all approved" (2026-10-08).
- D25 — 2026-10-08 — Sandbox tests are owed for this feature (testing guide § 2 and § 4: a scratch-home install, one sandbox start per CLI, a behavioural run of `metrics ingest`/`show` per CLI, reported as n of N) — by: operator, "Make sure to run sandbox tests" (2026-10-08).
- D28 — 2026-10-08 — Supersedes D18. Antigravity tool names and `run_command` command lines are read from `steps.metadata` (the tool-call column), never from `step_payload` (the output column that the 2026-09-30 decision at `internal/search/dream.go:203-206` keeps unread). Command lines are normalised by `Verb` in memory, exactly as on Claude; nothing is sent and nothing beyond the shape is stored. So C1, C2, C4 and per-tool B2 hold on both CLIs. The security review at Done checks the decoder — by: operator, "this needs to work for both claude and agy" (2026-10-08), with the probe of `steps.metadata` (tool name `view_file` present).
- D29 — 2026-10-08 — R10 parity rule for this feature: every metric, chart and verb works on both CLIs. A metric that one CLI's runtime cannot produce is listed in `specs/sh-agent/memory-loop.md` § Known differences with the probe that showed it; until the Phase 3 spike, the candidates are cost (Antigravity exposes none), and subagents, compactions and interrupts. `show` and `report` split every table and chart by vendor, and show "n/a (CLI)" rather than 0 — by: operator (2026-10-08), R10.
- D31 — 2026-10-08 — Supersedes D8 for `trailDates` only. It moves to a new driver-free package `internal/trail` (with `trailStatusArgs`), not into `search`, which holds the index. `dream` and `metrics` both import `trail`. Everything else in D8 stands — by: implementer (Phase 1).
- D32 — 2026-10-08 — A symlinked `~/.config/sh-agent` DIRECTORY is not refused, the same as search (a dotfiles setup may symlink it). Only the database and its companion files are refused when not regular files. The spec line saying "or directory" is corrected — by: implementer (Phase 1).
- D33 — 2026-10-08 — The verb normaliser models here-documents, comments, quotes, `${…}`, `$((…))`, subshells, case patterns, redirections and shell keywords. A program must be name-shaped. A short flag longer than three letters, or holding a digit, keeps only its first letter, because a glued value (`-psecret`) is an argument. Found by the real-data check in Phase 2, where the first tokenizer stored file names, `${GH_TOKEN:-}` and glued flag values as "programs" and flags. That was caught before any commit and recorded in the Phase 2 changes entry — by: implementer (Phase 2).
- D34 — 2026-10-08 — Verbs are normalised at ingest time and never re-normalised; the command is not stored (shape only, D1). A normaliser change applies to new rows only — by: implementer (Phase 2).
- D35 — 2026-10-08 — Antigravity spike (read-only, on copies of 15 local conversation databases):
  - The tool name is at `steps.metadata` path 4→2 (`run_command`, `view_file`, `write_to_file`, `invoke_subagent`, `manage_subagents`, `send_message`). The tool's arguments JSON is at 4→3; only `CommandLine` (run through `Verb`) and `AbsolutePath` (classified by `KindOf`) are read, in memory.
  - A step's time is at `metadata` 1→1, in seconds.
  - Step types: 14 is user input, 15 a model reply, 132 a tool. 101 and 17 hold IDs only and are not stored.
  - Status 7 with a non-empty `error_details` is a tool error.
  - A subagent is its own conversation. `conversation_summaries.agent_name` gives its type, `parent_conversation_id` its parent and `nesting_depth` its depth.
  - A skill load is a `view_file` of `…/skills/<name>/SKILL.md`.
  - Not observed: denials, interrupts and compactions. They are R10 known differences until a sample exists. Cost is not exposed.
  — by: implementer (Phase 3 spike, 2026-10-08).
- D36 — 2026-10-08 — The A/B kit ingests into its OWN store (`--home "$AB_OUT"`), never the real home's, and stamps the outcome after `accept.py` with a second `--outcome` ingest. That ingest sets the outcome on the label's sessions that have none, even when no transcript changed (found while wiring the kit: the offset skip would have dropped a late outcome) — by: implementer (Phase 5).
- D37 — 2026-10-08 — A session id ingested first without a label keeps that state when a later `--source` ingest of the same, unchanged files names a label: the offset skip means no relabelling. Sandbox homes have their own session ids, so this does not arise in the kits; the report test uses a renamed copy (`sandboxCopy`) — by: implementer (Phase 6).
- D38 — 2026-10-08 — `show`, `report` and `export` refuse a home with no store (exit 1, a pointer to `ingest`) instead of creating an empty one that prints zeros. `sh-agent`'s second-level verbs (`ingest show report export add`) are kept by `Verb` (both found by the sandbox runs) — by: implementer (§ Sandbox tests).
- Open: does the Claude Code `retro` matrix row's 0 of 2 (project-scope, no harness file) predate this feature? — default: a control run at 2c0f215 with the same model, before the next release; until then the row is reported as failing, not as passing.
- D30 — 2026-10-08 — The spec is agreed, and the R3 addendum (dream reads a numbers-only `metrics show --json` summary per project) is approved and applied with ADR-055a — by: operator, "go all approved" (2026-10-08, second).
- D26 — 2026-10-08 — Visual reports are a self-contained HTML file written by `sh-agent metrics report`: inline SVG drawn with stdlib, no script, no fetch (R4 stdlib only) — by: operator, "Self-contained HTML file" (2026-10-08).
- D27 — 2026-10-08 — Feedback goes back through the existing learning loop. Dream and `/retro` read `metrics show --json` as evidence; harness-target proposals already leave through `report`. No new map line, D1's "local only" stands for raw rows, and no R8 change — by: operator, "Via dream/retro proposals" (2026-10-08).
- Rejected: aggregates in `report.json` — a new line back on the map, which the operator did not choose (D27).
- Rejected: JSONL record + SQLite cache — the operator chose SQLite as the record (D1).
- Rejected: tables inside `search.db` — that file is deleted on a schema bump.
- Rejected: storing redacted arguments — the operator chose shape only (D1).
- Rejected: `internal/metrics` code inside `internal/search` — it keeps "one importer" true in wording only, and couples two stores.
- Rejected: decoding agy `step_payload` for anything — it holds tool output (2026-09-30 security decision).
