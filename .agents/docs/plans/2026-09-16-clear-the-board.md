---
owner: main
date: 2026-09-16
---
# Plan — clear the board: every open bug and Follow-up closed before styrman-crew's Session B   (2026-09-16)

## Why

The operator asked, before starting styrman-crew's Session B and argolis-devbox's Session C: "Fix everything before I do b and c. Nothing left, all tests, docs nothing in todo and roadmap. This project clean and fully documented and up to date." Six bug records stand open and eleven Follow-ups are filed across two plans and the newest handoff. One of them, BUG-046, is the mechanism Session C is about to install on a VM.

## Scope

In: every `status: open` record in `bugs/INDEX.md` (BUG-019, 020, 021, 022, 031, 046); every Follow-up in `plans/2026-09-16-attribution-write.md` and `plans/2026-09-16-linkshelf-simulation.md` not already applied; the four harness tasks the newest handoff lists; the docs each change touches; the full Definition of Done over the batch.
Out: cutting `v1.0.0` (the operator's own turn); the Mac; the two mirrors; anything in styrman-crew or argolis-devbox.
Already applied earlier today, not repeated here: G4 and G11 (retro Proposals 7 and 6), and three of the five attribution Follow-ups (fixed temp name, single size check, read-only refusal — e75e7f4, f85169c).

## Surfaces touched

`cmd/sh-agent/main.go`, `internal/harness/{apply.go,adopt.go,render.go,harness_test.go}`, `cmd/sh-agent/main_test.go`; `context/AGENTS.md`; `agents/verifier.md`; `skills/{bash,cloud-run-deploy,distill,docker,frontend,gcloud-setup,golang,init-project,preflight,review}/`; `README.md`, `AGENTS.md`, `.agents/docs/deploy/remote.md`.

## Interfaces

- `sh-agent update` gains a refusal: a dirty checkout is not installed on any path. Exit status non-zero, message names the checkout and what to do.
- No other command, flag or file layout changes.

## Failure modes

- A prose fix contradicts another restatement of the same rule elsewhere (BUG-003/007/019's own shape): every text change is followed by a subject-noun grep across `context/`, `skills/`, `agents/`, `README.md`, `AGENTS.md`.
- The `update` refusal breaks the operator's hourly timer's legitimate use (a clean checkout behind upstream): the test matrix covers clean-behind, clean-current, dirty-behind, dirty-current.
- A skill edit lands in a vendored skill: checked against `skills/VENDORED.md` first — none of the targets is vendored.

## Verification

`gofmt -l cmd internal` (not `gofmt -l .` — BUG-048, found while running this plan), `go vet ./...`, `go test -race -count=1 ./...`; `go run ./cmd/sh-agent install --home .agents/scratch/board/home` with `status` reporting no drift and `list` showing 27 skills and 7 agents; one headless start of each CLI from the sandbox home with its log grepped; the real-home smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI.

## Phases

### Phase 1 — Go: the guards, the writes, their tests

- [x] BUG-046: `update` refuses a dirty checkout on every path — `cmd/sh-agent/main.go` `update()`; `TestUpdateNeverInstallsDirtyTree` covers all four states plus `--check`; validated in the must-fail state (guard removed → the dirty-current leg fails), 2026-09-16
- [x] `writeAtomic` syncs before the rename — `TestWriteAtomicIsAtomicAndDurable` asserts content, mode and that no temp file is left on either outcome, 2026-09-16
- [x] `abandonTemp` and `readSettings` error paths tested — `TestAbandonTempReportsRemoveFailure` (a read-only parent forces the remove failure), `TestReadSettingsProblems` (missing, directory, unreadable; asserts no bytes escape), both skipped under root, 2026-09-16
- [x] `Adopt` stages the copy and renames it into place — `TestAdoptFailedCopyKeepsTheExistingSkill`; validated in the must-fail state (staging removed → the test fails), 2026-09-16
- [x] `RenderContext` symlink refusal tested directly — `TestRenderContextRefusesSymlinks`, both a linked `context/` and a linked file inside it, 2026-09-16
- [x] `TestNoTmpInHarnessProse` walks whole non-vendored skill directories, `commands/` and `rules/` — vendored skills exempt through a new `vendoredSkills` parser of `VENDORED.md`; it caught four real drifts (three in `golang-ecosystem.md`, BUG-019e, one in `bash-fundamentals.md`) and passes after fixing them; terraform's `/tmp` lines are vendored and correctly skipped, 2026-09-16
- [x] BUG-048, found in passing: the documented `gofmt -l .` walks `.agents/scratch` (the Go tool skips dot-directories, `gofmt` does not), so one stray scratch Go file fails the project's own gate — scoped to `gofmt -l cmd internal` in `AGENTS.md`, `README.md` and `deploy/local.md`, and named as a default caveat in `agents/verifier.md`, 2026-09-16
- [x] Checkpoint — `gofmt -l cmd internal`, `go vet ./...` and `go test -race -count=1 ./...` all clean; `changes/` entry, commit, push, 2026-09-16

### Phase 2 — Instruction text: the six open records and the rest of the Follow-ups

- [x] BUG-019 a: `skills/bash/references/makefile-patterns.md` — `REGION` is required of the operator, `REGISTRY` derives from it, `IMAGE_TAG` is the git short SHA, both deploy targets take `--region=$(REGION)`, and a `docker-build-amd64` target names what actually ships to Google; `skills/gcloud-setup/SKILL.md` region/zone likewise, 2026-09-16
- [x] BUG-019 b: `skills/docker/SKILL.md` — the daemon-start row states the standing-permission condition and otherwise hands the command back; the smoke test captures the container id and probes that container, with `docker logs`/`stop`, 2026-09-16
- [x] BUG-019 c: the three CI sections say the harness runs no CI and name the local form first (`bash/references/testing-linting.md` ×2, `gcloud-setup/SKILL.md` service-account keys reframed as a last resort, `frontend/references/accessibility.md` axe), 2026-09-16
- [x] BUG-019 d: Cobra+Viper, `samber/lo`, `samber/oops` and `goleak` (in both `golang-testing.md` and `golang-core.md`) are stated as dependencies to justify against a named stdlib alternative; the skill's own entry points already said so, which is why only these rule-lines had drifted, 2026-09-16
- [x] BUG-019 e: benchmark baselines go to `.agents/scratch/` — fixed in Phase 1 with the test that found them, 2026-09-16
- [x] BUG-020: `skills/cloud-run-deploy/SKILL.md` description now says ship/deploy/release in the user's own turn means `/preflight`, that this skill is the HOW it calls at its deploy row, and adds a Do NOT line for answering a ship request on its own, 2026-09-16
- [x] BUG-021: `skills/distill/SKILL.md` points at the `feedback` skill as the owner of the log's shape, says to include every field, and its example carries `project` with the reason (`/retro`'s exclusion filter) spelled out, 2026-09-16
- [x] BUG-022: `golang-ecosystem.md` scopes registration to a public library, forbids it for private or internal code, and says it is never done on the agent's own initiative, 2026-09-16
- [x] BUG-031: `context/AGENTS.md` — a headless or print-mode run ends its turn only after every dispatched agent has returned and its artifact is on disk, and names any that has not, 2026-09-16
- [x] G5: `skills/init-project/SKILL.md` says what filling a section means and what an empty one becomes (one line saying so and what would go there), 2026-09-16
- [x] G10: `skills/review/SKILL.md` re-dispatch brief now asks the agent to re-read each declaration a fix touched and every reference to it, not only the finding's lines, 2026-09-16
- [x] G12: `context/AGENTS.md` scratch rule — a service's runtime data path comes from configuration with a default outside the repository, or is gitignored in the same change; a plain `go run` never dirties the tree, 2026-09-16
- [x] G13: `agents/documentation.md` bounds the Setup section to the commands a contributor runs, tiers to the testing guide; the same bound added to `skills/init-project/assets/AGENTS.md`, 2026-09-16
- [x] G20: `agents/verifier.md` defines `<app>` in resolution order and forbids leaving it literal; `agents/documentation.md` restates the same resolution, 2026-09-16
- [x] `skills/preflight/SKILL.md` Step 4 already files what it defers ("file it under the active plan's Follow-ups, or the handoff's Next Steps when no plan is active") — applied in an earlier session, verified present, nothing to change, 2026-09-16
- [x] Restatement greps run over `context/`, `skills/`, `agents/`, `README.md`, `AGENTS.md` for each changed rule — two gaps found and closed (`<app>` resolution unstated for `documentation`, the Setup bound missing from the init-project template); the cloud-run triggers, the feedback-log shape and the print-mode wording had no other restatement, 2026-09-16
- [x] Checkpoint — Phase 2's text landed in 3c46356; its `changes/` entry is `changes/2026-09-16-clear-the-board-text.md`, written with the Phase 4 records (the docs review caught the tick claiming an entry that did not exist yet), 2026-09-16

### Phase 3 — Docs

- [x] `.agents/docs/deploy/remote.md` § The resync loop now states the refusal, why it exists, and shows the timer unit it makes safe; `README.md`'s command table and loop line agree, 2026-09-16
- [x] `README.md` and `AGENTS.md` updated (the `update` row, the loop line, the scoped `gofmt`), 2026-09-16
- [x] `bugs/INDEX.md`: BUG-019, 020, 021, 022, 031 and 046 closed and moved to Graduated with their new homes, and BUG-048, 049, 050 written and graduated — main table now holds zero `open` rows, verified by reading it, and the Graduated table re-sorted by ID, 2026-09-16
- [x] Checkpoint — Phase 3's docs landed in 3c46356 with the same `changes/` entry, 2026-09-16

### Phase 4 — Definition of Done over the batch

- [x] `verifier` over `e08fa65..HEAD` — PASS, `testing/2026-09-16-verify-clear-the-board.md`: pass 92, fail 0, skip 1 (from 79); install into a scratch home drift none; `sh-agent list` 27/7, 2026-09-16
- [x] `reviewer` code lens (APPROVE WITH NITS, three findings) and docs lens (REQUEST CHANGES, four findings) — `reviews/2026-09-16-review-clear-the-board-{code,docs}.md`, 2026-09-16
- [x] `security` over the whole range — `reviews/2026-09-16-security-clear-the-board.md`: two medium (the dirty check inheriting `status.showUntrackedFiles`; secret-shaped files deployed from a skill directory), one low TOCTOU, four further lows, 2026-09-16
- [x] Every finding fixed, rebutted or left with an owner in each report's `## Outcome` — two mediums fixed with tests and their own bug records (BUG-049, BUG-050), the TOCTOU accepted into ADR-027's Consequences, the root-step-gate finding rebutted against R3, govulncheck left to the operator, 2026-09-16
- [ ] Sandbox: install, one headless start of each CLI, log grepped
- [ ] Real home: install, the three smoke checks on each CLI
- [ ] `changelog`, `bugs`, commit, push, `/handoff`
- [ ] TODO.md and ROADMAP.md empty; `bugs/INDEX.md` main table empty of `open` rows

## Risks

- The batch is wide and shallow: the real risk is a text fix that contradicts a neighbour, which the subject-noun grep and the docs lens are there to catch.
- `update`'s new refusal changes behaviour the operator's timer depends on; the docs and the message must say what to do, and the timer keeps working on a clean checkout.

## Decisions

- D1. Every open record is either fixed here or, if a fix would be a design question, recorded as a decision with the operator rather than left open. Nothing stays `open` silently.
- D2. `update` refuses rather than warns: a timer that reinstalls a dirty tree defeats R6, and a warning nobody reads is what produced BUG-046.
