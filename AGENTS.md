# AGENTS.md — styrman-harness

Project rules for this repo. These override the global constitution, which this repo
happens to be the source of (`shared/AGENTS.md`).

## What this repo is

The agent harness: one constitution, six subagents, ~28 skills, and one Go binary that
installs and gates all of it, for both Claude Code and Google Antigravity.

## Layout

```
shared/AGENTS.md        the constitution — single source of truth
shared/skills/          skills, deployed to both vendors
claude/                 CLAUDE.md (generated), settings.json, agents/ (authored here)
antigravity/            GEMINI.md (generated), plugins/harness/ (agents generated)
cmd/harness/            the binary
internal/{trail,gate,hookio,install,dispatch}/
setup.sh                bootstrap only — runs `go run ./cmd/harness install`
```

## Generated files — never hand-edit

- `claude/CLAUDE.md` and `antigravity/GEMINI.md` are copies of `shared/AGENTS.md`.
- `antigravity/plugins/harness/agents/*.md` are derived from `claude/agents/*.md`.

Edit the source, then run `harness install` (or `go run ./cmd/harness install
--dry-run` to regenerate in-repo without deploying). CI fails on drift.

## Working on this repo

- `go test -race ./...` and `go vet ./...` before any completion claim.
- The installer must never be run against the operator's real `$HOME` during
  development. Use `--home /tmp/<scratch>` for every test.
- Adding a tool to an agent's `tools:` list means adding it to `toolMap` in
  `internal/install/derive.go`. An unmapped tool is a build error on purpose:
  Antigravity hangs a subagent that names a tool it does not recognise.

## Deployment is the operator's, not the agent's

This repo is deployed by hand. An agent working here changes repo files only — never
`~/.claude`, `~/.gemini`, or `~/.agents/bin`.
