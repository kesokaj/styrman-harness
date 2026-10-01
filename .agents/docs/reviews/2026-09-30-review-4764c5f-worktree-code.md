---
owner: reviewer
date: 2026-09-30
---
# Review — 4764c5f-worktree-code   (2026-09-30)

Range reviewed: working tree against commit 4764c5f (`git diff 4764c5f`; 16 files under cmd/sh-agent, internal/dream, internal/harness, internal/search)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Nothing is broken at runtime: gofmt, go vet and go test (cmd, internal) pass. The one blocking item is a stale description of a changed behaviour: `.agents/docs/deploy/remote.md:75` still says `update --check` "works on a dirty checkout", which this diff reverses (IN4); the operator runs that guide. The one security-relevant item is that the redaction relaxation in report.go lets short secret values through to the feedback repository; AGENTS.md / the list's own comment call any shape change there a design question for the operator.

## Findings
- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:75` — says `sh-agent update --check  # only reports; exit 2 when behind; works on a dirty checkout`. The diff (`cmd/sh-agent/main.go` update(), test "--check on a dirty tree refuses like the real run") makes `--check` on a dirty checkout return an error (exit 1) before the exit-2 path. The same file (line 92) recommends running `update --check` on the timer; on a dirty checkout that now errors, and no restatement says so. Fix: change line 75 to "refuses a dirty checkout like the real run (exit 1); exit 2 when behind", and say the same for the line-92 timer advice. README:183 and delivery.md:30 do not contradict it.
- **[medium / confidence: medium]** `internal/harness/report.go` keywordNeedles "secret assignment" (accept func) — redaction of `password|secret|token|api_key` assignments now needs a value of 8+ chars, or 6+ with a digit. Previously any value was redacted. Now `password=correct`, `token=abcdefg`, `password: hunter`, and `secret = my long passphrase` (value = `my`) all pass into the report sent off the machine. The rule comment on the list ("a shape past them is a design question for the operator, never an eighth batch: the list is a safety net") and AGENTS.md ("over-redaction in the safe direction is pinned there, never argued about") say this direction needs the operator's choice, not an audit-driven loosening. Tests added cover only the near-misses going through, not a pinned "these short values are now deliberately sent" row. Fix: put the loosening to the operator as a one-line choice; if kept, pin the short-secret cases in TestReportRedactShapes as accepted under-redaction.
- **[low / confidence: low]** `internal/search/dream.go` walk() — a length-delimited field is now tried as a nested message before as text. A prose leaf whose bytes happen to parse as a valid protobuf message (first byte a letter with wire type 0/1/5, e.g. `Y`+8 bytes, or chains of such) returns true with no string leaves and the prose is silently dropped; before, textField ran first and kept it. Probability is small for 20+ char prose with 3+ spaces (most letters give wire types 3, 4, 6, 7 and fail), but nothing measures the loss: TestProseLeaves prose samples both start with `T` (wire type 4, never parses), and TestAgyCorpusLeaves counts only leaked keys, never dropped prose, and is skipped without an env var. Fix: have the corpus measurement also compare leaf counts before/after, or fall back to text when the nested parse yields no leaves.
- **[low / confidence: low]** `internal/harness/apply.go` saveOnError — the union keeps every old-manifest entry that still exists by `Lstat`, which is true for a path the operator has since turned into a symlink (the main loop skipped it, IN3). `Uninstall` (apply.go ~line 322) `Lstat`s then `os.Remove`s with no symlink check, so after a mid-install failure an uninstall removes the operator's link, the thing IN3 set out to prevent. Fix: skip symlinks in the union, same as the resync loop. Not covered by TestInstallMidFailureKeepsManifestComplete (no link in its fixture).

## Hunk-by-hunk test check (fails without the change?)
- openCopy/copyFile var/copyRetryWait: TestCopyDBRetriesOnce (first-copy header torn, proves NOTADB in a subtest, asserts 2 copies) and TestCopyDBNeverOpensSource; callers ListProjects, AgyConversationText, conversationText all moved to openCopy, none left using copyDB+sql.Open (git grep). OK.
- walk order: TestProseLeaves new rows (id_rsa, hunter2) fail without the swap; see the low finding.
- search.go `have` keyed by chunk id: TestNoDuplicateHitAcrossRounds. Hits are built only in find(); ids are unique inside one index DB. OK.
- answer() returns (string, error): only caller is realRun; sc.Err() is checked only when no result event ended the loop early, as intended. Test in run_test.go added for it (not read in detail).
- removeCopies(conf, now): callers realRun and run_test (2 calls with explicit clocks). OK.
- contentText: TestSessionTextSkipsCommandOutput. A user line that legitimately begins `<command-` is dropped; accepted.
- trailDates: TestTrailDateFromGit. Paths are repo-root-relative from git but `filepath.Rel(p.Root, f)` is project-root-relative: a project that is a subdirectory of a repo, or a path git quotes (non-ASCII), falls back to mtime (the old behaviour, not a regression). git runs unbounded over the history of docs each night per project; no timeout. Not a failure now.
- prompt.go fence escape: replace is non-overlapping but cannot rebuild `>>>`; the label in `<<<data kind: label` is not escaped (file names only). Not tested in the parts I read.
- render indent: proposal text is now indented four spaces; skills/retro/SKILL.md:46 tells /retro to treat the text as data and re-derive it, so no reader depends on column 0. Drafts are written separately (not read in detail).
- Install symlink keep, resync skip, Status "foreign", Diff skip: TestInstallDoesNotFollowLinkedDest, TestResyncSkipsLink. Status/Diff tests for links not confirmed by name (only the two above and TestInstallMidFailure... seen in the test diff headings).
- update(): tests for dirty+behind, clean+behind exit 2, dirty up-to-date all added and changed assertions trace to the source change. Only doc restatement left stale is remote.md:75.
- IN12/IN13 (output to `out` and exit 2): memory_test.go TestDreamVerb rows.

## Criteria
1. Dead code — pass: `copyDB` is still used by openCopy and by the new test; nothing orphaned.
2. Flow tracing — pass with the two low findings above (walk fallthrough; saveOnError union vs Uninstall).
3. Standards conformance — fail: remote.md:75 restates the old `--check` behaviour; report.go loosening conflicts with the stated operator-decision rule for that list.
4. Directory layout — pass: no new files; tests sit beside code.
5. Simplicity (DIY-except-security) — pass: stdlib only; no new dependency; the `accept` func field is small. Redaction is not auth/crypto.
6. Decoupling — pass: `copyFile` as a package var is a test seam within search; no new cross-package reach-ins.

## Not examined
- Full bodies of run_test.go, dream_test.go (beyond the openCopy and walk tests), harness_test.go (beyond the test names and headers), search/store.go beyond isUnreadable.
- `-race` run and govulncheck; only `go test` without `-race`, gofmt and go vet were run.
- Behavioural / sandbox runs from AGENTS.md § Testing; the living-doc grep for `>>>`, `trailDates` and render indent restatements in specs/sh-agent/dream.md.
- The plan, the audit and any chat (out of scope by instruction).

## Outcome

Answered in the next round (see `2026-09-30-review-4764c5f-worktree-code-r2.md`); the fixes and their tests are in the phase commit. Every finding of this round was fixed, except the deploy-guide sentence, which the documentation pass of the phase fixes, and the nested-first prose-loss concern, rebutted by the corpus measurement (971 leaves and 0 argument keys before and after, `testing/2026-09-30-verify-audit-fix-phase-1.md`).
