---
owner: main
date: 2026-09-30
---
# Everything completed: the queued verifier proposal, the last two home-path globs, the redaction residue     (2026-09-30)

## Summary
The operator: "I want everything to be completed. All of it." and "go all" (2026-09-30), after a status list that named three items as the only ones not done. All three are done here. (1) The harness proposal queued since 2026-09-23, "The verifier never hands back with a background job alive": `agents/verifier.md`'s "Long suites, once" bullet now runs a suite that takes minutes in the foreground with the call's timeout raised to the tool's maximum and its output redirected to a file (never a pipe, BUG-086), keeps the background path only for a suite past that maximum, and ends every background job (`kill`, then `kill -0` fails) before the hand-back; the wording the struck attempt lost two rounds on (`tee`, `jobs`) is gone. Behavioural run `testing/2026-09-30-verifier-long-suites-behavioural.md`: new definition 0 of 3 background calls, control 2 of 3. (2) The last two `filepath.Glob` calls on a path built from the home — `LastRetro` (`internal/harness/state.go`) and `configView` (`internal/search/reconcile.go`) — list with `os.ReadDir` and a suffix match; `TestLastRetro` gained a home with `[`, `TestConfigViewListsProposalsByName` is new; no `filepath.Glob` remains in non-test source. (3) The redaction families recorded as accepted residue after the v1.14.0 gate each have one rule in `internal/harness/report.go` (jwt short, google refresh token, discord webhook, teams webhook, sendgrid key, azure account key, cookie header, vault token; prefixes for age and DigitalOcean), widened once for the near shapes the security check of this change named; 24 rows and 11 benign near-misses in `TestReportRedactShapes`; the comment at the end of `keywordNeedles` and `deploy/local.md:66` say a shape past these rules is a design question for the operator, never another batch — the list is closed for instance work. The redaction decision reverses the "accepted residue" line of the v1.13.1 handoff on the operator's words above; the reviewer's finding on that point is rebutted with those words in its report.

## Type
task (with the proposal's apply)

## Files changed
- agents/verifier.md — the long-suites bullet
- internal/harness/state.go, internal/search/reconcile.go — ReadDir instead of Glob
- internal/harness/report.go — ten rules, four prefixes, the closing comment
- internal/harness/harness_test.go, internal/harness/report_test.go, internal/search/global_test.go — tests
- .agents/docs/deploy/local.md — the coverage sentence, edited by the main agent (one sentence)

## Agents dispatched
reviewer (`reviews/2026-09-30-review-worktree-completion-code.md`, REQUEST CHANGES on the decision to cover the residue — rebutted with the operator's words in its Outcome — and one low fixed: two benign rows), security (`reviews/2026-09-30-security-worktree-completion.md`, two lows, both fixed: rules widened, guide reworded; Outcome), verifier (`testing/2026-09-30-verify-completion.md`, PASS: suite, scratch install, govulncheck, the search tier — rebuild 1.85 s, warm search 36 ms — and mutations a, b, c). No documentation dispatch: the one guide sentence was the main agent's edit.

## Verification
`testing/2026-09-30-verify-completion.md`; `testing/2026-09-30-verifier-long-suites-behavioural.md`. Owed at the release: the real-home smoke on both CLIs and one Antigravity run of the changed `verifier` (`AGENTS.md` § Release).

## Docs
`deploy/local.md:66`. The delivery spec's needle kinds already cover these by class.

## Requirements
n/a
