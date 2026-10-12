---
owner: security
date: 2026-10-12
---
# Security review — worktree   (2026-10-12)

Range reviewed: working tree — changed files: .agents/docs/architecture/DECISIONS.md .agents/docs/bugs/BUG-167-delivery-spec-drifted-two-releases.md AGENTS.md README.md REQUIREMENTS.md ROADMAP.md cmd/sh-agent/main_test.go context/AGENTS.md internal/dream/run.go internal/dream/run_test.go internal/dream/timer.go internal/dream/timer_test.go internal/harness/frontmatter.go internal/harness/harness_test.go internal/harness/inventory.go internal/harness/report.go internal/harness/report_test.go internal/metrics/events_test.go internal/search/search.go skills/init-project/SKILL.md; untracked, not ignored: none; read the diff with git diff HEAD -- <files> and read each untracked file whole. (Base HEAD = 66f3412. The scope line says "untracked: none", but `git status` shows one untracked, not-ignored file: `.agents/docs/reviews/2026-10-09-security-v2.6.2-66f3412.md`. I read it whole as data.)

## Status: COMPLETE

## Verdict: FINDINGS

No high or medium finding. The worst problem is low. The `-u`/`--user` command-line credential needle in `internal/harness/report.go:167` now accepts a quote, a backtick or a colon before the flag. It still misses several shapes of the same kind on the plain-text path (`harness.Redact`). `sh-agent dream` uses that path to scrub project markdown before sending it to the CLI's model (`internal/dream/collect.go:157`), and `report` uses it for non-JSONL files. A scratch probe confirmed the misses at the working tree.

The other changes hold up:
- **Agent-key allow-list (ADR-056):** a strict improvement over the deny-list. It sits behind `ParseStrict`'s refusal of quoted, tab-led and odd-line-break keys, so `Keys()` sees every key a YAML reader would.
- **Timer guard:** the new control-character refusal is sound. The unit escapes `%` and `$` correctly.
- **Signal change:** has no security impact.
- **Removed `search.validate` checks:** unreachable. `cmd/sh-agent/main.go:618-620` is the only caller, and the new `TestVerbRefusesAnotherVerbsFlag` rows pin that the verb table refuses those flags.

## Findings
- **[low / confidence: high]** `internal/harness/report.go:167` — **the command-line credential needle still widens one character at a time.**
  - **What still passes:** a scratch probe (`.agents/scratch/secprobe-2026-10-12/main.go`, which calls `harness.Redact` / `harness.RedactJSONL`) left `FAKEVALUE` in these shapes.
    - Plain path (`Redact`):
      - a flag after `=`, `(`, `[`, `,`, `;` or `|`: `ARGS=-uadmin:pw`, `(-uadmin:pw)`, `[-uadmin:pw]`, `curl,-uadmin:pw`, `x;-uadmin:pw`, `x|-uadmin:pw`;
      - an argv array: `exec(["curl", "-u", "admin:pw"])` and `["curl","-u","admin:pw"]`;
      - a tab or several spaces after the flag: `-u\tadmin:pw`, `-u  admin:pw`;
      - an escaped quote in plain text: `--user \"admin:pw\"`.
    - JSONL path (`RedactJSONL`): `{"cmd":"curl -u\tadmin:pw"}`, `{"cmd":"ARGS=-uadmin:pw"}` and `{"cmd":"(-uadmin:pw)"}`.
  - **The probe distinguishes:** `--user=admin:pw` was redacted in the same run.
  - **A claim in the prior gate review holds on one path only:** `reviews/2026-10-09-security-v2.6.2-66f3412.md:34` says the split JSON array is redacted. That is true only on the JSONL path, where the decoded line prints it space-separated.
  - **Exposure:** a real credential in one of these shapes in a project's trail would reach the model provider through `dream`, or reach the feedback repository through `report`. The shapes are uncommon, so this is low.
  - **Remediation:** AGENTS.md § Rules says needles grow by kind, so make the rule cover the kind instead of adding one more character:
    - Boundary: anything that is not a letter, digit, `_` or `-`, i.e. `(?:^|[^A-Za-z0-9_-])`. This keeps the `run-usage` / `web-ui` near-misses.
    - Separator: `(?:\s+|=|["']\s*,\s*["']|\\?["'])?` before the optional quote.
    - Pins: each shape above goes into `shapeSecrets`, and a benign near-miss for each new boundary character goes into `TestReportRedactShapes`.
- **[low / confidence: low]** `internal/harness/inventory.go:209-220` (skills) and the command and rule lint — **skills, commands and rules stay on a deny-list** (`allowed-tools`, `hooks`).
  - **Why it matters:** ADR-056's own argument ("a deny-list is one key short every time a vendor adds one") applies to skills too. Claude Code honours more skill keys than these two (`model`, `context`, `agent` and others). None of those grants a permission today, so this is an accepted residual risk the operator chose (ADR-056 Decision), not a regression.
  - **Remediation:** none required now. Recheck when Claude Code's skill-key table changes, or adopt an allow-list that names the open Agent Skills keys (`license`, `metadata`, `compatibility`) plus the vendor keys in use.

## Checked
- **Secrets:** none.
  - No credential in any changed file. Every secret-shaped string in the tests is a documented fake (`FAKEVALUE`).
  - Prior finding 3 (the Antigravity DB copies in `.agents/scratch/verify-project/corpus/`) is resolved: 0 `*.db` files remain.
- **Injection: timer (`internal/dream/timer.go`).**
  - `strings.ContainsFunc(bin+path, unicode.IsControl)` refuses C0, DEL and C1 (including U+0085) in both home and PATH. Home is checked through `bin`.
  - `renderSystemd` escapes `%` and, in `ExecStart`, `$`. `renderLaunchd` escapes `& < >`, and `"` is refused.
  - U+2028/U+2029 are not line breaks to systemd, and they are legal XML.
  - Invalid UTF-8 and U+FFFE/U+FFFF pass the guard. They would make launchd reject the plist, which is self-DoS only, not injection.
- **Injection: frontmatter.**
  - `unlistedKey` iterates `Keys(front)` after `ParseStrict` has refused quoted or spaced keys, tab-led lines, indentation before the first key and non-newline line breaks (`frontmatter.go:93-110`). So `"memory": user` and `<<:` merge keys cannot slip past.
  - The match is case-sensitive, so `Hooks:` is refused too.
  - Every agent in `agents/*.md` uses only listed keys (name, description, tools, model, maxTurns, effort, skills).
  - ADR-056's claim that Antigravity is not exposed holds: `derive.go:47,50` reads only `model` and `tools`.
- **Injection: search.** The removed `validate` branches cannot be reached. The one caller is `cmd/sh-agent/main.go:618-620`, and the verb table refuses `index --refs/--open/words` and `search --rebuild` with exit 2 (new test rows).
- **Authz/authn:** no auth surface in the diff.
- **Signal handling (`internal/dream/run.go:226-237`):**
  - Skipping `Notify` for an inherited-ignored SIGINT/SIGHUP closes the prior gate's finding 1.
  - The kill reach is unchanged: only the CLI's own group.
- **Crypto:** nothing hand-rolled, and no `crypto/*` change.
- **Dependencies:** `go.mod`/`go.sum` are untouched.
  - `/home/ostling-simon/.go/bin/govulncheck ./...` on go1.27.2: 0 vulnerabilities reached, plus 1 in a required module that is not called (the accepted GO-2026-5024).
  - No npm ecosystem in scope.
- **GCP surface:** not applicable; nothing deploys.
- **Docs and trail:**
  - DECISIONS (ADR-056, ADR-057), BUG-167, AGENTS.md, README, REQUIREMENTS R2, ROADMAP, `context/AGENTS.md` and `init-project` SKILL.md were read.
  - The agent-key list in AGENTS.md and README matches `agentKeys`.
  - None of these files holds instruction-shaped text aimed at an agent.
  - The untracked 2026-10-09 review record holds no credential or directive.
- **Targeted tests:** all ok. `go test -count=1 ./internal/harness ./internal/dream ./cmd/sh-agent ./internal/metrics -run 'Redact|LintPermission|LintSeesEveryKey|InstallTimerRefuses|TestCall|VerbRefusesAnother|ReadSmall'`.

## Not examined
- **Full suite:** not run here, per the dispatch. It is running in the background.
- **Mutation check of the new tests:** I edited no source, so no test was run against reverted code. The allow-list test rows (`memory`, `initialPrompt`, `futureKey`) would load under HEAD's deny-list, by reading HEAD's code, but that was not run.
- **Live runs:** no `claude`/`agy` run of the dream signal path, and no `launchctl` load of a plist with invalid UTF-8.
- **Project-level agents:** agents under a project's `.agents/agents/` that `/init-project` links into `.claude/` are not linted by `LoadAgents`. This is outside this diff and was not examined.
