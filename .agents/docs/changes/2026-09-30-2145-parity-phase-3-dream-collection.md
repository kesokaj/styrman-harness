---
owner: main
date: 2026-09-30
---
# Phase 3 of the parity plan: `sh-agent dream` collects, redacts and dry-runs     (2026-09-30)

## Summary
`sh-agent dream --all | --project DIR --dry-run` now shows, per project, what the background dreamer would send and which tool-free CLI it would run. The run itself is Phase 4.
- **Projects**: `search.ListProjects` reads the global index's project list from a copy. A scratch dry run found that a read-only open of a WAL database leaves `-wal`/`-shm` behind.
- **Antigravity transcripts**: `search.AgyConversationText` copies `conversation_summaries.db` and the matching conversation databases, maps them by `workspace_uris`, cuts by `last_modified_time` (D24), and keeps the prose runs of the step blobs (D8a probe). A smoke run on copies of real conversations returned the known prompts verbatim.
- **Collection**: `internal/dream` collects every memory file, the feedback lines, trail records and Claude Code transcript lines after the watermark, and the installed skills (global, so never "a change" alone). It redacts each through the new `harness.Redact` (report's rules, held to the same `shapeSecrets` list) and writes the home path as `~`.
- **Cap**: 256 KB. Transcripts are dropped first, then trail, feedback and skills; memory over the cap fails the project (D21).

Found on the way, each a test red before the fix:
- the ListProjects side files;
- a reversed cleanup order that left an empty `~/.config/sh-agent`;
- the driver handing a `datetime` column back as `time.Time`.

## Type
feature (phase 3 of 6)

## Files changed
- `internal/dream/{dream,collect,prompt,state}.go` and `collect_test.go` (new)
- `internal/search/dream.go` and `dream_test.go` (new)
- `internal/harness/report.go` (`Redact`) and `report_test.go` (`shapeSecrets`)
- `internal/memory/memory.go` (`Store`)
- `cmd/sh-agent/main.go` (the `dream` verb, `--all`, `--project` in `takesValue`, usage, doc comment) and `memory_test.go` (`TestDreamVerb`)
- the plan (D8a, ticks)

## Agents dispatched
verifier (`testing/2026-09-30-verify-parity-memory-dream-phase-3.md`: FAIL on branch coverage, PASS after the added tests).

## Verification
`go test -race -count=1 ./...` 426 pass. § 3d: rebuild 1.94 s, warm search 42 ms. A mutation of the skills rule failed its two tests. A dry run over two indexed fixture projects left the home unchanged. Not exercised: the CLI call (Phase 4).

## Requirements
R3's exception and R4's dream clause (ADR-047a).
