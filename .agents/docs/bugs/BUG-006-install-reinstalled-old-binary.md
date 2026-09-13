---
owner: bugs
date: 2026-09-13
id: BUG-006
status: fixed
area: installer/binary
---
# BUG-006: `sh-agent install` kept the operator's real `~/.local/bin/sh-agent` stuck on the 2026-09-03 build

## Symptom
The real home's `~/.local/bin/sh-agent` was still the 2026-09-03 21:05 build (no `--vendor` flag, `sh-agent version` printed `(devel)`, `state.json` had no `vendor` key) even though `install` had been run on 2026-09-08, 09-10 and 09-11 against a repo that had carried the `--vendor` flag since commit b16f5fb (2026-09-04).

## Repro
Recorded in `.agents/docs/testing/2026-09-13-verify-research-pass.md`: copy an old binary into a scratch home's `.local/bin/sh-agent`, run that binary's `install --home <scratch>` → `0 written, 256 unchanged`, md5 of the binary unchanged.

## Root cause
`install` copied the RUNNING executable (`selfExecutable()` in `cmd/sh-agent/main.go`) into `~/.local/bin`. When the running executable is already the installed copy, it reinstalls itself forever — the new source is never compiled in. Only `sh-agent update` (which ran `go run ./cmd/sh-agent install`) or an explicit `go run` actually rebuilt, but the README's daily loop and the retro/distill skills all told the operator to run plain `sh-agent install`. `go run` binaries also carry no VCS stamp, so `sh-agent version` said `(devel)` and nothing exposed the staleness.

## Fix
`install`/`sync` now always route through `harness.Reinstall`: it runs `go build` into `<repo>/.agents/scratch/sh-agent-build-<pid>/`, and that fresh build performs the install (guarded by env `SH_AGENT_CHILD=1` to prevent recursion). `status` now prints a `binary <rev>` line and warns `older than the repo` or `no build stamp`. Working tree today; changes entry `.agents/docs/changes/2026-09-13-installer-rebuilds-from-source.md`; plan `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` Phase 1.

## Prevention
A self-updating tool that copies "itself" runs its OLD code on the very tick meant to install the change — the same shape as styrman-crew's BUG-035 and the review-standards §2 "installers count as callers" paragraph. Any installer that copies its own running binary must instead build from source, and `status` must compare the running build's stamp against the repo. After any change to `cmd/` or `internal/`, run `sh-agent status` from the real home and read the `binary` line before calling the install done.

## Related
BUG-001 (status and the binary), review-standards §2 paragraph added 2026-09-10, styrman-crew BUG-035.
