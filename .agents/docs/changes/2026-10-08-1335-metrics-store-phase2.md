---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 2: the Claude Code ingester and `sh-agent metrics ingest|export`

- `internal/metrics/verb.go` — `Verb`, the command-shape normaliser. It strips here-documents and comments, and keeps quotes, `${…}` and `$((…))` unsplit. It handles case patterns, subshells, redirections and shell keywords. A program must be name-shaped, and a glued short-flag value keeps only its flag letter (D33).
- `internal/metrics/claude.go` — the all-projects walker, including `subagents/` and their `.meta.json` agent type. It writes one row per message (usage is the max over streamed lines) and one per tool call, folding each tool result's error and denial flags into its tool row. It also reads interrupts, compaction boundaries, the last `cost-state`, a Skill's name, and a Read's trail kind (the path is not stored). Ingest is incremental by byte offset, waits on a partial last line, and re-reads a replaced file from 0 without duplicating rows.
- `internal/metrics/ingest.go` — `Options`, with the transcript-root rule `--source` > `--home` > `CLAUDE_CONFIG_DIR` > home (D9); `Validate`; `Ingest`; `Result.Line`. `agy.go` is a placeholder until Phase 3 and says so.
- `cmd/sh-agent` — the `metrics ingest` and `metrics export` verbs, their usage text and package doc, and `metrics_test.go`. The tests cover the decoy `CLAUDE_CONFIG_DIR` being ignored under `--home`, export and its refusal, and exit 2 for every usage error with no store created.
- `internal/harness/report.go` — names `metrics.db*` as skipped ("the metrics store; never sent"). `TestCollectNeverSendsMetricsStore` and `TestUninstallKeepsMetricsStore` were added (D17).
- Docs:
  - README: verb rows and a state-table row.
  - The spec's § Normalised verb is rewritten.
  - `memory-loop.md` § Known differences gains three metrics rows.
  - Plan rows D3a, D33, D34 and D35, the Phase 3 Antigravity spike.
- **A defect caught before commit (not a bug record; it never reached a commit).** The first normaliser, run over a scratch copy of the real transcripts, stored content as "programs" and flags: file names, `${GH_TOKEN:-}`, case patterns, `cut -c1-200`. That broke the shape-only rule. The rewrite (D33) was audited again over the same copy: 117 programs and 190 flags, all name-shaped. The copy was deleted.
- Verification:
  - A scratch copy of the real transcripts ingested 85 sessions and 3667 steps. Its 1837 tool calls equal an independent Python count, with 0 differences, and a second ingest added 0 rows.
  - `go test -race -timeout 30m ./...` exit 0.
  - Mutation controls on copies went red as expected: keeping positionals, removing heredoc stripping, deleting the store in uninstall.
- Agents dispatched: none. The feature's one review and security pass run at the Definition of Done.
