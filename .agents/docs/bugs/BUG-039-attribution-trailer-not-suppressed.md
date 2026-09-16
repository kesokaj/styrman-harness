---
owner: bugs
date: 2026-09-16
id: BUG-039
status: fixed
area: context/Git + cmd/sh-agent (install, status)
---
# BUG-039: Claude Code headless commits carried a Co-Authored-By trailer despite the no-attribution rule

## Symptom
On Claude Code headless (Opus 5, sandbox), every commit of both linkshelf simulation runs carried
`Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>`, although `context/AGENTS.md` §
Git says "Author is the local git config only: no Co-Authored-By". Three of three commits in the
sandbox run.

## Repro
Sandbox `/home/kesokaj/sandbox/runs/sim-claude-linkshelf`, `git log --format=%B`; every commit shows
the trailer.

## Root cause
The trailer is Claude Code's own `attribution.commit` setting (docs.claude.com/docs/en/settings-reference),
not a behaviour the harness context can switch off with prose alone. Neither the sandbox home nor the
operator's real home sets `"attribution": {"commit": ""}`, and the context rule never named the CLI
setting or the command that would report it missing — so the rule had no mechanism behind it, only an
instruction competing with the vendor's own reminder.

## Fix
`context/AGENTS.md` § Git now names the CLI reminder and the setting that switches it off.
`cmd/sh-agent/main.go` `attributionHint`/`attributionOff` (install path `main.go:165`, status path
`main.go:319`) print an advisory note when `<home>/.claude/settings.json` does not set
`"attribution": {"commit": ""}` (bounded read, regular file only, no echo of file contents).
`TestAttributionNote` (9 subtests) covers both call sites; README and deploy guides updated. Verified
in `.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md` (all 9 subtests PASS; install and
status runs both printed the note against a scratch home lacking the setting).

## Prevention
A rule that overrides a CLI's built-in behaviour must name the CLI's own setting for it and the command
that reports the setting missing — prose alone does not hold against a vendor-injected reminder. Before
writing a "never do X" rule about agent-emitted output, check whether X is actually a tool/vendor default
rather than a choice the agent is making.

## Related
Reviews: `.agents/docs/reviews/2026-09-16-review-bug-g1-attribution-note.md`,
`.agents/docs/reviews/2026-09-16-security-bug-g1-attribution-note.md`. Testing:
`.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`. Sandbox re-run owed (not yet run):
one start of each CLI at the fixed tree, per
`.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § Re-runs.
