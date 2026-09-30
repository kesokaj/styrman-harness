---
owner: verifier
date: 2026-09-30
---
# Verification — retro2-apply   (2026-09-30)

Range reviewed: working tree — changed files: .agents/docs/deploy/local.md TODO.md skills/handoff/SKILL.md skills/init-project/assets/TODO.md skills/retro/SKILL.md; untracked, not ignored: none outside .agents/docs; read the diff with git diff HEAD -- <files>.

## Status: PASS

Scope: complete for the tiers this dispatch named (Go checks, harness self-check, scratch-home install/drift/copy-match, list sanity). Sandbox, real-home and Antigravity tiers intentionally not run — the dispatch reserves those for the main agent at the release gate. `testing/guide.md` already exists (dated 2026-09-17); no update needed since no steps changed.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race -count=1 ./...` → exit 0; 4 packages all `ok` (root, cmd/sh-agent, internal/harness, internal/search); output at `.agents/scratch/verify-retro2-apply/go-test-race.txt`
- `go test -race -v -count=1 ./...` (rerun for per-test counts per the Test counts rule) → exit 0; output at `.agents/scratch/verify-retro2-apply/go-test-race-v.txt`; counts: pass 374, fail 0, skip 2 at ae71bb6+dirty for `go test -race -v -count=1 ./...`; vs previous: `.agents/docs/testing/2026-09-30-verify-v1.15.0-gate.md` (counts at e69c56e), pass 374→374, skip 2→2 — no regression
- `go test -count=1 -run 'TestRenderedContextFitsAntigravityBlock|TestLint' ./internal/harness -v` → exit 0; both `--- PASS`; `TestRenderedContextFitsAntigravityBlock` reports 21996 bytes against the 22000-byte budget (unchanged by this prose-only diff, as expected — none of the changed files are `context/*.md` or `rules/*.md`)
- `go run ./cmd/sh-agent install --home .agents/scratch/home` → exit 0; "266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed"; built from `ae71bb6+dirty`
- `go run ./cmd/sh-agent status --home .agents/scratch/home` → exit 0; `drift      none — 265 files match the repo`
- `go run ./cmd/sh-agent list` → exit 0; prints 28 skills including `init-project`, 8 agents; scaffold parses fine
- `diff skills/handoff/SKILL.md .agents/scratch/home/.claude/skills/handoff/SKILL.md` → exit 0, no output (installed copy matches working tree)
- `diff skills/retro/SKILL.md .agents/scratch/home/.claude/skills/retro/SKILL.md` → exit 0, no output
- `diff skills/init-project/assets/TODO.md .agents/scratch/home/.claude/skills/init-project/assets/TODO.md` → exit 0, no output
- Final `git status --short --untracked-files=all` → only the five named files (M) plus new untracked files under `.agents/docs/` (my own log, plus `.agents/docs/changes/2026-09-30-1311-retro2-apply-and-docs.md` and `.agents/docs/reviews/2026-09-30-review-worktree-retro2-docs.md` — another agent's trail artifacts written while this run was in progress, treated as seen per the "tree under test" rule, not as the tree under test moving). No source, test or config path changed during the run.
- No background jobs were started during this run; none to end.

## Root cause
N/A — no failures.

## Coverage of changed code
Diff is prose-only (confirmed via `git diff HEAD -- <files>`, read in full before running commands):
- `.agents/docs/deploy/local.md` — a doc paragraph rewording the redaction-rule prose (footer-closing sentence reordered, enumeration paragraph restructured); no code path, no test needed.
- `TODO.md`, `skills/init-project/assets/TODO.md` — identical one-clause addition to the header sentence ("(other than a standing procedure the deploy guide records)"); template/scaffold text, no code.
- `skills/handoff/SKILL.md` — Step 3 clarification plus a new paragraph in the handoff template about not repeating deploy-guide standing procedures and not re-listing carried-over Next Steps items; prose only.
- `skills/retro/SKILL.md` — Harness-scope bullet reworded so an approved harness-scope proposal applies immediately when run from the harness checkout itself, rather than defaulting to `queued`; prose only.
No `.go` file in the diff, so gofmt/vet/race suite and the harness self-check are the applicable Coverage checks per the guide's own "skip only when no `.go` file is in the diff — say so" rule; nothing to skip here since they were run anyway and all green. No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) were added by this diff — confirmed by reading the diff above; the two pre-existing skips in the -v run are unrelated to these files (carried from the previous baseline at the same count).

## Not examined
- Sandbox CLI-start tier (`.agents/docs/testing/2026-09-03-verify-sandbox.md`), real-home smoke tier (`2026-09-03-verify-real-home.md`), and any Antigravity behavioral run — explicitly excluded by this dispatch; the main agent runs these at the release gate.
- `go run ./cmd/sh-agent install --home .agents/scratch/home` was NOT rerun against the operator's real `$HOME` (correctly, per the standing rule); only the scratch-home path was exercised.
- `.agents/docs/changes/2026-09-30-1311-retro2-apply-and-docs.md` and `.agents/docs/reviews/2026-09-30-review-worktree-retro2-docs.md` were not read/reviewed by this run — they are another agent's artifacts, out of scope for a verifier dispatch, noted only for the git-status accounting above.

## Environment
Linux 6.12.107+deb13-amd64, go1.27.1 linux/amd64, repo at ae71bb6 (+dirty, the five named files), styrman-harness project (no Docker/Colima involved — no Dockerfile changed).
