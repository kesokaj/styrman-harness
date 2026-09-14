---
owner: changelog
date: 2026-09-14
---
# Documentation steps, phase 2     (2026-09-14)

## Summary

Completed phase 2 of the documentation-steps plan: the verifier writes `testing/<app>/guide.md` on first use and maintained thereafter; review and preflight records gain an Outcome section appended by the task that acts on findings; retro proposals carry a Status line from day one and `--harness` lists still-queued ones first; handoff and retro templates reference pending REQUIREMENTS.md edits; deploy guides gain "Operator steps" sections naming the by-hand procedures on each target; the skill matrix gains a documentation row that exercises the new deploy step rule; all six queue files backfilled with per-proposal Status lines. The rendered context stays at 91 lines (under budget). Verification: lint, tests, `gofmt`/`vet`/`test` green; 27 skills, both sandbox CLIs parse and quote the new step wordings; `/retro --harness --dry-run` from this checkout against the real queue lists four queued 2026-09-04 proposals and one pending Known-limits edit; documentation matrix row ran on both CLIs with PASS; verifier created this repo's `testing/sh-agent/guide.md`. Reviewers (three rounds, cold read): backfill method for proposals, rejected status shape, documentation lane wording, guide file name in checklist, handoff template scope, verifier Records shape — all approved with minor trims. One agent-DoD review also approved.

## Type

feature (phase 2 of 2)

## Files changed

- `agents/verifier.md` — guide file creation and maintenance, Records section
- `agents/documentation.md` — guide file name, lane wording, Records section
- `agents/reviewer.md` — Records section
- `agents/security.md` — Records section
- `agents/explorer.md` — Records section
- `agents/bugs.md` — Records section
- `agents/changelog.md` — Records section, entry template with Docs line
- `context/AGENTS.md` — step 2 appended (Outcome section mechanics, Findings wording), step 3 DoD Docs clause clarified
- `skills/verifier/SKILL.md` — notes that guide is created on first use, not scaffolded
- `skills/review/SKILL.md` — step 4 and step 5 Outcome section rule
- `skills/preflight/SKILL.md` — step 4 `--docs` lens, step 6 NO-SHIP Outcome mechanics
- `skills/retro/SKILL.md` — step 1 `--harness` lists queued proposals first, step 5 Status line from day one
- `skills/handoff/SKILL.md` — "Requirements edits awaiting approval" section in template
- `.agents/docs/deploy/local.md` — Operator steps section
- `.agents/docs/deploy/remote.md` — Operator steps section
- `skills/init-project/assets/agents-tree/docs/deploy/local.md` — Operator steps placeholder
- `skills/init-project/assets/agents-tree/docs/deploy/remote.md` — Operator steps placeholder
- `.agents/docs/testing/skill-matrix.md` — documentation row (SHUTDOWN_TIMEOUT task)
- `~/.config/sh-agent/proposals/*.md` (six files) — Status backfill: all 2026-09-03 and 2026-09-04 proposals, P5 applied 67616cb, P6 applied e5833a3; all 2026-09-08 unchanged; all 2026-09-11 applied b64ffe8; 2026-09-10 blanket line replaced by per-proposal lines; 2026-09-04 and others: status determined by `git grep -F` of distinctive phrases in current harness
- `testing/sh-agent/guide.md` — new, verification procedures for this project (created by verifier under new rule)

## Agents dispatched

reviewer (3 cold-read rounds, 1 DoD), verifier, changelog

## Verification

Lint: `go run ./cmd/sh-agent list` ok, 27 skills, `render` 91 lines under budget. Tests: `go test -race ./...` green; `/tmp` prose test green; `gofmt -l . && go vet ./...` green. Sandbox: both CLIs parse context without error; `claude -p 'Based on Definition of Done, what are the new Docs step details?' --model claude-3-5-haiku` quoted step 3 ("it owns README.md, docs/specs/<app>/, docs/testing/<app>/, docs/deploy/{local,remote,gcp}.md, and the project AGENTS.md setup and deploy sections; a new command, flag, endpoint, or deploy step is not done until the file that describes it says so") and step 2 Outcome wording; both CLIs quote the `testing/<app>/guide.md` creation rule and "created by the verifier on first use, maintained by both" from documentation.md lane. Verifier: log created `testing/sh-agent/guide.md` on first run under new rule. Retro: `/retro --harness --dry-run` from this checkout lists four queued 2026-09-04 proposals (architecture revision, `/research` automation, agent description field pruning, derived-agent verification) and one pending Known-limits sentence from the newest handoff, writes nothing to queue. Documentation matrix row: `SHUTDOWN_TIMEOUT` task ran on both CLIs; Antigravity 1.2.2 agy-docs2 PASS (commit 818bd68 created testing/demo-app/guide.md, changelog with Docs section, verify log, handoff); Claude Code run written to `.agents/docs/testing/2026-09-14-verify-docs-matrix-row.md` by main agent.

## Docs

- `README.md` — no change (step 1 `/review` and `feature`/`bug` naming already done in phase 1)
- `.agents/docs/deploy/local.md` — added "Operator steps" section (standing procedures: `sh-agent update`, smoke checks, Go upgrade path)
- `.agents/docs/deploy/remote.md` — added "Operator steps" section (standing procedures: Mac update and checks, mirror refresh)
- `testing/sh-agent/guide.md` — new (verification procedure, created by verifier on first use)

## Requirements

R7 (project records include retro log with per-proposal Status from day one, harness state remains in `~/.config/sh-agent/`).
