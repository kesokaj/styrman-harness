---
owner: main
date: 2026-09-30
---
# Behavioural runs — the verifier and a suite that takes minutes   (2026-09-30)

Harness proposal "The verifier never hands back with a background job alive" (`~/.config/sh-agent/proposals/2026-09-23-retro.md` Proposal 2, queued since 2026-09-23 with the note "needs a behavioural run designed for it"). `agents/verifier.md`'s "Long suites, once" bullet now says: a suite that takes minutes runs once in the foreground with the call's timeout raised to the tool's maximum and its output redirected to a file (never a pipe, BUG-086); only a suite that cannot finish inside that maximum runs in the background; every background job is ended (`kill`, then `kill -0` fails) before the hand-back.

## Status: PASS on Claude Code — new definition 0 of 3 background calls, control 2 of 3

## Kit

`.agents/scratch/vkit/slowfix/` is `greetbox` plus `internal/slow`, whose one test sleeps 150 s, so `go test -race ./...` takes about three minutes — past the Bash tool's default two-minute timeout and inside its ten-minute maximum; its `AGENTS.md` says so. `run-slow.sh <new|control> <n>` clones it to `~/sandbox/vkit/runs/slow-<arm>-<n>` (outside every checkout), proves the fixture facts, and runs `claude -p --agents <json>` with the default setting sources (the real home's context loads), the `new` arm's JSON built from the working-tree `agents/verifier.md`, the `control` arm's from `git show HEAD:agents/verifier.md`. The dispatch asks for the project's commands, the early log and the git status at the end; it says the suite takes about three minutes and never states the changed rule. After the CLI exits the script counts processes still holding the clone's path (`pgrep -fa "go test|slow.test|go-build"`), then kills any.

Attribution from the subagent transcripts (`~/.claude/projects/-home-kesokaj-sandbox-vkit-runs-slow-<arm>-<n>/*/subagents/`): the new sentence "in the foreground, with the call" is present in 3 of 3 `new` transcripts and in 0 of 3 `control` transcripts. A suite call is a Bash tool call whose command holds `go test` without `-run`; background means `run_in_background: true` or a trailing `&`.

## Results

| Arm | Run | Background suite calls | Foreground suite calls (timeout) | Processes left after exit | Log Status |
|---|---|---|---|---|---|
| new | 1 | 0 | 3 (600000 ms each) | 0 | FAIL (untested packages, by design) |
| new | 2 | 0 | 3 (600000) | 0 | FAIL |
| new | 3 | 0 | 3 (600000) | 0 | FAIL |
| control | 1 | 1 (`run_in_background: true`) | 1 (600000) | 0 | FAIL |
| control | 2 | 0 | 2 (600000) | 0 | FAIL |
| control | 3 | 2 | 1 (600000) | 0 | FAIL |

The new definition never sent the suite to the background (0 of 3) and always raised the call's timeout to the maximum; the old one sent it to the background in 2 of 3 runs, then re-ran it in the foreground when the poll could not read the result. Every log's Status is FAIL for the fixture's untested packages, as intended; every log reports the suite's `ok`.

Why three foreground calls per new run: the first two forms (`mkdir -p … && go test … > file 2>&1`, then `go test … > file 2>&1`) were refused by the run's permission gate (`--allowedTools "Bash(go:*)"` matches neither a compound nor a redirected command in headless mode) and the third, bare `go test -race -v -count=1 ./...`, ran; the rule's "once" held as far as the gate allowed. A run from the real home with its own permissions has no such refusal (the v1.14.0 real-home runs).

## What was not exercised

- A suite longer than the ten-minute maximum (the rule's remaining background branch) and the `kill` / `kill -0` ending it prescribes: no fixture; the branch is text only.
- A live child at hand-back on this kit: headless `claude -p` reaps its children at exit, so "processes left" is 0 in both arms and does not measure the 2026-09-22 symptom (an interactive session's agent staying open); the measured variable is whether a background call is issued at all.
- Antigravity: no run here; the release gate owes one run of the changed agent there.
