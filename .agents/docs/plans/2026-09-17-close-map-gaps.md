---
owner: main
date: 2026-09-17
---
# Plan — close the dashed gaps on the harness map   (2026-09-17)

## Why

The map in `README.md` § How it works (R9, ADR-031) draws two gaps as dashed lines: nothing reads the reports other places send to `styrman-harness-feedback` (they are read by hand), and `/retro --harness` consumes only the proposal queue of the machine it runs on. So the Report back procedure ends at a repository nobody reads, and a proposal queued on the work Mac never reaches the harness. The operator asked: "make a plan to close the dashed gaps on the map", and confirmed the eight scope lines below with "go all" (2026-09-17).

## Scope

In:
- `/retro --harness` also consumes, from the feedback checkout, the newest report of every host other than this one: its queued proposals become findings; its `feedback.jsonl` and `retros.jsonl` feed only the promotion check (D1–D3, D6).
- The feedback checkout is brought current first with `git pull --ff-only`; when it is missing or the pull fails, the retro says so and continues with the local queue (D4).
- The harness repo's retro records are the cross-machine ledger: a consumed proposal is recorded with its source, and a proposal already recorded applied or rejected is skipped (D5).
- The map: both dashed lines become solid, the "Read reports by hand" node goes, with ADR-032, in ONE commit with the skill and docs (R9, D8).

Out: a timer; a new `sh-agent` verb; a dashboard; any Go code; writing to the feedback repository or to another machine's queue.

Follow-ups noted, not built: none.

## Surfaces touched

- `skills/retro/SKILL.md` — Usage block (line ~16), Step 1's `--harness` paragraph (line 45), Step 4's record template (lines ~100–107), Step 6 Harness scope bullet (line 132).
- `README.md` § How it works — Mermaid nodes `retroHarness` (15) and `readReports` (16), edges at 43–46, `linkStyle` 47–50, the legend sentence and the Report back bullet (57).
- `context/AGENTS.md:41` § Learning loop — the sentence naming `/retro --harness`.
- `skills/feedback/SKILL.md:88` — names the harness queue and `/retro`.
- `.agents/docs/deploy/local.md` § Report back and `.agents/docs/specs/sh-agent/delivery.md` § The loop back / § Where things stand — where the reports go and who reads them.
- `.agents/docs/architecture/DECISIONS.md` — ADR-032 appended.
- `cmd/sh-agent/main_test.go` `TestHarnessMap` — read only; it must still pass (the map keeps `sh-agent install`, `update`, `report`).
- The interactive map page (private claude.ai artifact, outside the repository) — updated from the README map after the commit (AGENTS.md rule 5).

## Interfaces

**Invocation.** `/retro --harness [--dry-run]`, unchanged. The checkout is always `~/.config/sh-agent/feedback` (D13). `--dry-run` keeps its meaning: findings in chat, nothing written — the pull in step 1 below still runs, because it changes only the feedback checkout, which is a copy.

**Step 1, `--harness`, in this order** (replaces the paragraph at `skills/retro/SKILL.md:45`; the local reads it lists stay):

1. `git -C <checkout> pull --ff-only --quiet`. The checkout missing, or the pull exiting non-zero: print `retro: feedback checkout <path> not read: <reason>` and continue with step 5 only.
2. This host: `hostname`, with `/` and `\` replaced by `_` — the same name `sh-agent report` gives its directory (`internal/harness/report.go` `hostDir`). `reports/<this host>/` is skipped: the local queue is newer.
3. For every other `reports/<host>/`, the newest run directory: the greatest name in byte order (`2026-09-17-110749-2` sorts after `2026-09-17-110749`). Older runs are not read.
4. From that run: every `proposals/*.md` proposal whose `**Status:**` is `queued` or absent is a finding, identified by its SOURCE — `<host> · <run> · proposals/<file> · <its "## From …" heading> · <its "### Proposal …" heading>`. Its `feedback.jsonl` and `retros.jsonl` lines are added to the inputs of Step 2's promotion check and nothing else.
5. The local reads, unchanged: `~/.config/sh-agent/retros.jsonl`, every `~/.config/sh-agent/proposals/*.md`, and `REQUIREMENTS.md` edits the newest handoff names.
6. Skip any report finding whose source identity, without the `<run>` part, appears on a `Source:` line in any `<docs-root>/retros/*.md` record of the harness checkout with a status of `applied …` or `rejected …`.

The chat table lists local queue findings first, then report findings grouped by host, oldest `## From` time first, each with a `Host` column.

**Ledger line.** In the retro record, a proposal consumed from a report gets one line under its `Status:` line:
`Source: <host> · proposals/<file> · <## From heading> · <### Proposal heading>`
(no run directory, so a later report carrying the same proposal matches it). Its status follows Step 6 as today: `applied <date> (<commit>)` or `rejected <date> — <reason>`. A local queue proposal keeps today's handling — its status is written in the local queue file — and carries no `Source:` line.

**Step 6, where statuses are written** (exact replacement text for the two bullets at `skills/retro/SKILL.md:131-132`):

- Harness scope bullet, its final clause "then, after that repo's review and verification, update the status in the queue file and in the record when it is on this machine." becomes: "then, after that repo's review and verification, update the status in the record, and in the local queue file for a proposal from this machine's queue. A proposal consumed from another place's report gets its status only in the record, on the proposal that carries its `Source:` line: its queue is a copy inside the feedback checkout, which is never edited, committed or pushed — the record is the ledger (D5, D11)."
- Rejected bullet "annotate `Status: rejected <date> — <reason>` in the record and the queue; never delete a proposal." becomes: "annotate `Status: rejected <date> — <reason>` in the record, and in the local queue for a proposal from this machine's queue — a report proposal's rejection lives only in the record, with its `Source:` line; never delete a proposal."

**A redacted copy.** A proposal whose text carries `[redacted by sh-agent report` cannot be applied verbatim: it is listed with `redacted — re-derive` in the table, and its diff is re-derived from its title and why, or it is rejected with that reason.

**The map** (`README.md`):
- Remove node `readReports` and its two edges (`feedbackRepo -.-> readReports`, `readReports -.-> edit`).
- `retroHarness` label: `/retro --harness<br/>this machine's queue and every place's newest report`.
- New solid edge `feedbackRepo -- "pull, newest run per place" --> retroHarness`; `configDir -. "same machine only" .-> retroHarness` becomes solid, `configDir -- "this machine's queue" --> retroHarness`.
- Edge count goes from 16 to 15; the `linkStyle` lines are renumbered so every report-back edge stays amber and the other colours keep their edges.
- The legend sentence becomes: `A dashed line would mark a gap — done by hand, or limited to one machine; the map has none today.`
- The Report back bullet's last sentence becomes: `Back at the harness, /retro --harness applies this machine's queue and the queued proposals in every other place's newest report, and records each one it consumes so it is not applied twice.`

## Failure modes

| When | What the operator sees |
|---|---|
| Feedback checkout absent (never reported from this machine) | `retro: feedback checkout ~/.config/sh-agent/feedback not read: absent — run sh-agent report once on this machine`; the retro continues with the local queue |
| `git pull --ff-only` fails (offline, diverged) | `retro: feedback checkout <path> not read: <git's first error line>`; continues with the local queue |
| `reports/` has only this host | No report findings; nothing printed beyond the table |
| A host's newest run has no `proposals/` | That host contributes nothing |
| A proposal's copy is redacted | Listed `redacted — re-derive`; never applied verbatim |
| The same proposal is in the local queue and a report | Cannot happen for this host (skipped, D6); for another host it is that host's |
| A report proposal already applied here in an earlier session | Skipped by the ledger (step 6); not listed |

## Verification

Tiers: skill text only, so no unit tier. The harness's own tier for `skills/` (project AGENTS.md § Testing): scratch install with no drift, `/review --docs` over the range, and a run of the changed skill.

- Install: `go run ./cmd/sh-agent install --home .agents/scratch/gaps-home` and `status` — drift none; the new `--harness` text present in both vendors' copies of `skills/retro/SKILL.md`.
- Regression: `go test -race -count=1 ./...` passes, `TestHarnessMap` included.
- E2E, real data only, nothing fabricated in the feedback repository or the trail (D13):
  - Now (edge, one host): with only `reports/dave/` in the checkout, `/retro --harness --dry-run` run in this checkout, following the working-tree skill text, pulls, skips `dave` as this host, lists no report findings, and writes nothing (`git status` of this checkout and of the feedback checkout unchanged apart from the pull).
  - Golden path (who: both): the operator runs `/retro` at another place until it queues a harness-scope proposal there, then `sh-agent report` from that place. Here, `/retro --harness --dry-run` lists that proposal with that place's host in the `Host` column, read from that host's newest run only.
  - Ledger (who: both): the operator applies or rejects that proposal here through `/retro --harness`; its record carries the `Source:` line and status; a following `/retro --harness --dry-run` no longer lists it, and the feedback checkout's `git status` shows no change made by the retro.
- Each criterion is observed in the run's output by the main agent and recorded in `.agents/docs/testing/<date>-close-map-gaps-e2e.md`.

## Phases

- [x] Cold-reader check — three rounds, `.agents/docs/reviews/2026-09-17-coldread-close-map-gaps.md`; never reached both tokens, ended at the cap. Every blocking finding resolved in the plan (D9, D11, D12); the operator resolved the cap by cutting the fixture E2E and its flag, the source of the last two blocking items: "Cut --feedback-repo + fixture E2E" (D13), 2026-09-17

### Phase 1 — skill, docs, map, ADR (one commit, R9)

- [x] Check the hostname assumption (Risks) — verify: `hostname` output against the directory names under `~/.config/sh-agent/feedback/reports/` (read only), who: AI — expect: this box prints `dave` and `reports/dave/` exists; any mismatch stops the phase and becomes a D-row before the skill text is written. Evidence: `hostname` prints `dave`; `~/.config/sh-agent/feedback/reports/` holds `dave/` only, 2026-09-17
- [x] `skills/retro/SKILL.md`: Usage line for `--harness` names the reports; Step 1 `--harness` paragraph replaced by the six-step order in Interfaces; Step 4 template shows the `Source:` line; Step 6's Harness scope and Rejected bullets replaced by the exact text in Interfaces § Step 6 — verify: read against Interfaces; `sh-agent install --home .agents/scratch/gaps-home` then `status`, who: AI — expect: drift none, the text in both vendors' copies — docs: `skills/retro/SKILL.md`. Evidence: Usage, Step 1 (six ordered steps and the table order), Step 4 `Source:` template line and both Step 6 bullets as Interfaces words them; `sh-agent install --home .agents/scratch/gaps-home` exit 0, `status` drift none, "the records are the ledger across machines" in `.claude/skills/retro/SKILL.md` and `.gemini/config/skills/retro/SKILL.md`, 2026-09-17
- [x] `context/AGENTS.md:41` and `skills/feedback/SKILL.md:88`: the sentences naming `/retro --harness` say it also consumes every place's newest report — verify: `grep -n "retro --harness" context/AGENTS.md skills/feedback/SKILL.md`, who: AI — expect: both name the reports. Evidence: context § Learning loop names every other place's newest `sh-agent report`, present in both installed context copies; `skills/feedback/SKILL.md:88` says a harness session finds other places' proposals through their newest report, 2026-09-17
- [ ] `.agents/docs/deploy/local.md` § Report back and `.agents/docs/specs/sh-agent/delivery.md`: who reads the reports — verify: `/review --docs` in the checkpoint, who: AI — expect: no living doc says reports are read by hand — docs: `.agents/docs/deploy/local.md` — in progress: edited by `documentation` (local.md § Report back, delivery.md § The loop back and § Where things stand); ticked at the checkpoint's `/review --docs`; blocker: none
- [x] E2E edge, one host — as Verification states, now — verify: e2e, who: AI — expect: `dave` skipped, no report findings, nothing written; log `.agents/docs/testing/<date>-close-map-gaps-e2e.md`. Evidence: PASS, `.agents/docs/testing/2026-09-17-close-map-gaps-e2e.md` § Edge — pull exit 0, `dave` skipped, no other host, the six local proposals already applied, 0 `Source:` lines, both checkouts' status unchanged, 2026-09-17
- [ ] E2E golden path — as Verification states — verify: e2e, who: both (the operator reports from another place with a queued harness proposal) — expect: the proposal listed with its host; appended to the same log — in progress: not started; blocker: needs another place to queue a harness-scope proposal and run `sh-agent report`
- [ ] E2E ledger — as Verification states — verify: e2e, who: both — expect: `Source:` line and status in the record, the proposal absent from the next dry run, the feedback checkout untouched; appended to the same log
- [ ] `README.md` map and prose per Interfaces; ADR-032 appended — verify: `go test -count=1 -run TestHarnessMap ./cmd/sh-agent/` and a mermaid-cli render to `.agents/scratch/map/map-gaps.png`, read once, who: AI — expect: test passes; the render shows no dashed line, 15 edges, amber report-back edges — docs: `README.md`
- [ ] Operator approves the redrawn map (R9), after the golden path and ledger E2E pass — nothing on the map turns solid before it works for every place (AGENTS.md rule 4) — verify: the operator's words quoted here, who: user — expect: approval in their own turn
- [ ] After the push: the interactive map page (private claude.ai artifact) updated from the committed README map, per AGENTS.md rule 5 — verify: the published page's edges and labels read against the README map, who: AI — expect: no dashed line, the same nodes and edge labels
- [ ] Checkpoint — expect: `gofmt`, `go vet`, `go test -race -count=1 ./...` clean; verifier log; `/review --docs` over the working tree with every finding fixed or rebutted; changelog entry read against the diff; one commit carrying skill, docs, map and ADR-032; pushed — evidence: <log paths, commit>, <date>

## Risks

Assumptions:
- `hostname` prints the same name `os.Hostname()` returns, so step 2 finds this host's directory — checked by comparing `hostname` with the directory `sh-agent report` created on this box (`reports/dave/`), first task of Phase 1; on a Mac, `hostname` can print `name.local` while `os.Hostname()` does too — unverified there, stated in the skill: when no directory matches, print `retro: no reports/<name>/ for this host` and read all hosts.
- An agent following skill text performs byte-order "newest" and identity matching reliably — checked only by the E2E fixture; a mistake there is a missed or duplicated proposal, never a leak.

## Decisions

Scope lines below are the numbered list the operator confirmed in chat on 2026-09-17 with "go all": 1 mechanism, 2 which reports, 3 what counts, 4 freshness, 5 status across machines, 6 this machine, 7 out of scope, 8 map and R9.
D1 — 2026-09-17 — Skill-only change to `skills/retro/SKILL.md`; no Go code, no new verb — by: the operator, "go all" on scope line 1
D2 — 2026-09-17 — Only the newest run per host is read; older runs are stale copies of the same queue — by: the operator, scope line 2
D3 — 2026-09-17 — Other hosts' queued proposals are findings; their `feedback.jsonl` and `retros.jsonl` feed only the promotion check — by: the operator, scope line 3
D4 — 2026-09-17 — `git pull --ff-only` first; missing or failing checkout is reported and the retro continues on the local queue — by: the operator, scope line 4
D5 — 2026-09-17 — The harness repo's retro records are the ledger, via a `Source:` line without the run directory; nothing is written to the feedback repository or another machine — by: the operator, scope line 5
D6 — 2026-09-17 — This host's own reports are skipped; its local queue is authoritative — by: the operator, scope line 6
D7 — 2026-09-17 — `--feedback-repo DIR` on `/retro --harness` for a fixture E2E — superseded by D13 — by: implementer
D8 — 2026-09-17 — Map, ADR-032, skill and docs land in one commit, after the operator approves the redraw — by: the operator, scope line 8, and R9
D9 — 2026-09-17 — superseded by D13: The E2E runs from `harness-copy/`, a scratch copy of this checkout, so the ledger fixture lives in the copy's `.agents/docs/retros/` and the real, immutable trail gets no fake record; no ledger-path option is added — by: cold read round 1, blocking question 1
D11 — 2026-09-17 — A report proposal's status (applied or rejected) is written only in the harness retro record beside its `Source:` line, never into the copy of its queue inside the feedback checkout; exact Step 6 text in Interfaces — by: cold read round 2, blocking question 1, following D5 and the Rejected row on a second write path back
D12 — 2026-09-17 — superseded by D13: The E2E names the fixture checkout by absolute path, because the run starts in `harness-copy/`; a golden-path run that does not list erin's A has not read the fixture and fails. ADR-032 is drafted during Phase 1 from D1–D11 in ADR-031's style — by: cold read round 3, blocking question 1 and non-blocking 2
D13 — 2026-09-17 — The cold read ended at its cap with one narrow item left, in the fixture E2E. Asked how to continue, the operator chose to cut a requirement and picked "Cut --feedback-repo + fixture E2E": no test flag, no fixtures; the E2E runs on real reports — the one-host edge now, the golden path and ledger after another place reports a queued proposal — and the single R9 commit waits for them. Supersedes D7, D9, D12 — by: the operator at the cold-read cap, 2026-09-17
D14 — 2026-09-17 — Until the golden path and ledger E2E pass on another place's report, the finished skill, docs, plan and logs live on branch `close-map-gaps`, pushed; `main` stays untouched, and the one R9 commit to `main` (skill, docs, map, ADR-032) lands after the E2E and the operator's map approval — by: the operator, "Push a branch", 2026-09-17
D10 — 2026-09-17 — "Newest run" is the greatest directory name in byte order; `runDir`'s unpadded suffix mis-orders only past nine same-second runs from one host (`-10` before `-9`), accepted — by: cold read round 1, non-blocking 3
Rejected: a timer, an `sh-agent` verb, a dashboard — none asked for; each adds a surface (scope line 7)
Rejected: marking proposals applied in the feedback repository — a second write path back, against R9's one way back
Rejected: reading every run of every host — each run repeats the whole queue, so older runs only add duplicates
