---
owner: main
date: 2026-10-08
---
# The leftover-process check moves into kit-preamble.sh; a harder seeded security set; the remote.md broker note

## Type

improvement (test kits), with BUG-163 follow-up

## Summary

This closes the v2.5.0 handoff's follow-ups (plan `plans/2026-10-08-kit-leftovers-and-hard-seeds.md`, the operator's "go").

- **`kit-preamble.sh`:**
  - The exit trap runs `kit_process_leftovers` first, from `/`.
  - **What matches:** a same-uid process whose physical cwd or executable is under the clone or `KIT_HOME`, or whose environment holds exactly `HOME=$KIT_HOME` (frozen at load). Each match gets SIGTERM and the run FAILs.
  - **What is skipped:** the trap's shell, its ancestors and its descendants.
  - **What is printed:** pid, match kind and a sanitised executable name, never argv.
  - **New line:** `kit: no processes left by the run`. Without `/proc` it reads `no /proc here: …`.
  - **Load-time refusals (exit 2):** a relative `KIT_CLONE`, and any root resolving to `/`, the real home or a parent of it, or the harness checkout or a parent of it. `KIT_HOME` is resolved without being created.
- **`TestKitPreambleKillsLeftoverProcesses`:**
  - orphans planted by double fork: cwd in the clone, executable under the clone, `HOME=$KIT_HOME`; each FAILs and is terminated
  - near misses `HOME=$KIT_HOME-x` and the kit shell's own child are left alone
  - a negative control
  - seven load-time refusals, each with its own message
  - mutations seen failing: the matcher disabled, and the real-home refusal removed
- **The three kits** drop their pasted `kit_leftovers`; their headers point at the preamble and name the SIGKILL gap.
- **`seed-security/run.sh <label> <effort> hard`:** three multi-step, fail-open flaws in one commit, scored `hard N/3`. The six-flaw default is unchanged.
- **Docs:**
  - `guide.md` § Behavioural runs: three exit lines, the refusals, nobody working inside a kit's clone
  - `remote.md` names where this box's credential broker is described
  - BUG-163's limit paragraph is closed, with a "Still not caught" list, plus its Prevention line and INDEX row

## Verification

- The preamble tests and the BUG-155 pin pass (`go test ./internal/harness -run 'TestKitPreamble|TestScriptsHaveNoEarlyExitPipe'`).
- **The first full suite FAILed** on `TestScriptsHaveNoEarlyExitPipe`: the new code piped environ into `grep -q` (BUG-155's pattern), now read from the file. The verifier's final full run is below.
- **E2E through the new preamble:**
  - `ab-kit` `control`: 25/25, twice; the second run was on the post-review preamble before the last load-time change
  - `agy-rel` `operate`: PASS
  - all five seed runs
  - every one printed the three lines
- **Seeded `security`, hard set** (`testing/2026-10-08-verify-seed-security-hard.md`): `high` 3/3 in 3/3 at about 150 s and $0.54; `max` 3/3 in 1/1 at 590 s and $1.75. A hand read of one `high` artifact confirmed all three plus five unplanted findings.
- **A mutation run went wrong.** With the checkout refusal removed, the trap treated the checkout as `KIT_HOME`. It wrote `<checkout>/bin/sh-agent`, deleted by name; nothing tracked was touched. It began a credential scan of the checkout, which was stopped. Its process check may have sent SIGTERM to same-uid processes under the checkout; this session's own were spared as ancestors.
- The verifier's log, run last, is under `testing/`.

## Agents dispatched

| Agent | Times | What for |
|---|---|---|
| reviewer | 2 | the plan cold read; the code+docs review (2 rounds by resume) |
| security | 1 | 2 rounds by resume; round 1 had a medium |
| verifier | 1 | last |
