---
owner: verifier
date: 2026-09-29
---
# Verification — post-release-v1.12.0   (2026-09-29)

Range reviewed: 0beb429..working tree

## Status: PASS

Scope: complete — every item the dispatch named ran, and `.agents/docs/testing/guide.md` exists (present, unchanged, not rewritten).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race -v -count=1 ./...` → exit 0; output at `.agents/scratch/verify-post-release-v1.12.0/go-test-race.txt`; counts: pass 171, fail 0, skip 0 at 0beb429+dirty for `go test -race -v -count=1 ./...`. vs previous: no earlier log ran this exact command with `-v` (prior gate logs, e.g. `2026-09-29-verify-v1.12.0-gate-r4.md`, ran `go test -race -count=1 ./...` without `-v`, reporting "counts: not printed" / package-level `ok` only) — vs previous: none (no comparable prior count).
- `go test -count=1 -run TestReinstallBuildsAndInstalls ./internal/harness` → exit 0; `--- PASS: TestReinstallBuildsAndInstalls (0.21s)`. This is the BUG-104 fix under test: the test now stats its own pid-named build directory (`sh-agent-build-<pid>`) instead of globbing `sh-agent-build-*`, so a concurrent `go test` process's in-flight build no longer flakes it. Confirmed no stray build dirs after the full suite: `ls .agents/scratch/sh-agent-build-* 2>/dev/null | wc -l` → 0.
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1 -v` → exit 0; `size_test.go:48: Antigravity rules block is 21952 bytes (context 21952, rules 0), budget 22000`; `--- PASS`. Byte count matches the prior gate log (`2026-09-29-verify-v1.12.0-gate-r4.md` line 27: "21952 bytes, under the 22,000-byte budget") — unchanged, as expected since the change set does not touch `context/`.
- `go run ./cmd/sh-agent list` → prints 28 skills, 8 agents — matches dispatch's expected counts exactly.
- Install tier into `.agents/scratch/verify-post-home` (removed after use):
  - `go run ./cmd/sh-agent install --home .agents/scratch/verify-post-home` → exit 0; `266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed`.
  - `go run ./cmd/sh-agent status --home .agents/scratch/verify-post-home` → exit 0; `repo has 28 skills, 8 agents, 0 commands, 0 rules`; `drift none — 265 files match the repo`.
  - `grep -c "not symlinks" .agents/scratch/verify-post-home/.claude/skills/init-project/SKILL.md` → 1
  - `grep -c "not symlinks" .agents/scratch/verify-post-home/.gemini/config/skills/init-project/SKILL.md` → 1
  - `tail -c 120 .agents/scratch/verify-post-home/.claude/AGENTS.md` → ends `Canary: this context ends with the word GANNET.`
  - Scratch home `.agents/scratch/verify-post-home` removed after checks (`rm -rf`).
- Real-home smoke: cited, not rerun per dispatch instruction. `.agents/docs/testing/2026-09-29-verify-real-home-v1.12.0.md` line 9: `## Status: PASS on both CLIs — the three smoke checks, plus the first runs of /dream, the resume clause, the dreamer and the linked-store row`.

## Root cause  (FAIL only)
N/A — all commands passed.

## Coverage of changed code
- `internal/harness/harness_test.go` (BUG-104 fix: pid-scoped stat instead of glob) — the change is itself a test; verified directly by running `TestReinstallBuildsAndInstalls` (PASS) and by confirming no stray `sh-agent-build-*` directories remain after the full suite. No production code changed.
- `skills/init-project/SKILL.md` (one added sentence, Step 1: every write refuses a symlinked target) — prose/skill instruction, no Go test framework applies; behavioural coverage is the `init-memory-linked-store` matrix row, already run and recorded PASS in `2026-09-29-verify-real-home-v1.12.0.md` (per the skill-matrix.md row's own citation).
- `AGENTS.md` (one added sentence, § Release: gate NO-SHIP rule for untested Antigravity runs) — project policy prose, not code; not test-coverable, no gap.
- `.agents/docs/testing/skill-matrix.md` (two rows corrected: `init-memory-linked-store` marked run+PASS with citation, plus dropped a stale `.gitignore` clause from its command) — a documentation correction reflecting an already-recorded run; no new behavior to cover.
- `.agents/docs/bugs/INDEX.md` (new BUG-104 row) — record only.
- New files `.agents/docs/bugs/BUG-104-reinstall-build-glob-flake.md`, `.agents/docs/changes/2026-09-29-1609-post-release-checks.md`, `.agents/docs/testing/2026-09-29-verify-real-home-v1.12.0.md` — trail records, not code; not test-coverable.

No suppressions, skips, or ignore directives (`//nolint`, `t.Skip`, `.skip(`, `.only(`, `@ts-ignore`, `eslint-disable`) were added anywhere in this change set (checked via `git diff 0beb429` on every changed file).

## Not examined
- Nothing in the dispatched scope was skipped. Out of dispatch scope and not examined here: any Antigravity-side behavioural run beyond what `2026-09-29-verify-real-home-v1.12.0.md` already recorded (cited, not rerun, per instruction); `govulncheck ./...` (not in this dispatch's command list); the full `.agents/docs/testing/skill-matrix.md` matrix beyond the one corrected row.

## Environment
Linux 6.12.107+deb13-amd64. go version go1.27.1 linux/amd64. HEAD at dispatch time: 0beb429 with the uncommitted working-tree diff listed in Range reviewed (working tree dirty throughout, as instructed — the tree was not moved).
