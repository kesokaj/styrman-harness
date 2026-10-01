---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-9   (2026-10-01)

Range reviewed: working tree on 2ba4a50 (Phase 9 fix round)

## Status: FAIL

Scope: every command of the dispatch (steps 1-9) ran; guide `.agents/docs/testing/guide.md` exists. One check fails (agents/changelog.md line 3). Not complete-PASS.

## Commands
- 1. `gofmt -l cmd internal` → exit 0, empty
- 1. `go vet ./...` → exit 0
- 1. `go test -race -v -count=1 ./...` (output .agents/scratch/verify-p9-gotest.txt) → exit 0; all 6 packages ok; counts: pass 518, fail 0, skip 3 at 2ba4a50+dirty for go test -race -v -count=1 ./...; vs previous: commit 2ba4a50 message says 517 tests, pass 517→518, skip not recorded there; skips are the opt-in TestAgyCorpusLeaves, TestRecall, TestHelperIndex (env-gated)
- 1. `go test -cover -count=1 ./cmd/... ./internal/...` → exit 0; cmd/sh-agent 86.1%, dream 89.7%, harness 89.0%, memory 89.1%, search 84.2%
- 2. `go test ./internal/harness -run TestRenderedContext -count=1 -v | grep` → "Antigravity rules block is 14992 bytes ... budget 22000"; "hand-written context is 14107 bytes (render 14992 with the inventory), target 15000". `go run ./cmd/sh-agent render | wc -c` → 14992
- 3. `go run ./cmd/sh-agent install --home .agents/scratch/verify-p9` → exit 0, 266 written; `status --home ...` → exit 0, `binary 2ba4a50351ab+dirty` (no warning), `drift none — 265 files match the repo`
- 4. built `.agents/scratch/verify-p9-bin`; `memory --home .agents/scratch/verify-p9 forget` → exit 2, prints "usage: sh-agent memory add --type T ... | sh-agent memory show"; `dream --help` prints "Usage of dream:" with flags (141 seen was SIGPIPE from `head`, not the binary)
- 5. terraform example copy: init -backend=false 0; validate "Success! The configuration is valid."; `fmt -check -recursive` 0; `grep -c ignore_changes main.tf` = 1
- 5. terraform reference assembly (11 hcl blocks of resources-gcp.md minus the backend block, example variables.tf/versions.tf, db_url variable): init 0, validate Success, `fmt -check` 0; `grep -c google_sql_user resources-gcp.md` = 3
- 6. go sample p9-migrate (block importing migrate/v4/database/pgx/v5): `go mod tidy` 0, `go vet ./...` 0, `go list -deps | grep -c lib/pq` = 0
- 7. `SH_AGENT_RECALL=1 go test ./internal/search -run 'TestRecall$' -v -count=1` → index mean open 1.000 over 2; index mean all 0.633 over 41; `grep -c '0.633' search.md` = 1
- 8. `grep -c 'applied in Phase 3 with ADR-048' PRIORITY_TODO.md` = 0; `requirements-edit-proposal` count PRIORITY_TODO.md 1, memory-loop.md 1; plan file exists; `write the artifact early` counts: explain-project 1, handoff 1, bug 2, feature 3, distill 1, retro 1, preflight 3
- 8. FAIL: `grep -n ROADMAP agents/changelog.md` → line 3 (description) ends "never writes specs, README, or ROADMAP; those belong to documentation." That says ROADMAP belongs to documentation; it contradicts line 34 ("you do NOT edit ROADMAP.md, TODO.md ... the main agent keeps those itself") and line 84.
- 9. `find skills -type f -perm /111` → empty; non-M/A/D/?? status lines → none

## Root cause
agents/changelog.md frontmatter `description:` still carries the older sentence "Never ... writes specs, README, or ROADMAP; those belong to documentation", left over beside the corrected body text. A restatement of a reversed claim was missed in the description field (one grep over key words would find it). Fix not applied (verifier does not edit).

## Coverage of changed code
- internal/harness/size_test.go, internal/dream/run_test.go, cmd/sh-agent/memory_test.go, internal/search testdata — are tests themselves; internal/dream/run.go covered by run_test.go; cmd/sh-agent/main.go covered by memory_test.go and the package tests. Whether the new `last_run.error` and exit-2 paths have named assertions was not read line by line.
- Skills/agents/docs: no Go test coverage beyond TestSkillLinksResolve/lint tests in internal/harness (passed); samples verified by steps 5-6.
- Suppressions: diff not scanned for added nolint/skip/.only (not examined).

## Tree check (step 8)
Status now vs step 1: only new paths are other agents' artifacts: `.agents/docs/reviews/2026-10-01-review-worktree-p9-docs.md` (new) and this log. No source, test or config path changed.

## Not examined
- Added-suppression scan of the diff (`//nolint`, `t.Skip`, `.only(`) was not run.
- Per-file coverage mapping beyond package coverage; behavioural runs on either CLI; govulncheck.
- Full read of testing/guide.md beyond existence.

## Environment
Linux 6.12.107 (Debian 13), Go from PATH, Terraform from .agents/scratch/tools, no Docker used.
