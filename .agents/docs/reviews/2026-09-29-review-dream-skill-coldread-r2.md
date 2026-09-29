---
owner: reviewer
date: 2026-09-29
---
# Review — dream-skill-coldread-r2   (2026-09-29)

Range reviewed: `.agents/docs/plans/2026-09-29-dream-skill.md` (cold-reader check, round 2, not a diff), against round 1's report (`.agents/docs/reviews/2026-09-29-review-dream-skill-coldread-r1.md`) and the code/docs surfaces the plan names: `REQUIREMENTS.md` R7, `skills/retro/SKILL.md`, `internal/harness/inventory.go`, `.agents/docs/testing/skill-matrix.md`, `AGENTS.md` § Testing, the operator's real Claude Code auto-memory store (`~/.claude/projects/-home-kesokaj-workspace-styrman-harness/memory/`), `claude --help` (2.1.284, this box), and the most recent behavioural-run precedent (`.agents/docs/testing/2026-09-29-prompts-lens-behavioural.md`, `2026-09-29-verify-prompts-lens.md`).

## Status: COMPLETE

## Verdict: design ready / implementation not ready

Round 1's two BLOCKING gaps are resolved: D7 cites a real, verifiable settings key and flag (`claude --help` on this box does list `--settings <file-or-json>`, confirmed by running it), and D8 replaces the impossible "one file, two `originSessionId`s" fixture with a two-file/two-session-id pair, which matches the real store's shape (every file here carries exactly one). D9 gives the matrix row a concrete seeding mechanism. But the plan inherits an unstated, currently-true infrastructure gap: Phase 1's install-tier task requires "one start of each CLI in the sandbox … with its log grepped for parse errors," verbatim from `AGENTS.md` § Testing's unconditional rule for any change to `context/`, `agents/`, or `skills/` — and the sandbox that rule points at, `/home/kesokaj/sandbox/home` (`.agents/docs/testing/2026-09-03-verify-sandbox.md`), does not exist on this machine (`ls` confirms `No such file or directory`), a fact the harness's own most recent behavioural run already recorded in writing the same day ("the sandbox home this repo's procedure names no longer exists on this box," `2026-09-29-prompts-lens-behavioural.md` § Not exercised). The plan does not decide whether Phase 1 rebuilds the sandbox (nontrivial: two CLI logins, `sh-agent install --home`, demo apps) or substitutes real-home smoke for it, so the install-tier checkpoint as literally written cannot be closed by a first-pass implementer without a decision the plan doesn't make.

## (1) Gaps

- **[high / confidence: high] Phase 1 install-tier task, plan line 106 ("one start of each CLI in the sandbox … with its log grepped for parse errors") and the identical clause implicit in Phase 2's matrix-row task (line 112, "the row's check command passes on the Phase 1 fixture clones," which presumes the matrix's own `Procedure` §1 sandbox home) — the sandbox this rule names does not exist on this box.** `ls /home/kesokaj/sandbox` → `No such file or directory`. The most recent behavioural-run record in this same repo, written the same day, says so explicitly: "The sandbox home this repo's procedure names no longer exists on this box; the real-home smoke checks after `sh-agent update` are the operator's, per `AGENTS.md` § Testing" (`2026-09-29-prompts-lens-behavioural.md` § Not exercised). `AGENTS.md` § Testing states the sandbox-start rule unconditionally for any `context/`/`agents/`/`skills/` change, with no fallback clause for the sandbox's absence (the real-home-smoke fallback it does name is scoped to the tier-4 *behavioural* run for `agents/*.md`, a different clause). Two implementers would do different things here: one rebuilds the sandbox from the 2026-09-03 procedure before Phase 1 can close; the other treats the rule as satisfied by real-home smoke alone (already a separate, later task in Phase 2) and skips the sandbox start. The plan's Phase 1 checkpoint ("the skill installed in a scratch home and proven on the fixture") doesn't name which happened, so the checkpoint's evidence file could look identical either way. Fix: one sentence deciding whether Phase 1 rebuilds the sandbox or the plan substitutes/defers that specific sub-clause, with a D-row.
- **[low / confidence: low] Phase 1 behavioural task, plan line 107 — the control run's N is unstated ("N=3 plus control").** A reasonable implementer runs the control once (a binary absence check), but the plan states "3 of 3" and "0 of 3" for the seeded run and doesn't say the control also gets 3 runs; non-blocking since either reading produces a verifiable pass/fail, but the phase-1 checkpoint expectation ("control 0 of 3") already assumes 3, so an implementer who runs the control once would need to reconcile the checkpoint wording against their own choice.

No other gap found: the store-path precedence question I checked (whether `--settings` is read before or after `.claude/settings.local.json`/`.claude/settings.json`/`~/.claude/settings.json` when several set `autoMemoryDirectory`) is real-world ambiguous in the plan's Interfaces text but doesn't affect Phase 1's fixture (a fresh scratch clone has none of the other three files, confirmed no such files exist in this repo's own `.claude/` either), so it is not load-bearing for the verification plan as scoped and I do not raise it as a gap. The `[[links]]` claim in the `merge` verdict row and the `asking-format.md`/context-line graduation claim in the E2E criteria both check out against the real store (`read-neighbours-before-inserting-rules.md` does carry a `[[fallback-paths-stay-simple]]` link; `asking-format.md`'s own "Why" section does say its rule is "now in the harness context (`context/AGENTS.md` Operator line)").

## (2) Scope

Unchanged from round 1: everything in Scope — In traces to `## Why`; Scope — Out items are each grounded in a D-row (D4, D6) or the sibling research doc. Nothing added between round 1 and round 2 (D7–D9) introduces new scope — they are citations and a fixture-mechanism decision, not new behaviour.

## (3) Readiness

1. **BLOCKING** — Does Phase 1 rebuild the sandbox at `/home/kesokaj/sandbox/home` before its install-tier task, or is that specific sub-clause of the required testing tier waived/substituted given its current absence? (Gap 1 above.) Two implementers ship different Phase 1 evidence.
2. NON-BLOCKING — How many runs does the behavioural control get? A reasonable implementer runs it once; the checkpoint wording's "0 of 3" is then read as "0 of however many the control ran."

Round 1's five questions are now closed: Q1 and Q2 (BLOCKING) by D7/D8; Q3 (non-blocking) by D9; Q4 and Q5 (non-blocking) were left as reasonable-default calls, which stands.

## Criteria

1. Dead code — not applicable (plan review, no source in scope).
2. Flow tracing — fail on the one BLOCKING gap: the Phase 1 checkpoint's evidence depends on a testing-tier sub-clause that points at infrastructure confirmed absent on this machine.
3. Standards conformance — pass; D7's citation practice (dated, sourced, cross-checked against `claude --help` on this box) matches the plan's own evidentiary bar; D5's "REQUIREMENTS line 31" citation is accurate (R7, `.agents/docs/retros/`).
4. Directory layout — pass; unchanged from round 1.
5. Simplicity — pass; unchanged from round 1.
6. Decoupling — pass; unchanged from round 1.

## Not examined

- `agents/*.md`, `commands/*.md`, `internal/harness/derive.go` — not touched by this plan (same as round 1).
- Live execution of `claude -p --settings` redirecting the auto-memory store — still not run; the CLI flag's existence is confirmed (`claude --help`) but its effect on `autoMemoryDirectory` is still an assumption the plan itself schedules as Phase 1's first checked task (Risks section), not something I independently verified.
- Whether the sandbox at `/home/kesokaj/sandbox/home` can in fact be rebuilt from the 2026-09-03 procedure on this box today (login state of either CLI, disk space) — I confirmed only that the directory is currently absent, not the cost or feasibility of recreating it.
- Full read of `.agents/docs/testing/skill-matrix.md` beyond the `review-prompts` row and § Procedure (read in full); other rows skimmed, not individually re-verified.

design ready / implementation not ready
