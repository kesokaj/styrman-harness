---
owner: main
date: 2026-09-29
---
# Plan — `sh-agent search`: release readiness   (2026-09-29)

## Why

The operator, 2026-09-29, after the search index plan (`plans/2026-09-29-search-index.md`) closed: "I want everything to be done before we cut a new version. Nothing waiting, nothing in todo roadmap or waiting." The open items are that handoff's Next Steps (`handoffs/2026-09-29-2102-search-index-done.md`): the main agent does not use `sh-agent search` in real work (0 of 28 sessions, D36 of the search plan), the Antigravity runs the release gate needs are incomplete, and the benefit on a project with a large trail was never measured over several sessions.

## Scope

In:
- A clean A/B of the whole harness per arm on Claude Code: `--setting-sources project,local` (measured 2026-09-29: user-level skills and context not loaded, project skills loaded), each arm's installed `.claude/{skills,agents,commands}` copied into the run's project, its rendered context appended once. This removes the first comparison's two caveats (v1.12.1 skills in both arms; the context loaded twice).
- Main-agent adoption: measured with the new skills in place; if the main agent still does not search, one change to the context or a skill (D2), measured again.
- A multi-session test on a large trail: three sessions of ordinary work in a clone of this repository at v1.12.1, per arm.
- Every Antigravity run `AGENTS.md` § Release requires: one completed run of each changed skill and agent.
- Reviews and verification of anything this plan changes; docs; the real-home smoke after the operator's install on this box.
- Trail closed: `TODO.md` and `ROADMAP.md` empty, the handoff with nothing queued, so `/preflight v1.13.0` is the only next step.

Out: the release itself (`/preflight` from the operator's turn); the known limits D25 (`x/sys` advisory, needs Go 1.25) and D35 (the edit-during-lock race), both documented decisions; installs on the other machines (after the release).

Follow-ups noted, not built: none yet.

## Surfaces touched

- `.agents/scratch/eff-kit/` (kit, gitignored): a clean-arm mode for `claude-session.sh` and `lane.sh`.
- `context/AGENTS.md` or one skill, only if D2 changes wording.
- `.agents/docs/testing/` logs; `testing/skill-matrix.md` if a scenario changes.

## Interfaces

Clean-arm Claude session: `cd <project>; PATH=<arm wrapper dir>:$PATH claude -p --setting-sources project,local --append-system-prompt "$(cat ctx-<arm>.md)" --permission-mode acceptEdits --allowedTools <the Phase 5 list> < prompt`, where `<project>/.claude/` holds a copy of `<arm home>/.claude/{skills,agents,commands}` from `sh-agent install --vendor claude --home <arm home>` of that arm's tree, excluded from git by `<project>/.git/info/exclude`. No `--agents`, no `--add-dir`.

Large-trail lane (per arm, one clone of this repository at v1.12.1 with its `.agents/docs` trail, the drift seed of the search plan's Phase 5, and a bare origin): s1 "resume from the trail, then fix the seeded drift: `sh-agent update --check` exits 2 when behind — make every living file agree, as a Task under the Definition of Done"; s2 "resume from the trail, then answer the five trail questions of the search plan's Phase 5 in ANSWERS.md with file:line, commit"; s3 "resume from the trail, then add one row to `testing/sh-agent-matrix.md` for `sh-agent version` on a zip place, with the evidence you can produce, commit". Same waivers text as the sample apps.

## Failure modes

- A run auto-denied on Antigravity: the command is added by name and the run repeated (project `AGENTS.md` § Testing).
- `agy -p` returns before a dispatched subagent: the run is repeated up to three times; still void after three, the gate is told in the log, and the operator decides (Open row).

## Verification

E2E criteria:
- The clean A/B's per-session table exists for Claude Code with both arms delivering; the main agent's `sh-agent search` count is reported for each arm.
- Every changed skill and agent has one completed Antigravity run on disk (artifact or answer) with no auto-denied line.
- After the operator's install on this box, the three real-home smoke checks pass on both CLIs, and `sh-agent search --refs ADR-040` answers from the installed binary.
- `TODO.md` and `ROADMAP.md` hold headers only; the newest handoff's Next Steps name only `/preflight v1.13.0`.

## Phases

- [ ] Cold-reader check — design ready and implementation ready, or the operator's same-turn override quoted with its date

### Phase 1 — Clean A/B on Claude Code, with the new skills
- [ ] Kit: clean-arm mode; dry run shows each arm's project `.claude/skills` holds that arm's skills (the search arm's `feature` skill contains `sh-agent search --refs`, the baseline's does not) — verify: AI — expect: both checks by `grep -c`
- [ ] Sample apps, both arms, 2 runs each, 3 sessions — verify: AI, `outcomes.py`, `metrics.py` — expect: tables per session, main-agent search count per arm
- [ ] Large-trail lane, both arms, 2 runs each — verify: AI — expect: drift fixed, answers scored, main-agent search count
- [ ] Checkpoint — expect: `testing/<date>-search-adoption.md` committed — evidence: <path>, <date>

### Phase 2 — Main-agent adoption (only if Phase 1's main-agent count is low)
- [ ] D2 decided with the operator: the change and its exact text — who: both
- [ ] The change, `TestRenderedContextFitsAntigravityBlock` PASS if context — verify: unit — expect: PASS
- [ ] Re-run of the probe (`search-first`, 3 with, 3 control) and one large-trail lane per arm — verify: AI — expect: n of N reported
- [ ] Checkpoint — expect: log and change committed; `/review --docs` over the change — evidence: <path>, <date>

### Phase 3 — Antigravity runs the release gate needs
- [ ] One completed run of each changed agent without one yet (`dreamer`; `security` on the final wording) and each changed skill (`review-standards`, `dream`, `handoff`, `bug`, `feature`, `retro`, `research`, `init-project`), in a sandbox home installed from this tree — verify: AI — expect: artifact or answer on disk, no auto-denied line, up to three attempts each
- [ ] Checkpoint — expect: `testing/<date>-search-agy-gate.md` committed — evidence: <path>, <date>

### Phase 4 — Real home and close
- [ ] Operator runs `sh-agent update` on this box — who: user — expect: `sh-agent status` binary line without a warning
- [ ] Real-home smoke on both CLIs (`testing/2026-09-03-verify-real-home.md`) plus `sh-agent search --refs ADR-040` — verify: AI — expect: PASS, GANNET
- [ ] Definition of Done: verifier if source changed, reviews for what changed, changes entries, plan closed, `TODO.md`/`ROADMAP.md` empty, handoff naming only `/preflight v1.13.0` — verify: AI — expect: all on disk and pushed

## Risks

Assumptions: `--setting-sources project,local` excludes the real home's skills and context and loads project skills — checked 2026-09-29 by a probe (a project skill listed, no harness skill, no GANNET). Headless Claude Code reads `.claude/` in the project (only writes are refused) — checked by the same probe.

## Decisions

D1 — 2026-09-29 — The clean A/B replaces the first comparison's setup for everything in this plan; the first comparison stays as recorded — by: implementer, from the probe.
Open: D2 — the main-agent adoption change, if Phase 1 shows it is needed — default: propose a step in the context's resume procedure with an equal cut, put to the operator before editing.
