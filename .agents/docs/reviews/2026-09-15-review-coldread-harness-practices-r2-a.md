---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r2-a   (2026-09-15)

Range reviewed: `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (blob c830573, at HEAD c71aacf), plus the surfaces it names at c192622 (unchanged at HEAD except ROADMAP.md)

## Verdict: REQUEST CHANGES

**design needs revision / implementation not ready.**

The worst problem is that Phase 6 has no way to close. Task 6.4 is one box covering 14 matrix rows on two CLIs. Its failure branch offers "a `[-]` box with the reason and an `Open:` row". The Checkpoint expects "6.2 to 6.4 pass", and the plan's own Step 4 rule (I-7a) says a missed expectation leaves the Checkpoint unticked and the next phase closed. A single row that fails for a vendor reason therefore blocks Done for good. The plan also gives no re-run policy for rows run by an LLM, which vary from run to run.

Five other gaps are decisions the plan does not make:
- The new handoff drift grep (I-9d) will flag this plan itself forever.
- The release gate's coverage rule (I-21a/b) passes, without checking anything, for verifier logs written from the current template.
- The widened docs lens (I-1a) produces findings on skills, agents and rules that only the `documentation` agent is routed to fix, and that agent's lane excludes those files.
- The G26 outcome-B criterion cannot tell "project wins" from "both load".
- The e2e row list leaves out the one row that exercises the Task path that I-8d and I-22a change.

The round-1 items are resolved. My in-memory rehearsal confirms every anchor:
- 82 edits parse against c192622, with 0 failures for both I-19a variants, across 26 files.
- At the end state, 5.3 finds exactly `skills/preflight/SKILL.md:40` and `:50`.
- The 5.5 expectations hold.
- `nolintlint` lands on `golang-ecosystem.md:323`.
- `sh-agent render` prints 93 lines today.

## Findings

- **[medium / confidence: high]** Plan Phase 6, task 6.4 and its Checkpoint (plan:987–988) vs I-7a Step 4 ("tick it only when the evidence meets its expect — a missed expectation leaves it unticked and the next phase closed"). **What is wrong:** 6.4 bundles `feature`, `bug`, `bug-no-repro`, `review`, `review-project`, `preflight`, `retro`, `distill`, `handoff`, `init-deploy`, `explorer`, `review-standards`, `golang` and `verifier-counts` on both CLIs into one box. Its failure branch ("a `[-]` box with the reason and an `Open:` row") and the failure mode "A matrix row fails because the demo or a CLI changed … it is not recorded as a pass" both leave the Checkpoint's expect ("6.2 to 6.4 pass") unmet. The Done phase then never opens. **Also missing:** how many times a failed row is re-run, and whether a green re-run counts, even though these are nondeterministic headless LLM runs. **Fix:** give each row its own box, or say that a `[-]` row with an `Open:` row meets the Checkpoint. Also state the re-run rule, for example one re-run with both logs kept.

- **[medium / confidence: high]** I-9d (`skills/handoff/SKILL.md:154`): `grep -l '^- \[ \]' <docs-root>/plans/*.md` matches checklist lines inside fenced blocks. **Evidence:** this plan has such lines at plan:172–174 (the I-7a template), plan:318, 326, 334, 342, 350, 358 and 366 (the I-11 `old` blocks), and plan:412–456 (the I-14/I-15 checklist text). The grep run during this review listed them. **Why it fails:** after the Done row closes every real box, every later handoff cold read in this repo reports this plan as "a plan with `[ ]` boxes that `ROADMAP.md` 'Now' does not name". That is a permanent false contradiction, and so is any future plan that quotes a template or checklist. Task 2.4 (plan:956) hides the problem because this plan is still expected to match at that point. **Fix:** make the check skip fenced blocks, for example `awk '/^```/{f=!f} !f && /^- \[ \]/{print FILENAME; nextfile}' <docs-root>/plans/*.md`, or limit it to the `## Phases` section. Then re-run the rehearsal.

- **[medium / confidence: high]** I-21a/I-21b (`skills/preflight/SKILL.md:50`, `:52`) and D14. **What is wrong:** a report counts as covering a commit when its not-examined section, named `Not examined`, `Not checked` or `What was not checked`, lists none of the commit's paths. Neither template has such a heading:
  - The verifier's log template (`agents/verifier.md:57–83`) has only a `Scope:` line.
  - The security template (`agents/security.md:52–69`) has `## Checked`.
  - Six existing logs in this repo use `## Not verified here`, for example `.agents/docs/testing/2026-09-13-verify-real-home.md`.
  - Eleven dated testing logs have none of the three headings, among them `2026-09-14-verify-docs-matrix-row.md` and `2026-09-13-verify-research-pass.md`.

  **Why it fails:** a log with no recognised section "names no path in the release diff" automatically, so the Verification row still fails open, which is the thing G20 is meant to stop. D14's claim that these three names are the ones "the agents and records use" is untrue for verifier logs. **Fix:** say that a report without one of the recognised sections covers nothing, or add a `## Not examined` heading to the verifier and security templates as edits in this plan.

- **[medium / confidence: high]** I-1a (`skills/review/SKILL.md:48`) together with I-4a (repo `AGENTS.md:46`), against the routing that stays unchanged:
  - `skills/review/SKILL.md:60`: "`--docs` → dispatch `documentation` with the list"
  - `context/AGENTS.md:16`: "its findings go to `documentation` in step 3"
  - I-12b: "a `--docs` list with `documentation`" is the in-flight ledger line

  **What is wrong:** the lens now searches `AGENTS.md`, `.agents/{rules,skills,agents,commands}/` and, here, `context/`, `agents/` and `skills/`, and it emits `dropped:` and `narrowed:` findings on them. `agents/documentation.md:20–32` limits that agent's lane to specs, the testing guide, deploy guides, README and ROADMAP. **Why it fails:** a `narrowed:` finding on a skill, which is BUG-018's mechanism and the reason G06 exists, goes to an agent that must decline it or edit outside its lane, and the `--fix` loop cannot finish. D7 covers only this plan's own I-4a/I-5 texts. **Fix:** in I-12a, or a new edit at `review:60`, route `--docs` findings on rules, skills, agents, context and `AGENTS.md` rule sections to the main agent, and send only lane files to `documentation`.

- **[medium / confidence: medium]** I-19 outcome rule (plan:694): "Outcome B holds only when both CLIs quote both control markers and both project markers." **What is wrong:** nothing requires the harness copy to be absent. If a CLI lists both same-named copies and the model quotes both descriptions, B is recorded. `new-B` then asserts "the project layer has priority on both vendors", which the run did not show. The research record's G26 row (`.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md:75`) already notes that Claude Code's docs report skill precedence differently. **Fix:** B also requires that neither CLI quotes the harness's own `review` skill or `reviewer` agent description, and the log records each quoted description verbatim.

- **[low / confidence: medium]** 6.4 row list vs I-8d and I-22a. **What is wrong:** the `documentation` matrix row (`skill-matrix.md:52`, "Do it as a task through the Definition of Done; do not ask me anything") is the only row that drives the Task class. I-8d makes a Task that meets a design question "a Feature … under the feature skill's gates", and I-22a makes a reached gate end a print-mode run. **Why it matters:** a regression where that run now stops at a scope gate and writes no docs would not be seen in Phase 6. **Fix:** add `documentation` to 6.4.

## Readiness questions

1. **BLOCKING (6.4 / Phase 6 Checkpoint):** when a matrix row fails for a vendor or demo reason, or does not reproduce, how many re-runs does it get, and can Phase 6 close with `[-]` rows? Two implementers would either leave the plan open or tick against an unmet expectation.
2. NON-BLOCKING (1.2): "dispatch `reviewer` there". An in-session subagent starts in the harness checkout. It reads `AGENTS.md` at the project root and resolves the docs root with `git rev-parse --show-toplevel`. The brief should give the control repository's absolute path and require `git -C` and a report path inside it; otherwise the harness `AGENTS.md`, with I-4a's paths, leaks into the control and the report lands in this repo's `reviews/`. Default: state the path in the brief.
3. NON-BLOCKING (1.6 / Failure modes): a sweep finding on a surface a later I-edit rewrites, such as feature Steps 1–4 replaced whole by I-7a, is lost if fixed before that edit and conflicts with the verbatim `new` if fixed after. Default: record it as a D-row and a new rehearsed edit applied after the I-edit.
4. NON-BLOCKING (I-24c): the content of `internal/greet/extra_test.go` and how the driver creates the uncommitted file are not given. A reasonable implementer writes one trivial passing test. The positive control is sound: the demo at 9ae4eff has no `counts: pass` or `vs previous:` line today.
5. NON-BLOCKING (5.1): "this tree installed into the sandbox home" happens at the start of Phase 5, before 6.2 reinstalls. Default: install at that moment and cite the install output in the G26 log.
6. NON-BLOCKING (5.2): "the render file" means re-rendering after I-19a, since `.agents/scratch/render.txt` is from 2.2.
7. NON-BLOCKING (I-18d): the column header of the `## Graduated` table is not given. The row shape `| BUG-NNN | area | → <home file> |` lets `bugs` derive it.

## Criteria

Plan review, so the criteria are applied to the plan's design rather than to a diff.

1. Dead code: pass. I-8e removes the context:71 sentence whose fact the generated inventory carries (`internal/harness/inventory.go:322`; render currently contains "can also be typed as" once). No orphaned surface found.
2. Flow tracing: fail. The routing of `--docs` findings (I-1a vs `review:60` and `context:16`), the release-gate coverage rule vs the actual log templates (I-21a/b), and the Phase 6 closure path (6.4 vs its Checkpoint) are broken, as described above. The derived Antigravity agent carries description and body verbatim (`internal/harness/derive.go:33–55`), so 4.5 and 5.7 are sound.
3. Standards conformance: pass with one exception. Frontmatter changes only in the `reviewer` description, which stays one quoted line. No Go code. No `/tmp` added per D15. The exception is that the I-9d grep will contradict the global resume rules for closed plans (finding 2).
4. Directory layout: pass. Scratch fixtures sit under `.agents/scratch/` (gitignored, `.gitignore:21`) and the sandbox. Records land in the paths the agents define.
5. Simplicity: pass. Edits are text-only and stdlib-only, with no new dependency.
6. Decoupling: fail on one seam. I-1a gives the reviewer's docs lens files that no fixer owns (finding 4).

## Not examined

- No CLI or sandbox run. Vendor load and precedence behaviour (G26), print-mode gate behaviour, and the matrix checks' live outcomes are judged from text only.
- The research record: only its group rows G06, G20–G22, G24 and G26 were read.
- Round-1 cold-read reports: finding headlines only, to confirm they are resolved. Not re-reviewed.
- `skills/init-project/assets/` templates. Not checked for describing text of the changed plan template or thresholds beyond the `git grep` over `skills`, which found none.
- `skills/retro/SKILL.md`, `skills/distill/SKILL.md` and `skills/handoff/SKILL.md`: only the edited regions and their anchors (rehearsal) were read, not every line.
- `agents/explorer.md`, `agents/security.md`, `agents/changelog.md`: only the anchor lines and section headings were read.
- The `research`, `ctx-deploy` and other matrix rows not named by 6.4: not traced against the new gate text.
