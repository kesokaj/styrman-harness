---
owner: main
date: 2026-10-12
---
# Review — worktree (merged)   (2026-10-12)

Range reviewed: working tree on 66f3412 — the resume items of handoff 2026-10-09-1921 approved 2026-10-12: the dream signal fix (BUG-170), the `-u` redaction anchor, the BUG-167 record line, the `readSmall` cap test, the agent-key allow-list (ADR-056), the timer's control-character refusal (row 21), the unreachable search refusals removed (row 23), and the stub list's home moved to the init-project skill (row 28 L2, ADR-057).

## Verdicts

- code + docs (`reviewer`): APPROVE WITH NITS — `reviews/2026-10-12-review-worktree-code-docs.md`
- security (`security`): FINDINGS, all low — `reviews/2026-10-12-security-worktree.md`

No FAIL. No finding is medium or higher.

## Findings

| # | Severity / confidence | Lens | Location | Finding |
|---|---|---|---|---|
| 1 | low / high | security | `internal/harness/report.go:167` | The `-u`/`--user` needle widened one character at a time: the flag after `= ( [ , ; \|`, argv arrays, a tab or several spaces after the flag, and an escaped quote in plain text still passed; on the JSON-lines path the tab form and the flag after `=` or `(`. |
| 2 | low / high | docs | `.agents/docs/specs/sh-agent/dream.md:32` | The spec says the CLI's group is killed on SIGINT, SIGTERM, SIGHUP or SIGQUIT, with no word on an ignored interrupt or hangup. |
| 3 | low / medium | docs | `.agents/docs/specs/sh-agent/delivery.md:77` | The paragraph names the v2.6.2 refused keys; the release that ships ADR-056 must update it, and `TestDeliverySpecCurrent` does not pin that clause. |
| 4 | low / low | code | `internal/harness/derive.go:30` | The doc comment names `isolation` and `permissionMode` as keys `DeriveAgy` drops; `LoadAgents` now refuses both first. |
| 5 | low / low | security | `internal/harness/inventory.go:209-220` | Skills, commands and rules stay on a deny-list (`allowed-tools`, `hooks`). |

## Not checked

- The full `go test -race ./...` by either lens agent (targeted runs only); the verifier runs it after these fixes.
- A live Claude Code run that refuses an agent with `memory` or `initialPrompt` at install (the lint is unit-tested; the release gate's real-home smoke installs the tree).
- The operator approvals ADR-056 and ADR-057 cite (they are in the session of 2026-10-12, not in the tree).

## Outcome

- 1 — fixed at `internal/harness/report.go:167`: the operator chose the class over a fourth batch (2026-10-12, "Widen to the class"): the flag after any character but a letter, digit, `_` or `-`, and whitespace, `=`, an argv-array `", "` or an escaped quote before the value. Every shape the finding names is in `shapeSecrets` or `TestReportRedactFlagAtValueStart`, red before the change (6 failures on the JSON-lines test alone), with a near-miss per new boundary. Over-redaction measured on the trail: 15 lines at v2.6.2, 36 now, all in review records quoting credential shapes, none in a living doc. A shape the class misses is a known limit from here, not a patch (the rule's comment says so).
- 2 — fixed at `.agents/docs/specs/sh-agent/dream.md:32`: an interrupt or hangup sh-agent started with ignored stays ignored and leaves the CLI alone.
- 3 — left: the release commit (`/preflight` Step 6), which updates `delivery.md` § Where things stand under BUG-167's prevention rule; the gate report names the clause.
- 4 — fixed at `internal/harness/derive.go:30`: the comment names the allow-list instead of two refused keys.
- 5 — rebutted: a trade-off recorded in ADR-056 § Decision (the Agent Skills standard keeps the skill, command and rule key sets open: `license`, `metadata`, `compatibility`).

Round 2 skipped: every finding is low, every code fix has a test that fails without it, and the verifier runs after the fixes (review skill Step 4). The release gate dispatches `reviewer` and `security` over the whole release range next.
