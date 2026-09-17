---
owner: security
date: 2026-09-17
---
# Security review — 012b4ca-worktree-final   (2026-09-17)

Range reviewed: 012b4ca..working tree (uncommitted): `git diff 012b4ca -- internal/harness/report.go internal/harness/report_test.go .agents/docs/deploy/local.md .agents/docs/specs/sh-agent/delivery.md`. The round-3 delta was isolated by diffing against the blobs the previous review pinned (report.go 9bb378e5, report_test.go 29c794b1, both kept in `.agents/scratch/secprobe-worktree/`).

**Blobs pinned** (`git hash-object`, the same at start and at end):
- `internal/harness/report.go` c97861a80070ace86273491039f95ffa749b5f6e
- `internal/harness/report_test.go` 5f8674a83467687e01d3f7fe76c43c522aed63a0
- `.agents/docs/deploy/local.md` d975776fb62186221ceb7e7343e6a8d61a454f8b
- `.agents/docs/specs/sh-agent/delivery.md` 07a2126ccb58fa8acd6d39e07b7b10927856926a
- Read for context: the prior review at 4ac15a924e681722e8cca4d55f74e43d7909b7cc, including its Outcome; plan D13 and D14; the `ReportRun` wiring in `cmd/sh-agent/main.go:242-245`.

**Probes.** Everything is under `.agents/scratch/secprobe-final/`:
- `harness/` is a copy of the working-tree package (report.go c97861a8, report_test.go 5f8674a8); its report tests pass (`copy-report-tests.log`).
- `r2/` is the round-2 package (report.go 9bb378e5), so every new probe runs before and after round 3 (`probe-r2.log`, `probe-current.log`; verdict tables `w-r2.txt`, `w-cur.txt`).
- Probe files:
  - `secprobe_final_test.go` is new. It is ASCII only, with every escape built at run time and byte-checked.
  - `secprobe_round3_test.go` holds the previous review's probes, re-run unchanged against the current code.
- **Fuzz:** `fuzz-label.log`, `fuzz-redact.log`, and the must-fail and must-pass checks `fuzz-label-mustfail-r2.log` and `fuzz-label-mustpass.log`.
- **govulncheck:** `govulncheck.log`.
- **Isolation:**
  - `TMPDIR`, `GOTMPDIR`, `GOCACHE` (fuzz) and `GIT_CEILING_DIRECTORIES` point inside the probe directory.
  - `GIT_CONFIG_GLOBAL` is an empty file, with `GIT_CONFIG_NOSYSTEM=1` and `GIT_SSH_COMMAND=false`.
  - Every remote is a local `file://` bare repository, and every value is fake.
  - `sh-agent report` was not run in any form.

## Verdict: FINDINGS (low only)

Every fix the Outcome claims holds under probe. For the code fixes (2b, 3, 4), each probed shape leaked at the round-2 blob and is caught now. Fix 6 is a test, and its mutant fails it. Fix 7 shows on both the dry run and a real run.
- **Regressions:** no fix introduced a leak. One fix introduced an over-redaction that fails safe (finding 3).
- **Worst remaining problem:** a written claim. `local.md` says a key block "stays open until a footer that actually names that header's label". The code applies the label test only on the header's own line, and only as a substring match. So two cases still close the block above the body, and the body is sent (finding 1):
  - a footer mention on the next line;
  - any private-key footer mention beside a PKCS#8 `-----BEGIN PRIVATE KEY-----` header.
- **Against D14:** all five findings are low and inside D14's accepted limits (redaction is a safety net), and none is outside them. In findings 1, 2 and 5 the defect is the written claim, not the net.

## Findings

- **[low / confidence: high]** `.agents/docs/deploy/local.md` § Report back ("A block opened on a line that merely mentions a footer … stays open until a footer that actually names that header's label, not the mention") against `internal/harness/report.go:112-118` (`blockEnd`). **The sentence claims more than the code does.**
  - **The code comment is accurate:** "On the line its header opened, a block closes only at a footer naming that header's label".
  - **The doc loses two limits:** it drops "on the line its header opened", and it does not say the label test is a substring match.
  - **Probed LEAKs** (`TestFinalBlockLabel`, a `proposals/*.md`). Each sent `FAKEBODY` while redactions were recorded:
    - **A footer mention on the line after the header** (`note: the body ends at -----END`). The label test does not run there (`j == opened`), so any footer-shaped text closes the block.
    - **A PKCS#8 header with another footer mentioned.** `-----BEGIN PRIVATE KEY-----` has the label `PRIVATE KEY`, and every private-key footer contains it, so a mention of `-----END RSA PRIVATE KEY-----` on the header's line closes the block.
    - **A mention that names the label itself** (`up to -----END OPENSSH PRIVATE KEY----- below`).
    - **A header whose label does not parse** (two spaces after `BEGIN`). It keeps the old any-footer rule, because the test requires `label != ""`.
  - **Not a regression:** all four leak identically at 9bb378e5. The fix narrowed the trigger; it did not create these.
  - **D14:** inside. It needs a key pasted raw into a proposal, plus a footer mention. The defect is the written claim.
  - **Remediation (no code):** write the code comment's scope into `local.md`:
    - on the header's own line, a footer mention closes the block only if it contains the header's label, and for `-----BEGIN PRIVATE KEY-----` any private-key footer does;
    - on any later line, any footer-shaped text closes it.
- **[low / confidence: high]** `.agents/docs/deploy/local.md` § Report back ("a line with a leading BOM, `NaN`, or a value shadowed by a duplicate key is checked only as raw text") and the `decoded` comment at `report.go:181-186`. **The list reads as closed, and it is not.** The code drops any value that fails to decode, whole, together with everything after it on the line.
  - **Probed LEAKs beyond the three stated** (`TestFinalDecoded`, JSONL, a needle behind a JSON unicode escape, backslash-u0041). Each of these was sent:
    - a leading non-BOM byte;
    - a log timestamp before the record;
    - a raw tab in another string of the same record;
    - a truncated record;
    - a record after a non-JSON byte;
    - `Infinity`;
    - nesting deeper than 10000.
  - **Wording:** the duplicate-key case is neither "checked only as raw text" nor "what cannot be decoded". The line decodes, and the shadowed value is discarded.
  - **Not a regression:** each of these leaks at 9bb378e5 too. The decoder swap turned six probed shapes from LEAK to CAUGHT, and none the other way.
  - **D14:** inside. It needs a needle hidden behind an escape on a malformed line.
  - **Remediation (no code):** reword the limit. A needle behind a JSON escape is missed in any value that does not decode, and in everything after it on that line. Name the probed cases: leading text or a BOM, `NaN` or `Infinity`, a raw control character, a truncated record, nesting deeper than 10000. A value shadowed by a duplicate key is never seen.
- **[low / confidence: high]** `internal/harness/report.go:63-74` (`blockLabel`, `labelAt`). **Introduced by round 3. It fails safe: an over-redaction, not a leak.**
  - **Mechanism:** the label is read greedily over `[A-Z0-9 ]`. A dashless header followed on its line by a space and capitals takes them into its label.
  - **Probed:** `{"q":"BEGIN PRIVATE KEY MIIEvFAKEbody END PRIVATE KEY"}`, a flattened key with its dashes stripped.
    - **Now:** the label becomes `PRIVATE KEY MIIE`, so the key's own footer on that line no longer matches. The block runs to end of file, and the two clean records after it are dropped.
    - **At 9bb378e5:** the same input redacted one line and kept both records.
  - **Visible:** the dry run and the real run both print `unterminated: everything after it in this file was dropped` for this line.
  - **D14:** inside ("over-redacts on purpose").
  - **Remediation:** none required. If wanted, a header with no closing dashes could stop its label at `PRIVATE KEY`.
- **[low / confidence: high]** `internal/harness/report.go:62` (`blockClose`, second alternative `END[A-Z ]*PRIVATE KEY`) against `local.md` "two concatenated keys kept open". **Inherited.**
  - **Mechanism:** the alternative is greedy over spaces and capitals.
  - **Failing shape:** a joint line with the dashes stripped and a space between the keys (`END RSA PRIVATE KEY BEGIN OPENSSH PRIVATE KEY`). Probed: the second key's body is sent.
    - The footer match swallows the second header.
    - The concatenation check then finds no open after it.
  - **Controls, caught:**
    - a dashed joint;
    - a dashed joint carrying a footer mention (round 3 fixed this one; it leaked at 9bb378e5);
    - a non-space rewritten joint (`~~`).
  - **Not a regression:** the same at 9bb378e5. `blockOpen` and `blockClose` are unchanged since 012b4ca.
  - **D14:** inside. It needs two keys pasted with their dashes stripped and the newline replaced by a space.
  - **Remediation, either of:**
    - qualify the claim: "unless the dashes were stripped and a space joins them". This fits D14;
    - make the alternative lazy (`END[A-Z ]*?PRIVATE KEY`, which RE2 supports), with this probe as a must-fail case.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/delivery.md:47` ("bearer and authorization headers"). **Authorization headers written as a call, with a non-Bearer scheme, are sent.**
  - **Probed** (`TestFinalAuthorization`, `proposals/*.md`). Each was sent with 0 redactions:
    - `req.Header.Set("Authorization", "Basic FAKEVALUE")`;
    - `headers["Authorization"] = "Basic FAKEVALUE"`;
    - `Authorization = "Basic FAKEVALUE"`.
  - **Controls, caught:** the colon forms, quoted or not; `Proxy-Authorization:`; JSONL.
  - **Scope:** `Bearer` values are caught by the `bearer` needle whatever the form.
  - **The claims:**
    - `local.md` names the shape exactly (`Authorization: <anything>`, and a quoted name before the colon), so it is accurate.
    - The spec's summary covers more than that.
  - **D14:** inside (item 8: unlisted shapes are sent).
  - **Remediation:** none required. Optionally, the spec could say "authorization headers written `Authorization: …`".

## Checked

- **Outcome claims, verified by probe** (verdict at 9bb378e5 → now):
  - **2b (a header line that mentions a footer):** holds.
    - **Probed, LEAK → CAUGHT, clean lines after kept:**
      - the text form;
      - the JSONL `note` form;
      - a dashless header with a mention of a shorter label;
      - a concatenated joint line with a mention.
    - **Mutant:** `r3-mutants/2b-label.log` fails the intended case.
    - **Fuzz:** `FuzzFinalLabel` ran 45 s and 4.99 M execs with no failure.
      - **Property:** a body between a labelled header and its own footer is never sent. This holds whatever else the header's line carries, when that line has no earlier `BEGIN` and no second copy of the label, and the label does not contain `END`.
      - **Must-fail:** the property's seeds fail at 9bb378e5 (exit 1, seed#0 "body sent").
      - **Must-pass:** they pass at the current blob.
  - **3 (a quoted `Authorization` name):** holds.
    - **Probed, LEAK → CAUGHT:** `"Authorization": "Basic …"`, `'Authorization': 'Basic …'` and `{"Authorization":"Basic …"}`.
    - **JSONL:** the wrappings are caught.
  - **4 (the decoded check):** holds.
    - **Probed, LEAK → CAUGHT:** `1e400` beside an escaped needle, two records, three records, and a trailing byte.
    - **Also `state.json`, LEAK → CAUGHT:** two documents, and `1e400`.
    - **No regression:** no decoded probe went from CAUGHT to LEAK.
    - **Why each case needs its own half of the fix:** keeping the decoder loop without `UseNumber` loses the whole `1e400` value, and `UseNumber` with a single `Decode` misses the second record. So the two test cases are independent must-fails, not one.
  - **6 (the decoded branch of `redactDocument`):** holds.
    - **Fixture bytes:** they carry a real backslash before `u0041` (`report_test.go:113,118,119,192`). No control byte or BOM was found in either Go file.
    - **Mutant:** its log fails the intended case.
    - **Probed, all caught:** an escaped key name, an escaped `ghp_`, and pretty-printed JSON with an escape.
  - **7 (the dry run shows dropped records):** holds on both paths.
    - **Output:** the dry run and a real run against a local bare repository both printed `feedback.jsonl:2 (-----BEGIN) 3 lines, unterminated: everything after it in this file was dropped`.
    - **What was pushed:** the commit carries no `FAKE` and no dropped record, and `report.json` has `"unterminated": true`.
    - **Why both paths print it:** the list is printed before the dry-run branch (read).
  - **1 (push claims, no code change):** the claims now match the behaviour. The previous push probes were re-run:
    - **Land in the feedback repository:** origin's pushurl, a second url, and origin's pushurl set through the environment.
    - **Land in the other repository:** all six URL-rewrite and named-remote variants.
    - **Claims checked against this:** the `FeedbackURL` and `commitAndPush` comments, `local.md` and `delivery.md`.
    - **Refusal:** it still never echoes the origin.
- **Stated limits, re-probed:** each leaks as `local.md` says:
  - a BOM, `NaN`, and a duplicate key;
  - escapes in text files;
  - the six next-line shapes.
- **Regression fuzz:** `FuzzProbeRedact` (the previous review's) ran 45 s and 1.23 M execs against the current code. There was no panic, and no needle survived on an output line or in a document.
  - **Blind spot:** by construction it cannot see a body with no needle. `FuzzFinalLabel` covers that for the label rule.
- **Secrets:** the added test lines carry only fake values (`FAKEVALUE`, `FAKEBODY`, `FAKE`).
- **Injection and robustness:**
  - **git:** no new git argv or exec.
  - **Slicing:** `labelAt` slices from an index that the preceding `blockOpen` match guarantees, so it cannot panic. Neither fuzzer panicked.
  - **Regexes:** RE2, so linear time.
  - **Loop:** `blockClose` has no zero-length match, so the new `continue` always advances.
- **Authz/authn:** unchanged since the prior review. The push destination is as D14 decided.
- **Crypto:** nothing is implemented or composed in the range.
- **Dependencies:** no imports were added in the range; three stdlib imports were removed. `govulncheck ./...` (go1.27.1, `GOTOOLCHAIN=local`) exited 0 with "No vulnerabilities found." There is no npm ecosystem, so `npm audit` does not apply.
- **GCP surface:** nothing touched.
- **Probe hygiene:**
  - **Real repository at the end:** HEAD 012b4ca, no commits ahead of upstream, tag `v1.0.0` only.
  - **Status:** it shows no path outside the other agents' existing edits and this report, and scratch is gitignored.
  - **Git repositories:** no `git init` ran in the shell. Every probe repository was created inside `t.TempDir()` under the probe `tmp/`.

## Not examined

- **Other agents' reviews:** `2026-09-17-review-012b4ca-worktree-code-docs.md` and `2026-09-17-review-015cf2d-012b4ca-code.md` were not read.
- **Uncommitted edits outside the requested diff:** not read. They are `README.md`, `REQUIREMENTS.md`, `context/AGENTS.md`, `remote.md`, `guide.md`, `DECISIONS.md`, BUG-053 to BUG-060, the round-2 fixes testing log, and the prior review file's own edits. Of the plan, only D13 and D14 were read.
- **Mutant sources:** the `r3-mutants/*` sources were not diffed; the diff command found no top-level Go files at the path it checked. Only their logs were read, and each fails its intended case. For 2b, 3, 4 and 6, the must-fail evidence cited here is this review's own before-and-after probe against 9bb378e5.
- **`FuzzFinalLabel` coverage:** the share of inputs skipped by its precondition was not measured. Its property excludes, by design, the residuals named in finding 1: a second copy of the label on the header line, and labels containing `END`.
- **Push:**
  - probed only with `file://` URLs;
  - not exercised: a real GitHub push, the operator's real global git config, macOS.
- **`report.json`:** fields beyond `redactions` were not re-probed. The `state.json` fields that reach `report.json` still leak, as D14 (5) accepted.
- **Unchanged code:** `Collect`'s handling of the proposals directory (symlinks, non-`.md` files), `BuildRecord`, `ResolveFeedback`, `runDir` and `hostDir` were not re-read. None changed in the round-3 delta.

## Outcome

Every finding is inside D14's accepted guarantees, and no fix leaked or regressed. Where a written claim said more than the code, the claim was narrowed; no code changed except one comment.

- Finding 1 (`local.md` overstates the key-block label rule): fixed — `local.md` states the label check at its real scope (the header's own line, as a substring) and names the early-close shapes as a limit.
- Finding 2 (the decoded-check limits read as a closed list): fixed — `local.md` and the `decoded` comment in `report.go` say a value that does not decode hides itself and everything after it on the line, list the probed cases as examples, and say a value shadowed by a duplicate key is never seen.
- Finding 3 (a flattened, dash-stripped key redacts to the end of the file): owned — fails safe, it is D14's deliberate over-redaction, and the dry run prints `unterminated` for it.
- Finding 4 (two dash-stripped keys joined by a space): fixed as a claim — `local.md` qualifies the concatenated-keys sentence and states the case as a limit; the lazy-regex alternative was not taken, per D14 (no chase for shapes).
- Finding 5 (`Authorization` set by a call): fixed as a claim — `delivery.md` names the colon form it catches; the call form is an unlisted shape, per D14 (8).
