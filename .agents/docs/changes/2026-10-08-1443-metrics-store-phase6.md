---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 6: `sh-agent metrics report`, the visual report

- `internal/metrics/report.go` renders `show`'s `Summary` as one self-contained HTML file, built with `html/template` and inline SVG drawn in Go. It has no script, no host, no `@import` or `url(`, and covers light and dark. The charts:
  - per-day lines (tool calls, tokens, cost, search:grep, errors) with one line per CLI, end labels nudged apart, and `<title>` hover values;
  - horizontal bars for tool calls by tool and for subagent tool calls by type;
  - a passed-of-k chart per sandbox label;
  - the summary table, and a data table under each chart.

  Colours follow the CLI (dataviz reference palette slots 1 and 2, validated pair) through CSS classes. `WriteReport` writes 0600 and refuses an existing file.
- `cmd/sh-agent`: the `metrics report` verb (`--since`, `--project`, `--label`, `--all-labels`, `--out`), plus usage and doc text and cmd tests.
- Tests: `TestReportSelfContainedAndConsistent` (no script or host; the summary row and the daily sums equal `show`), `TestReportEscapesNames` (its must-fail control on a copy, with escaping off, went red), `TestWriteReportRefusesOverwrite`, and `TestReportDrawsNoMissingValues`.
- **Browser check, from copies of both CLIs' real data.** The page was viewed in styrman-browser at 1100 px light and dark and at 400 px, with no console messages; screenshots are in `.agents/scratch/report-{light,dark,400}.png`. The first render showed three defects, all fixed under `TestReportDrawsNoMissingValues` (red, then green):
  - the legend swatches were sanitised to `ZgotmplZ` (an inline CSS variable);
  - Antigravity was drawn as a 0 cost line and as a 0 ratio on days with no grep call;
  - the end labels collided.

  Known limit: tick labels are about 6 px at 400 px; the tables carry the numbers at full size.
- Tests: the report fixture now uses a renamed sandbox copy (`sandboxCopy`), because the same transcript ingested unlabelled first is not relabelled by a later `--source`.
- Docs: spec § Report rewritten; README row.
- Verification: `go test -race -timeout 30m ./...` exit 0.
- Agents dispatched: none (review and security at the Definition of Done).
