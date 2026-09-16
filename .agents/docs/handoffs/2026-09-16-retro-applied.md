---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — retro proposals applied; nothing queued

Session ID: c50e0fd4-6273-4e35-b1ef-7639dff8bb07        (same machine: `claude --resume c50e0fd4-6273-4e35-b1ef-7639dff8bb07`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

Continuation of `2026-09-16-linkshelf-simulation-done.md` (read it for the simulation itself). Its one open gate was the retro's eight proposals; the operator answered "go all, apply the retro proposals". They are applied, reviewed, verified, installed here and in the sandbox, and pushed at 2b0dd38.

## Completed

- **All eight proposals applied in this repository** (`changes/2026-09-16-retro-proposals-applied.md`): `context/AGENTS.md` step 4 (TODO and ROADMAP EMPTY at every checkpoint), step 2 (two agents contradicting on one mechanism means the rule is the defect), Hard rules (a filter that fires is evidence only after its matching line is read); `agents/changelog.md` (`git show --stat` on every SHA named); `skills/review-standards/SKILL.md` § 3 and `agents/reviewer.md` criterion 3 (grep every restatement of a changed rule); `skills/golang/SKILL.md` (the HTTP and testing reference files named); `skills/init-project/SKILL.md` (Claude Code's headless loop refusal); BUG-029, BUG-040, BUG-043 moved to the new Graduated table of `bugs/INDEX.md`.
- **Two more restatements found by the reviews and fixed in the same commit:** `skills/handoff/SKILL.md` Step 3 ("when the task is done both are empty" → at every checkpoint) and `skills/review/SKILL.md` Step 5 (the contradiction clause in the round loop). Reviews: `reviews/2026-09-16-review-retro-2026-09-16-applied.md` (APPROVE WITH NITS), `reviews/2026-09-16-review-worktree-docs-retro.md` (REQUEST CHANGES, fixed); verifier `testing/2026-09-16-verify-retro-2026-09-16-applied.md` PASS, counts pass 64, fail 0, skip 1.
- **Sandbox:** installed from the working tree, one no-tools start per CLI, both quote the two new context sentences verbatim, Antigravity log `cli-20260916_094552.log` 0 `Failed to parse`; then installed again from 2b0dd38, `drift none`.
- **This box:** installed from 2b0dd38, `binary 2b0dd3859c1c`, `drift none`; reinstalled from the commit carrying this handoff so the `binary` line matches HEAD.
- Queue and index: `~/.config/sh-agent/proposals/2026-09-16-retro.md` all `applied 2026-09-16`; the crew's `2026-09-14-retro.md` P3 and P4 `applied 2026-09-16` (via P1 and P2 here); `retros.jsonl` last line statuses `applied`; `retros/2026-09-16-retro.md` statuses applied.

## Open / In Progress

- Nothing in flight; nothing queued in the harness proposals directory.
- Open records: BUG-019 to BUG-023 (skill sweep findings) and BUG-031 (a print-mode run can end before what it started returns).
- Requirements edits awaiting approval: none.

## Key Decisions

- Applying was the operator's call ("go all"); statuses in the queue and record were written before the review ran — the reviewer flagged it, and the rule for the next retro is status after review (its Outcome).
- Two restatements the proposals' own greps missed were fixed in the applying commit rather than filed: Proposal 5's rule is exactly for them.

## Corrections

- None: no feedback-log entry for this repository today.

## Agent notes (unverified)

- The verifier reported the sandbox home "does not exist on this box yet"; it does (`/home/kesokaj/sandbox/home`) — a cwd or permission blind spot in a dispatched agent, not a missing directory. The sandbox rows were run by the main agent.
- G20 (the verifier's `<app>` unresolved when `cmd/<app>` and the repository differ in name) is filed in the simulation plan's Follow-ups; the verifier definition was edited twice today, so it waits for a separate task.

## Environment State

- Branch: main, clean, pushed at 2b0dd38 plus this handoff.
- Running processes: none.
- Sandbox home: 2b0dd38, `drift none`; its `~/.claude/settings.json` holds `{"attribution": {"commit": "", "pr": ""}}`.
- This box: the commit carrying this handoff, `drift none`; `sh-agent status` prints the attribution note (Next Steps 1).
- Pending: nothing uncommitted.

## Artifacts to Read

- `handoffs/2026-09-16-linkshelf-simulation-done.md` — the simulation: what both CLIs built, the 19 gaps, BUG-039..043, the Follow-ups.
- `changes/2026-09-16-retro-proposals-applied.md` and `retros/2026-09-16-retro.md` — what was applied and why.
- `bugs/INDEX.md` — the new Graduated table (BUG-029, BUG-040, BUG-043).

## Next Steps

1. **On this box:** add `"attribution": {"commit": "", "pr": ""}` to `~/.claude/settings.json`; `sh-agent status` then prints no note.
2. **On the Mac:** `sh-agent update --check`, `sh-agent update`, `sh-agent status`, the three smoke checks of `testing/2026-09-03-verify-real-home.md`.
3. **The mirrors** (styrman-crew, argolis-devbox) are at a81ef4b: refresh from this HEAD by each `VENDORED.md` recipe; argolis pushes deploy, crew ships by release (v1.34.0 recommended in its handoff).
4. **Tasks when convenient:** G20 (define `<app>` in the verifier); the simulation plan's other Follow-ups (G5, G10, G12, G13); option C's filing clause in `/preflight` Step 4 and BUG-031's wait rule from earlier handoffs; adopt's temp-directory copy and the direct `RenderContext` symlink test.
