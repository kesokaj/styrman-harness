---
owner: explorer
date: 2026-09-30
---
# Recon — project   (2026-09-30)

## Status: COMPLETE

Range reviewed: the whole repository at HEAD (3478a1c), no diff.

## Question
Inventory and structure map of the repository: uniformity, orphaned and dead-code candidates.

## Summary
The Go tree (`cmd/sh-agent`, `internal/harness`, `internal/search`) shows no clear dead code from a quick exported-symbol sweep: every symbol that looked isolated (`AssetName`, `BuildRecord`, `Collect`, `LatestRelease`, `LoadAgents`, `LoadDocs`, `ResolveFeedback`, `RunningElsewhere`, `SwapReleaseSigner`) turned out to have a same-file, same-package, or cross-package test caller — this needs `go vet`/staticcheck-grade tooling to say more definitively (not run here). Skills are fully accounted for: 28 directories under `skills/`, matching `context/AGENTS.md`'s "Skills (28)" line and `skills/VENDORED.md`'s provenance table, so no orphaned skill. `commands/` and `rules/` are real but **empty** (each holds only a `.gitkeep`), even though `AGENTS.md` layout table describes them as holding "markdown prompts... deploy[ed] to both vendors verbatim" — this is either intentional scaffolding-for-future-use or a stale layout line; worth a one-line check with the operator. Agent frontmatter is fairly uniform (`name`, `description`, `tools`, `model`, `effort`, `maxTurns`) but `dreamer.md` has no visible `maxTurns` line in the excerpt read and `reviewer.md` alone carries a `skills:` list — not a defect, just the one agent with a skill dependency. Skill frontmatter is deliberately non-uniform: vendored skills (`anti-ui-slop`, `gke-basics`, `terraform`) carry `license`/`metadata` blocks per `VENDORED.md`'s own rule that provenance lives in the skill's own frontmatter, while home-grown skills carry only `name`/`description` (+ occasionally `argument-hint` or `allowed-tools`) — this is a convention, not an inconsistency to fix. Trail orphaning is largely explained by project convention: `ROADMAP.md` "Now"/"Next"/"Later" are all empty (the project's stated rule: empty at every checkpoint, work is closed elsewhere), so the fact that most of the 60 files in `.agents/docs/plans/` are not named by ROADMAP is expected, not a smell — genuine orphan-detection would need to check each plan against handoffs/changes/bug records, which was not done exhaustively here (see Not examined).

## Relevant files
- `context/AGENTS.md` — canonical "Installed" list (28 skills, 8 agents) used as the cross-check baseline against `skills/` and `agents/`.
- `skills/VENDORED.md` — provenance manifest for 4 table-listed vendored skills plus 2 Apache-2.0 skills (`gke-basics`, `terraform`) attributed in their own frontmatter; also documents two *pruned* skills (`grill-me`, `skill-repair`) as historical context, not currently present.
- `agents/*.md` (8 files) — Claude-format subagent definitions; Antigravity's agent.md is derived from these at install (per project AGENTS.md).
- `internal/harness/inventory.go` — `LoadAgents`, `LoadDocs` read `agents/*.md`, `commands/`, `workflows/`, `rules/` directly; confirms `commands/`/`rules/` are wired into the install pipeline even though currently empty.
- `internal/harness/release.go`, `report.go`, `pathguard.go` — home of the exported symbols checked for dead-code candidacy; all have live callers.
- `.agents/docs/plans/` (60 files), `.agents/docs/handoffs/` — trail; `ROADMAP.md` empty by design.
- `commands/`, `rules/` — both contain only `.gitkeep` (empty since at least commit `fc89509`, "Generic installer directories, hooks research, no attribution").

## Flows
- Install pipeline: `internal/harness/inventory.go` (`LoadAgents`, `LoadDocs`) reads `agents/`, `commands/`, `workflows/`, `rules/`, `skills/` off a `Source` (repo checkout, embedded zip, or plain dir) → `plan.go` decides what/where → `apply.go`/`derive.go` write Claude Code files and derive Antigravity's agent.md → `render.go` assembles `context/AGENTS.md` + `context/dev-setup.md` into the rendered context checked by `TestRenderedContextFitsAntigravityBlock` (`internal/harness/size_test.go`).
- Release/report: `release.go` (`AssetName`, `LatestRelease`, `SwapReleaseSigner`, `verifyRelease`) backs `/preflight`'s signed-tag flow; `report.go` (`Collect`, `BuildRecord`, `ResolveFeedback`) backs `sh-agent report`'s throwaway-clone push to the feedback repo — both exercised primarily through their own package's tests plus `cmd/sh-agent` integration tests (`zip_test.go`, `search_test.go`).
- Search: `internal/search/{search,reconcile,store,parse}.go` is the one package importing `modernc.org/sqlite` (per AGENTS.md's stdlib-only-with-one-exception rule); `cmd/sh-agent` wires the `search`/`index` verbs to it.

## Conventions
- Skill/agent addition is a file drop, never a code change (project AGENTS.md, "Rules for working here") — so an "orphan" skill would be a directory present but absent from `context/AGENTS.md`'s Installed count; none was found (28 = 28).
- Vendored skills keep upstream frontmatter verbatim (license/metadata blocks); home-grown skills use a lean `name`/`description` (+ `argument-hint` for slash-invoked ones, `allowed-tools` for `bash`). This split is intentional per `VENDORED.md`, not inconsistency.
- Agent frontmatter fields in common use: `name`, `description` (one line, double-quoted when it contains a colon per AGENTS.md rule), `tools`, `model`, `effort`, `maxTurns`; `reviewer` additionally declares `skills: [review-standards]`.
- Trail files are append-only and dated `$(date +%F[-%H%M])-<slug>.md`; ROADMAP/TODO are kept empty at every checkpoint by rule, so "not named by ROADMAP" is not itself evidence of an orphan in this repo.

## Risks
- `commands/` and `rules/` are empty but load-bearing in the install code path (`LoadDocs` reads them unconditionally); an empty-by-design directory that the layout table describes as holding content is a documentation/reality mismatch worth flagging to the operator, though functionally harmless (empty dirs just mean zero commands/rules deployed).
- The dead-code sweep here was a shallow heuristic (exported symbols, same-package/test-caller check only); it does not cover unexported helpers, nor does it run `staticcheck`/`deadcode` tooling. A symbol used only by its own package's tests (no non-test caller) is common in this codebase's TDD style and was not flagged as dead, but a stricter definition might disagree — this is a judgment call, not a settled fact.
- Trail-orphan detection was done at the ROADMAP level only; individual plan files were not each checked against handoffs/changes/bug INDEX for a lingering reference, so some of the 60 plan files could be genuinely stale (never closed) vs. closed-but-unlisted (normal). Not distinguished here.

## Open questions
- Is the `commands/`/`rules/` emptiness intentional scaffolding for future content, or should the AGENTS.md layout table be corrected to say they are currently unused? (Operator call.)
- Should dead-code detection be re-run with `staticcheck ./...` or `deadcode ./...` for a authoritative (not heuristic) answer on `internal/harness` and `internal/search`?
- Among the 60 `.agents/docs/plans/` files, which (if any) were abandoned mid-way rather than closed — this recon only confirmed ROADMAP is empty, not that every plan's checkboxes are ticked.

## Not examined
- No `staticcheck`/`deadcode`/`go vet -unusedresult` style tooling was run; the dead-code claim rests on manual grep only.
- Unexported (lowercase) function call graphs in `internal/harness` and `internal/search` were not swept for dead code — only exported symbols were checked.
- Each of the 60 files in `.agents/docs/plans/` was not individually opened to confirm its checkboxes are all ticked or that it is referenced by a `changes/` entry or bug record.
- `.agents/docs/changes/`, `.agents/docs/reviews/`, `.agents/docs/testing/`, `.agents/docs/bugs/` directories were not inventoried for orphans (only `plans/` and `handoffs/` were sampled).
- Full text of `dreamer.md` past line 8 was not re-checked to confirm whether it truly lacks `maxTurns` (only the first 8 lines were read for all agents; `dreamer`'s excerpt ended without one, but the file may continue past what was shown — verify before reporting as a finding).
- `context/dev-setup.md` and the full contents of every `SKILL.md` body (beyond frontmatter) were not read.
