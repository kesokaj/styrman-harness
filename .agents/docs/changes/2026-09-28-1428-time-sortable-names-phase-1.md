---
owner: main
date: 2026-09-28
---
# Time-sortable trail names, Phase 1: rule text     (date: 2026-09-28)

## Summary
Phase 1 of `plans/2026-09-28-time-sortable-trail-names.md`. Handoffs and change entries are now named `YYYY-MM-DD-HHMM-<slug>.md` (`date +%F-%H%M`, local time), so name order is time order in any clone. The rule is in the global context, together with how older date-only names sort. The handoff skill and the changelog agent write the new names. The changelog agent has no Bash, so its dispatcher passes the stamp (D3). The agent's impossible "run `date +%F` yourself" rule is deleted (D4). `/init-project` in update mode counts date-only names and warns, and never renames anything. This entry is itself the first change entry named in the new form.

## Type
feature

## Files changed
- `context/AGENTS.md` — § Project layout: the naming sentence (+169 B, context now 20,498 of 22,000 B).
- `skills/handoff/SKILL.md` — the save line and the artifact table row.
- `agents/changelog.md` — the description (the dispatcher passes the stamp), date rule 1 (the stamp names the file), the Bash-less rule 2 deleted, and the lane comment, entry heading and Records.
- `skills/preflight/SKILL.md` — the changelog dispatch passes the stamp.
- `skills/init-project/SKILL.md` — the update-mode count and warning; Step 7's summary carries it.
- `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` — the plan: cold-read rounds 1–3, D1–D4, Phase 1 ticked.
- `ROADMAP.md` — "Now" names the plan.
- `.agents/docs/reviews/2026-09-28-coldread-time-sortable-names-r{1,2,3}.md` — the cold reads.
- `.agents/docs/testing/2026-09-28-verify-time-sortable-names-phase-1.md` — verifier log.

## Agents dispatched
- `reviewer` ×3 (cold reads). Round 1: design needs revision, because the changelog agent has no Bash (fixed by D3). Round 2: design ready, implementation not ready, because of how the handoff-probe copy is named (fixed in the plan). Round 3: design ready, implementation ready; its note about rule 2 became D4.
- `verifier`: PASS.

## Verification
The verifier ran:
- gofmt, vet and the race suite: green;
- the size test: 20,498 B;
- a scratch install: drift none, and both installed contexts carry the sentence and end with GANNET;
- the describing-text sweep: timed forms only, with the positive control at HEAD printing the four old hits.

Behavioural runs, the sandbox start and E2E are Phase 2 and 3.

## Docs
The skill and agent texts above are the docs of this behaviour.

## Requirements
n/a
