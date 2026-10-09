---
owner: security
date: 2026-10-09
---
# Security review — dbd776d-3bffdcf   (2026-10-09)

Range reviewed: commits dbd776d..3bffdcf (git log dbd776d..3bffdcf); the diff is git diff dbd776d...3bffdcf.

## Status: COMPLETE

## Verdict: FINDINGS

The worst problem is in the new command-line credential needle (`internal/harness/report.go:151`). It puts `(^|\s)` in front of `sshpass` and the MySQL client names. That anchor fails the same way as the `\b` the project rule bans: the token is missed when a quote, a backtick, a `/` or a `key:` comes right before the command name. I ran a probe in a scratch copy at 3bffdcf (`.agents/scratch/sec-3bffdcf`, a test file inside package `harness`, since removed). With `redactLines(..., jsonl=true)`, `{"text":"sshpass -p FAKEVALUE ssh h"}` and `{"text":"mysql -u root -pFAKEVALUE db"}` both came through **unredacted**. The raw line has a `"` before the command name. The decoded line, printed by `fmt.Fprintln` of a Go map (`report.go:324`), has `text:` before it. `/usr/bin/mysql -uroot -pFAKEVALUE` was also not matched. This text reaches `sh-agent report` (feedback.jsonl, retros.jsonl, proposals) and the dreamer's prompt (`internal/dream/collect.go:157`, `:254`). The mcpServers refusal is correct and is reached on every install path. Its two test rows fail when the refusal is removed: I removed `"mcpServers"` in the scratch copy and `TestLintPermissionKeysByKind` failed on both rows. The usage-error change, the `Lint` removal and the process-group kill open no new hole. The process-group change does leave one gap: a child process can outlive the run (finding 6).

## Findings

- **[medium / confidence: high]** `internal/harness/report.go:151` — In the `command-line credential` needle, `(^|\s)` before `sshpass` and before `(?:mysql|mariadb|mysqldump|mysqladmin)` hides the token when a non-space character comes right before the command name. Four cases show it. (1) A JSON string value that starts with the command: the raw line has `"sshpass`, and `decoded()` (`report.go:315-325`) prints a map as `map[text:sshpass ...]`. Both were missed in the probe. (2) Markdown inline code, `` `mysql -pX` ``, which is common in prose, trail files and transcripts that the dreamer reads through `harness.Redact`. (3) A full path, `/usr/bin/mysql -pX` (missed in the probe). (4) A literal `\n` before the name in raw text, the case the AGENTS.md rule names. A password typed into such a command reaches the feedback repository or the model provider. This is the BUG-110/BUG-127 failure class. The `-u` and `--user` part is not affected, because the flag always comes after a space. **Remediation:** for `sshpass`, drop the anchor; the word is distinctive. For the MySQL family, use `(^|[^A-Za-z0-9_])`, which still spares `mkdir -p`. Add `{"text":"sshpass -p FAKEVALUE ssh h"}` (through the jsonl path), `` `mysql -pFAKEVALUE` `` and `/usr/bin/mysql -pFAKEVALUE` to `shapeSecrets`. Also check the older cookie rule's `(^|\s)(-b|--cookie)` (`report.go:138`) the same way. The same anchor is there, though `-b` normally comes after a space.
- **[medium / confidence: high]** `internal/harness/report.go:151-152` — The brief named Docker credentials. The two new needles cover only `curl -u` and `config.json`'s double-quoted `"auth": "<base64>"`. These forms were missed in the probe:
  - `docker login -u user -p FAKEVALUE reg`. `-p` is matched only after `sshpass` or the MySQL clients; `--password` is caught by the existing flag rule.
  - `'auth': '…'` (single-quoted) and YAML `auth: <base64>`.
  - A Kubernetes `.dockerconfigjson: <base64 of the whole config>`. Its inner `"auth"` is never visible as text.

  (`"identitytoken"` and `--docker-password=` are already caught by `secret assignment`.) **Remediation:** this is a new kind, so it is the operator's design question under the BUG-110 rule. Proposed: a `docker login` line carrying `-p`/`--password` as one rule, and `["']?auth["']?\s*[:=]\s*["']?[A-Za-z0-9+/=]{8,}` together with `\.?dockerconfigjson` as the Docker-auth kind. Each needs a near-miss.
- **[low / confidence: high]** `internal/harness/report.go:151` — The `command-line credential` needle has documented gaps beyond finding 1. The probe confirmed each of these as unmatched:
  - `-u` and `--user` forms: the attached `curl -uadmin:FAKEVALUE`, the empty-user `curl -u :TOKEN`, and a digit-first user `curl -u 12345:FAKEVALUE`. The last is excluded on purpose to spare `docker run -u 1000:1000`, but numeric account IDs exist.
  - Other command-line credentials: `redis-cli -a PASS`, `smbclient -U user%PASS`, `htpasswd -b file user PASS`, and the `.npmrc` form `_auth = <base64>`.

  **Remediation:** this is a decision for the operator. Recommended: accept `-u` with the value attached; allow an empty user (`-u :x`); add the `redis-cli -a` and `.npmrc _auth` kinds. Add a near-miss for each.
- **[low / confidence: high]** `internal/harness/report.go:152` — The `credentials key` value alphabet `[A-Za-z0-9+=_.-]{16,}` leaves out `/`. A base64 credential with a `/` in its first 16 characters is not matched. That is about 22% of random base64 values; `credentials: abc/defFAKEVALUE0123456789` was missed in the probe. **Remediation:** allow `/` but refuse a value that starts with `/`, `~` or `./`, or ends in `.json`, so `GOOGLE_APPLICATION_CREDENTIALS=/path` stays a near-miss. Add a slash-bearing shape to `shapeSecrets`.
- **[low / confidence: medium]** `internal/harness/inventory.go:261` and `internal/harness/plan.go:156` — The Claude agent file is deployed verbatim (`a.Raw`), and the lint uses a deny-list. Two keys from Claude Code's documented subagent table (`.agents/docs/plans/2026-09-30-research-audit-fix-formats.md:12`) still pass:
  - `memory`: per Claude Code's sub-agents documentation, it enables Read, Write and Edit automatically, which widens an agent's `tools:` allowlist without showing it.
  - `initialPrompt`: it is auto-submitted as the first user turn under `claude --agent`. Text the operator never typed would then arrive as "the operator's own turn", which the global context treats as approval.

  I did not re-verify either effect against a live CLI. Antigravity is not exposed, because `DeriveAgy` (`derive.go:41-62`) writes only an allowlist of keys. **Remediation:** replace the deny-list with an allowlist of agent keys (`name`, `description`, `tools`, `disallowedTools`, `model`, `maxTurns`, `skills`, `effort`, `color`), so a future vendor key is refused by default. Failing that, add `memory` and `initialPrompt` to the refused set.
- **[low / confidence: medium]** `internal/dream/run.go:211-217` — Two gaps came with the process-group change:
  - **Ctrl-C:** `Setpgid` takes the CLI out of the terminal's foreground group, and neither `cmd/sh-agent` nor `internal/dream` handles signals (grep for `signal.` found nothing). Ctrl-C on an interactive `sh-agent dream` kills sh-agent but not the CLI. The CLI then runs on as an orphan with the operator's credentials, sending the prompt to the model, with no deadline.
  - **Normal exit:** when the CLI exits before the deadline while a helper (an MCP or language server) keeps the pipes open, `Cancel` is never called. The helper is never killed, and `Run` returns `exec.ErrWaitDelay` after 30 s, so a successful run is reported as a failure. `TestCallReturnsAtTheDeadlineDespiteAGrandchild` covers only the deadline path.

  **Remediation:** after `Run` returns, on every path, send `syscall.Kill(-pid, SIGKILL)` to the group. In `runDream`, use `signal.NotifyContext(SIGINT, SIGTERM)` to drive the same context. Add a test for the early-exit-with-helper case.
- **[low / confidence: medium]** `internal/harness/report_test.go:184-197` — Two cases are over-redacted, which is safe, but neither is pinned in the over-redaction list as AGENTS.md requires: `mysql -P3306 -h db` (`(?i)` turns the port flag `-P` into a match) and `docker run -u node:node img`. In this repository's own trail, the three new needles match four lines (`.agents/docs/reviews/2026-09-30-security-cb799ff-f1ca525.md:38`, `2026-09-30-security-f1ca525-e2cffc1.md:27`, `2026-10-09-review-project-merged.md:39`, `2026-10-09-security-project.md:25`). These are review records quoting `curl -u user:pass`; they are redacted out of dream prompts, which does no harm. **Remediation:** pin both over-redactions in the test.

## Checked
- **Secrets:** none found. I grepped the added lines of the diff for AWS, GitHub, Google API, Slack, OpenAI-style, PEM and password-assignment shapes; the test fixtures use `FAKEVALUE`. New error output: the strings printed by `IsUsage` are constant (`search.go:98-128`, six `usageError` sites, none echoing the user's value). `(*ctx).vendor` now prints the `NormalizeVendor` error on `c.out`. That error can carry a manifest's `vendor` value, which is not secret and was already printed before, on stderr.
- **Injection:** no new command, SQL or path construction. Dream skills are now read with `harness.Split` and `harness.Parse` (`collect.go:284-290`). Block scalars are folded with spaces, so a hostile `SKILL.md` cannot inject a newline or a fake heading into the dreamer's prompt, which the old key-line reader could not either. The skill file is still read with `Lstat` and regular files only (`collect.go:212-221`). There is no length cap (`Parse` rather than `ParseStrict`), and skills are dropped from the budget last (`collect.go:169`). That was already the case with the old single-line reader, so it is not a new finding.
- **Authz/authn and the lint:** `mcpServers` is refused in `LoadAgents` (`inventory.go:261`). Every install path reaches it: `harness.Plan` → `LoadInventory` → `LoadAgents` (`plan.go:76`). The `--from-zip` path goes through `installFiles` → `harness.Plan` (`main.go:797`), and `status` and `report` also call `LoadInventory`. Mutation check: with the key removed in the scratch copy, the test failed on both new rows; restored, it passed. `Lint` was removed, but every loader already called `ParseStrict` (`inventory.go:190,241,299`), and no other `Lint` caller exists, so no load path lost the lint.
- **Process-group kill:** `Cancel` runs only after a successful `Start`, so `cmd.Process` is never nil there. Linux does not hand out a PID that is still in use as a process-group ID, so `Kill(-pid)` reaches only the CLI's own group: the CLI and the descendants that stayed in it. sh-agent sits in its own group and is never hit. There is a microsecond window between the leader being reaped by `Process.Wait` and `Wait` reading the context result. In that window the kill still targets the original group, or returns ESRCH.
- **`search.IsUsage`:** it matches only the unexported `usageError` (`errors.As`), which only `validate` builds. A database or filesystem error cannot be classified as a usage error.
- **Crypto:** nothing in the diff hand-rolls a primitive or protocol. No finding.
- **Dependencies:** no change to `go.mod`. `/home/ostling-simon/.go/bin/govulncheck ./...`: 0 vulnerabilities the code calls, 0 in imported packages, and 1 in a required module that the code does not call (the known GO-2026-5024 in AGENTS.md; I did not re-list it in verbose mode). There is no npm ecosystem in scope.
- **GCP surface:** the diff touches no IAM, service account, deploy flag or Secret Manager.
- **Bug history:** `sh-agent search --refs <file> --kind bug` for report.go (BUG-110, BUG-127), inventory.go (BUG-143, BUG-036), run.go (BUG-166, BUG-169), collect.go (BUG-164, BUG-131), main.go (BUG-168, BUG-130) and search.go (BUG-168, BUG-125). Findings 1–4 are in the BUG-110/127 class.

## Not examined
- The other 34 files of the commit, beyond their stat line: delivery_test, usage_test, zip_test, the `internal/trail` and `internal/memory` changes, `internal/metrics` (claude.go, events.go), release.go, the `internal/search` export and reconcile changes, and the docs and skill line edits.
- Any live Claude Code run to confirm the `memory` and `initialPrompt` effects. Finding 5 relies on the vendor documentation.
- The pre-existing skill-body dynamic shell injection (`` !`cmd` `` in a `SKILL.md`), and whether the skill lint should refuse it. It is out of the diff.
- Whether `report` sends anything else that is JSON-shaped beyond feedback.jsonl, retros.jsonl, state.json and proposals. I read only `report.go:380-435`.
- `govulncheck -show verbose`, to name the one required-module advisory.

## Round 2

Range reviewed: the round-1 fixes, uncommitted in the working tree, `git diff 3bffdcf` (report.go, report_test.go, dream/run.go, run_test.go, cmd/sh-agent/main.go, usage_test.go, delivery.md, search.md, skills/harness-audit/SKILL.md).

### Status: COMPLETE

### Verdict: FINDINGS (low only)

Round-1 findings 1–4, 6 and 7 are fixed. Finding 5 (agent keys `memory` and `initialPrompt`; allowlist) is still open: the review record says the operator decides, and `inventory.go` did not change. I ran a probe against a scratch overlay: `.agents/scratch/sec-r2` holds `git archive 3bffdcf` with the nine changed files copied over, and the probe file was removed afterwards.
- Every round-1 shape that slipped through is now redacted: the `"text":` JSON line on the jsonl path, `{"cmd":["sshpass","-p",…]}`, the backtick, the full path, `-uadmin:x`, `-u :x`, `-u 12345:x`, `docker login -p`, single-quoted and YAML `auth:`, `.dockerconfigjson`, `_auth =`, `redis-cli -a`, `smbclient -U user%pw`, `htpasswd -b`, and `credentials: abc/def…`.
- The fixes introduced one regression: a short Docker `auth` value is now missed (R2-2).
- The process-group fix closes the Ctrl-C and clean-exit gaps. It leaves SIGHUP and one small signal-drop window (R2-3, R2-4).

### Findings
- **[low / confidence: high]** `internal/harness/report.go:160` — **Regression.** The `auth key` rule needs 16 base64 characters, not counting the `=` padding, so a short Docker `auth` value that 3bffdcf's `"auth"…{8,}` caught now passes. Probe results:
  - `"auth": "YWRtaW46YWRtaW4="` (`admin:admin`): not matched.
  - `"auth": "YWRtaW46cHcx"` (`admin:pw1`): not matched.
  - `{"t":"{\"auth\":\"YWRtaW46YWRtaW4=\"}"}` on the jsonl path: not redacted.

  A weak default registry password is the likeliest one to be short. **Remediation:** keep the round-1 double-quoted JSON alternation `"auth"\s*:\s*"[A-Za-z0-9+/=]{8,}"` next to the 16-character rule for the other quotings (`"auth": "basic"`, at 5 characters, stays benign). Add the `admin:admin` shape to `shapeSecrets`.
- **[low / confidence: high]** `internal/harness/report.go:159` — The `-u` alternative now has nothing in front of it, so it matches inside any hyphenated word followed by `:` or `%`. It redacts 66 of 120,390 lines in this repository's trail and skills, against 4 at 3bffdcf. Examples:
  - `date -u +%FT%TZ`, a command the trail itself prescribes (`.agents/docs/architecture/DECISIONS.md:123`).
  - `anti-ui-slop/SKILL.md:3` (`.agents/docs/plans/2026-09-30-recon-audit-fix.md:81`).
  - `--unverified`:**, `X-User-Id:12345`, `https://…/my-url%20here`, `run-until:5`, `--user-agent:Mozilla`.

  The command alternative also matches any line that names mysql and later has a ` -a`/` -p` with a value (`mysql is down; ran ls -a /tmp`). All of this errs on the safe side and costs only dream inputs, but none of these cases is pinned on the over list. **Remediation:** put a non-word class in front of the flag only, `(^|[^A-Za-z0-9_-])(?:-u|--user|--username)`. Quotes, backticks and spaces still precede the flag there, and the command names keep no anchor. Leave `+` out of the user's first character. Pin `date -u +%FT%TZ` as benign.
- **[low / confidence: medium]** `internal/dream/run.go:226-227` — Only SIGINT and SIGTERM are forwarded. If the terminal closes or SSH drops during a hand-run `sh-agent dream`, SIGHUP kills sh-agent by default. The kernel sends SIGHUP only to the foreground group, which the CLI has left, and the shell does not know the CLI's group as a job. The CLI then runs on as an orphan with the operator's credentials and no deadline, which is the round-1 gap again by another signal. SIGQUIT and a SIGKILL or OOM kill of sh-agent leave the CLI running the same way. **Remediation:** add `syscall.SIGHUP` and `syscall.SIGQUIT` to `signal.Notify`. Optionally, on Linux, set `SysProcAttr.Pdeathsig` for the direct child. It fires when the creating OS thread exits, so it needs `runtime.LockOSThread` around `Start`, and it does not reach helpers.
- **[low / confidence: medium]** `internal/dream/run.go:239-241` — A signal that arrives after `close(waited)` and before `signal.Stop(sigs)` is buffered in `sigs` and never read. The Ctrl-C is swallowed, and sh-agent goes on to the next project. The window is microseconds. **Remediation:** call `signal.Stop(sigs)` first and then `close(waited)`. After that, drain `sigs` without blocking, and re-raise any signal found (kill the group, `signal.Reset`, `Kill(self)`).

### Checked
- **Kill reach after `Wait`** (`run.go:242`): `Kill(-pid)` runs after the leader is reaped. If helpers are still in the group, Linux allocates no new PID equal to a live process-group ID, so the kill reaches exactly those leftover helpers. If the group is empty, the call returns ESRCH. The only way to hit something else is the PID counter wrapping to exactly this value and the new process making itself a group leader, all within the microseconds between `Wait` returning and the kill. Even then only same-uid processes can be hit (EPERM otherwise). This is information, not a finding.
- **Kill reach on the signal path:** `killGroup` (the CLI's own group) and `Kill(os.Getpid(), sig)` (sh-agent itself). Nothing outside the dream job's children. A helper that called `setsid` is outside the group and survives, by design (the round-1 "left the group" test); WaitDelay only cuts its pipes.
- **Side effects of the signal path:** `signal.Reset` does not clobber any other handler; `signal.` appears nowhere else in `cmd/sh-agent`, `internal/dream` or `report.go`. dream.lock is a `flock` (`run.go:650`), so dying of the re-raised signal releases it.
- **The rest of the run.go fix:** `exec.ErrWaitDelay` is treated as success only after the deadline check, and Go returns it only for a zero exit status. The `Start` error path is equivalent to the old `Run`.
- **Declarations and references:** the needle renamed from `docker auth` to `auth key` has no other reference in Go or the living docs. `call`, `callWithin`, `killGrace` and `cliTimeout` are all still used (`run.go:130,188-189`). `installPassArgs` (`main.go:1113`) is called by update and by `runInstall`'s rebuild; forwarding an explicit `--vendor all` is correct, with no security effect. The fix introduces no new exported identifier.
- **Tests:** `go test ./internal/harness ./internal/dream ./cmd/sh-agent -run 'Redact|TestCall|InstallPass'` passed in the scratch overlay. The interrupt test drives a real SIGINT through a helper process, and it would fail without the handler: the default exit kills only the test binary.
- **Other surfaces:** the `skills/harness-audit/SKILL.md:36` govulncheck path and the `search.md:29` exit-status text are docs only, with no security effect. No secret shape appears in the added lines; the test shapes use `FAKEVALUE`.

### Not examined (round 2)
- The full test suite, `-race`, and the claim that removing the three needles fails 46 assertion lines; I ran only the targeted tests above.
- The "red before" claims for the new dream and usage tests.
- `usage_test.go` beyond its role as a caller of `installPassArgs`.
- The `delivery.md` sentence beyond the removed `mcpServers` mention.
- govulncheck was not re-run, since the fix changes no dependency.
