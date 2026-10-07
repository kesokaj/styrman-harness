---
owner: reviewer
date: 2026-10-07
---
# Review — worktree-2-code-docs   (2026-10-07)

Range reviewed: working tree (git diff HEAD over agents/security.md, context/AGENTS.md, skills/{bug,feature,preflight,review}/SKILL.md, .agents/docs/testing/guide.md) plus the untracked .agents/docs/testing/ab-kit/{run.sh,accept.py,metrics.py,prompt.md,prefix-full.md}, each read whole

## Status: COMPLETE

## Verdict: APPROVE (after round 2; was APPROVE WITH NITS)

Nothing blocks. The rule edits say one thing each and agree with each other (context DoD 1, feature Step 6, bug Step 5, preflight Step 1.1, review Step 5.5, security effort); `TestRenderedContext*` passes; run.sh honours the preamble's contract (no EXIT trap, existing dirs refused, link removed before the exit check; verified by reading, `bash -n`, shellcheck, and by running metrics.py on the two existing runs). The worst items are (1) preflight Step 2 still says "dispatch in parallel" for a verifier while the plan/D4 says Step 2 changes and the file is unchanged, (2) guide § 4 never states the broker off-switch that run.sh and the preamble's exit check need on a styrman box, and (3) accept.py runs agent-written code with the operator's full environment.

## Findings
- **[medium / confidence: medium]** `skills/preflight/SKILL.md:57,61` — plan `2026-10-07-lean-process-2` Surfaces and D4 name "preflight Step 2 (release gate: verifier after review and security)", ADR-054a lists only Step 1.1, and the diff leaves Step 2 untouched. Step 2 still says "Otherwise dispatch in parallel" and item 1 dispatches the verifier "over the release range" when no log exists, i.e. alongside reviewer and security, the very order the diff reverses elsewhere. A review/security fix then invalidates the log (Step 5 line 99 re-runs it), reproducing the second-verifier cost. Fix: either add one clause to Step 2 ("the verifier is dispatched after the reviewer's and security's fixes are committed") or strike Step 2 from the plan's Surfaces/D4 and say why Step 5's re-run rule is enough.
- **[low / confidence: medium]** `.agents/docs/testing/guide.md:168-170,181` vs `ab-kit/run.sh:16-19` — the guide's § 4 prescribes `HOME=<sb>` for the whole run and "a copy" of the credential, and line 181 now admits a symlink, but the broker off-switch (`env -u STYRMAN_CRED_BROKER STYRMAN_SKIP_CRED_PULL=1`) appears only in the kit header. On a styrman box any other kit written from the guide will have the login-shell broker (/etc/profile.d/zz-styrman-credentials.sh; it returns early on `STYRMAN_SKIP_CRED_PULL`, line 17, and does nothing when `STYRMAN_CRED_BROKER` is empty, line 30; both verified) write an SSH key, ADC and token into the sandbox home, and the preamble then FAILs the run. Add one sentence to § 4 (and say the symlink's deletion is by name, same as the copy). The "copy ... chmod 600 ... deleted by name" text at guide lines ~158-162 should also say the A/B kit's is a link.
- **[low / confidence: medium]** `.agents/docs/testing/ab-kit/accept.py:11` — the binary the agent wrote is built and run with `env={**os.environ, ...}` and the operator's real HOME, i.e. untrusted generated code gets every token in the operator's environment (`GH_TOKEN`, the broker's exports). Pass a minimal env (`PATH`, `PORT`, `DATA_FILE`, a scratch `HOME`). Same for `go build` in the clone (a generated `go.mod` could name a toolchain).
- **[low / confidence: medium]** `ab-kit/accept.py:34,3-9` — the malformed-JSON check shells out to `curl` and sends form content-type (a handler that demands `application/json` returns 415 and fails an arm for a reason the prompt does not state); missing `curl` raises, and an app that never listens makes `req` raise `URLError` (only `HTTPError` is caught), so a dead arm prints a traceback with no `SCORE`. Use `req`-style `urllib` with a raw body and catch `URLError` as a failed check, so the acceptance output always ends in `SCORE n/25`.
- **[low / confidence: high]** `ab-kit/metrics.py:20-23,29-30` — header says it reports "turns, tokens and subagents"; the token counter `mu` is accumulated and never printed (only `main_peak_ctx` is), and plan D7 says "whether SendMessage resumed it" per subagent: only one total `sendmessage_resumes` over the main transcripts is printed. Lines 23, 29, 30 call `json.loads` unguarded (a blank or half-written transcript line raises, `scan` guards it). Print `mu` or drop it from the header; guard the three loops.
- **[low / confidence: medium]** `context/AGENTS.md:15` vs `:17-18` — DoD 1 puts the verifier "last, after step 2's fixes" but step 3 (docs) and the changes entry follow it; a README or skill edit in step 3 moves the tree after the PASS, which the same sentence calls "a PASS that belongs to no commit" and which preflight Step 2 item 1 accepts only for records under the docs root. Say "after steps 2 and 3's source and docs edits; only records follow".
- **[low / confidence: low]** `ab-kit/run.sh:44-57` — the credential link is made at line 41 but `go env`, `cat` and date can fail between 41 and 54 under `set -e`; the preamble's FAIL then names the link (documented, fine). The script exits 0 whatever claude's `rc` was (final `echo`), so a caller chaining on `$?` sees success for a failed or timed-out arm; the message prints the code. shellcheck: only SC2015 info on line 40 (intended).

## Criteria
1. Dead code — one item: `mu` in metrics.py is computed and unused (finding above); nothing else left behind.
2. Flow tracing — run.sh → preamble → install → claude → rm link → exit trap traced; REPO/kit_repo relative paths resolve to the checkout; metrics.py ran on `control` and `lean2-full` and printed. Gap: preflight Step 2 (finding 1) and DoD 1 vs steps 3-4.
3. Standards conformance — restatement grep over "verifier" shows no other living text still saying "parallel"/"once per task, between phases" except preflight Step 2:57 and `skills/review/SKILL.md` Step 3 (`--tests` stays parallel by D4, consistent). Review Step 5.5 and `agents/reviewer.md`/`security.md` both `effort: high`; no other file names `effort: max` for security. Context size pins pass.
4. Directory layout — ab-kit under `.agents/docs/testing/` beside kit-preamble.sh; scripts live in the docs tree by plan choice; acceptable.
5. Simplicity (DIY-except-security) — stdlib only, no new dependency; no crypto/auth hand-rolled.
6. Decoupling — kit separate from the preamble (sourced, not copied); metrics/accept independent.

## Not examined
Did not run run.sh (it costs real tokens and touches the real credential); the Antigravity side of the changed rules (none of the edits is Claude-only except `SendMessage`/`ToolSearch` in review Step 5.5, which names a fallback); DECISIONS.md, BUG-160, plans and reviews (trail, excluded by the dispatch); README.md beyond the grep.

## Round 2   (2026-10-07; the last round)

Fixes read against the diff, each re-checked by running it where it can run. Evidence: `bash -n` and shellcheck on run.sh (only SC1091 and the intended SC2015 info), both Python files parse, accept.py run against a scratch app whose `main` exits at once (prints FAIL lines and `SCORE 0/25`, no traceback), metrics.py re-run on the new `control` and on `lean2-full` (prints), the new `control` run's result.json and an empty stderr.log (so the `env -i` launch authenticates), and a credential-pattern grep over `.agents/scratch/ab-homes/control` (only `bin/sh-agent.real`, an ELF the preamble skips; no link or token file left).

- Finding 1 (preflight Step 2) — **fixed** at `skills/preflight/SKILL.md:62` (item 1: the verifier is dispatched after the reviewer's and security's source-changing findings are fixed and committed). The line-59 "dispatch in parallel" sentence now governs reviewer and security only; item 1 is the one place the verifier is dispatched and it is specific. Plan D4 and the file now agree.
- Finding 2 (guide § 4 broker off-switch) — **fixed** at `testing/guide.md` § 4: names the broker, both switches, the `env -i` allowlist, and the link; it agrees with the kit header and run.sh:20-24, 50-54. One word of caution only: "no operator token reaches the agent" is true of the environment, not of the linked credential, which the next clause says.
- Finding 3 (accept.py environment) — **fixed** at `accept.py:5,8,13` (MINENV). Checked: build and run both use it; GOTOOLCHAIN=local and a throwaway HOME and GOCACHE under the work dir. Residual, low: `GOFLAGS=-mod=mod` lets `go build` rewrite `go.mod`/`go.sum` in the clone being measured; the build needs no module fetch (standard library only), so `-mod=readonly` would be the safer default.
- Finding 4 (curl check, URLError) — **fixed** at `accept.py:24,28-33,43`: `raw()` sends the body with `Content-Type: application/json`, `req()` and `raw()` return 0 on `OSError`; the dead-app run shows `SCORE 0/25`. Not introduced by the fix, still open and harmless: `starts = [x.get('start') for x in all_]` assumes dict elements.
- Finding 5 (metrics.py) — **fixed**: `mu` gone, `rows()` guards every transcript loop (`metrics.py:5-8,27,33,34`), git runs under `GITENV`. Output verified on two runs. The per-subagent `SendMessage` request is **rebutted** and accepted: a resume is a main-thread call, and plan D7 (line 82) now says a total, counted from the main transcripts. Note for the reader of the numbers: the new `control` run shows `commits: 0` because that agent never committed, not a metrics fault (git run by hand in the clone says the same).
- Finding 6 (DoD 1 "last" vs step 3) — **rebutted**, accepted: a docs-only edit after the verifier changes nothing it ran, and preflight Step 2 item 1 re-triggers the verifier for any source or testing-file change. Caveat left on the record: in this repository `context/` and `skills/` are the living docs and are pinned by `go test`, so a step-3 edit there is source for the gate and the release gate (not DoD 1) is what catches it.
- Finding 7 (run.sh exit code) — `run.sh` still exits 0 whatever claude's code was; the message names it. Left, low, unchanged.

New defects from the fixes: none that fails a run. `env -i` drops proxy and CA variables; this box has none set (count 0), and the control run authenticated, but a box that needs `HTTPS_PROXY` or `NODE_EXTRA_CA_CERTS` must add them to the allowlist in run.sh:50. The real Go caches are no longer exported, so the agent's Go cache is built under the sandbox home (252 MB after one control run), which the preamble scans without a false positive. `timeout --foreground` kills only claude on the 90-minute limit; a server the agent started in the background survives it.

## Final verdict: APPROVE

Every finding is fixed or rebutted with a reason I accept; the residuals above are low and none can fail a run silently.

## Outcome

Written by the main agent, 2026-10-07. Round 2 resumed the same agent with `SendMessage` and approved.

**Fixed:**
- **preflight Step 2:** `skills/preflight/SKILL.md` Step 2 item 1 now dispatches the verifier after the source fixes are committed.
- **guide § 4:** `testing/guide.md` § 4 gains the sentence on the broker off-switch, the `env -i` allowlist and the credential link.
- **accept.py's environment:** it builds and runs the agent's code with a minimal environment.
- **The smaller items:** the curl check replaced by `raw()`, `OSError` caught, the unused `mu` removed, and `rows()` guarding every transcript loop.

**Rebutted**, both accepted in round 2:
- per-subagent `SendMessage` attribution (plan D7 is amended to a total)
- DoD 1 "last" against step 3 docs

**After the cap.** These residuals were fixed by the main agent and not re-reviewed:
- `-mod=readonly` with `GOPROXY=off`
- `run.sh` exits with claude's exit code

A `control` kit run after them scored 25/25, printed both clean lines, and exited 0.

**Left (implementer):**
- `env -i` drops proxy and CA variables. None are set on this box; a box that needs them adds them to the allowlist, as the kit header says.
- A server the agent started survives a timeout (in the kit header).
