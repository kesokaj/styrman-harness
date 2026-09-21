---
owner: verifier
date: 2026-09-21
---
# Verification — preflight v1.5.0   (2026-09-21)

Range reviewed: v1.4.0..c7e178e

## Status: PASS

Scope: guide § 1, § 2 (list), § 3 (scratch-home install/status/canary/cmp), govulncheck ran; § 3a, § 3b, § 3c not run; § 4/5 cited from today's existing logs, not rerun (per dispatch)

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0
- `go test -race -count=1 ./...` → exit 0; all packages `ok` (root, cmd/sh-agent, internal/harness); run without -v per guide § 1 as written, so no per-test PASS/FAIL/SKIP lines; counts: not printed; output saved to .agents/scratch/verify-v150-go-test.txt
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1 -v` → exit 0; `--- PASS: TestRenderedContextFitsAntigravityBlock`; "Antigravity rules block is 20205 bytes (context 20205, rules 0), budget 22000" — matches expected 20205 of 22000 exactly; output saved to .agents/scratch/verify-v150-render-test.txt
- `govulncheck ./...` → exit 0; "No vulnerabilities found." (go1.27.1, so no known stdlib advisories expected); output saved to .agents/scratch/verify-v150-govulncheck.txt
- `go run ./cmd/sh-agent list` → exit 0; header confirms "### Skills (27)" and "### Subagents (7)" — matches expected 27 skills / 7 agents; output saved to .agents/scratch/verify-v150-list.txt
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v150` → exit 0; "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; settings.json created with attribution off; output saved to .agents/scratch/verify-v150-install.txt
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v150` → exit 0; `drift      none — 259 files match the repo` (259 vs 260 written is expected: settings.json is created-once, not drift-tracked); output saved to .agents/scratch/verify-v150-status.txt
- `tail -1 .agents/scratch/verify-v150/.claude/AGENTS.md` → "Canary: this context ends with the word GANNET." — present, matches BUG-085 canary check
- `cmp .agents/scratch/verify-v150/.claude/AGENTS.md .agents/scratch/verify-v150/.gemini/config/rules/AGENTS.md` → exit 0 (byte-identical)

## Root cause  (FAIL only)
None — all commands passed.

## Coverage of changed code
Range v1.4.0..HEAD touches `internal/harness/render.go` (context rendering/size budget) and `internal/harness/size_test.go` (the size test itself). Covering tests:
- `internal/harness/render_test.go` → `TestInventoryAndRender` covers `render.go`'s rendering logic (existing test, not rerun individually beyond the full-suite `go test ./...` pass above, which includes it).
- `internal/harness/size_test.go` → self-covering (`TestRenderedContextFitsAntigravityBlock`), run explicitly above.

Not covered by any Go test (prose-level, by design — the guide and REQUIREMENTS note test coverage stops at structure/byte-budget, not content):
- The prose content of `context/AGENTS.md` and the rules files beyond byte count (wording, correctness of instructions) — verified instead by the sandbox/real-home behavioural runs cited below, not by `go test`.
- The skill texts under `skills/*/SKILL.md` and `skills/*/references/*.md` — parsed/loaded correctness verified only by `sh-agent list` (structural: frontmatter present) and by the cited sandbox/skill-matrix runs, not by unit tests.

Per this dispatch's scope, the CLI-start and real-home behavioural tiers (guide § 4/5) are cited, not rerun:
- `.agents/docs/testing/2026-09-21-verify-sandbox-pruned-context.md`
- `.agents/docs/testing/2026-09-21-verify-real-home-smoke-5.md`
- `.agents/docs/testing/2026-09-21-context-consolidation-phase-2-behavioural.md`

## Not examined
- Guide § 3a (plain-copy install) and § 3b (`report --dry-run`) — not run; this dispatch's instructions scoped to § 1, § 2 (list), § 3 (scratch-home install/status/canary/cmp only), govulncheck, and citing § 4/5. No change in v1.4.0..HEAD to `IsCheckout`/rebuild logic or to `report.go` was checked for, so it is unconfirmed whether § 3a/3b apply to this range.
- Guide § 3c (verb × place-kind matrix) — not run.
- CLI starts (Claude Code, Antigravity) and real-home smoke checks — not run in this session; relied on citation of today's earlier logs per dispatch instructions ("never run agy or claude").
- Full skill-matrix (`.agents/docs/testing/skill-matrix.md`) — not run; dispatch did not request it and no "complete/all skills" claim is made here.
- Did not independently verify the content of the three cited logs beyond their filenames existing in the dispatch text; did not re-open them to confirm PASS status.

## Environment
Linux, go1.27.1 (user-level SDK per AGENTS.md § Setup commands), govulncheck at ~/.go/bin/govulncheck. Working tree had pre-existing untracked files at dispatch start (`.agents/docs/reviews/2026-09-21-review-v1.4.0-HEAD-code.md`, `.agents/docs/reviews/2026-09-21-security-v1.4.0-HEAD.md`, from earlier preflight steps in this session) plus this log itself — `sh-agent status`'s "uncommitted changes" / "+dirty" marker reflects those untracked doc files, not any change to tracked source. No file outside `.agents/scratch/verify-v150` or this log was created or modified; nothing deleted; real `$HOME` untouched.

