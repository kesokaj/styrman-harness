---
owner: security
date: 2026-09-30
---
# Security review — worktree-completion   (2026-09-30)

Range reviewed: working tree — changed files: .agents/docs/deploy/local.md agents/verifier.md internal/harness/harness_test.go internal/harness/report.go internal/harness/report_test.go internal/harness/state.go internal/search/global_test.go internal/search/reconcile.go; untracked, not ignored: none; read the diff with git diff HEAD -- <files>.

Base: HEAD e83c2d0, plus uncommitted edits to 8 files. The tree moved during the review: `report_test.go` gained one line of benign rows. `report.go`, the code under test, was byte-identical to the probed copy at the end (`cmp`), so the probe results stand.

## Status: COMPLETE

## Verdict: FINDINGS

Two findings, both low. Nothing blocks on security grounds.

**The diff does what it says.** All 19 shapes the v1.13.0 review listed are redacted, in four modes: text, JSONL, JSONL with `\/`-escaped slashes, and a `state.json` document.

**The worst problem is what still passes: near shapes of the same families.** These are sent verbatim:
- Discord's versioned API path (`discord.com/api/v10/webhooks/…`), which is what client libraries build.
- The Teams Workflows URL (`…/workflows/…/invoke?…&sig=…`), which replaces the connector URL.
- Azure keys outside the storage connection string's `AccountKey=` field.
- Cookies sent any way other than a `Cookie:` or `Set-Cookie:` header.
- A compact JWE with an empty key segment.
- An age post-quantum identity.

**`deploy/local.md:66` now overstates coverage.** It says each named family "has one rule", which reads as coverage of Discord and Teams webhooks that the rules do not give. That is the same unsafe-direction wording the v1.13.0 finding was about.

**No new rule matches this project's ordinary prose.** Across 89,887 lines of the repository's own docs, the only lines redacted by a new rule alone are six review lines that quote the shapes themselves.

**The two ReadDir changes read nothing outside the intended directory.** On a plain path they read exactly what the glob read. With a glob metacharacter in the path, they read the correct directory, which the glob missed.

## Findings

- **[low / confidence: high on the regex behaviour; medium on how common each form is]** `internal/harness/report.go:51,96-103` (the needles this diff adds): near shapes of the named families still pass.
  - **Probe.** Setup: an rsync copy of the working tree in `.agents/scratch/security-worktree-completion-probe/src/` (`report.go` identical by `cmp`), plus the scratch test `internal/harness/zz_probe_test.go`, which calls the real `redactLines` and `redactDocument`. Log: `.agents/scratch/security-worktree-completion-probe/probe.log`.
  - **Every entry below was SENT in all four modes:**
    - **Discord:**
      - `https://discord.com/api/v10/webhooks/<id>/<token>`: the versioned path Discord's API base and client libraries use.
      - `discord.com:443/api/webhooks/…`, a port.
      - `discord.com./api/webhooks/…`, a trailing dot.
      - `discord.com/api%2Fwebhooks%2F…`, an encoded path.
      - These are the same four shapes the Slack rule had to absorb in earlier rounds.
      - The `canary.` subdomain is caught.
    - **Teams:**
      - The Workflows (Power Automate) webhook URLs, `https://prod-NN.<region>.logic.azure.com:443/workflows/…/triggers/manual/paths/invoke?…&sig=…` and `…environment.api.powerplatform.com:443/powerautomate/automations/direct/workflows/…&sig=…`. These are the form that replaces Office 365 connectors (medium confidence on how widely they have replaced them by now).
      - The legacy `https://outlook.office.com/webhook/…` (likely dead URLs, low relevance).
    - **Azure:**
      - `SharedAccessKey=` in a Service Bus or Event Hubs connection string, which uses the same syntax as a storage one.
      - A SAS URL's `sig=`.
      - `AZURE_STORAGE_KEY=…`.
      - A JSON-keyed `"AccountKey": "…"`.
      - `AccountKey = …` with spaces.
    - **Cookies:**
      - `curl -b 'session=…'`.
      - `curl --cookie "sid=…"`.
      - A HAR pair `{"name":"Cookie","value":"session=…"}`.
      - `document.cookie = "session=…"`.
    - **JWT:** a compact JWE with direct key agreement, whose second segment is empty (`eyJhbGciOiJkaXIi….. <iv>.<ct>.<tag>`), as Auth.js session tokens are. `jwt` needs 8 or more characters in that segment and `jwt short` needs 2 or more, so neither rule matches it.
    - **age:** `AGE-SECRET-KEY-PQ-1…`, the post-quantum hybrid identity (medium confidence in this prefix). The needle pins `AGE-SECRET-KEY-1`, the X25519 form only.
    - **Legacy formats, too generic to match (accept as residue):**
      - Vault `s.<24>`.
      - Google `1/<token>`.
      - An unsigned `eyJ….e30.` JWT, which is not a credential.
  - **Why it stays low.**
    - The destination is the operator's private feedback repository.
    - `delivery.md` makes the operator's read of `~/.config/sh-agent/` before each report the guarantee.
    - `local.md:67` says anything unmatched is sent.
  - **Remediation: put it to the operator as one choice.**
    - This is at least the seventh pass over this mechanism. The global rule stops and asks after three rounds without convergence. The code review of this worktree also records that the operator had accepted these families as residue.
    - **(a) Recommended.** Five edits that keep the rule count, plus one new rule, each with one `TestReportRedactShapes` row:
      - discord → `(?i)discord(app)?\.com\S*webhooks`. This covers versioned paths, ports, trailing dots and encoded paths. The docs URL `discord.com/developers/docs/resources/webhook` still passes.
      - `jwt short`'s middle segment `{2,}` → `*`, which catches the empty-segment JWE.
      - `AGE-SECRET-KEY-1` → `AGE-SECRET-KEY-`.
      - azure → `(?i)(account|sharedaccess|storage)[_-]?key["']?\s*[:=]`.
      - cookie `:` → `[:=]`, which catches `document.cookie =`.
      - One new rule for the signed-URL kind: `(?i)[?&]sig=[^&\s]{16,}`. It covers Teams Workflows URLs, Azure SAS URLs and SAS tokens.
      - Name what remains (curl cookies, HAR pairs, the legacy formats) as residue.
      - These regexes are my proposal and are UNPROBED. Each needs its row to prove it before it counts.
    - **(b)** Keep the rules as they are, narrow the guide (finding 2), and record the list above as residue.
    - Either way this needs only stdlib `regexp`, no dependency.

- **[low / confidence: high]** `.agents/docs/deploy/local.md:66`: "the families a review named on 2026-09-30 (… Discord and Teams webhooks, … an Azure storage connection string, a cookie header, … a JWT with a short payload) each have one rule" reads as coverage of those families.
  - **Why this is a problem.** Per finding 1, the current Teams webhook form and Discord's versioned form are not matched. That is the direction the v1.13.0 finding 1 called unsafe, in the document the operator relies on when deciding what to read before a report.
  - **Also stale.** The paragraph's list of shapes matched without a leading word boundary names a JWT, an OpenAI key, an npm token and a Hugging Face token. It does not name the four new ones that are also boundary-less: SendGrid, Vault, the Google refresh token and the short JWT. The paragraph does say to read the list rather than itself, so this part is minor.
  - **Remediation.**
    - With (a), the sentence can stay.
    - With (b), name the forms actually matched:
      - Discord webhook URLs on the unversioned `/api/webhooks/` path.
      - Teams connector URLs (`webhook.office.com`, `/webhookb2/`), not Workflows URLs.
      - The `AccountKey=` field of an Azure storage connection string.
      - A `Cookie:` or `Set-Cookie:` header.
      - A JWT whose second segment is 2 characters or more.

## Checked

- **Secrets.** The 8 files contain only placeholder values (`FAKEVALUE…`, UUIDs of zeros). The probe used only `FAKE…` values. Nothing echoes a secret to logs or errors.

- **Injection and parsing.**
  - The new needles are Go RE2 regexes: linear time, no ReDoS.
  - The ReadDir changes remove the last two path-built glob patterns in the scratch-home paths (BUG-109's class). The remaining `fs.Glob` calls (`inventory.go:220,260`, `render.go:38`) run on an `fs.FS` with fixed relative patterns.

- **The two ReadDir changes, measured.** Setup: the scratch tests `internal/search/zz_readdir_test.go` and `internal/harness/zz_readdir_test.go` compare the real new code against HEAD's glob code, copied verbatim. Log: `.agents/scratch/security-worktree-completion-probe/readdir.log`.
  - **Plain path.** `configView` and `LastRetro` gave the same result as the glob on these name shapes:
    - `a.md`, `.hidden.md`, a bare `.md`
    - `x.MD`, `n.txt`, `b.md.bak`, `c.md~`
    - a directory `d.md`
    - a symlink `l.md` to a file outside `proposals/`
    - a dangling `z.md`
  - **With a metacharacter in the path.**
    - `c[1]`, an unclosed `u[x`: the old `configView` listed no proposals and the new one lists the right file.
    - `h*`: the same result from both.
    - The old `configView` never read a decoy sibling's content, because its `TrimPrefix` failed and the joined path did not exist. That was safe by accident.
    - The old `LastRetro` under `c[1]` stat'ed the decoy in `c1/…/proposals/` and returned its mtime (2026-08-01, not 2026-01-01). That is metadata only. The new code does not.
  - **So the new code reads files the glob did not only in the intended directory, when that directory's path holds a metacharacter.** It reads nothing outside that directory.
  - **Unchanged and pre-existing.** Both versions follow a symlinked `proposals/*.md` to a target outside the directory. `configView` indexes that target locally, by design per its comment. The report's `Collect` checks with `Lstat` and skips links (BUG-036).
  - **The new tests fail on HEAD's code.** For `configView` this is shown by case B: no proposals were listed under `c[1]`. For `LastRetro` it is by reasoning from case B: the glob resolves `u[1]` to a `u1` that does not exist, so `ok` is false.

- **Over-redaction by the new rules on synthetic prose.** Each of these is within the accepted over-redaction policy by analogy with the pinned cases:

  | Rule | Prose that fired it | Pinned case it matches |
  |---|---|---|
  | `cookie header` | any `cookie:` (`the cookie: banner`, `session: { cookie: { secure: true } }`, `fortune cookie: tasty`, `Cookie: none`) | "tokens: 3" |
  | `sendgrid key` | an identifier chain after an upper-case `SG` (`MSG.createInstanceOfClass.withSomethingLong`, `SSG.renderStaticPages…`) | a word ending in `sk-` |
  | `azure account key` | `AccountKey=<from Key Vault>` | a placeholder |
  | `teams webhook` | a bare `/webhookb2/` path | none; contrived |
  | `google refresh token` | a `gs://bucket1//0…` double-slash path | none; contrived |

  - One inconsistency, in the safe direction: `the cookie:` at the end of a line is sent as text but redacted in JSONL, where the closing quote counts as the value.
  - **Corpus result.** 1,074 files and 89,887 lines of this repository's own prose: the only lines redacted by a new rule and by no older rule are six review lines that quote the shapes. No ordinary prose fired.

- **Authz/authn.** Not touched.

- **Crypto.** Nothing added. `configView`'s `sha256` is a content key, stdlib used as designed. No hand-rolling.

- **Dependencies.**
  - None added; `go.mod` and `go.sum` are not in the diff.
  - `govulncheck ./...` (bare, go1.27.1) exited 0: "Your code is affected by 0 vulnerabilities", plus one vulnerability in a required module that is not called (the known GO-2026-5024, per AGENTS.md). Log: `.agents/scratch/security-worktree-completion-govulncheck.log`.
  - `npm audit` not run: no npm manifest is in the changed files.

- **GCP surface.** Not touched.

- **`agents/verifier.md`: the changed paragraph (long suites), read for instruction-shaped text.** None found.
  - Nothing directs network use, credential access, exfiltration, a change to permissions or config, or overriding another rule.
  - The new `kill <pid>`, then `kill -0 <pid>` fails, is scoped to the verifier's own background job. A wrong kill would need the PID reused, which means a full PID wrap between the suite's exit and the kill: negligible.
  - `> file 2>&1, never a pipe` is consistent with BUG-086.

- **Other agents' records.** `2026-09-30-review-worktree-completion-code.md` raises a high process finding: the operator's recorded acceptance is reversed without approval. That is not a security question. From the security side the added rules only add redaction, and they fired on no ordinary prose in the corpus.

## Not examined

- `agents/verifier.md` outside the changed paragraph (the brief scoped the read to the change).
- The remediation regexes proposed in finding 1. They were not probed.
- `report_test.go`'s final form, including the line added mid-review (`"see webhook.office.com for the docs", "the AccountKey field is required"`). It was not run. By the probed patterns both rows pass: no `/` after the host, no `=` after the key. The project's test suite was not run either; that is the verifier's job.
- `local.md`. I read it from the diff taken at the start and did not re-diff it after the tree moved; its stat line was unchanged.
- The operator's real `~/.config/sh-agent/`, deliberately not read because it is private data. The repository's own prose served as the corpus instead.
- The parts of the redaction pipeline the diff does not change (`blockEnd`, `decoded`), beyond reading them. Nothing of `report.go` past line 350: the push and the clone path.
- Scratch left for reproduction, disposable:
  - `.agents/scratch/security-worktree-completion-probe/`: a source copy, three `zz_*_test.go` probes, `probe.log`, `readdir.log`.
  - `.agents/scratch/security-worktree-completion-govulncheck.log`.

## Outcome
- 1 near shapes of the same families — fixed at `internal/harness/report.go`: the discord, teams (a `sig=` alternative), azure (three more key names), cookie (`-b`/`--cookie`, `document.cookie`), jwt-short (an empty middle segment) and age (`AGE-SECRET-KEY-` prefix) rules widened, the rule count unchanged; nine rows and three benign near-misses in `TestReportRedactShapes`, the mutation fails 18 lines. The list is closed for instance work: the comment at the end of `keywordNeedles` and `deploy/local.md:66` say a shape past these rules is a design question for the operator.
- 2 the guide implies full coverage — fixed at `.agents/docs/deploy/local.md:66`: "for the shapes the review named, widened once … a shape past those rules is sent".
