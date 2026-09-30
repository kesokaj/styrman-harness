---
owner: reviewer
date: 2026-09-30
---
# Review — agy-only-init   (2026-09-30)

Range reviewed: `git diff HEAD` (skills/init-project/SKILL.md, skills/dream/SKILL.md, context/AGENTS.md, .agents/docs/specs/sh-agent/memory-loop.md, .agents/docs/testing/skill-matrix.md) plus untracked init_project_test.go. The untracked verify log .agents/docs/testing/2026-09-30-1721-verify-agy-only-init.md is not part of the range and was not read.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

`go test ./internal/harness` FAILS: `TestRenderedContextFitsAntigravityBlock` reports the Antigravity rules block at 22050 bytes against the 22000 budget (REQUIREMENTS "Known limits", BUG-085). The 54-byte addition to `context/AGENTS.md:38` ("unless the harness is installed for Antigravity alone") pushes it over. The source file went from 18492 to 18546 bytes, and the rendered block was 21996 before, so HEAD passed. The rest of the change is coherent.

## Findings
- **[high / confidence: high]** `context/AGENTS.md:38` — the added clause takes the rendered Antigravity context to 22050 bytes; `internal/harness/size_test.go:45` fails ("cut 50 bytes"). Violates the REQUIREMENTS byte-budget limit, so it is blocking. Fix: cut at least 50 bytes elsewhere in the context, or drop or shorten the clause. The rule is fully stated in `skills/init-project/SKILL.md` Step 1, and the global paragraph only needs "linked for Claude Code" or a shorter qualifier. Recheck with `go test ./internal/harness`.
- **[low / confidence: medium]** `.agents/docs/testing/skill-matrix.md:36-37` — the `init-agy-only` and `dream-agy-only` rows are new, but the diff carries no run evidence. Other rows cite a verify log (init-memory-linked-store, for example). The row says "added", not "run", which is honest. The untracked verify log may cover it, so confirm the rows are cited as run only once a PASS log exists. This is not a defect on its own.
- **[low / confidence: low]** `init_project_test.go` — half the test reads skill text for the needle (a source-text assertion). The other half is behavioural: it writes real manifests through `SaveState` and checks the needle against them, and it fails if the manifest format or the skill needle drifts. This is acceptable, since the pin is its whole purpose and the header says so. The test does not cover an old manifest with an empty Vendor (which `LoadState` normalises to `all`, so the grep correctly does not match). Not required.

## Criteria
1. Dead code — pass. Nothing left over.
2. Flow tracing — pass. The needle `"vendor": "antigravity"` matches `json.MarshalIndent` output (`state.go:62`) and `VendorAntigravity` (`plan.go:22`), and the test proves it for all three vendors. Failing or missing manifests keep the `.claude/` writes, which is the safe default. Dream item 0 refers to item 1's conditions and stops when they fail. Step 5b items 2 and 3 are skipped on agy-only, consistent with Step 6's exclusion of the two `.claude/` lines. Step 5b item 1 still adds `.claude/settings.local.json` to `.gitignore` in its text, and the added paragraph says "without the `.claude/settings.local.json` line", so it is consistent.
3. Standards conformance — FAIL: the byte budget (finding 1). The restatements I grepped (memory-loop.md, README:214, context:38, dreamer.md) are consistent with the change. The spec updates both the "Who reads" and "Set up by" rows. `~/.gemini/config/rules/AGENTS.md` is the correct agy path (README:208).
4. Directory layout — pass. A root-level `init_project_test.go` sits beside `embed_test.go` and `plan_test.go` (package `styrman_test`), which is where the other skill-text tests live.
5. Simplicity — pass. One grep, no new dependency, nothing security-related.
6. Decoupling — pass. Both skills share one needle, and the test pins that they agree.

## Not examined
- The untracked verify log `2026-09-30-1721-verify-agy-only-init.md`.
- Behavioural runs of the new matrix rows (Antigravity is not available to me), so the agy-only skill flows themselves were reviewed by reading only.
- Other packages' tests: I ran only `go test . ./internal/harness` and `go vet .`. The root package passed.

## Outcome

1. fixed — `context/AGENTS.md:38` now reads "with `.claude/*` linked into it by `/init-project` (Claude Code only)", and the same line drops "in a project, " (the sentence already opens with "A project skill or agent"). Rendered Antigravity block: 21998 bytes of 22000; `TestRenderedContextFitsAntigravityBlock` fails on the reviewed diff (22050) and passes after the fix — the same test and arrangement the finding measured.
2. no change needed — the rows are marked "added"; both were run on Antigravity in a sandbox home installed `--vendor agy` (pre-fix run fails the `init-agy-only` check, post-fix runs pass both), evidence in `testing/2026-09-30-1721-verify-agy-only-init.md`.
3. no change needed — accepted as the reviewer states.
