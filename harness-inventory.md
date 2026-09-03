# Harness Inventory (Phase 0)

Research snapshot & runtime evaluation: **2026-09-03**.  
Harness: **styrman-harness** (Dual-target: Anthropic Claude Code & Google Antigravity / Gemini CLI).

---

## 1. Executive Summary: Tokens at Turn 0

Turn 0 represents the token tax paid before the user issues their first task prompt in a clean session.

| Model / Target | Fixed Context (AGENTS + Rules) | Skills Metadata (all registered) | Subagent Definitions (registered) | Tool Definitions | Total Fixed Turn 0 Overhead |
|---|---|---|---|---|---|
| **Claude Code** (Claude 3.5 Sonnet / Opus / Haiku) | 2,593 tok (`CLAUDE.md`) + 455 tok (Repo `AGENTS.md`) = 3,048 tok | 2,773 tok (28 skills frontmatter) | ~1,200 tok (subagent descriptions) | ~1,850 tok (native tools) | **~8,871 tokens** (~4.4% of 200k) |
| **Google Antigravity** (Gemini 2.5 / 3.0 / 3.8 Pro/Flash) | 2,593 tok (`GEMINI.md`) + 455 tok (Repo `AGENTS.md`) = 3,048 tok | 2,773 tok (28 skills frontmatter) | ~1,100 tok (plugin agent metadata) | ~2,100 tok (16 native tools) | **~9,021 tokens** (~0.9% of 1M) |
| **Resume Protocol / Cold Start** (when `harness resume` or trail is read) | + 2,161 tok (`REQUIREMENTS.md`) + 176 tok (`ROADMAP.md`) + 195 tok (`TODO.md`) + 2,785 tok (`DECISIONS.md`) + ~1,800 tok (3 newest changes) = **+7,117 tok** | — | — | — | **~15,988 – 16,138 tokens** |

---

## 2. Context Files Inventory

| Path | Lines | Est. Tokens | Bytes | Load Trigger | Survives Compaction? | Standard vs Tool-Specific | Notes |
|---|---|---|---|---|---|---|---|
| `shared/AGENTS.md` | 78 | 2,593 | 10,499 | Single source of truth | Yes | Standard (`agents.md`) | Canonical global constitution; authored here |
| `claude/CLAUDE.md` | 78 | 2,593 | 10,499 | Session start | Yes | Tool-specific name (`CLAUDE.md`) | Byte-for-byte generated copy of `shared/AGENTS.md` |
| `antigravity/GEMINI.md` | 78 | 2,593 | 10,499 | Session start | Yes | Tool-specific name (`GEMINI.md`) | Byte-for-byte generated copy of `shared/AGENTS.md` |
| `AGENTS.md` (root) | 43 | 455 | 1,820 | Session start | Yes | Standard (`agents.md`) | Project-specific rules & generated file rules |
| `REQUIREMENTS.md` | 152 | 2,161 | 8,644 | On demand / Resume Protocol | Yes (if in trail) | Standard Markdown | Requirements R1–R8, known limits, non-goals |
| `ROADMAP.md` | 19 | 176 | 706 | On demand / Resume Protocol | Yes (if in trail) | Standard Markdown | Now / Next / Later |
| `TODO.md` | 15 | 195 | 782 | On demand / Resume Protocol | Yes (if in trail) | Standard Markdown | Deferred work items |

---

## 3. Rules & Scoped Context

| Scope | Path | Status | Standard vs Tool-Specific |
|---|---|---|---|
| Path-scoped rules | `.claude/rules/*.md` / `.gemini/rules/` | None currently defined | Standard format |
| Repository overrides | `.agents/harness.json` | 22 lines (581 bytes), loaded by `harness` binary | Harness custom configuration |

---

## 4. Skills Inventory (28 Skills in `shared/skills/`)

All skills implement the Agent Skills open standard (`SKILL.md` with YAML frontmatter).

| Skill Name | Metadata Tok | Body Lines | Body Tok | Reference Tok | Scripts/Assets | Provenance / Pin | Standard vs Tool-Specific |
|---|---|---|---|---|---|---|---|
| `anti-ui-slop` | 104 | 47 | 810 | 17,450 | Manifest + Checksums | Vendored: UIZZE v1.2.13 (sha256) | Standard (has non-standard `reference/` dir) |
| `audit` | 92 | 49 | 809 | 0 | None | First-party | Standard |
| `bash` | 65 | 60 | 686 | 6,840 | None | First-party | Standard |
| `cloud-run-deploy` | 145 | 98 | 1,306 | 0 | None | First-party | Standard |
| `codebase-design` | 105 | 109 | 1,516 | 0 | 2 root md files | Vendored: mattpocock/skills (6654f6b) | Standard (has root `DEEPENING.md` & `DESIGN-IT-TWICE.md`) |
| `distill` | 82 | 156 | 1,608 | 0 | None | First-party | Standard |
| `docker` | 115 | 156 | 1,270 | 4,210 | None | First-party | Standard |
| `feedback` | 102 | 78 | 1,167 | 0 | None | First-party | Standard |
| `find-skills` | 89 | 136 | 1,364 | 0 | None | Vendored: vercel-labs/skills (435076e) | Standard |
| `frontend` | 128 | 151 | 1,844 | 12,850 | None | First-party | Standard |
| `gcloud-auth-verification` | 78 | 40 | 495 | 0 | None | Vendored: Google Apache-2.0 | Standard |
| `gcloud-setup` | 112 | 125 | 888 | 0 | None | First-party | Standard |
| `gke-basics` | 118 | 46 | 1,218 | 38,400 | 5 YAML assets | Vendored: Google Cloud Apache-2.0 | Standard |
| `golang` | 74 | 62 | 728 | 21,500 | None | First-party | Standard |
| `grill-me` | 42 | 1 | 39 | 0 | None | Vendored: mattpocock/skills (6654f6b) | Tool-specific (invokes Claude `Skill` tool) |
| `grilling` | 76 | 23 | 493 | 0 | None | Vendored: mattpocock/skills (6654f6b) | Standard |
| `handoff` | 116 | 179 | 1,889 | 0 | None | First-party | Standard |
| `init-project` | 98 | 231 | 1,546 | 0 | None | First-party | Standard |
| `json-api-style` | 102 | 45 | 626 | 5,200 | None | First-party | Standard |
| `mermaid-syntax` | 134 | 268 | 2,050 | 18,900 | 1 script, 2 examples | First-party | Standard |
| `preflight` | 154 | 79 | 2,510 | 0 | None | First-party | Standard |
| `python` | 92 | 57 | 821 | 8,900 | None | First-party | Standard |
| `research` | 75 | 49 | 403 | 0 | None | First-party | Standard |
| `retro` | 96 | 131 | 1,315 | 0 | None | First-party | Standard |
| `review-standards` | 95 | 50 | 1,269 | 14,200 | None | First-party | Standard |
| `skill-repair` | 78 | 40 | 444 | 0 | None | Vendored: Google Apache-2.0 | Standard |
| `terraform` | 105 | 384 | 3,629 | 26,900 | None | Vendored: Anton Babenko Apache-2.0 | Standard |
| `typescript` | 82 | 35 | 584 | 0 | None | First-party | Standard |
| **Total** | **2,773** | **2,805** | **30,477** | **175,350** | **1 script, 5 assets** | — | — |

---

## 5. Subagents Inventory (6 Subagents)

Source of truth: `claude/agents/*.md`. Generated target: `antigravity/plugins/harness/agents/*.md`.

| Agent Name | Model (Claude) | Model (Antigravity) | System Prompt Tokens | Tool Allowlist | Skills Mounted |
|---|---|---|---|---|---|
| `changelog` | `haiku` | `flash` | 1,368 tok | Read, Grep, Glob, Edit, Write | None |
| `documentation` | `sonnet` | `inherit` (unmapped) | 1,021 tok | Read, Grep, Glob, Edit, Write | None |
| `explorer` | `sonnet` | `inherit` (unmapped) | 873 tok | Read, Grep, Glob, Bash, Write | None |
| `reviewer` | `opus` | `pro` | 1,278 tok | Read, Grep, Glob, Bash, Write | `review-standards` |
| `security` | `opus` | `pro` | 1,042 tok | Read, Grep, Glob, Bash, Write | None |
| `verifier` | `sonnet` | `inherit` (unmapped) | 1,391 tok | Read, Grep, Glob, Bash, Write | None |

---

## 6. Hooks & Lifecycle Middleware

Implemented by the unified Go binary `~/.agents/bin/harness hook <subcommand>`.

| Target | Event | Matcher | Command | Timeout | Can Block? | Can Inject Context? |
|---|---|---|---|---|---|---|
| **Claude Code** | `PreToolUse` | `Bash` | `harness hook pretooluse` | 10s | Yes (credential gate) | No |
| **Claude Code** | `SessionStart` | — | `harness hook sessionstart` | 10s | No | Yes (drift alert) |
| **Claude Code** | `UserPromptSubmit` | — | `harness hook userpromptsubmit` | 10s | No | Yes (trail notice) |
| **Claude Code** | `PostToolUse` | `Edit\|Write` | `harness hook posttooluse` | 30s | No (formats file) | No |
| **Claude Code** | `SubagentStart` | 5 agents | `harness hook subagentstart` | 5s | No (stamps marker) | No |
| **Claude Code** | `SubagentStop` | 5 agents | `harness hook subagentstop` | 10s | Yes (holds once if artifact missing) | No |
| **Claude Code** | `Stop` | — | `harness hook stop` | 15s | Yes (blocks uncommitted changelog) | No |
| **Antigravity** | `PreToolUse` | `run_command` | `harness hook pretooluse` | 10s | Yes (credential gate) | No |
| **Antigravity** | `PreInvocation` | — | `harness hook agy-preinvocation` | 10s | No | Yes (trail notice) |
| **Antigravity** | `PostToolUse` | file tools | `harness hook posttooluse` | 30s | No (formats file) | No |
| **Antigravity** | `Stop` | — | `harness hook agy-stop` | 15s | Yes (blocks uncommitted changelog) | No |

---

## 7. Tools & MCP Inventory

- **Claude Code Built-ins:** `Read`, `Edit`, `Write`, `Bash`, `Glob`, `Grep`, `Skill`, `Agent`, `WebSearch`, `FetchURL`.
- **Antigravity Built-ins:** `view_file`, `replace_file_content`, `write_to_file`, `run_command`, `find_by_name`, `grep_search`, `list_dir`, `ask_question`, `invoke_subagent`, `define_subagent`, `manage_subagents`, `send_message`, `schedule`, `generate_image`, `read_url_content`, `search_web`.
- **MCP Servers:** 0 configured by default. No MCP overhead at turn 0.
- **Tool Mapping in `internal/install/derive.go`:**
  - `Read` ↔ `view_file`
  - `Grep` ↔ `grep_search`
  - `Glob` ↔ `find_by_name`
  - `Bash` ↔ `run_command`
  - `Write` ↔ `write_to_file`
  - `Edit` ↔ `replace_file_content`

---

## 8. Memory, State & Trail

All persistent state is file-based and append-only under `.agents/docs/`:
- `changes/`: Changelog entries with YAML frontmatter (`owner: changelog`, `date:`, `scope:`).
- `architecture/DECISIONS.md`: ADRs (ADR-001 through ADR-009).
- `handoffs/`: Mid-session snapshots.
- `plans/`: Recon and architectural plans.
- `reviews/`: Code reviews and security reviews.
- `testing/`: Verification logs and `dispatch-ledger.jsonl`.
- `~/.agents/feedback/log.jsonl`: Passive correction log.

---

## 9. Runtime Settings Per Target

### Claude Code (`claude/settings.json`)
- `effortLevel`: `"high"`
- `promptCacheTtl`: `"1h"`
- `subagentPromptCacheTtl`: `"1h"`
- `fileCheckpointingEnabled`: `true`
- `cleanupPeriodDays`: `90`
- `permissions.defaultMode`: `"auto"` (with high-risk bash commands and dotfile edits set to `"ask"`)

### Antigravity (`antigravity/plugins/harness/plugin.json` + environment)
- Models: `pro` (heavy analysis/review), `flash` (logging/changelog), `inherit` (general)
- Execution mode: Sandboxed first (Standard Sandbox Mode), bypass requires user approval.

---

## 10. Observability & Telemetry

- `harness check`: Evaluates trail freshness against `p.Scope()`.
- `harness doctor`: Checks drift between repository sources and deployed files across 10 verification points.
- `.agents/docs/testing/dispatch-ledger.jsonl`: Records subagent dispatch runtime, exit codes, and token costs.

---

## 11. Reference Freshness Report (TODAY = 2026-09-03)

| Reference Item | Source Date | Age (Days) | Grade at Snapshot | Grade Today | Status / Action |
|---|---|---|---|---|---|
| **AGENTS.md spec** | 2025-08-01 | 398 | [C3] | **D (>12m)** | **Crossed into Grade D**. Core principle valid (plain markdown, nested nearest wins); specific early extensions superseded by AAIF living spec. |
| **AAIF founding** | 2025-12-01 | 276 | [C3] | **C (6-12m)** | Valid standard foundation. |
| **Agent Skills living spec** | 2026-08-15 | 19 | [A3] | **A (<=3m)** | Fresh; follow strictly (`name`, `description`, <500 lines, `references/`). |
| **MCP 2026-07-28 spec** | 2026-07-28 | 37 | [A3] | **A (<=3m)** | Fresh. |
| **OpenTelemetry GenAI** | 2026-08-01 | 33 | [A5] | **A (<=3m)** | Fresh. |
| **OWASP LLM Top 10** | 2026-08-04 | 30 | [A3] | **A (<=3m)** | Fresh. |
| **OWASP Agentic Top 10** | 2026-05-01 | 125 | [B3] | **B (3-6m)** | Fresh. |
| **OpenAI Harness engineering** | 2026-02-11 | 204 | [C2] | **C (6-12m)** | Established reference for minimal context files & bash primitives. |
| **Anthropic cost & intelligence** | 2026-08-01 | 33 | [A2] | **A (<=3m)** | Fresh; justification for selective reasoning effort & caching. |
| **Anthropic context engineering** | 2026-07-24 | 41 | [A2] | **A (<=3m)** | Fresh; progressive disclosure over upfront instruction dump. |
| **Warp self-improving agents** | 2026-08-26 | 8 | [A4] | **A (<=3m)** | Fresh. |
| **LangChain agent environments** | 2026-08-25 | 9 | [A2] | **A (<=3m)** | Fresh. |
| **Dave Rensin Elephant-Goldfish** | 2026-04-27 | 129 | [B4] | **B (3-6m)** | Established design pattern for cold-start verification. |
| **Manus context engineering** | 2025-07-01 | 429 | [D4] | **D (>12m)** | Historical baseline for KV-cache stability; principle only. |
| **Don't Break the Cache** | 2026-01-15 | 231 | [C1] | **C (6-12m)** | Valid cache-stability guidelines. |

---

## 12. Goldfish Verification of Inventory

A cold reader holding only `harness-inventory.md` can determine:
1. Every token entering context at turn 0 and on-demand.
2. The exact trigger, lifecycle, and surviving semantics of every file.
3. The origin and deployment mechanism for all 28 skills and 6 subagents.
4. The exact command invocation and blocking/informing behavior of each lifecycle hook.
