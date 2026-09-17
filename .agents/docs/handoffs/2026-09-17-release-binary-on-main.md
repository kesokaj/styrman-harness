---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — release binary on `main`; only the release itself is left

Session ID: 28ca24dd-a2d8-4278-85aa-70e24bbed445
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at c9d24a3, clean and pushed. `release-binary` is kept at 24d08a8, merged and not deleted.

## Context

`.agents/docs/plans/2026-09-17-release-binary.md`: a place should be able to install the harness from one downloaded binary with the files built in, and `sh-agent report` should stop keeping a clone of the feedback repository on every machine. Phases 2 to 5 are done and landed on `main` as one commit, cd0bbcb, as R9 requires. What is left is the release: `/preflight v1.1.0`, which only the operator starts.

## Completed

- **cd0bbcb on `main`** — one commit, 54 files, squashing the branch's 22: the embedded source, `report`'s throwaway clone, `update` from a signed release, the redrawn map, README § Install without cloning, `REQUIREMENTS.md` R8 and R9 as approved, `AGENTS.md` rules 1 and 2, both deploy guides, the delivery spec, `/preflight` Step 6, ADR-033a, and BUG-066…072.
- Per-phase records: `.agents/docs/changes/2026-09-17-release-binary-phase-{2,3,4,5}.md`.
- Checkpoints: verifier PASS at each phase (Phase 4 twice — the first run is marked inconclusive, see BUG-072); `reviewer` APPROVE on Phase 2's code and REQUEST CHANGES on Phase 5's docs (both findings fixed); `security` on Phases 2, 3 and 4 — 1+2+1 medium and 6+6+6 low, fourteen fixed with tests seen failing first, the rest rebutted in writing in each report's `## Outcome`.
- One real `sh-agent report` from this box, authorised by the operator (D25), settled that GitHub serves a blob-less shallow clone over SSH and accepts a push from it (`09586ff`).
- Decisions D19–D28 recorded; D9 amended (BUG-069, BUG-070); D22 superseded by Phase 4.

## Open / In Progress

Two unticked boxes and three partial ones, all belonging to the release the operator runs:

- **E2E: `/preflight v1.1.0` from the operator's own turn**, then the work Mac golden path and edge in the plan's Verification. Not started by their choice ("Merge to main, stop before the release").
- **Checkpoint and Done** — the plan's last row, after the release.
- Partial (`[~]`), each owed before or with the release: the scratch REHEARSAL of the release procedure (build the three binaries, write and sign `SHA256SUMS`, run `UpdateFromRelease` against them through a fake `gh`); the dry build of v1.1.0's release notes for the operator to read; and one refused dry run of the D23 guard against a planted ignored file.
- **A release cut by the OLD procedure is uninstallable on purpose.** `update` refuses a release whose signed `SHA256SUMS` does not carry `# styrman-harness vX.Y.Z` (D26). `.agents/docs/deploy/remote.md` § Releases writes that line; anyone cutting v1.1.0 must follow that section, not memory.
- Requirements edits awaiting approval: none. R8 and R9 are applied as approved.
- Branch `close-map-gaps`: still waiting for another place's queued proposal; its ADR needs number 034 (033 and 033a are taken) and its reading clone moves to `~/.config/sh-agent/feedback-read` (D12).
- `/feedback` on this box: still to re-try after `sh-agent update` from `main` (BUG-065's check). This checkout is on `main` now, so that is possible — but `update` here also removes `~/.config/sh-agent/feedback`, the old kept clone, which is the intended behaviour (D9) and worth expecting rather than being surprised by.

## Key Decisions

D1–D28 in the plan; ADR-033 and ADR-033a in `DECISIONS.md`. The ones a fresh reader most needs: the embedded files are a fallback ONLY for a release binary (D14) and never after an explicit `--repo` (D15); a release binary's install replaces a checkout install of the same home (D24, the operator's choice); `update` never goes sideways or backwards and refuses a tag that is not `vX.Y.Z` (D27); the version lives inside the signed checksums (D26); the signer line is compiled in and the verifier is preferred from a system path (D28).

## Corrections

- none logged this session (`~/.config/sh-agent/feedback.jsonl` has no entry for it).

## Agent notes (unverified)

- The operator approved the map and the merge in their own turn and said "Approved" to the merge again afterwards. They have NOT approved cutting the release; do not run `/preflight` without a fresh instruction.
- Two hazards this session hit twice each, worth expecting: a generated Go string literal whose `\n` became real newlines (BUG-057 — write through a raw string and grep for the escape afterwards), and a test whose premise did not hold (BUG-071 — measure the setup, not just the assertion). A third: dispatching the verifier and then editing the tree it was reading (BUG-072).
- `.agents/scratch/` holds this session's probes and scratch homes: `p3-probe/` (git mechanics), `p3-mustfail/` and `p4-mustfail/mutate.py` (the mutation harnesses), `p4-probe/` (ssh-keygen signing), `sec-review/` (7 MB from a security agent, with a nested probe repository), and several install homes. All gitignored and deletable.

## Environment State

- `main` at c9d24a3, clean, pushed; only tag is `v1.0.0`. `release-binary` kept at 24d08a8.
- This box's INSTALLED harness is still ce8ed6e — older than `main` now. `sh-agent update` here would move it to `main` and remove the old kept feedback clone.
- Suite: 99 top-level tests, 85 subtests, none failing or skipped; `gofmt`, `go vet`, `go test -race`, `govulncheck` all clean on go1.27.1.

## Artifacts to Read

- `.agents/docs/plans/2026-09-17-release-binary.md` — the five partial/unticked boxes, D1–D28, § Verification's e2e criteria.
- `.agents/docs/deploy/remote.md` § Releases — the exact release procedure, including the version marker and the pre-build guard; `/preflight` follows it.
- `.agents/docs/changes/2026-09-17-release-binary-phase-{2,3,4,5}.md` — what landed, and how each guard was seen failing.
- `.agents/docs/reviews/2026-09-17-security-{f7ee817-21ba8dc,73f2c36-83f3cf6,dab2565-4d45155}.md` and `-review-edbddca-1a5578b-docs.md` — findings and outcomes.
- `.agents/docs/bugs/BUG-066` … `BUG-072` — the rules earned today.

## Next Steps

1. Nothing is owed until the operator wants the release. When they do: rehearse the procedure in `.agents/scratch/` (three builds, marker, signature, `UpdateFromRelease` through a fake `gh`), build v1.1.0's notes for them to read, then `/preflight v1.1.0` from their own turn, following `.agents/docs/deploy/remote.md` § Releases exactly.
2. After the release: the work Mac e2e from the plan's Verification (the `gh` one-liner, `sh-agent install`, `status` showing `source embedded v1.1.0`, `report` leaving nothing), then the plan's Checkpoint and Done row.
3. Independent of the release: re-try `/feedback` here after `sh-agent update` (BUG-065), and pick up `close-map-gaps` when the other place's proposal arrives.
