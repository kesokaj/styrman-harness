---
owner: reviewer
date: 2026-09-29
---
# Review — abbaddc-1634787-docs   (2026-09-29)

Range reviewed: abbaddc..1634787 (`git diff abbaddc...1634787`)

## Status: COMPLETE

## Verdict: APPROVE

Nothing in this range is broken. `.agents/docs/specs/sh-agent/search.md` was updated in step with the code: the new symlink/non-regular-file refusal row matches `prepareFiles` in `internal/search/store.go` exactly (error text and the `search: <path>:` wrapper from `openStore`, which is not misrouted into the rebuild-on-unreadable path since `isUnreadable`/`isBusy` only match sqlite error codes, not the plain `fmt.Errorf` `prepareFiles` returns); the 0600 file-mode claim matches; the rebuilt/smaller-index numbers (1.7 s, 21.7 MB, dropped `refs(oid)`) match the plan's measured D8 entry and the code diff (`refs(oid)` index removed from `store.go`'s schema, insert statements now prepared once per reconcile via the new `inserter` type, `parseRefs` prefiltered by `strings.Contains` before running the regexes). `context/AGENTS.md`'s resume order gained `sh-agent search --open`, and the Search first entry correctly dropped its own `--open` clause (it now only names `--refs` and `--global`); `skills/handoff/SKILL.md`'s cold-reader check, the one other place that restates the Resume Protocol's step list, picked up the matching `--open` addition. No other living doc (README.md, agents/*.md, skills/ non-vendored, `.agents/docs/specs/`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/`) restates the old resume order, the old Search-first wording, the old file counts/build time/size, or `refs(oid)`; the only remaining old numbers live in `changes/`, `reviews/` and the plan file itself, which are excluded from this check by design (historical record, not living doc). `ROADMAP.md`'s "Now" line pointing at the release-readiness plan is correct, not stale: the plan's own Checkpoint and Definition-of-Done boxes (lines 66, 78) are still unticked at this range's head, so the plan is genuinely still active and no `changes/` entry was owed yet.

## Findings
None.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- The vendored skills list (`skills/VENDORED.md`) was not diffed against `skills/` for drift — irrelevant to this range, which touches no skill directory besides `skills/handoff/SKILL.md` (checked).
- The two new review artifacts (`reviews/2026-09-29-review-aae16af-f41e5c0-code.md`, `reviews/2026-09-29-security-aae16af-f41e5c0.md`) and the coldread rounds were not fact-checked against the code they reviewed — they are trail records, explicitly out of scope for the staleness grep, and this review does not re-review a prior reviewer's findings.
- `.agents/docs/testing/2026-09-29-search-adoption.md`, `-search-agy-gate.md`, and `-verify-real-home-search.md` were read only for the numbers that also appear in `context/AGENTS.md`/`skill-matrix.md`/the plan; their internal procedural detail was not independently re-verified against a live run.
- `govulncheck`/`go test`/`gofmt` were not re-run; this is a docs-only review per the brief.
