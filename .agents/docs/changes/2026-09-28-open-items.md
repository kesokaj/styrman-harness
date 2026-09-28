---
owner: main
date: 2026-09-28
---
# Open follow-ups: three fixed, two struck, two deferred     (date: 2026-09-28)

## Summary
The operator said "go all" on the follow-ups carried since the 2026-09-21 handoffs, and on the `owner:` question from the OKF record. Each got the action recommended for it:
- **Fixed:**
  - the project `AGENTS.md` narrative trimmed to pointers;
  - BUG-011's index row moved to the Graduated table;
  - the unnumbered `cloud-run-deploy` pre-flight bullet numbered.
- **Struck:**
  - the `owner:` question (not a defect: across the trail, `owner:` always names the file's writer, and the init-project deploy templates say `documentation` because that agent writes them in a project);
  - an over-cap line in `sh-agent status` (`TestRenderedContextFitsAntigravityBlock` already blocks it in the repo);
  - the one headless `/init-project` run that made no commit (happened once; the matrix `init-deploy` row catches a repeat).
- **Deferred** (the handoff names the triggers):
  - the second-rules-file probe, with its `size_test.go` per-file allowance and retro-table row, until before the first `rules/` file;
  - the Antigravity `permissions.allow` restore on this box, until headless Antigravity checks with tools are wanted.

## Type
task

## Files changed
- `AGENTS.md`:
  - the toolchain paragraph under Setup commands is shortened to its rule with an ADR-029 pointer;
  - rule 6's v1.3.0 anecdote is replaced by a pointer to `retros/2026-09-18-retro-2.md`.
- `skills/cloud-run-deploy/SKILL.md` — § Pre-flight: the service-account bullet became item 6, **Service account**.
- `.agents/docs/bugs/INDEX.md` — BUG-011's row moved to Graduated (`context/AGENTS.md` § Hard rules, "Staging is a verification step too").
- `.agents/docs/testing/2026-09-28-verify-open-items.md` — verifier log, PASS.

## Agents dispatched
- `bugs`: moved the BUG-011 row.
- `verifier`: PASS.

## Verification
- **`verifier`** (`testing/2026-09-28-verify-open-items.md`):
  - gofmt, vet and `go test -race`: 251 pass, 0 fail;
  - scratch-home install and status: drift none;
  - the installed skill carries the new line;
  - `list` still parses `cloud-run-deploy`;
  - frontmatter untouched.
- **Not run:**
  - a CLI sandbox start, because only the skill's body text changed;
  - `/review --docs` over the skill change. Instead, the main agent grepped `skills`, `agents`, `context`, `commands`, `rules`, `AGENTS.md` and `README.md` for references to cloud-run-deploy pre-flight step numbers. The one hit, `agents/verifier.md:18`, refers to the verifier's own Pre-flight step 2, not this skill.
- **No `reviewer`:** 3 lines across 3 files. The main agent read the diff.

## Docs
`AGENTS.md` (project instructions) only, as above.

## Requirements
n/a
