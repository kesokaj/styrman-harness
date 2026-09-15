---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r1-a   (2026-09-15)

Range reviewed: the plan file `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (untracked, as on disk 2026-09-15) against the surfaces it names at HEAD `c192622` (working tree: only `ROADMAP.md` modified, not a named surface).

## Verdict: REQUEST CHANGES — `design needs revision` / `implementation not ready`

The worst problem: several verify and checkpoint expectations cannot be met as written. They would stop a first-pass implementer or send them back to the author.
- Task 5.3 will print `skills/preflight/SKILL.md:50`, a "50+ lines or 5+ files" line that I-18 never touches. So "every printed line also names the bug case" fails.
- The Phase 4 Checkpoint expects a verifier log that "applies the new rules". But in this checkout a dispatched verifier is the installed one from the real home, and that copy has no counts rule.
- The 1.2 positive control cannot tell a working `--docs` brief from a blind one.
- The G26 precedence test has a control marker for skills and none for agents.
- ADR-022 is appended one phase before ADR-021, in an append-only file whose author numbers entries in sequence.

The anchor work is otherwise sound: 67 of the 68 Old texts I checked occur exactly once, and `render` prints 93 lines today.

## Gaps (findings)

- **[high / confidence: high]** Plan task 5.3, I-18; `skills/preflight/SKILL.md:50`. The line reads: "dispatch the `reviewer` agent on the strongest model … when the release diff is 50+ lines or 5+ files". I-18 edits only preflight:40. So `git grep -n "50+ lines or 5+ files" -- context skills agents` still prints :50 without the bug case, and 5.3's expect fails. ADR-021's text names four surfaces and says nothing about the release-gate threshold.
  - **Decision needed:** does the release gate's reviewer threshold also get the bug-fix exception (edit :50 and ADR-021), or does 5.3's check exclude :50 with a stated reason?
- **[medium / confidence: high]** Phase 4 Checkpoint, the Done rows, D6; every in-session dispatch. This checkout has no `.claude/` directory. Only `.agents/skills -> ../skills` exists, as a git symlink (mode 120000). So Claude Code dispatches the agents and `/review` skill installed in the real home.
  - `~/.claude/agents/verifier.md` (installed 2026-09-14) has 0 lines matching `Test counts:`.
  - Scope Out and R6 rule out installing into the real home.
  - So "a `verifier` log over the phase range … which also applies the new rules (its log carries a counts line)" cannot come from a plain dispatch. The same applies to Done's "`/review --docs` over the same range", which would run the installed review skill without I-1.
  - Antigravity differs: it reads `.agents/skills` in this workspace (so it gets the checkout's skills), but there is no `.agents/agents`.
  - **Decision needed:** does each such dispatch carry the new I-15 or I-1 text verbatim in its brief? Or are those expectations dropped until the operator installs?
- **[medium / confidence: high]** Plan task 1.2 (positive control for I-1).
  - **The brief never reaches the new homes.** I-1's living docs are `README.md`, `AGENTS.md`, `.agents/{rules,skills,agents,commands}/`, specs, the testing guide and deploy. In this repo `skills/`, `agents/` and `context/` are not in that list, and `.agents/skills` is a symlink that `git grep` does not descend. I-2 adds them to this repo's `AGENTS.md`, but at 1.5, after 1.2. A reviewer following the brief literally never reads `skills/bug/SKILL.md:29` or `skills/feature/SKILL.md:15`. It flags every gate in the sentence 87a862d removed as `dropped:`, so the expect is met for the wrong reason.
  - **A careful reviewer can fail it instead.** One that does read `skills/` finds the repro gate has a home by meaning (bug:29) that lost only its autonomy qualifier, and may emit no `dropped:`. Failure modes would then "rework I-1" although the brief works.
  - **The brief has no category for BUG-018's mechanism.** A rule that moved but lost a qualifier is neither "still says the old thing" nor "no home by meaning".
  - **The brief does not say which revision to `git grep`** (HEAD or the range end).
  - **Decisions needed:**
    - add a negative control: the scope-confirmation gate, which landed with its clause at feature:15, must not be `dropped:`;
    - run 1.2 after I-2, or put the harness living-docs extension in the dispatch;
    - state whether a lost qualifier is reported as `dropped:` or as a named third finding type.
- **[medium / confidence: high]** Plan tasks 4.4 and 5.3; `.agents/docs/architecture/DECISIONS.md` (last entry ADR-020, line 157); `agents/changelog.md:78` ("Number sequentially").
  - 4.4 appends ADR-022 in Phase 4 and 5.3 appends ADR-021 in Phase 5, and the file is append-only (`context/AGENTS.md:69`).
  - Result: either ADR-020 → ADR-022 → ADR-021, or `changelog` renumbers the first append to ADR-021. In the second case 4.4's `grep -c '^## ADR-022'` prints 0 and ADR-021 appears twice.
  - **Fix:** swap the two ids, or land ADR-021 first.
- **[medium / confidence: medium]** Plan task 5.1, I-22, Risks bullet 3.
  - **The agent result is uncontrolled.** Risks says loading of "project skills and agents" is "checked by 5.1's control marker", but the only control is the skill `g26-control`.
  - **The test agent may be in the wrong format for Antigravity.** Antigravity's agent format in this harness is `agents/<name>/agent.md` (`internal/harness/plan.go:60`, `:151`), while the test drops a flat Claude-format `.agents/agents/reviewer.md`. No log under `.agents/docs/testing/` shows Antigravity loading a project agent; `skills/init-project/SKILL.md:66` only asserts it.
  - **So Outcome A can record a false cause.** If agy never loads that file, it quotes the harness description. Outcome A then writes "because the vendors resolve a clash differently" into the global context, a cause that was never observed.
  - **Decision needed:** add a uniquely named project-agent control marker, and a rule for what happens when that control fails.
- **[medium / confidence: high]** Verification, golden path, against I-23's `feature` check.
  - The criterion claims task lines carry `expect:`, phases end with a Checkpoint box, and the plan has a Cold-reader check box and `## Decisions`.
  - The check greps only `Checkpoint`, `## Decisions` and `- [ ]`. `expect:` and the Cold-reader box are never checked, and `Checkpoint` matches prose anywhere, not just a box.
  - **Fix:** add, for example, `grep -q -- '— expect:'`, `grep -q '^- \[.\] Cold-reader check'` and `grep -q '^- \[.\] Checkpoint'`, or narrow the criterion.
- **[medium / confidence: medium]** I-20 (bug and feature gate text), against the E2E edge case and `bug-no-repro`.
  - "A headless run with no user turn that reaches a gate stops there and writes the question … into the handoff's Next Steps." Both matrix drivers pass the prompt as a user turn (`claude -p "<prompt>"`, `agy -p`); what a headless run lacks is a reply.
  - A model can read the clause as not applying, ask in its final message, and write no handoff. The check `grep -rqi "repro" .agents/docs/handoffs/` then fails; the demos have an empty handoffs directory, so it really can fail.
  - **Decision needed:** word the trigger as "a run where no user reply can arrive (headless `-p`)", or similar.
- **[medium / confidence: medium]** I-19; `skills/preflight/SKILL.md:50-52`.
  - **Undefined term.** "A release path an artifact left unexamined" and "a release path Step 2 left uncovered" are never defined. "Release path" already means the ship-and-deploy process in this skill (description, :3).
  - **Unclear when an artifact is due.** A reviewer artifact is due only at 50+ lines or 5+ files, and security only for sensitive paths. Below threshold the main agent reads the diff itself and no artifact exists. Is that "a missing artifact" that gets a re-dispatch and then counts as blocking?
  - Runtime agents will gate small releases differently. **Decision needed:** define both terms.
- **[low / confidence: high]** I-9, `agents/changelog.md:94`. The Old text as written ends with a space after `` `TODO.md` ``, but the file has `` `TODO.md`, and any ``. It occurs 0 times. Failure modes catches this and asks for a D-row, but the plan's statement that its anchors are located by quoted text is false for this one.
  - **Fix:** drop the trailing space. The New text then reads correctly before ", and any".
- **[low / confidence: high]** I-24; `agents/bugs.md:69` and `:80` allow only a `status` edit to a past record.
  - I-24 asks `bugs` to record "the fix `skills/bug/SKILL.md` Step 1 (I-4) and the Phase 1 commit". BUG-018's `## Fix` section says "Not applied" and would keep saying so under `status: fixed`.
  - The Phase 1 commit SHA cannot sit inside that same commit.
  - **Decision needed:** status only, accepting the stale Fix text, or a named exception to the bugs agent's rule.
- **[low / confidence: medium]** I-21, the `CLAUDE.md` fold. The text says "never overwrite it … ask whether to fold its content into `AGENTS.md`, then make it that single line". A "no" answer has no branch: making the file that single line overwrites content that was not folded.
  - **Fix:** state the decline branch, for example keep the content and add `@AGENTS.md` as a line.
- **[low / confidence: medium]** I-10, `2026-09-13-learning-loop-project-level.md:136`. The "Otherwise it gains …" branch does not say the box becomes `[-]`. Left `[ ]`, it fails 2.4's expect ("only this plan's path"). Separately, rebuild :101 (`sh-agent update --check` from another machine after a push) is pointed at handoff 2026-09-14 Next Steps 1. That step reads "On the Mac: update, status, the three smoke checks on each CLI" and does not include the resync check.

## Questions an implementer would have to ask

1. **BLOCKING (I-18, 5.3):** does `skills/preflight/SKILL.md:50` get the bug-fix exception, or is it excluded from 5.3's check?
2. **BLOCKING (4.4, 5.3):** in which order and under which ids do ADR-021 and ADR-022 land?
3. **BLOCKING (Phase 4 Checkpoint, Done):** in-session `verifier` and `/review --docs` dispatches run the installed definitions. Does the brief carry the new I-15 and I-1 text, or are those expectations dropped?
4. **BLOCKING (1.2):**
   - Which negative control goes with the positive one?
   - Does 1.2 run after I-2, or get the harness living-docs extension in its dispatch?
   - Is a lost qualifier reported as `dropped:`?
5. **BLOCKING (5.1, I-22):**
   - What is the project-agent control marker?
   - What happens when it is not quoted?
   - May Outcome A assert a cause the test did not isolate?
6. NON-BLOCKING (I-23): extend the `feature` check to `expect:` and the Cold-reader box, or narrow the golden-path criterion. A verbatim copy works either way, but the criterion overclaims.
7. NON-BLOCKING (I-20): does "no user turn" cover `-p` runs? The text is given verbatim, but the edge-case E2E depends on how a model reads it.
8. NON-BLOCKING (I-19): what are "release path" and "missing artifact" below the reviewer threshold? The text is verbatim, but the runtime meaning is open.
9. NON-BLOCKING (I-21): what happens when the user declines the fold?
10. NON-BLOCKING (I-24): may `bugs` edit BUG-018's Fix section, and where does the Phase 1 SHA go?
11. NON-BLOCKING (I-10): does the "Otherwise" branch of :136 also flip the box to `[-]`? 2.4 implies yes.
12. NON-BLOCKING (I-9): the anchor fix follows the Failure modes D-row path.
13. NON-BLOCKING (I-14): the `golang-ecosystem.md` Old/New use mismatched backtick fences in the plan. The intended literal `` `unused`, `misspell`) `` is recoverable and occurs once at :323.
14. NON-BLOCKING (D7, I-2, I-5): `agents/documentation.md`'s lane (README, specs, testing guides, ROADMAP) and the global DoD step 3 ("AGENTS.md setup and deploy sections") do not include this repo's `AGENTS.md` "What this repo is" or Testing sections. Will the agent take the edit? The plan names exact text, so the main agent can apply it if the agent declines.
15. NON-BLOCKING (I-22): how is `agy -p` run "with no tools allowed"? The sandbox reference documents only `--mode accept-edits` plus allow-rules.
16. NON-BLOCKING (Phase 2): G12 lands in 2.1, yet the Phase 2 Checkpoint has no verifier log. Phase 3 names this, so it reads as intended.

## Checked and clean

- **Anchors:** 67 of 68 Old texts occur exactly once (script over every quoted anchor); only the I-9 anchor fails.
- **Anchor ordering:** I-20's review anchor lives inside I-12's new text, so the 3.1 → 5.5 order works. I-4's clause is replaced by I-20 as planned.
- **Render and list:** `go run ./cmd/sh-agent render` exits 0 and prints 93 lines. "can also be typed as" appears once, guarded by `internal/harness/harness_test.go:708` via `inventory.go:321`. `list` prints `### Skills (27)` and `### Subagents (7)`, so I-5's guide text is true.
- **Vendoring and unit tests:** `skills/VENDORED.md` lists none of the edited skills. No new text contains `/tmp` (`TestNoTmpInHarnessProse`). The reviewer description stays one quoted line.
- **Stale plan boxes:** the unticked boxes in the two stale plans are exactly the lines I-10 names (rebuild :58, :97–:99, :101, :113; learning-loop :136).
- **Matrix harness:** clones are plain `git clone` (`.agents/scratch/skilltest/run-one.sh:10`), so `@{u}` in the `bug-no-repro` check resolves. The demo repos' handoffs directories are empty, so that check can fail.
- **I-10 real-home step:** `status` (`cmd/sh-agent/main.go:227`) showed no write calls in its body, which fits the read-only step.

## Criteria

1. **Dead code:** pass. The plan removes the stale "Workflow skills" sentence (D4) and the TODO pointer in `changelog.md`; 2.5 greps for leftovers.
2. **Flow tracing:** fail. The threshold still at `preflight:50`, the ADR order, dispatches that load the installed rather than the checkout definitions, and the controls in 1.2 and 5.1 that cannot fail for the right reason.
3. **Standards conformance:** mostly pass. The BUG-010 anchor discipline holds except I-9. The global rule that silence from instrumentation is not evidence is not met by 1.2's control or the G26 agent result.
4. **Directory layout:** pass. Every edit lands in an existing file; new artifacts go under `testing/` and `reviews/`.
5. **Simplicity (DIY-except-security):** pass. Text edits only, no Go code, no dependency, no security surface.
6. **Decoupling:** pass with a note. The single-plan choice is justified by its ordering constraints (Rejected row). Ownership seams with `documentation` (Q14) and `bugs` (Q10) are not settled.

## Not examined

- The research record beyond its G09, G12–G14, G17 and G26 rows and lines 123–155. I did not check that each group's intent is fully carried by its I-block.
- `agents/security.md`, `agents/explorer.md` and `agents/reviewer.md` bodies beyond the I-16 anchor count. The six-agent read came back cut to a saved-output preview, and I read only the verifier template through targeted greps afterwards.
- `skills/init-project` assets, the `feature` Step 5 text, `context/dev-setup.md`, and the full `harness_test.go` prose assertions (only its `strings.Contains` literals were grepped).
- No CLI was run in the sandbox, and nothing was installed into any home.
