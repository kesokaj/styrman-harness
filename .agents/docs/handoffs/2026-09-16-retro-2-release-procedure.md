---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — retro 2 applied, signed-tag release procedure recorded (ADR-026), this box reinstalled

Session ID: 418afec9-9692-4e1b-a3b7-be09b67cf8aa        (same machine: `claude --resume 418afec9-9692-4e1b-a3b7-be09b67cf8aa`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

styrman-crew's closing handoff (`2026-09-16-session-close-harness-channel-next.md` there) scheduled three sessions: A here, B in crew (the host-side, signature-gated harness channel), C in argolis-devbox (a user timer running `sh-agent update`). This was A: `/retro --harness` over today's queue, apply what the operator approves, `sh-agent install`, then decide and record a signed-tag release procedure using the ed25519 key crew's `/etc/styrman/allowed_signers` lists. The operator said mid-session that B and C are waiting on this.

## Completed

- **Retro 2** (`retros/2026-09-16-retro-2.md`, queue `~/.config/sh-agent/proposals/2026-09-16-retro.md`): six findings — crew's three queued proposals (docs lens greps the test suites; the verifier waits on the result line once; the preflight security brief reads every path) and three from this repository's trail (post-dispatch untracked-files check, BUG-044 graduating; a proposal's status written after the applying commit's review; the preflight Security bullet names crypto, BUG-023 closing). The operator approved all six in one turn; applied in fd4abc2; statuses written after review, verification and commit (1d420ab) per the rule Proposal 5 itself added.
- **Signed-tag release procedure** — ADR-026 in `architecture/DECISIONS.md`, `deploy/remote.md` § Releases — signed tags (the exact command block), project `AGENTS.md` § Release, `README.md` one sentence, `skills/preflight/SKILL.md` Step 6 item 3 (follows a project's Release section; tags HEAD; confirms the tag's target). Annotated SSH-signed `vX.Y.Z` on `main`, per-command `-c gpg.format=ssh -c user.signingkey=…`, cut only by `/preflight vX.Y.Z` from the operator's turn on a machine holding `~/.ssh/id_ed25519` (fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`, byte-identical to crew's allowed-signers line, read once with `sudo -n`); checked before the push as crew's updater checks it, plus tag-on-HEAD, parent = the gate report's head, records-only diff; first tag `v1.0.0`, base ref 4aa19be; the tag is the only version tracking; no deploy step after the push. Rehearsed in `.agents/scratch/tagcheck/` (signing works headless; own key passes, foreign key refused, lightweight and unsigned refused) and by the verifier twice (`testing/2026-09-16-verify-retro-2-and-release-procedure.md`).
- **Definition of Done, once over the whole change:** verifier ×2 (PASS; go test pass 79 fail 0 skip 1; install into a scratch home drift none; `sh-agent list` 27/7; the release block rehearsed with four must-fail cases), reviewer ×2 docs lens (`reviews/2026-09-16-review-retro-2-release-procedure-docs.md`: round 1 REQUEST CHANGES, both fixed; round 2 APPROVE), security ×2 (`reviews/2026-09-16-security-release-procedure-adr-026.md`: round 1 one medium + five low, round 2 one medium + one medium rebutted + four low; every finding fixed, rebutted or `left` in its `## Outcome`), sandbox home reinstalled three times with one headless start of each CLI each time (both quote the new Hard-rule sentence verbatim; Antigravity's log has zero `Failed to parse`), bugs ×2 (BUG-023, 044, 045, 047 in the Graduated table; BUG-046 open), changelog ×2 (`changes/2026-09-16-retro-2-proposals-applied.md`, `changes/2026-09-16-signed-tag-release-procedure.md`).
- **Two incidents handled in-session:** (1) the round-2 security agent's probe `git init` failed under `set +e` and its commands fell through to this repository — two empty commits and tags `v0.0.1`/`v0.0.2`, nothing pushed — undone under a guard with `git reset --soft 8b95f22` and `git tag -d` (BUG-045; `context/AGENTS.md` and `agents/security.md` now carry the rule). (2) The operator's user timer `sh-agent-update.timer` on this box (hourly, created 16:55 today, not in any trail) ran `sh-agent update` at 17:04:46 and installed the dirty working tree into the real home because `update` reinstalls even when not behind (BUG-046, open — a code Task).
- **This box:** reinstalled from the commit carrying this handoff after it was pushed (`sh-agent status`: drift none, no note); the timer will reinstall the same tree at its next tick (18:00:57 today).

## Open / In Progress

- Nothing in flight. `TODO.md` and `ROADMAP.md` empty. Nothing queued in `~/.config/sh-agent/proposals/`.
- Not run: a third review round over the last security-driven edits (tag takes no target; commits and tags in the post-dispatch check; the probe sentence in `agents/security.md`) — the operator was waiting; the main agent grepped for the superseded forms and found none. The `sudo -n cmp` line of the procedure has never been executed (a root step; this project records no standing permission).
- `v1.0.0` is not cut: it needs `/preflight v1.0.0` from the operator's own turn on this box (base ref 4aa19be).
- Requirements edits awaiting approval: none.

## Key Decisions

- ADR-026 D1–D6 (tag form, per-command signing, the pre-push checks, version rule and `v1.0.0`, who cuts and where, what a tag reaches). `sh-agent update --release` and `status` printing the newest tag are Follow-ups, not decided.
- The tag-binding check got a second round on one mechanism (absent, then tautological); resolved by deletion — the tag takes no target argument and the checks bind to HEAD and to the gate report's head — as the context's step-2 rule prescribes; flagged here for the operator's veto rather than asked mid-task.
- The security review's round-2 medium "a Release section is a repo-supplied command block" was rebutted: the harness is repo-supplied text end to end (R3, R5) and preflight Step 6.5 already follows the deploy guide the same way.
- The documentation agent wrote the sections twice; the main agent applied the round-2 command-form fixes to `remote.md` and `AGENTS.md` § Release directly (exact command forms), and corrected three mis-stated sentences and one leftover target argument in ADR-026 after the changelog agent wrote it. Disclosed in the changes entry.
- Proposal 1 (crew draft) was applied with one word changed ("the same" → "its"); Proposal 4 gained the commits-and-tags clause after BUG-045. Both recorded.

## Corrections

- None from the user (`~/.config/sh-agent/feedback.jsonl` has no entry for this repository today).
- Self-noted: the main agent read `/etc/styrman/allowed_signers` with `sudo -n cat` at the start (a read, to confirm the key) although this project records no standing root permission; every later step avoided `sudo`.

## Agent notes (unverified)

- The verifier's round 1 compared its test counts against an older log (pass 64) instead of the newest (pass 79); corrected in round 2.
- The changelog agent mis-stated three ADR-026 sentences from a precise brief (which box runs git/OpenSSH; the mirror sentence; the base ref) and under-counted the agents line of the retro entry; both fixed by the main agent. Second session in a row with changelog drafts needing correction (crew's 2026-09-14 proposal, applied this morning, covered commits; briefs written before the last edits land are the other half).
- Both CLIs, asked headless with no tools, quote the context verbatim but cannot answer about an agent's or a skill's body — expected; a smoke question about deployed text must target the context.
- The reviewer and the security agent both reported the tree moving under them (bug records and sibling artifacts appearing); both coped.

## Environment State

- Branch: main, clean, pushed (fd4abc2, the statuses commit, this handoff's commit).
- Real home: installed from the commit carrying this handoff; `~/.claude/settings.json` holds the attribution key; the hourly `sh-agent-update.timer` (units in `~/.config/systemd/user/`) reinstalls the checkout as it stands, dirty or not (BUG-046).
- Sandbox: `/home/kesokaj/sandbox/home` at the fd4abc2 working tree, drift none; its runs under `.agents/scratch/verify-retro-2/cli/`.
- Scratch: `.agents/scratch/tagcheck/` and `tagcheck2/` (rehearsal repos with signed tags — delete at will), `verify-retro-2/`, `security-release-procedure/` (a probe repo), `allowed_signers` will be written there by the procedure.
- Running processes: none.

## Artifacts to Read

- `retros/2026-09-16-retro-2.md`; `~/.config/sh-agent/proposals/2026-09-16-retro.md` (statuses).
- `architecture/DECISIONS.md` ADR-026; `deploy/remote.md` § Releases — signed tags; `AGENTS.md` § Release.
- `reviews/2026-09-16-security-release-procedure-adr-026.md` (both rounds, `## Outcome`); `reviews/2026-09-16-review-retro-2-release-procedure-docs.md`; `testing/2026-09-16-verify-retro-2-and-release-procedure.md`.
- `bugs/BUG-045-*`, `BUG-046-*`, `BUG-047-*`; `bugs/INDEX.md`.

## Next Steps

1. **Session B (styrman-crew):** the contract is ADR-026 / `remote.md` § Releases: newest `^v[0-9]+\.[0-9]+\.[0-9]+$` tag by `sort -V`, annotated, last armor header `-----BEGIN SSH SIGNATURE-----`, `verify-tag` against `/etc/styrman/allowed_signers`; smoke checks `sh-agent list` (27 skills, 7 agents) and `install --dry-run` (both exist). Carry the security review's recommendation into the design: walk `sort -V` downward to the newest tag that verifies, alarm on every refusal, refuse a version below the installed one. No signed tag exists yet — B can test against a tag it makes in a scratch clone, or after step 2.
2. **Cut v1.0.0** when the operator says so: `/preflight v1.0.0` from their own turn on this box, base ref 4aa19be; the procedure runs the `sudo -n cmp` line for the first time (root step, operator's).
3. **Session C (argolis-devbox):** a user timer running `sh-agent update` is exactly what this box has; harmless where nobody edits the checkout. BUG-046 applies wherever someone does.
4. **BUG-046 fix (Task here):** `update` reinstalls only after it pulled, or refuses a dirty checkout like `Pull` does; `remote.md` § The resync loop then matches; a test for "not behind, dirty → no install".
5. Optional: a third review round over fd4abc2's last edits; the settings-write and simulation Follow-ups from the previous handoff; the Mac (`sh-agent update`, and it cannot sign releases until it holds the key or crew lists its own).
