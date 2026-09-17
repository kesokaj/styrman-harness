---
owner: reviewer
date: 2026-09-17
---
# Review — one-guide-per-project (ADR-032 / BUG-063)   (2026-09-17)

Range reviewed: working tree vs HEAD b093fd6 (`git diff HEAD`), plus the staged rename
`.agents/docs/testing/{sh-agent => }/guide.md` (`git diff --cached --stat`).
Files: `agents/documentation.md`, `agents/verifier.md`, `context/AGENTS.md`,
`skills/handoff/SKILL.md`, `skills/init-project/SKILL.md`,
`skills/init-project/assets/AGENTS.md`, `skills/review/SKILL.md`, plus the untracked
`.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`.
Blobs pinned: HEAD = `e9fa15fd91a34ab9950eccc55ad09bad054c879f` (this is also where
`.agents/docs/architecture/DECISIONS.md`'s ADR-032 entry and `BUG-062`/`BUG-063` records already
live — the uncommitted diff under review implements the decision ADR-032 already recorded);
base compared = `b093fd6`.

Correction during review: I initially misread `DECISIONS.md` (a `tail -10` on a truncated grep) as
lacking an ADR-032 entry and drafted a high-severity finding on that basis. A follow-up full grep
and direct read found `## ADR-032: One testing guide per project (2026-09-17) — Accepted` at
`.agents/docs/architecture/DECISIONS.md:259`, already committed in HEAD, with D1–D4 matching the
diff exactly (one guide at `testing/guide.md`; move-one/merge-several migration; `<app>` kept only
for `specs/<app>/`, defined in `documentation.md` alone; this repo's own guide renamed). That
finding is withdrawn; process was followed correctly (ADR before the implementing commit).

## Verdict: APPROVE WITH NITS

No correctness or process defect found once ADR-032 is accounted for. The change is coherent: the
verifier's per-app guide path and its `App:` log line are cleanly removed, `documentation.md` picks
up standalone ownership of `<app>` resolution for `specs/<app>/` without loss of fidelity, and every
cross-referencing skill (`handoff`, `init-project`, `init-project/assets/AGENTS.md`, `review`) is
updated in the same diff. The one substantive gap — the multi-file guide-merge branch of step 4 is
unexercised — is honestly disclosed in ADR-032's own Consequences and in the repro doc's "Not
examined" section, not hidden, so it is a nit rather than a blocker. A loose end worth closing before
calling the task done: no completed verifier log exists yet for this diff (see Findings).

## Findings

- **[low / confidence: high]** `.agents/docs/testing/2026-09-17-verify-one-guide.md` — the newest testing-directory entry, dated for this exact change, is `## Status: IN PROGRESS (being filled)` with `(filling in below)` placeholders under Commands and Not examined. No other, later verify log for this diff exists (checked `ls -t testing/*.md`). Definition of Done step 1 (verifier writes the log) has not actually closed for this change; either finish that run or remove the stub before the change is called done.

- **[low / confidence: low, but disclosed by the change itself]** `agents/verifier.md:39` step 4's several-files-merge branch ("when several [`testing/<name>/guide.md`] exist, write `testing/guide.md` with one section per file and leave them") has zero verification runs — `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`'s "Not examined" and ADR-032's own Consequences both say so plainly. Only the single-file `mv` branch (repro-12) and the open/fresh-write branch (repro-9/11) were exercised. Not a hidden gap — but the more error-prone half of a newly-written migration branch running zero times before shipping is worth a follow-up task, not just a disclosed footnote.

- **[low / confidence: medium]** `agents/verifier.md` steps 2 and 4 — step 2 reads as if `$docs_root/testing/guide.md` is already available to read for commands ("Read … `testing/guide.md` … for the project's own commands. Those are what you run."), while step 4 — two steps later — is where existence is first checked and, on absence, the guide is written "before running anything." On a first-ever run in a project, step 2's read is a no-op silently resolved by step 4; the ordering is correct in effect (AGENTS.md is still read in step 2, and step 4 backfills the guide before any command runs) but requires holding both steps together to see that. A first-time reader of the agent file could reasonably run step 2's read literally and treat "no guide" as `AGENTS.md`-only with no forward pointer to step 4. Not blocking; a one-clause cross-reference in step 2 would remove the ambiguity.

## Cross-checks performed

- Grepped `README.md`, `AGENTS.md`, `context/`, `agents/`, non-vendored `skills/` (per `skills/VENDORED.md`), `.agents/docs/deploy/`, `.agents/docs/specs/` for `testing/<app>`, `testing/<name>`, `testing/$app`, `testing/{app}`, and for the removed `App:` log-line template. Zero hits outside `agents/verifier.md`'s own migration-branch prose (which correctly names the old per-app shape only to describe the migration, not to prescribe it going forward) and dated logs/bug records/changes entries, which the `--docs` scope rule excludes as historical (e.g. `BUG-062-*.md`, `2026-09-17-verify-bug-062.md`, `2026-09-17-bug-062-verifier-app.md`, `INDEX.md:69`'s BUG-048 citation of the old path, `2026-09-17-bug-063-one-guide.md`) — all correctly describing what was true at the time.
- `.agents/docs/testing/sh-agent/guide.md` (old path) has no remaining references outside dated logs/reviews/plans/bugs/changes (all historical, correctly excluded).
- `agents/documentation.md` now defines `<app>` resolution for `specs/<app>/` standalone (no longer "as `agents/verifier.md` defines it" — verifier.md dropped that definition entirely, matching ADR-032 D3). Wording matches BUG-062's original three-branch order (`ls -d cmd/*/` singular → `AGENTS.md`-named → repo dir name) exactly — nothing lost in the move. No other file (`skills/golang`, `skills/docker`, `context/dev-setup.md`, `skills/init-project/assets/AGENTS.md`) claims to define the resolution algorithm; their `cmd/<app>` mentions are layout convention, not resolution logic.
- Verifier log template: `Scope:` line and `## Records` "Writes:" line both consistently say `testing/guide.md` (singular, no `<app>`), matching step 4's "Scope is never `complete`" rule and ADR-032 D1/D2. No lingering `App:` template line; its only surviving citations are dated BUG-062 records (correctly historical).
- `.agents/docs/architecture/DECISIONS.md:259` ADR-032 D1–D4 checked line-by-line against `agents/verifier.md` step 4 and `agents/documentation.md`'s new `<app>` paragraph — full agreement, no narrowing or drift.

## Criteria

1. Dead code — pass. The `App:` log-line field and the `<app>`-resolution paragraph in `verifier.md` are fully removed, not commented out; nothing orphaned.
2. Flow tracing — pass with one nit (step 2/step 4 ordering, low/medium confidence) and one disclosed-but-open gap (merge branch untested).
3. Standards conformance — pass. ADR-032 was written and accepted before this implementing diff, satisfying the harness's "propose → approval → ADR" rule; every restatement of the per-app-guide rule found by grep (context/AGENTS.md, skills/handoff, skills/init-project ×2, skills/review, agents/documentation.md) was updated in the same diff — none left stale.
4. Directory layout — pass. `testing/guide.md` at the docs root is the right place for a project-wide (not per-app) guide; the staged rename has no content diff (`git diff --cached --stat` shows `0 insertions, 0 deletions`).
5. Simplicity — pass, and an improvement: dropping the `<app>` guide-path branch removes a decision point that BUG-062/BUG-063's own evidence showed agents got wrong under load (repro-8, 1 of 5 runs). Consolidating to one guide is the simpler, more robust design.
6. Decoupling — pass. `documentation.md` and `verifier.md` each keep their own previously-owned pieces (verifier: guide existence/migration/writing; documentation: `<app>` definition for specs, shared-seam note for the guide) with no new cross-reach between them.

## Not examined

- `.agents/scratch/realhome-smoke/repro-*.jsonl` transcripts underlying the repro doc were not opened (grepped-only per instructions); I took the repro doc's own tables at face value against ADR-032's Consequences, which corroborate them. If the transcripts contradict the doc's claims, that is a separate, more serious finding I did not check for.
- `agents/verifier.md`'s other unrelated sections (What to run, suppressions, background-suite polling) were read but not re-reviewed for defects — unchanged by this diff, out of scope.
- No `go test`/`sh-agent install` re-run performed — not applicable, this diff touches only prose/markdown, no code.
- Whether `.agents/docs/testing/2026-09-17-verify-one-guide.md` is meant to be finished by a still-running verifier dispatch (in which case flagging it here is premature) versus abandoned — I could not tell from the file alone.

## Outcome

- Finding 1 (the verifier log was still IN PROGRESS): rebutted — it was being written during this review; `.agents/docs/testing/2026-09-17-verify-one-guide.md` now reads `Status: PASS`.
- Finding 2 (the several-guides merge branch unexercised): owned — stated in ADR-032's Consequences, the repro log's Not examined and the changes entry; no project with several per-app guides was available to run it against.
- Finding 3 (step 2 reads the guide that step 4 creates): fixed — step 2 reads the guide "when it exists (step 4 writes it when it does not)".
- Correction to this report's own note: ADR-032 was not committed at e9fa15f (`git show HEAD:…DECISIONS.md` has no ADR-032); it is in the uncommitted working tree, which is what the withdrawn finding should have been checked against. The withdrawal stands.
