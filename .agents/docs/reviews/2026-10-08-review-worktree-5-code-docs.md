---
owner: reviewer
date: 2026-10-08
---
# Review — worktree-5-code-docs   (2026-10-08)

Range reviewed: working tree vs HEAD (280c1df), `git diff HEAD --` over: `.agents/docs/testing/kit-preamble.sh`, `internal/harness/catalog_test.go`, `.agents/docs/testing/{ab-kit,agy-rel,seed-security}/run.sh`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/remote.md`, `.agents/docs/bugs/BUG-163-kit-run-left-a-server-listening.md`. The plan was not read (intent-after-own-read was optional and the diff sufficed).

## Status: COMPLETE

## Verdict: APPROVE WITH NITS (round 1; final verdict in ## Round 2)

Nothing blocking and no `REQUIREMENTS.md` violation found. The worst problem is a stale restatement: `agy-rel/run.sh:17` still tells the reader a row counts with two preamble lines, while the preamble and guide now require three (the new "no processes left" line), and `bugs/INDEX.md:115` still describes the removed `kit_leftovers` and "cwd in the clone" only. The trap itself holds up under `set -euo pipefail` (no unguarded abort found before the credential check), the hard fixture compiles, and `TestKitPreamble*` passes. The BUG-163 "Limit closed" paragraph claims more than the code covers and does not name what is still open.

## Findings

- **[medium / confidence: high]** `.agents/docs/testing/agy-rel/run.sh:17` — stale restatement of a changed rule. It says a row counts "only with PASS here plus the preamble's 'real config unchanged' and 'no credentials left in KIT_HOME' lines". The preamble header (`kit-preamble.sh:26-28`) and `guide.md:182` now say all three lines, including `kit: no processes left by the run` (or the no-/proc line). An operator following the agy-rel header would accept a row with no process-check line. Fix: add the third line there. Also stale: `.agents/docs/bugs/INDEX.md:115` (names `kit_leftovers`, which no longer exists, and "cwd in the clone" only; the INDEX row is the Graduated/prevention text a later reader sees first), `BUG-163…md` § Prevention ("cwd inside the clone"), and, weaker, `guide.md:161-162` and `ab-kit/run.sh:15`, which still present the credentials line as the line a run counts on.

- **[low / confidence: high]** `BUG-163-kit-run-left-a-server-listening.md:28` ("Limit closed") — a coverage claim without the gap named (review-standards § 7). The old text named a limit; the new text says "closed" but the code still misses: (a) a process that changes directory AND resets HOME (`env -i`, `HOME=/tmp`) with an executable outside the roots; (b) anything the agent starts through the reachable Docker socket (a container's process has neither the cwd, the exe nor the env); (c) a process that ignores SIGTERM, which is reported and fails the run but is not verified dead after the 1 s `sleep` (the 1 s is a pause, not a wait-until-gone); (d) the orphan-listener check the old limit paragraph proposed. Also the descendant skip means that when the kit itself is killed (trap runs while the CLI is still its child) that tree is never flagged or killed. Fix: name these in the paragraph and in guide.md:182 in one sentence.

- **[low / confidence: medium]** `kit-preamble.sh:164-166` (`kit_on_exit` `cd /` then `kit_process_leftovers`) — only `KIT_HOME` is required to be absolute (line 33); `KIT_CLONE` is not. With a relative `KIT_CLONE`, after `cd /` the root is dropped silently (`[ -d "$r" ]` false), no cwd/exe match is possible, and the trap still prints `kit: no processes left by the run`. All three kits pass absolute paths today, so no current failure; the line then reads as a pass for a check that did not look at the clone. Fix: require `/*` for `KIT_CLONE` beside the `KIT_HOME` check.

- **[low / confidence: medium]** `kit-preamble.sh:134-166` — false-positive surface, by design but unstated: any same-uid process whose cwd or exe is under the clone or `KIT_HOME` and is not an ancestor/descendant of the kit shell gets SIGTERM and fails the run: an operator's second terminal or `tail -f`/`less`/editor/language server open in the clone, a `git fsmonitor--daemon`/`gc --auto` the agent's git started (arguably a true leftover). An interactive bash ignores SIGTERM, so the cost is a failed run and a killed viewer, with the pid named in the message. Separately and unverified (I did not run a kit): if the Antigravity CLI leaves a daemon with `HOME=$KIT_HOME` alive by design, every agy-rel row would fail on the new HOME match; the first agy-rel run after this lands should be read for that. Fix: say in the preamble header that the operator keeps no shell or watcher in a kit's clone or home while it runs.

- **[low / confidence: medium]** `internal/harness/catalog_test.go:1853-1944` (`TestKitPreambleKillsLeftoverProcesses`) — gaps against its own claims. (1) The ancestor/descendant skip is not exercised: nothing in the test makes the skip the thing that prevents a match (the trap's subshells run with cwd `/` and no `HOME=$KIT_HOME`, so removing the skip would likely still pass); the BUG record says "mutations seen failing" but this test has no case for it. (2) `alive()` is `os.Stat("/proc/<pid>")`, which also succeeds for a zombie: on a box whose PID 1 never reaps (a container without an init) the killed orphan stays a zombie and the three kill cases fail although the check worked. Fix: read `/proc/<pid>/stat` and treat state `Z` as dead; add a case with a child of the shell that has its cwd in the clone and must not be killed (or state the skip as untested).
  What the cases do prove (verified by reading and by running `go test ./internal/harness -run TestKitPreamble -count=1 -v`, all PASS, 6.8 s for this test): each planted process is an orphan (inner subshell exits before the trap runs), the pid comes from `$!` of the backgrounded list that then `exec`s, so it is the process itself, output is on `/dev/null` so no pipe is held (cannot hang `CombinedOutput`), `t.Cleanup` kills leftovers, and the negative control is asserted alive. The three "wrong-detected" cases are one per match kind (cwd, exe, HOME). Not mutation-tested by me.

## Checks that came back clean (so the next reader need not redo them)

- Trap under `set -euo pipefail`: `kit_process_leftovers` is called as `… || failed=1`, so `set -e` is off for its whole body; `kit_ppid` and every `/proc` read end in `|| true`; `kill` is `|| true`; `$(…)` pipelines are wrapped in `{ …; } || true` and `tr | grep -q` under pipefail is guarded. `set -u`: `"${roots[@]}"` is safe on bash 5.2 (empty array only on bash < 4.4, which is macOS bash 3.2 where `/proc` is absent and the function returns first). No unguarded command found that aborts the trap before `kit_credential_leftovers`. Exit status: `exit "$rc"` is unchanged, and the three kits correctly lost their own `lo`/`verdict` plumbing (ab-kit ends `exit "$rc"`, seed-security ends with `echo`, agy-rel with the verdict test; the trap turns a leftover into exit 1 in each).
- Orphan reparented to a subreaper that is itself an ancestor of the kit shell (systemd --user, tmux): not skipped, since the skip set holds the ancestors' pids but the orphan is a different pid whose own chain never reaches `$$`.
- Prefix confusion: roots compare with `"$r"|"$r"/*`, so `clone2` does not match `clone`.
- seed-security hard mode: I extracted `main.go` + `hard.go` from the heredocs into `.agents/scratch/rev-hard/`, applied the `sed`, and ran `gofmt -l`, `go vet ./...` and `go build ./cmd/svc`: it compiles, vets and is gofmt-clean; the `grep -q 'hardRoutes(mux'` guard matches once. `HEAD~1..HEAD` is exactly the "sessions, raw file reads and token rotation" commit over the skeleton in hard mode (skeleton is `HEAD~1`); in six mode it stays the six-flaw commit. The three flaws are real in the code (`strings.Cut` found=false skips the HMAC; `HasPrefix` without a separator; `err != nil` admits). Scoring keywords are loose (`forge`, `unsigned`), which the header already says is read by hand. The fixture's `go test`-free, loopback-only design is unchanged. (Scratch directory `.agents/scratch/rev-hard/` is gitignored.)
- remote.md new sentence: `/etc/profile.d/zz-styrman-credentials.sh` exists on this box (outside the repo, as the guide already describes); `BUG-160-*.md` has `## Addendum` headings; the 2026-10-07 security reviews exist under `reviews/`. Names are real.
- guide.md § Behavioural runs vs the preamble header: both list the same three exit lines and the no-/proc variant (`kit: no /proc here: …`, which the log names as a gap); the test names cited in guide.md:182 exist.
- `restatement-grep.sh` for "no credentials left" listed `plans/2026-10-07-lean-process-2.md:41` (historical plan, leave), `ab-kit:15`, `agy-rel:17` (finding above), `guide.md:162`, `kit-preamble.sh:28` (updated), and the two tests. For "kit_leftovers": only `bugs/INDEX.md:115` (finding above). For "survives a timeout" and "outlived the run": nothing left standing.

## Criteria
1. Dead code — pass: the three pasted `kit_leftovers` copies and their `lo`/verdict plumbing are removed; no orphan helper left (`kit_ppid` is used). One stale name survives in `bugs/INDEX.md:115`.
2. Flow tracing — pass with notes: kits to preamble trap to exit status traced for all three; trap is robust under `set -euo pipefail`; the test's pid/orphan/negative-control claims hold, but the skip logic and zombie handling are untested (finding 5).
3. Standards conformance — fail (minor): stale restatements in agy-rel header, INDEX.md, BUG-163 text and a coverage claim without its gap (findings 1, 2).
4. Directory layout — pass: the test sits beside the other `TestKitPreamble*` tests in `catalog_test.go`; shared code moved into the shared preamble.
5. Simplicity (DIY-except-security) — pass: bash and `/proc` only, no dependency, no auth or crypto code written (the fixture's HMAC is seeded-flaw test data, not harness code).
6. Decoupling — pass and an improvement: one check in one file instead of three drifting copies.

## Not examined
Not run: any kit (as instructed), so the Antigravity-daemon and real-claude behaviours are unverified; no mutation of the preamble or test was run by me (the BUG record's "mutations seen failing" is unverified); the plan file `2026-10-08-kit-leftovers-and-hard-seeds.md` and `2026-10-08-verify-seed-security-hard.md` were not read; no check of macOS bash 3.2 behaviour; `go vet`/full `go test ./...` not run, only `TestKitPreamble*`.

## Round 2

Verified by reading the working tree and running `go test ./internal/harness -run 'TestKitPreamble|TestScriptsHaveNoEarlyExitPipe' -count=1` (ok, 40 s). No kit run; mutations not re-run by me.

Final verdict: APPROVE WITH NITS. All round 1 findings are fixed; no defect introduced by the fixes is blocking.

1. Stale restatements (was medium): fixed. `agy-rel/run.sh:17` names all three lines; `ab-kit/run.sh` header points at guide § Behavioural runs' three lines and now adds that a SIGKILLed kit never runs the trap (true; also in agy-rel); `bugs/INDEX.md` BUG-163 row and BUG-163 Prevention name cwd/exe/HOME and the preamble; `guide.md:162` marks the credentials line as one of three. `kit_leftovers` survives only as a historical fix bullet in BUG-163. Complete for the restatement-grep set.
2. "Limit closed" overclaim (was low): fixed. The "Still not caught" list names the env -i + cd orphan, Docker-socket containers, SIGTERM-ignoring process, the kit shell's descendants, and an operator's terminal in the clone/KIT_HOME. Matches the code.
3. Relative `KIT_CLONE` (was low): fixed in `kit-preamble.sh:33-45`: absolute check (exit 2), roots resolved once with `pwd -P` and made `readonly`, refused when `/`, the real home or a parent, or the checkout or a parent. Traced: `"$kit_real_home/"` against `"$kit_r"/*` refuses the home itself and its parents (a `$HOME/.` KIT_HOME passes the literal `"$HOME"` check and is caught here); `/` has its own test because `"//*"` would not match. The kits' legitimate roots (under `$OUT`, outside the checkout) are not refused; the full preamble tests pass.
4. Test gaps (was low): `alive()` now reads `/proc/<pid>/status` and treats `State:\tZ` as dead (zombie flake fixed). The descendant-skip case is a real test: without the skip the child with cwd in the clone would be killed and the case would expect exit 0 and fail. It is killed before the next case. Added near-miss `HOME=$KIT_HOME-x` (exact-entry match pinned) and refusal cases for `KIT_HOME=$HOME/.`, `KIT_HOME=/`, relative `KIT_CLONE`.
5. `grep -qzxF -- "HOME=$KIT_HOME" "$p/environ"`: correct (NUL-separated entries, whole-entry match, file read directly, no pipe, a failed read is a non-match under `if`). The FAIL line prints pid, kind and a sanitised exe basename only (`tr -cd`), which removes the argv-token leak and the escape-injection risk; the test still matches `"  <pid> ("`.

Remaining nits (low, none blocking):
- **[low / confidence: high]** The root refusals for "a parent of the real home" and "the harness checkout or a parent" have no test: only `$HOME/.` and `/` are covered, so the checkout branch could be deleted without a test failing (the coordinator's "mutation seen failing" is for the home case only as far as the test file shows). Fix: add two cases (`KIT_HOME` = parent of `home`, `KIT_HOME` = the repo root).
- **[low / confidence: medium]** `guide.md:182`, the list of what the preamble asserts and refuses, does not mention the new refusals (absolute `KIT_CLONE`, roots that are `/`, the real home or a parent, or the checkout or a parent); only the preamble's own comment and messages do. Fix: one clause in guide.md:182.
- **[low / confidence: low]** `mkdir -p "$KIT_HOME"` runs before the root refusals, so a refused root such as a path under another directory still gets created; harmless for `/` and an existing home, noted only for order.

## Not examined (round 2)
Real kit behaviour (one ab-kit run was in flight and not touched); the security review `2026-10-08-security-worktree-kit-leftovers.md` was not read, only the fixes it caused; mutation runs.

## Outcome

Written by the main agent, 2026-10-08. Round 2 resumed the same agent with `SendMessage` and approved with nits.

**Fixed in round 1:**
- the stale restatements: agy-rel :17, the ab-kit header, the INDEX row, BUG-163's Prevention, guide :162
- BUG-163's "Still not caught" list
- `KIT_CLONE` must be absolute
- the liveness probe is zombie-safe
- a descendant-skip case

**Fixed after the cap** (not re-reviewed; the main agent read the diff):
1. Refusal tests for a parent of the home, a link to the home, the checkout and a parent of the checkout, each asserting its own message.
2. `guide.md` names the refusals.
3. `KIT_HOME` is resolved without being created, and `mkdir` runs after the checks.

**Mutation evidence for the checkout refusal:** the run with that guard removed proceeded past the refusal, wrote its wrapper to `<checkout>/bin/sh-agent` and began a credential scan of the whole checkout, instead of exiting 2. That run was stopped, the wrapper deleted by name, and the preamble restored (`diff -q` identical). The mutation is not re-run: disabling that guard acts on the real checkout.
