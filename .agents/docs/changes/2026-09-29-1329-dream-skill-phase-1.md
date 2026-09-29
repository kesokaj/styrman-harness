---
owner: main
date: 2026-09-29
---
# `dream` skill, Phase 1: the skill, proven on a seeded store     (date: 2026-09-29)

## Summary
`skills/dream/SKILL.md` is written per the plan's Interfaces (`plans/2026-09-29-dream-skill.md`): it resolves the project's Claude Code auto memory store (the `autoMemoryDirectory` setting or the per-project default), reads every memory, gives each one verdict — `keep`, `graduated`, `stale`, `broken`, `merge`, with `promote` as a flag on a surviving memory whose fact came from two sessions — rewrites `MEMORY.md` from the files, puts the table to the operator (`--dry-run` stops there; headless runs apply nothing), records to `retros/<date>-dream.md` and `retros.jsonl`, and hands promotions to `/retro`. It never writes a rule. Doubt and any check that cannot run resolve to `keep`.

The plan went through three cold-reader rounds (`reviews/2026-09-29-review-dream-skill-coldread-r1.md` to `-r3.md`); at the cap the operator chose "Proceed, add the rows (Recommended)", recorded as D12. Decisions D4–D6 came from the operator's "Go all"; D7–D13 from the rounds and the probe.

## Type
feature (phase 1 of 2)

## Files changed
`skills/dream/SKILL.md` (new); the plan (D4–D13, Phase 1 ticked); `ROADMAP.md` "Now" (earlier commit 80bf5b6).

## Agents dispatched
- `reviewer` ×3, the cold-reader check (r1: two blocking gaps; r2: one; r3: one new item at the cap).
- `verifier` — `testing/2026-09-29-verify-dream-skill-phase-1.md`, PASS: Go suite, `sh-agent list` with the `dream` row, install into a scratch home (264 files), status 28 skills drift none, the skill under both vendors' skill directories, frontmatter quoted.
- No `reviewer` over the diff yet and no `documentation`: the skill is the only source file so far; the Definition of Done runs over the whole change at Phase 2's checkpoint.

## Verification
- `testing/2026-09-29-dream-phase-1-behavioural.md`: twelve `claude -p` runs on a fixture repository with a seeded six-file store redirected by `--settings`. Round 1 exposed two setup defects (a seed restating the global context's asking rule, correctly marked `graduated` citing `~/.claude/AGENTS.md:6`; controls that could still read the harness skill directory). Round 2: 3 of 3 skill runs on the store line, every seeded verdict, the `promote` flag, the guard `keep`, the load marker and an untouched store; 0 of 3 clean controls on the skill's vocabulary. Executed and refused calls counted separately.
- Not exercised: the apply path, Antigravity, the `feedback.jsonl` route to `promote`, `/retro` reading the verdicts — all Phase 2 or named in the log.

## Notes
- Headless load evidence is the slash-command expansion in the first user message, not a `Skill` tool call (D13).
- The memory docs (read 2026-09-29) also document a built-in `/doctor prompt-audit` in Claude Code 2.1.283+; noted for the `--prompts` lens shipped earlier today, not for this plan.
