---
owner: main
date: 2026-09-15
---
# Harness practices Phase 7, part 2: `/init-project` updates an older scaffold, and the Phase 7 review fixes     (2026-09-15)

## Summary

`/init-project` still scaffolds a new project, and it now also updates a project an older harness scaffolded: when `.agents/docs/changes/` has entries it adds only what is missing and never rewrites or deletes what the project wrote (plan D47, task 7.3). The Phase 7 code and docs reviews asked for changes: the installer's guidance for a hand-written `~/.claude/CLAUDE.md` now follows R2, the `--vendor all` docs name the command it applies to, and a mirror comment is corrected. Their three findings on the release gate's docs routing are the third round on that mechanism, so no patch was attempted: BUG-026 is reopened and the question is with the operator.

## Type

task

## Files changed

- `skills/init-project/SKILL.md` — the description, Usage and Step 1 name both modes. In an update, Step 1 lists what the run never edits (`REQUIREMENTS.md`, an existing `AGENTS.md` except the R2 fold, `DECISIONS.md`, existing deploy guides, trail files, `INDEX.md` rows, the project's rules, skills, agents, commands, workflows and scratch), shows a differing note under a title and takes the template's only on the user's yes, leaves out a file with uncommitted changes, and stops for a bare `docs/` trail. Step 3 fills only root files the run copied. Step 4 copies each missing template file instead of `cp -R`, and a `.gitkeep` only into an empty directory. Step 5 creates each missing `.claude` link instead of `ln -sfn` and reports any path it kept. Step 6 appends only missing `.gitignore` lines. Step 7 reports what an update added and kept, commits `init: update scaffold to current templates` in its own call, makes no commit when nothing changed, and asks before a push that deploys.
- `.agents/docs/testing/skill-matrix.md` — the `init-update` row: demo HEAD, an older scaffold; its check requires the added `.agents/docs/retros/.gitkeep`, an `init: update` commit, and no modified, deleted or renamed file except `.gitignore`.
- `internal/harness/apply.go` and `cmd/sh-agent/main.go` — the installer's `keep` line and `status`'s note tell the operator to move a hand-written `CLAUDE.md`'s content into a rule file under `rules/` and let install write the one-line import; the plain-copy comment says "image or VM build step". `internal/harness/harness_test.go` asserts the new guidance.
- `README.md` :32 and `.agents/docs/deploy/local.md` :12, :27 — `--vendor all` widens a single-vendor install through `install`, and `update` keeps the recorded vendor; the default is the vendor recorded by the existing install, else both.
- `.agents/docs/bugs/BUG-026-preflight-step4-routing-open.md` and `INDEX.md` — reopened with the three Phase 7 findings.
- `.agents/docs/reviews/2026-09-15-review-641eccd-worktree-{code,docs,merged}.md` — the two Phase 7 lens reports and the merged record.
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — D47, and task 7.3 extended to test both modes with an in-progress note.

## Agents dispatched

A workflow of nine agents surveyed seven older scaffolds read-only, designed an update path and critiqued it; the shorter rule shipped here keeps its survey evidence and drops the parts the critique broke. reviewer (code and docs lenses, Opus), documentation (the vendor wording), bugs (BUG-026 reopened). This entry was written by the main agent.

## Verification

- The Step 4 and Step 5 loops ran twice in scratch clones: on both sandbox demos the first run added only `.agents/docs/retros/.gitkeep` and the second nothing; on a bare clone at `8cbb128` the first run added all 18 template files and the four `.claude` links and the second nothing.
- The `init-update` check fails on a fresh demo clone, passes after a retros-only `init: update` commit, and fails when that commit also edits `TODO.md`. The code reviewer ran the snippets and the check again in its own scratch clones with the same results.
- `go test -race ./...` exit 0, `go vet ./...` exit 0, `gofmt -l` clean, a scratch install with `status` drift none and a `binary` line with no warning, `render` 93 lines, `sh-agent list` 27 skills and 7 agents.

Not run yet: the `init-update` and `init-deploy` rows on both CLIs, which need the sandbox after the Phase 6 re-runs.

## Open

- The release gate's docs routing (BUG-026): the operator is asked what the mechanism is for; the recommendation is to reuse Step 2's coverage rule for Step 4 in place of the rebuttal carry-forward.

## Requirements

n/a
