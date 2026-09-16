---
owner: reviewer
date: 2026-09-16
---
# Review — bug-038-docs   (2026-09-16)

Range reviewed: working tree — changed files: `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md` (docs lens per `skills/review/SKILL.md`'s `--docs` row)

## Verdict: APPROVE WITH NITS

The four agent files are internally consistent (identical wording, same location, same BUG-038 citation) and nothing already in the living docs now contradicts them. The gap is propagation: BUG-038's own Prevention/Related sections ask for this rule to also land in `context/AGENTS.md`'s subagent description — the Related section is explicit that it belongs there "as a global hard rule, not per-agent text" — and this diff does not touch `context/AGENTS.md` at all, so the rule exists only inside this repo's four agent definitions and nowhere in the text that describes what a subagent is in general.

## Findings

- **[medium / confidence: high]** `context/AGENTS.md` (Subagent description line, `Skill = knowledge... Subagent = judgement independent of the author...`) and `context/AGENTS.md` § Hard rules / § Git — neither mentions that a read-only subagent never runs a destructive git command. `.agents/docs/bugs/BUG-038-verifier-git-checkout-discarded-fix.md`'s own Prevention section says the rule "is a candidate for `context/AGENTS.md` where it describes what a subagent may do," and its Related section goes further: "recommend the rule above land in `context/AGENTS.md`'s subagent section as a global hard rule, not per-agent text; it is actionable ... and generalizable beyond this repo." The diff implements only the per-agent half of that recommendation; the stronger, explicitly-named alternative is unaddressed. Concretely this means: any consumer project reading `context/AGENTS.md` directly (rather than the harness's own `agents/*.md`) has no textual signal that a subagent must not run `git checkout --`/`restore`/`stash`/`clean`/`reset`/`apply` against the tracked tree — the global context's subagent paragraph ("it gets the artifact, never the conversation"; "An agent's artifact ... is trail, not a change") talks about what an agent writes, not what it may do to the tree it reads from. Fix: add one clause or pointer to `context/AGENTS.md`'s Hard rules or subagent paragraph, e.g. "a read-only subagent never runs a git command that rewrites the tracked tree — see the agent definitions for the exact list," in the same change, or explicitly rebut the Related section's recommendation in the changelog if per-agent text is judged sufficient.
- **[low / confidence: high]** `.agents/docs/bugs/BUG-038-verifier-git-checkout-discarded-fix.md` still shows `status: open` and its Fix section still reads "Not fixed. This is a prevention-rule gap in agent definitions, not a code defect." — that gap is now closed for the four named agents by this diff. Not a defect in the reviewed files (bug records are excluded from the docs-lens grep set by the skill's own instruction, and updating it is `bugs`'/the main agent's job, not this diff's), but worth flagging so the status gets updated in the same commit that lands this fix, per the project's own bug-record discipline.

No stale wording found: `git grep -n "checkout --\|git restore\|git stash\|git reset\|git apply"` across `README.md`, `AGENTS.md`, `.agents/{rules,skills,agents,commands}/`, `.agents/docs/specs/`, `.agents/docs/testing/*/guide.md`, `.agents/docs/deploy/` turned up nothing that describes subagent git usage in a way this rule now contradicts or narrows. No skill that dispatches these four agents (`review`, `preflight`, `bug`, `feature` — checked in `skills/review/SKILL.md`, `skills/preflight/SKILL.md`, `skills/bug/SKILL.md`, `skills/feature/SKILL.md`) instructs them to reproduce an earlier state via a destructive git command; all git usage named for them (`status`, `diff`, `log`, `show`, `describe`) is read-only and stays valid under the new rule.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Skills beyond `review`, `preflight`, `bug`, `feature` that may dispatch these four agents (e.g. `retro`, `distill`, `handoff`, `init-project`) were not individually grepped for git instructions to the agents — a quick `grep -rn "verifier\|reviewer\|security\|explorer" skills/*/SKILL.md` beyond the four named skills was not run. `internal/harness` code and any Go tests were not read (out of scope for a docs-lens review of markdown-only changes). The full `deploy/*.md` files were not opened line-by-line, only grepped.

## Outcome

- **[medium] `context/AGENTS.md` carries no mention of the rule** — `fixed` at `context/AGENTS.md:45`: one clause appended to the Git paragraph, so no rendered line is added. `sh-agent render` still prints 93 lines; a scratch install reports drift none with the sentence in both `.claude/AGENTS.md` and `.gemini/config/rules/AGENTS.md`; one start of each CLI in the sandbox, with no tools, quoted the sentence verbatim from its loaded context, and the Antigravity log that start created has `Failed to parse` 0 times.
- **[low] BUG-038's record still read open** — `fixed`: the record and its INDEX row are closed in this commit, the Fix section naming the four definitions and the context clause.
