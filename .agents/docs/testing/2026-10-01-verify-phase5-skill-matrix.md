---
owner: main
date: 2026-10-01
---
# Verification — Phase 5 skill-matrix rows on Claude Code, and the § 4 sandbox starts (2026-10-01)

Range reviewed: working tree on 3a1a861 (Phase 5 edits uncommitted at run time; the installed sandbox copy was the working tree as of the run's start — the later `restatement-grep.sh` fix, the `agents/documentation.md` spec-lane wording and the context's `project-memory` pointer were not in the copy the six rows ran against; the documentation run below was)

## Status: PASS

Procedure: `.agents/docs/testing/skill-matrix.md` § Procedure, run by `.agents/scratch/audit-fix/kit-phase5.sh <row>` (not committed, ADR-010), which sources `testing/kit-preamble.sh`, clones `~/sandbox/audit-fix/proj2` (the scaffolded demo; `proj` for `init-deploy`) into `~/sandbox/audit-fix/runs/p5-<row>` outside every checkout, runs `claude -p` from the sandbox home `~/sandbox/audit-fix/sbNew` with the working tree installed (`sh-agent install --home …`, drift none), and runs the row's check cell in the clone with `OUT=<stdout file>`. Credentials copied in for the run and deleted by name afterwards. `xargs -P 3`.

| row | pass | agent dispatched (from the transcript) | secs | note |
|---|---|---|---|---|
| `review-docs` | PASS | reviewer | 139 | report `reviews/2026-10-01-review-HEAD-1-HEAD-docs.md`, Status COMPLETE, README.md:7 vs greet.go:5 found; its `## Unverified` said the restatement grep "exited 1 with no output every time" → BUG-138, fixed in this change |
| `review-resume` | PASS | explorer | 187 | brief `plans/2026-10-01-recon-project-resume.md`, Status COMPLETE; names the missing `2026-09-20-ratelimit.md`, the empty Now, the placeholder path; no source touched |
| `distill-from-log` | PASS | dreamer | 186 | `retros/2026-10-01-distill-go-test-race.md` quotes the three seeded `gap` timestamps; nothing qualified as a skill (a correction, routed to /retro); no `feedback.jsonl` append by the distill |
| `handoff` | PASS | changelog, bugs | 221 | `handoffs/2026-10-01-0721-request-logging-middleware.md`; ran the cold-reader check itself (names the four missing files the resume order reads); the fixture's non-bare `origin` refused the push, recorded by the run as its BUG-001 — a fixture property, not a skill defect; the matrix row's prompt now says there is no remote |
| `init-deploy` | PASS | none | 207 | bare `proj` clone: `deploy/local.md` with `go run ./cmd/server`, `## Deploy` in AGENTS.md, gcp.md Not applicable, `init:` commit, no CLAUDE.md; Step 5b ran through sub-step 6 (settings file written) |
| `init-update` | PASS | none | 262 | `retros/.gitkeep` restored, `init: update` commit, no tracked file modified or deleted outside `.gitignore` |

## Agent definition change (documentation spec lane, `agents/documentation.md:26,46`)

`kit-agents.sh documentation <id>` (the Phase 4 kit, unchanged), one run per edit of the file today: p5d1 (spec lane widened to `specs/<app>/*.md`) dispatched=1 roadmap_untouched=1; p5d2 (`ls -d cmd/*/` → a Glob, no shell) dispatched=1 roadmap_untouched=1; p5d3 (Glob made file-shaped, `cmd/*/main.go`) dispatched=1 roadmap_untouched=0 later_section_added=1. The ROADMAP lane rule held in 2 of 3 runs today, 5 of 6 since the Phase 4 fix — BUG-139 (open). None of today's three edits touches that rule.

## § 4 sandbox starts

- Claude Code: the six rows above are starts from `sbNew`; every `p5-*.err` empty (no parse error, no permission refusal).
- Antigravity: `kit-agy-start.sh` from `sbAgy` (working tree installed `--vendor agy`, drift none; token copied, chmod 600, deleted by name on EXIT): the run listed the 29 harness skills by name plus Antigravity's three built-ins, ended its context answer with the canary word GANNET, `auto-denied` 0. `cli-*.log` grep for `parse|invalid|malformed|failed to load|yaml`: 6 lines, all `defaults.go:706] failed to get cs path: cs path /usr/bin/cs invalid` (Antigravity's own warning, unrelated to the harness files); 0 skill or rule parse errors.

## Guide § 1–3

`gofmt -l cmd internal` empty; `go vet ./...` clean; `go test -race ./...` all packages ok (including the new `TestRestatementGrepRunsOutsideHarness`, which failed before each of the two script fixes: exit 1 with no output, then exit 0 without the `README.md:1:` prefix); `sh-agent list` parses 29 skills, 8 agents; `render | wc -l` 86; `install --home .agents/scratch/x` then `status` drift none, 269 files.

## What was not exercised

Antigravity runs of the changed skills (Phase 9); the `review` (code lens), `preflight`, `init-memory`, `dream` and `dreamer` rows (unchanged by Phase 5 except where named above); `/review --fix`; the `--prompts` lens's skill-directory fallback on a project-level install (`.agents/skills/review-standards`); init-project 5b's `printf` fallback (no machine here lacks `python3`).
