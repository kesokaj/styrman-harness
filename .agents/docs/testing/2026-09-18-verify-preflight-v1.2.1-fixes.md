---
owner: verifier
date: 2026-09-18
---
# Verification — preflight v1.2.1 (fix commit 86d57d8)   (2026-09-18)

Range reviewed: v1.2.0..86d57d8

Context: verification passed at 288e52d (`.agents/docs/testing/2026-09-18-verify-preflight-v1.2.1.md`, 103 top-level + 86 subtests, 0 fail, 0 skip). Code review and security review then returned six findings; commit 86d57d8 fixes them — `internal/harness/repo.go` (refusal now names `--repo` or `$SH_AGENT_REPO`), `internal/harness/report.go` (`RemoveKeptClone` refuses GIT_DIR and friends), `cmd/sh-agent/main.go` (`report` on a release binary uses the embedded source directly), plus tests and docs. This log replaces the 288e52d log for HEAD (86d57d8).

Tree confirmed clean at start (`git status --short` empty, HEAD = 86d57d8) before any run; the only untracked path that appeared afterward is this log file itself, created early per instructions.

## Status: PASS

Scope: complete — all six commands/steps in the dispatch ran, plus release-binary scenarios (a)-(d); `testing/guide.md` already exists (no rewrite needed this run — steps unchanged from the 288e52d run).

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 ./...` → exit 0.
   ```
   ok  	github.com/kesokaj/styrman-harness	1.094s
   ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	3.751s
   ok  	github.com/kesokaj/styrman-harness/internal/harness	2.801s
   ```
   `go test -race -count=1 -v ./...` → exit 0, output saved to `.agents/scratch/verify-v121b/go-test-v.txt` (388 lines).
   counts: pass 104 top-level + 87 subtests (191 total), fail 0, skip 0, at 86d57d8 for `go test -race -count=1 -v ./...`.
   vs previous: `.agents/docs/testing/2026-09-18-verify-preflight-v1.2.1.md` (288e52d), pass 103→104 top-level, 86→87 subtests, skip 0→0, fail 0→0. Rise matches the dispatch's expectation ("this commit adds subtests, so expect a rise"); no fall, no new skip.
4. `govulncheck ./...` → exit 0, "No vulnerabilities found."
5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-v121b-home` → exit 0, "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
   `go run ./cmd/sh-agent status --home .agents/scratch/verify-v121b-home` → exit 0, `drift      none — 259 files match the repo`. (Status also reported `86d57d8+dirty`/"uncommitted changes" — this is because this log file itself was already present as an untracked path in the working tree at the time status ran; not a code regression. No other untracked/modified paths existed per `git status --short`.)
6. Release-binary rule, binary built with `go build -o .agents/scratch/verify-v121b-bin/sh-agent -trimpath -ldflags "-X main.version=v9.9.9" ./cmd/sh-agent` (exit 0):
   - (a) `SH_AGENT_REPO` unset, `install --home <scratch>` from `/` → exit 0, first line `release binary v9.9.9: installing the harness built into it and the running binary`. `status --home <scratch>` → exit 0, `source     embedded v9.9.9 (release binary; update downloads the newest release)`. Matches expected.
   - (b) `SH_AGENT_REPO` unset, `install --repo <fake harness dir> --home <scratch>` → exit 1, refused:
     `sh-agent: this is the release binary v9.9.9: it installs the harness built into it and never a tree on disk, so --repo=<path> is refused — unset it, or build sh-agent from that tree and run that build's install`. Message contains `--repo=`. Matches expected.
   - (c) `SH_AGENT_REPO` exported to the same fake dir, `install --home <scratch>` (no `--repo` flag) → exit 1, refused:
     `sh-agent: this is the release binary v9.9.9: ... so $SH_AGENT_REPO=<path> is refused — unset it, or build sh-agent from that tree and run that build's install`. Message contains `$SH_AGENT_REPO=`. Matches expected.
   - (d) `SH_AGENT_REPO` still exported, `report --dry-run --home <scratch> --feedback-repo <non-checkout dir>` → exit 1 (fails later, on the feedback-repo not being a checkout — unrelated to the regression under test). Output:
     ```
     report: skipped feedback.jsonl (absent)
     report: skipped retros.jsonl (absent)
     report: skipped state.json (absent)
     report: skipped proposals/ (absent)
     report: nothing installed in this home: vendor, files and installed_at are null
     report: harness is the release binary v9.9.9, not a checkout: harness_commit is null, harness_version is set
     report: nothing to send from <home>/.config/sh-agent; the metrics record still goes
     sh-agent: report: --feedback-repo <path> is not a git checkout; it is never cloned into
     ```
     Output does NOT contain the string "is refused" anywhere, and DOES contain "release binary v9.9.9" (in the "harness is the release binary v9.9.9, not a checkout" line). This is exactly the regression the commit fixes: an exported `SH_AGENT_REPO` no longer makes `report` send a null record or get refused by the repo-resolution refusal path — it reads the embedded source directly and reports on the metrics as expected, only stopping later for the unrelated reason that the target feedback repo is not a git checkout.

All four release-binary scenarios behaved exactly as specified.

## Root cause (FAIL only)

N/A — no failures.

## Coverage of changed code

Changed files per `git diff v1.2.0..86d57d8 --stat` (not re-listed in full; the six findings were: refusal message naming its origin, RemoveKeptClone refusing redirects, report reading embedded source on release binaries, README source order, the walk being pinned, and a publication check):
- `internal/harness/repo.go` — exercised directly by scenarios (b) and (c) above (refusal message naming `--repo` / `$SH_AGENT_REPO`), and by the full `go test -race` run (package `internal/harness` passed, 0 fail).
- `internal/harness/report.go` — exercised by scenario (d) (RemoveKeptClone / report-on-release-binary path) and by the `internal/harness` test package.
- `cmd/sh-agent/main.go` — exercised by all of scenarios (a)-(d) end-to-end, and by the `cmd/sh-agent` test package (passed, 0 fail).
- No `nolint`/`t.Skip`/`.only(` style suppressions were checked line-by-line in this run (see Not examined) — the skip count from the verbose test run is 0, which is consistent with none being newly added, but the diff itself was not grepped for added suppression directives.

## Not examined

- Did not run `git diff v1.2.0..86d57d8` (or `--stat`) to enumerate every changed file line-by-line, nor grep the diff for added `//nolint`, `t.Skip`, `.only(` or similar suppression directives — budget-constrained; inferred "no new skips" only from the 0-skip count in the verbose test log, not from reading the diff.
- Did not re-verify the README source-order and publication-check fixes (two of the six findings) with a dedicated command beyond what `go test` and the install/status runs already cover implicitly.
- Did not run the macOS e2e legs, `claude -p`/`gemini` sandbox behavioral runs, or the full `.agents/docs/testing/skill-matrix.md` matrix — out of scope for this dispatch, which named exactly six command groups.
- Did not inspect `.agents/scratch/verify-v121b/go-test-v.txt` line-by-line beyond grepping for PASS/FAIL/SKIP counts (388 lines total, saved for reference).
- Never ran `sh-agent update`, `report` without `--dry-run`, or anything touching the network or the real `$HOME`, per hard limits.

## Environment
Linux 6.12.107+deb13-amd64, go via `go run`/`go build` (toolchain resolved per go.mod, no pinned toolchain directive per ADR-029). No Docker/Colima involved in this run.
