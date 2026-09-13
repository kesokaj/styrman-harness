---
owner: main
date: 2026-09-13
---
# Handoff — 2026-09-13 — learning loop at project level shipped; improvement review awaiting the operator's go/no-go

Session ID: session_01CtMYM9AGp33edxBQGURNaB
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean, pushed at 5cdd4d6 (plus this handoff's commit)

## Context

The operator asked for a research pass over the harness, then: fix what it found, make `/retro` and `/distill` record at project level with per-project logs a later harness session can use, make every skill model-invocable, and then review the harness for improvements ("the agent storing files under .agent/docs is amazing"). Plan: `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` (Phases 1–5; all boxes ticked except the operator's own step).

## Completed

- Phase 1 (82fd4a6): `sh-agent install` rebuilds the binary from the checkout with `go build` and lets that build install; `status` prints a `binary` line and warns when the installed build is older than the repo or unstamped; `state.json` records `+dirty`; `cmd/sh-agent` has tests. BUG-006 (the real home ran the 2026-09-03 binary for ten days) is the reason.
- Phase 2 (bd5ed87): `disable-model-invocation` removed from `audit`, `distill`, `handoff`, `init-project`, `research`; safeguards in the descriptions; verified on both CLIs in the sandbox (all five listed, an unrelated prompt loads none).
- Phase 4 (fc7f195): `/tmp` gone from the repo docs (BUG-007); catch-up records for the four unrecorded retro commits (BUG-008); review rounds 1–3 on the checkout guard — a named path that is not a checkout is an error, `--dry-run` never builds, install from the crew mirror is refused; ADR-013, ADR-014; verifier PASS, reviewer APPROVE, security converged at low residuals (`changes/2026-09-13-review-rounds-2-3-checkout-guard.md`).
- Phase 3 (5cdd4d6): `skills/retro/SKILL.md` rewritten — scope by target, `rule_ref` vocabulary shared with `feedback`, project record in `.agents/docs/retros/`, harness queue in `~/.config/sh-agent/proposals/`, index in `~/.config/sh-agent/retros.jsonl`, `--harness` consumption; `distill` records the same way; ADR-015. Verified headless on both CLIs in fresh demo clones (records with `Scope:` tags, index lines, no parse errors).
- Phase 5: `.agents/docs/reviews/2026-09-13-improvement-review.md` — 20 numbered proposals from a 55-agent panel with adversarial verification, each marked `[fix]` or `[yours]`; presented to the operator in chat at the end of this session.

## Open / In Progress

- **Operator, this box and the Mac:** run `sh-agent update` (or `go run ./cmd/sh-agent install`) once — the installed 2026-09-03 binary cannot rebuild itself — then `sh-agent status` must show a `binary` line without a warning (`TODO.md`, high priority). Antigravity auto-updated to 1.2.2 during the session; the three smoke checks in `testing/2026-09-03-verify-real-home.md` have not been run on it from the real home.
- **Operator:** the go/no-go per item on the improvement review (`reviews/2026-09-13-improvement-review.md`). Nothing from it is applied. Items 15, 18, 19, 20 are `[yours]`.
- `TODO.md` medium: upgrade Go to ≥ 1.24.6 on both machines (govulncheck); pin the git root in `Inspect` (nested mirror inherits the enclosing repo's commit).
- The downstream mirrors (styrman-crew `image/harness/`, argolis-devbox `deploy/harness/`) are seven commits behind this repo and carry the BUG-006 binary; refreshing them is item 5 of the review and each project's own task (`styrman-crew/image/harness/VENDORED.md` has the recipe).

## Key Decisions

- The installed binary never copies itself: every install is a fresh `go build` of the named checkout (ADR-013); a checkout found only by walking up from the working directory is refused when another is installed; `--repo` or `$SH_AGENT_REPO` must name a real checkout.
- Every harness skill is model-invocable (ADR-014); the lock is kept only for third-party skills that must never fire on their own.
- Retro scope is decided by the target file; project records are the project's state, the harness queue and index stay under `~/.config/sh-agent/` so R7 holds as written (ADR-015). The cold-reader check took three rounds (nine reports under `reviews/2026-09-13-review-cold-read-learning-loop-*.md`); the round-3 residue was resolved by the operator's "continue with phase 3, use your recommendation".
- Rejected: editing `REQUIREMENTS.md` (a clarifying sentence is offered as review item 19); raising the `go` directive in `go.mod` for the govulncheck advisories (it would make every install download a toolchain).

## Corrections

- None logged this session in `~/.config/sh-agent/feedback.jsonl`. The operator's two turns were instructions, not corrections: "Good, this is what I want adding…" and "continue with phase 3, use your recommendation".

## Agent notes (unverified)

- `go run ./cmd/sh-agent status` always prints `binary (devel) — no build stamp` because it describes the running process; review item 17 fixes that. From the installed binary the line is right.
- The 09-04 retro's proposal P6 (feedback skill should log mild procedure corrections) is still unapplied; it is review item 11.
- The retro e2e on Claude Code observed that the old free-text `rule_ref` values in `feedback.jsonl` never match the new vocabulary — by design (no migration), so promotions start from entries written after today.

## Environment State

- Branch: main, clean; origin at HEAD. Real home on this box: content matches the repo at 5cdd4d6 only after the operator's `sh-agent update`; today's sessions did not install into the real home (R6).
- Sandbox home `/home/kesokaj/sandbox/home` has the working tree installed at 5cdd4d6; throwaway clones under `/home/kesokaj/sandbox/runs/{claude,agy}-retro-p3` hold the e2e retro records; `/home/kesokaj/sandbox/home/.config/sh-agent/retros.jsonl` has the two e2e index lines.
- Running processes: none.

## Artifacts to Read

- `.agents/docs/reviews/2026-09-13-improvement-review.md` — the 20 proposals; start here if the operator has answered
- `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` — decisions D1–D12 and the ticked phases
- `.agents/docs/changes/2026-09-13-*.md` — five entries, newest `retro-at-project-level` and `review-rounds-2-3-checkout-guard`
- `.agents/docs/bugs/INDEX.md` — BUG-006..009 added today
- `.agents/docs/testing/2026-09-13-verify-research-pass.md` and `2026-09-13-verify-learning-loop-phases-1-2.md`

## Next Steps

1. Operator: `sh-agent update` on this box and the Mac; `sh-agent status` shows `binary <rev>` with no warning.
2. Operator: answer the improvement review per number (go / no-no / go all). Then, in this checkout: each accepted `[fix]` item is a task under the Definition of Done (most are one commit each; 13 and 14 are an hour), with one sandbox start of each CLI after any change to `context/`, `agents/`, or `skills/`.
3. After the first `/retro` in a real project under the new skill: check the record in that project's `.agents/docs/retros/` and the index line, then run `/retro --harness` here to consume any queued harness-scope proposals — the loop's first live run.
