---
owner: main
date: 2026-09-15
---
# Harness practices Phase 7, part 1: the operator's decisions on R2 and the release gate, and the repo docs finalized     (2026-09-15)

## Summary

On 2026-09-15 the operator kept requirement R2 ("Yes. We should keep the AGENTS.md style. Do your recommendation.") and chose the recommended rule for docs findings at the release gate ("Do your recommendation."). Both land here, with ADR-023, and close the two round 2 findings the Done review left for the operator. The repo docs drift the 2026-09-15 sweep found in `README.md`, the deploy guides, `AGENTS.md` and `skills/VENDORED.md` (F11–F19) is fixed. Phase 7 of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` holds the rest of the operator's requests.

## Type

task

## Files changed

- `skills/preflight/SKILL.md` — :65 Step 4 records each docs finding with severity and confidence; a finding the newest prior gate report's `## Outcome` rebutted, the same statement at the same file:line, is cited and not raised again; a blocking gap is fixed after NO-SHIP as ordinary work and `/preflight` runs again; every other gap is `left` with an owner. :77 only a Step 4 REQUEST CHANGES verdict or high-severity docs finding blocks, not counting a rebutted one. The "never edits the tree" sentence is gone.
- `skills/init-project/SKILL.md` — :44 when the user keeps an existing `CLAUDE.md` with content, the summary says Antigravity never reads that content and Claude Code imports `AGENTS.md` only if the file has the line `@AGENTS.md`, so the project follows R2 once the content is folded.
- `README.md` — :24, :26, :29 the `--dry-run` and `--force` flags; :32 the `--vendor` default and how sh-agent finds its checkout; :53 an existing `~/.claude/CLAUDE.md` has its content moved into a rule file under `rules/` so `sh-agent install` writes the one-line import; :57 an unquoted `#` at the start of a frontmatter value.
- `.agents/docs/deploy/local.md` :27 and `.agents/docs/deploy/remote.md` :29, :33, :37 — the `--vendor` default, checkout discovery, how each mirror ships a refresh and builds `sh-agent`, and the Mac Antigravity check's 2026-09-04 pass instead of a `TODO.md` pointer.
- `AGENTS.md` :35 — a mirror's image or VM build step built its binary.
- `skills/VENDORED.md` :36 — `anti-ui-slop` ships its own LICENSE and NOTICE; no manifest file exists.
- `.agents/docs/architecture/DECISIONS.md` — ADR-023, the release gate weighs docs findings like review findings.
- `.agents/docs/bugs/BUG-026-preflight-step4-routing-open.md` (now fixed), `BUG-025-init-project-claude-md-requirement.md` and `INDEX.md`.
- `.agents/docs/reviews/2026-09-15-review-c192622-22c106c-merged.md` and `2026-09-15-review-sweep-merged.md` — outcome updates naming the fixed lines.
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — Phase 7 with tasks 7.1–7.6, decision rows D45, D46 and D48, tasks 7.1 and 7.2 ticked.

## Agents dispatched

documentation (README and deploy guides), bugs (BUG-026, BUG-025), changelog (ADR-023). A workflow surveying the seven projects with an older `.agents/` tree, for the `/init-project` update path of task 7.3, was running at this commit. This entry was written by the main agent.

Review: the Phase 7 range gets its `reviewer` and `--docs` review when Phase 7's changes are complete.

## Verification

The main agent ran `go test -race ./...`, `gofmt -l .`, `go vet ./...`, a scratch install with `status`, and `render`, each bare with its exit status read; the results are in the commit message. Not run: a sandbox start of each CLI for these changes, because the Phase 6 matrix rows are running in the sandbox home; Phase 7's checkpoint covers it.

## Requirements

n/a — the operator's decision keeps R2 unchanged.
