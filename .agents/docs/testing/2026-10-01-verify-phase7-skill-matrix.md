---
owner: main
date: 2026-10-01
---
# Verification — Phase 7 skill-matrix rows on Claude Code (2026-10-01)

Range reviewed: working tree on 677e8ee (Phase 7 edits uncommitted at run time)

## Status: PASS

Procedure: `.agents/docs/testing/skill-matrix.md` § Procedure, run by `.agents/scratch/audit-fix/kit-phase7.sh <row>` (not committed, ADR-010), which sources `testing/kit-preamble.sh`, clones `~/sandbox/audit-fix/proj2` (`proj` for `init-deploy`, its README removed in a fixture commit) into `~/sandbox/audit-fix/runs/p7-<row>` outside every checkout, runs `claude -p` from the sandbox home `~/sandbox/audit-fix/sbNew` with the working tree installed (34 skills, drift none), and runs the row's check cell in the clone with `OUT=<stdout file>`. No `gcloud`, `docker`, `terraform` or `npx` in the allowlist: a cloud command is denied, the mock. Credentials copied in for each batch and deleted by name after it. Skill loading is read from the transcript: a `Skill` tool use for a natural-language prompt, a `<command-name>/<skill>` block for a slash prompt (Claude Code expands `/name` without the Skill tool).

| row | pass | skill loaded (transcript) | cloud commands | secs | note |
|---|---|---|---|---|---|
| `testing-strategy` | PASS | `Skill: testing-strategy` | 0 | 23 | |
| `operate` | PASS | `Skill: operate-gcp` | 0 | 13 | |
| `api-auth` | PASS | `Skill: api-auth-gcp` | 0 | 21 | |
| `data` | PASS | `Skill: data-gcp` | 0 | 34 | |
| `spec` | PASS (check re-run by hand in the clone) | `<command-name>/spec` | 0 | — | `specs/server/hello-api.md` with `status: draft`, § Interfaces, § Verification; the run also wrote a handoff and a BUG-001 for the fixture's non-bare `origin` (fixture property). The kit's result line was lost (BUG-140) |
| `requirements` | PASS (check re-run by hand) | `<command-name>/requirements` | 0 | — | `.agents/scratch/requirements-draft.md` with `## Must`; `REQUIREMENTS.md` untouched; the draft presented as numbered items for approval. Result line lost (BUG-140) |
| `explain-project` | PASS (check re-run by hand) | `<command-name>/explain-project` | 0 | — | `.agents/docs/onboarding.md`, `owner:` frontmatter, ≥ 3 sections. Result line lost (BUG-140) |
| `init-deploy` (re-verified) | PASS | slash `/init-project` | 0 | — | first run FAIL on the extended check because the `proj` fixture had a README, which the skill rightly left alone (the row's base now removes it); re-run: README from the template with the Step 3 answers filled (`## Quick start` with the three commands), `testing/guide.md` copied as the template (12 placeholders left for the verifier), `deploy/local.md` filled, gcp.md Not applicable, `init:` commit |

## Agent definition changes (behavioural runs, `kit-agents.sh`)

- `agents/documentation.md`: p7b1 (lane tree gains `onboarding.md`) dispatched=1 roadmap_untouched=1; p7d2 (spec lane wording, guide wording) dispatched=1 roadmap_untouched=1; p7d3 (one `specs/<app>/` per binary) dispatched=1 roadmap_untouched=0 later_section_added=1. The ROADMAP lane rule held in 2 of 3 runs today (Phase 7), 4 of 6 today overall, 7 of 9 since the Phase 4 fix — BUG-139 stays open; none of today's edits touches that rule.
- `agents/verifier.md`: p7v1 (fills an untouched template guide) dispatched=1 step_numbered_commands=1 Status PASS — 1 of 1.

## What was not exercised

Antigravity runs (Phase 9); the GCP skills against a project (every gcloud command denied by the mock, as designed); `/spec` with `--app`; `/requirements` through to an approved write (the row forbids it by design).
