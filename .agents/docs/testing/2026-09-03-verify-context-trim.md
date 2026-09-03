---
owner: main
date: 2026-09-03
---
# Verification — context trim, procedure skills, reviewer tiering   (2026-09-03)

Range reviewed: 5063167..working tree

## Status: PASS

## Commands

- `gofmt -l .` → no output; `go vet ./...` → clean; `go test -race ./...` → `ok` (inventory test now checks the compact context and the full `list` table separately)
- `sh-agent install --home <scratch>` → `254 written … 254 files managed` (two new skills on two vendors)
- `sh-agent status --home <scratch>` → `repo has 30 skills, 7 agents, 0 commands, 0 rules`, `drift none — 253 files match the repo`
- deployed `~/.claude/AGENTS.md` → 89 lines, 1,533 words (was 200 lines, 3,867 words): 67 lines of context + Environment & Stack + a names-only inventory
- `feature` and `bug` skills deployed without `disable-model-invocation` (model-invocable); derived `reviewer/agent.md` carries the new description and no `model:` line (Sonnet → inherit on Antigravity)
- `sh-agent uninstall --home <scratch>` → clean
- sandbox: `sh-agent install --home /home/kesokaj/sandbox/home` → 254 files, no drift; `/home/kesokaj/sandbox/demo-app` builds and its test passes (`go vet ./... && go test ./...`)

## Coverage of changed code

- `RenderInventoryCompact`, `RenderContext` — `TestInventoryAndRender`
- context, skills, agents — install succeeds, no drift, skills and agents counted on both vendors

## Not verified here

Live loading in either CLI: the sandbox is prepared for that and waits for the operator's login.

## Environment

Linux (Debian 13, amd64), Go 1.24.4.
