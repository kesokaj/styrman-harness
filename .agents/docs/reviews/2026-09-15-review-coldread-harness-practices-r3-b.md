---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r3-b   (2026-09-15)

Range reviewed: the plan file `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (staged, uncommitted, on top of HEAD `55a860f`) and the surfaces it names, read at HEAD, where no named surface differs from `c192622`. Cold-reader check, round 3, reader b. Close reading started at I-12..I-18, ADR-021 and Phases 3 and 4, then covered the rest. Round 1 and 2 reports were deliberately not read.

## Verdict: REQUEST CHANGES — `design needs revision`, `implementation ready`

The mechanics hold up. I wrote my own parser and applied all 88 edits in order to copies of the files at `c192622`: 0 failures for both I-19a variants, 26 files. In that end state every Phase 3 and Phase 4 grep expectation gives the stated count. 5.3's grep prints exactly the two `skills/preflight/SKILL.md` lines, 2.5's grep prints nothing, and the I-9d awk prints only this plan (at HEAD it also prints the two old plans). No implementer has to ask the author anything, because every edit is verbatim text. The design still has holes, though.

The worst one: I-21d rewrites only the first sentence of preflight Step 4. The second sentence, "dispatch the `documentation` agent for them", survives. That routes every docs-lens finding to `documentation` at release. I-1a now produces findings in `.agents/{rules,skills,agents,commands}/`, which `documentation` does not own, and I-21e makes an unfixed Step 4 finding blocking. The release gate therefore ends NO-SHIP on findings that nothing in the gate routes to a fixer, and this contradicts D25.

## Findings

- **[medium / confidence: high]** `skills/preflight/SKILL.md:65` (I-21d, I-21e, D25): Step 4 still ends "Its findings are the docs gaps: dispatch the `documentation` agent for them and record each as a finding in the gate report". In the rehearsed end state that sentence sits right after the I-21d text (line 67).
  - I-21d switches Step 4 to the I-1a brief, which searches `.agents/{rules,skills,agents,commands}/`.
  - `agents/documentation.md:20-36` owns README, specs, the testing and deploy guides and ROADMAP Next/Later. Skills, agents, rules and the context are not in its lane.
  - I-21e makes "a Step 4 docs finding the `documentation` dispatch did not fix" blocking, and Step 6 allows no override on NO-SHIP.
  - D25 decided that docs-lens findings in skills, agents, rules or the context are the main agent's, and I-12d and I-12e say so for `/review` Step 5 and context line 16. The release gate was not changed to match.
  - Fix: extend I-21d's `old` to the whole Step 4 paragraph and route findings the way I-12d does. Findings in `documentation`'s files go to it; the main agent fixes the rest before the verdict. Reword I-21e's clause to "a Step 4 docs finding still unfixed". Add a D-row.
- **[medium / confidence: medium]** `skills/preflight/SKILL.md:50` (I-21a, D24): the new rule is "a report covers a commit only if its `Range reviewed:` line reaches it", but "reaches" is undefined for a working-tree report.
  - Definition of Done reviews run before the commit exists. `/review` Step 2 (`skills/review/SKILL.md:31`) writes `Scope: working tree — changed files: …`.
  - Most `Range reviewed:` lines in this repo's `reviews/` and `testing/` are `<sha>..worktree` or "working tree vs HEAD (<sha>)".
  - Preflight Step 1 already knows that "a line ending in `..worktree` names its base before the dots" (`skills/preflight/SKILL.md:44`). I-21a says nothing about whether such a report reaches the commit made from that tree.
  - Two runtime readers of the gate will diverge: one accepts a pre-commit review for its commit, the other treats every Definition of Done review as covering nothing and re-reviews the whole release on Opus.
  - Fix: state which one is meant, for example "a working-tree report reaches the first commit after its base whose paths it names", or "never; only range reports cover commits". Add a D-row.
- **[low / confidence: high]** `skills/retro/SKILL.md:43` (I-18i): the retro reads `Diagnosis missed:` lines "dated on or after the newest `<docs-root>/retros/` record".
  - `/distill` also writes into `retros/`, as `<date>-distill-<name>.md` (`skills/distill/SKILL.md:148`, `:174`).
  - A distill run between two retros moves the window forward. Bug records dated after the last retro but before the distill are then never read by any retro, which is exactly what G23 is meant to catch.
  - Fix: key the window on the newest `retros/*-retro*.md` record.
- **[low / confidence: high]** `.agents/docs/bugs/INDEX.md:3` and `skills/init-project/assets/agents-tree/docs/bugs/INDEX.md:3` (ADR-021, I-18d, I-18e): both header comments say "One row per BUG-NNN record in this directory; newest last".
  - After a graduation the main table no longer has one row per record, and a new table appears below it.
  - Neither file is in Surfaces touched, and `INDEX.md` is listed as "status only". This is the global describing-text hard rule (`context/AGENTS.md:55`).
  - Fix: add one edit per file, or a D-row saying the comment stays true because it counts both tables.
- **[low / confidence: medium]** `context/AGENTS.md:31` (I-18l) against `agents/bugs.md` (I-18e) and ADR-021: the graduation targets do not agree.
  - ADR-021 and I-18e count `context/AGENTS.md` as a graduation target.
  - I-18l's trigger names only "the project's `AGENTS.md`, a `.agents/rules/` file, or a test". I-18k adds a trigger only to the project-scope apply path (`skills/retro/SKILL.md:131`), not the harness-scope path at `:132`.
  - In this repo a rule graduating into `context/AGENTS.md` has no instruction telling anyone to dispatch `bugs`. Rules that graduate into `skills/` or `agents/`, which is where most of this repo's prevention rules land (BUG-014's lives in the agents' Records sections, BUG-018's in `skills/bug`), never qualify at all.
  - Either say so in ADR-021's Consequences, or make the trigger list match the target list.
- **[low / confidence: medium]** Failure modes, "A CLI log shows a parse error for a skill or agent: the phase stays open until fixed", against `AGENTS.md:46` and the Phase 1–5 checkpoints.
  - The repo's Testing rule verifies a change to `context/`, `agents/` or `skills/` by install, status and one start of each CLI.
  - No checkpoint in Phases 1–5 starts a CLI; the first CLI start is 6.3. So that failure mode cannot fire in the phases it names, and five phases of skill and agent edits are pushed without the repo's stated check.
  - Earlier prose-only verifier logs did pass without CLI starts (`testing/2026-09-14-verify-bug-012-fix.md`), so the checkpoints are reachable. The deferral is simply undecided.
  - Fix: add a D-row deferring the CLI start to 6.3, since only the `reviewer` description in Phase 5 touches frontmatter. Alternatively, add a CLI start to the Phase 5 checkpoint and reword the failure mode.

## Questions (none blocking)

1. Preflight Step 4 routing (finding 1): should the second sentence be rewritten, and how? Non-blocking, because a verbatim implementer leaves it as written; the design still needs the answer.
2. I-21a: does a working-tree `Range reviewed:` line reach the commit made from that tree (finding 2)? Non-blocking for the implementer, a design gap for the gate.
3. I-18i: does the "newest retros record" window count distill records? Non-blocking.
4. ADR-021: in this repo, does a rule graduating into `skills/`, `agents/` or `context/AGENTS.md` (through `/retro --harness`) move the row, and who dispatches `bugs`? Non-blocking.
5. Task 3.1: its `expect:` lists `1`, `1`, `1` for seven greps, four of which carry inline counts. A reasonable implementer reads every one as `1`, and all seven print `1` in the rehearsed end state. Non-blocking.
6. Task 3.1: the pattern ``"`in flight: <who>`"`` is written inside double quotes. Run as shown, the shell would treat the backticks as command substitution, so it has to be single-quoted. Non-blocking.

## Criteria

1. Dead code — pass. I-8e removes the line-71 sentence whose fact the generated inventory carries (D4). No edit leaves an orphaned rule; 2.5's grep of the rehearsed end state finds no stale wording.
2. Flow tracing — fail. The routing from the docs lens to a fixer is traced through `/review` Step 5 and context line 16 but not through preflight Step 4 (finding 1). The graduation trigger is not traced to the harness-scope retro path (finding 5). Everything else holds:
   - I-16a–d, I-17a–f and I-18a–l each have their `old` exactly once, as the rehearsal shows.
   - `derive.go:54` rewrites only `~/.claude/skills/`, so 4.5's derived agents keep the new body text.
   - ADR-021 is the next free number (the last is ADR-020, `DECISIONS.md:157`). Its "refused as speculative" claim is in `reviews/2026-09-13-improvement-review.md:61`.
3. Standards conformance — fail (low). Describing text for the index header is not updated (finding 4). The CLI-start deferral against `AGENTS.md:46` has no D-row (finding 6). No `REQUIREMENTS.md` violation: no Go code, no hooks, no settings edits, and every test runs against a scratch or sandbox home (R6).
4. Directory layout — pass. Scratch goes to `.agents/scratch/`, sandbox runs to `sandbox/runs/` as `skill-matrix.md:12` does, and new logs to `testing/`.
5. Simplicity (DIY-except-security) — pass. Text edits only, and the counts are taken with `grep`. `nolintlint` is already described at `golang-ecosystem.md:369`. No security-sensitive surface.
6. Decoupling — pass. Each phase edits its own files, and the shared files (`context/AGENTS.md`, `skills/review/SKILL.md`, `skills/bug/SKILL.md`) are edited in a fixed, rehearsed order.

## Not examined

- The rendered line count of the rehearsed end state. I did not re-run `sh-agent render` on it. I did confirm by reading that every context edit (I-8a–e, I-12e, I-18l, I-19a, I-20a) is single-line, so none adds a rendered line.
- The I-9d awk on macOS BSD awk (`{3,}` intervals, `nextfile`). It was run only with GNU awk 5.2.1 on this box.
- How either CLI actually behaves: G26 precedence, gate stops in print mode, the `bug-no-repro` and `verifier-counts` matrix rows. These are Phase 5 and 6 runtime outcomes a read cannot settle.
- A line-by-line semantic read of I-7a (feature Steps 1–4) and I-22a–e beyond their interaction with I-12. I-23a–b (init-project) and I-24a–d were read for consistency only.
- Round 1 and round 2 cold-read reports, unread on purpose to keep this read independent.
- The research record, beyond the G10, G11, G15–G18, G23 and G24 rows and its notes on them.
