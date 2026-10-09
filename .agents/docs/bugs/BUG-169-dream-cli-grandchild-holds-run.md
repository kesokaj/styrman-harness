---
owner: main
date: 2026-10-09
id: BUG-169
status: fixed
area: internal/dream (call, the CLI invocation)
---
# BUG-169: a grandchild of the dream CLI could hold the run, and the lock, past the deadline

## Symptom
Found by the project review of 2026-10-09 (code lens, medium), not observed in a run. `call` bounded the CLI with `exec.CommandContext` and a 10-minute deadline, but the kill reached only the direct child, and stdout and stderr were buffers, so `cmd.Run` waited for the pipe copies to see EOF. A helper the CLI started (a language or MCP server) that inherited the pipes kept them open, `Run` blocked for as long as it lived, `dream.lock` stayed held, and every later nightly run printed "already running" and exited 0: no error anywhere. The spec row "times out → the project fails, the others run" (`dream.md:82`) did not hold.

## Repro
`TestCallReturnsAtTheDeadlineDespiteAGrandchild` (`internal/dream/run_test.go`): a fake CLI that backgrounds `sleep 60` (and one that backgrounds it under `setsid`) and then sleeps, with a one-second deadline. Both subtests failed before the fix: "call did not return" after the test's 20-second limit.

## Root cause
`os/exec` documents both halves: `CommandContext` kills the process it started, not its descendants, and `Wait` waits for the I/O copying goroutines unless `WaitDelay` is set. `call` set neither a process group nor `WaitDelay`.

## Fix
`internal/dream/run.go`: `call` is `callWithin(argv, prompt, work, env, timeout, grace)` with the two bounds explicit. The CLI starts in its own process group (`Setpgid`), `cmd.Cancel` kills the whole group at the deadline, and `cmd.WaitDelay` (`killGrace`, 30 s) cuts loose a grandchild that left the group. The in-group case returns at the deadline; the escaped case after the grace. Spec: `dream.md` § Ask names both.

## Prevention
Every `exec` with a deadline whose output goes through pipes sets `WaitDelay`, and one that may start helpers kills its process group; its test runs a fake command that backgrounds a pipe holder and asserts the call returns. The dream call is the only deadline-bounded exec in the tree today (`grep -rn CommandContext cmd internal`).
