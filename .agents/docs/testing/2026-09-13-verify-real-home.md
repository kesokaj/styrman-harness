---
owner: main
date: 2026-09-13
---
# Verification — real home on the Linux box updated to a35b567 and smoke-checked   (2026-09-13)

Range reviewed: a35b567 installed (`sh-agent status`: 260 files, no drift, `binary a35b567582fc`)

## Status: PASS for both CLIs

Scope: `sh-agent update` run by the agent on the operator's instruction ("run sh-agent update and the smoke checks here"), then the three smoke checks of `2026-09-03-verify-real-home.md` from the real `$HOME` on each CLI, in fresh clones of the demo service under `/home/kesokaj/sandbox/runs/realhome-{claude,agy}`. Not run: the Mac; the full skill matrix; interactive starts.

## Update

- Before: `~/.local/bin/sh-agent` dated 2026-09-03 21:05, `sh-agent version` → `(devel)`.
- `sh-agent update` (the old binary) → `up to date with upstream (a35b567)`, `building sh-agent from /home/kesokaj/workspace/styrman-harness`, `73 written, 187 unchanged, 2 removed, 0 kept — 260 files managed`, exit 0. The old binary's `update` ran the new source's install, which rebuilt and installed itself (BUG-006 closed on this box).
- After: binary dated 2026-09-13 22:33, `sh-agent version` → `a35b567582fc`; `sh-agent status` → `installed … from commit a35b567 (vendor: all, 260 files …)`, `binary a35b567582fc` with no warning, `repo has 27 skills, 7 agents`, `drift none — 259 files match the repo`, exit 0.

## Runs (real `$HOME`, no settings edited — the real Antigravity settings already carry `toolPermission: always-proceed` and trust `/home/kesokaj`)

| Check | Claude Code 2.1.270 | Antigravity 1.2.2 |
|---|---|---|
| Context questions, no tools: first heading; skill and agent counts; the merged skill absent; Definition of Done step 4 on TODO/ROADMAP at Done; the grilling command | PASS — `# AGENTS.md — global agent context`, 27 skills, 7 subagents, `gcloud-auth-verification` absent (`gcloud-setup` named), step 4 quoted ("at Done both are EMPTY (headers only) — finished lines are deleted, the `changes/` entry is their record"), `/grilling` | PASS — same heading, 27 / 7, absent, step 4 quoted, `/grilling` (named `/grill-me` as an alias from habit) |
| Knowledge skill auto-load (`golang` on `greet.go`, no edits) | PASS — rules listed under the `golang` skill, no edits; its one low finding was the smoke-check comment itself | PASS — `golang` and `review-standards` named, no edits |
| Agent dispatch (`verifier` on an uncommitted comment edit in `internal/greet/greet.go`) | PASS — ran in the clone, `Status: PASS`, `Scope:` line present, log `.agents/docs/testing/2026-09-13-verify-greet-comment-edit.md` written in the clone | PASS — ran in the clone, `Status: PASS`, `Scope: complete`, log `.agents/docs/testing/2026-09-13-verify-greet.md` written in the clone |

Antigravity real-home log `cli-20260913_223419.log`: 0 `Failed to parse`; no denied tool in any run. Evidence files: `.agents/scratch/e2e-realhome/` (gitignored).

## Not verified here

The Mac (still on the 2026-09-03 binary until the operator runs `sh-agent update` there); Claude Code on the Mac has never been checked; the full matrix.

## Environment

Linux (Debian 13, amd64), Go 1.24.4, Claude Code 2.1.270, Antigravity CLI 1.2.2 (auto-updated during the day), harness a35b567.
