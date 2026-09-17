---
owner: main
date: 2026-09-17
---
# BUG-064: the verifier dispatch names a missing testing guide     (2026-09-17)

## Summary
On Claude Code the verifier's own rule to write `testing/guide.md` (step 4, BUG-063) was followed in some runs and not others — the real-home smoke check at e67eace skipped it, and a repro wrote it in 2 of 3 runs, one of which also logged `Scope: complete`. Asked how the guide should reliably get written, the operator chose "Dispatch says it": the main agent's Definition of Done step 1 now checks for `<docs-root>/testing/guide.md` before dispatching the verifier and, when it is absent, puts "testing/guide.md is absent: write it first" in the dispatch. The `review` skill's `--tests` row and `preflight`'s Verification row carry the same sentence.

## Type
bugfix

## Files changed
- `context/AGENTS.md` — Definition of Done step 1: the check and the dispatch sentence
- `skills/review/SKILL.md` — `--tests` dispatch text
- `skills/preflight/SKILL.md` — Verification dispatch text

## Agents dispatched
reviewer (APPROVE; one finding fixed, one rebutted), verifier (PASS), bugs (BUG-064 fixed; a note on BUG-063). No explorer: the operator chose the design; no changelog agent: this entry was written by the main agent from the diff.

## Verification
- Behaviour, `.agents/docs/testing/2026-09-17-bug-064-repro.md`: without the new step 1, the guide was written in 2 of 3 runs; with it (given via `--append-system-prompt`, the real home untouched), the dispatch carried the sentence 3 of 3 and the guide was written 3 of 3.
- `.agents/docs/testing/2026-09-17-verify-bug-064.md` PASS: gofmt, vet, `go test -race -count=1 ./...`; scratch install drift none; the sentence present in both vendors' context, review and preflight copies. The later "(`ls`)" wording in step 1 was re-installed by the main agent with drift none.
- Review `.agents/docs/reviews/2026-09-17-review-7960273-bug-064.md`.

## Docs
Context and skill text above; no README, spec or deploy guide describes the verifier dispatch.

## Requirements
n/a

## Not exercised
- Antigravity (it wrote the guide without the sentence in both real-home checks)
- `/review --tests` and `/preflight` run end to end with the new text
- More than three runs per arm; the installed context in a real home (the operator's `sh-agent update`)
