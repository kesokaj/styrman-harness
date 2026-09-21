---
owner: reviewer
date: 2026-09-21
---
# Cold-reader review — project-agents-md-native plan, round 2   (2026-09-21)

Range reviewed: `.agents/docs/plans/2026-09-21-project-agents-md-native.md` and the surfaces it names (`skills/init-project/SKILL.md`, `REQUIREMENTS.md`, `context/AGENTS.md`, `AGENTS.md`, `README.md`, `.agents/docs/architecture/DECISIONS.md` ADR-005, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/testing/2026-09-21-verify-agents-md-native.md`, `.agents/docs/bugs/INDEX.md` BUG-025/BUG-087, `internal/harness/size_test.go`, `internal/harness/plan.go`, `internal/harness/render.go`, `agents/changelog.md`, `ROADMAP.md`, `CLAUDE.md` at repo root). Round-1 report read for context only; judged fresh.

## Status: COMPLETE

## Verdict: implementation not ready

Round 1's two blocking gaps (R2 and skill-matrix.md left out of scope) are both fixed: `REQUIREMENTS.md` R2 is now in Surfaces touched, has exact approved replacement text (D4), and a Phase 1 task; the `init-deploy` matrix row is likewise in scope with a task to drop `CLAUDE.md` from the prompt and add `! test -e CLAUDE.md` to the check, and it checks out against the file on disk. One new blocking gap replaces them: the plan assigns the `REQUIREMENTS.md` edit to the `changelog` subagent, but both the file's own banner and that subagent's documented lane say agents don't touch it. A second, lower-severity defect: the Phase 1 verify clause for the `context/AGENTS.md` line cites a byte figure that is already stale on inspection (measured live below).

## Findings

- **[high / confidence: high]** Plan line 56 (Phase 1): "`REQUIREMENTS.md` R2 (Interfaces text) and ADR-036 by `changelog`, read against the diff, one commit." `REQUIREMENTS.md:3` is unambiguous: "> Immutable without the operator's explicit approval. **Agents never edit this file.**" `skills/init-project/SKILL.md:49` repeats it for the record: "`REQUIREMENTS.md` — leave the tables for the user; note that agents never edit it." `agents/changelog.md`'s own "Your lane" section lists exactly `changes/` and `architecture/DECISIONS.md`, and explicitly disclaims other doc surfaces ("You do NOT write specs, README, ROADMAP, or testing guides"); `REQUIREMENTS.md` never appears in its lane at all. The plan's own Decisions section (D4) records the operator's approval of the exact R2 replacement text, satisfying the global "propose the exact edit → approval → ADR" sequence — but it never says who performs the mechanical edit, and Phase 1 names the one agent whose own charter and the target file's own banner both say no. Fix: either the main agent makes the R2 edit itself (quoting D4's approval) and dispatches `changelog` only for ADR-036, or the plan states explicitly why this counts as within `changelog`'s remit despite both texts — either way this needs a decision before Phase 1 runs, not an implementer's guess.
- **[medium / confidence: high]** Plan line 54's verify clause: "`context/AGENTS.md` Project layout line (Interfaces) — verify: size test bare, AI — expect: PASS, **≤ 20,205 bytes**." Measured live: `go test ./internal/harness/ -run TestRenderedContextFitsAntigravityBlock -v` reports the *current, unedited* checkout at exactly 20,205 bytes (context 20,205, rules 0), against a budget of 22,000 (`size_test.go:19`, `renderedContextBudgetBytes = 22000`). The plan's own proposed replacement text for that line is measured at +104 bytes over today's parenthetical just for the fragment quoted in Interfaces (line 32: "shorter than today's parenthetical or equal" is asserted but false by inspection — the new clause adds a "before 2.1.277, Bedrock, telemetry off" enumeration the old one-clause parenthetical never had). The edit will push the total past 20,205 by design; citing "≤ 20,205 bytes" as the expected pass condition, immediately after the byte-count that is only true of the *unedited* file, risks an implementer reading a >20,205 (but <22,000, genuinely passing) result as a plan failure. There is headroom (1,795 bytes) and the real gate is 22,000 — the clause should say that, not the pre-edit baseline.
- **[low / confidence: medium]** D1 ("An update run leaves an existing project `CLAUDE.md` untouched") sits next to the new `SKILL.md` bullet text (Interfaces, line 31) which keeps the existing fold behaviour verbatim: "one that holds content other than `@AGENTS.md` is offered the fold as before" — and "as before" (`SKILL.md:48` today) means on "yes" the file *is* rewritten down to the single stub line. D1's "untouched" is accurate only as "not proactively deleted," which is the Out-of-scope line's actual claim ("removing stubs from existing projects" stays out); as written next to the fold-preserving bullet, D1 reads as a stronger, blanket claim than the behaviour it sits beside. Worth a one-clause tightening so an implementer doesn't wonder whether the fold offer survives this plan.
- **[low / confidence: low]** Plan Interfaces line 31 says "Step 2 loop becomes `for f in AGENTS.md REQUIREMENTS.md TODO.md ROADMAP.md`." The loop that lists `CLAUDE.md` today lives in `## Step 3: Root files` (`skills/init-project/SKILL.md:37-45`), not Step 2 (`## Step 2: Locate the templates`, which has no such loop). The plan's own "Surfaces touched" section cites the correct line numbers (`:42-44`), so an implementer going by line number is fine, but the prose mislabels the step — worth a one-word fix so a reader following the prose alone doesn't search Step 2 for it.

## Scope

Unchanged from round 1's finding: everything built (SKILL.md edits, this repo's `CLAUDE.md` removal, the `context/AGENTS.md` Project layout line, ADR-036, the `skill-matrix.md` row fix, v1.6.0 MINOR) traces to a D-row citing the operator's "Your recommendation. Approved." (D1, D2, D3, D4) or the implementer's own measurement (D5, now backed by a real testing log — `.agents/docs/testing/2026-09-21-verify-agents-md-native.md` exists on disk with five probes matching D5's summary, not chat memory as round 1 flagged). No unapproved scope creep found.

## Readiness

Could a first-pass implementer ship without asking? No — one BLOCKING question, carried from the new finding above:

- **BLOCKING**: Who makes the `REQUIREMENTS.md` R2 edit — the main agent (quoting D4's approval directly), or `changelog`? The plan names `changelog`, which conflicts with `REQUIREMENTS.md`'s own "Agents never edit this file" banner and with `changelog`'s documented lane (`changes/`, `architecture/DECISIONS.md` only). This needs a decision before Phase 1's `REQUIREMENTS.md`/ADR-036 task runs.
- **NON-BLOCKING**: Should the Phase 1 verify clause for the `context/AGENTS.md` line cite the 22,000-byte budget instead of the 20,205-byte pre-edit baseline, so a passing (but >20,205) result isn't misread as a failure?
- **NON-BLOCKING**: Does D1's "leaves an existing project `CLAUDE.md` untouched" need a qualifier ("not deleted; the fold offer on non-stub content still applies") to match the SKILL.md bullet it sits beside?

## Criteria (six-point lens, applied to the plan document)

1. Dead code — the plan correctly identifies `skills/init-project/assets/CLAUDE.md` for `git rm` and the description/loop/bullet lines that become dead text if left.
2. Flow tracing — traced correctly through `SKILL.md`'s description, update rule, Step 3 loop and bullet, `REQUIREMENTS.md` R2, `context/AGENTS.md`'s Project layout line, and the `skill-matrix.md` fixture; missed tracing the `changelog` agent's own lane against the task it's assigned (finding 1) and the live byte count against its own verify-clause number (finding 2).
3. Standards conformance — the plan itself is careful about the "every restatement of a changed rule" rule (R2, skill-matrix, ADR-005 superseding) but the *task assignment* violates the project's own "agents never edit REQUIREMENTS.md" standard as currently drafted.
4. Directory layout — n/a; no new files proposed outside existing conventions (ADR under `architecture/DECISIONS.md`, testing log under `testing/`, review under `reviews/`).
5. Simplicity — appropriately minimal; correctly rejects a `--claude-stub` flag as unneeded code (Rejected line).
6. Decoupling — n/a; single-skill, single-doc change with clean boundaries, aside from the agent-lane overlap in finding 1.

## Not examined

`internal/harness/harness_test.go`, `pathguard_test.go`, `TestHarnessMap` (global stub and map are Out of scope per the plan, confirmed untouched by `plan.go`/`render.go` inspection but not read end to end). `.agents/docs/testing/skill-matrix.md` rows other than `init-deploy` and `init-update`. `agents/verifier.md`/`agents/reviewer.md` definitions (not named by the plan). Did not run `/init-project`, `sh-agent install`, or any command that writes to the tree; the one test run (`go test -run TestRenderedContextFitsAntigravityBlock`) is read-only and was run to check a factual claim in the plan, not to exercise the change.
