---
owner: main
date: 2026-10-07
id: BUG-163
status: fixed
area: testing (committed run kits)
---
# BUG-163: an A/B kit run left an agent-written server listening on every interface for four hours

## Symptom
The v2.5.0 release security review (`reviews/2026-10-07-security-v2.5.0-release.md`, finding 1, medium) found pid 269650 running on the operator's box:
- the `crm` binary, built by the agent in the `ab-kit` `control-2` run
- listening on `*:18089`
- its parent PID 1, up 4 h 06 min
- its executable `/tmp/crmsmoke/crm` (deleted), its cwd `~/sandbox/ab/runs/control-2/app`

## Repro
In a kit run, an agent starts a server in the background to smoke-test its build and ends without stopping it. The kit's `claude -p` exits, and the server is reparented to PID 1.

## Root cause
The kits checked credentials and the real config after a run, never processes. Both 2026-10-07 security reviews named "a server survives a timeout" as a residual, and it was recorded only as a header note. No check ran, so nobody saw the process.

## Fix
- `readlink /proc/269650/cwd` confirmed the clone, then `kill 269650`. The port is closed, and no other PID-1 orphan from the runs remains.
- `kit_leftovers` was added to `testing/ab-kit/run.sh`, `testing/seed-security/run.sh` and `testing/agy-rel/run.sh`. After the CLI exits, any process whose cwd is inside the run's clone is killed and FAILs the run (Linux `/proc`).
- Validated both ways on a planted `sleep 300` with its cwd in a probe clone: it returned 1 and killed the process; with nothing planted it returned 0.

**Limit closed (2026-10-08, plan `2026-10-08-kit-leftovers-and-hard-seeds`):** the check moved into `kit-preamble.sh`'s exit trap (`kit_process_leftovers`), so every kit gets it. It matches by physical cwd or executable under the clone or `KIT_HOME`, or by `HOME=$KIT_HOME` in the process environment, which catches a server started from `/tmp` too. It skips the trap's own shell, its ancestors and its descendants. Pinned by `TestKitPreambleKillsLeftoverProcesses`, with mutations seen failing. The three kits' pasted copies are removed.

**Still not caught**, stated so no one reads the check as complete (code review 2026-10-08):
- an orphan that changed directory AND dropped `HOME` (`env -i`)
- a container started through the Docker socket
- a process that ignores SIGTERM: it is reported and fails the run, but the 1 s wait does not see it die
- the kit shell's own descendants, skipped by design, so a kit killed mid-run leaves its CLI tree unflagged
- an operator's terminal or editor sitting in the clone or `KIT_HOME`, which WOULD be sent SIGTERM; nobody should work inside a kit's clone while it runs

## Prevention
A kit that runs an agent checks, after the run, for processes it left behind (cwd or executable under the clone or `KIT_HOME`, or `HOME=$KIT_HOME`), kills them, and fails the run; since 2026-10-08 the preamble's exit trap does it for every kit. A residual a review names as "survives" gets a check, not a header note.

## Related
BUG-160 (credentials left in sandbox homes: the same class, "what the run left behind"). Plan `2026-10-07-lean-process-3`. Security reviews `2026-10-07-security-worktree-ab-kit.md` (R2-3) and `2026-10-07-security-worktree-agy-seed-kits.md` (R2 finding 1).
