---
owner: verifier
date: 2026-09-21
---
# Verification — preflight v1.6.0, second run (gated head moved by 2facb44)   (2026-09-21)

Range reviewed: v1.5.0..2facb44

## Status: PASS

Scope: re-verify after gated head advanced by one commit (2facb44: one sentence added to `skills/init-project/SKILL.md`, ROADMAP Now cleared). No Go source changed in this range — `git diff --stat 2facb44^..2facb44` shows only `ROADMAP.md` and `skills/init-project/SKILL.md` touched (1 insertion, 3 deletions total). testing/guide.md exists and was followed; earlier gate logs (`testing/2026-09-21-verify-preflight-v1.6.0.md`, plan `2026-09-21-project-agents-md-native.md`, and the prior `-verify-agents-md-native.md`) are cited, not rerun. Scope: complete.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output (clean)
- `go test -race -count=1 ./...` → exit 0; 3 packages, all `ok`, 0 FAIL. counts: pass N (subtests not itemized, no `-v`), fail 0, skip 0 at 2facb44 for `go test -race -count=1 ./...`; vs previous: testing/2026-09-21-verify-preflight-v1.6.0.md (prior run also all-ok, same command) — no regression.
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1 -v` → exit 0; `--- PASS: TestRenderedContextFitsAntigravityBlock`; log line "Antigravity rules block is 20329 bytes (context 20329, rules 0), budget 22000" — matches expected 20329 of 22000 exactly.
- `go run ./cmd/sh-agent list` → exit 0; 27 skills listed (table header "Skills (27)"), `init-project` present; 7 subagents listed (table header "Subagents (7)") — matches expected 27/7.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v160h` → exit 0; 258 written, 0 unchanged, 0 removed, 0 kept, 258 files managed.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v160h` → exit 0; `drift      none — 257 files match the repo` (no drift; repo has 27 skills, 7 agents, 0 commands, 0 rules).
- `grep -c "ancestor and nested" .agents/scratch/verify-v160h/.claude/skills/init-project/SKILL.md` → exit 1, count 0 — the literal phrase "ancestor and nested" does not appear in the deployed file; that exact wording is only in the commit message (2facb44), not the file text. The commit's actual added sentence reads "Native reading is wider than the stub was: Claude Code also loads an `AGENTS.md` in any directory above the project and, when it reads a file there, in any subdirectory — ...". Verified instead by `diff skills/init-project/SKILL.md .agents/scratch/verify-v160h/.claude/skills/init-project/SKILL.md` → IDENTICAL, confirming the new sentence deployed byte-for-byte from source. Treating this as the expected grep pattern being a paraphrase mismatch in the dispatch, not a deployment defect — the `status` drift-none result above already confirms the deployed tree matches the repo exactly.
- `test ! -e CLAUDE.md` (repo root) → exit 0 — confirmed absent, consistent with ADR-036 / REQUIREMENTS R2 (no CLAUDE.md stub in this project).
- `test -f .agents/scratch/verify-v160h/.claude/CLAUDE.md` → exit 0 — global stub present at the installed home, content is the generated header comment + `@AGENTS.md` (this is `~/.claude/CLAUDE.md`, the vendor-required stub outside any project, unaffected by R2/ADR-036 which concerns project-level CLAUDE.md files).

## Coverage of changed code
No Go source changed in v1.5.0..2facb44 (nor in the prior v1.5.0..HEAD range covered by the earlier v1.6.0 gate log). The only content changes are documentation/skill markdown (`skills/init-project/SKILL.md`) and `ROADMAP.md`. No suppressions (`//nolint`, `t.Skip`, `@ts-ignore`, `.skip(`, `.only(`) added in range.

## Not examined
- Did not rerun `govulncheck ./...` (no Go source changed since the prior gate run that covered it).
- Did not start Claude Code or Antigravity CLI (per dispatch instruction: never run agy or claude).
- Did not touch the real home (`~/.claude`, `~/.gemini`) — install/status ran only against `.agents/scratch/verify-v160h`.
- Did not re-verify the tier-3/tier-4 behavioural runs (`claude -p --agents ...`) for `init-project`'s changed sentence — out of scope for this dispatch and already covered, if at all, by the prior gate logs cited above.

## Environment
Go: user-level SDK ahead of `/usr/bin/go` per project AGENTS.md (Debian trixie has neither go1.24.6 nor go1.25.8; govulncheck not rerun this pass). OS: Linux 6.12.107+deb13-amd64. Scratch home: `.agents/scratch/verify-v160h` (gitignored, untouched real home).
