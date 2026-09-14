---
owner: main
date: 2026-09-14
---
# Handoff — 2026-09-14 — upstream items from styrman-crew landed (cceee54); the crew mirror refreshes to it

Session ID: session_01RUoU1zHUVbathNBFLo27xf        (a styrman-crew session working here cross-repo; same machine: `claude --resume`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

styrman-crew's v1.33.0 close-out plan (its `.agents/docs/plans/2026-09-14-v1.33.0-close-out.md`, Phase 5) carried three items for this repo since its v1.31.0 mirror refresh. They were done here under this repo's own Definition of Done, plus three security findings the review raised on the way.

## Completed

- Commit `cceee54` (pushed): `TestDryRunNeverBuilds` builds its own checkout-shaped fixture (BUG-016); `validName` at the parse boundary for skill/agent `name:` (one local identifier-like segment); `LoadState` drops non-local manifest entries (one guard for install/uninstall/status); the skill walk and `adopt`'s copy refuse symlinks (security high: a link's target was deployed under the link's mode); `context/dev-setup.md:8` conditions the root-steps sentence on the project's recorded permission and says what to do without it.
- Records: `.agents/docs/changes/2026-09-14-upstream-items-from-styrman-crew.md`; `.agents/docs/bugs/BUG-016-dry-run-test-assumed-its-own-checkout.md`; verifier `.agents/docs/testing/2026-09-14-upstream-items-verify.md` (PASS, two rounds); reviewer `.agents/docs/reviews/2026-09-14-review-upstream-items.md` (APPROVE WITH NITS, Outcome written); security `.agents/docs/reviews/2026-09-14-security-upstream-items.md` (three findings, all FIXED, residuals recorded).
- Verified from where it runs: sandbox `install --home .agents/scratch/x` rc 0, `status` drift none, the rendered `.claude/AGENTS.md` carries the new sentence; `sh-agent list --repo .` Skills (27) Subagents (7).

## Open / In Progress

- Not run: the sandbox CLI-start matrix (`.agents/docs/testing/2026-09-03-verify-sandbox.md`), the full skill matrix, and the real-home smoke checks — the change to `context/` is one sentence; the verifier log names the gap.
- Downstream: styrman-crew refreshes `image/harness/` to `cceee54` in its v1.33.0 (in progress in the same session); argolis-devbox's mirror refresh is a separate session with the same prompt shape as `.agents/docs/plans/2026-09-14-mirror-refresh.md`.
- Two harness-scope proposals from crew's retro are queued in `~/.config/sh-agent/proposals/2026-09-14-retro.md` (`agents/changelog.md`: never infer a commit's contents or a verifier's scope; `context/AGENTS.md` DoD step 4: "EMPTY means headers only, at every checkpoint") — for `/retro --harness`.
- Requirements edits awaiting approval: none.

## Key Decisions

- Names are ASCII identifier-like, not just path-shaped: the same value is rendered verbatim into every consumer's AGENTS.md.
- The manifest guard lives in `LoadState`, once, rather than in each reader.
- Symlinks in skill directories are refused, not dereferenced or copied as links.

## Corrections

- none in this repo this session (the crew session's corrections are in its own log).

## Agent notes (unverified)

- The changelog draft needed three wording corrections (called the work a "backport", misquoted the mode arithmetic, over-claimed round 2 before it ran) — the pattern the queued `agents/changelog.md` proposal addresses.

## Environment State

- Branch: main, clean at `cceee54` plus this handoff's commit; pushed.
- Scratch: `.agents/scratch/x`, `.agents/scratch/verify-2026-09-14*` (sandbox homes, gitignored).

## Artifacts to Read

- `.agents/docs/changes/2026-09-14-upstream-items-from-styrman-crew.md` — what changed and why, with the review dispositions.
- `.agents/docs/reviews/2026-09-14-security-upstream-items.md` — the three findings and the residuals left open.

## Next Steps

1. `/retro --harness` for the two queued proposals.
2. Operator: run the sandbox CLI-start matrix once on this commit if the one-sentence `context/` change should count as "complete" by this repo's Testing section.
3. argolis-devbox: refresh its mirror to `cceee54` in its own session.
