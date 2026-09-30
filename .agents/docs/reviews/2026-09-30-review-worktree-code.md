---
owner: reviewer
date: 2026-09-30
---
# Review — worktree-code (2026-09-30, round 3)

Range reviewed: working tree — changed files: AGENTS.md cmd/sh-agent/main.go cmd/sh-agent/main_test.go internal/harness/adopt.go internal/harness/apply.go internal/harness/harness_test.go internal/harness/release.go internal/harness/release_test.go internal/harness/report.go internal/harness/report_test.go internal/harness/state.go skills/VENDORED.md; untracked, not ignored: internal/harness/sweep.go internal/harness/sweep_test.go

## Status: COMPLETE

## Verdict: APPROVE

All four fixes named for this round were verified against their call sites, not only their own lines, and none introduces a new defect.

1. **`apply.go` `WriteAtomic`/`abandon` (round-2 medium finding)** — the ordering is now Write → Chmod → Sync → **Close** → **Rename**, matching the pre-refactor invariant: `error != nil` ⇒ `dest` unchanged, since a close failure now happens before any rename is attempted. `abandon(tmp, err)` is called from all five failure branches; on the post-close branches (Close itself, Rename) it calls `tmp.Close()` a second time (a no-op error, correctly ignored per its own comment) then `os.Remove(tmp.Name())`, which is safe since `File.Name()` doesn't require an open fd and, on both branches, the temp file is provably still at that path (Rename failure leaves the source in place). All five call sites (`Install`, `adopt.copyTree`, `state.SaveState`, `release.UpdateFromRelease`, `cmd/sh-agent`'s settings writer) now share the one implementation via the exported `WriteAtomic`, and `gofmt`/`go vet`/`go test ./internal/harness/... ./cmd/...` are clean. No stray reference to the removed `writeAtomic` (lowercase) remains.
2. **`report.go` needle additions** — `sk_test_`, `rk_test_`, `sk-svcacct-`, `sk-admin-`, `hooks.slack.com/workflows/` are in `prefixNeedles` (report.go:48), and the url-token regex now includes `%` in its character class (report.go:80, `://[A-Za-z0-9_.%-]{20,}@`). `TestReportRedactShapes` (report_test.go:96-97) exercises every one of these five new shapes plus a percent-escaped userinfo token, and passes.
3. **`release_test.go:654` sweep test** — `TestUpdateFromReleaseSweepStaysInHome` now creates a `release-young` directory alongside `release-old`/`release-keep`, leaves it at its natural (fresh) mtime, and asserts it survives the sweep alongside the sibling `release-keep`; `release-old` is still asserted removed. This closes the round-2 gap (an hour-threshold check with nothing young to test against).
4. **`skills/VENDORED.md:30`** — the `gke-basics` row now states the split to `gke-*` sub-skills happened at commit `5d9d054` on 2026-06-23, consistent with the surrounding provenance narrative (security finding 4's traced-by-diff justification).

No new findings. The `abandon` doc comment's claim that renaming a still-open temp file over `dest` risks a later `ETXTBSY` on `execve(dest)` is consistent with real Linux semantics (an inode open for writing anywhere blocks exec of it), not a stray justification — checked because it reads as a second, independent reason bolted onto the fix and could have been a misdescription.

## Criteria

1. **Dead code** — pass. `writeAtomic` (lowercase), `abandonTemp`, `globOrNil`, `sweepStaleZipDirs` all have zero remaining references (grepped clean); their replacements (`WriteAtomic`, `abandon`, `SweepStaleDirs`) are used at every former call site.
2. **Flow tracing** — pass. `WriteAtomic`'s error invariant is restored and consistent at all five call sites; `abandon`'s double-close-then-remove path is safe on both branches it's reachable from post-close. `TestAbandonNamesATempFileItCouldNotRemove`, `TestUpdateFromReleaseSweepStaysInHome`, `TestReportRedactShapes` all assert behavior that would fail if the corresponding fix were reverted.
3. **Standards conformance** — pass. `%w` error wrapping used in `abandon`/`pushError`; doc comments on `WriteAtomic` and the VENDORED.md table were updated with the code/data they describe, including an explicit note that round 2 had the close/rename ordering reversed.
4. **Directory layout** — pass, unchanged from round 2 (`sweep.go`/`sweep_test.go` beside their caller in `internal/harness/`).
5. **Simplicity / DIY-except-security** — pass. No new dependency; `abandon` is a small stdlib-only helper. No auth/crypto surface touched.
6. **Decoupling** — pass, unchanged from round 2 (`WriteAtomic`/`SweepStaleDirs` promoted to single shared implementations, no reach-ins).

## Not examined

- No new test directly exercises "a close failure leaves `dest` untouched" or "a rename failure leaves `dest` untouched" for `WriteAtomic` — the fix restores an invariant that held before the round-2 regression, verified here by code reading rather than a dedicated failure-path test. This is a gap, not a finding: the restored ordering matches the settings writer's long-standing, already-tested pattern.
- `skills/VENDORED.md`'s `gke-basics` split-commit claim (`5d9d054`, 2026-06-23) was checked for internal consistency only, not against a live fetch of `github.com/google/skills` (no network access in this review) — same caveat as round 2.
- `govulncheck`, `TestHarnessMap`, and the repo test suite outside `internal/harness`/`cmd/sh-agent` were not re-run (out of scope: no other packages appear in the diff).
