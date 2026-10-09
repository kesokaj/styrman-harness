---
owner: main
date: 2026-10-09
---
# Handoff — 2026-10-09 — project-review fixes done, verified; a third review round is the operator's call

Session ID: 728a4838-8b3c-4819-a62d-2b01e563a77f
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

`/review --project --all --high` at c267990 (`reviews/2026-10-09-review-project-merged.md`, 29 rows). The operator approved the recommendations by row. The fixes are `3bffdcf` and `6cec23b`, verified (PASS) and pushed. v2.6.2 stays the live tag; nothing here is released.

## Completed

- 21 rows fixed, 3 struck, 5 deferred: the record's § Outcome, one line per row. Bug records BUG-167 (delivery spec drift), BUG-168 (usage errors exit 1), BUG-169 (dream CLI grandchild).
- Review rounds 1 and 2 (reviewer code+docs on opus, security), the record's § Round 1 and § Round 2; every finding fixed with a test red first, except the agent-key allow-list (left: operator).
- Verifier: `testing/2026-10-09-verify-review-project-fixes.md`, PASS at `6cec23b`.
- Change entry: `changes/2026-10-09-1921-review-project-fixes.md`.

## Open / In Progress

- Nothing in progress. The review cap is reached: round 2 landed on the needle round 1 reshaped, its fixes are applied and verified, and no third review round has run (Next Steps 1).
- Requirements edits awaiting approval: none.

## Key Decisions

- Rows 7, 26, 27 struck with checkable reasons (the record's § Outcome); rows 21, 23, 24, 25, 28 deferred.
- `OriginKind` and `AgyConversation` stay exported (types of exported API).
- The delivery paragraph names only what a tag shipped; the release commit updates it, and `TestDeliverySpecCurrent` in the release procedure refuses the tag otherwise.
- The redaction needles grow by kind: command-line credential (the flag after whitespace), auth key, credentials key.

## Corrections

- None from the operator. Two from the reviews: the simplify row's "identical behaviour" claim for the install pass-through was false (fixed), and the unanchored `-u` rule over-matched the trail (fixed, measured: 83 lines before, 26 after, all review records quoting credential shapes).

## Agent notes (unverified)

- The verifier lists weak coverage: `readSmall`'s 1 MiB cap has no fixture (the cap predates the change); `search.IsUsage` and dream's `docsOf` are covered through the verb tests only.
- The security agent did not verify against a live CLI that Claude Code honours `mcpServers`, `memory` or `initialPrompt` in a user-level agent file.

## Environment State

- Branch: main, clean, pushed at the commit that carries this handoff. Tag v2.6.2 unchanged. `.agents/scratch/` holds the bare test logs and the scratch install (gitignored).

## Artifacts to Read

- `reviews/2026-10-09-review-project-merged.md`: the 29 rows, § Round 1, § Outcome, § Round 2 with the cap question.
- `testing/2026-10-09-verify-review-project-fixes.md`: the verifier's PASS and its coverage notes.

## Next Steps

1. Operator, at the review cap (recommendation first): (a) accept the round-2 fixes on their tests and the verifier, closing the task; (b) run one more round over `6cec23b` (reviewer and security, the round-2 fixed lines only), adding a `readSmall` 1 MiB fixture test in the same pass; (c) revert the three credential kinds to their round-1 shape.
2. Operator decisions carried from the review: the agent-key allow-list (security round-1 finding 5: `memory`, `initialPrompt` pass the deny-list; an allow-list refuses a future vendor key by default); a `PRIVATE_KEY=` redaction kind (not covered today); the deferred rows 21, 23, 24, 25, 28 (the R2 pointer on release and zip places, row 28 L2, is the one worth a look).
3. Operator decisions carried from the v2.6.2 handoff: whether `sh-agent dream` stops at the first agy login failure; whether the Mac's checkouts should be dreamt over (`sh-agent search <word>` once in each).
4. A release carries these fixes to the release and zip places (`/preflight` from the operator's turn); checkout places get them at their next `sh-agent update`.
