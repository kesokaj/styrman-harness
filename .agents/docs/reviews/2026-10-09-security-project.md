---
owner: security
date: 2026-10-09
---
# Security review — project   (2026-10-09)

Range reviewed: the whole repository at HEAD (c267990), no diff.

## Status: COMPLETE

## Verdict: FINDINGS

The worst problem: a harness agent file can still carry `mcpServers:`. Claude Code's subagent format documents that key (`.agents/docs/plans/2026-09-30-research-audit-fix-formats.md:12`), and the installer deploys the agent file byte for byte to `~/.claude/agents/` on every machine (`internal/harness/plan.go:156`). That gives a harness agent a way to start a process that the operator's permission mode never gates. The BUG-143 lint refuses only `hooks` and `permissionMode` for agents. The 2026-10-01 gate recorded this gap and left it open (`.agents/docs/changes/2026-10-01-1433-gate-lows-closed.md:18`), and it is still open at HEAD. The rest is low:
- the redaction test does not meet the project's own near-miss rule for every needle;
- a few credential kinds are not covered;
- the systemd unit text is not escaped against operator-supplied quotes and newlines;
- ADR-053's accepted residual risk is restated for the record.

The zip extraction, the release update, the SQL construction, the frontmatter line-grammar lint and the dream job's tool-free invocation held up in this reading.

## Findings
- **[medium / confidence: medium]** `internal/harness/inventory.go:261`: `permissionKey(front, "hooks", "permissionMode")` does not refuse `mcpServers`. An agent can declare an inline stdio MCP server, which is a command line that Claude Code starts when the agent is used. `plan.go:156` ships `a.Raw` unchanged. For Antigravity the risk is lower: `derive.go:47-50` writes only `model` and `tools`, so the key never reaches agy. Starting an MCP server is not a tool call, so the operator's permission mode never sees it. This is the confused-deputy shape BUG-143 closed for skills, reachable by anyone who lands a commit or a hostile zip in the harness tree, and then active on every place. The confidence is medium because whether Claude Code honours an inline `mcpServers` in a user-level agent file was not tested here (the 2026-10-01 review marked the same point low). Remediation: add `mcpServers` to the agent refusal list, and prefer an allow-list of the agent keys the harness uses (`name`, `description`, `tools`, `model`, `skills`, `color`, …) over a deny-list. Pin it with a row in `TestLintPermissionKeysByKind` that fails without the fix.
- **[low / confidence: high]** `internal/harness/report.go:97-100` against `internal/harness/report_test.go:161-176`: AGENTS.md requires "every needle gets a benign near-miss in `TestReportRedactShapes`". `bearer`, `authorization`, `url credentials` and most prefix needles (`AKIA`, `AIza`, `ghp_`, `glpat-`, `xoxb-`, …) have no benign near-miss. Only `ASIA` (as `ASIAN food`) and the regexp needles do. `bearer\s+\S` matches prose such as "the bearer of the token", so a whole line is dropped. That is the safe direction, but the behaviour is unpinned, so a later tightening of it is unreviewed. Remediation: add a benign near-miss per needle, with the over-redacting ones (`the bearer of`) pinned on the `over` list as `ASIAN food` is.
- **[low / confidence: medium]** `internal/harness/report.go:59-140`: these credential kinds match no needle, and the same rules are what `sh-agent dream` uses before it sends text to a model provider (`internal/dream/collect.go:160`):
  - CLI credential flags: `curl -u user:pass`, `--user user:pass`, `mysql -p<pw>`, `sshpass -p`.
  - A JSON `"auth": "<base64 user:pass>"` key, the shape of Docker's `config.json`.
  - A `credential`/`credentials` key (`CREDENTIALS=…`).

  Each is a kind, not an instance. Under the project's BUG-110 rule they are a design question for the operator, not a batch to add here. Remediation: decide per kind; if accepted, add one rule per kind with a near-miss each.
- **[low / confidence: high]** `internal/harness/release.go:180-197`, `:349`: the release binary and the `SHA256SUMS` it is checked against come from the same release, fetched by the same `gh release download`. The check proves the download is whole and belongs to that release, not who made it (ADR-053, stated at `:26-31`). Anyone with write access to `kesokaj/styrman-harness` releases (the account or a `contents:write` token) can publish a higher `vX.Y.Z`. `sh-agent update` on every release place then installs it and executes it (`:144`). The same holds for a checkout place's `git pull` followed by `go build` (`repo.go:280`, `:308`). This is the accepted model, recorded so the residual risk is on file, not as a defect against the ADR. If the threat model changes, a vetted signature check is needed (`gh attestation verify`, or `ssh-keygen -Y verify` against an embedded allowed-signers file); never a hand-rolled one.
- **[low / confidence: low]** `internal/dream/timer.go:86-99`: `renderSystemd` escapes `%` and `$` in the binary path, `--home` and PATH. It does not refuse or escape `"`, `\` or a newline. A newline in `--home` or in a PATH entry would add arbitrary directives, such as an extra `ExecStart`, to a user unit that runs nightly. A `"` ends the quoted word early. Every input belongs to the operator (their own `--home` argument and their shell's PATH), so this is self-inflicted unless something else sets PATH. Remediation: refuse control characters, `"` and `\` in all three before writing the unit, with a test per character.
- **[info]** `internal/harness/report.go:351-365`: `Collect` calls `Lstat` and then `os.ReadFile`. This is a small TOCTOU in which a file swapped for a link between the two calls is followed, and the read has no size cap. Both lie inside the operator's own `~/.config/sh-agent`, so there is no trust boundary. `ReadCapped` plus `O_NOFOLLOW` would make it uniform with the rest of the code.
- **[info]** `internal/dream/run.go:303-311`: on macOS the whole `~/Library/Keychains` folder is linked into agy's private HOME. Keychain item ACLs still apply per application. The child runs in plan mode with `request-review` and an empty allow list (`:316`), and any tool step voids the answer (`:228-231`). That is the same exposure as an interactive agy session. No finding.

## Checked
- **Secrets:** no credential in the tree was found in this reading. The installer refuses secret-shaped files in skills (`plan.go:120`, `secretShaped` at `:216`). The report's clone and push errors are scrubbed of URL credentials (`report.go:645-650`). Nothing in `cmd/` or `internal/` uses `log/slog` (grep: 0 calls), so there is no structured log of secrets. The dream job's private HOME and its database copies are 0600/0700 and swept after an hour (`run.go:377-398`).
- **Injection:** SQL is parameterised throughout. The only concatenations are constant table names (`reconcile.go:660`, `metrics/ingest.go:145`) and generated `?` lists (`search.go:77`). The SQLite DSN is escaped through `url.URL` against `?#%` pragma injection (`search/store.go:76`, `dream.go:40`). Release tags are validated before they reach argv (`release.go:86`). Every `git` call in the project root goes through `trail.Command` (prior review). The report refuses `GIT_DIR`-family redirects (`report.go:790`) and pushes to a constant URL.
- **Zip and update:** `update --from-zip` (`zip.go`) hashes and parses one buffer, so verify and unpack see the same bytes and there is no TOCTOU. Entries with `..`, absolute paths, backslashes or non-local names are refused, with depth and length bounds. Symlink and special entries are refused, writes use O_EXCL, sizes are counted by decompressing against per-file, total, entry and archive caps, and only the five trees are written. The release update has no downgrade or sideways move, the version marker binds the SUMS file, ambiguous hash lines are refused, the download goes to a random temp directory, and the binary is replaced by an atomic rename of the hashed bytes. `go build` runs with `GOTOOLCHAIN=local` (`repo.go:319`).
- **Authz (skill and agent lint):** `frontmatter.go` was read for flow maps, `<<` merge keys, `...` and `%` directives, BOMs, `key:value` without a space, tab-led and indented-first lines, and non-newline breaks. Each is refused or is not a top-level key to a YAML reader. Skills refuse `allowed-tools` and `hooks`, and symlinks inside skills are refused (`plan.go:108`). The agent gap is above.
- **Crypto:** only `crypto/sha256` as a checksum, plus `sha256` for a unit-name suffix. Nothing is hand-rolled.
- **Dependencies:** `go.mod` has `modernc.org/sqlite v1.46.1` and its indirect modules, as AGENTS.md permits. govulncheck is NOT on PATH (`which govulncheck` empty, no `~/go/bin`), so it was not run and the accepted-advisory list was not re-confirmed. There is no npm ecosystem. A vendored-skill grep for exfiltration and `| bash` found only catalog entries whose `curl | bash` lines `skills/google-catalog/catalog/MANIFEST.md:62-68` already flags.
- **GCP surface:** this repository deploys nothing to GCP and has no IAM bindings, service accounts or `--allow-unauthenticated`. Not applicable.
- **Dream:** the Claude call has no tools, `--strict-mcp-config`, `--setting-sources local` and no session persistence (`dream.go:181`). Every input is redacted before the prompt (`collect.go:131-134`, transcripts decoded before line redaction). Drafts are written to the proposals folder only and never installed. The lock is a non-blocking flock (`run.go:599-603`).
- **Memory:** the store must be gitignored before a write (`memory.go:131`), and writes are 0600.

## Not examined
- `internal/harness/pathguard.go` (the `--home` guard), `owned_unix.go`, `sweep.go`, `adopt.go`, `source.go`, and `apply.go` beyond its symlink checks.
- `internal/search/reconcile.go`'s `os.Root` use: the GO-2026-4970 Lstat guard was not re-verified. Also not read: `search/parse.go`, `export.go`, and `search/dream.go` beyond its DSN.
- `internal/metrics` readers (`agy.go`, `claude.go`, `ingest.go`) beyond the SQL grep, and the HTML report (`metrics/report.go`), where output escaping is unchecked.
- `cmd/sh-agent/main.go` beyond the `--from-zip` wiring, and `report.go`'s `RemoveKeptClone` and `writeAndPush`.
- The release procedure: `skills/preflight/SKILL.md` and `.agents/docs/deploy/remote.md` § Releases were not read.
- Pins in `skills/VENDORED.md` were not verified (no `git archive` diff), and `go.sum` was not checked.
- govulncheck was not run (not installed).
- The launchd plist renderer beyond its XML escaping.
- No test file was read except parts of `report_test.go`.
