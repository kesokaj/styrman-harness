---
owner: main
date: 2026-09-30
---
# Verification — real home on the Linux box before v1.15.0, both CLIs   (2026-09-30)

Range reviewed: v1.14.0..e69c56e, installed into the real home as e69c56e (`sh-agent update` from this checkout at 12:21: 3 written, 263 unchanged, 266 managed, binary e69c56e15042, `status` with no warning line; `agents/verifier.md` byte-identical between the tree and the installed copy). The checks run from the real `$HOME` in fresh clones of fixtures outside every checkout.

## Status: PASS on both CLIs — the three smoke checks, and one run of the changed `verifier` on each CLI with a suite that takes minutes

The rule: project `AGENTS.md` § Release. The only changed skill or agent in the range is `agents/verifier.md` (the long-suites rule); `AGENTS.md` at the root is this repository's own file, not deployed.

## Fixtures and kit

`~/sandbox/gate15/fixture` is the slow `greetbox` (`internal/slow`'s test sleeps 150 s, so `go test -race ./...` takes about three minutes; `AGENTS.md` says so); `~/sandbox/gate14/fixture` is the `greetbox` with the `filepath.Glob`-on-a-flag file for the skill-load check. `run.sh <cc|agy> <scenario> <n>` in each kit clones outside every checkout, proves the scenario's facts, pins every `sh-agent` call to a scratch home (BUG-107; `~/.config/sh-agent/` unchanged before and after), and for the verifier scenario drops a foreign artifact under `.agents/docs/reviews/` two seconds after the verifier's log appears, then counts processes still holding the clone after the CLI exits.

## Claude Code 2.1.284 (`claude -p`, real home)

| Check | Result |
|---|---|
| Context questions, no tools | PASS — "AGENTS.md — global agent context", 28 skills and 8 subagents, GANNET, `bug` |
| `golang` skill on the glob file (gate14 fixture, chat only) | PASS — `golang` loaded, the glob rule cited |
| `verifier` dispatch, slow suite, foreign artifact | PASS — the suite ran ONCE, in the foreground, `timeout: 600000`; log Status FAIL for the untested packages (intended) with the suite's `ok`; the foreign file "another agent's artifact that appeared while I ran … I list it by name. No source, test or config path changed"; 0 processes left after exit |

## Antigravity CLI 1.2.13 (`agy -p`, real home, `--mode accept-edits`, no skip flag)

No run's stderr has an `auto-denied` line.

| Check | Result |
|---|---|
| Context questions, no tools | PASS — the same four answers, GANNET |
| `golang` skill on the glob file (gate14 fixture) | PASS — `golang` named, `filepath.Glob` in the findings |
| `verifier` dispatch, slow suite, foreign artifact (the derived `agent.md`) | PASS — `go test -race -v -count=1 ./...` exit 0 recorded; the foreign file "(seen: peer subagent artifact)"; "No source, test, or config paths modified"; 0 processes left after exit |

## What was not exercised

- The Mac and the devbox: they pick the release up with `sh-agent update`; `darwin/arm64` and `linux/arm64` binaries were built by the gate verifier, not run.
- Interactive sessions on either CLI; every run was headless.
- A suite longer than ten minutes (the rule's background branch with `kill` / `kill -0`): no fixture.
- The Antigravity verifier's tool calls in detail (its conversation database was not opened); the log's text and the process count after exit are the evidence.
- A second attempt of any scenario: every scenario passed on its first run.
