---
owner: main
date: 2026-10-12
id: BUG-170
status: fixed
area: internal/dream (call, the CLI's signal handling)
---
# BUG-170: an ignored hangup or interrupt killed the dream CLI's group and left sh-agent running

## Symptom
Found by the security agent of the abandoned v2.6.2..66f3412 release gate (`reviews/2026-10-09-security-v2.6.2-66f3412.md`, finding 1, low), not observed in a run. `6cec23b` (unreleased) made `callWithin` watch SIGINT, SIGTERM, SIGHUP and SIGQUIT, kill the CLI's process group, then re-raise the signal so sh-agent dies of it. Under `nohup sh-agent dream --all &`, where SIGHUP arrives ignored, or as a background job of a non-interactive shell, where SIGINT does, a logout or an interrupt killed the project in flight. sh-agent itself did not die and went on to the next project. Before `6cec23b` the CLI shared sh-agent's ignored disposition and ran on.

## Repro
- `TestCallLeavesTheGroupOnAnIgnoredSignal` (`internal/dream/run_test.go`) starts the helper through `sh -c "trap '' HUP; exec …"` (and `INT`), which is what nohup does, then sends that signal. Before the fix both subtests failed: "the ignored hangup killed the CLI's group" and "sh-agent exited on the ignored hangup", and the same for interrupt.

## Root cause
`signal.Notify` takes over a signal that the process inherited as ignored. The Go runtime keeps an inherited SIG_IGN for SIGHUP and SIGINT only (`sigInstallGoHandler` in `runtime/signal_unix.go`), and Notify installs a handler over it. `signal.Reset` then restores the inherited ignore, so the re-raise in the goroutine did nothing: the group was already dead and sh-agent lived on.

## Fix
`internal/dream/run.go` `callWithin`: SIGTERM and SIGQUIT are always watched, and SIGINT and SIGHUP only when `signal.Ignored` is false at the call. `git grep 'signal\.Notify'` finds no other site in the tree.

## Prevention
A process that relays a signal checks `signal.Ignored` first: `signal.Notify` overrides an ignore the parent set on purpose (nohup, a background job), and `signal.Reset` brings the ignore back, so a re-raise after Reset cannot kill the process. A test for a signal path starts the child with the signal ignored (`trap '' SIG; exec`) as well as with it at default.
