---
owner: reviewer
date: 2026-10-07
---
# Review — worktree-4-code-docs   (2026-10-07)

Range reviewed: working tree vs HEAD (3293e25): `.agents/docs/testing/ab-kit/run.sh`, `.agents/docs/testing/agy-rel/run.sh`, `.agents/docs/testing/seed-security/run.sh`, `skills/review/SKILL.md` (Step 4); BUG-163 record read for intent. No run.sh executed.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The worst problem: `kit_leftovers` compares the readlinked (physical) cwd with `$KIT_CLONE` as written, so if the kit's output root (`AB_OUT`, `AGY_REL`, `SEED_OUT`, all overridable) contains a symlink component the check silently finds nothing and returns 0 — a false PASS, the "silence is not evidence" shape (reproduced in scratch). Otherwise the function does what the brief asks, its placement and exit-status wiring are correct in all three kits, it interacts correctly with the preamble's EXIT trap, and the Step 4 exemption is consistent with Step 5 and the context.

## Findings
- **[medium / confidence: high]** `.agents/docs/testing/ab-kit/run.sh:48-57` (identical copy at `agy-rel/run.sh:92-105` and `seed-security/run.sh:114-127`) — `case "$cwd" in "$KIT_CLONE"|"$KIT_CLONE"/*)` matches `readlink /proc/PID/cwd` (a physical path) against the unresolved `$KIT_CLONE`. Reproduced in `.agents/scratch`: a `sleep` with cwd in `real/clone` is found and killed when `KIT_CLONE=.../real/clone`, but with `KIT_CLONE=.../link/clone` (link -> real) the same process is missed and the function returns 0. `~/sandbox` is a real directory, so the default paths work; an operator-overridden root through a symlink turns the check into a silent pass. Fix: resolve once inside the function (`root=$(cd "$KIT_CLONE" && pwd -P)`) and compare against `$root`; fail, not pass, when the resolve fails.
- **[low / confidence: high]** the same function in `ab-kit/run.sh:45-57`, `agy-rel/run.sh:92-105`, `seed-security/run.sh:114-127` — pasted three times (verified byte-identical with `diff`). Every kit sources `kit-preamble.sh` just before, whose EXIT trap already runs the BUG-107/BUG-160 "left behind" checks; the fix above must be made three times, and guide § 4 kits written in scratch get no check. Move it into the preamble (one copy) and update the preamble header, which says "The exit check is that EXIT trap and nothing else".
- **[low / confidence: medium]** the same function — it sends SIGTERM and prints "killed" without checking the process died; one that ignores SIGTERM is reported killed while still holding its port. The run FAILs either way, so the verdict is right; only the claim is unverified. Optional: brief wait, `kill -0`, then `kill -9`.
- **[low / confidence: medium]** `skills/review/SKILL.md:94` — "never for a `security` finding of medium severity or higher, whose fix that agent re-checks": the relative clause can be read as a condition (skip allowed when the security agent would not re-check) rather than as the reason. Clearer: "never when the `security` agent reported a medium or high finding; round 2 resumes that agent (Step 5)". Meaning is otherwise consistent: the security agent's scale is high/medium/low so "or higher" is high; Step 5 round 2 already resumes the same lens agent; context DoD 2 ("at most two rounds, the second over the fixes only", `context/AGENTS.md:16`) and the feature/bug skills state no skip rule (grep finds the skip only at SKILL.md:94), so no restatement is left standing.

## Brief answers
- Finds every process with cwd in the clone without touching the kit's shell: yes. `/proc/[0-9]*` is globbed once; `$$` (the kit's main shell, which in seed-security has `cd "$KIT_CLONE"` at line 21, so it is the one process that must be skipped) is skipped; command-substitution subshells inside the loop are not in the snapshot. The prefix match is `"$KIT_CLONE"|"$KIT_CLONE"/*`, so a sibling `control-20` is not matched by `control-2`. A process outside the clone was left alone in scratch. Misses by design: a process that chdir'd out of the clone; other users' processes (readlink denied). Under `set -euo pipefail` every failure inside is guarded (`|| continue`, `|| true`, `[ ] && continue` as a non-final list item).
- Placement: ab-kit after claude exits, credential cleanup and the timeout note, before the summary echo; agy-rel after the check and agent-invocation verdict; seed-security right after the claude subshell, before the artifact lookup. All after the CLI, before the verdict. A concurrent kit's clone is a different directory, so a parallel run is not killed.
- Exit status: ab-kit `lo=0; kit_leftovers || lo=1` then `[ "$lo" = 0 ] || exit 1` (the `||` list keeps set -e quiet; exit 1 beats `exit "$rc"`). seed-security: `lo` is set before the early `exit 1` for a missing artifact (also 1); the final `[ "$lo" = 0 ]` is the script's last command so its status is the script's. agy-rel: the verdict becomes `FAIL (...)`, but the next line overwrites it with the denied/invalid text when that also fails (leftover text lost from the summary line; the FAIL line is already on stderr and status is non-zero) — cosmetic.
- Preamble EXIT trap: `kit_on_exit` takes `rc=$?`; the kit's `exit 1` gives rc 1 and the trap exits 1 (`failed` only adds). The kit sets no trap of its own, and the function neither sets nor needs one. No conflict.

## Criteria
1. Dead code — pass. The `[ -d /proc/self ]` branch is reachable on macOS (a supported platform) and says on stderr that the check did not run, returning 0: skipped and stated, not silent.
2. Flow tracing — pass with the symlink finding (a check that cannot emit in one arrangement). Validated in scratch both ways: planted `sleep` found, killed, rc 1; nothing planted, rc 0; plus the symlink miss. The agy leftover verdict flows to the exit status.
3. Standards conformance — pass, except the preamble header's "exit check is that EXIT trap and nothing else" now sits beside a second check (low finding). No other restatement of the round-2 skip rule.
4. Directory layout — pass; kits stay under `testing/<kit>/`.
5. Simplicity — pass; bash and `/proc` only, no dependency.
6. Decoupling — low finding: one function in three copies belongs in `kit-preamble.sh`.

## Not examined
Whether `agy` itself leaves a background process with cwd in the clone (which would FAIL every agy-rel row); the three run.sh were running and were not executed; the BUG-163 `INDEX.md` row and the other uncommitted files (verify log, security release review) are out of scope; the prefix-sibling case was argued from the pattern, not run (my scratch fixture for it failed to start).

## Outcome

Written by the main agent, 2026-10-07. No round 2: every fix is validated by a planted-process probe in both states, and the verifier runs after.

**Fixed:**
1. **`kit_leftovers` compares physical paths:** it resolves `root=$(cd "$KIT_CLONE" && pwd -P)` and matches against that, in all three kits. Validated:
   - a symlinked clone with a planted process → 1 (FAIL) and the process terminated
   - the symlinked clone with nothing planted → 0
   - a plain clone with a planted process → 1
3. **Wording:** the review Step 4 sentence now reads "never when the `security` agent reported a finding of medium severity or higher: round 2 then resumes that agent over its fixes".
- (low, from 2) The message now says "sent SIGTERM", not "killed".

**Left:**
- (low, 2) The function is still pasted three times instead of living in `kit-preamble.sh`. Moving it changes the preamble's "EXIT trap and nothing else" contract and every kit's sourcing order: a follow-up in the handoff.
