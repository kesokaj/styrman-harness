---
owner: changelog
date: 2026-09-13
---
# One review command on both vendors; `/audit` folded into `/review`     (2026-09-13)

## Summary
Built `skills/review/SKILL.md` as a unified review command on both Claude Code and Antigravity. Scope (working tree, range/ref, PR, paths, `--project`) is separate from lenses (`--code`, `--security`, `--simplify`, `--docs`, `--tests`, `--all`); the skill computes scope text itself and routes to agents (reviewer, security, verifier, explorer) with per-dispatch slugs so artifacts never collide. Merged record written when multiple lenses, `--project`, or `--fix` is used. `--fix` loops at most three rounds (apply findings, run checks, re-dispatch); `--project --fix` is refused (repo-wide edits need a plan). `--high` passes the strongest model to reviewer and security. Deleted `skills/audit/` and folded its behaviour into `/review --project --all`. Repointed references in context, review-standards, preflight, explorer, and skill-matrix. Skill count stays 27. Cold-read rounds 1–3 with nine reviewer artifacts; residual items from round 3 resolved by recommendation. BUG-012 (explorer artifact missing) opened during e2e; three round-1 findings applied in full.

## Type
feature

## Files changed
- `skills/review/SKILL.md` — new (≤140 lines, D1–D9 per plan, one command on both CLIs, scope-separate lenses, agents unchanged)
- `skills/audit/` — deleted; references repointed
- `context/AGENTS.md` — workflow-skills line `/audit` → `/review`
- `skills/review-standards/SKILL.md` — description and body: `/audit` → `/review`, "built-in /code-review" sentence removed
- `skills/preflight/SKILL.md` — description and Step 2 say `/review`, vendor-built-in reference removed
- `agents/explorer.md` — description names scan phase of `/review --project`
- `.agents/docs/testing/skill-matrix.md` — `audit` row → `review-project`, new `review` row with working-tree target

## Agents dispatched
reviewer ×9 (cold-read rounds 1–3), verifier ×2, reviewer (DoD), bugs ×2 (BUG-012, BUG-013), changelog (you)

## Verification
- `gofmt -l .` — no output
- `go vet ./...` — clean
- `go test -race -count=1 ./...` — all pass
- `go run ./cmd/sh-agent list` — 27 skills, `review` present, `audit` absent, zero parse errors
- `go run ./cmd/sh-agent render` — new workflow-skills line present
- `git grep -nE '(^|[ `(])/audit\b' -- . ':!.agents/docs' ':!skills/anti-ui-slop' ':!skills/gke-basics'` — no matches (vendored skills excluded)
- Verifier (1): `.agents/docs/testing/2026-09-13-verify-review-skill.md` — FAIL, first `/review --project --all` sandbox run hit 580-second wrapper timeout (exit 124); merged record was draft with security row pending and no explorer brief; main agent cited it anyway (BUG-013 recorded, headless allow-list permit rules logged as evidence for BUG-012: `bash`, `printf`, `echo`, `cd`, `set`).
- Verifier (2): re-launched detached with 1500-second budget; `claude-project2 exit 0` at 2026-09-13T21:34:06Z; clone `/home/kesokaj/sandbox/runs/claude-project2` holds all seven artifacts (`reviews/2026-09-13-review-project-{code,simplify,docs}.md`, `reviews/2026-09-13-security-project.md`, `testing/2026-09-13-verify-project.md`, `plans/2026-09-13-recon-project.md`, merged record `reviews/2026-09-13-review-project-merged.md` with full verdicts per lens and every link resolving); source untouched; explorer brief written in this run. See file `.agents/docs/testing/2026-09-13-verify-review-skill-2.md`.
- Reviewer (DoD): `.agents/docs/reviews/2026-09-13-review-review-skill.md` — REQUEST CHANGES on one medium (skill mentioned Claude's `/code-review`, removed per D8) and one low (D7: "four dispatches" corrected to four agents, six dispatches); both fixed.

## Requirements
R1 (one source, two vendors), R2 (the skill never mentions a vendor built-in), R5 (explicit scope text handed to agents unchanged)
