---
owner: reviewer
date: 2026-09-30
---
# Review — f1ca525-e2cffc1-code   (2026-09-30)

Range reviewed: f1ca525..e2cffc1 (round 2: the fixes for round 1's code review, plus the security and docs round-1 fixes in the same commit; source part `git diff f1ca525..e2cffc1 -- cmd internal`)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The code changes for all eight round-1 findings are correct as written, and I found no new high-severity defect. `go vet ./...` is clean, and `go test ./internal/dream ./internal/memory ./internal/search ./cmd/sh-agent` passes on this box. The blocking gap is test coverage, because under the project's Definition of Done a finding counts as fixed only when a test fails without the fix. Two Outcome claims have no such test:
- Round-1 finding 6: an Antigravity read error now fails the project and keeps its watermark, but no test drives that path.
- Round-1 finding 3, second half: a `saveState` failure now gives exit 2, but no test drives that path.

Reverting either fix leaves every test green. The other findings below are low: a false "under the lock" claim that lets a real run delete a concurrent dry run's copies, two targets() branches now unreachable, an unrecorded source for the step-type mapping, a symlink check that misses `.agents` itself, and a possible auth loss from `--setting-sources local`.

## Findings
- **[medium / confidence: high]** `internal/dream/collect.go:94-97`, and `collect.go:86-90` with `:361`: round-1 finding 6 is fixed in code but has no test.
  - No test in `internal/dream` creates an Antigravity database: a grep for `conversation_summaries` in `internal/dream/*_test.go` finds nothing, and the Outcome line names no test.
  - So putting back the old `c.Dropped = append(...)` at line 96 leaves every test green, and the watermark-advance defect could come back unnoticed.
  - The new `sessionText` → `sc.Err()` → `Dropped` path (a transcript line over 16 MB) is untested as well.
  - Fix: add a Run-level test with a corrupt `~/.gemini/antigravity-cli/conversation_summaries.db` in the fixture HOME. Assert exit 1, `st.Projects[root].IsZero()`, and a `fail reading Antigravity's conversations` log line. Add a collect test with an over-long transcript line that asserts the `claude transcript …` Dropped entry.
- **[medium / confidence: high]** `internal/dream/run.go:146-157`: Outcome 3 says "`saveState` errors are exit 2 at the end of the run", but no test makes the state write fail.
  - Changing the `stateErr` handling back to ignoring the error passes the whole suite.
  - Fix: add a test where the state write fails after the lock is taken. For example, pre-create `dream.lock`, `dream.log` and `dream-work`, then make `~/.config/sh-agent` 0o500 so that `WriteAtomic`'s temp file cannot be created.
  - Assert exit 2 with "watermarks could not be saved", and assert that the project's proposals were still written.
- **[low / confidence: high]** `internal/dream/run.go:295-297` (`removeCopies`): the comment says "It runs under the lock, so no live run owns them". That is false for `sh-agent dream --dry-run`.
  - The dry run takes no lock, yet it creates `dream-projects-*` (`internal/search/dream.go:55`) and `dream-agy-*` (`internal/search/dream.go:109`) in the same directory.
  - A real run (the timer, or a second shell) that starts during a dry run deletes those live copies, and the dry run then fails with a read error.
  - Fix: only remove directories older than a bound (for example `cliTimeout`, since no copy outlives a run), or have the dry run take the lock too. Also correct the comment and the spec row at `specs/sh-agent/dream.md:51`, which repeats "under the lock".
- **[low / confidence: high]** `internal/dream/dream.go:87-88` and `:102`: the `--all`-with-`--project` branch and the usage branch of `targets()` are now unreachable.
  - `Run` refuses both cases at `dream.go:64-68` before any caller reaches `targets()`.
  - `collect_test.go:227-229` pins the dead "not both" branch.
  - The dead usage text has also drifted from the live one: it lacks `--install-timer | --remove-timer`.
  - Fix: delete both branches and the test lines. `targets()` then only handles `o.Project` and `o.All`.
- **[low / confidence: low]** `internal/search/dream.go:221-223`: the comment attributes step types 14 and 15 to "the D8 probe", but D8a (`plans/2026-09-30-parity-memory-dream.md:296`) records no step types and no `step_type` column. The mapping has no recorded source anywhere in the trail I read.
  - The test fixture (`internal/search/dream_test.go:98-111`) encodes the same assumption, so it cannot catch a wrong mapping. The verifier's mutation (a) only proves that the filter exists.
  - Failure mode if the mapping is wrong: every Antigravity transcript comes back empty and silently, with no error. If a conversation database lacks the column, that project instead fails every night.
  - I could not check a real conversation database; the read was denied. See Not examined.
  - Fix: record the probe (the step-type counts from one real conversation) in D8a or in a testing log, and point the comment at that record.
- **[low / confidence: medium]** `internal/dream/collect.go:116-120` (`linked`): `os.Lstat` only tests the last path element. When the docs root is `.agents/docs`, a committed `.agents` symlink is followed.
  - For example, `.agents -> ../other/.agents` sends another checkout's trail to the model under this project's name.
  - `memory.Store` refuses exactly this shape for the store, so the two readers disagree.
  - Fix: also `Lstat` each element of `docs` (here `root/.agents`), and extend `TestCollectDoesNotFollowLinkedTrail` to cover it.
- **[low / confidence: low]** `internal/dream/dream.go:181`: `--setting-sources local` skips `~/.claude/settings.json`, including its `env` block and any `apiKeyHelper`.
  - An operator whose Vertex/Bedrock switch or key helper lives there would get an auth failure on every project, every night (exit 1).
  - The probe ran on a box where authentication comes from elsewhere.
  - Fix: confirm which settings the dream call needs. If user-level auth settings matter, state the limitation in `specs/sh-agent/dream.md`, or pass the auth environment explicitly.

### Round-1 Outcome lines checked against code and tests
1. agy prompt on stdin: **holds**. `command()` has no prompt in argv, and `call()` writes one `{"event":"user",…}` line. `TestRunAgyPrivateHomeAndStdin` asserts that argv lacks the prompt and that stdin decodes. The real agy probe was not re-run.
2. Pre-project exit 2 and record: **holds**. Usage is refused before the lock (`TestDreamVerb` asserts that nothing is written). Index failures go through `fail()` (`TestRunPreProjectFailureIsExit2AndRecorded`). An unreadable state logs a line and deliberately leaves the state unwritten; `StatusLine` reports it.
3. Corrupt state: **holds for `readState`** (`TestReadState`, `TestRunCorruptStateIsExit2AndKept`). **Untested for `saveState`** (finding 2).
4. Memory lock: **holds**. flock on a separate descriptor per call. `TestHeldLockRefusesAfterRetry` holds a flock on an empty file, which the old pid-file code would have judged stale, so that test fails without the fix. There is no concurrent-writer test, but flock semantics rule out the three races structurally.
5. Dream lock: **holds**. `TestRunLock` holds a real flock on an empty file, which the old code also treated as stale.
6. Antigravity read error: **code holds, untested** (finding 1). Claude Code transcript errors deliberately go to "left out" while the watermark advances. I accept that as stated, because a deterministic `ErrTooLong` would otherwise wedge the project for good, which was round 1's argv failure shape.
7. Timer flags with `--dry-run`: **holds**. `TestDreamTimerDryRun` fails without the fix, because `InstallTimer` returns an error when no binary is present and the remove path never prints "[dry-run] would".
8. Index link-target match: **holds** (`TestIndexMatchIsTheLinkTarget`). The pre-lock full-index refusal is re-checked under the lock.

## Criteria
1. Dead code: **fail** (low). The `targets()` branches at dream.go:87-88 and :102 are unreachable, and collect_test.go:227-229 pins one of them. The memory `lockAge`/`stale` code, the `strconv` import and `harness.Alive` in the lock paths were removed cleanly.
2. Flow tracing: **fail** on coverage (findings 1 and 2). Otherwise the paths traced cleanly:
   - The order lock → readState → finish/fail → saveState is consistent. The error paths return ExitError 2 or 1 as the spec table says, and `fail()` ignores a second saveState error only on a path that already exits 2.
   - `answer()` treats `denied_actions` as no answer, which becomes a failure plus a kept watermark (`TestRunAgyDeniedToolIsNoAnswer`).
   - The private HOME is removed by `defer cleanup()` on every path after it is created, and leftovers are removed at the next start (but see finding 3).
   - `call()` sets `cmd.Env = nil` for claude, so claude inherits the environment.
3. Standards conformance: **fail** (low). The `removeCopies` comment and `specs/sh-agent/dream.md:51` claim "under the lock" (finding 3). The step-type comment cites a probe that does not record the mapping (finding 5). Also, `specs/sh-agent/dream.md:24` still says "the child keeping sh-agent's own environment", while line 26 gives agy a private HOME; that one is for the docs lens. gofmt/vet are clean.
4. Directory layout: **pass**. Nothing new is outside `internal/dream`, `internal/memory`, `internal/search` and `cmd/sh-agent`, and the tests sit beside their code.
5. Simplicity (DIY-except-security): **pass**. There is no new dependency, and `syscall.Flock` is the stdlib lock. No auth or crypto code is written here: copying the operator's OAuth token file into a private HOME only consumes the CLI's own login and implements no protocol. `GOOS=windows go build` now fails on `syscall.Flock` in `internal/memory`, but Windows is not a release target (`release_test.go:488` expects no windows asset), so this is not a finding.
6. Decoupling: **pass**. `search` owns the sqlite reads, and `dream` owns the CLI call, the lock and the state. The two flock helpers are small near-copies. That is acceptable at this size, but a future change to one must be applied to the other.

## Not examined
- The real Antigravity conversation schema (the `step_type` column and its values). The read-only query against `~/.gemini/antigravity-cli/conversations/*.db` was denied by the permission classifier, so finding 5 stays at low confidence.
- The real `agy` stream-json exchange and the real `claude --setting-sources local` behaviour. I did not re-run either CLI; the Outcome and BUG-119 claim probes on this box.
- The docs, spec, README and skill text changes in the range (`.agents/docs/**`, `README.md`, `AGENTS.md`, `skills/{init-project,project-memory,retro}/SKILL.md`). That is the docs lens's job, and I read only the spec rows cited above.
- `internal/dream/timer.go`'s `%%` escaping on a real systemd, and the macOS launchd path, which does not use the `%` rule. Only the unit test was read.
- Whether Antigravity's OAuth token is file-based on macOS. `agyHome` requires the file, and if the Mac keeps the login elsewhere, every Mac agy run fails with exit 2 "no Antigravity login".

## Outcome
1. medium, the Antigravity read-error fix had no test: fixed. `TestRunAgyReadErrorKeepsWatermark` sets up a garbage `conversation_summaries.db`; undoing the fix compiled and failed it, restored byte-identical. `TestCollectNamesAnUnreadableTranscript` covers `sessionText`'s error path. The limit is a variable, `transcriptLineMax`, and the scanner's own 64 KiB floor means the test line exceeds it.
2. medium, the state-write failure had no test: fixed. `writeState` is a variable, and `TestRunStateWriteFailureIsExit2` makes it fail; dropping the error compiled and failed the test, restored byte-identical.
3. low, the start-of-run cleanup could delete a live dry run's copies: fixed. Only copies older than an hour are removed (`TestRemoveCopiesKeepsFreshOnes`).
4. low, unreachable branches in `targets()`: fixed. Removed; `Run` refuses both and neither before anything is written (`TestRunUsage`).
5. low, the source of step types 14/15 was unrecorded: fixed (plan D8b, from two real conversation databases).
6. low, the symlink check missed `.agents` itself: fixed. Every component from the root down is checked (`TestCollectDoesNotFollowALinkedAgentsDir`).
7. low, `--setting-sources local` may lose settings-based auth: left, documented. The spec's Known limit names it; a CLI login works, and a machine using `apiKeyHelper` sees the failure logged per project.
