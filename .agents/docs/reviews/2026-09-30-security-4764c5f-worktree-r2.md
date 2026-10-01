---
owner: security
date: 2026-09-30
---
# Security review — 4764c5f-worktree, round 2   (2026-09-30)

Range reviewed: 4764c5f..working tree (`git diff 4764c5f`, 16 files, +948/-95; the untracked trail files under `.agents/docs/` are not in the range)

## Status: COMPLETE

## Verdict: FINDINGS

5 findings, all low.

**One round-1 answer does not run.** The `--no-optional-locks` added for round-1 finding 6 is a global git option, but it was placed after `status`.
- git 2.47.3 rejects it (exit 129), so `trailDates` never learns which records are modified.
- There is no security exposure, but it is a regression that the test suite does not catch (F1).

**The other answers, measured:**
- Four hold: round-1 findings 1, 2, 3 and 7.
- Two hold in part: finding 4 (`isMeta` is still sent, F5) and finding 5 (the label is not escaped, F3).
- All four `harness.Redact` lines that round 1 measured as leaking are redacted now.

**New in this round:**
- Compound key names get the token threshold, so short passwords under them leak (F2).
- Uninstall's link guard has no test that fails without it (F4).

## Findings
- **[low / confidence: high]** `internal/dream/collect.go:227` — the answer to round-1 finding 6 does not run.
  - **Measured.** `--no-optional-locks` is a git global option. Placed after `status`, git 2.47.3 rejects it: `git -C . status --no-optional-locks --porcelain …` exited 129 with "error: unknown option `no-optional-locks'". The global form `git --no-optional-locks -C . status …` exited 0. Both were run on this repository and are read-only.
  - **Effect.** `.Output()` returns that error, so `trailDates` returns an empty `modified` map on every call.
    - A tracked trail record that was edited but not committed is dated by its last commit instead of its mtime.
    - That is the opposite of the doc comment (`collect.go:203-206`) and of the test header (`collect_test.go:468-469`): "an untracked or modified one by its mtime".
    - When that commit is older than the watermark, the edit is not sent until it is committed.
  - **Security effect: none.** The index-lock concern is gone only because git exits before it reads the index.
  - **Why no test fails.** `TestTrailDateFromGit` covers only the tracked, unmodified case. I ran it and it passes.
  - **Remediation.**
    - Put the flag before the subcommand: `git --no-optional-locks -c core.fsmonitor=false -C root status --porcelain --untracked-files=all -- docs`.
    - The `core.fsmonitor=false` matters because a status that actually runs brings back round 1's (b): a repository-local hook executed at night.
    - Add a test with a tracked record modified after its commit that must keep its fresh mtime. It fails on the current code.
    - `-z` (as `reconcile.go` uses) is correctness only, for quoted non-ASCII names.
- **[low / confidence: high]** `internal/harness/report.go:80-81` — the value threshold is chosen by the leftmost keyword the regex meets, not by the key's name.
  - **Mechanism.** `[a-z0-9_-]*` swallows the rest of the name. So `TOKEN_PASSWORD`, `api_key_password`, `API_KEY_SECRET` and `APIKEY_PASS` are judged as tokens (8 characters, or 6 with a digit), not as passwords or secrets (4).
  - **Measured.** Probe `.agents/scratch/sec-review-r2/probe/main.go` called `harness.Redact` on the working tree:
    - KEPT: `TOKEN_PASSWORD=letmein`, `api_key_password: letmein`, `API_KEY_SECRET: letmein`.
    - REDACTED: `password: letmein`.
  - **Before this range.** By reading the old pattern, all three kept lines were redacted at 4764c5f: it matched the key alone, whatever value followed.
  - **Related, from the same probe:**
    - An opening quote with no closing quote on the line falls to the unquoted rule: `api_key="AbCdEf` is KEPT.
    - `password: my dog rex 2019` is KEPT. This is the same class as the pinned `secret = my passphrase`, so it is accepted. But the comment at `report.go:66-74` states the rule only by example. The rule actually applied is "an unquoted value is judged by its first word".
  - **Remediation (no dependency).**
    1. Pick the threshold from the whole key name: the part of the match before `:`/`=`. Use four when it contains `pass` or `secret` anywhere.
    2. Add the three compound keys as must-redact rows in `TestReportRedactShapes`. Confirm each row fails on the current code.
    3. State the first-word rule in the needle comment.
- **[low / confidence: high]** `internal/dream/prompt.go:51` — the fence escape reaches the text but not the label.
  - **Mechanism.**
    - `in.Label` goes onto the `<<<data` line unescaped.
    - A trail label is `dir/<file name>` (`collect.go:94`).
    - `mdFiles` takes any regular `*.md` entry, so a name carrying `\n>>>\n<<<data memory: …` is taken as is. git stores and checks out such a name on Linux and macOS.
  - **Measured through the real `prompt()`.** The probe `.agents/scratch/sec-review-r2/zz_secprobe_r2_test.go` was added to the package only through `go test -overlay`; `git status` was unchanged afterwards.
    - Two inputs produced three lines equal to `>>>`.
    - A `<<<data memory: forged` line appeared at column 0.
    - The text escape itself holds: a run of five `>` came out as `> > > > >`.
  - **Impact.** A repository the operator works in can still end its own block and open one labelled with another kind, which is what SD8 set out to prevent. It stays low for round 1's reasons: the fence is a soft, model-level boundary, and `/retro` re-derives dream output.
  - **Also still open** from round 1's minimum remediation: `<<<data` inside the text is not escaped.
  - **Remediation.**
    - Fold the label the way `line()` does (no line breaks or control characters).
    - Escape `<<<` runs the same way as `>>>` runs.
    - Add a test row with a newline in a label.
    - *Robust alternative (unchanged):* a per-run nonce in the fence from `crypto/rand.Text()` (stdlib).
- **[low / confidence: high]** `internal/harness/apply.go:327` — Uninstall's link guard has no test that fails without it. That guard is the one line that closes both triggers round-1 finding 2 named: a failed install followed by uninstall, and an uninstall against a manifest written before this change.
  - **Mutation.**
    - I copied `apply.go` to `.agents/scratch/sec-review-r2/apply_m1.go` and changed line 327 to `if false && fi != nil {`. It compiled, and `diff` shows the one line.
    - `go test -overlay` swapped the copy in, and the whole `internal/harness` package passed (`ok`, exit 0).
    - The unmutated baseline also passed.
  - **Why.** In `TestUninstallAndLinkedDirs` the reinstall leaves the files under the linked directory out of the new manifest. By the time `Uninstall` runs, no listed path is under the link, so the guard is never reached.
  - **The code itself is right as read.**
  - **Remediation.** Write a test that:
    1. writes a manifest listing a path that is a symlink, and one that sits under a linked directory;
    2. runs `Uninstall`;
    3. asserts that the links and their targets survive.

    Confirm it fails with line 327 mutated.
- **[low / confidence: medium]** `internal/dream/collect.go:423-455` (`contentText`, called from `sessionText`) — round-1 finding 4 is answered for tags but not for `isMeta`.
  - **Measured.** `.agents/scratch/sec-review-r2/tags.py` read this machine's 413 Claude Code transcripts and printed tag names and counts only, never content.
    - The four prefixes now drop 629 user string contents.
    - No user string content, and no text part, that opens with or contains a harness tag is still sent.
  - **Still sent:** 529 `isMeta` string records and 306 `isMeta` text parts. These are injected, not typed: expanded skill bodies and caveats, per round 1's reading.
  - This is harness text rather than the operator's command output, so the exposure is small. It is the part of round-1 finding 4 left with neither a fix nor a rebuttal.
  - **Remediation:** skip a record with `isMeta: true` in `sessionText` (one condition), or record the rebuttal.

## Round-1 answers, as measured
1. **report.go secret assignment: fixed.**
   - The four lines are redacted: `export DB_PASSWORD="correct horse battery staple"`, `PASSWORD='my dog rex 2019'`, `password: letmein` and `api_key: "AbCdEf"`.
   - The controls are unchanged: `password=hunter2` and `curl -b session=abc123` are redacted; `## Secrets:` and `curl -b cookies.txt` are kept.
   - The pinned under-redaction behaves as pinned: `token=abcdefg` and `secret = my passphrase` are kept.
   - New gap: F2.
   - Still unnamed in the comment, as round 1 asked it to be: the line-based limit. `secret: |` and `password:` with the value on the next line are both KEPT (measured). This is the same at 4764c5f.
2. **apply.go, links at the last component: fixed in code.** Every write, chmod and remove site in `apply.go` sits behind the guard:
   - the write loop;
   - the unchanged-branch `os.Chmod` at :90;
   - the resync `os.Remove` at :143;
   - `saveOnError`'s union;
   - `Uninstall` at :327/:338.

   `WriteAtomic`'s temp file sits in the guarded directory. `pruneEmptyDirs` (:437) is reached only after a guarded removal, so every directory it walks passed `underLink`. Test gap: F4.
3. **apply.go, links at any depth: fixed.** `underLink` covers parent-directory links, dangling ones included (`Lstat`). The printed claim now holds.
4. **Transcript filter: tags fixed on this corpus.** `isMeta` is still sent (F5).
5. **Fence: the text is fixed** (`>{3,}` becomes spaced; runs of five and eight were measured). The label is not escaped, and neither is `<<<data` in the text (F3).
6. **`git status` lock: the answer does not run** (F1).
7. **`indent()` and bare CR: fixed.** Measured through the overlay probe: 0 lines left unindented across CR, CRLF and LF, and no bare CR left.

## Also noticed (not security)
- `internal/harness/apply.go:405-411` — `underLink` was inserted between `pruneEmptyDirs`'s doc comment and its function. Godoc now attributes "pruneEmptyDirs removes now-empty directories…" to `underLink`, and `pruneEmptyDirs` has no comment.
- **A linked vendor root.** When `~/.claude` or `~/.gemini` is itself a link (a dotfiles setup), the under-link rule refuses the whole vendor tree.
  - Every file prints "keep … (a symlink; …)" and nothing is installed.
  - `Status` does not call `underLink`. It reads through the link and reports drift that install never fixes.
  - This machine's home has real directories (`ls -ld`, read-only), so it does not bite here.
  - Round 1 named the gentler option: `os.OpenRoot`, which still allows links inside the home. Operator choice.
- The keep line says "(a symlink; …)" also for a regular file under a linked directory. This is wording only.

## Checked
- **Secrets:** the range adds no credential; the test values are fake. The new `answer()` error wraps the scanner error with no stream content and goes only to the dream log.
- **Injection:**
  - The `trailDates` argv keeps round 1's shape: no shell, `root` after `-C`, `docs` after `--`.
  - The fence and label are covered in F3.
  - The unchanged `needleIn` loop is line-based, as the probe confirms: a next-line value is kept. RE2 means no backtracking blow-up.
- **Authz/authn:** not touched by the range.
- **Crypto:** nothing added or changed; no hand-rolled primitive.
- **Dependencies:**
  - `go.mod` and `go.sum` are not in the range.
  - `govulncheck ./...` on go1.27.1 exited 0: 0 vulnerabilities in called code. With `-show verbose`, one module-level item: GO-2026-5024 in golang.org/x/sys@v0.37.0, not called. This is the known item in AGENTS.md. Output: `.agents/scratch/sec-review-r2/govulncheck.txt`.
  - There is no `package.json`, so `npm audit` does not apply.
- **GCP surface:** none in the range.
- **Tests run:**
  - `go test -count=1 ./internal/harness/`: ok, as the baseline for F4.
  - `go test -count=1 -run 'Trail|Date|Modified|Git' ./internal/dream/`: ok.
  - Every probe went through `-overlay` or a separate `go run` file. `git status --short` was unchanged afterwards.
- **Scope since round 1:** I compared the saved round-1 patch with this one file by file. Only `collect.go`, `prompt.go`, `run.go`, `apply.go`, `report.go` and their tests changed. `cmd/sh-agent/main.go`, `internal/search/search.go`, `internal/search/dream.go` and their tests are byte-identical to what round 1 read.

## Not examined
- **Test bodies**, beyond `TestTrailDateFromGit` and `TestUninstallAndLinkedDirs`.
  - `TestInstallMidFailureKeepsManifestComplete`, `TestInstallDoesNotFollowLinkedDest` and `TestResyncSkipsLink` ran only as part of the whole-package run.
  - The rest of `harness_test.go`, `run_test.go`, `collect_test.go`, `dream_test.go`, `main_test.go`, `memory_test.go` and `search_test.go` were not read.
- **`saveOnError`'s symlink/underLink filter** was not mutated. Uninstall's guard would mask such a mutation, so it would say nothing.
- **Unrun remediation:** F1's corrected command was not run against a project with a modified record.
- **Model behaviour:** whether the dreamer model acts on a forged `<<<data` block was not measured.
- **Other labels:** the memory-store (`memory.Store`), Antigravity transcript and skill labels were not traced for newline-bearing names. Only the trail label path was read.
- **`isMeta` records** were counted, not read. "Harness-injected text" comes from round 1's reading and the field's meaning.
- **`harness.Redact` callers** other than report and dream `scrub` were not traced (as in round 1).
- **Still open, outside the range (round 1):** `Uninstall` removes a listed file the operator has since overwritten; it does not check `IsOurs`.

## Outcome

Answered in the next round (see `2026-09-30-security-4764c5f-worktree.md`); the fixes and their tests are in the phase commit. Every finding of this round was fixed, except the deploy-guide sentence, which the documentation pass of the phase fixes, and the nested-first prose-loss concern, rebutted by the corpus measurement (971 leaves and 0 argument keys before and after, `testing/2026-09-30-verify-audit-fix-phase-1.md`).
