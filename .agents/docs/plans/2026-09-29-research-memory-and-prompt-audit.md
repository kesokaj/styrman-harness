---
owner: main
date: 2026-09-29
---
# Research — memory for Antigravity, and `prompt-audit` for both CLIs   (2026-09-29)

## Questions

The operator asked, 2026-09-29: (1) "Do some broader research on internet and github if you find anything for antigravity" — memory, knowledge and "dreaming" features or plugins for the Antigravity editor and the `agy` CLI; (2) "do a check if there are a ./claude-api prompt-audit that we can implement for agy and claude in this harness." Light mode: research, compare, recommend. Follow-on to `plans/2026-09-29-dream-skill.md`.

## Answer in brief

1. **Antigravity has no memory feature the harness can use, and nothing to adopt.** The editor generates Knowledge Items (auto-extracted at conversation end, viewable in the Agent Manager, stored under `~/.gemini/antigravity/knowledge/`), but the feature has no docs page (the one Google staff linked now 404s), two forum threads report it never firing, and the CLI does not expose it: no docs, no help entry, no settings key, and on this box `~/.gemini/antigravity-cli/knowledge/` holds only a lock file. Every community memory project for agy runs on hooks, a daemon, or a cloud service, which the harness rules out; the two hook-free ones re-do what the trail already does. Recommendation: nothing; the dream plan stays Claude Code only, and the trail is agy's memory.
2. **`prompt-audit` exists and can run on both CLIs.** It is a subcommand of Anthropic's `claude-api` skill: one 47 KB reference file, Apache-2.0, that audits prompts and agent configuration files (`AGENTS.md`, rules, skills, commands, subagents) for dated patterns and stale or contradicting instructions, and produces a report plus a proposed diff, never applying. On Claude Code it is bundled and runs today as `/claude-api prompt-audit`; agy has no copy. Recommendation: vendor the reference into a harness skill `prompt-audit` (option 2 below) so both CLIs run the same audit over `context/`, `agents/`, `skills/`, `commands/` and `rules/`, with the bug index as provenance.

## Part 1 — Antigravity memory

### What exists (verified 2026-09-29)

| Surface | Finding | Evidence |
|---|---|---|
| Editor (Antigravity 2.0) Knowledge Items | Auto-extracted from conversations; no setup; viewable in the Agent Manager; on disk as `~/.gemini/antigravity/knowledge/<item>/` with `metadata.json` and `artifacts/` | launch blog 2025-11-18; Google staff on the forum 2026-03-03; two export tools that read the directory |
| Editor KIs, reliability | Two forum threads (Dec 2025 to May 2026) report KIs never generated; escalated as a possible bug, no resolution posted | discuss.ai.google.dev threads 111234, 126774 |
| Editor KIs, docs | No docs page; `antigravity.google/docs/knowledge` is a 404; the docs nav has no Knowledge or Memory entry | fetched 2026-09-29 |
| CLI (`agy` 1.2.13) | No memory feature: docs overview and changelog say nothing; `agy --help` has no memory verb; `settings.json` has no memory key. The binary links the editor's Knowledge Base code (`KnowledgeBaseItem`, `GetUserMemories`, `GetCascadeMemories`, `knowledge_base_enabled`); `~/.gemini/antigravity-cli/knowledge/` exists with only `knowledge.lock` (created 2026-09-17, never written); `cli.log` never mentions knowledge | this box, `strings $(which agy)`, `agy changelog` |
| CLI, its own summaries | `~/.gemini/antigravity-cli/conversation_summaries.db` exists; one community project (`agy-dreamer`) reads it as its consolidation input | this box; github.com/j-steve/agy-dreamer |
| CLI plugin system | `agy plugin install <path>` (local path; the help also says `plugin@marketplace`), `import` from gemini or claude, `link <mp> <target>`; a plugin is `plugin.json` plus `skills/`, `agents/`, `rules/`, `mcp_config.json`, `hooks.json`. The marketplace is in-app (editor v2.18.1, 2026-09-28); no CLI marketplace verb documented; no memory plugin named in it | antigravity.google/docs/plugins; `agy plugin help` |
| Google announcements 2026 | I/O 2026 introduced the CLI; changelog through editor v2.18.1 and agy 1.2.13 has hooks (v2.11.0, agent-declared hooks v2.17.0) and the marketplace, and zero entries on knowledge, memory, learning or consolidation. No "dreaming" exists or is announced | changelog, I/O post |

### Community projects (GitHub, web-only survey; stars as shown on 2026-09-29)

| Project | Mechanism | Needs | Note |
|---|---|---|---|
| thedotmack/claude-mem (94.9k★) | `~/.gemini/config/hooks.json` (5 events) + MCP + `GEMINI.md` injection | hooks + local Node worker daemon | agy hooks target 1.2.1; several open schema-breakage issues; editor support closed "not planned" |
| Mem0 plugin v0.3.1 | PreInvocation / PostToolUse / Stop hooks + MCP `search_memories` + 6 skills | hooks, cloud API key, Python 3.10+ | no consolidation |
| mahirocoko/agy-memory-layer (2★, v1.23.0 2026-09-27) | plugin with PreInvocation / PreToolUse hooks; skills `/dream`, `/remember`, `/palace`; git-backed `~/.gemini/memory/` | hooks | `/dream` is regex by default, LLM reflection only by hand; human-gated |
| atlaso-labs/antigravity (0★) | SessionStart / PreInvocation / Stop hooks + MCP + skill + rules | hooks, Atlaso cloud | — |
| Digital-Process-Tools/claude-remember (186★) | hooks; Haiku-compressed session → daily → weekly summaries | hooks | no SessionEnd analogue on agy |
| HydStAn/agy-memory-engine (34★) | MCP + hooks → SQLite FTS5 | hooks + worker every 5 min | — |
| j-steve/agy-dreamer (0★) | cron 03:00 reads `conversation_summaries.db`, writes `~/.gemini/config/dreaming/` | cron daemon | closest to "dreaming"; three commits |
| nowledge-co, blackmoon87/antigravity-mem | hooks + MCP + daemon | server or daemon | — |
| thevisionhub/agent-sleep (1★) | MCP tool `agent_sleep_consolidate` | MCP server only | hook-free |
| feodus/antigravity-memory-bank (9★, 2025-12) | rules + markdown files autoloaded in planning mode | nothing | hook-free; the trail (`handoffs/`, `bugs/`, `retros/`) already is this |

### Comparison

| Approach | Pros | Cons | Complexity | Cost | Maturity |
|---|---|---|---|---|---|
| A. Nothing; the trail is agy's memory; re-check the CLI store after updates | no new store, no hooks, no daemon; matches project rule 1 | agy carries nothing across sessions beyond what a session reads from the trail | none | none | the trail: in use since 2026-09-03 |
| B. A hook-based plugin (claude-mem, Mem0, agy-memory-layer, claude-remember) | real capture across sessions | hooks and a daemon or cloud key; a second store beside the trail and Claude's auto memory; agy hook schema still breaking | medium | daemon or API key per machine | claude-mem mature on Claude Code, fragile on agy |
| C. A hook-free consolidation (agent-sleep MCP, or a harness skill reading `conversation_summaries.db`) | no hooks | reads an undocumented sqlite the CLI owns; an MCP server to run; nothing writes back into a store agy loads | medium | low | one to three commits each |
| D. Wait for the CLI's own Knowledge Base | first-party, code already linked in | undocumented, no switch, editor version unreliable | none now | none | not shipped |

### Recommendation: A

Nothing to adopt. Every real option needs hooks (project rule 1 and the global no-hooks rule) or a daemon; the hook-free ones duplicate the trail. The dream plan's Out section stands: the skill reads Claude Code's store only. Two lines to keep: re-check `ls ~/.gemini/antigravity-cli/knowledge/` after each agy update (already in the dream plan's Risks), and note that `conversation_summaries.db` is the input a future CLI-side consolidation would read, if the operator ever wants one.

## Part 2 — `prompt-audit`

### What it is

- A subcommand of Anthropic's `claude-api` skill (github.com/anthropics/skills, `skills/claude-api/`, Apache-2.0 per its `LICENSE.txt`). The bare invocation `/claude-api prompt-audit` reads `shared/prompt-audit.md` and executes it. On Claude Code the skill is bundled with the CLI (session-scoped extraction; this session's copy is at Claude Code 2.1.283). agy has no copy: it is not a plugin and not in any agy marketplace.
- The reference is one file: 47,233 bytes bundled, 48,983 bytes upstream at commit `8a1541c` (2026-09-29; the upstream adds Claude Sonnet 5.5 rows). It cites `shared/model-migration.md` twelve times and `shared/prompt-caching.md` once, all for API request code (which model 400s on which parameter). For a configuration repository those rows are `not applicable`; the file says so itself ("a configuration repository mostly Groups 1 and 2").
- What it does, in order: establish scope and target model without asking; inventory every file that reaches the model (`AGENTS.md`, rules, skills and their references, commands, subagent definitions); `git blame` for provenance ("which failure, on which model, did this prevent"); classify every line ("could the model already know this?"); scan four groups: dated prompt text (pressure language, scaffolds, over-specification, fossils, prohibition clusters), brittle configuration files (volatile paths and flags, contradicting files, history narratives, recency-trap rules), tool descriptions, request config; then a report (`file:line`, evidence, pattern, why obsolete, confidence, action) and a proposed diff, one finding per hunk, never applied unless the request said so. A keep list protects context, exact scripts for fragile operations, prohibitions against demonstrated failures, and working redundancy.

### Fit with this harness

- **The shape matches.** The harness is a configuration repository: `context/AGENTS.md` (21,539 of 22,000 bytes, ADR-035), seven agents, 27 skills, commands and rules. Group 2 ("brittle skill and configuration files": volatile paths, contradicting files, one session's stumble as a permanent rule) is the harness's known failure mode: BUG-095 and the 2026-09-28 retro were four review rounds over rules that contradicted their neighbours.
- **Provenance is on disk.** The audit's central question per prohibitive line is which failure it prevents. Here every such line is meant to trace to a `bugs/INDEX.md` row or an ADR; keep-list item 5 then protects it, and a line with no row is the finding. A quick signal count: `never` appears 29 times in `context/AGENTS.md` and 23 in `agents/verifier.md`; caps-lock boosters (MUST, NEVER, ALWAYS, CRITICAL, IMPORTANT) zero times in any of them.
- **Two target models, not one.** The reference reasons about Claude's documented behaviour. The files it would audit here are read by Claude Code and by Gemini on agy. A finding grounded only in Claude's behaviour ("current models are proactive by default") is `flag`, not `remove`, until it holds for both; a stale-fact or conflict finding (Group 2) is model-independent and stays actionable. This is the one adaptation the harness wording must add.
- **What it cannot do here:** Group 4 (request config, API fossils) has no surface in this repository; the twelve `model-migration.md` cross-references are unreachable in a vendored copy and matter only for that group.

### Comparison

| Approach | Pros | Cons | Complexity | Cost | Maturity |
|---|---|---|---|---|---|
| 1. Claude Code only: run `/claude-api prompt-audit` as bundled, record the report in `reviews/` | zero harness change; always the current Anthropic text | not on agy; the bundled copy moves with the Claude Code version, so two machines audit against two texts; the two-model caveat lives nowhere | none | none | shipped |
| 2. Vendor the reference as a harness skill `prompt-audit` (thin `SKILL.md`, `references/prompt-audit.md` verbatim at a pin, `LICENSE.txt`, row in `VENDORED.md`) | same audit on both CLIs; pinned text, updated on purpose; the harness `SKILL.md` adds scope, provenance and the two-model rule; report goes to `<docs-root>/reviews/<date>-prompt-audit.md` | a 47 KB reference to re-vendor at each pin; a fork by the vendoring rule (the harness `SKILL.md` is the fork, the reference stays verbatim) | one sitting: file drop, `VENDORED.md` row, skill-matrix row, one behavioural run per CLI | none | upstream: active, Apache-2.0 |
| 3. A `/review --prompts` lens dispatching `reviewer` with the reference | reuses the review record and the report-only contract | `reviewer` runs on sonnet with a fixed report template the audit does not fit; the audit is judgment-heavy over the harness's own rules and wants the main model; adds a lens to a 66-line skill | medium | none | — |
| 4. Install the `claude-api` plugin on agy (`agy plugin import claude`, or a marketplace link) | whole skill, not one file | 1.8 MB, 79 files, most of it SDK code guides the harness never needs; `agy plugin install` takes a local path only per the docs; plugin import is untested here and the harness rules say a skill is deployed by `sh-agent install`, not by a vendor's plugin system | medium | none | untested |

### Recommendation: 2 — decided: option 3 (operator, 2026-09-29)

The operator chose option 3, a `/review --prompts` lens; built the same day (`changes/2026-09-29-1251-review-prompts-lens.md`). The reference is vendored under `skills/review-standards/references/prompt-audit/`, not a new skill.

Vendor `prompt-audit.md` verbatim (pin `8a1541c`, Apache-2.0, `LICENSE.txt` beside it) under `skills/prompt-audit/references/`, with a harness `SKILL.md` that says: scope is `context/`, `agents/`, every skill under `skills/` not in `VENDORED.md`, `commands/` and `rules/`, or the paths the user names; provenance comes from `bugs/INDEX.md`, `DECISIONS.md` and `git blame`, and a prohibition with a row is keep-list item 5; the target is both vendors, so a Claude-only rationale is `flag`; Group 4 is `not applicable` here; the report is written to `<docs-root>/reviews/<date>-prompt-audit.md` and nothing is applied. Run it by hand, and from `/retro` when a model release lands (the reference's own re-audit trigger). It is the only option that gives agy the audit at all, and the harness already has the vendoring path for it (`VENDORED.md`, the `anti-ui-slop` precedent for a bundled licence).

What approving it would do: a Task under the Definition of Done (one sitting, no design question left): the skill directory, the `VENDORED.md` row, a `skill-matrix.md` row, `sh-agent install` into a scratch home, one behavioural run per CLI on a fixture with a seeded contradiction and a seeded dead path (expect both found, nothing applied), `/review --docs` over the range, then a first real run over `context/` with the report in `reviews/`, put to the operator hunk by hunk.

## Not proposed

- Running `/claude-api prompt-audit` over this repository now: the operator asked whether it can be implemented, not for the audit; its findings over `context/AGENTS.md` would be a review round in their own right and belong to the approved task.
- Any change to the dream plan's design: Part 1 confirms its Out section.
