---
owner: main
date: 2026-09-15
---
# Harness practices Phase 6: end to end in the sandbox on both CLIs, with Phase 7's targeted rows     (2026-09-15)

## Summary

The rows of `.agents/docs/testing/skill-matrix.md` ran on Claude Code and on Antigravity against this plan's harness in the sandbox home. The log is `testing/2026-09-15-verify-harness-practices-e2e.md`.

Every row passes on both CLIs, directly or on its one re-run, with four exceptions, all on Antigravity: `bug`, `review-project`, `preflight` and `documentation`. Each ran into the 20-minute print timeout twice, while five CLI processes ran at once. Each then passed twice under a lighter load: once with c192622's harness, and once with this tree under that same load. D50 records these four as caused outside this plan.

The skills that changed after the matrix passes ran again at fcf3b79:
- `init-project`, through the `init-update` and `init-deploy` rows (7.3);
- `preflight`, on both CLIs;
- `bug`, `bug-no-repro`, `feature` and `handoff`, for 641eccd's edits.

Two harness findings stay open, both in the plan's Follow-ups:
- BUG-031: a print-mode run can end its turn while an agent it dispatched is still running.
- Option C's filing clause, which the Antigravity `preflight` run did not follow.

## Type

task

## Files changed

- `.agents/docs/testing/2026-09-15-verify-harness-practices-e2e.md` — the Phase 6 log:
  - 6.2 install and drift;
  - 6.3 CLI starts;
  - the 6.4 table per row and CLI, with the D50 classification and the read-back of each C1 pass;
  - the Phase 7 rows at fcf3b79 (I1, P7, T7);
  - the evidence columns and their limits;
  - the edit-to-row map;
  - what no row examined.
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`:
  - 6.4, the Phase 6 Checkpoint and 7.3 ticked with evidence;
  - Follow-ups for option C's filing clause on Antigravity, and for `gofmt -l .` walking the gitignored scratch.
  - D50 and the other Follow-ups (BUG-031, BUG-029's rule) landed in fcf3b79.

## Agents dispatched

- `bugs`: BUG-029, BUG-030 and BUG-031 (fcf3b79).
- `reviewer`: twice, the code and docs lenses of the Phase 7 review round 2 (2a2e6a1).

Inside the sandbox clones each matrix run dispatched its own agents; the log's evidence table lists them.

## Verification

The detached kit in `.agents/scratch/phase6-kit/` ran each row in a fresh clone. Before each CLI started, it ran the row's check once as a control that must fail. After each run, the main agent read what the run produced, not only its check.

- **Passes 1 and 2** ran at `22c106c+dirty`, 19:48–22:06: 30 row–CLI pairs.
  - 26 pass: directly, on the corrected check of BUG-024, or on their re-run.
  - 4 failed twice.
- **Comparison (B1)** at c192622, 22:09–22:26, with the plain-copy install at rc=0: the four passed in 323, 230, 504 and 774 s.
- **Same-conditions control (C1)** at fcf3b79, 22:26–22:42, after a reinstall with rc=0, drift none and a binary line with no warning: the four passed in 375, 209, 340 and 713 s.
- **Phase 7 rows at fcf3b79**, every control failing before its run:
  - `init-update`: Claude Code 121 s, Antigravity 107 s. Each made one commit that only adds, and no file was deleted, modified or renamed.
  - `init-deploy`: Claude Code 147 s, Antigravity 115 s. Only added files; Claude Code's headless guard refused the `.claude` links, and the run said so.
  - `preflight` on Claude Code (P7): 975 s, `VERDICT: SHIP`, with the Step 4 advisory section.
  - Antigravity T7: `bug-no-repro` 94 s, stopped at the gate with the question in its handoff; `feature` 250 s, a plan in the new template with the cold read and no source changes; `handoff` 90 s.
  - Claude Code T7: `handoff` 300 s, deferred items moved from `TODO.md` into the handoff. Claude Code T7 also ran `bug` (981 s): an `explorer` diagnosis, the fix, BUG-002 and BUG-003, `reviewer` and `security` reports, a verifier log, `bugs` and `changelog`, committed. It ran `bug-no-repro` (144 s), which stopped at the gate with the question in its handoff, and `feature` (494 s), which wrote a plan in the new template with the cold read and no source changes.
  - Antigravity `bug` and `documentation` at this tree are their C1 runs.
- **No push:** `demo-app` and `demo-app-agy` have no commit or reflog entry after 2026-09-14.
- **No parse errors:** the newest Antigravity logs print `Failed to parse` 0 times.
- **`gofmt -l .` walked into gitignored scratch** and printed five Go files from this session's finished `init-project` survey workflow: `.agents/scratch/adv-init`, 108M of project copies, and `.agents/scratch/lintcheck`. Nothing tracked them, referenced them or was using them, so both were removed. `gofmt -l .` now prints nothing.

## Docs

The log and the plan are this phase's docs. `README.md`, the testing guide and the deploy guides are unchanged: the matrix spec already names every row that ran.

## Requirements

n/a
