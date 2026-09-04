---
owner: main
date: 2026-09-04
---
# Verification — install into the operator's real home (Mac, Antigravity only)   (2026-09-04)

Range reviewed: b16f5fb installed (`sh-agent status`: 128 files, no drift, vendor: antigravity)

## Status: PASS for Antigravity

Scope: Antigravity CLI 1.1.26 exercised from the real `$HOME` on macOS (Apple Silicon) following `--vendor antigravity` install. Claude Code was omitted as requested; `~/.claude/` was not created.

## Runs (real `$HOME`, cwd `/Users/spgo/Workspace/styrman-harness`)

| Check | Antigravity CLI (`agy`) | Result |
|---|---|---|
| Context questions (title, skills, subagents, pruned names) | `agy -p "What is the title of the global AGENTS.md file, and what are the names of your subagents?"`<br>`agy -p "List all 28 skills available in your context. Is grill-me or skill-repair present?"` | PASS — `# AGENTS.md — global agent context` read directly from `~/.gemini/config/rules/AGENTS.md`; all seven subagents listed; 28 skills listed; `grill-me` and `skill-repair` confirmed absent. |
| Skill auto-load (`golang` on `cmd/sh-agent/main.go`) | `agy -p "Review cmd/sh-agent/main.go according to the golang skill..."` | PASS — auto-loaded `golang` skill and its references (`~/.gemini/config/skills/golang/references/golang-io-web.md`), reviewed code against stdlib-first and CLI conventions. |
| Subagent dispatch (`explorer`) | `agy --agent explorer -p "Give a 3-sentence brief of this repository"` | PASS — `explorer` subagent ran in isolated context with tools enabled, inspected repository, and returned a concise 3-sentence brief. |
| Log inspection | `grep -i "failed to parse" ~/.gemini/antigravity-cli/log/cli-*.log` | PASS — zero frontmatter parse errors (BUG-004 prevention verified). |
| Vendor isolation | `test -d ~/.claude` | PASS — `~/.claude` does not exist; only `~/.gemini/config/` and `~/.local/bin/sh-agent` were written. |

## Commands

`go run ./cmd/sh-agent install --vendor antigravity` → `128 written, 0 unchanged, 0 removed, 0 kept — 128 files managed`; `sh-agent status` → `drift none — 127 files match the repo` (vendor: `antigravity`).

## Environment

macOS (Darwin arm64), Antigravity CLI 1.1.26, Go 1.24, harness b16f5fb.
