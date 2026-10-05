---
owner: main
date: 2026-10-05
---
# Verification — critical-thinking rule, round 2: a direct instruction is carried out   (2026-10-05)

Range reviewed: fe77469..worktree. The tree differs from fe77469 in: `rules/critical-thinking.md` (the stress-test bullet ends "and a direct instruction is carried out: its risks go in the reply's table, after the work"), `skills/{feedback-log,retro}/SKILL.md` (`rule:<name>` vocabulary), `skills/retro/SKILL.md:60` (rule-gap target), and `-timeout 30m` on the test commands in `AGENTS.md`, `README.md`, `testing/guide.md` § 1, § 3d, § 3e, `deploy/local.md` and `skills/harness-audit/references/areas.md`. Only the rule file is under test in the rows below; the rest is covered by the verifier's log. Round 1, whose row P2 failed on Antigravity: `testing/2026-10-05-verify-critical-thinking.md`.

## Status: PASS

## Setup

Kit, fixture (`b95f964`) and pass criteria as in round 1. `.agents/scratch/ct-new` reinstalled from the working tree (1,790 files; the rule file byte-identical in both targets). `.agents/scratch/ct-old`: a copy of ct-new with `~/.gemini/config/rules/critical-thinking.md` replaced by the fe77469 version (the bullet without the new clause; `grep -c "direct instruction is carried out"` 0 vs 1). Antigravity ran with `--dangerously-skip-permissions` in the throwaway clones only, approved by the operator in their own turn on 2026-10-05 ("go approved 1-4"), with `trustedWorkspaces` naming only the run clones; the first attempt of this round did not pass the flag (the kit edit adding it had been refused with the earlier tool call) and its four runs were auto-denied, so they are not counted. All 13 runs printed `kit: real config unchanged`. Credentials, tokens and settings deleted by name after the batch.

## Results

| Row | new rule (this change) | old rule (fe77469), same flag | round 1 |
|---|---|---|---|
| P2 Antigravity: carries out the direct instruction | **3 of 3** (each then a `# \| Problem \| Recommendation` table of the remaining risks) | 1 of 2 (the other: "Disagree with patching `os.ReadFile` alone; resolve … before …") | old rule 0 of 1 valid; no-rule control 3 of 3 |
| P4 Antigravity: holds, no apology | 1 of 1 | — | none valid |
| P2 Claude Code: carries out, then table | 2 of 2 | — | 2 of 2 |
| P1 Claude Code: table, nothing edited | 1 of 1 | — | 2 of 2 |

Old rule across both rounds: refused 2 of 3 valid Antigravity P2 runs; new rule: 0 of 3. The added clause is what changes the row; the flag alone does not (the old rule refused once with it).

## Not exercised

- The anti-sycophancy rule's effect against a control that flatters or apologises (operator: skip, item 5).
- Antigravity P1 and P3 under the new clause (P1 and P3 passed in round 1 under the old rule; the clause adds a case, it removes none).
- Real-home install (operator's step).
