---
owner: bugs
date: 2026-09-16
id: BUG-038
status: fixed
area: agents (read-only agents and git)
---
# BUG-038: a read-only agent ran a destructive git command against the tracked working tree

## Symptom
Dispatched to verify the uncommitted BUG-036 fix, `verifier` ran
`git checkout -- internal/harness/inventory.go internal/harness/render.go internal/harness/pathguard_test.go`
to reproduce the pre-fix failure in place. Its own log records: "reproducing this
required temporarily reverting the working-tree fix with `git checkout -- …`, which —
because no stash was taken first — discarded the uncommitted fix from the working
tree." It restored the content from a captured diff plus a saved copy of the test
file and confirmed `git diff HEAD` was byte-identical afterward
(`.agents/docs/testing/2026-09-16-verify-bug-036-symlink-reads.md` lines 95-105).
The main agent independently confirmed no work was lost: three files modified,
`readRegular` defined once with its five call sites, gofmt/vet/`go test -race ./...`
clean.

## Repro
Dispatch `verifier` over an uncommitted change and ask it to confirm a new test
fails without the fix. The direct way to do that in place is `git checkout -- <file>`,
which discards the change with no edit tool involved.

## Root cause
`agents/verifier.md` says the agent never writes or edits tests or source, and the
dispatch said to reproduce in `.agents/scratch/` only; neither names git. "Never
edits source" reads as "never open an editor tool", while `git checkout --`,
`git restore`, `git stash` and `git clean` rewrite tracked files without any edit
tool call. `git stash` (with no pop) would have been just as lossy here. The same
gap exists in `agents/reviewer.md`, `agents/security.md` and `agents/explorer.md` —
all read-only in the same sense and all have Bash.

## Fix
`agents/verifier.md`, `agents/reviewer.md`, `agents/security.md` and `agents/explorer.md`
each gained one line in `## Records`, after "Never: source, tests, config, or another
agent's records": never a git command that changes the working tree —
`git checkout --`, `restore`, `stash`, `clean`, `reset`, `apply` — even to reproduce an
earlier state; copy what is needed into `.agents/scratch/` and change the copy, or read
old content with `git show <rev>:<path>`. `bugs`, `changelog` and `documentation` are
left out — they have no Bash and cannot run git. `context/AGENTS.md:45` carries the same
rule as a clause on the Git paragraph, so every subagent on either CLI reads it from the
context too; rendered context stays at 93 lines. `sh-agent list` parses 27 skills and 7
agents; scratch installs show drift none with the new line and clause present in both
vendors' copies. In the sandbox, one start of each CLI quoted the agent line from its
installed copy and one start of each quoted the context sentence, both verbatim, with
zero `Failed to parse` in the Antigravity logs. Code review
(`.agents/docs/reviews/2026-09-16-review-bug-038-code.md`) APPROVE with no findings;
docs review (`…-bug-038-docs.md`) APPROVE WITH NITS, its medium finding fixed by the
context clause. Changes entry:
`.agents/docs/changes/2026-09-16-bug-038-agents-git.md`.

## Prevention
An agent that does not own the working tree never runs a git command that changes
it — `checkout --`, `restore`, `stash`, `clean`, `reset`, `apply` — even to reproduce
a prior state. To compare against another revision it copies what it needs into
`.agents/scratch/` and reverts the copy, or reads old content with
`git show <rev>:<path>`. Reading commands (`status`, `diff`, `log`, `show`, `grep`)
stay allowed. This belongs in each read-only agent definition (`verifier`,
`reviewer`, `security`, `explorer`) and is a candidate for `context/AGENTS.md` where
it describes what a subagent may do.

## Related
BUG-036 is the fix that was discarded and restored here. BUG-013 and BUG-029 share
the theme that an agent's own account of what it did is checked, not assumed. Three
records now touch git-against-tracked-tree safety (BUG-011, BUG-036's own near-miss,
this one) — recommend the rule above land in `context/AGENTS.md`'s subagent section
as a global hard rule, not per-agent text; it is actionable (grep each agent
definition and dispatch prompt for git write-commands) and generalizable beyond this
repo.
