---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 7: feedback into the learning loop; sandbox runs on both CLIs

- `internal/metrics/numbers.go`:
  - `Snapshot` opens a read-only copy of the store (no `-wal`/`-shm` beside it, no migration) and refuses another schema.
  - `Summary.Numbers` reduces a summary to the package's own metric keys and numbers; a metric a CLI cannot produce is absent.
  - Tests: `TestNumbersOnly` (its must-fail control on a copy, leaking tool names into keys, went red) and `TestSnapshotReadOnly`.
- `internal/dream`: each project's prompt gains a `metrics` input, `Numbers()` for the project's window and the window of the same length before it, read from a snapshot (R3 addendum, ADR-055a).
  - It never marks a project changed and is dropped first over the cap.
  - The instructions treat a worsened metric as evidence only when other inputs show its cause.
  - Tests: `TestPromptCarriesMetricsNumbersOnly` (no tool, command, path or session id in the block; nothing left beside the store) and `TestMetricsDroppedFirstAndNeverChanges`, both red before the code.
- `skills/retro/SKILL.md` Step 1 source 11, **metrics evidence**: `sh-agent metrics show --json --project` for the retro window and the one before. A worsened metric is a finding only with a cause in sources 1–9, n/a is never a regression, and a missing store gets one line.
- `cmd/sh-agent`: `show`, `report` and `export` refuse a home with no store, with exit 1 and a pointer to `ingest` (D38, `TestMetricsWithoutStore`). The usage and doc text for dream now describe the metrics.
- `internal/metrics/verb.go`: `sh-agent`'s second-level verbs (`ingest show report export add`) are kept (D38).
- Docs: `specs/sh-agent/dream.md` (the metrics block), the metrics spec (a no-store row, the normaliser rule), README (dream bullet and row), and `testing/sh-agent-matrix.md` (four `metrics` rows).
- **Sandbox runs** (`testing/2026-10-08-verify-metrics-sandbox.md`): a scratch kit per CLI, kit-preamble first, linked credentials, clean exit checks.
  - Every metrics check passes: 2 of 2 on Claude Code, 1 of 1 on Antigravity.
  - Both CLIs' `/retro` ran source 11 and cited the metrics.
  - The skill-matrix `retro` row check is **0 of 2 on Claude Code**: both runs scoped the `-race` rule to the project under Step 2, so no harness file was queued. It is 1 of 1 on Antigravity. A control at the pre-feature commit is Open.
- Verification: `go test -race -timeout 30m ./...` exit 0, twice. The last `verb.go` change passed `go test -race ./internal/metrics`; the verifier's run at the Definition of Done covers the whole tree.
- Agents dispatched: none yet (review, security and verifier follow at the Definition of Done).
