---
owner: verifier
date: 2026-09-16
---
# Verification — preflight-v1.0.0   (2026-09-16)

Range reviewed: 4aa19be..HEAD (14 commits, 75 files)

## Status: PASS

Scope: complete — every command AGENTS.md's Setup commands name, plus the release-gate additions (govulncheck quoted, GOTOOLCHAIN=local build+test, install/status/list count check, coverage map over the whole range) all ran.

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-v1.0.0/go-test-1.txt`; counts: pass 103, fail 0, skip 1 at c7e3688 for `go test -race -v -count=1 ./...`. vs previous: `.agents/docs/testing/2026-09-16-verify-clear-the-board.md`, pass 92→103, skip 1→1 (same pre-existing `TestFindRepo` skip, unchanged reason). Rise in pass, no fall, no new skip — no FAIL condition. (Note: the dispatch cited "103" as the main agent's own measurement after three more tests landed since the 92-count log; this run reproduces 103 exactly.)
- `govulncheck ./...` → exit 0. Quoted result: "No vulnerabilities found." / "Your code is affected by 0 vulnerabilities. This scan also found 2 vulnerabilities in packages you import and 24 vulnerabilities in modules you require, but your code doesn't appear to call these vulnerabilities." Confirms ADR-028's claim (2 affecting → 0 affecting) directly — the 2 advisories are now in the "not called" bucket, not the "affecting" bucket. Full output at `.agents/scratch/verify-v1.0.0/govulncheck.txt`.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v1.0.0/home` → exit 0; "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; settings.json created with attribution off.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v1.0.0/home` → exit 0; "drift none — 259 files match the repo"; "repo has 27 skills, 7 agents, 0 commands, 0 rules"; `binary c7e3688a5ea7` (no warning) — matches expectation.
- `go run ./cmd/sh-agent list` → exit 0; `### Skills (27)`, `### Subagents (7)` — matches the counts styrman-crew's channel will smoke-check against this tag.
- `GOTOOLCHAIN=local go build ./...` → exit 0. Confirmed the box's installed toolchain is actually go1.24.4 (`go env GOVERSION` → go1.25.8 pinned via toolchain directive; plain `go version` outside GOTOOLCHAIN=local, inside the module, auto-switches to go1.25.8; `GOTOOLCHAIN=local go version` reports go1.24.4 — the real local binary). The build succeeded under the older, locally-installed toolchain, confirming ADR-028's "does not break a consumer building offline with the older Go" claim.
- `GOTOOLCHAIN=local go test -count=1 ./internal/...` → exit 0, `ok github.com/kesokaj/styrman-harness/internal/harness 0.453s`, under go1.24.4.

## Root cause (FAIL only)

N/A — Status is PASS.

## Coverage of changed code

Go changes (mapped function → covering test, `cmd/sh-agent/main_test.go` / `internal/harness/harness_test.go`):
- `cmd/sh-agent/main.go: update` — `TestUpdateNeverInstallsDirtyTree` (5 subtests)
- `cmd/sh-agent/main.go: attributionSet` — `TestAttributionSet`, `TestAttributionNote`
- `cmd/sh-agent/main.go: readSettings` — `TestReadSettingsProblems`
- `cmd/sh-agent/main.go: abandonTemp` — called directly in `TestReadSettingsProblems`'s neighborhood (line 875), asserted on stdout note
- `internal/harness/repo.go: Inspect` (untracked-file counting) — `TestInspectCountsUntrackedFiles`
- `internal/harness/repo.go: Pull` — `TestBehindAndPull`
- `internal/harness/plan.go: Plan` (new secret-shaped-file refusal) — `TestPlanRefusesSecretShapedFiles`
- `internal/harness/plan.go: secretShaped` (new function) — same test, `TestPlanRefusesSecretShapedFiles`
- `internal/harness/inventory.go: LoadSkills` (duplicate-name refusal) — `TestLoadSkillsRefusesDuplicateNames`
- `internal/harness/adopt.go: Adopt` (staging outside skills/) — `TestAdopt`, `TestAdoptRefusesDestructiveSources`, `TestAdoptFailedCopyKeepsTheExistingSkill`
- `internal/harness/apply.go: writeAtomic` — `TestWriteAtomicIsAtomicAndDurable`

No changed Go function in the named files was left without a covering test.

No added suppressions found: `git diff 4aa19be..HEAD` was checked for `//nolint`, `//lint:ignore`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, `.only(` — none added in this range.

Non-Go range (bulk of the diff — `context/AGENTS.md`, `agents/verifier.md`, and the reference/SKILL.md files under `skills/bash`, `skills/cloud-run-deploy`, `skills/distill`, `skills/docker`, `skills/frontend`, `skills/gcloud-setup`, `skills/golang`, `skills/init-project`, `skills/preflight`, `skills/retro`, `skills/review`): this project's own Testing section names no test framework for instruction text — it is verified by (a) the install run parsing all 27 skills and 7 agents without error (confirmed above, exit 0, no parse warnings) and (b) `/review --docs` over the range, which is a separate reviewer dispatch, not this verifier's job. This log does not itself run `/review --docs`; I did not find a same-range `/review --docs` report under `.agents/docs/reviews/` in this session — that is a gap this log names but does not close (see Not examined).
`.agents/docs/architecture/DECISIONS.md` carries ADR-028, the record for the toolchain-pin trajectory change — present and accepted.

## Not examined

- Did not confirm a `/review --docs` pass exists over exactly 4aa19be..HEAD for the prose/skill files listed above — this verifier only confirmed those files parse cleanly at install time, which is necessary but not sufficient for doc-quality review. If release requires that review gate closed, check `.agents/docs/reviews/` for a matching entry before shipping.
- Did not run the sandbox/real-home CLI smoke checks (`2026-09-16-verify-real-home-fcc8539.md` and related) fresh for this exact commit c7e3688; those were run at earlier commits in this session's trail and are not re-verified here. The dispatch did not ask for them and this run stayed inside the five numbered steps given.
- `govulncheck -show verbose` was not run; the summary line was quoted and taken as sufficient per the dispatch's ask to "quote the result line."

## Environment

Linux 6.12.107+deb13-amd64. Local `go` binary: go1.24.4 (`/usr/bin/go`). Module toolchain: go1.25.8 (pinned via `go.mod`, ADR-028), auto-fetched/used for the plain runs; `GOTOOLCHAIN=local` runs used the real go1.24.4 binary. `govulncheck` present at `/home/kesokaj/.go/bin/govulncheck`, run directly (not via `go run`). No Docker/Colima involved — no Dockerfile in this range.
