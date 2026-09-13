# Roadmap — styrman-harness

<!-- Division of labor: ROADMAP.md = where the project is going; TODO.md = small deferred items; .agents/docs/plans/ = how. -->

Last reviewed: 2026-09-13

## Now

`.agents/docs/plans/2026-09-13-learning-loop-project-level.md`: Phases 1–5 done and pushed (82fd4a6, bd5ed87, fc7f195, 5cdd4d6): the installer rebuilds from source on every install, every skill is model-invocable, `/retro` and `/distill` record at project level with a harness queue and a cross-project index. Open: the operator runs `sh-agent update` on both machines, then answers the improvement review (`.agents/docs/reviews/2026-09-13-improvement-review.md`, 20 numbered items) per number; accepted `[fix]` items become the next tasks here.

## Next

- Fold the accepted items from the Phase 5 improvement review, once the operator has picked them, into `context/`, skills, agents, or this roadmap.
- Re-run the skill matrix (`.agents/docs/testing/skill-matrix.md`) after the accepted improvement-review items land (Phase 3 changed `context/`, `agents/` and workflow skills; the sandbox spot checks ran, the full matrix did not).
- Open question: `~/.config/sh-agent/` is per machine, so a retro run on one box is invisible to a harness session on another until the project record is pulled there — raised in the improvement review, not yet resolved.

## Later / Ideas

- A `sh-agent sandbox` subcommand if re-running the matrix by hand becomes routine (ADR-010 leaves the door open).
- Nothing else scheduled. A rejected option is a decision, not a plan — see `REQUIREMENTS.md` Non-goals.
