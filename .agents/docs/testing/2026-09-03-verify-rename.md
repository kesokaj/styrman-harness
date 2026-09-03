---
owner: main
date: 2026-09-03
---
# Verification — rename to sh-agent and format fixes   (2026-09-03)

Range reviewed: e75bf90..working tree

## Status: PASS

## Commands

- `gofmt -l .` → no output
- `go vet ./...` → clean
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/internal/harness 1.16s` (block-list frontmatter case added to `TestParseFrontmatter`)
- `go build -o <scratch>/sh-agent ./cmd/sh-agent`; `sh-agent install --home <scratch>` → `244 written, 0 unchanged, 0 removed, 0 kept — 244 files managed` (two fewer than before: the deleted script, once per vendor)
- `sh-agent status --home <scratch>` → `drift none — 243 files match the repo`
- `<scratch>/.local/bin/sh-agent` and `<scratch>/.config/sh-agent/` exist; `~/.claude/CLAUDE.md` is the marker line + `@AGENTS.md`; `~/.gemini/config/rules/AGENTS.md` is 104 lines
- Derived `reviewer` agent for Antigravity: quoted description, `subagent: true`, `mainAgent: false`, `model: pro`, mapped tools
- `sh-agent list` → 34 table rows (28 skills + 6 agents)
- `sh-agent uninstall --home <scratch>` → `removed 244 files and the manifest`; no files left under the scratch home
- Skill validator (Python, in the session): 28 skills, names match directories, descriptions ≤ 1024 characters, one non-spec key found and fixed
- `git ls-files | grep -E '\.(sh|py|…)$|scripts/'` → nothing; `git ls-files -s | awk '$1=="100755"'` → nothing

## Coverage of changed code

- `internal/harness/frontmatter.go` — `TestParseFrontmatter` (block list added)
- `cmd/sh-agent/main.go` — end-to-end run above
- context, agents, skills — install succeeds, status reports no drift

## Environment

Linux (Debian 13, amd64), Go 1.24.4, agy 1.1.24 (not signed in; live load deferred by the operator).
