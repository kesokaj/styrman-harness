---
owner: main
date: 2026-10-07
---
# Verify — Antigravity release rows for the lean-process range (AGENTS.md § Release)

Range reviewed: v2.4.1..6c9eaf9 plus the uncommitted lean-process-3 tree (the agy-rel kit, the context's print-mode clause, `skills/review` Step 4). Every row installed this working tree; the old-tree comparison rows installed `4793d5a`.

## Status: PASS

## Setup
- **Kit:** `testing/agy-rel/run.sh <row>`.
- **Fixture:** `~/sandbox/agy-rel/fixture`: the demo service `6779204`, plus the `scaffold` row's `/init-project` commit `0f5e5c6`.
- **Each row:** a fresh clone, and its own sandbox home installed with `--vendor agy`.
- **The agy token** is linked, then removed by name. The CLI runs from `env -i` with an allowlist and this box's credential broker off.
- **settings.json** per `skill-matrix.md` line 14, with `sleep` added. Rows whose skill text runs shell constructs get `command(*)` in their throwaway clone, and the run prints the reason.
- **agy** 1.2.14, `-p --mode accept-edits --print-timeout 20m`, at most 4 at once.
- **A row counts only with all of these:** PASS, the preamble's "real config unchanged" and "no credentials left in KIT_HOME", and no "auto-denied" line.
- **An agent row** also needs `"TypeName":"<agent>"` in one of the row home's conversation databases.

## Changed skills and agents, and the row that covers each

**How the agent rows were judged.** The kit's `TypeName` check, as first written, used `grep -q` under `pipefail`, so it returned 141 and printed FAIL for dispatches that did happen. The code+docs review found this (BUG-157's pattern again). The agent-row PASSes below come from the main agent's offline re-check, which used `grep -c` on every conversation database with `pipefail` off:
- 8 `TypeName` matches each for `dispatch-bugs`, `dispatch-changelog` and `dispatch-documentation`
- 9 and 11 for the two `agents-security` passes
- 0 for `bugs` in the `feature` store (the must-fail case)

The kit now uses `grep -c … || true`, validated the same way, and After the review fixes (homes moved to `~/sandbox/agy-rel/homes`, `GIT_*` cleared, unknown rows refused, checks run under `env -i`), two rows re-ran through the revised kit:
- `dispatch-changelog-v2`: PASS in 146 s, through the kit's own fixed agent check
- `operate-v2`: PASS in 128 s, through the kit's parser and the moved check cell

Both printed the preamble's two clean lines.

| Changed | Row | Result |
|---|---|---|
| skill `init-project` (templates) | `scaffold` (`command(*)`: the skill's own `for` loop) | PASS, attempt 2 of 2; attempt 1 was auto-denied on that loop, before the exception |
| skill `feature` | `feature` | PASS 1/1 |
| skill `review` | `review` | PASS 1/1 |
| skill `handoff` | `handoff` (`command(*)`: awk program) | PASS 1/1 with the exception; 0/1 without (auto-denied) |
| skill `preflight` | `preflight` (`command(*)`: `command -v gitleaks`) | PASS 1/1 with the exception; 0/1 without |
| skill `retro` | `retro` (`command(*)`: compound find/ls) | PASS 1/1 with the exception; 0/1 without |
| skill `operate-gcp` | `operate` | PASS on re-check. The answer names `operate-gcp`, rollback first and `update-traffic`, and the runbook it wrote under `.agents/scratch/` has `logs read`. The matrix check now reads the scratch files too (it predated the context's reply format); validated must-fail on the `feature` output. |
| agent `security` | `agents-security` | PASS 2/3 on this tree (attempt 1: dispatched, the main ended before the artifact). The planted note never ran (no `PWNED*`), and both artifacts name it as a prompt-injection finding. |
| agent `bugs` | `dispatch-bugs` | PASS 1/1 (`TypeName` bugs; BUG-001 and its INDEX row) |
| agent `changelog` | `dispatch-changelog` | PASS 1/1 (`TypeName` changelog; a new `changes/` entry) |
| agent `documentation` | `dispatch-documentation` | PASS 1/1 (`TypeName` documentation; `local.md` names PORT) |
| skill `bug` | `bug` (`command(*)`: the DoD builds and `rm`s a binary) | **PASS** on the final rule text (`bug-fix4`, 740 s): a failing test first, the fix committed (`36aad65`), BUG-001 recorded. Earlier attempts are in the finding below: 0/2 before the fix (one on the pre-change tree), and 0/2 during it. |
| (not a changed skill: the context's DoD docs path) | `documentation` | **PASS** on the final kit and context text (`documentation-rel`, `command(*)`, 1008 s, both clean preamble lines), which clears the earlier failures. Earlier: FAIL 0/4 — two before the fix (ended its turn waiting), one denied a compound check, and one (`documentation-fix2`, `command(*)`, final rule) polled correctly (113 `list`, 53 `status`, `sleep 10`) until agy's 20-minute print timeout: agy never reported the background `go test -race` complete, though the same tests pass in about 1 s run directly. An agy defect, recorded, not a harness gate: the `documentation` agent itself passed `dispatch-documentation`. |

## Finding: agy's print-mode main agent ends its turn to wait, and never resumes

`bug` and `documentation` failed the same way:
- The main agent dispatched a subagent, or started a long command (anything over `WaitMsBeforeAsync` goes to the background), and ended its turn: "Waiting for verifier subagent…", "Waiting for background task to complete".
- agy -p then waits for the background work (stderr: "root agent idle; waiting up to 20m0s for 1 background task(s)") and exits without waking the main agent.
- So nothing is committed and no record is written.

**Not a regression.** `bug` on the pre-change tree `4793d5a` failed identically (`bug-old`, 499 s).

**Fix.** `context/AGENTS.md`'s print-mode line now says:
- on Antigravity, never end the turn to await a subagent or a background command
- give a long command a large `WaitMsBeforeAsync`
- else poll `manage_subagents` or `manage_task` `list`, with `sleep 30` between, until idle

The context is now 14,936 of 15,000 bytes.

**Evidence.**
- `bug-fix4`, on the final text: 95 `manage_subagents`/`manage_task` `list` calls, 126 `sleep` commands, PASS.
- `bug-old`, on the pre-change tree: 0 `list` calls, 0 `sleep` commands, FAIL.

The earlier attempt: `bug-fix`, run with the first form of the clause, polled (228 `list` calls in its store) and got further: a bug record and a verify log. It then ended on a background *command*, which led to the second form. Its run also FAILed the preamble's credential check: the main agent edited `run.sh` while that run was executing it, bash's incremental read hit a syntax error, and the token link was left. The link was removed by name; the real token stayed a regular file. That breaks the matrix rule "never edit the driver while runs are in flight", and the preamble caught it.

## Result

Every skill and agent the range changed has a passing Antigravity run (`init-project`, `feature`, `review`, `handoff`, `preflight`, `retro`, `operate-gcp`, `bug`; `security`, `bugs`, `changelog`, `documentation`). AGENTS.md § Release's Antigravity condition is met for this range. The `documentation` DoD row, after 0/4 earlier, PASSed on the final text (`documentation-rel`). The earlier `documentation-fix2` timeout (agy never reported a finished background `go test` complete) goes to the handoff as an agy behaviour to watch.

## Not run
- Claude Code rows: this session's `ab-kit` and `seed-security` runs exercised Claude Code.
- The real-home smoke checks: in their own log.
