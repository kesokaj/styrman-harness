---
owner: main
date: 2026-09-03
---
# Verification — first install into the operator's real home (Linux box)   (2026-09-03)

Range reviewed: cebd2c9 installed (`sh-agent status`: 256 files, no drift)

## Status: PASS for Claude Code; Antigravity BLOCKED on login

Scope: Claude Code exercised from the real `$HOME` in a throwaway clone of the demo service; Antigravity not exercised — the real home has no OAuth token, and the login is an interactive browser flow only the operator can complete. This is a smoke set, not the full matrix (`skill-matrix.md`), which ran against the sandbox home earlier the same day on the same files.

## Runs (real `$HOME`, cwd a fresh clone of the demo service)

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools: first heading, skill count and pruned names, subagents, DoD step 3, the grilling command | PASS — `# AGENTS.md — global agent context`, 28 skills, `grill-me`/`skill-repair` absent, seven agents, step 3 quoted with the deploy guides, `/grilling` | BLOCKED — `agy -p` stops at "paste the authorization code"; its log shows no skill loading before auth |
| Knowledge skill auto-load (`golang` on `greet.go`, no edits) | PASS — named the skill, no findings, no edits | — |
| Agent dispatch (`verifier` on an uncommitted comment edit) | PASS — ran in the clone itself, `Status: PASS`, `Scope:` line present, log written under `testing/` | — |

## Commands

`go run ./cmd/sh-agent install` → `256 written … 256 files managed`; `sh-agent status` → `drift none — 255 files match the repo`; `~/.claude/CLAUDE.md` is the generated two-line import.

## Not verified here

Antigravity from the real home (needs `agy` run once interactively to log in, then the same three checks); the Mac; the full matrix from the real home.

## Environment

Linux (Debian 13, amd64), Claude Code 2.1.259 from `~/.local/bin`, harness cebd2c9.
