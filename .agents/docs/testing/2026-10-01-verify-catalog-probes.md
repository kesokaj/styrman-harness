---
owner: main
date: 2026-10-01
---
# Verification — Phase 0 probes for the Google skills catalog, before anything is vendored (2026-10-01)

Range reviewed: 425a444 plus a scratch tree (`.agents/scratch/import/tree2`): the checkout's five source directories with one added skill, `skills/google-catalog/` — a draft router `SKILL.md`, the 141 candidate entries under `catalog/<name>/` at pins 2c39459 (google/agents-cli) and d5d9052 (google/skills), exec bits cleared, and a generated `catalog/INDEX.md` (36,435 bytes)

## Status: PASS — no leak on both CLIs; 2 of 3 routed tasks read a catalog entry on each, the third correctly took the harness's own skill

Setup: `go run ./cmd/sh-agent install --repo <scratch tree> --home ~/sandbox/audit-fix/sbImp2` (35 skills listed by `sh-agent list --repo`, render 15,155 bytes, no install refusal); a fresh clone of the demo service outside every checkout as the working directory; credentials copied into the sandbox home for the batch and deleted by name after it (`ls`: no such file, both).

| probe | Claude Code | Antigravity 1.2.14 |
|---|---|---|
| Skill list, no tools ("list every skill; is `gke-storage`, `bigquery-basics`, `google-catalog` in it; quote `typescript`'s description") | 48 names = 35 harness + 13 built-in; every harness skill present; `google-catalog` yes; `gke-storage` no, `bigquery-basics` no; `typescript`'s description quoted in full | 37 names = 35 + 2 built-in; every harness skill present; `google-catalog` yes; the two entry names no; description quoted in full |
| Task 1: "BigQuery reservation slots cost too much" (`--allowedTools Read Grep Glob Skill Bash(grep:*) Bash(ls:*)`; agy grants `read_file(*)`, `command(*)`) | `Skill` tool: `google-catalog`; transcript `Read` of `catalog/bigquery-optimization/…` and `catalog/bigquery-slot-cost-optimizer/…` | conversation database: `catalog/bigquery-optimization/SKILL.md` and a reference, `catalog/bigquery-slot-cost-optimizer/SKILL.md` and its script |
| Task 2: "GKE pod access to a bucket without a key file" | `Skill` tool: `gke-basics` (the harness's fork), its own `references/gke-security.md`; no catalog read — the harness skill that covers the task won, as the router says | the same: `gke-basics`, no catalog read |
| Task 3: "scaffold a new ADK agent project with agents-cli" | `google-catalog`; `Read` of `catalog/google-agents-cli-scaffold/…` | `catalog/google-agents-cli-scaffold/SKILL.md` |
| Denials | none | none (`auto-denied` 0 with those grants); a first probe with per-name command grants was denied a `find` (BUG-141's shape), which is a grant matter, not a read refusal |

Observed in the runs' own records: the Claude Code transcript (`sbImp2/.claude/projects/<clone>/*.jsonl`, tool_use entries) and Antigravity's conversation database (`strings <db> | grep google-catalog/catalog/`), not in the models' replies.

Also seen: Claude Code used `Grep` on `INDEX.md` three times and `Read` on it twice — the draft router says never to read the index whole; the final wording gives the grep command first and says what reading it costs.

Earlier the same day, as global skills (the design this replaces): Antigravity listed 109 of 173 and Claude Code showed descriptions for 41 of 173 (`plans/2026-10-01-research-google-skills-catalog.md`).

## What was not exercised

Interactive sessions (whether a default interactive session prompts before a read under `~/.claude/skills/…` or `~/.gemini/config/skills/…`: headless runs need the read grants above); the Mac; the four other families Phase 4 samples; the `allowed-tools` entry (`catalog-no-grant`, Phase 4).
