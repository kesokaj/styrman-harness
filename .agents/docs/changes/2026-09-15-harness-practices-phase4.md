---
owner: main
date: 2026-09-15
---
# Harness practices Phase 4: agents and records (G07, G16 verifier half, G17, G18, G23, G24)     (2026-09-15)

## Summary

Phase 4 of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`: edits I-16a–d, I-17a–f and I-18a–n, applied verbatim from the plan's old/new blocks, and ADR-021. The verifier runs every command on its own, records test counts against the previous log, marks inherited failures and lists added suppressions. The agents that can be cut off count a cut tool result as not reached, and the reviewer, security and verifier templates carry a Not examined section. A bug record can carry a `Diagnosis missed:` line that the retro reads. A bug's index row moves to a Graduated table when its prevention rule graduates.

## Type

task

## Files changed

- `agents/verifier.md` — I-16a–d and I-17b. What to run gains four rules. Every command on its own: a chain is split into its commands, each run in the directory the chain would run it in, and `gofmt -l` fails by printing. Test counts: each Go test command runs once with `-v`, its output saved under `.agents/scratch/verify-<slug>/`, and the log writes a `counts:` line and a `vs previous:` line against the newest earlier log with counts; a fall in pass or a rise in skip is a FAIL unless the diff or the environment explains it. Inherited failures: a failing test the change set did not touch is marked `pre-existing?` and Status stays FAIL. Suppressions added: every added suppression, skip or `.only` is listed as file:line, and an added `.only(` or a suppression with no reason is a FAIL. The log template's Go line records exit status, counts and vs previous, and the template gains `## Not examined`. A cut tool result counts as not reached.
- `agents/reviewer.md` — I-17a and I-17e: a cut tool result (a preview, a saved-output path, a line limit) counts as not reached until the rest is read; the report template gains `## Not examined` for files, criteria and commits, or `nothing`.
- `agents/security.md` — I-17c and I-17f: the same sentence; the findings template gains `## Not examined` for paths, surfaces and tools, or `nothing`.
- `agents/explorer.md` — I-17d: the same sentence.
- `agents/bugs.md` — I-18a–g: it keeps both tables of `INDEX.md`, and a new id is the highest in either table plus one; a record carries a `Diagnosis missed:` line when the dispatch names one; when a dispatch says a prevention rule graduated in this repository, the bug's row moves to a `## Graduated` table (`| ID | Area | Graduated to |`) and the record file stays untouched, while a rule promoted from a project to the harness leaves the project's row in place; the prevention recommendation says whether the rule is actionable and generalizable, as information and never a filter; moving a row edits the index, not a record.
- `skills/bug/SKILL.md` — I-18h: when Step 2 ran and the reviewer, the verifier or the fix surfaced a cause the diagnosis did not rank, the `bugs` dispatch carries a `Diagnosis missed:` line naming the line of the Step 2 brief that would have found it.
- `skills/retro/SKILL.md` — I-18i–k: the retro reads `Diagnosis missed:` lines in bug records dated on or after the newest retro record; a bug-derived proposal that targets an `AGENTS.md` says whether its rule is actionable and generalizable; when an applied rule is a bug's prevention rule, `bugs` moves that bug's row to the Graduated table in the same task.
- `context/AGENTS.md` — I-18l: a recurring rule can also graduate into a skill or an agent definition, and the commit that writes it has `bugs` move the matching index rows to the Graduated table. The rendered context stays at 93 lines.
- `.agents/docs/bugs/INDEX.md` and `skills/init-project/assets/agents-tree/docs/bugs/INDEX.md` — I-18m and I-18n: the header says a row lives in the main table, or in the Graduated table once its prevention rule lives in an `AGENTS.md`, a rule, a skill, an agent definition or a test.
- `.agents/docs/architecture/DECISIONS.md` — ADR-021, "Bug index rows move to a Graduated table when their rule graduates" (task 4.4).
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — boxes 4.1–4.5 and the Phase 4 checkpoint ticked with evidence; task 4.3's label corrected from I-18a–l to I-18a–n, since I-18m and I-18n are the index edits its check already covered.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-phase4.md` — the Phase 4 verifier log.

## Agents dispatched

changelog (ADR-021); verifier, dispatched with the repo's new What to run section pasted, because the session loads the installed definition (plan D9). This entry was written by the main agent from the phase diff.

Review: no `reviewer` ran over the Phase 4 diff. The plan's Done list runs `reviewer` and a `--docs` review once over `c192622..HEAD`.

## Verification

The verifier log reads Status: PASS, with `counts: pass 31, fail 0, skip 1 at f3013a6+dirty for `go test -race -v ./...` and `vs previous: none (no earlier `testing/*.md` log carries a `counts:` line; grep confirmed).`. The main agent also ran, each bare with its exit status read: `go test -race ./...` exit 0; `gofmt -l .` printed nothing; `go vet ./...` exit 0; a scratch install into `.agents/scratch/phase4-checkpoint-home` exit 0, then `status` with drift none; `render` 93 lines. In that home the four derived Antigravity agents carry the not-reached sentence, the derived verifier carries "Every command on its own", and each changed agent has one `## Records` section in both vendors' copies. The checks of boxes 4.1–4.5 met their expects, and every pattern printed 0 before the edits, in the tree and over `c192622`.

Not run: the sandbox start of each CLI (testing guide §4), deferred to the plan's Phase 5 checkpoint (D36).

## Docs

The changed living docs are the agent, skill, context and index texts above. `README.md`, the testing guide and the deploy guides are unchanged in this phase.

## Requirements

n/a
