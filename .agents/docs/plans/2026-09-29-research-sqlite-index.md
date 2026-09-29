---
owner: main
date: 2026-09-29
---
# Research — an embedded SQLite index (FTS5 + sqlite-vec) for the agent   (2026-09-29)

## Question

The operator asked, 2026-09-29, through `/research` (full mode): should the harness give the agent a local index inside a Git-tracked project, an embedded SQLite file using FTS5 for lexical search and sqlite-vec for dense vectors? Constraints: no daemon or server, a small storage footprint, and no binary index churn in Git. Wanted: a verdict (ADOPT, REJECT or CONDITIONAL); the boundary against Git and how branch switches, rebases and merges invalidate it; retrieval quality against LanceDB and against ripgrep, git grep and ctags; binary distribution and cross-compilation; failure modes; a decision matrix; and an exact adoption rule.

## Answer in brief

**Verdict: CONDITIONAL. SQLite is the right engine if the harness ever needs an index, but the harness does not need one today, so build nothing.**

- **The engine question is settled in SQLite's favour.** A pure-Go driver (`modernc.org/sqlite` v1.60.1) ships FTS5 and bundles sqlite-vec v0.1.9. Both were verified on this box: no cgo, a clean cross-compile to darwin/arm64, a single file, no daemon. LanceDB needs cgo over a Rust library through a thin Go SDK, and it grows a new version on every write.
- **The index question fails on this project's own numbers.** A full-tree `git grep` takes about 11 ms here. A query-time reconcile against Git (4 ms) plus an FTS5 query (under 1 ms) saves about 6 ms per search. It costs a dependency that R4 forbids, a new `sh-agent` verb, a cache that goes stale without hooks (R3 forbids them), and 3 to 23 MB of disk per checkout.
- **Dense vectors add a second, larger cost.** An embedding model must run somewhere. In-process means a model file of tens to hundreds of MB, and cgo for the fast runtimes. Otherwise it needs a daemon (Ollama) or a network API. Anthropic reports that Claude Code dropped its early local vector database for agentic grep. Published 2026 benchmarks find a small gain from fusing embeddings with lexical search, and no single retriever that wins.
- **The adoption rule is in § 4.** Adopt FTS5 when a measured search latency or recall failure exists and the operator approves the R4 edit. Add sqlite-vec only when a fixed query set on the project shows that fusion beats FTS5 plus grep by a stated margin.

## Measured on this repository (2026-09-29, Linux amd64, 16 cores, warm cache, go1.27.1)

Corpus: 1,033 tracked files, 7.9 MB of text. 859 of the files sit under `.agents/`, the trail.

| Measurement | Result |
|---|---|
| `git grep -n 'BUG-0'`, whole tree | 11 ms |
| `grep -rl` over every tracked file | 21 ms |
| `rg --hidden` / `git grep`, literal (web agent's run) | 13–14 ms / 10–12 ms |
| `git ls-files -s` (index blob IDs, 1,033 entries) | 1.6 ms |
| `git status --porcelain=v2` | 4 ms |
| FTS5 build, porter and trigram tables, Python sqlite 3.46.1 | 0.98 s |
| FTS5 BM25 query (three phrase and term queries) | 0.1–0.4 ms |
| FTS5 index, porter, content stored | 12.3 MB |
| FTS5 index, porter, contentless (`content=''`) | 3.4 MB |
| FTS5 index, trigram, content stored | 31.8 MB |
| FTS5 index, trigram, contentless | 22.9 MB |
| Go binary, stdlib only (baseline) | 2.3 MB |
| + `modernc.org/sqlite` (FTS5 trigram works) | 9.9 MB, 25 modules in the graph |
| + `modernc.org/sqlite/vec` (`vec_version()` = v0.1.9) | 10.7 MB |
| Cross-compile of both, `GOOS=darwin GOARCH=arm64 CGO_ENABLED=0` | succeeds |
| sqlite-vec, 20,000 random 384-dim vectors, insert | 0.41 s |
| sqlite-vec, brute-force KNN k=10 over those vectors | 8–9 ms |
| sqlite-vec file after that insert | 32 MB database + 32 MB WAL, not yet checkpointed |
| `asg017/sqlite-vec-go-bindings/ncruces` v0.1.6, its README example | panics at startup ("i32.atomic.store invalid as feature threads is disabled"); fails to compile against current `ncruces/go-sqlite3` v0.35.6 |

The probe files were deleted after the run. The chunk count for real embeddings of this corpus was not measured. At 1 to 1.5 KB per chunk it would be roughly 5,000 to 8,000 vectors, which is under the 20,000 measured above.

## 1. Architectural role and Git hygiene

**The boundary.** Git (HEAD plus the working tree) is the only source of truth. The database is a derived, disposable cache: it is never committed, never read by a human, never the only copy of anything, and deleting it loses nothing but rebuild time. That rules out the "persistent storage" role for it. Anything worth keeping (a decision, a bug, a summary) belongs in a Markdown file in the trail, where Git versions it and review sees it. A cache that holds facts not in Git is a second memory store, which the 2026-09-29 memory research already rejected.

**Where it lives.** Under `.agents/scratch/index/`, which is already gitignored in every project `/init-project` scaffolds, so nothing new reaches `.gitignore`. It is not in `.git/`, which Git owns, and not in `~/.config/sh-agent/`, which R7 reserves for harness state. `git clean -fdx` deletes it, which is correct for a cache. A worktree gets its own index, because a worktree is a different working tree.

**Consistency across checkout, rebase and merge.** Key every indexed chunk by the Git blob ID of its file's content, never by path or mtime.

- The file table maps path to blob ID for the current view: the index blob IDs from `git ls-files -s` (1.6 ms here), plus a `git hash-object` of each path that `git status` reports dirty (4 ms here).
- Chunks and vectors are stored once per blob ID. A branch switch, rebase or merge changes the path-to-blob map. It re-chunks only blobs the database has never seen; switching back to a branch already visited costs nothing.
- Rebases rewrite commits but mostly keep blobs, so they are nearly free. A merge creates new blobs only for files it changed.
- Blobs no path references any more are garbage-collected on a size threshold, not on every switch.

**Invalidation is done at query time, never by a watcher.** There is no daemon (the constraint) and no hook (R3), so nothing learns about a checkout when it happens. Every query first runs the reconcile above, about 6 ms here, and indexes the difference before it answers. An index that answers without reconciling returns results from another branch. That is the one failure the design must make impossible, not merely unlikely.

**Why not mtime.** `git checkout` sets mtime to the checkout time, a same-second edit can keep the recorded mtime (Git's own "racy git" problem), and a copy with `cp -p` or a restore from an archive keeps old mtimes on new content. Blob IDs have none of these problems.

## 2. Retrieval capability against overhead

### What each option retrieves

| Query shape | grep / git grep / rg | ctags / repo map | FTS5 BM25 | FTS5 + sqlite-vec | LanceDB |
|---|---|---|---|---|---|
| Exact identifier, error string, BUG number | exact, fastest | symbols only | exact (trigram) or token | same as FTS5 | BM25 via Lance FTS |
| Regex (`func [A-Z]\w+Test`) | yes | no | no | no | no |
| Ranked "most relevant files for these words" | no ranking; the agent reads hit counts | ranks by reference graph | BM25 | BM25 fused with vector score | BM25 + vector |
| Paraphrase with no shared words ("where do we refuse unsafe archives" for the zip guard) | agent must guess synonyms | no | no, unless stems match | yes, the one real gain | yes |
| Definition and callers of a symbol | approximate | yes (ctags defs; Aider-style map adds refs) | no | no | no |

### Against LanceDB (a file-based columnar vector format)

- **Distribution.** The official Go SDK is v0.1.2 (2025-09-30), 61 stars. It is cgo over a prebuilt Rust static library you download separately. A cross-compile to darwin/arm64 from Linux needs that library for each target plus a C cross toolchain. The harness today ships a stdlib binary built with `go build`.
- **Growth.** Every write creates a new table version. Old versions stay until `optimize()` prunes them, after seven days by default. An index updated on every query-time reconcile writes often, so it grows until something compacts it.
- **What it wins.** Columnar storage and ANN indexes pay off from millions of vectors, and on multimodal data. Neither applies to a code and Markdown corpus of thousands of chunks.
- **Conclusion.** No property of this corpus needs LanceDB, and its distribution cost is the highest of the three.

### Against plain CLI tools (ripgrep, git grep, ctags)

- **Latency.** 11 to 21 ms for the whole tree here. The Linux kernel, about 70,000 files, searches in about 0.35 s (burntsushi.net/ripgrep, 2016). Grep latency stops mattering to an agent whose turn takes seconds only at kernel scale and above.
- **Freshness.** Grep reads the working tree, so it is never stale and has no invalidation logic. This is the property Boris Cherny names for Claude Code: "security, privacy, staleness, and reliability" (x.com/bcherny/status/2017824286489383315, 2026-02-01).
- **Availability.** `git grep` ships with Git. Both CLIs already expose a grep tool to the agent. `rg` is not installed on this box, but Claude Code bundles its own.
- **ctags.** It adds symbol definitions. For the Go code here, `go doc` and `gopls` already cover that, and the Go source is 31 files, 18 of them non-test.
- **What grep lacks** is ranking and paraphrase. The agent compensates by reformulating and reading, which costs tokens and turns, not latency.

### Binary distribution and process overhead

| Route | cgo | Cross-compile to darwin/arm64 from Linux | Binary delta | Status |
|---|---|---|---|---|
| `modernc.org/sqlite` + `/vec` | no | plain `go build` (verified) | +8.4 MB | FTS5 built in; sqlite-vec v0.1.9 bundled since 2026-03-15; verified working |
| `mattn/go-sqlite3` + sqlite-vec cgo bindings | yes | needs a C cross toolchain and a macOS SDK (osxcross or zig) | not measured | FTS5 only with a build tag |
| `ncruces/go-sqlite3` | no | plain `go build` | not measured | FTS5 opt-in per connection; own `ext/vec1` (IVF, "testing is insufficient"); the sqlite-vec binding for it is broken (measured) |
| LanceDB Go SDK | yes, over Rust | per-target prebuilt library + C toolchain | not measured | v0.1.2, thin |

Process overhead: SQLite runs inside the calling process, so there is nothing to start or keep alive. An embedding model is different. In-process options are a pure-Go MiniLM at about 4.5 ms per embedding (github.com/clems4ever/all-minilm-l6-v2-go), `hugot` on a pure-Go backend or on ONNX Runtime through cgo, or llama.cpp loaded as a shared library through purego. Each adds a model file: bge-small-en-v1.5 is 33M parameters, and jina-embeddings-v2-base-code is 161M. The alternatives are Ollama, which is a daemon, or a hosted API, which sends source off the machine.

## 3. Failure modes and bottlenecks

- **Concurrent access.** The harness runs parallel subagents and, on one machine, several sessions per checkout. WAL mode allows many readers and one writer across processes on one host. Every reconcile write must use `BEGIN IMMEDIATE` with a `busy_timeout`. A read transaction that later tries to write after another process committed fails with `SQLITE_BUSY_SNAPSHOT`, not a wait. WAL does not work on a network filesystem (sqlite.org/wal.html), so a checkout on NFS or a synced folder must fall back to grep. A reader inside a long read transaction also sees a snapshot older than the working tree another session just edited.
- **The verifier freeze.** The Definition of Done freezes the tree while the `verifier` runs. An index that writes at query time does not touch the tracked tree, so the freeze holds. An index stored inside the tree but outside `.agents/scratch/` would dirty `git status` and break the `sh-agent update` dirty-tree refusal.
- **Chunk invalidation on rapid edits.** The file is the unit: an edit re-chunks its file. For FTS5 that is sub-millisecond. For vectors it is the chunk count times the embed time, for example 50 chunks times 4.5 ms, about 0.23 s per file, paid on the next query after each edit. An agent that edits a file between every search pays it every time. AST-aware chunking adds a parser: the maintained tree-sitter Go bindings use cgo, and `go/ast` covers Go only. Chunking by heading and by blank line covers Markdown, which is 859 of 1,033 files here.
- **Embedding drift.** Vectors from two models, two model versions or two chunkers are not comparable, and mixing them degrades results silently rather than failing. The database needs a meta row (model ID, dimensions, chunker version) that forces a full drop and rebuild on any mismatch. A model upgrade becomes a full re-embed of every blob in the cache, not only of HEAD.
- **File growth.** FTS5 accumulates segments and tombstones under incremental updates until `optimize` runs. vec0 deletes do not shrink the file without `VACUUM`. The WAL grows until checkpointed: the probe left a 32 MB WAL beside a 32 MB database. Blob-keyed caching keeps every branch ever visited unless garbage-collected. The simplest bound: rebuild from scratch past a size threshold, since a full FTS5 build of this tree takes about one second.
- **Stale answers.** An index is wrong silently. Grep that finds nothing is honest about the tree it read. An index that skipped the reconcile, or that crashed half-way through one, answers from the past with the same confidence. This is the harness's own rule on silent instrumentation: "silence from your own instrumentation is not evidence".
- **Two stores of the same fact.** The trail already has curated indexes: `bugs/INDEX.md`, `DECISIONS.md`, and date-ordered file names. A search index over the trail competes with them for "where do I look first".

## 4. Recommendation and decision matrix

### Decision matrix

Scores run 1 to 5, higher is better, scored for this harness at this scale.

| Option | Latency | Portability | Git cleanliness | Search quality (lexical + semantic) | Maintenance overhead (5 = least) | Fits R3/R4 today |
|---|---|---|---|---|---|---|
| A. Agent-driven `git grep` / `rg` (today) | 4 (11–21 ms, never stale) | 5 (ships with Git and both CLIs) | 5 (no artifact) | 3 (exact and regex strong; no ranking, no paraphrase) | 5 | yes |
| B. A + ctags or a Go symbol map | 4 | 4 (universal-ctags not installed here) | 4 (a tags file to ignore) | 3.5 (adds definitions) | 4 | yes |
| C. SQLite FTS5 via `modernc.org/sqlite`, reconcile at query time | 5 (under 1 ms query + 6 ms reconcile) | 5 (pure Go, verified cross-compile) | 5 (gitignored scratch) | 3.5 (BM25 ranking; trigram for substrings) | 3 (invalidation code, growth, locks) | no: a dependency (R4); a new verb |
| D. C + sqlite-vec + an in-process model | 3 (embed on each changed file; 9 ms KNN at 20k) | 3 (pure Go possible; model file per machine; fast runtimes need cgo) | 5 | 4 (paraphrase recall; fusion gain small in published work) | 2 (model pin, drift, re-embed, growth) | no: dependencies, model distribution |
| E. LanceDB | 3 | 2 (cgo + per-target Rust library) | 4 (versions accumulate until pruned) | 4 | 2 | no |

### The adoption rule

Use **option A** (grep driven by the agent) while every one of these holds. All hold today.

1. A warm whole-tree `git grep` for a literal returns in under 500 ms. It is 11 ms here.
2. No named retrieval failure is on record: no BUG record whose root cause is "the agent did not find the file".
3. The corpus is text a model can grep for, not binary or generated material.

Use **option C** (SQLite FTS5 only) when all of these hold.

1. A trigger is measured, not assumed: condition 1 or 2 above fails, and the failure is recorded as a BUG with its measurement.
2. The operator approves an R4 edit that allows exactly one dependency, `modernc.org/sqlite`, with an ADR. cgo drivers stay out, because they break the plain `go build` cross-compile the release depends on.
3. The index is a derived cache under `.agents/scratch/`, keyed by blob ID, reconciled against `git ls-files -s` and `git status` before every query, with no watcher and no hook, and rebuilt from scratch in under 10 s.
4. The index is local only: WAL mode, `BEGIN IMMEDIATE`, `busy_timeout`, and a fallback to grep when the database is locked past the timeout, corrupt, or on a network filesystem.
5. Grep stays available and remains the tool for exact strings and regex. The index answers ranked and fuzzy queries only.

Add **sqlite-vec** to C (option D) only when all of these also hold.

1. A fixed query set of at least 30 of this project's own questions, including paraphrases with no shared vocabulary, shows BM25 fused with vectors beating BM25 plus agent grep by at least 10 points of recall@10. The query set is committed, and the result is recorded under `testing/`.
2. The embedding model runs in-process with no cgo, no daemon and no network, and it is pinned by name and hash. Otherwise the operator explicitly accepts which of the three the harness takes on.
3. The vector count stays under about 100,000. The sqlite-vec author puts brute-force latency at about 75 ms there (alexgarcia.xyz, 2024-08-01). Above that, wait for ANN indexes to ship stable (DiskANN and IVF are in v0.1.10-alpha only).
4. A meta row pins the model and chunker, and any mismatch forces a full rebuild.

Use **LanceDB** only when vectors reach the millions or the corpus is multimodal, and the harness has accepted cgo. Neither is on this project's horizon.

### Comparison (the research skill's table)

| Approach | Pros | Cons | Complexity | Cost | Maturity |
|---|---|---|---|---|---|
| A. Agent grep (keep) | never stale; zero code; 11 ms here; what Claude Code and Codex CLI chose | no ranking, no paraphrase; more agent turns on vague questions | none | none | shipped in both CLIs |
| B. + symbol map | definitions and references | Go already has `go doc`/`gopls`; 31 Go files here | low | none | mature |
| C. SQLite FTS5 cache | sub-ms ranked search; pure Go; single file | a dependency against R4; invalidation code; saves ~6 ms here | medium | +7.6 MB binary; 3–23 MB per checkout | FTS5 mature; driver active (v1.60.1, 2026-09-29) |
| D. C + sqlite-vec | paraphrase recall | model distribution, drift, re-embed; sqlite-vec pre-1.0 after a 16-month gap between stable releases | high | model file + vectors per checkout | sqlite-vec v0.1.9 (2026-03-31), maintained again since 2026-03 |
| E. LanceDB | ANN at scale, multimodal | cgo over Rust; thin Go SDK; version growth | high | Rust library per target | Go SDK v0.1.2 |

### Recommendation: A, keep agent-driven grep; no index is built

C is the engine to reach for if the rule's triggers ever fire, and nothing else on the list comes close under R3 and R4. At today's scale, 11 ms searches and no recorded retrieval failure, an index would add a dependency, a cache that can lie, and a new verb, to save milliseconds the agent does not notice.

**Full mode stops here.** The research skill's Step 4 writes an integration spec and hands into the `feature` skill for the recommended option. The recommended option builds nothing, so no spec and no plan were written. If the operator chooses C or D instead, that is a Feature with three design questions for the operator first:

1. The R4 edit, as exact text for approval.
2. How the agent calls the index: a new `sh-agent search` verb, or an MCP server. An MCP server needs a vendor settings entry, which the non-goals forbid.
3. Whether the index covers the tracked tree only or also the trail.

## Sources

Web facts come from a research subagent's report, checked 2026-09-29. Dates are the sources' own. Two benchmark papers were located by search and not read in full: Agent Retrieval Bench (huggingface.co/papers/2607.24882, July 2026) and CORE-Bench (arxiv.org/abs/2606.11864).

- sqlite-vec releases, the v0.1.7 "back after a long hiatus" notes and issue #226: github.com/asg017/sqlite-vec
- sqlite-vec brute-force benchmarks: alexgarcia.xyz/blog/2024/sqlite-vec-stable-release (2024-08-01; figures as summarised, charts not read)
- `modernc.org/sqlite/vec`: pkg.go.dev/modernc.org/sqlite/vec (verified by build and run here)
- `ncruces/go-sqlite3` and `ext/vec1`: github.com/ncruces/go-sqlite3, sqlite.org/vec1/doc/trunk/doc/vec1.md
- FTS5 (BM25, trigram, external content, index sizes): sqlite.org/fts5.html
- WAL and network filesystems: sqlite.org/wal.html, sqlite.org/useovernet.html; busy handling: sqlite.org/rescode.html, sqlite.org/lang_transaction.html
- LanceDB Go SDK and versioning: github.com/lancedb/lancedb-go, docs.lancedb.com/tables/versioning
- Claude Code on RAG against agentic search: x.com/bcherny/status/2017824286489383315 (2026-02-01)
- Cursor's server-side index and Merkle sync: cursor.com/blog/secure-codebase-indexing (2026-01-27)
- Continue's deprecated `@Codebase`: docs.continue.dev/reference/deprecated-codebase
- Sourcegraph Cody dropping embeddings: sourcegraph.com/blog/how-cody-understands-your-codebase (2024-02-15)
- Augment on grep and SWE-bench: jxnl.co/writing/2025/09/11/why-grep-beat-embeddings-in-our-swe-bench-agent-lessons-from-augment
- Aider repo map: aider.chat/docs/repomap.html
- Pure-Go MiniLM throughput: github.com/clems4ever/all-minilm-l6-v2-go

## Not proposed

- Any change to `REQUIREMENTS.md`, `sh-agent`, a skill or the context. The recommendation is to build nothing.
- A skill that teaches the agent better grep habits. It was not asked for. The CLIs' own search tools and the trail's indexes already cover it.
