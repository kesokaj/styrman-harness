---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — retro applied; trail closed

Session ID: 4e712c10-b55b-4126-ae29-579c7d94ac09
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

Same session as `.agents/docs/handoffs/2026-09-17-delivery-rework-done-r9-map.md`, which still holds the full picture of the delivery rework and R9. After it, the operator ran the retro and approved all six proposals.

## Completed

- `fd07a29` — the six retro rules: an escape must survive a file write, checked by grepping for it (BUG-057, context § Hard rules); a finding is fixed only when every trigger it names is tested (BUG-058, § Definition of Done step 2); dispatches to reviewer, security and verifier carry a tool-call budget and "write the artifact early" (context § Work, and the dispatch text in `skills/review` Step 3 and `skills/preflight` Step 2 — the docs review's one finding); the changelog agent states only what the dispatch or diff says and the main agent reads its records against the diff (BUG-061, `agents/changelog.md` § Rules and § Definition of Done step 4); a rule the operator wants never to drift is proposed as a `REQUIREMENTS.md` edit first (§ Hard rules). BUG-057, 058, 061 moved to the index's Graduated table.
- Record `.agents/docs/retros/2026-09-17-retro.md`; queue `~/.config/sh-agent/proposals/2026-09-17-retro.md` and `~/.config/sh-agent/retros.jsonl` marked applied (fd07a29).

## Open / In Progress

- Nothing queued in the trail: no unticked plan box, `TODO.md` and `ROADMAP.md` empty.
- Requirements edits awaiting approval: none.

## Key Decisions

- All six retro proposals applied as drafted; see the record's `## Not proposed` for BUG-059, BUG-060 and BUG-053–056, deliberately left as records.

## Corrections

- none logged this session.

## Agent notes (unverified)

- The escape conversion also hit tool-call INPUT: a dispatch prompt and a shell command both had the escape turned into the letter before they ran, so the verifier's escape check was void. The main agent re-checked with a script that builds the needle from character codes. A Claude Code bug report draft for this is queued locally for the operator's `/feedback`.
- The changelog agent's entry for fd07a29 misstated where sentences went and quoted text the diff does not have; corrected against the diff before commit, per the new rule.

## Environment State

- Branch: main, pushed at fd07a29 plus this handoff's commit. Tags: `v1.0.0` only. No running processes.

## Artifacts to Read

- `.agents/docs/changes/2026-09-17-retro-applied.md` — what changed and what was not exercised.
- `.agents/docs/handoffs/2026-09-17-delivery-rework-done-r9-map.md` — the rest of this session.

## Next Steps

1. Operator: at each place, `sh-agent update` (context, agents and skills changed today), `sh-agent status`, then the three real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`).
2. Operator: read what was added under `~/.config/sh-agent/` since the last report, `sh-agent report --dry-run`, then `sh-agent report` — the new report code has not pushed to the real feedback repository yet.
3. Operator: `/feedback` to review and send the queued Claude Code bug report about escapes converted by the file tools.
4. If wanted: a plan to close the map's dashed gaps (reading other places' reports into `/retro --harness`); under R9 it needs approval, an ADR and one commit with the map.
