---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r2-c   (2026-09-15)

Range reviewed: `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (staged, round 2), against the surfaces it names at `c192622`

## Verdict: REQUEST CHANGES — `design needs revision` / `implementation not ready`

The worst problem is in I-21e. It adds "a reviewer REQUEST CHANGES verdict" and "a reviewer finding of high severity" to the blocking list for Steps 2–4. Step 4's `--docs` reviewer keeps the reviewer template, verdict line included. Step 4 also sends that reviewer's findings to `documentation` in the same run, and nothing re-reviews afterwards. The plan does not say whether a docs-lens REQUEST CHANGES blocks the release, or whether the `documentation` fix clears it. Two careful implementers would ship different gates (BLOCKING Q1).

The G26 precedence test has a gap. Outcome B only requires the project markers to be quoted; it never checks that the harness copy's description is absent. A CLI that loads both same-named copies would therefore keep line 37's priority claim with no evidence of priority.

I-22b adds a new bug-skill gate (a diagnosis that disagrees with your hypothesis). No decision row records it, and the research adaptation for G21 does not include it. It can end the existing `bug` matrix row, which 6.4 expects to pass.

The rest holds up. I re-ran a copy of the rehearsal against `c192622`: 82 edits parsed, 0 failures, 26 files, for both I-19a variants. On the rehearsed tree every grep expectation in 5.2, 5.4, 5.5, 5.6 and 6.1 prints what the plan says.

## Findings (gaps)

- **GAP-1 [medium / confidence: high]** plan:665 (I-19, 5.1, D12 at plan:1022) against `context/AGENTS.md:37`.
  - **What is wrong:** "Outcome B holds only when both CLIs quote both control markers and both project markers." The prompt asks for "the description you have for the skill named review". A CLI that lists both same-named copies, or a model that quotes both, gives a project-marker hit. That records B and keeps "the project layer has priority on both vendors", but it does not show priority.
  - **Why it matters:** the research record gives the doubt this test exists for (`2026-09-15-research-…md:75`, "Claude Code's docs report as untrue for skills"). Its source is a documentation summary, not a live test (research :155).
  - **Fix:** record the harness `review` / `reviewer` description too, quoted or not. B requires every project marker quoted and every harness description not quoted; anything else is A.
- **GAP-2 [medium / confidence: high]** plan:788 (I-21e) and D2 (plan:1012), against the rehearsed `skills/preflight/SKILL.md` Step 4, `agents/reviewer.md:73` and `skills/review/SKILL.md:48`.
  - **What is wrong:** the reviewer template has `## Verdict: … REQUEST CHANGES` (`agents/reviewer.md:73`), and the `--docs` brief says "Keep your report template" (`skills/review/SKILL.md:48`). So the Step 4 docs review emits a verdict. Blocking now covers "Steps 2–4" and includes "a reviewer REQUEST CHANGES verdict".
  - **Why it fails:** Step 4 still says to dispatch `documentation` for the docs gaps and record them. It does not say whether that fix resolves them, whether a re-review is needed, or whether the gate goes NO-SHIP. A stale README found at release time would ship under one reading and block under the other.
  - **Decide:** either docs-lens findings are blocking only at high severity (or never), or they block until a re-run of Step 4 comes back clean.
- **GAP-3 [medium / confidence: medium]** plan:810 (I-22b), against `skills/bug/SKILL.md:37`, research record :70, and `.agents/scratch/research/wf2-round1-details.txt:133,135`.
  - **What is wrong:** today Step 2 lets the agent re-investigate and choose. I-22b makes "a diagnosis that disagrees with your hypothesis when the evidence does not settle it" a second gate that "always waits for the user". The research adaptation for G21 names only "a gate section for bug (missing repro, inferred-repro path)", and no D-row records the addition.
  - **Consequence:** in print mode a reached gate "ends the work there". The existing `bug` matrix row (`skill-matrix.md:23`, a `claude -p` / `agy -p` run with "Do not ask me anything") can now stop without a fix. 6.4 (plan:987) expects it to exit 0 on both CLIs, and the Failure modes section has no entry for a row stopping at a gate by design.
  - **Fix:** add a D-row taking the Step 2 gate into scope (or drop it), and state what 6.4 records when a row stops at a designed gate.
- **GAP-4 [low / confidence: high]** plan:834 (I-22e), against the rehearsed `skills/review/SKILL.md:60`.
  - **What is wrong:** the cap options include "run one more round", but the same line still opens Step 5 with "At most three rounds." The rehearsed line contains both.
  - **Fix:** reword to "at most three rounds without the user's word", or drop the option.
- **GAP-5 [low / confidence: high]** plan:976 (5.3 expect), against the rehearsed `skills/preflight/SKILL.md:40`.
  - **What is wrong:** the expect says "the Step 1 line ending 'or is a bug fix.'". In the rehearsed file that phrase sits mid-paragraph, followed by "Then commit with a descriptive message…".
  - **Why it fails:** Step 4 of the feature skill ticks a box only when the evidence meets its expect, so a literal reading fails the check.
  - **Fix:** say "containing".
- **GAP-6 [low / confidence: high]** plan:905 (golden path) against plan:857 (I-24a).
  - **What is wrong:** the criterion claims "`expect:` on its task lines" and "a Checkpoint box per phase". The check greps for one `expect:` anywhere and one `- [.] Checkpoint` anywhere, so a three-phase plan with one Checkpoint box passes. This is describing text that overclaims coverage (global hard rule).
  - **Fix:** narrow the criterion wording, or have the check compare the Checkpoint box count with the `### Phase` heading count.
- **GAP-7 [low / confidence: medium]** plan:987–988 (6.4 and the Phase 6 checkpoint), against `AGENTS.md:46` ("A verification log states what it did not exercise").
  - **What is wrong:** the chosen rows do not exercise several interfaces:
    - I-21's fail-closed gate: the `preflight` check `grep -q "SHIP"` (`skill-matrix.md:26`) matches NO-SHIP too.
    - I-12 and I-22e's `--fix` ledger and cap.
    - I-18's move of index rows to the Graduated table.
    - I-23b's CLAUDE.md fold: the demo's CLAUDE.md is already the single line.
    - I-19b–d.
  - **Gap:** the checkpoint expect only says "pass"; it does not require the e2e log to name these.
- **GAP-8 [low / confidence: medium]** plan:992 (Done `reviewer` row) against plan:32 and D9 (plan:1019).
  - **What is wrong:** D9 says every dispatch whose behaviour the plan changes carries the changed text. The reviewer's behaviour changes through I-14, I-15, I-17a and I-20c. The Done verifier row cites D9, but the reviewer row pastes nothing.
  - **Consequence:** the final review runs on the installed, old `review-standards`, without the new describing-text and suppression items.
  - **Fix:** paste the new `review-standards` §1–§3 and I-17a text into that brief, or add a D-row saying the Done review runs on the installed definitions by choice.

## Questions

- **Q1 BLOCKING (GAP-2):** does a Step 4 docs-lens reviewer's REQUEST CHANGES verdict, or a high-severity docs finding, block the release? If it does, what clears it: the `documentation` dispatch, or a clean re-run of Step 4?
- **Q2 NON-BLOCKING (GAP-1):** should outcome B also require that the harness copy's description is not quoted? The plan states a rule, so implementers converge, but on a test that cannot tell "both loaded" from "project wins".
- **Q3 NON-BLOCKING:** I-20d (plan:728) says "or is a bug fix" for uncommitted work in preflight Step 1. How does the gate tell? Suggested default: the dirty tree adds a `bugs/BUG-*.md` record, or this session ran the `bug` skill on it.
- **Q4 NON-BLOCKING:** I-21a/b say a report's not-examined section "names none of that commit's paths". Does naming a directory or an area count as naming every path under it? Suggested default: yes.
- **Q5 NON-BLOCKING:** I-22a and I-22b say "neither is the skill invocation" is a waiver. Do waiver words inside the invocation's arguments count? Examples are the `feature` row's "Scope is confirmed exactly as written" (`skill-matrix.md:27`) and this plan's start prompt (plan:8). The research adaptation asked for example phrases that match the matrix prompts (`wf2-round1-details.txt:135`), and the plan gives none. Suggested default: the words count; only the bare invocation does not.
- **Q6 NON-BLOCKING (GAP-3):** in 6.4, when a row ends at a designed gate, is that a pass, a bug record, or `[-]`? And what does "a `[-]` box" mean for one failing row inside a single 14-row task?
- **Q7 NON-BLOCKING:** the `verifier-counts` base does not give the content of `internal/greet/extra_test.go`; a reasonable implementer writes any passing test.

## Checked and holding

- **Rehearsal:** re-ran a copy of `.agents/scratch/rehearsal/rehearse.py` into `.agents/scratch/coldread-r2-c/`. Both variants: `edits parsed: 82; failures: 0; files touched: 26`, exit 0.
- **Expectations on the rehearsed tree (variant A):**
  - 5.2 prints `1`; `2` and `1`.
  - 5.4 prints `1` ×4.
  - 5.5 prints `1`/`1`, `1`, `0`, `1`.
  - 5.6 prints `1` ×2.
  - 6.1 prints `1` ×2.
  - 5.3's `git grep` prints exactly the two preflight lines (40 and 50); only the expect's wording is off (GAP-5).
- **Threshold wording:** `git grep` for "50+ lines / 5 files" outside dated records finds only the four I-20 surfaces plus `preflight:50` (D8). The precedence claim appears only at `context/AGENTS.md:37`.
- **ADR numbering:** `DECISIONS.md` ends at ADR-020, so ADR-021 and ADR-022 are free.
- **Not-examined section names:** `reviewer.md:23`, `verifier.md:21` and `security.md:21` all use "Not examined", which matches D14.
- **Sandbox demo:** both demo repos track only `.gitkeep` under `handoffs/`, so the `bug-no-repro` handoff grep cannot match stale content. Their `testing/` logs have no `counts:` line.
- **Verifier anchor:** `agents/verifier.md` has a `## What to run` section (line 40), so D9's paste target exists, and I-16a lands inside it.

## Criteria

1. **Dead code** — pass. I-8e removes the workflow-skills sentence, whose fact lives in the generated inventory (D4). No orphaned surfaces.
2. **Flow tracing** — fail on two hops.
   - The preflight blocking definition does not trace to the Step 4 docs reviewer and its `documentation` hand-off (GAP-2).
   - The new bug Step 2 gate does not trace to the `bug` matrix row and 6.4 (GAP-3).
   - Holding: the rehearsal anchors, and I-2a → I-22c ordering.
3. **Standards conformance** — mostly pass. BUG-010 discipline holds (rehearsal). BUG-004 exposure is one quoted description line. Two gaps against global rules:
   - Describing text overclaims in GAP-5 and GAP-6.
   - The not-exercised list is missing (GAP-7, `AGENTS.md:46`).
4. **Directory layout** — pass. Logs go under `testing/`, temporary trees under `.agents/scratch/`, and the sandbox run under `runs/`.
5. **Simplicity (DIY-except-security)** — pass. No code and no dependency. The fail-closed re-dispatch is one step deep.
6. **Decoupling** — pass. `skills/bug/SKILL.md` is edited in Phases 1, 2, 4 and 5, but always in application order. Documentation's lane is handled by D7.

## Not examined

- I-1 to I-18 were read and anchor-checked through the rehearsal, but not judged in depth against each group's research intent. Of the research record I read only the G20–G26 rows, lines 123–155, and the EG-07 notes.
- `sh-agent render | wc -l` on the rehearsed tree: not re-run; D15's 93 is taken as reported.
- No live CLI run. Vendor precedence, and how Antigravity loads `.agents/agents/<name>/agent.md`, are unverified by me.
- The docs-lens control design (1.2), `agents/bugs.md`, `agents/changelog.md`, `README.md`, and the handoff, retro and distill skills beyond their anchor lines.
- The round-1 reports were read only through grep, to avoid re-raising settled items.
