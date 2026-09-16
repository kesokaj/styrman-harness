---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — install writes the attribution setting; nothing queued

Session ID: c50e0fd4-6273-4e35-b1ef-7639dff8bb07        (same machine: `claude --resume c50e0fd4-6273-4e35-b1ef-7639dff8bb07`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

Third handoff of the day, after `2026-09-16-linkshelf-simulation-done.md` (the simulation) and `2026-09-16-retro-applied.md` (the eight proposals). The operator asked whether the installer should write the attribution setting itself ("so when I install it fixes this?"), approved the two REQUIREMENTS edits it needs, and said "keep going, notify me when it's done". Done: `sh-agent install` now writes the key when absent, at 42ba2f5.

## Completed

- **Feature** (`changes/2026-09-16-attribution-write.md`, plan `plans/2026-09-16-attribution-write.md`, every box ticked): `install`/`sync` for vendor all or claude writes `"attribution": {"commit": "", "pr": ""}` into `<home>/.claude/settings.json` when the `attribution` key is absent — file created when missing; a present key of any value left alone (`status` keeps its note when the file is missing or sets `commit` to anything else); a textual insertion re-parsed before landing, so every other key, the formatting, the mode and the line endings stay; `O_EXCL` on create, `Sync` before the rename; a symlinked, non-object, malformed or over-1-MiB file refused with one advisory line; `--dry-run` writes nothing; `uninstall` leaves the key. `TestAttributionSet`, twelve subtests, three shown to fail without their code.
- **Requirements, with the operator's approval:** R4 gained "`install` also sets one key … only when the key is absent …"; the Non-goals bullet gained the one exception; ADR-025 records it and ADR-010 is superseded on that clause only — the permission write it refused stays refused. README and both deploy guides say `install` does it.
- **Reviews:** three cold-read rounds (round 3 `implementation ready`; the operator's override at the cap for one doc-line residue), code lens APPROVE WITH NITS (3 fixed), security FINDINGS three low (all fixed; it confirmed four earlier defects closed), docs lens REQUEST CHANGES (2 fixed). Verifier PASS three times, the last on the final tree: counts pass 79, fail 0, skip 1.
- **Installed:** this box and the sandbox home from 42ba2f5, `drift none`; sandbox end to end — key removed, `install` created it, `status` 0 notes, a headless Claude Code commit with no trailer, one start per CLI, 0 Antigravity parse errors.
- **BUG-044** recorded: the security review's probe files landed in the repository root (seven files named by their contents), removed by name; the bugs agent flags the post-dispatch `git status --untracked-files=all` check as a graduation candidate (third record of the shape: BUG-014, BUG-038, BUG-044).

## Open / In Progress

- Nothing in flight; nothing queued in `~/.config/sh-agent/proposals/`.
- Follow-ups in the plan: `internal/harness/apply.go:240-263` writes managed files without `Sync` (the same pattern the security review flagged here); `abandonTemp`'s error paths have no forcing test. The simulation plan's Follow-ups (G5, G10, G12, G13, G20) stand.
- Open records: BUG-019 to BUG-023, BUG-031, BUG-044.
- Requirements edits awaiting approval: none.

## Key Decisions

- D1–D9 in the plan: only an absent key is written; textual insertion over re-encoding (key order); no state of the file fails install; `sync` shares the code path; non-object files get one line from `attributionHint`; no error return.
- The cold read's cap was reached on a doc-line residue; the operator chose "Proceed, residue as Open rows"; both Open rows resolved with their defaults (a shared bounded reader; the test drives `run()`).
- The verifier's symlink flag and the security review's `O_EXCL`/`Sync`/single-line findings were fixed in the same change, each with a test.

## Corrections

- None: no feedback-log entry for this repository today.

## Agent notes (unverified)

- The security review noted the tree moved under it (the verifier's and the code review's fixes landed mid-review); it re-pinned and finished against the later blobs. A rule for the next retro: hold source edits while a review is pinned, or re-dispatch after them.
- The verifier reported the sandbox home "does not exist on this box yet" in an earlier run — a dispatched agent's cwd blind spot, not a missing directory.

## Environment State

- Branch: main, clean, pushed at 42ba2f5 plus this handoff.
- Running processes: none.
- This box: installed from the commit carrying this handoff (`drift none`); its `~/.claude/settings.json` holds the attribution key (added by hand earlier today; the installer now finds it present and stays silent).
- Sandbox home: 42ba2f5, `drift none`; its `~/.claude/settings.json` is the installer-written one (`{\n  "attribution": {"commit": "", "pr": ""}\n}`); the previous hand-written copy is at `.agents/scratch/sim-kit/results/attr-e2e-sandbox/settings.json.before`.
- Pending: nothing uncommitted.

## Artifacts to Read

- `changes/2026-09-16-attribution-write.md` — what changed and how it was verified.
- `plans/2026-09-16-attribution-write.md` — decisions D1–D9, Open rows, Follow-ups.
- `architecture/DECISIONS.md` ADR-025 — the exception to the vendor-settings non-goal and its limits.
- `reviews/2026-09-16-security-attribution-write.md` — the 20-input probe of the write path and its Outcome.
- `handoffs/2026-09-16-linkshelf-simulation-done.md` and `handoffs/2026-09-16-retro-applied.md` — the rest of today.

## Next Steps

1. **On the Mac:** `sh-agent update --check`, `sh-agent update`, `sh-agent status` — the install writes the attribution key there (expect the `settings: … created with …` line, then no note); then the three smoke checks of `testing/2026-09-03-verify-real-home.md`.
2. **The mirrors** were refreshed to 4aa19be later on 2026-09-16 ("refresh the two mirrors from HEAD"): styrman-crew 6a0b39e (+ handoff 9dab61d, `verify-commit.sh` OK, R5 amended and ADR-072 on the operator's word), argolis-devbox e772c28 (+ handoff bdb23a3, ADR-013 on the operator's word). Both committed locally, NOT pushed: crew ships by release (v1.34.0 recommended), an argolis push deploys to the VM. Their security reviews filed three low findings here — the Follow-ups in `plans/2026-09-16-attribution-write.md` (fixed temp name, one size check, refuse a read-only file).
3. **Tasks when convenient:** the `Sync` in `internal/harness/apply.go`; G20 (define `<app>` in the verifier); the simulation plan's other Follow-ups; option C's filing clause in `/preflight` Step 4 and BUG-031's wait rule from earlier handoffs; adopt's temp-directory copy and the direct `RenderContext` symlink test.
4. **Next retro candidates:** the post-dispatch `git status --untracked-files=all` check (BUG-044, third of its shape); "hold source edits while a review is pinned".
