# styrman-harness

A personal agent harness for **Claude Code** and the **Antigravity CLI** (`agy`): one global context file in the [agents.md](https://agents.md) standard, 28 skills in the [Agent Skills](https://agentskills.io) `SKILL.md` standard, six subagents, and `styrman` — a small Go binary that installs, resyncs, inspects, and removes all of it. No hooks, no gates, no CI.

## Install

```bash
git clone git@github.com:kesokaj/styrman-harness.git ~/workspace/styrman-harness
cd ~/workspace/styrman-harness
go run ./cmd/styrman install     # deploys everything and copies itself to ~/.local/bin/styrman
styrman status                   # should end with "drift none"
```

Make sure `~/.local/bin` is on your `PATH`. Requires Go 1.24+ and git; nothing else.

## Daily commands

| Command | What it does |
|---|---|
| `styrman status` | Shows what is installed, from which commit, and every file that drifted. Exit 1 on drift. |
| `styrman diff` | Unified diff of every drifted file (installed vs repo). |
| `styrman install` (or `sync`) | Deploys the repo. Idempotent; removes files a previous install wrote that the repo no longer has. |
| `styrman update` | `git fetch` the harness repo; if behind, `pull --ff-only`, rebuild, and reinstall. `--check` only reports (exit 2 when behind). |
| `styrman uninstall` | Removes exactly the files it installed (from its manifest) and nothing else. |
| `styrman list` | Prints the skill and agent inventory. |
| `styrman render` | Prints the assembled context file. |
| `styrman adopt <name\|path>` | Copies a skill installed elsewhere (for example by `npx skills add -g`) into `skills/`. |

Every command accepts `--home DIR` (install somewhere other than `$HOME`, used by the tests) and `--repo DIR`. Without `--repo`, styrman finds the checkout from `$STYRMAN_REPO`, the working directory, or the path it recorded at the last install — so `styrman update` works from anywhere.

## The loop: edit on one machine, pull on another

```bash
# machine A
vim context/dev-setup.md && styrman install && git commit -am "context: tweak" && git push
# machine B
styrman update      # fetches, pulls, rebuilds, reinstalls
```

## What gets installed where

| Source in this repo | Claude Code | Antigravity CLI |
|---|---|---|
| `context/AGENTS.md` + `context/dev-setup.md` + generated inventory | `~/.claude/AGENTS.md`, imported by a one-line `~/.claude/CLAUDE.md` (`@AGENTS.md`) | `~/.gemini/config/rules/AGENTS.md` |
| `skills/<name>/` | `~/.claude/skills/<name>/` | `~/.gemini/config/skills/<name>/` |
| `agents/<name>.md` (Claude format) | `~/.claude/agents/<name>.md` | `~/.gemini/config/agents/<name>/agent.md` (tool and model names translated at install) |
| the binary | `~/.local/bin/styrman` | same |
| manifest | `~/.config/styrman/state.json` | same |

An existing `~/.claude/CLAUDE.md` that styrman did not write is left alone; add the line `@AGENTS.md` to it yourself.

## Adding a skill by hand

Any directory with a `SKILL.md` dropped into `skills/` is a skill. From [skills.sh](https://skills.sh):

```bash
cd ~/workspace/styrman-harness
npx skills add vercel-labs/skills -s find-skills -y   # lands in skills/ via the .agents/skills symlink
styrman install
git add skills/find-skills skills-lock.json && git commit -m "skill: find-skills" && git push
```

Record third-party provenance in `skills/VENDORED.md`. A skill you already installed globally with `npx skills add -g` can be pulled in with `styrman adopt <name>`.

## Layout

```
context/     AGENTS.md (the global context) and dev-setup.md (environment, stack, GCP defaults)
agents/      six subagents in Claude Code format — the single source; Antigravity's flavour is derived
skills/      28 skills, each a directory with SKILL.md (+ references/, assets/, scripts/)
cmd/styrman  the binary; internal/harness is its logic and tests
.agents/     this repo's own trail (docs/) and a skills -> ../skills symlink for the skills CLI
```

## Working on the harness

```bash
gofmt -l . && go vet ./... && go test -race ./...
go run ./cmd/styrman install --home /tmp/styrman-test   # never test against your real $HOME
```

See `AGENTS.md` for the repo rules and `.agents/docs/plans/2026-09-03-rebuild.md` for the rebuild checklist.
