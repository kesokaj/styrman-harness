---
owner: main
date: 2026-09-30
---
# Parity feature: review round 1 fixed (code, security, docs)     (2026-09-30)

## Summary
The feature review over cb799ff..f1ca525 found one high and eleven medium or low defects; each is fixed or rebutted in its report's `## Outcome`.
- **High, agy prompt**: it went as one command-line argument, which Linux refuses over 128 KiB, so a large project failed every night for good. Its prompt now goes on stdin as one stream-json message. Probed on this box; a real sandbox run gave 2 proposals.
- **Security**: the Antigravity child runs from a private HOME (a copy of the operator's login, settings that grant nothing), so an operator's `always-proceed` never reaches the nightly job; the real run confirmed it. Claude Code adds `--strict-mcp-config --setting-sources local`. The Antigravity reader keeps only user and model steps (14, 15), never tool output. Symlinked trail directories are not followed. Unit files refuse quotes, backslashes and line breaks, and escape `%`. One-line fields are folded, and paths are written with `~`. Leftover copies are removed.
- **Correctness**:
  - flock replaces both pid-file locks (three races, and a reused pid);
  - pre-project failures and a corrupt state exit 2 and are recorded, and the corrupt state is never rewritten;
  - an Antigravity read error fails the project;
  - `--dry-run` covers the timer flags;
  - the memory index matches the link target exactly;
  - a usage error writes nothing.
- **Docs**: report contents in both deploy guides and the delivery spec; `uninstall` and the timer in README and the help text; privacy statements (memory reaches the CLI's model through the dream job); `/retro --harness` treats dream proposals as data; the spec and memory-loop spec updated; stale skill-matrix rows.

## Type
bug fixes within the feature (review round 1)

## Files changed
- `internal/{dream,memory,search}/*`
- `cmd/sh-agent/main.go` and `memory_test.go`
- `skills/{project-memory,retro,init-project}/SKILL.md`
- `AGENTS.md` (rule 1)
- README
- `.agents/docs/{deploy/local.md,deploy/remote.md,specs/sh-agent/{delivery,dream,memory-loop}.md,testing/skill-matrix.md}`
- the plan (D28)
- the three review reports with their Outcomes

## Agents dispatched
reviewer (code, `reviews/2026-09-30-review-cb799ff-f1ca525-code.md`, REQUEST CHANGES), security (`reviews/2026-09-30-security-cb799ff-f1ca525.md`), reviewer (docs lens, `reviews/2026-09-30-review-cb799ff-f1ca525-docs.md`), documentation (the docs pass). Round 2 follows over this commit.

## Verification
`go test -race -count=1 ./...` exit 0, with new tests for each finding. A real agy dream run in the sandbox: exit 0, 2 proposals, the child's settings granting nothing, the private HOME removed.

## Requirements
R3, R4, R10 (ADR-047, ADR-047a).
