---
owner: main
date: 2026-09-16
---
# Plan — linkshelf simulation: both CLIs build one project under the harness   (2026-09-16)

## Why
The harness has been verified one skill at a time (`testing/skill-matrix.md`). Nobody has watched it carry a project from an empty repository through `/init-project` and `/feature` to the Definition of Done on both CLIs and compared the results against its own rules. The operator asked for that simulation, every gap turned into a harness fix under the bug skill or a Follow-up here, then `/retro`, a reinstall into this box and `/handoff`.

## Scope
In: two empty repositories under the sandbox, one per CLI; `/init-project` then `/feature` headless with the brief below; a comparison of both results against `context/AGENTS.md`, the `feature`, `init-project` and `golang` skills; a fix or a Follow-up per gap; `/retro` (harness scope); install into this box; `/handoff`.
Out: the Mac, the two mirrors, GCP anything.
Follow-ups noted, not built:
- G4 — Claude Code's headless permission classifier refuses the init-project skill's `for … $d` loops (`Contains simple_expansion`) and its `cp -n` chain (multiple operations), so the agent re-did the copies with Write; outcome equal. Note the refusal next to the `.claude` guard remark in the skill's Step 2/5 when the skill is next touched — found in: init run, Claude Code
- G5 — Antigravity's scaffold left `AGENTS.md` § Testing as the template's two comment lines and restated the overview under § Architecture; the skill's "fill the sections" does not say what an empty section becomes — found in: init run, Antigravity
- G10 — a `store.Store` interface stayed exported with no caller after the reviewer had it declared in the consumer; the one-round review after fixes (ADR-022) does not re-read — found in: feature run, Antigravity
- G11 — the golang skill's reference table was not consulted headless (only `SKILL.md` read; `golang-io-web.md` names method-pattern routing, `DisallowUnknownFields`, timeouts); consider naming the two files every HTTP service needs in `SKILL.md` — found in: feature run, Antigravity
- G12 — default data path `data/links.json` inside the repository and not gitignored; a plain `go run` dirties the tree; no harness rule names runtime files — found in: feature run, Antigravity
- G20 — the verifier's `testing/<app>/guide.md` leaves `<app>` unresolved when the binary (`cmd/server`) and the repository (`demo-app`) are named differently: one of three runs at BUG-043's second wording wrote `testing/guide.md`; define `<app>` in the verifier (the `cmd/<app>` name; the repository name when there is none) — found in: real-home check 3b, Claude Code
- G13 — `AGENTS.md` § Setup Commands grew from one test command to a six-tier list (`go test -cover`, `bash scripts/e2e_test.sh`); the documentation agent's lane allows it, no rule bounds it — found in: feature run, Antigravity

Brief given to both CLIs, verbatim: *a local-only Go HTTP service, linkshelf, that stores bookmarks — POST, GET and DELETE /links, JSON, an in-memory store saved to a file — with tests and the deploy guide filled in.*

## Surfaces touched
Sandbox home `/home/kesokaj/sandbox/home` (installed at HEAD before the runs); run repositories `/home/kesokaj/sandbox/runs/sim-<cli>-linkshelf` with bare remotes `/home/kesokaj/sandbox/runs/sim-remotes/`; driver kit `.agents/scratch/sim-kit/` (not committed, ADR-010); this repository's `context/`, `skills/`, `agents/` for the fixes.

## Interfaces
Kit: `setup.sh` (empty repos and remotes), `run.sh <cli> <init|feature>` (one headless run, status and artifacts under `sim-kit/results/<cli>-<step>/`), `drive.sh` (both lanes: init then feature, detached, status file per lane), `compare.sh` (the checklist below, run in each repo, one line per check).

## Failure modes
A run hits a skill gate (three-round cold read, review cap): the run ends with a handoff and that is itself a comparison finding, not a kit failure. A run exceeds its timeout: exit 124, recorded as such, never cited as a result (BUG-013). A tool denial ends an Antigravity run with empty stdout: the stderr notice is the record.

## Verification
Tiers: kit dry run against a throwaway repo (setup, one no-op run) before the real runs; the runs themselves; the comparison checklist; each harness fix under the bug skill carries its own failing test or reproduction and verifier log.
E2E criteria: both repositories hold a service that answers POST, GET and DELETE on /links and survives a restart, built without a human reply; the comparison log names, per CLI, every harness rule the result kept or broke, with the file that shows it.

## Phases
- [x] Cold-reader check — not dispatched: a testing plan the operator specified in one turn; every task here is a run or a read, and the fixes it produces go through the bug skill's own review (2026-09-16)
### Phase 0 — Sandbox at HEAD and the kit
- [x] `sh-agent install --home /home/kesokaj/sandbox/home` at HEAD — verify: `sh-agent status --home …` — expect: `binary cad463d…`, `drift none` — done 2026-09-16 (1 written, 259 unchanged)
- [x] Kit written and dry-run on a throwaway repo — verify: `bash setup.sh -dry && bash run.sh <cli> noop -dry` — expect: two empty repos with `origin` set, a status line per run, no run directory left behind — done 2026-09-16: both CLIs answered `ok` plus the context's first heading (claude 3s, agy 14s), Antigravity log 0 `Failed to parse`, dry repos removed
- [x] Checkpoint — expect: kit ready, both lanes launched detached — evidence: `sim-kit/results/status-claude` and `status-agy` both `init running 2026-09-16T07:49:12+02:00`, `run.sh <cli> init` and both `-p` processes seen with `pgrep`, 2026-09-16
### Phase 1 — The runs
- [x] `/init-project` on Claude Code and Antigravity, answers in the prompt — verify: `compare.sh init` — expect: scaffold committed and pushed in both repos, `deploy/local.md` filled, `gcp.md` and `remote.md` not applicable — done 2026-09-16: Claude 34/39 (symlink guard ×4, trailer), Antigravity 39/39; both pushed
- [x] `/feature` on both, scope gate waived in the prompt, through the Definition of Done — verify: exit status read first, then `compare.sh feature` — expect: a service that passes `go vet ./... && go test -race ./...` in both repos and a trail (plan, reviews, testing, changes, handoff) — done 2026-09-16: Claude rc 0, 4508 s, 6 commits, 90/101 checks; Antigravity rc 0, 1573 s, 5 commits, 94/101; both trails complete
- [x] Checkpoint — expect: both lanes `done`, artifacts under `sim-kit/results/` — evidence: `status-claude` and `status-agy` both `done`, `summary.tsv` rows `claude feature 0 4508s 6`, `agy feature 0 1573s 5`, 2026-09-16
### Phase 2 — Comparison
- [x] Checklist run in both repos: layout, tests, logging, deps, docs trail, review and verifier artifacts, git hygiene — verify: `compare.sh` output saved — expect: one PASS/FAIL line per check per CLI — done 2026-09-16: `sim-kit/results/{claude,agy}-feature/compare.txt`, 90/101 and 94/101
- [x] Code read of both services against the golang skill and review-standards; an independent `reviewer` dispatch per repository — verify: reports under `reviews/` — expect: severity-ranked findings per CLI — done 2026-09-16: `reviews/2026-09-16-review-linkshelf-claude-code.md` (APPROVE), `…-agy-code.md` (APPROVE WITH NITS)
- [x] Transcript audit: skills loaded, agents dispatched, commands denied, per CLI — verify: `audit.py` over the Claude transcript and the Antigravity conversation dbs — expect: a table in the testing log — done 2026-09-16 (the log's § Transcript audit)
- [x] Testing log `testing/2026-09-16-linkshelf-simulation.md` written — expect: both CLIs' results side by side, every gap numbered — done 2026-09-16, G1–G19; § Re-runs filled as the rows land
- [x] Checkpoint — expect: gap list closed, each gap tagged fix or Follow-up — evidence: the testing log § Gaps (7 fixed under BUG-039..043, 6 Follow-ups, 6 informational), 2026-09-16
### Phase 3 — Fixes
- [x] One bug-skill pass per fixable gap (failing test or reproduction, smallest fix, reviewer, verifier, bugs record), committed and pushed each — expect: BUG rows in `INDEX.md`, one commit per fix — done 2026-09-16: BUG-039 (7b8743f), BUG-040 (ff83cf3), BUG-041 (10c2155), BUG-042 (140650c), BUG-043 (7a6ea90); every owed sandbox row run at the fixed tree (testing log § Re-runs)
- [x] Follow-ups for the rest under Scope — expect: each names the gap and why it is not fixed now — G4, G5, G10–G13 filed 2026-09-16; the Claude-only gaps became BUG-043, no further Follow-up
- [x] Checkpoint — expect: every gap fixed or filed — evidence: `git log` 10c2155..7a6ea90 and the testing log § Gaps and § Re-runs, 2026-09-16
### Phase 4 — Close
- [x] `/retro` at harness scope — expect: a record in `retros/` and a line in `~/.config/sh-agent/retros.jsonl` — done 2026-09-16: `retros/2026-09-16-retro.md`, eight proposals queued in `~/.config/sh-agent/proposals/2026-09-16-retro.md`, index line appended; applying them is the gate this session stops at
- [ ] Install into this box (`sh-agent install`), `sh-agent status` — expect: `binary <HEAD>`, `drift none`; the three real-home smoke checks
- [ ] `/handoff` — expect: a cold-readable handoff, `TODO.md` and `ROADMAP.md` empty
- [ ] Checkpoint — expect: pushed, trail closed — evidence: `git status` clean, <date>

## Risks
Assumptions: the sandbox home is still logged in to both CLIs — checked by the dry run, first task of Phase 0. A full feature run may exceed an hour on either CLI — the driver is detached with `setsid nohup` and a status file (BUG-013).

## Decisions
D1 — 2026-09-16 — one empty repository per CLI, each with a bare local `origin`, so the Definition of Done's push can succeed — by: implementer (the brief says "from an empty repository"; a push target is part of the harness's git workflow)
D2 — 2026-09-16 — the brief is passed verbatim, with the DELETE shape and the file format left to the CLI, so the design question reaches the feature skill — by: implementer
D4 — 2026-09-16 — Phase 3 started on the gaps the init runs and the Antigravity lane had settled while the Claude lane was still running, because a reinstall of the sandbox home mid-run would have changed the agents that lane loads; each fix commit names its owed sandbox row — by: implementer (the operator's "keep going")
D3 — 2026-09-16 — Claude's allow-list adds `docker`, `sleep`, `kill`, `rm` and the other shell verbs a service e2e needs; whether a Dockerfile is written is then the agent's choice, not the permission system's — by: implementer
