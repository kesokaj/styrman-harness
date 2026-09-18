---
owner: verifier
date: 2026-09-18
---
# Verification — preflight v1.2.0 gate (fixes commit)   (2026-09-18)

Range reviewed: v1.1.0..a1d2ee8

Commit under test: a1d2ee80db5580e9f499ed0bd1dc941ee36e4eda ("v1.2.0 gate fixes: the release block uploads one asset per call, the kept-clone settling command covers the same refs as its count and pastes no path, preflight's non-blocking rule moved before the tag")

This log replaces .agents/docs/testing/2026-09-18-verify-preflight-v1.2.0.md (which described bba384d, not HEAD) as the gate's verification row, per the reviewer's REQUEST CHANGES on report.go / report_test.go and the three doc fixes in a1d2ee8. Tree was clean and frozen for this run (git status --short at start: clean).

## Status: PASS

Scope: complete for the six commands the dispatch named. Did not run the full behavioural CLI matrix (`.agents/docs/testing/skill-matrix.md`) or the two-CLI real-home smoke checks — those are out of scope for this gate run per the dispatch (which named exactly the six commands below).

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 -v ./...` → exit 0.
   - Packages: `github.com/kesokaj/styrman-harness` ok (1.092s), `.../cmd/sh-agent` ok (3.688s), `.../internal/harness` ok (2.775s).
   - counts: pass 102 top-level, 86 subtests (188 `--- PASS:` lines total), fail 0, skip 0, at a1d2ee8 for `go test -race -count=1 -v ./...`.
   - vs previous: .agents/docs/testing/2026-09-18-verify-preflight-v1.2.0.md (bba384d) reported 102 top-level tests + 86 subtests, 0 fail, 0 skip. This run reproduces that exactly: 102 top-level → 102, 86 subtests → 86, no fall in pass, no rise in skip.
4. `go test -race -count=1 -run 'RemoveKeptClone' ./internal/harness/ -v` → exit 0. All 11 subtests of `TestRemoveKeptClone` PASS, none skipped:
   - clean_clone_is_removed — PASS
   - dirty_clone_is_kept_and_named — PASS
   - clone_with_unpushed_commits_is_kept_and_named — PASS
   - another_repository_is_kept_and_named — PASS
   - unpushed_work_on_a_branch_with_no_upstream_is_kept — PASS
   - a_stale_tracking_ref_keeps_the_clone_without_claiming_lost_work — PASS (did not skip; its reproduction setup succeeded)
   - an_uncountable_clone_is_kept — PASS
   - stashed_work_is_kept — PASS
   - a_symlink_there_is_kept — PASS
   - a_relative_home_still_resolves — PASS
   - nothing_there_says_nothing — PASS
5. `govulncheck ./...` → exit 0, "No vulnerabilities found."
6. `go run ./cmd/sh-agent install --home .agents/scratch/verify-v120b-home` → exit 0. "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed." Then `go run ./cmd/sh-agent status --home .agents/scratch/verify-v120b-home` → exit 0, `drift  none — 259 files match the repo` (the 260th is the manifest/state.json itself, not diffed against itself). No stray writes outside `.agents/scratch/`: `git status --short --untracked-files=all` after both commands showed only this report file as untracked; no other path, commit, or tag appeared.

## Root cause (FAIL only)
N/A — no failures. This run confirms the a1d2ee8 fixes (report.go / report_test.go changes plus the three doc updates) did not regress test counts, vet, gofmt, or vulnerability scan versus the bba384d baseline, and the previously-flagged RemoveKeptClone edge case (stale tracking ref) still passes rather than skipping.

## Coverage of changed code
Per the reviewer's REQUEST CHANGES round, a1d2ee8 touched `internal/harness/report.go`, `internal/harness/report_test.go`, and three docs. `report_test.go` is exercised by `go test -race -count=1 -v ./...` (package `.../internal/harness`, ok, no failures) as part of the 102/86 count above; no test file was skipped or excluded. Doc-only changes carry no test coverage requirement.

## Not examined
- The full behavioural CLI matrix (`.agents/docs/testing/skill-matrix.md`) — not part of this dispatch's six named commands.
- The two-CLI real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`) — not part of this dispatch.
- `/review --docs` over the doc changes in a1d2ee8 — that is the reviewer's/documentation agent's lane, not verifier's, and was not requested here.
- No behavioural `claude -p --agents ...` run against the changed `internal/harness/report.go` logic beyond unit tests — not requested in this dispatch's command list.

## Environment
- OS: Linux 6.12.107+deb13-amd64 (Debian trixie box named in this repo's AGENTS.md).
- Go: user-level SDK ahead of `/usr/bin/go` per this repo's documented setup (govulncheck reported zero advisories, consistent with go1.24.6/1.25.8+ being in effect).
- Docker/Colima: not exercised in this run (not part of the six named commands).
- Real `$HOME` untouched; all scratch under `.agents/scratch/verify-v120b-home` and `.agents/scratch/verify-v120b-*.txt`; no network reached; no git command altered the working tree.
