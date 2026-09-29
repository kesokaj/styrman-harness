---
owner: reviewer
date: 2026-09-29
---
# Review — dream-skill-coldread-r3   (2026-09-29)

Range reviewed: `.agents/docs/plans/2026-09-29-dream-skill.md` (cold-reader check, round 3, not a diff), against round 2's report (`.agents/docs/reviews/2026-09-29-review-dream-skill-coldread-r2.md`) and round 1's, and the code/docs surfaces the plan names: `REQUIREMENTS.md` R7, `skills/retro/SKILL.md` (Steps 1–6 read in full), `skills/feedback-log/SKILL.md`, `skills/distill/SKILL.md` (for the `retros.jsonl` `kind` field precedent), `.agents/docs/testing/skill-matrix.md` (`review-prompts` row and § Procedure), `.agents/docs/testing/2026-09-29-prompts-lens-behavioural.md` § Not exercised, `AGENTS.md` § Testing, the operator's real Claude Code auto-memory store (`~/.claude/projects/-home-kesokaj-workspace-styrman-harness/memory/`, all three memory files read), and version checks of `agy` (1.2.13) and `claude` (2.1.284) on this box.

## Status: COMPLETE

## Verdict: design ready / implementation not ready

D10 and D11 correctly close round 2's blocking and non-blocking questions (the sandbox substitution matches the same-day `--prompts` precedent verbatim, and is independently confirmed by `2026-09-29-prompts-lens-behavioural.md`'s own "Not exercised" note). But a new gap surfaces on this pass: no verification tier the plan schedules — not the behavioural runs, not the real-home smoke, not the must-fail guard — ever exercises the skill's one destructive code path (deleting a memory file, rewriting `MEMORY.md`, merging two files into one). Every scripted run is `--dry-run` or headless (which the plan's own Failure-modes table says never deletes), while E2E criterion (1) explicitly requires the non-dry-run "go" path to run and its effects to be observed. No Phase task's verify/expect clause performs that run. The plan is otherwise sound: the store-path and verdict logic are well-cited, the two round-2 items are closed, and scope stays grounded in `## Why`.

## (1) Gaps

- **[high / confidence: high] No Phase task ever runs `/dream` past `--dry-run`, so the apply path (delete, index rewrite, merge, broken-memory rewrite) is never exercised by anything the plan schedules, despite E2E criterion (1) requiring it.** Verification's behavioural tier (plan line 95) types `"/dream --dry-run"` and states the expectation in the same sentence: "nothing deleted (dry run)." The real-home smoke task (line 97, restated in Phase 2's second-to-last task, line 113) also runs only `/dream --dry-run`. The Failure-modes table (line 87) states as a matter of design that a headless run (`claude -p`, which is what every behavioural tier uses per `AGENTS.md` § Testing) "computed and recorded, nothing deleted; every verdict stays `proposed`" — meaning the headless methodology this plan otherwise reuses is structurally incapable of ever reaching the apply branch, regardless of the `--dry-run` flag. Yet E2E criterion (1) (line 100) is unconditional: "after 'go' the file and its index line are gone and `MEMORY.md` has one line per remaining file" — a claim about real, observable behaviour that no scripted or checklisted task produces. Two implementers diverge here: one closes Phase 2's Checkpoint ("Definition of Done closed") once the scripted dry-run smoke passes, permanently leaving E2E-1 unverified; another insists the operator run `/dream` for real in this repository and type "go" before the checkpoint can close. The plan decides neither. Fix: add an explicit Phase 2 task — e.g. "operator runs `/dream` for real in this repository, approves the table, confirms `asking-format.md` and its `MEMORY.md` line are gone" — as the interactive verification E2E-1 needs, since it cannot be produced by any headless run this plan's own Failure-modes table says never deletes.
- **[low / confidence: low, carried forward from round 1, not newly raised as blocking] The verdict table's `project`/`reference`/`user`/`feedback` type vocabulary (plan lines 46–53) is still uncited.** D7 cites the docs for storage location and load limits but not for the `type` field's enumerated values; every memory file on this box carries `type: feedback` only (`asking-format.md`, `fallback-paths-stay-simple.md`, `read-neighbours-before-inserting-rules.md`, all confirmed by direct read). Round 1 (non-blocking Q5) and round 2 both left this as a reasonable natural-language default rather than a coding interface; I confirm the same reasoning still holds and do not raise it as a new blocking item.
- **[low / confidence: low] The Verification section's Docs tier ("`/review --docs` over the range," line 98) has no matching Phase checkbox.** `AGENTS.md` § Testing already makes `/review --docs` mandatory for any `context/`/`agents/`/`skills/` change independent of what a plan's Phase list states, so this is unlikely to be skipped in practice, but neither Phase 1 nor Phase 2 names it as a checkable task the way the install, behavioural, and real-home tiers each get one. Not blocking on its own.

## (2) Scope

Unchanged from rounds 1–2: every `## Scope — In` item traces to `## Why` (the untended-memory problem, and the specific `asking-format.md` example already graduated into context). `## Scope — Out` items are each grounded in a D-row (D4, D6) or the sibling research doc. Nothing added since round 2 (D10, D11) introduces new behaviour — both are verification-methodology decisions, not new scope. The standalone `/dream` invocation (not only via `/retro`) is named in `## Why` implicitly (memory drifts *between* retros) and explicitly in Scope — In's first bullet ("typed as `/dream`, `--dry-run`"), so it is not unrequested scope.

## (3) Readiness

1. **BLOCKING** — Does Phase 2's Checkpoint require an actual, non-dry-run `/dream` run in this repository (operator-approved "go") before Definition of Done closes, given that E2E criterion (1) needs it and no scripted/headless tier can ever produce it? (Gap 1 above.) Two implementers ship different closing conditions for the same checkpoint.
2. NON-BLOCKING — Exact `type:` vocabulary for memory files (carried forward, reasonable-default per rounds 1–2).
3. NON-BLOCKING — Whether `/review --docs` needs its own Phase checkbox, given `AGENTS.md` § Testing already mandates it unconditionally for this diff's surfaces.

Round 2's items are now closed: Q1 (BLOCKING, sandbox) by D10, independently confirmed against `2026-09-29-prompts-lens-behavioural.md`; Q2 (non-blocking, control run count) by D11.

## Criteria

1. Dead code — not applicable (plan review, no source in scope).
2. Flow tracing — fail on the one BLOCKING gap: the apply/delete/merge/rewrite code path is traced through Interfaces and Failure modes but never reaches a verify/expect clause in any Phase task.
3. Standards conformance — pass; D7's and D9's citation and fixture-seeding practices hold up against direct checks on this box (`agy --version` 1.2.13 and empty `~/.gemini/antigravity-cli/knowledge/` matching the Risks assumption; `claude --version` 2.1.284 matching D7's `--settings` citation; real memory files matching the `originSessionId`-per-file and `modified:` claims).
4. Directory layout — pass; unchanged from prior rounds (`skills/dream/SKILL.md` is a plain file drop per project rule 5).
5. Simplicity — pass; unchanged from prior rounds; no new dependency, no hand-rolled auth/crypto surface.
6. Decoupling — pass; unchanged from prior rounds (dream never writes a rule, per D1; retro stays the single writer of rule destinations).

## Not examined

- `internal/harness/inventory.go:346` (the `## Installed` count generation) — read only by line reference in the plan, not opened; low risk since it is an existing, unmodified mechanism.
- Live execution of `claude -p --settings` redirecting `autoMemoryDirectory`, and of `/dream`/`/retro` themselves — this plan schedules both as Phase 1/2 tasks; neither was run by this review (a plan review, not an implementation).
- Whether the `/home/kesokaj/sandbox/home` sandbox can be rebuilt at all on this box (feasibility/cost) — still unchecked, as round 2 noted; moot given D10's substitution.
- Full text of `context/AGENTS.md` § Learning loop beyond the quoted insertion point, and full read of `skill-matrix.md` beyond the `review-prompts` row and § Procedure (already read in full by round 2).

design needs revision — implementation not ready
