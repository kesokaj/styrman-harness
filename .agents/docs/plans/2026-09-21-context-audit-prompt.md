---
owner: main
date: 2026-09-21
---
# Agent context audit

## Role

Audit the instruction surfaces of the agent harness you are running in (or the repo I point you at) against the standard below.

- READ-ONLY. Do not edit, create, move, or delete anything. Propose changes instead.
- Count with commands (`find`, `wc -l`, `wc -m`, `wc -c`, `grep`). Never estimate what you can measure.
- Label every statement MEASURED, DOCUMENTED (cite the source), or INFERRED.
- If a surface doesn't exist, write "n/a". If you cannot determine something, state your assumption and continue.
- Work through the steps in order.

## Target design

Sort content by when it loads, and write for the tightest limit across harnesses.

- **L0 Enforcement.** Anything that must hold every time lives in hooks, permission rules, sandboxing, linters, or CI, not in prose.
- **L1 Always-on.** One root file: build and test commands, conventions that differ from tool defaults, and pitfalls the agent cannot read from the code. With several harnesses, AGENTS.md is canonical, CLAUDE.md is `@AGENTS.md` plus Claude-only lines, and GEMINI.md and .cursor/rules hold only harness-specific deltas. With one harness, use its native root file.
- **L2 Scoped rules.** One topic per file, triggered by glob, path, or description.
- **L3 Skills.** Multi-step procedures and domain knowledge, loaded by description match or explicit invocation. One canonical skills directory; other harness directories link to it.
- **L4 References.** Design docs, READMEs, schemas, API docs, and scripts, loaded only when read and never imported into L1. A design doc holds the problem, the plan, rejected alternatives, and the file-by-file change list.

Two tests decide length. The right length passes both.

- **Pruning (necessity):** keep a line only if removing it would cause the agent to make mistakes.
- **Goldfish (sufficiency):** a fresh agent with zero history can act correctly from the file alone.

Practice: build context in one session, write it down, then execute and review in fresh sessions fed only that document.

## Step 0 - Map the loading model

Identify the harness(es) in use. For each instruction surface, record when it enters context:

- ALWAYS: every session.
- SCOPED: triggered by glob, path, or description.
- ON-DEMAND: read only when referenced or invoked.

Use the harness's docs, config, or source. For a custom harness, read the code that assembles the prompt. Also record: context window; truncation or byte caps; skill-listing budget, if any; what survives compaction; whether fresh-context subagents exist; deterministic controls (hooks, permission rules, sandbox, CI).

## Step 1 - Verify the reference limits

Collected September 2026. With web access, fetch each source and mark each limit CONFIRMED, CHANGED (give the new value), or UNVERIFIABLE; without it, mark all UNVERIFIABLE. Use CHANGED values in Step 3. Apply only the rows for the harnesses found in Step 0; for a custom harness, treat all rows as the portable standard.

- R1 Claude Code memory: CLAUDE.md target <200 lines per file; loads in full up to 4 MiB; @imports load at launch and save no context; AGENTS.md is read only if imported; MEMORY.md loads only its first 200 lines or 25 KB. https://code.claude.com/docs/en/memory
- R2 Claude Code skills: description + when_to_use cut at 1,536 chars; listing budget ~1% of the context window, least-used descriptions dropped first; after compaction only the first 5,000 tokens per skill return (25,000 combined). https://code.claude.com/docs/en/skills
- R3 Antigravity: rules and workflow files capped at 12,000 chars each; global rules in ~/.gemini/GEMINI.md; workspace rules in .agents/rules; activation: Manual, Always On, Model Decision, Glob; skills in .agents/skills and ~/.gemini/config/skills. Check whether a root AGENTS.md is read. https://antigravity.google/docs/rules-workflows/ and https://antigravity.google/docs/skills/
- R4 Cursor: rules <500 lines; .mdc required in .cursor/rules; types: Always Apply, Apply Intelligently, file globs, manual; AGENTS.md at root and nested. https://cursor.com/docs/rules
- R5 Codex: AGENTS.md chain (root to cwd) shares 32 KiB (project_doc_max_bytes), truncated silently; skills list capped at 2% of the context window or 8,000 chars, descriptions shortened first, skills possibly omitted. https://developers.openai.com/codex/guides/agents-md and https://developers.openai.com/codex/skills
- R6 Agent Skills spec: name 1-64 chars (lowercase, digits, single hyphens, not at either end) equal to the directory name; description 1-1,024 chars stating what and when; SKILL.md <500 lines and <5,000 tokens; references one level deep. https://agentskills.io/specification
- R7 Evidence: more simultaneous instructions, lower adherence (arXiv:2507.11538); context files cost 20%+ more without improving success (arXiv:2602.11988); adherence decays within a session; file length 25-500 lines had no effect (arXiv:2605.10039).
- R8 Elephant-Goldfish Model (Rensin, 2026), source of the Step 5 test. https://drensin.medium.com/elephants-goldfish-and-the-new-golden-age-of-software-engineering-c33641a48874

## Step 2 - Inventory

Find every surface:

- Root and home files: AGENTS.md, AGENTS.override.md, CLAUDE.md, CLAUDE.local.md, GEMINI.md in the repo, parent dirs, ~/.claude, ~/.codex, ~/.gemini, ~/.cursor
- Rules: .claude/rules, .cursor/rules, .cursorrules, .agents/rules, .agent/rules, workflow files
- Skills: {.agents,.agent,.claude,.cursor,.codex}/skills, ~/.agents/skills, ~/.gemini/config/skills, /etc/codex/skills
- Subagent definitions, commands, MEMORY.md, custom system-prompt files
- References: design docs, READMEs, and any doc an instruction file points to

Per file: path, lines, chars, bytes, tokens (bytes / 4, or a real tokenizer), layer, load mode, trigger if scoped. Resolve imports and @mentions and count them where they load.

Per skill, also: name, name equals directory?, description length and its first 80 chars, body lines and tokens, reference depth, frontmatter keys.

Totals:

- (a) ALWAYS-on lines, chars, and tokens
- (b) discrete imperative instructions in ALWAYS-on content (one instruction = one independently checkable directive)
- (c) sum of name + description chars over model-invocable skills
- (d) each total as % of the context window

If the harness has its own accounting (/context, /doctor, /skill-doctor, /status), ask me to run it and reconcile.

## Step 3 - Score against the standard

Score each check PASS, WARN (over target), or FAIL (over ceiling). Targets are judgment calls to tune with Step 7; ceilings are documented limits from Step 1.

List DEFECTS first: anything silently truncated, dropped, or never loaded (wrong extension, wrong directory, dead glob, missing description, over a byte cap).

- A1 Single always-on file: target <=150 lines; ceiling 200 lines or 12,000 chars.
- A2 All always-on content combined: target <=300 lines or ~4,000 tokens; ceiling 32 KiB.
- A3 Discrete always-on instructions: target <=100.
- A4 Scoped rule file: one topic; target <=100 lines; ceiling 500 lines or 12,000 chars; trigger works (glob matches a real file; description present if model-decided).
- K1 Skill name and description meet R6; description in third person, says what AND when, trigger words first; target <=250 chars.
- K2 SKILL.md body: ceiling 500 lines and 5,000 tokens; critical instructions first.
- K3 Skill listing, total (c): ceiling 8,000 chars; also report it against 1% and 2% of the context window. Flag skill pairs a human could not choose between.
- K4 References one level deep from SKILL.md; reference files >100 lines open with a table of contents.
- K5 Multi-harness skills use only spec frontmatter keys: name, description, license, compatibility, metadata, allowed-tools.
- E1 Every rule that must always hold (MUST, NEVER, security, data loss) has a deterministic backstop: hook, permission rule, linter, or CI. List those that don't.
- P1 With several harness surfaces: one source of truth, as in L1. Report duplicated or contradictory text, with both locations.
- P2 With several harnesses: one canonical skills directory, the others link to it. Report drifted copies.
- D1 READMEs and design docs are ON-DEMAND, never imported into ALWAYS-on files. Flag AI-generated READMEs or context files that were not human-verified.
- D2 Work in progress has a written intent and state doc a new session could resume from. Design docs follow L4.

## Step 4 - Pruning review (necessity)

Classify every instruction or paragraph in ALWAYS-on content as exactly one of:

- KEEP: specific, verifiable, not derivable from the repo, relevant to most sessions
- DERIVABLE: layouts, dependency lists, overviews
- VAGUE: cannot be checked
- COMMON-KNOWLEDGE: the model already does this
- MOVE-TO-SCOPED: applies to part of the codebase
- MOVE-TO-SKILL: multi-step procedure
- MOVE-TO-ENFORCEMENT: must hold every time
- DUPLICATE or CONFLICT: quote both locations
- STALE: path, command, or tool no longer exists (verify)

Report counts per class and the projected ALWAYS-on size if only KEEP remained.

## Step 5 - Goldfish test (sufficiency)

Every new session is a goldfish: it knows only what the files put in front of it. Test each ALWAYS-on file, each skill, and any design doc or README I name.

First state how many runs you plan and a rough token cost. Above 15 runs, test the ALWAYS-on files plus the 5 largest and 5 most recently changed skills, and list what you skipped.

For each file, start a fresh subagent with NO conversation history, given only that file plus the files it references. If the harness has none, give me the prompts to run in a new session. Three separate checks:

- G1 Comprehension: "Using only this file and what it references, explain what it is for, how the related system works, and what you must and must not do." Compare with the repo; list every wrong or missing point.
- G2 Critic: "Review this as an expert. List ambiguities, contradictions, faulty assumptions, and missing edge cases." Keep only findings that would change agent behavior.
- G3 Readiness: "Could you do <task> correctly on the first pass with only this? List every question you would need answered." Pick a representative task from recent git history and name it. Each question is a gap in the file.

PASS = no material G1 errors and no blocking G3 questions. Propose each fix in the file where a goldfish would look for it, then re-check it against Step 4. Per file, report tokens, G1 errors, G2 findings kept, G3 questions, and two verdicts: SUFFICIENT (goldfish) and MINIMAL (pruning).

## Step 6 - Target layout

Propose where each piece of content should live, one row per move: content (file:lines) | current layer | proposed layer | reason | check fixed | token change.

Then give:

- projected totals per layer against the Step 3 budgets
- wiring: imports or links to add, duplicate files to retire
- rules to move into L0, naming the hook event, permission rule, linter rule, or CI check
- design docs or READMEs worth adding, and where

## Step 7 - Propose, do not run, an adherence test

One trivial machine-detectable canary rule; three representative tasks, at least one modifying existing code; N fresh headless sessions with the current files versus the proposed layout. Measure canary compliance early versus late in the session, task success, tokens, and skill trigger rates on should-trigger and should-not-trigger prompts. Give exact commands and an estimated cost. Wait for my approval.

## Output

1. Defects
2. Loading-model map
3. Reference-limit verification
4. Inventory with totals
5. Scorecard A1-D2
6. Pruning counts and projected ALWAYS-on size
7. Goldfish results per file
8. Target layout
9. The 10 highest-value changes, each with file:line, check fixed, and expected effect
10. Adherence test plan

Make no edits until I approve specific items.
