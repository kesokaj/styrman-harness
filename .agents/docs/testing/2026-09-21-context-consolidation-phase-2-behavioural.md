---
owner: main
date: 2026-09-21
---
# Verification — behavioural runs of the changed `agents/bugs.md`, Phase 2   (2026-09-21)

Range reviewed: working tree over 1df57f3 — `agents/bugs.md` § Rules gains one bullet (a recommended permanent home in the context names the existing sentence the rule would replace or merge into, or the rule goes to the owning skill); `skills/retro/SKILL.md` gains a target-table row and a size sentence in Step 6. Recipe: this repo's `AGENTS.md` § Testing (BUG-083).

## Status: PASS — 3 of 3 (after a failed first arrangement, BUG-087)

## Setup

`HOME=.agents/scratch/sandbox-home` (logged in, harness reinstalled with the Phase 2 edits, drift none). The agent definition was built from the working-tree `agents/bugs.md` into `--agents` JSON (`description`, `prompt` = the body, `tools: Read, Grep, Glob, Edit, Write`, `model: sonnet`); the dispatch text did not repeat the changed rule. Prompt on stdin; `--permission-mode acceptEdits --max-turns 25 --allowedTools Agent Read Write Edit Glob Grep Bash(git:*) Bash(ls:*) Bash(cat:*) Bash(date:*) Bash(grep:*)`. Each run in its own subshell with its own clone (the Phase 1 sandbox log's caveat).

## First arrangement — demo-service clones, failed (BUG-087)

Three clones of the demo service with three fixture records sharing one prevention rule; dispatch: "for the fourth time, a commit message named a subagent internal id — the body of the commit at HEAD carries the line agentId: a1b2c3". The fixture's HEAD had no such line. Runs 1 and 3: the main agent ran `git log`, found no such line, and asked the user (no dispatch, no record). Run 2: dispatched; the record said "not reproduced at HEAD" and recommended a project `.agents/rules/` entry or a lint — never the global context, because in a project repo the bugs agent does not propose `context/AGENTS.md` as a home. A second fixture with the line committed was not run: it could still not exercise the changed rule. The arrangement the rule measures is this harness repo, where graduated rules live in the context.

## Second arrangement — three clones of this repo at 1df57f3, subject BUG-086 (a real failure of the day)

Dispatch: record the install piped into `head -1` that died on SIGPIPE; "the record must recommend where the prevention rule should live permanently".

| Run | Record | Recommended home | Names the sentence it widens / an owning skill | Subagent tool order (transcript) |
|---|---|---|---|---|
| 1 | `BUG-086-install-piped-into-head-sigpipe.md` | widen the existing Hard rules bullet "A verification command is never piped into a filter…"; rejects the `bash` skill, the project AGENTS.md and a `.agents/rules/` file with reasons; cites the 22,000-byte pin | yes | Read Glob Read Write Edit |
| 2 | same name | widen the same bullet "rather than add a new bullet"; "one clause, not a new sentence"; cites ADR-035, 20,012 / 22,000, ~2,000 bytes headroom; `bash` skill only as fallback | yes | Read Glob Read Grep Grep Grep Grep Read Glob Write Edit |
| 3 | same name | widen `context/AGENTS.md` line 61 into one bullet covering three shapes; cites 16,569 source bytes, 20,012 rendered, 1,988 headroom, "confirms against `TestRenderedContextFitsAntigravityBlock` before committing"; also found BUG-011's row never moved to Graduated | yes | Read Glob Read Grep Grep Grep Read Grep Write Edit Edit Grep … Read Edit Edit Edit Read |

Main-agent transcripts show `Agent(bugs)` once per run; no run committed. The subagent used no tool outside its definition's list. What only the new bullet asks for — naming the existing sentence and the byte cost — appears in 3 of 3 records.

## Not exercised

Antigravity (the recipe is Claude Code's `--agents`; Antigravity's derived agent was checked by install, `status` and a CLI start only). A control run with the OLD definition (the baseline is inferred from the first arrangement's run 2, which named no sentence, and from the definition text).
