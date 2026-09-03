---
owner: main
date: 2026-09-03
---
# Generic installer directories, attribution removed, hooks research     (date: 2026-09-03)

## Summary

- **Attribution removed.** The three pushed commits carried `Co-Authored-By: Claude …` and `Claude-Session:` trailers. Rewrote their messages without the trailers and force-pushed `main` with a lease. No trailers from now on.
- **Old traces deleted.** `~/.agents/` on the Linux box (feedback log, old changelog) is gone; the TODO items about old-harness leftovers are dropped.
- **Installer is generic.** Besides `skills/` and `agents/`, the repo now has `commands/` (alias `workflows/`) and `rules/`. A markdown file dropped into either deploys to both vendors (`~/.claude/commands|rules/<name>.md`, `~/.gemini/config/workflows|rules/<name>.md`) with no code change; the inventory in the deployed context lists them; `sh-agent status` counts them. `/init-project` also symlinks `.claude/rules`. `LoadInventory` replaces the separate loaders in the CLI.
- **Hooks researched, not built.** `.agents/docs/plans/2026-09-03-research-hooks.md` compares both vendors' hook systems and designs the deploy gate, the wrong-model warning, and the loop protector as `sh-agent hook` subcommands, with costs and five decisions for the operator.
- **Answered:** a project needs `AGENTS.md` plus a one-line `CLAUDE.md`; Antigravity reads `AGENTS.md` natively, so no `GEMINI.md`.

## Type

task

## Files changed

- `internal/harness/inventory.go`, `render.go`, `plan.go`, `harness_test.go`, `cmd/sh-agent/main.go` — `Inventory`, `LoadDocs`, commands and rules in the plan and the inventory
- `commands/.gitkeep`, `rules/.gitkeep` — the new directories
- `README.md`, `AGENTS.md`, `REQUIREMENTS.md` (R5), `context/AGENTS.md`, `skills/init-project/SKILL.md`, `TODO.md`, `.agents/docs/plans/2026-09-03-rebuild.md` — docs
- `.agents/docs/plans/2026-09-03-research-hooks.md` — new

## Verification

`gofmt`, `go vet`, `go test -race` clean (fixture now has a command and a rule; deployed paths and inventory rows asserted). End-to-end on a scratch home: 244 files, `repo has 28 skills, 6 agents, 0 commands, 0 rules`, no drift, clean uninstall. `git log --format=%B 20a5abd..HEAD` contains no trailer lines.

## Requirements

R5 (hand-added customizations), R3 (no hooks yet — research only).
