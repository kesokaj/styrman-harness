---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — release-binary: Phase 2 done, Phase 3 next

Session ID: 28ca24dd-a2d8-4278-85aa-70e24bbed445
Project: /home/kesokaj/workspace/styrman-harness
Branch: release-binary (pushed, clean; `main` untouched at 22b92ff)

## Context

Continuing `.agents/docs/plans/2026-09-17-release-binary.md`: a place should be able to install the harness from a downloaded release binary with the files built in, and `sh-agent report` should stop keeping a clone of the feedback repository on every machine. Phase 2, the embedded source, is complete through its checkpoint.

## Completed

- 280dca7 — every reader takes a `Source{FS, Dir, Version}` (`internal/harness/source.go`); `inventory.go`, `render.go`, `plan.go` read through `io/fs`; BUG-036's link refusals run for `Dir != ""` and are unchanged in meaning; `internal/harness/source_test.go` (MapFS plan, nine refusals, a linked `skills/<name>` followed, a dangling one refused).
- 3f27ce2 — `plan_test.go` `TestEmbeddedPlanEqualsCheckout`: the embedded plan equals the checkout plan in paths, bytes, modes and KeepForeign, all three vendors.
- 21ba8dc — `ResolveSource` replaces `FindRepo` with the release-only embedded fallback (D14, D15); `Install` and `BuildRecord` take a `Source`; `state.json` records `embedded:vX.Y.Z`; `status` prints `source     embedded vX.Y.Z`; `adopt` and `update` refuse; `harness_version` added to the report record; D20–D22.
- 73f2c36 — checkpoint: the three security fixes (see below), verifier/reviewer/security artifacts, changes entry `.agents/docs/changes/2026-09-17-release-binary-phase-2.md`.
- fef34dd — BUG-066, BUG-067, BUG-068 with `INDEX.md` rows; D24.
- Phase 2's four task boxes and its Checkpoint box are ticked with evidence in the plan.

## Open / In Progress

- The plan's first unticked box is Phase 3's first: `report`'s throwaway partial clone (sparse `reports/<host>/`, deferred removal, retry once on rejection, `--feedback-repo` unchanged) — not started, no blocker. Phase 3 then rewrites the four clone-pinning tests in `internal/harness/report_test.go` (each must-fail against HEAD's `report.go`) and adds the old-kept-clone cleanup to `update`, then a checkpoint with a `security` dispatch.
- Phase 3's first task carries an assumption to settle before it is called done (plan § Risks): GitHub must serve `--filter=blob:none` over SSH and accept a push from a `--depth 1` clone. Against local bare repositories only, until the operator agrees to one real `report` from this box.
- Requirements edits awaiting approval: none — the R8 sentence and the R9 bullets are approved as quoted in the plan's Interfaces and are applied in Phase 5's single commit to `main` (D11).
- Branch `close-map-gaps`: still waiting for another place's queued proposal; its planned ADR needs number 034, and its reading clone moves to `~/.config/sh-agent/feedback-read` (D12).
- `/feedback` on this box: still to re-try after `sh-agent update` from `main` (BUG-065's interactive check) — switch this checkout to `main` first.

## Key Decisions

- D1–D24 in the plan. New this session: **D20** `Source` carries `Version`, set only for the embedded files; **D21** `status` reads an installed release binary's version from its build info's module version, because `-trimpath` keeps `-ldflags` out of the build info (measured with go1.27.1); **D22** `update` refuses on an embedded source until Phase 4 replaces it with the release flow; **D23** a release must be built from a tree with nothing untracked or ignored under the five embedded directories — pinned by a test now, refused by `/preflight` in Phase 5; **D24** a release binary's install replaces a checkout install of the same home, chosen by the operator ("Replace, as it does now").
- Rejected: hiding `Source.Dir` behind an accessor to make the link refusals harder to bypass — `DirSource` is the only constructor and both fields are documented; it buys nothing against a caller that ignores the constructor (security finding 7, rebutted in that report's `## Outcome`).

## Corrections

- none logged this session (`~/.config/sh-agent/feedback.jsonl` has no entry for it).

## Agent notes (unverified)

- The security review's remaining low findings are rebutted in writing in its `## Outcome`; the empty-but-set `--repo ""` / `SH_AGENT_REPO=` fall-through is pre-existing behaviour the harness's own tests rely on, and changing it would be its own task.
- Dispatch rules that held again: a tool-call budget plus "write the artifact early" in every dispatch — all three agents had their file on disk. The `bugs` agent's two drafts needed correction against the diff (a repro credited to the wrong probe, and the wrong mutation cited), so read an agent-written record before committing it (BUG-061).
- `.agents/scratch/sec-review/` holds the security agent's 7 MB of probes, including a nested probe git repository. Gitignored, harmless, deletable.

## Environment State

- Branch `release-binary` at fef34dd, clean, pushed; only tag is `v1.0.0`.
- No running processes. This box's installed harness is still ce8ed6e — do NOT run `sh-agent update` while on this branch.
- Scratch homes from this session: `.agents/scratch/{source-home,embed-home,embed-cmp,embed-bin,verify-phase2-*}` (gitignored).

## Artifacts to Read

- `.agents/docs/plans/2026-09-17-release-binary.md` — Interfaces (§ Report throwaway clone is Phase 3's spec), Failure modes, Verification, D1–D24, the unticked boxes.
- `.agents/docs/changes/2026-09-17-release-binary-phase-2.md` — what Phase 2 changed and how each guard was seen failing.
- `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md` — findings and the `## Outcome` (three fixed, five rebutted, one decided by the operator).
- `.agents/docs/reviews/2026-09-17-review-f7ee817-21ba8dc-code.md` — code review, APPROVE, 0 findings.
- `.agents/docs/testing/2026-09-17-release-binary-phase-2.md` — verifier log, and what it did not exercise.
- `.agents/docs/bugs/BUG-066…068` — the embed/manifest/mutation rules to apply in Phase 3.

## Next Steps

1. Phase 3, task 1: the throwaway partial clone in `internal/harness/report.go` (`ResolveFeedback`'s clone branch, sparse `reports/<host>/`, `os.RemoveAll` deferred on every path, one retry from a fresh clone on a rejected push), against local bare repositories holding two hosts' reports — the other host's blobs must never be in the clone, checked through a test hook before removal.
2. Phase 3, tasks 2–3: the four clone-pinning tests rewritten, each seen failing against HEAD's `report.go` first; the old-kept-clone cleanup in `update` (clean, origin-verified, not ahead → removed and named; otherwise left and named with the reason).
3. Phase 3 checkpoint: verifier, `security` over the phase (a credential-bearing directory and pushes to a remote), changes entry, commit, push.
4. Operator decision to ask before Phase 3 is done: whether one real `sh-agent report` from this box into `styrman-harness-feedback` may run, to settle the partial-clone-over-SSH assumption.
