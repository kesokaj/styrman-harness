---
owner: main
date: 2026-09-03
---
# Rename to sh-agent, ask-for-GCP rule, zero scripts, vendor format check     (date: 2026-09-03)

## Summary

Follow-up to the rebuild after the operator's review.

- Binary renamed `styrman` → `sh-agent` everywhere except history docs and the repo name: `cmd/sh-agent`, `~/.local/bin/sh-agent`, `~/.config/sh-agent/state.json`, `SH_AGENT_REPO`, the generated-file marker, README, AGENTS.md, REQUIREMENTS.md, skills that mention the tool.
- GCP rule: `context/dev-setup.md` no longer stores defaults. Work is local by default; when the user says GCP, the agent asks for project, region, Artifact Registry repo, service name, and visibility, then records them in the project's `AGENTS.md`. `context/AGENTS.md` operator line updated to match.
- Zero scripts: deleted `skills/mermaid-syntax/scripts/validate-mermaid.sh` (and its reference in the skill) and the old `~/.agents/bin` hook scripts with their backups on the Linux box. Executable bits removed from two `SKILL.md` files.
- Documentation-based format check of context, skills, and agents for both vendors (report in `reviews/`). Three fixes: `mermaid-syntax` `version` moved under `metadata`; `reviewer` `skills` written as a block list; the installer's frontmatter parser now reads YAML block lists.

## Type

task

## Files changed

- `cmd/sh-agent/main.go` (moved from `cmd/styrman`), `internal/harness/{frontmatter,plan,render,state,repo,apply,inventory}.go`, `internal/harness/harness_test.go` — rename, block-list parsing, test
- `context/AGENTS.md`, `context/dev-setup.md` — GCP rule, tool name
- `agents/reviewer.md` — `skills` block list
- `skills/mermaid-syntax/SKILL.md`, `skills/mermaid-syntax/scripts/` (deleted), `skills/{distill,find-skills,retro}/SKILL.md`, `skills/VENDORED.md` — script removal, tool name
- `README.md`, `AGENTS.md`, `REQUIREMENTS.md`, `TODO.md`, `.gitignore`, `.agents/docs/plans/2026-09-03-rebuild.md` — tool name, follow-up checklist
- `.agents/docs/architecture/DECISIONS.md` — ADR-006
- `.agents/docs/reviews/2026-09-03-vendor-format-check.md`, `.agents/docs/testing/2026-09-03-verify-rename.md` — new

## Verification

`gofmt`, `go vet`, `go test -race` clean; end-to-end install/status/uninstall of the renamed binary against a scratch home (244 files, no drift, clean removal). See `testing/2026-09-03-verify-rename.md`.

## Requirements

R2 (standards), R3 (no scripts beyond the binary), R4 (binary naming), R6 (scratch home).
