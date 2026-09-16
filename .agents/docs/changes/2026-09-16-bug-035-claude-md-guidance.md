---
owner: main
date: 2026-09-16
---
# BUG-035: the CLAUDE.md guidance would have sent personal notes into every mirror     (2026-09-16)

## Summary

c514936 taught the installer's keep line, `sh-agent status`'s note and `README.md` to tell the owner of a hand-written `~/.claude/CLAUDE.md` to move its content into `rules/` in the harness repo. Nothing said that `rules/` is deployed to every install and copied into every downstream mirror, including styrman-crew's runtime image, so a reader following the advice with hostnames or tokens in that file would have shipped them everywhere. styrman-crew's security review of the a81ef4b refresh found it (finding 5); it was recorded as BUG-035.

All three sites now say the same thing in one sentence: move only rules meant for every machine into `rules/`, which ships to every install and mirror, never credentials or personal or host-specific notes; then remove the file and run install. The two tests that assert the wording moved with it.

## Type

bug

## Files changed

- `internal/harness/apply.go:37` — the keep line.
- `cmd/sh-agent/main.go:309` — `status`'s note.
- `README.md:53` — the prose, edited by `documentation` with the exact sentence.
- `internal/harness/harness_test.go` — `TestInstallStatusResyncUninstall` asserts the new phrase, and its failure message now describes what it checks instead of the retired wording.
- `cmd/sh-agent/main_test.go` — `TestStatusNamesForeignClaudeMD` asserts the new note.
- `.agents/docs/bugs/BUG-035-claude-md-guidance-ships-personal-content-everywhere.md` and `INDEX.md` — status fixed.
- `.agents/docs/testing/2026-09-16-verify-bug-035-claude-md-guidance.md` — the verifier log, second run.
- `.agents/docs/reviews/2026-09-16-review-bug-035-code.md` (APPROVE) and `…-docs.md` (round 1 APPROVE WITH NITS, round 2 APPROVE), each with its `## Outcome`.

## Agents dispatched

- `documentation`: the README sentence.
- `verifier`: twice, the second after the wording was aligned.
- `reviewer`: the code lens once (ADR-022), the docs lens twice.
- `bugs`: closed BUG-035.

No goldfish diagnosis: the record names all three sites by file:line and each change is one sentence, the bug skill's exemption; the code reviewer agreed.

## Verification

- **Test first.** With the assertions tightened and the messages unchanged, both tests failed; with the messages changed, both pass. The verifier reproduced that against `git show HEAD:<path>` in scratch.
- **Verifier, Status PASS:** counts pass 54, fail 0, skip 1, unchanged, since no test was added; `gofmt`, `go vet` and `go test -race ./...` clean; `list` 27 skills and 7 agents; `render` 93 lines; a scratch install with drift none, and a second with a planted hand-written `CLAUDE.md` whose keep line and `status` note print the final wording word for word.
- **Reviews.** Code: APPROVE, each message with its own failable assertion. Docs round 1: APPROVE WITH NITS — the Go messages had dropped "host-specific" against the README, a `narrowed:` finding, fixed by two words in both messages and both assertions. Docs round 2: APPROVE, the three sites stating the same rule with the same qualifiers.

## Docs

`README.md:53` is the living doc that changed, and it now matches the two messages verbatim. The docs lens found no other living doc with the old wording. Both downstream mirrors receive the new text at their next refresh; today their `rules/` hold only `.gitkeep`.

## Requirements

n/a
