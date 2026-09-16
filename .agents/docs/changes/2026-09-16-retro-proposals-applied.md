---
owner: changelog
date: 2026-09-16
---
# Retro proposals applied (2026-09-16)

## Summary
Eight harness-scope proposals from `.agents/docs/retros/2026-09-16-retro.md` queued in `~/.config/sh-agent/proposals/2026-09-16-retro.md` were applied on the operator's instruction ("go all, apply the retro proposals"). Five proposals bring standing rules into skills and agents; three graduate bug-prevention rules from the index. Proposals 1, 2, 3, 5, 6, 7 touch context/AGENTS.md or skills prose to address rule restatements scattered across the harness; Proposal 4 graduates BUG-029 to a Hard-rules sentence; Proposal 8 (bugs agent) moved BUG-040 and BUG-043 to the Graduated table. Reviewed, verified, and installed; owed: one sandbox start of each CLI with parse-error grep (recorded in the handoff).

## Type
task

## Files changed
- `context/AGENTS.md` — step 2 contradiction clause (two dispatched agents contradicting means fix the rule), step 4 EMPTY at every checkpoint, Hard-rules filter-that-fires sentence
- `agents/changelog.md` — commits clause (run `git show --stat` on named SHAs, say what the commit IS)
- `agents/reviewer.md` — criterion 3 restatement of the changed rule on two-agent contradictions
- `skills/review-standards/SKILL.md` — § 3 checklist item for grepping every restatement of a changed rule
- `skills/review/SKILL.md` — Step 5 round loop contradiction clause
- `skills/handoff/SKILL.md` — Step 3 item 3 restatement (both empty at every checkpoint, not just at Done)
- `skills/golang/SKILL.md` — reference-file reading sentence (HTTP service reads golang-io-web.md and golang-testing.md first)
- `skills/init-project/SKILL.md` — Step 2 headless remark (Claude Code may refuse loops with `$var` expansion)
- `.agents/docs/bugs/INDEX.md` — BUG-029, BUG-040, BUG-043 moved to Graduated table
- `.agents/docs/retros/2026-09-16-retro.md` — applied status recorded (after commit)
- `~/.config/sh-agent/proposals/2026-09-16-retro.md` — applied status recorded in each row

## Agents dispatched
reviewer ×2 (code lens and docs lens over the six prose paths), verifier (sandbox install/status/grep, go test), bugs (INDEX.md graduation), changelog (this entry).

## Verification
- `go test -race -v -count=1 ./...` → pass 64, fail 0, skip 1 at 93807be+dirty.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-retro-applied` → exit 0; 260 files written, no errors.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-retro-applied` → drift none, 259 files match.
- `go run ./cmd/sh-agent list` → 27 skills, 7 agents.
- Grep confirmation: new text present in all six changed files and in their rendered copies (`.agents/scratch/verify-retro-applied/.claude/`).
- Owed: one start of each CLI in the sandbox (recorded in the handoff; not run in this dispatch).

## Docs
`context/AGENTS.md`, `agents/changelog.md`, `agents/reviewer.md`, `skills/review/SKILL.md`, `skills/handoff/SKILL.md`, `skills/review-standards/SKILL.md`, `skills/golang/SKILL.md`, `skills/init-project/SKILL.md` all changed (step 2 and step 4 of the Definition of Done; Hard-rules filter; agent rules; skill prose and checklists). `.agents/docs/bugs/INDEX.md` updated (three rows graduated). No project README, specs, testing guides, or deploy guides changed — harness-only, no user-facing interface or deploy step change.

## Requirements
n/a — no interface, config, or deploy path changed in a user-facing direction.
