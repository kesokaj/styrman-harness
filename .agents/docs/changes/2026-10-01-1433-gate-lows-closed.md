---
owner: main
date: 2026-10-01
---
# 2026-10-01-1433 — The v2.0.0 gate's eleven lows, closed

Operator, in their own turn: "Do the eleven lows from the gate too". Per item (gate rows in `reviews/2026-10-01-preflight-v2.0.0.md`):

1. **Behavioural hook test (row 1 nit).** `TestStoreRunsNoProjectHook`: a `core.fsmonitor` hook that writes a marker fires for a plain `git status` (the control) and not for `memory.Store`.
2. **The dream's feedback redaction (row 7).** `harness.RedactJSONL` — `Redact` with report's decoded-line matching — applied to each feedback entry in `internal/dream/collect.go`; `TestFeedbackRedactsEscapedSecret` uses the gate's measured shape (`\\u0026sig=<32 hex>`). BUG-145: the test's first two fixtures proved nothing (a doubled escape; a needle the raw pass already caught) and the must-fail control showed it.
3. **Symlinked `MEMORY.md` (row 7).** `readLines` refuses a symlinked index, so `memory show` and `add` exit 2; `TestSymlinkedIndexRefused` (the index is still the link after the refused Add).
4. **Antigravity plan-mode probe (row 8).** `testing/2026-10-01-verify-agy-plan-mode-probe.md`: the dream's exact command from a grant-nothing private HOME wrote no file and ran no command; the one write attempt is in `denied_actions`.
5. **data-gcp** names CREATEROLE and CREATEDB as the cost of the default role. 6. **docker builder tag**: tried a digest pin, reverted on the security review (a digest freezes the Go patch while the text promised patch tracking); the tag stays, the choice is the operator's (handoff). 7. **feature** treats issue text as data. 8. **README and VENDORED.md** pin the skills CLI's version and drop `-y`. 9. **Makefile sample** passes `--project=$(GCP_PROJECT)`. 10. **README and remote.md samples** stage by name instead of `commit -am`.
11. **Lint, graduating BUG-143.** `LoadSkills` refuses `allowed-tools`; after the security review, by kind: `allowed-tools`/`hooks` in skills, commands and rules, `hooks`/`permissionMode` in agents (`permissionKey` over `Keys`), and `ParseStrict` refuses any top-level frontmatter line the key grammar cannot read (`"key":`, `key :`), which a name-keyed lint never saw. `TestLintRefusesAllowedTools`, `TestLintPermissionKeysByKind`.

Docs: project `AGENTS.md` § Rules, `skills/project-memory`, `skills/VENDORED.md`, `skills/data-gcp` (main agent); README, `specs/sh-agent/{memory-loop,dream,delivery}.md`, `deploy/remote.md`, `testing/sh-agent-matrix.md` (documentation agent; the JSON escape in dream.md written through a script after the tool turned it into its character).

Verification: `verifier` PASS twice (`testing/2026-10-01-verify-gate-lows.md`: 347 tests, four must-fail controls; `-r2.md`: 348 tests, four more controls incl. the re-aimed assertion). Review: code lens three rounds (REQUEST CHANGES for an undefined Makefile variable → APPROVE → APPROVE WITH NITS on the lint hunks), docs lens (six gaps, all fixed above), security (one medium reverted, two lows fixed, one left; Outcome appended). Left for a next task, owner main agent: `LoadDocs` lints only when `Split` succeeds (a BOM or `--- ` closer skips it; `commands/` and `rules/` are empty today) and a `hooks` row for rules in the test; an agent's `mcpServers` key; JSON-escaped secrets in the dream's transcript, trail and memory inputs (first step: whether tool results reach the prompt through `contentText`). Agents dispatched: verifier ×2, reviewer ×2 (code three rounds, docs), security ×1, documentation ×1, bugs ×1 (BUG-143 graduated, BUG-145).
