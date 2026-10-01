---
owner: bugs
date: 2026-10-01
id: BUG-140
status: graduated 2026-10-01 to testing/guide.md § 4 § Behavioural runs (a kit separates the run from the check); cause of the lost lines unknown
area: testing (skill-matrix run kits)
---
# BUG-140: three of five matrix kit processes ended without their result line or EXIT-trap listing

## Symptom
Phase 7 skill-matrix batch (`.agents/scratch/audit-fix/kit-phase7.sh`, five rows under `xargs -P 3`, 2026-10-01 ~09:00): `spec`, `requirements` and `explain-project` wrote no line to `~/sandbox/audit-fix/p7-results.txt` and no kit-preamble EXIT-trap listing to the driver log. Their Claude runs completed: each `.out` holds the final answer, each clone the expected artifact; transcripts show `<command-name>/spec` etc., so the skills loaded. `testing-strategy` and `init-deploy` wrote their lines. Checks re-run by hand in the clones all PASS. Phase 5 (6 rows) and Phase 6 (11 rows), same shape, lost no line.

## Repro
Not reproduced on demand. The three affected rows were the slash-command rows whose runs themselves ran `/handoff` and committed in the clone (requirements, explain-project) or wrote a bug record and handoff (spec).

## Root cause
Unknown. The kit (`set -euo pipefail`, `set +e` after the CLI call, check, `tee -a` of one line, EXIT trap from `testing/kit-preamble.sh`) reached neither the echo nor the trap, which points at the process being killed rather than exiting. Structural cause of the lost evidence: the result line is produced by the same process that ran the CLI, so losing that process loses the check.

## Fix
Not fixed. Mitigation: check cells re-run by hand in the surviving clones, recorded in `testing/2026-10-01-verify-phase7-skill-matrix.md`.

## Prevention
1. A run kit separates RUN from CHECK: the check is a second script (or a `--check-only` mode) that re-reads the clone and transcript and can be re-run after the driver finishes; a result line is derived from state on disk, not from the process that produced it.
2. `grep ... | wc -l || echo 0` under `pipefail` prints two zeros (grep's exit 1 fails the pipeline, so the `||` also fires; the kit's `cloud_cmds` field shows it). Use `grep -c ... || true` for a count.

## Graduated
2026-10-01: the prevention rule is in `.agents/docs/testing/guide.md` § 4 § Behavioural runs (a kit separates the RUN from the CHECK; the check re-reads the saved output and the clone and can be run again; counts are `grep -c … || true`), written 2026-10-01. The cause of the three lost result lines is still unknown (status stays open). The rule was used twice on 2026-10-01 to re-judge a row from its saved record when a kit's verdict script failed.

## Related
BUG-107, BUG-123 (kit-preamble); `testing/skill-matrix.md` § Procedure step 3. Cause open: next time, capture the exit status of each `xargs` child and `dmesg`/OOM evidence before re-running.
