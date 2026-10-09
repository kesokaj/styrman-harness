---
owner: reviewer
date: 2026-10-09
---
# Review — worktree-code-docs   (2026-10-09)

Range reviewed: working tree vs HEAD f7c607a — .agents/docs/specs/sh-agent/search.md, internal/search/reconcile.go, internal/search/search.go, internal/trail/trail.go, internal/trail/trail_test.go; untracked internal/search/hardening_test.go (read whole)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The change itself (search reconcile moved off `git status` onto `ls-files` plus byte hashing, through the hardened `trail.Command`) is sound: build, vet, gofmt clean, `go test` of internal/search and `-race` of internal/trail pass, and the `--debug` stat parser matches real git 2.47.3 output. Worst problem: README.md:187 still tells the operator that search reconciles "against `git ls-files` and `git status`", which the diff makes false (docs lens). Next, a behaviour regression: a tracked `.sh-agent-noindex` marker OUTSIDE the prefix and deleted from the work tree still counts, because the old status pass handled it before the `under()` check and the new loop only visits in-prefix entries.

## Findings
- **[medium / confidence: high]** `README.md:187` — stale: "before every answer it is reconciled against `git ls-files` and `git status`". The diff removes `git status` (spec search.md:42 now says it is never run). Fix: "reconciled against `git ls-files` and the files' own bytes (never `git status`)". Spec search.md:35 and 42-53 were updated; README was not (restatement grep of "git status"/"porcelain"/"ls-files" across README, skills, context, agents, testing guide found this one live copy).
- **[medium / confidence: high]** `internal/search/reconcile.go:206-224` (`gitView` tracked loop over `order`) — dropped: the old status pass ran `delete(markers, p)` for a marker deleted in the work tree BEFORE the `!under(p)` skip, so a tracked root marker removed from disk stopped covering the checkout even when `prefix` (global: `.agents/docs/`) excluded it. `order` now holds only in-prefix entries, so the loop never visits an out-of-prefix marker: it keeps counting and the global reconcile leaves the whole checkout out until the deletion is committed. Direction is fail-safe (less indexed) but contradicts the spec line "one at the checkout's root leaves the whole checkout out" semantics for deleted markers and the code comment kept at reconcile.go ("a marker deleted from the work tree no longer counts"). Fix: Lstat each marker in `markers` (the set is small) after the index loop, independent of `prefix`; add a test with prefix set and a deleted tracked root marker.
- **[low / confidence: high]** `internal/trail/trail.go:113` — comment claims "Every package that runs git in a repository builds the call here (BUG-165)"; false: `internal/memory/memory.go:63,85,128` (`rev-parse`, tracked-files listing, `check-ignore` in the project root) and `internal/harness/repo.go:280,335,342` still use bare `exec.Command("git", …)`. `check-ignore`/`ls-files` read the index and can start the configured fsmonitor hook, the BUG-164/165 class. BUG-165 stays "open" in INDEX.md:117 and its record, which is right only if the memory calls are still to be moved; either route memory through `trail.Command` (a marker test, as `memory_test.go:359` already sets up for filters) or narrow the comment to "search and trail", and update BUG-165 status/prevention in the same change when the fix lands.
- **[low / confidence: medium]** `internal/search/reconcile.go:82-146` (`indexStats`) — no test exercises the parser or the fast path. If git's `--debug` format shifts, `stats` is empty and every tracked file is hashed: all tests stay green while the spec's "about 15 ms" claim (search.md:46) quietly becomes false. `TestGitViewRacilyCleanEditIsHashed` only proves the slow side. Add a test that a clean tracked file's stat entry is parsed (call `indexStats` directly, assert size/mtime for a known file), and one that a path with a newline in its name does not poison its neighbour's entry. I verified by hand that the format parses today; the 15 ms figure itself was not re-measured.
- **[low / confidence: medium]** `internal/search/hardening_test.go:11-86` — `TestGitViewRunsNoCleanFilter` proves the filter control (plain `git status` runs it) and the absence in `gitView`, and checks states; good. Its coverage gap: only `*.md` under a clean filter and no `fsmonitor` hook case (the other half of BUG-165's text), and `os.WriteFile`/`os.Chtimes` errors are ignored in setup. Name the gap in the test comment or add an fsmonitor marker case; BUG-164's trail tests did both.

## Docs lens
- Changed claims: reconcile uses `git status` (removed); modified-file detection; hardened builder. Restatement grep key words: git status, porcelain, ls-files, untracked-files, status --porcelain.
- Stale: `README.md:187` (above). Not stale: search.md (updated, steps renumbered 1-6, no inbound step-number references found), `testing/guide.md` and `sh-agent-matrix.md` (no search-reconcile `git status` claims), skills/, context/, agents/.
- `stale test:` none; `internal/search/security_test.go:316` and `search_test.go:170` use `git status` only to confirm fixtures, legitimately.
- Unverified: the "about 15 ms" cost claim (search.md:46) was not re-measured; BUG-165 record text ("Fix: Open") still says open and proposes `ls-tree`+`hash-object`, while the code chose `ls-files`+own-bytes and a stat shortcut; the fix commit should update it (bugs/INDEX row too).

## Criteria
1. Dead code — pass. `os/exec` import still used in search.go (ExitError); no leftover status parsing; `x` var removed.
2. Flow tracing — one fail: out-of-prefix marker deletion no longer observed (finding 2). Error paths: `indexStats` failures fall back to hashing everything (safe); `ls-files --others` error propagates. Racy check (`st.sec < index mtime`) is stricter than git's, safe.
3. Standards conformance — fail on describing text (README:187, trail.go:113 claim). gofmt/vet clean; errors wrapped.
4. Directory layout — pass. Test beside code; `Command` exported from internal/trail as the shared builder.
5. Simplicity — pass. Stdlib only; the stat shortcut is justified in the comment (perf) though it adds ~70 lines of ad-hoc parsing; no new dependency, no crypto.
6. Decoupling — pass with note: search now imports trail (a git-hardening seam, acceptable); memory and harness remain on bare git (finding 3).

## Not examined
Full `go test -race ./...` (only internal/search without -race and internal/trail with -race were run); the rest of reconcile.go below line ~260 (marker-handling tail, `readAll`, `blobID`) and search.go beyond the changed hunk; behaviour on git older than 2.47; Windows/macOS path behaviour of `rev-parse --git-path index`; timing of the fast path.

## Round 2

Range reviewed: working tree vs HEAD f7c607a after the round-1 fixes (README.md, internal/memory/memory.go and memory_test.go, internal/search/reconcile.go, internal/trail/trail.go; hardening_test.go read whole again). search.md and the BUG-165 record skipped as instructed.

Verdict: APPROVE WITH NITS. Every round-1 finding is closed and the fixes introduced no defect I could find. `go build ./...`, `go vet ./internal/...` and `gofmt -l` are clean; `go test -race` passes for internal/memory, trail, dream and metrics; the new and touched search tests pass.

Per fix:
- README.md:187 (finding 1) — fixed. The row says "each file's own bytes (never `git status` …; BUG-165)". No other living-doc copy of the old claim (re-grepped README, skills, context, agents, testing guide, deploy).
- Out-of-prefix marker (finding 2) — fixed. The pre-loop `for p := range markers` Lstat check runs on the whole checkout's marker set, before the prefix-limited loop; the in-loop delete is correctly gone (nothing relied on it: an untracked marker is added later in the `--others` loop, and the Lstat-error semantics match the old `y == 'D' || err != nil`). `TestGitViewDeletedMarkerOutsidePrefix` asserts one entry under the prefix with a deleted tracked root marker, so it fails if the pre-check is removed. The earlier `TestNoIndexMarkerDeletedCountsNoMore` still passes.
- trail.go comment and memory (finding 3) — fixed. The three memory calls (`Root`, `trusted` via `storeTracked`, `Add`'s check-ignore) all use `trail.Command`; no bare git remains in internal/memory, search, dream or metrics. The only remaining `exec.Command("git"…)` are internal/harness/repo.go (the operator's own harness checkout), as the new comment states. The re-declared helper left nothing orphaned: `storeTrackedArgs` and the lowercase `command` have no remaining references (grep over Go sources and tests; the `command` hits in internal/dream are an unrelated function). `TestStoreTrackedArgs` pins the new full argument list. `TestAddRunsNoProjectHook` ran (not skipped here): its control proves plain `check-ignore` fires the fsmonitor hook, then `Add` does not. Its `t.Skip` on a failed control means it proves nothing on a git that does not run the hook; that is the right fallback, not a defect.
- `indexStats` without `-z` (security finding) — correct. Git quotes names with newline, tab, quote or backslash, so a header line cannot carry another entry's stat lines, and a quoted name is skipped (hashed). Stat lines start with a space, headers with a mode digit, so a line classifies unambiguously. `TestIndexStatsIgnoresForgedNames` doubles as the direct parser test: if `indexStats` returned nothing, `a.md`'s size would read 0 and the test would fail. `TestGitViewRunsNoFsmonitor` ran (control held) and covers `gitView` plus `guardProject`.

Findings
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/search.md:44` (not yet updated, skipped by instruction) still says `git ls-files -s --debug -z`; the code no longer passes `-z`. Make sure the spec update lands in the same commit, with the reason (quoted names are skipped).
- **[low / confidence: medium]** `internal/search/reconcile.go` `indexStats` — consequence of dropping `-z`, not a bug: tracked files whose names git quotes (non-ASCII under the default `core.quotePath`, or any name with a tab, quote or backslash) never take the stat shortcut and are hashed on every reconcile. Correctness is unaffected; the spec's "about 15 ms" is for ASCII names. Name this in the spec sentence rather than leaving it implied.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/metrics.md:259` — its "git in project roots" bullet lists dream and metrics callers; it now has two more (search, memory). Align it in the same change if it enumerates callers (not read to the end of the line; unverified).

Criteria (round 2): dead code pass (no orphaned `command`/`storeTrackedArgs`); flow tracing pass (marker pre-check order, `indexStats` fallback to full hashing on any error); standards pass (describing text fixed, BUG-165 record/INDEX row still to be updated to open->fixed in the fix commit); layout pass; simplicity pass; decoupling pass (memory now imports trail, a one-way edge: trail imports neither memory nor search).

Not examined (round 2): the full `go test -race ./...` and `go test ./cmd/...`; internal/search with `-race`; the search.md and BUG-165 edits (not yet in the tree); behaviour on a git other than 2.47.3; the tail of reconcile.go below the untracked loop (unchanged since round 1).

## Outcome

Written by the main agent, 2026-10-09.

Round 1:
1. README still names `git status`: **fixed** at `README.md` (the search row).
2. A deleted marker outside the prefix is still counted: **fixed** in `internal/search/reconcile.go` (every marker is checked before the loop), under `TestGitViewDeletedMarkerOutsidePrefix` (red before the fix).
3. The `trail.Command` comment overclaims: **fixed**. The comment is narrowed to project repositories, and memory's three calls are moved onto the builder (`internal/memory/memory.go`), under `TestAddRunsNoProjectHook`.
4. `indexStats` has no direct test: **fixed** by `TestIndexStatsIgnoresForgedNames`.
5. No fsmonitor case in search: **fixed** by `TestGitViewRunsNoFsmonitor`.

Round 2 nits:
6. The spec still says `-z`: **fixed** in `specs/sh-agent/search.md`.
7. Quoted names are hashed every time: **fixed**, stated in the same spec step.
8. The `metrics.md` caller list: **fixed**; the allowed-command list now includes `ls-files`, `check-ignore` and `config --get`.
