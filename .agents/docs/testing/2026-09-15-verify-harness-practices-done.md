---
owner: verifier
date: 2026-09-15
---
# Verification — harness-practices-done (2026-09-15), re-run after BUG-032 fix

Range reviewed: c192622..8202d6a, which includes the Phase 7 range
641eccd..8202d6a (the Phase 7 checkpoint) as a subrange. This log is both the
Definition of Done's verifier row and the Phase 7 checkpoint's verifier log.
This is a re-run over the previous run's FAIL, which found `cmd/sh-agent/
main.go:309`'s `status` note (hand-written `~/.claude/CLAUDE.md`) had no
covering test; that FAIL log is committed at 8202d6a and kept by git for the
record. 8202d6a adds `TestStatusNamesForeignClaudeMD`
(`cmd/sh-agent/main_test.go`) and BUG-032.

Procedure: `.agents/docs/testing/sh-agent/guide.md`, §1–§3.

## Status: PASS

Scope: §1 (Go checks), §2 (harness self-check), §3 (scratch-home
install/status + agent Records parity) all ran to completion. §3a does not
apply (checked against the whole-range diff, see below). §4 (sandbox CLI
starts) and the full matrix are cited from
`.agents/docs/testing/2026-09-15-verify-harness-practices-e2e.md`, whose
Status and Not examined are read in full below. §5 (real-home smoke) is the
main agent's own separate log, not run here — named as not run.

## Commands

- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0. Output saved to
  `.agents/scratch/verify-harness-practices-done-2/go-test-1.txt`.
  counts: pass 32, fail 0, skip 1 at 8202d6a for `go test -race -v ./...`
  vs previous: `.agents/docs/testing/2026-09-15-verify-harness-practices-done.md`
  (the FAIL log at a33978d, committed in 8202d6a), pass 31→32, skip 1→1. The
  rise of one pass is expected: `TestStatusNamesForeignClaudeMD` is the one
  added test (confirmed present and PASS in the log at line 19-20). No fall
  in pass, no rise in skip.
- `go run ./cmd/sh-agent list` → exit 0. "Skills (27)" / "Subagents (7)"
  headers, matching the plan's expectation.
- `go run ./cmd/sh-agent render | wc -l` → 93 (< 100, matches plan
  expectation of 93).
- `rm -rf .agents/scratch/verify-harness-practices-done-home-2 && mkdir -p ...
  && go run ./cmd/sh-agent install --home
  .agents/scratch/verify-harness-practices-done-home-2` → "260 written, 0
  unchanged, 0 removed, 0 kept — 260 files managed", exit 0.
- `go run ./cmd/sh-agent status --home
  .agents/scratch/verify-harness-practices-done-home-2` →
  `repo       /home/kesokaj/workspace/styrman-harness (main @ 8202d6a, clean)`,
  `binary     8202d6a5105b`, `drift      none — 259 files match the repo`,
  exit 0.
- `## Records` parity: for each of the 7 `agents/*.md` files, grep-counted
  `^## Records` in the source and in the derived
  `.gemini/config/agents/<name>/agent.md` under the scratch home. All 7 pairs
  show source=1, derived=1 (bugs, changelog, documentation, explorer,
  reviewer, security, verifier).

## §3a applicability

Checked: `git diff c192622...8202d6a --stat -- '*.go'` shows exactly four
changed Go files across the whole range: `cmd/sh-agent/main.go` (+3/-3),
`cmd/sh-agent/main_test.go` (+30/-0, new in this fix commit),
`internal/harness/apply.go` (+1/-1), `internal/harness/harness_test.go`
(+3/-0). Read the full diff of all four: `main.go`'s hunks are a comment
sentence edit and the `status()` foreign-CLAUDE.md note string; `main_test.go`
adds `TestStatusNamesForeignClaudeMD` (the BUG-032 fix); `apply.go`'s hunk
changes the `keep` line's guidance string; `harness_test.go` adds one
assertion matching that string. `IsCheckout` and the rebuild-decision logic
in `cmd/sh-agent` and `internal/harness/repo.go` are untouched across the
whole range. §3a (plain-copy install) does not apply — no rebuild-decision
logic changed anywhere in c192622..8202d6a, only message strings and their
matching test assertions.

## Coverage of changed code (BUG-032 fix verification)

Two production strings were reworded in c514936, in two different files, and
each now has its own assertion — the per-file coverage map BUG-032's
prevention rule calls for:

- `internal/harness/apply.go:37` — the `keep` line's guidance string
  ("move its content into a rule file under rules/ in the harness repo,
  remove the file, and run install again") — covered by
  `internal/harness/harness_test.go`'s `TestInstallStatusResyncUninstall`,
  which asserts the substring `"move its content into a rule file under
  rules/"` on the install output (added in this fix commit, alongside the
  pre-existing `"keep    ~/.claude/CLAUDE.md"` assertion for the same test).
- `cmd/sh-agent/main.go:309` — the `status()` foreign-CLAUDE.md note string
  ("note       ~/.claude/CLAUDE.md is yours; move its content into a rule
  file under rules/ in the harness repo, remove it, and run install to write
  the one-line import") — now covered by
  `cmd/sh-agent/main_test.go`'s new `TestStatusNamesForeignClaudeMD`, which
  writes a hand-authored `~/.claude/CLAUDE.md`, runs `install` then `status`
  against a temp home, and asserts the string via `strings.Contains`.

Read `main.go:309` directly and diffed it character-for-character against the
test's `want` string: the two match exactly (test's `want` is the full line
text of `main.go:309` as printed). The two assertions are distinct strings in
distinct files with distinct tests — no test stands for two strings, closing
the BUG-032 gap.

**Can the new test fail?** Read its body: it asserts
`strings.Contains(out.String(), want)` after running `install` then `status`
against a temp home pre-seeded with a hand-written `CLAUDE.md`, so it only
passes if `status()` actually prints that exact note text when
`report.Counts["foreign"] > 0`. The BUG-032 record itself states the
fix commit's own evidence: the test passed as written, failed (exit 1) with
the note's string swapped out in `main.go`, and passed again after `main.go`
was restored — this is recorded as the main agent's / fix-commit's evidence
in `.agents/docs/bugs/BUG-032-status-note-uncovered-string.md`, not
independently reproduced by me in this dispatch (the instructions direct me
to record it as such rather than re-run it outside the tracked tree). The
test's structure (temp home, real `install`+`status` round trip, exact-string
`Contains` check with no wildcard) supports that it is a genuine, failable
assertion rather than a tautology.

No other changed `.go` file in the whole range lacks coverage: `main.go`'s
comment-only hunk is prose, not code, and needs no test.

## Inherited failures / added suppressions

- Inherited failures: none. `go vet` and `gofmt -l` both clean; `go test
  -race -v ./...` is 32 pass / 0 fail / 1 skip, matching the pre-existing
  skip count from earlier logs (unchanged, so not a new inherited failure).
- Added suppressions: scanned the whole `c192622...8202d6a` diff for added
  `//nolint`, `//lint:ignore`, `@ts-ignore`, `@ts-expect-error`,
  `eslint-disable`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`,
  `.only(`. All hits are prose in markdown (skill/agent docs, review
  findings, plan text) describing the suppression rules themselves — none is
  an actual suppression or skip directive added to source or test code. No
  FAIL line here.

## Cited: e2e / matrix log

`.agents/docs/testing/2026-09-15-verify-harness-practices-e2e.md` — read its
Status and Not examined sections in full:

Status: PASS. "Every row passes on both CLIs, directly or on its one re-run,
except four Antigravity rows that D50 records as caused outside this plan.
Every Phase 7 row at fcf3b79 passes on both CLIs."

Not examined (reproduced): `/review --fix` end to end (no row passes
`--fix`); a release shipped through a remote (no remote in the demos, so
preflight's Step 6 tag/push did not run); the operator's real home (6.2–6.4
ran in the sandbox home only; `testing/2026-09-15-verify-real-home.md` covers
the real home); the Mac; Antigravity's interactive prompts (every run
headless); Claude Code's `documentation` row at fcf3b79 (covered by proxy
rows); a long list of I-edits each with its own reason (harness-repo-only
files not installed, no matrix row exercises them, `--fix`-only edits,
preflight rows that cannot ship with no remote, etc — full list in that log).

The real-home tier (§5) is the main agent's own separate log for this date;
not run or read in full by me — named here as not run by this dispatch.

## Environment

- OS: Linux 6.12.107+deb13-amd64 (per `uname`, matches e2e log's
  environment).
- Go: as installed in this shell; `go test -race`, `go vet` both ran clean.
- No Docker/Colima work — no Dockerfile changed in this range.
- Repo state at dispatch start: clean, HEAD 8202d6a, branch main.

## Not examined

- §4 (sandbox CLI starts) and §5 (real-home smoke checks) and the full
  skill-matrix — explicitly out of scope per dispatch; cited above from
  `2026-09-15-verify-harness-practices-e2e.md` (its Not examined section
  reproduced above) and the real-home log (the main agent's, not read in
  full here).
- The BUG-032 record's own repro (test fails when the string is reverted) is
  taken as the fix commit's/main agent's evidence per the dispatch
  instructions, and was not independently reproduced by me inside
  `.agents/scratch/`.
- Full task-by-task cross-check of
  `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` against the
  diff — read the plan's existence, the Follow-ups entry on BUG-032, and the
  guide/e2e logs, but did not walk every individual task box against the
  diff line by line.
