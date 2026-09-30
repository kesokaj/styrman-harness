---
owner: reviewer
date: 2026-09-30
---
# Review — e2cffc1-a275b96-code   (2026-09-30)

Range reviewed: e2cffc1..a275b96 (source: `git diff e2cffc1..a275b96 -- internal/`; the round-2 Outcomes of `reviews/2026-09-30-review-f1ca525-e2cffc1-code.md` and `reviews/2026-09-30-security-f1ca525-e2cffc1.md`, checked against code and tests)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem: the security round-2 medium fix, "a reply's tool-call arguments are dropped", does not hold on real data. The filter drops a run only when it starts with `{"`. In a real protobuf payload, the byte right before a length-delimited string is its varint length, and for many lengths that byte is itself printable. `textRun` then joins it to the front of the run, so the run starts with, for example, `` `{"`` and passes the filter. I confirmed this on the probe copies already in `.agents/scratch/`, which include the gate18 conversation the Outcome cites. The test fixture cannot catch it, because it writes a fixed non-text length prefix whatever the payload's length. The other round-2 Outcome claims hold. There is one stale spec row.

## Findings
- **[medium / confidence: high]** `internal/search/dream.go:240` (`strings.HasPrefix(s, `{"`)`), with `textRun` at `:92`: Antigravity tool-call arguments still reach the other vendor's model.
  - **Mechanism.** A protobuf string field is `tag, varint length, bytes`. `textRun`'s first character class is `[\x20-\x7e\xc2-\xf4]`, and its later classes also admit `\x80-\xbf`. Two length ranges put a printable byte in front of `{"`:
    - 33 to 126 bytes: the one-byte varint length is printable ASCII. A length of 32 is a space, which `TrimSpace` removes.
    - 4096 to 16383 bytes: the second varint byte is `0x20`-`0x7e`. A first byte of `0xc2`-`0xf4` can even start the run, and `ToValidUTF8` then deletes it, leaving one printable character in front.

    The run then begins `<char>{"…`, fails the prefix check, and is sent whenever it has 3 or more spaces. That covers short shell commands (`CommandLine`) and whole written files of 4 to 16 KB (`CodeContent`), which are exactly what the finding was about.
  - **Evidence.** I ran a read-only scan (`mode=ro&immutable=1`, metadata only printed) with the same regex and filter over step types 14 and 15 of the conversation copies under `.agents/scratch/**/conversations/`:
    - `gate18-agy/…/15776d24….db`: a run starting `0x7e 0x7b 0x22` (`~{"`, a length of 126) holds `{"AbsolutePath":…`, a view_file argument. It is kept.
    - `agysmoke/…/e85bf990….db`: this is the database with "32 dropped", matching the Outcome's count. A run starting `` 0x60 `` (a length of 96) holds `{"Action":"list"…`. It is kept, and is counted among the Outcome's "21 prose runs".
  - **Test.** The fixture `blob()` (`internal/search/dream_test.go:119-126`) always writes `0x0a 0x8f 0x01` before each part. The byte `0x01` is non-text, so every fixture run starts cleanly at `{"`. This is not the arrangement a real payload has, so `TestAgyConversationText` passes with the leak present.
  - **Describing text now overclaims.** `specs/sh-agent/dream.md:39` says "a reply's tool-call arguments … are dropped". Plan D29 (`plans/2026-09-30-parity-memory-dream.md:342`) records the operator's "keep both, filtered" choice, and names the residual risk as only "an unknown field type". That choice was made on the premise that this filter works.
  - **Fix.**
    - Make `blob()` write the part's real varint length, and add a 33-126-byte argument case and a 4096-16383-byte argument case. Both must fail against the current code.
    - Then either drop a run whose `{"` appears within its first 2 characters after trimming, or, cleaner, walk the protobuf wire format (tag, varint and length are stdlib-simple) so that each string field is taken whole and its prefix is exact.
    - Correct spec :39 in the same change, and put the D29 premise back to the operator.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/dream.md:60`: the table row still says a killed run's copies are removed as "any leftover … at the start of the next". Since this range, `removeCopies` (`internal/dream/run.go:301-314`) removes only copies older than an hour.
  - Round-1 code finding 3 asked for the spec row to be corrected along with the comment. The comment was fixed; the row was not.
  - Fix: state the one-hour bound, and why (a dry run takes no lock), in that row.

### Round-2 Outcome claims checked against code and tests
Code review (f1ca525..e2cffc1) Outcome:
1. Antigravity read error test: **holds.**
   - `TestRunAgyReadErrorKeepsWatermark` feeds a garbage `conversation_summaries.db`. Putting back the `Dropped` append would give exit 0 and an advanced watermark, so the test can fail.
   - `TestCollectNamesAnUnreadableTranscript`: `bufio.Scanner`'s limit is max(80 KiB, 64 KiB cap), so a 100 KiB line gives `ErrTooLong` as the test claims.
2. State-write failure: **holds.** `writeState` is the only write path used by `saveState`. The test asserts exit 2 and "could not be saved".
3. Cleanup spares a live dry run: **holds in code.** `removeCopies` has one call site (`run.go:85`, under the lock). The dry run never calls it, and no copy outlives `cliTimeout` (10 min), far below the hour. The spec row is stale (finding 2).
4. Unreachable `targets()` branches: **holds.** They are removed. `Run` refuses both and neither before `realRun` or `targets` (`dream.go:63-68`), and `TestRunUsage` asserts exit 2 with nothing written.
5. Step types 14/15 source: **holds.** Plan D8b records it.
6. `.agents` link: **holds.** `linked` now `Lstat`s each component. `TestCollectDoesNotFollowALinkedAgentsDir` would fail under the old two-path check, because `Lstat(root/.agents/docs)` resolves the intermediate link. Files are also `Lstat`ed (`mdFiles`, `readFile`).
7. `--setting-sources local`: left and documented at `specs/sh-agent/dream.md:37`. Accepted.

Security review (f1ca525..e2cffc1) Outcome:
1. Tool-call arguments dropped: **does not hold** (finding 1).
2. `.agents` link: **holds** (as code item 6 above).
3. Remaining gaps:
   - Absolute-only PATH: **holds.** `TestInstallTimerKeepsOnlyAbsolutePath` fails without `absolutePath`.
   - SUCCESS required: **holds.** The real stream outputs in `.agents/scratch/sj-*.out` carry `"status":"SUCCESS"` and `"status":"ERROR"`, and the fake at `run_test.go:143` emits SUCCESS.
   - Explicit `toolPermission`: **present in code.** The sandbox claim was not re-run.

## Criteria
1. Dead code: **pass.** The dead `targets()` branches and the test lines that pinned them are gone, and the `errors` import was removed with them. `transcriptLineMax` and `writeState` are test seams with callers.
2. Flow tracing: **fail** (finding 1).
   - The rest traces cleanly. The usage refusal happens before the lock. The read-error path is collect → `fmt.Errorf` → project failure with the watermark kept. The state-write error becomes exit 2. `removeCopies` ignores an `Info()` error, which keeps the data.
   - The package-level test seams are safe, because no test in `internal/dream` uses `t.Parallel`.
3. Standards conformance: **fail** (low).
   - `specs/sh-agent/dream.md:39` restates a guarantee the code does not give (finding 1), and `:60` is stale (finding 2).
   - `removeCopies`' comment is now accurate.
4. Directory layout: **pass.** Changes stay in `internal/dream` and `internal/search`, with the tests beside the code.
5. Simplicity (DIY-except-security): **pass.** There is no new dependency, and no auth or crypto code. Walking the protobuf wire format for finding 1 is stdlib byte parsing, not a security primitive.
6. Decoupling: **pass.** `search` still owns the sqlite reads, and `dream` owns the run, the lock and the state.

## Not examined
- I did not run `go test`, `go vet`, or any mutation. The claims that mutations were "restored byte-identical" were judged only by reasoning about whether each test can fail, not by re-running them.
- The real `agy` runtime with the explicit `toolPermission: request-review` (the security Outcome's sandbox run), `claude` under the new flags, and macOS launchd.
- Protobuf field layouts beyond the 30 conversation copies in `.agents/scratch/`. Only step types 14 and 15 were scanned. JSON found deep inside long prose runs (for example `{"autoMemoryDire…` at offset 11582, and the dream's own `{"proposals"…` answers) was not classified.
- The trail files in the range: BUG-116..122, the INDEX rows, the change entry, and both verifier logs. The plan diff was read only at D8b and D29, and the spec only in the rows cited.

## Outcome
1. medium, the JSON-argument filter did not hold on real data: fixed. The scan is replaced by a protobuf wire-format decoder (plan D30); see `reviews/2026-09-30-security-e2cffc1-a275b96.md` § Outcome 1 for the tests and the 95-database measurement. The fixture's `blob()` now writes real framing.
2. low, the spec row on leftover copies: fixed (only copies older than an hour; a dry run takes no lock).
