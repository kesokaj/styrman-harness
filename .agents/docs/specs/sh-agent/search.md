---
owner: main
date: 2026-09-29
---
# `sh-agent search` and `sh-agent index`

Status: Phase 1 (project index) built 2026-09-29; Phase 2 (global search index) specified below and built when plan `plans/2026-09-29-search-index.md` ticks it. Requirement: R4 (one dependency, `modernc.org/sqlite`), ADR-045.

## What it is

A derived SQLite index the agent queries instead of reading whole files. Git and the files in `~/.config/sh-agent/` stay the source of truth: the index is a cache, never committed, never sent by `report`, and deleting it loses nothing but the next query's rebuild time (about two seconds for this repository's 1,050 files).

## Commands

```
sh-agent search [--limit N] [--kind K] [--global] <words...>
sh-agent search --refs <BUG-NNN|ADR-NNN|path> [--kind K] [--global]
sh-agent search --open [--kind K] [--global]
sh-agent index [--global] [--rebuild]
```

- **Words** are each quoted for FTS5, so `BUG-049`, `a.go` or `NOT` are searched as text, never as query syntax, and ANDed; with no hit the query is retried once as OR and the header says so. Porter stemming: `verify` finds `verifies`. Hits are ranked by BM25 with a heading weighted twice the body, at most 2 per file, `--limit` (1–50, default 10) in all.
- **Each hit** prints `path:line  kind[ ID status][  date]` and a one-line snippet (at most 200 characters, matched terms in `»«`). The line is the first line of the chunk holding a matched term.
- **`--refs X`** prints the ID's own record (`ADR-040  <path>:<line>  Superseded by ADR-041  <date>  <title>`, or `no record`), or a path's state (`tracked`, `modified`, `untracked`, `absent`, or `directory` for a trailing `/`), then `referenced by N:` and one entry per file and record that names X, newest record date first, with `×N` when it names X more than once. A path matches as written and relative to the docs root (the trail writes `plans/x.md` for `.agents/docs/plans/x.md`); no basename matching.
- **`--open`** lists bugs with `status: open` and plans with an unticked `- [ ]` outside fenced blocks, newest first.
- **`--kind`** is one of `bug adr plan change handoff review testing retro spec deploy trail code doc feedback proposal` and filters every mode by the kind of the printed path. `trail` is everything under the docs root; `feedback` and `proposal` exist only in the global search index.
- **`index`** prints `files`, `chunks`, `records`, `refs`, `changed` (by this reconcile), `bytes` (database and WAL) and `ms`. `--rebuild` deletes the database and builds it again.

Exit status: 0 with hits or without; 1 on a refusal or an error.

## Where the indexes live

| Index | Path | Covers |
|---|---|---|
| Project | `<checkout>/.agents/scratch/search.db` (+ `-wal`, `-shm`) | every tracked, modified and untracked-not-ignored file under 1 MiB with no NUL byte in its first 8 KiB; symlinks and submodules skipped |
| Global search index | `~/.config/sh-agent/search.db` (+ `-wal`, `-shm`), under `--home` when given | `feedback.jsonl`, `retros.jsonl`, `proposals/*.md`, and the docs root of every checkout that has run `search` or `index` |

The global file is named "the global search index", never "cross-project index": R7 already gives that name to `retros.jsonl`.

## Freshness: reconcile before every answer

1. `git rev-parse --show-toplevel` from the working directory.
2. `git ls-files -s -z`: path and blob ID of every index entry.
3. `git status --porcelain=v1 -z --untracked-files=all`: modified and untracked paths, their IDs from one `git hash-object --stdin-paths`; deleted paths leave.
4. Everything derived from content (chunks, records, references) is stored once per blob ID, so a branch switch re-reads only content the index has never seen. Every query joins through the current path-to-blob view, which is what keeps another branch's content out of an answer.
5. When nothing changed, no write lock is taken; the reconcile costs about 15 ms here. A blob no view references is deleted after 14 days unseen.

No watcher, no service, no hook (R3, R8).

## Trail records

The docs root is `.agents/docs` when it exists, else `docs`.

| Record | Read from | ID, status |
|---|---|---|
| Bug | `bugs/BUG-NNN*.md` | frontmatter `id:` (else the file name), `status:`, `date:`, first `# ` heading |
| ADR | `architecture/DECISIONS.md` headings `## ADR-NNN: <title>     (<date>) — <status>` | one record per heading; status is everything after the dash |
| Plan | `plans/*.md` | `open: N unticked` or `closed`; line of the first unticked box |
| Other trail file | any `.md` under the docs root | its docs-relative path; `date:` |

References are record IDs (`BUG-NNN`, `ADR-NNN`) and path-shaped tokens on every line of every indexed file, with a trailing `:NN`, trailing punctuation and a leading `./`, `a/` or `b/` removed.

## Failure modes

| Case | Output | Exit |
|---|---|---|
| Not inside a git checkout, or no `git` | `search: not inside a git checkout; use grep` | 1 |
| `.agents/scratch/search.db` not gitignored | `search: .agents/scratch/search.db is not gitignored; add ".agents/scratch/*" to .gitignore (/init-project does)`; nothing created | 1 |
| `.agents` or `.agents/scratch` is a symlink | `search: .agents/scratch is a symlink; refusing to write through it` | 1 |
| Another process holds the write lock past 10 s | `search: index busy (another sh-agent is indexing); retry, or use grep`; no answer printed | 1 |
| Unreadable database, or a schema other than 1 | `note: index at <path> was unreadable (or an older schema); rebuilt`, then the answer | 0 |
| Argument errors (`--limit` out of range, unknown `--kind`, two modes at once, no mode) | `search: <what is wrong>; see sh-agent help` | 1 |
| Global: a registered checkout is gone | `note: dropped <root> (no longer a git checkout)` | 0 |
| Global rebuild | `note: N project registrations dropped; each returns on its next search` | 0 |
| Network filesystem | not supported: SQLite's WAL needs one host's shared memory; the SQLite error is printed | 1 |

`sh-agent report` never collects `search.db`; its dry run lists it as skipped. `sh-agent uninstall` leaves `~/.config/sh-agent/search.db`, which is not in the install manifest; delete it by hand if wanted.

## Measured (this repository, 2026-09-29, Linux amd64)

| Measurement | Result |
|---|---|
| Full build, 1,050 files | 2.2 s |
| Warm `search`, nothing changed, end to end | 23–67 ms |
| Warm `--refs` / `--open`, end to end | 47–60 ms / 20 ms |
| Database size | 23 MB (10,340 chunks, 913 records, 39,560 references) |
| Binary size, linux/amd64 | 13.7 MB |
