---
owner: verifier
date: 2026-09-28
---
# Verification — update-from-zip Phase 1   (2026-09-28)

Range reviewed: 1cb2608..worktree (uncommitted Phase 1 work: `internal/harness/zip.go`, `zip_test.go`, `Source.Origin` in `source.go`, `Origin` precedence in `apply.go`, zip-manifest skip and not-found prefix in `repo.go`; per plan D2 there is no commit at this checkpoint)

## Status: PASS

Scope: complete — all four commands in the dispatch ran, testing/guide.md exists (used, unchanged), coverage of the new code mapped, uncovered branches named, must-fail mutation demonstrated and restored.

## Commands
1. `gofmt -l cmd internal` → exit 0, no output (clean).
2. `go vet ./...` → exit 0.
3. `go test -race ./...` (run as `go test -race -count=1 -v ./...` to get counts; bare exit status read) → exit 0. `ok github.com/kesokaj/styrman-harness 1.085s`, `ok .../cmd/sh-agent 6.474s`, `ok .../internal/harness 2.840s`. Output saved to `.agents/scratch/verify-update-from-zip/go-test-full.txt`.
   - counts (top-level `^--- (PASS|FAIL|SKIP):`): pass 139, fail 0, skip 0
   - counts (incl. indented subtests `^    --- (PASS|FAIL|SKIP):`): pass 127 additional subtest lines (266 total `--- PASS:` lines, matching 266 `=== RUN` lines)
   - counts: pass 139, fail 0, skip 0 at 1cb2608+dirty for `go test -race -count=1 -v ./...` (top-level); pass 266 total incl. subtests
   - vs previous: `.agents/docs/testing/2026-09-28-verify-v1.7.0-gate-3.md` (pass 132, fail 0, skip 0 at 56621b6 for the same command, top-level count; that log separately noted 251 total incl. subtests), pass 132→139 top-level, 251→266 incl. subtests, skip 0→0. Rise in pass count from the new `zip.go`/`zip_test.go` and related tests (7 new top-level tests: `TestUnpackZipGitHubTopFolder`, `TestUnpackZipRootLayout`, `TestUnpackZipRefuses`, `TestUnpackZipLimits`, `TestUnpackZipNotHarness`, `TestInstallRecordsZipOrigin`, `TestResolveSourceSkipsZipManifest`, plus their 8 subtests). No fall in pass, no rise in skip.
4. `go test -race -count=1 -v -run 'TestUnpackZip|TestInstallRecordsZipOrigin|TestResolveSourceSkipsZipManifest' ./internal/harness` → exit 0. Output saved to `.agents/scratch/verify-update-from-zip/go-test-targeted.txt`.
   - `TestUnpackZipGitHubTopFolder` — PASS
   - `TestUnpackZipRootLayout` — PASS
   - `TestUnpackZipRefuses` — PASS, subtests: `parent_element`, `absolute_name`, `backslash_name`, `symlink_inside_a_tree`, `device-like_mode_inside_a_tree` — all PASS
   - `TestUnpackZipLimits` — PASS, subtests: `too_many_entries`, `one_file_too_large`, `total_too_large` — all PASS
   - `TestUnpackZipNotHarness` — PASS
   - `TestInstallRecordsZipOrigin` — PASS
   - `TestResolveSourceSkipsZipManifest` — PASS
   - counts: pass 7 top-level + 8 subtests = 15, fail 0, skip 0 at 1cb2608+dirty for the targeted `-run` command; vs previous: none (no earlier log ran this exact targeted `-run` expression).
5. `go test -coverprofile=.agents/scratch/zipcov.out ./internal/harness` → exit 0, `coverage: 88.5% of statements`.
   `go tool cover -func=.agents/scratch/zipcov.out | grep zip.go`:
   - `UnpackZip` 87.8%
   - `unpackFile` 80.0%
   - `zipNameSafe` 90.9%
   - `zipRoot` 100.0%
6. Must-fail mutation (review-standards §2): saved `internal/harness/zip.go` to `.agents/scratch/zip.go.orig` (sha256 `88afc74408b4e4bb8d221fc48aa5edff36e01c3793d30fb5ef53a2bce4cb820c`), changed the unique text `if n > maxZipFileBytes {` (line 129) to `if n > maxZipFileBytes*1000 {`. `go build ./...` → exit 0 (mutation compiles). `go test -count=1 -run TestUnpackZipLimits ./internal/harness` → test binary FAIL: `--- FAIL: TestUnpackZipLimits/one_file_too_large`, `zip_test.go:203: err = <nil>, want a per-file size refusal naming big.md`. Named subtest failed as expected. Restored `internal/harness/zip.go` from the saved copy; `git diff --stat internal/harness/zip.go` shows no change (file remains untracked/new, identical to before the mutation); `sha256sum` of the restored file and the saved copy both `88afc74408b4e4bb8d221fc48aa5edff36e01c3793d30fb5ef53a2bce4cb820c` (match). Re-ran `go test -count=1 -run TestUnpackZipLimits ./internal/harness` → `ok` (passes again, confirming restoration).
7. (bonus, not in dispatch but run for context) `govulncheck ./...` → "No vulnerabilities found."

## Root cause  (FAIL only)
N/A — Status is PASS.

## Coverage of changed code
- `internal/harness/zip.go` — covered by `internal/harness/zip_test.go` (`TestUnpackZipGitHubTopFolder`, `TestUnpackZipRootLayout`, `TestUnpackZipRefuses`, `TestUnpackZipLimits`, `TestUnpackZipNotHarness`), 87.8-100% per function; uncovered branches identified by reading the source against the zero-count lines in the coverage profile:
  - `UnpackZip` line 59-60: `os.MkdirTemp` failure path — not exercised (would need an unwritable parent directory).
  - `UnpackZip` line 79-82: `os.MkdirAll(dest, 0o755)` failure inside the directory-entry case — not exercised.
  - `UnpackZip` line 98-99: the post-unpack `IsRepo(dir)` false branch ("not a harness zip: the unpacked tree has no context/AGENTS.md or skills/") — not exercised; every test archive that passes `zipRoot`'s pre-check also satisfies `IsRepo` after unpack, so this second check has no test that reaches it via a mismatch between the two.
  - `unpackFile` line 110-111: `os.MkdirAll(filepath.Dir(dest))` failure — not exercised.
  - `unpackFile` line 114-115: `f.Open()` failure (a corrupt/truncated zip stream) — not exercised.
  - `unpackFile` line 120-121: `os.OpenFile` `O_EXCL` failure (a duplicate entry name in the archive) — not exercised; the doc comment calls this out ("a name that appears twice is refused, never overwritten") but no test constructs a zip with a duplicate path.
  - `unpackFile` line 127-128: `w.Close()` returning an error after a successful `io.Copy` (e.g. a full disk) — not exercised.
  - `zipNameSafe` line 151-152: the `!filepath.IsLocal(...)` branch ("not a local path") — not exercised; every refused-name subtest in `TestUnpackZipRefuses` (`parent element`, `absolute name`, `backslash name`) trips one of the earlier explicit checks first, so this fallback check has no test proving it fires on its own.
- `internal/harness/source.go` (`Source.Origin`) — covered by `TestInstallRecordsZipOrigin` and `TestResolveSourceSkipsZipManifest`.
- `internal/harness/apply.go` (`Origin` precedence) — covered by `TestInstallRecordsZipOrigin`.
- `internal/harness/repo.go` (zip-manifest skip, not-found message prefix) — covered by `TestResolveSourceSkipsZipManifest` and `TestUnpackZipNotHarness`.
- `ROADMAP.md` — prose/plan-tracking change, not code; no test applicable.

No added `//nolint`, `//lint:ignore`, `t.Skip`, `.skip(`/`.only(`, or similar suppressions found in the diff (`zip.go`, `zip_test.go`, `source.go`, `apply.go`, `repo.go`).

## Not examined
- Phases 2+ of the update-from-zip plan (not part of this checkpoint).
- The full `skill-matrix.md` / sandbox CLI-start smoke checks (out of scope for this dispatch, which named only the four commands above).
- The eight uncovered branches listed above were read, not exercised by a new test (finding them is this run's job; adding tests is not).
- `TestInstallRecordsZipOrigin` and `TestResolveSourceSkipsZipManifest` internals were not read line-by-line beyond confirming they pass and target `source.go`/`apply.go`/`repo.go`; a full code review of those files is `reviewer`'s job, not verifier's.

## Environment
Linux 6.12.107+deb13-amd64, go1.27.1 linux/amd64 (user-level SDK per AGENTS.md). `govulncheck ./...`: "No vulnerabilities found." Working tree dirty as described in Range reviewed (Phase 1 uncommitted); tree was frozen for this run — the only transient change was the mutation in step 6, restored and hash-verified before this log was finalized.
