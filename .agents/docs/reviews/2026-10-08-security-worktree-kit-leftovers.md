---
owner: security
date: 2026-10-08
---
# Security review — worktree-kit-leftovers   (2026-10-08)

Range reviewed: HEAD (280c1df)..working tree, uncommitted. Files: `.agents/docs/testing/kit-preamble.sh`, `internal/harness/catalog_test.go`, `.agents/docs/testing/{ab-kit,agy-rel,seed-security}/run.sh`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/remote.md`, `.agents/docs/bugs/BUG-163-kit-run-left-a-server-listening.md`.
The tree changed while I was reviewing it. `kit-preamble.sh` went from 253 to 255 insertions: the environ match changed from a pipe into `grep -q` to a `penv` variable. This review covers the state at my last read: `kit-preamble.sh` lines 133-167 as quoted below, and `git diff HEAD --stat` showing 255 insertions and 61 deletions.

## Status: COMPLETE

## Verdict: FINDINGS (round 2: low only — the medium and every round-1 low are fixed; see ## Round 2)

The worst problem: `kit_process_leftovers` trusts its two roots without checking them. It sends SIGTERM to every same-uid process whose physical cwd is under `KIT_CLONE` or `KIT_HOME`, or whose executable is under either. The only guard on those roots is a literal string comparison at `kit-preamble.sh:32` (`"$HOME"|"$HOME/"`). So a `KIT_HOME` of `$HOME/.`, `$HOME//`, a symlink to the home, `/home`, or `/` gets through. `pwd -P` at line 143 then resolves it to a broad directory, and the exit trap terminates most of the operator's session: editors, shells, language servers, the browser. With `/`, the exe match hits every process. Only the trap's own ancestors are spared. The three kits in this change are not exposed today: each one derives its roots from `$OUT/homes/<label>` and `$RUN/app` and refuses a path that already exists. The exposure is the next kit built from guide § 4, because the preamble is the shared guard and it has just become destructive.

The second problem is newer: the trap now reads the full environment of EVERY same-uid process into a shell variable. That includes the operator's other shells and their tokens, and it happens even when the process already matched. Under `bash -x` that variable is echoed into the run's output.

The answers to the brief:
- **A HOME that merely contains `KIT_HOME` is not matched.** The match is `grep -x`, an exact line.
- **This Claude session is not killed.** It is an ancestor, so it is skipped. If it is not an ancestor it still does not match: its cwd is the harness checkout and its HOME is the real one.
- **Environ is never printed** by the code as written.
- **The seed-security hard fixture is never started by the kit.** It binds `127.0.0.1:0`.

## Findings
- **[medium / confidence: high (mechanism), medium (likelihood)]** `.agents/docs/testing/kit-preamble.sh:32`, `:143`, `:153-155` — The roots of a destructive check are not validated.
  - **The weakness:** the real-home guard is a literal string comparison, and `$HOME/.`, `$HOME//`, a symlink to the home, a parent of the home (`/home`), or `/` all pass it. `KIT_CLONE` is only checked for having a `.git` and a parent outside a checkout, so a dotfiles home whose parent is not a checkout passes too. A kit typo also gets through, for example `KIT_HOME=$SCRATCH/home` with `SCRATCH` unset gives `/home`.
  - **Why it bites:** the roots are resolved physically with `pwd -P`, and the matches at lines 154-155 (`"$r"|"$r"/*` on cwd, `"$r"/*` on exe) then SIGTERM every same-uid process under them. The roots are also re-read from the mutable `$KIT_CLONE` and `$KIT_HOME` at exit time, so a kit that reassigns either one after sourcing changes what gets killed.
  - **Remediation:** resolve both roots with `pwd -P` once, when the preamble is sourced, and keep them in `readonly` variables that the trap uses. Refuse with exit 2 before the CLI starts when a resolved root is `/`, equals the physical real `$HOME`, is an ancestor of it, or equals or contains the harness checkout. Pin each refusal in `TestKitPreambleKillsLeftoverProcesses` with `$HOME/.`, a symlink to the home, and a parent of the home.
- **[low / confidence: high]** `.agents/docs/testing/kit-preamble.sh:158` (`penv=$(tr '\0' '\n' < "$p/environ" ...)`) — Every same-uid process's full environment is copied into a shell variable.
  - **The weakness:** the copy happens for every process, including the operator's unrelated shells, even when `hit` is already set.
  - **Why it bites:** no kit sets xtrace today. But under `bash -x run.sh` or a `set -x` added while debugging, bash prints `+ penv='<whole environment>'` for each process, which puts GH_TOKEN-style values into the run's output. Kit output is quoted into committed `docs/testing/` logs. The earlier pipe form never held the data in a variable.
  - **Remediation:** match without buffering, and only when `hit` is still empty: `grep -qzxF -- "HOME=$KIT_HOME" "$p/environ" 2>/dev/null`. `-z` treats NUL as the record separator. grep reads the file itself, so there is no pipe and no 141 from BUG-155. The environ value then never sits in a variable.
- **[low / confidence: medium]** `.agents/docs/testing/kit-preamble.sh:162` — Up to 80 bytes of a leftover process's argv go to stderr unsanitised.
  - **Tokens:** argv can carry a token, for example `--token=…` or `-H 'Authorization: Bearer ya29…'`. Kit stderr is quoted into verification logs.
  - **Forged lines:** the agent under test controls that argv, including `exec -a`, newlines and ANSI escapes. It can forge a `kit: no processes left by the run` line, or erase the FAIL line on a terminal. In a file log the FAIL line and exit 1 still stand.
  - **Remediation:** print the pid, the match kind and `basename` of the exe. If argv is kept, mask it with `kit_cred_text` and strip non-printables first (`LC_ALL=C tr -c '[:print:]' '?'`).
- **[low / confidence: high]** `.agents/docs/testing/seed-security/run.sh:10`, `ab-kit/run.sh:25`, `agy-rel/run.sh:14` — The headers now say a server the agent leaves running "is killed by the preamble's exit check". That holds only when the kit shell's EXIT trap runs.
  - **The gap:** a kit killed with SIGKILL (TaskStop, an outer tool timeout, OOM) leaves an agent-started server listening, and nothing says so any more. For seed-security that server is the deliberately RCE fixture.
  - **Remediation:** keep one line saying that after a killed kit you run `ss -ltnp` and kill by pid.
- **[low / confidence: high]** Restatements left stale, an integrity issue rather than a security one.
  - `agy-rel/run.sh:16-17`: the "a row counts only with…" rule does not list the new `kit: no processes left by the run` line.
  - `BUG-163-…md:31`: Prevention still says "(cwd inside the clone)".
- **[low / confidence: medium]** `internal/harness/catalog_test.go` `TestKitPreambleKillsLeftoverProcesses` — The test is missing some cases.
  - **Near-misses:** it has no negative cases for the exact-match claims, so nothing pins them: `HOME="$KIT_HOME"x`, `HOME="$KIT_HOME/sub"`, and a cwd in a sibling `"$KIT_CLONE"-other` (the separator in `"$r"/*`).
  - **Ancestor skip:** it has no case for skipping an ancestor.
  - **Broad roots:** it has no case for refusing a broad root (finding 1).

## Checked
- **Secrets:** no credential or token in the diff. `remote.md:175` names the broker script path `/etc/profile.d/zz-styrman-credentials.sh`, an internal identifier, not a secret. Environ is never echoed; argv is (findings 2 and 3).
- **Injection:** `$label`, `$effort` and `$mode` in seed-security are restricted to a fixed set or used only in paths behind the existence refusal. The fixture `sed` uses a fixed pattern. The CLI and the checks run under `env -i`. No new eval.
- **Authz / process control:** SIGTERM only, never SIGKILL. The check is same-uid by the real uid in `/proc/<pid>/status`. It skips the trap shell, its ancestors and its descendants, and it fails closed: a match FAILs the run even when the process survives. PID reuse between the read and the kill is a theoretical race and negligible. Matching on the HOME environment is exact, with `-x`.
- **Crypto:** none in harness code. The seed-security hard fixture (`run.sh:133` onward) hand-rolls an HMAC session cookie and fails open on purpose. It is test data inside a heredoc, written to `$SEED_OUT/runs/<label>/app` outside the harness module, and never compiled or started by the kit. It is not a finding against the harness. It also has unscored flaws: an empty key when `SVC_COOKIE_KEY` is unset, and the token sent unescaped in a query string.
- **Fixture exposure:** the skeleton binds `127.0.0.1:0` (`run.sh:40`, `:72`). Only the agent under test (`--permission-mode auto`) could start it. Its `HOME=$KIT_HOME` (from `env -i`) and its cwd in the clone mean the trap kills it, as long as the trap runs (finding 4).
- **Dependencies:** `go.mod` and `go.sum` are unchanged. `govulncheck` is not installed here: it is not on PATH and there is no `~/go/bin`, so it did not run. npm is not applicable.
- **GCP surface:** not touched.

## Not examined
- The untracked files outside the dispatch's scope: `.agents/docs/plans/2026-10-08-kit-leftovers-and-hard-seeds.md`, `.agents/docs/reviews/2026-10-08-review-kit-leftovers-plan-cold-read.md`, `.agents/docs/reviews/2026-10-08-review-worktree-5-code-docs.md`, `.agents/docs/testing/2026-10-08-verify-seed-security-hard.md`.
- The unchanged bodies of `ab-kit/run.sh` and `agy-rel/run.sh`. I read only their diffs and grepped their root and launch lines. Of `kit-preamble.sh`, lines 1-130 and 133-230 were read.
- `guide.md` and `remote.md`: only the changed paragraph or line.
- Any edit made after my last read, given that the tree moved during the review.
- Not run, on instruction: any kit, and `TestKitPreambleKillsLeftoverProcesses`. `govulncheck` is unavailable.

## Round 2 (2026-10-08)

Range reviewed: HEAD (280c1df)..working tree at round 2 (`git diff HEAD --stat`: 9 files, 324 insertions, 66 deletions). I verified the fixes only. Read: `kit-preamble.sh:28-55` (the root guards), `:146-184` (`kit_process_leftovers`), the three kit headers, agy-rel:14-18, BUG-163 Prevention, the INDEX row, and `catalog_test.go:1930-1995`. Ran `go test ./internal/harness -run 'TestKitPreamble|TestScriptsHaveNoEarlyExitPipe' -count=1`: all PASS, including `TestKitPreambleKillsLeftoverProcesses` (16.3 s). No kit was run.

### Fix verification
- **Round-1 medium (roots not validated): FIXED.**
  - **What the code does now:** `KIT_CLONE` must be absolute. Both roots are resolved once with `pwd -P` when the preamble is sourced and made `readonly`. `kit_process_leftovers` uses only those two roots (line 157). The trap is installed after the guards, so a refused load kills nothing.
  - **Traced by hand:**
    - `case "$kit_real_home/" in "$kit_r"/*)` matches both equality (`*` can be empty) and a parent. `$kit_r` is quoted, so it is literal.
    - `/` misses both case patterns (`//*`) but the explicit `[ "$kit_r" != / ]` catches it.
    - `$HOME/.` and a symlink to the home resolve to the physical home and are refused.
    - The harness checkout, or any parent of it, is refused.
  - **`mkdir -p "$KIT_HOME"` running before the refusal is not a defect.** Every refused path is `/`, the home, the checkout, or a parent of one of those, so it already exists and mkdir creates nothing.
  - **Mutation B as described is consistent with the code:** the "/." fixture's home is neither `/` nor a parent of the checkout, so only the real-home refusal catches it.
- **Round-1 low (environ in a variable): FIXED.** Line 175 is `grep -qzxF -- "HOME=$KIT_HOME" "$p/environ"`. It runs only when there is no cwd or exe hit. No variable holds the environment and nothing is piped. The exact NUL-delimited record match is pinned by the `HOME="$KIT_HOME-x"` near-miss, which stays alive.
- **Round-1 low (argv printed): FIXED.** The FAIL line now carries the pid, the match kind and `${exe##*/}` passed through `LC_ALL=C tr -cd '[:alnum:]._+-' | cut -c1-40`. That cannot print argv, a newline or an escape. The pipeline has no early-exit reader, so pipefail is safe.
- **Round-1 lows (docs):**
  - **FIXED:** the three headers now say a SIGKILLed kit never runs the trap and to check `ss -ltnp`.
  - **FIXED:** the agy-rel count rule lists all three lines.
  - **FIXED:** BUG-163 Prevention and the INDEX row describe cwd, exe and HOME.
- **Round-1 low (test gaps): PARTLY FIXED.**
  - **Added:** the HOME near-miss, the descendant skip, and refusals for "/.", `/` and a relative clone.
  - **Still open:** finding R2-1 below.

### Round 2 findings
- **R2-1 [low / confidence: high]** `internal/harness/catalog_test.go:1960-1975` — Some refusals have no test, and the assertion is loose.
  - **Untested refusals:** a parent of the real home (for example `KIT_HOME=<tmp>`, the fixture home's parent), a symlink to the home, and the harness-checkout refusal. A mutation of `"$kit_r"/*` to `"$kit_r"/` (equality only) would drop parent detection and still pass every case.
  - **The loose assertion:** the check is `"refusing" || "absolute"`, so any refusal satisfies any case.
  - **Remediation:** add a parent-of-home case and a symlink-to-home case, and assert each case's own message.
- **R2-2 [low / confidence: medium], residual by design, recorded rather than fixed:**
  - **Child directories:** a root that is a child of the home but holds operator work passes the guard. Examples: `KIT_CLONE=$HOME/other-project` with a `.git` of its own, or `KIT_HOME=$HOME/Documents`. Its processes would get SIGTERM. The guard cannot tell a throwaway clone from a real project. The kits keep roots under `$SEED_OUT`, `$OUT` and `$BASE` and refuse a path that already exists, which is the real protection.
  - **Unreadable environ:** a same-uid process whose `/proc/<pid>/environ` cannot be read (non-dumpable: ssh-agent, gpg-agent) is passed over on the HOME channel, which fails open. The alternative is flagging the operator's own agents on every run.
  - **HOME match at exit:** the HOME match still reads the mutable `$KIT_HOME` at exit, unlike the readonly roots. A kit that reassigns it after sourcing only shifts an exact-string match. Low.

### Not examined (round 2)
- The unchanged remainder of the files, beyond what round 1 read.
- `guide.md` round-2 wording.
- The ab-kit run in flight: not touched, and no kit was run.
- `govulncheck`: still not installed.

## Outcome

Written by the main agent, 2026-10-08. Round 2 resumed the same agent with `SendMessage` (required: round 1 had a medium). No medium is open.

**Fixed:**
1. **(medium)** The kill roots are resolved once, read-only, and refused at `/`, the real home or a parent of it, and the checkout or a parent of it. Tested; mutation B (the real-home refusal removed) fails the `/.` case.
2. **(low)** `environ` is matched in the file with `grep -qzxF`.
3. **(low)** The FAIL line prints pid, kind and a sanitised executable name.

**Fixed after the cap** (not re-reviewed):
- **R2-1:** each refusal now has its own test, asserting its own message: a parent of the home, a link to the home, the checkout, a parent of the checkout.
- **R2-2 (partly):** the HOME match reads a read-only `kit_home_env` frozen at load.

**Left, by design:**
- **R2-2:** a root that is a child of the home but holds real work (`$HOME/other-project`) passes. The kits' own roots are under `~/sandbox`, and the guide says nobody works inside a kit's clone or home while it runs.
- **R2-2:** a non-dumpable process's environment cannot be read, so it is not matched by HOME; its cwd or executable still is.
