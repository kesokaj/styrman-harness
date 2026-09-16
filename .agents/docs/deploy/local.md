---
owner: main
date: 2026-09-15
target: local
---
# Deploy — local (this machine)

"Deploying" the harness means installing the context, skills, and agents into a home directory. Development never touches the real `$HOME` (`REQUIREMENTS.md` R6); deploying to it is the operator's action.

## Prerequisites

Go 1.24+ and git. `~/.local/bin` on `PATH`. Claude Code and/or the Antigravity CLI installed; both are optional — the installer writes for the vendor recorded by the existing install, else both, or selectively with `--vendor`.

## Build and install

```bash
cd ~/workspace/styrman-harness
gofmt -l . && go vet ./... && go test -race ./...   # before deploying anything
go run ./cmd/sh-agent install                        # builds the binary from this checkout, writes ~/.claude and ~/.gemini/config, installs the build to ~/.local/bin/sh-agent
# or for Antigravity only:
# go run ./cmd/sh-agent install --vendor antigravity
sh-agent status                                      # must end with "drift none"
```

## Configure

Nothing to configure. `--vendor VENDOR` selects the target CLI: `antigravity` (or `agy`) or `claude`; the default is the vendor recorded by the existing install, else `all` — passing `--vendor all` explicitly widens a single-vendor install to both through `install`; `update` keeps the recorded vendor instead. `--repo DIR` overrides where the checkout is; `$SH_AGENT_REPO` does the same persistently. `--home DIR` installs elsewhere and is what every development run uses (a scratch home under `.agents/scratch/`, never `/tmp`). From its own git checkout, `install` rebuilds the binary with `go build` first; the running binary is installed as is when `go` is missing (with a warning) or when the tree is a plain copy with no `.git` — a downstream mirror (BUG-015).

## Verify

- `sh-agent status` shows the current commit, a `binary` line carrying the same revision (a warning there means the installed binary is older than the checkout or unstamped: run `sh-agent update`), and no drift. `install` writes `"attribution": {"commit": "", "pr": ""}` in `~/.claude/settings.json` when the key is absent (creating the file when needed, every other key untouched) and prints a `settings:` line saying so; `status` prints an advisory note when the file is missing or the key is present with a different value.
- `sh-agent render | head -3` prints the assembled context.
- Start each CLI once and check it read the files: Claude Code `/memory` lists `~/.claude/CLAUDE.md` importing `AGENTS.md`, `/skills` and `/agents` list the harness entries; Antigravity `/skills` and `/agents` do the same, and `grep 'Failed to parse' ~/.gemini/antigravity-cli/log/cli-*.log` is empty (BUG-004).

## First session after install

Claude Code asks once per session before the first read or write under `~/.config/sh-agent/` (feedback log, proposals, the retro index), because it is outside the project. Approve it, or add the directory to `permissions.additionalDirectories` in `~/.claude/settings.json` yourself — the installer touches no vendor setting except the attribution key (ADR-025). Antigravity prompts the same way for paths outside the workspace.

## Operator steps

Standing by-hand procedures on this box, never pending items (those live in the newest handoff):

- After a push from another machine: `sh-agent update`, then `sh-agent status` shows a stamped `binary` line and no drift.
- After any change to `context/`, `agents/`, or `skills/` reaches this home: the three smoke checks of `.agents/docs/testing/2026-09-03-verify-real-home.md` on each CLI, from a fresh clone of the demo service.
- When `govulncheck ./...` names a toolchain-fixed advisory: upgrade Go with the system package manager (this box: apt), never by raising the `go` directive in `go.mod`.

## Roll back

`sh-agent uninstall` removes exactly the files the manifest lists. To return to an earlier harness commit: `git checkout <commit> && go run ./cmd/sh-agent install`.

## Sandbox instead of the real home

`.agents/docs/testing/2026-09-03-verify-sandbox.md` records how a throwaway `HOME` with both CLIs logged in is used to test a change end to end before it reaches the real home.
