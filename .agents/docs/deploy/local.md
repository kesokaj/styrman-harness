---
owner: main
date: 2026-09-03
target: local
---
# Deploy — local (this machine)

"Deploying" the harness means installing the context, skills, and agents into a home directory. Development never touches the real `$HOME` (`REQUIREMENTS.md` R6); deploying to it is the operator's action.

## Prerequisites

Go 1.24+ and git. `~/.local/bin` on `PATH`. Claude Code and/or the Antigravity CLI installed; both are optional, the installer writes for both.

## Build and install

```bash
cd ~/workspace/styrman-harness
gofmt -l . && go vet ./... && go test -race ./...   # before deploying anything
go run ./cmd/sh-agent install                        # writes ~/.claude and ~/.gemini/config, copies itself to ~/.local/bin/sh-agent
sh-agent status                                      # must end with "drift none"
```

## Configure

Nothing to configure. `--repo DIR` overrides where the checkout is; `$SH_AGENT_REPO` does the same persistently. `--home DIR` installs elsewhere and is what every development run uses.

## Verify

- `sh-agent status` shows the current commit and no drift.
- `sh-agent render | head -3` prints the assembled context.
- Start each CLI once and check it read the files: Claude Code `/memory` lists `~/.claude/CLAUDE.md` importing `AGENTS.md`, `/skills` and `/agents` list the harness entries; Antigravity `/skills` and `/agents` do the same, and `grep 'Failed to parse' ~/.gemini/antigravity-cli/log/cli-*.log` is empty (BUG-004).

## Roll back

`sh-agent uninstall` removes exactly the files the manifest lists. To return to an earlier harness commit: `git checkout <commit> && go run ./cmd/sh-agent install`.

## Sandbox instead of the real home

`.agents/docs/testing/2026-09-03-verify-sandbox.md` records how a throwaway `HOME` with both CLIs logged in is used to test a change end to end before it reaches the real home.
