---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — harness practices G01–G26 done

Session ID: b12af1f5-5a11-41f7-8a72-867bf23337a4        (same machine: `claude --resume b12af1f5-5a11-41f7-8a72-867bf23337a4`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The operator asked for groups G01–G26 of the 2026-09-15 research record (beads, elephant-goldfish, goldfish) to become harness practice. They then added Phase 7:
- keep R2;
- decide how the release gate handles docs findings;
- install on this box;
- prepare the two downstream mirrors;
- let `/init-project` update an older scaffold;
- finalize all documentation, as the Definition of Done.

That work is done, and the plan is closed.

## Completed

- **Plan:** `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — Phases 1–7 and the Done list ticked with evidence. `ROADMAP.md` "Now" and `TODO.md` are empty.
- **Phases 1–5** (cec4d4e, 32c8427, f3013a6, b43dde1, 22c106c): text edits I-1 to I-25 to the context, skills and agents. Each phase has its changes entry and a verifier log.
- **Done review (641eccd) and Phase 7** (fe21671, c514936, e31c3a2):
  - R2 kept (D45).
  - `/init-project` updates an older scaffold (D47).
  - The release gate's docs review is advisory (D49, ADR-024).
  - The installer and `status` tell the operator where a hand-written `CLAUDE.md` belongs.
  - The Phase 7 review round 2 is APPROVE (2a2e6a1).
- **Phase 6** (a33978d): the matrix on both CLIs, logged in `testing/2026-09-15-verify-harness-practices-e2e.md`.
  - 26 of 30 row–CLI pairs pass.
  - Four Antigravity rows are recorded as caused outside this plan (D50).
  - Every row for skills that changed after the matrix ran passes at fcf3b79.
- **Done verifier:** FAIL at a33978d on a missing test, fixed in 8202d6a (BUG-032), then PASS at 8202d6a (`testing/2026-09-15-verify-harness-practices-done.md`).
- **7.5** (a81ef4b): the harness is installed into this box's real home, and the three smoke checks pass on both CLIs (`testing/2026-09-15-verify-real-home.md`). After a81ef4b it was reinstalled from the clean checkout: `binary a81ef4be3cb8`, drift none.
- **7.6:** styrman-crew — refreshed from cceee54 to a81ef4b in 5418f19, with handoff 4957dfe. Its verifier passed, including the rebuilt runtime image (integration 94 passed, 0 failed, 3 skipped) and the mirror's own tests (32, 0, 1). Its code and docs lenses returned APPROVE, and its security review left 2 medium and 3 low findings to this repository; argolis-devbox — refreshed from 8046e4d to a81ef4b in 9380a63, with handoff 2d83282. Its verifier passed, its code review returned APPROVE WITH NITS with the docs finding fixed, and its security review left 2 medium and 4 low findings to this repository. Neither is pushed.

## Open / In Progress

- Nothing in the plan. Its Follow-ups list what was noticed and not built; two need the operator (Next Steps 6 and 7).
- Open bug records:
  - BUG-019 to BUG-023: sweep findings in skills and references.
  - BUG-031: a print-mode run, or an agent, can end its turn before what it started returns.
  - BUG-033 to BUG-037: the two mirrors' security findings on harness text and code (Next Step 4).
- Requirements edits awaiting approval: none.

## Key Decisions

- **D45:** R2 stays; a project's `CLAUDE.md` is the single line `@AGENTS.md`.
- **D47:** in update mode, `/init-project` adds only missing template files, `.claude` links and `.gitignore` lines, and never rewrites what the project wrote.
- **D49 / ADR-024:** the release gate's docs review is advisory. A docs finding blocks only at the task's own `/review --docs`.
- **D50:** a row that failed twice under five concurrent CLIs, then passed both with c192622 and with this tree under a lighter load, is caused outside this plan. It is not counted as a pass.
- **D51:** the mirrors are refreshed from a81ef4b. Among the paths a mirror copies, the closing commit changes only `ROADMAP.md`'s "Now" line.
- **Rejected:** patching option C's filing clause in this plan. A fifth round on BUG-026's mechanism goes to the operator instead.

## Corrections

- None logged: `~/.config/sh-agent/feedback.jsonl` has no entry from this session.

## Agent notes (unverified)

- **The four Antigravity timeouts** look like load. B1 and C1 fit that, but no run separated local CPU from the service side.
- **Claude Code's headless runs** could not create `.claude/` links or read skills' `references/` from outside the clone. An interactive session prompts for both.
- **Scratch checks and guards:** two more of this session's own checks matched their own command line (a `pgrep` count and a live-process scan). Both were caught before anything relied on them; that is BUG-029's rule.

## Environment State

- **Branch:** main, clean and pushed after the handoff commit.
- **Running processes:** none. The sandbox kit, the smoke script and both mirror sessions have exited.
- **Pending:** nothing here. styrman-crew and argolis-devbox each hold local commits, not pushed (Next Steps 2 and 3).

## Artifacts to Read

- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — the closed plan, its Decisions D1–D51 and its Follow-ups.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-e2e.md` — the matrix per row and CLI, D50's evidence, and what no row examined.
- `.agents/docs/testing/2026-09-15-verify-real-home.md` — the real-home install and smoke checks.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-done.md` — the Done verifier, PASS at 8202d6a.
- `.agents/docs/changes/2026-09-16-harness-practices-close.md` — the last entry: BUG-032's test, 7.5, 7.6 and the close.
- `.agents/docs/bugs/INDEX.md` — this plan's records, BUG-024 to BUG-032.
- `~/.config/sh-agent/proposals/2026-09-14-retro.md` — two harness proposals from styrman-crew, still queued.
- styrman-crew and argolis-devbox — each repository's newest handoff and changes entry, for its mirror refresh.

## Next Steps

1. **On the Mac:**
   - `sh-agent update --check`, `sh-agent update`, then `sh-agent status`: expect drift none and a binary line with no warning.
   - Claude Code: `/memory` and `/skills`. Antigravity: `/skills` and `/agents`.
   - The three smoke checks of `.agents/docs/testing/2026-09-03-verify-real-home.md`, on each CLI.

   On this box, run `sh-agent update` once more; it only moves the binary's build stamp to the closing commit, since the installed files already match.
2. **argolis-devbox:** read its local commits 9380a63 and 2d83282 and its handoff (it asks for the recipe's literal `diff -rq` before the push, the VM check by DEPLOY.md's Agent harness block, and a decision on a `verify_mirror` guard that refuses anything but regular files and directories), then push. A push to `main` deploys to the VM within five minutes; verify it on the VM by `.agents/docs/testing/devbox/DEPLOY.md`.
3. **styrman-crew:** its local commits 5418f19 and 4957dfe refresh `image/harness/`. Ship them through a release, by `.agents/docs/specs/styrman/versioning.md`: push, then `/preflight` (its handoff recommends `v1.34.0`, a MINOR). Only a release makes running projects re-seed. In that gate, treat a missing or partial security row as blocking (crew's merged review, finding 2).
4. **Harness security session** over BUG-033 to BUG-037 and the plan's 7.6 Follow-ups. Start with the `sh-agent adopt` record: `adopt .` or `..`, or adopting a skill onto itself, together with `--force`, removes `skills/` or the whole checkout, `.git` and unpushed work included. Its fix: make `from` absolute, run `validName` on its base, refuse when `from` and `dest` overlap, and copy into a temporary directory before renaming it over `dest`. Until it lands, never pass `--force` to `adopt`.
5. **In this checkout:** run `/retro --harness` for the two proposals from styrman-crew queued in `~/.config/sh-agent/proposals/2026-09-14-retro.md`.
6. **Decide option C's filing clause** (plan Follow-ups). Antigravity's `/preflight` recorded its docs gaps but filed them nowhere a resume reads.
   - Recommendation: one sentence in `skills/preflight/SKILL.md` Step 4. The gate appends its docs gaps to the newest handoff's Next Steps, and writes a handoff when the project has none. That is what the Claude Code run did.
7. **BUG-031.** Recommendation: one sentence in `context/AGENTS.md`, beside the `/goal` line: a headless run waits for every agent it dispatched before it ends its turn. Put the same rule in the agent definitions, since crew's `verifier` ended its turn the same way. Verify it with the `documentation` row on Claude Code.
8. **For the next `/retro`:**
   - whether BUG-029's and BUG-032's prevention rules graduate;
   - the setup command `gofmt -l .`, which walks `.agents/scratch/` (plan Follow-ups).
