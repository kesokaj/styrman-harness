---
owner: main
date: 2026-09-17
---
# Verification — real-home smoke checks after the second retro, both CLIs   (2026-09-17)

Range reviewed: harness `ce8ed6e` installed in the real `$HOME` (`sh-agent update`: 9 written, 251 unchanged; `status`: binary `ce8ed6e8ab1a`, drift none)

## Status: PASS — three checks, both CLIs

Procedure as in `2026-09-17-verify-real-home-smoke-3.md`: fresh demo clones `.agents/scratch/realhome-smoke/smoke4-c` (Claude Code) and `smoke4-a` (Antigravity), one uncommitted doc-comment edit, no testing guide. The context question added "(7) who may write the changes/ entry according to Definition of Done step 4". Antigravity's two tool checks used the per-run `--dangerously-skip-permissions` flag, approved by the operator for these runs; `settings.json` not modified.

| Check | Claude Code | Antigravity |
|---|---|---|
| Context, no tools | PASS, exit 0 — heading, 27 skills and 7 subagents, the guide at `docs/testing/guide.md`, (6) "testing/guide.md is absent: write it first", (7) the `changelog` agent or the main agent from the diff | PASS, exit 0 — the same answers for (6) and (7); 0 `Failed to parse` lines |
| Skill load (`golang`, change nothing) | PASS, exit 0 — `Skill:golang`; the doubled doc comment; no edits | PASS, exit 0 — `golang` named; no edits |
| Agent dispatch (`verifier`) | PASS, exit 0 — the dispatch carried "testing/guide.md is absent"; `Status: PASS`; `Scope: complete — testing/guide.md was absent, written … before running commands`; guide on disk | PASS, exit 0 — `Status: PASS`; `Scope: complete (guide written to testing/guide.md …)`; guide on disk |

Both demos afterwards: only the dispatch-target edit, one log and the guide; no new commits.

## Observation, no bug number

Claude Code's verifier named its log `2026-09-17-greet-doc-comment.md`; `agents/verifier.md` names the log `<date>-verify-<slug>.md`. Cosmetic — the log is found and cited — recorded here in case it recurs.

## Not examined

The full skill matrix; the Mac; the other places.

## Environment

Linux (Debian 13, amd64), real `$HOME`, harness `ce8ed6e`, Claude Code 2.1.274, Antigravity 1.2.5.
