---
owner: main
date: 2026-10-08
id: BUG-164
status: fixed
area: internal/dream, internal/trail, internal/metrics (git in project roots)
---
# BUG-164: `git log` in a project root ran that repository's `gpg.program`

## Symptom
The security review of the metrics range (`reviews/2026-10-08-security-2c0f215-166aa39.md`, finding 1, medium) noted that a repository's own config can make `git log` start a program. The trigger is `log.showSignature=true` plus `gpg.program` (or `gpg.ssh.program`) and one signed commit. The nightly dreamer has run `git log` in every indexed project root since v2.0.0 (`trailDates`, `internal/dream/collect.go`), turning off only `core.fsmonitor`. The metrics store's events join added a second caller, in every root a transcript names, sandbox and A/B work dirs included.

## Repro
`TestDatesRunsNoSignatureProgram` in `internal/trail/trail_test.go` builds the repository:
- a commit signed through a fake `gpg.program`;
- `log.showSignature=true` in the repo config;
- the fake program writes a marker when asked to `--verify`.

The fixture first proves that a plain `git log` writes the marker. Before the fix, `Dates` wrote it too: the test was red ("git log ran the repository's gpg.program").

## Root cause
The hardening after BUG-135's class (2026-10-01: no index lock, no fsmonitor) covered `git status` and the flag order, but not the other config keys that make a read-only `git` command run a program. `log.showSignature` is one of them. Each caller built its own git argv, so there was no single place to harden.

## Fix
- One helper, `internal/trail`: `hardened(root)` gives `--no-optional-locks -c core.fsmonitor=false -c log.showSignature=false -C root`.
- `Log(root, …)` adds `--no-show-signature`.
- `StatusArgs` and the new `Toplevel` use the same options.
- `Dates` (the dreamer's watermark, and the metrics events) and `metrics.commitsIn` call `trail.Log`. No other `git log` runs in a project root (`internal/search` runs none).
- Tests: `TestDatesRunsNoSignatureProgram` and `TestLogRunsNoSignatureProgram`, both red before the fix. `TestTrailStatusArgs` pins the new argv.

**Round 2 of the same review (same class):** `git status`, which `Dates` used to tell modified records, runs a repository's clean filter (`.gitattributes` plus `filter.x.clean`), and a partial clone can lazy-fetch through its own `core.sshCommand`. The fix, by the operator's choice ("Fix both in this feature", 2026-10-09), covers the dreamer's own read too:
- `status` is gone. Modified means the working file's raw bytes (`git hash-object --no-filters --stdin-paths`) differ from HEAD's blob (`git ls-tree -r`). A filtered file reads as modified and keeps its mtime.
- Every call carries `GIT_NO_LAZY_FETCH=1` through one builder, `trail.command`.
- Tests: `TestDatesRunsNoCleanFilter` (its fixture forces status to hash through the filter by moving an mtime, and proves a plain status trips the marker) and `TestCommandsNeverLazyFetch`, both red before the fix.

## Prevention
A git call into a directory the operator did not just choose (a project root taken from the index or a transcript) goes through `internal/trail.command`, and nothing else. Only object- and file-reading commands run there: `log`, `ls-tree`, `hash-object --no-filters`, `rev-parse`. A new git command there gets a marker test against a repository whose config or attributes would run a program. Not covered: config keys that make other git commands run programs (`diff.external`, textconv, `core.pager` on a tty, hooks). Today the harness runs none of those commands in a project root.
