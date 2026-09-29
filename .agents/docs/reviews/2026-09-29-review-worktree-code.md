---
owner: reviewer
date: 2026-09-29
---
# Review — worktree-code (2026-09-29)

Range reviewed: working tree (unstaged + staged) — `.agents/docs/testing/skill-matrix.md`, `skills/VENDORED.md`, `skills/review-standards/SKILL.md`, `skills/review/SKILL.md`, `skills/review-standards/references/prompt-audit/LICENSE.txt` (new), `skills/review-standards/references/prompt-audit/prompt-audit.md` (new)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing broken. This wires a new `/review --prompts` lens onto a verbatim-vendored Anthropic reference file: the vendoring is correctly pinned (digest matches, file untouched), the installer places both new files under the right skill directory on both vendors (verified with a scratch install), the dispatch text's paths for both CLIs are correct against `internal/harness/plan.go`'s `AgyRoot`/`ClaudeRoot`, and `go build`/`go vet`/`go test ./internal/harness/...` are clean. The one real gap is in the new skill-matrix test row itself: its check script only asserts 3 of the 4 fixture facts the scenario sets up, so a regression in the most basic pattern group (Group 1a pressure language) could pass this scenario silently.

## Findings

- **[medium / confidence: medium]** `.agents/docs/testing/skill-matrix.md:31` (`review-prompts` row) — the scenario commits four fixture facts (a Group 1a pressure-language line `CRITICAL: You MUST ALWAYS run gofmt...`, a stale path `cmd/nothing/main.go`, a Group 2 conflicting-instructions pair `make test` vs `go test ./...`, and a provenance-protected BUG-001 line), but the `check` column only asserts three: `## Proposed diff`, `cmd/nothing`, `make test`, plus the negative BUG-001/remove check. Nothing in the check requires the report to mention or act on the `CRITICAL: ... gofmt` line at all. If the `--prompts` lens implementation (or the vendored file's Group 1a section) regressed and stopped detecting pressure language, this scenario would still PASS. Add an assertion (e.g. `grep -qi "CRITICAL" "$f"` plus a check that it's flagged, not just quoted) so the row actually exercises the fixture it sets up. This is the AGENTS.md "changed code without covering test is a FAIL" rule applied to a test spec that doesn't cover part of what it claims to.
- **[low / confidence: low]** `skills/review-standards/SKILL.md:14` — the new `prompt-audit` bullet is appended under the "Language deep-dives (load only for the languages present in the diff)" heading, but it is not a language deep-dive and is explicitly *not* loaded that way — its own text says `/review --prompts` hands it to the reviewer and "nothing else loads it". The bullet's own wording resolves the ambiguity for a careful reader, so this doesn't cause a wrong load, but the heading's scoping rule no longer covers the list under it. Consider a second one-line heading (e.g. "Special-purpose references, loaded only when named:") rather than folding it into the language list.

## Vendoring check (as scoped)

- `skills/VENDORED.md`'s new row pins `github.com/anthropics/skills` at `8a1541c`, digest `c33200da…1878f2`. `sha256sum` of the working-tree file: `c33200da181e10f1f31bd4625059cd30bc6922940aac49017d8ab0c2721878f2` — matches the row's prefix/suffix exactly.
- `LICENSE.txt` beside it is the stock Apache-2.0 text (spot-checked header); not reviewed further per instructions.
- No network access in this sandbox, so the fetched-content-matches-upstream direction could not be independently re-verified against `raw.githubusercontent.com`; only the local digest-matches-file check was possible. Noting this as a scope limit, not a finding.
- Not reviewed: the prose of `prompt-audit.md` itself (per instructions) — noted only that it internally references `shared/model-migration.md` and `shared/prompt-caching.md`, two sibling files from the same upstream skill that are not vendored here. The `/review --prompts` dispatch brief in `skills/review/SKILL.md` only ever points the reviewer at `prompt-audit.md` and never asks it to open those two files, so nothing in this diff's own flow breaks on their absence — flagged here only as a fact for the operator, not a finding.

## Criteria

1. **Dead code** — pass. No removed callers, no orphaned exports; the new files are the sole additions and both are referenced (from `review-standards/SKILL.md` and `review/SKILL.md`).
2. **Flow tracing** — pass, with the one test-coverage gap above. Traced: `review/SKILL.md`'s new `--prompts` row (Step 3) → `review-standards/SKILL.md`'s new bullet → `references/prompt-audit/prompt-audit.md` → both vendor paths (`~/.claude/skills/review-standards/...`, `~/.gemini/config/skills/review-standards/...`), confirmed against `internal/harness/plan.go`'s `ClaudeRoot`/`AgyRoot` constants and a live scratch install. The `--fix` step's `--prompts` handling (apply hunks, run `TestRenderedContextFitsAntigravityBlock` only for `context/`/`rules/` edits) matches the project AGENTS.md's own scoping of that test. The new skill-matrix row's check under-asserts one of its four fixtures (finding above).
3. **Standards conformance** — pass. Frontmatter in `skills/review/SKILL.md` is double-quoted (contains a colon), no unquoted `#` (BUG-004). Grepped `README.md`, `skills/preflight/SKILL.md`, `context/AGENTS.md` for other restatements of the lens list (`--docs`, `--simplify`, etc.) — none found, so no other restatement needed updating for the new `--prompts` flag.
4. **Directory layout** — pass. `references/prompt-audit/{prompt-audit.md,LICENSE.txt}` sits under the existing `review-standards/references/` convention beside `go-review.md`/`typescript-review.md`; `skills/VENDORED.md`'s header rule ("never edit vendored content in place") is honored — the file is untouched from its pinned fetch.
5. **Simplicity (DIY-except-security)** — pass. No new dependency; this is a markdown reference file plus two doc/skill edits wiring it in. No security-relevant code touched.
6. **Decoupling** — pass. The new lens is additive to `review/SKILL.md`'s existing per-lens table/dispatch pattern; no existing lens's brief or table row was restructured, only appended to.

## Not examined

- Did not re-fetch `prompt-audit.md` from `raw.githubusercontent.com` to confirm the local copy matches upstream byte-for-byte (no network in this environment) — verified only that the local file's digest matches the row's recorded digest, and that the file is unedited relative to that pin (no diff-in-place indicators).
- Did not execute the `review-prompts` skill-matrix scenario itself (that is a live-CLI behavioural run per `AGENTS.md` § Testing, out of scope for a static review of the spec row).
- Did not review the prose content of `prompt-audit.md` beyond the passages needed to confirm it does not embed obviously wrong internal references — per the dispatch's explicit instruction not to review vendored prose.

## Outcome
- `skill-matrix.md:31` check asserts 3 of 4 seeds — fixed: the `review-prompts` row now greps `CRITICAL`, `reviewer-go` and `reviewer-ts` as well (seven greps); run against the four round-2 clones in `.agents/scratch/prompts-probe/runs/{4,5,6,ctl}`, PASS in each (`testing/2026-09-29-prompts-lens-behavioural.md`).
- `skills/review-standards/SKILL.md:14` bullet under the language-deep-dives heading — fixed: the bullet now sits under its own line, "Lens reference (not a language; loaded only when `/review --prompts` names it in the reviewer's brief)".
- Not examined, upstream byte check — done by the main agent: `sha256sum` of the vendored file equals the file fetched from `raw.githubusercontent.com/anthropics/skills/8a1541c4a3ffa5a20a5a91de0dcf3f0bab1d1ef4/skills/claude-api/shared/prompt-audit.md` (`c33200da…1878f2`, 48,983 bytes).
