---
owner: reviewer
date: 2026-09-29
---
# Review — dream-code   (2026-09-29)

Range reviewed: commits `80bf5b6..HEAD` plus the working tree (`git diff 80bf5b6`) over `skills/dream/SKILL.md` (new), `skills/retro/SKILL.md`, `context/AGENTS.md` (one sentence), `.agents/docs/testing/skill-matrix.md` (one row).

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing destructive or unguarded: `dream`'s doubt-resolves-to-keep rule, its `--dry-run`/headless "nothing applied" guard, and its "store is not under git, quote every deletion in full" undo path are all present and consistent with `AGENTS.md`'s destructive-action rule. The one real defect is a cross-file citation bug: `skills/retro/SKILL.md` refers to `dream`'s step numbers twice, and both references are wrong against `dream`'s actual step numbering, in ways that omit dream's index-consistency check from `/retro`'s evidence gathering and mislabel which step performs the index rewrite.

## Findings

- **[medium / confidence: high]** `skills/retro/SKILL.md:44` — cites "`dream`'s orientation and verdict table (its `--dry-run`, Steps 1–3 and 5 of that skill)" for retro's Memory evidence. `dream`'s actual steps are: 1 store, 2 orientation, 3 verdicts, 4 the index, 5 the gate, 6 the record (`skills/dream/SKILL.md:20,29,33,49,53,59`). This citation skips Step 4 ("The index"), which is exactly the step that adds the `index` verdict rows ("An index line with no file, or a file with no line, is a row of the table below with verdict `index`", `skills/dream/SKILL.md:51`) — a verdict explicitly tracked in dream's own record schema (`"index":0`, `skills/dream/SKILL.md:63`). An agent following retro's citation literally reproduces `keep`/`graduated`/`stale`/`broken`/`merge`/`promote` but never computes the `index` rows, so `/retro`'s `## Memory` section would not match dream's promised "same content" (`skills/dream/SKILL.md:66`, "From `/retro`: ... the same content"). Also, the same line says the non-`promote` verdicts "are applied at Step 6 with the proposals" — but in `dream`, application happens in Step 5 ("Apply only the numbers the user names", `skills/dream/SKILL.md:57`); Step 6 is only the record. Fix: cite Steps 1–4 for evidence-gathering and Step 5 for application, or renumber the citation to match `dream`'s actual steps.
- **[medium / confidence: high]** `skills/retro/SKILL.md:138` — "apply the `dream` skill's edits (its Steps 5 and 6: the memory files, the index rewrite)". In `dream`, "the memory files" (delete/merge/rewrite content) is Step 3 and "the index rewrite" is Step 4 (`skills/dream/SKILL.md:33,49`); Step 5 is "The gate" (present + apply) and Step 6 is "The record" — neither step performs the index rewrite. This is a second, independent mislabeling of the same seam (the `dream`/`retro` interface), inconsistent with the first (line 44 implies apply-at-6, this line implies apply/index-rewrite-at-5-and-6, and dream itself puts application at 5 and the index write at 4). Fix: cite Steps 3–5 (edits computed at 3–4, applied at 5), and drop the "Steps 5 and 6" framing.

These two are the same defect in two places: `retro`'s references to `dream`'s step numbers do not track `dream`'s actual numbering, so grepping "Step" citations against `dream/SKILL.md`'s headers (review-standards §3, "every other restatement... changed too") turns up a genuine mismatch rather than a stale copy of old numbering — both files are new/changed together in this range, so the numbers were never aligned to begin with.

No other correctness, security, or destructive-guard problems found: the `--dry-run`/headless paths in both skills correctly stop before any write; the `skill-matrix.md` row's fixture and grep checks match `dream`'s Step 3/4 vocabulary (`graduated`, `broken`, `merge`+`promote`, `keep`, `index`) exactly; the `context/AGENTS.md` sentence is consistent with `retro`'s new Step 1 item 7 and with plan decision D4 (no new trigger/timer added); frontmatter in `skills/dream/SKILL.md` is properly quoted (no unquoted `#`); paths named in the skill (`<docs-root>/bugs/INDEX.md`, `<docs-root>/architecture/DECISIONS.md`, `<docs-root>/handoffs/`, `<docs-root>/changes/`) all exist in this repo.

## Criteria

1. Dead code — pass; no code touched, no orphaned prose or files.
2. Flow tracing — fail on the two findings above (retro's citations of dream's steps don't match dream's own numbering); everything else traced cleanly (dry-run/headless guards, doubt-resolves-to-keep, the record's full-text-of-deletions requirement, the matrix row's checks against the skill's vocabulary).
3. Standards conformance — pass otherwise: frontmatter quoting correct, dates specified as `date +%F`/`date -u +%FT%TZ`, no unquoted `#`; the changed rule (Learning loop sentence) has its one other restatement (`retro`'s Step 1 item 7) updated in the same range.
4. Directory layout — pass: `skills/dream/SKILL.md` lands where every other harness skill lives; no stray files.
5. Simplicity — pass: pure prose changes, no new dependency, no hand-rolled auth/crypto surface.
6. Decoupling — the two findings are exactly a decoupling defect at the `dream`/`retro` seam: the interface (dream's step numbers) that `retro` depends on is cited incorrectly, which is the kind of thing that causes two engineers editing the two files to drift further apart rather than staying in sync.

## Not examined

`.agents/docs/plans/2026-09-29-dream-skill.md`, `.agents/docs/testing/2026-09-29-dream-phase-1-behavioural.md`, and `.agents/docs/retros/2026-09-29-dream.md` were read only as context (per the dispatch, plan and logs under `.agents/docs/` are not in scope) — not reviewed as review targets in their own right. `ROADMAP.md` and `TODO.md` were not read. No Go code is in this range, so the `golang`/`typescript` skill checklists do not apply.

## Outcome
- `skills/retro/SKILL.md:44` cites dream's Steps 1–3 and 5 and "applied at Step 6" ambiguously — fixed: the item now reads "its Steps 1 to 5: the store, orientation, one verdict per memory, the index rows, the table" and "applied at Step 6 below" (retro's own Step 6).
- `skills/retro/SKILL.md:138` "its Steps 5 and 6: the memory files, the index rewrite" — fixed: "the edits the `dream` skill's verdicts imply (its Step 3: delete, rewrite, merge) and its index rewrite (its Step 4), as its Step 5 would on a go".
