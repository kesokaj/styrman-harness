---
owner: main
date: 2026-09-29
---
# Research — memory and dreaming as subagents, the same on Claude Code and Antigravity   (2026-09-29)

## Question

The operator, 2026-09-29, while the dream skill was in its cold-reader check: "Cant we build dreaming and memory as a subagents that triggers when something in the project/conversation highlights to generate a skill/rule or procedure? Think about this later when you are done. I want them to work the same for claude and agy." Light mode: what the harness already has, what a subagent design would add, three ways to build it, one recommendation.

## Answer in brief

1. **The loop exists; its capture side is in-thread and its store is Claude-only.** Capture: `feedback-log` (a correction → one line in `~/.config/sh-agent/feedback.jsonl`) and Claude Code's auto memory (Claude writes files under `~/.claude/projects/<slug>/memory/`, Claude-only). Consolidation: `/dream` (auto memory → verdicts → promotions), `/distill` (third repeat → a skill draft), `/retro` (corrections and trail → proposals). Everything but the auto memory store already runs the same on both CLIs, because skills deploy verbatim and agents are derived for agy.
2. **A trigger "when something highlights" is the main agent noticing.** The harness forbids hooks and timers, and neither CLI fires a subagent on its own. Every trigger the harness has today is a rule in the context that the main agent applies (`feedback-log` on a correction, `/distill` on the third repeat, `/retro` on two corrections). A subagent design keeps that: the main agent dispatches when a rule fires; the subagent brings independent judgement and its own artifact, never the conversation (context § Work).
3. **Recommendation: one new subagent, `dreamer`, plus a vendor-neutral store; keep capture in-thread.** The consolidation judgement (is this a rule, a skill, a procedure, or nothing?) is the part that benefits from an author-independent reader, and it is heavy enough to earn a dispatch. Capture is a one-line append at the moment of the correction; a dispatch per correction would cost more than it adds and the subagent could not see the exchange anyway. Vendor parity comes from where the store lives, not from who writes it: point Claude Code's auto memory at a harness-owned per-project directory (`autoMemoryDirectory` in `.claude/settings.local.json`, written by `/init-project`), so Claude writes natively and agy reads the same files under the Resume Protocol and the `dreamer` reads them on both.

## What each CLI offers (verified 2026-09-29)

| | Claude Code 2.1.284 | Antigravity CLI 1.2.13 |
|---|---|---|
| Native memory | auto memory: `~/.claude/projects/<slug>/memory/`, `MEMORY.md` index loaded at start (200 lines / 25 KB), topic files on demand; `autoMemoryDirectory` from any settings scope (`code.claude.com/docs/en/memory.md` § Storage location) | none (`plans/2026-09-29-research-memory-and-prompt-audit.md` Part 1) |
| Subagents | `agents/*.md`, tools mapped by `internal/harness/derive.go` | the same definitions, derived at install (`agent.md`, `subagent: true`) |
| Subagent memory | a subagent may keep its own auto memory (`memory` field, docs § How it works) — Claude-only | none |
| Event triggers | hooks (forbidden here) | hooks (forbidden here) |
| What a subagent sees | its dispatch text and the files it reads; never the conversation | the same |

## Three ways to build it

| Approach | Pros | Cons | Complexity | Cost | Maturity |
|---|---|---|---|---|---|
| A. Status quo plus `/dream` (skills in-thread; store stays Claude's) | nothing new; every skill already runs on both CLIs | agy has no store to consolidate; the consolidation judgement is the author's own | none | none | shipped today |
| B. `dreamer` subagent + harness-owned store via `autoMemoryDirectory` (recommended) | independent judgement on what becomes a rule, skill or procedure; one store both CLIs read, Claude writes natively; no hooks; a file drop (agent) plus one `/init-project` edit | the store's location is a per-project setting `/init-project` must write (`.claude/settings.local.json`, gitignored, honoured only in a trusted workspace); agy sessions write to the store only through the agent, since agy has no native writer | one sitting plus behavioural runs on both CLIs | none | agents on both vendors: proven by every harness agent |
| C. `memory` and `dreamer` subagents, store in the project trail (`.agents/docs/memory/`, committed) | cross-machine, reviewable in git | memory about the operator (`user` type) lands in every project's history; a capture dispatch per correction is slow and blind to the exchange; duplicates `feedback.jsonl` and the handoffs | medium | none | — |
| D. Subagents that write rules and skills directly | "generate a skill/rule" literally | no gate: a rule an agent changes on its own is the weaker home (global Hard rules on `REQUIREMENTS.md`); bypasses `/retro`'s promotion check and the cold read | low | none | rejected |

## Recommendation: B — decided: B (operator, "Option B", 2026-09-29); plan `plans/2026-09-29-dreamer-memory-loop.md`, spec `specs/sh-agent/memory-loop.md`; the store's location was then moved inside the project (`.agents/memory/`, gitignored) on the operator's "Option 2" the same day (plan D17), so the `~/.config/sh-agent/memory/<slug>/` path below is superseded

- **Store.** `~/.config/sh-agent/memory/<project slug>/`, the auto memory layout as is (`MEMORY.md` plus one file per fact, the frontmatter Claude writes). `/init-project` writes `{"autoMemoryDirectory": "~/.config/sh-agent/memory/<slug>"}` into `.claude/settings.local.json` (gitignored) so Claude Code reads and writes there natively; the context's Resume Protocol gains one clause so an agy session reads `MEMORY.md` from that directory at start; `sh-agent report` does not ship it: `internal/harness/report.go` sends an allowlist (`feedback.jsonl`, `retros.jsonl`, `proposals/`), and memory stays out by design — the `dreamer` artifact in `retros/` is what carries learning into the project (corrected 2026-09-29 while planning; the first draft said the report would ship it). This keeps `/dream` unchanged (it resolves `autoMemoryDirectory` first).
- **`dreamer` agent** (`agents/dreamer.md`, tools Read, Grep, Glob, Bash, Write; sonnet). Input: the store path and the trail. Output: one artifact under `<docs-root>/retros/<date>-dreamer.md` listing what it would make — a rule (target file and the exact text), a skill (a draft `SKILL.md` under `.agents/scratch/`), or a procedure (a `references/` page) — each with the memories and trail rows that justify it, and what it left alone. It applies nothing; `/retro` and `/distill` consume its artifact as evidence (retro Step 1 source 8; distill's candidate list). Dispatched by `/retro` and `/dream`, and by the main agent when the distill trigger fires (third repeat) — a rule in the context's Learning loop, which the size budget (21,636 of 22,000 bytes) can take only if a sentence pays for it.
- **Capture stays in-thread.** `feedback-log` and Claude's own auto-memory writes are the capture; on agy the main agent writes the same memory file shape into the store when a rule fires (one Write, the frontmatter Claude uses), which is the only vendor gap left and closes it without a dispatch.
- **Not built:** a `memory` subagent (blind to the conversation, one dispatch per correction), timers or counters, any hook.

What approving it would do: a Feature under the feature skill — plan, cold-reader check, then phases: (1) `agents/dreamer.md` with behavioural runs on both CLIs (the agy sandbox procedure or the real-home smoke, per `AGENTS.md` § Testing), (2) the store move (`/init-project`, the resume clause, `/dream` unchanged), (3) `/retro` and `/distill` reading the artifact, (4) docs and the map are untouched (no delivery path changes; rule 3 does not fire).

## Not proposed

- Replacing `/dream` with the agent: the skill's verdict logic is what the agent would run; keep the skill as the in-thread path and let the agent load it.
- Auto-applying anything: every output stays a proposal until `/retro` or the operator says go.
