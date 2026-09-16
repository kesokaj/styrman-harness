---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — session done: simulation, retro applied, installer writes the attribution key, mirrors refreshed

Session ID: c50e0fd4-6273-4e35-b1ef-7639dff8bb07        (same machine: `claude --resume c50e0fd4-6273-4e35-b1ef-7639dff8bb07`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

One session, four pieces, each with its own handoff earlier today — this one closes the session. The operator asked for a harness simulation (both CLIs build one project from an empty repository under the harness), every gap fixed or filed, a retro, an install into this box; then said "go all" to the retro's eight proposals; then asked the installer to write Claude Code's attribution setting itself; then asked for both downstream mirrors refreshed from HEAD; then for the prompts to ship them.

## Completed

- **Simulation** (`handoffs/2026-09-16-linkshelf-simulation-done.md`): 19 gaps, BUG-039..043 fixed and re-verified in the sandbox, BUG-044 later; log `testing/2026-09-16-linkshelf-simulation.md`, plan `plans/2026-09-16-linkshelf-simulation.md` (Follow-ups G4, G5, G10–G13, G20).
- **Retro applied** on "go all" (`handoffs/2026-09-16-retro-applied.md`, 2b0dd38): eight proposals, two more restatements caught by the reviews, BUG-029/040/043 graduated; both queue files marked applied.
- **Installer writes the attribution key** (`handoffs/2026-09-16-attribution-write-done.md`, 42ba2f5, ADR-025; R4 and § Non-goals amended on the operator's word): `install` writes `"attribution": {"commit": "", "pr": ""}` when absent; proven end to end in the sandbox (install created the key, a headless commit carried no trailer). This box's `~/.claude/settings.json` holds the key (added by hand earlier, so the installer stays silent here).
- **Mirrors refreshed to 4aa19be** in their own repositories, each under its own Definition of Done and each needing the operator's approval of a "no settings files" rule: styrman-crew 6a0b39e + handoff 9dab61d (R5 amended, ADR-072, entrypoint logs the `settings:` lines, `verify-commit.sh` OK), argolis-devbox e772c28 + handoff bdb23a3 (ADR-013, DEPLOY.md smoke check changed). Committed locally, NOT pushed. The prompts to ship them were given in chat and are restated in Next Steps.
- This box installed at 4aa19be then reinstalled from later docs-only commits (`sh-agent status`: `drift none`, no note); sandbox home at 42ba2f5 with the installer-written settings.json.
- 19 commits since cad463d, every one with its `changes/` entry (five bug fixes, the simulation, the retro, the feature, the mirrors' entries are in their repositories); handoffs and cold reads for each piece.

## Open / In Progress

- Nothing in flight here; nothing queued in `~/.config/sh-agent/proposals/`; `TODO.md` and `ROADMAP.md` empty.
- Follow-ups in `plans/2026-09-16-attribution-write.md`: a fixed temp name for the settings write (a SIGTERM under argolis's `timeout 60` orphans `.settings-*.json`); one size check before either reader; refuse a read-only settings.json; `internal/harness/apply.go` lacks `Sync` before its rename; `abandonTemp`'s error paths untested. In `plans/2026-09-16-linkshelf-simulation.md`: G5, G10, G12, G13, G20 (define `<app>` in the verifier).
- Open records: BUG-019 to BUG-023, BUG-031, BUG-044. Retro candidates flagged by the bugs agent: a post-dispatch `git status --untracked-files=all` check (BUG-014/038/044), "hold source edits while a review is pinned".
- Requirements edits awaiting approval: none.

## Key Decisions

- The retro's proposals were applied on the operator's word; their statuses were written before the review ran — the next retro writes status after review (the code review of 2b0dd38 flagged it).
- The settings write: only an absent key; textual insertion over re-encoding; no state of the file fails install; `O_EXCL` on create, `Sync` before rename, symlinks refused; `uninstall` leaves the key (ADR-025, plan D1–D9).
- The cold read's cap was reached on a doc-line residue; the operator chose "Proceed, residue as Open rows".
- Each mirror's "no settings files" rule was amended by the operator rather than worked around with a harness opt-out flag.

## Corrections

- None: `~/.config/sh-agent/feedback.jsonl` has no entry for this repository today.

## Agent notes (unverified)

- Antigravity never opened the golang reference files during its build; Claude Code did. Retro Proposal 6 names them in the skill either way.
- Two reviewers this session reported the tree moving under them; both re-pinned. Whether to hold edits during a pinned review is a retro question, not a rule yet.
- The Go 1.24.4 toolchain on this box carries stdlib advisories per `govulncheck` (reported by three security reviews); an environment item.

## Environment State

- Branch: main, clean, pushed at 6072174 plus this handoff; this box reinstalled from the commit carrying this handoff.
- Running processes: none (`pgrep -f sim-kit` empty; crew's `verify-commit.sh` finished).
- Sandbox: `/home/kesokaj/sandbox/home` at 42ba2f5, `drift none`; its settings.json installer-written; the sim repositories and kits under `.agents/scratch/sim-kit/` and `sandbox/runs/`.
- Pending: nothing uncommitted here. styrman-crew and argolis-devbox each hold four local commits, not pushed.

## Artifacts to Read

- The three earlier handoffs of the day, named above — each is the full record of its piece.
- `plans/2026-09-16-attribution-write.md` § Follow-ups and `plans/2026-09-16-linkshelf-simulation.md` § Follow-ups — the deferred work.
- `retros/2026-09-16-retro.md` and `changes/2026-09-16-retro-proposals-applied.md` — what changed in the rules and why.
- `bugs/INDEX.md` — BUG-039..044 and the new Graduated table.

## Next Steps

1. **Ship argolis-devbox** (a push deploys to the VM): in that repository, resume from `handoffs/2026-09-16-harness-mirror-refresh-4aa19be.md`, push `main`, then the DEPLOY.md "Agent harness" block on the VM — expect the `settings: … created with` line in the install log, `repo has 27 skills, 7 agents`, no note, `drift none`, COMMIT 4aa19be, `grep -c '"attribution"' ~/.claude/settings.json` = 1, `doctor` green; record it as `testing/<date>-verify-vm-harness-refresh-4aa19be.md`.
2. **Ship styrman-crew** (a release): in that repository, resume from `handoffs/2026-09-16-harness-refresh-4aa19be.md`, push `main`, cut v1.34.0 through `/preflight` per `versioning.md` (the release entry names COMMIT 4aa19be7eeb5b15e4b876832bd81807cd5d91b37; a missing or partial security row is blocking); then on one re-seeded project expect the boot log's `agent harness deployed (… 4aa19be)` and `settings:` lines, the key in that home, and no trailer on a Claude Code commit there.
3. **The Mac:** `sh-agent update --check`, `sh-agent update`, `sh-agent status` (the install writes the attribution key there: expect the `settings:` line, then no note), then the three smoke checks of `testing/2026-09-03-verify-real-home.md`.
4. **Harness tasks when convenient:** the three settings-write Follow-ups (one task: fixed temp name, single size check, read-only refusal, with tests), `apply.go`'s `Sync`, G20, the other simulation Follow-ups, option C's filing clause in `/preflight` Step 4 and BUG-031's wait rule, adopt's temp-directory copy, the direct `RenderContext` symlink test.
5. **Next retro candidates:** the post-dispatch untracked-files check (BUG-044, third of its shape); status after review; holding edits during a pinned review.
