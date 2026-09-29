---
owner: explorer
date: 2026-09-29
---
# Recon — project/global SQLite index feature   (2026-09-29)

## Status: COMPLETE

## Question
Recon for adding `sh-agent index` / `sh-agent search` (pure-Go modernc.org/sqlite + FTS5), a
per-project index (`.agents/scratch/index.db`) and a global index (`~/.config/sh-agent/index.db`),
plus exposure via `context/AGENTS.md` (22,000-byte budget) and a skill.

## Summary
The biggest blocker is not mechanical: R4 in REQUIREMENTS.md ("`sh-agent` is a single Go binary
with no dependencies") and the project AGENTS.md ("Stdlib only in the Go code. The binary is a
file copier with a manifest; it does not need a dependency") both make `modernc.org/sqlite` a
REQUIREMENTS.md change requiring the operator's explicit approval + an ADR before any code lands —
this is not a "just add the import" task. Mechanically the rest is straightforward: verb dispatch
is a flat switch in `cmd/sh-agent/main.go` `run()` sharing one `flag.FlagSet`; `TestHarnessMap`
(cmd/sh-agent/main_test.go:1350) only requires that verbs the mermaid map in README.md names exist
in `usage`, and that install/update/report appear in the map — adding `index`/`search` verbs does
NOT by itself force a README-map edit unless the map's prose is changed to mention them (they are
not one of the four R9 procedures). `~/.config/sh-agent/` already holds state.json/feedback.jsonl/
retros.jsonl/proposals/ with no project registry; `report.go`'s `Collect()` is a closed allowlist
of exactly those names, so a new `~/.config/sh-agent/index.db` is automatically NOT sent by
`sh-agent report` — good, but worth an explicit skip-list entry/comment for defense in depth.
CGO/cross-compile is not a blocker: the release build already sets `CGO_ENABLED=0` for
linux/amd64, darwin/arm64, linux/arm64 (deploy/remote.md:141), which is exactly why modernc.org's
pure-Go driver (not mattn/go-sqlite3) is the right choice if a dependency is approved — but go.mod
currently has no `go.sum` at all (zero deps), and `TestFilesMatchTree`/`TestEmbedHoldsOnlyTrackedFiles`
(embed_test.go) plus the release contamination guard only check that embedded context/skills/etc
files match git-tracked ones — they do not police go.mod, so no test currently blocks adding a
dependency, but doing so is still an R4 violation on its face. The context-size budget is tight:
`context/AGENTS.md` + `context/dev-setup.md` are already 21,031 bytes combined against a hard
22,000-byte ceiling pinned by `TestRenderedContextFitsAntigravityBlock`
(internal/harness/size_test.go:19, comment: "never raised by an implementer... a rule that does
not fit is paid for by a deletion or a merge, or lives in a skill") — under 1,000 bytes of margin,
so "a few lines" for index/search usage will need real pruning elsewhere in context/AGENTS.md, and
note context/AGENTS.md IS the literal global context (same text shown to every session via
`~/.claude/AGENTS.md` → this repo), so any edit there is global-blast-radius, not scoped to this
feature's own project.

## Relevant files
- `cmd/sh-agent/main.go:1-20` doc comment and `:38-53` `usage` const — BOTH enumerate every verb
  and must be edited together for `index`/`search`; not test-enforced against each other (only the
  README mermaid map is pinned to `usage`, not the doc comment).
- `cmd/sh-agent/main.go` `run()` (~86 onward) — one shared `flag.FlagSet` parses `--repo`, `--home`,
  `--vendor`, `--dry-run`, `--check`, `--force`, `--feedback-repo`, `--from-zip` for every command
  before the verb switch; then `switch cmd { case "version": ...; case "uninstall": ...;
  case "report": ...}` handles verbs needing no `harness.ResolveSource`, followed by
  `src, err := harness.ResolveSource(...)` and a second `switch cmd { case "install","sync": ... }`
  for verbs that need a resolved harness tree. A `--query`/`--limit` flag for `search` would sit in
  this same global FlagSet (visible to every verb) unless a per-command sub-parsing pattern is
  added — none exists today; every flag is global. Positional args after the verb (e.g. the search
  query string) aren't handled anywhere yet — `reorder(rest)` (main.go, referenced at parse time)
  reorders flags ahead of positionals (BUG-002's fix), so a bare positional query arg is plausible
  but unprecedented in this codebase; would need new code, not just a new case.
- `cmd/sh-agent/main_test.go:1345-1377` `TestHarnessMap` — reads README.md `## How it works` §,
  extracts the mermaid block, regexes `sh-agent ([a-z]+)`, requires each named verb to appear in
  `usage` as `"\n  "+verb+" "`, and requires install/update/report to appear in the map. New verbs
  do not force a map edit unless the map's own prose grows to mention them.
- `internal/harness/state.go:71-137` — `ConfigDir(home)` = `~/.config/sh-agent` (created 0700 by
  `TightenConfigDir`, called from `apply.go:30` at install time). Comment at 71-75 currently
  enumerates what lives there ("the feedback log, the retro proposals and, briefly, a clone of the
  private feedback repo") — needs updating for a new `index.db`. Comment at 132-137 already uses
  the phrase "cross-project index ~/.config/sh-agent/retros.jsonl" for a DIFFERENT thing (a
  retros.jsonl-derived concept referenced by `harnessProposalInputs` or similar, feeding
  `/retro --harness`) — a literal `index.db` should not overload this existing term; reconcile
  naming with the operator before writing code/docs.
- `internal/harness/report.go:220-295` `Collected`/`Collection`/`Collect(home)` — closed allowlist:
  `feedback.jsonl`, `retros.jsonl` (both via `redactLines(...,true)`/jsonl mode), `state.json` (via
  `redactDocument`), and `proposals/*.md` (via text redaction). A new `index.db` at
  `~/.config/sh-agent/` is NOT in this list and thus never collected/sent by `sh-agent report` —
  matches the scope's "must not be sent" requirement by construction. No code change needed to
  exclude it, but recommend adding a comment/explicit skip entry (`c.Skipped = append(...,
  "index.db (binary index, never sent)")`) for auditability, matching the pattern already used for
  absent/non-regular files.
- `go.mod` — `module github.com/kesokaj/styrman-harness`, `go 1.24.4`, no `require` block, **no
  go.sum file exists** (zero dependencies today). Per project AGENTS.md, no `toolchain` directive
  (ADR-029, `.agents/docs/deploy/local.md:101` — "a pin makes the module unbuildable on a host that
  cannot fetch a toolchain... never raise the `go` directive or add a `toolchain` line").
- `.agents/docs/deploy/remote.md:139-144` — release cross-compile: `for target in linux/amd64
  darwin/arm64 linux/arm64; do GOOS=... GOARCH=... CGO_ENABLED=0 GOTOOLCHAIN=local go build
  -trimpath ...`. CGO is already off for every release target, which is consistent with (and
  arguably requires) a pure-Go SQLite driver like modernc.org/sqlite if one is approved — a cgo
  driver (mattn/go-sqlite3) would break this build outright.
- `embed_test.go` (repo root, package `styrman`) — `TestFilesMatchTree` and (grep hit only, not
  read in full) `TestEmbedHoldsOnlyTrackedFiles` police that `context/`, `skills/`, `agents/`,
  `commands/`, `rules/` on disk exactly match what's embedded — irrelevant to go.mod dependencies,
  no test polices "no dependencies" mechanically; that constraint is REQUIREMENTS.md-only.
- `internal/harness/size_test.go:9-49` — `antigravityRulesBlockBytes = 24000` (Antigravity's real
  cutoff, BUG-085), `renderedContextBudgetBytes = 22000` (this repo's self-imposed ceiling, ADR-035,
  8% margin, "never raised by an implementer"). `TestRenderedContextFitsAntigravityBlock` sums
  `RenderContext(...)` output + every `rules/<name>.md` raw byte length and fails over budget.
- `context/AGENTS.md` (18,448 bytes) + `context/dev-setup.md` (2,583 bytes) = 21,031 bytes measured
  by `wc -c` (a rough proxy for the rendered/summed total the test checks; not identical since
  `RenderContext` may add wrapper bytes and rules/*.md also count) — margin to 22,000 is under
  ~1,000 bytes on this proxy, likely tighter once rules/*.md are included. "A few lines" of
  index/search guidance will require pruning elsewhere.

## Longest paragraphs in context/AGENTS.md (candidates to shrink if room is needed)
- Line 33 (~1,500 chars): "Skill = knowledge the main thread loads while working. Subagent =
  judgement independent of the author..." — dense, multi-topic (skill vs subagent, per-dispatch
  budget repetition rule, TaskStop handling).
- Line 51 (~1,450 chars): the git hygiene hard rule (heredoc/numstat, escape-character check,
  `rm -r` guard, post-delete/dispatch stop-sign checks) — four distinct sub-rules crammed into one
  bullet.
- Line 55 (~1,350 chars): "A behaviour change includes the text that DESCRIBES it..." — the
  description/coverage/vocabulary/claim-reversal rule, itself already dense.
- Line 57 (~1,050 chars): the "Silence from your own instrumentation is not evidence" mutation/
  monitor-validation rule, cross-references `review-standards` § 2 (partially already delegated to
  a skill — a template for delegating more).
- Line 5 (~750 chars): the operator-preferences paragraph (SHORT LIST / ONE NUMBERED TABLE rules) —
  stable, unlikely to be a pruning target since it governs interaction style broadly.
Not read in detail against actual redundancy — a real pruning pass needs the operator's read, not
just length ranking; flagging candidates only, per scope item 5's ask.

## Flows
- Verb dispatch: `os.Args[1:]` → `run(args, out)` → `cmd, rest := args[0], args[1:]` → shared
  `flag.FlagSet` parses `reorder(rest)` → `--from-zip` special-cased early (return before the main
  switch) → switch on `cmd` for verbs not needing a resolved source (`version`, `uninstall`,
  `report`) → `harness.ResolveSource(*repoFlag, home, buildVersion())` → second switch for
  `install`/`sync`, and (not fully traced past line ~180) presumably further cases for `status`,
  `diff`, `update`, `list`, `render`, `adopt`. New `index`/`search` verbs would each need a case in
  one of these two switches (likely the first, since building/querying an index needs no vendor
  resolution) plus a handler function mirroring existing ones (e.g. `report()` at main.go, `report`
  case).
- Report/redaction flow: `sh-agent report` → `Collect(home)` (report.go:235) reads a fixed name
  list under `~/.config/sh-agent/`, redacts each, builds a `Collection` → (not traced further)
  presumably pushed via a throwaway git clone per REQUIREMENTS R8/R9. A global `index.db` sits
  beside these files but outside the allowlist, so it is structurally excluded without code changes
  — confirm this is deliberate and durable (e.g. a comment or explicit test) rather than incidental.

## Conventions
- Every verb's usage line, doc comment, and (if it appears in the README mermaid map) the map
  itself are kept in sync by convention/test, not fully cross-tested; `TestHarnessMap` only checks
  map→usage, not usage→doc-comment or doc-comment→map.
- New `~/.config/sh-agent/` content gets documented in `ConfigDir`'s doc comment (state.go:71-75)
  and, if it should never be reported, an explicit `Skipped` entry pattern is already used in
  `Collect()` for absent/non-regular files — mirror that for a deliberately-excluded `index.db`.
- Trail record formats (from bugs/INDEX.md, a BUG file, and DECISIONS.md headings observed):
  - `bugs/INDEX.md` table: `| ID | Date | Area | Status | Symptom | Prevention |` with rows like
    `| BUG-105 | 2026-09-29 | skills/preflight (release gate verdict) | open | <symptom> |
    <prevention> |`; a second "Graduated" table exists elsewhere in the same file (not fully read)
    for prevention rules that made it into AGENTS.md/a rule/skill/agent/test.
  - BUG file frontmatter: `owner: bugs`, `date: 2026-09-29`, `id: BUG-105`, `status: open`,
    `area: skills/preflight (release gate verdict)`; body sections `## Symptom`, `## Repro`,
    `## Root cause`, `## Fix`.
  - `DECISIONS.md` ADR heading pattern: `## ADR-NNN: <title>     (<date>) — <Status>`, statuses
    seen: `Accepted`, `Superseded by ADR-041`, `Accepted; its one-file refusal superseded by
    ADR-041`, `Accepted; superseded in part by ADR-044 (...)`, `Accepted, amends ADR-042`.
  - `handoffs/` and `changes/` filenames: `$(date +%F-%H%M)-<slug>.md` per project AGENTS.md.
  - Trail artifacts consistently start with YAML frontmatter `owner:` + `date:`.

## Risks / constraints (REQUIREMENTS.md, read in full for R1–R9 + Known limits + Non-goals)
- **R4 (blocking):** "`sh-agent` is a single Go binary with no dependencies." Project AGENTS.md:
  "Stdlib only in the Go code. The binary is a file copier with a manifest; it does not need a
  dependency." Adding `modernc.org/sqlite` directly contradicts both. Per the global rule,
  REQUIREMENTS.md is "immutable without the operator's explicit approval": propose the exact edit →
  approval → ADR → update REQUIREMENTS.md and ROADMAP.md in one commit, BEFORE implementation.
  This is the single largest open item — flag to the operator as a SHORT LIST choice (amend R4 vs.
  find a stdlib-only approach vs. defer) rather than routing around it.
- **R7:** "The harness's own state (manifest, feedback log, the harness proposal queue, the
  cross-project index) lives in `~/.config/sh-agent/` and nowhere else." A global `index.db` there
  is consistent with R7's location rule in spirit, but R7 already uses the term "the cross-project
  index" for something else (apparently the retros.jsonl-derived concept referenced in state.go)
  — naming needs reconciling with the operator, not necessarily blocking.
- **R8:** "no vendored mirror inside a consumer project, no per-project copy, no service." The
  per-project `.agents/scratch/index.db` is scratch — gitignored, ephemeral, explicitly rebuildable
  from git per the feature scope ("reconciled against `git ls-files -s` + `git status` before every
  query") — this reads as consistent with R8 (it is not a copy of harness state, it is a derived
  cache of the project's own trail), but confirm with the operator since R8 is specifically about
  the harness's OWN delivery channel, and this feature is adjacent but not identical territory.
- **Project AGENTS.md rule 6** ("Changing a path is a design question... lands in ONE commit with
  the code, the docs and the map"): this rule is about the delivery map (install/update/report/
  release paths), not about adding a verb; `index`/`search` do not appear to be "a line on the map"
  under R9's four procedures, so the map likely does not need editing for the verbs themselves —
  only if the feature's rollout is later folded into report/update (unlikely, out of scope here).
- **Known limits / renderedContextBudgetBytes:** the 22,000-byte ceiling is explicitly "never
  raised by an implementer" (size_test.go:17-18 comment) — a hard constraint, not a target to
  negotiate; any context/AGENTS.md addition is paid for by deletion/merge/moving text to a skill.
- **Blast radius of editing context/AGENTS.md:** this file (verified identical in content to the
  `~/.claude/AGENTS.md` shown for THIS session) is the single global context deployed to every
  project via `sh-agent install`; editing it to mention index/search changes behavior for every
  project using the harness, not just this feature's target repo. Treat any edit there as a
  "Feature" under the project's own gates (design question), and prefer putting most detail in a
  skill (`skills/<name>/SKILL.md`) with only a short trigger line in context/AGENTS.md — this
  matches the existing pattern (`review-standards` is already referenced this way from a dense
  hard rule at line 57).

## Not examined
- Full `cmd/sh-agent/main.go` switch statement past ~line 180 (status/diff/update/list/render/adopt
  cases) — only install/sync case partially read; a new-verb author should read the rest before
  writing a case, but the pattern (case → dedicated handler function, e.g. `report()`) is clear
  enough from what was read.
- `internal/harness/repo.go`, `pathguard.go`, `inventory.go` — likely relevant to how a project's
  tracked files / `.agents/docs/` tree would be walked for the per-project index, not read.
- `.agents/docs/testing/2026-09-03-verify-sandbox.md` and `skill-matrix.md` in full — only
  `2026-09-16-linkshelf-simulation.md`'s header/setup/checklist tables were read. That file does
  show a usable pattern for the scope's "later" sandboxed token-measurement runs: detached lanes
  per CLI with a status file, `audit.py` (transcripts) for post-hoc analysis, and `--allowedTools`
  lists that already include `sqlite3` in at least one prior sandbox run — a precedent worth
  reusing rather than re-deriving.
- Exact regex/parsing details for BUG frontmatter fields beyond what one file showed, and the
  bugs/INDEX.md "Graduated" table's exact columns (not scrolled to).
- `internal/harness/derive.go` (Claude→Antigravity agent translation) and `internal/harness/plan.go`
  (what installs where) — relevant if a skill for index/search needs agent-visible tool grants
  (`toolMap` in derive.go, per project AGENTS.md "Rules for working here") — not read this pass.
- Whether any existing test asserts go.mod has zero `require` entries (none found via search, but
  search was pattern-based, not exhaustive).

## Open questions
- Does the operator want a REQUIREMENTS.md R4 amendment (ADR) approved BEFORE any implementation,
  or should recon instead explore a stdlib-only alternative (e.g. shelling out to a system
  `sqlite3` binary, which is itself a dependency of a different kind and likely no better under R4;
  or a much simpler non-SQLite grep/JSON index) as a competing option to present?
- Where would `--query`/`--limit`-type flags for `search` live given the single shared FlagSet
  parsed once for every verb — a new per-verb sub-FlagSet pattern, or positional args via
  `reorder()`? No precedent exists in this codebase today.
- How should the global `index.db`'s existence be reconciled with R7's existing term "the
  cross-project index" (currently = retros.jsonl-derived), to avoid two different things sharing
  one name in the trail and in context/AGENTS.md.
