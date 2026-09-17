---
owner: reviewer
date: 2026-09-17
---
# Review — worktree changes over 012b4ca, final pass   (2026-09-17)

Range reviewed: working tree vs `012b4ca`, limited to the files named by the dispatch:
`internal/harness/report.go`, `internal/harness/report_test.go`, `README.md`,
`context/AGENTS.md`, `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`,
`.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/sh-agent/guide.md`,
`.agents/docs/plans/2026-09-17-delivery-rework.md`.

A prior review of an earlier state of this range is
`.agents/docs/reviews/2026-09-17-review-012b4ca-worktree-code-docs.md`; the tree moved on
since it was written (D14, a third security round) and this pass re-derives from the
current bytes rather than trusting that report.

Blobs read (pinned via `git hash-object`, current worktree):
- `internal/harness/report.go` = `c97861a80070ace86273491039f95ffa749b5f6e`
- `internal/harness/report_test.go` = `5f8674a83467687e01d3f7fe76c43c522aed63a0`
- `README.md` = `2e282ef77a484bdb33360730dd7f2b32682bf2cb`
- `context/AGENTS.md` = `e1aae1c9899b9901c5d4fdf1b1fd952178557a51`
- `.agents/docs/deploy/local.md` = `d975776fb62186221ceb7e7343e6a8d61a454f8b`
- `.agents/docs/deploy/remote.md` = `85629dd209398a795a297b017ec2985e7adc329b`
- `.agents/docs/specs/sh-agent/delivery.md` = `07a2126ccb58fa8acd6d39e07b7b10927856926a`
- `.agents/docs/testing/sh-agent/guide.md` = `92fc106a33d353cd93e08e2159b9fa12433a76ff`
- `.agents/docs/plans/2026-09-17-delivery-rework.md` = `e4d06bd5f9eced5becabc2da6e52b61aec78102b`

Base commit blobs (at `012b4ca93838373bd1823672eee3a5a8fd584de2`) hashed separately for the
diff; not repeated here.

## Verdict: APPROVE

This is D14's fix for the third security round: it narrows the push-target claim
("cannot be redirected" → "cannot be redirected by origin's own config; the operator's own
git rewrites and named remotes still apply, trusted rather than refused"), closes the
header-mentions-a-footer PEM boundary bug, widens the `authorization` needle to accept a
quote before the colon, fixes the JSON-escape decoder to read every value on a line
(closing a `1e400`/two-records-on-one-line regression D13 had introduced), and adds a
must-fail test for `redactDocument`'s decoded branch and dry-run visibility for dropped
records. `gofmt -l`, `go vet ./...`, and `go test -race ./...` (both packages) are clean.
Every code-carrying fix in D14 has a mutant log in `.agents/scratch/r3-mutants/` showing
the specific assertion fail when the mechanism is reverted (`2b-label`, `3-authz`,
`4-decoded`, `6-document`, `7-dryrun`) — I did not re-run the mutants myself (budget), but
the logged failures name the exact test and assertion the fix protects, which is the
right shape of evidence. The docs (README, local.md, delivery.md, testing guide, the
plan's D9–D14 rows) were cross-checked against the code and against each other and agree,
including the specific overclaim the dispatch asked me to hunt for: the push-target text
in `report.go`, `local.md`, and `delivery.md` all now say the push is not redirectable by
`origin`'s config specifically, while a URL rewrite or an identically-named remote — the
operator's own configuration — still works, matching D14 item 1's resolution exactly. No
stray escape-sequence corruption (checked with `grep -P '[\x00-\x08\x0e-\x1f]'` across
every file in the range; `git diff --numstat` shows no binary flag).

I found nothing that blocks. Two low-confidence observations below are worth a look but
not fixes.

## Findings

- **[low / confidence: low]** `internal/harness/report.go:108-124` (`blockEnd`) — the
  header-footer label check (D14 item 2b) only applies `j == opened` — i.e. only guards a
  footer that appears on the SAME line the header opened on. A real footer with a
  *mismatched* label on a **later** line (`-----BEGIN RSA PRIVATE KEY-----\nbody\n-----END
  DSA PRIVATE KEY-----\n`) still closes the block there, unchecked, because the label
  comparison never runs once `j != opened`. This is consistent with what D14 item 2b and
  the docs literally scope the fix to ("a header line that merely mentions a footer"), so
  it is not a doc/code mismatch — the guarantee was never claimed for cross-line label
  mismatches, which barely occur in a real PEM anyway (a well-formed key's footer always
  names its own header's label). Flagging only because the fix's own doc comment
  ("on the line its header opened, a block closes only at a footer naming that header's
  label") could be misread as a general label-matching rule if read out of context; the
  code and the surrounding prose are consistent, so no action needed.

- **[low / confidence: low]** `internal/harness/report.go:78-83` (`labelAt`) —
  `s[strings.Index(s, "BEGIN"):]` panics if `Index` returns `-1` (slicing at `-1`). Both
  call sites pass a slice that starts exactly at a `blockOpen` match, which always
  contains the literal `"BEGIN"`, so this is unreachable today. Noting it because it's an
  implicit invariant enforced only by caller discipline, not by the function's own
  signature or a guard — a future caller passing an arbitrary string would panic instead
  of getting `""`. Not a finding against the current diff's correctness (no test can
  reach the panic), just worth the one-line guard (`if i := strings.Index(s, "BEGIN"); i
  >= 0 { s = s[i:] } else { return "" }`) if `labelAt` ever gets a second caller.

## Criteria

1. **Dead code** — pass. No leftover references to the struck `walkJSON`/`redactJSON`/
   `Redaction.Path` (confirmed by grep in the previous review and re-confirmed here: no
   hits). `labelAt`/`blockLabel`/`marker`/`decoded` are all called from the diff's own new
   code paths.
2. **Flow tracing** — pass. `redactLines`/`redactDocument`/`marker`/`decoded`/`labelAt` all
   agree on the `([]byte, []Redaction)` signature (no `error` return, correctly dropped —
   nothing in this path can fail); `Collect`, `SendReport`, `commitAndPush` all updated in
   step. `commitAndPush(dir, url, paths, msg, out)` receives `r.URL` from `SendReport`,
   and `verifyFeedback` still gates entry before any write, with the trust split now
   correctly documented (`verifyFeedback` guards the commit's location, `commitAndPush`
   guards the push's destination — report.go:418-421). Every new/changed test I sampled
   (`TestReportRedactShapes`'s new cases, `TestReportRedactStructuralGaps`'s new cases,
   `TestSendReportDryRunShowsDroppedRecords`) has a corresponding must-fail mutant log
   under `.agents/scratch/r3-mutants/` demonstrating the assertion fails when the
   mechanism is reverted.
3. **Standards conformance** — pass. Stdlib only, `gofmt -l` clean, `go vet ./...` clean.
   The testing guide's claim that `-run 'TestReport|TestSendReport|TestResolveFeedback'`
   plus `-list` covers every report test was independently re-verified: `go test
   ./internal/... -list 'TestReport|TestSendReport|TestResolveFeedback'` lists exactly the
   13 tests present in the file, including the two new ones
   (`TestSendReportPushesOnlyToFeedbackURL`, `TestSendReportDryRunShowsDroppedRecords`).
4. **Directory layout** — pass, no new files in the named range.
5. **Simplicity / DIY-except-security** — pass. Net simplification versus the D9 JSON walk
   (fewer branches, no `bytes`/`sort`/`strconv` imports); no crypto/auth surface touched.
6. **Decoupling** — pass. `marker`, `decoded`, `labelAt`, `blockEnd` remain small,
   independently testable functions; the push-target trust decision stays at the caller
   (`SendReport`/`commitAndPush`), not smeared into `verifyFeedback`.

Docs lens (`/review --docs`): `README.md`, `local.md`, `remote.md`, `delivery.md`, and the
plan's D13/D14 rows describe the current whole-line/whole-document/decoded-JSON redaction,
the URL-not-origin push with its narrowed claim, and the "redaction is a safety net, read
what changed before every report" framing consistently with the code and with each other.
D13 and D14 match the code precisely, including the specific items I checked against the
diff: D13 item 2 (`push -u <url> HEAD:main`, origin check stays for the commit) matches
`report.go:590` (`gitOut(dir, "push", "--quiet", "-u", url, "HEAD:main")`); D14 item 1 (push
claims narrowed, no code) matches the `FeedbackURL`/`commitAndPush` doc comments, which now
name `insteadOf`/`pushInsteadOf`/a same-named remote as trusted rather than refused; D14
item 2b (label-matching footer fix) matches `blockEnd`'s new `label`/`opened` tracking;
item 3 (authorization needle allows a quote before the colon) matches the widened
`authorization` regex; item 4 (decoded reads every value, numbers as written) matches
`decoded`'s loop with `dec.UseNumber()`; item 6 (must-fail test for `redactDocument`'s
decoded branch) matches the two `redactDocument` decoded-fallback assertions in
`TestReportRedactDocument`; item 7 (dry-run prints `lines`/`unterminated`) matches
`SendReport`'s output loop and `TestSendReportDryRunShowsDroppedRecords`. No description
claims more than the code does — the previous review's target overclaim ("nothing can
redirect it") is gone from every file in the range, replaced by the hedged version
everywhere I checked. No silently-converted escape sequences found in any file in the
range (checked with a control-byte grep and `git diff --numstat`, both clean).

## Not examined

- `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/bugs/INDEX.md`,
  `.agents/docs/reviews/2026-09-17-security-6af54f7-015cf2d.md`, `REQUIREMENTS.md`, and the
  eight new `BUG-053..060` files plus the four new/updated `2026-09-17-review-*`/
  `2026-09-17-security-*` files under `.agents/docs/reviews/`: all shown as changed/
  untracked in `git status` but not in the dispatch's named file list, so not diffed here.
  `REQUIREMENTS.md` in particular gained a new "R9 — Four procedures" section in this
  working tree; I skimmed it only far enough to confirm it references the same "Report
  back" procedure the reviewed code implements (consistent on that one point), not a full
  review of R9 itself or its approval trail.
- I did not re-execute the mutant tests in `.agents/scratch/r3-mutants/`; I read the five
  logs there and confirmed each names the test/assertion the corresponding fix protects,
  but did not independently rebuild the mutant binaries to re-confirm the fail-then-pass
  sequence myself.
- No live end-to-end run against the real `styrman-harness-feedback` GitHub repository —
  out of scope for a diff review; the project's own testing guide assigns that tier to the
  main agent, not a bare reviewer dispatch.
- `cmd/sh-agent` and other `internal/harness/*.go` files outside the named diff were not
  read beyond confirming `go build`/`go vet`/`go test -race ./...` pass for the whole
  module.

## Outcome

- Finding 1 (the label guard applies only on the header's own line): rebutted — scoped that way on purpose (D14 2b). The defect was a footer merely MENTIONED on the header's line closing the block above the body; a footer on a later line is a real footer, and requiring its label to match would let a typo or a rewritten label run redaction to the end of the file. `local.md` states the fix at that scope.
- Finding 2 (`labelAt` would panic on an index of -1): rebutted — both call sites slice from a `blockOpen` match, and every alternative of `blockOpen` contains `BEGIN`, so the index is never -1; a guard would be code for a state the regex cannot produce.
