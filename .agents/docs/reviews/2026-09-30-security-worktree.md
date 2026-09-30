---
owner: security
date: 2026-09-30
---
# Security review — worktree, round 3   (2026-09-30)

Range reviewed: Scope: working tree — changed files: AGENTS.md cmd/sh-agent/main.go cmd/sh-agent/main_test.go internal/harness/adopt.go internal/harness/apply.go internal/harness/harness_test.go internal/harness/release.go internal/harness/release_test.go internal/harness/report.go internal/harness/report_test.go internal/harness/state.go skills/VENDORED.md; untracked, not ignored: internal/harness/sweep.go internal/harness/sweep_test.go; read the diff with git diff HEAD -- <files> and read each untracked file whole.

## Status: COMPLETE

## Verdict: FINDINGS

Three findings, all low.
- **The fixes work.** All four round-2 fixes are correct, and none introduced a new defect. Three are pinned by a mutation that compiled and failed its named test.
- **The worst problem is the redactor's deny-list, again.** Two items from round 2's first finding were neither fixed nor recorded as `left: operator`: an upper-case Slack host, and a bare URL userinfo token of 16-19 characters. A probe on the live code also shows three more shapes from the same classes still pass. This is the fourth pass over this list (BUG-053, the project review, round 2, now round 3), and each pass finds more. By the global rounds rule, stop and ask the operator rather than add needles a fourth time.
- **Untested fix.** The close-before-rename order in `WriteAtomic` is right, but no test would catch it being reverted.
- **Stale guide.** The operator's guide still lists the needle set from before this range.

## Findings
- **[low / confidence: high]** `internal/harness/report.go:46,49,74,80` — what is left of round 2's finding 1, plus shapes from the same classes. **Probe:** `.agents/scratch/sec-worktree-r3/mut/internal/harness/zz_probe_test.go` calls `redactLines` on a copy equal to the live code, in text mode and in JSONL mode. Its positive control was redacted and its negative control (`ssh://git@…`) was not. Each of these came back **unredacted in both modes**:
  - `https://HOOKS.SLACK.COM/services/…`. Prefix needles match case-sensitively (`needleIn`, `strings.Contains`), and the host is not case-sensitive. Carried from round 2.
  - A 17-character bare userinfo, `https://FAKEVALUEFAKEVALU@host`. `url token` needs 20 or more characters. Carried from round 2.
  - A bare userinfo carrying `=` padding, `https://…FAKE==@host`. `=` and `+` are outside the `url token` class, and both appear in base64 tokens. New this round.
  - `https://hooks.slack.com/triggers/…`, Slack's workflow-trigger webhook path. New this round; my confidence that this path is live is medium.
  - An OpenAI `sk-None-…` key, the prefix of user keys created without a project. The `openai key` regex needs 32 alphanumerics straight after `sk-`. New this round; my confidence in the vendor format is medium.

  Each would reach the feedback repository verbatim. **Remediation:** put it to the operator as one choice.
  - **(a) Recommended: close the classes, not the instances.** This leaves one needle fewer than today, each with a `TestReportRedactShapes` case.
    - Replace both Slack prefixes with one keyword needle, `(?i)hooks\.slack\.com/`. It covers every path and every letter case.
    - Change `url token` to `://[^/\s:@]{16,}@`. That is any userinfo alphabet and 16 or more characters; `ssh://git@` still passes.
    - Add `sk-None-` as a prefix needle.
  - **(b) Record `left: operator — deny-list residue accepted`.** `.agents/docs/deploy/local.md:67` already tells the operator "Anything that matches no shape IS sent".
- **[low / confidence: high]** `internal/harness/apply.go:331-342` — the close-before-rename fix is correct by reading, but no test pins it.
  - **Mutation O1** moved the `Rename` block back above the `Close` block, as round 2 had it. It compiled, and no test in `internal/harness` or `cmd/sh-agent` failed.
  - The one failure was `TestRelativeHomeResolvesAgainstTheInvokingDirectory`, which fails the same way in the unmutated copy (`git clone` of a copy that is not a checkout). It is environmental and passes on the live tree.
  - Under the global Definition of Done, a finding counts as fixed only when a test fails without the fix. Pinning this order needs a close seam in `WriteAtomic`, which is more code than a window of one syscall warrants.

  **Remediation:** record `left: operator — order pinned by reading; no seam`. Separately, a comment nit: `abandon`'s inline comment (`apply.go:350`) says the file is "already closed on the rename branch", but after this fix it is also closed on the close-failure branch.
- **[low / confidence: high]** `.agents/docs/deploy/local.md:66` — the guide's "The shapes:" list still describes the needle set from before this range.
  - **What it lists:** vendor prefixes only up to `xoxs-`, plus the four keyword shapes.
  - **What it omits:** every shape this range added:
    - the prefixes `sk_live_`, `rk_live_`, `whsec_`, `sk_test_`, `rk_test_`, `sk-svcacct-` and `sk-admin-`;
    - both PyPI prefixes;
    - the Slack `services/` and `workflows/` hosts;
    - the `jwt`, `openai key`, `npm token`, `hugging face token` and `url token` regexes.
  - **Impact:** it understates coverage, which is the safe direction, so nothing leaks. It is still the operator's reference for what `report` catches, and it is an enumeration of a vocabulary this change grew without updating (the global "every enumeration of it" rule).

  **Remediation:** replace the list with a pointer to `prefixNeedles` and `keywordNeedles` in `report.go`, which is shorter and cannot drift; or update it in the same change. Route it to `documentation`.

## Fix by fix
- **Prefixes `sk_test_`, `rk_test_`, `sk-svcacct-`, `sk-admin-` and `hooks.slack.com/workflows/` (`report.go:49`): closed.**
  - Mutations P1-P5 each deleted one of the five needles. Each compiled, and each time `TestReportRedactShapes` failed.
  - A side effect, not a finding: the new prefixes also match ordinary identifiers, and the probe drops the whole line for `see task-admin-panel`, `network_test_helpers.py` and `framework_test_utils.py`. That is the safe direction.
    - Unlike `ASIAN food`, no such case is pinned in the benign list.
    - The comment at `report.go:36-38` ("a bare `sk-` would match every "task-"") implies that `task-` words survive, and `sk-admin-` breaks that.
    - If this is accepted, pin one case among the benign lines so that a change to it is deliberate.
- **`%` in `url token` (`report.go:80`): closed.** Mutation P6 removed `%` from the class. The code compiled and `TestReportRedactShapes` failed. The change can only add matches, so it cannot cause a miss.
- **`WriteAtomic` closes before the rename (`apply.go:331-342`): correct, untested** (the second finding).
  - **Every error path leaves `dest` alone.** Write, chmod, sync, close and rename all go through `abandon`.
  - **The second `Close` is harmless.** In `abandon`, `Close` on a file that is already closed returns `os.ErrClosed`, which is discarded. `tmp.Name()` is still valid after the close.
  - **Every caller reports failures truthfully again.** The five callers are `Install` (`apply.go:79`), `UpdateFromRelease` (`release.go:170`), `SaveState` (`state.go:68`), `copyTree` (`adopt.go:175`) and `attributionSet` (`main.go:907`). With this order, their "failed" messages are true again: round 2's false failures are gone.
  - **The ETXTBSY window is closed.** The inode is closed before it is renamed into place.
- **A young `release-*` survives (`release_test.go:654`): closed.** Mutation A1 set the age at `release.go:147` to `time.Duration(0)`. The code compiled and `TestUpdateFromReleaseSweepStaysInHome` failed. The test now pins the threshold between the test's own duration and 2 h: a threshold of 3 h would leave `release-old` in place and fail too.
- **`VENDORED.md` split date: closed.** The blob-less upstream clone (`.agents/scratch/sec-worktree/gskills.git`) gives these commits:
  - `5d9d054`, dated 2026-06-23 (author and committer), "Refactor GKE basic skill into individual skills".
  - The first commit touching `skills/cloud/gke-basics` after the pin `6f0b877` (2026-04-23) is `626910e` (2026-05-11), so "unchanged upstream through 2026-05-10" is true.

## Checked
- **Secrets:** the diff adds only placeholder values (`FAKEVALUE`, `FAKEPW`). No secret-shaped path is in scope or untracked. `git status` after all my runs is identical to its state at the start, with no new commit and no new tag.
- **Injection:** no new exec, SQL or template. The regexes are RE2.
  - `SweepStaleDirs` matches names from `ReadDir` by prefix and skips symlinks and non-directories.
  - The parent directory is `~/.config/sh-agent`, created 0700, so `RemoveAll` targets only the owner's own entries (the rule from BUG-037).
- **Authz/authn:** no endpoint or credential flow is touched.
  - **Settings write:** `attributionSet` now writes through `WriteAtomic`. The Lstat symlink refusal and the non-regular refusal still come first, and the rename replaces the path without ever writing through a link. Only the owner can write to `~/.claude`.
  - **`pushError`:** it scrubs both branches, and `errors.Is(err, errPushRejected)` still works because that branch keeps `%w`. The other branch no longer wraps its cause, and no caller unwraps it: `report.go:761` and `:787` test only `errPushRejected`. `TestPushErrorScrubsCredentials` covers both branches.
- **Crypto:** nothing composed here.
- **Dependencies:** `go.mod` and `go.sum` are unchanged. `govulncheck ./...` exited 0 with "Your code is affected by 0 vulnerabilities", plus 1 vulnerability in a required module that the code does not call. `AGENTS.md` names that one as GO-2026-5024; the default output does not print the ID. `npm audit` does not apply because there is no `package.json`.
- **GCP surface:** none in the range.
- **Tools on the live tree (go1.27.1):**
  - `gofmt -l cmd internal` printed nothing.
  - `go vet ./...` exited 0.
  - `go test -race -count=1 ./...` exited 0, with all four packages ok. This closes round 2's gap for the whole suite with `-race`.
- **Mutations:** all ran in `.agents/scratch/sec-worktree-r3/mut`, a tar copy of `git ls-files --cached --others --exclude-standard`, with `GIT_CEILING_DIRECTORIES` set so that git inside the copy cannot reach the real repository.
  - **Copy checked:** `cmp` found six files equal to the working tree.
  - **Baselines:** the targeted tests passed there first. The unmutated full baseline fails only the environmental BUG-080 test.
  - **Each mutation:** the diff was checked non-empty and the code compiled (`go test -run '^$'`) before the named test ran, and each file was restored and checked with `cmp` afterwards.
  - **Summary:** P1-P6 and A1 failed their named tests; O1 survived.
- **Past bugs (`sh-agent search --refs … --kind bug`):**
  - BUG-053 covers the needle list, the same mechanism as the first finding.
  - BUG-079 covers binary replacement, consistent with the ETXTBSY reasoning.
  - BUG-080 is the test that is environmental in the copy.
  - BUG-037 is the `RemoveAll` validation rule, which `SweepStaleDirs` meets.
  - `sweep.go`, `VENDORED.md` and `state.go` have no bug references.

## Not examined
- **Other filesystems:** `WriteAtomic` on NFS or other network filesystems. The claim "an error always means dest is unchanged" was checked by reading, for local filesystems only.
- **JSONL that fails to decode:** a line that fails to decode partway was not probed; the probe fed well-formed JSONL.
- **BUG-080 test under O1:** `TestRelativeHomeResolvesAgainstTheInvokingDirectory` cannot run in the scratch copy, so for `cmd/sh-agent`, O1's survival rests on that package's other tests.
- **`VENDORED.md` rows:** the `terraform` row, unchanged since round 1 checked it. For `gke-basics`, only the split date and the upstream log between the pin and the split were re-verified.
- **Zip path:** `updateFromZip`'s `parent` and the zip sweep's call site were not re-read; they are unchanged since round 2.
- **Search:** `internal/search` beyond the `-race` run.
- **Other agents' artifacts:** `review-worktree-code` and `verify-review-fixes` were grepped for `left: operator` records, not read whole.
- **Scratch left on disk:** `.agents/scratch/sec-worktree-r3/` (the mutation copy, the probe, the logs and the govulncheck output) and `.agents/scratch/sec-worktree-r3-live-test.log`. Both are gitignored and can be deleted.
