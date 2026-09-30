---
owner: changelog
date: 2026-09-30
---
# Review project: seven findings fixed, one accepted risk     (2026-09-30)

## Summary
`/review --project --all` at 3478a1c found seven findings (one medium, five low, one deferred). Six were fixed: atomic write factored from two copies to one `WriteAtomic` helper; release-directory sweep switched from glob to `ReadDir` + prefix to stay in parent; sixteen redaction needles added (ten vendor prefixes, five length-and-alphabet patterns, a bare userinfo token); gke-basics and terraform vendored forks given source pins; push errors wrapped in credential scrubbing. Finding 3 (unsigned checkout tip) left to operator and recorded as ADR-046 (accepted risk, R9's contract). Three review rounds fixed a round-1 side effect (WriteAtomic order) in round 3. Tests added: TestSweepStaleDirsStaysInParent, TestUpdateFromReleaseSweepStaysInHome, TestPushErrorScrubsCredentials, TestVendoredRowsCarryAPin, TestSettingsWriteUsesHarnessWriteAtomic, TestAbandonNamesATempFileItCouldNotRemove, TestWriteAtomicClosesBeforeRename, rows in TestReportRedactShapes. Exports added: WriteAtomic at apply.go (replacing writeAtomic+attributionSet), SweepStaleDirs at sweep.go (replacing sweepStaleZipDirs and globOrNil). Removed: abandonTemp.

## Type
task

## Files changed
- internal/harness/apply.go — WriteAtomic exported, abandoned-file diagnostic restored
- internal/harness/sweep.go — SweepStaleDirs, replaces glob with ReadDir + prefix + age check
- internal/harness/release.go — calls SweepStaleDirs instead of inline sweep
- cmd/sh-agent/main.go — calls SweepStaleDirs instead of sweepStaleZipDirs; attributionSet calls WriteAtomic; abandonTemp removed
- internal/harness/report.go — deny-list expanded (see Summary); pushError classifies and scrubs a failed push
- internal/harness/adopt.go, internal/harness/state.go — callers renamed to WriteAtomic
- internal/harness/sweep_test.go — new
- skills/VENDORED.md — gke-basics pinned to google/skills 6f0b877; terraform pinned to antonbabenko/terraform-skill v1.6.0 2271bc4
- AGENTS.md — vendoring rule names the forks
- internal/harness/release_test.go, internal/harness/report_test.go, internal/harness/harness_test.go, cmd/sh-agent/main_test.go — tests for findings 1, 2, 4, 5, 6 above

## Agents dispatched
explorer (project recon), reviewer (code, simplify, docs, prompts lenses; three rounds on fixes), security (project, then three rounds on fixes), verifier (project, then three rounds in single log), documentation (redaction shapes docstring), changelog, bugs

## Verification
- `.agents/docs/reviews/2026-09-30-review-project-merged.md` — merged gate report, seven findings, three review rounds
- `.agents/docs/reviews/2026-09-30-review-project-code.md` — code lens APPROVE
- `.agents/docs/reviews/2026-09-30-review-project-simplify.md` — simplify lens APPROVE WITH NITS (medium: atomic write duplication)
- `.agents/docs/reviews/2026-09-30-review-project-docs.md` — docs lens APPROVE
- `.agents/docs/reviews/2026-09-30-review-project-prompts.md` — prompts lens clean
- `.agents/docs/reviews/2026-09-30-security-project.md` — security lens FINDINGS, all low
- `.agents/docs/testing/2026-09-30-verify-project.md` — verifier PASS (366 pass, 0 fail, 2 skip with reasons)
- `.agents/docs/plans/2026-09-30-recon-project.md` — explorer recon
- `.agents/docs/reviews/2026-09-30-review-worktree-code.md` — code review of the fixes, rounds 1–3 (overwritten each round; final: APPROVE)
- `.agents/docs/reviews/2026-09-30-security-worktree.md` — security check of the fixes, rounds 1–3 (final: three lows, put to the operator)
- `.agents/docs/testing/2026-09-30-verify-review-fixes.md` — verifier log of the fixes, rounds 1–3 appended (PASS each)

## Docs
`documentation` updated `.agents/docs/deploy/local.md:66` (the redaction shapes described by kind, with a pointer to `prefixNeedles`/`keywordNeedles` in report.go and to the known gaps in the merged record) and `.agents/docs/specs/sh-agent/delivery.md:56` (the needle kinds). The delivery model is untouched. Finding 7 (main.go layout) left deferred; finding 3 (unsigned checkout) becomes ADR-046.

## Requirements
n/a

## Operator choices
From Outcome section, two items left to operator (recommended answers given):
- Finding 4 remainder, after three rounds: an upper-case `HOOKS.SLACK.COM` host, a userinfo token of 16–19 characters, `hooks.slack.com/triggers/`, an `=`-padded token, an `sk-None-` key. Choice: (a) generalise to rules per kind — recommended, or (b) accept and record.
- Finding 5 remainder: whether gke-basics and terraform forks should be renamed to satisfy "differently named skill" rule. Choice: keep names — recommended, rule's purpose (no re-fetch overwrites local text) now stated in VENDORED.md.

## PATCH candidate
This change fixes defects only in released behaviour (the update sweep, the report redaction, the push error, the settings write's diagnostic); no new behaviour added. Candidate: v1.13.1.
