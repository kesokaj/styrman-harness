# Harness Portability Map (Phase 0)

Evaluation Date: **2026-09-03**  
Target Environments: **Claude Code** and **Google Antigravity / Gemini CLI**.  
Core Principle: Portable core (Agent Skills spec, `AGENTS.md`, `.agents/` output contract, Go binary gates) + thin per-vendor adapters.

---

## 1. Mechanism Cross-Vendor Mapping

| Mechanism | Open Standard | Claude Code Implementation | Gemini CLI / Antigravity Implementation | Portable Core Location | Adapter Location |
|---|---|---|---|---|---|
| **Constitution / Context** | `AGENTS.md` | `claude/CLAUDE.md` (generated verbatim) | `antigravity/GEMINI.md` (generated verbatim) | `shared/AGENTS.md` | `claude/`, `antigravity/` |
| **Project Overrides** | `AGENTS.md` | `AGENTS.md` (repo root) | `AGENTS.md` (repo root) | `AGENTS.md` (root) | None (both read natively) |
| **Skills** | Agent Skills (`SKILL.md`) | `~/.claude/skills/` | `~/.gemini/config/skills/` | `shared/skills/` | Mirrored by `harness install` |
| **Subagents** | Isolated agent spec | `claude/agents/*.md` | `antigravity/plugins/harness/agents/*.md` | `claude/agents/*.md` (source of truth) | `internal/install/derive.go` |
| **Hooks / Lifecycle** | Event hooks | `claude/settings.json` | `antigravity/plugins/harness/hooks.json` | `internal/gate/` (shared Go logic) | `claude/settings.json`, `hooks.json` |
| **Artifact Trail** | `.agents/docs/` | `.agents/docs/` | `.agents/docs/` | `.agents/docs/` | None (universal contract) |
| **Gate Enforcement** | Single binary | `~/.agents/bin/harness` | `~/.agents/bin/harness` | `cmd/harness/`, `internal/` | None (compiled binary) |

---

## 2. Core vs Adapter Architecture

```
                       ┌───────────────────────────────┐
                       │         PORTABLE CORE         │
                       │   shared/AGENTS.md            │
                       │   shared/skills/              │
                       │   .agents/docs/ (artifacts)   │
                       │   cmd/harness/ & internal/    │
                       └──────────────┬────────────────┘
                                      │
                   ┌──────────────────┴──────────────────┐
                   ▼                                     ▼
        ┌───────────────────────┐             ┌───────────────────────┐
        │  CLAUDE CODE ADAPTER  │             │ ANTIGRAVITY ADAPTER   │
        │  claude/CLAUDE.md     │             │ antigravity/GEMINI.md │
        │  claude/settings.json │             │ plugin.json           │
        │  claude/agents/*.md   │             │ hooks.json            │
        │  ~/.claude/skills/    │             │ derived agents/*.md   │
        │  ~/.claude/agents/    │             │ config/plugins/       │
        └───────────────────────┘             └───────────────────────┘
```

### Analysis of Current Split:
1. **Constitution (`shared/AGENTS.md`):** Perfectly separated. `shared/AGENTS.md` is the single authoring point; `claude/CLAUDE.md` and `antigravity/GEMINI.md` are generated copies. CI and `harness doctor` enforce zero drift.
2. **Skills (`shared/skills/`):** Mirrored to both `~/.claude/skills/` and `~/.gemini/config/skills/`. Almost entirely portable, with minor exceptions:
   - `grill-me` invokes Claude-specific `Skill` tool.
   - Non-standard subdirectories in `anti-ui-slop` (`reference/` vs `references/`) and `codebase-design` (root markdown files).
3. **Subagents:** Authored in `claude/agents/*.md` and derived for Antigravity by `internal/install/derive.go`.
   - Tool names mapped (`Read` -> `view_file`, etc.).
   - Model aliases mapped (`opus` -> `pro`, `haiku` -> `flash`).
   - `sonnet` currently unmapped (defaults to `inherit`).
4. **Gates and Hooks:**
   - Unified Go binary executes identical gate checks (`commitBoundary`, `SubagentStopDecision`, `CredentialGate`).
   - Claude adapter registers hooks in `settings.json`.
   - Antigravity adapter registers hooks in `hooks.json`.

---

## 3. Duplication and Drift Audit

1. **Drift Risk in Generated Files:**
   - The repository guarantees zero drift between `shared/AGENTS.md`, `claude/CLAUDE.md`, and `antigravity/GEMINI.md` via `harness doctor` and `harness install`.
2. **Duplicated Prose in Agent Prompts:**
   - Both `reviewer.md` and `verifier.md` duplicate identical shell snippets for resolving the docs root:
     ```bash
     if [ -d .agents/docs ]; then docs_root=.agents/docs
     elif [ -d docs ]; then docs_root=docs
     else docs_root=.agents/docs; fi
     ```
     This can be standardized or simplified since `.agents/docs` is the established root.
3. **Vendor-Specific Tooling in Skills:**
   - 4 skills contain `agents/openai.yaml` (`anti-ui-slop`, `codebase-design`, `grill-me`, `grilling`), which is an artifact from OpenAI-targeted packaging and inert on both Claude Code and Antigravity.
   - `anti-ui-slop` contains npm/git packaging files (`MANIFEST.json`, `CHECKSUMS.sha256`, `MODIFICATIONS.md`, `NOTICE`, `LICENSE`) which add metadata noise without functional value.

---

## 4. Portability Improvements & Plan

1. **Normalize Skill Layouts to Agent Skills Living Spec (2026-08-15):**
   - Consolidate `anti-ui-slop`'s `reference/` into `references/`.
   - Move `codebase-design`'s `DESIGN-IT-TWICE.md` and `DEEPENING.md` into `references/`.
   - Prune dead packaging metadata (`agents/openai.yaml`, checksums) to reduce repository surface and token leakage.
2. **Explicit Antigravity Model Mapping:**
   - Update `internal/install/derive.go` to explicitly handle `sonnet` (e.g., mapped to `inherit` or `pro` with documented rationale).
3. **Cross-Vendor Graceful Fallbacks in Skills:**
   - Update `grill-me` with instructions that work across both Claude Code (`/grill-me`) and Antigravity (delegating cleanly to `grilling`).
