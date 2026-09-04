# Roadmap — styrman-harness

<!-- Division of labor: ROADMAP.md = where the project is going; TODO.md = small deferred items; .agents/docs/plans/ = how. -->

Last reviewed: 2026-09-04

## Now

Real use. The rebuild is verified on both the Linux box (both CLIs) and the Mac (Antigravity with `--vendor antigravity`, `.agents/docs/testing/2026-09-04-verify-mac-real-home.md`), installed in the operator's real home on both machines. No tuning before two weeks of use (ADR-009).

## Next

- After two weeks: `/retro` over the feedback log; fold accepted proposals into `context/`, skills, or agents; count agent dispatches from the changelogs (ADR-009 cost check).
- Re-run the skill matrix (`.agents/docs/testing/skill-matrix.md`) after any change to `context/`, `agents/`, or a workflow skill.

## Later / Ideas

- A `sh-agent sandbox` subcommand if re-running the matrix by hand becomes routine (ADR-010 leaves the door open).
- Nothing else scheduled. A rejected option is a decision, not a plan — see `REQUIREMENTS.md` Non-goals.
