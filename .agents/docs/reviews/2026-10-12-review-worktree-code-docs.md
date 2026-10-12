---
owner: reviewer
date: 2026-10-12
---
# Review — worktree-code-docs   (2026-10-12)

Range reviewed: working tree — changed files: .agents/docs/architecture/DECISIONS.md .agents/docs/bugs/BUG-167-delivery-spec-drifted-two-releases.md AGENTS.md README.md REQUIREMENTS.md ROADMAP.md cmd/sh-agent/main_test.go context/AGENTS.md internal/dream/run.go internal/dream/run_test.go internal/dream/timer.go internal/dream/timer_test.go internal/harness/frontmatter.go internal/harness/harness_test.go internal/harness/inventory.go internal/harness/report.go internal/harness/report_test.go internal/metrics/events_test.go internal/search/search.go skills/init-project/SKILL.md; untracked, not ignored: none of source (untracked are review records only: 2026-10-09-security-v2.6.2-66f3412.md, 2026-10-12-security-worktree.md, this file). Read with `git diff HEAD`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing in the code is broken. One docs statement no longer matches the code: `specs/sh-agent/dream.md:32` still says the CLI's group is killed "on SIGINT, SIGTERM, SIGHUP or SIGQUIT", but an interrupt or hangup the parent had ignored (nohup, a background job) is now left alone (`internal/dream/run.go`). Everything else checked out: the agent frontmatter allow-list, the timer control-character refusal, the redaction flag-prefix fix, the signal fix and the removed `search.validate` branches are consistent with their callers, and their tests pass and can fail.

## Findings
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/dream.md:32` — stale. It says the process group is killed "on SIGINT, SIGTERM, SIGHUP or SIGQUIT (then sh-agent dies of the signal as before)". After this change (`internal/dream/run.go`, `callWithin`) SIGINT and SIGHUP are watched only when the parent did not ignore them at start; under `nohup sh-agent dream --all &` they neither kill the group nor sh-agent. Fix: add "unless sh-agent started with it ignored (nohup, a background job), which stays ignored" to that sentence.
- **[low / confidence: medium]** `.agents/docs/bugs/BUG-167-delivery-spec-drifted-two-releases.md:21` — the edit is a record correction, accurate on its facts (`git show 3bffdcf` added `mcpServers` to delivery.md twice; `6cec23b` removed it). Not a defect. One consequence to keep in mind: `delivery.md:77` says the installer refuses "`allowed-tools`, `hooks`, `permissionMode`, per kind", which is true of the tag v2.6.2 and, per its own rule, must not name the unreleased allow-list; the release that ships ADR-056 has to update that sentence. `TestDeliverySpecCurrent` does not pin this clause, so nothing will force it. Listed so the next release remembers; no change now.
- **[low / confidence: low]** `internal/harness/derive.go:30` — the DeriveAgy doc comment lists `isolation` and `permissionMode` as Claude-only keys it drops; after ADR-056 `LoadAgents` refuses both before `DeriveAgy` can see them. The function's behaviour is unchanged and the comment is not false, only no longer a case that can occur. Reword or leave; not blocking.

## Unverified
- ADR-056's "The operator chose the allow-list on 2026-10-12" and ADR-057's "The operator approved the exact text on 2026-10-12" (REQUIREMENTS.md R2 is immutable without that approval). No chat record is in this context; the ADRs, `REQUIREMENTS.md` and `ROADMAP.md` are in the same uncommitted change set, which is the order AGENTS.md asks for. The approval itself I could not check.
- No `changes/` entry for this work exists yet (newest is `2026-10-09-1921-review-project-fixes.md`); expected for an uncommitted tree, flagged so the entry is not forgotten.

## Evidence
- Allow-list: every `agents/*.md` sets only keys in `agentKeys` (awk over all 8 files: name description tools model [skills] effort maxTurns). `TestLintPermissionKeysByKind` and `TestLintSeesEveryKeyAYAMLReaderSees` pass; they refuse `memory`, `initialPrompt`, quoted `"memory"`, `isolation`, `futureKey`, and load the nine listed keys.
- Signal fix: `TestCallLeavesTheGroupOnAnIgnoredSignal` passes on the tree. Must-fail: in a scratch copy (`.agents/scratch/mut`, removed afterwards) I replaced the `watched` list with the old unconditional four signals; the file compiled and both subtests failed ("the ignored hangup killed the CLI's group", "sh-agent exited on the ignored hangup"). Go keeps only SIGHUP and SIGINT ignored at start, which matches the two the code tests; a signal ignored and never `Notify`d stays ignored across the repeated `callWithin` calls of a `--all` run.
- Timer: `bin` is `filepath.Join(home, BinPath)`, so the home is covered through `bin`; `TestInstallTimerRefusesEveryControlCharacter` and `TestInstallTimerRefusesUnsafePaths` pass. No living doc or test quotes the old refusal text (grep).
- Redaction: `TestReportRedactShapes` and `TestReportRedactFlagAtValueStart` pass. The new near-miss `` see `-ui:8080` above `` sits in the accepted-over-redaction list (asserts redacted), not the benign list, which matches the regex and the test comment's intent; the benign additions (`"-u" flag`, `{"flag":"-u"}`, `opts:-u root`) do not match.
- Search: `index` has no `positional`/`--refs`/`--open`/`--kind`/`--rebuild`-on-search flags in the `verbs` table, so the removed `validate` branches were unreachable; `runSearch` is the only `search.Options` constructor (grep). `./internal/search`, `TestVerbRefusesAnotherVerbsFlag`, `TestReadSmallStopsAtOneMiB` pass. The search spec's exit-2 sentence still holds.
- R2 pointer move: `grep` for `R2`, `one home of that list`, `2.1.281` over context, skills, agents, rules, README, specs, testing guide and deploy found the list only in `skills/init-project/SKILL.md:61` and REQUIREMENTS.md R2 (which now points at the skill). `context/AGENTS.md` pointer updated. No go test pins the old wording.
- Agent-key restatements: README.md:233 and AGENTS.md:63 both carry the allow-list and ADR-056; no other living doc lists the refused agent keys except `delivery.md:77` (see finding 2) and the derive.go comment (finding 3).
- Targeted tests run: dream (3 groups), harness (lint, redaction), metrics, cmd/sh-agent verb refusal, search. All pass. The full suite was running elsewhere and was not started by me.

## Criteria
1. Dead code — pass. The removed `validate` branches were unreachable and are gone with a test pinning the CLI refusal; no new unused symbols (`unlistedKey` and `agentKeys` are used; `permissionKey` is still used by the skill/command/rule path).
2. Flow tracing — pass. Signal, timer, lint, redaction and verb-table paths traced to their callers; new tests fail without the fix (signal test mutated; the others are pure additions to refusals that the old code accepted).
3. Standards conformance — pass with one stale restatement (dream.md:32) and one comment (derive.go:30). Restatement grep run for agent keys, the R2 list and the timer refusal text.
4. Directory layout — pass. No new files in source; tests sit beside the code.
5. Simplicity (DIY-except-security) — pass. Stdlib only (`slices`, `unicode`, `os/signal`); no new dependency; no auth/crypto code touched.
6. Decoupling — pass. `agentKeys` lives with the loader in `internal/harness`; the search change shrinks the package's duplicate validation rather than adding a reach-in.

## Docs lens
1. dream.md:32 — stale (finding 1).
2. delivery.md:77 — correct for the tag it describes; the release that ships the allow-list must update it (finding 2).
3. README.md:233, AGENTS.md:63, ADR-056, ADR-057, REQUIREMENTS R2, context/AGENTS.md:64, init-project SKILL.md:61 — agree with each other and with the code.
4. testing/guide.md — no statement about agent-key lint to update; § 1 commands unchanged.
5. deploy/*.md — nothing changed that they describe.
6. Unverified: operator approvals (above); the `--docs` check of the `TestHarnessMap` map was not needed (no install/update/report/release path changed).

## Not examined
- Full `go test -race ./...` and `govulncheck` (the caller runs the suite; not started here).
- A real-home or Antigravity run of `init-project` after the skill-text change (the release rule in AGENTS.md § Release applies at `/preflight`, not here).
- `skills/init-project` behaviour under a behavioural run; I read the text change only.
- Security depth on the redaction regex beyond the cases in the tests; `2026-10-12-security-worktree.md` is the security lane's record and was not read.
