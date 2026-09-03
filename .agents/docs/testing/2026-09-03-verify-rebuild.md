---
owner: main
date: 2026-09-03
---
# Verification — rebuild   (2026-09-03)

Range reviewed: 20a5abd..working tree (the whole rebuild, uncommitted at the time of the run)

## Status: PASS

## Commands

- `gofmt -l .` → no output
- `go vet ./...` → clean
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/internal/harness 1.15s` (8 tests: frontmatter parsing, split errors, derive, inventory/render, install→status→resync→uninstall, own CLAUDE.md stub, adopt, find-repo, git behind/pull against a local bare remote)
- `go run ./cmd/styrman install --home <scratch>` → `246 written, 0 unchanged, 0 removed, 0 kept — 246 files managed`
- `styrman status --home <scratch>` → `drift none — 245 files match the repo`, exit 0 (245 = 246 minus the binary, which status does not compare)
- append a line to `<scratch>/.claude/skills/golang/SKILL.md`; `styrman status` → `1 differ`, exit 1; `styrman diff` → unified diff with `-# tampered`; `styrman install` → `2 written, 244 unchanged`; `status` → `drift none`
- `styrman adopt <path> --home <scratch>` → `adopted into …/skills/my-skill`; second run without `--force` → refused, exit 1; with `--force` → replaced
- `styrman version` → `styrman 20a5abde73bc+dirty`
- `styrman update --check --home <scratch>` → `up to date with upstream (20a5abd)`, exit 0 (real fetch from origin over ssh)
- `styrman uninstall --dry-run` → `[dry-run] would remove 246 files and the manifest`; `styrman uninstall` → `removed 246 files and the manifest`; left behind: only `~/.claude`, `~/.config`, `~/.gemini/config`, `~/.local/bin` (empty directories)
- `styrman list` → 28 skills, 6 agents; derived Antigravity agent for `reviewer` has `model: pro`, `tools: [view_file, grep_search, find_by_name, run_command, write_to_file]`, quoted description

## Bugs found by this run, fixed before it was repeated

1. `status` reported `~/.local/bin/styrman` as stale because status plans without the binary. Fixed in `internal/harness/apply.go` (BinPath excluded from the stale check) with a regression test.
2. `styrman adopt <path> --home DIR` failed with a usage error: the flag package stops at the first positional argument. Fixed in `cmd/styrman/main.go` (`reorder` moves flags ahead of positionals).

## Coverage of changed code

- `internal/harness/*.go` — covered by `internal/harness/harness_test.go`
- `cmd/styrman/main.go` — no unit test; covered by the end-to-end run above (every command exercised)
- `context/`, `agents/`, `skills/` — verified by the install succeeding, `status` reporting no drift, and the rendered inventory listing 28 skills and 6 agents

## Not verified here

- That the Antigravity CLI reads `~/.gemini/config/rules/AGENTS.md`, `agents/<name>/agent.md`, and `skills/` from the global root: agy is not signed in on this box. Operator task on the Mac.
- That Claude Code loads `~/.claude/AGENTS.md` through the `@AGENTS.md` import: documented behaviour, not exercised, because the installer was deliberately not run against the real home.

## Environment

Linux (Debian 13, amd64), Go 1.24.4, git 2.x, agy 1.1.24 (not signed in), Docker present but unused.
