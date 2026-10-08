---
owner: main
date: 2026-10-08
status: draft
---
# Spec — sh-agent: the metrics store   (2026-10-08)

## Purpose

`sh-agent metrics` keeps the **shape** of every Claude Code and Antigravity session in `~/.config/sh-agent/metrics.db`. Shape means:
- the tools each session called;
- the normalised command verbs;
- tokens, timings and errors;
- events joined from git and the trail.

The operator queries the store for agent evaluations (`show`, `report`). The nightly dreamer and `/retro` read it as evidence for proposals.

It is not:
- a transcript archive: it holds no arguments and no prose;
- a cache: it is the only copy once a transcript expires (R4);
- the "metrics record" that `report` sends (R8), which is a different thing.

Decisions: ADR-055. Plan: `plans/2026-10-08-metrics-store.md`.

## Interfaces

### Verbs

The verbs reuse the global `--home` and `--vendor` flags.

```
sh-agent metrics ingest [--home H] [--vendor claude|agy|all] [--source DIR --label KEY=VALUE [--outcome pass|fail]]
sh-agent metrics show   [--home H] [--since YYYY-MM-DD] [--project DIR] [--label KEY=VALUE | --all-labels] [--json]
sh-agent metrics report [--home H] [--since YYYY-MM-DD] [--label KEY=VALUE | --all-labels] [--out FILE]
sh-agent metrics export [--home H] FILE
```

| Flag | Type | Default | Meaning |
|---|---|---|---|
| `--vendor` | `claude`, `agy`, `all` | `all` (not the manifest's) | Which CLI's sessions `ingest` reads. |
| `--source` | directory | none | A sandbox home: `DIR/.claude/projects` and `DIR/.gemini/antigravity-cli`. Requires `--label`. |
| `--label` | `KEY=VALUE` | none | On `ingest`: stamped on every session added from `--source`. On `show`/`report`: shows only that label. |
| `--outcome` | `pass`, `fail` | none | The run's acceptance result, supplied by the kit. Used for pass^k. |
| `--all-labels` | bool | false | Include labelled (sandbox) sessions. Without it, and without `--label`, they are excluded. |
| `--since` | date | 30 days ago | Start of the window. |
| `--project` | directory | all | Only sessions whose `project_root` is this directory. |
| `--json` | bool | false | Print the `## Output` object instead of tables. |
| `--out` | file | `~/.config/sh-agent/reports/<YYYY-MM-DD-HHMM>-metrics.html` | Where the report is written. An existing file is refused. |

### Where transcripts are read from

The first rule that applies wins:
1. `--source DIR`
2. `--home H` → `H/.claude/projects`, `H/.gemini/antigravity-cli`
3. `CLAUDE_CONFIG_DIR`, only when neither flag is given
4. the real home

### `ingest` output and exit status

On success it prints this line and exits 0:

```
metrics: N sessions, M steps, K events added; U lines skipped (unknown shape); F sources failed
```

- Exit 1 when F > 0. Each failed source is named on its own line, and every other source is still ingested.
- Exit 2 on a usage error.

### Normalised verb

`Verb(cmd string) []string` reduces a command line to its shape:
- It splits on `;`, `&&`, `||`, `|` and newlines. Quoted text is not split.
- In each segment it drops:
  - leading `VAR=value` words;
  - `sudo`, `env`, `command` and `time`;
  - `go run <pkg>`, keeping the package's last element as the program.
- It keeps, in order:
  - the program's basename;
  - the subcommand, when the program is `git go gh docker gcloud kubectl npm npx sh-agent terraform` and the next word is not a flag;
  - each flag NAME (`--x=v` becomes `--x`).
- Positional words and values are never kept. The one exception is a bare `.` after `git add`.

### Shared tool vocabulary

Vendor tool names map to one common name. Unknown names are kept as `<vendor>:<name>`.

| Common | Claude Code | Antigravity |
|---|---|---|
| Bash | `Bash` | `run_command` |
| Read | `Read` | `view_file` |
| Grep | `Grep` | `grep_search` |

The Phase 3 spike completes this table from real conversation databases.

A **skill load** is:
- on Claude Code, a `Skill` call (its `skill` value is stored);
- on Antigravity, a read of `skills/<name>/SKILL.md` (only `<name>` is stored).

### Output (`show --json`)

Keys are snake_case, like `report.json`. Every metric is split by vendor. A metric one CLI cannot produce is `null` and is named in `notes`, never shown as 0.

```json
{
  "window": {"since": "2026-10-01", "until": "2026-10-08"},
  "sessions": {"claude": 12, "agy": 3},
  "subagents": {"claude": {"reviewer": 4, "verifier": 2}, "agy": {}},
  "A": {"dod_closed": {"claude": "3 (ambiguous 1)"}, "verifier": {"PASS": 2, "FAIL": 0},
        "review_findings": {"high": 1, "medium": 3, "low": 2}, "escaped_bugs": 0,
        "corrections": 2, "pass_k": {"arm=lean": {"k": 3, "pass_k": 1}}},
  "B": {"tool_calls": {"claude": 900, "agy": 40},
        "tool_calls_by_tool": {"claude": {"Bash": 600}, "agy": {"Read": 20}},
        "tokens": {"claude": {"in": 0, "out": 0, "cache_read": 0, "cache_write": 0}, "agy": null},
        "cost_usd": {"claude": 12.3, "agy": null}, "wall_min": {}, "active_min": {},
        "peak_context": {}, "compactions": {}},
  "C": {"sqlite_verbs": {"claude": 38}, "search_grep_ratio": {"claude": 0.045},
        "forbidden": {"git add -A": 0, "git add .": 0, "git push --force": 0, "--no-verify": 0, "commit chained to test": 0},
        "first_reads": {}, "skills": {"claude": {"feature": 2}}},
  "D": {"tool_error_rate": {}, "denials": {}, "interrupts": {}, "retry_loops": {}, "subagent_tool_calls": {}},
  "notes": ["agy cost unavailable (CLI exposes none)"]
}
```

The values above illustrate the shape. They are not fixtures.

### Report (`metrics report`)

One HTML file, mode 0600:
- built with `html/template` and SVG drawn in Go;
- no `<script>` and no URL to any host;
- light and dark themes through `prefers-color-scheme`;
- readable at a 400 px width.

Its charts:
- sessions, cost and tokens per day (one line per vendor);
- tool calls by tool;
- the search:grep ratio over time;
- the error and denial rate per week;
- subagents by type;
- pass^k by label;
- the A-group table.

Every number in the report equals the same key in `show --json` over the same window.

### Metric definitions

| Id | Metric | Definition |
|---|---|---|
| A1 | DoD closure | A main session's window is `[started_at, ended_at + 30 min]`. The session is closed when the window holds at least one commit in its project AND a handoff whose `TrailDates` time falls in it. |
| A2 | Verifier status | The `## Status:` heading of `testing/*-verify-*.md` files in the window. Main sessions only. |
| A3 | Review findings | `**[high\|medium\|low / confidence: …]**` lines in `reviews/*.md` in the window. Main sessions only. |
| A4 | Escaped bugs | A `bugs/BUG-*.md` whose text contains a hex prefix of 7 or more characters of one of the session's commits. |
| A5 | Corrections | `feedback.jsonl` entries in the window whose project is the session's root. |
| A6 | pass^k | Per label, 1 when all k `--outcome` runs passed. k and the pass count are shown. |
| B1 | Tokens and cost | Claude: the maximum of each usage field per message id, summed. Cost and duration come from the session's last `cost-state` line. Antigravity: `gen_metadata` usage. |
| B2 | Time and volume | Wall time = `ended_at - started_at`. Active time is the same, less gaps of 5 min or more between steps. Also turns and tool calls. |
| B3 | Context | Peak context = max(input + cache read + cache write) of one message. Also the number of compactions. |
| B4 | Subagents | Count by type, the tool calls of each, and SendMessage resumes. |
| C1 | Search vs grep | `sh-agent search\|index` verb segments divided by `grep`/`rg` segments plus Grep tool calls. |
| C2 | Forbidden shapes | `git add -A`, `git add .`, `git push --force\|-f`, `--no-verify`, and `git commit` in the same Bash call as `go test`. |
| C3 | First reads | The trail kind of a session's first three Read targets. |
| C4 | Skills | Skill loads per session. |
| D1 | Errors | Tool error rate, denials and interrupts. |
| D2 | Retry loops | The same `cmd_hash` in two consecutive tool steps. |
| D3 | Subagent volume | Tool calls per subagent session. |

## Data

`~/.config/sh-agent/metrics.db` has mode 0600 in a 0700 directory. A symlinked file or directory is refused. It uses WAL with a 10 s busy timeout. Every time is stored as unix milliseconds, UTC. Writer: `internal/metrics` only.

```sql
meta(key TEXT PRIMARY KEY, value TEXT)            -- 'schema', 'salt' (32 random hex)
sessions(id TEXT PRIMARY KEY, vendor TEXT, project_root TEXT, parent_id TEXT, agent_type TEXT,
         model TEXT, label TEXT, outcome TEXT, started_at INTEGER, ended_at INTEGER,
         cost_usd REAL, duration_ms INTEGER)
steps(session_id TEXT, ext_id TEXT, ts INTEGER, kind TEXT, tool TEXT, skill TEXT, target_kind TEXT,
      is_error INTEGER, denied INTEGER, tok_in INTEGER, tok_out INTEGER, cache_read INTEGER,
      cache_write INTEGER, cmd_hash TEXT, PRIMARY KEY(session_id, ext_id))
step_verbs(session_id TEXT, ext_id TEXT, n INTEGER, verb TEXT, PRIMARY KEY(session_id, ext_id, n))
events(session_id TEXT, ts INTEGER, kind TEXT, ref TEXT, value TEXT, PRIMARY KEY(session_id, kind, ref))
ingest_state(source TEXT PRIMARY KEY, size INTEGER, mtime INTEGER, offset INTEGER)
```

### Invariants

- **Identity.**
  - Claude: one `steps` row per assistant message (`ext_id` = message id) and one per tool call (`ext_id` = `tool_use.id`).
  - A tool result sets `is_error` and `denied` on its tool row. It never gets a row of its own.
  - A user turn's `ext_id` is the line's `uuid`.
  - Antigravity: `ext_id` = `<conversation id>:<step idx>`.
  - Re-ingesting anything adds no rows.
- **Path columns.** Only `sessions.project_root` and `ingest_state.source` may hold a path. No column holds an argument or prose.
- **Identifiers.** `steps.skill`, `sessions.agent_type` and the `git add .` verb are the only identifiers taken from a tool's input.
- **`events` payloads.**

  | kind | ref | value |
  |---|---|---|
  | `commit` | full SHA | `''` or `ambiguous` |
  | `handoff` | file basename | `''` |
  | `verifier` | file basename | the status word |
  | `review` | file basename | `high=N,medium=N,low=N` |
  | `bug` | `BUG-NNN` | the linked SHA |
  | `correction` | the entry's `ts` | `''` |

- **`cmd_hash`** = `hex(sha256(salt || command))[:16]`. It supports equality checks only.
- **Fixed markers.**
  - `denied`: a tool result starting with one of Claude Code's permission-refusal prefixes, or an Antigravity step status the spike names.
  - `interrupt`: a user line starting `[Request interrupted by user`.
  - `compaction`: Claude's compact-boundary line.
- **Antigravity sources.** Antigravity databases are read only from copies. `step_payload` is never read; tool names and command lines come from `steps.metadata`.
- **Migrations.** They run forward only, in one transaction. A store with a newer schema is refused untouched. The binary never deletes the store, and `uninstall` keeps it.

## Failure modes

| Trigger | Observable |
|---|---|
| `metrics.db` or its directory is a symlink | `metrics: refusing symlinked <path>`, exit 1 |
| The store's schema is newer than the binary's | `metrics.db schema N is newer than this binary (M); update sh-agent`, exit 1, store unchanged (`ingest`, `show`, `report`, `export`) |
| A transcript line is not JSON or has an unknown shape | Skipped and counted in U, exit 0 |
| A file shrank or was replaced | Re-read from offset 0, with no duplicate rows |
| An Antigravity copy or read fails | That conversation is named and counted in F; the others are ingested; exit 1 |
| An Antigravity token field is absent | Tokens are NULL; `notes` says `agy tokens unavailable for N sessions` |
| A project root is gone or not git | Session rows are still written; no commit events; a note |
| Two main sessions overlap in one root | Commits in the overlap are `ambiguous`; A1/A4 are shown as `n (ambiguous m)` and excluded from rates |
| Still locked after the 10 s busy timeout | The source counts as failed (`database is locked`), exit 1; the next run picks it up |
| `export FILE` or `report --out FILE` exists | `refusing to overwrite FILE`, exit 1 |
| `--source` without `--label` | Usage error, exit 2 |
| Ingest fails inside `dream` | A line in `dream.log`; the dream's exit status and `last_run.error` are unchanged |

## Non-goals

- Storing or sending arguments, prose, file contents or tool output.
- Sending metrics through `report`. Feedback leaves only as dream or retro proposals (ADR-055).
- A server, a live dashboard or a JavaScript chart library.
- Rebuilding the store from transcripts that no longer exist.

## Open questions

- The Antigravity field numbers for tool name, command line, subagent link, compaction and denial. Default: unknown items are R10 known differences, listed in `specs/sh-agent/memory-loop.md` § Known differences (Phase 3 spike).

## Verification

| What | Test |
|---|---|
| Store, migrations, refusals | `internal/metrics/store_test.go` (`-run 'Store\|Migrat'`) |
| `Verb` | `verb_test.go` (`-run Verb`) |
| Claude ingest, idempotence, transcript root, canary with its must-fail control | `claude_test.go` over `testdata/claude/` |
| Antigravity ingest from copies (read-only source directory) | `agy_test.go` over `testdata/agy/` |
| Shared vocabulary on both vendors | `tools_test.go` |
| Events and A-group definitions | `events_test.go` over a git fixture in `t.TempDir()` |
| `show --json` keys and values | `show_test.go` golden |
| Report: no script, no host, numbers equal to `show` | `report_test.go` golden, plus a browser check (plan Phase 6) |
| `report` never sends the store; `uninstall` keeps it | `internal/harness` `-run 'Report\|Uninstall'` |
| Dream ingests inside the lock, never in a dry run | `internal/dream` `-run Ingest` |
| E2E on both CLIs | testing guide § 4 runs, plan § Sandbox tests |
