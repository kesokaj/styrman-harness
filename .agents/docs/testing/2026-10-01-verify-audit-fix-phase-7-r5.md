---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-7-r5   (2026-10-01)

Range reviewed: working tree on 677e8ee (Phase 7, round 5: final)

## Status: PASS

Scope: complete (steps 1-6 of the dispatch ran; testing/guide.md exists).

## Commands
- `gofmt -l cmd internal` -> exit 0, empty output
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0; all 6 packages ok; counts: pass 492, fail 0, skip 3 at 677e8ee+dirty for go test -race -v -count=1 ./... (output .agents/scratch/verify-p7r5/go-test.txt); vs previous: r4 log not compared
- `go run ./cmd/sh-agent list` -> exit 0; "### Skills (34)"
- `go run ./cmd/sh-agent render | wc -c` (redirected to file) -> exit 0; 14969 bytes (< 15,000)
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-p7r5` -> exit 0 (266 written); `status --home ...` -> exit 0; "drift none — 265 files match the repo"
- Job command: block extracted to `.agents/scratch/verify-p7r5-job.sh`, run with the gcloud stub and env -> exit 0; exactly two gcloud lines: `gcloud run jobs deploy app-migrate --image img:t --region r --args migrate --service-account sa --set-cloudsql-instances p:r:app-db --set-secrets DATABASE_URL=app-db-url:1 --max-retries 0 --task-timeout 10m` and `gcloud run jobs execute app-migrate --region r --wait`; no "command not found".
- SQL on postgres:18 (script `.agents/scratch/verify-p7r5/sql.sh`; outputs sql-out.txt, sql-out-b.txt). Model (a) superuser postgres: block 1 GRANT ok, blocks 2 and 3 as app-migrate ok; CREATE TABLE t ok; has_table_privilege INSERT = t; has_schema_privilege public CREATE = f; app-run CREATE TABLE u -> "permission denied for schema public" (fails as required). Model (b) Cloud SQL-shaped (cloudsqlsuperuser, pg, db app owned by pg): identical results, block 1 run as pg. Note: the first model-(b) attempt failed on a flaw in my own script (CREATE DATABASE inside a multi-statement -c); fixed in the script and model (b) rerun alone, result above. Container vp7 removed after each model; leftovers 0.
- Greps: GRANT CREATE ... =1; ALTER SCHEMA public OWNER =0; GRANT USAGE ON SCHEMA public TO =0; no password exists anywhere =0; backslash-comment `\\\\   #` =0; Database roles (deploy/gcp.md) =1; without an end-user login =1; grants the migrate user CREATE =1. All as expected.
- `find skills -type f -perm /111` -> empty. Unexpected status lines filter -> empty (grep exit 1); 11 untracked paths besides this log are the pre-existing trail files.

## Coverage of changed code
Changed files are skills/docs markdown; no Go source changed. Covered by the Go skill/lint tests (go test ok), the installer parse (list 34, install, drift none), and the behavioural executions above (job command, SQL in two models). No suppressions added in Go/TS source.

## Not examined
Real Cloud SQL / IAM authentication (modelled only by the non-superuser container); run kit behavioural runs on Antigravity; vs-previous test counts against r4 log not computed.

## Environment
Linux 6.12 (Debian 13), Docker native, postgres:18 image; tree at 677e8ee with uncommitted changes (tree did not change in source during run; only this log and scratch files were added).
