---
owner: bugs
date: 2026-09-16
id: BUG-040
status: fixed
area: skills/init-project (deploy template)
---
# BUG-040: init-project's deploy/local.md template shipped invented facts for a project that didn't exist yet

## Symptom
Antigravity's `/init-project` on an empty repository committed `deploy/local.md` with `PORT 8080`
under Configure and `curl -i localhost:8080/links` under Verify, for a service that did not yet exist,
and with no test command — because the template's example values read as facts rather than
placeholders, and it had no `## Test` section. Claude Code's run, working from the same template,
blanked the examples on its own and added a Test section unprompted.

## Repro
Sandbox `sim-agy-linkshelf`, commit `c520712`, `.agents/docs/deploy/local.md`.

## Root cause
`skills/init-project/assets/agents-tree/docs/deploy/local.md` carried literal example values
(`PORT`/`8080`/`healthz`) and commands in its Configure/Verify sections and had no Test section, even
though the skill instructs filling the guide "from the build/run/test answers" — the template gave no
signal for what was a fillable placeholder versus a shipped fact, so an agent that took the template
literally propagated the examples as the actual configuration.

## Fix
`skills/init-project/assets/agents-tree/docs/deploy/local.md`: Build/Run/Test blocks are
`<placeholders>` (Test section added); the Configure row is three placeholders under an HTML comment
saying to delete the row when the code reads no setting; Verify carries no example command.
`skills/init-project/SKILL.md` Step 4 now states that a placeholder the answers do not fill is
deleted, never replaced with an invented value.

## Prevention
A template ships no value that could be mistaken for a fact: an example is either a `<placeholder>`
or wrapped in a comment, and the filling instruction states explicitly what happens to a placeholder
the answers don't cover (delete it, don't invent a value).

## Related
Review: `.agents/docs/reviews/2026-09-16-review-bug-g2-deploy-template.md`. Testing:
`.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`. Sandbox re-run owed (not yet run):
per-CLI `/init-project` re-run at the fixed tree for G2/G3, per
`.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § Re-runs.

Graduated 2026-09-16 to `skills/init-project/SKILL.md` Step 4 and `skills/init-project/assets/agents-tree/docs/deploy/local.md`.
