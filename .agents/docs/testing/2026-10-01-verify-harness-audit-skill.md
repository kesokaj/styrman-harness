---
owner: main
date: 2026-10-01
---
# Verification — the `harness-audit` skill on both CLIs, and the catalog rows for v2.2.0 (2026-10-01)

Range reviewed: v2.1.0..f3c3ee4 — the skill's text is that of f3c3ee4 (installed into a fresh sandbox home from the checkout at that commit; `diff -rq skills/harness-audit` against both vendors' installed copies empty). Plan `plans/2026-10-01-gate-lows-and-harness-audit.md`, E2E criteria (c), (d), (e); rows `harness-audit`, `harness-audit-elsewhere`, `harness-audit-reaudit` of `testing/skill-matrix.md`.

## Status: PASS — 5 of 5 rows at f3c3ee4 (Claude Code 3 of 3, Antigravity 2 of 2); 15 of 15 over the three commits the rows ran on

Each row is a fresh clone (of this repository for the audit rows, of the demo service for the "elsewhere" row), the kit preamble sourced first, the run saving the CLI's output and the check re-reading the clone and that output from disk. The area agent was given 20 tool calls (25 for the re-audit). Claude Code ran with an allow-list (`Read Edit Write Agent Glob Grep Skill` and named shell commands); Antigravity with `command(*)`, `read_file(*)`, `write_file(*)`: the audit of area AG reads no vendored or third-party text, and the skill needs to run the project's own checks.

| CLI | row | result | what the check found | denied | seconds |
|---|---|---|---|---|---|
| Claude Code | elsewhere (`/harness-audit --area AG` in the demo service) | pass=1 | paths_written=0 said_why=3 | denied 0 | 6 |
| Claude Code | first audit (`/harness-audit --area AG`) | pass=1 | finding_rows=9 report=.agents/docs/reviews/2026-10-01-harness-audit.md plan_phases=7 changed_outside_docs=[] | denied 0 | 487 |
| Claude Code | re-audit (`/harness-audit --reaudit plans/2026-09-30-recon-audit-fix.md --area AG`) | pass=1 | report=.agents/docs/reviews/2026-10-01-audit-4-AG.md ids=15 missing=[] changed_outside_docs=[] | denied 0 | 390 |
| Antigravity | elsewhere | pass=1 | paths_written=0 said_why=1 | denied 0 | 15 |
| Antigravity | first audit | pass=1 | finding_rows=12 report=.agents/docs/reviews/2026-10-01-harness-audit.md plan_phases=5 changed_outside_docs=[] | denied 0 | 439 |

What "pass" means per row: elsewhere — `git status --porcelain --untracked-files=all` is empty in the clone and the answer says it is not a harness checkout; first audit — a recon brief with at least one row carrying an `AG<n>` ID, a severity and a `path:line`, a report, an action plan with at least one phase, and nothing changed outside `.agents/docs/`; re-audit — `reviews/<date>-audit-4-AG.md` with "Status of every ID" and a row for each of the 15 AG IDs of the named brief.

The same five rows ran on the two earlier commits of the gate, on the skill text of each: 5 of 5 at a62eb9d, 5 of 5 at 768efce (`~/sandbox/audit-fix/ha-results-a62eb9d.txt`, `ha-results-768efce.txt`). The audits were not make-believe: the first Claude Code run numbered its findings on from AG15, the last ID of the 2026-09-30 audit, and reported that four tests of the dream package fail when `CLAUDE_CONFIG_DIR` is set in the caller's environment, which was true and is fixed in 768efce. The clones with those briefs, reports and plans are in the sandbox (`~/sandbox/audit-fix/runs/ha-*`), not in this repository: a real audit is the operator's to start.

## Catalog rows

The range changes two catalog notes and the refresh script; the router is unchanged. The router's rows ran at 768efce, whose router, manifest and index are what f3c3ee4 ships (`git diff 768efce..f3c3ee4 --stat -- skills/google-catalog` names only the script): Claude Code 9 of 9 (`catalog-no-leak`, seven routed tasks, `catalog-precedence`), Antigravity 8 of 8 routed and precedence rows plus the no-leak cell by two identical self-listings of 38 names (36 harness skills, `harness-audit` among them, and its two built-ins; no catalog entry). The manifest note was read before the entry in 14 of 14 routed runs.

Credentials copied into the sandbox homes were deleted by name after each batch; `find` over the sandbox finds none.

## Not checked

`--reaudit` on Antigravity (the first-audit row ran there). A full audit over all eight areas. The audits' findings themselves were not judged here beyond the one that was reproduced.
