---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — release-binary: Phases 2 and 3 done, Phase 4 next

Session ID: 28ca24dd-a2d8-4278-85aa-70e24bbed445
Project: /home/kesokaj/workspace/styrman-harness
Branch: release-binary at c7da4fb (pushed, clean; `main` untouched at 22b92ff)

## Context

`.agents/docs/plans/2026-09-17-release-binary.md`: install the harness from a downloaded release binary with the files built in, and stop `sh-agent report` keeping a clone of the feedback repository on every machine. Phase 2 (the embedded source) and Phase 3 (`report` keeps nothing) are complete through their checkpoints. Phase 4 — `update` downloading a release and verifying its signed checksums — has not started.

## Completed

- Phase 2, commits 280dca7, 3f27ce2, 21ba8dc, 73f2c36: `Source{FS, Dir, Version}` for every reader, the embedded plan proven equal to the checkout plan, `ResolveSource` with the release-only embedded fallback, `embedded:<version>` in the manifest, `harness_version` in the report record. Changes entry `.agents/docs/changes/2026-09-17-release-binary-phase-2.md`.
- Phase 3, commits 0eec14a, 83f3cf6, dab2565: the throwaway blob-less sparse clone with removal on every path out, one retry on a rejected push, `RemoveKeptClone` in `update`, `.agents/docs/deploy/local.md` § Report back rewritten. Changes entry `.agents/docs/changes/2026-09-17-release-binary-phase-3.md`.
- One real `sh-agent report` from this box, authorised by the operator (D25), settled the last Phase 3 assumption: a `--filter=blob:none --no-checkout --depth 1` clone over SSH works against GitHub and a push from it is accepted (`09586ff`), leaving nothing in `~/.config/sh-agent/`. Log: `.agents/docs/testing/2026-09-17-release-binary-phase-3-real-report.md`.
- Both checkpoints: verifier PASS both times; `reviewer` APPROVE on Phase 2 with no findings; `security` on both phases — 1 medium + 6 low, then 2 medium + 6 low. Eight findings fixed with tests seen failing first, eight rebutted in writing in each report's `## Outcome`.
- BUG-066 to BUG-071 recorded with `INDEX.md` rows. D19–D25 added to the plan; D9 amended for BUG-069 and BUG-070.

## Open / In Progress

- The plan's first unticked box is Phase 4's first: `internal/harness/release.go` — latest tag via `gh`, download, `ssh-keygen -Y verify` of `SHA256SUMS.sig`, the per-binary checksum, rename into place, then the new binary's `install`. Not started, no blocker. The five cases the plan's Verification names (good signature, bad signature, bad checksum, missing asset, missing `gh`) drive it, against a fake `gh` on PATH and a test key.
- Phase 4's second box: the signer line compiled in from `.agents/docs/deploy/remote.md` § Releases' allowed-signers line, with a unit test that the constant equals that line.
- `update` on an embedded install currently refuses with a placeholder message (D22); Phase 4 replaces it with the release flow. It calls no git helper in that state.
- Requirements edits awaiting approval: none. The R8 sentence and the R9 bullets are approved as quoted in the plan's Interfaces and land in Phase 5's single commit to `main` (D11).
- Phase 5 carries three boxes added this session: `/preflight` must refuse a release when anything is untracked or ignored under the five embedded directories (D23, from BUG-066); the `AGENTS.md` `govulncheck` note names `LoadSkills`, whose call path now goes through `os.DirFS`; and the README/map/AGENTS/REQUIREMENTS/delivery-spec edits that still describe a kept clone and the unpacked bundle.
- Branch `close-map-gaps`: still waiting for another place's queued proposal; its ADR needs number 034 and its reading clone moves to `~/.config/sh-agent/feedback-read` (D12).
- `/feedback` on this box: still to re-try after `sh-agent update` from `main` (BUG-065's check) — switch this checkout to `main` first.

## Key Decisions

- D1–D25 in the plan. New this session: D20 (`Source.Version`), D21 (`status` reads a release binary's version from the module version, because `-trimpath` drops `-ldflags` from the build info — measured), D22 (`update` refuses on an embedded source until Phase 4), D23 (a release is built only from a tree with nothing untracked or ignored under the five directories), D24 (a release binary's install REPLACES a checkout install of the same home — the operator's choice), D25 (the one authorised real report run).
- Rejected: hiding `Source.Dir` behind an accessor; making an empty-but-set `--repo`/`$SH_AGENT_REPO` an error (pre-existing behaviour the tests rely on); defending against a pid collision across pid namespaces. All three are argued in the security reports' `## Outcome`.

## Corrections

- none logged this session (`~/.config/sh-agent/feedback.jsonl` has no entry for it).

## Agent notes (unverified)

- The operator was asked whether to continue into Phase 4 or stop; at the time this handoff was written they had not answered. Phase 4 must not begin without that answer.
- Phase 4 is the signature-verification phase. Two things measured this session bear on it: `-trimpath` keeps `-ldflags` out of the build info, so a release binary's version is readable only from the module version or from the `-X` variable at runtime (D21); and `ssh-keygen -Y verify` has not been checked on this box yet — it is the first task of Phase 4 in the plan's Risks.
- `.agents/scratch/` holds this session's probes: `p3-probe/` (the git mechanics, including `probe.sh`), `p3-mustfail/mutate.py` (the mutation harness that asserts a unique target and restores from a saved copy), `sec-review/` (7 MB from the Phase 2 security agent, including a nested probe repository), and several scratch homes. All gitignored, all deletable.

## Environment State

- Branch `release-binary` at c7da4fb, clean, pushed. Only tag is `v1.0.0`.
- This box's installed harness is still ce8ed6e — do NOT run `sh-agent update` while on this branch.
- `~/.config/sh-agent/feedback` still exists here, deliberately: removing it is `sh-agent update`'s job in the operator's own home, not this session's.
- No running processes.

## Artifacts to Read

- `.agents/docs/plans/2026-09-17-release-binary.md` — Interfaces (§ `update` on a binary-only place and § Release assets are Phase 4's spec), Failure modes, Verification, D1–D25, the unticked boxes.
- `.agents/docs/changes/2026-09-17-release-binary-phase-2.md` and `-phase-3.md` — what landed and how each guard was seen failing.
- `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md` and `-73f2c36-83f3cf6.md` — findings and their `## Outcome`.
- `.agents/docs/bugs/BUG-066` … `BUG-071` — the rules that apply to Phase 4: assert both directions of a generated artifact, grep every reader of a field that gains a new kind of value, a destructive guard keeps data on error, place a side effect below a read-only mode's returns, measure a test's premise, and make a mutation's target unique.
- `.agents/docs/deploy/remote.md` § Releases — the allowed-signers line Phase 4 compiles in.

## Next Steps

1. Ask the operator to confirm Phase 4 before starting it, if they have not answered.
2. Phase 4, first task: `ssh-keygen -Y 2>&1 | head -1` on this box (the plan's Risks asks for it first), then `internal/harness/release.go` with the five Verification cases against a fake `gh` and a test key.
3. Phase 4, second task: the compiled-in signer line with its unit test; then the checkpoint (verifier, `security` over signature verification and binary replacement, changes entry, commit, push).
4. Phase 5 as the plan orders it, landing on `main` as ONE commit with code, docs, map, R8/R9 and ADR-033 (D11, R9), then `/preflight v1.1.0` and the work Mac e2e.
