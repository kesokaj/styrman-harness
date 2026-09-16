---
owner: bugs
date: 2026-09-16
id: BUG-045
status: fixed
area: agents/security (probe repositories)
---
# BUG-045: failed probe-repo init fell through to the real checkout

## Symptom
During round 2 of the security review of the release procedure
(`.agents/docs/reviews/2026-09-16-security-release-procedure-adr-026.md`),
the agent's probe `git init -c …` failed under `set +e` and its later
commands ran in the real repository: two empty commits ("one" 5d5172e,
"two" 5815594) on `main` over 8b95f22, and two annotated unsigned tags
`v0.0.1`, `v0.0.2` — names matching the consumer's
`^v[0-9]+\.[0-9]+\.[0-9]+$` selection. Nothing pushed (origin/main stayed
8b95f22); working tree, index, and git config untouched. Detected only
because the agent reported it; the main agent's post-dispatch check
(`git status --short --untracked-files=all`, graduated from BUG-044 the
same day) was blind to commits and tags.

## Repro
Not re-run. Inferred from the debris: a `git init -c …` invocation whose
exit status was not checked (`set +e` in effect) followed by `git commit`
and `git tag` with no `-C <dir>`/`cwd` binding to the intended probe
directory, so they executed against whatever directory the shell was
already in — the real checkout.

## Root cause
Probe commands were not bound to the probe directory; a failed `git init`
left the shell in the checkout the agent was dispatched from, and every
following command ran there instead of erroring out.

## Fix
Main agent undid it after a guard (both commits empty, subjects
"one"/"two", exactly those two tags): `git reset --soft 8b95f22`,
`git tag -d v0.0.1 v0.0.2`; verified HEAD 8b95f22, 16 working-tree paths
unchanged, index empty.

## Prevention
Both halves written in the same commit as this record:
`context/AGENTS.md` § Hard rules, "Ask git what actually changed" bullet
now also runs `git log --oneline @{u}..HEAD` and `git tag` after any
dispatch that can write, treating a commit or tag the dispatch did not
name as a stop sign; `agents/security.md` now says a probe repository is
created with `git init <dir>` whose exit status is read, and every later
command names that directory (`git -C <dir>`).

## Related
Fourth of this shape after BUG-014, BUG-038, BUG-044 (an agent's own
scratch/isolation rule existed and was not followed, or in this case not
enforced against a failure path); this is the first where the
consequence reached the real branch history rather than the working
tree. Recommend the "bind every probe command to its own directory,
check init's exit status" rule stay in `agents/security.md` (it is
actionable — a specific `-C` binding — and generalizable to any agent
that inits a scratch git repo).
