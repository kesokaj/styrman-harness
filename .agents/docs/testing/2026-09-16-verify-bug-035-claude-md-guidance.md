---
owner: verifier
date: 2026-09-16
---
# Verification — bug-035-claude-md-guidance (2026-09-16)

Re-run after a docs review round. The two Go messages and their assertions changed
by two words: they now say "never credentials or personal or host-specific notes",
matching `README.md:53`. This log supersedes the previous run of the same slug,
which quoted the shorter wording.

Range reviewed: uncommitted working tree — `git diff HEAD -- README.md cmd/sh-agent/main.go cmd/sh-agent/main_test.go internal/harness/apply.go internal/harness/harness_test.go` (BUG-035)

## Status: PASS

Scope: §1 Go checks, §2 harness self-check, §3 scratch-home install + drift on two
homes (plain and foreign-CLAUDE.md), independent reproduction of the must-fail
assertion against HEAD's old strings. §3a does not apply. §4 (sandbox CLI starts)
and §5 (real-home smoke) not re-run — cited as prior art, not exercised this run.
Cites `.agents/docs/testing/sh-agent/guide.md`.

## Commands

- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0.
- `go test -race -v ./...` → exit 0. Output saved to
  `.agents/scratch/verify-bug-035-2/go-test-1.txt`.
  counts: pass 54, fail 0, skip 1 at 703c6cd+dirty for `go test -race -v ./...`
  (38 top-level `--- PASS:` + 16 indented subtest `--- PASS:` lines, all 4
  subtests of `TestAdoptRefusesDestructiveSources`; 0 `--- FAIL:`; 1
  `--- SKIP:` — `TestFindRepo`).
  vs previous: `.agents/docs/testing/2026-09-16-verify-bug-035-claude-md-guidance.md`
  (this file's own prior version, pass 54, fail 0, skip 1) → unchanged, pass 54→54,
  skip 1→1. No test was added in this round's diff (only string literals changed
  inside the two existing assertions), so an unchanged count is expected.
- `go run ./cmd/sh-agent list` → exit 0; Skills (27), Subagents (7), both sections
  print full tables with no parse errors.
- `go run ./cmd/sh-agent render | wc -l` → `93` (expected 93 per §2).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug-035-home-2`
  (fresh home, no pre-existing `~/.claude/CLAUDE.md`) → `260 written, 0 unchanged,
  0 removed, 0 kept — 260 files managed`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug-035-home-2` →
  `drift      none — 259 files match the repo` (the count line's own file is
  excluded from its own comparison, as in the previous run; not a regression).
- Foreign-home tier: planted `.agents/scratch/verify-bug-035-home-foreign-2/.claude/CLAUDE.md`
  with the content `my own personal notes` before install, then
  `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug-035-home-foreign-2`
  → printed the keep line, quoted verbatim:
  > keep    ~/.claude/CLAUDE.md (not written by sh-agent; move only rules meant for
  > every machine into a rule file under rules/ in the harness repo — rules/ ships
  > to every install and mirror, never credentials or personal or host-specific
  > notes — then remove the file and run install again)
  `259 written, 0 unchanged, 0 removed, 1 kept — 259 files managed`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug-035-home-foreign-2`
  → `foreign  ~/.claude/CLAUDE.md` and the note line, quoted verbatim:
  > note       ~/.claude/CLAUDE.md is yours; move only rules meant for every machine
  > into a rule file under rules/ in the harness repo — rules/ ships to every
  > install and mirror, never credentials or personal or host-specific notes — then
  > remove it and run install to write the one-line import
  `drift      none — 258 files match the repo`.
  Both the keep line and the status note match `README.md:53`'s final wording
  ("never credentials or personal or host-specific notes") word for word.

## §3a Plain-copy install

Does not apply: this change touches only message strings in `cmd/sh-agent/main.go`
and `internal/harness/apply.go`, not `IsCheckout` / rebuild-decision logic in
`cmd/sh-agent` or `internal/harness/repo.go`. No plain-copy/mirror behavior is
affected.

## §4 / §5 not re-run

Sandbox CLI starts (§4, `.agents/docs/testing/2026-09-03-verify-sandbox.md`) and
real-home smoke checks (§5, `.agents/docs/testing/2026-09-03-verify-real-home.md`,
`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`) were not started this
run, per the dispatch's instruction not to start either CLI or run the matrix.
Cited as prior art from earlier verification rounds in this repo's history, not
exercised here.

## Independent must-fail reproduction

Read the pre-change source with `git show HEAD:<path>` (no working-tree mutation)
into `.agents/scratch/verify-bug-035-2/oldcheck/old-main.go.txt` and
`old-apply.go.txt`, then grepped for the new strings:

```
grep -n "move its content" old-main.go.txt old-apply.go.txt
  → old-main.go.txt:309: ...move its content into a rule file under rules/...
  → old-apply.go.txt:37: ...move its content into a rule file under rules/...
grep -n "never credentials" old-main.go.txt old-apply.go.txt
  → (no matches)
```

`HEAD`'s messages contain neither "move only rules meant for every machine" nor
"never credentials or personal or host-specific notes" — the strings the new
assertions in `cmd/sh-agent/main_test.go` (`TestStatusNamesForeignClaudeMD`) and
`internal/harness/harness_test.go` (`TestInstallStatusResyncUninstall`) require.
The live install/status runs above show the working-tree code prints the full new
wording, and it is not present in `HEAD`'s source, so the two updated tests would
fail against the unfixed code. This was checked by source diff, not by re-running
a build of `HEAD` (a rebuild of HEAD-only code was unnecessary given the exact
string is provably absent from the source both tests grep for).

## Root cause

N/A — Status is PASS, not a bug-hunt log.

## Coverage of changed code

- `cmd/sh-agent/main.go` (`status`'s foreign-CLAUDE.md note) — covered by
  `cmd/sh-agent/main_test.go: TestStatusNamesForeignClaudeMD`, one assertion
  checking the full updated string including "never credentials or personal or
  host-specific notes".
- `internal/harness/apply.go` (`Install`'s keep-line message) — covered by
  `internal/harness/harness_test.go: TestInstallStatusResyncUninstall`, one
  assertion (two `Contains` checks combined with `||`) covering both the
  "move only rules meant for every machine" phrase and the "never credentials or
  personal or host-specific notes" exclusion.
- `README.md:53` — prose only, no test; its wording is the source of truth the two
  Go messages were changed to match, confirmed identical by direct comparison
  above.
- No added `//nolint`, `//lint:ignore`, `t.Skip`, `.only(`, `.skip(`, or similar
  suppressions in this diff.

## Not examined

- §4 sandbox CLI starts (Claude Code, Antigravity) — not started this run, per
  dispatch instruction; last known-good runs are the dated logs cited above, not
  re-verified against this exact wording.
- §5 real-home smoke checks — not run this run; last known-good run predates this
  change.
- Full skill-matrix (`.agents/docs/testing/skill-matrix.md`) — not run; this change
  does not touch `context/`, `agents/`, or `skills/` content in a way that affects
  matrix scenarios, only two Go message strings and their tests plus README prose.
- `.agents/docs/reviews/2026-09-16-review-bug-035-code.md` and
  `-docs.md` — present in the working tree as untracked files but not opened or
  evaluated by this verification run (review is a separate agent's scope).

## Environment

Linux 6.12.107+deb13-amd64, repo at `703c6cd` + uncommitted changes (the BUG-035
diff under verification). `go` and `git` on PATH; no Docker/Colima involved (no
Dockerfile in this change).
