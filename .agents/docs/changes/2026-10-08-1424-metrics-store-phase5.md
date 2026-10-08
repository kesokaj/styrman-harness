---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 5: dream ingests first; the map; the A/B kit

- `internal/dream/metrics.go` and `run.go`: a real `dream` run calls `sh-agent metrics ingest` store-wide inside `dream.lock`, before choosing projects, and logs the summary line. A failure is logged as `metrics ingest failed: …` and changes neither the exit status nor `last_run.error` (D16). A dry run prints `[dry-run] would ingest metrics into <path>` and writes no store (`dream.go`). New tests: `TestRunIngestsMetricsFirst`, `TestDryRunDoesNotIngest`, `TestIngestFailureDoesNotFailDream`, each red before the code.
- **README map (AGENTS.md rule 3, in this commit with the wiring).** The configDir node names `metrics.db (session shape, never sent)`, and the dream→configDir edge is labelled `proposals/dream/, metrics.db`. No line was added or rerouted. The Dream bullet and the `dream` verb row describe the ingest. `TestHarnessMap` passes.
- `specs/sh-agent/dream.md` § What it does: the ingest step.
- `internal/metrics/ingest.go`: an `--outcome` ingest stamps the label's sessions that have no outcome yet, even when no transcript changed (D36; `TestOutcomeStampedAfterTheRun` was red, then green). Spec flag row updated.
- `testing/ab-kit/run.sh` ingests the arm into the kit's own store (`--home "$AB_OUT"`, never the real home) and prints the `--outcome` command to run after `accept.py`. `metrics.py` is removed (D24); `metrics show --home "$AB_OUT" --all-labels` replaces it. Checked with `bash -n` only: a real kit run costs tokens and was not run.
- `testing/guide.md` § 3g (new): metrics tests, the must-fail canary control, the real-data check from copies only, and the report's browser check. The project AGENTS.md points at § 3g.
- Verification: `go test -race -timeout 30m ./...` exit 0.
- Not done here: the interactive claude.ai page of the map (AGENTS.md rule 5) follows from the map at the end of the feature.
- Agents dispatched: none (review and security at the Definition of Done).
