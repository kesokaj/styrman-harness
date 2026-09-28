---
owner: reviewer
date: 2026-09-28
---
# Review — update-from-zip-docs   (2026-09-28)

Range reviewed: 1cb2608..worktree

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The delivery-model rewrite (REQUIREMENTS R8/R9, ADR-038, `AGENTS.md` rule 1, `delivery.md` § The model, README's map and prose, `local.md`, `remote.md`) is internally consistent and each restatement of "two ways in" / "no third" was updated in the same diff — I could not find a surviving unqualified copy of the old wording. The CLI `usage` string and README's `## Install from a zip` section agree on flags and behaviour. The one real gap is that ADR-034's superseded claim ("a release binary installs the harness built into it, wherever it is run and whatever trees are on the machine") survives unqualified in one place README does not route through `delivery.md`, and the map draws the new zip edge solid despite the E2E verification recording it as run on only one machine.

## Findings

- **[medium / confidence: high]** `README.md:185` — "A RELEASE binary uses none of that: it installs the harness built into it, wherever it is run and whatever trees are on the machine, and refuses `--repo` rather than appearing to accept one (ADR-034)." This is the pre-ADR-038 rule verbatim. `.agents/docs/architecture/DECISIONS.md:340` (ADR-038 **Supersedes**) states plainly: "A release binary installs a tree on disk only when the operator names a zip with `--from-zip`; a recorded or discovered path is still refused." `delivery.md:32` carries the same carve-out. README's line was not updated to match: read alone it tells the operator a release binary can never install a tree on disk, which `--from-zip` on that same binary now does. `--repo` is still correctly refused (that half is unaffected), but the "whatever trees are on the machine" clause is stale. narrowed: `.agents/docs/architecture/DECISIONS.md:340` and `delivery.md:32` (both updated) → `README.md:185` (not updated).

- **[medium / confidence: medium]** `README.md` mermaid map, `zipPlace` node and the `harnessRepo -- "zip, downloaded by hand..." --> zipPlace` edge (README.md ~line 39) — drawn solid teal (in the `linkStyle 0,1,5,6,7` group), the same weight as the checkout and release-binary paths. `AGENTS.md` rule 4: "A dashed line is done by hand or limited to one machine; it turns solid in the commit that closes the gap, and nothing is drawn solid before it works for every place." `.agents/docs/testing/2026-09-28-verify-update-from-zip-e2e.md` § Not exercised lists "macOS" and "Antigravity reading the installed files" as unexercised, and the golden-path run was on the one Linux box only; `remote.md`'s four places include two Macs. If "works for every place" means verified there, the new edge does not yet qualify and per rule 4 should be dashed until a Mac run closes the gap (or the report states it was verified on one and generalizes safely, which it does not — it explicitly disclaims macOS).

- **[low / confidence: medium]** `README.md` § Daily commands, the `sh-agent update` table row (README.md:177) — describes checkout and release-binary behaviour and `--check`, but does not mention `--from-zip`, even though `main.go`'s usage string now documents it as a first-class mode of `update` and the table is the CLI's canonical command-by-command reference. The behaviour is documented elsewhere (§ Install from a zip), so this is a completeness gap in one table rather than a wrong statement.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- `.agents/docs/plans/2026-09-28-update-from-zip.md`, `2026-09-28-recon-from-zip.md`, and the three `2026-09-28-coldread-*` review files, `BUG-093`'s own record body, and `2026-09-28-verify-update-from-zip-phase-1.md` — plans/reviews/bugs are excluded from the grep instruction and were only skimmed for orientation, not checked line-by-line against the code.
- `internal/harness/zip_test.go` and `cmd/sh-agent/zip_test.go` bodies (existence and general shape confirmed via the phase-2 verification log; not read in full).
- `.agents/docs/testing/guide.md` beyond a grep for "zip"/"delivery"/"two ways" — the file has no hits, so no zip-specific update procedure text was found there to check for staleness; whether the guide *should* gain a zip-place row was not evaluated against the full file.
- Exhaustive edge-index arithmetic in the mermaid `linkStyle` groups (i.e., whether every edge index still points at the intended edge after insertion) was not verified beyond a read-through.
- macOS/Antigravity behaviour itself — only the verification doc's own "not exercised" claim was read, not re-tested.
