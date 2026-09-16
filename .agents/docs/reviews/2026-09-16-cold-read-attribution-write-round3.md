---
owner: reviewer
date: 2026-09-16
---
# Cold read round 3 (cap) — `.agents/docs/plans/2026-09-16-attribution-write.md`

Range reviewed: plan file

Prior findings checked: round 1's BLOCKING 1–2 (Non-goals contradiction, ADR-010 left standing) are closed by D5 and Phase 1 task 4 (Non-goals exception sentence, ADR-010 heading marked "superseded in part"). Round 1 NON-BLOCKING 3 (`sync`) closed by D6. Round 1 NON-BLOCKING 4 (README:16) closed by Phase 1 task 5 (reworded sentence named verbatim). Round 2's BLOCKING (the `TestAttributionNote` regression) is closed by Phase 1 task 2, which now decides the new expectations for both cases per command leg and adds the `status`-alone case (D9). Round 2 NON-BLOCKING carried (non-object JSON) is closed by D7 and the new Interfaces bullet (plan line 26). Round 2 NON-BLOCKING (error contract) is closed by D8 (no error return).

## (1) Gaps

- **NEW, NON-BLOCKING — a third, unnamed restatement of "the installer never edits/touches vendor settings (ADR-010)" survives at `.agents/docs/deploy/local.md:37`, in a section Surfaces does not name.** That file's `## First session after install` (line 37, not `## Verify`, lines 29–34) reads: "add the directory to `permissions.additionalDirectories` in `~/.claude/settings.json` yourself — the installer never edits vendor settings (ADR-010)." Surfaces touched names only "`.agents/docs/deploy/local.md` § Verify," not § First session after install. This is the same subject-noun restatement round 1 caught at README:16 (closed) and the plan already fixes at README:16 and cites ADR-010's "superseded in part" heading — but this third instance, in a different section of a different file, is not in any Phase 1 task. After this ships, ADR-010 is superseded in part and the sentence is stale for the general claim (the installer does now edit one vendor-settings key) even though it remains true for the specific `permissions.additionalDirectories` case it's talking about — confusing, not decided. `remote.md:19` ("sh-agent never writes it") is the parallel instance and is correctly named and covered by Surfaces ("remote.md § First install"); `local.md:37` is the one left out.

## (2) Scope

Nothing found beyond `## Why` plus the operator-approved additions (D2, D5) and the cold-read-sourced D6–D9. No unnamed scope expansion.

## (3) Readiness

1. NON-BLOCKING — does `deploy/local.md:37` ("the installer never edits vendor settings (ADR-010)") get reworded alongside README:16 and local.md's § Verify paragraph, or is it left as a now-imprecise restatement? Surfaces names the wrong section of this file; a first-pass implementer editing "§ Verify" as instructed would plausibly miss it.
2. NON-BLOCKING (carried, deferred by the plan itself) — new helper vs. three-state `attributionOff` (Interfaces line 27) — implementer's explicit choice, recorded as such.
3. NON-BLOCKING — `TestAttributionSet` (Phase 1 task 1) is specified only by behavior ("missing file created," "malformed and FIFO untouched," etc.), not by whether it drives `attributionSet` directly (unit-level, matching the Interfaces section's exact signature) or through `run(["install", ...])` (black-box, matching `TestAttributionNote`'s style). Either produces the same coverage; a reasonable implementer picks one and it doesn't change behavior, so this doesn't block.

No BLOCKING questions remain.

## Verdict

**design needs revision** (one new gap: the third, unnamed `local.md:37` restatement, same root cause as the two already fixed, not yet decided — low stakes, but the plan's own Surfaces list for that file points at the wrong section).
**implementation ready** (zero BLOCKING questions — all three rounds' BLOCKING items are now resolved; the remaining gap is docs-accuracy, non-blocking, and does not stop a first-pass implementer from building and shipping the write path, tests, R4, and ADR-025 correctly).

## Not examined

`ADR-025`'s text (not yet written, nothing to check). `context/AGENTS.md`'s claimed 93-line render after the edit — not independently re-rendered. `deploy/remote.md` and `deploy/local.md` read in full this round (confirmed no further unnamed restatements beyond the one cited above); `README.md` read in full for "vendor setting"/"attribution" hits — no others found beyond the two the plan already names. Did not re-verify `go test` compiles/passes against the plan's exact task ordering (Phase 1 task 1 writing `TestAttributionSet` against an as-yet-unimplemented `attributionSet` would fail to compile rather than fail red on assertions if the test calls the helper directly — cosmetic TDD-sequencing point, not filed as a gap since a stub function is the obvious implementer move and the plan doesn't preclude it).
