---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-6-r3   (2026-10-01)

Range reviewed: working tree on 8be6c06 (Phase 6, round 3)

## Status: FAIL

Scope: every commanded step ran (guide: .agents/docs/testing/guide.md exists). One literal grep expectation not met (see Root cause); everything else passed.

## Commands
- `gofmt -l cmd internal` -> exit 0, empty
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0, 6 packages ok; counts: pass 492, fail 0, skip 3 at 8be6c06+dirty for go test -race -v -count=1 ./...; vs previous: r2 log (not compared to a counts line here; r2 log not re-read)
- `go run ./cmd/sh-agent list` -> exit 0; Skills (27), Subagents (8)
- `go test ./internal/harness -run TestRenderedContext -count=1` -> exit 0, ok
- `sh-agent install --home .agents/scratch/verify-p6r3` -> exit 0, 240 written; `status --home ...` -> exit 0, "drift none — 239 files match the repo"
- terraform example copy: init -backend=false 0; validate Success; fmt -check -recursive 0
- terraform resources-gcp.md assembled (11 hcl blocks, 10 kept) + variables.tf/versions.tf + db_url: init 0; validate Success; fmt -check 0
- go samples: block1.go (service_main) `go vet ./...` 0; block2.go (service_gcplog) `go vet ./...` 0
- SignupForm.tsx into fe-app: `npx tsc -b` exit 0
- greps: `golang:1.N.N` empty (ok); `"latest"` in skills/terraform empty (ok); `golang:1.24.4` (go.mod line) empty (ok); `^\*$\|^!` count in docker setup.md = 6 (ok); immutable_tags = true 1 in each file (ok); fieldViolations 2 lines (ok); `export function SignupForm` 1 line (ok); compute.googleapis.com count 1 (ok)
- `grep -n 'pgdata:/var/lib/postgresql$' skills/docker/references/setup.md` -> exit 1, EMPTY (expected 1 line) -> FAIL
- `find skills -type f -perm /111` -> empty; `git status --short skills/terraform/assets/example` -> four ` A` lines (main.tf, outputs.tf, variables.tf, versions.tf)

## Root cause
The pgdata check is a literal `$`-anchored grep. Line 51 of skills/docker/references/setup.md is `      - pgdata:/var/lib/postgresql   # postgres:18 moved PGDATA ... the old /data path would leave the named volume empty`; the trailing comment means `$` never matches. Substance is met (one pgdata mount, path /var/lib/postgresql, no `/data` mount); the check pattern or the line shape needs reconciling by the caller. Not patched.

## Coverage of changed code
Changed files are skills/docs/terraform/markdown samples, no Go source changed since round 2; samples covered by the extraction checks above. No suppressions checked in the Go diff (none in scope).

## Not examined
Real-home install; behavioural runs on either CLI; docker build; the r2 counts line comparison (not read); `git diff` of Go suppressions.

## Environment
Linux 6.12 (Debian), terraform from .agents/scratch/tools, Go per go.mod 1.24.4 line; no Docker used. Tree status at end identical to start (diff empty) apart from this log.
