---
owner: main
date: 2026-10-07
---
# Plan — lean process: fewer review rounds, cold agents and long sessions   (2026-10-07)

## Why
The operator asked whether the harness wastes turns and tokens (2026-10-07). An A/B test (`.agents/scratch/2026-10-07-0842-harness-ab-crm-report.md`; not committed) built the same CRM booking service three ways:

| Arm | Cost | Wall time | Acceptance test | Main turns | Subagents |
|---|---|---|---|---|---|
| No harness | $0.54 | 2.4 min | 25/25 | | |
| Harness loaded | $1.11 | 5.0 min | 25/25 | | |
| Harness, full process | $8.46 | 30.8 min | 25/25 | 87 | 13 |

- About 80% of the full-process run's wall time was process: three plan cold-reader rounds, a verifier per phase, separate documentation and bugs agents, and review rounds.
- The operator's own transcripts show the same pattern: a release session ran 28 agents, and main threads ran 220 turns, one resumed at a 560k context.

The operator approved recommendations 1–6 and 8 and struck 7 ("go your recommendations", 2026-10-07).

## Scope
**In:**
1. **One cold-reader round.** The plan's cold-reader check runs once. Its residue becomes `Open:` rows with defaults. The check is skipped for a plan with one phase that touches fewer than 5 files.
2. **No verifier per phase.** The main thread runs the full test command bare at each checkpoint. One `verifier` per task, at DoD step 1. Review fixes are re-checked by the main thread's bare test run, not by a new verifier dispatch.
3. **The main thread writes small records itself.** That covers the `changes/` entry, the bug record and docs edits of 30 lines or fewer. `changelog`, `bugs` and `documentation` are dispatched only above that, or when the main thread lacks the evidence.
4. **Review rounds capped at 2,** not 3. Round 2 checks the fixes only: it resumes the same agent where the CLI can (Claude Code), and is otherwise a new dispatch carrying the fix lines only.
5. **One reviewer dispatch covers the code and docs lenses together.** `security` is still separate, on its existing triggers.
6. **A new task starts in a fresh session.** After `/handoff`, the agent recommends a fresh session (`/clear`, or a new `agy` conversation) that resumes from the handoff, and does not take up a new task in a session that has closed one.
8. **A BUG-NNN only for what escaped a gate,** or for a user correction. A defect that review catches before commit is a review finding, recorded in the review's `## Outcome`, not a bug.

**Out:** 7 (a light scaffold): the unscaffolded path in `context/AGENTS.md` § Work already is the light tier.

**Follow-ups noted, not built:**
- the preflight gate's own round count beyond what the review skill sets
- dream and retro cadence

## Surfaces touched
- `context/AGENTS.md` — DoD 1-5, the failure line (:33), the subagent line (:35). Must not grow: 14,727 of 15,000 bytes (cold read Q4); detail moves into the skills.
- `skills/feature/SKILL.md` — Gates (:15), template box (:57, :61), Step 3 (:87-93), Step 4 (:99, :101), Step 6 (:111-113).
- `skills/review/SKILL.md` — Step 3 lens table (code+docs as one dispatch), Step 4 Rounds (:94) and merged record, Step 5 (:100, :108, :110, :112).
- `skills/preflight/SKILL.md` — Step 1.1 (:44 dispatches), Step 2 code review and Step 4 docs lens (:59, :78) as one dispatch, round cap.
- `skills/bug/SKILL.md` — Gates (:15), Step 0 "what counts" (:21), Step 5 (:55).
- `skills/handoff/SKILL.md` — Step 3.2 (:65), Step 7 (:182), cross-reference (:191).
- `skills/operate-gcp/SKILL.md:83` and `skills/init-project/assets/**/bugs/INDEX.md` — the bug trigger list, wherever it is restated.
- `agents/bugs.md` (:3, :20, :78), `agents/changelog.md`, `agents/documentation.md` — triggers and "dispatch once per task" wording.
- `README.md:259` — only if it states a changed rule.
- Every other restatement is found by `bash skills/review/scripts/restatement-grep.sh` before review.

## Interfaces
None: these are prose rules. No code, flag, or binary behaviour changes.

## Failure modes
- **Over-cutting:** a defect that a removed round would have caught ships.
  - Mitigation: the reviewer still runs once per task, and security still runs on its triggers.
  - Measure: re-run the full-process experiment; the acceptance test must stay 25/25.
- **Size pins:** a context edit breaks `TestRenderedContextUnderTarget`, whose ceiling is 15,000 bytes. Mitigation: the context only shrinks or holds its size.

## Verification
- **Unit:**
  - `gofmt -l cmd internal && go vet ./... && go test -race -timeout 30m ./...` passes, size pins included.
  - `go test ./internal/harness -run TestSkillListingBudget` passes when a description changes.
- **Install:** `go run ./cmd/sh-agent install --home .agents/scratch/home-lean` succeeds.
- **Behavioural:** re-run the full-process arm of the A/B kit against the scratch-home install. Expect:
  - the acceptance test at 25/25
  - at most 1 cold-reader round, 1 verifier and 1 reviewer, plus security if triggered
  - cost and wall time below the $8.46 / 30.8 min baseline (one run, n=1)
- **E2E criteria:** a fresh `claude -p` session given only the new harness builds the CRM spec. It records no BUG-NNN for review-caught defects, and dispatches no verifier before the DoD.

## Phases
- [x] Cold-reader check — one round per D1: `reviews/2026-10-07-review-lean-process-plan-cold-read.md` (design needs revision / not ready). Every gap and blocking question is answered by D6-D14 below, none needed the operator, and there is no second report.
### Phase 1 — rule edits
- [x] Edit `context/AGENTS.md` (items 2, 3, 5, 6, 8) — verify: `go test ./internal/harness -run TestRenderedContext -count=1` — expect: PASS
- [x] Edit `skills/feature/SKILL.md` (items 1, 2, 3) — verify: read the diff — expect: no "three-round" or per-phase verifier left
- [x] Edit `skills/review/SKILL.md` and `skills/bug/SKILL.md` (items 4, 8) — verify: grep `three` and `third round` — expect: no cap of three left
- [x] Edit `skills/handoff/SKILL.md` (items 3, 6) and the agent descriptions (item 8) — verify: `go test ./internal/harness -run 'TestAgentDescriptionRules|TestSkillListingBudget|TestLint' -count=1` — expect: PASS
- [x] Restatement grep over the key words of each changed claim — expect: every hit updated or a reason written in this plan. Also updated: `skills/preflight`, `skills/retro`, `skills/operate-gcp`, `agents/changelog.md`, `agents/documentation.md` and three init-project templates. The project `AGENTS.md`'s own `/review --docs` lines are kept: the combined dispatch satisfies them.
- [x] Checkpoint — expect: the full test command passes, run bare — evidence: `go test -race -timeout 30m ./...` exit 0, 6 packages ok (before the review fixes, which touched only markdown; `go test ./internal/harness` passed after them; the context is 14,719 bytes, down from 14,727), 2026-10-07
### Phase 2 — measure
- [x] Scratch-home install, then re-run the full-process A/B arm against it — verify: `accept.py` plus `metrics.py` — expect: 25/25, fewer agents and rounds than the baseline; numbers recorded in the changes entry
- [x] Checkpoint — expect: the measurement recorded — evidence: `lean-full` against `.agents/scratch/home-lean`, 25/25, $7.27 / 29.3 min / 6 agents / 62 main turns, against $8.46 / 30.8 / 13 / 87, with 0 bug records (4 before); numbers in the changes entry, 2026-10-07
### Phase 3 — Done
- [x] One `reviewer` over the range (code and docs lenses), fixes, ADR-054, the `changes/` entry, commit, push, `/handoff` — review `reviews/2026-10-07-review-worktree-code-docs.md`: 2 rounds, approved, with an `## Outcome` section

## Risks
Assumptions:
- Claude Code can resume a finished subagent with `SendMessage`. Checked by this session's own tool list.
- Antigravity has no equivalent, hence the fallback in item 4.

## Decisions
- D1 — 2026-10-07 — This plan gets one cold-reader round, under the rule it introduces — by: the operator's approval "go your recommendations".
- D2 — 2026-10-07 — Item 7 is struck: the unscaffolded path is already the light tier — by: the operator's approval.
- D3 — 2026-10-07 — Item 6 is a fresh-session rule, not a context-size threshold, because the model cannot measure its own context — by: implementer, put to the operator in the scope table.
- D4 — 2026-10-07 — The cutoff for the main thread writing docs itself is 30 changed lines — by: implementer.
- D5 — 2026-10-07 — The cold-reader skip is for a plan with one phase that touches fewer than 5 files — by: implementer.
- D6 — 2026-10-07 — Item 1 (cold read Q1, BUG-042): exactly one report and never a second, so the count stays one a headless run cannot argue past. After that report, every gap and BLOCKING question is either answered in the plan (a D-row by the implementer, or a rebuttal) or, when only the user can settle it (scope or a requirement), becomes the user gate. At that gate print mode ends with the question in the handoff, as today. A skipped check is ticked `[-] Cold-reader check — skipped: one phase, <5 files`. The three-round cap leaves the gates list; "no code until the check is resolved" stays — by: implementer, cold read round 1.
- D7 — 2026-10-07 — Item 5 (Q2): `--code` and `--docs` in one `/review` run go to ONE `reviewer` dispatch with slug `<scope>-code-docs`. Its brief is the code brief followed by the Docs brief under its own heading, and the artifact is `reviews/<date>-review-<scope>-code-docs.md`, recorded as two lenses in the merged record. Preflight's Step 2 code review and Step 4 docs lens become that one dispatch; docs findings stay advisory there — by: implementer.
- D8 — 2026-10-07 — Item 8 (Q3): what is excluded is a defect that the review, tests or verifier of the same uncommitted change find and that is fixed before that change's commit. It is recorded in the review's `## Outcome`, not as a BUG. Everything else on the existing list still counts, and "a review finding that was a real defect" becomes "…on committed work". Every restatement is updated — by: implementer.
- D9 — 2026-10-07 — Size pin (Q4): the context must not grow. DoD 2's docs routing and DoD 4's agent detail move to the review and feature skills, and the measured byte count goes in the changes entry — by: implementer.
- D10 — 2026-10-07 — Item 3 (Q5): the 30-line cutoff applies to docs only. The `changes/` entry and the bug record are always the main thread's, unless it lacks the evidence (then `changelog` or `bugs`) — by: implementer.
- D11 — 2026-10-07 — Item 4 (Q6): the cap of 2 is restated at review :94/:100/:112, bug :15 and feature :15, all updated. Resuming a finished subagent with `SendMessage` is verified in this session (the tool is offered to the main thread); every other CLI falls back to a new dispatch carrying only the fix lines — by: implementer.
- D12 — 2026-10-07 — Item 6 (Q7): recommend, never refuse. "Closed a task" means this session wrote a handoff whose Next Steps queue nothing for the agent. The recommendation goes in the handoff's Step 7 report and as one clause in DoD 5 — by: implementer.
- D13 — 2026-10-07 — Non-blocking findings:
  - The E2E criterion is checked as "the run's `bugs/` holds no record of a defect that its own reviews' `## Outcome` lists as fixed before commit", so it cannot pass vacuously.
  - n=1 is stated with the numbers.
  - The A/B kit stays in `.agents/scratch/ab-kit/` (see D14).
  - The Antigravity run of the changed skills is owed before the next release (AGENTS.md § Release) and goes in the handoff's Next Steps — by: implementer.
- D14 — 2026-10-07 — supersedes D13's kit clause — The A/B kit is not committed. It never sourced `kit-preamble.sh`, so the harness arms' `sh-agent` calls registered the scratch app in the real global search index (BUG-162). It is rebuilt on the preamble, with fixtures under a non-checkout parent, as a follow-up in the handoff's Next Steps — by: implementer.
- Rejected: a numeric context threshold (150k), because the agent cannot observe it.
- Rejected: removing the verifier entirely, because an independent run is the only check that the tree, not the author's shell, passes.
