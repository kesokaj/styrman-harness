---
owner: main
date: 2026-09-21
---
# Verification — does Claude Code read `AGENTS.md` natively, at which scope   (2026-09-21)

Range reviewed: no code; the vendor behaviour behind `plans/2026-09-21-project-agents-md-native.md`. Claude Code 2.1.278 on this box (`claude --version`); the documented threshold is v2.1.277 (memory page § AGENTS.md and the 2.1.277 changelog entry, both fetched 2026-09-21).

## Status: PASS — the plan's premise holds at project scope; the global stub is required

Every probe: `claude -p '<question>' --max-turns 1`, no tools, the answer read from the output. Credential copies for a probe home were `chmod 600` and deleted after the run (guide § 4); every probe directory was deleted after.

| # | Arrangement | Result |
|---|---|---|
| 1 | Real home; `~/workspace/agentsmd-probe-project` with `AGENTS.md` only (probe word HERON-4412); no `CLAUDE.md` in or above it | run 1: "No `AGENTS.md` or `CLAUDE.md` from the project directory appears in my context" — the documented first-session-after-upgrade case; run 2: "The probe word is HERON-4412, from …/AGENTS.md" — **read natively** |
| 2 | Same directory with the harness stub `CLAUDE.md` = `@AGENTS.md` | "HERON-4412, from the project's AGENTS.md (loaded via the @AGENTS.md line)" — the import form still works |
| 3 | Probe home OUTSIDE `/home` (the session scratchpad under `/tmp`), `~/.claude/AGENTS.md` only, no `~/.claude/CLAUDE.md`; project `AGENTS.md` under that home | runs 1 and 2: both files loaded — probe word right, global context ends with GANNET, first heading `# AGENTS.md — global agent context`, both paths named — the global file is read as an ancestor `.claude/AGENTS.md` |
| 4 | Arrangement 3 plus a project `CLAUDE.md` with content | only the project `CLAUDE.md` given; "No global context ending in GANNET" — **the global context is dropped** (CLAUDE.md-only mode) |
| 5 | Arrangement 3 with the one-line `~/.claude/CLAUDE.md` stub restored | global context loads (GANNET, heading) — the user-level stub is exempt from the check |

Invalid arrangements, recorded so they are not repeated (BUG-087's rule): a probe home under `/home/kesokaj` (`~/workspace/sbhome`, or `.agents/scratch/sandbox-home`) — the real `~/.claude/CLAUDE.md` is then an ancestor `.claude/CLAUDE.md` of the probe project, Claude switches to CLAUDE.md-only mode, and its `@AGENTS.md` import (external to the project) is not expanded headless; every such run reported `/home/kesokaj/.claude/CLAUDE.md` as its only context whatever `HOME` said.

## Conclusion

Project scope: a project needs only `AGENTS.md` on Claude Code 2.1.277+ where feature flags are fetched (not Bedrock, Vertex, Foundry, not telemetry disabled), from the second session after an upgrade. Global scope: `~/.claude/AGENTS.md` alone loads only as an ancestor file and is lost the moment a project carries a `CLAUDE.md` or sits outside the home directory; the one-line `~/.claude/CLAUDE.md` stub stays (plan D3, D5).

## Not exercised

Antigravity (unchanged: reads `AGENTS.md` natively at both levels, ADR-005); Claude Code before 2.1.277; the devbox.
