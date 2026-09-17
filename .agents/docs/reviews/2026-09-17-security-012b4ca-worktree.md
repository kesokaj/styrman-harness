---
owner: security
date: 2026-09-17
---
# Security review — 012b4ca-worktree   (2026-09-17)

Range reviewed: 012b4ca..working tree (uncommitted): `git diff 012b4ca -- internal/harness/report.go internal/harness/report_test.go .agents/docs/deploy/local.md .agents/docs/specs/sh-agent/delivery.md README.md`, and its blast radius (`repo.go` `gitOut`/`IsCheckout`, `state.go` `LoadState`, the `report` wiring in `cmd/sh-agent/main.go`).

**Blobs pinned** (`git hash-object` of the working-tree files, first at the start and again at the end):
- `internal/harness/report.go` 9bb378e5005e99c0fea1eb8cfb79fff9ac0f11e7 — the same at start and end.
- `internal/harness/report_test.go` 06b295af2041085d8f1e65389fcb5c0f280d503d — the blob probed.
  - **It changed mid-review** to 29c794b18dea852bb6d9f464d91fc184dc76ac05 through another agent's edit. The edit touched one test comment and the fixture of "needle behind a JSON unicode escape".
  - That diff was read in full, and both blobs were run against the mutants (finding 6).
  - It was still 29c794b1 at the end.
- `.agents/docs/deploy/local.md` 820a9a070b641246be03fc2c86c7568847c84ce8, `.agents/docs/specs/sh-agent/delivery.md` c4d0cf16643a5516a12692efd1a8c5c18527f66c, `README.md` 2e282ef77a484bdb33360730dd7f2b32682bf2cb — the same at start and end.
- Read for blast radius, both unchanged since 012b4ca: `repo.go` 653ab38c47d21108e984c607a99251dbc2489d4e, `state.go` 1cd44a807791e6c3721ec482a270c1d88d8e4c7d.
- The prior review, `2026-09-17-security-015cf2d-012b4ca.md`, was read at 046e282d9126d4d5b6af54ab762c80f05ed68755.

**Probes.** Everything is under `.agents/scratch/secprobe-worktree/`:
- **Copy:** `harness/` is a copy of the working-tree package. It hashes as pinned, and its 12 report tests pass (`copy-report-tests.log`).
- **Base:** `base-012b4ca/` was taken with `git show 012b4ca:`, to separate what this change introduced from what it inherited (`probe-base.log`).
- **Mutants:** `mut1/` and `mut2/` (`mut*-*.log`).
- **Probes and logs:** the probe file is `harness/secprobe_test.go`. Results are in `probe-run1.log`, `fuzz-run.log`, `guide-list.log` and `govulncheck.log`.
- **Isolation:**
  - `TMPDIR`, `GOTMPDIR`, `GOCACHE` for the fuzz, and `GIT_CEILING_DIRECTORIES` all pointed inside `tmp/` there.
  - `GIT_CONFIG_GLOBAL` pointed at an empty file, with `GIT_CONFIG_NOSYSTEM=1`.
  - `GIT_SSH_COMMAND=false` means no probe could reach a network remote.
- **Remotes and values:** every probe remote is a local `file://` bare repository, and every value is fake. `sh-agent report` was not run in any form.

## Verdict: FINDINGS

The worst gap is the push target, the third round on that mechanism.
- **What the fix closes, probed:** `remote.origin.pushurl`, a second `origin` url, and a `GIT_CONFIG_*` origin pushurl now all land only in the feedback repository.
- **What still redirects:** git applies `url.<base>.insteadOf` and `pushInsteadOf` to a URL argument. It also resolves a `[remote "<that URL>"]` section by the URL used as a remote's name.
- **Probed:** six variants, in the checkout's own config and in the environment. Each time `SendReport` returned success, the feedback repository got 0 commits and the other repository got 1.
- **The exact constant:** checked with `git ls-remote --get-url` (no network). A remote named `git@github.com:kesokaj/styrman-harness-feedback.git`, or an `insteadOf` on it, resolves to the other URL.
- **What it contradicts:** the `FeedbackURL` comment, `local.md`, `delivery.md` ("never to a remote") and D13 (3), which strikes `GIT_CONFIG_*` as moot.

**Redaction.** The round-2 fixes for findings 2 and 3 hold for every shape probed, and a 45 s fuzz found no needle surviving on any output line. Still sent:
- a secret on a later line than the needle that located it;
- a key body after a header line that also mentions a footer (round 2's "second trigger", which the Outcome's "fixed" did not cover);
- a quoted `"Authorization": "Basic …"`, which the new `authorization` needle misses;
- escape-hidden needles on JSONL lines that do not decode as exactly one JSON value; one of these (a number beyond float64) regressed in this change;
- `state.json` values copied into `report.json`.

**Coverage.** `redactDocument`'s decoded branch has no must-fail test.

**Severity.** All eight findings are low: each needs a deliberate local git setting, or a record or proposal shape the skills do not normally write.
- **Introduced by this change:** the float64 regression and the `authorization` gap. The pinned test blob also carried the mis-encoded fixture, since fixed.
- **Inherited:** the rest, from 012b4ca or earlier.

**Recommendation.** Both mechanisms, the push target and the redaction needles, are on their third round. So the remediations below lead with narrowing the written claims. Ask the operator what each mechanism must guarantee before more code goes in (Definition of Done, step 2).

## Findings

- **[low / confidence: high]** `internal/harness/report.go:16-22` (`FeedbackURL`) and `:583-586` (`commitAndPush`); `local.md` § Report back ("Where it goes"); `delivery.md:47`; D13 (3). **The push can still be redirected to another repository, while `report` reports success.**
  - **Mechanism:**
    - `git push <url>` rewrites the URL through `url.<base>.insteadOf` and `url.<base>.pushInsteadOf`.
    - Before that, it looks up a configured remote whose NAME is the URL string, using that remote's `url` or `pushurl`. A subsection name may contain `/` and `:`; git refuses only a leading `/`, with the warning "config remote shorthand cannot begin with '/'" (probed).
    - `verifyFeedback` reads only `remote.origin.url`, "as configured".
  - **Probed** (`TestProbePushTarget`, git 2.47.3, `file://` URLs):
    - **Redirected** (`err=nil feedback=0 elsewhere=1`), in checkout-local config:
      - `url.<E>.pushInsteadOf=<O>`
      - `url.<E>.insteadOf=<O>`
      - `remote.<O>.url=<E>`
      - `remote.<O>.pushurl=<E>` with no url
    - **Redirected** (`err=nil feedback=0 elsewhere=1`), through the environment with `GIT_CONFIG_COUNT/KEY_0/VALUE_0`:
      - `url.<E>.pushInsteadOf=<O>`
      - `remote.<O>.pushurl=<E>`
    - **Controls, not redirected** (`feedback=1 elsewhere=0`): none set; `remote.origin.pushurl`; a second `origin` url; an env `remote.origin.pushurl`.
    - **The real constant string:** with `git ls-remote --get-url`, which resolves the URL without connecting, both `remote.<constant>.url` and `url.<x>.insteadOf=<constant>` resolved to the other URL.
  - **Contradicts:**
    - `FeedbackURL`: "the only place it pushes to … no pushurl or second url in a checkout's config can redirect it". A `pushurl` under a remote named by the URL is such a pushurl.
    - `commitAndPush`: "a URL argument only itself".
    - `local.md`: "no `pushurl` or second `url` in the checkout's config can redirect it".
    - `delivery.md`: "pushes to the feedback URL itself, never to a remote".
    - D13 (3): "`GIT_CONFIG_*` injection: struck — … (2) makes a pushurl moot".
    - The prior review's "Residual, not raised" covered only a global `pushInsteadOf` that keeps the repository path. These variants point at another repository.
  - **Why low:**
    - `.git/config` never comes from the remote.
    - Nothing sets these keys or exports `GIT_CONFIG_*` by accident.
    - Whoever can set them already controls this user's git.
    - The commit itself still lands only in the verified checkout.
  - **Remediation.** This is the third round on this mechanism, so ask what it must guarantee.
    - **Shortest:** narrow the four claims to what holds, and restore D13 (3)'s environment clause to match.
      - `origin`'s own `pushurl` and extra `url` values cannot redirect.
      - URL rewrites, and a remote named after the URL, can, as for any `git push`.
    - **Code alternative** (stdlib, via `gitOut`). Before the push, refuse when either holds:
      - `git config --get-regexp '^url\..*insteadof$'` prints a value that is a prefix of the URL;
      - `git config --get-regexp '^remote\..*\.(url|pushurl)$'` prints a key whose subsection equals the URL.

      Then add the six probed variants as must-fail cases to `TestSendReportPushesOnlyToFeedbackURL`.
- **[low / confidence: high]** `internal/harness/report.go:121-157` (`redactLines`), `:89-109` (`blockEnd`) — **a needle fires and a redaction is recorded, but the secret it located is sent, because the secret is not on the needle's line.** Inherited: base 012b4ca sends the same inputs (`probe-base.log`).
  - **(a) Value on the following lines** (`TestProbeValueOffTheNeedleLine`, in a `proposals/*.md`). Each recorded 1 redaction and sent `FAKEVALUE`:
    - a YAML block scalar (`password: |`);
    - YAML folded (`api_key: >-`);
    - a YAML plain continuation (`token:` then an indented value);
    - an HCL heredoc (`secret = <<EOT`);
    - a shell continuation (`export GITHUB_TOKEN=\`);
    - pretty-printed JSON (`"password":` then the value). In a JSONL file this happens whenever a record spans lines.
  - **(b) Footer mention on the header's own line.** `blockEnd` closes the block at any `-----END` or `END…PRIVATE KEY` after the header on the same line, so a key body on the following lines stays outside the block.
    - **Probed:**
      - text `pasted -----BEGIN OPENSSH PRIVATE KEY----- up to -----END below` followed by body lines;
      - a JSONL record whose `note` mentions `-----END`.

      Both sent `FAKEBODY` with 2 redactions recorded.
    - **Round-2 history:** this is round 2 finding 2's "second trigger". The Outcome marks finding 2 fixed "by deleting the mechanism", but the trigger lives in `blockEnd`, which did not change.
  - **Contradicts:**
    - `delivery.md`, the design's own test: "Any span shorter than a structural unit leaks while reporting success". A YAML key's unit includes its block.
    - `local.md`: "never just the match".
  - **Why low:** it needs a proposal quoting multi-line config, or a key pasted after a line that mentions its footer. JSONL records written correctly are one line each.
  - **Remediation:**
    - **Shortest:** state both limits in `local.md`, beside "Anything that matches no shape IS sent".
    - **Code for (b):** close a block on its own header line only at a footer carrying the header's label, and add the probe as a must-fail case.
    - **Code for (a):** this needs a notion of continuation (a trailing `|`, `>`, `\` or `<<`, or an empty value). That is a design question for the operator, not a patch.
- **[low / confidence: high]** `internal/harness/report.go:50` — **the `authorization` needle added in this change misses a quoted header name.**
  - **Mechanism:** `(?i)authorization\s*:\s*\S` needs the colon directly after the name. The secret-assignment needle allows a closing quote (`["']?`); this one does not.
  - **Probed** (`TestProbeDecodedGaps`):
    - **Sent, 0 redactions:**
      - `headers: {"Authorization": "Basic FAKEVALUE"}` in a `proposals/*.md`;
      - `requests.get(u, headers={'Authorization': 'Basic FAKEVALUE'})` in a `proposals/*.md`;
      - the JSON form on a JSONL line that does not decode.
    - **Controls, caught:** `curl -H 'Authorization: Basic FAKEVALUE'`, and a JSONL record that decodes, because `fmt.Sprint` drops the quotes.
    - **Scope:** `Bearer` values are caught by the `bearer` needle; `Basic` and `token` are what is sent.
  - **Contradicts:** `local.md` ("`Authorization: <anything>`") and `delivery.md` ("bearer and authorization headers").
  - **Remediation:**
    - Change the pattern to `(?i)authorization["']?\s*:\s*\S`.
    - Add the two text forms to `TestReportRedactShapes`' `text` branch, and watch them fail before the change.
- **[low / confidence: high]** `internal/harness/report.go:145-148, 159-169` (`decoded`) — **the decoded check sees only what `json.Unmarshal` into `any` returns without error, and text files get none.** On such a line, a needle visible only after unescaping is sent: quoted JSON `{\"password\": \"…\"}`, `postgres:\/\/app:…@db`, or `\u0041KIA…`.
  - **Probed** (`TestProbeDecodedGaps`). Each case sent `FAKEVALUE` with 0 redactions:
    - two records on one line (3 needle forms);
    - a trailing byte;
    - `NaN`;
    - `1e400`;
    - a leading BOM;
    - nesting deeper than 10000;
    - a duplicate key whose later value shadows the escaped needle;
    - quoted JSON in a `proposals/*.md`.
    - **Controls:** each needle form in one well-formed record is caught.
  - **Regression in this range:** `1e400`, a number beyond float64. At 012b4ca the walk decoded with `UseNumber()` and caught it (`probe-base.log`). `json.Unmarshal` into `any` now fails on it, so `decoded` returns "".
  - **Inherited:** every other case here leaks at 012b4ca too.
  - **Contradicts:** nothing written. `local.md`'s "raw text OR JSON-decoded text" is accurate for decodable lines and silent about the rest.
  - **Why low:** agents append one well-formed record per line, and an escape hides a needle only in JSON-in-JSON or `\/` text.
  - **Remediation (stdlib):**
    - In `decoded`, read the line with a `json.Decoder` with `UseNumber()`, looping over every value. That closes concatenated records, bytes after a complete value, and the `1e400` regression.
    - State the rest as a limit in `local.md`: `NaN`, a BOM, a shadowed duplicate, and text files.
- **[low / confidence: high]** `internal/harness/report.go:171-183` (`redactDocument`), `:285-309` (`BuildRecord`) — **a needle in `state.json`'s `installedAt` or `vendor` replaces the `state.json` copy whole and records a redaction, while `report.json` sends both values verbatim.**
  - **Mechanism:** `BuildRecord` reads the original manifest through `LoadState`, and nothing runs `needleIn` over `report.json`.
  - **Probed** (`TestProbeStateFieldsInRecord`): `report.json` carried `"installed_at":"ghp_FAKEVALUE1"` and `"vendor":"token=FAKEVALUE2"`, beside `redactions: [{file: state.json, pattern: ghp_}]`.
  - **Why low:** the likelihood is negligible. `sh-agent install` writes both fields (a timestamp and `--vendor`), so this needs a hand-edited manifest.
  - **Inherited:** `report.json` was never redacted. The comment added here ("nothing report.json needs is lost") states the design without the caveat.
  - **Remediation:** run `needleIn` over the marshalled record. On a hit, null `vendor` and `installed_at` and add a note, or refuse to send.
- **[low / confidence: high]** `internal/harness/report_test.go` — **coverage of the decoded check.**
  - **At the pinned blob (06b295af):**
    - The case "needle behind a JSON unicode escape" held a literal `AKIAFAKEVALUE`, not `\u0041KIAFAKEVALUE` (checked with `od -c`).
    - The raw needle therefore caught it, and the case never reached the decoded path its name claims.
    - Mutant `mut1`, which has no decoded check in `redactLines`, failed only `TestReportRedactShapes`' JSON-in-JSON case.
  - **At 29c794b1:** the mid-review edit gives the fixture its `\u0041` escape, and `mut1` now fails both cases.
  - **Still open at 29c794b1:** mutant `mut2`, which removes `redactDocument`'s decoded check, passes all 12 report tests at both blobs.
  - **Probable cause (confidence: medium):** in this review, two of three tool writes turned escapes into literal characters: `\u0041` became `A`, and `\ufeff` became a raw BOM. This report was hit too, and was repaired and byte-checked with `grep` and `od -c`.
  - **Remediation:**
    - Add a `state.json` document whose needle shows only when decoded (`{"x":"\u0041KIAFAKE"}`) to `TestReportRedactDocument`, and watch it fail against `mut2`.
    - After any tool writes an escape-bearing fixture, check its bytes. Worth a bug record with that prevention rule.
- **[low / confidence: medium]** `.agents/docs/deploy/local.md` § Report back — **the round-2 finding 7 deferral is honest, but the signal it names is not in the dry run.** "`unterminated: true` in `report.json` shows it fired" is true only of `report.json`. The dry-run list prints `file:line (pattern)` without `lines` or `unterminated`, so a pre-send read cannot see that later records were dropped. This fails safe; it is not a leak.
  - **Remediation:** print `lines` and `unterminated` in `SendReport`'s list, or say in `local.md` that only `report.json` shows them.
- **[low / confidence: medium]** Round-2 finding 8's deferred shapes — **the deferral is honest.**
  - **Unlisted shapes:** `ghr_`, `npm_`, `--password VALUE` and `passphrase:` are all sent (probed). None is a listed shape, and `local.md` says anything unmatched is sent.
  - **Inside a listed shape:** `postgres://app:FAKE/VALUE@db` is also sent, while `scheme://user:pass@` is listed. A raw `/` in userinfo is not a valid RFC 3986 URL, so the listed shape arguably excludes it.
  - **Remediation:** none needed, beyond a clause in `local.md` if the operator wants it exact.

## Checked

- **Round-2 claimed fixes, verified by probe:**
  - **F1, push target:**
    - **Holds for:** an `origin` pushurl, a second `origin` url, and an env `origin` pushurl.
    - **Does not hold for:** URL rewrites, or a remote named by the URL (finding 1).
    - **Upstream after `-u <url>`:** `branch.main.remote` is set to the un-rewritten URL.
    - **Detached HEAD:** from a detached HEAD, `HEAD:main` pushed successfully to the same repository, so the destination does not change.
  - **F2, split key:** holds.
    - **Probed shapes:** array values, object fields, and a key spread over three records (one 3-line block marker) all go.
    - **Not covered:** the second trigger is finding 2(b).
  - **F3, secret-named key:** holds.
    - **JSONL shapes, each removed whole:** beside another hit, nested in an array, and an escaped key name (`pass\u0077ord`).
    - **`state.json`, each replaced whole:** a duplicate key, a nested password, an escaped AKIA, and a non-JSON document.
  - **F4, origin echo:** holds.
    - **Probed:** refusing an origin `https://x-access-token:ghs_FAKETOKENVALUE@…` names only the expected URL.
    - **Read:** no other path in `report.go` prints the checkout's origin.
  - **F5, test filter:** holds. `go test ./internal/... -list 'TestReport|TestSendReport|TestResolveFeedback'` lists all 12 report tests.
  - **F6, dry-run guidance:** holds. `README.md` and `local.md` say "before every report", and plan line 25 carries the timer condition.
  - **F7:** deferred in `local.md` and D13 (6); see finding 7.
  - **F8:** `redis://:FAKE@`, `Authorization: Basic` and `authorization: token` are caught; see finding 8.
- **Fuzz:** `FuzzProbeRedact`, 45 s, 2.60 M execs, over `redactLines` (JSONL and text) and `redactDocument`.
  - **Result:** no panic. No output line carries a needle in its raw text or, for JSONL, in its decoded text, and no needle survives in a document's output.
  - **Blind spot:** by construction it cannot see a secret that carries no needle (findings 2 and 4).
- **Declarations:**
  - **Removed:** `redactJSON`, `walkJSON`, `Redaction.Path`, and the imports `bytes`, `sort` and `strconv`. No reference remains in `internal/` or `cmd/`.
  - **New:** `marker` and `decoded`, both called.
  - **Unchanged:** `cmd/`, `go.mod` and `go.sum` are unchanged since 012b4ca.
- **Secrets:**
  - **Fixtures:** the added lines hold only fake values (`FAKEVALUE`, `FAKEBODY`, `ghs_FAKETOKEN`, `hunter2`, and the AWS documentation key).
  - **Constant and errors:** `FeedbackURL` carries no credential. `gitOut` errors quote git's argv and stderr, and the only URL in the argv is that constant.
- **Injection:**
  - **Git calls:** every call goes through `gitOut` to `exec.Command`, with no shell. The new push argv is `push --quiet -u <constant> HEAD:main`. Its only other inputs are `-C dir` and pathspecs under `reports/<host>/…`.
  - **Regexes:** RE2, so linear time.
  - **JSON:** `json.Unmarshal` stops at depth 10000.
- **Authz/authn:** the credential model is unchanged (the operator's SSH key). The push destination is finding 1.
- **Crypto:** nothing is implemented or composed here. The range has no signing, hashing or TLS.
- **Dependencies:**
  - **Imports:** none added; three stdlib imports removed.
  - **govulncheck:** `govulncheck ./...` (go1.27.1, `GOTOOLCHAIN=local`) exited 0 with "No vulnerabilities found."
  - **npm:** no npm ecosystem, so `npm audit` does not apply.
- **GCP surface:** nothing touched.
- **Residuals not raised:**
  - Proposal file names and the hostname are never needle-checked. They become paths, the commit message and `report.json` fields. Negligible, and inherited.
  - Unlisted shapes, for example a Kubernetes or Cloud Run `name: DB_PASSWORD` / `value: …` pair, are sent. That is what `local.md` says of any unmatched shape.
- **Probe hygiene:**
  - **Real repository at the end:** HEAD 012b4ca, no commits ahead of upstream, tags `v1.0.0` only.
  - **Status:** it lists only other agents' edits and untracked files, plus this report.
  - **The one shell `git init`:** its exit status was read, and `rev-parse --show-toplevel` confirmed that the probe repository was its own top.

## Not examined

- **Other agents' reviews:** `2026-09-17-review-012b4ca-worktree-code-docs.md` and `2026-09-17-review-015cf2d-012b4ca-code.md`, not read, to keep this review independent.
- **Uncommitted edits outside the requested diff:**
  - not read: `remote.md`, `bugs/INDEX.md`, BUG-053 to BUG-056, `context/AGENTS.md`, the prior review file's edits, and `testing/2026-09-17-report-round2-fixes.md` (it appeared mid-review);
  - read only in part:
    - the plan: lines 25 and 51-57, D9 and D13;
    - `guide.md`: only § 3b.
- **Test blob 29c794b1:** only diffed and run against the two mutants, not put through the full probe set. Any edit after the final hash was not seen.
- **Push probes:** these used `file://` URLs. The real `git@github.com:` constant was checked only by URL resolution (`ls-remote --get-url`), never by a push.
- **Not exercised:**
  - the operator's real global git config;
  - macOS (case-insensitive paths in the top-level check);
  - a real GitHub push;
  - `sh-agent report` in any form.
- **Redaction:** probed for shapes, not for every field the skills write.

## Outcome

Third round on both mechanisms, so the operator was asked what each must guarantee, in one numbered table with a recommendation per item, and took every recommendation as written (plan D14): redaction is a safety net and the guarantee is the operator's read before each report; the push trusts the operator's own git configuration. Must-fail for every code fix: one mutant per fix in `.agents/scratch/r3-mutants/`, each failing its own test; must-pass in the tree.

- Finding 1 (URL rewrites and a remote named after the URL redirect the push): rebutted as a defect, fixed as a claim — the operator's own configuration redirects every push on the machine and is trusted, not refused. The `FeedbackURL` and `commitAndPush` comments, `local.md` and `delivery.md` now say origin's pushurl and extra urls cannot redirect and URL rewrites can; D13 (3) is corrected by D14. No code.
- Finding 2 (secret not on the needle's line): (a) owned as a stated limit in `local.md` — multi-line values are sent; (b) fixed — on a header's own line a block closes only at a footer naming the header's label (`TestReportRedactStructuralGaps` "header line that mentions a footer", mutant `2b-label`). The round-2 Outcome's "fixed" for this trigger was wrong; bug record dispatched.
- Finding 3 (quoted `Authorization` name): fixed — `authorization["']?\s*:`; two text forms and their JSONL wrappings in `TestReportRedactShapes`, mutant `3-authz`.
- Finding 4 (lines that do not decode as one value; `1e400` regression): fixed for the regression and for several values on a line — `decoded` reads every value with `UseNumber()` up to the first non-JSON byte (`TestReportRedactStructuralGaps` "number beyond float64 …" and "two records on one line", mutant `4-decoded`); a BOM, `NaN`, a shadowed duplicate and escapes in text files are stated limits in `local.md`.
- Finding 5 (`state.json` values in `report.json`): struck — needs a hand-edited manifest; `install` writes both fields.
- Finding 6 (`redactDocument` decoded branch untested; escape mangling): fixed — an escaped needle in `state.json` in `TestReportRedactDocument`, mutant `6-document`; fixtures now written through a script with explicit escapes and byte-checked; bug record dispatched.
- Finding 7 (dry run does not show an unterminated block): fixed — the dry-run list prints `N lines` and `, unterminated: everything after it in this file was dropped`; `TestSendReportDryRunShowsDroppedRecords`, mutant `7-dryrun`.
- Finding 8 (unlisted shapes): struck — `local.md` already says anything unmatched is sent.
