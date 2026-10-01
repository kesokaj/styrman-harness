---
owner: main
date: 2026-10-01
---
# Verification — install into the operator's real home and the three smoke checks on both CLIs (2026-10-01)

Range reviewed: 2ba4a50 installed (`sh-agent status`: 266 files managed, drift none — 265 files match the repo, binary 2ba4a50351ab without a warning)

## Status: PASS for both CLIs

Procedure: `.agents/docs/testing/guide.md` § 5 (`testing/2026-09-03-verify-real-home.md` is the first run of it). The install was the one-time exception D15 of plan `2026-09-30-audit-fix.md` (the operator's words, 2026-09-30). The three checks ran from the REAL `$HOME` (credentials and token in place, nothing copied) in a fresh clone of the demo service per check (`~/sandbox/audit-fix/runs/rh-<cli>-<check>`, outside every checkout), driven by `.agents/scratch/audit-fix/kit-realhome.sh` (sources `testing/kit-preamble.sh`; not committed, ADR-010).

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools (first heading, skill count and three names, the subagents, Definition of Done step 3, the closing word) | PASS — `# AGENTS.md — global agent context`, 34 skills, the eight agents, step 3 restated, `GANNET` | PASS — the same five answers, `GANNET` |
| Knowledge skill auto-load (`golang` on `greet.go`, change nothing) | PASS — named the skill, no edits (`git diff` clean) | PASS — named the skill, no edits |
| Agent dispatch (`verifier` on an uncommitted comment edit) | PASS — log under `testing/` with `## Status`, 47 s | PASS — log under `testing/` with `## Status`, 86 s; the first Antigravity run of a derived agent since the Phase 4 rewrite |

`auto-denied` lines in stderr: 0 on every run. Commands: `go run ./cmd/sh-agent install` (151 written, 115 unchanged, 52 removed — the struck skills and the replaced references), `hash -r; type sh-agent` → `/home/kesokaj/.local/bin/sh-agent`, `sh-agent status` as above; `grep -c GANNET ~/.claude/AGENTS.md ~/.gemini/config/rules/AGENTS.md` → 1 each.

## Not verified here

The Mac (no access this session; `deploy/remote.md` § Operator steps covers its update); the full matrix from the real home (the matrix runs from the sandbox homes: `testing/2026-10-01-verify-phase{5,6,7}-skill-matrix.md` and the Antigravity run of this phase); Antigravity's interactive prompts (headless only).

## Environment

Linux (Debian 13, amd64), Claude Code and `agy` from `~/.local/bin`, harness 2ba4a50, Go 1.27.1.
