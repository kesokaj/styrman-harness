---
owner: verifier
date: 2026-09-18
---
# Verification — sh-agent-works-everywhere-phase-3   (2026-09-18)

Range reviewed: tree e7fec88

## Status: PASS

Scope: re-ran `.agents/scratch/matrix/run.sh` on tree e7fec88 (after `rm -rf` of the
prior fixtures, keeping `run.sh`) and checked `.agents/docs/testing/sh-agent-matrix.md`
(written at tree 3209cf3) against the new `run.log`, cell by cell, plus the round-2
0700 `.config/sh-agent` check. Complete: every check the dispatch named ran and passed.
The matrix file itself was not edited (out of scope).

## Commands
- `rm -rf .agents/scratch/matrix/{checkout,plain,rel,dev,home-*,widget,nogh,run.log}` → left only `run.sh`
- `bash .agents/scratch/matrix/run.sh` → exit 0, ~7s wall time (`time` output: real 0m7.237s)
- Produced `.agents/scratch/matrix/run.log`, 485 lines, 53 `[exit N]` lines (matches the matrix's "53 cells")

## Checks performed

1. **Exit codes.** All 53 `[exit N]` lines read and matched what the matrix implies for
   each cell (0 for exercised/refused-but-non-fatal cells that the matrix calls
   "exercised"/"refused" without stating an exit code; 1 for every `status` after drift,
   `status` after uninstall, `update --check`/`update` on plain and release,
   `adopt` on release, `install --repo` and `$SH_AGENT_REPO` on release — these are the
   cells the matrix itself marks `refused`). No exit code was unexpected for its cell.
2. **Quoted messages.** Every message quoted verbatim in the matrix table was grepped
   against the new log and found verbatim, including:
   - `plain copy, not its own git checkout; update is refused` (status)
   - `plain copy, not its own git checkout: the running binary is installed as is` (install --dry-run/install)
   - `source embedded v9.9.9 (release binary; update downloads the newest release)` (status)
   - `sh-agent: not its own git checkout; a plain copy cannot be updated with git — \`install\` still deploys it as it is; to move forward, use a checkout (\`git clone git@github.com:kesokaj/styrman-harness.git\`) or a release binary (README § Install without cloning)` (plain update --check/update)
   - `sh-agent: update: this install has no checkout; install gh and run gh auth login, or download a newer release by hand (README § Install without cloning)` (release update --check/update)
   - `adopted into …/skills/widget` + the two different `next:` lines for checkout vs plain
   - `sh-agent: adopt needs a harness checkout (…); this install is a release binary and carries no tree to adopt into` (release adopt)
   - `report: harness is a plain copy, not a git checkout: harness_commit is null` / `report: harness is the release binary v9.9.9, not a checkout: harness_commit is null, harness_version is set` (report --dry-run)
   - `sh-agent: this is the release binary v9.9.9: it installs the harness built into it and never a tree on disk, so --repo=… is refused — unset it, or build sh-agent from that tree and run that build's install` / the `$SH_AGENT_REPO=…` variant
   All matched. The only difference from the matrix's quoted text is the commit SHA
   (matrix ran at 3209cf3, this run at e7fec88 — `main @ e7fec88, clean`, `from commit
   e7fec88`, `up to date with upstream (e7fec88)`), which is the expected consequence of
   re-running at a newer tree, not a discrepancy.
3. **No blank cells.** All 53 log entries have output and an exit line; every matrix row
   maps to a cell that ran. `report` (real, non-dry-run) and the `linux/arm64` /
   darwin-signature / real-`gh`-on-this-box remainder rows are marked `remainder` in the
   matrix and correctly not run by this script (per the dispatch's hard limit: never run
   `report` without `--dry-run`).
4. **Fixture assumptions (plan § Risks).**
   - Plain copy `status` → `repo …/plain (plain copy, not its own git checkout; update is
     refused)` — confirmed, line 203 of run.log.
   - Release-stamped binary `status` → `source embedded v9.9.9 (release binary; update
     downloads the newest release)` — confirmed, line 359 of run.log.
5. **0700 check (round-2, not in the matrix).** After each kind's `install`,
   `stat -c '%a' home-<kind>/.config` = `700` for all three homes (checkout, plain,
   release) — confirmed by direct `stat` after the run (the `.config/sh-agent`
   subdirectory itself was later removed by each kind's own `uninstall` step in the
   script, but `.config` was tightened to `700` during install and the script never
   loosens it back, so the mode observed on the parent directory is the mode `sh-agent`
   left it in). `sh-agent-matrix.md` is silent about this check — it predates the
   round-2 change (commit e7fec88's own log line) and does not claim to cover it; that
   silence is expected, not a gap I am flagging as a discrepancy.

No discrepancy found between the matrix and the new log. The matrix file was not
edited (dispatch says "do not edit the matrix"); it would need a refresh of the tree
SHA in its "How the cells were run" line (3209cf3 → e7fec88) but that is the main
agent's / documentation's job, not verification.

## Root cause  (FAIL only)
N/A — Status is PASS.

## Coverage of changed code
This dispatch re-runs an existing script against frozen Go sources; it is a
reproducibility/consistency audit of the matrix document against a fresh run, not new
source coverage. No source changed in this dispatch.

## Not examined
- Did not check every one of the ~260 individual file lines inside `install`/`uninstall`
  output beyond the `head -14` truncation the script itself applies to each `run()`
  call — the log itself only keeps the first 14 lines per command, so full listings past
  that point (e.g. all 260 `remove` lines of uninstall --dry-run) are not visible in
  run.log and were not independently re-verified line-by-line; this matches what the
  original matrix run itself could see, so it is not a new gap introduced by this
  re-run.
- Did not attempt Phase 4 (real machines) — out of scope for this dispatch.
- Did not modify or refresh `sh-agent-matrix.md`'s SHA/date line — out of scope per the
  dispatch ("do not edit the matrix").
- Did not re-verify the `linux/arm64`, darwin-signature, or real-`gh`-on-this-box
  remainder rows — the matrix and this dispatch both treat them as out of scope for this
  box.

## Environment
Linux devbox (Debian trixie), this checkout at /home/kesokaj/workspace/styrman-harness,
tree e7fec88. Script fixtures under .agents/scratch/matrix/ (gitignored).
