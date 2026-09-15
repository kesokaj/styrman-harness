---
owner: main
date: 2026-09-15
---
# Harness practices Phase 5: gates, release gate, init and names (G20, G21, G22, G25, G26)     (2026-09-15)

## Summary

Phase 5 of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`: the G26 precedence test on both CLIs, then edits I-19a–e with outcome C, I-20a–e with ADR-022, I-21d–g, I-25a–c, I-22a–c, I-22e–f and I-23a–b, applied verbatim from the plan's old/new blocks. The context no longer promises that the project layer wins a name clash. Every fix under the `bug` skill gets the reviewer. Preflight's docs lens uses the `/review` brief verbatim, routes its fixes and defines what blocks a release. The documentation agent's trigger and lane cover deploy paths and the project `AGENTS.md` setup and deploy sections. The feature, bug and review skills state their gates the same way. `/init-project` derives the build, run and test commands and never overwrites an existing `CLAUDE.md`.

## Type

task

## Files changed

- `context/AGENTS.md` — I-19 with outcome C (D41): "A project skill or agent may not load on every vendor, and which copy wins a name clash is not guaranteed: in a project, a skill or agent never reuses a harness name — extend one with a differently named skill or a `.agents/rules/` file", replacing "Both load; the project layer has priority on both vendors"; `/skills` and `/agents` "show what a vendor loaded" rather than listing both layers. I-20a: Definition of Done step 2 dispatches `reviewer` for 50+ lines, 5+ files, "or any fix under the `bug` skill". The rendered context stays at 93 lines.
- `skills/retro/SKILL.md` and `skills/distill/SKILL.md` — I-19: a proposed project skill or agent is never named after a harness one (`sh-agent list`).
- `agents/reviewer.md` — I-20c: the description dispatches it for a diff of 50 lines or 5 files "or fixes a bug".
- `skills/bug/SKILL.md` — I-20: the reviewer runs whatever the fix's size, since Step 0 already exempted the trivial and a fix under this skill is never a docs-only change. I-20e: `/review --docs` over the fix when it changed an interface, command or deploy step (sweep F29). I-22: a `## Gates that need the user` section with two gates, a missing repro and a review loop at its three-round cap (sweep F28, D40); a gate holds under any "work autonomously" instruction, and only the user's own words in the same turn, naming the gate, waive it; the options come with a recommendation first; when no reply can arrive, a gate that is reached ends the work and its question goes into the handoff's Next Steps. Step 1 points at the gate, and a repro is inferred only when the user picks that option, printed before Step 2 runs.
- `skills/preflight/SKILL.md` — I-20: Step 1 also asks for the reviewer when the release diff is a bug fix, which is a `bugs/BUG-*.md` record added or the `bug` skill run on it this session. I-21d–g: Step 4 runs the `--docs` lens with the brief exactly as `/review` Step 3 gives it, sends findings in the `documentation` agent's files to it, fixes the rest itself and records each with its fix, and those fixes ride in the release commit (sweep F77); Step 5 defines blocking as a verifier FAIL or BLOCKED, a high-severity security finding, a reviewer REQUEST CHANGES verdict or high-severity finding, a `REQUIREMENTS.md` violation, or a Step 4 docs finding still unfixed, and a project's `AGENTS.md` may add conditions but not remove these; Step 6 commits Step 4's docs fixes.
- `agents/documentation.md` — I-25a–c: the description names the deploy guides and the project `AGENTS.md` setup and deploy sections, and it is dispatched when behavior, interfaces, the plan, or a deploy path changed; the trigger paragraph says the same; the lane adds `AGENTS.md` at the root, "the project's Setup commands and Deploy sections only" (sweep F03, F04).
- `skills/feature/SKILL.md` — I-22a: the gates paragraph uses the same waiver wording; at a cap the options are to clarify or cut a requirement, to proceed with the residue as `Open:` rows while the phases that depend on them wait, or to abandon; a gate reached with no possible reply ends the work and writes its question into the handoff's Next Steps.
- `skills/review/SKILL.md` — I-22: the three-round cap is a gate with the same waiver rule; at the cap the options are to accept the remaining findings as `left` with an owner, run one more round, or abandon the fix; a headless `--fix` run applies up to the cap and, when no reply can arrive, stops and writes the question into its record and the handoff.
- `skills/init-project/SKILL.md` — I-23a–b: when `go.mod`, a `web/` or `ui/` `package.json`, or a Dockerfile exists, it derives the build, run and test commands, prints them and asks to confirm or correct, and otherwise asks; it deletes every template command line the user did not confirm; an existing `CLAUDE.md` with content is never overwritten, and the user chooses between folding it into `AGENTS.md` and keeping it with `@AGENTS.md` added as its first line.
- `.agents/docs/architecture/DECISIONS.md` — ADR-022, "Every bug fix gets the reviewer, whatever its size" (task 5.3).
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — boxes 5.1–5.7 and the Phase 5 checkpoint ticked with evidence; D41 records outcome C; 5.1's expect now names outcome C, which D36 added.
- `.agents/docs/testing/2026-09-15-verify-g26-precedence.md` — the precedence test, one line per CLI and layer, answers verbatim.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-phase5-sandbox.md` — the checkpoint's sandbox start of each CLI.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-phase5.md` — the Phase 5 verifier log.

## Agents dispatched

changelog (ADR-022); verifier, dispatched with the repo's What to run rules (plan D9). This entry was written by the main agent from the phase diff.

Review: no `reviewer` ran over the Phase 5 diff. The plan's Done list runs `reviewer` and a `--docs` review once over `c192622..HEAD`.

## Verification

The verifier log reads Status: PASS, with `counts: pass 31, fail 0, skip 1 at b43dde1+dirty for `go test -race -v ./` and `vs previous: `.agents/docs/testing/2026-09-15-verify-harness-practices-phase4.md`, pass 31→31, skip 1→1 (no fall in pass, no rise in skip)`. The main agent also ran, each bare with its exit status read: `go test -race ./...` exit 0; `gofmt -l .` printed nothing; `go vet ./...` exit 0; a scratch install into `.agents/scratch/phase5-checkpoint-home` exit 0, then `status` with drift none; `render` 93 lines; the new reviewer and documentation descriptions in both vendors' installed copies. In the sandbox, this tree installed with drift none, and one start of each CLI quoted Definition of Done step 2 with "or any fix under the `bug` skill"; the newest Antigravity log counts 0 "Failed to parse" after I-20c changed the reviewer's frontmatter. The precedence test ran before I-19a and gave outcome C. The checks of boxes 5.2–5.7 met their expects, and every pattern printed its before value first.

Not run: the skill-matrix rows, which are Phase 6.

## Docs

The changed living docs are the context, agent and skill texts above. `README.md`, the testing guide and the deploy guides are unchanged in this phase.

## Requirements

n/a
