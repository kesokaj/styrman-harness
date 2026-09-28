---
owner: reviewer
date: 2026-09-28
---
# Cold-read review, round 3 (final) — `.agents/docs/plans/2026-09-28-update-from-zip.md`   (2026-09-28)

Range reviewed: plan file `.agents/docs/plans/2026-09-28-update-from-zip.md` and the code surfaces it names
(`internal/harness/{source.go,repo.go,apply.go,state.go}`, `cmd/sh-agent/main.go`, `embed.go`, `REQUIREMENTS.md`
R8/R9, `.agents/docs/architecture/DECISIONS.md` ADR-030/ADR-031/ADR-037). Rounds 1 and 2's reports read only after
forming my own view, per the dispatch.

## Status: COMPLETE

## Verdict: design needs revision — one new blocking gap; rounds 1 and 2's blockers are fixed

Rounds 1 and 2's high-severity findings are resolved in the current plan text: D6 scopes the symlink/mode/size
checks to the five extracted trees while keeping name-safety archive-wide (round 1's `.agents/skills` contradiction);
D7 explicitly decides that a zip-installed dev build with no checkout gets an improved *error* naming the zip
install, not a working `ResolveSource` success path (round 2's gap); D8 settles the plain-copy-vs-zip-place docs
relationship (round 1 finding 4 / round 2 finding 3); the Phase-3 grep in task 2 now includes `only way the
harness|\bways in\b|one of the ways` and a positive control pinned to the exact two lines (`README.md:182`,
`deploy/remote.md:61`) round 2 found the narrower pattern missing. All four of rounds 1-2's open items are closed.

But the plan still has an unresolved, self-contradicting edit to an immutable, operator-locked requirement: it adds
a fifth bulleted item to R9 without deciding what happens to R9's title and its own "no others" clause, both of
which the operator personally dictated as "vital," "set in stone," "rock solid" (ADR-031).

## (1) Gaps

- **[BLOCKING]** R9 is titled "Four procedures, set in stone" and ADR-031 (D1) records the operator's own words
  fixing "the four procedures... Install, Update, Report back, Release — and no others," calling this "vital,"
  wanting it "set it stone... rock solid." The plan's Interfaces section directs: "R9 **Install/Update** gains:
  '**Zip place:** the operator downloads the repository zip by hand and runs `sh-agent update --from-zip <file>`....'"
  — quoted with the same bold-label-plus-colon format as the four existing bullets (`**Install:**`, `**Update:**`,
  `**Report back:**`, `**Release:**`). Read as written, this inserts a fifth bulleted procedure into a section whose
  own title says four and whose own text says "and no others." The plan:
  - never states whether R9's title becomes "Five procedures" (and if so, on what approval — R9's title is not
    among the "REQUIREMENTS text quoted under Interfaces" the Why section says was approved verbatim, only the R8
    first sentence and the R9 bullet text are; the title and the "no others" clause are not shown as edited anywhere
    in the plan);
  - never states whether ADR-031 (D1: "fixes the four procedures... and no others") is amended or superseded
    alongside ADR-038 — ADR-038's own description ("supersedes ADR-034 for this one path only... amends the 'two
    ways, no third' claim in `delivery.md`") names ADR-034 and `delivery.md` but not ADR-031 or R9 itself;
  - never states an alternative reading where the zip material is folded as prose into the existing `Install:`
    and/or `Update:` bullets instead of becoming a labeled fifth bullet, which would keep R9's count at four and
    avoid the contradiction — the literal quoted text (a standalone bold-labeled paragraph) argues against this
    reading, but the plan doesn't rule it out either.

  `TestHarnessMap` (ADR-031 D4) checks that install/update/report are named on the map and pins verbs, not R9's
  bullet count or its title text — nothing in the test suite would catch a shipped REQUIREMENTS.md whose heading
  says "Four... and no others" sitting above five bullets. Two careful implementers diverge here: one pastes the
  approved bullet in and ships R9 self-contradicting itself (literal execution of the approved text, title left
  untouched since not itself approved); the other treats "set in stone... and no others" as still binding, retitles
  the section and treats that retitling as needing its own explicit operator approval before Phase 3 — which the
  plan's Why section, scoped to "the REQUIREMENTS text quoted under Interfaces," does not currently grant. This must
  be decided — and, per this project's own rule, decided by the operator explicitly, before Phase 3 writes to
  `REQUIREMENTS.md` — not inferred by whoever implements the phase.

- **[non-blocking, low]** `internal/harness/state.go:77`'s `ConfigDir(home)` returns `(string, error)`, not a bare
  string. The plan's "Order of work" step 2 pseudocode, `harness.UnpackZip(file, harness.ConfigDir(home))`, reads as
  passing the two-value call directly — not valid Go against `UnpackZip(zipPath, parent string)`. Trivial for an
  implementer to fix (call it, handle the error) and not a design question; noted so nobody copies it literally.

## (2) Scope

Nothing in `## Scope: In` is unattributed: every item traces to the operator's "approve as written, go all but
number 6," and item 7 (the matrix) is explicitly D5, read from the same turn. No second shape for `update`, no
unrequested endpoint, and no configuration surface for a value that never varies was found — `--from-zip`,
`--home`, `--vendor`, `--dry-run` are all flags every command already accepts, not new surface invented for this
feature. Nothing to flag.

## (3) Readiness

1. **BLOCKING** — Does R9's title stay "Four procedures, set in stone," change to five, or does the zip material
   fold into the existing `Install:`/`Update:` bullets instead of becoming a fifth labeled one — and, either way,
   does ADR-031's D1 ("fixes the four procedures... and no others") get amended alongside ADR-038? (Gap above.
   Two implementers ship different REQUIREMENTS.md text, and one of the two ships a document that visibly
   contradicts its own heading and the operator's own "set in stone... and no others" instruction — on a file this
   project treats as immutable without the operator's explicit approval.)

Everything else checked resolves without a question: D6-D8 close rounds 1-2's blocking and non-blocking items;
`apply.go:117-123`'s `Origin`-before-`Dir` precedence, `status()`'s `LoadState` reordering, and `repo.go:75`'s
`ZipPrefix` skip are all now spelled out explicitly enough for a first-pass implementer to follow without guessing;
`reorder()`'s `takesValue` table needing `--from-zip` is named as a task; the five embedded trees
(`context, skills, agents, commands, rules`) match `embed.go:9`'s `//go:embed` directive exactly — the `workflows`
in `source.go:19`'s doc comment is a stale, pre-existing reference to Antigravity's vendor-side deploy layout, not
a source tree, and outside this plan's scope; ADR-038 is the next free ADR number.

## Not examined

- `.agents/docs/testing/skill-matrix.md` and `reviews/2026-09-28-review-v1.7.0-381e86d-code.md` findings 1-3
  (Phase 4) — not read in detail; downstream of the Phase 1-3 design question this round focused on, and rounds 1-2
  didn't flag it either.
- `README.md` § How it works (the map) itself, and its dashed/solid line convention — not opened in any of the
  three rounds.
- `.agents/docs/specs/sh-agent/delivery.md`'s full text (only grep hits read across rounds).
- `cmd/sh-agent/main_test.go` existing test conventions, `zip_test.go` (does not exist yet — new file).
- ADR-034's full text (only the portion quoted in `repo.go`'s comments was read, across all three rounds).

design needs revision
