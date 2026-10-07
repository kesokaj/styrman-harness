---
owner: security
date: 2026-10-07
---
# Security review — worktree-ab-kit   (2026-10-07)

Range reviewed: HEAD (f82457b)..working tree. Diffs read: `agents/security.md` and `.agents/docs/testing/guide.md`. Read whole: the untracked `.agents/docs/testing/ab-kit/{run.sh,accept.py,metrics.py,prompt.md,prefix-full.md}`. Read as blast radius: `.agents/docs/testing/kit-preamble.sh`, `/etc/profile.d/zz-styrman-credentials.sh`, `/etc/profile.d/zzz-styrman-motd.sh`, `/usr/local/bin/styrman-git-credential` (lines 1-16 plus a grep of the rest), and the BUG-160 working-tree diff.

## Status: COMPLETE (rounds 1 and 2)

## Verdict: FINDINGS. Final, after round 2: low findings only, none medium or high open (see ## Round 2)

The worst problem: the kit isolates the headless agent only by moving `HOME` and setting two environment switches. The agent runs with `--permission-mode auto`, and the prompt tells it to approve its own choices. It still has three things it should not have:
- the operator's live `GH_TOKEN`, plus the parent Claude Code session's messaging token and socket. These are inherited, because `run.sh:47` removes only one variable from the environment.
- passwordless root. `sudo -n -l` on this box shows `(ALL) NOPASSWD: ALL`.
- network access to the credential broker. The broker serves every credential to any caller in the container, and its clients send no credential of their own.

The only thing between the agent and these is Claude Code's auto-mode classifier. The broker switches do one thing well. I read every script that could write: the profile script, the MOTD script and the git helper. With the switches set, none of them writes a file, which fixes BUG-160's leftovers. But the switches are not a security boundary.

There is a second problem. `accept.py` builds and runs code the model wrote, and it passes that code the operator's full environment and real `HOME`. `metrics.py` runs `git` in the repository the model wrote, with the same environment.

## Findings

- **[medium / confidence: high]** `.agents/docs/testing/ab-kit/run.sh:47` — **The environment is a denylist of one, so the operator's tokens reach the unattended agent.**
  - **What leaks.** `env -u STYRMAN_CRED_BROKER` stops the broker's *file* tier inside the child. The *env* tier has already run, in the login shell that launches the kit (`zz-styrman-credentials.sh:43-57`), and everything it exported flows through to the child. These names are present in this box's environment (I read names only, never values):
    - `GH_TOKEN`
    - from the launching Claude Code session: `CLAUDE_CODE_MESSAGING_TOKEN`, `CLAUDE_CODE_MESSAGING_SOCKET`, `CLAUDE_CODE_SESSION_ID`, `CLAUDE_CODE_CHILD_SESSION`, `CLAUDECODE`, `CLAUDE_EFFORT`
    - `STYRMAN_GIT_NAME` and `STYRMAN_GIT_EMAIL`

    `GIT_DIR`, `GIT_WORK_TREE` and `GIT_INDEX_FILE` are not cleared either. If the caller ever exports them, the agent's `git add` or `checkout` would act on another repository or work tree.
  - **Why it matters.** The agent can use the token through `gh` or `curl` without any turn from the operator. It may be able to post into the operator's attended session through the inherited messaging handle; that impact is not verified, and the inheritance is certain. A dump of the environment written outside `KIT_HOME` is never scanned: `kit-preamble.sh:130` scans `KIT_HOME` only, not the clone or `$RUN/{result.json,stderr.log}` under `~/sandbox/ab`. BUG-160's own addendum of 2026-10-05 records a CLI writing `GH_TOKEN` to disk.
  - **The record misses it too.** BUG-160's new Prevention line ("turns the broker off for the CLI it runs") does not cover the env tier.
  - **Remediation.**
    - Start the CLI from `env -i` with an allowlist: `PATH`, `LANG`/`LC_*`, `TERM`, `HOME`, `CLAUDE_CONFIG_DIR`, `STYRMAN_SKIP_CRED_PULL=1`, the `GIT_AUTHOR_*`/`GIT_COMMITTER_*` identity, and `GOCACHE`/`GOMODCACHE` only if they are kept (see the next finding).
    - Extend the exit scan to `$RUN` and the clone.
    - Amend BUG-160's prevention rule to say "environment allowlist", not "broker off".

- **[medium / confidence: high]** `.agents/docs/testing/ab-kit/run.sh:44-49`, `prompt.md:21`, `guide.md:181` — **The "sandbox" has no confinement, and the guide now calls it the sandbox procedure.** The agent runs as the operator's uid with:
  - **Root.** Passwordless `sudo`, confirmed.
  - **The real home.** Its working directory is under the real home (`AB_OUT` defaults to `~/sandbox/ab`, `run.sh:28`). Every other checkout and `~/.claude/.credentials.json` are reachable by absolute path.
  - **The real Go caches.** `GOCACHE` and `GOMODCACHE` are exported into the run (`run.sh:44`). The build cache trusts its own entries, so anything written there feeds the operator's later builds, including the `go build` of `sh-agent`.
  - **The broker.** It is reachable on the container network. Its clients send no credential (`zz-styrman-credentials.sh:44,61`, `styrman-git-credential:29`), and its address is derivable from the subnet (`zz-styrman-credentials.sh:19-21`). So the switches stop automatic materialization, not access.
  - **Self-approval.** `prompt.md:21` tells it to take its own recommendation wherever it would ask. That turns every operator-consent gate in the installed context into self-approval.
  - **No deny rules.** `CLAUDE_CONFIG_DIR` moves the config away from the operator's real settings, and the `settings.json` written at `run.sh:38` carries no `permissions`.

  **Why it matters:** the run is unattended and costs up to $8, and the guide now endorses it. A model that misbehaves, or content it reads that injects instructions, has root and the operator's credentials.

  **Remediation:** take the env allowlist above, then layer:
  1. In the `settings.json` that `run.sh:38` writes, add `permissions.deny` for `Bash(sudo:*)`, `Bash(gh:*)` and reads of the real home. This is cheap, though it can be bypassed.
  2. For real confinement, use either Claude Code's own sandbox (`sandbox.enabled`) or a dedicated uid without sudo. The Claude Code sandbox needs `bwrap`, which is not installed here (`socat` is). Either one needs a root step from the operator.
  3. Point `GOCACHE` at `KIT_HOME`, or record that the real caches are shared on purpose.

- **[medium / confidence: high]** `.agents/docs/testing/ab-kit/accept.py:6,11`; `metrics.py:32,37-38` — **The scoring scripts execute agent-controlled code with the operator's full environment and real `HOME`.**
  - `accept.py` runs `go build` on the module the agent wrote. That module's `go.mod` can use `replace` to point at any local path. It then runs the binary with `env={**os.environ, 'PORT': ..., 'DATA_FILE': ...}`, which includes `GH_TOKEN` and the parent session's messaging handle. That code runs outside even the kit's `HOME` redirect.
  - `metrics.py` runs `git rev-list` and `git ls-files` (`shell=True`, constant strings, so no injection) inside the agent's clone. The clone's `.git/config` is agent-controlled, and a `core.fsmonitor` value there runs on commands that read the index. That second point has medium confidence.

  Remediation:
  - Run the built binary with a minimal environment: `PATH`, `PORT`, `DATA_FILE`, and `HOME` set to the work dir.
  - Run git with `-c core.fsmonitor=false -c core.hooksPath=/dev/null`, `GIT_CONFIG_NOSYSTEM=1` and a scratch `HOME`.
  - Better still, run both scripts under the same env allowlist as the CLI.

- **[medium / confidence: medium]** `.agents/docs/testing/ab-kit/run.sh:41-54` (and `:20`) — **Stopping the kit does not stop the agent, and cleanup is not on the signal path.**
  - **How a stop goes wrong.** GNU `timeout` 9.7 without `--foreground` moves itself and `claude` into their own process group.
    - Ctrl-C at a terminal does not reach `claude`. bash waits for its child, so the interrupt is deferred for up to 90 minutes.
    - A plain `kill` of `run.sh` signals none of its children. The preamble's EXIT trap then reports the link as left over (`kit-preamble.sh:143-151`) while the agent keeps running unwatched. `run.sh:54`, the only line that removes the link, never runs.
  - **Other exits.** Any `set -e` exit between `:41` and `:54` (`:42-44`) also leaves the link.
  - **A full copy can be left behind.** If the CLI replaced the link with a regular file (the case `:53` warns about), a kill leaves a full copy of the refreshed credential in scratch.
  - **Untested.** Not tested here: how `TaskStop` kills a background Bash, and bash's trap order on SIGTERM.
  - **Remediation.**
    - Create the link immediately before the run, after `:42-44`.
    - Start `timeout … claude` in the background and record its PID, which is also its PGID.
    - Add a `trap` on INT, TERM and HUP that runs `kill -TERM -- -<pgid>`, then `rm -f` the link. The guide forbids only a kit's own EXIT trap.
    - Alternatively, give the preamble a delete-by-name list that its EXIT trap removes before it checks.

- **[low / confidence: high]** `.agents/docs/testing/guide.md:181` — **The changed sentence contradicts itself.** It still opens with "The real home is never touched", then endorses the A/B kit's symlink. That symlink exists so that the CLI writes the real `~/.claude/.credentials.json` through it (`run.sh:13-14`). The kit also writes its clone and results under the real home (`~/sandbox/ab`).
  - **The shared token.** The kit and the operator's live sessions use one refresh token, so a refresh by one can invalidate the other. This is an availability risk, not a confidentiality one.
  - **Remediation:** state both exceptions and the shared-token consequence where the sentence makes its claim.

- **[low / confidence: medium]** `agents/security.md:6` — **The `effort: max` → `high` change needs no compensating text, but its risk is not measured.**
  - **No compensating text needed.** No rule text depends on effort. A grep finds "max" only in ADR-054a and the plan, which record the change. `derive.go:30` already drops `effort` for Antigravity, so Antigravity runs are unchanged.
  - **The measure misses the risk.** The plan measures the risk as "the re-run's security still reports, and the acceptance test stays 25/25" (`plans/2026-10-07-lean-process-2.md:46`). That cannot detect a missed finding, because the acceptance test scores the app, not the review.
  - **The behavioural run is n=1.** It is D7's `lean2-full` run, launched with the launcher's `CLAUDE_EFFORT` inherited (see the first finding). I have not checked whether that variable affects a subagent; low confidence.
  - **Remediation:** a seeded-weakness check. Plant a known weakness in a fixture diff, dispatch `security` at `high`, and report it as n of N found. Also clear `CLAUDE_EFFORT` through the allowlist.

## Checked
- **Secrets:** no credential, token or key in any file in the range; the files hold only the model ID and placeholder emails.
  - **Leftovers now.** No `.credentials.json` link remains in any arm. BUG-160's four files are gone. A grep for token shapes over `~/sandbox/ab` and `.agents/scratch/ab-homes` (text files, paths only) found nothing.
  - **Remaining non-credentials.** `ab-homes/control-1/.gitconfig` holds the operator's name and email, written by the broker before the fix. Go telemetry `upload.token` files are in every arm. Neither is a credential.
- **Broker off-switch, file writes:** verified complete on this box.
  - `zz-styrman-credentials.sh:17` returns before its first write when `STYRMAN_SKIP_CRED_PULL` is set. With only the broker unset it reaches `:32`, which only removes the marker.
  - `zzz-styrman-motd.sh` writes nothing.
  - `styrman-git-credential:14` exits when `STYRMAN_CRED_BROKER` is unset.
  - `/etc/gitconfig`'s helper is the only git-level path.
  - Reach is the second finding, not this line.
- **Injection:**
  - `run.sh` quotes every expansion, and the prompt is passed as a single argument.
  - `metrics.py` uses `shell=True` only with constant strings, and its output goes through `json.dumps`.
  - `accept.py` uses argument lists.
  - `run.sh:24` accepts any `*-full` arm, including `/` and `..`, so `KIT_HOME` and `RUN` can leave scratch. Only the operator types that argument, and `:31` refuses an existing path, so this is not a boundary and not a finding.
  - No file in the range holds text addressed to the reviewer.
- **Authz/authn:** the agent's permission gate is the auto-mode classifier alone; see the second finding.
- **Crypto:** nothing hand-rolled; none in the range.
- **Dependencies:** none added. The scripts use the Python stdlib (run with `-I`), bash, `curl` and `timeout`.
  - `govulncheck` is not on PATH (`command -v` failed), so it did not run. The range changes no Go code and no `go.mod`.
  - `npm audit` does not apply; there is no JS ecosystem in the range.
- **GCP surface:** no IAM, service account or deploy change. gcloud ADC no longer lands in the sandbox. The broker still serves it to any process in the container (second finding).

## Not examined
- **Not read end to end:** `.agents/docs/testing/guide.md` beyond the changed line `:181` and the hunk around it; the rest of `/usr/local/bin/styrman-git-credential` past line 16 (only grepped); `plans/2026-10-07-lean-process-2.md` and ADR-054a beyond the grep hits.
- **Claude Code internals, not verified:**
  - what the `auto` permission mode allows under `-p`
  - what a holder of `CLAUDE_CODE_MESSAGING_TOKEN`/`SOCKET` can do
  - whether the CLI writes credentials through the link or replaces it
  - whether `--strict-mcp-config` excludes the operator's claude.ai account connectors (one is live in this session). To check: `--output-format stream-json`, then read the `mcp_servers` list in the run's `init` event.
- **Process behaviour not tested:** how `TaskStop` kills a background Bash, and bash's EXIT trap on SIGTERM. Both come from documented behaviour, not a probe.
- **Not read, only scanned for token shapes:** the three runs' `result.json`, `install.log` and clones, and the installed harness under `ab-homes/lean2-full`.
- **Broker server side:** not visible from here; only its clients were read.
- **Tools:** `govulncheck` did not run (not on PATH).

## Round 2

Status: COMPLETE.

**Re-read:**
- whole, after the fixes: `ab-kit/run.sh`, `ab-kit/accept.py`, `ab-kit/metrics.py`
- the `guide.md` diff (now `:158` and `:181`)
- `guide.md:185`, the accepted-design sentence the rebuttal cites. It is at HEAD.
- `BUG-160:23,48` and the rows of `skills/VENDORED.md`

**Evidence gathered:**
- In the three runs' transcripts, tool names only. No WebFetch or WebSearch call in any run. `lean2-full` loaded only `feature`, `init-project` and `handoff`.
- A scratch probe at `.agents/scratch/sec-probe-fsmon-114822` (git 2.47.3). Its hook only appends a line to a marker file. The directory is gitignored; delete it by name.

### Final verdict: FINDINGS. Low only; no medium or high is open.

The token exposure is closed:
- `run.sh:50-53` starts the CLI from `env -i` with an explicit allowlist. Nothing token-bearing passes: no `GH_TOKEN`, no `CLAUDE_CODE_*`, no `CLAUDE_EFFORT`, no `GIT_DIR` family, no broker URL. `STYRMAN_SKIP_CRED_PULL=1` is kept, so a login shell inside the run hits `zz-styrman-credentials.sh:17`, and `styrman-git-credential:14` exits without the URL. Nothing is re-exported.
- The real Go caches are no longer exported.
- `accept.py` builds and runs the agent's code with `MINENV` only.
- `metrics.py`'s git gets no tokens.

I accept the no-confinement rebuttal for text the harness wrote. It rests on the operator's recorded acceptance at `guide.md:185`, which predates this range. Because that acceptance covers the Docker socket, which is root-equivalent, sudo and the broker fall inside it.

What remains is low and the operator's call: hardening (R2-1 to R2-4), two doc fixes, and one deferral.

### Round-1 findings
1. **Env inheritance (medium): fixed.** For the agent's environment, see above.
   - **Not done: the scan.** The exit scan still covers only `KIT_HOME`. The one secret left in reach is the Claude credential through the link, and a copy of it written into the clone or `$RUN` would not be scanned. Low now.
   - **Not done: the record.** `BUG-160:48`'s Prevention still names only the broker switch (R2-5).
2. **No confinement (medium): rebutted, accepted in part.** `guide.md:185`: "A sandbox home is a `HOME` override, nothing more: the CLI under test keeps the operator's uid, environment, Docker socket and reachable checkouts". `run.sh:22-23` now says the kit is not a confinement. The same sentence limits itself to "text the harness wrote" (R2-1).
3. **Scoring scripts get the full environment (medium): fixed for tokens.** See `accept.py:5,8,13` (`MINENV`) and `metrics.py:37-38` (`GITENV`). The `curl` subprocess is replaced by `urllib` (`accept.py:28-33`). The git-config execution half is still open (R2-2).
4. **Stop path (medium): partly fixed.**
   - **Fixed:** `timeout --foreground` keeps `claude` in `run.sh`'s process group, so Ctrl-C or a kill of that group reaches it.
   - **Still open, accepted:** cleanup stays off the signal path, because `run.sh:58` runs only on a normal return. A run stopped by a signal leaves the credential link, and the preamble FAILs on it by name. `run.sh:24` documents this as accepted. The link grants nothing the uid lacks.
   - **New:** R2-3.
5. **`guide.md:181` "The real home is never touched" (low): still open.** It is now contradicted inside the same section by the new `:158`, which says the kit "links the credential … so a token refresh lands in the real file".
6. **Effort `high` (low): deferred by the operator.** The seeded-weakness check is in the handoff. The `CLAUDE_EFFORT` inheritance it named is fixed by the allowlist.

### New in round 2 (worst first)
- **R2-1 [low / confidence: high]** `run.sh:53` against `guide.md:185` — **The A/B kit runs outside the acceptance it cites, for third-party text.**
  - **What the acceptance says.** It "covers text the harness wrote". A row that loads vendored or third-party text "runs with read-only grants … never `command(*)`".
  - **What the kit does.** It runs `--permission-mode auto`, not read-only grants. It does not disallow WebFetch or WebSearch. Its harness arms install third-party skills verbatim: `grilling` and `codebase-design` (mattpocock/skills), `find-skills` (vercel-labs/skills), `gke-basics` (derived from google/skills) and `google-catalog/catalog`. All are present in `ab-homes/lean2-full/.claude/skills`.
  - **Not exercised.** None of this text was loaded in the three runs on disk, so this is a policy gap, not an incident.
  - **Remediation.** Disallow WebFetch and WebSearch (`--disallowedTools` is variadic like `--allowedTools`; mind the order, guide § 4). Then either have the operator extend `:185` to the A/B harness arms in one sentence, or have `metrics.py` count Skill and web-tool uses and fail a run that loaded third-party text.
- **R2-2 [low / confidence: high]** `metrics.py:37-38` — **`GITENV` removes the tokens, but git still runs agent-controlled config.** Verified both ways by the probe:
  - A repo-local `.git/config` `core.fsmonitor` ran on `git ls-files -co --exclude-standard`, twice. `git rev-list --count HEAD` did not run it.
  - The fix's `HOME`=the clone makes a working-tree `<clone>/.gitconfig` a global config. Its `core.fsmonitor` also ran twice. The fix added this second agent-controlled source.
  - `git -c core.fsmonitor=false …` ran it zero times in both cases.

  The impact equals what `accept.py` already does by design, which is to run the agent's binary under `MINENV`. But this execution happens after the run, outside the auto-mode classifier and after the preamble's exit checks. Remediation: add `-c core.fsmonitor=false` to both commands, and set `HOME` to a scratch directory outside the clone, or set `GIT_CONFIG_GLOBAL=/dev/null`.
- **R2-3 [low / confidence: high on the change, low on the impact]** `run.sh:53` — **`--foreground` changes what the 90-minute timeout kills.** `timeout --help` says: "in this mode, children of COMMAND will not be timed out". Only `claude` gets SIGTERM now; before, timeout signalled the whole group.
  - **Why it matters.** A server or test the agent left running would survive the run unattended. The prompt's port default is 8080, and a Go `":"+port` listener binds every interface.
  - **Not verified.** Whether Claude Code reaps its children on SIGTERM.
  - **Remediation.** After the subshell returns, kill what is left that runs in the clone, or record the gap beside `:60`.
- **R2-4 [low / confidence: high]** `accept.py:5` — **`GOFLAGS=-mod=mod` fails open.** It lets `go build` download modules the agent's code imports but `go.mod` does not declare. It can also rewrite `go.mod`/`go.sum` in the clone before `metrics.py` counts its files. The prompt allows the standard library only. Remediation: use `-mod=readonly` with `GOPROXY=off`. Any third-party import then fails the build, and the build stays offline. Downloads would be checksum-verified anyway, so this is about scope and integrity, not tampering.
- **R2-5 [low / confidence: high]** `guide.md:158`, `BUG-160:48` — **The allowlist is advisory everywhere else.** The guide says "best from `env -i`", so for every other kit (skill-matrix, agy sandboxes) the allowlist is a suggestion. The 2026-10-05 addendum records Antigravity writing `GH_TOKEN` to disk, and BUG-160's Prevention line still names only the broker switch. Because `:185` accepts that the CLI "keeps the operator's … environment", this is the operator's call: make it "starts from `env -i` with an allowlist", or leave it advisory knowingly.

### Checked in round 2
- **Allowlist contents:** `PATH` (with the preamble's `--home` wrapper first), `TERM`, `LANG`, `USER`, `LOGNAME`, `SHELL`, `STYRMAN_SKIP_CRED_PULL=1`, `HOME`, `CLAUDE_CONFIG_DIR`, and the git identity. None of these carries a token.
- **Secrets in the changed files:** none.
- **`run.sh:22-24` header:** accurate against `guide.md:185`.

### Not examined in round 2
- **Not reproduced:** no kit run; no re-run of the coordinator's environment count (3 vs 0) or the control run (25/25). I took both as reported.
- **Claude Code signal behaviour:** whether it dies by SIGINT or exits normally (which decides whether Ctrl-C leaves the link), and whether it reaps background children on SIGTERM.
- **Operator intent:** whether `guide.md:185` was meant to cover the A/B harness arms.
- **Tools:** `govulncheck` is still not on PATH; there is no Go code in the range.

## Outcome

Written by the main agent, 2026-10-07. Round 2 resumed the same agent with `SendMessage`: no medium or high finding is still open.

**Fixed:**
1. **The agent's environment.** The CLI starts from `env -i` with an allowlist. Verified both ways: 3 token-like variable names without it, 0 with it. A `control` run authenticated under it, scored 25/25 and printed both clean lines.
2. **The real Go caches** are no longer exported to the agent (`ab-kit/run.sh`).
3. **`accept.py` and `metrics.py`** run with a minimal environment.

**Rebutted:**
- **No confinement** (sudo, the real home reachable by path, the broker reachable). This is the operator's accepted design for text the harness wrote (`guide.md` § Behavioural runs, 2026-10-01). Round 2 accepted it.

**After the cap.** These were fixed by the main agent and not re-reviewed:
- **R2-1:** `--disallowedTools WebFetch WebSearch`. The re-run's transcript shows no web tool use.
- **R2-2:** `metrics.py`'s git runs with `HOME=/nonexistent`, `GIT_CONFIG_GLOBAL=/dev/null`, `-c core.fsmonitor=false` and `-c core.hooksPath=/dev/null`.
- **The build:** `-mod=readonly` with `GOPROXY=off`.
- **Two sentences:** `guide.md` line 181, and BUG-160's prevention line.
- **The probe directory** `.agents/scratch/sec-probe-fsmon-114822` was deleted by name.

**Left:**
- **R2-1's grant question** (read-only grants for runs that load vendored skills): for the operator, in the handoff.
- **R2-3:** a server the agent started survives a timeout (in the kit header).
- **A signal leaves the credential link.** The preamble's exit check flags it by name.
- **A seeded-weakness check of `security` at effort high:** for the operator, in the handoff.
