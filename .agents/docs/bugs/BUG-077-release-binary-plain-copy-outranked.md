---
owner: bugs
date: 2026-09-18
id: BUG-077
status: fixed
area: internal/harness/repo.go (ResolveSource source order), delivery model
---
# BUG-077: release binary silently installed a stale plain-copy tree instead of its embedded harness

## Symptom
On `argolis-devbox` (hostname `devbox`), which had an old unpacked-bundle harness at
`/opt/devbox/deploy/harness` recorded in its manifest, the v1.2.0 release binary installed
that old tree's content instead of the harness built into it, and reported nothing wrong:
`plain copy, not its own git checkout: the running binary is installed as is`,
`0 written, 260 unchanged`, and `status` showing `repo /opt/devbox/deploy/harness (plain
copy…)` with `drift none — 259 files match the repo` — measured against the stale tree, not
the release. Confirmed stale: asked to quote the destructive-action guard rule, that
machine's Claude Code quoted the `rm -r` rule, not the fail-closed rule graduated into
`context/AGENTS.md` on 2026-09-17 and shipped in v1.2.0. The machine was also stuck: `update`
refused with "not its own git checkout; a plain copy cannot be updated with git", and there
was no flag to say "use what you carry".

## Repro
In `.agents/scratch/devbox-repro`: write a `state.json` whose `repo` points at a
harness-shaped directory with no `.git`, then run a v1.2.0 release binary's
`install --home <that home>` from a directory outside any tree. It installs from the
recorded copy and `status` reports it as the repo.

## Root cause
`ResolveSource`'s source order (D4 in `.agents/docs/plans/2026-09-17-release-binary.md`:
explicit path, `$SH_AGENT_REPO`, working directory and its parents, recorded path, embedded
files) was written for a development build, which has nothing of its own to install. Applied
unchanged to a release binary — which carries the published repository at the tagged commit —
it let any tree on the machine, including a dead-end plain copy, silently outrank the thing
the operator downloaded and verified.

## Fix
`internal/harness/repo.go` `ResolveSource`: a release binary (`IsReleaseVersion(version)`)
now installs only the harness built into it; an explicit `--repo` or `$SH_AGENT_REPO` is
REFUSED with a message naming the release and the alternative, never silently ignored.
Development builds keep the old lookup order. ADR-034. Verified with a release-shaped binary
in four situations: the devbox shape, standing inside a plain copy, standing inside the
operator's real checkout, and `--repo` at a real tree — the first three install the embedded
files, the fourth is refused.

## Prevention
When a program can obtain the same input from several sources, decide the order per KIND of
build, not once for all of them — an artifact that ships its own copy of an input must never
silently prefer a copy it found on the machine, and a source that cannot be updated (a plain
copy) is never preferred over one that can. State which source won whenever more than one was
possible and the choice changes what the tool then deploys.

## Related
ADR-034; `.agents/docs/plans/2026-09-17-release-binary.md` D4, D14, D15; BUG-046 (dirty-tree
message), BUG-006 (Reinstall path).
