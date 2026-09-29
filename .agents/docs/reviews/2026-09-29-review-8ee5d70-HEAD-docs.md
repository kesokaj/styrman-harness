---
owner: reviewer
date: 2026-09-29
---
# Review — 8ee5d70-HEAD-docs   (2026-09-29)

Range reviewed: 8ee5d70..HEAD (docs lens; `git diff 8ee5d70...HEAD`)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The bulk of this range's documentation is precise and internally consistent — REQUIREMENTS.md R4, README.md's command table, the new `specs/sh-agent/search.md`, and the per-consumer lines added to seven agents and eight skills all match `internal/search`'s actual message strings, flags, and ranking behaviour exactly (spot-checked against `search.go`, `store.go`, `parse.go`). The three "three sentences cut or moved" items named in the brief all land correctly and have no stale copy left behind (`git grep` over the living-docs set found none). Two claims about the new dependency are wrong, though: one names the wrong pinned version, and one asserts a `govulncheck` result that no longer holds now that the dependency is in `go.mod`. Neither is a design defect, both are two-word/one-clause fixes.

## Findings

- **[medium / confidence: high]** `AGENTS.md:29,32` — "`govulncheck ./... # clean on go1.25.8+`" and "On a Go older than go1.24.6 or go1.25.8, `govulncheck` reports two standard-library advisories … the fix is upgrading the machine's Go" is no longer true once this range's dependency (`modernc.org/sqlite`, pulling in `golang.org/x/sys` v0.37.0) is in `go.mod`. Verified by running `govulncheck ./...` on `go1.27.1` (above the 1.25.8 floor the doc names): it reports "1 vulnerability in modules you require" (GO-2026-5024 in `x/sys`, per plan D25/`.agents/docs/plans/2026-09-29-search-index.md` D25), not the "0 vulnerabilities" the doc's "clean" claim promises. This advisory is orthogonal to the local Go version — D25 says the fix needs `x/sys` v0.44.0, which needs `go 1.25.0` in the module's `go` directive, ruled out by D21 — so it will keep showing on every machine, including go1.25.8+, until that constraint lifts. `AGENTS.md` was edited in this range (the `go mod download` line was added right above) but the govulncheck lines were not updated to say a third, dependency-pinned advisory is now expected and non-blocking. Fix: add a clause noting the `x/sys` GO-2026-5024 advisory as expected/non-reachable (D25), same treatment as the two stdlib ones.

- **[low / confidence: high]** `.agents/docs/architecture/DECISIONS.md` ADR-045 — "The pure-Go `modernc.org/sqlite` v1.60.1 was measured here" states the wrong version. `go.mod:16` and `go.sum:22-23` pin `v1.46.1`, matching the plan's own D21 ("`modernc.org/sqlite` is pinned at v1.46.1 … `go get` of v1.60.1 raised this module's `go` directive to 1.26.0 … a Go 1.24 machine … could no longer build"). v1.60.1 was the version *rejected* by D21, not the one shipped; ADR-045's sentence reads as if v1.60.1 is what got pinned. Fix: correct the version number to v1.46.1 in ADR-045's Context paragraph (or add the D21 rejection detail there).

- **[low / confidence: medium]** `.agents/docs/architecture/DECISIONS.md` ADR-045 — "It adds 7.6 MB to the binary and nine indirect modules" undercounts: `go.mod`'s `require` block lists ten modules besides `modernc.org/sqlite` itself (`dustin/go-humanize`, `google/uuid`, `mattn/go-isatty`, `ncruces/go-strftime`, `remyoudompheng/bigfft`, `golang.org/x/exp`, `golang.org/x/sys`, `modernc.org/libc`, `modernc.org/mathutil`, `modernc.org/memory`), not nine. Minor and doesn't change the decision, but a reader recounting the tree gets ten, not nine.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Verified matches (no finding)
- REQUIREMENTS.md R4 rewrite ("stdlib only" → "one dependency") matches `go.mod`'s actual dependency set and `internal/search/store.go`'s import of `modernc.org/sqlite`; the installer (`internal/harness`) does not import it (checked via `git grep` — no hit).
- `sh-agent search`/`index` rows added to README's command table match `cmd/sh-agent/main.go`'s flag definitions and `internal/search/search.go`'s literal error/message strings (`errNotCheckout`, `errNotIgnored`, `--limit` range, `--kind` list, busy/unreadable/hint messages) word for word.
- Word-search ranking claim ("chunks holding every word rank first, filled with any-word hits") matches `search.go:566-577` (`sc.find` two-pass AND-then-OR, header line format) and plan D27.
- `--refs` bare-filename hint matches `search.go:693` and plan D28.
- BM25 "heading weighted twice the body" matches `search.go:492` (`bm25(chunks_fts, 2.0, 1.0)`); `perFileCap = 2` matches the "at most 2 per file" claim.
- The three brief-named context cuts: "Read the index before touching an area it names" was folded into the new "Search first" entry's `--refs` clause and left no stale copy in any living doc (only in `plans/`/`reviews/` history, which is out of scope); "`/skills` and `/agents` show what a vendor loaded…" moved verbatim to `skills/init-project/SKILL.md:110` and was removed from `context/AGENTS.md`; the `find-skills`/`npx skills add` procedure in `context/AGENTS.md` was shortened as described, with the full `npx skills add` mechanics still correctly documented in the project `AGENTS.md`, `README.md`, and `skills/VENDORED.md`.
- ADR-002 and ADR-029 headings correctly append "superseded in part by ADR-045" pointers; `TestRenderedContextFitsAntigravityBlock`'s 22,000-byte pin is not violated (`context/AGENTS.md` is 18,484 bytes raw).
- `testing/guide.md` § 3d, `sh-agent-matrix.md`, and `skill-matrix.md` rows for the new verbs are present and consistent with the spec and code.
- `.agents/docs/deploy/local.md`'s "needs Go, and the Go module proxy or a filled module cache for the one dependency" addition matches ADR-045's Consequences paragraph on offline-build reopening.
- `.agents/docs/deploy/remote.md` and `.agents/docs/deploy/gcp.md` have no diff in this range and needed none (no release, remote, or GCP-facing change here).
- `ROADMAP.md` "Now" naming the search-index plan is correct: the plan (`.agents/docs/plans/2026-09-29-search-index.md`) still has 9 unticked boxes (Phase 3's own checkpoint at line 234 is unticked — awaiting exactly this kind of `/review --docs` pass plus `reviewer`/`security`/verifier — and all of Phase 5), so the plan is genuinely in progress, not stale.

## Not examined
- `internal/search/reconcile.go` was not read line-by-line against the spec's "Freshness: reconcile before every answer" five-step description; only `search.go`, `store.go`, and `parse.go`'s ADR-heading regex were checked directly. The reconcile steps (git rev-parse, ls-files -s -z, status --porcelain, hash-object, 14-day blob GC) are taken on the spec's word against the commit messages and plan D-rows, not re-derived from the reconcile code itself.
- `cmd/sh-agent/search_test.go` and the full `internal/search/*_test.go` suite were not run or read in full; test-vs-doc consistency for edge cases (network filesystem refusal, corruption/rebuild path) rests on the failure-mode table's string matches already checked, not on executing those paths.
- `agents/security.md`'s and `agents/reviewer.md`'s new pre-flight `sh-agent search --refs` lines were read in diff form only, not exercised against a live dispatch.
- Binary-size delta ("7.6 MB") and the recall/measured-numbers table in `search.md` § Measured were taken as reported, not independently re-measured (would require a fresh `index --rebuild` and a before/after binary build, outside this review's budget).

## Outcome

- 1. `AGENTS.md` govulncheck line promised a clean result — fixed: `AGENTS.md` § Setup commands now names GO-2026-5024 (`golang.org/x/sys` via the driver, not reachable) and ADR-045a records why it stays.
- 2. ADR-045 named v1.60.1 — fixed: `DECISIONS.md` ADR-045a, an addendum (the file is append-only), states v1.46.1 and why.
- 3. ADR-045 counted nine indirect modules — fixed: ADR-045a says ten, counted after `go mod tidy` (which also moved `modernc.org/sqlite` from `// indirect` to a direct requirement; the `go` directive stays 1.24.4).
