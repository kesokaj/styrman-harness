---
owner: main
date: 2026-09-21
---
# Handoff — 2026-09-21 — context consolidation done; v1.5.0 and the Mac smoke are the operator's

Session ID: 22858f28-e5eb-4616-977a-d19bb810a17f        (same machine: `claude --resume 22858f28-e5eb-4616-977a-d19bb810a17f`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 470162b (this handoff's commit follows it), clean and pushed. Tags `v1.0.0` … `v1.4.0`; v1.5.0 not yet cut.

## Context

The operator handed over an "agent context audit" prompt and asked whether it would help. Measuring the harness against it found a live defect: the Antigravity CLI loads only the first 24,000 bytes of the deployed global context (27,114 bytes), silently, in every release since v1.2.0 (BUG-085). The operator approved the assessment's recommendations, then the consolidation Feature that fixes it: prune the context under a size-pinned budget, prove the end of the file with a canary, stop the regrowth, goldfish the always-on files.

## Completed

- Assessment: `.agents/docs/plans/2026-09-21-context-audit-assessment.md` (the prompt beside it, `2026-09-21-context-audit-prompt.md`); the defect probed and recorded; REQUIREMENTS § Known limits bullet (operator-approved) and ADR-035.
- Plan `.agents/docs/plans/2026-09-21-context-consolidation.md`, all boxes closed (Phases 0–4); cold read three rounds, cap on one item, operator's "approved".
- Context 27,163 → 20,205 bytes (`context/AGENTS.md`, `context/dev-setup.md`; 44 operator-approved actions, `plans/2026-09-21-context-pruning-review.md`); pinned by `internal/harness/size_test.go` (22,000-byte budget, must-fail mutation verified); canary last in the rendered context (`internal/harness/render.go`); smoke docs, README, guide, this repo's `AGENTS.md` updated.
- Regrowth stop: `skills/retro/SKILL.md` (delete/merge row, size check at apply), `agents/bugs.md` (a home in the context names what pays for it) — 3 of 3 behavioural runs.
- Three procedures moved into `skills/review-standards` § 2 and `skills/golang`; goldfish gaps closed in `context/`, `skills/cloud-run-deploy`, project `AGENTS.md`.
- Real home on this box at 94a6596 (`sh-agent update`, on the operator's word), three smoke checks PASS on both CLIs, Antigravity's store shows zero truncation markers (`testing/2026-09-21-verify-real-home-smoke-5.md`). BUG-085 fixed; BUG-086, BUG-087 recorded (fixed).
- Changes entries: `2026-09-21-context-audit-assessment.md`, `-context-consolidation-phase-0.md` … `-phase-4.md`.

## Open / In Progress

- Nothing in progress; `ROADMAP.md` "Now" and `TODO.md` are empty.
- Two operator actions (Next Steps): the release, and the Mac smoke after its `sh-agent update`.
- Follow-ups from the plan, not scheduled: whether a second Antigravity rules file loads when the block is under the cap (needs a probe with a small context; test before any `rules/` file); the retro table's `rules/` row lands with the first rule file; a `status` line for an over-cap context; the project `AGENTS.md`'s two narrative spans the MINIMAL check named (the ADR-028/029 toolchain history, the v1.3.0 anecdote in rule 6); BUG-011's index row still sits in the main table though its rule graduated (found by the behavioural runs; `bugs` moves it when a graduating commit next touches that bullet).
- Requirements edits awaiting approval: none.

## Key Decisions

- Plan D1–D17 (`plans/2026-09-21-context-consolidation.md` § Decisions): budget 22,000 bytes (observed cap 24,000); canary emitted last by `render.go`; pruning approved by number; no `rules/` file in this plan (the cap is per block — a 45-byte second file was dropped and cost the first one 342 bytes, so a split fixes nothing); the size test measures context plus every `rules/*.md`; the retro size check sits at the apply step; the bugs rule needs no tool the agent lacks; MOVE-TO-SCOPED rows kept in shortest form; v1.5.0 is a MINOR (the retro skill gains behaviour, every user's context changes).
- Rejected: pinning the documented 12,000 characters; dropping the inventory; an install-time size check; the audit prompt's enforcement layer (R3, ADR-007), Cursor/Codex rows, goldfish over all skills, Step 7 experiment.
- Versioning answer given: one version for the whole harness — the signed tag; context and binary ship together; checkout places follow `main`, binary-only places move at the next tag.

## Corrections

- None this session. The one `feedback.jsonl` entry dated 2026-09-21 (10:20) belongs to another project's session (`ranta-rankare`, a deploy-ordering correction), not to this one.

## Agent notes (unverified)

- The real `~/.gemini/antigravity-cli/settings.json` today holds only `trustedWorkspaces`, no `permissions.allow`; the 2026-09-16 log recorded six narrow entries then. Headless Antigravity checks that need tools were run with the per-invocation `--dangerously-skip-permissions` inside throwaway clones; the settings file was not modified.
- The hourly `sh-agent-update.timer` refuses a dirty checkout (by design, BUG-046); the 13:01 run failed for that reason. Leave the tree clean when a timer run should move the real home.
- Parallel shell calls share one working directory: a `cd` in one call was overtaken by another, so two smoke checks ran in the wrong clone once (recorded in `testing/2026-09-21-verify-sandbox-pruned-context.md`). Wrap each CLI run in its own subshell.
- The Antigravity cap number comes from `<truncated N bytes>` markers in its conversation store, twice at exactly 24,000; the vendor page says 12,000 characters. A later CLI may change it; the size test pins 22,000.

## Environment State

- Branch: `main`, clean; 470162b pushed. Real home on this box: 94a6596, drift none (the next timer run will move it to this handoff's commit, docs only).
- Running processes: none. Sandbox home `.agents/scratch/sandbox-home` remains (gitignored) without credentials; scratch clones under `.agents/scratch/sandbox-runs/` and `.agents/scratch/realhome-smoke/smoke5-*`.
- Mac and devbox: last known v1.3.0 / v1.4.0; both truncate the context on Antigravity until they update to v1.5.0.

## Artifacts to Read

- `.agents/docs/plans/2026-09-21-context-consolidation.md` — the closed plan, D1–D17, every box with its evidence.
- `.agents/docs/changes/2026-09-21-context-consolidation-phase-4.md` — the feature totals and the last phase.
- `.agents/docs/bugs/BUG-085-antigravity-global-context-truncated-24000-bytes.md` — the defect, the probe, the fix.
- `.agents/docs/plans/2026-09-21-context-audit-assessment.md` — why the audit prompt was adopted in part and what was struck.
- `.agents/docs/testing/2026-09-21-verify-real-home-smoke-5.md` — the e2e proof on the real home.

## Next Steps

1. **Operator:** `/preflight v1.5.0` from your own turn (MINOR: the context changes for every user and the retro skill gains behaviour). The gate runs over `v1.4.0..HEAD`.
2. **Operator, on the Mac:** `sh-agent update` once v1.5.0 is published, then from the real home ask Antigravity, with no tools, what word the global context ends with — expect GANNET — and `strings` the newest `~/.gemini/antigravity-cli/conversations/*.db` for no `<truncated` line inside the `<RULE[...AGENTS.md]>` block. The devbox the same.
3. Next session, only if wanted: the follow-ups above, starting with the second-rules-file probe before any `rules/` file.
