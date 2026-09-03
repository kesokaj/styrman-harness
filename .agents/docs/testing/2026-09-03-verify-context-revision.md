---
owner: main
date: 2026-09-03
---
# Verification — context revision, bugs agent, scratch   (2026-09-03)

Range reviewed: 204f0d9..working tree

## Status: PASS

## Commands

- `gofmt -l .` → no output; `go vet ./...` → clean; `go test -race ./...` → `ok` (adds `TestLastRetro`)
- `sh-agent install --home <scratch>` → `250 written, 0 unchanged, 0 removed, 0 kept — 250 files managed` (6 more than before: `bugs` agent on two vendors, `INDEX.md` and `scratch/.gitkeep` in the init-project assets on two vendors)
- `sh-agent status --home <scratch>` → `retro none yet — run /retro …`, `repo has 28 skills, 7 agents, 0 commands, 0 rules`, `drift none — 249 files match the repo`
- `ls ~/.claude/agents` and `~/.gemini/config/agents` → bugs, changelog, documentation, explorer, reviewer, security, verifier
- deployed `~/.claude/AGENTS.md` → 200 lines, 3,867 words, 13 sections
- `sh-agent uninstall --home <scratch>` → `removed 250 files and the manifest`
- `git ls-files | grep -E '\.(sh|py)$|/scripts/'` → nothing

## Coverage of changed code

- `internal/harness/state.go` `LastRetro` — `TestLastRetro`
- `cmd/sh-agent/main.go` status line — end-to-end run above
- context, agents, skills, assets — install succeeds, no drift, agents listed on both vendors

## Environment

Linux (Debian 13, amd64), Go 1.24.4.
