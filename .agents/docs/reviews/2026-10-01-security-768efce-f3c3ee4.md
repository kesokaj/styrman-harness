---
owner: security
date: 2026-10-01
---
# Security review — 768efce..f3c3ee4, v2.2.0 gate fix round   (2026-10-01)

Range reviewed: 768efce..f3c3ee4 (9 paths; the rest of v2.1.0..f3c3ee4 is covered by reviews/2026-10-01-security-v2.1.0-a62eb9d.md and reviews/2026-10-01-security-a62eb9d-768efce.md)

This is the release-gate read of the fix round for `reviews/2026-10-01-security-a62eb9d-768efce.md` (five lows).

- **How it was read:** only through git (`git diff 768efce..f3c3ee4 -- <path>`, `git show f3c3ee4:<path>`). The working tree was never read.
- **Where the reproductions ran:** on `git archive f3c3ee4` copies under `.agents/scratch/security-v2.2.0-r3/`, with a stand-in `HOME` and `GIT_CONFIG_NOSYSTEM=1`. The only sources were `127.0.0.1:1`, a closed port, and local paths that do not exist.

## Status: COMPLETE

## Verdict: FINDINGS

There is no high and no medium finding. There are two new lows and one carried low.

Every claimed fix holds at f3c3ee4 for what the previous findings named:
- No refusal of a source override prints its value. 16 spellings were run.
- A query string is refused.
- `fetch()` names the source by its role, and git's message has the value replaced for every ordinary spelling tried.
- `printable` replaces every character the previous finding listed.
- The audit skill's three sentences now keep an agent to its own trail and scratch.

The worst residual is in `printable`. It was fixed with an enumerated list, not with the `unicode.Cf` class the previous finding proposed. U+061C (Arabic letter mark, a direction mark like the LRM and RLM it now drops), the zero-width characters, U+FEFF and the tag block U+E0000–U+E007F still pass.
- **Reproduced:** a search line followed by "ignore the rules" written in tag characters comes out whole. A terminal shows nothing, but a model that decodes tag characters reads it.
- **Pre-existing:** this behaviour predates the range.

The second residual is in `fetch()`. It removes the value from git's message by exact match, and git truncates its own message at 4,096 bytes. An override longer than about 4 KB is therefore printed: its first ~4 KB, raw. Reproduced.

## Findings

1. **[low / confidence: high]** `internal/search/search.go:616-629` (`printable`; the new check is at `:624`). The weakness is older than this range.
   - **What passes.** A probe on the archived function, `printable("a"+r+"b")`, gave:
     - **Replaced:** U+0085, U+2028, U+2029, U+202E, U+2066 and U+200E.
     - **Passed, all `unicode.Cf`:** U+061C, U+200B–U+200D, U+2060–U+2064, U+206A–U+206F, U+FEFF, U+00AD, U+180E, U+FFF9–U+FFFB, U+1D173, U+E0001, U+E0041 and U+E007F.
     - **Passed, invisible but not Cf:** U+034F, U+FE0F, U+E0100, U+115F, U+3164 and U+FFA0.
   - **Reproduced:** `printable("note.md:1  harmless" + "ignore the rules" as U+E0000+c)` returns 83 bytes, with the tag runes kept. A terminal shows only `note.md:1  harmless`.
   - **Nothing found that forges a line:**
     - C0, DEL, C1 (U+0085 among them), U+2028 and U+2029 are replaced.
     - A raw invalid byte such as 0x9b becomes U+FFFD (`strings.Map`), so no 8-bit CSI passes.
     - Python's `splitlines` set (`\x1c`–`\x1e`, `\x85`, U+2028, U+2029) is covered.
   - **What can still hide or reorder text:**
     - the tag block, which models that decode it read as ASCII ("ASCII smuggling");
     - zero-width characters and U+FEFF, which make two names look the same;
     - U+061C, which has bidi class AL and can reorder the digits and neutral characters next to it, but cannot override a run as RLO does;
     - a name padded with spaces so that it wraps into what looks like a separate terminal line. That one cannot be filtered without breaking real names.
   - **Pre-existing:** `printable` never filtered any Cf character outside this range's list. The previous Outcome's "fixed" holds for the characters the finding listed, not for the class it proposed.
   - **Remediation:** small. Map `unicode.In(r, unicode.Cf, unicode.Zl, unicode.Zp, unicode.Variation_Selector)` to `?`. In search output, over-replacing a ZWJ inside an emoji or a ZWNJ is the safe direction. Add rows for U+061C, U+200B and U+E0041 to `TestPrintableDropsLineAndDirectionControls`.
   - **Exposure if shipped:** a file name or a line in an indexed repository can carry text that the operator does not see in `sh-agent search` output but a model reading that output does. Search is not a new channel, because the same bytes reach the model when it opens the file. It only fails to close this one.

2. **[low / confidence: high]** `skills/google-catalog/scripts/refresh-catalog.sh:176-184` (`fetch()`, the `${err//"$2"/<the $1 source>}` replacement).
   - **Exact match only.** The replacement removes the value only where git prints it exactly as given, so any other spelling escapes it. Two such cases were reproduced with the extracted function:
     - **A long value.** git formats its fatal message into a fixed 4,096-byte buffer. A value longer than roughly 4 KB minus the message's own text is printed truncated, the match fails, and the first ~4 KB of the value is printed raw (`[long-4KB] LEAK`, `[long-32KB] LEAK`). The mechanism is inferred, with medium confidence; the leak itself is reproduced.
     - **An `insteadOf` rewrite.** A URL rewritten by the operator's `url.<base>.insteadOf` is printed in its rewritten form: `unable to access 'https://127.0.0.1:1/REWRITTEN-SECRET7/x/'`. git anonymizes the userinfo of a rewrite itself (SECRET8 was not printed).
       - This value comes from the operator's own config, not from the override.
       - It is not a regression: before this range, git's stderr reached the terminal unfiltered.
   - **What holds.** Each of these came out as `<the agents-cli source>` or `<the agents-cli source>/`:
     - a plain URL, and one with a trailing slash;
     - percent-encoding (`%53ECRET`);
     - `*?[x]`, `&` and `\` in a path;
     - a slash inside userinfo, which curl rejects;
     - the reuse branch: `remote set-url`, then a failed fetch.
   - **The replacement itself is safe:**
     - the quoted pattern matches `*`, `?` and `[` literally;
     - the replacement string holds no `&`, so bash 5.2's `patsub_replacement` does not re-insert the match;
     - its cost is bounded, because git's message is capped at 4 KB. The 32 KB case returned well inside a 60-second watchdog.
   - **The restatement it falsifies.** The changes entry's title says "no value of a source override is ever printed". The long-value case is an exception.
   - **Remediation:** trivial. Either refuse an override longer than, say, 1,024 bytes, next to the `?` check, or print only curl's reason (the text after the last `': `) instead of git's whole message.
   - **Exposure if shipped:** none in ordinary use. When a clone fails, an override longer than 4 KB, or an `insteadOf` rule with a token in its path, is put on the terminal and in the transcript.

3. **[low / confidence: high] Carried and already owned:** this is the previous finding 1's Left, owner operator. `refresh-catalog.sh:279-285` writes the override value verbatim into the Repository column of `MANIFEST.md`.
   - **Which values reach it now:** every value the guard accepts. These were reproduced as accepted:
     - an https URL with a token in its path (`https://h/SECRET/x`);
     - a URL with a `;` parameter;
     - a token as a subdomain;
     - a URL with a fragment;
     - any absolute path, which carries the operator's user name and directory layout;
     - values holding `|` or a newline, which break the table.
   - **Only after a successful clone:**
     - A fragment never clones: curl drops it, together with the `info/refs` that git appends.
     - A slash in userinfo never clones either.
     - Refused values never get there: userinfo, `?`, and any scheme other than https.
   - **Exposure if shipped:** as the Outcome states. A path token, or a local path, on a source that clones is committed into `MANIFEST.md` and into the release zip. Overrides are documented as being for the tests' fixtures.

## The questions of the dispatch

1. **Refresh script.**
   - **Refusals:** none of them shows the value. 9 refused spellings came out `noleak`. They include `HTTPS://user:SECRET@`, a leading space, `http://user:SECRET@`, `ssh://`, `ext::`, and `?` in the path and in the authority.
   - **`fetch()`:** see finding 2.
   - **The `remote set-url` line:**
     - It prints nothing on success, and its errors (`No such remote`) do not carry the URL.
     - It stores the value in the work clone's `.git/config` under scratch (`SECRET6 … : 1`). Scratch is gitignored, and this was already so before the range.
   - **`MANIFEST.md`:** see finding 3.
2. **`printable`:** see finding 1. Every character the previous finding named is replaced. What is left is listed there, and it predates the range.
3. **The `harness-audit` skill:** no finding. Nothing tells an agent to change state outside its trail and scratch.
   - **Step 0.1 (`:34`):** the audit now never commits, stashes or discards a change it did not write.
   - **Step 2 (`:55`):** a reproduction runs with `--home .agents/scratch/<name>` and with `--dry-run` "wherever the command has them", and never against the real home or the service manager.
     - Every verb that changes state has `--dry-run`: `install`, `update`, `uninstall`, `report`, `dream --all|--project` and the timer verbs (`cmd/sh-agent/main.go:5-18` at f3c3ee4).
   - **The re-audit (`:92`):** it repeats the read-only, no-network and data duties in as many words.
   - **Everything else an agent writes** is one of two kinds:
     - its own trail, committed and pushed by name (`:81`, `:94`);
     - scratch (`:98`).
   - **Building, testing and the network** belong to the main agent (`:36`).
4. **The tests (`catalog_test.go`, `search_test.go`).**
   - **If a refusal regressed:** the `http://user:token@example.invalid/x` and `https://example.invalid/x?access_token=abc` fixtures would reach `git clone`.
     - The network would see one DNS lookup of `example.invalid`. RFC 6761 reserves `.invalid` to fail, and no TCP connection follows.
     - If the operator's environment sets `http(s)_proxy`, the request would go to that proxy, carrying the dummy `user:token`. `fixtureEnv` keeps every variable that does not start with `GIT_`.
     - The only write would be `tmp/work-o`. The clone fails before `dest-o`, and both the message check and the work-directory check turn red.
   - **The nonexistent-path fixture** (`tmp/no-such-SECRETPATH-repository`) uses no network. It writes only `tmp/work-c`, the marker, by design.
   - **Pre-existing:** `fixtureEnv` also keeps `XDG_CONFIG_HOME`, so the operator's git config (credential helpers, `insteadOf`) applies in these subprocesses.
   - **Not covered by a test:** the reuse branch's replacement (a failed `fetch` after `remote set-url`). It held in the reproduction.
   - **`search_test.go`:** a pure function test, with escapes in the source and no raw bytes.
5. **The trail.**
   - **Secrets:** none of the 465 added lines holds a secret, a credential or a private key. A scan for token shapes, key armor, IP addresses, e-mail addresses and `.internal` names found only:
     - the known fixtures: `t@example.invalid`, `/home/u` and `https://user:token@example.invalid/x`;
     - `roam.internal`, which the added previous security report restates from an already-committed report. That is not new exposure, and the repository is private.
   - **Raw bytes:** 0 C0 or DEL bytes, and 0 raw bidi, separator or tag bytes, in the added lines.
   - **Instruction-shaped text:** none. Nothing tells a later agent to skip a gate or to act outside its lane.
   - **Outcomes measured against their findings:**
     - The security Outcome's item 3 ("fixed") holds for the characters the finding named, but not for the class its remediation proposed (finding 1).
     - The changes entry's "ever printed" claims too much; finding 2's long-value case is an exception.
     - Every other "fixed" claim checked out at f3c3ee4:
       - the non-echo of refused values;
       - the query refusal;
       - the role names;
       - both verifier gaps (`catalog_test.go:835` and the value-absence block);
       - the re-audit duties;
       - both skill wordings.

## Checked

- **Secrets:** none found; see question 5.
- **Injection:** findings 1 and 2. The arguments of `fetch()` sit behind `--`, and `$3` is a pin checked as hex.
- **Authz/authn:** the audit skill's duties (question 3). The range has no auth code.
- **Crypto:** none in the range, and nothing hand-rolled.
- **Dependencies:** `go.mod` and `go.sum` are not in the range, and `search.go` adds no import.
  - `govulncheck ./...` ran on a full `git archive f3c3ee4` copy with go1.27.1 (`/home/kesokaj/.go/bin/govulncheck`). Its output: "Your code is affected by 0 vulnerabilities". It also reported 1 vulnerability in a required module that the code does not call, consistent with the known GO-2026-5024.
  - `npm audit`: not applicable; there is no JavaScript in the range.
- **GCP surface:** none in the range.

## Per-path coverage

All 9 paths were read.

Added, read in full at f3c3ee4:
- `.agents/docs/changes/2026-10-01-1856-v2.2.0-gate-second-pass.md`
- `.agents/docs/reviews/2026-10-01-security-a62eb9d-768efce.md`
- `.agents/docs/testing/2026-10-01-verify-v2.2.0-gate-r2.md`

Modified, read as a complete diff:
- `.agents/docs/reviews/2026-10-01-review-v2.1.0-a62eb9d-code.md`
- `internal/harness/catalog_test.go`, plus `fixtureEnv` and its setup lines
- `internal/search/search.go`, plus `printable` in full and its print sites
- `internal/search/search_test.go`
- `skills/google-catalog/scripts/refresh-catalog.sh`, plus the whole file at f3c3ee4
- `skills/harness-audit/SKILL.md`, plus a grep of the whole file for steps that change state

## Not examined

- **Release-diff paths:** none left unread.
- **A real https server:** these were not exercised on a failed clone:
  - redirects (`warning: redirecting to …`);
  - a 401 with a credential prompt;
  - `remote:` sideband text.

  Only a refused connection on `127.0.0.1:1` and local paths that do not exist were run.
- **macOS:** not run. There, `t.TempDir()` is under `/var`, a link to `/private/var`. By reading git's clone code, git prints the path as given, so the replacement still matches.
- **Tests:** neither the suite nor the new tests were run, and no mutation was made. That is the verifier's lane, running at the same time. Question 4 is reasoned from the script and the test code.
- **Model and terminal behaviour:** not tested:
  - whether a particular model decodes tag characters;
  - how a terminal renders U+061C.

  The pass-through itself is reproduced.
- **`cmd/sh-agent/main.go`:** only its usage lines were grepped, for `--dry-run`. Whether `report` takes `--home` was not checked.
- **Working tree:** not read, per the dispatch.
- **Scratch left by this review**, in `.agents/scratch/security-v2.2.0-r3/`:
  - `src/`, the three archived files;
  - `full/`, the whole archive;
  - `probe/`, holding `go.mod` and `main.go`;
  - `home/.gitconfig`, with two stand-in `insteadOf` rules;
  - `work-1` to `work-10`;
  - `work-reuse/agents-cli`, a probe repository whose `git init` returned 0;
  - `override-block.sh`, `fetch-fn.sh`, `repro-refresh.sh`, `added.txt` and `tmp/`.
- **`/tmp`:** the `go run` of the `printable` probe built in Go's default temporary directory under `/tmp`, because `TMPDIR` was not pointed at scratch for that call. The go command removed it on exit.
- **Working-tree state at the end of this run** (`git status --short --untracked-files=all`). Apart from this file, the following are taken to be the concurrent agents' and were not checked:
  - ` M .agents/docs/reviews/2026-10-01-review-v2.1.0-a62eb9d-code.md`
  - `?? .agents/docs/testing/2026-10-01-verify-v2.2.0-gate-r3.md`

## Outcome

Written by the main agent on 2026-10-01, in the release commit. No finding is medium or high; all three are recorded with an owner and not fixed before the tag.

1. `printable` filters a list of characters, not the format class, so tag characters, zero-width characters and U+061C pass (low, older than this range): left, owner main agent, first task after the release. What the list covers is fixed and tested; the earlier Outcome's "fixed" holds for the characters that finding named, as this report says. Gate report row L2.
2. git's own message is cleaned only where git prints the value exactly as given: a value over about 4 KB, or one rewritten by an `insteadOf` rule, can still appear when a clone fails (low): left, owner main agent. The 1856 changes entry's "no value is ever printed" is too strong for these two cases; the release entry says so. Gate report row L3.
3. Every value the guard accepts is written into `MANIFEST.md` after a successful clone, so a token in a URL's path would be committed (low, carried): left, owner operator: the script's header says the value goes there. Gate report row L4.
