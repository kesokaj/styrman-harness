# styrman-harness

A personal agent harness for **Claude Code** and the **Antigravity CLI** (`agy`): one global context file in the [agents.md](https://agents.md) standard, skills in the [Agent Skills](https://agentskills.io) `SKILL.md` standard, subagents, and `sh-agent` — a small Go binary that installs, resyncs, inspects, and removes all of it. No hooks, no gates, no CI.

## Install

```bash
git clone git@github.com:kesokaj/styrman-harness.git ~/workspace/styrman-harness
cd ~/workspace/styrman-harness
go run ./cmd/sh-agent install     # builds the binary from the checkout, deploys everything, installs the build to ~/.local/bin/sh-agent
sh-agent status                   # should end with "drift none"
```

Make sure `~/.local/bin` is on your `PATH`. Requires Go 1.24+ and git; nothing else. `install` writes Claude Code's attribution setting, `"attribution": {"commit": "", "pr": ""}` in `~/.claude/settings.json`, when the key is absent, so no commit carries a Co-Authored-By trailer; `status` prints a note when the file is missing or sets `commit` to anything else (see context § Git).

Harness state (feedback log, retro proposals and index) lives in `~/.config/sh-agent/`, outside any project, so Claude Code asks once per session before the first read or write there. To stop the prompt, add the directory yourself to `permissions.additionalDirectories` in `~/.claude/settings.json`; the installer touches no vendor setting except the attribution key (ADR-025).

## Daily commands

| Command | What it does |
|---|---|
| `sh-agent status` | Shows what is installed, from which commit, the running binary's build (with a warning when it is older than the checkout), a note when `~/.claude/settings.json` is missing or sets `attribution.commit` to anything but `""` (install writes the key when it is absent), and every file that drifted. Exit 1 on drift. |
| `sh-agent diff` | Unified diff of every drifted file (installed vs repo). |
| `sh-agent install` (or `sync`) | From its own git checkout, rebuilds the binary with `go build`, then deploys the repo; from a plain copy (no `.git`) deploys the running binary as is (`--vendor antigravity\|claude\|all`). Idempotent; removes files a previous install wrote that the repo no longer has. `--dry-run` prints what would change without changing anything. |
| `sh-agent update` | `git fetch` the harness repo; refuses a checkout with uncommitted changes, because it deploys the working tree rather than the fetched commit (BUG-046); if behind, `pull --ff-only`, rebuild, and reinstall. `--check` only reports (exit 2 when behind) and works on a dirty checkout. |
| `sh-agent uninstall` | Removes exactly the files it installed (from its manifest) and nothing else. `--dry-run` prints what would be removed without removing anything. |
| `sh-agent list` | Prints the skill and agent inventory. |
| `sh-agent render` | Prints the assembled context file. |
| `sh-agent adopt <name\|path>` | Copies a skill installed elsewhere (for example by `npx skills add -g`) into `skills/`. `--force` replaces an existing skill of the same name. |
| `sh-agent report` | Sends a redacted copy of `~/.config/sh-agent/` (feedback log, retro index, proposals, manifest) and a `report.json` metrics record to the private feedback repository, one directory per machine per run; clones that repository into `~/.config/sh-agent/feedback` on first use. `--dry-run` prints every file and redaction and writes nothing; `--feedback-repo DIR` names a checkout to use instead. See `.agents/docs/deploy/local.md` § Report back. |
| `sh-agent version` | Prints the build version. |

Every command accepts `--vendor VENDOR` (target `antigravity` / `agy` or `claude`; default is the vendor recorded by the existing install, else `all` — `--vendor all` widens a single-vendor install to both through `install`; `update` keeps the recorded vendor instead), `--home DIR` (install somewhere other than `$HOME`, used by the tests) and `--repo DIR`. Without `--repo`, sh-agent finds the checkout from `$SH_AGENT_REPO`, then the working directory or any of its parents, then the path it recorded at the last install — so `sh-agent update` works from anywhere.

## The loop: edit on one machine, pull on another

```bash
# machine A
vim context/dev-setup.md && sh-agent install && git commit -am "context: tweak" && git push
# machine B
sh-agent update      # fetches, refuses a dirty checkout, pulls, rebuilds, reinstalls
```

`sh-agent` is the only way the harness reaches a machine (R8): a checkout kept current with `update`, or an unpacked copy of the repository with no `.git`, installed with `go run ./cmd/sh-agent install` from that directory (it installs the running binary and `update` refuses it — unpack a newer copy to move forward). No project vendors a copy.

And back: `sh-agent report --dry-run` lists every file and every redaction it would make — what was caught, never what was missed, so before every report, not just the first, read what was added under `~/.config/sh-agent/` since the last one — then `sh-agent report`, so each machine's corrections and retro proposals reach the feedback repository where the harness is improved from them.

## What gets installed where

| Source in this repo | Claude Code | Antigravity CLI |
|---|---|---|
| `context/AGENTS.md` + `context/dev-setup.md` + generated inventory (names only; `sh-agent list` has the descriptions) | `~/.claude/AGENTS.md`, imported by a one-line `~/.claude/CLAUDE.md` (`@AGENTS.md`) | `~/.gemini/config/rules/AGENTS.md` |
| `skills/<name>/` | `~/.claude/skills/<name>/` | `~/.gemini/config/skills/<name>/` |
| `agents/<name>.md` (Claude format) | `~/.claude/agents/<name>.md` | `~/.gemini/config/agents/<name>/agent.md` (tool and model names translated at install) |
| the binary | `~/.local/bin/sh-agent` | same |
| manifest, feedback log, retro proposals, `retros.jsonl` (cross-project retro/distill index), and the feedback checkout `sh-agent report` clones (`feedback/`) | `~/.config/sh-agent/` | same |

An existing `~/.claude/CLAUDE.md` that sh-agent did not write is left alone. Move only rules meant for every machine into a rule file under `rules/` in this repository — `rules/` is deployed to every install on every machine, so credentials and personal or host-specific notes never go there — then remove `~/.claude/CLAUDE.md` and run `sh-agent install`, which writes the one-line import (`@AGENTS.md`).

## Adding things by hand — no code

Every customization is a file in a directory. Drop it in, run `sh-agent install`, commit. Nothing to register. The installer refuses a `SKILL.md` without a `description`, and any frontmatter value that a strict YAML parser would read as a comment (an unquoted `#` at the start of the value, or after a space or tab) — Antigravity silently drops such a skill, so quote it.

| Put it in | What it is | Claude Code gets | Antigravity gets |
|---|---|---|---|
| `skills/<name>/SKILL.md` | a skill (Agent Skills standard); the model loads it when the description matches, and `/name` types it | `~/.claude/skills/<name>/` | `~/.gemini/config/skills/<name>/` |
| `agents/<name>.md` | a subagent in Claude Code format (`name`, `description`, `tools`, `model`, …) | `~/.claude/agents/<name>.md` | `~/.gemini/config/agents/<name>/agent.md`, tools and model translated |
| `commands/<name>.md` (or `workflows/`) | a plain markdown prompt invoked as `/name`; frontmatter `description` optional | `~/.claude/commands/<name>.md` | `~/.gemini/config/workflows/<name>.md` |
| `rules/<name>.md` | an always-on or path-scoped rule (Claude: `paths:` frontmatter; Antigravity: `trigger: glob` + `globs:`; both ignore the other's keys) | `~/.claude/rules/<name>.md` | `~/.gemini/config/rules/<name>.md` |
| `context/<name>.md` | more global context, appended after `AGENTS.md` | part of `~/.claude/AGENTS.md` | part of `~/.gemini/config/rules/AGENTS.md` |

A skill is the portable way to write a command: both vendors run any skill as `/name` without translation, and the model can load it itself when the description matches — every harness skill works both ways, so the agent can run `/handoff` or `/retro` for you. Add `disable-model-invocation: true` only to a skill that must never fire on its own (none of the harness's own do). Prefer a skill over `commands/` for anything new.

### From skills.sh

Any directory with a `SKILL.md` dropped into `skills/` is a skill. From [skills.sh](https://skills.sh):

```bash
cd ~/workspace/styrman-harness
npx skills add vercel-labs/skills -s find-skills -y   # lands in skills/ via the .agents/skills symlink
sh-agent install
git add skills/find-skills skills-lock.json && git commit -m "skill: find-skills" && git push
```

Record third-party provenance in `skills/VENDORED.md`. A skill you already installed globally with `npx skills add -g` can be pulled in with `sh-agent adopt <name>`.

## Layout

```
context/     AGENTS.md (the global context) and dev-setup.md (environment, stack, GCP rule)
agents/      subagents in Claude Code format — the single source; Antigravity's flavour is derived
skills/      skills, each a directory with SKILL.md (+ references/, assets/); `feature` and `bug` are the work protocols and load themselves; `review` is the review command on both CLIs
commands/    markdown prompts invoked as /name (empty; prefer skills)
rules/       extra always-on or path-scoped rules (empty)
cmd/sh-agent  the binary; internal/harness is its logic and tests
.agents/     this repo's own trail (docs/, incl. bugs/INDEX.md), scratch/ (gitignored), and a skills -> ../skills symlink for the skills CLI
```

## Working on the harness

```bash
gofmt -l cmd internal && go vet ./... && go test -race ./...   # gofmt takes the source dirs: unlike the Go tool it walks .agents/scratch
# govulncheck: clean on go1.25.8+; an older Go reports two stdlib advisories, fixed by a newer Go on the machine, never a pin here (ADR-029)
go run ./cmd/sh-agent install --home .agents/scratch/home   # never test against your real $HOME; scratch lives in the repo, never /tmp
```

See `AGENTS.md` for the repo rules and `.agents/docs/plans/2026-09-03-rebuild.md` for the rebuild checklist. To cut a release, see `.agents/docs/deploy/remote.md` § Releases — signed `vX.Y.Z` tags via `/preflight`. A tag is a release marker only; every machine follows `main` with `sh-agent update` (ADR-030). How the harness reaches each machine and reports back: `.agents/docs/deploy/remote.md` § The four places.
