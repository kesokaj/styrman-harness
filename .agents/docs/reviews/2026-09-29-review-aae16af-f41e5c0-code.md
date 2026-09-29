---
owner: reviewer
date: 2026-09-29
---
# Review — aae16af..f41e5c0   (2026-09-29)

Range reviewed: aae16af..f41e5c0

## Status: COMPLETE

## Verdict: APPROVE

Clean, well-tested diff. It lands three things: (1) search-index build/size optimizations (prepared-statement `inserter` reused per reconcile, a substring prefilter in `parseRefs` before the regexes run, and dropping the unused `refs(oid)` index), each backed by a real equivalence test or measured numbers in the spec doc; (2) the two low-severity findings from the prior review (`2026-09-29-review-search-index-code-r2.md`) — the locale-dependent `dropReason` match and the stale `printable` doc-comment — both fixed with a targeted test (`TestDropReasonIgnoresLocale`) or a corrected comment; (3) a context/skill/ROADMAP update adding `sh-agent search --open` to the resume reading order, applied consistently in both `context/AGENTS.md` and `skills/handoff/SKILL.md`. `gofmt -l cmd internal`, `go vet ./...`, and `go test ./internal/search/... -race -count=1` are all clean; nothing broken.

## Findings

None. Nothing meets the bar for correctness/security/requirement-violation.

Two things considered and not written up as findings:
- `internal/search/store.go` drops `CREATE INDEX refs_oid ON refs(oid)`. `refs.oid` is still read in `store.gc()`'s `DELETE FROM refs WHERE oid IN (...)` and in the `refs`/`files` JOIN in `search.go:758` and `search.go:717`, but in both SELECT paths the join is driven from the indexed `refs_target`/`records_id` side into `files` via `files_oid`, so the dropped index isn't load-bearing there; only `gc()`'s bulk DELETE loses an index, and that runs over a small, aged subset. The commit message states this was measured (24.1 MB → 21.7 MB). Not a finding — a deliberate, apparently-verified trade-off.
- `context/AGENTS.md`'s "Search first" bullet lost its explanation of `--open` entirely (it's now only named, unexplained, in the new Resume-order clause), while `README.md` and `.agents/docs/specs/sh-agent/search.md` still fully document the flag. Given `context/AGENTS.md` is byte-capped at 22,000 (`TestRenderedContextFitsAntigravityBlock`, currently 21,996) and the commit message names this exact trade (`paid for by the --open clause of the Search first entry`), this reads as a deliberate, budgeted trim rather than a stale/contradictory restatement — not a finding.

## Criteria

1. **Dead code** — pass. `store.insertContent` was fully replaced by `inserter.content`; `grep -rn insertContent internal/` finds nothing left. No other orphaned code in the diff.
2. **Flow tracing** — pass. `applyLocked`'s new `inserter` is created lazily (`in == nil` guard), `defer in.close()` registers exactly once, and `newInserter` closes any partially-prepared statements on its own error path (no leak). `parseRefs`'s new prefilter is proven equivalent to the un-prefiltered regex pass by `TestRefsPrefilterIsEquivalent`, which independently recomputes `want` from the same regexes with no shortcuts and asserts `reflect.DeepEqual` against `parseRefs`'s actual output — a real test that fails if the prefilter drops or duplicates a match (verified: fixture has 10+ len(refExtensions) references, one line per accepted extension). `TestDropReasonIgnoresLocale` plants a fake `git` that answers in German unless `LC_ALL=C` is set, and fails without the `cmd.Env` fix (Go's documented last-value-wins dedup for `Cmd.Env` makes the `append(os.Environ(), "LC_ALL=C", ...)` override correct even when the parent already has `LC_ALL` set).
3. **Standards conformance** — pass. `go vet ./...` and `gofmt -l cmd internal` are clean. The `--open` resume-order addition is applied identically in `context/AGENTS.md` and `skills/handoff/SKILL.md` (grepped both); the spec doc's measured numbers (1.7s build, 21.7 MB, 40,850 refs) match the commit message's claims.
4. **Directory layout** — pass. No new files under `internal/search`; tests stay beside the code they cover (`search_test.go`, `security_test.go`).
5. **Simplicity (DIY-except-security)** — pass. No new dependency. The `LC_ALL=C`/`LANGUAGE=` env fix and the prepared-statement reuse are both plain stdlib (`os/exec`, `database/sql`), not new abstractions.
6. **Decoupling** — pass. `inserter` is a small, package-private helper scoped to `reconcile.go`; no new cross-package reach-ins.

## Not examined

- The trail-only files in the diff (`.agents/docs/handoffs/*`, `.agents/docs/plans/2026-09-29-search-index.md` and `-search-release-readiness.md`, three `coldread-*` files, the `review-search-index-code-r2.md` review artifact, `.agents/docs/testing/2026-09-29-search-efficiency.md`, two `verify-search-index-final*` files) — outside the dispatch's named file list, and several are other agents' own trail records, not source under review.
- Behavioural/CLI smoke testing beyond the unit test suite (no manual `sh-agent search`/`sh-agent index` run against a real home) — not needed to judge this diff's correctness given the automated coverage, but not performed either.
- `internal/search/search.go` outside the two diffed regions (the `refs()` query path and the `printable` comment) — unchanged in this range, not re-reviewed.
