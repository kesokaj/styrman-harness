---
owner: bugs
date: 2026-09-03
id: BUG-002
status: fixed
area: cli/flags
---
# BUG-002: `sh-agent adopt <path> --home DIR` failed with a usage error

## Symptom
`sh-agent adopt /path/to/skill --home /tmp/x` printed `usage: sh-agent adopt <skill-name-or-path> [--force]` instead of adopting.

## Repro
Any subcommand with a positional argument followed by a flag.

## Root cause
Go's `flag` package stops parsing at the first non-flag argument, so `--home` became a second positional and `NArg()` was 3.

## Fix
`cmd/sh-agent/main.go`: `reorder` moves flags (and the values of `--repo`/`--home`) ahead of positionals before `Parse`.

## Prevention
When a command takes positionals, test it with flags after the positional; the stdlib parser does not interleave.

## Related
`.agents/docs/testing/2026-09-03-verify-rebuild.md`
