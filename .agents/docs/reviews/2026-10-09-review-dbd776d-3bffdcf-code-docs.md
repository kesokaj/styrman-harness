---
owner: reviewer
date: 2026-10-09
---
# Review — dbd776d-3bffdcf-code-docs   (2026-10-09)

Range reviewed: commits dbd776d..3bffdcf (git log dbd776d..3bffdcf); the diff is git diff dbd776d...3bffdcf.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem is a regression in `sh-agent install`. When install rebuilds the binary, it now builds the child's arguments with `installPassArgs`, which drops an explicit `--vendor all`. On a home whose state records one vendor, `sh-agent install --vendor all` then installs only that one vendor, and nothing reports it. No test covers this.

The second problem is in the delivery spec, the one home of the delivery claim (AGENTS.md rule 6). It now says v2.6.2 ships the `mcpServers` refusal, but that refusal landed after the tag.

The third problem is in the dream job. The new `WaitDelay` handles a CLI that hits its deadline. It also changes what happens when the CLI exits 0 but a helper it started keeps the output pipes open: that run is now logged as a failure after 30 s, its finished answer is thrown away, and a helper still in the CLI's process group is never killed. No test drives that path.

Everything else holds:
- The seven exit-2 sites, `TestMainExitStatus` and `TestVersionVerb` are correct and tested.
- The dream block-scalar fix, `trail.DocsRoot`, the `metricsSummary` and `readSmall` refactors and the `mcpServers` lint are correct and tested.
- The unexported identifiers are clean.

## Findings

- **[medium / confidence: high]** `cmd/sh-agent/main.go:664` — **`install --vendor all` is lost on the rebuild path.**
  - **Before:** `runInstall` passed `--vendor <resolved>` whenever `c.vendorFlag != ""`, so an explicit `all` reached the child.
  - **Now:** it calls `installPassArgs(c.homeFlag, c.home, c.vendorFlag)` (main.go:1111). That function adds `--vendor` only when the resolved vendor is not `all`.
  - **How it fails:** a home whose state records `claude` runs `sh-agent install --vendor all`. The parent resolves `all`, passes no `--vendor`, and `harness.Reinstall` runs the fresh binary. The child resolves the vendor from state (`resolveVendor`, main.go:878) and installs claude only. The operator's request to widen to both vendors is silently dropped. Traced by reading; not run.
  - **No test:** no test drives the rebuild path's arguments with an explicit vendor.
  - **Fix:** keep the old rule for install (pass `--vendor` whenever the flag was given), or give `installPassArgs` a "flag given" input. Add a test where state says `claude`, the flag is `--vendor all`, and the child args contain `--vendor all`.
  - **Not introduced here:** `update` has had the same drop since before this range, through the same helper. It is out of range, noted here only.

- **[medium / confidence: high]** `.agents/docs/specs/sh-agent/delivery.md:77` — **the spec claims v2.6.2 ships a refusal it does not have.**
  - The "Built and shipped (newest tag v2.6.2, 2026-10-09)" paragraph now lists `mcpServers` among the frontmatter keys the installer refuses.
  - v2.6.2 does not refuse it: `git show v2.6.2:internal/harness/inventory.go | grep -c mcpServers` prints 0. The refusal landed in 3bffdcf, after the tag.
  - This file is the one home of what a tag delivers (rule 6), so it now misstates what the shipped tag does. `TestDeliverySpecCurrent` only checks the tag string and the verb names, so it cannot catch this.
  - **Fix:** take `mcpServers` out of the shipped paragraph until the next release commit writes it, or mark it as on `main`, unreleased.

- **[medium / confidence: medium]** `internal/dream/run.go:215-217` (with `realRun`, run.go:129-138) — **a CLI that exits 0 but leaves a helper holding the pipes is now a failed run.**
  - **Mechanism:** the `WaitDelay` timer also starts when `Wait` sees the child exit normally. If the CLI exits 0 while a helper it started still holds stdout or stderr, `Run` returns `exec.ErrWaitDelay` after 30 s. `ctx.Err()` is nil, so `callWithin` returns `"<cli>: exec: WaitDelay expired before I/O complete"`.
  - **Effect:** `realRun` counts the project as failed and keeps the failure. The finished answer already in `stdout` is thrown away.
  - **Leak:** `cmd.Cancel` only runs when the context ends. A helper still in the CLI's process group is therefore never killed on a normal exit, and it outlives the run.
  - **Before:** this path hung, so the change improves the lock. But BUG-169's own premise (helpers that hold the pipes) makes this path real, and `TestCallReturnsAtTheDeadlineDespiteAGrandchild` covers only the deadline.
  - **Doc:** the new sentence in `specs/sh-agent/dream.md:34` describes only the deadline case.
  - **Fix:** after `Run` returns, kill the group regardless of how the CLI exited (`syscall.Kill(-pid, SIGKILL)`, ignore ESRCH). Treat `errors.Is(err, exec.ErrWaitDelay)` with `ctx.Err() == nil` as success. Add a test case with the script `sleep 60 & echo ok; exit 0` that expects `ok` and no error.

- **[low / confidence: high]** `internal/dream/run.go:216` — **Ctrl-C on a hand-run `sh-agent dream` leaves the CLI running.**
  - With `Setpgid: true` the CLI is no longer in the terminal's foreground process group, so Ctrl-C reaches only sh-agent.
  - Nothing in `cmd/` or `internal/` handles signals (`grep signal.Notify` finds nothing). sh-agent dies and `claude -p` / `agy` keeps running with no timeout, because the timeout lived in the dead parent. It keeps sending the project's content to the model and nothing reads the answer.
  - The systemd timer path is unaffected: stopping the unit kills its whole cgroup.
  - **Fix:** use `signal.NotifyContext` in `runDream` and pass that context into `callWithin` instead of `context.Background()`, so `Cancel` kills the group on SIGINT or SIGTERM.

- **[low / confidence: high]** `skills/harness-audit/SKILL.md:36` — **narrowed: one restatement still uses the bare command.** The new rule is to run govulncheck by its full path, because "GOPATH/bin is not on every shell's PATH, and two review agents skipped it" (`AGENTS.md:30`). The same skill's `references/areas.md:36` was updated to match, but the skill's own baseline step at SKILL.md:36 still restates the setup command as bare `govulncheck ./...`. The audit's main step therefore keeps the form the change says gets skipped. **Fix:** use the full path there too.

- **[low / confidence: medium]** `.agents/docs/specs/sh-agent/search.md:29` — **the spec's exit-status line does not name the cases that moved to exit 2.**
  - The spec says "1 on a refusal or an error; 2 on a usage error (a flag the verb does not define, a stray argument)".
  - These argument errors used to exit 1 and now exit 2 (`internal/search/search.go:98-128`, `cmd/sh-agent/main.go:624`): `--limit` out of range, an unknown `--kind`, no query, words combined with `--refs` or `--open`, `--rebuild` on `search`, and `index` given words.
  - The spec names none of them, so a reader would take them for refusals (exit 1).
  - **Fix:** name them in the parenthesis.

- **[low / confidence: medium]** `internal/harness/report.go:150-152` with `report_test.go:157-197` — **known over-redactions are not pinned in the test, as AGENTS.md requires.** AGENTS.md says over-redaction in the safe direction "is pinned there". I checked the three new regexps, copied into scratch, against these strings. Each one is redacted:
  - `docker exec -u root:root c sh`
  - `chown -u www-data:www-data`
  - `sudo -u postgres:x`
  - `"auth": "required"`
  - `credentials: application-default`

  None of them is in the pinned over-redaction list, and the comment at report.go:140-149 names only `-u 1000:1000` as staying. A later change could reverse any of them without anyone deciding to.

  For context: running the three regexps over this repo's tracked `.agents/docs`, `skills` (catalog excluded), `context`, `agents`, `README.md` and `AGENTS.md` gave 4 hits, all review records quoting `-u user:pass`. No ordinary prose in the trail over-matches.

  **Fix:** add `-u root:root` and `"auth": "required"` (at least) to the `over` list.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:126-134` with `main.go:883` — **a bad vendor in the state file is now reported as a usage error (exit 2).** `resolveVendor` also normalises the vendor recorded in the home's state. A corrupt recorded vendor now prints `sh-agent <verb>: unknown vendor …` and exits 2, although no argument was wrong. **Fix:** return exit 2 only when the error came from the flag; let a bad state value stay an ordinary error (exit 1).

## Criteria

1. **Dead code — pass.**
   - Removed together with their last callers: `frontField` and `unquote` (dream), `Lint`, `IsBusy`, and the metrics copy of `docsRootOf`. `search.docsRootOf` keeps a caller.
   - The unexported `loadDocs`, `latestRelease` and `refusal` are named by no living doc or test.
   - The test failure messages in `internal/harness/harness_test.go` at lines 89, 102, 2379 and 2386 still say "Lint" for a function that no longer exists. That is cosmetic, not a finding.
2. **Flow tracing — fail.**
   - Install's rebuild arguments (finding 1). The dream normal-exit path (finding 3).
   - Exit-2 sites, each with a test that fails without the change: the seven sites are covered by `TestUsageErrorsExit2` (`errors.As` would fail on a plain error) and by the zip and search tests. `metrics show` and `metrics report` share `metricsSummary`; `Summary` is plain data, so closing the store first is safe.
   - The dream block-scalar fix: `TestSkillsReadsBlockScalarDescription` would fail on the old `>`.
   - The `mcpServers` lint: `TestLintPermissionKeysByKind` covers both the block and the flow form, and would load the agent without the change.
   - The process group: with `WaitDelay` or `Setpgid` removed, the deadline test blocks for 60 s against its 20 s guard. This was checked by reading, not mutated.
3. **Standards conformance — fail (docs).** Findings 2, 5, 6 and 7. `gofmt -l cmd internal` is empty and `go vet ./...` is clean (both run).
4. **Directory layout — pass.** `trail.DocsRoot` sits in `internal/trail` and the new tests sit beside their code.
5. **Simplicity (DIY-except-security) — pass.** No new dependency and no hand-rolled security code; the duplicated helpers were consolidated.
6. **Decoupling — pass.** `metrics`, `search` and `dream` share `trail.DocsRoot` without importing each other. `dream` now reads frontmatter through `harness.Split` and `harness.Parse`, and already imported `harness`.

## Docs lens

- `restatement-grep.sh -i` was run for `permissionMode`, `govulncheck ./`, `971`, `LatestRelease`, `LoadDocs`, `IsBusy`, `Lint(`, `Refusal{` and `frontField`.
- Living hits that were read:
  - `delivery.md:77`, `AGENTS.md:63` and `README.md:233` for `mcpServers`. All three are updated; only delivery.md is wrong, and only about which release has it (finding 2).
  - `skills/harness-audit/SKILL.md:36` (finding 5). Its neighbours `agents/security.md:30`, `skills/operate-gcp/SKILL.md:91` and `skills/init-project/assets/.../gcp.md:83` are generic guidance for any project, not this repo's setup command. They are left as they are; low confidence that they need the full path.
  - The corpus figure: `guide.md:141`, `AGENTS.md:72` and `dream.md:47` agree (971 on one box; dream.md already names the box).
- Dropped claim checked: init-project's description no longer says "no CLAUDE.md: 2.1.277+ reads AGENTS.md natively". The claim still lives in `skills/init-project/SKILL.md:61`, so it is not dropped.
- `git grep` at 3bffdcf found none of the reclassified usage-error messages in `.agents/docs/testing`, `specs`, `README.md` or `skills`.
- The release step is a documented deploy step: `deploy/remote.md:99` now runs `TestDeliverySpecCurrent`.
- `README.md:189` and `README.md:198` already said every usage error exits 2, which the code now matches.
- Stale tests: none found.

## Tests run (this review)

- `go test ./cmd/sh-agent -run 'TestUsageErrorsExit2|TestVersionVerb|TestMainExitStatus|FromZip|TestDeliverySpecCurrent|TestSearchVerb' -count=1`: ok.
- `go test ./internal/{dream,harness,trail,...} -run 'TestCallReturns|TestSkillsReads|TestReportRedact|TestLint|TestDocsRoot|TestParseStrict|TestAgentDescription' -count=1`: ok.
- `gofmt -l cmd internal`: empty. `go vet ./...`: clean.

## Not examined

- The full `go test -race -timeout 30m ./...` and `govulncheck` were not run.
- No mutation was run: every must-fail claim above comes from reading the test against the old code.
- The bug records BUG-167, BUG-168 and BUG-169 and the `INDEX.md` rows were not read: they are records, out of the docs lens.
- Not verified on a real system: findings 1, 3 and 4 were traced by reading, with no reproduction in a scratch home or with a real CLI.
- `.agents/docs/testing/*.sh` kits were not grepped for the changed messages; only `.agents/docs/testing` was searched as a whole for the error strings.

## Round 2

Range reviewed: the round-2 fixes, uncommitted in the working tree on top of 3bffdcf (`git diff 3bffdcf`: 9 files, +207/-28), reviewed on 2026-10-09.

### Status: COMPLETE

### Verdict: REQUEST CHANGES

Round 1's findings 1 to 6 are fixed, and each fix has a test that fails without it.

The rewrite of the redaction needles brought in a new regression, which is the reason for this verdict:
- The command-line needle lost its leading anchor, and the space after `-u` became optional.
- It now matches `-u` inside any hyphenated word that is followed by `:` or `%`. That includes `path/foo-usage.md:72`, `web-ui:8080`, `sh-updater:v1` and `date -u +%FT%TZ`; the last is a command the harness's own `dream`, `distill` and `retro` skills write into every project's trail.
- Over this repository's tracked living trail, the needle hits 49 lines (round 1: 4, all real quotes).
- `redactLines` replaces each matched line with a marker, and `dream` runs the same `Redact` over the trail it reads (`internal/dream/collect.go:157`). Both `sh-agent report` and the nightly dreamer therefore lose those lines.

### Findings

- **[medium / confidence: high]** `internal/harness/report.go:159` — **the command-line needle over-matches ordinary paths, hosts and commands.**
  - Matched strings, checked against a scratch copy of the three regexps:
    - `see internal/harness/run-usage.md:72` → `-usage.md:72`
    - `open http://web-ui:8080` → `-ui:8080`
    - `image: ghcr.io/x/sh-updater:v1` → `-updater:v1`
    - `stamp=$(date -u +%FT%TZ)` → `-u +%FT%TZ)`
    - `mysql is up; ls -a /var/lib` → the whole tail
    - `echo "mysql" && mkdir -p .agents/scratch` → the whole tail
  - **Scale:** over `git ls-files` of `.agents/docs`, `skills` (catalog excluded), `context`, `agents`, `README.md` and `AGENTS.md`, the needle matches 49 lines. Most are path:line citations to a file whose name holds `-u…` (the `*-usage.md:NN` family), `date -u +%FT%TZ` (8 lines) and `-url:` names (6 lines).
  - **Effect:** each matched line becomes a redaction marker, in the report bundle and in the dreamer's input alike.
  - **Not pinned:** none of these shapes is in the pinned over-redaction list in `report_test.go:169-177`. AGENTS.md requires that over-redaction is pinned there; these are not.
  - **Causes:**
    - The `-u` branch dropped both its `(^|\s)` prefix and the required `[ =]` separator. The attached form (`-uadmin:x`) cannot be told apart from a hyphenated word.
    - `%` as a separator takes `date`'s format string as `user%pass`.
    - The client branch (`mysql\S*` … `[^\n]*?\s-[pa]`) lets any later `-a` or `-p` on the same line count as the client's password flag.
  - **Fix:** in front of `-u`/`--user`, require a character that is not part of a word: whitespace, a quote, a backtick, `(`, `=`, or a literal `\n`. Keeping the rule's "no `\b`" intent does not require leaving the needle with no anchor at all. Also:
    - Do not let the user part start with `+`.
    - Bound the client branch to the client's own arguments (stop at `;`, `&` or `|`).
    - Add `run-usage.md:72`, `web-ui:8080` and `date -u +%FT%TZ` as benign near-misses.
  - **Out of scope here:** whether this weakens or strengthens the security lens's round-1 findings on the same lines is the security lens's call.

- **[low / confidence: medium]** `internal/dream/run_test.go:779-781` — **`TestCallKeepsTheAnswerWhenAHelperOutlivesTheCLI` can flake.** It checks `alive(pid)` once, straight after `callWithin` returns. `killGroup` has only just sent SIGKILL, and `kill(pid, 0)` also succeeds on a zombie that its new parent (init or a subreaper) has not reaped yet. The interrupt test retries for 3 s; this one does not.
  - **Fix:** use the same retry loop.
  - Both tests also assume a PID 1 that reaps. In a container without an init that reaps, they fail.

### Round-1 findings, re-checked

1. **Fixed — install's `--vendor`.** `installPassArgs` (`cmd/sh-agent/main.go:1113`) now forwards `--vendor` when the flag was given, or when the resolved vendor is narrower than all.
   - It has three callers: the rebuild path of `runInstall` (main.go:664), the release-binary path of `runUpdate` (main.go:709) and the checkout path of `update` (main.go:1086). All three want an explicit vendor forwarded: an explicit `--vendor all` on update now widens a narrowed home, as on install. This also closes the pre-existing update case noted in round 1.
   - `TestInstallPassArgsForwardsAnExplicitVendor`: the "explicit all" case gets `""` from the old code, so it fails without the fix.
   - The test does not drive `runInstall` or `update` themselves. They call the helper directly (read), so I accept that.
2. **Fixed — delivery.md.** `delivery.md:77` no longer names `mcpServers` for v2.6.2. `README.md:233` and `AGENTS.md:63` describe the checkout and keep it, which is correct.
3. **Fixed — the dream CLI's process group.** `callWithin` (`internal/dream/run.go:221-249`):
   - **Every return path returns the collected output:**
     - A `Start` failure returns `""`, `""`; nothing was collected.
     - The timeout, error and success paths all return `out` and `errb`.
     - After `ErrWaitDelay` the output is complete: the standard library closes the pipes and then waits for the copying goroutines (`$GOROOT/src/os/exec/exec.go:1011-1016`), so there is no race on the buffers.
   - `ErrWaitDelay` is returned only when the CLI exited 0 (a non-zero exit returns `ExitError`). The deadline is checked first.
   - **The group is killed after `Wait` on every path.** A helper is therefore killed after a clean exit, but only once the grace has run out, so such a run takes up to 30 s longer. That is acceptable.
   - **The signal goroutine cannot outlive the call.** It returns on `waited`, which is closed right after `Wait`, or on a signal. On a signal it kills the group and then kills its own process with the default disposition restored.
   - **Two windows remain, each microseconds wide:**
     - A signal delivered after the goroutine took `waited` and before `signal.Stop` is buffered and never acted on, so that Ctrl-C is lost.
     - A Ctrl-C between `Start` and `Notify` kills sh-agent and orphans the CLI.

     Both are noted, not findings.
   - **Tests:**
     - `TestCallKeepsTheAnswerWhenAHelperOutlivesTheCLI`: the old code returned `ErrWaitDelay`.
     - `TestCallKillsTheGroupOnInterrupt`: the old code had no handler, so the helper survived.

     Both fail without the fix, by reading. Both pass with `-race` here (run).
4. **Fixed — Ctrl-C on a hand-run `sh-agent dream`.** Same change and same test as item 3.
5. **Fixed — harness-audit.** `skills/harness-audit/SKILL.md:36` now uses the full path, matching `AGENTS.md:30` and `references/areas.md:36`.
6. **Fixed — search.md.** `search.md:29` now names the argument errors that exit 2. `1..50` matches `maxLimit = 50` (`internal/search/search.go:35`).
7. **Superseded — the over-redactions.** `-u 1000:1000`, `-u node:node`, `mysql -P3306`, `mysql -p …` and `credentials: application-default` are now pinned, and `"auth": "required"` is a near-miss. The rewrite, however, adds the over-matches in the finding above.
8. **Not addressed — exit 2 for a bad recorded vendor (low).** `ctx.vendor` still returns exit 2 for a bad vendor read from the state file. It is not in this fix set, and I record it only as still open.

### Criteria (round 2)

1. **Dead code — pass.**
   - Nothing was moved or re-declared. `call`, `callWithin` and `killGrace` are unchanged in shape.
   - `installPassArgs` keeps its three callers.
   - The new test helpers `readPid`, `alive` and `killPidFile` are each used.
2. **Flow tracing — pass**, apart from the redaction flow into `dream` (finding above).
3. **Standards conformance — fail.**
   - AGENTS.md's redaction rule, "over-redaction in the safe direction is pinned there": the new over-matches are not pinned.
   - `gofmt -l cmd internal` is empty, and `go vet` over `cmd`, `internal/dream` and `internal/harness` is clean (run).
4. **Directory layout — pass.**
5. **Simplicity (DIY-except-security) — pass.** Stdlib only (`os/signal`).
6. **Decoupling — pass.**

### Docs lens (round 2)

- `delivery.md:77`, `search.md:29` and `harness-audit/SKILL.md:36` were re-read and are correct.
- The comment above the needles (`report.go:140-158`) says `docker run -u 1000:1000` "goes too, pinned". It does not mention the path-citation and `date -u` over-matches; it will need the same edit as the finding.
- `specs/sh-agent/dream.md:34` still says only "cut loose 30 s later". It does not mention the group kill after a clean exit or on an interrupt.
  - **[low / confidence: high]** `.agents/docs/specs/sh-agent/dream.md:34` — **stale.** Add "killed whole after the CLI exits and on SIGINT or SIGTERM; a clean exit whose helper kept the pipes keeps its answer."

### Tests run (round 2)

- `go test -race -count=1 ./internal/dream -run TestCall`: ok.
- `go test -count=1 ./internal/harness -run TestReportRedact`: ok.
- `go test -count=1 ./cmd/sh-agent -run 'TestInstallPassArgs|TestUsageErrorsExit2'`: ok.
- `gofmt -l cmd internal`: empty. `go vet` over `cmd`, `internal/dream` and `internal/harness`: clean.

### Not examined (round 2)

- No mutation was run: the claims that tests fail without the fixes come from reading.
- The full `go test -race ./...` was not run.
- Scanned for over-matching: the auth-key and credentials-key needles over the trail. That scan found 0 hits.
- Not scanned for over-matching: transcripts and feedback logs.
- The security lens's round-1 findings on the needles were not read. Whether the rewrite closes them is that lens's question.
- The interrupt path was exercised only through the test's helper binary, never against a real `claude` or `agy`.
