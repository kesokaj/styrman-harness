---
owner: verifier
date: 2026-09-15
---
# Verification — harness-practices-done (2026-09-15)

Range reviewed: c192622..a33978d, which includes the Phase 7 range 641eccd..a33978d
(the Phase 7 checkpoint) as a subrange. This log is both the Definition of Done's
verifier row and the Phase 7 checkpoint's verifier log.

Procedure: `.agents/docs/testing/sh-agent/guide.md`, §1–§3.

## Status: FAIL

One line of changed code has no covering test (see Coverage below); everything
else — Go checks, harness self-check, scratch install/status, agent Records
sections, test counts vs. previous, and the cited e2e/matrix log — passes.

Scope: §1 (Go checks), §2 (harness self-check), §3 (scratch-home install/status +
agent Records parity) all ran to completion. §3a does not apply (checked against
the diff, see below). §4/§5 (sandbox CLI starts, real-home smoke) explicitly out
of scope for this dispatch — cited from other logs, not run here.

## Commands

- `gofmt -l .` → exit 0, no output (printed nothing → pass).
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0. Output saved to
  `.agents/scratch/verify-harness-practices-done/go-test-1.txt`.
  counts: pass 31, fail 0, skip 1 at a33978d for `go test -race -v ./...`
  vs previous: `.agents/docs/testing/2026-09-15-verify-harness-practices-phase5.md`,
  pass 31→31, skip 1→1 (no change; b43dde1+dirty there vs a33978d here).
- `go run ./cmd/sh-agent list` → exit 0. 27 skills, 7 agents listed (table
  headers "Skills (27)", "Subagents (7)"), matching plan expectation.
- `go run ./cmd/sh-agent render | wc -l` → 93 (< 100, matches plan expectation).
- `rm -rf .agents/scratch/verify-harness-practices-done-home && mkdir -p ... &&
  go run ./cmd/sh-agent install --home .agents/scratch/verify-harness-practices-done-home`
  → "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed", exit 0.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-harness-practices-done-home`
  → "drift      none — 259 files match the repo", exit 0. (259 vs 260 written is
  the binary line reported separately as `binary a33978d7ab9e`, not drift.)
- Agent `## Records` parity: for each of the 7 `agents/*.md` files, grep-counted
  `^## Records` in the source and in the derived
  `.gemini/config/agents/<name>/agent.md` under the scratch home. All 7 pairs
  show source=1, derived=1 (bugs, changelog, documentation, explorer, reviewer,
  security, verifier).

## §3a applicability

Checked: `git diff c192622..a33978d --stat -- '*.go'` shows exactly three
changed Go files: `cmd/sh-agent/main.go` (+3/-3), `internal/harness/apply.go`
(+1/-1), `internal/harness/harness_test.go` (+3/-0). Read the full diff: the
`main.go` hunk changes a comment sentence ("its image built the running
binary" → "its image or VM build step built the running binary") and the
`status()` foreign-CLAUDE.md note string; the `apply.go` hunk changes the
`keep` line's guidance string; `harness_test.go` adds one assertion matching
the new string. `IsCheckout` and the rebuild-decision logic in `cmd/sh-agent`
and `internal/harness/repo.go` are untouched in this range. §3a (plain-copy
install) therefore does not apply — no rebuild-decision logic changed, only
message strings and their matching test assertion.

## Coverage of changed code

- `internal/harness/apply.go` — covered: `TestInstallStatusResyncUninstall`
  in `internal/harness/harness_test.go` asserts the new keep-line string
  ("move its content into a rule file under rules/").
- `internal/harness/harness_test.go` — is itself the test; not applicable.
- `cmd/sh-agent/main.go`:
  - The comment-only edit (line ~124) is not code and needs no test.
  - The `status()` foreign-CLAUDE.md `note` line (new wording "move its
    content into a rule file under rules/ in the harness repo, remove it,
    and run install to write the one-line import") is changed `fmt.Fprintln`
    output but **has no test asserting its text**. Grepped
    `cmd/sh-agent/main_test.go` for "note" / "foreign" / the new string:
    no match. The scratch-home install run above used a home with no
    pre-existing foreign `~/.claude/CLAUDE.md`, so `report.Counts["foreign"]`
    was 0 there too and the line never printed in this run either — it was
    not observed executing anywhere in this verification.
  - **FAIL line: `cmd/sh-agent/main.go`'s `status()` foreign-note string is
    NOT COVERED** by any automated test and was not exercised in this
    verification's scratch run (no foreign CLAUDE.md present).

## Inherited failures / added suppressions

None found: `go vet` clean, `gofmt -l` clean, no `//nolint` or `//lint:ignore`
additions in the diff (checked by inspection of the three changed .go files
above — no suppression comments in any hunk).

## Cited: e2e / matrix log

`.agents/docs/testing/2026-09-15-verify-harness-practices-e2e.md` —
Status: PASS. "Every row passes on both CLIs, directly or on its one re-run,
except four Antigravity rows that D50 records as caused outside this plan.
Every Phase 7 row at fcf3b79 passes on both CLIs." Its own Not examined
section (read in full): `/review --fix` end to end; a release that ships
through a remote (no remote in the demos); the operator's real home (covered
separately by `testing/2026-09-15-verify-real-home.md`); the Mac; Antigravity
interactive prompts (all headless); Claude Code's `documentation` row at
fcf3b79 (covered by proxy rows); a named list of I-edits with per-edit
reasons (no matrix row exercises them, harness-repo-only files, `--fix`-only
edits, etc. — see that log for the full list).

The real-home tier (§5) is the main agent's own log, written separately for
this same date; not run or read in full by me — named here as not run by
this dispatch.

## Environment

- OS: Linux 6.12.107+deb13-amd64 (per `uname`, matches e2e log's environment).
- Go: as installed in this shell (`go test -race` and `go vet` both ran clean).
- No Docker/Colima work — no Dockerfile changed in this range.
- Repo state: clean, HEAD a33978d, branch main (confirmed via `git status --short`
  at dispatch start).

## Not examined

- §4 (sandbox CLI starts) and §5 (real-home smoke checks) and the full
  skill-matrix — explicitly out of scope per dispatch; cited from
  `2026-09-15-verify-harness-practices-e2e.md` above (its own Not examined
  section is reproduced above) and the real-home log is the main agent's,
  not read in full here.
- Whether the uncovered `main.go` status-note string is otherwise checked
  by a sandbox/matrix row that supplies a real foreign `~/.claude/CLAUDE.md`
  in the real-home or e2e tiers — not checked; the e2e log's scope list does
  not name a `status`-with-foreign-file row, so this is flagged as an open
  gap rather than assumed covered elsewhere.
- Full contents of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`
  task-by-task — read only the plan's existence and the guide/e2e logs; did
  not cross-check every individual task box against the diff.
