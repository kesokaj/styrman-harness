---
owner: reviewer
date: 2026-09-30
---
# Review — worktree-retro-docs   (2026-09-30)

Range reviewed: working tree — changed files: AGENTS.md agents/verifier.md skills/bug/SKILL.md skills/feature/SKILL.md skills/golang/SKILL.md skills/review-standards/SKILL.md; untracked, not ignored: none; read the diff with git diff HEAD -- <files>.

## Status: COMPLETE

## Verdict: APPROVE

Clean. All eight additions in this diff are same-day rule-graduations from the four bugs recorded today (BUG-109, BUG-110, BUG-111, BUG-112 in `.agents/docs/bugs/INDEX.md`) and from the operator's 2026-09-29 correction (feature skill). Each new sentence was grepped for a living twin — none found with different wording, none found dropped or narrowed, and none contradicts the four texts the brief named (the frozen-tree rule, the Rounds rule, `review-standards` § 2's test rules, `preflight`'s Step 2 parallel dispatch). The `agents/verifier.md` addition actually resolves a latent tension rather than creating one: `preflight` Step 2 dispatches `verifier`, `reviewer` and `security` in parallel, so `reviewer`/`security` writing their reports under `<docs-root>` while `verifier` runs would previously look indistinguishable from "the tree moved" under the frozen-tree rule; the new bullet correctly scopes "moved" to source/tests/config and treats a docs-root or `.agents/scratch/` write as another agent's artifact, not tree movement.

## Findings

None. No stale restatement, no dropped rule with no new home, no narrowed rule, and no contradiction between two living texts.

## Checks performed (per new sentence)

- **AGENTS.md — vendored-skill pin verification** (`git archive <pin> <path> | diff -r`, `TestVendoredRowsCarryAPin`, BUG-111): matches `skills/VENDORED.md`'s already-updated second table (pins, "Local changes" column, provenance note) and `bugs/INDEX.md` BUG-111 row. No other file restates the old "never edit vendored content in place" rule in conflicting terms.
- **AGENTS.md — redaction needles grow by KIND** (BUG-110, `TestReportRedactShapes`, no leading `\b`): consistent with the pre-existing, more detailed enumeration in `.agents/docs/deploy/local.md` ("matched without a leading word boundary, so a token after a literal `\n` or an underscore is still caught") and `.agents/docs/specs/sh-agent/delivery.md`'s redaction paragraph — same fact, not a contradiction. Matches `changes/2026-09-30-0910-redaction-per-kind-rules.md` and the BUG-110 index row.
- **agents/verifier.md — tree-under-test is source/tests/config only**: no other living doc restates the frozen-tree rule (`context/AGENTS.md` § Definition of Done step 1 is the only place it lives); this bullet narrows "moved" correctly rather than contradicting it (see Verdict above).
- **skills/bug/SKILL.md — class fix, `git grep` the same call** (BUG-109): matches `bugs/INDEX.md` BUG-109 row ("grep for the same call elsewhere once one site is fixed") and `skills/golang/SKILL.md`'s new glob rule (narrower, file-specific instance of the same lesson — consistent, not conflicting).
- **skills/feature/SKILL.md — capability consumers** (operator's 2026-09-29 correction): no existing "consumers" rule elsewhere states this differently; `skills/terraform/SKILL.md:281` and `delivery.md:46` use "consumers" in unrelated senses (Terraform output docs, a different capability question) — no conflict.
- **skills/golang/SKILL.md — no untrusted path in `filepath.Glob`**: matches BUG-109's prevention rule verbatim in substance; no other Go-review text (`review-standards/references/go-review.md` not touched by this diff) states a different rule for glob use.
- **skills/review-standards/SKILL.md — order-pinning exception for unobservable order**: matches BUG-112's prevention rule ("pin any order no behavioural test can observe") and does not weaken the adjacent, unchanged must-fail-mutation bullet — it is a narrow, explicitly-scoped exception to the "every test can FAIL" bullet, not a restatement of the mutation protocol.
- **skills/review-standards/SKILL.md — consolidated step-sequence order check**: matches BUG-112 and its `changes/` entry (`WriteAtomic` reorder). No other criterion or checklist item restates step-sequence ordering differently.

Also checked: `REQUIREMENTS.md` for glob/redaction/vendor/consumer/order/consolidation terms — no requirement touches these, no conflict. The 22,000-byte context-size test (`TestRenderedContextFitsAntigravityBlock`) scopes to `context/*.md` and `rules/*.md` only, per `AGENTS.md`'s own Setup commands section — not implicated by additions to the root `AGENTS.md` or `skills/`.

## Criteria

All six marked `n/a — docs lens` per the brief; this is a docs-consistency review, not a code review.

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- `README.md`'s Install/table sections were checked only for terms directly touched by this diff (redaction, glob, consumers, vendoring); I did not re-read `README.md` end to end for unrelated drift.
- `.agents/docs/testing/guide.md` was grepped for the specific terms this diff touches, not read end to end.
- The four bug records themselves (BUG-109 through BUG-112) were read only via `bugs/INDEX.md` rows, not their full `docs/bugs/BUG-NNN-*.md` files (if any exist beyond the index) — the diff's new sentences were checked against the index summaries, which matched in every case, so this gap is unlikely to hide a contradiction but is not fully ruled out.
- No code was reviewed (out of scope for a docs lens); this report covers only the six changed markdown/skill files against the living docs named in the brief.
