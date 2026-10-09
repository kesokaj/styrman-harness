---
owner: explorer
date: 2026-10-09
---
# Recon - project   (2026-10-09)

## Status: COMPLETE

Range reviewed: the whole repository at HEAD (c267990), no diff.

## Question
Inventory and structure map: uniformity, orphaned and dead-code candidates (six deliverables from the dispatch).

## Summary
Layout in AGENTS.md matches disk except: `internal/trail` (3 files, 517 lines) exists and is not named; root files `PRIORITY_TODO.md`, `embed.go`, `embed_test.go`, `init_project_test.go`, `plan_test.go`, `testmain_test.go`, `go.sum`, `README.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md` are not in the Layout block (only README/REQUIREMENTS/embed are named elsewhere). `PRIORITY_TODO.md` (1 line, a header) is almost certainly a leftover: commit 1b8e89b said "PRIORITY_TODO emptied and its rule removed" but the file was kept; nothing outside dated trail records mentions it. No dead files or unused go.mod requirements found; dead-code candidates are ~28 exported identifiers used only in their defining file (could be unexported). Trail is healthy: 166 bug records = 166 INDEX rows, one open bug (BUG-139), no plan with real unticked boxes. `.agents/scratch/` holds ~200 untracked-by-design leftovers (gitignored).

## 1. Inventory
Root (not named in Layout): AGENTS.md (17,396 B, named by definition), README.md (48,625 B), REQUIREMENTS.md (65 lines), ROADMAP.md (10 lines; "Last reviewed: 2026-10-01", Now and Next empty), TODO.md (5 lines, empty), PRIORITY_TODO.md (1 line: "# PRIORITY TODO - audit fix (operator, 2026-09-30)"), embed.go (the `styrman` package, `//go:embed all:context all:skills all:agents all:commands all:rules`), root tests embed_test.go, init_project_test.go, plan_test.go, testmain_test.go (package at repo root, 4 test files), go.mod/go.sum, .gitignore.
Dirs, all named in Layout: context/ (AGENTS.md, dev-setup.md), agents/ (8 files), skills/ (38 entries: 36 skills + VENDORED.md + google-catalog catalog; 885 tracked files), commands/ (only .gitkeep, by design), rules/ (critical-thinking.md + .gitkeep), cmd/sh-agent (7 files, 4,586 lines), internal/ (89 tracked files), .agents/docs (1,523 tracked), .agents/scratch (gitignored, 1 tracked .gitkeep), .agents/skills (symlink to ../skills, tracked as a symlink).
Layout names `internal/{harness,search,memory,dream,metrics}` only; disk also has `internal/trail`. Layout says `.agents/docs/` has "plans, changes, ADRs, reviews, testing, deploy/{local,remote,gcp}.md, bugs/INDEX.md"; disk also has handoffs/, retros/, specs/sh-agent/, architecture/DECISIONS.md (ADR home). Not on disk but named: none found.
.gitignore: `/.claude/` and `/agent/` have no matching path today (fine, they are created by tools); `.agents/memory` and `.claude/settings.local.json` likewise absent. Secret-shaped block is pinned to plan.go secretShaped() per its comment.

## 2. Structure map of the binary (counts: go files / lines incl. tests)
- internal/harness 35 / 14,073; ~70 exported funcs/methods; imported by cmd/sh-agent and the root package tests. Installer, update, report, release, zip, state, lint, derive.
- internal/search 14 / 5,068; ~12 exported; imported by cmd, metrics (agy.go, claude.go), dream (collect.go). Owns SQLite DSN helpers (`search.DSN` used by metrics/store.go:71).
- internal/memory 3 / 762; 6 exported; imported by cmd, dream.
- internal/dream 12 / 3,750; 11 exported; imported by cmd only.
- internal/metrics 21 / 5,009; 16 exported; imported by cmd, dream.
- internal/trail 3 / 517; 6 exported; imported by dream/collect.go, metrics/events.go, memory/memory.go, search/reconcile.go, search/search.go. No SQLite.
Flow: cmd/sh-agent/main.go:84 main() -> verb table (one FlagSet per verb; entries at main.go ~146-200: version, uninstall, report, memory, dream, metrics, search, index, install, sync, status, diff, update, list, render, adopt) -> `run` func or inline closure. install/sync -> runInstall -> harness plan/apply (plan.go, apply.go, derive.go); update -> repo.go/release.go/zip.go; report -> harness/report.go; search/index -> search.Search/Reconcile (git-reconciled SQLite); memory -> memory add/show; dream -> dream.Run (collect via search + memory, run CLI, timer install); metrics -> ingest/show/report via metrics store. (Per-verb leaf tracing below main.go's table was not read line by line.)

## 3. Uniformity
- Errors: `%w` wrapping is concentrated in harness (apply 7, release 7, repo 6, zip 6, report 10, inventory 7); dream, memory, search, metrics wrap sparsely (dream/state.go 1, dream/collect.go 1); metrics and search files show none by this grep. Uneven but not contradictory.
- log/slog: used only in internal/harness/harness_test.go; the CLI prints with fmt to stderr (main.go "sh-agent:" prefix). Global context says Go services use slog JSON; this is a CLI, so divergence is by kind, not a defect. Say so if a plan asks for uniform logging.
- Home resolution: os.UserHomeDir/HOME only in cmd/sh-agent/main.go, internal/harness/pathguard.go, internal/dream/run.go (dream also needs the agy private HOME). No other package reads HOME directly; they take a home argument. Good.
- Atomic writes: os.Rename in harness/apply.go and harness/adopt.go only; os.WriteFile in harness/report.go and dream/run.go (non-atomic). memory and metrics writes not located by this grep (SQLite or other helper); unverified.
- SQLite: sql.Open at internal/metrics/store.go:71 (via search.DSN), internal/search/store.go:187, and internal/search/dream.go:266 (read-only DSN on copies). Matches AGENTS.md (search and metrics are the driver importers; dream goes through search). Close discipline not audited.
- Tests: each of 8 test dirs has its own testmain_test.go (dream uses main_test.go): the root, cmd, harness, metrics, trail, search, memory (dream differs in name). 38 test files use t.TempDir or a home helper.
- Frontmatter: agents all `description:` double-quoted, tools lists (bugs/changelog/documentation: Read, Grep, Glob, Edit, Write; dreamer/explorer/reviewer/security/verifier: Read, Grep, Glob, Bash, Write), models sonnet x6, haiku (changelog), opus (security); maxTurns 20/30/40/40/40/40/40/45 (bugs, changelog, documentation, dreamer, explorer, reviewer, security, verifier). Skills: no non-catalog SKILL.md uses `disable-model-invocation`; description quoting mixed (some quoted with `"`, some plain; first-char scan was inconclusive).

## 4. Orphan and dead-code candidates
Exported but referenced only in their defining file (no other non-test file): harness Alive, AssetName, BuildRecord, Canary, Collect, Collection, Doc, Keys, LatestRelease, Lint, LoadAgents, LoadDocs, OriginKind, ParseOrigin, ReadSettings, Record, Redaction, ReleaseRepo, ResolveFeedback, RunningElsewhere; search AgyConversation, IsBusy; memory Refusal; metrics Adherence, Efficiency, ErrNewerSchema, Reliability, RenderReport. Many are used by tests (e.g. LoadAgents 9) and some are method names (Collect, Keys, Alive) where the per-name grep over-matches. Treat as "could be unexported", not dead: LatestRelease, LoadDocs, OriginKind, AgyConversation, IsBusy, Refusal have 0 test refs and no cross-file refs - best true candidates (verify with `go vet`/staticcheck-like check before acting).
Files no package imports: none found; every internal package is imported by cmd or another package. go.mod: only `modernc.org/sqlite` direct; indirect set is what sqlite pulls (go-humanize, uuid, isatty etc. not imported by repo code, as expected).
Agents with no dispatcher: none; all 8 appear in context/AGENTS.md (dreamer, explorer once each) and in 4-12 skills.
Skills nothing names (no context mention, no other skill mentions): explain-project, find-skills, google-catalog (router; named in Installed list at install time), harness-audit. They are triggered by description only; likely deliberate.
Scripts: kit-preamble.sh, seed-security/run.sh, ab-kit/{run.sh,accept.py}, agy-rel/run.sh, review/scripts/restatement-grep.sh, google-catalog/scripts/refresh-catalog.sh all have other referrers (basename match, so run.sh counts are inflated: 81 refs is the shared name). accept.py has 7.
Docs: not all untracked-by-link docs checked (see Not checked). `.agents/docs/testing/` has 349 files, most dated records; undated: guide.md, sh-agent-matrix.md, skill-matrix.md, kit-preamble.sh.

## 5. Trail health (tracked files)
architecture 1 (DECISIONS.md), bugs 167 (166 records + INDEX; BUG-001 .. BUG-166), changes 264 (2026-09-03 .. 2026-10-09-1651), deploy 3, handoffs 114 (2026-09-03 .. 2026-10-09-1652), plans 79 (22 recon/research briefs; 2026-09-03 .. 2026-10-08), retros 16 (2026-09-14 .. 2026-09-30; none since), reviews 525, specs 5 (specs/sh-agent), testing 349. Uncommitted: .agents/docs/testing/2026-10-09-verify-project.md (untracked, from another agent, not mine).
Open: `sh-agent search --open` reports 1 open record: BUG-139 (documentation agent roadmap lane held in 5 of 6 runs, open since 2026-10-01). Plans: only 2026-09-15-harness-practices-g01-g26.md matches `- [ ]` (4 hits; the matches at lines 180-183 are a template inside the plan, not work; ROADMAP does not name it). INDEX vs record: 166 records and 166 rows, ids all present; status agreement not conclusively checked (graduated rows use a different column layout, my check was inconclusive).
ROADMAP.md "Last reviewed: 2026-10-01" is 8 days old with Now empty - consistent with no active plan.

## 6. Risks and open questions for a fix plan
- PRIORITY_TODO.md: leftover (1b8e89b says its rule was removed) or deliberate? Deleting it is a root-file change; also check no test pins it.
- Is internal/trail meant to be in AGENTS.md Layout? (It exists since plan 2026-10-08-metrics-store D31.) Adding is a docs-only edit but AGENTS.md is the project rule file.
- Root-level Go tests (init_project_test.go, plan_test.go, embed_test.go, testmain_test.go) sit in package `styrman` beside embed.go: deliberate (they test embedded content) but unnamed in Layout.
- Unexporting the six zero-test candidates is an API change within internal/ only; low risk but staticcheck was not run, so confirm first.
- .agents/scratch holds ~200 untracked artifacts (release dirs, logs, search.db); gitignored, not a repo concern, but disk clutter.
- retros stopped 2026-09-30 despite many corrections since? Not determined.
- No skill uses `disable-model-invocation`; unclear whether intended (TestSkillListingBudget is the stated mechanism).

## Not checked / Not examined
Per-verb call chains below main.go's verb table; full exported-surface listing via `go doc`; error-wrapping and atomic-write audit for memory, metrics, search, dream beyond a grep; SQLite Close discipline; skill description quoting per file; INDEX status vs record status (check inconclusive); undated docs under .agents/docs with no inbound link; test naming conventions; vendored catalog contents; staticcheck/deadcode tooling (not run).
