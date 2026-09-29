---
owner: main
date: 2026-09-29
---
# Memory loop, Phase 1: one store for both CLIs     (date: 2026-09-29)

## Summary
The project's auto memory moves to a harness-owned directory both CLIs can read, `~/.config/sh-agent/memory/<slug>/` (spec `specs/sh-agent/memory-loop.md`; plan `plans/2026-09-29-dreamer-memory-loop.md`, the operator's option B):

- **`/init-project` Step 5b** (new and update paths): computes the slug the way Claude Code names the project, creates the store, writes `autoMemoryDirectory` into `.claude/settings.local.json` (gitignored; an empty object is overwritten, any other file gets one inserted line and a parse check with restore on failure; a tracked file is named for `git rm --cached`), and moves an existing default store's files in — move, never copy; both-hold-files and no-files cases report and move nothing.
- **The Resume Protocol** (`context/AGENTS.md`) ends with the store's `MEMORY.md` on Antigravity; Claude Code loads it itself.
- **`feedback-log` rule 7**: on Antigravity a correction or stated preference is also written as a memory file in Claude Code's shape, since agy has no native writer.
- **`/dream` Step 1** counts the settings store and the default path on its mandatory first line and stops on a split store.

## Type
feature (phase 1 of 2)

## Files changed
`skills/init-project/SKILL.md`, `context/AGENTS.md`, `skills/feedback-log/SKILL.md`, `skills/dream/SKILL.md`, `.agents/docs/specs/sh-agent/memory-loop.md` (new), the plan, `ROADMAP.md` (earlier commit).

## Agents dispatched
- `reviewer` ×3, the cold-reader check (`reviews/2026-09-29-review-dreamer-memory-loop-coldread-r1.md` to `-r3.md`); the cap was reached and the residue resolved as D13–D14 under the operator's "Go all on the cold-read gaps and keep building".
- `verifier` — `testing/2026-09-29-verify-dreamer-memory-loop-phase-1.md`, PASS.
- No `reviewer` over the diff yet and no `documentation`: the Definition of Done runs over the whole change at Phase 2's checkpoint.

## Verification
`testing/2026-09-29-memory-loop-phase-1-behavioural.md`: `/init-project` on an older scaffold with a seeded default store, 3 of 3 on store, migration and `.gitignore`, the settings write refused by Claude Code's headless `.claude/` guard in every run with the exact content reported; the absent-store message 2 of 2; the setting honoured headless 2 of 2 (the memory landed in the new store); the split-store stop 3 of 3 after nine void runs — the real home had been installed from 4426e58 mid-session and its user-level `dream` shadowed the project link (plan D15; the kits now pass `--setting-sources project,local`); a regression run kept the verdicts. Size test 21,787 of 22,000 bytes.

## Notes
- The interactive settings write, the Antigravity side and a multi-key settings file are Phase 2 or the operator's smoke (the log's Not exercised).
- The cleanup of the runs' real-home directories matched wider than the seeds (empty auto-created directories of earlier scratch probes); nothing held files.
