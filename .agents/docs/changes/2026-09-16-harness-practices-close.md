---
owner: main
date: 2026-09-16
---
# Harness practices close: the missing status test, the real-home install, the downstream mirrors and the Done list     (2026-09-16)

## Summary

These are the last steps of `plans/2026-09-15-harness-practices-g01-g26.md`, which is now closed.

- **The Done verifier's failure.** Its first run, at a33978d, failed on one line: `sh-agent status`'s note for a hand-written `~/.claude/CLAUDE.md`, reworded in c514936, had no covering test, although that commit said it came with one. 8202d6a adds `TestStatusNamesForeignClaudeMD` and BUG-032, and the re-run passes.
- **7.5.** The harness is installed into this box's real home, and its three smoke checks pass on both CLIs.
- **7.6.** Both downstream mirrors are refreshed from a81ef4b and committed in their own repositories, not pushed (D48, D51).
- **The close.** ROADMAP "Now" and `TODO.md` are empty.

## Type

task

## Files changed

- `cmd/sh-agent/main_test.go`: `TestStatusNamesForeignClaudeMD` (8202d6a).
- `.agents/docs/bugs/BUG-032-status-note-uncovered-string.md` and `INDEX.md`: BUG-032 (8202d6a).
- `.agents/docs/testing/2026-09-15-verify-harness-practices-done.md`: the Done verifier. The FAIL version from a33978d is kept in 8202d6a; the PASS version from 8202d6a is in a81ef4b.
- `.agents/docs/testing/2026-09-15-verify-real-home.md`: 7.5 (a81ef4b).
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`:
  - 7.5, 7.6, the Phase 7 Checkpoint and the Done list ticked;
  - D51;
  - Follow-ups for BUG-032's graduation, the mirrors' low security findings, and BUG-031's subagent form.
- `ROADMAP.md`: "Now" named 7.6 next in a81ef4b and is emptied at the close.
- `.agents/docs/bugs/` and `INDEX.md`: open records BUG-033 to BUG-037, from the two mirrors' security reviews.
- `.agents/docs/handoffs/2026-09-16-harness-practices-done.md`: the handoff.
- Outside this repository, committed there and not pushed: styrman-crew `image/harness/` (32 files, including `COMMIT` and `SHA256SUMS`) with its changes entry, verifier log and reviews (5418f19, 4957dfe); argolis-devbox `deploy/harness/` (39 files) with its changes entry, verifier log, reviews, bug records and testing guide (9380a63, 2d83282).

## Agents dispatched

- `verifier`: the Done row, twice.
- `bugs`: BUG-032, and BUG-033 to BUG-037 for the mirrors' security findings.
- `explorer`: the handoff's cold read.

Each mirror ran as a headless Claude Code session in its own repository with that repository's own agents; its changes entry names them.

## Verification

- **The new test.** `TestStatusNamesForeignClaudeMD` passed as written. It failed, exit 1, with `main.go`'s note string replaced, and passed again with `main.go` restored from git. Then:
  - `gofmt -l .` printed nothing;
  - `go vet ./...` exited 0;
  - `go test -race ./...` exited 0.
- **The Done verifier at 8202d6a: PASS.**
  - Counts: pass 32, fail 0, skip 1. The pass count rose from 31 with the added test.
  - `sh-agent list`: 27 skills and 7 agents. `render`: 93 lines.
  - Scratch install: drift none, and every `## Records` section matches.
  - § 3a does not apply.
- **The real home.**
  - Install at a33978d exited 0 with 45 files written; `status` showed drift none and `binary a33978d7ab9e`.
  - The three smoke checks pass on both CLIs, `Failed to parse` prints 0, and Antigravity's settings are restored identical to the backup.
  - A reinstall at `8202d6a+dirty` wrote only the binary. The final reinstall, after a81ef4b from the clean checkout, shows `binary a81ef4be3cb8` and drift none.
- **The mirrors.** The main agent checked each mirror commit in its own repository. `COMMIT` reads a81ef4be3cb8e4a933a386d9c87c5a0ccbd701a8. `diff -rq` against this clean checkout at a81ef4b prints exactly three lines: `Only in …: COMMIT`, `SHA256SUMS` and `VENDORED.md`. `sha256sum -c` exits 0 over 155 entries. `go vet ./...` and `go build ./cmd/sh-agent` exit 0. `sh-agent list` prints Skills (27) and Subagents (7). Both repositories are clean and two commits ahead of `origin/main`, with nothing pushed. argolis-devbox's own session could not run `diff -rq` under its sandbox (its BUG-008) and used an rsync dry-run instead; the main agent's run is the recipe's check.

## Docs

`README.md`, the testing guide and the deploy guides are unchanged: no command, flag or deploy step changed after the Phase 7 docs review. `deploy/remote.md` § Downstream mirrors already describes the refresh.

## Requirements

n/a
