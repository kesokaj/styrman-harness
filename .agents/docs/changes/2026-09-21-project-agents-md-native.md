---
owner: main
date: 2026-09-21
---
# Projects use `AGENTS.md` alone: no project `CLAUDE.md` stub     (2026-09-21)

## Summary

`/init-project` no longer writes a project `CLAUDE.md`: Claude Code 2.1.277+ reads `AGENTS.md` natively when no `CLAUDE.md` sits in or above the project (measured on 2.1.278, `testing/2026-09-21-verify-agents-md-native.md`). This repo's own `CLAUDE.md` is gone. REQUIREMENTS R2 says so in the operator's approved words, with ADR-036 superseding the project half of ADR-005. The user-level `~/.claude/CLAUDE.md` stub stays, measured as required: a user-level `AGENTS.md` alone loads only as an ancestor file and is dropped by any project carrying its own `CLAUDE.md`. Plan `plans/2026-09-21-project-agents-md-native.md`, cold-read three rounds.

## Type

feature (skill behaviour, requirements, docs)

## What changed

- `skills/init-project/SKILL.md` (description, template loop, the `CLAUDE.md` bullet: not written; an existing one never removed; the by-hand stub named for Claude before 2.1.277, Bedrock, Vertex, Foundry, telemetry off); `skills/init-project/assets/CLAUDE.md` deleted.
- `CLAUDE.md` at this repo's root deleted. `context/AGENTS.md` Project layout: a root is `AGENTS.md` alone (size 20,328 of 22,000).
- `REQUIREMENTS.md` R2 rewritten (approved text) with ADR-036 in one commit; `.agents/docs/testing/skill-matrix.md` `init-deploy` row no longer asks for `CLAUDE.md` and checks it is absent.

## Verification

`testing/2026-09-21-project-agents-md-native.md` (`verifier`) PASS; `testing/2026-09-21-verify-agents-md-native.md` (the five probes); e2e: this checkout answered from its own `AGENTS.md` twice; sandbox home reinstalled, drift none, both CLIs quoted the new Root line and the canary, Antigravity store no truncation marker; headless `/init-project` scaffolded a probe with no `CLAUDE.md` (that one run made no init commit — one observation, outside this plan). After the docs review's punctuation fix: size test and scratch install re-run bare, green.

## Review

Source diff under 50 lines, read by the main agent; docs lens `reviews/2026-09-21-review-v1.5.0-HEAD-docs.md` APPROVE WITH NITS, the nit fixed.

## Agents dispatched

`reviewer` ×4 (three cold-read rounds, docs lens), `changelog` (ADR-036), `verifier`.
