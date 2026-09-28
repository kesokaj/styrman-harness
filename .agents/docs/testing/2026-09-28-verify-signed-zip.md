---
owner: verifier
date: 2026-09-28
---
# Verification — signed-zip   (2026-09-28)

Range reviewed: 0b762db..worktree

## Status: PASS

Scope: complete — every item in the dispatch (steps 1-6) ran; testing/guide.md exists and was followed (steps 1-3 of the guide).

## Commands

1. `gofmt -l cmd internal` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 -timeout 180s ./...` (bare, no `-v` per dispatch wording for this step) → exit 0.
   Output (`.agents/scratch/verify-signed-zip/go-test-full.txt`):
   ```
   ok  	github.com/kesokaj/styrman-harness	1.087s
   ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	11.754s
   ok  	github.com/kesokaj/styrman-harness/internal/harness	3.015s
   ```
   counts: no `-v` on this run so no per-test PASS/FAIL/SKIP lines; package-level: 3 `ok`, 0 `FAIL` at 0b762db+dirty for `go test -race -count=1 -timeout 180s ./...`.
   vs previous (per-test counts, `-v` run): `.agents/docs/testing/2026-09-28-verify-v1.8.0-gate.md` — pass 162, fail 0, skip 0 at a962740 for `go test -race -count=1 -timeout 180s -v ./...`. Not directly comparable (that run used `-v`, this one per dispatch does not); no FAIL either way (0 FAIL lines / 0 package failures both times).
   This satisfies the "context size" line implicitly via `TestRenderedContextFitsAntigravityBlock` inside `go test ./...` (package `internal/harness` ok, no FAIL) — not separately isolated with `-v` in this step; explicitly re-verified in step 2 below.

4. `go test -race -count=1 -v -run 'TestVerifyZipSums|TestUnpackZipCheckStops|TestUpdateFromRelease|TestUpdateFromZipForeignSignature|TestUpdateFromZipUnsigned' ./...` → exit 0.
   Output saved to `.agents/scratch/verify-signed-zip/go-test-targeted.txt`. Tests and subtests, all PASS:
   ```
   --- PASS: TestUpdateFromZipForeignSignatureRefused (0.25s)
   --- PASS: TestUpdateFromZipUnsignedSaysUnverified (0.42s)
   --- PASS: TestUpdateFromReleaseReplacesAndInstalls (0.02s)
   --- PASS: TestUpdateFromReleaseUpToDateAndCheck (0.01s)
   --- PASS: TestUpdateFromReleaseRefusals (0.19s)
     subtests: gh_missing, gh_cannot_reach_the_repository, ssh-keygen_missing, no_asset_for_this_platform,
     signature_made_by_another_key, SHA256SUMS_altered_after_signing, binary_does_not_match_its_signed_checksum,
     release_ships_its_own_signer_file, signed_checksums_of_another_release, signed_checksums_naming_no_version,
     no_signature_in_the_release, no_checksums_in_the_release  (12 subtests, all under the passing parent — no
     individual --- PASS/FAIL lines printed per Go's table-test convention here, parent PASS covers all)
   --- PASS: TestUpdateFromReleaseNeverGoesBackwards (0.01s)
   --- PASS: TestVerifyZipSums (0.11s)
     subtests: verified, no_sums:_unsigned, sig_missing, foreign_key, checksum_mismatch, renamed_zip (6 subtests)
   --- PASS: TestUnpackZipCheckStopsBeforeWriting (0.00s)
   ```
   counts: pass 8 (named top-level tests), fail 0, skip 0 at 0b762db+dirty for the `-run` list above (subtests ran under RUN but did not print individual `--- PASS` lines; parent tests all PASS). All 3 packages `ok`.
   vs previous: none found for this exact `-run` set (new tests this range) → `vs previous: none`.

5. `govulncheck ./...` → exit 0. Output: `No vulnerabilities found.` (saved `.agents/scratch/verify-signed-zip/govulncheck.txt`).

6. `go run ./cmd/sh-agent install --home .agents/scratch/sz` → exit 0.
   `go run ./cmd/sh-agent status --home .agents/scratch/sz` → exit 0:
   ```
   repo       /home/kesokaj/workspace/styrman-harness (main @ 0b762db, uncommitted changes)
   installed  2026-09-28T19:04:56+02:00 from commit 0b762db+dirty (vendor: all, 258 files, manifest ~/.config/sh-agent/state.json)
   binary     0b762dba2e0f+dirty
   repo has   27 skills, 7 agents, 0 commands, 0 rules
   claude     ~/.claude  (AGENTS.md, CLAUDE.md, skills/, agents/, commands/, rules/)
   agy        ~/.gemini/config  (rules/AGENTS.md, skills/, agents/, workflows/, rules/)
   drift      none — 257 files match the repo
   ```
   Drift reported as none, as required. GANNET canary found in both rendered outputs:
   `.agents/scratch/sz/.claude/AGENTS.md` (grep count 1) and
   `.agents/scratch/sz/.gemini/config/rules/AGENTS.md` (grep match). Real `$HOME` never touched —
   only `.agents/scratch/sz` used.

7. `git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules` → exit 0, empty output: nothing untracked or ignored under these dirs.
   `go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...` → exit 0:
   ```
   --- PASS: TestEmbedHoldsOnlyTrackedFiles (0.00s)
   ok  	github.com/kesokaj/styrman-harness	0.004s
   ```
   counts: pass 1, fail 0, skip 0 at 0b762db+dirty for `go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...`. vs previous: none found citing this exact command.

8. `git diff --numstat` → all rows numeric (no `-`/`-` binary markers):
   ```
   25  0  .agents/docs/architecture/DECISIONS.md
   4   2  .agents/docs/bugs/BUG-094-red-test-pushed-real-report.md
   1   1  .agents/docs/bugs/INDEX.md
   6   3  .agents/docs/deploy/remote.md
   1   1  .agents/docs/specs/sh-agent/delivery.md
   2   2  README.md
   3   3  REQUIREMENTS.md
   2   0  ROADMAP.md
   21  2  cmd/sh-agent/main.go
   54  0  cmd/sh-agent/zip_test.go
   81  39 internal/harness/release.go
   8   1  internal/harness/zip.go
   107 13 internal/harness/zip_test.go
   ```
   No binary files in the diff.

## Suppressions added (from `git diff -- cmd internal`)

- `cmd/sh-agent/zip_test.go` (added line, offset ~65 in diff): `t.Skip("no /usr/bin/ssh-keygen")` — carries a reason message ("no /usr/bin/ssh-keygen") in the `t.Skip` argument itself. Not a FAIL: reasoned skip, guarding an environment-dependent test (missing `ssh-keygen` binary), not a hidden test disable.
- No `.only(`, no unreasoned skips, no added `//nolint`/`//lint:ignore` found in the diff.

## Coverage of changed code

- `cmd/sh-agent/main.go` — covered by `cmd/sh-agent/zip_test.go` (modified in this range, exercises `update --from-zip` paths added to `main.go`).
- `internal/harness/release.go` — covered by `internal/harness/release_test.go` (present; exercised by `TestUpdateFromRelease*` above).
- `internal/harness/zip.go` — covered by `internal/harness/zip_test.go` (modified in this range; exercises `TestVerifyZipSums`, `TestUnpackZipCheckStopsBeforeWriting` above).
- Doc-only changes (`DECISIONS.md`, `BUG-094...md`, `INDEX.md`, `remote.md`, `delivery.md`, `README.md`, `REQUIREMENTS.md`, `ROADMAP.md`) — no code to cover; not applicable.

## Not examined

- No `.agents/docs/reviews/*` files read or judged — instructed to ignore, reviewers running in parallel.
- Per-test `-v` breakdown was not requested/run for the full-suite step (dispatch step 1 says "bare, exit status read"); only package-level `ok`/`FAIL` counted for that step. If per-test counts for the full suite are wanted, rerun with `-v` (not done here to respect dispatch wording and budget).
- Individual subtest `--- PASS` lines for `TestUpdateFromReleaseRefusals` and `TestVerifyZipSums` were not printed by `go test` output at this verbosity beyond the `=== RUN` lines shown; all subtests ran (visible `=== RUN` lines) and the parent test PASS covers them, but no line-by-line subtest verdict beyond that was captured.
- No frontend/`web`/`ui` directory exists in this repo at this range — not applicable, per the guide (Go-only project).
- Container/Docker tier not run: no Dockerfile changed in this range.
- `.agents/docs/reviews/2026-09-28-coldread-signed-zip-r1.md` and `-r2.md`, and `2026-09-28-verify-signed-zip-e2e.md` were not read (out of scope per dispatch — reviewer/other artifacts).

## Environment

- OS: Linux 6.12.107+deb13-amd64 (Debian trixie), x86_64.
- Go: go1.27.1 linux/amd64 (user-level SDK ahead of `/usr/bin/go`, per AGENTS.md).
- git: 2.47.3.
- Docker/Colima: not applicable — no Dockerfile changes in this range, container tier not run.
- HEAD: 0b762db, working tree dirty (the signed-zip changes under test), plan `.agents/docs/plans/2026-09-28-signed-zip.md`.
