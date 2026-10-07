---
owner: main
date: 2026-10-07
---
# Plan — lean process 3: Antigravity release runs, security seeded check, real-home install, open choices   (2026-10-07)

## Why
The operator said "Go fix all" (2026-10-07), covering the handoff `2026-10-07-1217-lean-process-2.md` Next Steps:
1. Every skill and agent changed by ADR-054 and ADR-054a needs a run on Antigravity before the next release (AGENTS.md § Release).
2. The real home must be installed.
3. Three operator choices, each taken at its recommendation:
   - security round 2 stays optional, and says so
   - a seeded-weakness check of `security` at effort `high`
   - the kit's no-web-tools flag answers security R2-1, so guide § 185 is not extended

## Scope
**In:**
- **A. `skills/review` Step 4:** state that outside `--fix` (the Definition of Done's review), round 2 may be skipped when every fix has a test that fails without it and the verifier runs after, and that the Outcome then says so.
- **B. A seeded-weakness check.** It runs on a fixture Go service whose one commit plants 6 known flaws:
  1. shell injection through `sh -c`
  2. path traversal in a file handler
  3. a hard-coded admin password compared in code
  4. an admin delete endpoint with no auth check
  5. a non-constant-time token compare
  6. `math/rand` session IDs

  `security` runs at `high` 3 times and at `max` once, as the reference. The score is flaws found, out of 6. The runs go through `testing/ab-kit/`'s isolation; the script is `testing/seed-security/run.sh`.
- **C. The Antigravity runs.** The demo service is rebuilt: a three-file Go HTTP service with `/healthz` and `/hello`, scaffolded, as `skill-matrix.md` § Procedure 1 describes. It lives at `~/sandbox/agy-rel/fixture`.
  - **Skill rows** from `skill-matrix.md`, run on agy only: `feature`, `bug`, `review`, `preflight`, `handoff`, `retro` and `operate`.
  - **Agent rows,** one each, for the agents that have no matrix row of their own:
    - `security`: the `agents-data-not-instructions` security row
    - `documentation`: the `documentation` row
    - `bugs`: covered by the `bug` row, through its record
    - `changelog`: "Dispatch the changelog agent for the last commit …", checked by its `changes/` file
  - **The sandbox home** is installed with `--vendor agy` from this tree. agy's token is linked, never copied. The CLI starts from `env -i` with the broker off. `settings.json` follows `skill-matrix.md` line 14, and the kit sources the preamble.
  - **The result** is each row's check, as n of N, in `testing/2026-10-07-verify-agy-release-rows.md` (`owner: main`).
- **D. Real home:** `sh-agent install` from this checkout, then guide § 5's smoke checks on both CLIs, from a clone of the demo fixture.

**Out:**
- a release (no tag)
- the Claude Code matrix rows (already run in this session's kit runs)

## Surfaces touched
- `skills/review/SKILL.md` Step 4
- new `.agents/docs/testing/seed-security/run.sh`
- new `.agents/docs/testing/agy-rel/run.sh`
- the testing logs
- the real home (`~/.claude`, `~/.gemini`, `~/.local/bin`) in D

## Interfaces
- `bash .agents/docs/testing/seed-security/run.sh <label> <effort>` prints `found N/6`, judged by grep for each flaw's file:line in the artifact.
- `bash .agents/docs/testing/agy-rel/run.sh <row>` runs one row and prints PASS or FAIL.

## Failure modes
- A row fails. It is recorded as FAIL with its output, then fixed as a bug or reported. It is never re-run silently.
- agy auth fails in the sandbox. That is reported as BLOCKED, and the release rule stays unmet.

## Verification
- `go test ./...` and the size pins.
- B scores in the log.
- C as n of N per row.
- D: guide § 5 PASS on both CLIs.

## Phases
- [x] Cold-reader check — one report: `reviews/2026-10-07-review-lean-process-3-plan-cold-read.md` (design needs revision). Answered by D3–D9; none needed the operator.
### Phase 1 — rule edit and seeded check
- [x] A: edit review Step 4 — verify: `go test ./internal/harness -count=1` — expect: PASS. Result: PASS.
- [x] B: seeded check — verify: 4 runs — expect: high finds at least what max finds, or the gap is reported. Result: high 6/6 in 5/5 runs, max 6/6 in 1/1 (`testing/2026-10-07-verify-seed-security.md`).
### Phase 2 — Antigravity rows
- [x] Fixture plus kit — verify: one smoke row — expect: agy authenticates, and the preamble lines are clean. Result: `scaffold` PASS on attempt 2, with `command(*)` for the skill's own loop.
- [x] The rows — expect: n of N per row in the log. Every changed skill and agent has a PASS (`testing/2026-10-07-verify-agy-release-rows.md`). Found and fixed: agy's print-mode main agent ended its turn to wait (a context clause plus `sleep`; `bug` passed after it). The `documentation` DoD row fails on an agy background-task defect (in the handoff).
### Phase 3 — Done
- [x] Review, verifier last, records, commit and push; then D (the real-home install and § 5); handoff. The reviewer and security each took 2 rounds by resume; the verifier ran last and passed (559/0). D is recorded in its own log, before `/preflight`.

## Risks
Assumptions:
- agy authenticates from a sandbox `HOME` with a linked token. Checked by the first smoke row.

## Decisions
- D1 — 2026-10-07 — The real-home install runs after the commit, so the real home gets the committed tree — by: implementer.
- D2 — 2026-10-07 — The three choices are taken at the handoff's recommendations — by: the operator's "Go fix all".
- D3 — 2026-10-07 — The agent rows (cold read 1): `bugs`, `changelog` and `documentation` each get an explicit dispatch row in the kit (`dispatch-bugs`, `dispatch-changelog`, `dispatch-documentation`), and `security` gets `agents-security`. Each row passes only when its artifact exists AND the agy conversation store shows that agent invoked (`strings` of the row home's newest `conversations/*.db` names the agent beside `invoke_subagent`). `skill-matrix.md` gains one `agent-dispatch` row pointing at the kit, and line 7 says the release-row and seeded kits are committed (ADR-010's "scripts are not committed" is about the matrix kit) — by: cold read round 1.
- D4 — 2026-10-07 — The seeded check (cold read 2):
  - `max` is set by `sed` on the sandbox home's installed `security.md` only, and asserted by grep before the run.
  - The scorer is validated must-pass and must-fail before any run, on a synthetic artifact naming all six flaws and one naming none.
  - Every artifact is also read by the main agent, so a keyword match on "this is fine" is caught.
  - With 3 runs at `high` and 1 at `max`, it is reported as a measurement, never as a pass.
  — by: cold read round 1.
- D5 — 2026-10-07 — agy secrets and web tools (cold read 3):
  - The token link is checked as in `ab-kit/run.sh`: still a link after the run, and the real file still a regular file.
  - Web tools on agy are not granted: `settings.json` allows only `command(<name>)`, `read_file` and `write_file`. A web call is auto-denied, and an auto-denied line FAILs the row. So choice 3 holds on both CLIs.
  — by: cold read round 1.
- D6 — 2026-10-07 — The fixture is rebuilt today at `~/sandbox/agy-rel/fixture` (`6779204`). The `scaffold` row's commit is fetched into it before the other rows run. One home per row — by: cold read round 1.
- D7 — 2026-10-07 — Item A stays general ("outside `--fix`, round 2 may be skipped when…"). It is consistent with `context/AGENTS.md:16`'s "at most two rounds", which never required a second round. The text goes in Step 4's Rounds paragraph, which Step 5 already cites — by: cold read round 1.
- D8 — 2026-10-07 — D (the real home) adds `sh-agent status` (no warning on the binary line), the post-update check (`hash -r`, `sh-agent version`, no shadowed line) and the agy store check of guide § 5. The authority is the operator's "Go fix all" in their own turn (2026-10-07), recorded in the log — by: cold read round 1.
- D9 — 2026-10-07 — Per-row seeds:
  - `bug`: the TrimLeft commit
  - `review`: one uncommitted function appended
  - `retro`: three identical tooling corrections in the row home's `feedback.jsonl`
  - `agents-security`: a planted comment commit
  - `dispatch-changelog`: one empty commit

  All of them are in the kit — by: implementer.
- D10 — 2026-10-07 — A new finding goes in this change: agy's print-mode main agent ends its turn to await a subagent or a background command. The same happens on the pre-change tree. The fix is one context clause: never end the turn to wait; use a large `WaitMsBeforeAsync` or poll `manage_subagents`/`manage_task` `list` with `sleep 30`. The context is now 14,936 of 15,000 bytes — by: implementer (the operator's "solve and prove the last issues").
