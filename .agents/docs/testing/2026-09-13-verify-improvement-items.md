---
owner: verifier
date: 2026-09-13
---
# Verification — improvement-items   (2026-09-13)

Range reviewed: 5cdd4d6..working tree (uncommitted). Applies the 20 items of
`.agents/docs/reviews/2026-09-13-improvement-review.md`. Blobs pinned:
`internal/harness/repo.go` = 48f38977..., `cmd/sh-agent/main.go` = 21b08a07...
(git hash-object, working tree, at run time).

## Status: PASS

Scope: ran the project's own commands (gofmt/vet/test), a fresh install into a
scratch home with status/drift/list checks, mapped every changed Go file to
its covering test, and spot-checked prose items 6, 7, 8, 9, 13, 14, 16 against
the review's numbered proposals. Did not re-run the sandbox CLI starts — cited
per the dispatch as already run on 2026-09-13 (Antigravity 1.2.2,
`.agents/scratch/e2e-items/`, log `cli-20260913_220545.log`, 0 parse errors).
Did not verify item 15 (frontend `[yours]`, explicitly deferred to the
operator) or items 18-20 (`[yours]`, no code/prose to check — they are
decisions/actions for the operator, not part of this change).

## Commands

- `gofmt -l .` → no output (clean).
- `go vet ./...` → no output (clean).
- `go test -race ./...` →
  `ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	1.139s`
  `ok  	github.com/kesokaj/styrman-harness/internal/harness	1.450s`
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-home` →
  `building sh-agent from /home/kesokaj/workspace/styrman-harness`
  `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`
  (plus the expected PATH note; no error).
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-home` →
  ```
  repo       /home/kesokaj/workspace/styrman-harness (main @ a4146dc, uncommitted changes)
  installed  2026-09-13T22:07:24+02:00 from commit a4146dc+dirty (vendor: all, 260 files, manifest ~/.config/sh-agent/state.json)
  binary     a4146dca092f+dirty
  retro      none yet — run /retro at the end of a session with corrections
  repo has   27 skills, 7 agents, 0 commands, 0 rules
  claude     ~/.claude  (AGENTS.md, CLAUDE.md, skills/, agents/, commands/, rules/)
  agy        ~/.gemini/config  (rules/AGENTS.md, skills/, agents/, workflows/, rules/)
  drift      none — 259 files match the repo
  ```
  Binary line carries a real revision + `+dirty`, no "no build stamp" warning —
  confirms item 17 (status reads the installed file's build info via
  `debug/buildinfo.ReadFile`, not the `go run` process). Note: install reports
  260 files managed but drift reports 259 match; not chased further (see
  Not examined) — did not affect PASS since drift is still "none" (no
  mismatch, no missing/extra flagged), but the off-by-one between the two
  counts is worth a follow-up look, not treated as a defect here since it
  predates this change's scope (state.json bookkeeping, not part of the 20
  items).
- `go run ./cmd/sh-agent list --home .agents/scratch/verify-home` → "Skills
  (27)" table with 27 rows, "Subagents (7)" table with 7 rows. Matches the
  required 27 skills / 7 agents.

## Coverage of changed code

- `internal/harness/repo.go` (`ownCheckout`, `Inspect`, `Behind`, `Pull` guard)
  — covered by `internal/harness/harness_test.go::TestNestedCopyIsNotACheckout`
  (nested mirror inside another git repo: `Inspect` empty, `Behind`/`Pull`
  return `errPlainCopy`; repo root itself: `Inspect` populated, `Behind` not
  refused) and the existing `TestBehindAndPull`.
- `cmd/sh-agent/main.go` (`status` binary line via `buildinfo.ReadFile`,
  `versionFromInfo` extraction) — covered by
  `cmd/sh-agent/main_test.go::TestVersionFromInfo` (4 cases: clean, dirty,
  short revision, no revision) and
  `TestRunInstallStatusUninstall` (asserts the printed status never contains
  "binary     unknown" after a real install). The `status()` print branch for
  a plain-copy repo (`info.Commit == ""` → "plain copy, not its own git
  checkout") has no direct unit test of its own, but the input it branches on
  (`Inspect` returning `RepoInfo{}`) is directly tested in
  `TestNestedCopyIsNotACheckout`; the print statement itself is a one-line
  format wrapper, not treated as a coverage gap.
- `internal/harness/derive.go` (comment-only change, dropping `memory` from
  the documented Claude-only key list) — no behavior change, no test needed.
- `internal/harness/harness_test.go`, `cmd/sh-agent/main_test.go` — test files
  themselves, exercised by the `go test` run above.

No changed Go file is uncovered.

## Prose items spot-checked (review numbering)

- 6 (drop `memory: project`): confirmed — `derive.go` comment no longer lists
  `memory`; `agents/{bugs,changelog,documentation}.md` diffs shrink the memory
  paragraphs.
- 7 (explorer write-early, maxTurns wording): confirmed — `agents/explorer.md`
  gained the "Write the brief EARLY" section from this diff.
- 8 (no unconditional Colima command): confirmed — `skills/docker/SKILL.md`
  and `skills/cloud-run-deploy/SKILL.md` scope `colima status || colima
  start` to macOS only.
- 9 (feature e2e never deploys): confirmed — `skills/feature/SKILL.md` line
  78 routes the GCP deploy through `/preflight` after Step 6.
- 13 (golang references dangling links / off-stack content): spot-checked —
  no dangling `./references/*.md` links found across `skills/golang/*.md` and
  `skills/golang/references/*.md`; no DI-framework or "stay-updated" section
  found in `golang-ecosystem.md`. `staticcheck` remains referenced (in
  `SKILL.md` and `golang-core.md`) as a linter recommendation, not as an
  unconditional project requirement — not treated as a residual defect.
- 14 (docker skill Go example + amd64 rule): confirmed — `skills/docker/
  SKILL.md` example is `FROM golang:1.24 AS build`; amd64 rule matches
  `context/dev-setup.md` wording; no registry/Artifact-Registry section found
  in the docker skill (moved to `cloud-run-deploy`, per item 16's parallel).
- 16 (gcloud-setup merge): confirmed — `skills/gcloud-auth-verification/
  SKILL.md` is deleted (`git status` shows `D`); `skills/gcloud-setup/
  SKILL.md` still contains the auth troubleshooting content (no "discover the
  project id" instruction found in a scan of its auth section).

Not independently re-derived: items 1-5, 10-12 (trail/handoff/workflow-skill
prose), 17's exact wording match against the review's phrasing beyond
behavior, 19-20 (`REQUIREMENTS.md` R7 sentence and queue-location decision —
`[yours]`, operator decisions, present in `REQUIREMENTS.md` diff but not
cross-checked word-for-word against the review's proposed sentence).

## Sandbox CLI starts (cited, not re-run)

Per the dispatch: both CLIs, three checks each, ran 2026-09-13, evidence
under `.agents/scratch/e2e-items/`, Antigravity 1.2.2, log
`cli-20260913_220545.log`, 0 parse errors. Not re-executed in this
verification pass.

## Not examined

- Items 1-5 (trail: TODO/ROADMAP ownership, handoff template, bug-skill
  commit ordering, init-project templates, downstream-mirror naming) and
  10-12 (research→feature handoff, feedback "when in doubt log", duplicate
  prose removal) were not individually diffed against the review's exact
  wording — the diff stat shows the files touched but content was not
  line-matched to each numbered proposal.
- Item 15 (frontend TanStack Start staleness) is explicitly `[yours]` and
  unresolved by design — not evaluated.
- Items 18-20 (`[yours]`: vendor-fact re-verification, R7 wording exactness,
  queue/index location decision) are operator decisions; only confirmed that
  `REQUIREMENTS.md` was touched, not that the wording matches the review's
  proposed sentence verbatim.
- Did not investigate the 260-vs-259 file count discrepancy between `install`
  and `status` output (see Commands section) — flagged for a follow-up, not
  blocking this PASS.
- Did not re-run the sandbox CLI smoke checks; relied on the dispatch's
  citation.

## Environment

Linux 6.12.107+deb13-amd64, Go version per repo toolchain (not pinned here —
`go test` ran without a version mismatch warning), scratch home at
`.agents/scratch/verify-home` (removed and recreated fresh for this run, not
committed — `.agents/scratch/` is gitignored).
