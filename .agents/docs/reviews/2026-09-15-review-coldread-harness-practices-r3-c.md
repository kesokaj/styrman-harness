---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r3-c   (2026-09-15)

Range reviewed: plan file `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` at `55a860f` (cold-reader check, round 3; the working tree is clean), plus the surfaces it names at that commit. Close reading started at I-19 to I-24, ADR-022, Failure modes, Verification, Phases 5 and 6, Done, Risks and Decisions, then covered the rest.

## Verdict: REQUEST CHANGES — `design needs revision` / `implementation not ready`

The worst problem is in the release gate. Preflight Step 4 still sends every `--docs` finding to `documentation` (`skills/preflight/SKILL.md:65`, second sentence, which I-21d leaves in place). After I-1a the brief also hunts in skills, agents, rules and the context, and `documentation` does not own those files. I-21e then makes "a Step 4 docs finding the `documentation` dispatch did not fix" a blocking finding. A release with one stale skill sentence therefore reaches NO-SHIP, and no preflight step can clear it. This contradicts the plan's own routing in I-12d/I-12e/D25. The global describing-text rule forces the implementer to rewrite that sentence without text from the plan, so this question is blocking.

Two more gaps change what ships:
- I-19a `new-A` still says "Both load". The plan expects an outcome-A result in which the project layer did not load at all.
- The I-21a coverage test fails open for a not-examined section that names commits or criteria instead of paths, and I-17e's own template invites exactly that.

Round 2's findings are otherwise absorbed: fence-aware awk, ledger states, the G26 harness-description condition, chain directory prefixes, the scope waiver in the `feature` row.

## Findings

- **[medium / confidence: high]** I-21d, I-21e vs I-12d, I-12e, D25 — `skills/preflight/SKILL.md:65`. **What is wrong:** I-21d replaces only Step 4's first sentence. "Its findings are the docs gaps: dispatch the `documentation` agent for them" stays. The I-1a brief that Step 4 now uses "exactly as `/review` Step 3 gives it" searches `.agents/{rules,skills,agents,commands}/`. `documentation`'s lane (`agents/documentation.md:22-32`) is README, specs, testing guides, deploy guides and ROADMAP. **How it fails:** under I-21e, any finding in a skill or agent file is "a Step 4 docs finding the `documentation` dispatch did not fix", so the result is NO-SHIP. Step 6 allows no override. A second problem stays undecided: a docs fix made during Step 4 edits files outside the gated tree, but Step 6.3 and :90 say "nothing else rides along". **Fix:** give I-21d a replacement for the second sentence that routes like I-12d. Define the blocking condition as "a Step 4 finding its owner did not fix before the release commit", and say that those fixes ride in the release commit. Add `skills/preflight/SKILL.md` to 2.5's describing-text grep.
- **[medium / confidence: high]** I-19a `new-A` vs Failure modes ("A G26 control is not quoted on a CLI: outcome A (D12); the log says which layer did not load") and D12 ("cause-neutral outcome A"). **What is wrong:** `new-A` begins "Both load, but do not rely on which copy wins". The rest of `context/AGENTS.md:37` is untouched and still says "`/skills` and `/agents` list both layers". **How it fails:** if 5.1 finds that Antigravity does not load `.agents/agents/<name>/agent.md` at all (the Risks section lists this assumption, "checked by 5.1's control markers"), Phase 5 ships a context sentence that the test just disproved. That breaks the hard rule on describing text. **Fix:** add a third alternative for an unquoted control that drops "Both load" and the "list both layers" clause and names the layer that did not load. Or add a D-row saying that such a result stops at 5.2 for the user.
- **[medium / confidence: high]** I-21a, I-21b vs I-17e and I-16d. **What is wrong:** a report covers a commit when its not-examined section "nam[es] none of that commit's paths". I-17e's reviewer template tells the agent to list "files, criteria, commits", and I-16d's verifier template says "tiers, packages, paths". **How it fails:** a report whose Not examined section says "commits 4–9 of the range", "criterion 2", "the e2e tier" or "the rest of the diff after the first 40 files" names no path, so it counts as full coverage. That is the fail-open case G20 and D2 ("a release path no report examined") exist to close. **Fix:** a report covers only when its not-examined section is `nothing` or names paths only. Any entry that is not a path (a commit, a criterion, a count) makes the report cover none of the commits it could refer to.
- **[medium / confidence: medium]** I-21c "gets one re-dispatch over those paths". **What is wrong:** the scope of the re-dispatch is not defined. A path scope in `/review` is "the named files … read whole — not a diff" and goes on the `Range reviewed:` line (`skills/review/SKILL.md:34,37`). Under I-21a that line "reaches" no commit, so the re-dispatched report covers nothing and the path is still unexamined, which makes every re-dispatch end blocking. For the verifier, "over those paths" has no meaning, because it runs the project's commands (`agents/verifier.md:40-47`). **Fix:** name the re-dispatch scope, for example `Scope: commits <base>..HEAD limited to <paths>; the diff is git diff <base>...HEAD -- <paths>`. Say that I-21a reads this report's paths together with the first report. For the verifier, re-dispatch over the full range.
- **[medium / confidence: medium]** I-24b `bug-no-repro` prompt and the E2E edge criterion, vs I-22b. **What is wrong:** the row passes only if a handoff is written, which I-22b requires only "when no reply can arrive — a print-mode run … or a user who said they will not answer". The prompt says only "Do not ask me anything", and I-22b itself calls that phrase "not a waiver". Nothing in the prompt tells the model it runs in print mode. Compare I-24d, which gives the `feature` row "I will not answer follow-ups". **How it fails:** a model that correctly holds the gate asks the repro question in chat and stops without writing a handoff, so `grep -rqi "repro" .agents/docs/handoffs/` fails. 6.4 bars this row from the "outside this plan" route, so the failure is recorded as a harness defect in text that behaved correctly. **Fix:** add "I will not answer follow-ups" to the row prompt. It names no gate, so it is not a waiver.
- **[medium / confidence: high]** Phase 6 Checkpoint: "Not examined section names what no row exercises end to end: I-21…, I-12…, I-22e, I-18…, I-23b, I-19b–d". **What is wrong:** the list is presented as complete, but it is not. Gaps:
  - I-20 and ADR-022: the `bug` row's check (`skill-matrix.md:23`) tests only the fix and a BUG file, and Procedure step 3's transcript audit sets no expectation for a `reviewer` dispatch.
  - I-9a–d: the `handoff` row checks only that a file exists, `skill-matrix.md:31`.
  - I-16b: inherited failures and added suppressions; `verifier-counts` adds neither.
  - I-13a: rebuttal reasons.
  - I-17a–d: the cut-result rule.
  - I-14 and I-15.
  - I-23a: the `init-deploy` prompt supplies the commands (`skill-matrix.md:22`), so derive, print and confirm never runs.
  - I-8b and I-8d.

  The global hard rule says claims of test coverage name what is not covered. **Fix:** complete the list, or phrase it as "includes at least" and have the log enumerate every I-edit against the row that exercises it.
- **[low / confidence: medium]** 6.4 and Failure modes: "a harness defect (a bug record, a fix, and the row re-run until it passes)". **What is wrong:** headless LLM rows are nondeterministic, and nothing caps the fix-and-re-run loop. For `bug-no-repro` and `verifier-counts` the only non-pass exits are a harness defect or a gate the row does not test. **How it fails:** a flaky row loops without end, against the three-round convergence rule (`context/AGENTS.md:16`). **Fix:** cap it at three fix rounds per row, then stop at the user's gate with the row recorded as Open.
- **[low / confidence: medium]** I-20a vs `context/AGENTS.md:21`, "Docs-only changes of ten lines or fewer need steps 4 and 5 only". **What is wrong:** in this repository, skill and context prose is source (D18). A ten-line prose fix run under the `bug` skill, which is BUG-018's own class, is both "any fix under the `bug` skill" (reviewer owed) and "docs-only ≤10 lines" (no review). Which rule wins is not stated. **Fix:** one clause in I-20a or ADR-022.

## Questions

1. **BLOCKING** — Preflight Step 4 (`skills/preflight/SKILL.md:65`): once I-12d routes `--docs` findings by owner, what text replaces "dispatch the `documentation` agent for them", and do Step 4's fixes ride in the release commit? The describing-text hard rule obliges the implementer to change that sentence at 5.4, and the plan gives no text, so two implementers write different gates.
2. **BLOCKING** — Task 5.2: if 5.1 records an unquoted control (a layer that did not load), does the implementer ship `new-A` ("Both load…") as written, write a new variant, or stop for the user? The plan anticipates this result but supplies text that states the opposite.
3. NON-BLOCKING — I-21a: does a not-examined entry that names a commit or criterion rather than a path remove coverage? As written it does not, and the text ships the same either way, so this is recorded as a gap.
4. NON-BLOCKING — I-21c: what scope does the one re-dispatch carry, for the reviewer and for the verifier? Also a gap. A reasonable default is the range limited to those paths.
5. NON-BLOCKING — I-22 is lettered a, b, c, e. Is an I-22d edit missing, or is the letter simply unused? The rehearsal count of 88 matches the headings present, so this is probably the latter.
6. NON-BLOCKING — 6.4's c192622 comparison needs "a second sandbox home" logged in to both CLIs. The logins live only in `/home/kesokaj/sandbox/home`. Default: `cp -a` that home and install c192622 from a `git worktree`. `internal/harness/repo.go:76` Lstats `.git`, so a worktree counts as a checkout.
7. NON-BLOCKING — Do Phase 6 harness-defect fixes get the per-bug `reviewer` of I-20 even though the session runs the installed context (D9 covers agent briefs, not main-agent triggers)? Default: yes, follow the repository's text.
8. NON-BLOCKING — Is a reasonless skip that `verifier-counts` does not seed exercised anywhere? Default: name it under Not examined.

## Criteria

1. Dead code — pass. I-8e drops the line-71 inventory sentence; its fact is generated at `internal/harness/inventory.go:321` (verified). No orphaned files: the scratch control repository and `runs/g26-precedence` are uncommitted by design.
2. Flow tracing — fail. The routing of `--docs` findings is changed at `/review` Step 5 and in the context, but not at preflight Step 4. The I-21c re-dispatch does not connect to the I-21a coverage test. I-19a `new-A` does not follow its own failure mode.
3. Standards conformance — fail on the global describing-text rule. The Phase 6 coverage claim is incomplete, and `new-A` can state a disproved claim. Otherwise pass: no Go code, frontmatter changes stay on one quoted line (`agents/reviewer.md:3`), and no `/tmp` is added.
4. Directory layout — pass. Every edited file is where `AGENTS.md` § Layout puts it, and scratch paths are under `.agents/scratch/`.
5. Simplicity (DIY-except-security) — pass with a note. The fail-closed coverage rule adds real procedure (heading names, a path test, a re-dispatch). Its gaps above are where the complexity has not yet been paid for. There are no security surfaces.
6. Decoupling — pass. The skills own their gates and agents own their templates. The only seam problem is finding 1, where the review skill and the preflight skill diverge on who owns a docs finding.

## Not examined

- I did not re-run the rehearsal (D15). `old` uniqueness was spot-checked only for I-19b–d, I-20a–d, I-21a–e, I-22c, I-24a and I-24d, and the count of 88 edit headings matches D15.
- I did not read in full: `agents/bugs.md`, `agents/changelog.md`, `agents/explorer.md`, `agents/security.md` (headings only), `skills/handoff/SKILL.md` beyond lines 20–50, `skills/retro/SKILL.md` beyond lines 56–58 and 125–135, `skills/distill/SKILL.md` beyond line 98, `skills/review-standards/` and its references, the `skills/golang/` files, the wording of `README.md` and repo `AGENTS.md`, the research record beyond grep hits, and the round-1 reports. Of the round-2 reports I read only their finding lines.
- Phases 1–4: task checks were read for the logic of their expectations, not executed. I-7a's template was read but not traced against every Phase 1–4 expectation.
- The I-9d awk command was run with gawk 5.2.1 only. It prints the two old plans and this plan, which matches 2.4's positive control. It was not run with macOS BWK awk or mawk, which the Mac or a downstream box may use.
- Nothing was run in the sandbox. I only checked the demo repositories: empty `handoffs/`, no `counts:` lines in `testing/`, `demo-app-agy` has an `origin` and `demo-app` has none, so neither new matrix check passes by default. `render | wc -l` prints 93 at HEAD (verified).
