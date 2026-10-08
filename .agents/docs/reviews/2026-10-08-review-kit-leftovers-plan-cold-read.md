---
owner: reviewer
date: 2026-10-08
---
# Review — kit-leftovers plan, cold read   (2026-10-08)

Range reviewed: `.agents/docs/plans/2026-10-08-kit-leftovers-and-hard-seeds.md` (untracked), with `.agents/docs/testing/kit-preamble.sh`, the three kits' `run.sh`, `internal/harness/catalog_test.go:1543-1655`, `testing/guide.md:158-182`, `deploy/remote.md:175`, BUG-160/163, `reviews/2026-10-07-security-v2.5.0-release.md:34`.

## Status: COMPLETE

## Verdict: REQUEST CHANGES (design needs revision; implementation not ready)

The worst problem is part B. `seed-security/run.sh:139` dispatches `security` over `git diff HEAD~1..HEAD`, which is the last commit only. A "second commit of three harder flaws" therefore takes the original six out of the diff the agent sees. The plan's D2 ("6/6 history stays comparable") cannot hold unless the dispatch, the budget or the number of dispatches changes, and the plan decides none of that. The second problem is part A. The process match is specified without the self-match rule (a command-substitution subshell of the exit trap, with its cwd in the clone, kills itself or fails the run), the grace period for a CLI that is still exiting, or the order against the credential scan. Part C is sound: the security review at line 34 names exactly this.

## Findings

- **[high / confidence: high] `seed-security/run.sh:139`, plan B and D2** — Dispatch scope. Today the agent sees one commit holding six flaws. With a second commit, either it gets `HEAD~2..HEAD` (nine flaws, one 20-tool-call budget, so the `found N/6` numbers are no longer comparable with the 6/6 history), or it gets two dispatches per run (twice the cost, and "3 runs at high and 1 at max" doubles). The commit order also decides which six are visible. Decide it in the plan. Recommendation: keep the six-flaw commit and dispatch exactly as before, so `found N/6` is untouched. Add a separate `hard` mode, selected by an argument, that builds the six-flaw commit plus the hard commit and dispatches over the hard commit only (`HEAD~1..HEAD`). Scored `hard M/3`, it is then comparable only with itself, and D2 is true by construction.

- **[high / confidence: high] plan A "What matches" / Failure modes, against `kit-preamble.sh:132-154`** — The scan must not match its own subtree. `seed-security/run.sh:20` does `cd "$KIT_CLONE"` in the kit's main shell, so every process the exit trap spawns inherits a cwd inside the clone. Examples are a `$(...)` subshell (where `$$` is still the parent's pid, so a `$$` skip does not cover it), `readlink`, `tr` and `cut`. The preamble's own idiom is `left=$(kit_credential_leftovers)` (line 143); copying it for the new check makes the scanner a descendant of `$$`, not an ancestor, and the plan skips only "the kit's own shell and its ancestors". The old copies avoided this only because they ran in the main shell. Decide: skip `$BASHPID`, `$$` and every process whose parent chain reaches `$BASHPID`, or run the scan without a subshell and read each `/proc/<pid>/stat` ppid for the walk. Also decide that "same uid" is `[ -O "$p" ]`, not `stat -c` (GNU only; the plan promises a Mac path).

- **[medium / confidence: high] plan A, nothing stated about timing** — A CLI that has just exited (the `claude` node children, an `agy` language server) can still be tearing down when the trap runs. With the new HOME and exe matches, that is a false FAIL on a normal run, and the plan has no grace. The only evidence the cwd-only version was quiet is two runs. Decide: poll up to about 2 s for matches to disappear before SIGTERM, then SIGTERM, then a short wait and SIGKILL (a process that traps TERM stays alive, and the plan only says "sent SIGTERM"). Also decide the order against the credential scan: kill and wait FIRST, then run `kit_credential_leftovers` and the fingerprint. A process still alive can write transcripts or tokens into KIT_HOME after the scan has passed. In the old kits `kit_leftovers` ran before the trap, so moving it must not reverse that order.

- **[medium / confidence: high] plan A "The three kits drop their copies", `agy-rel/run.sh:142-146`** — The verdict line goes stale. In agy-rel the leftover result is folded into `verdict`, and the kit prints `row N: PASS|FAIL…` and exits on `${verdict%% *}`. After the move, the kit prints "row … PASS", then the trap prints the FAIL and exits 1. The printed verdict contradicts the exit status, and the operator or any grep reads PASS. The same applies to ab-kit's "finished (claude exit rc)" line and to seed-security's `found N/6` line, which print before the trap. Decide: the guide rule becomes "a run counts only at exit 0 with the three lines; a kit's own verdict line is provisional", and state it in guide § Behavioural runs and in the preamble header (lines 24-27). Recommend also that agy-rel keep its verdict honest by reading the exit status in its caller, or by printing "provisional" in the line.

- **[medium / confidence: high] plan Verification** — The tests as written cannot fail cleanly.
  1. `CombinedOutput` waits for pipe EOF, and a planted `sleep` holds the pipe open. In the mutation run (check disabled) and in any failing run, the test hangs until the `sleep` ends. Decide: the planted `sleep` is started with `>/dev/null 2>&1 </dev/null`, has a short duration, and a `t.Cleanup` kills it by pid (read from a pid file).
  2. The test must prove the condition holds before the verdict: read `/proc/<pid>/cwd`, and for the HOME case read `/proc/<pid>/environ`, and `t.Skip` with a reason when `/proc` is missing (macOS: the plan's own Mac path would otherwise FAIL the test) or `environ` is unreadable (hidepid, containers). The plan's "Checked by the test" is not a design.
  3. It must assert the process is dead afterwards (poll `kill -0`, because SIGTERM is async), not only that the FAIL line printed.
  4. The plan lists only the cwd and HOME cases. The exe case (copy `sleep` into the clone, run from elsewhere) is specified but untested. There is no negative control either: a same-uid `sleep` with cwd outside both trees, no match in environ, must survive and the run must exit 0. The false-positive guard is the plan's stated protection, and nothing tests it.
  5. Symlinked KIT_HOME: the existing test at `catalog_test.go:1626-1635` exercises it. `/proc/<pid>/cwd` and exe are physical, but `HOME=` in `environ` is the literal string. The match needs both the physical and the literal form of KIT_HOME and of KIT_CLONE (the old copies only did `pwd -P` for the clone). A trailing slash on KIT_HOME is allowed by `kit-preamble.sh:31`. Decide how `environ` is compared: split on NUL and compare whole entries `HOME=<value>` exactly, so `…/kit-home2` does not match `…/kit-home`.
  6. The mutation must follow `review-standards` § 2: count the target text first, restore from a saved copy, read `git diff` after, and show the shell still parses (it is a `.sh` file, so "compiled" means `bash -n`).

- **[medium / confidence: medium] plan A, `kit_on_exit` runs under `set -euo pipefail` (`kit-preamble.sh:28`)** — A vanished pid inside the trap (`tr … < /proc/$pid/cmdline`; the `<` redirect error is the shell's, `2>/dev/null` on `tr` does not cover it) returns non-zero in an assignment or a pipeline. That aborts the trap function in the middle and skips the remaining exit lines, and the run is then "FAIL" for the wrong reason or prints neither the credential line nor the process line. Decide: every `/proc` read in the scan is `|| true` and the scan runs in a block that cannot propagate. Also say what happens when the scan itself errors (fail closed: print a FAIL naming the error, per `review-standards` § 8 first bullet). The plan only covers "no /proc".

- **[medium / confidence: medium] plan Interfaces: "A run counts only with this line plus the two existing ones"** — On a machine with no `/proc` the new line is `kit: no /proc here: …`, not `kit: no processes left by the run`, so by the plan's own rule the Mac run never counts, yet "does not fail the run". Decide: either the Mac run counts with the no-/proc line (then say so in guide § Behavioural runs and mark such runs as not covering BUG-163), or it does not count. Recommend the first, with the claim in the guide naming the gap. The Mac path (`[ -d /proc/self ]`) is not tested anywhere on Linux; allow an override (`KIT_PROC_ROOT`, default `/proc`) so the test can point at a missing directory.

- **[medium / confidence: medium] plan D1 and Failure modes: "HOME=$KIT_HOME … which no operator process carries"** — The claim is not checked. The operator debugging a run with `HOME=$KIT_HOME bash` or `claude` carries it, and an inspection shell `cd`'d into the clone is matched by cwd. Both are killed with SIGTERM. The same-uid guard does not help. Also unchecked: that neither CLI keeps a legitimate helper under `HOME=$KIT_HOME` that outlives the CLI's exit (an Antigravity language server or daemon). The old check could not hit that case; the new HOME match can, and both CLIs are only exercised once each in the E2E. Decide: the E2E reads the printed pids and commands of any match on a clean run (expected: none) for both `claude` and `agy`; document in the guide that an operator shell with HOME or cwd in a live run is killed.

- **[low / confidence: high] stale restatements the plan's Surfaces list misses**
  - `seed-security/run.sh:9-10` says "A server the agent itself starts survives a timeout: check `ss -ltnp`". It is the kit being edited, and it is stale once the check is in the preamble.
  - `BUG-163`'s `## Prevention` ("cwd inside the clone") and the `Related` text; the plan lists BUG-163 for its limit paragraph only.
  - `BUG-160` line 20 ("A run counts only with 'kit: no credentials left in KIT_HOME'"), and `guide.md:161-162`, which restates that counting rule. Both gain a third line.
  - `guide.md:182`, the preamble sentence in § 4 ("lists the real … and FAILS the kit when it changed … and puts a wrapper first"), which enumerates the preamble's effects and is the sentence the plan means, but the plan names "the sentence describing the preamble's exit lines" without a line number; and `TestKitPreambleAssertsCommit`'s guide mention.
  - The preamble header lines 18-27, which the plan does name.
  Run `scripts/restatement-grep.sh` over "BUG-163", "outlived", "left by the run", "run counts only" before the first review round.

- **[low / confidence: medium] plan B scorer** — Why item 2 criticises "scored by keyword", and the hard scorer is again a keyword grep (`HasPrefix`, `fail-open`). A keyword for flaw 2 matches an artifact that merely quotes the line. Decide that the hard pass/fail comes from the by-hand read of each artifact, with the keywords as a hint, and state it in the kit's header. Also say which budget the dispatch gets: a prompt of "20 tool calls" for nine flaws.

- **[low / confidence: medium] plan Phases** — The E2E row requires one `ab-kit` and one `agy-rel` run, each printing all three lines. Neither can show the kill path (the test covers it); say so. Nothing in the plan verifies the kits after the copies are removed, other than their E2E. `bash -n` on the three kits is cheap. The kits each use `$$`-keyed locals (`lo`, `kit_leftovers || lo=1`) that need removing together: `ab-kit/run.sh:82,94`, `seed-security/run.sh:145,152` (`[ "$lo" = 0 ]` is the seed kit's exit status; replace, do not leave `lo` unset under `set -u`).

## (2) Scope

Not named in `## Why`, but I consider them justified: the HOME and exe match (the Why names the cwd-only limit and R2-1's server-from-/tmp case). Not justified by the Why: SIGKILL escalation and grace polling are my additions, not the plan's. The plan itself adds nothing beyond its Why. "Out" leaves the `GIT_*` clearing alone, correctly.

## Criteria
1. Dead code — pass: the three `kit_leftovers` copies are removed; the `lo` variables and the stale seed header sentence must go with them (see the last finding).
2. Flow tracing — fail: the exit-trap ordering against the credential scan, the stale verdict lines in agy-rel/ab-kit/seed, `set -e` in the trap, and a scan that can match its own subtree.
3. Standards conformance — fail: restatements in BUG-160, BUG-163, guide.md:161-162 and 182, and the seed header are not on the plan's surface list; the tests do not prove their precondition (`review-standards` § 2); the mutation protocol is not stated.
4. Directory layout — pass: the surfaces are where the project puts them.
5. Simplicity — pass: bash and `/proc`, no dependency. The grace and SIGKILL escalation are the only additions, and they are needed for correctness.
6. Decoupling — pass: the check lives in one place, which is the point of A.

## Questions (answers recommended)
1. BLOCKING — dispatch scope for the hard flaws: separate `hard` mode, dispatch over the hard commit only (see finding 1).
2. BLOCKING — how the scan avoids matching its own subtree: skip `$BASHPID`, `$$` and descendants of `$BASHPID`; no `$( )` around the scan.
3. BLOCKING — order of kill, wait, then credential scan; grace of about 2 s, SIGTERM then SIGKILL, and a final "still alive" FAIL naming the pid.
4. BLOCKING — what the Mac run prints and whether it counts: the no-/proc line counts, guide names the gap.
5. BLOCKING — fate of the kit-printed verdict lines: provisional; the counting rule is exit 0 plus three lines.
6. NON-BLOCKING — compare `environ` by whole NUL-separated entry; both physical and literal KIT_HOME and KIT_CLONE.
7. NON-BLOCKING — the planted-process tests: redirected output, pid-file cleanup, precondition read from `/proc`, skip when `/proc` or `environ` is unavailable, a negative control and the exe case.
8. NON-BLOCKING — the scan runs with `|| true` on every `/proc` read and fails closed on its own error.
9. NON-BLOCKING — add the restatement files above to Surfaces and run the restatement grep first.
10. NON-BLOCKING — hard scorer is a hint; hand read decides.

design needs revision
implementation not ready

## Not examined
`internal/harness/size_test.go` and the context size pins (the plan changes no `context/` file); `ab-kit`/`agy-rel` lines outside the leftover blocks and their tail; `skills/` and `commands/` (not touched by the plan); whether Antigravity's CLI leaves a daemon with `HOME=$KIT_HOME` after exit (unknown, no run made; no CLI was started in this review); the contents of the seeded Go fixture beyond the commit structure.
