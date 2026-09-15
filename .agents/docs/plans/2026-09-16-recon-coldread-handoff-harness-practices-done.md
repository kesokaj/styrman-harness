---
owner: explorer
date: 2026-09-16
---
# Recon — coldread-handoff-harness-practices-done   (2026-09-16)

## Question
Cold-read the 2026-09-16 handoff (`.agents/docs/handoffs/2026-09-16-harness-practices-done.md`) against `AGENTS.md`, `TODO.md`, `ROADMAP.md`, the active plan, the three newest `changes/` entries and `bugs/INDEX.md`; answer six specific checks.

## Summary
The trail is internally consistent: the plan `2026-09-15-harness-practices-g01-g26.md` is fully ticked including its Done list, `ROADMAP.md` "Now"/"Next" and `TODO.md` are empty (matching the handoff's claim), and every artifact the handoff names exists. No contradictions found between handoff, ROADMAP, TODO and the plan. The fence-aware awk check finds zero unticked boxes outside fenced blocks in any plan file. All commits since the previous handoff (88f9771..HEAD) that touch paths outside `.agents/docs/` are covered by named `changes/` entries per phase (phase1–phase5, done-review, phase7-decisions/init-update/gate-docs-advisory, phase6-e2e, close). Both mirror repos (styrman-crew, argolis-devbox) exist locally with the exact commits the handoff names, unpushed as stated. Nothing is left to guess beyond what the handoff itself flags as open (BUG-031, the option-C filing clause, operator-only Next Steps 1–3).

## Relevant files
- `.agents/docs/handoffs/2026-09-16-harness-practices-done.md` — the handoff under cold read.
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — closed plan; Phases 1–7 and Done list all `[x]`, evidence inline, 51 Decisions (D1–D51).
- `ROADMAP.md`, `TODO.md` — both empty (headers/no "Now" line), consistent with handoff.
- `.agents/docs/bugs/INDEX.md` — BUG-001..037; open rows BUG-019..023, BUG-031, BUG-033..037 match what the handoff lists as open.
- `.agents/docs/changes/2026-09-16-harness-practices-close.md` and the other 9 phase/close entries — cover all 13 commits from cec4d4e through a81ef4b (88f9771 is the prior handoff commit, d1ff612 is the plan-cap commit before it).
- `/home/kesokaj/workspace/styrman-crew` (HEAD 4957dfe, local commit 5418f19) and `/home/kesokaj/workspace/argolis-devbox` (HEAD 2d83282, local commit 9380a63) — both mirrors exist, both hold exactly the commits the handoff names, both unpushed.

## Flows
Resume order followed as instructed: AGENTS.md → TODO.md (empty) → handoff → ROADMAP.md Now (empty) → plan's unticked boxes (none) → three newest changes entries → bugs/INDEX.md. No missing link in the chain.

## Conventions
n/a for this recon — no code was read.

## Answers to the six questions

1. **Status / next step, as the trail states them:** the plan is closed, Done list fully ticked, main box of work finished. The single next step per the handoff's own framing is operator action: Next Steps 1–3 (Mac install/mirror pushes) are explicitly "operator" work; the first thing an AI resuming this session could act on unprompted is Next Step 4, the harness security session over BUG-033–037 (starting with the `adopt --force` deletion bug), since 5–8 are lower priority/recommendations.

2. **Files named that do not exist:** none found. Checked all testing/, changes/, reviews/, DECISIONS.md (ADR-024 at line 189) paths named in the handoff's "Artifacts to Read" list plus both mirror repos — all present with the exact commits cited.

3. **Contradictions among handoff / ROADMAP / TODO / plan:** none found. ROADMAP "Now" and "Next" are empty, TODO is headers-only, the plan's Done list closing box confirms both are empty as of 2026-09-16 — all three agree with each other and with the handoff.

4. **Plans with an unticked box outside fenced blocks not named by ROADMAP "Now":** the awk check returns no output — zero unticked boxes outside fenced blocks across all 21 files in `.agents/docs/plans/*.md`. (ROADMAP "Now" being empty is therefore consistent, not a gap.)

5. **Commits since 88f9771 touching paths outside `.agents/docs/` with no covering changes entry:** none. All 13 commits (cec4d4e, 32c8427, f3013a6, b43dde1, 22c106c, 641eccd, fe21671, c514936, e31c3a2, fcf3b79, 2a2e6a1, a33978d, 8202d6a, a81ef4b) touching `agents/`, `context/`, `skills/`, `cmd/sh-agent/`, `internal/harness/`, `AGENTS.md`, `README.md` are each covered by name in the handoff's Completed section and by one of: phase1–phase5, done-review, phase7-decisions, phase7-init-update, phase7-gate-docs-advisory, phase6-e2e, or the uncommitted close entry.

6. **Left to guess:** nothing beyond what the handoff itself already flags as open/unverified — specifically (a) whether the four Antigravity timeouts (D50) are load-caused, called "unverified" by the handoff's own Agent notes; (b) BUG-031's headless-run-ends-early mechanism, open; (c) option C's filing clause (Next Step 6), explicitly deferred to the operator; (d) whether BUG-029/032's prevention rules graduate (Next Step 8), explicitly deferred to the next `/retro`.

## Risks
None found in the trail itself. Operationally: Next Step 4 warns that `sh-agent adopt --force` can delete `skills/` or the whole checkout including `.git` and unpushed work (BUG-037) — worth flagging loudly to whoever picks up next if `adopt` is used before that fix lands.

## Open questions
None outstanding from this recon — all six checks resolved cleanly against the working tree as found.

## Not examined
- Contents of styrman-crew's and argolis-devbox's own handoffs/changes entries beyond the `git log` one-liners shown above (not needed to answer the six questions).
- The full text of every `changes/` entry body (only newest-3 plus existence/coverage-by-commit was checked, per the dispatch's scope).
- `~/.config/sh-agent/proposals/2026-09-14-retro.md` in full (only tail read, to confirm it exists and holds the two queued proposals the handoff cites).
