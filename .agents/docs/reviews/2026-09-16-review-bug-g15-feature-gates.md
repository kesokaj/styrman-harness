---
owner: reviewer
date: 2026-09-16
---
# Review — bug-g15-feature-gates   (2026-09-16)

Range reviewed: working tree, skills/feature/SKILL.md
Blob read: 6d388154b62745564418aeca80cd800c7a1120ad

## Verdict: APPROVE WITH NITS

The diff itself is internally consistent — the new "Answer three things" brief, the "three-round cap" language in Gates, and the Step 2 scope-discipline sentence all agree with each other and with the plan template. The one real problem is outside the diff: it leaves a sibling file's describing text stale, which the global context's describing-text rule treats as a finding, not a nit.

## Findings

- **[medium / confidence: high]** `agents/reviewer.md:28` — describes the cold-reader dispatch as "answer the two questions the dispatch asks." `skills/feature/SKILL.md:77` now asks **three** things (Gaps, Scope, Readiness), not two. This line is untouched by `reviewer.md`'s own diff (checked `git diff -- agents/reviewer.md`), so it is stale relative to the feature-skill change under review. Two readers of `reviewer.md` alone would expect a two-item answer and could treat a Scope-only finding as out of contract. Fix: change "the two questions" to "the questions" (or name the count) in the same change that ships the feature-skill edit.

- **[low / confidence: medium]** `skills/feature/SKILL.md:97` — introduces a new tag `supersedes D<n>` for a review fix that reverses a D-row/default/Failure-modes row, but the Decisions template at line 66 (`D<n> — <date> — <decision> — by: …`) has no field or example showing where `supersedes D<n>` goes. Two careful implementers could place it differently (prefix on the decision text vs. a separate clause vs. abandoning the old D-row entirely rather than superseding it). Not blocking — the convention is legible from context — but the template could show one worked example to remove the ambiguity.

## Checks specific to the dispatch brief

- **Gates / Step 1 / Step 2 / Step 3 / Step 6 mutual consistency:** no contradiction found. "Four points always wait for the user" (Gates, unchanged) still enumerates exactly scope confirmation, the Step 3 no-code gate, the Step 3 three-round cap, and the review-loop cap — all four still exist after the edit. The new Step 2 sentence ("Scope In holds what the request names… settled by ONE shape… `Rejected:`") matches the Decisions template's existing `Rejected:` row and matches the new Step 3 item (2) Scope check, which treats an unrequested extra shape as a gap "unless a D-row says who asked for it." No new clause tells an implementer to proceed past a gate without the user.
- **Numbering/count match:** "Answer three things" is followed by exactly three numbered items (1)(2)(3) — count matches. The final "design ready … implementation ready" criteria still read "zero gaps" / "zero BLOCKING questions"; Scope items are explicitly folded into "gaps" by the new item (2)'s last sentence, so the pass/fail criteria are unchanged in meaning.
- **"Three rounds" / "three-round cap" staleness sweep:** `grep -rn "three rounds\|three-round"` across `context/`, `skills/`, `agents/`, `commands/`, `README.md` shows `skills/feature/SKILL.md:15` (Gates), `skills/review/SKILL.md:62`, and `skills/bug/SKILL.md:15` all still say "three round(s)" and all agree with the new Step 3 wording ("the third report that is not both tokens 'ready' ends the check … there is no fourth round") — same cap, no drift.
- **"Answer two things" / "simplest reading" sweep:** no other file in `context/`, `skills/`, `agents/`, `commands/`, `README.md` contains "Answer two things" (so no other stale copy of the old brief) or "simplest reading" outside this file (the phrase is new, introduced only here).
- **Print-mode readability:** the new Step 3 sentence ("however small the remaining items look… stop and ask the user") still routes to the unchanged Gates paragraph, which explicitly handles the no-reply case (`claude -p`/`agy -p`) by ending the work and writing to the handoff's Next Steps rather than proceeding. No new phrase in the diff (e.g. "your recommendation first") reads as standing permission to continue without a reply in print mode.

## Criteria

1. Dead code — n/a (prose file, no code removed or orphaned).
2. Flow tracing — pass: the new Step 3 brief text is dispatched verbatim to `reviewer`; the reviewer's own six-criteria template is unaffected, but see finding above on `reviewer.md:28` going stale.
3. Standards conformance — fail (one item): global context's describing-text rule ("a behaviour change includes the text that DESCRIBES it… in the same change") is violated by leaving `agents/reviewer.md:28` unupdated.
4. Directory layout — n/a.
5. Simplicity — pass: no new machinery, the added sentences are policy clarifications reusing existing template fields (`Rejected:`, `## Decisions`, D-rows).
6. Decoupling — n/a for a single prose file; the new `supersedes D<n>` convention touches only the plan's own `## Decisions` section.

## Not examined

Did not read the full history of `.agents/docs/reviews/2026-09-15-*` or `2026-09-13-*` review records that discuss the three-round cap's earlier drafting (found via grep, not opened) — they are prior review artifacts, not living docs, so out of scope for the describing-text sweep. Did not check `skills/bug/SKILL.md` or `skills/review/SKILL.md` beyond the grep match shown above (their own review loops were not otherwise re-read line by line). No other files in the working tree were reviewed, per the dispatch's single-file scope.

## Outcome

- Finding 1 (medium) — fixed: `agents/reviewer.md` § Scope now says "answer the questions the dispatch asks, in its numbering", so the brief's count lives in one place (the feature skill).
- Finding 2 (low) — fixed: the Decisions template in `skills/feature/SKILL.md` gained the row form `D<n> — <date> — <decision> — supersedes D<m> — by: <review finding, report path | implementer>`.
- Round 1 closed 2026-09-16; no second round owed (a stale restatement and a template gap, no mechanism finding).
