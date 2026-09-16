---
owner: reviewer
date: 2026-09-16
---
# Cold read — `.agents/docs/plans/2026-09-16-attribution-write.md`

Range reviewed: plan file (`.agents/docs/plans/2026-09-16-attribution-write.md`), against `cmd/sh-agent/main.go` (install case, `attributionHint`, `attributionOff`), `cmd/sh-agent/main_test.go` (`TestAttributionNote`), `REQUIREMENTS.md`.

## (1) Gaps

- **`REQUIREMENTS.md` Non-goals contradicts the plan, and nothing addresses it — file:line `REQUIREMENTS.md:43`.** The Non-goals section reads: "Editing `settings.json` or any vendor setting: the harness is files in the customization directories, nothing else." The plan's D2/R4 sentence has `install` write a key into `~/.claude/settings.json` — the exact thing this bullet forbids. The plan's "Surfaces touched" list names only `REQUIREMENTS.md` R4, not the Non-goals section, and the operator's recorded approval (Why, D2) covers only "the one-sentence change to R4," not any edit to Non-goals. Per this repo's own hard rule, `REQUIREMENTS.md` is immutable without explicit operator approval for the *exact edit* — the plan has approval for one addition, not for resolving this direct contradiction. Implemented as scoped, the file ships self-contradictory. **BLOCKING.**
- **`ADR-010` is left standing and now false — `DECISIONS.md:81`.** ADR-010's Decision text says "`REQUIREMENTS.md` names editing vendor settings a non-goal" and its Consequences lean on "nothing installed outside the customization directories." The plan adds ADR-025 but does not mark ADR-010 superseded (the file's own convention, e.g. `## ADR-023: … — Superseded by ADR-024`). A reader of DECISIONS.md after this ships finds two ADRs asserting opposite things with no cross-reference. Not in Surfaces touched. **BLOCKING** (same root cause as the Non-goals gap — resolving one likely resolves both, but the plan decides neither).
- **A second, unnamed restatement of "never writes vendor settings" in README — `README.md:16`.** The paragraph the plan clearly intends to fix is line 14 (the attribution note, in the same section named in Surfaces touched). But line 16, in the same `## Install` section, separately says "the installer never touches vendor settings (ADR-010)" about `permissions.additionalDirectories`. This is a second, independent claim that becomes false and cites the same ADR-010 that gap 2 above leaves unresolved. Likely to be caught by an implementer re-reading the whole section, but it is not named, so it is a gap under review-standards §3 (grep the subject noun "vendor settings" across the repo) rather than a certainty. **NON-BLOCKING** (a competent implementer editing "README § Install" would likely see it), but worth calling out since it shares the ADR-010 root cause.
- **`sync` is not mentioned anywhere in the plan, but shares the exact code path.** `main.go:110` is `case "install", "sync":` — one switch arm. `attributionHint` already fires for both today; the plan's Interfaces section says `attributionSet` is "called from the install case for vendor all or claude" and Scope says "In: `install` (vendor all or claude)," never mentioning `sync`. Precedent (attributionHint firing for both) suggests the intent is for `sync` to write the key too, but the plan never says so, and the verification/E2E criteria only exercise `install`. **NON-BLOCKING** — the existing shared switch arm and `attributionHint` precedent give a reasonable implementer a default (treat `sync` the same as `install`), but it should be named so the test matrix covers it explicitly.
- **Detection of "no attribution key" vs. "not a JSON object" is under-specified.** The Interfaces section's insert branch requires "valid JSON object, no `attribution` key"; the fallback ("malformed... write nothing") doesn't explicitly say whether a syntactically-valid-JSON-but-non-object top level (e.g. `[]`, `"x"`, `42`) counts as "malformed." In practice `json.Unmarshal` into a `map[string]json.RawMessage` errors on those inputs, so the natural implementation folds this into "malformed" — but the plan doesn't state the decoding target, so this is asserted rather than checked. **NON-BLOCKING** (a reasonable implementer's natural choice matches the advisory-only spirit of D4).

## (2) Scope

Nothing found beyond what `## Why` and the Decisions ask for. The plan explicitly rejects the two scope expansions a reviewer would otherwise flag (forcing a value over a custom one; removing the key on `uninstall`), each with a D-row. The "new helper vs. extend `attributionOff`" choice (Interfaces, line 26) is explicitly left to the implementer and recorded as such — not silent scope creep.

## (3) Readiness

Questions a first-pass implementer would have to ask (or guess) before/while building:

1. **BLOCKING** — Does the Non-goals bullet in `REQUIREMENTS.md` get edited (and did the operator actually approve that edit), or does the plan's R4 addition ship alongside a bullet that flatly contradicts it?
2. **BLOCKING** — Does ADR-010 get a "Superseded by ADR-025" marker, consistent with how this file already handles superseded decisions (ADR-023)?
3. NON-BLOCKING — Does `sh-agent sync` also get the key written (same switch arm, same vendor gate), or only `sh-agent install`? (Reasonable default: same as `install`, following the `attributionHint` precedent already in the code.)
4. NON-BLOCKING — README's second "vendor settings" claim (line 16, ADR-010 citation) — reasonable implementer catches it while editing the section the plan already names.
5. NON-BLOCKING — New helper function vs. three-state `attributionOff` — plan explicitly defers this to the implementer.

## Verdict

**design needs revision** (two gaps: the Non-goals contradiction and the unresolved ADR-010).
**implementation not ready** (2 BLOCKING questions, both stemming from the same root cause: the plan never proposes the Non-goals edit or the ADR-010 supersession, so it cannot yet claim the operator approved them).

## Not examined

`context/AGENTS.md` full render/line-count claim ("render stays 93 lines") was not independently re-rendered; `deploy/local.md` and `deploy/remote.md` full files were only grepped for the attribution string, not read end to end for other stale restatements; `ADR-025`'s not-yet-written text obviously not reviewed since it doesn't exist yet.
