---
owner: changelog
date: 2026-09-30
---
# Retro: apply seven proposals to project and harness rules   (2026-09-30)

## Summary
Applied all seven proposals from `.agents/docs/retros/2026-09-30-retro.md` (findings 1–5 and 7; finding 6 stays queued). Two project-scope rules (vendor pinning verification, redaction by kind) added to `AGENTS.md` § Rules for working here; five harness-scope rules added to skills and verifier agent for Go file globbing, class-fix grep scope, consolidation order checking, frozen-tree semantics, and feature-capability consumer audit. Reviewer found two low-severity length contradictions in feature and bug skills (misaligned step-level and capability-level sentence bounds); both fixed before commit. Verifier behavioural run confirms the new verifier sentence correctly classifies foreign artifacts: 9 of 9 `new` runs and 1 of 8 `control` runs (control only when the global frozen-tree rule was readable in round 3) call the tree frozen.

## Type
task

## Files changed
- AGENTS.md — two new project rules (BUG-111 vendor pin verification, BUG-110 redaction by kind)
- agents/verifier.md — tree-under-test scoped to source/tests/config; treats docs-root and scratch artifacts as seen (finding 5, no bug record)
- skills/bug/SKILL.md — class fix implies `git grep` the same call everywhere (BUG-109)
- skills/feature/SKILL.md — capability consumer audit scope clarified to numbered list (operator 2026-09-29 correction)
- skills/golang/SKILL.md — operator paths never passed to `filepath.Glob` (BUG-109)
- skills/review-standards/SKILL.md — order unobservable by test gets source pin; consolidation checked against each original's step order (BUG-112)

## Agents dispatched
reviewer (code), reviewer (docs), verifier (`testing/2026-09-30-verify-retro-apply.md`, PASS; the behavioural run is the main agent's kit), changelog, bugs (graduating BUG-109, BUG-110, BUG-111, BUG-112 from records to rules)

## Verification
- **Code review:** `.agents/docs/reviews/2026-09-30-review-worktree-retro-code.md` APPROVE WITH NITS; nits (feature consumer list form, bug class-fix naming) fixed at the locations named
- **Docs review:** `.agents/docs/reviews/2026-09-30-review-worktree-retro-docs.md` APPROVE
- **Verifier behavioural:** `.agents/docs/testing/2026-09-30-verifier-foreign-artifact-behavioural.md` — three rounds, 17 total runs (9 new definition, 8 control); new definition 9 of 9 classify foreign artifact as seen and 0 of 9 call tree not frozen; control 1 of 8 call tree not frozen (only when global frozen-tree rule was readable in round 3)
- **Compile and test:** `go test -race ./...` pass; `TestRenderedContextFitsAntigravityBlock` pass (context-size test); `go run ./cmd/sh-agent install --home .agents/scratch/x` succeeds, `status --home .agents/scratch/x` reports no drift
- **Headless skill load:** `claude -p --setting-sources project,local` started in scratch fixture; all four changed skills listed, no parse errors; transcript confirmed each skill definition was consulted

## Docs
No docs changed — proposal application is a rules graduation, not a command, interface or deploy-step change.

## Requirements
n/a

## Notes
Verifier's next run (on Antigravity per `AGENTS.md` § Release) and one run each of the four changed skills owed before the next release cuts (v1.13.2 or later). Every place picks up the rules with `sh-agent update`.
