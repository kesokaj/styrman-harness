---
owner: main
date: 2026-09-03
---
# Verification — first install into the operator's real home (Linux box)   (2026-09-03)

Range reviewed: cebd2c9 installed (`sh-agent status`: 256 files, no drift)

## Status: PASS for both CLIs

Scope: both CLIs exercised from the real `$HOME` in a throwaway clone of the demo service, three checks each (context, skill load, agent dispatch). Antigravity ran after the operator logged in; its headless runs need tool permissions that settings must grant, so a temporary `permissions.allow` set (`command(*)`, `read_file(*)`, `write_file(*)`, `edit_file(*)`, `list_dir(*)`) was written to `~/.gemini/antigravity-cli/settings.json` for the six minutes of the checks and the original file restored from a backup afterwards; interactive use prompts instead. This is a smoke set, not the full matrix (`skill-matrix.md`), which ran against the sandbox home earlier the same day on the same files.

## Runs (real `$HOME`, cwd a fresh clone of the demo service)

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools: first heading, skill count and pruned names, subagents, DoD step 3, the grilling command, the canary word (expect GANNET — the last line of the rendered context; on Antigravity also `strings` the newest `~/.gemini/antigravity-cli/conversations/*.db` and expect no `<truncated` line inside the `<RULE[...AGENTS.md]>` block, BUG-085) | PASS — `# AGENTS.md — global agent context`, 28 skills, `grill-me`/`skill-repair` absent, seven agents, step 3 quoted with the deploy guides, `/grilling` | PASS — same heading, 28 skills, both pruned names absent, seven agents, step 3 quoted; named `/grill-me` from habit while noting `/grilling` is what the harness has; zero `Failed to parse skill` lines in its log |
| Knowledge skill auto-load (`golang` on `greet.go`, no edits) | PASS — named the skill, no findings, no edits | PASS — `golang` via `review-standards`, no findings, no edits |
| Agent dispatch (`verifier` on an uncommitted comment edit; Antigravity also `explorer`) | PASS — ran in the clone itself, `Status: PASS`, `Scope:` line present, log written under `testing/` | PASS — `explorer` mapped the three packages and their imports; `verifier` ran in the clone, `Status: PASS`, `Scope: complete`, log written under `testing/` |

## Commands

`go run ./cmd/sh-agent install` → `256 written … 256 files managed`; `sh-agent status` → `drift none — 255 files match the repo`; `~/.claude/CLAUDE.md` is the generated two-line import.

## Not verified here

The Mac; the full matrix from the real home; Antigravity interactive prompts for `~/.config/sh-agent` (only the headless form was seen).

## Environment

Linux (Debian 13, amd64), Claude Code 2.1.259 from `~/.local/bin`, harness cebd2c9.
