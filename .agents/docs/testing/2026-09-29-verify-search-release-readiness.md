---
owner: verifier
date: 2026-09-29
---
# Verification — search-release-readiness   (2026-09-29)

Range reviewed: aae16af..1634787

## Status: PASS

Scope: complete — every command AGENTS.md § Setup commands and the plan's section 3d named ran, plus the four named tests, plus coverage mapping of changed non-test .go files and a suppression scan of the diff. `.agents/docs/testing/guide.md` exists (not modified; no step changed).

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go version` → go1.27.1 linux/amd64; `grep '^go ' go.mod` → `go 1.24.4` (matches expectation; no `toolchain` directive, per ADR-029).
- `go test -race -count=1 -v ./...` → exit 0. Full output `.agents/scratch/verify-search-release-readiness/go-test-1.txt`. All 4 packages `ok`: root, `cmd/sh-agent`, `internal/harness`, `internal/search`.
  counts: pass 362, fail 0, skip 2 at 1634787 for `go test -race -count=1 -v ./...`
  vs previous: `.agents/docs/testing/2026-09-29-verify-search-index-final-2.md` — pass 359, fail 0, skip 2 at `4c9a327` (confirmed an ancestor commit inside this range: `git merge-base --is-ancestor 4c9a327 1634787` → true, and `4c9a327` appears in `git log aae16af..1634787`). pass 359→362 (+3), skip 2→2 (unchanged). No fall in pass, no rise in skip.
  explained: `git diff 4c9a327..1634787 -- '*_test.go' | grep '^+func Test\|^-func Test'` shows exactly 3 new top-level test functions added and none removed: `TestRefsPrefilterIsEquivalent`, `TestRefuseSymlinkedDatabase`, `TestDatabaseIsPrivate` — matches the +3 exactly.
  The 2 skips, both carrying a reason: `TestRecall` (`SH_AGENT_RECALL=1` gate) and `TestHelperIndex` (subprocess-helper pattern) — confirmed via `grep -E '^\s*--- SKIP' go-test-1.txt`.
  Note (process, not a finding): an initial attempt at this command used `run_in_background` with an internal shell `&`, which detached from the tool's own wait and returned a truncated 8-line log (only the root package's 3 fast tests). Discarded; re-run correctly via the tool's own backgrounding (no internal `&`), which waited for full completion and reported the counts above.
- `govulncheck ./...` → "No vulnerabilities found" in code/imports; 1 vulnerability in required-but-uncalled modules: `GO-2026-5024` (`golang.org/x/sys`, `NewNTUnicodeString` integer overflow, windows-only, fixed in v0.44.0, this module at v0.37.0) — matches expectation (only GO-2026-5024 expected on this Go version).
- Cross-builds (release targets), `CGO_ENABLED=0 GOTOOLCHAIN=local go build -trimpath ./cmd/sh-agent`, one call for all three:
  `linux/amd64` exit 0, `darwin/arm64` exit 0, `linux/arm64` exit 0. Binaries under `.agents/scratch/verify-search-release-readiness/builds/sh-agent-{linux-amd64,darwin-arm64,linux-arm64}`.
  (Superseded an earlier, wrong first attempt that built linux/darwin/windows all at amd64; that attempt's output is not used in this report.)
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -v -count=1` → PASS. `size_test.go:48: Antigravity rules block is 21996 bytes (context 21996, rules 0), budget 22000` — matches the expected 21,996 exactly. Run twice (once before, once after the coordinator's correction), identical result both times.
- Install/status into `.agents/scratch/verify-rr`:
  `go run ./cmd/sh-agent install --home .agents/scratch/verify-rr` → `266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed`.
  `.agents/scratch/verify-rr/.local/bin/sh-agent status --home .agents/scratch/verify-rr` →
  `repo       ... (main @ 1634787, uncommitted changes)`; `binary     163478749717+dirty` (no "warning" text present — the "+dirty" reflects two untracked files sitting in the working tree during this run, `.agents/docs/reviews/2026-09-29-review-abbaddc-1634787-docs.md` and this verification log itself, not any change to a tracked source file; the tracked tree is unmodified); `drift      none — 265 files match the repo`.
  `grep -c 'sh-agent search --open' .claude/AGENTS.md .gemini/config/rules/AGENTS.md` → 1 and 1 (present in both installed contexts).
  `tail -c 8` of both installed `AGENTS.md` files → `GANNET.` in both.
- Section 3d, binary built into `.agents/scratch/verify-search-release-readiness/sh-agent`, `--home .agents/scratch/verify-rr-home`:
  `index --rebuild` × 3 on this checkout: files 1084, chunks 10621, records 944, refs 41044 (identical counts across all 3 runs); ms 1718 / 1720 / 1741 (≈1.7–1.74 s); `.agents/scratch/search.db` byte size 21,786,624 / 21,790,720 / 21,782,528 (≈20.8 MiB / 21.79 MB reported by the tool's own `bytes` line 21,790,800–21,798,992) — matches the plan's claim of "about 1.7 s and 21.7 MB" closely.
  Warm `search symlink guard --home …` × 5, timed end to end (includes process exec overhead): 36, 33, 28, 31, 31 ms — all under 100 ms.
  `stat -c %a` of `.agents/scratch/search.db` → `600`. `stat -c %a` of `.agents/scratch/verify-rr-home/.config/sh-agent/search.db` (global index) → `600` — present and private (this checkout is harness-scaffolded and registers itself on search/index).
  `SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v -count=1` → PASS (4.04s). Both arms' means: grep mean all 0.379 over 41 (words 0.269/26, drift 0.353/10, refs 1.000/3, open 1.000/2); index mean all 0.649 over 41 (words 0.731/26, drift 0.359/10, refs 1.000/3, open 0.500/2). The indexed arm scores higher overall (0.649 vs 0.379), driven by the `words` category (0.731 vs 0.269); `refs` ties at 1.000 for both; `open` is grep's strong category (1.000 vs 0.500) on only 2 queries.
- Named tests, individually, verbose, confirmed RAN and PASSED (`.agents/scratch/verify-search-release-readiness/named-tests.txt`):
  `TestRefuseSymlinkedDatabase` PASS (0.05s); `TestDatabaseIsPrivate` PASS (0.05s); `TestRefsPrefilterIsEquivalent` PASS (0.00s); `TestDropReasonIgnoresLocale` PASS (0.00s).

## Root cause  (FAIL only)

Not applicable — Status is PASS.

## Coverage of changed code

Changed non-test `.go` files in `aae16af..1634787` (git diff --name-only, excluding `_test.go`): `internal/search/parse.go`, `internal/search/reconcile.go`, `internal/search/search.go`, `internal/search/store.go`.

- `internal/search/parse.go` — `refExtensions`, `hasRefExtension`, the rewritten `pathToken` regexp built from that one list, and the new prefilter branches inside `parseRefs` (skip a line with neither `BUG-`/`ADR-` nor a ref extension) — covered by `TestRefsPrefilterIsEquivalent` (`internal/search/search_test.go:585`), which runs the prefiltered `parseRefs` over every tracked file and asserts the same 41,044 references as an unfiltered reference implementation, and by mutation testing recorded in the commit message (`b2e0d5a`: mutations O1–O3 each fail it). Exercised indirectly by every `index --rebuild` run above (1084 files, 41044 refs, stable across 3 runs).
- `internal/search/reconcile.go` — the `git()` function's new `cmd.Env = append(os.Environ(), "LC_ALL=C", "LANGUAGE=")` line is covered by `TestDropReasonIgnoresLocale` (`internal/search/security_test.go:214`), confirmed RAN and PASSED above (installs a fake `git` that only succeeds in the C locale). The new `inserter` type (`newInserter`, `close`, `content`, replacing `insertContent`, prepared statements instead of ad-hoc `Exec`) is exercised by every existing reconcile test that inserts content — `TestReconcileBranchSwitch`, `TestReconcileDirtyAndUntracked`, `TestReconcileDeletedFileLeaves`, `TestReconcileStagedRename`, `TestReconcileGC` (`internal/search/search_test.go`) — all of which passed in the full-suite run (part of the 362 pass) and by the `index --rebuild` runs above (944 records, 10621 chunks inserted identically across 3 runs, proving the prepared-statement path produces the same output as before).
- `internal/search/search.go` — the only change is a doc-comment correction on `printable` (describing existing, unchanged behavior more precisely); no logic changed, so no new test is required; the function itself remains covered by the existing `printable`-adjacent tests in `security_test.go`.
- `internal/search/store.go` — `prepareFiles` (refuses a non-regular-file database or WAL/SHM companion, creates a new database 0600, tightens loose permissions) is covered by `TestRefuseSymlinkedDatabase` and `TestDatabaseIsPrivate` (`internal/search/security_test.go:230`, `:267`), both confirmed RAN and PASSED above, and by the perms checks in section 3d (`stat -c %a` → 600 on both the local and global index files). The dropped `CREATE INDEX refs_oid ON refs(oid)` (unused, served only garbage collection) is exercised by `TestReconcileGC` continuing to pass and by the unchanged reference counts across three rebuilds.

No changed function in these four files is without a covering test.

Suppression scan of the diff (source/test files only, excluding `.agents/docs/` prose that merely quotes past findings): `git diff aae16af..1634787 -- internal/search/parse.go internal/search/reconcile.go internal/search/search.go internal/search/search_test.go internal/search/security_test.go internal/search/store.go context/AGENTS.md skills/handoff/SKILL.md ROADMAP.md | grep -nE '(nolint|lint:ignore|ts-ignore|ts-expect-error|eslint-disable|t\.Skip|SkipNow|go:build ignore|\.skip\(|\.only\()'` → no matches. No suppression, skip, or `.only(` added by this range's source or test changes.

## Not examined

- The A/B adoption experiment itself (sample-app and large-trail lane sessions, cost/quality comparisons, `metrics-clean.py`/`score-trail.py` outputs) — those are the plan's own AI-verified phases (Phase 1, Phase 2), already marked done in the plan with evidence paths; this run did not re-run or re-score those sessions, only the code/test/build verification named in the dispatch.
- The Antigravity gate runs (Phase 3: `dreamer`, `security`, and the changed-skill runs) — not re-run here; the plan cites `testing/2026-09-29-search-agy-gate.md` as evidence.
- The real-home smoke checks on both CLIs (Phase 4) — not re-run here (would touch the operator's real `$HOME`, out of scope for this dispatch, which used `.agents/scratch/verify-rr` and `.agents/scratch/verify-rr-home` throughout).
- `TestHelperIndex` (skipped, subprocess-helper pattern) was not run directly with its enabling condition; only its presence and reason were confirmed in the full-suite output.
- Coverage of `context/AGENTS.md`, `skills/handoff/SKILL.md`, `ROADMAP.md` changes — these are prose/context files, not `.go` code; verified only via the byte-budget test, the installed-context grep/tail checks, and the suppression scan; not otherwise "tested" since they carry no test framework of their own.
- The plan's own open Phase 2 and Phase 4 checkpoint/Definition-of-Done boxes (unticked in the plan file as of this range) are outside this verification's scope, which is the code range `aae16af..1634787`, not plan closure.

## Environment

OS: Linux 6.12.107+deb13-amd64. Go 1.27.1 (module declares `go 1.24.4`, no `toolchain` directive, per ADR-029). Working tree: two untracked files present throughout this run — `.agents/docs/reviews/2026-09-29-review-abbaddc-1634787-docs.md` (not created by this verifier; appears to be a concurrent `/review --docs` artifact) and this log itself — no tracked file was modified during the run; `git diff` and `git diff --staged` against tracked files were empty throughout.
