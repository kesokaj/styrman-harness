---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-3   (2026-10-01)

Range reviewed: 329a99b..working tree (uncommitted)

## Status: FAIL

Scope: guide § 1, § 2, § 3 and bash -n ran; § 4 (sandbox start of each CLI) and the behavioural run are NOT in this log (main agent runs them with .agents/scratch/audit-fix/kit-size-tier.sh); not complete. guide.md exists.

## Commands
- `gofmt -l cmd internal` → exit 0, printed nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 1 (output .agents/scratch/verify-audit-fix-phase-3/go-test-1.txt); counts: pass 484, fail 1, skip 3 at 329a99b+dirty for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-10-01-verify-audit-fix-phase-1-r2.md, pass 484→484
  - FAIL root package TestEmbeddedPlanEqualsCheckout (plan_test.go:42): `~/.claude/skills/review/scripts/restatement-grep.sh mode -rw-r--r-- ... from the embedded files, -rwxr-xr-x from the tree` (same for gemini/antigravity). Other packages ok (cmd/sh-agent, dream, harness, memory, search).
- `go run ./cmd/sh-agent list` → exit 0 (29 skills, 8 agents)
- `go run ./cmd/sh-agent render | wc -c` (render to file, wc -c) → exit 0; 14766 bytes
- `go test -count=1 ./internal/harness -run 'TestRenderedContextFitsAntigravityBlock|TestRenderedContextUnderTarget' -v` → both PASS; Antigravity block 14766 bytes (context 14766, rules 0), budget 22000; under-target 14766 vs target 15000
- `go test -count=1 ./cmd/sh-agent -run TestHarnessMap -v` → PASS
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-p3` → exit 0, 270 written; `status --home ...` → exit 0, `drift none — 269 files match the repo`; installed restatement-grep.sh is -rwxr-xr-x in the scratch home (source-tree install keeps the bit); scratch home removed.
- `bash -n .agents/docs/testing/kit-preamble.sh skills/review/scripts/restatement-grep.sh` → exit 0

## Root cause  (FAIL)
Root-package TestEmbeddedPlanEqualsCheckout fails: the embedded harness (embed.FS) reports restatement-grep.sh as 0644 while the checkout (git mode 100755) yields 0755; plan.go:137-138 derives mode from the file info (`Perm()&0o755|0o600`), and embed.FS carries no exec bit. So a release binary would install the new executable review script without its exec bit (the source-tree install above is fine). Not pre-existing in a passing sense: the file entered in 329a99b (phase 2) and is untouched by this range, but the earlier log's 484 pass did not include this fail count either; the previous log's fail count was not compared. Needs fix (embedded plan must set exec for it, or the script must be invoked as `bash <path>`), not patched here.

## Coverage of changed code
- internal/harness/size_test.go — the only Go change (+27 lines): TestRenderedContextUnderTarget, PASS at 14766. From reading: it renders the checkout context via RenderContext and t.Fatalf when len(ctx) > 15000, so it FAILS above 15,000 and passes at/below; I did not mutate-run it. A must-fail run is not recorded.
- Non-Go changes (context/, skills/, agents/bugs.md, docs) are covered by § 2/§ 3 checks above only.
- Suppressions added: none checked beyond the Go diff (no added skips/nolint in size_test.go).

## Not examined
- § 4 sandbox start of each CLI; behavioural run of agents/bugs.md and context changes; Antigravity runs; /review --docs; govulncheck; .agents/docs/testing/skill-matrix.md.
- Must-fail mutation of TestRenderedContextUnderTarget.
- Tree at end vs start: only new untracked artifacts appeared (reviews/2026-10-01-review-329a99b-worktree-docs.md from another agent, and this log); no source changed.

## Environment
Linux 6.12 (Debian 13), go from PATH (version not recorded), no Docker used.
