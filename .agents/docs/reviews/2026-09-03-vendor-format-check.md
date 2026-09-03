---
owner: main
date: 2026-09-03
---
# Vendor format check — context, skills, subagents   (2026-09-03)

Range reviewed: the working tree after the rebuild (e75bf90 plus the rename and GCP-rule changes)

Documentation-based check that what `sh-agent install` deploys is in the shape each vendor documents. The live check on the operator's Mac is a later task.

## Verdict: PASS with three fixes applied

## Sources

| Vendor | Source | What it settled |
|---|---|---|
| Agent Skills spec | https://agentskills.io/specification | `SKILL.md` fields: `name` (1–64, lowercase, hyphens, **must match the directory name**), `description` (1–1024), `license`, `compatibility`, `metadata` (string map), `allowed-tools` (space-separated). Optional `scripts/`, `references/`, `assets/`. Keep `SKILL.md` under 500 lines. |
| Claude Code skills | https://code.claude.com/docs/en/skills.md | User-level `~/.claude/skills/<name>/SKILL.md`; command name comes from the directory; extra fields `disable-model-invocation`, `user-invocable`, `argument-hint`, `arguments`, `allowed-tools`, `disallowed-tools`, `model`, `effort`, `context`, `agent`, `background`, `hooks`, `paths`, `shell`, `when_to_use`; `metadata` is kept but not acted on; `description` + `when_to_use` are capped at 1,536 characters in the listing; `~/.claude/commands/*.md` still works but is merged into skills. |
| Claude Code subagents | https://code.claude.com/docs/en/sub-agents.md | User-level `~/.claude/agents/<file>.md`, identity from `name`; `tools` is a comma-separated string; `skills` is a YAML block list; `model` ∈ sonnet/opus/haiku/fable/inherit/full id; `memory` ∈ user/project/local; `isolation: worktree`; `effort` ∈ low…max; `maxTurns` integer; `description` single line. |
| Claude Code memory | https://code.claude.com/docs/en/memory.md (via the claude-code-guide agent) | `~/.claude/CLAUDE.md` supports `@path` imports, relative to the file, up to 4 hops; user-level imports load without a prompt. |
| Antigravity CLI 1.1.24 | Embedded customization guide and changelog in the `agy` binary (`agy changelog`; `strings` on the binary) | Global customization root `~/.gemini/config/` (skills at `skills/<name>/SKILL.md`, agents created by `/agents` there, `skills.json`/`plugins.json`, `mcp_config.json`, `hooks.json`); workspace root `.agents/` with `rules/*.md`, `skills/`, `agents/<name>/agent.md`; workflows deprecated in favour of skills; skills injected by name + description (progressive disclosure); markdown agents support `mainAgent`, `subagent`, `hidden`, `inheritMcp`, `inheritCustomizations`, `commandExecutionPolicy`, `model` (`flash`/`pro`/`inherit`), plus `tools`, `skills`, `rules`, `plugins`, `mcpServers` keys; markdown agents inherit ambient skills and rules by default; `disable-model-invocation` is a recognised skill key. |

## Skills (28)

Checked by script: every `name` equals its directory name and matches `[a-z0-9]+(-[a-z0-9]+)*`; every description is non-empty and under 1,024 characters; no `SKILL.md` exceeds 500 lines; every frontmatter key is either a spec field or a documented Claude Code field.

- **Fix:** `mermaid-syntax` had a top-level `version: 1.5.0`, which is not a spec field. Moved to `metadata: {version: "1.5.0"}`.
- `bash` and `docker` use `allowed-tools: Bash` (spec field, experimental; Claude honours it, Antigravity ignores unknown fields).
- `audit`, `distill`, `grill-me`, `handoff`, `init-project`, `research`, `retro` use `disable-model-invocation: true` + `argument-hint`; both vendors know the first, Antigravity ignores the second.
- Vendored skills (`anti-ui-slop`, `gke-basics`, `gcloud-auth-verification`, `skill-repair`, `terraform`, `codebase-design`, `grilling`, `grill-me`, `find-skills`) carry `license`/`metadata` per spec; untouched.
- Directory layout: `references/`, `assets/`, `examples/` only. No `scripts/` remain.

## Subagents (6)

- Frontmatter keys used: `name`, `description`, `tools`, `model`, `effort`, `maxTurns`, plus `memory` (changelog, documentation), `skills` (reviewer), `isolation` (verifier). All documented for Claude Code.
- **Fix:** `reviewer` had `skills: review-standards` as a scalar; the docs show a block list. Changed to `skills:\n  - review-standards`. The installer's frontmatter parser was taught YAML block lists so the inventory and derivation read both spellings.
- Descriptions are one line, double-quoted, 317–446 characters, free of unescaped colons.
- Antigravity derivation emits `name`, quoted `description`, `subagent: true`, `mainAgent: false`, `model` (`opus`→`pro`, `haiku`→`flash`, `sonnet` omitted = inherit), `tools` mapped to `view_file`, `grep_search`, `find_by_name`, `run_command`, `write_to_file`, `replace_file_content`; Claude-only keys dropped. `skills` is dropped because markdown agents inherit the global skills.
- The agy binary carries `description` YAML keys, so a frontmatter `description` is read; what its changelog calls "H1-delimited system prompts" could not be pinned down from the binary. The derived body keeps `##` headings only, so the whole body is one prompt. Left for the live check.

## Global context

- Claude Code: `~/.claude/CLAUDE.md` = marker line + `@AGENTS.md` → `~/.claude/AGENTS.md`. Documented import form, one hop, user level.
- Antigravity: `~/.gemini/config/rules/AGENTS.md`. The global root is confirmed; the `rules/` subdirectory is the workspace convention (`.agents/rules/*.md`) applied to the global root, and the plugin guide recommends `rules/AGENTS.md` as the consolidated rule file. Left for the live check. If it fails, the alternative is a plugin (`~/.gemini/config/plugins/sh-agent/{plugin.json,rules/AGENTS.md}`), which needs no registration when placed in the global root.
- Content: 104 lines deployed (context + environment + inventory), no vendor-specific instructions except the discovery note.

## What was not checked

- Live loading on either CLI (deliberately deferred by the operator).
- Antigravity IDE paths (out of scope).
- The exact agy tool-name list beyond the six mapped names.
