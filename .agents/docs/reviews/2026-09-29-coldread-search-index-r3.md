---
owner: reviewer
date: 2026-09-29
---
# Cold-read review (round 3) — sh-agent search plan   (2026-09-29)

Range reviewed: `.agents/docs/plans/2026-09-29-search-index.md` (design-document cold read, no code yet) plus the code/text surfaces it names: `cmd/sh-agent/main.go` (doc comment, `usage`, `reorder`/`takesValue`), `internal/harness/state.go` (`ConfigDir`), `internal/harness/report.go` (`Collect`'s explicit per-name allowlist), `internal/harness/size_test.go` (the 22,000-byte budget), `REQUIREMENTS.md` R1–R9 (not just R4 — read the whole file this round, since D4 explicitly reasons about a term R7 already uses), a real `BUG-*.md` file's frontmatter, `skills/dream/SKILL.md`, `skills/bug/SKILL.md`, `skills/feature/SKILL.md`, `agents/dreamer.md` (D13's cited consumer lines, spot-checked), `.agents/docs/specs/sh-agent/` and `.agents/docs/bugs/` directory listings. Also read both prior cold-read rounds (`reviews/2026-09-29-coldread-search-index-r1.md`, `-r2.md`) in full to confirm their findings are actually closed in the current text rather than take the D-rows' word for it.

## Status: COMPLETE

## Verdict: design needs revision (two new NON-BLOCKING gaps) — implementation ready (zero BLOCKING questions this round; both prior rounds' BLOCKING gaps stay genuinely closed)

Round 1's two BLOCKING gaps (FTS5 sync, query-time `files` join) and round 2's one BLOCKING gap (`--kind` enum missing `feedback`/`proposal`) are all fixed in the current text and I re-verified each independently: `chunks_fts`'s AFTER INSERT/DELETE triggers are present in the schema (Storage, lines 102-107); "Every query reads through `files`" is spelled out for `search`, `--refs` and `--open` alike (line 116); `--kind`'s line 58 enum now lists `feedback` and `proposal` alongside the other 13 values, matching Storage's kind assignment (line 118); `--global` now self-registers and reconciles the invoking project first (Global registry, line 151), closing round 2's non-blocking item too. I found no BLOCKING gap this round. I found two small NON-BLOCKING gaps neither prior round touched.

## (1) Gaps

- **NON-BLOCKING — R7 and the new feature both use the word "index" for different things, and the plan never touches R7's text.** `REQUIREMENTS.md` R7 (`REQUIREMENTS.md:29`, unchanged by this plan's Surfaces touched list) reads: "The harness's own state (manifest, feedback log, the harness proposal queue, **the cross-project index**) lives in `~/.config/sh-agent/` and nowhere else" — an existing reference to `retros.jsonl` (D4 says so explicitly: "R7 already calls `retros.jsonl` 'the cross-project index'"). The plan's own vocabulary for the new feature is "the global **index**" throughout Scope and Interfaces, its verb is literally `sh-agent index`, and its file is `~/.config/sh-agent/search.db` — landing in the exact directory R7's sentence enumerates, unmentioned by that enumeration. D4 used the R7 collision only to justify the *file name* (`search.db`, not `index.db`); it never revisits R7's prose itself. A reader of `REQUIREMENTS.md` after this ships, with no other context, cannot tell from R7 alone whether "the cross-project index" now means `retros.jsonl`, the new `search.db`, or both — and Surfaces touched names only R4 for `REQUIREMENTS.md` edits, so nothing in the plan closes this. A one-clause fix in Phase 1's R4/D1 task (naming `search.db` in R7's parenthetical, or renaming R7's phrase) settles it; a reasonable implementer would not silently invent a fix on their own since `REQUIREMENTS.md` is immutable without the operator's approval (project `AGENTS.md`) — this has to be an explicit edit, decided before Phase 1's ADR-045/R4/AGENTS.md commit task, not inferred.
- **NON-BLOCKING — `--refs`'s output shape for a directory target is unstated.** Interfaces § Commands says: "For a path, the first line is `<path> tracked|modified|untracked|absent` and then the references" (line 86). Separately, § Refs says a trailing `/` on `--refs X` means X is a directory, matching every ref under that prefix (line 147). A directory has no git blob, so it has no `tracked|modified|untracked|absent` state — the "first line" rule as literally written does not compose with the directory case. Two implementers could diverge: one omits the state line entirely for a directory target and prints only the reference list, one tries to compute a state for the directory itself (and either errors or defaults to something arbitrary). No unit test in Phase 1's `TestRefs` checklist item names a directory-target case, so nothing forces this decision to surface during implementation. A reasonable implementer most likely skips the state line for a directory (the natural reading, since directories were never given a state anywhere else in the spec), so this is non-blocking, but it is worth one clause before `TestRefs` is written.
- Everything else re-checked against real repo state matches the plan's claims exactly: a real `BUG-001` file's frontmatter (`id:`, `status:`, `date:`, `# ` title) matches Records § Bug precisely; the docs-root directory listing (`architecture bugs changes deploy handoffs plans retros reviews specs testing`) matches D17's directory-to-kind table exactly, with no directory left unmapped; `.agents/docs/specs/sh-agent/` already holds `delivery.md` and `memory-loop.md`, so a new `search.md` there (Phase 3) matches the existing layout convention with no surprise; `report.go`'s `Collect` remains a closed per-name allowlist with no wildcard scan of `~/.config/sh-agent/`, so the Phase 1 task to add an explicit `search.db` skip entry is necessary (not redundant) for the Failure-modes row to actually print it under "skipped"; D13's cited consumer lines (`skills/dream/SKILL.md:41`, `skills/bug/SKILL.md:33`, `skills/feature/SKILL.md:19`, `agents/dreamer.md:29`) each land on the paragraph D13 describes, though `skills/dream/SKILL.md:48` (cited alongside `:41` for "graduation check") is the "doubt resolves to `keep`" sentence rather than the graduation-verdict row itself (that's line 41) — a loose but defensible pairing, not a wrong citation.

## (2) Scope

Re-confirming rounds 1 and 2, independently: every Scope § In line and every Phase deliverable still traces to a D-row or a direct operator quote in `## Why`; the global index → D5, the recall/grep/token baselines and vector gate → D6/D8, `--open` and the metadata graph → Scope bullet 1's own words, the Phase 5 larger-project task → D11 (an implementer inference from the operator's words, attributed, not silent). No config knob for a value that plausibly never varies; no endpoint or output shape without a D-row. Scope is clean.

## (3) Readiness

1. NON-BLOCKING (new) — Does `REQUIREMENTS.md` R7's "the cross-project index" get a clarifying word once `search.db` also lives in `~/.config/sh-agent/`, or does it stay pointed only at `retros.jsonl`? (Gap 1.) State it in the R4/R7 edit task.
2. NON-BLOCKING (new) — For `--refs` on a directory target, is the `<path> tracked|...` first line printed at all? (Gap 2.) A reasonable implementer omits it; state the default.
3. Rounds 1 and 2's BLOCKING questions (FTS5 sync, query-time `files` join, `--kind` vs. global kind values) are closed by D15/D16/D18 and the current Interfaces text — re-verified independently against the schema and prose, not taken on the D-rows' word.
4. Round 1's and round 2's non-blocking questions (kind-mapping table, duplicate-oid path choice, DSN syntax, `--global` self-registration) are closed by D14/D17/D18 and the current text.

No Verification/E2E criterion bakes in an unstated implementation choice; no verify-clause fixture fails to produce the condition it claims (this round re-checked `TestFTSFollowsGC`'s `INSERT INTO chunks_fts(chunks_fts) VALUES('integrity-check')` call, which is the documented way to validate an external-content FTS5 table's shadow-table consistency and genuinely fails if the triggers are wrong — it exercises what it claims to).

## Criteria

Design-document cold read, six criteria applied loosely (no code exists yet):
1. Dead code — n/a.
2. Flow tracing — the chunks↔chunks_fts and query↔files seams (rounds 1-2) and the `--kind`/global-kind-value seam (round 2) stay closed under re-verification; this round's two gaps are lower-stakes cross-references (a REQUIREMENTS.md term, an output-format edge case) rather than a broken data-flow seam.
3. Standards conformance — R4's D1 replacement text is internally consistent; the new gap is R7, a sibling requirement the plan's Surfaces touched list never names, which is exactly the kind of unaddressed restatement the project's own rule on changed-rule restatements flags.
4. Directory layout — `internal/search` outside `internal/harness`, and `specs/sh-agent/search.md` beside the existing `delivery.md`/`memory-loop.md`, both match convention on inspection.
5. Simplicity — pure-Go single dependency, no vendoring, no service; no hand-rolled crypto/auth surface.
6. Decoupling — unchanged from round 2's assessment; no new coupling issue found.

## Not examined

- The scratch research/audit briefs (`plans/2026-09-29-recon-project-index.md`, `plans/2026-09-29-audit-search-consumers.md`, `plans/2026-09-29-research-sqlite-index.md`) were not opened this round either.
- Most of D13's individual consumer-line citations beyond the four spot-checked here (`skills/review-standards/SKILL.md:48`, `skills/handoff/SKILL.md:154`, `skills/retro/SKILL.md`, `agents/bugs.md`, `agents/changelog.md`, `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `skills/research/SKILL.md`) were not opened this round; rounds 1-2 spot-checked a different, disjoint pair.
- `modernc.org/sqlite`'s actual DSN pragma syntax, FTS5 support and cross-compile behaviour were not built or run (no module cache/network access from this review), same limitation both prior rounds stated.
- Phase 4/5 mechanics (`.agents/scratch/eff-kit/`) — nothing exists on disk yet to check.
- `internal/harness/size_test.go`'s exact current byte count was not re-measured this round (round 1 measured 21,952/22,000; not re-run here).

design needs revision, implementation ready
