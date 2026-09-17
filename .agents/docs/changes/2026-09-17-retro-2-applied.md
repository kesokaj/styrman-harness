---
owner: main
date: 2026-09-17
---
# Second retro applied: dispatch text for must-not-skip duties, changes entries, escapes in input, behavioural runs     (2026-09-17)

## Summary
The operator approved all four proposals of `.agents/docs/retros/2026-09-17-retro-2.md`. Three change `context/AGENTS.md`: a duty a subagent must not skip is repeated in the dispatch text (the verifier's own rules failed intermittently in BUG-062 and BUG-063; the dispatch sentence held in BUG-064); the main agent may write the `changes/` entry itself from the diff; and a check for an escape builds its needle from character codes, because the conversion of BUG-057 also hits prompts and shell commands. The fourth changes the project `AGENTS.md` § Testing: an agent-definition change gets a behavioural run with `claude -p --agents`, and a `context/` change meant to alter the main agent with `claude -p --append-system-prompt`, the real home untouched. The docs review's findings carried the step 4 allowance into the bug, feature and preflight skills.

## Type
task

## Files changed
- `context/AGENTS.md` — § Work (dispatch text for must-not-skip duties), § Definition of Done step 4 (self-written entry), § Hard rules (escape needle from character codes)
- `AGENTS.md` — § Testing (behavioural runs)
- `skills/bug/SKILL.md`, `skills/feature/SKILL.md`, `skills/preflight/SKILL.md` — "`changelog` agent or yourself"

## Agents dispatched
reviewer (docs lens, APPROVE WITH NITS, three findings fixed), verifier (PASS). This entry was written by the main agent from the diff.

## Verification
- `.agents/docs/testing/2026-09-17-verify-retro-2.md` PASS: gofmt, vet, `go test -race -count=1 ./...`; scratch install drift none; the three context sentences in both vendors' copies.
- After the review fixes: re-installed into `.agents/scratch/retro2-home`, drift none; `go test -count=1 ./...` passes.
- Review `.agents/docs/reviews/2026-09-17-review-ec22479-retro-2-docs.md`.

## Docs
Context, project AGENTS.md and three skills, above.

## Requirements
n/a

## Not exercised
- A CLI start in a sandbox home; the real-home smoke checks with this context (the operator's `sh-agent update`).
