---
owner: main
date: 2026-10-08
---
# Plan — kit leftovers in the preamble, a harder seeded security check, the remote.md broker note   (2026-10-08)

## Why
The operator said "go" (2026-10-08) to the v2.5.0 handoff's follow-ups (`handoffs/2026-10-07-1829-v2.5.0-released.md` Next Steps 2.1–2.3):
1. **`kit_leftovers` is pasted in three kits** and matches by working directory only (code review worktree-4, low; security v2.5.0 R2-1, low; BUG-163's stated limit).
2. **The seeded `security` check uses one-line textbook flaws** scored by keyword; it does not test multi-step or fail-open flaws (security v2.5.0, low 3).
3. **`remote.md:175`** does not name the reviews that describe this box's credential broker (security v2.5.0, low).

## Scope
**In:**
- **A. The leftover check moves into `kit-preamble.sh`'s exit check** (`kit_on_exit`), so every kit sourcing the preamble gets it.
  - **What matches:** any process of this uid, other than the kit's own shell and its ancestors, whose physical cwd or executable is under the clone or `KIT_HOME`, or whose environment carries `HOME=$KIT_HOME`. The last one catches a server started from `/tmp` by the agent: the kits run the CLI with `HOME=$KIT_HOME`.
  - **What happens:** each match gets SIGTERM and the run FAILs. The line is "kit: FAIL processes left by the run (BUG-163)", followed by pid and command. The success line is "kit: no processes left by the run".
  - The three kits drop their copies.
  - The preamble header's exit-check sentence names the new check.
  - A Go test, `TestKitPreambleKillsLeftoverProcesses`, is added beside `TestKitPreambleRefusesLeftoverCredentials`.
- **B. The `seed-security` kit gains a second commit of three harder flaws**, reported apart from the six (`hard N/3`):
  1. A signed-cookie check that accepts a cookie with no signature part (fail-open on the missing branch).
  2. A path check `strings.HasPrefix(filepath.Clean(p), dataDir)` with no trailing separator (`/srv/files-evil` passes).
  3. An auth middleware that lets the request through when the token store returns an error (fail-open on error).

  The scorer is validated must-pass and must-fail as before; then 3 runs at `high` and 1 at `max`, and each artifact is read by hand.
- **C. `remote.md:175`** gets one sentence naming the BUG-160 addendum and the 2026-10-07 security reviews as where this box's credential broker is described.

**Out:**
- a release (not asked)
- moving the `GIT_*` clearing into the preamble

## Surfaces touched
- `.agents/docs/testing/kit-preamble.sh` (`kit_on_exit`, header)
- `internal/harness/catalog_test.go` (the new test)
- `.agents/docs/testing/{ab-kit,agy-rel,seed-security}/run.sh`
- `.agents/docs/testing/guide.md` (the sentence describing the preamble's exit lines)
- `.agents/docs/deploy/remote.md`
- BUG-163 (its limit paragraph, now closed)

## Interfaces
- **The preamble exit output gains one line:** `kit: no processes left by the run`, or `kit: FAIL processes left by the run (BUG-163): <pid> <cmd>` (one per process, cmd cut to 80 characters), followed by exit 1.
- **A run counts** only with this line plus the two existing ones (guide § Behavioural runs).
- **On a system with no `/proc`:** `kit: no /proc here: the leftover-process check did not run`, which does not fail the run (the Mac).
- **The seed kit prints** `found N/6` and `hard M/3`.

## Failure modes
- **A false positive** kills a process the operator started. Guard: same uid only; a process counts only if its cwd or exe is under the clone or `KIT_HOME`, or its `HOME` is `KIT_HOME`, which no operator process carries.
- **The kit's own shell or the CLI's still-exiting parent** is never matched: the kit skips `$$` and its ancestors.

## Verification
- `go test ./internal/harness -run TestKitPreamble -count=1`:
  - a planted `sleep` with its cwd in the clone FAILs and is killed
  - a planted `sleep` started from the temp dir with `HOME=$KIT_HOME` FAILs and is killed
  - a clean run prints the no-processes line
  - mutation: with the check disabled, the planted cases pass, so it must fail
- The full suite: `gofmt -l cmd internal && go vet ./... && go test -race -timeout 30m ./...`.
- **E2E:** one `ab-kit` `control` run and one `agy-rel` `operate` row through the preamble, each printing all three exit lines; the seed runs as above.

## Phases
- [x] Cold-reader check — one report: `reviews/2026-10-08-review-kit-leftovers-plan-cold-read.md` (design needs revision, 5 BLOCKING). All answered by D3–D7; none needed the operator.
### Phase 1 — preamble and kits
- [x] A: the preamble check, the test, the kits' copies removed — verify: `go test ./internal/harness -run TestKitPreamble -count=1` plus the mutation — expect: PASS, and FAIL with the check disabled. Result: PASS. Mutations: the matcher disabled fails the 3 planted cases; the real-home refusal removed fails the `/.` case. The full suite caught a BUG-155 pipe in the new code, now fixed.
- [x] Checkpoint — expect: the full suite passes, run bare — evidence: the first run FAILed (`TestScriptsHaveNoEarlyExitPipe`: an environ pipe into grep -q); after the fix, the pins passed; the verifier's full run is the final evidence, 2026-10-08
### Phase 2 — the hard seeds and the docs note
- [x] B: the hard flaws plus their scorer, validated both ways; 4 runs — expect: a score per run, every artifact read. Result: high 3/3 in 3/3, max 3/3 in 1/1, six-mode 6/6 unchanged (`testing/2026-10-08-verify-seed-security-hard.md`).
- [x] C: the remote.md sentence
- [x] E2E: `ab-kit` control and `agy-rel` operate print the three exit lines. Result: both did, 25/25 and PASS, and `ab-kit` control did again on the final preamble, 25/25.
- [x] Review (code+docs), security, verifier last, records, commit, handoff. The reviewer and security each took 2 rounds by resume; the verifier `testing/2026-10-08-verify-kit-leftovers.md` PASSed (560/0/3).

## Risks
Assumptions:
- `/proc/<pid>/environ` is readable for the same uid. Checked by the test.

## Decisions
- D1 — 2026-10-08 — Match on cwd, exe or `HOME=$KIT_HOME`, not "any listener new since the start" (security R2-1's alternative). A listener diff would also catch the operator's own servers started meanwhile — by: implementer.
- D2 — 2026-10-08 — The hard flaws are reported apart from the six, so the 6/6 history stays comparable — by: implementer.
- D3 — 2026-10-08 — The hard flaws (cold read 1): `seed-security/run.sh <label> <effort> [hard]`.
  - `hard` builds the skeleton commit, then ONE commit of the three hard flaws, and dispatches over `HEAD~1..HEAD`, which is that commit only. It prints `hard N/3`.
  - The default mode is unchanged (six flaws, `found N/6`), so D2's comparability holds.
  — by: cold read round 1.
- D4 — 2026-10-08 — The process match (cold read 2):
  - `kit_on_exit` runs `cd /` first, so its own subshells are never in the clone.
  - The check skips `$$`, every ancestor of `$$`, and every process whose PPid chain reaches `$$` (descendants).
  - Every `/proc` read is `|| true`; the function returns only its verdict, and a read error never aborts the trap under `set -e`.
  - It runs BEFORE the credential scan, and waits 1 s after the SIGTERMs so a killed process cannot write past the scan.
  - The kits run the CLI synchronously, so it has exited before the trap.
  — by: cold read round 1.
- D5 — 2026-10-08 — Lines and counting (cold read 3):
  - A run counts only with `kit: no processes left by the run`, plus the existing two lines, whatever the kit printed before the trap.
  - Without `/proc` (the Mac), the run counts with `kit: no /proc here: the leftover-process check did not run`, and its log names that gap.
  - The guide's sentence (:161), the `seed-security` header, BUG-160's addendum and BUG-163's limit paragraph are updated in the same change.
  — by: cold read round 1.
- D6 — 2026-10-08 — The test (cold read 3):
  - Planted processes run with `>/dev/null 2>&1 </dev/null`, their pids are written to a file, and `t.Cleanup` kills them, so the mutation run cannot hang `CombinedOutput`.
  - **The cases:** cwd in the clone; an executable copied into the clone and run from elsewhere; `HOME=$KIT_HOME` from the temp dir. Each FAILs and is gone after.
  - **The negative control:** a `sleep` in the temp dir with the test's own HOME stays alive and the run passes.
  - **The precondition:** `/proc/self/environ` is readable, else `t.Skip`.
  - **The mutation:** the check disabled, so the planted cases pass, so the test must fail.
  — by: cold read round 1.
- D7 — 2026-10-08 — The surfaces add `seed-security/run.sh`'s header, `bugs/BUG-160-*.md`, `bugs/BUG-163-*.md` and `guide.md:161` — by: cold read round 1.
