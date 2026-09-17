---
owner: reviewer
date: 2026-09-17
---
# Review — retro docs (context/AGENTS.md + agents/changelog.md)   (2026-09-17)

Range reviewed: `git diff HEAD -- context/AGENTS.md agents/changelog.md` (uncommitted worktree changes; six inserted sentences)

## Verdict: APPROVE WITH NITS

Nothing broken and no contradiction found across the six inserted sentences, the three cited bug records, or the byte-level escape claim. The one real gap: the new rule requiring a dispatch to `reviewer`/`security`/`verifier` to name a tool-call budget and repeat "write the artifact early" is stated only in `context/AGENTS.md` — the two skills that actually build those dispatch prompts (`skills/review/SKILL.md` Step 3, `skills/preflight/SKILL.md` Step 2) still describe their dispatch text without either element, so the rule is unenforced by its primary callers.

## Findings

- **[medium / confidence: high]** `context/AGENTS.md:33` (new sentence: "A dispatch to `reviewer`, `security` or `verifier` names a tool-call budget below the agent's `maxTurns` and repeats 'write the artifact early'...") vs. `skills/review/SKILL.md` Step 3 ("Each dispatch gets the scope text, today's date, its slug, its brief, and the sentence 'Write your artifact with the Write tool; report only means no source edits.'") and `skills/preflight/SKILL.md` Step 2 ("dispatch in parallel, each with the range `<base>..HEAD`, today's date, and the instruction to record a `Range reviewed:` line"). Neither skill's dispatch-text template names a tool-call budget or repeats "write the artifact early" — the two places that automate these dispatches don't implement the rule the diff just added. As written today, `/review` and `/preflight` dispatches will keep omitting both elements. Fix: add the budget + early-write reminder to both Step-3/Step-2 dispatch-text templates in the same change that adds the rule (or note in the new sentence that the skills still need updating, so the gap isn't silently discovered later, similar to BUG-061's "not done" follow-up note in the handoff).

## By inserted sentence

1. `agents/changelog.md:82` (BUG-061, factual-claim tracing) — accurate: matches BUG-061's record and the handoff's explicit "not done" follow-up asking for this graduation into `agents/changelog.md`. No contradiction with the rest of the Rules block or `skills/review`'s Outcome-writing text.
2. `context/AGENTS.md:16` DoD step 2 (BUG-058, "every trigger... test that fails without the fix") — accurate against BUG-058's record (round-2 Outcome marked a two-trigger finding fixed on its headline case only). Additive, doesn't conflict with the existing "fixed or rebutted in writing" sentence or `skills/review/SKILL.md`'s Outcome definitions.
3. `context/AGENTS.md:18` DoD step 4 (BUG-061, main agent reads changelog/ADR against diff before committing) — accurate against BUG-061 and the handoff's Next Steps line naming this exact graduation as pending. Consistent with step 5's commit ordering.
4. `context/AGENTS.md:33` (dispatch tool-call budget + early-write repetition) — see Finding above: accurate as a standalone instruction (matches how this very review was dispatched: "about 20 tool calls" against `reviewer`'s `maxTurns: 40`), but not yet wired into the skills that build these dispatches.
5. `context/AGENTS.md:51` (BUG-057, escape-to-character conversion) — accurate against BUG-057's record; complements rather than restates the existing binary-conversion sentence in the same bullet (opposite failure direction, correctly framed as "the opposite happens too"). Byte-check passed (see below).
6. `context/AGENTS.md:52` (REQUIREMENTS.md as the home for rules that must not drift) — consistent with the existing "immutable without approval... propose the exact edit → approval → ADR" sentence in the same bullet, and with `skills/init-project/SKILL.md:49` ("agents never edit it") and `skills/retro/SKILL.md:45`'s REQUIREMENTS.md-edit-awaiting-approval handling. No conflict found in `skills/feature/SKILL.md`.

## Bug IDs

- BUG-057, BUG-058, BUG-061 all exist under `.agents/docs/bugs/` and each cited sentence's attribution matches its record's Symptom/Root cause/Prevention text (checked verbatim above).
- `context/AGENTS.md` is deployed to every project (per its own header), yet the file already cites this-harness-only bug IDs (BUG-012, BUG-031, BUG-038, BUG-011, BUG-014/038/044/045) before this diff. The three new citations (BUG-057/058/061) follow the same existing precedent — not a new problem introduced by this diff, just continuing a pattern that was already there.

## Byte check (Hard rules / BUG-057 sentence)

`grep -o '.\{0,5\}u0041.\{0,10\}' context/AGENTS.md | od -c` shows the literal six bytes `\`, `u`, `0`, `0`, `4`, `1` — confirmed NOT converted to the character `A`. The sentence's own claim ("no byte count shows [the conversion]") is demonstrated true by this check having to go byte-level to catch it.

## Length / restatement check

None of the six sentences is a pure restatement of an existing rule — each adds a new, checkable criterion (a specific "fixed" bar, a specific tracing bar, a specific dispatch requirement, a specific byte-check trigger, a specific REQUIREMENTS.md-first trigger) grounded in a named bug or dated incident. No candidate for trimming found.

## Not examined

- Full contents of `skills/VENDORED.md` cross-check beyond the grep for the exact subject nouns used above; other non-vendored skills not grepped individually beyond `feature`, `init-project`, `retro`, `review`, `bug` (as the dispatch named) and `preflight` (checked because it shares the reviewer/security/verifier dispatch mechanism with `review`).
- `agents/bugs.md`, `agents/documentation.md`, `agents/explorer.md` — not read; no cited sentence names them directly, but they weren't checked for stray restatements of the six subjects either.
- Did not verify whether other in-flight/uncommitted files besides `context/AGENTS.md` and `agents/changelog.md` reference the same six subjects (scope was limited to the named diff).

## Outcome

- Finding 1 (the tool-call budget rule is not in the skills that build the dispatch text): fixed — `skills/review/SKILL.md` Step 3 and `skills/preflight/SKILL.md` Step 2 now add a tool-call budget below the agent's `maxTurns` and "write the artifact early, then update it" to every dispatch, pointing at the context rule. Same commit; re-verified by a scratch install.
