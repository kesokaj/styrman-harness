---
owner: reviewer
date: 2026-09-29
---
# Review — search-index-code   (2026-09-29)

Range reviewed: c04f2d7..HEAD

## Status: COMPLETE

## Verdict: APPROVE

The diff adds `sh-agent search`/`index` (a new `internal/search` package backed by `modernc.org/sqlite`), wires it into `main.go`, teaches `internal/harness/report.go` to skip the three search-index cache files on `report`, and threads "search first" guidance through seven agents and eight skills plus `context/AGENTS.md`. `go vet ./...` and `gofmt -l cmd internal` are clean, `deadcode` finds nothing new (its one hit, `SwapReleaseSigner`, is pre-existing and outside this range), and the 31 new tests in `internal/search` exercise real behaviour (branch switches, dirty/untracked files, renames, binary/symlink/submodule skips, GC, FTS-trigger sync, corrupt/stale-schema rebuild, busy-lock refusal, concurrency, global registration/collision/vanished-project drop, `--refs`/`--open`/`--kind` filtering, argument validation) rather than re-asserting constants. The one new dependency (`modernc.org/sqlite`) was approved and justified in the commit that is this diff's *base* (c04f2d7, R4/ADR-045), not introduced by this range, and `internal/harness` still never imports it (confirmed by grep and by the package's own doc comment). I found no correctness, security or requirement-violation defect in the range.

## Findings
None. No high/medium/low findings — the diff is clean against the six criteria below.

## Criteria
1. **Dead code** — pass. `deadcode ./cmd/sh-agent/` reports only a pre-existing, out-of-range hit (`internal/harness/release.go:422 SwapReleaseSigner`, untouched by this diff). No unused functions/types were introduced in `internal/search`, `main.go`, or `report.go`.
2. **Flow tracing** — pass. `search`/`index` wiring in `main.go` → `search.Search`/`search.Index` → `prepare()` → per-project and global `store` is consistent end to end; resource cleanup (`st.close()`) is present on every error path in `prepare`, `reconcileGlobal`, and `dropProject`; transactions (`BEGIN IMMEDIATE`/`COMMIT`/`ROLLBACK`) are paired correctly, including the busy-lock (`errBusy`) and stale-schema/corruption rebuild paths. `report.go`'s new `c.Skipped` entries are consumed by the existing dry-run print loop at report.go:671, same as every other skip reason. Tests can fail (behavioural assertions on printed output/DB state, not tautologies).
3. **Standards conformance** — pass. Go errors are wrapped with `%w` and context throughout; `gofmt`/`go vet` clean. `context.Context` is not threaded through `internal/search`, but this matches the rest of the codebase (`internal/harness` also does not use `context.Context` for its own single-shot CLI operations) — not a deviation introduced by this diff. The "restate the changed rule everywhere" check: the new "search first" rule is now present in `context/AGENTS.md`, 7 agents and 8 skills as the commit claims (spot-checked bugs.md, changelog.md, documentation.md, dreamer.md, explorer.md, reviewer.md, security.md, and bug/dream/feature/handoff/init-project/research/retro/review-standards SKILL.md); `context/AGENTS.md`'s trim of the `npx skills add`/`find-skills` detail in "Global vs project" is still fully covered by README.md (confirmed present there), so no orphaned rule.
4. **Directory layout** — pass. `internal/search/*.go` with tests beside the code (`*_test.go`), `cmd/sh-agent/search_test.go` beside `main.go`; nothing new at the repo root.
5. **Simplicity (DIY-except-security)** — pass. The one new dependency is pre-approved (R4/ADR-045, predates this range) and confined to `internal/search` as its doc comment states; SQL is stdlib `database/sql` + parameterized queries throughout (no string-built literals for values, only structural clause assembly in `refs()`, verified parameterized and covered by `TestRefs`/`TestGlobalRegistersProject`).
6. **Decoupling** — pass. `internal/search` is the only importer of the SQLite driver; `internal/harness` references the three cache file names as plain strings rather than importing `search`, avoiding a dependency edge from the installer into the new package. Clean seam.

## Not examined
- `internal/search/search_test.go` and `internal/search/recall_test.go` were read in large part but not exhaustively line-by-line past the sections covering: helper setup, reconciliation/GC/busy/corruption tests, recall-arm construction and `TestRecall`/`TestRecallQuestionsAreWellFormed`. No issues found in what was read; the remaining assertions in `TestSearchRanksAndCaps`, `TestRefs`, `TestOpen`, `TestConcurrentReconcile` bodies were not individually re-verified line-by-line (only their names/shape and general structure).
- `.agents/docs/specs/sh-agent/search.md`, the plan file `2026-09-29-search-index.md`, the handoffs, and the testing/recall-baseline docs under `.agents/docs/testing/` and `.agents/docs/bugs/BUG-106*`, `BUG-107*` were read only partially (bug records read; spec/testing docs not read) — these are process/trail records, not source, and were out of the dispatch's Go-source focus, but are flagged here as unexamined for completeness.
- Did not run `go test ./...` (only `go vet` and `gofmt -l`), per the reviewer's role of reading rather than re-running the full suite; the verifier's own log is the authority for pass/fail, not restated here.
- `internal/search/testdata/recall.tsv` contents were not read in detail.
