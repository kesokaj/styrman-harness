---
owner: reviewer
date: 2026-10-07
---
# Review — lean-process-2 plan, cold read   (2026-10-07)

Range reviewed: `.agents/docs/plans/2026-10-07-lean-process-2.md` (untracked, HEAD f82457b) and the surfaces it names: `agents/security.md`, `context/AGENTS.md` DoD, `skills/{feature,bug,preflight,review}/SKILL.md`, `.agents/docs/testing/kit-preamble.sh`, `testing/guide.md` § 4, the scratch kit `.agents/scratch/ab-kit/{run.sh,metrics.py,accept.py(head),prefix-full.md}`, BUG-162, ADR-054.

## Status: COMPLETE

## Verdict: REQUEST CHANGES (design needs revision; 5 blocking questions, all answerable by D-rows)

Worst first: the plan's headline measurement cannot show what it claims, and its sandbox kit has two failure modes that touch the operator's real login and real config. (a) The credential symlink (D2) is safe only if Claude Code refreshes a token by writing through the link; if it writes a temp file and renames, the symlink is replaced by a regular file in KIT_HOME and the real token is the stale one — the exact harm D2 says it avoids. Nothing in the plan checks it, and the preamble fingerprint covers only `~/.config/sh-agent`, never `~/.claude/.credentials.json`. (b) Overriding `HOME` drops `~/.gitconfig`, so the arm's agent cannot commit, and the run is then not comparable with the earlier one. (c) Round 2 "by resume or none" and "security shorter" are outcomes of one stochastic run, not conditions the verification produces, and `metrics.py` as it exists cannot report them. (d) The restatement sweep misses `skills/review/SKILL.md:41` and `skills/preflight/SKILL.md:46,57`, which keep the verifier parallel with review.

## Findings

1. **[high / confidence: high] Restatements left standing, verifier "parallel".** The plan edits `skills/preflight` Step 1.1 only.
   - `skills/preflight/SKILL.md:44` is the Step 1.1 text to change.
   - `skills/preflight/SKILL.md:46` (Step 1.3) says "Make no further edits between dispatching the rows and committing; if the tree moves, re-run them". With verifier last, review fixes are edits after the review and before the verifier; the sentence must be reworded or it contradicts the new order.
   - `skills/preflight/SKILL.md:57` (Step 2) dispatches reviewer, security and verifier "in parallel". The plan does not say whether the release gate stays parallel (fixes there mean NO-SHIP and a fresh gate, so parallel may be right, but that is a decision).
   - `skills/review/SKILL.md:41` (Step 3), "one agent per lens, all in parallel", includes the `--tests` verifier row at `:54`; under `--all --fix` that verifier runs on the pre-fix tree. Step 5.3 (`:107`) runs checks bare but never re-dispatches the verifier.
   - `skills/feature/SKILL.md:113` says "one `verifier` log over the whole feature" with no order, and `skills/bug/SKILL.md:55` lists `verifier` first. Both need the order, as the plan says.
   - Fix: add the review Step 3 `--tests` row, preflight Steps 1.3 and 2 to Surfaces, each with a decision. Run `skills/review/scripts/restatement-grep.sh -i -- verifier parallel dispatched` first.

2. **[high / confidence: medium] D2's credential symlink is an unchecked assumption that can log the operator out.** Plan lines 40 and 68. Claude Code's refresh may write atomically (temp file plus rename); that replaces the symlink with a regular file holding the new token, and the real file keeps a rotated-out refresh token. `kit-preamble.sh:74-93` fingerprints only `~/.config/sh-agent`, so "real config unchanged" says nothing about `~/.claude/.credentials.json`. The plan also gives no path for a killed or failed run: guide § 4 forbids the kit its own EXIT trap (`kit-preamble.sh:24-27`), and under the preamble's `set -e` a nonzero `claude` skips the deletion line, leaving the symlink for the exit check to FAIL on.
   - Fix: `run.sh` records `cksum`+`stat` of the real credential before the run. After `claude ... || true` it tests `[ -L "$KIT_HOME/.claude/.credentials.json" ]` before deleting. A regular file there prints `kit: FAIL credential was replaced, real login may be stale` and says so loudly. The Phase 1 control run asserts the real file's cksum is unchanged. State the symlink-or-copy result in D2 only after that measurement.

3. **[high / confidence: high] `HOME=KIT_HOME` changes the environment the A/B compares.** Plan line 39 keeps only `GOCACHE`/`GOMODCACHE`. The old kit (`.agents/scratch/ab-kit/run.sh`) ran with the real `HOME`, so the $7.27 / 29.3 min baseline came from a different environment.
   - No `~/.gitconfig`: the arm's agent cannot commit (the DoD ends in commit), or commits fail with "tell me who you are". Needs `GIT_CONFIG_GLOBAL` pointing at a file written into KIT_HOME, or `GIT_AUTHOR_*`/`GIT_COMMITTER_*`.
   - `GOPATH`, `GOFLAGS`, `GOTOOLCHAIN` and `XDG_CACHE_HOME` default under HOME. The Go SDK on PATH is the user-level one (project AGENTS.md "Setup commands"). The kit must read `go env GOCACHE GOMODCACHE` before overriding HOME, or the scan walks gigabytes.
   - The preamble's wrapper puts an `sh-agent` first on PATH, so the `control` arm (no harness) now has the command; the plan must say that is accepted.
   - Fix: list the exported variables in Interfaces and say the old $7.27 / 29.3 min figures are not an environment-matched baseline (see 4).

4. **[high / confidence: high] The Phase 2 verification does not produce the condition it claims.** Plan lines 49-50 and 59.
   - "less than $7.27 / 29.3 min" compares one stochastic run with one earlier run, in another environment (3), another harness build, and with `claude-opus-5-5` hardcoded in the old `run.sh` but unspecified here. One run cannot show a 4% improvement against run-to-run variance of 2-3x (ADR-054's own table: $1.11 vs $8.46).
   - "round 2 by resume or none", "one verifier" and "security shorter" are things the arm's model chooses to do. `metrics.py` (`.agents/scratch/ab-kit/metrics.py:25-34`) reports only counts and totals per agent type. It has no per-subagent wall time, per-subagent tool-call count, dispatch order or `SendMessage`/`ToolSearch` tool names (`tools` is the main thread's top 8 only), so it cannot show the three things the diagnosis was about.
   - Fix: reword the expectation to what the run can show: 25/25; sub_types has exactly one `verifier`; its first timestamp is later than the last `reviewer` and `security` timestamps; `security` wall time and tool calls are reported (diagnosis baseline: 19 of 29 min, 11 and 6 calls) with the effort value recorded from the installed `agents/security.md`; `SendMessage` appears in the tools (or its absence is reported as "no round 2"). Add those fields to `metrics.py` and name them in Interfaces. Call the result "one sample, not a comparison". Record the measured tree (`git rev-parse HEAD` plus a diff hash) in the changes entry, since the run precedes review and fixes (Phase 2 order).

5. **[medium / confidence: high] No behavioural check for the agent and context changes the project requires.** Project `AGENTS.md` § Testing and guide § 4 § Behavioural runs (`guide.md:176`): a change to `agents/*.md` gets a `claude -p --agents` run reported as n of N; a `context/` rule that changes what the main agent does gets an `--append-system-prompt` run plus the must-fail control (outside this repo); `/review --docs` over the range; and, for `agents/` and `skills/`, the sandbox start of each CLI. The plan's verification names none of these. Its one LLM run is the A/B, which is neither n of N nor a probe of the definition. The Antigravity exclusion lists only the release run; the sandbox start is separate and cheap.
   - Fix: add a `security` behavioural row (n of N) that checks the installed definition's `effort: high` took effect (the transcript or meta). A SendMessage probe: does `ToolSearch select:SendMessage` load it, and does it resume an agent that has already reported (3 of 3)? And the guide § 4 sandbox start for both CLIs, or an explicit "not run".

6. **[medium / confidence: medium] `SendMessage` against a stopped agent.** Context line 35 requires stopping an agent that reported while its task still runs (`TaskStop`); `skills/review/SKILL.md:108` then resumes it. Whether a stopped or finished agent accepts `SendMessage` is unstated and unchecked (the plan only verified it is deferred). If it cannot, the rule needs "else re-dispatch", which is already the fallback; the load step must say what to do when ToolSearch finds nothing, and that it is Claude-only (Antigravity has neither tool).

7. **[medium / confidence: medium] Harness arm install path unverified until the $7 run.** Phase 1 verifies only `control`. The harness arm adds: `sh-agent install --home KIT_HOME --vendor claude` (exists: `cmd/sh-agent/main.go:205`); the installed skills, `AGENTS.md` and `sh-agent` binary under KIT_HOME; and the preamble's credential scan (`kit-preamble.sh:104-130`) over all of it. The scan's `kit_cred_text` includes `"refresh_token"`, `AKIA...` and `sk-ant-`; installed skill or context text that quotes them (preflight Step 1.5 lists them) would FAIL the exit check, and it is not an ELF so it is not skipped. A $7 run could end with a FAIL line. Fix: Phase 1 also runs `run.sh harness` in an install-only, no-CLI mode (or runs the preamble's `kit_credential_leftovers` over an installed KIT_HOME), expecting the no-credentials line.

8. **[medium / confidence: medium] Surfaces list omits `testing/guide.md` and the contradiction it carries.** D2 adds a sentence to § 4, but the Surfaces list does not name the file. Guide `:169-170` already says a scratch `HOME` with `CLAUDE_CONFIG_DIR` and a credential copy works; `:181` says "A CLI started from a scratch `HOME` cannot authenticate". The plan's Risks section treats the first as an unproven assumption. The control run settles it, and the guide must then be made consistent in one edit, with the second sentence reworded (copy vs symlink, and the CRM A/B kit's location). Add `guide.md:169-170,181` and the testing index (if any) to Surfaces.

9. **[medium / confidence: medium] Context headroom is 281 bytes.** `TestRenderedContextUnderTarget`: hand-written context 14,719 of 15,000 (checked now). DoD 1 already says "once per task (between phases and after fixes, run the tests bare yourself)" and ADR-054 already says the verifier "runs at the Definition of Done". "Dispatched last" must fit in the margin or replace words in place. Write the sentence before Phase 1 and measure.

10. **[medium / confidence: medium] Kit run controls undecided.**
    - Re-running an arm: the old kit refused if `runs/<arm>` existed; the new `KIT_HOME` and clone paths need the same rule or a clean-up rule. The preamble does `mkdir -p KIT_HOME/bin`, so a reused home accumulates transcripts and `metrics.py` would add the old ones.
    - The credential deletion must be inline, after `claude`, under `|| true` (see 2). A timeout or `--max-budget-usd`-style cap for a $7 run is not mentioned, and the old kit has none.
    - `metrics.py` takes `<claude config dir> <app>` but also reads `wall_seconds` from the run directory (old code: `{W}/wall_seconds`); the new signature gives it no path.
    - `metrics.py:loc()` builds shell strings from the app path (`shell=True`, f-string); AB_OUT with a space or quote breaks it.
    - The preamble needs an absolute KIT_HOME, and the plan's `.agents/scratch/ab-homes/<arm>` is relative; resolve with `$PWD` or the checkout root.

11. **[low / confidence: medium] No ADR amendment.** ADR-054 (`DECISIONS.md:606-628`) states "One verifier per task. It runs at the Definition of Done" and the resume rule without ToolSearch. The plan changes both and a DoD ordering; AGENTS.md says a structural decision gets an ADR. D1-D3 are plan rows only. Recommend ADR-054a (or a superseding ADR) in the same commit.

12. **[low / confidence: high] Plan Phase 2 mixes the measurement with the tree it changes.** "Run `lean2-full` against this tree" occurs before review and fixes, so the measured harness is not the shipped one. Either measure after fixes (recommended: one more cheap decision, but then the verifier is last and the A/B is last) or state that it measures the pre-review tree.

## Questions (BLOCKING or NON-BLOCKING, with recommended answers)

| # | Q | Class | Recommended answer |
|---|---|---|---|
| Q1 | Do preflight Steps 1.3 and 2 and review Step 3 stay parallel for the verifier? | BLOCKING | Step 1 (dirty tree): reorder, verifier after review fixes, and reword 1.3. Step 2 (release gate): stay parallel, since a finding there is NO-SHIP and the gate is re-run; say so in a D-row. Review Step 3: `--tests` joins the parallel round for a report-only review; with `--fix` the verifier runs after Step 5.3. |
| Q2 | Symlink or copy for the credential, and how is a refresh detected? | BLOCKING | Keep the symlink only with the `[ -L ]` check plus real-file cksum before/after in `run.sh`; a regular file there means FAIL and a loud message; control-run asserts unchanged. |
| Q3 | Which environment variables does the arm get once HOME is KIT_HOME (git identity, Go dirs, model)? | BLOCKING | `GIT_CONFIG_GLOBAL` to a file written in KIT_HOME (user name and email only, `commit.gpgsign=false`), `GOCACHE`/`GOMODCACHE` taken from `go env` before the override, model pinned to what the baseline used, all listed in Interfaces. |
| Q4 | What does the Phase 2 measurement assert, and what fields does `metrics.py` add? | BLOCKING | See finding 4: per-subagent first and last timestamps, tool-call count, order, SendMessage use; assertions limited to 25/25 and one verifier after review/security; durations reported, not asserted; "one sample". |
| Q5 | Which behavioural/sandbox checks run for the `security.md` and context changes? | BLOCKING | `security` n of N (3) behavioural row, a SendMessage-after-ToolSearch probe (3), guide § 4 sandbox start both CLIs or named "not run"; `/review --docs` over the range (already in Phase 2 review dispatch). |
| Q6 | Does a harness-arm install pass the preamble's credential scan? | NON-BLOCKING | Add a no-CLI check to Phase 1 (finding 7). |
| Q7 | Does SendMessage reach a stopped agent; what if ToolSearch has no result? | NON-BLOCKING | Probe it (Q5); rule text: "else re-dispatch as today"; Claude-only. |
| Q8 | Rerun policy, timeout and budget cap, `wall_seconds` path, absolute KIT_HOME, `loc()` quoting. | NON-BLOCKING | Refuse an existing `runs/<arm>` and `ab-homes/<arm>` as the old kit did; `timeout 45m`; metrics takes the run directory; resolve KIT_HOME against the checkout root; `subprocess` lists without a shell. |
| Q9 | Context headroom (281 bytes): the sentence. | NON-BLOCKING | Draft and measure it first; rewrite in place. |
| Q10 | ADR for the change? | NON-BLOCKING | ADR-054a in the same commit. |
| Q11 | Measure before or after the fixes? | NON-BLOCKING | After the fixes, before the verifier, and record the tree. If before, say so. |
| Q12 | `guide.md` Surfaces and the `:169-170` vs `:181` contradiction. | NON-BLOCKING | Add to Surfaces, edit in one pass. |

## Scope: what the plan builds that `## Why` did not name

- Why names "the kit is unsafe" (4). It does not name: the `HOME` override and its environment (`GIT_CONFIG_GLOBAL`, Go variables), the committed `ab-kit/` set of five files (the old kit is explicitly "not committed" in plan 1's D14, so committing it reverses that D-row: say `supersedes D14`), the credential-symlink design (D2), `metrics.py`'s new interface, a guide change, and the `testing/` placement that makes the kit a "script under the docs root" which preflight Step 2 item 1 requires committed before a verifier runs.
- Why names a diagnosis of 19 of 29 min in `security` but no per-agent timing in the measuring tool (4).
- The new `control` run is justified as the kit's own check, but it is a spend ($0.54 per ADR-054) outside the Why; fine, name it.

## Criteria (applied to the plan as a design)

1. Dead code — n/a for the plan itself; the scratch kit is superseded by `ab-kit/`, and the plan does not say the scratch copy is deleted (scratch is gitignored; non-issue).
2. Flow tracing — fail: restatements of the verifier order (finding 1), the credential flow on failure (2), the environment (3).
3. Standards conformance — fail: behavioural/sandbox checks missing (5), ADR (11), guide contradiction (8).
4. Directory layout — pass: `testing/ab-kit/` beside `kit-preamble.sh` follows the guide. The mode of `run.sh` (0644 like `kit-preamble.sh`, run with `bash`) should be said.
5. Simplicity (DIY-except-security) — pass: stdlib Python, bash, no new dependency; the credential handling uses the CLI's own login, nothing hand-rolled.
6. Decoupling — pass with a note: `metrics.py` reaching into the CLI's transcript layout is a seam only the kit uses; fine.

## Readiness

design needs revision
implementation not ready

(Ready after Q1-Q5 are answered by D-rows and the Surfaces list is extended with the files in finding 1 and 8.)

## Not examined

- `accept.py` beyond its first 25 lines (25/25 contract not re-derived); `prompt.md` not opened.
- `skills/bug/SKILL.md` and `skills/feature/SKILL.md` read at the verifier-bearing lines only; `skills/preflight/SKILL.md` Steps 3-6 not read.
- No run of any kit or CLI: whether Claude Code refreshes a credential through a symlink, whether `SendMessage` resumes a stopped agent, and whether the harness arm passes the credential scan are unknown, all raised as questions.
- The earlier cold-read report for plan 1 (`.agents/docs/reviews/2026-10-07-review-lean-process-plan-cold-read.md`) was not read.
- `internal/harness` kit-preamble tests were located (`catalog_test.go:1543,1655`), not read.
