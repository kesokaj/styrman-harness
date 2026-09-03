---
owner: main
date: 2026-09-03
---
# Handoff — 2026-09-03 — harness verified on both CLIs, installed in the real home; the Mac remains

Session ID: 77314da2
Project: /home/kesokaj/workspace/styrman-harness
Branch: main (origin up to date)

## Context

The rebuilt harness (context, 28 skills, 7 subagents, the `sh-agent` installer) had to be proven live before real use. The operator wanted every skill and agent exercised on both CLIs, deploy guides added to the convention, everything documented, and the harness installed for real.

## Completed

- Sandbox home logged in to both CLIs; both driven headless against a demo Go service. Three harness bugs found and fixed: BUG-003 (duplicate scratch rule), BUG-004 (Antigravity's YAML dropped `bug`/`feature` — installer now lints frontmatter), BUG-005 (`verifier` `isolation: worktree` hid uncommitted work). `/home/kesokaj/workspace/styrman-harness/.agents/docs/testing/2026-09-03-verify-sandbox.md`
- Full matrix: 31 mock scenarios × 2 CLIs, 62/62 pass; spec `/home/kesokaj/workspace/styrman-harness/.agents/docs/testing/skill-matrix.md`, log `2026-09-03-verify-sandbox-skills.md` beside it.
- Deploy guides: `docs/deploy/{local,remote,gcp}.md` per project — templates in `init-project`, named in the context, `/preflight` and `cloud-run-deploy` use them, `documentation` owns them, this repo has its own three under `.agents/docs/deploy/`.
- `/preflight` base-ref rule tightened; `verifier` orients first, reports BLOCKED, logs carry a `Scope:` line; `grill-me` and `skill-repair` pruned (28 skills); ADR-010 (no vendor-settings writes; matrix is documentation) and ADR-011 (GCP defaults recorded in `deploy/gcp.md`).
- Real home on the Linux box: installed at `ff58afd`, 256 files, no drift; both CLIs pass context, skill-load, and agent-dispatch checks from `$HOME` (`testing/2026-09-03-verify-real-home.md`).
- Documentation audited and brought current (changelog `changes/2026-09-03-doc-audit.md`); one retro proposal set written and applied (`~/.config/sh-agent/proposals/2026-09-03-retro.md`).

## Open / In Progress

- The Mac: `sh-agent update`, then one interactive start of each CLI (`/skills`, `/agents`, a context-only question). Plan section 12 and `TODO.md` high priority. Nothing else is open.
- No tuning before two weeks of use (ADR-009); then `/retro`.

## Key Decisions

- Symlinked `.claude/*` → `.agents/*` kept; copies rejected after a discovery test showed symlinks work and the only blocker was Claude Code's write guard on `.claude/`, which copies would hit too.
- Installer never edits vendor settings; the per-session approval for `~/.config/sh-agent` is accepted and documented (ADR-010). Headless `agy -p` therefore needs a temporary `permissions.allow` set in `~/.gemini/antigravity-cli/settings.json` for tests — add with a backup, restore after.
- Skill matrix stays scratch scripts plus a committed spec; a `sh-agent sandbox` subcommand is possible later.

## Corrections

- "Have you run the full tests in the sandboxes. Test all skills" — feedback log 2026-09-03T21:01:31+02:00. The sandbox had been reported complete at 3 of 30 skills; rule applied: verification logs carry a Scope line and "complete" means the full matrix.

## Agent notes (unverified)

- `/preflight` verdicts differ by reviewer model strength (Opus found blocking findings Flash did not); not a harness defect, but expect NO-SHIP more often on Claude Code.
- Antigravity within a trusted workspace reads sibling checkouts; keep one project per workspace when testing.

## Environment State

- Branch: main, clean after this handoff commit; real home installed at `ff58afd`.
- Running processes: none. Sandbox: `/home/kesokaj/sandbox/home` (both CLIs logged in), `demo-app`, `demo-app-agy`; test clones removed.
- Pending: nothing uncommitted.

## Artifacts to Read

- `/home/kesokaj/workspace/styrman-harness/.agents/docs/plans/2026-09-03-rebuild.md` — section 12 is the current checklist
- `/home/kesokaj/workspace/styrman-harness/.agents/docs/bugs/INDEX.md` — BUG-001..005, all fixed, prevention rules
- `/home/kesokaj/workspace/styrman-harness/.agents/docs/architecture/DECISIONS.md` — ADR-001..011
- `/home/kesokaj/workspace/styrman-harness/.agents/docs/deploy/local.md` and `remote.md` — how the harness reaches a machine

## Suggested Skills

- `handoff` — read this first on the Mac session
- `feedback` then `/retro` — after two weeks of use
- `golang` — for any change under `internal/harness/`

## Next Steps

1. On the Mac: `cd ~/workspace/styrman-harness && sh-agent update && sh-agent status`, start `claude` and `agy` once, run `/skills` and `/agents`, ask "which slash command grills a plan" (expect `/grilling`); tick plan section 12 and `TODO.md`.
2. Use the harness on a real project for two weeks; log corrections as they happen.
3. Then `/retro`, fold accepted proposals, and re-run the matrix if `context/` or `agents/` changed.
