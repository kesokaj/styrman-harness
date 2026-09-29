---
owner: reviewer
date: 2026-09-29
---
# Cold-read — `.agents/docs/plans/2026-09-29-search-index.md`   (2026-09-29)

Range reviewed: the plan file `.agents/docs/plans/2026-09-29-search-index.md` and the code surfaces it names (`cmd/sh-agent/main.go`, `internal/harness/state.go`, `internal/harness/report.go`, `internal/harness/size_test.go`, `REQUIREMENTS.md` R4, `skills/handoff/SKILL.md:154`, `.agents/docs/architecture/DECISIONS.md`, `README.md` § How it works). No git diff — this is a design-document read, per the dispatch. The "R4 approved by the operator" gate is treated as approved on D1's text, per the dispatch note.

## Status: COMPLETE

## Verdict: design needs revision (2 BLOCKING gaps), implementation not ready

Both blocking gaps are in the FTS5/query design (`internal/search`, Phase 1), the newest and least-checked surface. Everything else the plan touches (main.go's verb switch, `reorder`/`takesValue`, `state.go`, `report.go`'s explicit `Collect` allowlist, the ADR-heading regex against real `DECISIONS.md` text, the handoff skill's fence-scan line, the context-budget arithmetic) was verified against the actual code and text and matches what the plan claims.

## (1) Gaps

- **BLOCKING — `chunks_fts` external-content sync is undecided.** Interfaces § Storage declares `CREATE VIRTUAL TABLE chunks_fts USING fts5(heading, body, content='chunks', content_rowid='id', ...)`. An `external content` FTS5 table does **not** update itself when `chunks` is written — SQLite requires either triggers on `chunks` or explicit `INSERT INTO chunks_fts(chunks_fts, rowid, ...) VALUES('delete', ...)`/insert calls issued by the writer, or queries against it can return stale results or fail once the backing row is gone. Reconcile step 5 ("for every oid not in blobs, read the file and insert chunks, records and refs") and the GC clause ("delete blobs … not in any project's view whose last_seen is older than 14 days") never mention writing to `chunks_fts` at all — not on insert, not on the GC delete. Two implementers would genuinely diverge here: one adds triggers in the schema, one issues manual shadow-table writes in Go, a third forgets it and ships a search that returns nothing (or, after the first GC, a corrupt/erroring FTS index once the underlying `chunks` rows it points at are gone). This has to be decided in Interfaces § Storage/Reconcile before `internal/search`'s schema task (Phase 1, first checklist item) starts.
- **BLOCKING (inferable, but unstated) — the query-time join back to the current `files` view is never written down.** Storage's prose says "Everything derived from content is keyed by the git blob ID (`oid`)… `files` is the current view" — implying every branch/commit's content is chunked once into a project-wide, cross-branch cache, and that `search` must restrict hits to oids present in the *current* `files` row for the project. But the Interfaces § Commands / Query description of `search` never states this join explicitly, and it is exactly what the first E2E criterion tests ("after `git switch` to a branch that lacks a word, `search` returns no hit from the other branch"). An implementer who reads only the Commands section and queries `chunks_fts` (a global, non-per-project table with no `project` column) without joining through `chunks.oid = files.oid AND files.project = ?` will produce stale/cross-branch hits and fail that criterion silently in normal use (it only shows up when two branches diverge on a word). Worth one explicit sentence in Query, not left to be inferred from Storage's aside.
- **Non-blocking, but worth a one-line default before Phase 1**: the mapping from a docs-root directory name to a `--kind` value is stated for the one irregular case ("`architecture` → `adr`") but not for the regular ones (`bugs`→`bug`, `changes`→`change`, `handoffs`→`handoff`, `reviews`→`review`, `retros`→`retro`, `specs`→`spec`) — low ambiguity, a reasonable implementer singularizes correctly, but it is technically unstated.
- **Non-blocking**: when an oid (content) is current at more than one path in the same project (duplicate file content, or a rename mid-reconcile), `records`/`refs` are parsed once per oid but the parsing rules are path-dependent (a Bug's ID comes from "the file name's `BUG-NNN`"; a Plan's kind comes from living under `<docs>/plans/`). The plan does not say which path is used when an oid has more than one current path. A rare edge case; a "first path seen this reconcile" default is fine and an implementer would plausibly pick it without asking.

## (2) Scope

Everything under Scope § In traces to a D-row attributing it to the operator's own words or to the research/audit records the operator commissioned (D1–D13). I found no endpoint, output shape, or configuration surface added beyond what a D-row covers:
- The global index, the recall/grep/token baselines, and the vector gate are each tied to specific operator quotes or to the Phase 4 gate the operator explicitly conditioned ("if it gives the agent … more memory", "audit … which services benefit", the eight-row recommendation table).
- `--open` and the metadata graph are the concrete shape of Scope § In bullet 1 ("a metadata graph of the trail … status … cross-references"), not an unrequested addition.
- No config knob is added for a value that plausibly never varies (busy timeout, GC age, per-file cap, snippet length are all fixed constants in the design, not flags).
Scope is clean.

## (3) Readiness

Questions a first-pass implementer would have to ask, or guess at:

1. **BLOCKING** — How does `chunks_fts` stay in sync with `chunks` on insert and on GC delete (triggers vs. manual shadow-table writes)? (Gap 1 above.)
2. **BLOCKING** — Does every `search`/`--refs`/`--open` query join through `files` to restrict hits to the current project's oids, and is that also true for `--global`'s per-project sub-queries? (Gap 2 above.)
3. NON-BLOCKING — directory-name-to-`--kind` mapping for the regular cases (Gap 3): a reasonable implementer singularizes without asking.
4. NON-BLOCKING — which path is used to parse records/refs when an oid has more than one current path (Gap 4): rare, obvious default.
5. NON-BLOCKING — exact `modernc.org/sqlite` DSN syntax for the stated pragmas (the plan writes `journal_mode=WAL, busy_timeout=10000, foreign_keys=off` as intent, not as a literal connection string; the driver takes `_pragma=` query parameters). An implementer translates this without needing to ask.
6. NON-BLOCKING — the Phase 3 byte budget is razor-thin: `go run ./cmd/sh-agent render` on this checkout today measures 21,952 of the 22,000-byte budget, i.e. 48 bytes of headroom before any of the three named cuts are applied. The plan already names this exact contingency ("If these do not free enough, Phase 3 stops and puts further cuts to the operator; the budget is never raised."), so it is not an undecided gap — just flag it as the first thing to measure in Phase 3, since the margin is essentially zero.

I did not find a Verification/E2E criterion that bakes in an implementation choice beyond what Interfaces already commits to, and I did not find a verify-clause fixture that fails to produce the condition it claims (the two candidates I checked — `TestBusyRefusesNotStale`'s "second connection" and `TestConcurrentReconcile`'s "twice as subprocesses" — are both achievable as described: SQLite's file-level locking works across two connections in one process as well as across two processes, so "second connection" doesn't need cross-process to hold a real write lock).

## Criteria

This is a design-document cold-read, not a diff review; the six review-standards criteria apply loosely:
1. Dead code — n/a (no code yet).
2. Flow tracing — fail on the two BLOCKING gaps above (chunks→chunks_fts sync; query→files join); everything else traced (verb switch placement, `Collect`'s explicit allowlist needing a new skip entry, `reorder`/`takesValue`, size-test budget) matches the actual source.
3. Standards conformance — R4's proposed text (D1) correctly threads the "stdlib only" rule change through to the project `AGENTS.md` bullet (Phase 1 task); D9 justifies the one dependency and the no-vendor decision with a measured number, matching the security-exemption/stdlib-first bar. No unaddressed restatement found (checked `REQUIREMENTS.md` R4/R9, project `AGENTS.md`, `README.md` § How it works — R9's four procedures and the map are about delivery paths, not query verbs, so search/index correctly don't need a map change).
4. Directory layout — `internal/search` as a new package, kept out of `internal/harness` so the installer's package doesn't gain the driver, matches the project's Go layout rule and the decoupling concern the plan itself names.
5. Simplicity — pure-Go driver, no vendoring (measured 137 MB, D9), no service/watcher/hook (R3/R8), reconcile-at-query-time (D3) — all justified against the stdlib-first/DIY-except-security bar; no hand-rolled crypto or auth surface in this plan.
6. Decoupling — the oid-keyed, cross-branch content cache decoupled from the per-project `files` view is a reasonable seam, but it is exactly the seam the two BLOCKING gaps sit on: the plan draws the boundary but doesn't finish specifying the maintenance obligation (FTS sync) or the read obligation (current-view join) at that seam.

## Not examined

- The scratch research/audit briefs the plan cites (`plans/2026-09-29-recon-project-index.md`, `plans/2026-09-29-audit-search-consumers.md`, `plans/2026-09-29-research-sqlite-index.md`) were not opened; D-row attributions were taken at the plan's word.
- `agents/*.md` and `skills/*/SKILL.md` files D13 names for one-line additions were not individually opened to confirm the cited line numbers (`skills/review-standards/SKILL.md:48`, `agents/dreamer.md:36,29`, `skills/dream/SKILL.md:41,48`, `skills/bug/SKILL.md:33`, `skills/feature/SKILL.md:19`) — only `skills/handoff/SKILL.md:154` was checked and confirmed accurate.
- `modernc.org/sqlite`'s actual DSN pragma syntax and its FTS5 support surface were not built/run — Risks already flags this as checked only for linux/amd64 (run) and darwin/arm64 (build), with linux/arm64 deferred to the Phase 1 build task; taken as stated.
- Phase 4/5 sandboxed efficiency-run mechanics (cost, kit contents) were read but not evaluated against the actual `.agents/scratch/eff-kit/` (does not exist yet — nothing to check).
