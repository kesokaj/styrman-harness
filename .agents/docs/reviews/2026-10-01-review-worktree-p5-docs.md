---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p5-docs   (2026-10-01)

Range reviewed: working tree on 3a1a861 — changed files: context/AGENTS.md, agents/documentation.md, skills/VENDORED.md, skills/bug/SKILL.md, skills/distill/SKILL.md, skills/dream/SKILL.md, skills/feature/SKILL.md, skills/feedback-log/SKILL.md, skills/handoff/SKILL.md, skills/init-project/SKILL.md, skills/preflight/SKILL.md, skills/research/SKILL.md, skills/retro/SKILL.md, skills/review-standards/SKILL.md, skills/review/SKILL.md, skills/review/scripts/restatement-grep.sh, .agents/docs/bugs/INDEX.md, .agents/docs/plans/2026-09-30-audit-fix.md, .agents/docs/testing/skill-matrix.md; untracked, not ignored: internal/harness/restatement_script_test.go, .agents/docs/bugs/BUG-138-restatement-grep-silent-outside-harness.md, .agents/docs/testing/2026-10-01-verify-phase5-skill-matrix.md (read with `git diff HEAD -- <files>`; untracked files read whole). Docs lens.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No stale restatement, stale test or dropped rule found. The living docs (README, AGENTS.md, context/, agents/, specs, testing guide, deploy guides, other skills) agree with every changed claim I could grep. All findings are low: one unresolvable pointer inside the new `--resume` brief, two spots where the new `--resume` lens or the fix-loop cap text is less explicit than what it replaced, one untested changed branch of the script, and two plan clauses that do not match what was built. `go test ./internal/harness` is green on the working tree (size pins and the new script test included).

## Findings
- **[low / confidence: medium]** `skills/review/SKILL.md:83` — the Resume brief's question 4 says "the awk program in the `handoff` skill § Cold-reader check". The recipient is `explorer`, a read-only subagent given only the brief; every other brief in this file that points at a skill file gives the lookup paths (`~/.claude/skills/...`, `~/.gemini/config/skills/...`, `.agents/skills/...`; docs brief item 3, prompts brief item 2). Here there is no path, so the agent may not find the program when `sh-agent search --open` refuses. Fix: put the awk program in the brief, or give the same three-path lookup.
- **[low / confidence: high]** `skills/review/SKILL.md:27` (Step 1) — `--resume` "takes no scope token and ignores `--fix`", but it is not said what a scope token does (error, as for a bad token, or ignored), nor what `--resume` combined with another lens does. Step 4 (line ~92) then branches on "one lens" (artifact is the record) versus several (a merged record), and that branch is undefined for `--resume --code`. Fix: one sentence, e.g. "`--resume` with a scope token or another lens is an error".
- **[low / confidence: medium]** narrowed: old `skills/review/SKILL.md` Step 5 (HEAD, last sentences) "A headless run with `--fix` applies up to the cap; when no reply can arrive there ... it stops and writes the question" → new `skills/review/SKILL.md:112` ("The cap" paragraph): "When no reply can arrive ... stop there and write the question". The explicit statement that a headless run keeps applying rounds up to the cap, rather than stopping at round 1, is now only implied. Fix: restore "A headless run with `--fix` applies up to the cap;" at the start of that sentence.
- **[low / confidence: high]** `internal/harness/restatement_script_test.go` (whole file) vs `skills/review/scripts/restatement-grep.sh:57` — the `-H` change was made on both branches, but both subtests run the default (any-word) branch only. The `--all` branch strips `${line#*:*:}` assuming `path:line:text`, which is exactly what `-H` guarantees for a one-file result; nothing runs `--all`. Changed code without a covering test. Fix: one subtest in the no-VENDORED.md project calling `run(dir, "--all", "--", "greets", "name")` and asserting the `README.md:1:` hit.
- **[low / confidence: high]** `.agents/docs/plans/2026-09-30-audit-fix.md:140` — the Phase 5 box names the doc obligation "README skills table (the review row names every lens)". README has no skills table and does not enumerate lenses (it says only "`review` is the review command on both CLIs", README.md:263), so nothing in README is stale, but the clause cannot be met as written. Fix: strike the README half of the clause when the box is ticked, or add the lens list deliberately.
- **[low / confidence: high]** `.agents/docs/plans/2026-09-30-audit-fix.md:48` — the plan text says distill's `--from-log` is "implemented as Step 1c" and "Step 6.3 deleted"; `skills/distill/SKILL.md` has no Step 1c: the paragraph sits unnumbered under Step 1 (line 40) and the old 6.3 is replaced by a one-line prohibition (line 158). The edited plan line only updated the `--resume` clause. Fix: reword to "a paragraph under Step 1".

## Checked, no finding (per assignment items 3 and 4)
- `--resume` / `--all`: `skills/review/SKILL.md` description, argument-hint, usage, Step 1, table row and brief agree; `--all` (every lens but `--resume`) is restated nowhere else in the living docs (restatement-grep over `--all`, `every lens`, `--resume`: remaining hits are plans, `dream --all`, the script's own `--all`). `explorer`'s template names `plans/<date>-recon-<slug>.md`, so slug `project-resume` matches; the global context's resume line already treats `recon-*` as briefs.
- Review Step 3 docs and prompts briefs, Step 5, preflight Step 1 / Step 2 / Step 5 / item 6, handoff cold-reader check, init-project 5b: each compared clause by clause with `git show HEAD:<file>`; every old condition has a home (docs brief: vendored exclusion, test-suite paths, `stale test`/`dropped`/`narrowed` shapes, "records never searched"; prompts brief: the none-found stop, provenance, Group 4 roster; preflight: staging list, secret scan, tree-frozen rule; init-project: symlink and tracked-file trust checks, JSON validation and revert, migration conditions, headless guard). Only the headless-cap wording above is weaker.
- Preflight `## Status: COMPLETE` rows match the reviewer and security templates ("## Status: IN PROGRESS | COMPLETE", agents/reviewer.md:67, agents/security.md:54); `grep -c 'Status: COMPLETE' skills/preflight/SKILL.md` = 2 (the plan's expectation). Step 1's "documentation ... has no shell" matches `agents/documentation.md` tools (no Bash). No other living file restates the preflight blocking list.
- distill: `--from-log` and `--global` appear in the argument-hint and Usage; `feedback-log` (lines 14, 164 of distill) and `signal: gap` (feedback-log:50, 79) agree; no living doc still says a distill appends to `feedback.jsonl`; reuse of `retros/<today>-dreamer.md` agrees with retro Step 1 source 8 and memory-loop.md:26.
- dream's Antigravity row agrees with memory-loop.md:14 and :33 (one store at `<root>/.agents/memory`). feedback-log rule 1 agrees with rule 7 (the memory refusal) and handoff:36/:105 (log timestamps quoted). research's `specs/<app>/<topic>.md` and `--kind spec` agree with `internal/search/parse.go:250` and `specs/sh-agent/search.md:26`; `integration-spec` survives only in agents/documentation.md (updated). VENDORED notes: find-skills no longer names `skill-repair` (no hit), grilling's note matches prompts brief item 4 and `skills/grilling/SKILL.md:26`. Context resume line names `project-memory`; `TestRenderedContext*` size pins pass.
- BUG-138 record and INDEX row agree with the script (fix at lines 43, 57, 62) and with the test; `skills/review/scripts/restatement-grep.sh` mode is 100644 (BUG-137 rule kept). Plan boxes WS8/WS13 grep (`styrman-harness|BUG-1xx` in the eight skills) returns 0.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- Antigravity behaviour of every changed skill (Phase 9; the verify log names it as not exercised).
- Unverified: that the review `--docs` lens now works end to end with the fixed script. The `review-docs` matrix row ran against the copy that still had BUG-138 (the verify log says so); only the new Go test exercises the fix.
- Unverified: init-project 5b's `printf` fallback (no machine here lacks `python3`; the log says the same) and the prompts brief's `.agents/skills/review-standards` fallback.
- Not read in full: the Phase 5 hunks of `.agents/docs/plans/2026-09-30-audit-fix.md` beyond line 48 and the Phase 5 boxes at lines 138-150; `skills/review-standards/SKILL.md` and `skills/feature/SKILL.md` diffs were read but their restatements were checked only by grep on the key phrases (`review-standards` loaded-by wording, the verifier/reviewer budget clause). There is no Phase 5 `changes/` entry yet (newest is `2026-10-01-0707-audit-fix-phase-4.md`); expected while the Phase 5 checkpoint is unticked.

## Outcome
- [low] Resume brief's awk pointer has no lookup path — fixed at `skills/review/SKILL.md:83` (three-path lookup, as the other briefs)
- [low] `--resume` with a scope token or another lens undefined — fixed at `skills/review/SKILL.md:27` (runs alone; any other token is the error "--resume runs alone"; its artifact is Step 4's one-lens record)
- [low] narrowed: headless `--fix` applies up to the cap — fixed at `skills/review/SKILL.md:112` (sentence restored)
- [low] `--all` branch of restatement-grep untested — fixed at `internal/harness/restatement_script_test.go:66-74` (same fix as the code lens's round 1)
- [low] plan line 140 names a README skills table that does not exist — fixed at `.agents/docs/plans/2026-09-30-audit-fix.md:140` (docs clause now names `skills/review/SKILL.md` § Usage and says README has no skills table)
- [low] plan line 48 "Step 1c" / "6.3 deleted" — fixed at `.agents/docs/plans/2026-09-30-audit-fix.md:48` (a paragraph under Step 1; 6.3 replaced by a one-line prohibition)
