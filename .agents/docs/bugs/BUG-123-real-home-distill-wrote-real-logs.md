---
owner: bugs
date: 2026-09-30
id: BUG-123
status: fixed
area: testing (real-home gate kit)
---
# BUG-123: a real-home gate run of /distill wrote test lines into the real ~/.config/sh-agent logs

## Symptom
During the v1.17.0 gate (`testing/2026-09-30-verify-real-home-v1.17.0.md` § Kit defect), the Antigravity `/distill` run on a fixture (`~/sandbox/gate19/runs/agy-distill-1`) appended one line to the REAL `~/.config/sh-agent/feedback.jsonl` and one to `retros.jsonl`, both naming the fixture path. `report` would have sent them to the feedback repository.

## Repro
Run `/distill` under the real home on a fixture project, then list `~/.config/sh-agent` before and after. Found exactly that way: the kit's before/after listing.

## Root cause
The kit's BUG-107 wrapper pins only `sh-agent` calls to a scratch home. Skills that write `~/.config/sh-agent` with file tools (feedback-log, retro, distill) bypass the wrapper and write the real directory when the run uses the real home.

## Fix
The two lines were removed by exact match: feedback 16 -> 15, retros 17 -> 16, `grep -c sandbox/gate19` is 0 in both. The feedback-log gate run was moved to the agy-only sandbox home.

## Prevention
A gate run of a skill that writes `~/.config/sh-agent` (feedback-log, retro without `--dry-run`, distill) runs in a sandbox home, never the real one. The kit lists `~/.config/sh-agent` before and after every real-home batch and reads any change before the log is written.

## Related
BUG-107 (the `--home` wrapper covers `sh-agent` calls only), BUG-115 (sandbox homes). Both precede this; the wrapper rule was followed but its scope was too narrow.
