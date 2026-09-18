---
owner: reviewer
date: 2026-09-18
---
# Review — worktree-retro-2-docs (2026-09-18)

Range reviewed: working tree, `git diff HEAD -- context/AGENTS.md agents/reviewer.md agents/security.md agents/explorer.md skills/preflight/SKILL.md AGENTS.md README.md .agents/docs/deploy/remote.md .agents/docs/specs/sh-agent/delivery.md` (round 1); round 2 re-checked the same working-tree scope after the fixes below landed.

## Verdict (round 2): APPROVE

Both round-1 findings are fixed and verified. `README.md:59` and `:61` now point at `delivery.md` § The model instead of restating the delivery claim; `README.md:19`'s mermaid label ("main: what checkout places follow") is a map label, not a restatement, which project `AGENTS.md` rule 6 (now amended: "the README map draws the paths and its labels name them, which is not a restatement") explicitly excepts. No qualifier from the old README sentences — "refuses a dirty checkout", "replaces `~/.local/bin/sh-agent`", "so that place changes only when you cut a release" — went missing: each is still stated verbatim or near-verbatim in a named home (see Round 2 verification below).

## Findings (round 1, resolved — see Outcome)

- **[medium / confidence: high]** `README.md:61` — "Checkout places follow `main`; binary-only places update to the newest release." was a full restatement of the delivery model, left un-pointed while this same diff fixed the identical claim eight lines below it at `README.md:245`. RESOLVED: now "what each kind of place does with it is `.agents/docs/specs/sh-agent/delivery.md` § The model."
- **[low / confidence: medium]** `README.md:59` — the Update bullet fully restated per-place mechanics instead of pointing at § The model. RESOLVED: now "what that means for a checkout place and for a binary-only place is `.agents/docs/specs/sh-agent/delivery.md` § The model. Nothing pushes to a machine."

## Round 2 verification — qualifiers from the old README sentences

Checked with plural-form greps (`checkout places? follow`, `binary-only places? update`) across `README.md`, `AGENTS.md`, `agents/`, `context/`, `skills/`, `.agents/docs/specs/`, `.agents/docs/deploy/`: no stale restatement found anywhere except the mermaid map label at `README.md:19`, which rule 6's amended text now names as an intentional exception, not a restatement.

- "refuses a dirty checkout" → present in `delivery.md:74` ("`update` refusing a dirty checkout on every path"), `README.md:179` (Daily commands table), `.agents/docs/deploy/remote.md:73-74`, `.agents/docs/deploy/local.md:28`. Not narrowed.
- "replaces `~/.local/bin/sh-agent`" → present in `.agents/docs/deploy/local.md:43`, `.agents/docs/specs/sh-agent/delivery.md:37`, `README.md:163` (Daily commands table), `.agents/docs/deploy/remote.md:59`. Not narrowed.
- "so that place changes only when you cut a release" → present as "changes only when a release is cut" in `delivery.md:37` and `.agents/docs/deploy/remote.md:59`. Not narrowed.

No qualifier is now stated nowhere in the living docs; every one has at least one home among `delivery.md` § The model / § How a machine gets the harness, `deploy/local.md`, `deploy/remote.md`, or the README § Daily commands table.

Also spot-checked the items the coordinator flagged as changed since round 1: `AGENTS.md` rule 1 and rule 6 read consistently with `delivery.md` § The model's content (rule 6 now also carves out the map-label exception, which matches `README.md:19` correctly); `delivery.md:36-37` (`§ How a machine gets the harness`, items 1-2) point at § The model and add no new unrepointed restatement; did not re-diff the three agents' `## Status` lines in this round (out of scope for the delivery-claim check this round asked for — see Not examined).

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens; both restatement chains now resolve to a single home, traced above
3. Standards conformance — pass, round 2: no stale exact-old-wording or unrepointed restatement of the delivery claim found anywhere in the living-docs set (README, AGENTS.md, agents/, context/, skills/, `.agents/docs/specs/`, `.agents/docs/deploy/`)
4. Directory layout — n/a — docs lens
5. Simplicity — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- The three agents' `## Status` line changes mentioned by the coordinator as "also changed since your read" — not diffed in round 2; the round-2 brief scoped the check to the README qualifiers and plural-form grep, not this.
- `skills/preflight/SKILL.md`'s "showing is not a gate" change and the write-early agent changes — covered in round 1 (see prior findings text below), not re-examined in round 2.
- `.agents/docs/architecture/DECISIONS.md`, `changes/`, `handoffs/`, `bugs/INDEX.md`, `retros/` — intentionally excluded both rounds per the dispatch's own instruction (append-only/dated records).
- No command run in either round — pure text review as scoped ("report only", no source edits, no network, no real `$HOME`).

## Round 1 findings text (superseded, kept for trail)

- `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md` "write early at a third of budget" changes — cross-checked against `agents/verifier.md`, `agents/bugs.md`, `agents/changelog.md`, `agents/documentation.md`; `verifier.md` already exceeds the new rule with its own per-command write discipline (not a gap), and `bugs.md`/`changelog.md`/`documentation.md` never carried this rule so nothing fell out of sync.
- `skills/preflight/SKILL.md`'s "showing is not a gate" change — no internal contradiction found in the surrounding numbered steps.

## Outcome

- Finding 1 (`README.md:61` restatement): fixed — pointer to `delivery.md` § The model.
- Finding 2 (`README.md:59` Update bullet): fixed — "what that means for a checkout place and for a binary-only place is … § The model"; "Nothing pushes to a machine" kept.
- Root cause of both: the sweep grep matched only the singular ("checkout place follows"); it now includes plural forms.
