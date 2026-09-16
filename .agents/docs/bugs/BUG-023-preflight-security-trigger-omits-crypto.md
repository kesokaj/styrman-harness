---
owner: bugs
date: 2026-09-15
id: BUG-023
status: fixed
area: skills/preflight
---
# BUG-023: preflight's security trigger list omits crypto, so a crypto-only release skips security review

## Symptom
`skills/preflight/SKILL.md:51` (Step 2): "dispatch the `security` agent whenever a changed path in the release diff touches auth, IAM, secrets, input parsing, or dependencies" — crypto is absent. `context/AGENTS.md:16` and `agents/security.md:3` both include crypto in the set of surfaces that trigger a security dispatch. A release that changes only crypto-related code (and none of auth/IAM/secrets/input-parsing/dependencies) would be read as having "no security row" under preflight's own wording, skipping the review. Found by the 2026-09-15 docs-lens sweep, F76.

## Repro
Not reproduced as an executed skip; found by review — the trigger list at `skills/preflight/SKILL.md:51` was compared against `context/AGENTS.md:16` and `agents/security.md:3`.

## Root cause
The security trigger list is restated locally in `skills/preflight/SKILL.md` instead of referencing the canonical list, and drifted when crypto was added elsewhere.

## Fix
2026-09-16: `skills/preflight/SKILL.md` Step 2's Security bullet now reads
"...touches auth, crypto, secrets, IAM, input parsing, or dependencies...",
matching `context/AGENTS.md` § Definition of Done step 2 and
`agents/security.md:3`.

## Prevention
A trigger-list change greps every restatement of it (here, `context/AGENTS.md`, `agents/security.md`, and `skills/preflight/SKILL.md` all needed the same edit); a restated list should instead point at the canonical source.
Graduated: this rule has lived in `skills/review-standards/SKILL.md` § 3's
last bullet ("a rule, ban, threshold or trigger the diff changes has every
other restatement of it in the repository changed too...") since the
2026-09-16 commit applying Proposal 5 of `retros/2026-09-16-retro.md`.

## Related
2026-09-15 docs-lens sweep, `.agents/docs/reviews/2026-09-15-review-sweep-merged.md` F76.
`.agents/docs/retros/2026-09-16-retro.md` Proposal 5.
