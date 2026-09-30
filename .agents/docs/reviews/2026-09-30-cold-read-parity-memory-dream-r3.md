---
owner: reviewer
date: 2026-09-30
---
# Review — cold read r3, parity-memory-dream plan   (2026-09-30)

Range reviewed: `.agents/docs/plans/2026-09-30-parity-memory-dream.md` (working-tree version, the uncommitted diff over d6e396a), checked against `REQUIREMENTS.md`, `AGENTS.md`, the recon brief, and the code surfaces the plan names: `cmd/sh-agent/main.go` (usage, `run`, `status`), `internal/search/{store,search,reconcile,parse}.go`, `internal/harness/{report,state}.go`, `context/AGENTS.md:9`, and `skills/{feedback-log,init-project,dream,retro,distill}`, `agents/dreamer.md`, `specs/sh-agent/memory-loop.md`, and README.

## Status: COMPLETE

## Verdict: REQUEST CHANGES — design needs revision / implementation not ready

The worst gap is the prompt cap. The plan caps the prompt at 256 KB and says what to drop only for transcript text. In this repository the trail alone (input 3) written in the last 24 hours is 583,878 bytes across 90 files, and 3.9 MB across 572 files over 14 days (the first-run window). So on the harness's own busiest project, every nightly run goes over the cap before any transcript is added, and the plan does not say what happens then. Two careful implementers would build different things: one truncates the trail, the other fails the project.

The second gap is sequencing. The plan builds the `report` change (D19) and the README map edits (configDir `:30`) in Phase 4, and the map's dream job in Phase 5. That is two commits. R9 (`REQUIREMENTS.md:47`) and AGENTS.md:15 (rule 3) require the operator's approval, an ADR, and ONE commit that carries the code, the docs and the map. Phases 3–4 also use the SQLite driver outside the search index (R4, `REQUIREMENTS.md:19`) before D12 is approved, unless Phase 0 is a hard prerequisite, and the plan does not say whether it is.

## Findings

### Gaps

- **[high / confidence: high]** Plan :135 (Prompt) — the 256 KB cap has a drop order only for transcripts ("the oldest transcript text is dropped first"). Measured here: 24 h of `{bugs,reviews,retros,changes,handoffs}` = 583,878 bytes; 14 days = 3,945,260 bytes (the `find -newermt` counts above). Memory is "sent whole" (:126) on top of that. Decide what happens after all transcript text is gone and the prompt is still over the cap. For example: drop the oldest trail files next and list them in the `inputs` frontmatter, and fail the project only when the memory store alone exceeds the cap. BLOCKING.
- **[high / confidence: medium]** Plan Phases 3–5 against `REQUIREMENTS.md:47` (R9: "one commit carrying the code, the docs and the map", with approval and an ADR) and AGENTS.md:15 (rule 3: "A change to install, update, report or release behaviour whose diff does not touch the map is incomplete").
  - Phase 4 (:261) changes `report`'s bundle (D19) and edits the map's `configDir` node (README :30, inside the map :5-60) under Phase 4's docs clause (:260).
  - Phase 5 (:267) adds the dream job to the map in a separate commit.
  - D13 (map approval) is still "awaiting the operator", and nothing gates Phase 4 on it.
  - Fix: put the D19 `Collect` change, the `configDir` label, the dream job on the map and `TestHarnessMap` in one commit, after D13 is approved, with its ADR.
- **[medium / confidence: high]** Plan :233-234 and :265 — the phase gating is ambiguous.
  - Phase 0's second box (D3/D12/D13 approved) and its Checkpoint come before Phase 1, which suggests they are prerequisites.
  - But Phase 5 alone carries "(after D3/D12 approval)", which suggests Phases 1–4 may run before approval.
  - If they may, Phase 3's `search.AgyConversationText` uses `modernc.org/sqlite` for something other than the search index, which R4 (`REQUIREMENTS.md:19`, "used only by the search index") forbids until D12's sentence is in. Phase 4 changes what `report` sends (see the finding above).
  - State that Phase 0 box 2 blocks Phase 3, or put the gate on each affected box. BLOCKING.
- **[medium / confidence: medium]** Plan :211-212, :216, :247 — the "both" setup is "the real home, vendor all", and E2E 1 must run there in the Phase 2 checkpoint. That needs the working-tree harness (the new `project-memory` skill, rule 7, the resume line) installed into the real home during development. R6 (`REQUIREMENTS.md:27`) makes deploying to the real `$HOME` "the operator's action", and AGENTS.md:65 says "The real home is never touched" and requires the three real-home smoke checks after any real-home install. The plan names no `who: user` on that step and no smoke checks. Mark the real-home install and the "both" runs `who: user`, or run "both" in a sandbox.
- **[medium / confidence: medium]** Plan :126 and :133 against recon `plans/2026-09-30-recon-memory-dream.md:33` — "Changed ... for transcripts, the line or step time", and "the other inputs only the changed parts".
  - Antigravity's `steps` table has no time column (`idx, step_type, status, …, metadata, …, step_payload`). The only times are per conversation, in `conversation_summaries.last_modified_time` / `last_user_input_time`.
  - So agy text cannot be cut at the watermark per step unless the probe finds a time inside a protobuf blob.
  - Decide: send a changed conversation whole (re-sending old text every night it changes), or have the probe prove a per-step time. BLOCKING, because it sets the `AgyConversationText(home, root, since)` contract.
- **[medium / confidence: medium]** Plan :148-159 (Written on success) — harness-target proposals are "written to a second file, `proposals/dream/harness/…`, same shape". The plan does not say whether they also stay in the project file, whose rule is "one section per proposal".
  - If they are in both: each harness proposal has two independent `**Status:**` lines. `/retro` in the project and `/retro --harness` can then set them differently, and project `/retro` (:262) would re-see harness items.
  - If they are in the harness file only: say so, and say what project `/retro` does with a file that has no project-target sections.
  - BLOCKING.
- **[medium / confidence: medium]** Plan :104 (`originSessionId` not written) against `skills/dream/SKILL.md:47`. The promote test counts sessions by `originSessionId`, with "a file without the field counts as one session". That wording is ambiguous once every file lacks the field. It can mean each field-less file is its own session (so every merge of two verb-written files promotes), or that all field-less files together are one session (so a merge never promotes). Also, the second route ("a feedback.jsonl entry from a session other than the memory's `originSessionId`") has nothing to compare against. `skills/dream/SKILL.md` appears in the Surfaces only for "Step 1 and description" (:74), not for :47. Decide the promote test for verb-written memory and add :47 to the Surfaces. BLOCKING, because two implementers of the `/dream` edit would diverge.
- **[medium / confidence: medium]** Plan :92 (Root = `git rev-parse --show-toplevel`) — in a linked worktree (`.claude/worktrees/`, which Step 6 gitignores), this resolves to the worktree. The notes then go to a gitignored store that is deleted with the worktree. Today, Claude's native writes go through the absolute `autoMemoryDirectory` (`skills/init-project/SKILL.md:113`) to the main checkout's store from any worktree. So the verb is a regression for worktree sessions, and the plan does not decide it: show-toplevel, or the main worktree via `git rev-parse --git-common-dir`. BLOCKING.
- **[medium / confidence: medium]** Plan :120 and D20 (:310) — the rollout hint fires only "when `<root>/.claude/settings.local.json` exists and lacks …". A project scaffolded before Step 5b existed, or one where the headless guard refused the write (`skills/init-project/SKILL.md:116`), has no such file. On a Claude-installed machine, native memory is then on (writing to `~/.claude/projects/<slug>/memory`) and the hint stays silent. That silently splits the stores, which is the case D20 exists for. Include "file absent, and the manifest vendor includes claude".
- **[medium / confidence: medium]** Plan :216 (E2E 1) — it has no must-fail control proving the fresh session's "Friday" came through `sh-agent memory show`. The plan's risk (:278) claims only that `autoMemoryEnabled: false` stops native *writes*. If Claude Code still loads `MEMORY.md` natively through `autoMemoryDirectory`, E2E 1 passes on Claude even when the resume step never calls the verb. The same holds for an agy session that simply reads the file. Add a control: the store present with the resume line removed (or the verb's output made distinguishable), and read the transcript for the `sh-agent memory show` call.
- **[medium / confidence: medium]** Plan :159 and D19 against `REQUIREMENTS.md:31` (R7: the memory store is "never sent by `sh-agent report`").
  - Harness-target proposals carry an `evidence` field the model writes from memory and transcript text.
  - D19 ships them through `report`, and redaction covers credential shapes only (:311, "Open … default: none").
  - A proposal that quotes a memory note verbatim sends that memory's content to the feedback repository.
  - Either constrain harness-target `evidence` in the prompt to trail `file:line` references and paraphrase, with no verbatim memory or transcript text, or record an operator decision that such quotes are acceptable.
- **[low / confidence: high]** Plan :224 and :266 (`systemd-analyze verify` on the rendered files in scratch) — measured on this box (systemd 257):
  - With the service's `ExecStart` at a scratch-home BinPath that holds no file, `systemd-analyze verify --user` exits 1 with "Command … is not executable: No such file or directory".
  - With an executable placed there, it exits 0.
  - Files rendered by a unit test into a temp home have no binary, so the clause fails for a fixture reason, not a unit defect. The verify step needs an executable at `<scratch>/.local/bin/sh-agent`, for example after `install --home <scratch>`.
  - Scratch probe: `.agents/scratch/r3-sysd/`.
- **[low / confidence: medium]** Plan :125 and :163 — the watermark is stored as `<RFC3339>` (whole seconds), and "Changed" compares file mtime (nanoseconds) against it. A fixture written in the same second as, but before, run 1's `t0` has an mtime above the truncated watermark, so run 2 reports it changed. E2E 3's "second run … logs `unchanged` twice" is then flaky when the script writes fixtures and runs immediately. Store RFC3339Nano, or compare at second precision in both places.
- **[low / confidence: high]** Plan Phase 5 (:266) has no docs clause for `cmd/sh-agent/main.go:47` (usage: `uninstall [--dry-run]   remove exactly the files install wrote, nothing else`) or README:180 (the uninstall row), both of which become false once `uninstall` removes the timer. The Surfaces name "`uninstall` row" but the Phase 5 box lists only "README.md uninstall row". The usage lines for `--install-timer`/`--remove-timer` and `specs/sh-agent/dream.md` (timer section) are also missing from that box.
- **[low / confidence: medium]** `.agents/docs/deploy/remote.md:66` describes what `report` sends ("feedback log, retro proposals and index … never the project's own auto memory store"). D19 changes that set, and remote.md is not in the plan's Surfaces.

### Scope (against `## Why`)

- Every addition I found has a D-row:
  - `--project` and `--dry-run`: D15.
  - the rollout hint: D20.
  - `proposals/dream/harness/`: D19.
  - no-tools CLI: D4.
  - opt-in timer: D5.
  - `project-memory` naming: D14.
- **[low / confidence: medium]** `autoMemoryDirectory` is kept ("as today", :194) and folded into D13's Non-goals exception, although with `autoMemoryEnabled: false` Claude Code writes nothing to it. No D-row says why it stays. Plausible reasons: the Step 5b migration gate at `skills/init-project/SKILL.md:115`, and keeping an operator's deliberate `true` pointed at the right store. Record which one, or drop the key and ask for one Non-goals key instead of two.
- `dream-failed/<date>-<project>.txt`: this is a failure artifact and not scope creep. But the plan does not say whether its raw text is redacted (redaction is stated only under "Written on success", :161), and it has no retention rule (see the non-blocking questions).

### Readiness — questions I would have to ask

BLOCKING (two careful implementers would build different things):
1. What is dropped, or what fails, when the prompt is still over 256 KB after all transcript text is gone? (Finding 1.)
2. May Phases 1–4 start before D3/D12/D13 are approved? In particular, may Phase 3's driver use and Phase 4's `report`/map change land before then? (Findings 2–3.)
3. Do harness-target proposals appear in both files or only in `harness/`? If both, which `Status:` is authoritative?
4. `AgyConversationText(home, root, since)`: is `since` applied per conversation or per step, when steps carry no time column?
5. How does the `/dream` promote test treat verb-written memories without `originSessionId` (`skills/dream/SKILL.md:47`)?
6. In a linked worktree, which root does `memory add|show` use?

NON-BLOCKING (the plan states a default, or a reasonable implementer decides alone):
7. Should the rollout hint also fire when `settings.local.json` is absent on a Claude-installed machine? (Recommend: yes.)
8. Step 5b's existing "empty object → overwrite with the one-key object" branch (`skills/init-project/SKILL.md:115`): how does it combine with two keys? (Default: overwrite with the two-key object.)
9. `--project <dir>`: must the root be registered in `search.db`, and is a relative path made absolute or resolved to the git toplevel? (Default: absolute toplevel, registration not required.)
10. What happens to the child CLI's stderr? Discarded, or kept beside `dream-failed/`? Verification per AGENTS.md:69 needs agy's stderr ("auto-denied").
11. Is `dream-failed/` text redacted before it is written, and is it ever pruned?
12. Drafts directory `<date>-<project>-<name>/` on a same-day rerun: overwrite, or `-2` as the proposal file does?
13. Does the Claude transcript reader honour `CLAUDE_CONFIG_DIR`? It matters only for the claude-only sandbox.
14. `ListProjects` with a missing db: open read-only, so the file is not created (`:123` says "A missing db means an empty list").
15. Two projects with the same basename collide on `<date>-<project name>.md`. The plan's `-2` rule handles it, but the `project` frontmatter must then carry the root, not the name.

## Against REQUIREMENTS.md and AGENTS.md § Testing

- R3 (`:15`), R4 (`:19`), and Non-goals (`:64`): the departures wait on D3/D12/D13, and the design follows their stated text. The gating ambiguity (Finding 3) is the one place where the design depends on something those rows leave undecided.
- R6 (`:27`): the "both" setup on the real home, with no `who: user` (Finding 4).
- R7 (`:31`): memory content can leave through harness-proposal `evidence` (Finding 11).
- R9 (`:47`): map and `report` change split across commits and phases (Finding 2).
- R10 (`:51`): the design routes both CLIs through `sh-agent memory`. The remaining asymmetries are acknowledged (D16, D17). One more: agy transcripts cannot be cut at the watermark per step while Claude's can (Finding 5), so the two CLIs' inputs differ in granularity. That belongs in the R10 known-difference list if the probe confirms it.
- AGENTS.md:65:
  - § 3d for `internal/search` is present (:255);
  - the `agents/*.md` behavioural run, n of N, is present (:245);
  - `/review --docs` is present (:272);
  - a scratch install plus one sandbox start of each CLI for the skill and context changes is covered only by reference ("sandbox rules per AGENTS.md § Testing", :247);
  - the real-home smoke checks after a real-home install are missing (Finding 4).
- AGENTS.md:72-73: the Phase 4 checkpoint (E2E 3, 4, 6) does not cite § Testing, unlike Phase 2. Fixture projects must live outside every checkout. Any E2E session that runs the context's `sh-agent search --open` needs the `--home` wrapper, or it registers the fixture in the real global index (BUG-107), and the real nightly dream job then processes the fixture.
- AGENTS.md:77: an Antigravity run of each changed skill and agent before the tag is present (:273).

## Criteria
(The dispatch is a plan cold read. These six criteria are applied to the design the plan describes.)
1. Dead code: not applicable (no code yet). `autoMemoryDirectory` may become inert configuration (Scope).
2. Flow tracing: fail. The cap has no path after transcripts are exhausted; agy `since` has no step time; the harness-proposal copies diverge; the `/dream` promote test depends on a field the verb no longer writes.
3. Standards conformance: fail on R9 / AGENTS.md rule 3 (split commit), and on the unlisted restatements (`main.go:47`, `deploy/remote.md:66`).
4. Directory layout: pass. `internal/memory` and `internal/dream` are new packages; the driver stays in `internal/search` (AGENTS.md:59), and the Layout update is planned.
5. Simplicity (DIY-except-security): pass. Stdlib plus the existing driver; no crypto; the redaction is reused, not re-implemented.
6. Decoupling: pass. `uninstall` calls `dream.RemoveTimer` so `internal/harness` does not import `internal/dream`, and the agy reader stays behind `internal/search`.

## Not examined
- `internal/harness/apply.go:88-109` (the manifest's stale deletion, which the plan cites for rejecting the manifest option): not read. The claim is taken from the plan and recon.
- The current byte count of `TestRenderedContextFitsAntigravityBlock`: not run. The plan's 21998/22000 is not re-measured.
- `skills/retro/SKILL.md:118`, the `LastRetro` interplay beyond `state.go:137-157`, and the bodies of `main_test.go` / `report_test.go`, apart from the `Collect` skipped-list assertion at `report_test.go:293`, which will change with D19.
- The recon's claims about agy's schema and the `claude -p` flags were not re-verified; Findings 5 and 10 rest on recon :33-34.
- The cold-read r1/r2 reports were not read beyond r2's verdict line, so this is an independent read.
- Claude Code's documented semantics of `autoMemoryEnabled: false` (whether native *loading* also stops) were not checked. Finding 10 is the consequence either way.
