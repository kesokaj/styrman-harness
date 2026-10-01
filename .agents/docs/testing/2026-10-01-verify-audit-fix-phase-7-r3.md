---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-7-r3   (2026-10-01)

Range reviewed: working tree on 677e8ee (Phase 7, round 3)

## Status: FAIL

Scope: steps 1-8 all ran; one expected-output mismatch (step 7 issuer grep). Guide: .agents/docs/testing/guide.md (exists).

## Commands
- `gofmt -l cmd internal` → exit 0, empty output
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0; counts: pass 492, fail 0, skip 3 at 677e8ee+dirty for go test -race -v -count=1 ./...; all 6 packages ok; vs previous: .agents/docs/testing/2026-10-01-verify-audit-fix-phase-7-r2.md, pass 492→492, skip 3→3 (r2 counts line as recorded there)
- `go run ./cmd/sh-agent list` → exit 0; status shows "34 skills, 8 agents"
- `go run ./cmd/sh-agent render | wc -c` (run to file, bytes counted) → exit 0; 14969 bytes (<15,000, margin 31)
- `go test ./internal/harness -run TestRenderedContext -count=1` → exit 0 ok
- `install --home .agents/scratch/verify-p7r3` → exit 0, 266 written; `status --home ...` → exit 0, "drift none — 265 files match the repo"
- Go sample: extracted the go block with schedulerSA and func New( (1 of 4 go blocks) to .agents/scratch/go-samples/p7-api-auth-gcp-2/block.go; `gofmt -l .` → empty; `go vet ./...` → exit 0; `grep -c 'payload.Issuer\|errors.New'` → 2
- Terraform (10 of 11 hcl blocks, backend block excluded) in .agents/scratch/verify-p7r3-tfref: init -backend=false → 0; validate → Success; `fmt -check` first run exit 3 ONLY on my appended `variable "db_url"` block (my assembly wrote `type = string` unaligned; diff showed nothing else); after aligning my own appended lines, fmt -check → exit 0, validate Success again. The reference blocks themselves were fmt-clean.
- gcloud help greps (gcloud in .agents/scratch/tools): `secrets create --replication-policy` → 7; `secrets add-iam-policy-binding --member` → 3; `sql users set-password --prompt-for-password` → 2 (all ≥1)
- `docker run --rm postgres:18 postgres --version` → PostgreSQL 18.6
- SQL roles check on throwaway postgres:18 (p7r3pg): 3 blocks extracted from data-gcp § Roles (<app>→app, <project>→proj; files .agents/scratch/verify-p7r3/role{1,2,3}.sql). Block 1 as postgres: ALTER SCHEMA, exit 0. Blocks 2 and 3 as app-migrate@proj.iam: exit 0 each. CREATE TABLE t as migrate role OK. has_table_privilege(run,'t','INSERT') = t; has_schema_privilege(run,'public','CREATE') = f. Cleanup: `docker rm -f p7r3pg` ran; 0 containers named p7r3pg remain.
- Greps: REVOKE cloudsqlsuperuser FROM (data-gcp) 0 ok; 'all three checks' (api-auth-gcp) 0 ok; **'issuer' in terraform/references/resources-gcp.md = 0, expected ≥2 → MISMATCH**; 'func New(' go-samples 4 ok; '2>/dev/null ||' data-gcp 0 ok; 'read -rs' data-gcp 1 ok; 'Data and access' requirements 1 ok; 'Secret Manager: <name>' init-project README 1 ok; 'app_name' research 0 ok; 'older per-app guide' init-project 1 ok
- `find skills -type f -perm /111` → empty; status filter for non-M/A/D/?? entries → empty

## Root cause (FAIL)
The only miss is the step-7 grep `issuer` in skills/terraform/references/resources-gcp.md: 0 matches vs expected ≥2. Lines 229-230 (the Public IAM-form comment) say the route verifies "signature, audience, AND that `email` is the scheduler SA and `email_verified` is true"; the word "issuer" appears nowhere in the file (case-insensitively either). Either the round-3 edit to that file did not add the issuer wording the dispatch expected, or the expectation is stale. Not patched. Everything else in the plan passed.

## Coverage of changed code
Changed files this round are skills/markdown only (no Go changes); Go covered by existing suite (492 pass). Skill samples were compiled/validated (Go block, terraform, SQL) as above. No suppressions added in Go/TS source examined (markdown only).

## Not examined
Real-home install/Antigravity runs (guide § 5), behavioural runs per changed skill, govulncheck. The second go sample blocks (3 of 4 go blocks) were not compiled this round, only the schedulerSA/func New block. Terraform fmt check on my own appended variable block needed an alignment fix of my assembly (not the reference).

## Environment
Linux 6.12 (Debian 13), Docker available (postgres:18.6 pulled), terraform/gcloud from .agents/scratch/tools. Tree during run: only new path was this log (start vs end git status diff), no source change.
