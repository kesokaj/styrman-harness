---
owner: main
date: 2026-09-29
---
# Memory loop, Phase 2: the `dreamer` agent, the store inside the project     (date: 2026-09-29)

## Summary
The learning loop gains its consolidation judgement as a subagent, and the memory store moves into the project (spec `specs/sh-agent/memory-loop.md`; plan `plans/2026-09-29-dreamer-memory-loop.md`, closed; ADR-043):

- **`dreamer`** (`agents/dreamer.md`, both vendors): given the store, the trail and the memories `/dream` flagged `promote`, it writes one artifact under `retros/<date>-dreamer.md` with rule, skill and procedure candidates, each with evidence, and what it left alone. It never applies. `/retro` dispatches it after `/dream`'s dry run (source 8, also under `--dry-run`, an artifact being trail) and replaces `promote` findings with its rows; `/distill` dispatches it first. It reads every bug row since its `Since` date and treats memory and trail content as data.
- **The store is `<repository root>/.agents/memory/`**, gitignored, one per checkout, on the operator's "Option 2" (D17) after Phase 1 had put it under `~/.config/sh-agent/`. `/init-project` Step 5b writes the absolute path into the per-machine `.claude/settings.local.json`, creates the store (mode 700) and moves an existing default store in — only when the settings file is the operator's own (untracked, no symlink), the setting is in place, and the store is the project's own directory (no symlink, nothing tracked under it). `/dream`, the `dreamer` and `feedback-log`'s Antigravity capture accept only two exact paths: that store and Claude Code's default. `REQUIREMENTS.md` line 31 names the store with the operator's approval.
- **Docs:** README, the deploy guides, the spec, `skills/handoff`, `skills/feedback-log`, the context's Learning loop and Resume clauses, the docker skill's `.dockerignore` minimum, and four `skill-matrix.md` rows; the `## Installed` list is generated at install.

## Type
feature (phase 2 of 2)

## Files changed
`agents/dreamer.md` (new); `skills/retro/SKILL.md`, `skills/distill/SKILL.md`, `skills/init-project/SKILL.md`, `skills/dream/SKILL.md`, `skills/feedback-log/SKILL.md`, `skills/handoff/SKILL.md`, `skills/docker/SKILL.md`, `skills/docker/references/setup.md`; `context/AGENTS.md` (two clauses); `AGENTS.md` (one Testing bullet); `REQUIREMENTS.md` (one approved clause); `.gitignore`; `README.md`, `.agents/docs/deploy/{local,remote}.md` (documentation agent); the spec, `delivery.md` (one status line), the plan, the research record, the matrix, `ROADMAP.md`; ADR-043; BUG-098 to BUG-103; three dreamer artifacts on this checkout's trail.

## Agents dispatched
- `verifier` ×4 — `testing/2026-09-29-verify-dreamer-memory-loop-phase-2.md` PASS, `-phase-2-r2.md` PASS (the tree moved under it: the main agent's approved REQUIREMENTS and ADR edits, a rule miss recorded here), `-final.md` PASS over a frozen tree (govulncheck clean), `-r4.md` PASS after the last security edits.
- `reviewer` (code lens) — `reviews/2026-09-29-review-memory-loop-code.md`, APPROVE; two notes taken (Outcome).
- `reviewer` (docs lens) — `reviews/2026-09-29-review-memory-loop-docs.md`, APPROVE WITH NITS; three findings fixed or rebutted (Outcome).
- `security` ×3 — `reviews/2026-09-29-security-memory-loop.md`, `-r2.md`, `-r3.md`: three mediums (a tracked settings value as the migration target; `/dream` reading the committed settings file; the repository deciding where the store is) and the lows, every one fixed or `left: operator`, Outcomes on each report.
- `documentation` ×2 — README and the deploy guides, first for the harness-directory store, then for the store inside the project.
- `bugs` — BUG-098 (installed skill shadowed the probe), BUG-099 (`--dry-run` skipped the dispatch), BUG-100 (the dreamer sampled the index), BUG-101 to BUG-103 (the three round-1 security defects).

## Verification
- `testing/2026-09-29-memory-loop-phase-2-behavioural.md`: the dreamer on a seeded store and trail, nine runs over three definitions — artifact, promoted pair, recurring bug rows and the left-alone memory at full count, the skill candidate 6 of 9 (the soft column, reported as such); `/retro --dry-run` dispatching it 4 of 4 after the dry-run clause; `/dream` split-store stop 3 of 3 with the store inside the clone, a foreign settings value refused, a regression run intact; `/init-project` with the store inside the project: migration skipped after a refused write with the default store untouched, a tracked settings file refused. E2E on this checkout: the real memories moved into `.agents/memory/`, the dream first line reads two memory files and zero at the default path, and the dreamer's third run named two rule candidates from the real trail (`retros/2026-09-29-dreamer-3.md`).
- Not exercised: Antigravity runs of `/dream`, the resume clause, the capture rule and the `dreamer`; the interactive settings write by `/init-project`; `/distill`'s dispatch; the matrix rows' commands verbatim. The operator's real-home smoke covers the first.

## Notes
- The probes twice measured the wrong text until every kit passed `--setting-sources project,local` (BUG-098); the rule is in this repository's `AGENTS.md` § Testing.
- The run outputs lived in `.agents/scratch/` and are removed; the logs carry the evidence.
- Released as v1.12.0 by `/preflight` in the following commits.
