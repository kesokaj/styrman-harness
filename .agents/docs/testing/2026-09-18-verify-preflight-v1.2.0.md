---
owner: verifier
date: 2026-09-18
---
# Verification — preflight v1.2.0

Range reviewed: v1.1.0..bba384d

## Status: PASS

Scope: complete — tiers 1-3 from `testing/guide.md` plus the dispatch's step 6
deployed-copy checks for `context/AGENTS.md` and `skills/preflight/SKILL.md` all ran.
Tier 3a/3b (plain-copy install, `report`) were not in scope for this range (no change
to `IsCheckout`/rebuild-decision logic; `report.go`'s change is covered by tier 1's
`go test` and did not warrant a separate e2e `report --dry-run` run since the change
is a message-text fix, not new collection/redaction/push logic). Tier 4 (sandbox CLI
starts) not run — see Not examined.

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -v -count=1 ./...` → all three packages `ok`
   (`github.com/kesokaj/styrman-harness` 1.094s, `.../cmd/sh-agent` 3.689s,
   `.../internal/harness` 2.788s); log at
   `.agents/scratch/verify-v120/go-test.txt`.
   counts: pass (top-level 102, subtests 86, total 188), fail 0, skip 0
   at bba384d for `go test -race -v -count=1 ./...`
   vs previous: `.agents/docs/testing/2026-09-17-verify-preflight-v1.1.0-fixes.md`
   (101 top-level, 85 subtests, 0 fail, 0 skip) → pass 101→102 top-level,
   85→86 subtests, skip 0→0. Rise matches the dispatch's stated expectation
   (one release test added, one `TestRemoveKeptClone` subtest added); no fall,
   no new skip.
4. `govulncheck ./...` → exit 0, "No vulnerabilities found." (go1.27.1 toolchain in
   use, ahead of the go1.25.8 fix threshold in `AGENTS.md`, so the two
   stdlib advisories noted there do not apply here).
5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-v120-home` → exit 0,
   "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed". Then
   `go run ./cmd/sh-agent status --home .agents/scratch/verify-v120-home` → exit 0,
   "drift none — 259 files match the repo" (259 vs 260 written is the settings.json
   file created at install and not diffed for drift — consistent with normal
   status behaviour, not a finding).
   Note: status reported "commit bba384d+dirty" / "uncommitted changes" — this is
   because this verification log itself (`.../2026-09-18-verify-preflight-v1.2.0.md`)
   was an untracked file in the tree at the time install ran (written early per this
   run's own instructions); `git status --short --untracked-files=all` before and
   after showed only that one untracked path, nothing else. Not a release-tree
   finding.
6. `go run ./cmd/sh-agent list` → exit 0, table printed for all 27 skills and 7
   subagents, parses cleanly. `go run ./cmd/sh-agent render` → exit 0, 96 lines
   (under the guide's 100-line ceiling).
   Deployed-copy grep of the scratch home for the two changed context files:
   - `.claude/AGENTS.md`: contains "IN FLIGHT" (1), "KEEPS the data" (1),
     "BOTH directions" (1) — present.
   - `.gemini/config/rules/AGENTS.md`: contains "IN FLIGHT" (1), "KEEPS the data"
     (1), "BOTH directions" (1) — present.
   - `.claude/skills/preflight/SKILL.md` and
     `.gemini/config/skills/preflight/SKILL.md`: both contain
     "upload each artifact as its OWN call" — present in both deployed copies.

## Coverage of changed code

Range v1.1.0..HEAD, 19 files changed (7 commits): mostly `.agents/docs/` trail files
(bug records, changes, handoff, plan, retro, testing log — not code). Code/behaviour
files changed:
- `internal/harness/report.go` (8 lines, message-text fix per BUG-074) — covered by
  `internal/harness/report_test.go`, which changed in the same commit (+48 lines).
- `context/AGENTS.md`, `skills/preflight/SKILL.md`, `README.md`, `ROADMAP.md`,
  `.agents/docs/deploy/{local,remote}.md` — doc/context files, verified via the
  install+grep check in step 6 above (both deployed copies carry the new text).

Suppressions/skips added in this range (`git diff v1.1.0..HEAD -- '*.go'`):
- `internal/harness/report_test.go`: one added `t.Skip("this git still resolves
  @{upstream}; the shape under test needs it to fail")` — carries a reason on the
  call itself. Not a FAIL.
No added `//nolint`, `//lint:ignore`, `@ts-ignore`, `.only(`, or unexplained skip.

## Not examined

- Tier 4 sandbox CLI starts (`.agents/docs/testing/2026-09-03-verify-sandbox.md`) and
  the full skill-matrix (`skill-matrix.md`) were not run in this pass — this log does
  not claim "complete" against that matrix; scope above states which tiers did run.
  Given the range touches `context/AGENTS.md` and `skills/preflight/SKILL.md`, the
  guide calls for one start of each CLI asking a context question that quotes the new
  wording; this run substituted the deployed-file grep in step 6 instead, on the
  dispatch's explicit instruction set (which specified the grep check, not a sandbox
  CLI start, for this release gate). Flagging this substitution explicitly since the
  guide's own tier 4 was not exercised.
- Tier 3a (plain-copy install) and 3b (`report --dry-run` behavioural run) not run —
  judged out of scope per the "Scope" note above; no rebuild-decision-logic or
  report-wiring change in this range beyond the text/format fix already covered by
  `go test`.
- No `docker build` run — no Dockerfile in this repo/range.

## Environment

- OS: Linux 6.12.107+deb13-amd64 (Debian trixie), Linux box.
- Go: go1.27.1 (user-level SDK ahead of `/usr/bin/go`, per repo `AGENTS.md`).
- govulncheck: run against go1.27.1, clean, no advisories.
- No Docker/Colima involved (no Dockerfile in range).
- Scratch used: `.agents/scratch/verify-v120/` (test/list/render logs),
  `.agents/scratch/verify-v120-home/` (install/status scratch home). Real `$HOME`
  untouched. No `sh-agent update` or `sh-agent report` run (network commands not
  invoked at all in this pass). No working-tree-mutating git command run.
