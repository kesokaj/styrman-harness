---
owner: main
date: 2026-09-28
---
# Handoff — 2026-09-28 — OKF researched: not now (A approved)

Session ID: f538b90a-5aa7-44d1-8311-a160b5bb0167
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main`, starting from 914dc90. This handoff's commit carries the record, the change entry, the verify log and the frontmatter fix.

## Context

The operator asked whether the harness implements OKF (Open Knowledge Format, Google Cloud's markdown knowledge-bundle format, June 2026) and whether it should, and wanted a report to approve.

## Completed

- The research record, `.agents/docs/plans/2026-09-28-research-okf.md`:
  - **Answer:** OKF is not implemented.
  - **Recommendation:** A, watch with named triggers. The operator approved it and answered §6: "no BigQuery work on other machines".
- The record's found-along-the-way #2 is fixed, as the operator asked: `plans/2026-09-21-context-audit-prompt.md` gained `owner:`/`date:` frontmatter.
- The verifier passed: `testing/2026-09-28-verify-research-okf.md`.
- The change entry is `changes/2026-09-28-research-okf.md`. It names the 8 workflow agents and the verifier.

## Open / In Progress

- Nothing is in progress, and `ROADMAP.md` "Now" and `TODO.md` are empty.
- Requirements edits awaiting approval: none.
- The research record's found-along-the-way items #1, #3, #4 and #5 were not approved and stay open (Next Steps).

## Key Decisions

- **OKF: not now, not never.** The reasons:
  - no CLI reads it;
  - no project here produces or consumes it;
  - the spec is unstable (one committer, "0.2" changed in place);
  - every adopting option costs more than it returns.
- **Rejected options:**
  - vendoring a community skill breaks the zero-scripts rule, ADR-006/007;
  - a harness skill has no user;
  - a conformant trail adds nothing, because the directories already type every file;
  - an `sh-agent` validator conflicts with R4;
  - a formal non-goal claims more than the evidence supports.
- **No ADR:** nothing structural changed.
- **The trigger that would reopen it:** T2, a project that must import, publish or maintain an OKF bundle. The first step would then be a skill in that project, with no harness change.

## Corrections

- None this session. `feedback.jsonl` has no entry dated 2026-09-28.

## Agent notes (unverified)

- The research scratch in `.agents/scratch/research/okf/` is gitignored and not committed. It holds the spec clone at `ad30107`, the saved changelogs and the lens notes, and it can be deleted any time.

## Environment State

- `main`, clean after this commit is pushed. No processes running.

## Artifacts to Read

- `.agents/docs/plans/2026-09-28-research-okf.md` — the record: approaches, tables, triggers, and found-along-the-way.
- `.agents/docs/changes/2026-09-28-research-okf.md` — what was dispatched and verified.

## Next Steps

1. **Operator, only if wanted:** the research record's found-along-the-way items. Each is a one-line decision.
   - **#1:** does `owner:` name the writer or the lane? This repo's deploy guides say `main`, the init-project templates say `documentation`, and `testing/guide.md` says `verifier`.
   - **#3:** unifi-home-network's config snapshot lives only in its 2026-09-14 handoff; fix it in that project.
   - **#4:** ranta-rankare's BUG-005 rule is half applied (`app.py:43-45`, `:251`); fix it in that project.
   - **#5:** valkompass's README names a deleted revision; defer it.
2. **Carried from the 2026-09-21 handoffs:**
   - the Mac canary check on Antigravity;
   - the devbox `claude --version` check and `sh-agent update`.
