---
owner: reviewer
date: 2026-09-28
---
# Cold-read review — `.agents/docs/plans/2026-09-28-update-from-zip.md`   (2026-09-28)

Range reviewed: plan file `.agents/docs/plans/2026-09-28-update-from-zip.md` and the code surfaces it names
(`internal/harness/source.go`, `repo.go`, `apply.go`, `cmd/sh-agent/main.go`, `embed.go`, `REQUIREMENTS.md` R8/R9,
`.agents/docs/specs/sh-agent/delivery.md`, `AGENTS.md`, `.agents/docs/deploy/{local,remote}.md`, `DECISIONS.md` ADR list,
`.agents/docs/plans/2026-09-28-recon-from-zip.md`).

## Status: COMPLETE

## Verdict: REQUEST CHANGES (design needs revision)

The worst problem: the plan's own refusal rule for `UnpackZip`, read literally, refuses the exact fixture its Golden-path
E2E criterion depends on. `.agents/skills` is a real, git-tracked symlink in this repository (`120000 blob`, confirmed
with `git ls-files -s .agents/skills`), and GitHub's "Download ZIP" and zipball API both preserve a tracked symlink as
a zip entry carrying the Unix symlink mode bit. The plan's refusal list says the whole zip is refused when "any entry
... is a symlink" — unscoped — then separately says entries outside the five extracted trees are only "name-checked."
Those two sentences do not agree on whether `.agents/skills` triggers the refusal, and the plan never says which
reading is correct. Two implementers reading this plan would build differently, and the E2E golden-path test target
named in the plan (`gh api repos/kesokaj/styrman-harness/zipball/main`) is a strong candidate to fail as soon as it is
built literally.

## Findings

- **[high / confidence: high]** Plan `## Interfaces` § `harness.UnpackZip`, bullets under "It refuses the whole zip ...
  when any entry:" vs. "Extracted: ... Every other entry is still name-checked, but never written." — internal
  contradiction over whether the symlink/mode/size refusals apply to every entry in the archive or only to entries
  under the five embedded trees. `.agents/skills` is tracked as a git symlink (`git ls-files -s .agents/skills` →
  `120000 blob 42c5394a...`), and GitHub's zip export (both the web button and `zipball/main`, the exact fixture named
  in Phase 3's e2e task and the Golden-path criterion) encodes a tracked symlink as a zip entry with the symlink mode
  bit set. If the "is a symlink" refusal is read as applying to every entry (as its own sentence says, unscoped), the
  real repository's own zip is refused outright and the Golden-path E2E and Phase-3 e2e task both fail against the
  plan as written. If it is read as scoped to the five trees (as "every other entry is still name-checked" implies),
  a symlink placed inside `context/`, `skills/`, `agents/`, `commands/` or `rules/` would need the mode check anyway,
  which the plan does say ("is a symlink or any non-regular, non-directory mode" is listed first, unscoped, before the
  Extracted/Layout section even appears) — so the ordering of the two sections in the plan does not resolve which
  wins. This must be decided (scope the symlink/mode/size checks to the five extracted trees; name-check only,
  everywhere else) before `UnpackZip` is implemented, or the first Phase-3 e2e run burns a cycle discovering it.
- **[medium / confidence: medium]** Plan `## Interfaces` § "Order of work" step 4 (`src := harness.Source{FS: ...,
  Dir: dir, Origin: ...}`) vs. `internal/harness/apply.go:117-123`, which currently branches purely on
  `if src.Dir != ""` to decide whether the manifest records `src.Dir` (a tree's path) or `EmbeddedPrefix+src.Version`.
  For a zip install, `Source.Dir` is set (to the temporary unpack directory) *and* `Origin` is set. The plan's
  Interfaces text says the manifest should record `Origin` "instead of Dir, when set," but the "Surfaces touched"
  list only says apply.go:117-123 changes, without stating that the existing `if src.Dir != ""` branch order must be
  changed to check `Origin` first. Read literally against the current code, an implementer who edits only what's
  described could leave the manifest recording `src.Dir` — the ephemeral temp directory that Order-of-work step 6
  deletes moments later — rather than `zip:sha256:<hex>`, which would break `status`'s zip detection and
  `ResolveSource`'s skip. Likely resolved correctly by inference, but the plan does not spell out the precedence,
  and this is exactly the kind of place where "obviously right" silently regresses (the temp dir is gone by the time
  anyone reads the manifest, so a wrong write here is not caught by output, only by a later `status` or `update` that
  is quietly wrong).
- **[medium / confidence: medium]** Plan `## Interfaces` § `status`: the added "zip" source-line suffix
  ("this home was installed from a zip ... run update --from-zip <file> again") needs `st` (`harness.LoadState`)
  to decide what to print on the "source"/"repo" line (`cmd/sh-agent/main.go:346-368`), but in the current code
  `LoadState` is not called until line 370, *after* that line is already written (`main.go:352`/`364`/`366`). The
  plan's Surfaces-touched entry ("status (346-472, the from line 382-385 and the drift block 444-471)") does not
  flag that this specific addition requires reordering (load state before printing the source line, or holding the
  suffix and appending it after). Not a design ambiguity — a reasonable implementer resolves it one way — but it is
  an unstated refactor inside a function the plan otherwise treats as a set of independent edits.
- **[low / confidence: low]** `## Why` / `## Scope` do not address the existing "plain copy" path documented in
  `delivery.md` § How a machine gets the harness ("A plain copy with no `.git` — the unpacked zip — still installs
  ... it is the fallback for a machine that can reach neither of [a checkout or a release binary]") and
  `deploy/local.md:55`. That paragraph already describes manually downloading and unzipping the repo zip, then
  running `sh-agent install --repo <dir>` from a dev build — materially the same operator problem this plan solves,
  via a different, already-shipped mechanism that leaves the unpacked tree on disk permanently instead of deleting
  it. The plan lists `delivery.md` as a doc surface to touch but never states whether the old manual-unzip-and-install
  path is now deprecated in favour of `--from-zip`, kept as a separate fallback, or how the docs should describe the
  relationship. A first-pass doc edit could leave both described as live options with no guidance on which to use,
  or silently contradict the "it is the fallback" framing once a third, cleaner path exists. Low confidence this
  rises to a design blocker — the `documentation` step in Phase 3 could reasonably resolve it inline — but it is an
  undecided scope edge the plan should have named.

## Gaps (per the task's checklist)

- Symlink/mode/size-refusal scope (five trees vs. whole archive) — **undecided, blocking** (finding 1 above).
- `Origin` vs. `Dir` precedence in `apply.go`'s manifest-recording branch — **under-specified, non-blocking**
  (finding 2).
- `status`'s source-line ordering needs `LoadState` moved earlier — **unstated refactor, non-blocking** (finding 3).
- Relationship between `--from-zip` and the pre-existing "plain copy / unpacked zip" install path in the docs —
  **undecided, non-blocking, worth a one-line decision** (finding 4).
- `harness.ConfigDir(home)` returns `(string, error)`; the plan's step-2 pseudocode
  (`harness.UnpackZip(file, harness.ConfigDir(home))`) reads as passing the two-value call directly, which is not
  valid Go against `UnpackZip(zipPath, parent string)`. Almost certainly just plan-prose shorthand, not a real
  ambiguity — noted, not a finding.
- Checked and found sound, not gaps: ADR numbering (next free is ADR-038, matches the plan); `-ldflags
  "-X main.version=..."` matches `var version string` in `main.go:55-57`, so the Phase-3 e2e release-versioned build
  command is correct; the five embedded trees named in the plan (`context, skills, agents, commands, rules`) exactly
  match `embed.go:9`'s `//go:embed` directive (the `workflows` mentioned in `source.go:19`'s doc comment is a stale,
  pre-existing comment — there is no `workflows/` source directory in this repo — not something this plan touches or
  needs to touch); `repo.go:75`'s existing `EmbeddedPrefix` skip is exactly the line the plan says to extend for
  `ZipPrefix`, and the pattern is a direct copy of working code; the `grep -iE 'two ways|no third'` check in Phase 3
  would catch every restatement found in this review (`AGENTS.md:13`, `delivery.md:34`, `deploy/remote.md:56`,
  `deploy/local.md:55`).

## Scope

Everything in `## Scope: In` traces to a D-row or to the operator's "approve as written, go all but number 6."
Item 7 (the matrix) is explicitly attributed to D5, read from the same approval. No unrequested endpoint, flag, or
configuration surface was found — the plan does not add a second shape for `update` or a knob for a value that never
varies. `--vendor`/`--dry-run`/`--home` on the new flag are the existing ones every command already takes, not new
surface. Nothing to flag here.

## Readiness

Questions a first-pass implementer would have to ask, or guess at:

1. **BLOCKING** — Do the symlink/non-regular-mode/size refusals in `UnpackZip` apply to every entry in the zip, or
   only to entries under the five extracted trees? (Finding 1. Two implementers diverge, and one produces a
   Golden-path E2E that fails against the plan's own named fixture.)
2. **NON-BLOCKING** — In `apply.go`, does the manifest-recording branch check `Origin` before `Dir`, given a zip
   `Source` sets both? (Finding 2. A competent implementer reading "the manifest records `Origin` instead of `Dir`,
   when set" gets this right, but the plan doesn't say so against the current `if src.Dir != ""` code.)
3. **NON-BLOCKING** — Does `status()` need `LoadState` moved earlier, or does it hold the zip-suffix text and append
   it after the existing state load? (Finding 3. Either resolves the same output; a first-pass implementer picks one
   without needing to ask.)
4. **NON-BLOCKING** — How should the docs describe the relationship between `--from-zip` and the pre-existing
   "plain copy from an unpacked zip" install path? (Finding 4. The `documentation` step in Phase 3 can reasonably
   decide this inline; it doesn't block writing code.)

## Not examined

- `internal/harness/report.go`, `release.go`, `state.go` (beyond `ConfigDir`'s signature) — not named as touched
  surfaces by the plan and not read in full.
- `cmd/sh-agent/main_test.go` and any existing `zip_test.go`-adjacent fixtures — there is no `zip.go`/`zip_test.go`
  yet (both are new files per the plan), so nothing existed to read.
- `.agents/docs/testing/skill-matrix.md` and the three findings in `reviews/2026-09-28-review-v1.7.0-381e86d-code.md`
  that Phase 4 must close — confirmed the referenced review file exists on disk; its findings 1-3 were not read in
  detail, since Phase 4 (matrix + release) is downstream of the design questions in Phases 1-3 and this cold read
  focused there given the tool budget.
- README.md § How it works (the map) itself was not opened; only grep-matched for "two ways" style text. The map's
  drawing conventions (dashed vs. solid, per AGENTS.md rule 4) were not checked against what a zip line should look
  like.
- `ADR-034`'s full text (only the portion quoted in `repo.go`'s comments was read) — the plan's claim that ADR-038
  "supersedes ADR-034 for this one path only" was not verified against ADR-034's full wording.

design needs revision
