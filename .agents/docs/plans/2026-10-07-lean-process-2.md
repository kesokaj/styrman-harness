---
owner: main
date: 2026-10-07
---
# Plan — lean process 2: security effort, verifier last, round-2 resume, a safe A/B kit   (2026-10-07)

## Why
Plan `2026-10-07-lean-process` (ADR-054) re-measured at $7.27 / 29.3 min, against $8.46 / 30.8 min. The operator approved the follow-ups and asked to "fix everything or at least figure out what the issues are. Go approved" (2026-10-07). The diagnosis comes from the `lean-full` run's transcripts:

1. **Security effort.** `security` took 19 of 29 minutes. `agents/security.md:6` sets `effort: max`, so each step thought for 21–180 s before acting: round 1 made 11 tool calls in 9 min, round 2 made 6 in 8 min. Its context peaked at about 94k, so reading was not the cost.
2. **A second verifier.** The verifier was dispatched in parallel with review and security, so after the fixes the main thread re-ran it "so its PASS matches what ships". The cause is the order of dispatch, not disobedience.
3. **No resume.** Round 2 was a fresh dispatch. `SendMessage` was available to the `-p` session, but only as a deferred tool that needs `ToolSearch` first, which the rule does not say.
4. **The kit is unsafe.** The A/B kit skipped `kit-preamble.sh` (BUG-162).

## Scope
**In:**
- `agents/security.md`: `effort: max` → `high`, matching `reviewer`.
- **Verifier last:** the one verifier is dispatched after review and security and their fixes, so one PASS covers what ships. The text goes in context DoD 1, `skills/feature` Step 6, `skills/bug` Step 5 and `skills/preflight` Step 1.1.
- **`skills/review` Step 5:** on Claude Code, load `SendMessage` with `ToolSearch` (`select:SendMessage`) when it is deferred, then send.
- **The kit:** rebuilt on `kit-preamble.sh` and committed as `.agents/docs/testing/ab-kit/`, then used to re-measure the full-process arm.

**Out:**
- the real-home install (the operator's job, project AGENTS.md)
- the Antigravity release run (owed before the next tag; feasibility is checked and reported)

## Surfaces touched
- `agents/security.md:6`
- `skills/preflight/SKILL.md` Step 2 (the release gate: verifier after review and security), `.agents/docs/testing/guide.md:181` (stale "never `HOME=<scratch> claude -p`", contradicting :169-170), ADR-054a
- `context/AGENTS.md` (DoD 1)
- `skills/feature/SKILL.md` (Step 6)
- `skills/bug/SKILL.md` (Step 5)
- `skills/preflight/SKILL.md` (Step 1.1)
- `skills/review/SKILL.md` (Step 5.5)
- the new `.agents/docs/testing/ab-kit/{run.sh,accept.py,metrics.py,prompt.md,prefix-full.md}`

## Interfaces
`bash .agents/docs/testing/ab-kit/run.sh <arm>`, where `<arm>` is `control`, `harness` or `<name>-full`.
- The arm's clone is `${AB_OUT:-$HOME/sandbox/ab}/runs/<arm>/app`, outside every checkout.
- `KIT_HOME` is `.agents/scratch/ab-homes/<arm>`. A harness arm gets `sh-agent install --home KIT_HOME --vendor claude` from this checkout.
- The CLI runs with `HOME=KIT_HOME` and `CLAUDE_CONFIG_DIR=KIT_HOME/.claude`, and `GOCACHE`/`GOMODCACHE` are kept from the real environment.
- The credential is a symlink to the real `.credentials.json`, not a copy, because a token refresh in a copy could stale the real home's login. The kit deletes it by name before the preamble's exit check, which must print `kit: no credentials left in KIT_HOME` and `kit: real config unchanged`.
- The scripts are invoked as `python3 -I accept.py <app> <work>` and `python3 -I metrics.py <result.json> <claude config dir> <app>`.

## Failure modes
- A kit run whose exit check FAILs does not count.
- Lower security effort could miss a finding. Measure: the re-run's security still reports, and the acceptance test stays 25/25.

## Verification
- `gofmt -l cmd internal && go vet ./... && go test -race -timeout 30m ./...`
- **The kit:** a `control` run first (cheap). Expect the preamble lines, a 25/25 build, and the real `~/.config/sh-agent` unchanged.
- **Measure:** a `lean2-full` run. Expect 25/25 and less than the $7.27 / 29.3 min it is compared with: one verifier, round 2 by resume or none, security shorter.

## Phases
- [x] Cold-reader check — one report: `reviews/2026-10-07-review-lean-process-2-plan-cold-read.md` (design needs revision / not ready, 5 BLOCKING). All five are answered below by D4–D10; none needed the operator.
### Phase 1 — rules and kit
- [x] Edit the six rule surfaces — verify: `go test ./internal/harness -count=1` — expect: PASS, context under 15,000. Result: PASS, context 14,722 bytes. Also edited: preflight Step 2 item 1 and guide § 4, both from the review.
- [x] Write `ab-kit/` — verify: `bash ab-kit/run.sh control` — expect: the preamble's ok, unchanged and no-credentials lines, and 25/25. Run 1 FAILed the credential check (the box's broker, BUG-160 addendum). Runs 2–4 passed: 25/25 and both clean lines, after the broker off-switch, then `env -i`, then the post-review flags.
- [x] Checkpoint — expect: the full test command passes, run bare — evidence: the verifier's `go test -race -timeout 30m ./...` exit 0 (559 pass, 0 fail), `testing/2026-10-07-verify-lean-process-2.md`, 2026-10-07
### Phase 2 — measure and Done
- [x] Run `lean2-full` against this tree — verify: `accept.py` and `metrics.py` — expect: 25/25, numbers in the changes entry. Result: $4.51 / 14.4 min / 4 agents, 25/25 on the corrected check.
- [x] Review (code and docs in one dispatch), fixes, verifier last, changes entry, commit, push, handoff. The reviewer and security each took 2 rounds by resume; the verifier ran last and passed.

## Risks
Assumptions:
- Claude Code runs with `HOME` overridden when `CLAUDE_CONFIG_DIR` holds the credential. Checked by the `control` kit run, which is the first task.

## Decisions
- D1 — 2026-10-07 — `effort: high` for `security`, keeping `model: opus` — by: the operator's "fix everything … go approved", from the diagnosis above.
- D2 — 2026-10-07 — The credential in the sandbox is a symlink, not the guide's copy, so a refresh writes to the real file. `testing/guide.md` § 4 gets one sentence on it — by: implementer.
- D3 — 2026-10-07 — The DoD step numbers stay, and "dispatched last" is said in DoD 1, so no "step N" reference elsewhere moves — by: implementer.
- D4 — 2026-10-07 — Verifier last (Q1) applies wherever a verifier and fixes can follow a review: DoD 1 (and so feature and bug), and `/preflight` Step 1.1 and Step 2. `/review --tests` stays parallel, because it is a report lens with no fix loop unless `--fix`. Under `--fix` its re-dispatch is round 2 like any lens — by: cold read round 1.
- D5 — 2026-10-07 — The credential symlink (Q2), from evidence in this session:
  - In the `lean-full` run, Claude Code refreshed the token through the symlink at 09:15. The real `~/.claude/.credentials.json` stayed a regular file with a new mtime, and the link was still a link at cleanup (`find -type l` matched it).
  - The kit checks that the real file is a regular file before and after, and that the link is still a link after.
  - A killed run leaves the link. The preamble's exit check flags it by name (FAIL), which is the loud outcome wanted, and the kit's header says to delete it by name.
  — by: cold read round 1.
- D6 — 2026-10-07 — `HOME=KIT_HOME` (Q3):
  - The git identity is set by `GIT_AUTHOR_*` and `GIT_COMMITTER_*`.
  - `GOCACHE` and `GOMODCACHE` keep the real values; the other XDG directories move with HOME, as § 4 intends.
  - The baseline was measured with the real HOME. The comparison is therefore stated as approximate, n=1, never as a proof.
  — by: cold read round 1.
- D7 — 2026-10-07 — Evidence (Q4): `metrics.py` reports for each subagent its type, its start and end time, and its tool calls, plus the run's total of `SendMessage` resumes. A resume is a main-thread call, so the total is counted from the main transcripts (amended after the code+docs review, round 1). The `lean2-full` run is the behavioural run of `security.md` (n=1 of 1: the per-turn thinking gaps show the effort). The context change is checked by grep of the installed file. No separate sandbox start — by: cold read round 1.
- D8 — 2026-10-07 — A run is wrapped in `timeout 90m`. A timed-out or failed run is recorded as such, not re-run silently, and a re-run needs the run directory moved away (the kit refuses an existing one) — by: cold read round 1.
- D9 — 2026-10-07 — ADR-054a amends ADR-054: security effort `high`, the verifier dispatched last — by: implementer.
- D10 — 2026-10-07 — Context headroom (281 bytes): the DoD 1 edit replaces words rather than adding them, and the measured byte count is recorded — by: implementer.
