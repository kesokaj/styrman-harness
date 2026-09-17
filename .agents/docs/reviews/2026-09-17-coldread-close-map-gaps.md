---
owner: reviewer
date: 2026-09-17
---
# Cold-reader check — close-map-gaps plan (2026-09-17)

Range reviewed: `.agents/docs/plans/2026-09-17-close-map-gaps.md` (plan document, not a diff — per dispatch, the cold-reader check)

## Round 1

Read: the plan in full; `REQUIREMENTS.md` R7–R9; `internal/harness/report.go` (`hostDir`, `runDir`, `Collect`, redaction); `skills/retro/SKILL.md` in full; `README.md` § How it works mermaid block and legend; `cmd/sh-agent/main.go` (confirms `os.Hostname()` feeds `hostDir`); a real report tree at `~/.config/sh-agent/feedback/reports/` (only host present: `dave`, three runs today).

### (1) Gaps

- **G1 — BLOCKING interface gap: no way to point Step 6's ledger scan at a fixture.** Step 6 (Interfaces, line 46) and the Verification E2E both require a "harness retro record fixture" that has erin's proposal B recorded `Source: … applied …`, so the golden-path run can prove B is skipped. But Step 6's scan target is fixed: `<docs-root>/retros/*.md` of *this* harness checkout — i.e. the real `.agents/docs/retros/`. The plan gives `--feedback-repo DIR` to redirect the feedback checkout (D7) but no equivalent override for the ledger path, and the Phase 1 E2E task (line 96) only says the fixture is built "exactly as Verification states" without naming its location or how the skill run is made to read it instead of the real retros directory. Two careful implementers diverge here: one drops synthetic `erin`/`frank`/`dave` records into the real `.agents/docs/retros/` (temporarily or permanently — this repo's own audit trail, meant to be immutable and append-only), the other fakes the check some other way that no longer exercises the actual skill text. Needs a decision before Phase 1's E2E task starts: either add a `--docs-root DIR` (or similar) override to Step 6, or specify exactly where the fixture ledger record lives and how the run is pointed at it without touching real `.agents/docs/retros/`.

- **G2 — Traceability gap: D-row scope-line citations don't match the visible Scope section.** D7 and D8 both cite "scope line 8," and the Rejected block cites "scope line 7," but the Scope section's visible enumeration is 5 `In:` bullets + one `Out:` line + one `Follow-ups:` line — 7 groups, and no visible "line 8" exists. D8's own subject (the map) is the *5th* `In:` bullet by visible position, not the 8th. D1–D3 are all attributed to "scope line 1" although they cover three different sub-clauses of that one bullet (report consumption, findings, promotion-check-only), while D6 (skip own reports) is also a sub-clause of bullet 1 but is cited as "scope line 6." This reads as if the operator's original "go all" was against an 8-item enumeration that isn't reproduced in this file — a cold reader cannot verify each Decision's claimed operator approval against what's on the page. Notably D7 (`--feedback-repo`) is the one D-row NOT attributed to "the operator" (it says "by: implementer") even though `--feedback-repo` is visibly one of the operator-confirmed `In:` bullets — worth confirming that flag really was operator-approved and not an implementer add-on riding on the map/docs commit.

- **G3 — minor, low-severity: byte-order "newest run" tie-break has an unstated edge case.** `runDir` (report.go:565-584) suffixes collisions un-padded: `-2`, `-3`, … `-10`. Step 3's "greatest name in byte order" (line 43) is correct for single-digit suffixes and cross-day/cross-second comparisons, but at 10+ reports from one host in the same wall-clock second, `"…-9"` sorts *after* `"…-10"` in byte order, so the rule would silently pick a stale run. Practically unreachable (needs ≥10 `sh-agent report` runs inside one second on one host) — flag only because the plan's Risks section names one "checked only by the E2E fixture" assumption about byte-order matching but doesn't mention this failure shape; not worth blocking on.

- **G4 — no checkbox/verification for the interactive map artifact.** Surfaces touched lists "the interactive map page (private claude.ai artifact, outside the repository) — updated … after the commit," but no Phase task, verification step, or owner appears in Phases for it. Low severity since it's explicitly out-of-repo and covered by a standing global rule, but it is a named surface with no corresponding checkbox, unlike everything else in Surfaces touched.

None of G1–G4 are misreadings of `report.go`'s mechanics: `hostDir`'s `/`→`_`, `\`→`_` replacement (report.go:588-594), the redaction marker prefix `[redacted by sh-agent report` (report.go:26, 161), and `Collect`'s inclusion of `feedback.jsonl`, `retros.jsonl`, `state.json`, `proposals/*.md` (report.go:260-290) all match what the plan asserts about them.

### (2) Scope

Every "In" item traces to a D-row citing "the operator" or (for D7) "implementer, to make [operator-approved] scope line 8's commit verifiable" — no unattributed extra endpoint, shape, or config-for-values-that-never-vary found. `--feedback-repo` mirrors an existing flag (`sh-agent report`'s own `--feedback-repo`), so it isn't a novel surface. The one open question is G2 above: whether the *citation* for `--feedback-repo`'s approval is trustworthy, not whether the feature itself is in scope — the bullet is visibly in the operator-confirmed `In:` list either way.

### (3) Readiness

Questions a first-pass implementer would have to ask:

1. **BLOCKING** — Where does the "harness retro record fixture" (with erin's `Source:`/`applied` line) live, and how is the retro run pointed at it instead of the real `.agents/docs/retros/`, given Step 6 names no override? (G1)
2. **NON-BLOCKING** — Do the D7/D8 "scope line 8" citations need reconciling with the visible 5-bullet `In:` list before Phase 1 starts? A reasonable implementer can proceed straight from the Interfaces section regardless (the flag and behavior are unambiguous there) and treat this as a documentation/traceability cleanup, not a behavior question. (G2)
3. **NON-BLOCKING** — Byte-order tie-break at ≥10 same-second runs: the plan's stated rule ("greatest name in byte order") is a clear enough default that an implementer just follows it as written. (G3)
4. **NON-BLOCKING** — The interactive map artifact's update: covered by the standing global rule (AGENTS.md rule 5) even without its own Phase checkbox. (G4)

## Criteria

Not a diff review (cold-reader check of a plan document); the six-criteria review-standards lens does not apply file-by-file here. Flow tracing was done against `report.go` and `skills/retro/SKILL.md` as named by the dispatch; directory-layout, dead-code and decoupling criteria are not meaningfully assessable against a plan with zero code changes (D1: skill-only).

## Not examined

- `.agents/docs/deploy/local.md` § Report back and `.agents/docs/specs/sh-agent/delivery.md` (named in Surfaces touched) — not read; cannot confirm today's text actually says "read by hand" or what exactly Phase 1's docs task must change there.
- `context/AGENTS.md:41` and `skills/feedback/SKILL.md:88` — not read directly; took the plan's line citations on faith.
- `.agents/docs/architecture/DECISIONS.md` (for ADR-032's eventual shape/numbering conflicts) — not read.
- Full `README.md` mermaid block beyond the lines grepped (node/edge definitions 1–14, full edge list, exact current edge count) — only the retro/report-relevant lines were pulled; did not independently recount "16 edges → 15."
- `cmd/sh-agent/main_test.go` `TestHarnessMap` — not read; took the plan's claim that it only needs `install`/`update`/`report` to keep passing at face value.
- No git diff exists for this dispatch (plan-only, no code yet) — nothing to review under the six-criteria diff lens beyond what's above.

design needs revision
implementation not ready

## Round 2

Read (fresh, before reading Round 1): the plan in full; `internal/harness/report.go` (`Collect`, `BuildRecord`, `SendReport`, `runDir`, `hostDir`); `skills/retro/SKILL.md` in full (Steps 1, 4, 5, 6); `README.md` § How it works (mermaid + prose); `REQUIREMENTS.md` R7–R9; project `AGENTS.md` § How the harness works and § Testing; `context/AGENTS.md` § Learning loop; `skills/feedback/SKILL.md` (Where the loop's outputs land); `.agents/docs/deploy/local.md` § Report back; `.agents/docs/specs/sh-agent/delivery.md` § The loop back / § Where things stand; a real `~/.config/sh-agent/feedback/reports/dave/` tree (three runs, one `proposals/` dir, six proposals — all already `Status: applied`, none queued).

### (1) Gaps

- **NEW, BLOCKING — Step 6's "apply" instructions don't distinguish a report-sourced proposal from a local-queue one, and the natural reading writes back into the feedback checkout.** Current `skills/retro/SKILL.md` Step 6, "Harness scope" bullet: "...apply there under that repo's own rules (...), then, after that repo's review and verification, **update the status in the queue file** and in the record...". The plan's Surfaces touched (line 26) lists this bullet as needing an edit for "report findings and the ledger," but Interfaces gives no replacement text for it — unlike every other Step it touches (1, 4, 5 all get exact wording). For a proposal consumed from another host's report, "the queue file" is a copy of that host's `proposals/*.md` sitting inside the local clone of the feedback checkout (`~/.config/sh-agent/feedback/reports/<host>/<run>/proposals/*.md`). The plan's own Rejected list says "marking proposals applied in the feedback repository — a second write path back, against R9's one way back," and D5 says nothing is written to the feedback repository — but that rejection lives only in the plan, not in the skill text an implementer follows during Phase 1. Two careful implementers diverge: one recognizes the R8/R9 one-way-back rule and writes status only into the local ledger record's `Source:` line (per Interfaces' "Ledger line" spec); the other follows Step 6's existing pattern literally and writes/commits status into the feedback checkout's copied file (which the current text does for every other harness-scope proposal). This needs its exact replacement text before Phase 1's first task, the same way Steps 1, 4 and 5 got theirs.
- **Self-inconsistency, low/medium — Risks section claims a check that no task performs.** Risks: "`hostname` prints the same name `os.Hostname()` returns... checked by comparing `hostname` with the directory `sh-agent report` created on this box (`reports/dave/`), **first task of Phase 1**." Phase 1's actual first task (line 95) is only the `skills/retro/SKILL.md` text edit, verified by "read against Interfaces" and an install/status check — it does not mention running `hostname` or comparing it against `reports/dave/`. Confirmed independently: this box's real `~/.config/sh-agent/feedback/reports/` has exactly one host directory, `dave`, matching `hostname` here — so the assumption happens to hold, but nothing in the Phase 1 checklist as written verifies it. Either add the comparison as an explicit verify step on task 1, or drop the "first task of Phase 1" claim from Risks.
- **Confirmed resolved from Round 1:** G1 (ledger-scan fixture location) — D9 now runs the whole E2E from `harness-copy/` with its own `.agents/docs/retros/`, so `<docs-root>` resolves inside the copy and the real trail is never touched; this is a sound fix given the retro skill's own docs-root rule is cwd-relative. G2 (scope-line citation mismatch) — Decisions now opens with the operator's own 8-item enumeration ("1 mechanism, 2 which reports, ... 8 map and R9"), and each D-row's citation now matches that list. G3 (byte-order tie-break) — D10 explicitly accepts the >9-same-second-runs edge case rather than silently ignoring it. G4 (interactive map artifact has no Phase checkbox) — still true; still non-blocking per Round 1's own reasoning (covered by the standing global AGENTS.md rule).

### (2) Scope

Nothing found beyond what `## Why` and the eight confirmed scope lines ask for. `--feedback-repo DIR` (D7) mirrors an existing `sh-agent report` flag and is attributed, even though to "implementer" rather than "the operator" directly — Round 1 flagged this same attribution oddity (G2) as non-blocking, and I agree: the flag itself is inside the operator's confirmed scope-line-8 group ("map and R9" / interfaces), and D7's stated reason (making the E2E verifiable without touching the real feedback repo) is sound engineering, not an unrequested feature.

### (3) Readiness

Questions a first-pass implementer would have to ask:

1. **BLOCKING (new)** — What is the exact replacement text for Step 6's "Harness scope" bullet, so that a report-sourced proposal's status is recorded only via the local ledger's `Source:` line and never written into the feedback checkout's copy? Two implementers build different things here: one preserves R8/R9's one-way-back rule, the other silently violates it by reusing the existing "update the status in the queue file" instruction verbatim.
2. **NON-BLOCKING** — Should the Risks section's "first task of Phase 1" claim about checking `hostname` against `reports/dave/` be turned into an actual verify line on task 1, or dropped? Either way an implementer can proceed (the assumption already holds on this box); it's a documentation-accuracy gap, not a design question.
3. Round 1's G4 (no Phase checkbox for the interactive claude.ai map artifact) — still open, still non-blocking, as Round 1 found.

design needs revision
implementation not ready

## Round 3

(in progress)

Read fresh (before reading Rounds 1–2): the plan in full; `skills/retro/SKILL.md` in full; `internal/harness/report.go` (`SendReport`, `runDir`, `hostDir`, redaction marker); `README.md` § How it works (mermaid + prose, node/edge lines); `REQUIREMENTS.md` R7–R9; a real `~/.config/sh-agent/feedback/reports/dave/` tree (three runs, `proposals/*.md` — six proposals, all `**Status:** applied ...`, none queued) confirming this box's real hostname is `dave`, matching `reports/dave/` — the plan's own Risks example; `cmd/sh-agent/main_test.go` `TestHarnessMap`.

### (1) Gaps

- **NEW, BLOCKING — the E2E golden-path command's `--feedback-repo` path does not resolve from the stated working directory.** Verification says the E2E's working directory is `harness-copy/` (a `git archive`-derived copy under `.agents/scratch/`, so it has no nested `.agents/scratch/` of its own — that directory is gitignored and untracked, so `git archive HEAD` cannot carry it). The golden-path command is then given as `/retro --harness --dry-run --feedback-repo .agents/scratch/fb-fixture`, a path relative to the repo root, not to `harness-copy/`. Run literally from `harness-copy/`, that relative path does not exist, so Step 1's own failure text fires: `retro: feedback checkout <path> not read: absent`. That is exactly the wording and behaviour the plan's separate **Edge** case is designed to test ("naming a directory that does not exist prints `not read:` and still lists the local queue") — so the golden path, taken literally, degenerates into the edge case: it would report success (`not read:`, local queue only) without ever having parsed `fb-origin`/`fb-fixture`'s `erin`/`frank`/`dave` data, and nobody would notice, because both paths print consistent-looking output. This is the "silence from your own instrumentation is not evidence" failure shape from the global rules — a green E2E run here would not prove the report-consumption code path was exercised at all. Two careful implementers diverge: one passes an absolute path (or `../fb-fixture` / `$PWD`-anchored) so the golden path actually reads the fixture; another follows the literal relative path as written, gets the "absent" branch, and never notices the fixture went unread. Needs an exact path (absolute, or explicitly "resolved from the repo root regardless of cwd") before the E2E task starts.
- **NON-BLOCKING — ADR-032 has no exact text, unlike every other Interfaces item.** The plan gives verbatim replacement text for the SKILL.md steps, the map edges, and even the Step 6 bullets, but for ADR-032 only "ADR-032 appended" (Scope, Surfaces touched) with no drafted content. Given ADR-030/031's precedent (short, R-number-anchored, quoting the approved text) and the plan's own rich Decisions section to draw from, a competent implementer can draft this alone — flagging only because it is the one Surfaces-touched item with zero specified content.

### (2) Scope

Unchanged from Rounds 1–2: every `In:` item still traces to a D-row citing the operator (or, for D7/D9/D11, the implementer/cold-read against an operator-approved scope line), and the scope-line citations now visibly match the operator's own 8-item enumeration opening the Decisions section. No unattributed extra endpoint, shape, or never-varying config found. `--feedback-repo` continues to mirror the existing `sh-agent report` flag, so it is not a new surface.

### (3) Readiness

1. **BLOCKING (new)** — Is the golden-path `--feedback-repo` value an absolute path (or otherwise anchored to the repo root), given the stated cwd is `harness-copy/` and the value as literally written does not resolve there? Without this, two implementers get different E2E outcomes, and one of them is a false pass.
2. **NON-BLOCKING** — Should ADR-032's text be drafted in the plan itself, or left to Phase 1? A reasonable implementer drafts it from the Decisions section and prior ADRs' style.
3. Round 1's G4 (no Phase checkbox for the interactive claude.ai map artifact) — **resolved**: Phase 1 now has an explicit task ("After the push: the interactive map page ... updated ... per AGENTS.md rule 5").

### Round 1 and 2 items, checked against the current plan text

- G1 (ledger-scan fixture location, blocking) — **resolved**: D9 runs the whole E2E from `harness-copy/` with its own `.agents/docs/retros/`, keeping the real trail untouched.
- G2 (D-row scope-line citations vs. visible Scope enumeration) — **resolved**: Decisions now opens with the operator's own 8-item list and every D-row's citation matches it.
- G3 (byte-order tie-break at ≥10 same-second runs) — **resolved**: D10 explicitly accepts the edge case rather than leaving it silent.
- G4 (no Phase checkbox for the interactive map artifact) — **resolved**: see Readiness item 3 above.
- Round 2 new-1 (Step 6 exact replacement text for report-sourced status, one-way-back) — **resolved**: Interfaces § "Step 6, where statuses are written" gives verbatim replacement text for both bullets, and D11 records the decision.
- Round 2 new-2 (Risks' "first task of Phase 1" hostname claim didn't match the checklist) — **resolved**: Phase 1's first task is now literally the hostname-vs-`reports/<host>/` check, gating the rest of the phase on a mismatch.

## Criteria

Not a diff review (cold-reader check of a plan document); dead-code, directory-layout and decoupling do not meaningfully apply to a skill-text-only plan (D1). Flow tracing was done against `report.go` and `skills/retro/SKILL.md` as named by the dispatch, plus the real report tree and `TestHarnessMap`.

## Not examined

- `.agents/docs/deploy/local.md` § Report back, `.agents/docs/specs/sh-agent/delivery.md`, `context/AGENTS.md:41`, `skills/feedback/SKILL.md:88` — not re-read this round; Round 2 already read these and found no conflict, taken on that basis.
- `.agents/docs/architecture/DECISIONS.md` (ADR-030/031 exact wording, for ADR-032 numbering/style precedent) — not opened this round; ADR-032 finding above is inferred from the plan's own citations of ADR-030/031.
- Full current `README.md` mermaid edge list beyond the retro/report-relevant lines (did not recount "16 edges → 15" independently).
- No actual dry run of `/retro --harness` was performed (not possible outside a live session); the path-resolution finding above is a textual/logical read of the Verification section, not an observed failure.

design needs revision
implementation not ready
