# Roadmap — styrman-harness

<!-- Division of labor: ROADMAP.md = where the project is going; TODO.md = small deferred items; .agents/docs/plans/ = how. -->

Last reviewed: 2026-09-13

## Now

`.agents/docs/plans/2026-09-13-learning-loop-project-level.md`: Phases 1–2 done (installer rebuilds from source on every install, all skills model-invocable) and committed as 82fd4a6 and bd5ed87. Phase 3 (learning loop at project level: per-project retro records, scope rule, global index) waits on the operator's answer to the plan's D11 residue (the `rule_ref` vocabulary and grep-spacing items). Phase 4 (docs and trail) is closing. Phase 5 (an improvement review of the whole harness) produces a numbered go/no-go list for the operator, not more code.

## Next

- Fold the accepted items from the Phase 5 improvement review, once the operator has picked them, into `context/`, skills, agents, or this roadmap.
- Re-run the skill matrix (`.agents/docs/testing/skill-matrix.md`) after Phase 3 lands.
- Open question: `~/.config/sh-agent/` is per machine, so a retro run on one box is invisible to a harness session on another until the project record is pulled there — raised in the improvement review, not yet resolved.

## Later / Ideas

- A `sh-agent sandbox` subcommand if re-running the matrix by hand becomes routine (ADR-010 leaves the door open).
- Nothing else scheduled. A rejected option is a decision, not a plan — see `REQUIREMENTS.md` Non-goals.
