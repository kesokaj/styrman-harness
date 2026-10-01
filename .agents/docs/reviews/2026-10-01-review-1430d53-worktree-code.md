---
owner: reviewer
date: 2026-10-01
---
# Review — 1430d53-worktree-code   (2026-10-01)

Range reviewed: working tree against 1430d53, `git diff 1430d53 -- internal/harness/derive.go internal/harness/inventory.go internal/harness/frontmatter.go internal/harness/harness_test.go`

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing blocking. The new checks work on the real tree (the four named tests pass; the longest real description is under 800 characters). Worst problem: `Lint` measures a block-scalar description 2 to 3 characters longer than `Parse` and the Agent Skills limit see it, so the limit is enforced slightly too tight for `|`/`>` blocks. Five lows in total, none a defect that fires on today's files.

## Findings
- **[low / confidence: high]** `internal/harness/frontmatter.go:111` — `joined := TrimSpace(unquote(e.val))` starts from `e.val`, which for a block scalar is the indicator itself (`|`, `>-`). The loop then yields `"| first second"`, so the count includes the indicator and a space (2 to 3 runes), while `Parse` (`strings.Join(cont, " ")`) and the spec measure only the text. A block description of 1,023 to 1,022+ characters is refused wrongly. Fix: start `joined` empty when `isBlockIndicator(e.val)`. The test margin (1,100 vs 1,024) hides it; add a case at 1,023 as a block.
- **[low / confidence: high]** `internal/harness/frontmatter.go:99-108` — the `descriptionLimit` const was inserted between `Lint`'s doc comment and `func Lint`, with no blank line, so the "Lint rejects ..." comment now documents the const and `Lint` has no doc comment; neither text mentions the new description rule. Describing-text rule: move the const above the comment and add the limit to Lint's comment.
- **[low / confidence: high]** `internal/harness/inventory.go:250` — `strings.Contains(fm["description"], "\n")` is dead: `Parse` joins continuation lines with spaces and `unquote` does not decode `\n`, so `fm["description"]` can never hold a newline. `descriptionIsMultiLine` carries the whole check. Delete the first operand. (`descriptionIsMultiLine` itself is correct: it uses the same `entries()` as `Lint`, so a `|` block, a `>` block and a wrapped plain or quoted scalar all show `val` as indicator or `len(cont)>0`; blank continuation lines are skipped by `entries` in both.)
- **[low / confidence: high]** `internal/harness/inventory.go:289` (and `LoadSkills` 189, `LoadAgents` 234) — `Lint` is called for skills (yes: the claim holds), agents, AND commands/rules, so the 1,024 limit also refuses a command or rule description. AGENTS.md:59 says "any skill or agent description". Harmless today, but the restatement is narrower than the behaviour; say "any description" or limit the check to skills and agents.
- **[low / confidence: medium]** `internal/harness/harness_test.go` `TestAgentDescriptionRules` — no control case: the fixture agent `x` is never loaded with a valid description, so a refusal from another cause (name, tools) would pass the test; error text is not asserted either. Add one accepted `x` agent or match the error substring.

## Per-hunk tests (fails without the hunk?)
- `inventory.go` multi-line check: TestAgentDescriptionRules "two lines" (block joins with spaces, only `descriptionIsMultiLine` catches it); banned-phrase loop: the other two cases. Yes.
- `frontmatter.go` Lint limit: TestLintDescriptionLimit (quoted 1,029, block ~1,101 refused; 1,019 accepted). Yes, no boundary case at 1,024/1,025.
- `derive.go`: comment only, no code change; the claim about the probe (3/3, 2/2) is not verifiable from the diff and is cited to an untracked log. TestDeriveAgy's added `maxTurns:` assertion pins existing behaviour (DeriveAgy builds from a whitelist) and would pass without a code change; TestDeriveAgyRealReviewer likewise pins existing rewrite on the real file. Both can fail if the behaviour regressed; they do not need a source hunk.
- Rune count: `len([]rune(joined))` counts characters, matching the spec's "characters"; invalid UTF-8 bytes count one each (conservative).

## Criteria
1. Dead code — one dead operand (inventory.go:250), see findings.
2. Flow tracing — pass: Lint reaches skills, agents, commands/rules; `descriptionIsMultiLine` and `Lint` share `entries`; errors wrapped with the agent/skill name. Off-by-indicator in Lint (finding 1).
3. Standards conformance — nit: Lint's doc comment displaced; AGENTS.md:59 restatement narrower than behaviour.
4. Directory layout — pass: code and tests beside each other in internal/harness.
5. Simplicity — pass: stdlib only, no new dependency.
6. Decoupling — pass: no new cross-package reach.

## Not examined
Everything outside the four Go files (agents/*.md, skills, AGENTS.md beyond the one restated line, the untracked verify log, the probe evidence quoted in the derive.go comment). Full `go test ./...`, gofmt and vet not run (only the four named tests, which pass). No mutation runs.

## Outcome

All five lows fixed in the phase: the `|`/`>` indicator is no longer counted (an under-the-limit block case added to TestLintDescriptionLimit); `descriptionLimit` sits above `Lint`'s doc comment; the dead newline check is gone; AGENTS.md § Rules says the limit covers skills, agents, commands and rules; TestAgentDescriptionRules has a valid-fixture control. The derive hunk is comment-only by design (D5b: the key is not emitted; the probe is the evidence).
