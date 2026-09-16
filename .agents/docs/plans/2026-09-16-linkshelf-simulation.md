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
Follow-ups noted, not built: (filled as gaps are found)

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
- [ ] `/init-project` on Claude Code and Antigravity, answers in the prompt — verify: `compare.sh init` — expect: scaffold committed and pushed in both repos, `deploy/local.md` filled, `gcp.md` and `remote.md` not applicable
- [ ] `/feature` on both, scope gate waived in the prompt, through the Definition of Done — verify: exit status read first, then `compare.sh feature` — expect: a service that passes `go vet ./... && go test -race ./...` in both repos and a trail (plan, reviews, testing, changes, handoff)
- [ ] Checkpoint — expect: both lanes `done`, artifacts under `sim-kit/results/` — evidence: status files and `summary.tsv`, <date>
### Phase 2 — Comparison
- [ ] Checklist run in both repos: layout, tests, logging, deps, docs trail, review and verifier artifacts, git hygiene — verify: `compare.sh` output saved — expect: one PASS/FAIL line per check per CLI
- [ ] Code read of both services against the golang skill and review-standards; an independent `reviewer` dispatch per repository — verify: reports under `reviews/` — expect: severity-ranked findings per CLI
- [ ] Transcript audit: skills loaded, agents dispatched, commands denied, per CLI — verify: `jq` over the Claude transcript, `sqlite3` over the Antigravity conversation — expect: a table in the testing log
- [ ] Testing log `testing/2026-09-16-linkshelf-simulation.md` written — expect: both CLIs' results side by side, every gap numbered
- [ ] Checkpoint — expect: gap list closed, each gap tagged fix or Follow-up — evidence: the testing log, <date>
### Phase 3 — Fixes
- [ ] One bug-skill pass per fixable gap (failing test or reproduction, smallest fix, reviewer, verifier, bugs record), committed and pushed each — expect: BUG rows in `INDEX.md`, one commit per fix
- [ ] Follow-ups for the rest under Scope — expect: each names the gap and why it is not fixed now
- [ ] Checkpoint — expect: every gap fixed or filed — evidence: `git log`, <date>
### Phase 4 — Close
- [ ] `/retro` at harness scope — expect: a record in `retros/` and a line in `~/.config/sh-agent/retros.jsonl`
- [ ] Install into this box (`sh-agent install`), `sh-agent status` — expect: `binary <HEAD>`, `drift none`; the three real-home smoke checks
- [ ] `/handoff` — expect: a cold-readable handoff, `TODO.md` and `ROADMAP.md` empty
- [ ] Checkpoint — expect: pushed, trail closed — evidence: `git status` clean, <date>

## Risks
Assumptions: the sandbox home is still logged in to both CLIs — checked by the dry run, first task of Phase 0. A full feature run may exceed an hour on either CLI — the driver is detached with `setsid nohup` and a status file (BUG-013).

## Decisions
D1 — 2026-09-16 — one empty repository per CLI, each with a bare local `origin`, so the Definition of Done's push can succeed — by: implementer (the brief says "from an empty repository"; a push target is part of the harness's git workflow)
D2 — 2026-09-16 — the brief is passed verbatim, with the DELETE shape and the file format left to the CLI, so the design question reaches the feature skill — by: implementer
D3 — 2026-09-16 — Claude's allow-list adds `docker`, `sleep`, `kill`, `rm` and the other shell verbs a service e2e needs; whether a Dockerfile is written is then the agent's choice, not the permission system's — by: implementer
