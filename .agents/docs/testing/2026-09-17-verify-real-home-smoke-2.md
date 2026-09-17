---
owner: main
date: 2026-09-17
---
# Verification — real-home smoke checks after ADR-032, both CLIs   (2026-09-17)

Range reviewed: harness `e67eace` installed in the real `$HOME` (`sh-agent update`: 15 written, 245 unchanged; `status`: binary `e67eace9bcbc`, drift none)

## Status: PASS with one defect — five checks clean, Claude Code's dispatch skipped the guide and said so

Procedure as in `2026-09-17-verify-real-home-smoke.md`: fresh clones of the demo service (`.agents/scratch/realhome-smoke/smoke2-c` for Claude Code, `smoke2-a` for Antigravity), one uncommitted doc-comment edit, no testing guide. Antigravity's two tool checks used the per-run `--dangerously-skip-permissions` flag, approved by the operator for these runs; `~/.gemini/antigravity-cli/settings.json` not modified. Claude Code's dispatch ran with an allowlist including `Write`, `mkdir`, `echo` and `tee`; its transcript shows no denied call.

| Check | Claude Code | Antigravity |
|---|---|---|
| Context, no tools | PASS, exit 0 — heading, 27 skills and 7 subagents, the guide at `docs/testing/guide.md`, one per project; the escape and dispatch-budget rules | PASS, exit 0 — the same five answers; 0 `Failed to parse` lines |
| Skill load (`golang` on `greet.go`, change nothing) | PASS, exit 0 — `Skill:golang`; the doubled doc comment and the missing package comment; no edits | PASS, exit 0 — `golang` named; the same finding; no edits |
| Agent dispatch (`verifier`) | PASS on the dispatch, DEFECT on the guide — `Agent:verifier`, `Status: PASS`, log written; no `testing/guide.md` written, no attempt, no denial; Scope `guide not written: no testing/guide.md existed for this project and this run did not add one … Scope is not complete` | PASS — waited for its subagent; `Status: PASS`, `Scope: complete`, `.agents/docs/testing/guide.md` written |

## Defect

BUG-063's reporting half holds: the missing guide is named and Scope is not complete. Its duty half does not hold on Claude Code in this run: step 4 says to write the guide before running anything, nothing refused the write, and the verifier did not write it. Earlier today with the same definition passed by `--agents`, repro-11 wrote it. Recorded as a bug.

## Not examined

The full skill matrix; the Mac; the other places.

## Environment

Linux (Debian 13, amd64), real `$HOME`, harness `e67eace`, Claude Code 2.1.274, Antigravity 1.2.5.
