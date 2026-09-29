---
owner: main
date: 2026-09-29
---
# Handoff — 2026-09-29 — `/review --prompts` shipped; dream-skill plan proposed, awaiting the operator

Session ID: 8a02b8f1-6422-4650-9dfd-3ea55bba9b1b
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 11b23b6, pushed; this handoff's commit follows it.

## Context

The operator asked whether memory or "dreaming" plugins exist for Claude Code and Antigravity, then for a consolidation-skill proposal, then for broader research on Antigravity memory and on whether Anthropic's `/claude-api prompt-audit` could be implemented for both CLIs in this harness. They chose option 3 of that research, a `/review` lens, and it shipped this session.

## Completed

- **Research record** `.agents/docs/plans/2026-09-29-research-memory-and-prompt-audit.md`: Antigravity has no usable memory feature (editor Knowledge Items undocumented and unreliable; agy exposes none, only latent Knowledge Base code and an empty `~/.gemini/antigravity-cli/knowledge/`); every community plugin needs hooks or a daemon, so nothing to adopt. `prompt-audit` exists (one Apache-2.0 reference file in github.com/anthropics/skills) and fits this configuration repository.
- **`/review --prompts`** (commit 11b23b6, `changes/2026-09-29-1251-review-prompts-lens.md`): a lens row in `skills/review/SKILL.md` dispatching `reviewer` with the prompt-audit reference over the instruction files in scope; two harness rules in the brief (a Claude-only rationale is `flag` because Gemini reads the same files; a prohibition backed by a `bugs/INDEX.md` row or an ADR is a keep). The reference is vendored verbatim at pin `8a1541c` under `skills/review-standards/references/prompt-audit/` with `LICENSE.txt`; row in `skills/VENDORED.md`; matrix row `review-prompts` in `.agents/docs/testing/skill-matrix.md`.
- Definition of Done for that change: `verifier` PASS twice (`testing/2026-09-29-verify-prompts-lens.md`, `-r2.md`), code and docs lenses APPROVE WITH NITS with every finding fixed or rebutted in each report's `## Outcome`, behavioural runs 3 of 3 twice on a seeded fixture (`testing/2026-09-29-prompts-lens-behavioural.md`), changelog, commit, push.
- **Dream-skill plan** `.agents/docs/plans/2026-09-29-dream-skill.md` (commits 21142d3, 0cd185c): a proposal for a harness skill that consolidates Claude Code's auto memory and hands promotions to `/retro`. Status proposed; not in `ROADMAP.md` "Now" because it is not in progress.

## Open / In Progress

- The dream-skill plan awaits the operator's answer on its three Open rows (O1 trigger: retro only recommended; O2 standalone record in `retros/` recommended; O3 this project's store only recommended). Nothing is being built. Its unticked boxes are the plan's future phases, deliberately outside "Now" until a go.
- Requirements edits awaiting approval: none.

## Key Decisions

- **Lens, not skill** (operator, 2026-09-29): prompt-audit lands as `/review --prompts` (option 3) rather than a standalone vendored skill (option 2, the recommendation). Recorded on the research record.
- **Reference lives under `review-standards/references/`**, the skill the reviewer already preloads; the brief spells both vendors' installed paths because skills deploy verbatim and only agent bodies get the `~/.claude/skills/` to `~/.gemini/config/skills/` rewrite (`internal/harness/derive.go`).
- **Antigravity memory: nothing adopted** (research recommendation A). The dream plan stays Claude Code only; re-check `ls ~/.gemini/antigravity-cli/knowledge/` after each agy update.
- Rejected for the `--prompts` docs-lens finding on `--add-dir`: the other `/review` matrix rows are not under-permissioned; the reviewer's `skills:` preload carries `review-standards` into context, while this brief makes it Read a file outside the project.

## Corrections

- None. The feedback log has no entry for 2026-09-29; the operator's choice of option 3 over the recommendation is a decision, recorded above, not a correction.

## Agent notes (unverified)

- The control run with the old Group 4 wording still found the seeded redundant-agent pair (1 of 1): the reviewer follows the reference over a contradicting paraphrase in its brief. One run, not a rate.
- In every behavioural run the reviewer first tried to write its report with a Bash heredoc, was refused, then used the Write tool. Harmless here; a candidate observation for a future retro on `agents/reviewer.md`, not acted on.
- `TestEmbedHoldsOnlyTrackedFiles` fails on an untracked vendored file until it is staged; expected behaviour of the embed test.
- The real home on this box is still installed from 0d7dc32 (v1.11.0): `/review --prompts` is not yet available there until the operator runs `sh-agent update`.

## Environment State

- Branch: `main`, clean, pushed at 11b23b6 (plus this handoff).
- Running processes: none. Background probe runs finished.
- Scratch: `.agents/scratch/prompts-probe/` (fixture, kit, six runs and a control), `.agents/scratch/prompts-home/`, `.agents/scratch/verify-prompts-home/` — gitignored, removable.
- Real home untouched.

## Artifacts to Read

- `.agents/docs/changes/2026-09-29-1251-review-prompts-lens.md` — what shipped, agents dispatched, what was not exercised.
- `.agents/docs/plans/2026-09-29-research-memory-and-prompt-audit.md` — both research questions, the comparison tables, the decision line.
- `.agents/docs/plans/2026-09-29-dream-skill.md` — the proposed consolidation skill and its three Open rows.
- `.agents/docs/testing/2026-09-29-prompts-lens-behavioural.md` — the seeded-fixture runs, executed and refused calls counted separately.
- `.agents/docs/reviews/2026-09-29-review-worktree-docs.md` — the docs lens with the Group 4 finding and its Outcome.

## Next Steps

1. **Operator:** decide the dream-skill plan (go per Open row or go all, or strike); on a go, the feature skill starts at its cold-reader check and `ROADMAP.md` "Now" names the plan.
2. **Operator:** `sh-agent update` on this box and the other machines, then the three smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI; only then is `/review --prompts` available from a real home.
3. **Optional first real use:** `/review context/ agents/ skills/ --prompts` in this checkout, findings taken hunk by hunk through `--fix` or a task; an Antigravity run of the lens has not happened and would be the first.
4. Nothing else is queued. `TODO.md` and `ROADMAP.md` are empty.
