---
owner: main
date: 2026-10-09
---
# Preflight — v2.6.1   (2026-10-09)

Range reviewed: v2.6.0..5e5dc29

This is a PATCH release: the sandbox kits' credential scan no longer fails a run on AWS's published example key. No installed file changes except the zip's copy of the testing kit. No skill, agent, rule, context or Go source outside a test changed.

## Failures found during the gate (all fixed before the tag)

| # | Failure | Severity / confidence | Resolution |
|---|---|---|---|
| 1 | f6d9854 held the scan's matches in a shell variable: under `bash -x` a real key reached the kit log | medium / high | fixed in eaa029a (AWS keys counted through a pipe, never captured); pinned by "a key found under bash -x" |
| 2 | f6d9854: the scan's own status could trip the EXIT trap when the last file held only the example | low / high | fixed in eaa029a (`return 0`); pinned by the one-file case |
| 3 | f6d9854: a shell function named `grep` blinded the scan (reviewer) | low / medium | fixed in eaa029a (`command grep`), widened in 399c1c6 (`unset -f find head od tr grep`); pinned by three function cases |
| 4 | eaa029a: `|| true` let a grep read error pass a real AWS key, a regression introduced by the round-1 fix | low / high | fixed in 399c1c6 (exit 2 reports the file); pinned by "a grep that cannot read for -o". At the review cap the operator chose "Fix both, round 3" (this session, 2026-10-09) |
| 5 | 399c1c6's comment claimed protection against "a kit or an agent"; a kit can still redefine builtins | low / high | the comment was fixed in 5e5dc29 to state the real scope; the mechanism is rebutted (the agent under test runs in its own process), recorded in the security report's Round 3 Outcome |

Every fix has a test case that fails against the version before it. The verifier's copy control shows v2.6.0's and eaa029a's scripts failing the suite.

## Step 2: review, security, verification

- **Shadowing:** neither `.claude/agents/` nor `.agents/agents/` exists.
- **Reviewer:** `reviews/2026-10-09-review-v2.6.0-f6d9854-code-docs.md`, opus, code and docs, COMPLETE, approve with nits; Outcome written.
  - Its range ends at f6d9854. The commits after it are fixes from the review rounds. Under `review` Step 4, round 2 went back to `security`, the agent with the medium finding, which reviewed eaa029a and 399c1c6.
  - 5e5dc29 is a comment and a removed `local`, read by me.
- **Security:** `reviews/2026-10-09-security-v2.6.0-f6d9854.md`, COMPLETE, rounds 1–3 over v2.6.0..399c1c6.
  - Every release-diff path was read, and none is named under Not examined.
  - No high finding; the one medium is fixed (row 1).
  - Its Round 3 Outcome covers 5e5dc29.
- **Verification:**
  1. **Log reaching HEAD:** `testing/2026-10-09-verify-v2.6.1-release.md`, `owner: verifier`, range v2.6.0..5e5dc29, PASS.
     - gofmt clean, vet clean on host and the three targets, `bash -n` clean;
     - `go test -race -v -count=1 ./...` 646/0/3, with all 22 subtests of `TestKitPreambleRefusesLeftoverCredentials` and `TestScriptsHaveNoEarlyExitPipe` passing;
     - govulncheck: no reachable finding;
     - the must-fail control on a copy against the v2.6.0 and eaa029a scripts.

     Its earlier stopped runs (HEAD moved three times) are recorded in its `Previous run:` line.
  2. **What the range changed was run:**
     - `testing/guide.md` changed only the credential paragraph's wording; its commands are unchanged.
     - The kit itself is exercised by its behaviour test above.
     - No plan was added or modified, and no skill or agent changed, so no sandbox, behavioural or real-home tier is owed.
  3. **No failure is left standing:** the verifier's log had no FAIL.

  The verifier noted one deviation from my dispatch: two FAIL-expected guard cases also pass on v2.6.0's script, which is their purpose. The example-key-only pass cases were shown failing against v2.6.0's scan before f6d9854 (`changes/2026-10-09-1324-kit-aws-example-key.md`).

## Step 3: trail freshness

| Check | Result |
|---|---|
| Newest `changes/` entry describes the last work (`…-kit-scan-comment-scope.md`) | pass |
| No future `date:` | pass |
| Nothing deleted under `testing/` | pass |
| `changes/` and `DECISIONS.md` not modified | pass |
| ROADMAP and TODO | untouched since v2.6.0; Now and Next are empty, TODO has no items (report-only) |

## Step 4: docs freshness

The reviewer's docs lens was clean. `testing/guide.md` names the exception. No open docs gap.

## Accepted, with reasons

- A kit script that redefines a shell builtin (`printf`, `read`, `[`, `unset`) can blind the scan. The kit is the operator's own committed script, and the comment says so (row 5).

## What was not checked

- **macOS (bash 3.2, BSD grep):** the scan was not run there. The pipe-and-count form avoids the here-string that older bash writes to a temp file. `grep -o`, `-c`, `-v`, `-x` and `-F` are POSIX or long-standing BSD flags, but this was not checked on a Mac.
- **A live sandbox kit run** on a home holding the example key, after the fix. The behaviour test plants it in a transcript and a binary database instead.

Neither is a blocking condition: no skill or agent changed, so the project's Antigravity rule does not apply.

## Bridge check

v2.6.0 is not a bridge release. No bridge closes with this tag.

VERDICT: SHIP
