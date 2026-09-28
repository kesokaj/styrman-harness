---
owner: reviewer
date: 2026-09-28
---
# Review — universal-cleanup   (2026-09-28)

Range reviewed: 381e86d..worktree (agents skills)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The diff is a citation-stripping pass (drops trailing `(BUG-NNN[, ADR-NNN])` / dated-event parentheticals from 11 files) and every rule's substantive text, conditions and thresholds are preserved verbatim or generalised without loss. Two leftover repo-specific citations were missed in files this same diff otherwise cleaned, which is the only real defect: the cleanup is inconsistent within its own scope, not wrong in what it removed.

## Findings

- **[medium / confidence: high]** `agents/verifier.md:37` — same file's line ~35 had `(ADR-032, BUG-062)` stripped from "one guide per project, never one per app", but line 37 still reads "A project from before ADR-032 may hold `testing/<name>/guide.md` files...". `ADR-032` is styrman-harness's own architecture-decision numbering (confirmed via `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`), not a concept any consuming project's verifier agent can resolve — this file is deployed as-is to every project. Same class of leftover the rest of the diff removed; strike the citation (the "project from before X may hold..." migration guidance reads fine without naming the ADR).
- **[medium / confidence: high]** `skills/review/SKILL.md:62` — the same file had `(BUG-033)` stripped at line 40, but line 62 still cites `(G10 of the 2026-09-16 simulation)`, a dated event from this repo's own trail (`.agents/docs/testing/2026-09-16-linkshelf-simulation.md`, a fictional "linkshelf" consumer-project rehearsal). Meaningless to any other project the skill deploys to; same fix as above.

No other changed line drops a condition, qualifier or reason a reader needs — checked each hunk individually (see below).

## Judgment on the stated deviation (no behavioural run for the five `agents/*.md` edits)

Acceptable. Read each of the five diffs (`changelog.md`, `explorer.md`, `reviewer.md`, `security.md`, `verifier.md`) hunk by hunk: every edit removes only a trailing `(BUG-NNN[, ADR-NNN])` parenthetical from the end of an existing sentence; no instruction, condition, threshold or example step was reworded, added or removed. `verifier.md` has the largest count (five citations dropped: `ADR-032, BUG-062`; `BUG-063`; `BUG-048`; `BUG-014`; `BUG-038`) but each is the same shape — trailing citation only, rule text byte-identical apart from the parenthetical. Nothing in these five files changes what any agent does, so `AGENTS.md § Testing`'s behavioural-run requirement (aimed at catching a definition change that alters tool-call order or output) has nothing to exercise here; skipping it does not hide a risk. (This finding is orthogonal to the two leftover citations above, which are in `skills/`, not the five `agents/*.md` files under this deviation.)

Spot-checked the two skills/*.md edits that do more than delete a citation, to confirm neither is a behaviour change in disguise:
- `skills/feature/SKILL.md:88` — `(never … — BUG-033/041 in styrman-crew both shipped a red test that way)` → `(never …: that chain guards `tail`, and has shipped a red test)`. Adds a mechanism explanation, drops the incident citation; the banned pattern and its consequence are unchanged.
- `skills/golang/SKILL.md:20` — "once deleted the only copy of it, BUG-069" → "can delete the only copy of it". Tense generalised from a specific past incident to a standing capability; the guard rule itself (check errors → keep, never fall through to delete) is untouched.
- `skills/preflight/SKILL.md` items 1 and 5 — drop a dated anecdote (`2026-09-18: a run paused here…`) and an incident narrative (GitHub's `HTTP 500` batch-upload failure), in both cases leaving the governing rule itself (`the ship request was the approval…`; `one create-with-every-asset call rolls the whole release back…`) stated as fact, not weakened to a hedge.

## Repository-specificity check on deployed files (especially `skills/preflight/SKILL.md`)

Read `skills/preflight/SKILL.md` in full (96 lines, post-diff). No repo-specific content remains — deploy paths, project IDs, tag forms and doc locations are all referenced generically (`<docs-root>`, "the project's `AGENTS.md`", "the deploy guide it points to"), and the diff's own citation removals (`BUG-033`, `BUG-034`, `BUG-083`) leave no residue in this file. The two leftovers found are in `agents/verifier.md` and `skills/review/SKILL.md`, not in `preflight/SKILL.md`.

## Criteria

1. Dead code — pass; no code in this diff (all `.md`), no orphaned references created by the removals (each stripped parenthetical was purely explanatory, nothing else in the same files points at the removed citations).
2. Flow tracing — n/a; no functions/modules changed, prose-only diff.
3. Standards conformance — fail (low-severity within itself): the global context's describing-text rule ("a rule ... has every other restatement of it ... changed too") applies to this diff's own stated goal (strip repo-specific citations) and two instances were missed in files the diff otherwise touched; see Findings.
4. Directory layout — pass; no new files, no relocations.
5. Simplicity — pass; the diff is a strict subtraction, no new abstractions or dependencies.
6. Decoupling — pass; each file edited independently, no cross-file coupling introduced.

## Not examined

- Files outside `agents/` and `skills/` (e.g. `context/AGENTS.md`, `README.md`) were only spot-checked via grep to confirm the *already-generic* rules there (destructive-guard, verification-piping) match the newly-generalised skill wording — not fully read line-by-line, since they are out of the stated diff scope.
- Did not run `gofmt`/`go vet`/`go test` — no `.go` files in this diff, not applicable.
- Did not check the untracked files listed in git status (`.agents/docs/plans/2026-09-28-recon-from-zip.md`, three `reviews/*381e86d*` files) — out of scope (dispatch names `agents skills` only).

## Outcome

- **Finding 1** (`agents/verifier.md:37`, ADR-032): fixed. The sentence now reads "An older project may hold per-app `testing/<name>/guide.md` files".
- **Finding 2** (`skills/review/SKILL.md:62`, "G10 of the 2026-09-16 simulation"): fixed, by removing the reference.
- **Found by the widened sweep** (ADR numbers, dated events): `skills/init-project/SKILL.md:50` "(REQUIREMENTS R2, ADR-036)" was removed. `init-project`'s `DECISIONS.md` template `ADR-001` is a placeholder and stays.
- **Deviation** (no behavioural run for the `agents/*.md` edits): accepted by this review.
