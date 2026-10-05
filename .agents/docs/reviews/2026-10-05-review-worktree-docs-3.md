---
owner: reviewer
date: 2026-10-05
---
# Review — worktree-docs-3   (2026-10-05)

Range reviewed: working tree vs HEAD 507244a: `agents/dreamer.md:32`, `skills/dream/SKILL.md:30,39`; untracked `.agents/docs/testing/2026-10-05-verify-dreamer-global-rules.md` (read whole). Docs lens.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The code-side claims are accurate. `internal/harness/plan.go:62,179-184` deploys `rules/<name>.md` to `~/.claude/rules/` and `~/.gemini/config/rules/`, so the new read paths exist. One living doc still restates the dreamer's "what already exists" read list without the global rule files.

## Findings
- **[low / confidence: high]** narrowed: `agents/dreamer.md:32` → `.agents/docs/specs/sh-agent/memory-loop.md:23`. The spec's **Reads** row lists what the dreamer and dream read as "the global context, the project `AGENTS.md` and `.agents/rules/`, the installed skills' descriptions". The scope adds the global rule files (`~/.claude/rules/*.md`, `~/.gemini/config/rules/*.md`) to both the agent's read list and the skill's Step 2 and `graduated` test. The spec was not updated, so it states the old read set. Fix: add "the global rule files" after "the global context" on that line. The `graduated` clause is not restated in the spec, so only the Reads line is affected.

## Checked, no finding
- Restatement grep and `git grep` over README, AGENTS.md, context/, rules/, agents/, skills/ (vendored catalog excluded), specs, testing/guide.md, deploy/ and `*_test.go`: no other living file lists the dreamer's or `/dream`'s read set.
- `skills/retro/SKILL.md:55` and `skills/feedback-log/SKILL.md:58` already define `rule:<name>` as covering the harness's global `rules/<name>.md`, so they agree with the change.
- `README.md:240` and `plan.go:62` agree with the `~/.gemini/config/rules/*.md` path. Antigravity rule files sit beside `AGENTS.md`, which matches "beside that `AGENTS.md`" in the new text.
- The verification log's claims (3 changed lines, Antigravity not exercised, probe does not discriminate on Claude Code) match the diff. It names what was not covered.

## Unverified
- `.agents/docs/specs/sh-agent/dream.md:13-18` (the nightly job's collected inputs) lists no global context or rule files. I did not check `internal/dream/collect.go` to see whether the headless job reads them. The scope does not touch that path, so no finding.
- The verification log's run results (2 of 2 per definition) are not re-run here.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Code under `internal/dream/` (collect.go and the prompt it builds), `commands/` (empty by design), the `TestAgent*` and size-pin tests, and a live `/dream` or Antigravity run.

## Outcome

- [low] narrowed: `agents/dreamer.md:32` → `.agents/docs/specs/sh-agent/memory-loop.md:23` — fixed by `documentation`: the Reads row lists the global rule files.
