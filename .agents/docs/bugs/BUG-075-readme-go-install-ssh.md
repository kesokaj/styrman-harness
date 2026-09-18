---
owner: bugs
date: 2026-09-18
id: BUG-075
status: fixed
area: README.md § Install without cloning
---
# BUG-075: documented `go install` line for v1.1.0 fails without an SSH rewrite the reader was never told to add

## Symptom
README's "Install without cloning" section shipped this for v1.1.0:

    # or go install, over SSH, for v1.1.0 or newer (v1.0.0 has no harness built in)
    GOPRIVATE=github.com/kesokaj/* go install github.com/kesokaj/styrman-harness/cmd/sh-agent@v1.1.0

Run as written on a fresh machine it fails:
`go: github.com/kesokaj/styrman-harness/cmd/sh-agent@v1.1.0: invalid version: git ls-remote -q --end-of-options https://github.com/kesokaj/styrman-harness … fatal: could not read Username for 'https://github.com': terminal prompts disabled`

## Repro
```
GOPATH=<scratch> GOPRIVATE='github.com/kesokaj/*' \
  go install github.com/kesokaj/styrman-harness/cmd/sh-agent@v1.1.0
```
on a machine whose git has no `insteadOf` rewrite for github.com — reproduced on this box.

## Root cause
The comment says "over SSH" but the command gives git no way to use SSH: the module path is an https URL to a private repo, and git needs a `url."git@github.com:".insteadOf https://github.com/` rewrite (or stored https credentials) to authenticate. That rewrite existed only in the environment the instruction was authored in, so it was never missing when it was written. The release gate could not execute the command before the release — the tag it names does not exist until the release creates it — and the docs review over the release range read the command against `release.go`, where the code lives, not against a real `go install` run.

## Fix
README's Install without cloning section now gives a working one-liner that supplies the rewrite scoped to that single command via `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_0`/`GIT_CONFIG_VALUE_0` (changing no git config), with the permanent `git config --global url."git@github.com:".insteadOf` form mentioned as the reader's own choice. Verified today: `go install` succeeds, the resulting binary reports `sh-agent v1.1.0`, and installing from it writes 260 files with `status` showing `source embedded v1.1.0` and drift none. The other install path in that section (`gh release download`) was also run verbatim and works.

## Prevention
A command a doc tells the reader to RUN is executed as written, from the state the reader will be in, before the release that documents it — never approximated with different flags or read against the code instead. When it cannot be run before the release because it needs the tag the release itself creates, the release's Done row runs it immediately after and the doc is corrected in that same task; a docs review that can only read such a command says so under `## Not examined` rather than passing it.

## Related
BUG-030 (a detached script's setup had to be run once in the foreground before being trusted) — same family: a procedure is not verified by reading it. `.agents/docs/changes/2026-09-17-release-binary-phase-5.md`.
