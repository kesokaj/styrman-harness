---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r3-a   (2026-09-15)

Range reviewed: `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (staged round-2 revision, on HEAD `c71aacf`), together with the surfaces it names at `c192622`. `git diff --stat c192622..HEAD` shows no change to `context/`, `skills/`, `agents/`, `AGENTS.md`, `README.md`, the testing docs, `DECISIONS.md` or the bug index. Cold-reader check, round 3, reader a.

## Verdict: REQUEST CHANGES — `design needs revision` / `implementation not ready`

The edits still apply mechanically. My own copy of the rehearsal against `c192622` parsed 88 edits with 0 failures across 26 files, for both I-19a variants (`.agents/scratch/coldread-r3-a/`).

**Worst problem — D25 was not carried into the release gate.** I-21d rewrites only the first sentence of preflight Step 4. The second sentence stays: "Its findings are the docs gaps: dispatch the `documentation` agent for them" (rehearsed `skills/preflight/SKILL.md:67`). I-21e then makes blocking "a Step 4 docs finding the `documentation` dispatch did not fix" (`:79`).
- I-1a and I-4a now put skills, agents, rules, the context and `AGENTS.md` inside the docs lens.
- `documentation`'s lane excludes all of those (`agents/documentation.md:20-32`).
- So any docs-lens finding on a skill, even a low `narrowed:`, cannot be fixed by that dispatch, and every such release is NO-SHIP by definition. The same finding in `README.md` clears.

**Blocking question — task 1.6's expect cannot be met.** It asks one 40-turn reviewer (`agents/reviewer.md:9`) to classify, with file:line proof, "every living doc". Here that is about 625 KB and 15,500 lines. The plan does not say what happens when the report comes back partial. Because sweep findings become new edits in later phases, two implementers would ship different content.

## Findings (gaps)

- **[medium / confidence: high]** I-21d/I-21e → rehearsed `skills/preflight/SKILL.md:67`, `:79`; D25; `agents/documentation.md:20-32`.
  - **What is wrong:** Step 4 still sends every docs-lens finding to `documentation`. Blocking includes "a Step 4 docs finding the `documentation` dispatch did not fix". I-1a (`.agents/{rules,skills,agents,commands}/`, `AGENTS.md`) and I-4a (`context/`, `agents/`, `skills/`) widen the lens onto files outside that agent's lane.
  - **Failure:** on this repo, any `/preflight` whose range touched a rule leaves a skill, agent or context finding unfixed, at any severity, so the verdict is NO-SHIP. I-12d and I-12e route these findings to the main agent in `/review` and the context, but the gate never does. Nothing in Step 4 or Step 6 gives the main agent a path to fix them inside the gate. (`skills/preflight/SKILL.md:40`'s no-edits sentence covers only Step 1's pre-gate commit.)
  - **Fix:** in I-21d, route as I-12d does: lane files go to `documentation`, and the rest are blocking findings for a follow-up task. Otherwise state in I-21e that findings outside the lane block only at high severity, and say which reading D30 means.

- **[medium / confidence: high]** Task 1.6 (plan Phase 1) → `agents/reviewer.md:9` (`maxTurns: 40`).
  - **What is wrong:** one reviewer dispatch must "classify each living doc as accurate, outdated or redundant with the file:line that proves it". The surface is 60 non-vendored skill markdown files plus the agents, `context/`, `README.md`, `AGENTS.md`, the testing guide and three deploy guides: about 625 KB, or 15,500 lines by `wc`.
  - **Failure:** the reviewer's early-write rule turns a cut run into a partial report with a Not examined section. The expect ("a classification of every living doc") is then unmet, and I-7a Step 4 leaves the box unticked. Failure modes covers only "many findings". Some implementers will tick against a partial sweep; others will re-dispatch or split it. Sweep findings on edited surfaces become new rehearsed edits, so the shipped content differs between them.
  - **Fix:** split 1.6 into dispatches by area (context + root docs, agents, skills in named batches), or state that one re-dispatch covers the remainder and then a named Not examined meets the expect.

- **[medium / confidence: medium]** Checkpoints of Phases 1–5 ("a `verifier` log … with Status PASS") → `AGENTS.md:46` (Testing) and BUG-004's prevention rule (`.agents/docs/bugs/INDEX.md`).
  - **What is wrong:** the repo's testing bar for any change to `context/`, `agents/` or `skills/` includes one start of each CLI in the sandbox. No task in Phases 1–4 starts a CLI, and the verifier brief is not told that this is deferred to 6.3.
  - **Precedent:** past verifier logs handled the same gap differently. `2026-09-13-verify-learning-loop-phases-1-2.md:24` gave PASS by citing an older record. `2026-09-14-verify-docs-steps.md:13` left it PENDING. So Status PASS is not a stable expect.
  - **Also:** the plan's only frontmatter change (I-20c, task 5.3) comes after 5.1, the last CLI start before Phase 6. The Phase 5 checkpoint therefore pushes it to `main` before any CLI has parsed it. A `sh-agent update` on another machine between Phases 5 and 6 would install it unchecked, which BUG-004's rule forbids.
  - **Fix:** add a D-row giving the per-phase verifier brief ("the CLI start runs at 6.3; list it under Not examined; not a FAIL"). Add a CLI parse check after 5.7, before the Phase 5 commit.

- **[medium / confidence: medium]** I-21a/b/c against I-16d.
  - **What is wrong:** a report covers a path only if its not-examined section names none of the release paths. An owed artifact that "leaves a path in the release diff unexamined" gets one re-dispatch and then blocks. I-16d's new verifier section asks for "tiers, packages, paths" that did not run.
  - **Failure:** on this repo an honest verifier names the sandbox CLI start for changed skills (the `AGENTS.md:46` bar). The path is then "unexamined", the re-dispatch cannot run that tier either, and the result is blocking. Every harness release goes NO-SHIP unless the verifier writes `nothing`, which rewards a less honest log.
  - **Undecided:** does a not-run tier on a path count as leaving that path unexamined?
  - **Fix:** define "unexamined" as not read and not run at all, or name the tiers the gate owes.

- **[low / confidence: high]** I-12e → `context/AGENTS.md:16`, compared with I-12d and `context/AGENTS.md:17`.
  - **What is wrong:** the always-loaded context sends docs findings to `documentation` "except findings in skills, agents, rules or the context". I-12d also keeps "other `AGENTS.md` sections" with the main agent, and `:17` gives `documentation` only the setup and deploy sections of `AGENTS.md`.
  - **Failure:** a finding in an `AGENTS.md` rules section, such as this repo's "Rules for working here", goes to an agent whose lane does not include it. That is the mechanism D25 closed, reopened in the context.
  - **Fix:** add "or `AGENTS.md` sections other than setup and deploy" to I-12e.

- **[low / confidence: medium]** I-10a (`skills/bug/SKILL.md:45`) and I-10b (`agents/changelog.md:94`) → `context/AGENTS.md:18-19`.
  - **What is wrong:** with no active plan, noticed or deferred work goes to "the handoff's Next Steps". The unchanged DoD reserves Next Steps for "an action only the operator can take" (`:18`), and says Done is "a handoff on disk that says nothing is queued" (`:19`).
  - **Failure:** a bug fix that notices a refactor now closes with a handoff that queues work, against two DoD sentences the plan leaves in place.
  - **Fix:** amend `:18-19` in I-8, or give that work another home (for example a Follow-ups line in the `changes/` entry).

- **[low / confidence: high]** Task 6.4, the "cause outside this plan" branch.
  - **What is wrong:** that branch needs "the harness at `c192622` installed in a second sandbox home". CLI logins live inside the home (`/home/kesokaj/sandbox/home/.claude/.credentials.json`, `.claude.json`, `.gemini/`).
  - **Undecided:** how the second home gets logged in. An interactive login is the operator's step, and copying credential files is a secrets action that no plan text authorizes. The plan also does not say whether reinstalling `c192622` into the one home is acceptable (skill-matrix Procedure step 4 forbids changes while runs are in flight).
  - **Fix:** name the procedure. For example: after all rows finish, reinstall `c192622` into the same home, re-run the failing rows, reinstall HEAD, and record both installs.

- **[low / confidence: medium]** I-1a against I-4a, with `.agents/skills -> ../skills` and `AGENTS.md` ("Third-party skills are vendored verbatim … never edit them in place").
  - **What is wrong:** I-1a searches `.agents/skills/` "at its target", which here is all of `skills/`, vendored skills included. I-4a names "every skill directory under `skills/` not listed in `skills/VENDORED.md`".
  - **Failure:** the two sentences define different sets. A stale-wording finding in a vendored skill cannot be fixed in place, and the `--fix` loop cannot close it.
  - **Fix:** I-4a says vendored skills are excluded even when reached through the symlink.

- **[low / confidence: medium]** E2E edge case / I-24b against I-22b.
  - **What is wrong:** the `bug-no-repro` prompt says "Do not ask me anything". I-22b calls that phrase "not a waiver" but never says whether it counts as "a user who said they will not answer". A `claude -p` model is not told it runs in print mode.
  - **Failure:** a run that asks for the repro in its final output and writes no handoff fails `grep -rqi "repro" .agents/docs/handoffs/`. 6.4 would then record a harness defect caused by the prompt.
  - **Fix:** add "I will not answer follow-ups" to the row prompt (as the `feature` row has), or have I-22b name "do not ask me anything" as a no-reply signal.

- **[low / confidence: high]** I-11b–g against I-7a's `[-]` grammar (rehearsed `skills/feature/SKILL.md:71`).
  - **What is wrong:** the old boxes close as "moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list". The grammar allows only `dropped:`, `superseded by <plan § task>`, or `moved to <handoff path> Next Steps (<date>)`. The plan that introduces the grammar writes closures outside it.
  - **Fix:** use `superseded by plans/2026-09-15-harness-practices-g01-g26.md § Done (last row)`.

## Questions

1. **BLOCKING (1.6):** when the sweep report comes back partial, or cut at its turn limit, is Phase 1's box ticked on it? Or is the sweep split or re-dispatched until every living doc is classified? The answer decides which docs get sweep fixes as new edits, so it decides the shipped content.
2. **NON-BLOCKING (1.2, and the in-session dispatch paragraph):** "the I-1a brief verbatim". I-1a's `new` block is a fragment that starts at "report only." and ends with a table `|`. Default: paste the whole rehearsed `--docs` brief cell (`skills/review/SKILL.md:48` after I-1a), without the table pipes.
3. **NON-BLOCKING (checkpoints 1–5):** what does the verifier brief say about the sandbox CLI start? Default: deferred to 6.3, listed under Not examined, not a FAIL (gap 3).
4. **NON-BLOCKING (6.4):** how is the second sandbox home logged in? Default: reuse the one home once no rows are in flight (gap 7).
5. **NON-BLOCKING (Done):** "Commit, push, `/handoff`" cannot truthfully be ticked before the commit it describes, as I-7b requires. Once "Close the plan" has emptied `ROADMAP.md` "Now", I-9d's drift check flags this plan if the box stays open. Default: tick it in the handoff's own commit, staged by name.
6. **NON-BLOCKING (Done `/handoff`):** the in-session handoff skill is the installed one (four questions, no I-9d drift check), and D9 names no paste for it. Default: paste the rehearsed `skills/handoff/SKILL.md` § Cold-reader check into the cold-read brief.

## Criteria

Plan review: each criterion is applied to the design the edits ship.

1. Dead code — pass. I-8e removes the hand-typed workflow-skills sentence. Its fact is in the generated inventory (`internal/harness/inventory.go:321`, which is conditional on `slashNote`; D15 reports that the render check holds). No orphaned surface.
2. Flow tracing — fail. D25's routing does not reach preflight Step 4 (gap 1). The new coverage rule does not trace to the verifier's own Not examined content or to this repo's testing bar (gaps 3, 4). I-10's handoff Next Steps conflicts with DoD step 5 (gap 6). Traced and holding:
   - I-21d → `/review` Step 3.
   - I-2a → I-22c ordering.
   - The 1.2 control's expected findings follow from I-1a's text.
   - The review table row keeps 5 pipes after I-1a.
3. Standards conformance — pass with gaps. BUG-010 discipline holds (rehearsal: 88 / 0). Frontmatter changes only in the one quoted description, but it lands in Phase 5 with no CLI parse check before the push (gap 3, BUG-004). No Go test asserts on the changed prose: `TestInventoryAndRender` checks fixture names only, and `TestNoTmpInHarnessProse` checks `/tmp` only.
4. Directory layout — pass. The control repository and the rehearsal stay under the gitignored `.agents/scratch/`, the sandbox project under `runs/`, and the ADRs in `DECISIONS.md` via `changelog`.
5. Simplicity (DIY-except-security) — pass. Text edits only, no dependency. The one security-adjacent step is gap 7's credential question.
6. Decoupling — fail on one seam. Docs-lens findings still have no owner at the release gate (gap 1), and the context's routing sentence differs from the review skill's (gap 5).

## Not examined

- No CLI or sandbox run. G26 precedence, how print mode is recognised, and matrix outcomes are reasoned from text.
- No `go test`, scratch install, `status` or `render` on the rehearsed tree. D15's 93-line render and the round-2 grep expectations for tasks 2.x–6.1 are taken from the round-2 reports and D15, not re-run by me. The only exception is the rehearsal itself (88 edits, 0 failures, both variants).
- The research record, and the round-1 reports beyond the round-2 reports' summaries.
- `skills/retro/SKILL.md`, `skills/distill/SKILL.md`, `skills/init-project/SKILL.md` and `skills/handoff/SKILL.md` beyond their anchor lines and the regions the edits change. `agents/security.md`, `agents/explorer.md` and `agents/changelog.md` beyond their anchors and tool lists.
- `skills/golang/references/golang-ecosystem.md`, `references/go-review.md` and `references/typescript-review.md`: anchors only, through the rehearsal.
- On the macOS awk: I-9d's awk was run only with GNU awk 5.2.1 on this box. Its `{3,}` interval and `nextfile` are unverified on the operator's Mac.
- I-23b (the `CLAUDE.md` fold) has no rule for when no reply can arrive. Not judged beyond noting it.
