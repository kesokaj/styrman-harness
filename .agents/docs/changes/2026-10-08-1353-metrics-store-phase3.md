---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 3: the Antigravity ingester

- `internal/metrics/agy.go` reads every Antigravity conversation that changed since its last ingest, from copies (`search.OpenCopy`); the live databases are never opened. What it reads:
  - user, reply and tool steps;
  - a tool's name at `metadata` 4→2, mapped to the common vocabulary (`agyTools`; an unknown name is kept as `agy:<name>`);
  - from the arguments JSON at 4→3, only `CommandLine` (through `Verb` and the salted hash) and `AbsolutePath` (trail kind and skill name), both in memory;
  - errors from status 7;
  - subagents from `conversation_summaries` (`agent_name`, `parent_conversation_id`);
  - per-generation tokens, the model and the context from `gen_metadata` (D3a).

  `step_payload`, the tool output, is never read (D28).
- `internal/metrics/store.go`: migration 2 adds `steps.context`, and `TestMigrationOneToCurrent` checks a schema-1 store upgrading with its rows. `internal/search/export.go`: `AgyTime`.
- Docs:
  - README's `metrics ingest` row now covers Antigravity.
  - The spec gains the full tool-vocabulary table, where each CLI's shape is read, and the `context` column.
  - The plan gains D35 (the spike) and D3a (exact usage paths), and the Phase 3 boxes.
- Verification:
  - `TestAgyIngestShape` runs with the source directories at 0o500 and the fixture in WAL mode. A mutation that opens the live file in place fails both conversations, so the test proves the copy path.
  - The canary is absent from the store, its WAL and the export.
  - A scratch copy of the real Antigravity data matched a direct count of the databases: 190 tool, 11 user, 194 reply and 196 generation steps; 15 sessions, 4 of them subagents. A second ingest added 0. The copy was deleted.
  - `go test -race -timeout 30m ./...` exit 0.
- Agents dispatched: none (review and security at the Definition of Done).
