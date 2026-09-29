---
owner: main
date: 2026-09-29
---
# Plan — `sh-agent search`: a project and a global SQLite index for the agent   (2026-09-29)

## Why

The operator, 2026-09-29, after the research record `plans/2026-09-29-research-sqlite-index.md`: "I have no problems to add this if it gives the agent a better understanding and more memory what it is doing. That is the whole point. Drift, searching through multiple files." Then, on the eight-row recommendation table (full versions in the conversation's scratch file, restated as D-rows below): "Go with your recommendations, start the phases", "make sure to add to the context how to use the database. Because having a database feels like a powerhouse in larger projects", "When everything is done I want you to run sandboxed projects that test the efficiency when building sample apps", and "audit the full project which services that will benefit of using the sqlite" (brief: `plans/2026-09-29-audit-search-consumers.md`).

Measured motivation (this project's transcripts, 2026-09-29): subagents spent 59% of their tool-output volume on file reads, 41% of that on trail files, and 583 of 827 reads loaded a whole file. Grep cannot rank, cannot find a claim restated in other words (BUG-027, BUG-097), and cannot answer "which records name this file" or "is this ADR superseded".

## Scope

In:
- A per-project index `.agents/scratch/search.db` over every tracked, modified and untracked-not-ignored text file of the git checkout, with a metadata graph of the trail (record IDs, status, date, cross-references, paths named).
- A global index `~/.config/sh-agent/search.db` over `feedback.jsonl`, `retros.jsonl`, `proposals/*.md` and the trail (docs root only) of every project that has run `sh-agent search` or `sh-agent index`.
- Two verbs: `sh-agent search` and `sh-agent index`. Both reconcile against git before answering. No service, no watcher, no hook.
- Instructions: a context entry (the operator's explicit ask) paid for by cuts elsewhere; one line in each agent and skill D13 names.
- A recall test (committed question set, a Go test behind an env var) that measures the index against a grep baseline and gates vectors.
- Sandboxed efficiency runs building sample apps, with and without the index, on both CLIs.
- The R4 edit, its ADR, the project `AGENTS.md` "stdlib only" restatement, docs.

Out:
- Vectors and an embedding model, unless the Phase 4 gate passes. If it passes, the model's delivery is a new path on the map and goes to the operator first (project rule 3); this plan then ends at the gate result.
- A service, watcher, daemon, MCP server or hook (R3, R8, project rule 1).
- A release. `/preflight` from the operator's own turn cuts it.
- Moving any record out of Markdown or JSONL. The databases are caches.

Follow-ups noted, not built: none yet.

## Surfaces touched

Recon brief: `plans/2026-09-29-recon-project-index.md`. Consumer audit: `plans/2026-09-29-audit-search-consumers.md`.

- `go.mod`, new `go.sum` — the one dependency.
- `internal/search/` — new package: open, reconcile, chunk, records, refs, query, global registry. Kept out of `internal/harness` so the installer's package does not import the driver.
- `cmd/sh-agent/main.go` — doc comment (lines 1-15), `usage` (lines 38-58), the first verb switch (before `harness.ResolveSource`, line ~138: searching needs no harness tree, and a release-binary or zip place has none), new flags in the shared FlagSet and in `reorder`'s `takesValue` (line ~710).
- `internal/harness/state.go:71-75` — `ConfigDir` comment names `search.db`.
- `internal/harness/report.go:235-295` — `Collect` names `search.db` as skipped, never sent.
- `context/AGENTS.md` — the search entry and the cuts that pay for it; `size_test.go` budget unchanged.
- `agents/*.md`, `skills/*/SKILL.md` — one line each in the consumers D13 names.
- `REQUIREMENTS.md` R4, `AGENTS.md` (project) "Stdlib only" bullet and Setup/Layout, `.agents/docs/architecture/DECISIONS.md` (ADR), `README.md` verb list, `.agents/docs/deploy/local.md` (first build fetches modules), `.agents/docs/specs/sh-agent/` (a `search.md` spec), `.agents/docs/testing/guide.md`, `.agents/docs/testing/skill-matrix.md` (rows for the changed agents and skills).

## Interfaces

### Commands

```
sh-agent search [--limit N] [--kind K] [--global] <words...>
sh-agent search --refs <BUG-NNN|ADR-NNN|path> [--kind K] [--global]
sh-agent search --open [--kind K] [--global]
sh-agent index [--global] [--rebuild]
```

- `<words...>` are joined with spaces and passed to FTS5 as a query of quoted terms ANDed (each word double-quoted, so `-`, `:` and `.` in `BUG-049` or `release.go` never reach FTS5 syntax). No hits for the AND query → retry once as OR and say so in the header line.
- `--limit` default 10; a value below 1 or above 50 is refused. At most 2 hits per file.
- `--kind` one of `bug adr plan change handoff review testing retro spec deploy trail code doc feedback proposal`. `feedback` and `proposal` exist only in the global index (Storage); `retro` covers both a trail `retros/` file and the global `retros.jsonl`. `trail` is every file under the docs root. It filters the hits of all three modes (`<words>`, `--refs`, `--open`) by the kind of the path each hit is printed at. For `--refs`, the ID's own record line is printed whatever the kind.
- `--refs X`: X matching `^(BUG|ADR)-\d{3}[a-z]?$` is an ID; anything else is a path, relative to the checkout root (a trailing `/` means a directory and matches every ref under it).
- `--open`: list records whose status is open — bugs with `status: open`, plans with an unticked box — one line each, `<path>:<line>  <kind> <id> <status> <date> <title>`, newest first, no limit. `--open` and `--refs` exclude each other and `<words>`.
- `--global`: query the global index instead. Every mode's output prefixes each path with the project's name, `<name>:<path>:<line>`; the files of `~/.config/sh-agent/` print as `~config:<file>:<line>`. `--refs X` with `--global` matches X in every project.
- `--rebuild` deletes this index's database files (the three names below, at this path only) and builds from scratch. `--global --rebuild` also drops the project registry, which lives in that file; each project re-registers on its next `search` or `index`, and the command prints how many registrations were dropped.
- `--limit`, `--kind`, `--refs` take values and join `takesValue` in `reorder`. `--global`, `--rebuild` are booleans. All five sit in the shared FlagSet, as every flag does today; the other verbs ignore them.
- Exit status: 0 with hits or with none (a "no hits" line is a result); 1 on any refusal or error.

Output of `search` (plain text, one hit per two lines; the snippet is one line, at most 200 characters, matched terms wrapped in `»«`):

```
search: 7 hits for "zip verify" (index current: 1,033 files, reconciled 4 changed in 38 ms)
.agents/docs/architecture/DECISIONS.md:355  adr ADR-039 Accepted; its one-file refusal superseded by ADR-041  2026-09-28
    … Release zips are signed, and `--from-zip` »verifies« them …
internal/harness/zip.go:88  code
    … // »verify« the »zip« against SHA256SUMS …
```

Output of `search --refs`:

```
ADR-040  .agents/docs/architecture/DECISIONS.md:383  Superseded by ADR-041  2026-09-28  An unsigned zip needs `--unverified`
referenced by 6:
.agents/docs/architecture/DECISIONS.md:404  adr ADR-041 Accepted  2026-09-28
    … supersedes ADR-040 …
REQUIREMENTS.md:37  doc
    … (refused unless `--unverified`; ADR-040) …
```
For a path, the first line is `<path>  tracked|modified|untracked|absent` and then the references; for a directory (trailing `/`), it is `<dir>/  directory`.

Output of `index`: one line per count — `files`, `chunks`, `records`, `refs`, `changed` (this reconcile), `bytes` (database + WAL), `ms`.

### Storage

Per project: `<root>/.agents/scratch/search.db` (+ `-wal`, `-shm`). Global: `~/.config/sh-agent/search.db` (+ `-wal`, `-shm`), honouring `--home`. Driver `modernc.org/sqlite`, opened as `sql.Open("sqlite", "file:<path>?_pragma=journal_mode(WAL)&_pragma=busy_timeout(10000)&_pragma=foreign_keys(0)")`.

```sql
CREATE TABLE meta(key TEXT PRIMARY KEY, value TEXT);            -- schema=1, docs_root
CREATE TABLE projects(root TEXT PRIMARY KEY, name TEXT, last_seen INTEGER); -- global only; per project: one row, itself
CREATE TABLE files(project TEXT, path TEXT, oid TEXT, state TEXT, PRIMARY KEY(project, path)); -- state: tracked|modified|untracked
CREATE TABLE blobs(oid TEXT PRIMARY KEY, last_seen INTEGER);
CREATE TABLE chunks(id INTEGER PRIMARY KEY, oid TEXT, line INTEGER, heading TEXT, body TEXT);
CREATE INDEX chunks_oid ON chunks(oid);
CREATE VIRTUAL TABLE chunks_fts USING fts5(heading, body, content='chunks', content_rowid='id', tokenize='porter unicode61');
CREATE TRIGGER chunks_ai AFTER INSERT ON chunks BEGIN
  INSERT INTO chunks_fts(rowid, heading, body) VALUES (new.id, new.heading, new.body);
END;
CREATE TRIGGER chunks_ad AFTER DELETE ON chunks BEGIN
  INSERT INTO chunks_fts(chunks_fts, rowid, heading, body) VALUES ('delete', old.id, old.heading, old.body);
END;
CREATE TABLE records(oid TEXT, id TEXT, line INTEGER, title TEXT, status TEXT, date TEXT);
CREATE INDEX records_id ON records(id);
CREATE TABLE refs(oid TEXT, line INTEGER, target TEXT);
CREATE INDEX refs_target ON refs(target);
```

Everything derived from content is keyed by the git blob ID (`oid`), so a blob seen on any branch is chunked once. `files` is the current view. A chunk row is never updated, only inserted and deleted; the two triggers are the only writers of `chunks_fts`, so the full-text index follows `chunks` on insert and on GC delete.

**Every query reads through `files`.** `search` joins `chunks_fts` → `chunks` (on `rowid = id`) → `files` (on `oid`, with `files.project` = the queried project; for `--global`, every registered project and `~config`). `--refs` joins `refs` → `files` on `oid`, and `--open` joins `records` → `files` on `oid`, the same way. A chunk, ref or record whose oid is in no current `files` row is never printed, which is what keeps another branch's content out of the answer. An oid current at two paths yields a hit at each path; the per-file cap applies per path. The printed path is always `files.path`.

Kind is derived from `files.path` at query time. Under the docs root, from its first directory: `bugs` → `bug`, `architecture` → `adr`, `plans` → `plan`, `changes` → `change`, `handoffs` → `handoff`, `reviews` → `review`, `testing` → `testing`, `retros` → `retro`, `specs` → `spec`, `deploy` → `deploy`, any other → `trail`. Outside it: `code` for extensions `go ts tsx js jsx mjs py sh rs java kt swift rb c h cc cpp cs sql tf`, else `doc`. `--kind trail` matches every kind under the docs root. The global index stores `feedback.jsonl`, `retros.jsonl` and `proposals/*.md` under project `~config` with their SHA-256 as `oid`, kind `feedback`, `retro`, `proposal`.

### Reconcile (before every query and every `index`)

1. `git rev-parse --show-toplevel` from the working directory → root. Failure → refuse (Failure modes).
2. `git ls-files -s -z` → path, mode, oid for tracked entries. Mode `120000` (symlink) and `160000` (submodule) are skipped.
3. `git status --porcelain=v1 -z --untracked-files=all` → modified and untracked paths (BUG-049: untracked files count). Their oids come from one `git hash-object --stdin-paths` call. Deleted paths leave the view.
4. A path whose content is binary (a NUL byte in its first 8 KiB) or larger than 1 MiB is skipped and not in `files`.
5. In one `BEGIN IMMEDIATE` transaction: replace `files` for this project with the new view; for every oid not in `blobs`, read the file and insert chunks, records and refs; set `last_seen` on every oid in the view; delete blobs (and their chunks, records, refs) not in any project's view whose `last_seen` is older than 14 days. After a transaction that changed anything: `PRAGMA wal_checkpoint(TRUNCATE)`.

### Chunking

- Markdown (`.md`): a chunk starts at each ATX heading; a chunk longer than 40 lines is split every 40 lines; `heading` is the nearest heading text.
- Everything else: windows of 40 lines, `heading` empty.
- `line` is the chunk's first line, 1-based. A hit's reported line is the first line in the chunk containing a matched term (found by a case-insensitive scan of the chunk for the query words), else the chunk's first line.

### Records (trail files only; docs root = `.agents/docs` if it exists, else `docs`)

Records and refs are parsed once per oid, when the oid is first indexed, using the first path at which the oid appears in that reconcile (`git ls-files` order, then `git status` order); the parser is chosen by that path. They are reported at every current path of the oid (D14).

- **Bug:** `<docs>/bugs/BUG-NNN*.md` → id from frontmatter `id:` else the file name's `BUG-NNN`; status from frontmatter `status:`; date from `date:`; title from the first `# ` heading.
- **ADR:** every line of `<docs>/architecture/DECISIONS.md` matching `^## (ADR-\d{3}[a-z]?): (.+?)\s+\((\d{4}-\d{2}-\d{2})\) — (.+)$` → id, title, date, status (the whole text after the dash, e.g. `Accepted; superseded in part by ADR-044 (…)`).
- **Plan:** `<docs>/plans/*.md` → id is the docs-relative path; status `open: N unticked` when it has N lines matching `^- \[ \]` outside fenced blocks, else `closed`; `line` is the first such line. Fences follow the handoff skill's scan (`skills/handoff/SKILL.md:154`): a run of three or more backticks opens a block, a line of only a run of the same length closes it.
- **Other trail files:** id is the docs-relative path, date from frontmatter `date:`, status empty.

### Refs (every indexed file)

- IDs: `\b(BUG|ADR)-\d{3}[a-z]?\b`, target the ID.
- Paths: tokens matching `[A-Za-z0-9_.@-]+(/[A-Za-z0-9_.@-]+)+/?` or `[A-Za-z0-9_-]+\.(go|md|ts|tsx|js|json|yaml|yml|sh|py|toml|mod|sum|sql|tf)`, with a trailing `:NN` or `:NN-NN` and trailing `.,;:)` removed and a leading `./` or `a/`, `b/` (diff prefixes) removed. Stored as written. A record's references to its own ID are kept.
- At query time `--refs X` matches targets equal to X, or equal to X without the docs-root prefix (the trail writes `plans/x.md` for `.agents/docs/plans/x.md`). A directory X (trailing `/`) matches targets with that prefix. No basename matching.

### Global registry

Every `search` and `index` run inside a git checkout, with or without `--global`, reconciles that checkout's own index and upserts `(root, base name, now)` into the global `projects` table first, so a project's first run being `--global` still includes itself. Run outside a git checkout, `--global` skips that step and does not refuse; plain `search` and `index` refuse there as Failure modes says. `--global` then reconciles `~config` from `~/.config/sh-agent/`, then each registered project: a root that no longer exists or is no longer a git checkout loses its registration and its rows (a note line says so); for the others, steps 2-5 limited to paths under the docs root. Two registered roots with the same base name are shown as `name`, `name~2` in registration order.

### Instructions (text, Phase 3)

Context entry, under `## Work: classify, then load the skill`, at most 330 bytes, final text chosen in Phase 3 and recorded as a D-row:

> **Search first:** in a git checkout, `sh-agent search <words>` ranks every file and trail record (`path:line`, status, snippet); `--refs <BUG-NNN|ADR-NNN|path>` lists every record naming it; `--open` lists open bugs and plans; `--global` spans projects. Search before opening trail files; exact strings, regex and past revisions stay with grep.

Cuts that pay for it (the rendered context must stay under 22,000 bytes, ADR-035; each cut names where the rule still lives):
1. `context/AGENTS.md:37`: the `find-skills`/`npx skills add` procedure shrinks to "a skill found elsewhere goes into the harness checkout's `skills/` (`sh-agent adopt`), then `sh-agent install`, never only into a vendor directory." The `npx` route stays in the project `AGENTS.md` Layout, where it applies.
2. `context/AGENTS.md:37`: "`/skills` and `/agents` show what a vendor loaded; when something is missing, check the symlinks first." moves to `skills/init-project/SKILL.md`, which creates the symlinks.
3. `context/AGENTS.md:31`: "Read the index before touching an area it names." becomes the `--refs` clause of the new entry and is deleted here.
If these do not free enough, Phase 3 stops and puts further cuts to the operator; the budget is never raised.

## Failure modes

| Case | What the user sees | Exit |
|---|---|---|
| Not inside a git checkout, or `git` missing | `search: not inside a git checkout; use grep` | 1 |
| `.agents/scratch/search.db` not gitignored (`git check-ignore -q` fails) | `search: .agents/scratch/search.db is not gitignored; add ".agents/scratch/*" to .gitignore (/init-project does)`; no file created | 1 |
| `.agents` or `.agents/scratch` is a symlink | `search: .agents/scratch is a symlink; refusing to write through it` | 1 |
| Another process holds the write lock past 10 s | `search: index busy (another sh-agent is indexing); retry, or use grep`; no stale answer is printed | 1 |
| Database unreadable (`SQLITE_CORRUPT`, `SQLITE_NOTADB`) or `meta.schema` ≠ 1 | `note: index at <path> was unreadable (or an older schema); rebuilt` and the answer | 0 |
| A file vanishes between `git status` and the read | the file is dropped from this reconcile; the next one picks up its state | 0 |
| `--refs` names an ID with no record | `<ID>  no record` then any references | 0 |
| `--limit` below 1 or above 50; `--kind` not in the list; `--refs` or `--open` with words, or with each other; `search` with no words and neither flag | `search: <what is wrong>; see sh-agent help` naming the valid values | 1 |
| The global index is rebuilt or found unreadable | the registry is gone with it; `note: N project registrations dropped; each returns on its next search` | 0 |
| Global: a registered root is gone | `note: dropped <root> (no longer a git checkout)` | 0 |
| The checkout is on a network filesystem | SQLite's own error, prefixed `search:`; documented as unsupported in the spec | 1 |
| `sh-agent report` | `search.db` is never collected; the dry run lists it under skipped | — |
| `sh-agent uninstall` | leaves `~/.config/sh-agent/search.db` (not in the manifest, R4); the spec says to delete it by hand | — |
| Checkout build with no module cache and no network (install's rebuild, BUG-006) | Go's own "cannot find module" error; `deploy/local.md` says the first build of a version needs the module proxy | 1 |

## Verification

Tiers: unit (Go tests with `-race`, temp git repositories), integration (the binary against this checkout and a scratch clone), behavioural (headless CLI runs per the project `AGENTS.md` § Testing), e2e (the sandboxed sample-app runs, Phase 5).

E2E criteria:
- In a checkout, after `git switch` to a branch that lacks a word, `sh-agent search <word>` returns no hit from the other branch; after switching back it returns the hit.
- An uncommitted edit and an untracked new file are found by the next search with no other command run.
- `sh-agent search --refs ADR-040` in this repository shows ADR-040's status as superseded by ADR-041 and lists ADR-041 among the references.
- `sh-agent search --global <a phrase from another project's bug record>` run in this repository returns that record.
- A fresh headless session on each CLI, asked where a past decision was made, runs `sh-agent search` before reading trail files (n of N reported).
- The rendered context stays under 22,000 bytes and ends with GANNET.

## Phases

- [x] Cold-reader check — design ready and implementation ready (round 1 `reviews/2026-09-29-coldread-search-index-r1.md`: two blocking gaps, fixed as D15 and D16; round 2 `reviews/2026-09-29-coldread-search-index-r2.md`: one blocking gap, fixed as D18; round 3 `reviews/2026-09-29-coldread-search-index-r3.md`: implementation ready, design needs revision on two non-blocking items) — cap reached; the operator, 2026-09-29: "Proceed with my fixes (Recommended)", resolved as D19, or the operator's same-turn override quoted with its date
- [x] R4 approved by the operator on the exact text in D1 — who: user — expect: the operator's own words quoted in D1 — approved 2026-09-29: "Approve as written (Recommended)", on the text shown verbatim

### Phase 0 — Baselines that need no dependency
- [x] Recall question set: `internal/search/testdata/recall.tsv`, at least 30 rows `question<TAB>gold path[,gold path…]`, taken from this repository's history, at least 10 of them drift questions (a claim and its restatements in other words, BUG-027/BUG-095/BUG-097 as models), at least 5 structured (`--refs`) questions — verify: AI, `awk -F'\t' 'NF!=2' internal/search/testdata/recall.tsv | wc -l` and every gold path `git ls-files --error-unmatch` — expect: `0`, and every gold path tracked — done 2026-09-29: 41 questions (26 ordinary, 10 drift, 3 refs, 2 open), `0`, every gold path tracked; structured questions carry a `refs X:` or `open K:` prefix and are scored over the whole list (D20)
- [x] Grep baseline scorer: the recall test's grep arm — for each question, the words (stop words removed) each run through `git grep -l -i -F`, files ranked by distinct words matched then total hits; recall@10 per question — verify: AI, `SH_AGENT_RECALL=1 go test ./internal/search -run TestRecall -v` prints the grep arm — expect: one line per question and a mean — done 2026-09-29: `testing/2026-09-29-recall-baseline.md`, mean recall 0.379 (ordinary 0.269, drift 0.353)
- [x] Token baseline: this project's transcripts classified by tool kind (the probe of 2026-09-29) re-run and saved as `testing/2026-09-29-token-baseline.md`, the probe embedded in the log (ADR-010: a procedure is documented, not committed as a script) — verify: AI, the file exists with both tables — expect: main and subagent tables matching the 2026-09-29 conversation's numbers within today's new sessions — done 2026-09-29: 23 main and 365 subagent transcripts; reads 33.8% and 59.2% of tool output
- [x] Checkpoint — expect: question set, scorer and baseline committed — evidence: `gofmt -l cmd internal` empty, `go vet ./...` exit 0, `go test -race ./...` exit 0 (four packages ok), 2026-09-29; verifier deferred to the Phase 1 range (D20)

### Phase 1 — Project index (needs R4 approved)
- [x] ADR-045 "The harness has one dependency, a pure-Go SQLite driver, for a derived search index"; `REQUIREMENTS.md` R4 per D1; `ROADMAP.md` Now; project `AGENTS.md` "Stdlib only" bullet and Setup (`go mod download` note) — one commit — verify: AI, `git show --stat` — expect: the four files in one commit
- [x] `go get modernc.org/sqlite@v1.60.1`, `go mod tidy`; `go.sum` committed; no vendor directory — verify: AI, `go build ./... && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/sh-agent && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o /dev/null ./cmd/sh-agent` — expect: exit 0 for all three — done 2026-09-29 at v1.46.1, not v1.60.1: v1.46.2 and later raise the module's `go` directive (D21); three targets built with `CGO_ENABLED=0 GOTOOLCHAIN=local`, and `/usr/bin/go` 1.24.4 builds it
- [x] `internal/search`: schema, open with the pragmas, schema check and rebuild — verify: unit, `go test -race ./internal/search -run 'TestOpen'` — expect: PASS, including a corrupt file rebuilt with the note and a schema-2 file rebuilt — done 2026-09-29: `TestOpenRebuildsCorruptAndOldSchema` PASS
- [x] Reconcile per Interfaces — verify: unit — expect: PASS for tests `TestReconcileBranchSwitch` (no hit from the other branch, for `search`, `--refs` and `--open` alike; switching back re-chunks nothing, asserted by a chunk-insert counter), `TestFTSFollowsGC` (after GC deletes a blob, an FTS query for its unique word returns no row and `INSERT INTO chunks_fts(chunks_fts) VALUES('integrity-check')` succeeds), `TestSameContentTwoPaths` (one oid at two paths: one chunk set, hits at both paths), `TestReconcileDirtyAndUntracked`, `TestReconcileDeletedFileLeaves`, `TestReconcileSkipsBinaryLargeSymlinkSubmodule`, `TestReconcileGC` (a blob unseen 15 days is gone, one seen 13 days ago stays) — done 2026-09-29: the five named tests plus `TestReconcileStagedRename` PASS; mutations M1, M2, M5, M6 (M6 equivalent, D22)
- [x] Chunking, records, refs per Interfaces — verify: unit — expect: PASS for `TestChunkMarkdownHeadingsAndWindows`, `TestRecordsBugADRPlan` (fixtures copied from this repository's real BUG file, `DECISIONS.md` heading shapes including "superseded in part" and "amends", a plan with 2 unticked boxes), `TestRefsIDsAndPaths` (a `file.go:12-30`, a `./x`, a docs-relative `plans/x.md`, a diff `a/` prefix) — done 2026-09-29: PASS; mutation M7
- [x] Query: search (quoted terms, AND then OR, per-file cap, kind filter, snippet, line), refs (ID, path, directory, docs-relative) and open — verify: unit — expect: PASS for `TestSearchRanksAndCaps`, `TestSearchOperatorCharactersAreLiteral` (`BUG-049`, `a.go`, `NOT`, `"`), `TestRefs` (including a directory target, whose first line is `<dir>/  directory`), `TestOpen` (an open bug listed, a fixed one not; a plan whose only `- [ ]` sits in a fenced block is closed, as the handoff skill's awk scan says) — done 2026-09-29: PASS; mutations M8, M9; refs grouped per file and record (D23)
- [x] Guards per Failure modes — verify: unit — expect: PASS for `TestRefuseOutsideCheckout`, `TestRefuseNotIgnored` (asserts no file created), `TestRefuseSymlinkedScratch`, `TestBusyRefusesNotStale` (a second connection holding `BEGIN IMMEDIATE`; the search exits 1 and prints no hit) — done 2026-09-29: PASS, plus `TestArgumentErrors`; mutations M3, M4
- [x] Concurrency: two processes reconciling the same checkout at once both succeed — verify: unit, `TestConcurrentReconcile` runs the test binary twice as subprocesses — expect: both exit 0, one index — done 2026-09-29: `TestConcurrentReconcile` PASS (202 files, one index)
- [x] Verbs in `cmd/sh-agent`: `search`, `index`, flags, `reorder` `takesValue`, doc comment, `usage` — docs: `README.md`, `.agents/docs/specs/sh-agent/search.md` — verify: unit, `go test ./cmd/sh-agent -run 'TestSearchVerb|TestReorderSearchFlags|TestHarnessMap|TestUsageMatchesDocComment'` — expect: PASS; `TestUsageMatchesDocComment` is new and asserts every verb in `usage` appears in the doc comment and back — done 2026-09-29: PASS; README rows and `specs/sh-agent/search.md` written; mutation M10
- [x] `report` skips `search.db`; `ConfigDir` comment — verify: unit, `TestCollectNeverSendsSearchDB` — expect: PASS, and the skipped list names it — done 2026-09-29: PASS; mutation M11
- [x] Integration on this checkout: `go run ./cmd/sh-agent index --rebuild`, then `search` for three trail questions, `--refs ADR-040`, `--refs internal/harness/zip.go` — verify: AI — expect: build under 10 s; a warm `search` with nothing changed under 100 ms end to end (`time`); ADR-040 shown superseded; `git status --short` clean afterwards — done 2026-09-29: build 2.2 s; warm search 23–67 ms end to end after D24; ADR-040 shown superseded by ADR-041; `--refs internal/harness/zip.go` 33 references; `git status` shows no `.agents/scratch` path
- [x] Checkpoint — expect: `gofmt -l cmd internal` empty, `go vet ./...` and `go test -race ./...` exit 0, verifier log PASS — evidence: `testing/2026-09-29-verify-search-index-phase-1.md` PASS over c04f2d7..7c24503 (198 pass, 2 skips with reasons; rebuild 2.1 s; warm search 25–32 ms, 5 of 5; three cross-builds; `go 1.24.4`), 2026-09-29

### Phase 2 — Global index
- [x] Registry and `~config` sources per Interfaces — verify: unit — expect: PASS for `TestGlobalRegistersProject`, `TestGlobalDropsVanishedProject`, `TestGlobalIndexesFeedbackRetrosProposals`, `TestGlobalNameCollision`, all with `--home` in a temp dir — done 2026-09-29: the four named tests plus `TestGlobalRunsOutsideACheckout` and `TestGlobalRebuildDropsRegistry` PASS; mutations G1–G5 each compiled and failed a named test
- [x] `--global` for `search` and `--refs` — docs: `.agents/docs/specs/sh-agent/search.md` — verify: integration, two scratch clones under `.agents/scratch/`, a bug record only in the second, searched from the first with `--home .agents/scratch/home-global` — expect: the record found with the second clone's name — done 2026-09-29: from `.agents/scratch/g1`, `search --global flamingo cache` found `g2:.agents/docs/bugs/BUG-001-flamingo.md:6  bug BUG-001 open`; the fixture was asserted first (g2 holds the word, g1 does not); spec updated
- [x] Checkpoint — expect: full test command exit 0, verifier log PASS — evidence: `testing/2026-09-29-verify-search-index-phase-2.md` PASS over bfcedcc..7fd56d5 (204 pass, 2 skips with reasons; `TestGlobal*` 18 of 18 across three `-race` repeats; hand integration found the other checkout's open bug), 2026-09-29

### Phase 3 — Instructions (context, agents, skills) and docs
- [x] Context entry and the three cuts (Interfaces § Instructions); final text as a D-row — verify: unit, `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock` and `go run ./cmd/sh-agent render --home .agents/scratch/ctx | tail -c 7` — expect: PASS, `GANNET` — done 2026-09-29: PASS at 21,988 of 22,000 bytes, `GANNET.` last; final text D31
- [x] One line in each consumer D13 names, naming the exact query it runs and that a past revision stays with `git grep <rev>` — verify: AI, `sh-agent install --home .agents/scratch/x` and `status` — expect: install succeeds, no drift — done 2026-09-29: 16 files; `install --home .agents/scratch/home-p3` exit 0, `status` drift none; the four lines that were weak in round 1 rewritten as ordered steps (D29)
- [ ] Docs: `README.md` (verbs, what the index is), `specs/sh-agent/search.md` (interfaces, failure modes, network filesystem limit, uninstall leaves the global file), `deploy/local.md` (first build fetches modules; `govulncheck` covers them), `testing/guide.md` (the recall test, the env var), `skill-matrix.md` rows — docs: those files — verify: AI, `/review --docs` over the range — expect: no open finding — in progress: README, spec, `deploy/local.md`, `testing/guide.md` § 3d, `sh-agent-matrix.md` and `skill-matrix.md` rows written 2026-09-29; `/review --docs` runs at Done; blocker: none
- [x] Behavioural runs, context: `claude -p --append-system-prompt <new context>` in a scratch clone of a small repository with a trail (fixture asserted by the run script before the CLI starts), prompt "where did this project decide X" — verify: both, 3 runs on Claude Code; must-fail control without the entry, 3 runs — expect: n of 3 run `sh-agent search` before any trail read with the entry, reported as n of N; the control's count reported — done 2026-09-29 (`testing/2026-09-29-search-index-behavioural.md`): with the final entry 2 of 3 ran `sh-agent search` as the first command; control 0 of 3; round 1 (entry without its ordering clause) 0 of 3
- [x] Behavioural runs, each changed agent: `claude -p --agents <json>` per project `AGENTS.md` § Testing — verify: both, 3 runs each — expect: n of 3 call `sh-agent search`, from the subagent transcript — done 2026-09-29: round 2 (final text, clones outside the checkout, BUG-106): `reviewer`, `security`, `explorer`, `dreamer` each 3 of 3 ran it; `documentation` 3 of 3 fixed the stale doc, `bugs` 3 of 3 cited BUG-001, `changelog` 3 of 3 cited BUG-001 or ADR-003 from the dispatched output
- [x] Antigravity runs: sandbox procedure `testing/2026-09-03-verify-sandbox.md` with a run-specific `settings.json` granting `command(sh-agent)`, `command(git)` — verify: both, 1 run of the context question and 1 of each changed agent — expect: no "auto-denied" line, no "invalid allow entry" line; whether search was called, reported — done 2026-09-29 (`testing/2026-09-29-search-index-behavioural.md`): no auto-denied line; completed subagent runs searched (`reviewer` 2 of 2, `security` 1 of 1, `explorer` 1 of 2, final wording 1 of 1); `dreamer` and one `security` run void (print mode ended before the subagent returned); context question 0 of 2 on agy
- [ ] Real-home smoke per `testing/2026-09-03-verify-real-home.md` after the operator installs — who: user — verify: the three smoke checks on each CLI — expect: PASS, GANNET present
- [ ] Checkpoint — expect: reviewer and security reports with `## Outcome`, verifier PASS, `/review --docs` clean — evidence: <paths>, <date>

### Phase 4 — Recall test and the vector gate
- [x] Index arm of `TestRecall`: each question through `search` with `--limit 10` (structured questions through `--refs`) — verify: AI, `SH_AGENT_RECALL=1 go test ./internal/search -run TestRecall -v` — expect: recall@10 per question and the mean for both arms, saved to `testing/<date>-search-recall.md` — done 2026-09-29: `testing/2026-09-29-search-recall.md`: index 0.653 against grep 0.379 overall, 0.692 against 0.269 on ordinary questions, 0.376 against 0.353 on drift
- [x] Gate decision: vectors proceed only if a probe of BM25 fused with vectors (all-MiniLM-L6-v2, 384 dimensions, embedded by the pure-Go `github.com/clems4ever/all-minilm-l6-v2-go`, stored with `modernc.org/sqlite/vec`, fused by reciprocal rank fusion with k=60, in a separate module under `.agents/scratch/` only, never committed) beats the index arm plus the grep arm by at least 10 points of mean recall@10 — verify: AI — expect: the numbers and the decision as a D-row; a pass becomes an `Open:` row for the operator (model delivery), a fail closes vectors — done 2026-09-29: not passed (D30); vectors closed
- [x] Checkpoint — expect: recall log committed, D-row written — evidence: `testing/2026-09-29-search-recall.md`, D30, 2026-09-29

### Phase 5 — Sandboxed efficiency runs (after Phases 1-4 are committed)
- [ ] Kit in `.agents/scratch/eff-kit/` (not committed, ADR-010): per-arm sandbox homes, baseline arm installed from the v1.12.1 release binary, index arm from this checkout; run-specific agy `settings.json`; the run script asserts each fixture fact before the CLI starts — verify: AI, a dry run of the kit — expect: both homes install with no drift; `sh-agent search` absent from the baseline home's context and present in the index arm's
- [ ] Two sample apps, three sessions each in one project: session 1 builds the app through `/feature` (gates waived by name in the prompt, as in the 2026-09-16 simulation), session 2 (fresh) adds a feature, session 3 (fresh) fixes a seeded bug through `/bug`. App A: linkshelf (the 2026-09-16 brief, verbatim). App B: a Go CLI `tasklog` (add, list, done, JSON file store, tests). Claude Code: 2 runs per arm per app. Antigravity: 1 run per arm per app — verify: both — expect: every run ends with tests passing and a handoff, or the failure is recorded
- [ ] Larger-project task: in a scratch clone of this repository at v1.12.1, a fresh session per arm answers five trail questions and fixes one seeded doc drift; Claude Code 2 runs per arm, Antigravity 1 — verify: both — expect: answers scored against gold
- [ ] Metrics from transcripts per session: input tokens (fresh, cache read, cache write), output tokens, tool calls by kind, read volume, `sh-agent search` calls, subagent tokens, wall time, success; Antigravity reports what its logs hold and names what they do not — verify: AI — expect: `testing/<date>-search-efficiency.md` with medians and ranges per arm, n of N, and what was not measured
- [ ] Checkpoint — expect: the efficiency log committed — evidence: <path>, <date>

### Done
- [ ] Definition of Done: verifier over the whole range, reviewer (diff only), security (a dependency and input parsing are touched), documentation, `/review --docs`, changes entries, bug records for anything that went wrong, plan closed, `ROADMAP.md`/`TODO.md` emptied, commit, push, handoff — verify: AI — expect: every artifact on disk and pushed

## Risks

Assumptions:
- `modernc.org/sqlite` v1.60.1 builds and runs FTS5 on linux/amd64 and cross-compiles for darwin/arm64 and linux/arm64 with `CGO_ENABLED=0` — checked 2026-09-29 for linux/amd64 (run) and darwin/arm64 (build); linux/arm64 checked by the Phase 1 build task.
- Claude Code subagents receive the global context — observed by the audit's own dispatch on 2026-09-29 (`plans/2026-09-29-audit-search-consumers.md`); Antigravity subagents unverified, which is why D13 puts a line in each consumer's own file; Phase 3's agy runs report it.
- `git check-ignore` reports `.agents/scratch/search.db` ignored in every project scaffolded by `/init-project` — checked by recon (`skills/init-project/SKILL.md:110-122`); Phase 1's guard test covers a project where it is not.
- The Linux box's `claude` and `agy` can run headless in a sandbox home with a copied token — checked by the 2026-09-03 and 2026-09-16 procedures; re-checked by Phase 5's kit dry run.
- Phase 5 cost: about 30 headless sessions of up to two hours each. Stated to the operator with the R4 question.

## Decisions

D1 — 2026-09-29 — `REQUIREMENTS.md` R4 changes to the text below — by: the operator, 2026-09-29, in their own turn, choosing "Approve as written (Recommended)" on this exact text, with the no-vendoring trade-off (D9) and the Phase 5 cost stated in the question.

> ## R4 — One binary, one dependency, easy on and off
>
> `sh-agent` is a single Go binary. Its one dependency is `modernc.org/sqlite`, a pure-Go SQLite driver used only by the search index, pinned in `go.mod` and `go.sum` and not vendored; any other dependency, and any driver that needs cgo, is a requirements change. The search indexes (`.agents/scratch/search.db` in a project, `~/.config/sh-agent/search.db`) are caches derived from files git or `~/.config/sh-agent/` already holds: never committed, never sent by `report`, rebuilt when deleted. `install` is idempotent … *(the rest of R4 unchanged)*

D2 — 2026-09-29 — The databases are derived indexes; Markdown and JSONL stay the source of truth — by: the operator, "Go with your recommendations" (row 1).
D3 — 2026-09-29 — No service: reconcile against git at query time — by: the operator (row 2).
D4 — 2026-09-29 — Per-project file under `.agents/scratch/`; named `search.db`, not `index.db`, because R7 already calls `retros.jsonl` "the cross-project index" — by: the operator (row 3); name by implementer from recon.
D5 — 2026-09-29 — A global index at `~/.config/sh-agent/search.db` over feedback, retros, proposals and registered projects' trails — by: the operator (row 4).
D6 — 2026-09-29 — Phase 1 is full-text search plus the metadata graph; vectors only through the Phase 4 gate — by: the operator (row 5).
D7 — 2026-09-29 — Instructions go in the context, paid for by cuts, plus one line per HIGH-ranked agent or skill; no new skill — supersedes recommendation row 6's "a new skill" — by: the operator, "make sure to add to the context how to use the database".
D8 — 2026-09-29 — The embedding model's delivery is deferred to the Phase 4 gate — by: the operator (row 7).
D9 — 2026-09-29 — Not vendored: `go mod vendor` measured 137 MB (modernc `libc` 71 MB, `sqlite` 57 MB), which would enter every clone and every release zip; a checkout's first build of a version needs the module proxy — by: implementer, measured 2026-09-29; stated to the operator with D1.
D10 — 2026-09-29 — The verbs sit in the first switch, before `ResolveSource`, so a release-binary or zip place with no harness tree can search — by: implementer, from recon.
D11 — 2026-09-29 — Phase 5 adds a larger-project task (a clone of this repository) beside the two sample apps, because the operator's reason is "a powerhouse in larger projects" and a new sample app has almost no trail — by: implementer, from the operator's words.
D12 — 2026-09-29 — Snippet text, not whole files: `search` prints at most 200 characters per hit and at most 50 hits, so a search costs a few hundred tokens — by: implementer, from the token baseline.
D14 — 2026-09-29 — An oid current at several paths is parsed once, by the first path seen; its records and refs are reported at every current path — by: cold read round 1 (`reviews/2026-09-29-coldread-search-index-r1.md`, readiness question 4, its default).
D15 — 2026-09-29 — `chunks_fts` is kept in sync by two triggers on `chunks` (insert, delete); chunks are never updated — by: cold read round 1, blocking question 1.
D16 — 2026-09-29 — Every query joins through `files` for the queried project(s); nothing keyed only by oid is printed without that join — by: cold read round 1, blocking question 2.
D17 — 2026-09-29 — The directory-to-kind table and the driver's DSN are written into Interfaces; the Phase 3 context edit measures the budget first — by: cold read round 1, readiness questions 3, 5 and 6.
D18 — 2026-09-29 — Round 2 fixes: `feedback` and `proposal` join the `--kind` list; a `--global` run inside a checkout registers and reconciles that checkout first; argument errors and the registry's loss on a global rebuild are Failure modes rows; the vector probe names its library and fusion — by: cold read round 2 (`reviews/2026-09-29-coldread-search-index-r2.md`) and implementer's own read before round 3.
D19 — 2026-09-29 — Round 3 residue: R7 already calls `retros.jsonl` "the cross-project index", so every doc, comment and message calls the new file "the global search index" and never "cross-project index"; R7 is not edited. `--refs <dir>/` prints `<dir>/  directory` first. D13's dream citation is line 41 — by: the operator, "Proceed with my fixes (Recommended)", on the cap question.
D20 — 2026-09-29 — Phase 0 adds only a measurement test that skips unless `SH_AGENT_RECALL=1`, so its verifier run is the Phase 1 dispatch, whose range starts before Phase 0. Reference and open questions are scored over the complete list, not the top ten, because those modes return exhaustive lists — by: implementer.
D21 — 2026-09-29 — `modernc.org/sqlite` is pinned at v1.46.1 (2026-02-18), the newest release whose `go.mod` says `go 1.24.0`: `go get` of v1.60.1 raised this module's `go` directive to 1.26.0, which ADR-029 and the project `AGENTS.md` forbid, and a Go 1.24 machine building with `GOTOOLCHAIN=local` could no longer build — by: implementer, measured.
D22 — 2026-09-29 — The Phase 1 tests were written after the code, against the feature skill's test-first rule. Compensation: every key test was checked by the mutation protocol (review-standards § 2): eleven mutations, each compiled; ten failed their named tests; M6 (skipping a rename's original path removed) passed, an equivalent mutation, because the unskipped original path is parsed as a nonexistent path and dropped. The rename test itself pins the path order: with the skip in place it passes only when `git status -z` lists the new path first. Recorded in the Phase 1 change entry — by: implementer.
D23 — 2026-09-29 — `--refs` groups references per file and record, not per file, so ADR-041's "supersedes ADR-040" is its own entry and not folded into ADR-040's heading; groups are ordered newest record date first, then path; a group naming the target more than once shows `×N`. The schema adds `blobs.skipped` (a skipped binary or large file is not re-read on every reconcile) and `projects.docs_root` (the global search index classifies each project's paths) — by: implementer, from `TestRefs` and the design.
D24 — 2026-09-29 — Word search materializes the FTS5 match first (`WITH m AS MATERIALIZED`) and makes snippets only for printed hits: joined directly, SQLite drove the loop from `files` and evaluated the match per row, 0.6–2.5 s per query on this repository; now 23–67 ms — by: implementer, measured.
D25 — 2026-09-29 — `govulncheck` reports GO-2026-5024 in `golang.org/x/sys` v0.37.0 (pulled in by the driver), not reachable from this code. Its fix, v0.44.0, needs `go 1.25.0`, which D21 rules out, so it stays until the module's `go` directive moves; the Done review names it — by: implementer, from the Phase 1 verifier log.
D26 — 2026-09-29 — `documentation`, `bugs` and `changelog` have no shell: the main agent runs `sh-agent search --refs <changed paths>` and puts the output in their dispatch, and each of the three files says to use it; their tool lists do not change — by: the operator, "Main agent runs it (Recommended)".
D27 — 2026-09-29 — Word search ranks chunks holding every word first, then fills the remaining slots with chunks holding any word, and the header says how many had every word; supersedes the Interfaces bullet "No hits for the AND query → retry once as OR". The recall test's first index run scored 0.077 on ordinary questions because an AND of many words matched two chunks and stopped — by: implementer, from the recall test.
D28 — 2026-09-29 — `--refs <bare file name>` that matches no path prints `hint: did you mean --refs <path>` for up to five tracked paths ending in that name; still no basename matching. Found in the behavioural runs: two reviewer runs passed `store.go` and got nothing — by: implementer.
D29 — 2026-09-29 — Behavioural runs go outside every checkout (`~/sandbox/search/`, BUG-106), and the lines that round 1 showed weak are ordered steps: `explorer` "your first command, before any Grep, Glob or Read"; `reviewer` a first pre-flight step with the full path; `dreamer` "run `sh-agent search` too (the Grep tool replaces neither)"; the context entry "Run it before Grep" (its "in a git checkout" dropped to pay for it) — by: implementer, from round 1 (0 of 3 for `explorer`, `dreamer` and the context; 1 of 3 for `reviewer`).
D30 — 2026-09-29 — Vector gate not passed: fusing MiniLM vectors into the index adds 6.4 points of mean recall@10 overall and 3.0 on drift against the index alone, under the 10-point bar; the plan's literal baseline, the RRF of index and grep, scores below the index alone, and a gate a weaker baseline passes measures the baseline. The delivery condition fails too: the "pure Go" MiniLM library loads ONNX Runtime as a shared library. Vectors are closed; the operator can reopen them — by: implementer, from `testing/2026-09-29-search-recall.md`.
D31 — 2026-09-29 — The context entry, final: "**Search first:** `sh-agent search <words>` ranks files and trail records; `--refs <BUG-NNN|ADR-NNN|path>` lists every record naming it (run it before touching an area); `--open` lists open bugs and plans; `--global` spans projects. Run it before Grep; grep stays for exact strings, regex, past revisions." — by: implementer, measured at 21,988 bytes.
D32 — 2026-09-29 — Security round 1 (`reviews/2026-09-29-security-search-index.md`): the DSN is built with `net/url`; only a checkout with `.agents/docs` registers in the global search index (supersedes the Interfaces § Global registry line "Every `search` and `index` run inside a git checkout … upserts"); modified and untracked files are hashed from their own bytes (no `git hash-object`); printed text carries no control characters; directories the index creates are 0700; secret-shaped names are never indexed (the installer's list, exported as `harness.SecretShaped`) — by: security review, each fix with a test that failed first.
D33 — 2026-09-29 — Docs review (`reviews/2026-09-29-review-8ee5d70-HEAD-docs.md`): ADR-045's version and module count corrected by addendum ADR-045a; `go mod tidy` makes the driver a direct requirement (the `go` directive stays 1.24.4); the project `AGENTS.md` names the known `x/sys` advisory — by: docs review.
D34 — 2026-09-29 — Security round 2: registration stays automatic for a checkout with `.agents/docs`; the residual (a crafted clone that commits `.agents/docs` and that the operator runs `search` in joins `--global` and `/retro` evidence until `index --global --rebuild`) is accepted — by: the operator, "Keep automatic (Recommended)". The global search index reads only `.agents/docs` and drops a registration that lost it; hint and note lines are printable; hashed bytes are the indexed bytes; the secret list ignores case and adds `.git-credentials` and `credentials` (the installer's list and `.gitignore` too) — by: security review round 2.
D13 — 2026-09-29 — Consumers that get their own line, from the audit (`plans/2026-09-29-audit-search-consumers.md`): `skills/review-standards/SKILL.md:48` and `agents/dreamer.md:36` (restatements: `search <key words>` beside the literal grep), `skills/dream/SKILL.md:41` (graduation check), `skills/handoff/SKILL.md:154` (`search --open --kind plan` first, the awk scan when search refuses), `agents/dreamer.md:29`, `skills/bug/SKILL.md:33`, `skills/feature/SKILL.md:19` (`search --refs <path>` and `--open --kind bug` instead of reading all of `bugs/INDEX.md`), `skills/retro/SKILL.md` (`search --global` for "happened in another project"). The context's restatement rule (`context/AGENTS.md:55`) is covered by the new entry. Added on the operator's question "why not all places that generate lots of md files?" — the audit ranked steps that already grep, and the writers have none because they write without looking: `agents/documentation.md` (`search <changed command or flag>` for every doc that describes it, before editing), `agents/bugs.md` (`search <symptom words> --kind bug` and `--refs <path>` before assigning a new number, to find a duplicate or an existing prevention rule), `agents/changelog.md` (`--refs <path>` for the bug and ADR IDs a change entry should cite; `search` for the ADR a new decision supersedes), `agents/explorer.md` (`search` and `--refs` as its first recon step), `agents/reviewer.md` and `agents/security.md` (`--refs <path> --kind bug` for the files in the diff: past failures, never the author's intent), `skills/research/SKILL.md` Step 1 (`search <topic> --kind plan` for earlier research). Not to use it, and their lines unchanged: `agents/verifier.md` (live-tree, diff-scoped checks), `skills/preflight/SKILL.md:40` (secret scan), the retro skill's exact `rule_ref` greps, and any grep at a past revision (`/review --docs` at the scope's head) — by: implementer, from the audit.
Rejected: vendoring (D9); a watcher or service (D3); a trigram tokenizer (23 MB here against 3.4-12 MB for porter; exact substrings stay with grep); basename matching in `--refs` (`main.go` would match every command); a `--json` output (no consumer named); a separate skill (D7); storing records in the database as the source (D2).
