---
owner: reviewer
date: 2026-10-08
---
# Review — metrics-store plan, cold read   (2026-10-08)

Range reviewed: `.agents/docs/plans/2026-10-08-metrics-store.md` (a plan file, not a diff), plus the code surfaces it names: `internal/search/store.go`, `internal/search/dream.go`, `internal/search/parse.go`, `internal/dream/{collect,run,dream}.go`, `cmd/sh-agent/main.go`, `internal/harness/{state,report}.go`, `README.md`, `REQUIREMENTS.md`, `AGENTS.md`, `testing/ab-kit/{run.sh,metrics.py}`.

## Status: COMPLETE

## Verdict: design needs revision / implementation not ready

Worst problem first: the schema in `## Interfaces` cannot compute several catalogue metrics the plan promises (A6 pass^k has no pass/fail source or run identity; C4 skills loaded has no column; C2 `git add .` is undetectable because positional words are dropped; D1 permission denials have no column), while `## Scope` says no paths are stored and the schema stores several. Second, the plan's "Code to reuse" does not survive contact with the code: the helpers are unexported, `proseLeaves` returns prose only, `dream` would have to import `metrics` while `metrics` needs `dream`'s `trailDates` (an import cycle), and `--vendor` is already a global flag on every verb, so redeclaring it panics at start-up. Third, two verify fixtures do not produce the condition they claim (see G8, G9). The 7 BLOCKING questions below must be answered before a first pass.

I could not check the plan's spike claims (`cost-state.totalCostUSD`, `subagents/*.meta.json`) against real transcripts: the probe of `~/.claude/projects` was refused by the permission layer and I did not retry it. They are unverified by me, not disproved. Nothing else in the repo mentions `cost-state` (grep over `.agents/docs`, `internal`, `skills`: only this plan).

## Findings (Gaps, question 1)

- **G1 [high / high]** Schema vs catalogue vs "no paths" (`plan:22-29, 46-67, 86-98`).
  - A6: "pass rate per label, with pass^k" needs a per-run pass/fail and a run identity. The schema has neither (`sessions.label` only; `events.kind` has no `pass`). Pass/fail today lives in `ab-kit/accept.py`, which is not ingested.
  - C4 "Skills loaded" needs the Skill tool's `skill` argument, B4 needs `subagent_type`; no column holds either, and `## Scope` Out says "Arguments... in any row".
  - C2 `git add -A`/`git add .`: `Verb` keeps flag names only and drops positional words, so `git add .` is indistinguishable from `git add file`.
  - D1 "permission denials" and "interrupts": `steps.kind` has `interrupt`, but nothing for denial; how either is detected (text match on tool_result prose) is undecided.
  - "Out: file paths in any row", yet `sessions.project_root`, `sessions.source`, `ingest_state.source` are paths, and `events.ref`/`events.value` are untyped. A5 "corrections" from `feedback.jsonl` and A3/A4 events could hold prose or file:line unless stated. Fix: a D-row listing exactly what each `events` kind stores in `ref`/`value`, and which path-valued columns are permitted (the canary test would also fail if the canary sat in a directory name).
- **G2 [high / high]** `--vendor` is already registered for every verb: `cmd/sh-agent/main.go:~203` (`c.fs.StringVar(&c.vendorFlag, "vendor", ...)`, values `all, antigravity (or agy), claude`, resolved against the manifest by `resolveVendor`). The plan's verb flag `--vendor claude|agy|all` on `ingest` would call `StringVar` twice on one FlagSet, which panics ("flag redefined: vendor"). Decide: reuse the global flag, and say what an empty value means (manifest's vendor, as `resolveVendor`, or all). `ctx` is one flat struct (`main.go:~100-108`), so each new flag is a field there.
- **G3 [high / high]** The reuse list is wrong, and the dependency direction is unspecified.
  - `internal/search`: `tryOpen`, `prepareFiles` (store.go:~144-171), `openCopy`, `copyDB`, `kindOf` (parse.go:270), `proseLeaves`/`walk` are all unexported. The plan says only "exported if needed" for `kindOf`. Two driver-importing packages either duplicate the security-reviewed `prepareFiles` or share a helper; undecided.
  - `walk` (search/dream.go:~316-403) returns text leaves only, discarding field numbers. D3 ("tokens read by field number") and Phase 3 ("type 132 tool name only") need a different decoder returning (field, value); not in the plan or the spike.
  - `claudeTranscript` (dream/collect.go:355-433) lists only `projects/<dir>/*.jsonl` and filters per line by `cwd` under one root; it does not walk `<session>/subagents/*.jsonl` and is not usable for an all-projects ingest. `trailDates` (collect.go:207) is in package `dream`, which `dream/run.go` would make import `metrics`; `metrics` importing `dream` is a cycle. These helpers need a home (an exported package, or `metrics` owns copies).
- **G4 [high / medium]** Reverses a recorded security decision without a security step. `search/dream.go:203-206` reads only step types 14 and 15 and says 132 "holds file and command output ... (security review 2026-09-30)". The plan reads 132 for the tool name. Nothing in the phases dispatches `security` (global DoD: input parsing, secrets; also a transcript parser feeding a persistent store), nor the reviewer, verifier, changes entry or handoff. The step-132 tool-name field was not in the Phase 0 spike (the spike lists only `gen_metadata`).
- **G5 [high / medium]** Idempotency model is inconsistent. `unique(session_id, ext_id)` with `ext_id` = "tool_use.id or message id": a tool_use step and its tool_result share the same `tool_use_id`, so the unique index collides. `steps` PK is `(session_id, seq)`; after "size < offset, re-read from 0" a replaced file reassigns `seq`, so PK conflicts or silently drops rows even with ext_id dedupe. Streaming assistant messages repeat one message id across lines (the reason `ab-kit/metrics.py` dedupes by `m['id']`); which line's `usage` wins is not stated. The unit of `ts`, `started_at`, `ended_at` (s or ms) is not stated. For agy, `ext_id` and the incremental offset of a SQLite source are not defined (`ingest_state.offset` for a DB).
- **G6 [high / high]** Metric definitions that two implementers would compute differently: A1 "commit plus handoff in the session window" (handoff is committed after the session ends in the DoD; window = started_at..ended_at?); A4 "BUG records dated after a session commit" (any bug after any commit is meaningless without a link; bug date is a day, commit time is seconds); "ambiguous" overlap "counted separately" (what does `show` print?); project root choice when `cwd` changes within a session (collect.go filters per line); A2/A3 window for subagent sessions.
- **G7 [high / high]** R6 / sandbox hazard with `CLAUDE_CONFIG_DIR`. `claudeTranscript` gives the env var priority over `<home>/.claude` (collect.go:358-361). `metrics ingest --home .agents/scratch/metrics-home` and `--source DIR` would read the operator's real transcripts whenever the variable is set. The dream tests defeat this with `os.Unsetenv("CLAUDE_CONFIG_DIR")` (`internal/dream/main_test.go:9-16`); the plan's Verification does not, and `--source` makes it a product bug, not a test bug. Decide: `--source`/`--home` ignore the env var, or the env var applies to the real home only. Also the golden path copies real transcripts (may hold secrets) into `.agents/scratch`; say so and delete it after.
- **G8 [medium / high]** Phase 3 verify fixture: "the live DB path is never opened (the test makes it unreadable)". Copy-then-read must `open(2)` the source to copy it; a file at mode 0 fails that with EACCES (and root bypasses it), so the test would fail the copy rather than prove the driver never opened the live path. The existing proof is different: `TestCopyDBNeverOpensSource` (`internal/search/dream_test.go:360-377`) makes the source DIRECTORY read-only (0o500), so SQLite cannot create `-wal`/`-shm` beside the original, and the copy succeeds. Use that.
- **G9 [medium / high]** Verification names implementation choices and under-specifies conditions.
  - Golden path: "an independent Python count" of the SQLite-verb count needs a second implementation of `Verb` (the plan's own baseline says "verbs counted at command position"); a mismatch cannot be told from a normaliser difference. `--json` shows `"A":{…}` with no key names, so "reports the same tool-call count" has no key to read. Whether subagent calls are in the count (baseline 1,619) is not stated.
  - Canary: where it sits is unstated. In a positional word it is dropped; as a flag name (`--SECRETCANARY123`) `Verb` keeps it by design; in `cwd` it lands in `project_root`. The test must place it in each field the ingester reads and prove it can fail (a deliberate leak makes it red) — silence is not evidence.
  - "Migration: a binary carrying a v2 test migration" is an implementation choice (an injectable `migrations`); state the criterion (a v1 store with rows opens under a newer schema list and keeps every row, and a newer store is refused untouched).
- **G10 [medium / high]** `dream` wiring. `Run` splits at `dream.go:69-80`: `realRun` or `dryRun`. The plan says `--all` and `--project` ingest first, without saying a `--dry-run` must not (global rule: read-only mode resolves the same guards and performs no side effect, and its test asserts that). Ingest is store-wide, but `--project DIR` implies one root; undecided. Whether the ingest sits inside `dream.lock` (run.go:~53, taken before `o.targets()`) or before it: undecided, as is whether an ingest failure counts toward `last_run.error`.
- **G11 [medium / high]** Map rule 3 and the delivery files. AGENTS.md rule 3: a change to a path lands "in ONE commit with the code, the docs and the map". The plan puts ADR-055 and the REQUIREMENTS edit in Phase 0, the README map label in Phase 5, the code between. Also: `report` collection is an explicit list and names `search.db*` in `Skipped` so a dry run shows it (`internal/harness/report.go:~424-430`); the plan adds `metrics.db` and `-wal`/`-shm` without a "never sent" test or a `Skipped` line. "Metrics record" already means `report.json` in R8/R9 and `delivery.md:54`; the plan's new `metrics` verb and store reuse the word without a sentence that separates them. `uninstall` (R4: removes exactly what its manifest lists): what it does with `metrics.db`, the only copy (D5), is not decided. The ADR file is `.agents/docs/architecture/DECISIONS.md` (the plan says `DECISIONS.md`; ADR-054 is the last, so 055 is free).
- **G12 [medium / medium]** R10 parity. Antigravity has no subagents, compaction records or cost in the plan's sources; B1 tokens may be Claude-only. R10 says what a CLI's runtime alone controls is listed in `.agents/docs/specs/sh-agent/memory-loop.md` § Known differences; the plan puts it in the new spec only, and lists none of B3, B4, cost as Claude-only. The Risks claim "checked by ... first task of Phase 3" but Phase 3 has no such task, and Verification's "one Antigravity session run" tier has no phase task either.
- **G13 [medium / medium]** Spike evidence. Phase 0 spike is ticked with "field list recorded in D2, D3", but D2 and D3 are one line each and no probe artifact is named; `cost-state` semantics (cumulative vs per-line, per subagent, reset on resume) are not stated and the repo's existing metric script takes cost from the CLI result JSON, not `cost-state` (`ab-kit/metrics.py:3, 14`). Risks lists "Claude keeps writing cost-state" as an unchecked assumption while D2 treats it as a spike result.
- **G14 [low / high]** Docs clauses. `export` (Phase 1) and `--vendor agy` (Phase 3) change the command surface with no `docs:` clause; the `usage` text at `cmd/sh-agent/main.go:~56` and its tests are not listed. `README.md` verb table row, the state table rows 217-218 and `delivery.md`/`dream.md` appear only in some tasks.
- **G15 [low / medium]** Concurrency and failure edges left open: "the second waits or names the lock" (pick one); an Antigravity copy failure makes the whole source fail while "the others are still ingested" (per source, per conversation?); whether a non-git root or a deleted root still gets session rows; a store written by a newer binary on one place and read by an older `show`.

## Scope (question 2)

Built but not named in `## Why` (the request: use SQLite for trajectory data and agent/sandbox metrics; suggest metrics):
- `metrics export` and its JSONL shape: a second output shape beside `show --json`, no D-row says who asked (the only operator answers quoted in D1 are "SQLite as store of record", "New verb + nightly dream", "Shape only", "Local only for now"). A dated nightly export is already a follow-up, so the verb has no consumer in `## Consumers`. Gap unless a D-row or strike.
- Replacing `ab-kit/metrics.py` and `run.sh` (Phase 5) is a consumer in the plan (Consumers 3) but not in D1; acceptable if the operator approves it as a D-row, since it changes a verification kit.
- `--project`, `--since`, `--all-labels` and `--label` filters on `show`: a plausible reading of "metrics from all sessions"; low risk, but each is surface. `meta.salt` and `cmd_hash` serve D2 (retry loops) only; D6 records the choice.
- The A-group joins (git, trail, feedback) go beyond "trajectory data" but are the "metrics for agent evaluations" the operator asked for suggestions on; the operator approved the plan's shape in plan mode per D1, so I do not count them.

## Readiness (question 3)

BLOCKING (two careful implementers would build different things):
1. Where does pass/fail and run identity for A6 pass^k come from? (G1)
2. How are C2 `git add .`, C4 skill names, B4 subagent type and D1 denials stored given "no arguments"? Which exact columns/`events` payloads, and which path columns are allowed? (G1)
3. `--vendor`: reuse the global flag (and its meaning) or a new name? (G2)
4. Which package owns the shared SQLite open rules, the agy copy, the protobuf field decoder, `kindOf` and `trailDates`, given `dream -> metrics`? (G3)
5. `ext_id`/`seq` rule for tool_use vs tool_result rows, replaced files, streamed message ids, and the units of every timestamp. (G5)
6. Definitions of A1 window, A4 (which bug belongs to which commit), and what "ambiguous" does to the numbers. (G6)
7. Do `--home` / `--source` ignore `CLAUDE_CONFIG_DIR`? (G7)

NON-BLOCKING (a reasonable implementer decides, or a default can be stated):
- Fixture for "never opens the live DB" (use the read-only-directory pattern) (G8); canary placement and JSON key names for the golden path (G9) (a reasonable implementer fixes these, but the plan's own acceptance should name them).
- Dry-run behaviour of dream ingest; lock placement; failure-counting (G10).
- Single-commit rule for the map; `Skipped` line and "never sent" test for `metrics.db`; `uninstall` behaviour (G11) (ask the operator in Phase 0 when proposing R3/R4/R7).
- R10 known-differences location; the missing Phase 3 check task (G12).
- Spike evidence and `cost-state` semantics (G13).
- Docs clauses, usage text (G14).
- Concurrency wording; per-conversation failure (G15).
- `export`: keep or strike (operator decision).
- No `security`, `reviewer`, `verifier`, changes entry or handoff step in the phases (G4).

## Check against REQUIREMENTS.md and AGENTS.md

- R3: the dream exception says it "writes only under `~/.config/sh-agent/`" and runs a CLI over what it read. An ingest step that writes a persistent store from both CLIs' transcripts is a new write the exception text does not name; the plan lists the R3 edit (Phase 0). OK as proposed, pending operator approval.
- R4: "one dependency ... used only by the search index"; D4 makes `internal/metrics` a second importer; the plan lists the R4 edit and the AGENTS.md Layout and Rules lines (`AGENTS.md:45, 63`). Not a violation if approved first; Phase 0 orders it correctly. R4 also says search indexes are caches, never sent: `metrics.db` is a store of record, so the edit must say it is not a cache and is never sent (G11).
- R6: ingest and tests use `--home`, but G7 shows the env-var path to the real transcripts.
- R7 edit listed. R8/R9: the "metrics record" wording collision (G11).
- README map rules 1-6: rule 1 (no new path back): met, nothing goes through `report`, but no test pins it. Rule 3: ADR and the map must share one commit with the code (G11). Rule 4: the dream→configDir label gains `metrics.db`; the plan keeps the line solid, and it works for every place only if the dream timer or `metrics ingest` runs there; the first real ingest is a user task in Phase 2. Rule 5: `TestHarnessMap` is named in Phase 5. Rule 6: the plan does not restate the delivery claim.
- Driver/dependency: no new dependency; the existing driver only. OK.

## Criteria (applied to a plan)

1. Dead code: pass. The plan removes `ab-kit/metrics.py` and adds nothing it leaves unused; `kindOf`'s export is conditional.
2. Flow tracing: fail. Schema vs catalogue (G1), reuse vs code (G3), flag collision (G2), `ext_id` rule (G5), dream dry-run (G10).
3. Standards conformance: fail. Rule 3 single commit (G11); DoD steps absent (G4); R10 parity location (G12); docs clauses (G14).
4. Directory layout: pass. `internal/metrics/` with tests beside code, spec under `specs/sh-agent/`, ADR at `architecture/DECISIONS.md`.
5. Simplicity (DIY-except-security): mostly pass. Stdlib plus the existing driver; `sha256` with a random salt is stdlib use, not DIY crypto; `export` is a second shape without an operator request (Scope).
6. Decoupling: fail until G3 is decided; `dream -> metrics -> (dream|search internals)` has no stated seam, and two packages would duplicate or reach into each other's unexported code.

## Not examined

- Real Claude/Antigravity transcripts and conversation databases (the probe was refused by the permission layer; the spike's field claims are unverified by me, and the agy step-132 payload and `gen_metadata` field numbers were not read).
- `internal/search/store.go` beyond lines 100-180 and `internal/dream/collect.go` lines 1-195 and 470-end; `internal/dream/state.go`, `timer.go`, the `testing/guide.md` § 3e text, `sh-agent-matrix.md`, `kit-preamble.sh`, `accept.py`, `specs/sh-agent/dream.md`/`search.md`/`memory-loop.md` bodies; `TestHarnessMap`; the remainder of `run.sh` after line 50.
- No code was run.

design needs revision
implementation not ready
