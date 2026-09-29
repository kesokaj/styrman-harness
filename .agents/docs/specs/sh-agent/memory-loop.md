---
owner: main
date: 2026-09-29
---
# The memory loop — one store for both CLIs, and the `dreamer` agent

Status: specified 2026-09-29 (plan `plans/2026-09-29-dreamer-memory-loop.md`); nothing below is built until that plan's boxes say so.

## The store

- **Path:** `~/.config/sh-agent/memory/<slug>/`, one directory per project, where `<slug>` is the repository root's absolute path with every character other than a letter, a digit or a hyphen replaced by a hyphen — the name Claude Code gives the project under `~/.claude/projects/` (a `.` in the path becomes `-` too). Machine-local, never committed, never sent by `sh-agent report` (which ships only `feedback.jsonl`, `retros.jsonl` and `proposals/`, `internal/harness/report.go`).
- **Layout:** Claude Code's auto memory layout, unchanged: `MEMORY.md` (one line per memory, `- [Title](file.md) — hook`, under 200 lines and 25 KB) and one file per memory with the frontmatter Claude writes (`name`, `description`, `metadata.type` in `user` / `feedback` / `project` / `reference`, `metadata.originSessionId`, `metadata.modified`).
- **Who writes:** Claude Code natively, because the project's `.claude/settings.local.json` sets `autoMemoryDirectory` to that path (`~/` form; the docs accept an absolute path or `~/`). On Antigravity the main agent writes a memory file in the same shape when the `feedback-log` skill fires (its rule), since agy has no native writer.
- **Who reads:** Claude Code loads `MEMORY.md` at start (its own rule, 200 lines / 25 KB). An Antigravity session reads `MEMORY.md` from the store at the end of the Resume Protocol (context § Resume). `/dream` consolidates it on either CLI (its Step 1 resolves `autoMemoryDirectory` first). The `dreamer` agent reads it on either CLI.
- **Set up by** `/init-project`, new and update paths: write the key into `.claude/settings.local.json` (create the file when absent; add the key when the file exists without it; keep and report a different value), add `.claude/settings.local.json` to `.gitignore`, `mkdir -p` the store, and move an existing default store's files (`~/.claude/projects/<slug>/memory/*`) into it when the new store is empty — move, never copy, one store.
- **Trust:** Claude Code honours a project settings file only in a trusted workspace. Where it does not, Claude writes to the default store while the setting names the other; `/dream` detects a split (memory files in both) and stops with both paths for the operator to merge.

## The `dreamer` agent

- **What it is:** the consolidation judgement of the learning loop, run by an author-independent reader: given the store and the trail, what should become a rule, a skill or a procedure, and what should be left alone. It proposes; it never applies.
- **Definition:** `agents/dreamer.md` — tools `Read, Grep, Glob, Bash, Write`; `model: sonnet`; `effort: medium`; `maxTurns: 40`; derived for Antigravity by the installer like every agent.
- **Dispatch text (the contract):** `Store: <path>. Docs root: <docs-root>. Since: <date of the newest retros/ record, or "all">. Promoted: <the memory files /dream flagged promote, or none>.` plus the standard artifact and budget sentences. Never the conversation. The dreamer does not re-derive `/dream`'s verdicts: the promoted files are its rule-candidate inputs, one candidate row per memory file.
- **Reads:** the store (`MEMORY.md` and every memory file); the trail since the date — `handoffs/` (newest three), `bugs/INDEX.md` rows and their prevention rules, `changes/` entries; what already exists — the global context, the project `AGENTS.md` and `.agents/rules/`, the installed skills' descriptions (`sh-agent list` or the skill directories), so a candidate is never something a rule already says.
- **Candidate tests:** a **rule** for each file on the dispatch's `Promoted:` line (given, never re-judged: `/dream` owns the two-memory test), for a bug row whose prevention rule is recorded nowhere loaded and that recurs (two rows, or a row and a memory on the same fact), and for a memory of any type that a trail row restates (a handoff or `changes/` line on the same fact) when no loaded rule holds it; a **skill** when a memory or the trail shows a procedure of three or more steps done more than once (`/distill`'s bar); a **procedure** (a `references/` page for an existing skill) when the steps belong to a skill that exists; **none** otherwise, listed under "Left alone" with the reason.
- **Writes:** one artifact, `<docs-root>/retros/<date>-dreamer.md` (`owner: dreamer`, `date:`), written early with `Status: IN PROGRESS` and finished with `COMPLETE`: Scope, Summary, a Candidates table (`#`, kind, target file, the exact rule text or the draft's path, evidence as `file:line` and memory file names), "Left alone", "Not examined". Drafts for skills go under `.agents/scratch/dreamer/<date>/<name>/SKILL.md`. No other file, ever; no edit to source, config, rules, skills or memory.
- **Consumers:** `/retro` dispatches it (Step 1 source 8) and turns rule candidates into proposals through its own Step 2 (target, scope, promotion check) — the one place that decides where a rule lives; `/distill` dispatches it first and takes its skill and procedure candidates as candidates. No `retros.jsonl` line of its own; the retro or distill record cites the artifact.

## Error handling

| Situation | Behaviour |
|---|---|
| No store, or an empty one | `/dream` and the `dreamer` say so in one line; the resume clause reads nothing |
| Split store (files in both the settings store and the default) | `/dream` stops with both paths; the `dreamer` reads the settings store and names the other under "Not examined" |
| `~/.config/sh-agent/` denied to the session | the reader says so; the `dreamer` artifact still writes into the project trail |
| `.claude/settings.local.json` exists with another `autoMemoryDirectory` | `/init-project` keeps it and reports; the store is wherever it points |
| Agent cut at its turn limit | a partial artifact with `Status: IN PROGRESS` (its early-write rule); the caller says so |
