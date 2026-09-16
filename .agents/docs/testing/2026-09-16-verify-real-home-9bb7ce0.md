---
owner: main
date: 2026-09-16
---
# Verification — install into the operator's real home, commit 9bb7ce0   (2026-09-16)

Range reviewed: 8b95f22..9bb7ce0 (retro 2's six proposals, ADR-026's release procedure, BUG-045's two prevention halves)

## Status: PASS for Claude Code, PARTIAL for Antigravity

Scope: the install itself and the three smoke checks of `2026-09-03-verify-real-home.md` from the real `$HOME`, in a throwaway Go service at `.agents/scratch/realhome-smoke/demo` (its own git repository). Claude Code ran all three. Antigravity ran the first; its other two need tool permissions written into `~/.gemini/antigravity-cli/settings.json`, and this session's permission classifier refused that write twice ("Create Unsafe Agents", then "Self-Modification") — correctly, since it changes a CLI's permission file. Those two are the operator's to run; the commands are below. Not the full matrix (`skill-matrix.md`), which last ran against the sandbox home.

## Install

`go run ./cmd/sh-agent install` → `building sh-agent from …`, `7 written, 253 unchanged, 0 removed, 0 kept — 260 files managed`.
`go run ./cmd/sh-agent status` → `repo … (main @ 9bb7ce0, clean)`, `installed 2026-09-16T17:35:25+02:00 from commit 9bb7ce0`, `binary 9bb7ce059892`, `repo has 27 skills, 7 agents`, `drift none — 259 files match the repo`. No advisory note: `~/.claude/settings.json` already holds the attribution key.

## Runs (real `$HOME`, cwd the throwaway service)

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools: first heading, skill and agent counts, DoD step 3, and the Hard-rule sentence this commit changed | PASS — `# AGENTS.md — global agent context`, 27 skills and 7 subagents, step 3 quoted, and the new "After any dispatch that can write …" sentence quoted verbatim including `git log --oneline @{u}..HEAD`, `git tag` and the BUG-045 clause. Asked whether the global Git rule mentions signed release tags it correctly answered no — ADR-026 is this project's rule, not a harness-global one | PASS — same heading, same counts, step 3 quoted with the deploy guides, the same new sentence quoted verbatim; `grep -c 'Failed to parse' ~/.gemini/antigravity-cli/log/cli-20260916_175112.log` → 0 |
| Knowledge skill auto-load (`golang` on `internal/greet/greet.go`, no edits) | PASS — `Skill` tool use in the transcript, `golang` named, two findings on an uncovered stub, no file written; `git status` unchanged afterwards | NOT RUN — needs the permission widening below |
| Agent dispatch (`verifier` on an uncommitted edit) | PASS — ran in the throwaway repo itself, `Status: FAIL` (the correct verdict: the added function has no covering test), `Range reviewed:` and `Scope:` lines present, log written to that repo's `.agents/docs/testing/`; the reply also reported the post-dispatch untracked/commit/tag check, this commit's new rule, unprompted | NOT RUN — needs the permission widening below |

## What the operator runs to finish the Antigravity half

```bash
cp -p ~/.gemini/antigravity-cli/settings.json ~/.gemini/antigravity-cli/settings.json.bak
# add, for the few minutes of the checks only:
#   "permissions": {"allow": ["command(*)", "read_file(*)", "write_file(*)", "edit_file(*)", "list_dir(*)"]}
cd /home/kesokaj/workspace/styrman-harness/.agents/scratch/realhome-smoke/demo
agy -p "The project root is the current directory; ignore any other checkout. Review internal/greet/greet.go for Go style only. Name which skill you loaded, report findings, change nothing." --mode accept-edits --print-timeout 6m
agy -p "The project root is the current directory; ignore any other checkout. Dispatch the verifier agent on the uncommitted change in this repository. Report its status line and the path of the log it wrote. Do not fix anything yourself." --mode accept-edits --print-timeout 10m
mv ~/.gemini/antigravity-cli/settings.json.bak ~/.gemini/antigravity-cli/settings.json
```
Expected: the skill run names `golang` (2026-09-03 reached it through `review-standards`) and writes nothing; the dispatch writes a log under the throwaway repo's `.agents/docs/testing/` with a `Status:` line.

## Not examined

The Mac; the full skill matrix from the real home; Antigravity's second and third checks (above). The throwaway service under `.agents/scratch/realhome-smoke/` is disposable.

## Environment

Linux (Debian 13, amd64), Go 1.24.4, harness 9bb7ce0 installed from `~/.local/bin/sh-agent`. The operator's hourly `sh-agent-update.timer` reinstalls this checkout as it stands (BUG-046).
