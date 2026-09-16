---
owner: reviewer
date: 2026-09-16
---
# Review — bug-043-guide-wording   (2026-09-16)

Range reviewed: working tree, agents/verifier.md
Blob sha (working tree content): 0c97f0b300c2fdc98472ad98d5b0b3cd26ac8842

## Verdict: REQUEST CHANGES

The edited line (line 54) does close the loophole — the duty is now unconditional on `guide.md`'s own absence, with "however many earlier logs `testing/` holds" explicitly overriding the old evidence-of-a-first-time reading. But the file was not read whole before this diff was cut: line 97, in the same file's own **Records** section, still says the guide is written "when the procedure is new or its steps changed" — the exact ambiguous framing BUG-043 fixed, left standing eight lines below its own log entry as a second, contradicting statement of the same duty. An agent that treats Records (the section literally titled "what you write") as authoritative can read "the procedure is new" as false — this project has earlier logs — and skip the guide even though it is absent, reproducing BUG-043 from inside the same file that just fixed it.

## Findings
- **[high / confidence: high]** `agents/verifier.md:97` — "Writes: ... and `<docs-root>/testing/<app>/guide.md` when the procedure is new or its steps changed" contradicts the fixed rule at line 54 ("whenever `$docs_root/testing/<app>/guide.md` does not exist ... however many earlier logs `testing/` holds"). "The procedure is new" is not defined and reads exactly like the pre-fix "first time" trigger that let a run with an older log skip the write. Fix: state the same condition once (file absent, or steps changed) so the two sections cannot diverge — e.g. "and `<docs-root>/testing/<app>/guide.md` when it does not exist or its steps changed".
- **[low / confidence: medium]** `agents/documentation.md:35` and `skills/init-project/SKILL.md:75` — both still describe the verifier's duty as writing the guide "the first time it runs a project's procedure," the pre-fix phrasing. Not part of this diff and not verifier.md itself, so not a contradiction inside the reviewed file, but a describing-text hit: these two locations were not updated alongside the behaviour change and can seed the same misreading in an agent that consults documentation.md or the init-project skill instead of verifier.md directly.
- **[low / confidence: low]** `context/AGENTS.md:17` — "(created by the verifier on first use, maintained by both)" is a shorter parenthetical using "first use," same family of phrasing but terse enough that it reads more like a label than a trigger condition; flagged for completeness, not treated as a live contradiction.

## Checks requested by the dispatch
- Unconditional reading of the diffed line: line 54 alone is unconditional on absence — pass. But line 97 in the same file reintroduces the condition — see high finding above, so the file as a whole still contains a reading under which the guide can be skipped.
- Contradiction with `context/AGENTS.md` step 3 ("created by the verifier on first use, maintained by both"): no direct contradiction — that line is compatible with "write it whenever absent," it just doesn't state the "however many earlier logs" clarification. Not a blocking contradiction.
- Contradiction with `agents/documentation.md`: yes, in spirit — see the describing-text finding above (line 35, "the first time it runs a project's procedure").
- `go run ./cmd/sh-agent list` — ran; output lists all 7 subagents including `verifier`, and 27 skills including `review-standards`. Frontmatter still parses.

## Not examined
`agents/verifier.md` was read in full; grep covered `context/`, `agents/`, `skills/`, `README.md` for "first use", "first time", "guide.md" (README.md had no hits). Did not examine other agent files beyond the described/grep scope, and did not run the verifier live against a sandbox project to observe behaviour end to end — this review is text-only, as scoped by the dispatch.

## Outcome

- Finding 1 (high) — fixed: `agents/verifier.md` § Records now reads "whenever it is absent or its steps changed", matching § What to run.
- Finding 2 (low) — fixed: `agents/documentation.md` § Your lane ("whenever it runs a project's procedure and finds the guide absent") and `skills/init-project/SKILL.md` Step 4 ("the first time it runs the project's procedure and finds it absent").
- Finding 3 (low) — fixed as well: `context/AGENTS.md` step 3 now says "written by the verifier whenever it is absent" (rendered context stays at 93 lines).
- No restatement of the old duty remains (`grep -rnE 'first use|first time it runs|procedure is new' context/ agents/ skills/*/SKILL.md README.md` is empty). Round 1 closed 2026-09-16; this is the recurrence the retro's Proposal 5 addresses.
