---
owner: bugs
date: 2026-09-30
id: BUG-115
status: fixed
area: testing (agy sandbox home)
---
# BUG-115: first agy-only sandbox runs for BUG-113 used a copy of the real home's Antigravity settings

## Symptom
While testing the BUG-113 fix, the main agent built a sandbox home (`.agents/scratch/agyonly-home`) and wrote its `~/.gemini/antigravity-cli/settings.json` as a copy of the real home's file (keys `artifactReviewPolicy` and `toolPermission` "always-proceed"), with only `trustedWorkspaces` narrowed to the fixture. AGENTS.md § Testing says a sandbox home for headless agy gets its own settings file "written for the run with the fewest grants the run needs ... never a copy of the real home's file, which approves every tool". Found by the main agent itself while reading § Testing at the release gate. The runs (fixtures `~/sandbox/agyonly/proj`, `proj2`) did not count as gate evidence.

## Repro
Read the sandbox settings file of the first BUG-113 runs: it carried `toolPermission` "always-proceed".

## Root cause
The main agent built the sandbox from memory of an older procedure instead of reading AGENTS.md § Testing and `testing/skill-matrix.md` step 1 (the agy settings line) before the first sandbox run.

## Fix
No source change. The gate reran both rows (`init-agy-only`, `dream-agy-only`) in a fresh sandbox home `.agents/scratch/gate17-agy` with the matrix's named grant list, no `toolPermission` key (`grep -c toolPermission` gave 0), `trustedWorkspaces` = the clone only. Named grants alone got every run auto-denied: agy writes multi-line scripts whose first token is a loop (`for d in ...`) or an assignment (`root=$(git rev-parse ...)`), which no `command(<name>)` grant covers; adding `command(for)` did not help. The session's permission classifier refused `--dangerously-skip-permissions` until the operator approved it in their own turn ("Ok, do all of these"), as testing/guide.md § 5 allows inside a throwaway clone. Rerun: both CHECKs exit 0, no auto-denied line, 0 invalid allow entries. The copied token and settings file were deleted by name after each batch. The skill-matrix rows `init-agy-only` and `dream-agy-only` now name the skip flag and why.

## Prevention
1. Before building any sandbox home, read AGENTS.md § Testing and the skill-matrix agy settings line and write the settings file from them, never from the real home.
2. An agy scenario whose skill runs compound shell scripts (/init-project, /dream) needs `--dangerously-skip-permissions` in its throwaway clone, decided before the run.
The first rule already existed in AGENTS.md § Testing; it was not followed.

## Related
BUG-113 (the fix under test); BUG-087, BUG-090, BUG-091, BUG-106, BUG-107 (behavioural-run setup rules in AGENTS.md § Testing).
