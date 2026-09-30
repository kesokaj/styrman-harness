---
owner: reviewer
date: 2026-09-30
---
# Review — cold read r2, parity-memory-dream plan   (2026-09-30)

Range reviewed: the working-tree file `.agents/docs/plans/2026-09-30-parity-memory-dream.md` (uncommitted edits over d6e396a), and the code surfaces it names

## Status: COMPLETE

## Verdict: REQUEST CHANGES — design needs revision / implementation not ready

r2 fixes most of r1: the manifest and the timer, PATH, BinPath, the dream lock, the watermark value, `check-ignore`, the skill name, the D-rows for cost and disclosure. The worst problem left is the Antigravity half of the dream job. Its CLI runs in `~/.config/sh-agent/dream-work/`, which is not in the real home's agy `trustedWorkspaces`. The Phase 4 probe runs in a sandbox whose settings file the harness writes for the run (AGENTS.md:69), so it cannot show whether the real run works. If trust is needed, the only fix is a vendor-setting edit, which Non-goals forbid. Next:
- The timer verbs' real `systemctl`/`launchctl` paths have no test that stays off the real user manager (R6).
- Three project rules in AGENTS.md § Testing and § Release are not scheduled: the behavioural run for `agents/dreamer.md`, the § 3d search measurement, and the Antigravity run of each changed skill.
- Five semantic questions change what gets built: a zero-proposal result, whether a skills-only change counts as "changed", how `--home` relates to the child CLI's environment, rollout to projects already scaffolded, and harness-target proposals on places without a checkout.

## Findings

### Gaps (question 1)

- **[high / confidence: medium]** `plan:137-138`, Phase 4 probe `plan:250`, D7 `plan:281`. The agy run's cwd is `dream-work/`. The real home's `~/.gemini/antigravity-cli/settings.json` has `trustedWorkspaces` listing real checkouts only (read on this box). The probe runs "in an agy-only sandbox home", and AGENTS.md:69 says such a home gets a settings file "written for the run … the scratch clone as its only trusted workspace". Whether the probe's file lists `dream-work/` is not stated. If it does, the probe proves a condition the real home does not have. If trust is required, the only fix is adding `dream-work` to the real `trustedWorkspaces`, a vendor-setting edit, which `REQUIREMENTS.md:64` (Non-goals) forbids. D7's fallback covers "cannot answer without a tool", not "refused as untrusted". Fix: the probe's settings file must leave `dream-work/` untrusted, so it matches the real home; decide now what happens when agy refuses; and name that outcome as an R10 known difference in D7.
- **[medium / confidence: high]** `plan:168-181`, Phase 5 verify `plan:257`, vs R6 (`REQUIREMENTS.md:27`). `--home` is "accepted by every command" (`cmd/sh-agent/main.go:61`). The plan does not say whether `--install-timer`, `--remove-timer` and the `uninstall` removal honour `--home` for the unit-file path. A test that runs them against a scratch home would drive the real user manager: `systemctl --user disable --now sh-agent-dream.timer` names the same unit as the operator's real timer. Only the rendered text and `uninstall --dry-run` are tested. The non-dry-run branches have no test: keeping a hand-edited file, "absent is fine", `systemctl` failing (exit 2 naming the files left behind), and uninstall with no manifest. Fix: resolve `systemctl`/`launchctl` on PATH so a test can put a fake first (the Phase 4 fake-CLI pattern). State that the unit directory follows `--home`. Test each branch against the fake.
- **[medium / confidence: high]** AGENTS.md:65, "A change to `agents/*.md` also gets a behavioural run: `claude -p --agents <json …>` … several runs … reported as n of N". The plan edits `agents/dreamer.md:18,28,32` (`plan:43,75`), but no Phase 2 task or checkpoint schedules that run. AGENTS.md:77 also makes the release NO-SHIP until "one run of each changed skill and agent on Antigravity" exists. That covers dreamer, feedback-log, init-project, dream, retro, distill and the new project-memory. Phase 6 names only `/preflight`. Fix: add both runs as tasks, with their fixtures, in Phase 2 and Phase 4.
- **[medium / confidence: high]** AGENTS.md:65: "A change under `internal/search` … also runs `.agents/docs/testing/guide.md` § 3d (the recall measurement, a rebuild under 10 s, a warm search under 100 ms) and adds its rows to `testing/sh-agent-matrix.md`." Phase 3 adds `ListProjects` and `AgyConversationText` to `internal/search` (`plan:242-243`), but its checkpoint does not name § 3d.
- **[medium / confidence: high]** `plan:21` vs D4 (`plan:278`) and D15 (`plan:295`). The plan says both tables are "quoted as the operator approved them", but dreamer item 5 at `plan:21` contains neither "overridable with `--cli`" (cited by D15) nor "Read/Grep/Write with a dedicated agy settings file" (cited by D4). Either the table is not verbatim, or two D-rows cite approved text that does not exist. `--cli` stands on D15 alone (see Scope). Fix: quote item 5 whole, or drop the citations.
- **[medium / confidence: high]** Harness-target dream proposals on a place without a checkout (r1 blocking question 4, second half, still open). `report` does not send `proposals/dream/`: D6, and `report.go:331-339` reads only top-level `.md`. `/retro --harness` runs only in the harness checkout. So on a binary-only or zip place (R8, `REQUIREMENTS.md:35`), a `target: harness` proposal never reaches the harness. R8 says `report` sends "this machine's corrections, proposals". Decide one: `report` sends the harness-target dream files, or the plan and README:214 state that dream proposals stay local and the dreamer's harness output is only useful on the checkout machine.
- **[medium / confidence: high]** `plan:125` ("Changed means any input below newer than the watermark") includes input 4, the installed skills, which are global rather than per project. `install` rewrites only changed files (`apply.go:62-80`). So an `update` that changes any skill makes every indexed project "changed" and triggers one CLI call per project, even with no project activity. The plan does not say whether input 4 counts toward "changed". Two implementers would build different skip logic, and the nightly cost (dreamer item 6) depends on it.
- **[medium / confidence: high]** `plan:145-159`, E2E 3 `plan:212`. When the model returns `{"proposals":[],"skills":[]}`, the plan does not say whether a proposal file is written or whether the watermark advances. The watermark advances "only after the proposal file is written" (`plan:159`), so if no file is written, the project is re-sent in full every night.
- **[medium / confidence: medium]** E2E 2 (`plan:211`): the must-fail control may not produce its condition. In the claude setup the harness is installed with `--vendor claude`, which means the new `project-memory` skill plus feedback-log rule 7 ("the verb on every CLI"). Those tell the model to call `sh-agent memory add`, so a native write carrying `originSessionId` depends on the model choosing Claude Code's own writer. The key is also honoured only in a trusted workspace (`skills/init-project/SKILL.md:115`), which a headless `-p` session may not count as. The plan does not say what happens when the control produces no native file. Fix: run the control with the pre-change harness (or no harness) installed, and say that a control without a native file makes E2E 2 "not exercised", not passed.
- **[medium / confidence: high]** Rollout to projects already scaffolded. `autoMemoryEnabled: false` reaches a project only when `/init-project` runs there again (`plan:188-193`, D11). Until then, Claude Code keeps writing natively in every existing project, beside the verb. That breaks memory item 1 ("both CLIs write the same four types by the same rule") on the "both" setup. No task, handoff line or `status` hint tells the operator to run `/init-project` again in each project.
- **[medium / confidence: medium]** AGENTS.md:13 (rule 1: "…no timer as the mechanism (R8, R9…)"). D13 draws a nightly timer job on the README map. The rule means the delivery and report paths, but a reader checking the map against rule 1 will read the new timer as a violation. The plan's AGENTS.md edits cover Layout and :59 only (`plan:85,263`). Fix: add rule 1 to the grep of restatements, and say in the same commit that the dream timer is neither a delivery path nor a report path.
- **[medium / confidence: medium]** Which home E2E 3, 4 and 6 run in, and what `--home` does to the child CLI. `dream --all` reads `<home>/.config/sh-agent/search.db`, so in the "both" setup (the real home) E2E 3 would dream over every real project and write real proposals and watermarks. The plan does not say whether the child `claude`/`agy` gets `HOME=<--home>`. A scratch HOME cannot authenticate (AGENTS.md:65). The agy-only sandbox needs the sandbox HOME for its token. Fix: E2E 3/4/6 name their home, and the Interfaces section states the child's environment (inherited HOME, or HOME and CLAUDE_CONFIG_DIR set).
- **[low / confidence: medium]** `plan:145` vs `plan:138`: the raw text of a failed run is kept in `dream-work/<date>-<project>.txt`, but `dream-work/` is described as the CLI's empty cwd. If the D7 probe ends with an agy mode that still has read tools, one project's run can read another project's kept raw output. Fix: keep failures in a sibling directory.
- **[low / confidence: medium]** `plan:243`: D8's probe runs on "a copy of one real conversation db in scratch". The real run opens the live `conversation_summaries.db`, which is in WAL mode (`-wal` and `-shm` present on this box) and which agy may be writing. The copy does not reproduce that condition, and the open mode is unstated: `mode=ro` or copying first (never `immutable=1` on a live file).
- **[low / confidence: high]** Docs clauses: no task names the README verbs table, or the `main.go` usage const at `:39`, for `dream` and its flags. Phase 1 names the verbs table for `memory` only, and Phase 5 the uninstall row only. The claude-only kit (CLAUDE_CONFIG_DIR plus a copied credential) is a new testing procedure, and no task writes it into `testing/guide.md`. Also, `deploy/local.md` should state two limits:
  - systemd user timers do not fire without a running user manager (no linger); `Persistent=true` catches up at the next login;
  - rerunning `--install-timer` from a shell with another PATH renders different text, so the installed unit is "kept as hand-edited", and the operator must `--remove-timer` first.

### Scope (question 2)
- `--cli claude|agy`: it rests on D15's quote of dreamer item 5, which is not in the approved table at `plan:21` (see Gaps). It is also a second shape for a value the global `--vendor` flag already carries (`main.go:62`: "default: from manifest"). The plan does not say what `dream --vendor agy` does. This is a gap until a D-row quotes the operator.
- `--dry-run` on `dream`: no D-row names who asked for it. It is useful, and the global read-only-mode rule shapes it, but by the dispatch's rule it needs a D-row.
- `--project <dir>`: D15, by the implementer, for the e2e and for reruns. A D-row names who asked, so it is accepted.
- Everything else traces to the Why tables or to R10: `memory add|show`, the timer and its removal, drafts, the `status` line, the per-project cap.

### Readiness (question 3)
BLOCKING (two careful implementers would build different things):
1. When agy refuses `dream-work/` as an untrusted workspace, what happens: a known difference, another cwd, or a request to change Non-goals?
2. Do the timer verbs follow `--home`, and how are the `systemctl`/`launchctl` branches tested without the real user manager (R6)?
3. Does a change to input 4, the installed skills, make a project "changed"?
4. On a zero-proposal result, is a file written, and does the watermark advance?
5. What environment does the child CLI get under `--home`, and in which home do E2E 3, 4 and 6 run?
6. Do harness-target dream proposals leave a non-checkout place (through `report`), or are they local only (R8)?

NON-BLOCKING (a default is stated, or a reasonable implementer decides it alone):
- Updating a memory finds the existing `MEMORY.md` line by its link target `(<slug>.md)`. Native lines use a title, not the slug (`.agents/memory/MEMORY.md` on this box).
- A newline in `--description` breaks the index line: refuse it or fold it.
- Mode of the memory files: native files here are 0664 inside a 0700 store.
- Stale-lock removal race in `memory add`: two processes each remove the other's fresh lock. The practical risk is low.
- How `status` decides "timer installed": the file exists, or `systemctl --user is-enabled`.
- `claude -p` in `dream-work/` loads the 22 KB global context on every call; `--setting-sources` could trim it.
- The glob `~/.claude/projects/*/*.jsonl` excludes subagent transcripts; accept that.
- D16 names only headless runs, but interactive Claude Code also prompts for `sh-agent memory show` at every session start. `~/.claude/settings.json` on this box has no `permissions` key, while agy's `toolPermission` approves. So the two CLIs behave differently on this machine. Record it beside D16.
- AgyConversationText open mode (`mode=ro` over the WAL file, or copy first).

## Requirements and project-rule check
- R3 (`REQUIREMENTS.md:15`), R4 (`:19`), Non-goals (`:64`): each departure waits on D3, D12 or D13, and the design follows their stated text. The agy trust question may need another Non-goals exception, which none of the three covers (Gaps, first finding).
- R6 (`:27`): the timer tests; E2E 3 in the real home (Gaps).
- R7 (`:31`): the dream records question is settled by D10; the store stays `.agents/memory`. Pass.
- R8 (`:35`): "sends … proposals" vs D6 for harness-target dream proposals (Gaps).
- R9 (`:47`) and AGENTS.md:15: the map change is D13, in one commit with an ADR. Pass, pending approval.
- R10 (`:51`): the agy trust question, the interactive prompt difference, and existing projects keeping native writes (Gaps).
- AGENTS.md § Testing: `:65` (behavioural run for the agent, § 3d, `/review --docs`, which Phase 6 has), `:69` (the agy sandbox settings vs the probe), `:72` (fixtures outside checkouts: the Phase 2 checkpoint cites the § Testing sandbox rules, pass), `:73` (the `--home` wrapper: the memory verbs do not use home, pass), and `:77` (the Antigravity run of each changed skill and agent before the tag).

## Criteria
1. Dead code — n/a for a plan. The one no-op task, the `Collect` pin, is now honestly labelled a pin (`plan:252`). Pass.
2. Flow tracing — the named line references hold: main.go:39 and :508, report.go:199 and :331-339, store.go:41, apply.go:88-109, `BinPath` at plan.go:16, context 21998/22000. Four flows are open: the child CLI's environment, the timer's command paths, harness proposals leaving the machine, and rollout to existing projects.
3. Standards conformance — fails: three AGENTS.md § Testing and § Release duties are unscheduled, and AGENTS.md:13 is not reconciled.
4. Directory layout — `internal/memory` and `internal/dream`, with the Agy reader in `internal/search`, fit AGENTS.md:43,59 and the D12 text. Pass.
5. Simplicity (DIY-except-security) — stdlib plus the existing driver. Redaction reuses `redactLines`. No crypto. `--cli` duplicates `--vendor`. Otherwise pass.
6. Decoupling — `uninstall` calls `dream.RemoveTimer` from `main.go`, so `internal/harness` stays free of `internal/dream`. Pass.

`design needs revision` — `implementation not ready`

## Not examined
- `specs/sh-agent/memory-loop.md` beyond the grep hits; the recon brief `plans/2026-09-30-recon-memory-dream.md`; `DECISIONS.md` ADR-047.
- The README map (:5-60) and `TestHarnessMap`; `skills/distill/SKILL.md:38` beyond one grep line.
- The scratch file `.agents/scratch/2026-09-30-parity-operator-items.md` with the operator items D3, D12 and D13 cite.
- Nothing was run against a live CLI. The agy trust behaviour for an untrusted cwd is inferred from the settings key and AGENTS.md:69, not observed.
