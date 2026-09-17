---
owner: main
date: 2026-09-17
---
# Release binary and a report that keeps nothing, landed on `main` — Phase 5     (2026-09-17)

## Summary

The release-binary work reached `main` as one commit, `cd0bbcb`, as R9 requires: code, docs, the map, `REQUIREMENTS.md` R8 and R9, and the ADRs together. What changes for a reader of this repository: there are two ways in, a checkout and a release binary with the harness built into it; a tag is now a release rather than only a marker, because it carries the binaries a binary-only place updates to; and `sh-agent report` keeps nothing of the feedback repository on a machine. The unpacked bundle leaves the docs and the map and stays as the fallback for a machine that can reach neither way in — its code path is untouched.

Phase 5 itself was documentation and the merge. `README.md` § How it works carries the redrawn map (a Releases node, the release binary as the second way in, `report` drawn as a throwaway clone) and a new § Install without cloning with the three ways to get the binary, the macOS quarantine note, and what `update` does at such a place. `.agents/docs/deploy/remote.md` § Releases gained the whole asset procedure — the pre-build refusal of anything untracked or ignored under the five embedded directories, the three cross-compiled binaries, `SHA256SUMS` with its version marker, the signature, the verify-before-upload step, `gh release create` — and a § Release notes section. `skills/preflight/SKILL.md` Step 6 gained release notes and a generic release-artifacts step, which every project's preflight now follows when its deploy guide names artifacts. `.agents/docs/deploy/local.md` gained § Update a release-binary install and turned its bundle section into § From a plain copy. `AGENTS.md` rules 1 and 2, the delivery spec, and the report-record documentation follow the same change.

## Type

feature

## Files changed

The merge commit `cd0bbcb` carries 54 files. Phase 5's own edits: `README.md` (map, § Install without cloning, § The loop, the `install`/`update`/`report` rows, the state table), `REQUIREMENTS.md` (R8's first sentence, R9's four bullets), `AGENTS.md` (rules 1 and 2, the install-rebuild note, the `govulncheck` call-path note), `skills/preflight/SKILL.md` (Step 6 items 1, 5 and the renumbering), `.agents/docs/deploy/remote.md` (§ Releases, § Release notes, § The four places, rollback), `.agents/docs/deploy/local.md` (§ Update a release-binary install, § From a plain copy, § Report back), `.agents/docs/specs/sh-agent/delivery.md` (how a machine gets the harness, the feedback paragraph, the follow-up now built), `.agents/docs/architecture/DECISIONS.md` (ADR-033a appending D19–D28). Phases 2–4 are covered by their own entries.

## Agents dispatched

`reviewer` with the docs lens over `edbddca..1a5578b` — REQUEST CHANGES, 1 high and 1 low, `.agents/docs/reviews/2026-09-17-review-edbddca-1a5578b-docs.md`. The high finding was a paragraph in README § The loop that still framed R8 as "a checkout … or an unpacked copy", contradicting what the same commit had rewritten five other places; the low was an incomplete `install` row. Both fixed, and the sweep they prompted found a third instance in `local.md` § Report back, where a release-binary place's record — the one that sets the new `harness_version` — was undocumented. Outcome written in that report. This is the fourth time a rename has left a stale restatement here (BUG-003, BUG-007, BUG-027), and the rule those recorded is what found the third one: grep the subject NOUNS across every living doc, not the old wording.

## Verification

`gofmt -l cmd internal *.go` silent, `go vet ./...` clean and `go test -race -count=1 ./...` green on the staged merge before it was committed, and again after; `TestHarnessMap` passes and caught a map node labelled `sh-agent binaries`, which reads as a command the binary does not have. A scratch install from this tree reports drift none and `sh-agent list` parses every skill, which is what covers the `skills/preflight/SKILL.md` edit. The suite stands at 99 top-level tests and 85 subtests, none failing or skipped, with `govulncheck` clean (go1.27.1). Not exercised, and owed before or with the release: the scratch rehearsal of the release procedure (build the three binaries, sign the checksums, run `UpdateFromRelease` against them through a fake `gh`), the dry build of v1.1.0's release notes for the operator to read, and the e2e on the work Mac — all of which belong to `/preflight v1.1.0`, which is the operator's to run and which they deliberately have not started.

## Docs

This entry IS the docs change. No living doc now presents the unpacked bundle as a way in, describes a kept feedback clone, or says a tag delivers nothing; the map, R9 and `AGENTS.md` rules 1 and 2 agree, and the docs review checked that agreement rather than taking it on trust.

## Requirements

R8's first sentence and R9's four procedures now read as the operator approved them on 2026-09-17, each carrying its approval note and ADR-033; the map they refer to is the one the operator approved in their own turn ("Approve the redrawn map"). R4 holds: stdlib only, `gh`, `git` and `ssh-keygen` shelled out. R9's requirement that the procedures, the code, the docs and the map change in ONE commit is satisfied by `cd0bbcb`.
