---
owner: reviewer
date: 2026-09-16
---
# Review — bug-035-code (2026-09-16)

Range reviewed: working tree diff (`git diff HEAD --`) over `README.md`, `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`, `internal/harness/apply.go`, `internal/harness/harness_test.go`

## Verdict: APPROVE

Nothing is broken. All three guidance sites (`README.md:53`, `cmd/sh-agent/main.go:309`, `internal/harness/apply.go:37`) now state, in the same sentence as the "move into `rules/`" instruction, that `rules/` ships to every install and mirror and that credentials/personal notes never go there — matching BUG-035's Prevention rule verbatim in spirit. Each still ends with the actionable "remove the file, run install" step, so `status`/the `keep` line still tells the reader how to clear the flag. Both Go call sites have their own assertion (`main_test.go:347` for the `status` note in `main.go:309`; `harness_test.go:257` for the `keep` line in `apply.go:37`), each pinned to source text that would fail if that specific message were reverted independently of the other. The updated failure message at `harness_test.go:258` ("keep line lacks the rules/ guidance or its exclusion of credentials and personal notes") now names what the two `Contains` checks verify, replacing the stale "R2" reference. `gofmt -l .`, `go vet ./...`, and `go test ./...` all pass clean. No stale "move its content" wording remains in live code, comments, doc-comments, or test names — the only surviving occurrences are in append-only historical records under `.agents/docs/reviews/` and `.agents/docs/testing/`, which correctly describe past state, not current behavior.

On the goldfish-diagnosis skip: I agree it was allowed. The bug record (`BUG-035`) names all three sites by file:line and the fix at each is a one-sentence text edit with a matching one-line test-assertion update — squarely inside the bug skill's exemption for "a one-to-three-line fix in a clearly identified place."

## Findings
None.

## Criteria
1. Dead code — pass: no orphaned code; the stale "R2" wording in the failure message was replaced, not left as a comment.
2. Flow tracing — pass: `main.go:309` and `apply.go:37` each have an independent, failable test; text edits have no callers to break, and both changed messages are traced to their exact test.
3. Standards conformance — pass: matches BUG-035's Prevention rule (same-sentence exclusion) and REQUIREMENTS.md R4/R5's hand-off model; describing-text rule satisfied, no comment left pointing at removed behavior.
4. Directory layout — pass: no files moved or added, all edits in place.
5. Simplicity — pass: pure text edits, no new abstraction, no new dependency.
6. Decoupling — pass: the two Go messages and their tests stay in their own files/packages (`cmd/sh-agent` vs `internal/harness`); no cross-package reach-in introduced.

## Not examined
- BUG-035's `status:` field and `.agents/docs/bugs/INDEX.md` row still say "open" — expected, since the bug-record update and INDEX row are the `bugs` agent's job later in the Definition-of-Done sequence, and neither file is in this diff's scope.
- The downstream mirrors (styrman-crew, argolis-devbox) that BUG-035 says will need the same guidance text — out of scope for this repo's diff; not verified here.
- No security or `/review --docs` pass was run as part of this review; the dispatch scoped this to code-criteria review only.

## Outcome

- **No findings.** The one carried note, BUG-035's record still open, is closed in this commit. After this review the docs lens found the Go messages narrowed against the README by the words "or host-specific"; both messages and their assertions were aligned (`reviews/2026-09-16-review-bug-035-docs.md`, round 1), and each message keeps its own failable assertion.
