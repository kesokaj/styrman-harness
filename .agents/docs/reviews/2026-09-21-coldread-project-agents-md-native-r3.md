---
owner: reviewer
date: 2026-09-21
---
# Cold-reader review — project-agents-md-native plan, round 3   (2026-09-21)

Range reviewed: `.agents/docs/plans/2026-09-21-project-agents-md-native.md` and the surfaces it names (`skills/init-project/SKILL.md`, `REQUIREMENTS.md`, `context/AGENTS.md`, `AGENTS.md`, `README.md`, `.agents/docs/architecture/DECISIONS.md` ADR-005/ADR-035, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/testing/2026-09-21-verify-agents-md-native.md`, `.agents/docs/bugs/INDEX.md` BUG-025/BUG-087, `internal/harness/size_test.go` (run live), `internal/harness/plan.go:85-92`, `internal/harness/render.go:13-19`, `agents/changelog.md`, `skills/init-project/assets/CLAUDE.md`, repo-root `CLAUDE.md`). Round-1 and round-2 reports read for context only; judged fresh against the plan text and the live tree.

## Status: COMPLETE

## Verdict: design ready, implementation ready

Both prior blocking gaps are fixed and verified against the live tree: the `REQUIREMENTS.md` R2 edit is now assigned to the main agent (line 56: "replaced by the main agent with the operator-approved Interfaces text"), matching `REQUIREMENTS.md`'s own "Agents never edit this file" banner and `changelog`'s lane (`changes/`, `architecture/DECISIONS.md` only — confirmed at `agents/changelog.md:26-31`); `changelog` is now scoped to ADR-036 only. The size-test verify clause (line 54) now cites the real 22,000-byte budget instead of the stale 20,205 pre-edit baseline — confirmed live: `go test ./internal/harness/ -run TestRenderedContextFitsAntigravityBlock -v` reports 20,205/22,000 today, matching the plan's own figures. D1 now carries the fold-behaviour qualifier round 2 asked for. The testing log `.agents/docs/testing/2026-09-21-verify-agents-md-native.md` exists on disk with five probes backing D5 and the `Why` claim, not chat memory. Every D-row (D1-D5) traces to an operator quote or a measured probe; no unapproved scope creep. Every citation checked against the live tree (ADR-005's exact superseded sentence, BUG-025/BUG-087, `plan.go:89`, `render.go:17`, the asset and root `CLAUDE.md` contents, ADR-036 being the next unused number) matches.

## Findings (non-blocking, both carried from round 2 and left unfixed)

- **[low / confidence: high]** Plan line 32 (Interfaces): "…before 2.1.277, Bedrock, telemetry off), …" — shorter than today's parenthetical or equal; the size test stays green." Measured: today's parenthetical fragment `` (`CLAUDE.md` is the line `@AGENTS.md`; no `GEMINI.md`)`` is 54 bytes; the plan's own proposed replacement fragment is 158 bytes — +104 bytes, not "shorter … or equal." This is the same false claim round 2 flagged; it was not corrected, only the downstream verify clause (line 54, now correctly citing the 22,000-byte budget) was. An implementer reading line 32 alone still gets a wrong expectation, though the operative gate (line 54) is right and the size test does pass with headroom (1,795 bytes today). Fix: drop or correct "shorter than today's parenthetical or equal" in line 32 to match line 54's "a few bytes longer or shorter."
- **[low / confidence: high]** Plan line 31: "Step 2 loop becomes `for f in AGENTS.md REQUIREMENTS.md TODO.md ROADMAP.md`". The `for f in ... CLAUDE.md` loop is in `## Step 3: Root files` (`skills/init-project/SKILL.md:42`), not Step 2 (`## Step 2: Locate the templates`, which has no such loop, confirmed by reading the live file). Round 2 flagged this exact mislabel; it is unchanged in round 3. Not blocking — "Surfaces touched" (line 21) cites the correct line numbers (`:42-44`), so an implementer working off line numbers lands in the right place, but a reader following the prose alone would search Step 2 and find nothing.

## Scope

Unchanged from rounds 1-2: everything built (SKILL.md edits, this repo's own `CLAUDE.md` removal, the `context/AGENTS.md` Project layout line, `REQUIREMENTS.md` R2, ADR-036, the `skill-matrix.md` row fix, the testing log, v1.6.0 MINOR) traces to a D-row citing the operator's approval (D1-D4) or the implementer's own measurement (D5, now backed by the testing log on disk). No unapproved scope creep found on this fresh read either.

## Readiness

Could a first-pass implementer ship without asking a question? Yes.

- **NON-BLOCKING**: Should plan line 32's "shorter than today's parenthetical or equal" be corrected to match line 54's accurate "a few bytes longer or shorter" (+104 bytes measured), so the Interfaces prose doesn't contradict the Phase 1 verify clause sitting 22 lines below it?
- **NON-BLOCKING**: Should line 31's "Step 2 loop" be corrected to "Step 3 loop" to match the file it describes and the plan's own correct line-number citation at line 21?

Neither withholds a decision implementation needs — both are prose accuracy nits in the plan document itself, not omissions the implementer must resolve before acting, and both were already surfaced (and left non-blocking) in round 2.

## Criteria (six-point lens, applied to the plan document)

1. Dead code — the plan correctly identifies `skills/init-project/assets/CLAUDE.md` for `git rm` (confirmed present on disk, content `@AGENTS.md`) and the description/loop/bullet lines that become dead text if left.
2. Flow tracing — traced correctly through `SKILL.md`'s description, update rule, Step 3 loop/bullet, `REQUIREMENTS.md` R2, `context/AGENTS.md`'s Project layout line, ADR-005's exact sentence, the `skill-matrix.md` fixture, and the agent-lane assignment (now correct); the two byte/step-number inconsistencies above are the only unresolved trace mismatches.
3. Standards conformance — the plan now correctly respects "agents never edit `REQUIREMENTS.md`" (fixed from round 2) and "every restatement of a changed rule is changed with it" (R2, skill-matrix, ADR-005 superseding, all present).
4. Directory layout — n/a; no new files proposed outside existing conventions (ADR under `architecture/DECISIONS.md`, testing log under `testing/`, review under `reviews/`).
5. Simplicity — appropriately minimal; correctly rejects a `--claude-stub` flag (Rejected line) and removing stubs from existing projects on update.
6. Decoupling — n/a; single-skill, single-doc change; the `changelog`/main-agent lane split is now clean (round 2's overlap fixed).

## Not examined

`internal/harness/harness_test.go`, `pathguard_test.go`, `TestHarnessMap` (global stub and map are Out of scope per the plan; `plan.go:85-92` and `render.go:13-19` inspected directly and confirm the global stub write is unconditional and untouched by this plan). `.agents/docs/testing/skill-matrix.md` rows other than `init-deploy`. `agents/verifier.md`/`agents/reviewer.md` definitions (not named by the plan). Did not run `/init-project`, `sh-agent install`, or any command that writes to the tree; the only command run (`go test -run TestRenderedContextFitsAntigravityBlock`) is read-only, used to check a factual claim in the plan against the live tree.
