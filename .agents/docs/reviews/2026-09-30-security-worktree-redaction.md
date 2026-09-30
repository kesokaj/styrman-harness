---
owner: security
date: 2026-09-30
---
# Security review — worktree-redaction   (2026-09-30)

Range reviewed: Scope: working tree — changed files: internal/harness/report.go internal/harness/report_test.go; read the diff with git diff HEAD -- internal/harness/report.go internal/harness/report_test.go.

## Status: COMPLETE

## Verdict: FINDINGS

Three findings, all low. Nothing leaks that leaked before.

- **The fix works.** A probe of the live patterns shows that each of the five shapes round 3 named is now redacted, in text mode and in JSONL mode. Mutations M1-M5 each removed or weakened one new rule. Each compiled, and each was caught by `TestReportRedactShapes`.
- **The over-redaction question: yes, the `url token` rule goes past what the written policy covers.** It does so in two ways, both in the safe direction:
  - it redacts usernames of 16 or more characters;
  - it redacts text that is not a userinfo at all, because its match can run out of a URL that has no path and into the text beside it.
- **The second kind is not new.** HEAD's unchanged `url credentials` rule already does the same thing.
- **The cost:** a whole line of a proposal, or a whole correction record in `feedback.jsonl`, is dropped.
- **Residuals.** Five shapes in the same two classes still pass: three Slack host variants, a GovSlack host, and a token with an empty password. This is the fifth pass over this deny-list, so it is an operator choice. Two small edits, one of them shorter than today's rule, close all five; a copy of the code with those edits was probed.
- **Stale comment.** One comment in `report.go` still describes a list entry this diff moved.
- **The working tree moved during this review, but not the two files in scope.** Someone else edited `.agents/docs/deploy/local.md:66` to describe the new rules, and a verifier log appeared. I made neither change. After every run, `report.go` and `report_test.go` were still `cmp`-equal to the copy I probed.

## Findings
- **[low / confidence: high]** `internal/harness/report.go:76-81` — `url token` (`://[^/\s:@]{16,}@`) over-redacts beyond the written policy.
  - **The written policy.** Three places state it: `report_test.go:122-123` (the `ASIAN` pin), `report.go:70-71` (a word that ends in `sk-`), and `deploy/local.md:66` ("tokens: 3", "ASIAN"). All three name only lexical collisions: text that looks like a needle.
  - **Probe.** Each of these benign lines was redacted, in both modes:
    - **(1) A lexical collision: a username of 16 or more characters in the userinfo.**
      - `clone https://simon-ostling-dev@bitbucket.org/ws/repo.git` (Bitbucket clone URLs carry the username);
      - `https://contoso-engineering@dev.azure.com/…` (Azure DevOps clone URLs carry the organization);
      - `ssh://firstname.lastname@host/repo.git`;
      - `postgres://reporting_readonly_user@db/x`;
      - `docker://ubuntu-with-extra-tools@sha256:abc`.

      This is the same kind of over-redaction as `ASIAN`. But the comment says "a username (`ssh://git@`) stays", which is true only below 16 characters. The old rule already dropped usernames of 20 or more characters. This change widens that to 16-19 characters, and to usernames with characters outside `[A-Za-z0-9_.%-]`.
    - **(2) Structural: the match leaves the URL.** The class admits `"`, `,`, `?`, `#`, `)`, `]` and any non-ASCII character. A userinfo holds none of them unescaped. RE2's `\s` is ASCII-only, so a non-breaking space pasted from a chat is not a space to it. A match therefore starts in a URL with no path and runs on to the next `@` in the text beside it:
      - `origins = ["https://app.example.com","ops@example.com"]`
      - `ALLOWED=https://app.example.com,ops@example.com`
      - `https://example.com?email=me@example.com`
      - `see [docs](https://docs.example.com)@simon`
      - `the url https://example.com#frag@anchor`
      - three words joined by U+00A0 non-breaking spaces

      The written policy names no case like these. It is not new in kind, though: HEAD's unchanged `url credentials` rule (`report.go:63`) crosses the same way. The probe redacted `{"url":"https://app.example.com","by":"ops@example.com"}` under `url credentials`.
  - **What the comment gets wrong.** "Anything a userinfo may hold" and "a username stays" both overstate: the class is wider than a userinfo alphabet, and only short usernames stay.
  - **Not over-redacted:** `ssh://git@…`, `ssh://ostling.simon@gmail.com@source.developers.google.com:2022/…`, `…/golang.org/x/sys@v0.20.0`, `git+https://github.com/org/repo@v1.0`, `https://unpkg.com/react@18`, and a URL followed by a space before an email.
  - **Remediation:** put it to the operator as one choice.
    - **(a) Recommended: accept it.**
      - Reword `report.go:76-78` to say what else goes: a username of sixteen characters or more, and text after a URL with no path when an `@` follows with no space between.
      - Pin one case of each among the benign lines, expected redacted, as `ASIAN` is pinned.
      - Add the username case to the over-redaction examples in `local.md:66`.
    - **(b) Narrow the class to a token alphabet, `[A-Za-z0-9._~%+=-]`.** Not recommended. `url credentials` crosses the same way and would need the same narrowing, so two rules get longer to protect usefulness, not secrecy. It also still drops long usernames.
- **[low / confidence: medium]** `internal/harness/report.go:81,83` — residual shapes in the same two classes. The probe of the live code left each of these **unredacted in both modes**:
  - `https://hooks.slack.com:443/services/T0/B0/…` (an explicit port);
  - `https://hooks.slack.com./services/…` (a trailing-dot fully qualified name);
  - `url=https%3A%2F%2Fhooks.slack.com%2Fservices%2F…` (a webhook embedded in another URL or in a form body, where the `/` after `com` is `%2F`);
  - `https://hooks.slack-gov.com/services/…` (the GovSlack webhook host; medium confidence in the vendor format);
  - `https://TOKEN:@host/r` (a userinfo token with an empty password). `url credentials` needs at least one password character, and `url token` stops at the colon. Some registries document this form, for example Gemfury's `https://TOKEN:@pypi.fury.io/…` (medium confidence).

  Not counted as residuals:
  - a 15-character userinfo, which falls under the threshold by design;
  - a scheme-relative `//TOKEN@host`, which is rare. JSONL mode catches it only by accident, because `decoded` prints `key:value` and so creates a `://`; text mode misses it;
  - a JSONL line that fails to decode and carries `https:\/\/` escapes. That limit is older than this diff, applies to both URL rules, and is `decoded`'s documented limit.

  **Remediation:** this is the fifth pass over this list, so the global rounds rule applies: put it to the operator as one choice.
  - **(a) Recommended: two edits.** Each gets a `TestReportRedactShapes` row.
    - `slack webhook` becomes `(?i)hooks\.slack`, dropping `\.com/` (five characters shorter).
    - `url token` becomes `://[^/\s:@]{16,}:?@` (two characters longer).

    I probed both edits in the scratch copy only:
    - all five residuals were then redacted in both modes;
    - `TestReportRedactShapes` still passed, every positive and benign row;
    - the benign URL lines stayed unredacted.

    After these edits, record any further residue as `left: operator — deny-list residue accepted`.
  - **(b) Record all of it as accepted residue now.** `local.md:66` already calls the list "a safety net, never a guarantee".
- **[low / confidence: high]** `internal/harness/report.go:43-44` — the prefix-list comment says its entries include "a Slack webhook URL". This diff removed both Slack prefixes from that list; the webhook is now the `slack webhook` rule at `:83`. **Remediation:** delete "and a Slack webhook URL" from the comment in the same change (the global rule on text that describes a behaviour).

## Checked
- **Secrets:** the diff adds only `FAKEVALUE` placeholders. `git status` shows no secret-shaped path.
- **Injection:** the only new logic is RE2 regexes, which run in linear time, with no exec, SQL or template.
  - **What `report.json` records:** `needleIn` returns the rule name (`slack webhook`) or the needle literal (`sk-None-`), never the matched text (BUG-054 still holds).
  - **Callers:** only `Collect` (`report.go:286-314`) calls the redactors, before `--dry-run` output and before the push. It sends:
    - `feedback.jsonl` and `retros.jsonl` in JSONL mode;
    - `state.json` whole, written by `MarshalIndent`, so whitespace separates its values;
    - `proposals/*.md` in text mode.
- **Superset check:** every line that contains one of the two removed prefixes also contains `hooks.slack.com/`. Only the recorded pattern name changes.
- **Authz/authn:** nothing touched.
- **Crypto:** nothing composed.
- **GCP surface:** none.
- **Dependencies:** `go.mod` and `go.sum` are unchanged.
  - `govulncheck ./...` exited 0: "Your code is affected by 0 vulnerabilities", plus 1 vulnerability in a required module that the code does not call (GO-2026-5024 per `AGENTS.md`).
  - `npm audit` does not apply: there is no `package.json`.
- **Tools on the live tree (go1.27.1):** `gofmt -l cmd internal` printed nothing; `go vet ./...` exited 0; `go test -race -count=1 ./...` exited 0, with all four packages ok.
- **Probe:** `zz_probe_test.go` calls `redactLines` in text mode and in JSONL mode, in a `tar` copy of `git ls-files -co --exclude-standard`. `GIT_CEILING_DIRECTORIES` was set, and `git rev-parse` inside the copy failed as it should. `cmp` showed `report.go`, `report_test.go`, `go.mod` and `go.sum` equal to the working tree.
  - The positive control (`AKIA…`) was redacted, and the negative control (`ssh://git@…`) was not.
  - Every round-3 shape was redacted by the intended rule: `slack webhook` ×2, `url token` ×2, `sk-None-`.
- **Mutations:** each was applied to a copy, its diff checked non-empty, compiled with `go test -run '^$'`, then run against `TestReportRedactShapes`, and the file restored and checked equal with `cmp`.
  - **Killed, each by the intended row:**
    - M1 dropped `sk-None-`;
    - M2 dropped `(?i)`;
    - M3 raised the threshold to `{17,}`;
    - M4 removed `=` from the class;
    - M5 deleted the Slack rule;
    - M7 lowered the threshold to `{10,}` (the benign `deploy-bot` row caught it);
    - M8 changed the Slack rule to `(?i)slack\.com/` (the benign `slack.com/apps` row caught it).
  - **Survived:** M6 (`{11,}`). The benign pin is 10 characters long, so the threshold can fall to 11 unnoticed. That is the over-redaction direction, so it is not a finding; a benign username of exactly 15 characters would pin the threshold exactly.
- **Past bugs** (`sh-agent search --refs`):
  - BUG-053 and BUG-110 cover this deny-list mechanism.
  - BUG-056, BUG-058 and BUG-059 cover structure and the decoder, which this diff leaves unchanged.
  - BUG-057 covers escapes in fixtures; the new rows hold none.
  - Round 3's side note, pinning one case like `task-admin-panel`, was not taken up. It is outside this diff and not a finding.

## Not examined
- **Dry-run tier:** the end-to-end `sh-agent report --dry-run` tier (`testing/guide.md` § 3b), which the verifier's log `testing/2026-09-30-verify-redaction-rules.md` reports running. I did not re-run it, and I read only that log's header.
- **Concurrent `local.md` edit:** read as a diff only, not reviewed as a document.
- **Other living docs:** a grep found no other description of the changed thresholds. `delivery.md` § The loop back was not re-read.
- **Coverage beyond round 3:** vendor webhook and token formats other than Slack and OpenAI, and any deny-list kind round 3 did not name.
- **Slack host behaviour:** whether `hooks.slack.com` accepts an explicit port or a trailing dot. I reasoned from URL syntax and did not probe the vendor.
- **Scratch left on disk:** `.agents/scratch/sec-worktree-redaction/` (the copy, the probe, the pristine `report.go`, `live-test.log`, `govulncheck.log`). It is gitignored and can be deleted.

## Outcome
- 1 over-redaction of `url token` — fixed: accepted in the safe direction, the comment reworded at `report.go` (`url token`), one case of each pinned in `TestReportRedactShapes` (a long username, a pathless URL with an unrelated `@`).
- 2 five residual shapes — fixed at `report.go`: `(?i)hooks\.slack` and `://[^/\s:@]{16,}:?@`, the rules this report proposed; five rows in `TestReportRedactShapes`, mutation fails ten lines.
- 3 stale comment — fixed at `report.go` (prefix list comment).
