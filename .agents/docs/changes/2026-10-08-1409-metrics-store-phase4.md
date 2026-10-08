---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 4: events and `sh-agent metrics show`

- `internal/trail` (new, driver-free): `Dates` and `StatusArgs`, moved from `internal/dream/collect.go` (D31). Dream calls them there, so `metrics` never imports `dream`.
- `internal/metrics/events.go`: `joinEvents` derives each main session's events. Each event's key is (session, kind, ref), so re-joining adds nothing, and a value changes in place when a later session makes a commit ambiguous. The events:
  - commits in the window `[start, end + 30 min]`, from git with no index lock and no fsmonitor; ambiguous when two windows overlap;
  - handoffs, by trail time;
  - verifier status (`## Status:`);
  - review findings by severity;
  - bugs, linked by a 7+ hex prefix of one of the session's commits;
  - corrections from `feedback.jsonl`.

  `Ingest` runs the join and prints a note for a project root that is gone or has no git history.
- `internal/metrics/show.go`: `Show` builds the A–D summary and the `daily` series per CLI, with label, project and window filters. `WriteJSON` and `WriteText` print it. A metric a CLI cannot produce is null / n/a with a note.
- `cmd/sh-agent`: the `metrics show` verb, with `--since`, `--project`, `--label`, `--all-labels` and `--json`, plus usage and doc text and cmd tests.
- Docs: a README row; the spec's `## Output` rewritten to the real shape; B2 (main sessions only); the missing-root failure row.
- **Found by the E2E and fixed under a test.** Antigravity active time exceeded wall time, because subagent gaps were counted. The new `TestShowActiveTimeMainOnly` was red, then green.
- Verification:
  - Copies of both CLIs' real data: 100 sessions. `show --json` `B.tool_calls_by_tool` equals the independent counts (Claude 1899, Antigravity 190), a second ingest added 0 rows, and 91 events were joined. The copies were deleted.
  - `go test -race -timeout 30m ./...` exit 0.
- Agents dispatched: none (review and security at the Definition of Done).
