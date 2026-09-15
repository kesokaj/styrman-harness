---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r1-c   (2026-09-15)

Range reviewed: plan file `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (untracked, read in full, 417 lines) against the surfaces it names at HEAD `c192622` (working tree: only `ROADMAP.md` modified). Cold-reader check, round 1.

## Verdict: REQUEST CHANGES — `design needs revision` / `implementation not ready`

The worst problem is that three checkpoint or verify expectations cannot be met as written. Task 5.3's grep will print `skills/preflight/SKILL.md:50`, a release-threshold line I-18 never touches. Phase 4's checkpoint expects the verifier to follow the new rules, but no dispatch from this checkout can load the edited agent definitions. I-10's likely branch leaves a `- [ ]` box that task 2.4 says must be gone. Two verifications also cannot fail when they should. The 1.2 positive control's answer is already in `bugs/INDEX.md`, which the reviewer reads before anything else. The `bug-no-repro` matrix check never looks at untracked files, and its `@{u}` clause is empty on the demo, which has no remote.

The plan's mechanics are otherwise sound. All 66 quoted old texts occur exactly once in their files (checked with `grep -o -F | wc -l` for each). The `render | wc -l` baseline really is 93. The ADR numbers, BUG-018 record, handoff path and Mac verification log it cites all exist.

## Findings

### Gaps (must be decided before implementation)

- **[high / confidence: high]** Plan:373 (5.3) vs `skills/preflight/SKILL.md:50`. The check `git grep -n "50+ lines or 5+ files" -- context skills agents` expects every printed line to name the bug case. I-18 (plan:248-251) edits `preflight:40` but not `preflight:50` ("when the release diff is 50+ lines or 5+ files"). That line still prints without the bug case, so 5.3 cannot pass. The plan also never decides whether the release gate must review a sub-threshold bug fix. Fix: either add a preflight:50 edit to I-18 and ADR-021's list of changed files, or narrow 5.3's expect to exclude the release-diff threshold and write that as a D-row.
- **[high / confidence: high]** Plan:368 (Phase 4 checkpoint), plan:360, plan:346 and D6 (plan:412). The plan assumes the repo's edited agents and skills are what its own dispatches load. That is false on both vendors:
  - The repo has no `.claude/` directory and no `.agents/agents`. The only project link is `.agents/skills -> ../skills`, which Antigravity reads and Claude Code does not.
  - `~/.claude/agents/verifier.md` has 0 "Test counts" lines, and installing into the real home is out of scope (plan:16, R6).
  - The Phase 4 verifier therefore runs on the old definition, and "its log carries a counts line" cannot be expected from it.
  - The same applies to 1.6's `/review --project --docs` "with the I-1 brief". On Claude Code the `/review` skill loaded is the installed copy, whose `--docs` brief is the old one.
  - It also applies to the Phase 3 note "(G12 is in the skill from Phase 2)".

  Fix: state how the new text reaches each dispatch (paste I-1 or I-15 verbatim into the brief, or dispatch from a sandbox session on the scratch home), or drop those expectations.
- **[high / confidence: high]** Plan:176 (I-10, learning-loop :136) and plan:353 (2.4).
  - In the "Otherwise" branch the line gains a note but stays `- [ ]`, so 2.4 ("only this plan's path") fails.
  - "Otherwise" is the near-certain branch. `binaryOlder` (`cmd/sh-agent/main.go:398-406`) warns whenever the installed rev is not a prefix of repo HEAD. By Phase 2, HEAD already contains the Phase 1 commit, which the real-home binary cannot have been built from.
  - The same branch moves the item to "this plan's closing handoff", a file that does not exist yet.

  Fix: close the box as `[-]` in both branches, and name an existing handoff or say that the closing handoff's Next Steps gets the item at Done.
- **[medium / confidence: high]** Plan:342 (1.2 positive control) and failure mode plan:315. The control cannot tell whether the I-1 brief can emit a `dropped:` finding:
  - The reviewer's pre-flight step reads `.agents/docs/bugs/INDEX.md` (`agents/reviewer.md:32`), whose BUG-018 row already records that the bug skill's missing-repro gate lost its autonomy clause in 87a862d. A `dropped:` finding can come from that row rather than from the grep procedure. That is the "green run that never hit the thing under test" the global hard rules forbid.
  - The expected finding is also vague. The removed sentence is `87a862d^:context/AGENTS.md:105` ("Four gates need the user and are never skipped…"). Part of it survives in `skills/feature/SKILL.md:15`, so a reviewer can reasonably call it homed. The failure mode would then send I-1 back for rework even though I-1 behaved correctly.

  Fix: use a removal the trail does not record (a seeded removal on a scratch branch), and name the exact old file:line the finding must cite.
- **[medium / confidence: high]** Plan:302 (`bug-no-repro` row) and plan:332 (E2E edge case). The check cannot see the behaviour it names, and the gate text does not decide what the run must do:
  - (a) `git diff --quiet -- cmd internal` ignores untracked files. A test-first file created under `internal/` and left uncommitted passes the check.
  - (b) `[ -z "$(git log --oneline @{u}..HEAD -- cmd internal)" ]` is always true on the demo. `/home/kesokaj/sandbox/demo-app` has no remote, so `@{u}` fails on stderr and stdout stays empty. A committed code change also passes.
  - (c) I-20's clause covers only "a headless run with no user turn" (plan:274, :280). A `claude -p` or `agy -p` prompt is a user turn. A model that reads the gate literally asks in its reply and writes no handoff, and the check's `grep -rqi "repro" .agents/docs/handoffs/` fails. The demo's `handoffs/` is empty today.

  Fix: check with `git status --porcelain -- cmd internal` plus `git log --oneline <base>..HEAD -- cmd internal`, and define whether a `-p` run counts as headless with no user turn.
- **[medium / confidence: high]** Plan:331 vs plan:301. The golden-path E2E criterion claims four properties: task lines carry `expect:`, phases end with a Checkpoint box, a Cold-reader check box exists, and a `## Decisions` section exists. The updated `feature` row check greps only `Checkpoint`, `## Decisions` and `- [ ]`. The `expect:` clause and the Cold-reader box are untested, yet "yes when the feature row's check exits 0". Fix: add `grep -q "— expect:"` and `grep -q "Cold-reader check"`, or drop those two claims from the criterion.
- **[medium / confidence: high]** Plan:280 (I-20 bug gates) vs `skills/bug/SKILL.md:29`. The new gate offers "let me infer one (the inferred repro is printed before Step 2 runs)". Step 1 keeps "Never invent a repro; a wrong-shaped guess wastes the diagnosis", and no task edits that sentence. The skill would ship contradicting itself, which breaks the describing-text hard rule. Fix: add an edit to Step 1's sentence ("never invent one unless the user chose 'infer'") or drop the infer option.
- **[medium / confidence: medium]** Plan:293-296 (I-22, G26 test). The only control is a skill (`g26-control`); there is none for the agent layer.
  - Antigravity's agent form is `agents/<name>/agent.md` (`internal/harness/plan.go:60`, `:151`). A Claude-format `.agents/agents/reviewer.md` may simply go undiscovered by agy.
  - Undiscovered looks identical to "the harness copy wins". Outcome A's text would then assert a cause the test cannot establish ("because the vendors resolve a clash differently").
  - I have not checked whether Claude Code shows subagent descriptions to the model when the run allows no tools.

  Fix: add a control agent `g26-control-agent`. Write Outcome A's "because" clause only when both controls are quoted; otherwise use cause-neutral wording.
- **[medium / confidence: medium]** Plan:265-268 (I-19).
  - "A missing artifact … gets one re-dispatch; still uncovered, it is a blocking finding" does not say which artifacts are owed. The reviewer row is threshold-gated ("Below threshold, read the diff yourself", `skills/preflight/SKILL.md:50`) and the security row is conditional (`:51`). One reading makes every small release re-dispatch and then go NO-SHIP.
  - "release path" is undefined.
  - "its Not examined section" matches only the agent templates. The merged record calls it "Not checked" (`skills/review/SKILL.md:56`) and the gate report "What was not checked" (`skills/preflight/SKILL.md:75`).

  Fix: define the owed rows, define release path as `git diff --name-only <base>..HEAD`, and name all three section titles.
- **[low / confidence: high]** Plan:366 (4.4, ADR-022 in Phase 4) and plan:373 (ADR-021 in Phase 5). Both go into an append-only `DECISIONS.md`. The `changelog` agent reads numbering from the newest entries and numbers sequentially (`agents/changelog.md:12`, `:78`), so it may renumber the Phase 4 ADR to 021, and 4.4's expect fails. Otherwise the file permanently reads 020, 022, 021. Fix: swap the two numbers or the order the phases append them.
- **[low / confidence: high]** Plan:350 (2.1 verify, "at least 7"). `grep -c` counts lines, not patterns. Block A's Checkpoint line (plan:111) matches two patterns, and `## Step 3: Cold-reader check —` already matches today. With `^## Decisions` missing the count is still 8, and with two patterns missing it is still ≥7, so the check cannot catch a partial write. Fix: count each pattern separately.
- **[low / confidence: medium]** Surfaces (plan:22) miss some text that describes changed behaviour:
  - `skills/preflight/SKILL.md:65` restates the `--docs` lens's file list, which I-1 widens.
  - The BUG-018 record will cite I-4's wording as the fix, and I-20 deletes that wording in Phase 5 (plan:282).

  Both break the describing-text hard rule once the phases land.

### Readiness questions

See the structured list. Five are BLOCKING: Q1 (preflight:50), Q2 (which definitions the dispatches load), Q3 (I-10's else branch), Q4 (whether a `-p` run is headless with no user turn), Q5 (the 1.2 control's pass criterion). The rest are NON-BLOCKING:
- Q6: the G26 agent control.
- Q7: the ADR order.
- Q8: whether Phases 1 and 2 need a verifier log, given that Step 4 of Block A ties it to "changed source".
- Q9: I-15's "full Go test command" in a non-Go, non-frontend project.
- Q10: the malformed code spans at plan:42 (I-2) and plan:206 (I-14). The literal is recoverable from 1.5's grep and the file, and it occurs once.

## Criteria

1. Dead code — n/a (plan review). I-7 removes one sentence from context line 71 on purpose (D4). Its replacement sentence exists in `internal/harness/inventory.go:321` and is covered by `harness_test.go:708`. No orphan.
2. Flow tracing — **fail**. The edited agent and skill definitions never reach the dispatches that are supposed to exercise them (no `.claude/`, no `.agents/agents`). The I-18 threshold edit misses `preflight:50`. I-10's branch and 2.4 disagree.
3. Standards conformance — **fail** on the global hard rules. Two controls cannot fail (1.2, `bug-no-repro`), and describing text goes stale (`bug:29` vs the infer option, `preflight:65`). Agent frontmatter rules (one-line quoted `reviewer` description) are respected.
4. Directory layout — pass. Every new artifact path (testing logs, `reviews/`, `DECISIONS.md`, the Graduated table in `INDEX.md`) matches the repo's layout.
5. Simplicity (DIY-except-security) — pass. Text edits only, no Go, no dependencies, and there is no security-sensitive surface.
6. Decoupling — pass with a note. Task ownership is clear (documentation, changelog, bugs). I-4 then I-20 rewrite the same clause in two phases, which is a deliberate order constraint (plan:416).

## Not examined

- Whether the plan's group mapping G01–G26 matches `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md`. The research record was not read.
- The full text of `skills/retro/SKILL.md`, `skills/handoff/SKILL.md`, `skills/init-project/SKILL.md` and the review-standards references beyond each anchor's uniqueness and immediate context. Also the init-project template's "template command lines" that I-21 deletes.
- Vendor behaviour that cannot be run from here: whether agy discovers Claude-format `.agents/agents/*.md`, and whether Claude Code with no tools allowed exposes subagent descriptions.
- `sh-agent render | wc -l` after the edits. Not simulated; the 93 baseline was checked at HEAD.
- `.agents/docs/deploy/remote.md` § Downstream mirrors; ADR-009's text; `internal/harness/apply.go:173` (`os.CreateTemp("", …)`, whether `status` or `diff` writes to `/tmp`).
- `.agents/docs/plans/2026-09-15-coldread-handoff-2026-09-15.md`.
