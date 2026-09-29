---
owner: main
date: 2026-09-29
---
# `/review --prompts`: the prompt-audit lens     (date: 2026-09-29)

## Summary
`/review` gains a `--prompts` lens. It dispatches the `reviewer` agent with Anthropic's prompt-audit reference and audits the instruction files in the scope (`AGENTS.md`, rules, skills, agents, commands; in this checkout `context/`, `agents/`, `skills/`, `commands/`, `rules/`) for dated prompt patterns, stale paths and flags, contradicting instruction files, redundant sub-agents, and writes findings in the reference's report shape plus a proposed diff, one hunk per finding, never applied. Two rules the brief adds for this harness: the audited files are read by Claude and by Gemini, so a finding whose only reason is one model's documented behaviour is `flag`; and a prohibition that traces to a `bugs/INDEX.md` row or an ADR is the reference's keep-list item 5, not a finding. `--fix` on the lens has the main agent apply the diff's hunks.

The reference (`skills/claude-api/shared/prompt-audit.md` from github.com/anthropics/skills, Apache-2.0) is vendored verbatim at pin `8a1541c` under `skills/review-standards/references/prompt-audit/` with its `LICENSE.txt`, listed in `skills/VENDORED.md`. The operator chose this shape (a lens, option 3 of `plans/2026-09-29-research-memory-and-prompt-audit.md`) over a standalone skill.

## Type
feature

## Files changed
`skills/review/SKILL.md` (description, argument-hint, usage line, Step 1 flag list, lens table row, Step 5 `--fix` mapping); `skills/review-standards/SKILL.md` (one "Lens reference" line); `skills/review-standards/references/prompt-audit/{prompt-audit.md,LICENSE.txt}` (new, vendored); `skills/VENDORED.md` (one row); `.agents/docs/testing/skill-matrix.md` (row `review-prompts`); the research record's decision line. README, `context/AGENTS.md`, the deploy guides and the specs enumerate no lens list and are unchanged (docs lens, below).

## Agents dispatched
- `verifier` ×2 — `testing/2026-09-29-verify-prompts-lens.md` (PASS) and `-r2.md` (PASS, over the tree after the review fixes): Go suite, `sh-agent list`, render under 100 lines, install into a scratch home with drift none, the reference present under both vendors' skill directories, digest `c33200da…` as `VENDORED.md` states.
- `reviewer` (code lens) — `reviews/2026-09-29-review-worktree-code.md`, APPROVE WITH NITS; two findings fixed (the matrix check now asserts every seed; the reference bullet has its own heading). `## Outcome` on the report.
- `reviewer` (docs lens, the `/review --docs` of this change) — `reviews/2026-09-29-review-worktree-docs.md`, APPROVE WITH NITS; one finding fixed (the brief's Group 4 clause had dropped the reference's sub-agent roster-check exception), one rebutted (`--add-dir` in the matrix row is needed because the brief makes the reviewer Read outside the project). `## Outcome` on the report.
- `documentation`: not dispatched — no file it owns changed; the `review` skill file is the documentation of `/review`.

## Verification
- Behavioural, main agent, `testing/2026-09-29-prompts-lens-behavioural.md`: the brief handed to the working-tree `reviewer` via `claude -p --agents` over a seeded fixture (`.agents/scratch/prompts-probe/`): round 1, 3 of 3 runs found the dead path, the contradiction and the pressure line, kept the bug-backed prohibition, wrote `## Proposed diff` and touched nothing; round 2 with the corrected brief and a near-duplicate agent pair, 3 of 3 found the pair as well; the matrix row's check passes on every round-2 clone. The one control run with the old Group 4 wording also found the pair, recorded as such. Refused tool calls per run (a heredoc write redone with the Write tool) are counted separately in the log.
- Not exercised: an Antigravity run of the lens (no sandbox home on this box), the shipped reference path in the real home (working-tree path substituted; install tier shows placement), the `--prompts` token through the skill's own parser, `--fix` on the lens. After `sh-agent update` on a real home, the three smoke checks of `testing/2026-09-03-verify-real-home.md` are the operator's.

## Notes
- `TestEmbedHoldsOnlyTrackedFiles` fails on an untracked vendored file until it is staged; staging by name resolved it (expected behaviour of the embed test, not a defect).
- First real use is up to the operator: `/review context/ agents/ skills/ --prompts` in this checkout would be the first audit of the harness's own text; its findings then go through `--fix` or a task, hunk by hunk.
