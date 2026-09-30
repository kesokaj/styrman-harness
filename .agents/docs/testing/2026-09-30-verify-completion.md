---
owner: verifier
date: 2026-09-30
---
# Verification — worktree completion   (2026-09-30)

Range reviewed: working tree — changed files: .agents/docs/deploy/local.md agents/verifier.md internal/harness/harness_test.go internal/harness/report.go internal/harness/report_test.go internal/harness/state.go internal/search/global_test.go internal/search/reconcile.go; untracked, not ignored: none outside .agents/docs; read the diff with git diff HEAD -- <files>.

## Status: PASS

Scope: complete — all setup commands (AGENTS.md), testing guide § 3d (search index tier), mutation protocol (a)/(b)/(c), and the `filepath.Glob` sweep ran; the verifier-change behavioural tier is cited from the existing `2026-09-30-verifier-long-suites-behavioural.md` log per the dispatch, not rerun. `testing/guide.md` exists (dated 2026-09-17, unchanged — no step in it needed updating for this change).

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0, run in the foreground with output redirected to `.agents/scratch/verify-completion/go-test-1.txt` (no pipe), all 4 packages `ok`:
  `ok github.com/kesokaj/styrman-harness 1.108s`, `ok .../cmd/sh-agent 14.742s`, `ok .../internal/harness 3.343s`, `ok .../internal/search 6.965s`.
  counts: pass 374, fail 0, skip 2 at e83c2d0+dirty for `go test -race -v -count=1 ./...` (grep `-cE '^\s*--- (PASS|FAIL|SKIP)'`, indented subtests included).
  vs previous: `.agents/docs/testing/2026-09-30-verify-v1.14.0-gate.md`, pass 373→374, skip 2→2. explained: +1 pass matches the one new test this change adds (`TestConfigViewListsProposalsByName`); harness_test.go and report_test.go additions ran inside their existing test functions (`TestLastRetro`, `TestReportRedactShapes`) rather than adding new `Test...` entries, so only one net new top-level/subtest line appears. Skips unchanged (`TestRecall`, `TestHelperIndex`, both with reasons, same as before).
- `go run ./cmd/sh-agent install --home .agents/scratch/home` → exit 0. `266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed`.
- `go run ./cmd/sh-agent status --home .agents/scratch/home` → exit 0. `binary e83c2d028b46+dirty` (no warning), `drift none — 265 files match the repo`.
- `govulncheck ./...` → exit 0. "No vulnerabilities found... 1 vulnerability in modules you require, but your code doesn't appear to call these" (the known, pre-existing, non-reachable finding AGENTS.md already documents — not new).
- `go test -race ./internal/search ./cmd/sh-agent` (guide § 3d unit tier) → exit 0, both `ok`.
- `SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v` → exit 0, PASS, mean recall 0.624 over 41 (informational, not a gate, per the guide).
- `.agents/scratch/sh-agent index --rebuild --home .agents/scratch/s` (this checkout) → 1148 files, 1853 ms — under the 10 s bound.
- `.agents/scratch/sh-agent search --home .agents/scratch/s worktree completion` (warm) → 36 ms wall time — under the 100 ms bound.
- `git grep -n 'filepath.Glob' -- '*.go' | grep -v _test` → no output (exit 1) — confirms no `filepath.Glob` remains in non-test Go source, matching the fix's rule ("listed by name, never through a pattern built from the [home/config] path").

### Mutation protocol (review-standards § 2), each restored exactly

- (a) `internal/harness/state.go` `LastRetro` reverted to `filepath.Glob(filepath.Join(home, ".config", "sh-agent", "proposals", "*.md"))` (with the now-unused `strings` import dropped so it compiles) → `go build ./...` succeeded, `go test ./internal/harness/... -run TestLastRetro -v` FAILED: `harness_test.go:841: a proposal under a home with [ in its path should be found`. Restored via `cp` from a pre-mutation backup; `git diff --stat -- internal/harness/state.go` back to `13 ++++++++++++-` (unchanged from before the mutation).
- (b) `internal/search/reconcile.go` `configView` reverted to the `filepath.Glob(filepath.Join(dir, "proposals", "*.md"))` listing → `go build ./...` succeeded, `go test ./internal/search/... -run TestConfigViewListsProposalsByName -v` FAILED: `global_test.go:180: proposals should be listed by name, .md only, got "feedback.jsonl retros.jsonl"`. Restored; `git diff --stat` back to `15 +++++++++++----` (unchanged).
- (c) `internal/harness/report.go`: removed the four new prefix needles (`AGE-SECRET-KEY-`, `dop_v1_`, `doo_v1_`, `dor_v1_`) and the nine new keyword needles (jwt short, google refresh token, discord webhook, teams webhook, sendgrid key, azure account key, cookie header, vault token, plus the comment) → `go build ./...` succeeded, `go test ./internal/harness/... -run TestReportRedactShapes -v` FAILED on many of the newly-added shapes (age key, DigitalOcean tokens, jwt short, discord webhook, azure keys, cookie header, vault token all logged as unredacted, `recs=[]`). Restored; `git diff --stat` back to `15 +++++++++++++++` (unchanged).

After restoring all three, `go build ./...` and `go vet ./...` both clean, and `git status --short` shows exactly the same eight tracked files as at the start plus this log and the pre-existing untracked review/testing files — the tree did not move under the run.

## Root cause
N/A — all commands passed; no FAIL to explain.

## Coverage of changed code

- `internal/harness/state.go` (`LastRetro`) — covered by `internal/harness/harness_test.go` `TestLastRetro`; mutation-confirmed above.
- `internal/search/reconcile.go` (`configView`) — covered by `internal/search/global_test.go` `TestConfigViewListsProposalsByName`; mutation-confirmed above.
- `internal/harness/report.go` (new needles) — covered by `internal/harness/report_test.go` `TestReportRedactShapes`; mutation-confirmed above.
- `.agents/docs/deploy/local.md`, `agents/verifier.md` — prose describing the same redaction/verifier behaviour; not code, no test applies. The verifier.md behavioural change (foreground long-suite runs, killing background jobs before hand-back) is proven by the existing `.agents/docs/testing/2026-09-30-verifier-long-suites-behavioural.md` (3 of 3 new-definition runs foreground, 0 of 3 background, cited per dispatch, not rerun here). This run itself ran `go test -race -v -count=1 ./...` in the foreground with `> file 2>&1` (no pipe) per that same updated rule, which is a second live demonstration of the new behaviour.
- Suppressions: `git diff HEAD` over the changed files shows no added `//nolint`, `//lint:ignore`, `t.Skip`, `.only(`, or `.skip(` in this change set.

## Not examined

- The full `sh-agent-matrix.md` (§ 3c) was not re-run/re-rowed: this change does not add or change a verb's behavior or a refusal message (it fixes internal path-listing inside the existing `search`/`index`/`report` machinery), so per the guide's trigger ("a new or changed verb also gets its rows") no new matrix rows were added. Decision recorded here rather than editing that file.
- `report --dry-run` (guide § 3b) was not rerun against the real feedback wiring in this pass — the redaction-shape change is covered by the unit test `TestReportRedactShapes` (mutation-confirmed) and the guide marks a real end-to-end `report` run as an e2e-tier check outside a bare verifier dispatch.
- Sandbox and real-home CLI tiers (guide §§ 4–5) — explicitly out of scope per the dispatch.
- `go run ./cmd/sh-agent list` / `render` (guide § 2) were not run: no file under `context/`, `agents/`, or `skills/` changed in this diff.
- The by-hand install chains and `2026-09-17-release-binary-e2e.md` tier — not applicable to this change, not run.

## Environment

- Linux, `go version go1.27.1 linux/amd64` (the user-level SDK ahead of `/usr/bin/go`, per AGENTS.md).
- `git version 2.47.3`.
- HEAD at `e83c2d0`, working tree dirty with exactly the eight named files (plus the pre-existing untracked review/testing docs and this new log).
- No Docker/Colima involved (no Dockerfile in this change).
