---
owner: main
date: 2026-10-05
---
# Handoff — 2026-10-05 — critical-thinking rule and operator approval protocol

Session ID: 74aca75d-0194-4616-a64e-486c3ded9a6c
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

The operator asked for two global-context improvements: critical thinking (anti-sycophancy) and how the agent asks them questions (short explanation and recommendation, a numbered table highest priority first, no work until approved by number). Both landed in `da620c3`, approved in their own turn ("Go all approved with your recommendations"). Not released.

## Completed

- `rules/critical-thinking.md` — the harness's first always-on rule file (both CLIs).
- `context/AGENTS.md:5` — Operator paragraph carries the approval protocol; "disagree before complying" moved from Hard rules into the rule file; SHORT LIST struck.
- Restatements aligned: `skills/{feature,dream,requirements}/SKILL.md`, `README.md:240,278`, `AGENTS.md` Layout line.
- ADR-050; verifier PASS; docs review APPROVE WITH NITS with `## Outcome`; behavioural log (FAIL, see below); changes entry.

## Open / In Progress

- **Antigravity P2 (behavioural log row, Status FAIL):** in the one valid run, Antigravity under the new rule refused a direct instruction ("handle the error from os.ReadFile") and stress-tested it instead; the control home carried it out 3 of 3. Claude Code ran it 2 of 2. Likely cause: the rule's "Stress-test before agreeing … For a plan, decision, claim or correction" read as covering an instruction.
- The anti-sycophancy rule showed no measurable difference from the control on Claude Code (the model already pushes back without praise or apology on these probes); its effect is unproven, not disproven.
- Antigravity headless runs under named grants were auto-denied on compound shell lines; a `--dangerously-skip-permissions` re-run was refused by this session's permission classifier. P4 has no valid Antigravity run.
- Requirements edits awaiting approval: none.

## Key Decisions

- ADR-050: the anti-sycophancy text lives in `rules/` (Antigravity block 16,815/22,000), not the context (14,700/15,000, goal 14,500); the approval protocol lives in the context.
- Numbered TABLE kept over a list (operator: "Keep the numbered table, I like the table").
- A direct instruction is its own approval; unattended runs and subagents follow their dispatch (so `claude -p`, `/goal`, the nightly dreamer never wait on "go").

## Corrections

- None logged (`~/.config/sh-agent/feedback.jsonl` does not exist on this machine).

## Agent notes (unverified)

- This box has ONE CPU: `cmd/sh-agent` takes 635-790 s under `-race`, past Go's 10-minute default; the verifier ran with `-timeout 30m`. The project's setup line in `AGENTS.md` has no timeout, so it fails here as written.
- Antigravity sandbox runs read outside their clone through `command(git)` (incl. `git diff --staged` in this checkout); read-only, already accepted by guide § 4.
- `rule:<name>` in the feedback-log/retro vocabulary names a project `.agents/rules/` file only; a correction about the global `rules/critical-thinking.md` has no exact ref.

## Environment State

- Branch: main, clean, pushed (`da620c3`). No running processes.
- Scratch (gitignored): kit `.agents/scratch/ct/kit/` (run1.sh, check.sh), homes `ct-new`, `ct-ctl` (credentials deleted); fixture and clones under `~/sandbox/ct/`.

## Artifacts to Read

- `.agents/docs/testing/2026-10-05-verify-critical-thinking.md` — behavioural rows, the P2 finding, what was not exercised.
- `.agents/docs/architecture/DECISIONS.md` ADR-050 — the decision and its open consequence.
- `.agents/docs/reviews/2026-10-05-review-worktree-docs.md` — docs review and Outcome.
- `.agents/docs/changes/2026-10-05-0932-critical-thinking-and-approval-protocol.md` — the entry.

## Next Steps

Put to the operator as one table, each with the recommendation:

1. Antigravity P2 — recommend: add to `rules/critical-thinking.md` "A direct instruction is carried out; its risks go in the reply's table, after the work", then re-run P2 on Antigravity (needs `command(*)` or the skip flag, operator-approved, since named grants are auto-denied).
2. Anti-sycophancy probes that discriminate — recommend: defer until a probe is found where the control flatters or apologises.
3. Test timeout on one-CPU boxes — recommend: add `-timeout 30m` to the setup line in `AGENTS.md` and guide § 1.
4. `rule:<name>` vocabulary for a global rule file — recommend: `rule:<name>` covers both, with a one-line note in feedback-log and retro.
5. Real-home install (`sh-agent install`) to pick up the rule — operator's step, after item 1.
