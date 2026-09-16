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

`gofmt -l .`, `go vet ./...`, `go test -race -count=1 ./...`; `go run ./cmd/sh-agent install --home .agents/scratch/board/home` with `status` reporting no drift and `list` showing 27 skills and 7 agents; one headless start of each CLI from the sandbox home with its log grepped; the real-home smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI.

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

- [ ] BUG-019 a: `skills/bash/references/makefile-patterns.md` — region and registry come from variables the operator supplies, amd64 for Google, `${REGION}-docker.pkg.dev/...:${GIT_SHA}`
- [ ] BUG-019 b: `skills/docker/SKILL.md` — the root step follows `dev-setup.md` (standing permission recorded, else hand the command back); the health probe probes the container it started
- [ ] BUG-019 c: the three CI sections (`skills/bash/references/testing-linting.md`, `skills/gcloud-setup/SKILL.md`, `skills/frontend/references/accessibility.md`) say the harness runs no CI and name the local form
- [ ] BUG-019 d: the four dependency defaults in `skills/golang/references/` are marked as choices that need justifying against stdlib first
- [ ] BUG-019 e: `skills/golang/references/golang-ecosystem.md` benchmark baselines go to `.agents/scratch/`
- [ ] BUG-020: `skills/cloud-run-deploy/SKILL.md` names `/preflight` as the entry point for ship/deploy/release words from the operator's own turn
- [ ] BUG-021: `skills/distill/SKILL.md` points at the `feedback` skill's canonical entry shape instead of restating it, and its example carries `project`
- [ ] BUG-022: `skills/golang/references/golang-ecosystem.md` scopes third-party registration to public work and names what must never be sent
- [ ] BUG-031: `context/AGENTS.md` — a headless run waits for every agent it dispatched before ending its turn, and says so when one has not
- [ ] G5: `skills/init-project/SKILL.md` says what an empty template section becomes
- [ ] G10: `skills/review/SKILL.md` — the round after fixes re-reads the changed declarations, not only the findings
- [ ] G12: a runtime data file defaults outside the repository or is gitignored (`context/AGENTS.md` or the golang skill)
- [ ] G13: `AGENTS.md` § Setup Commands stays the commands a contributor runs; tiers live in the testing guide (`agents/documentation.md` or the init-project template)
- [ ] G20: `agents/verifier.md` defines `<app>` (the `cmd/<app>` name; the repository name when there is none)
- [ ] `skills/preflight/SKILL.md` Step 4: option C files what it defers (handoff Follow-up)
- [ ] After each: grep the subject nouns across `context/`, `skills/`, `agents/`, `README.md`, `AGENTS.md` for restatements
- [ ] Checkpoint — `changes/` entry, commit, push

### Phase 3 — Docs

- [ ] `.agents/docs/deploy/remote.md` § The resync loop describes what `update` now does
- [ ] `README.md` and `AGENTS.md` where a command's behaviour changed
- [ ] `bugs/INDEX.md`: every record closed moves to Graduated with its new home (the `bugs` agent)
- [ ] Checkpoint — `changes/` entry, commit, push

### Phase 4 — Definition of Done over the batch

- [ ] `verifier` over the whole range
- [ ] `reviewer` (code lens) and `reviewer` (docs lens) over the whole range
- [ ] `security` (the `update` refusal and the skill text touching root steps, keys and external services)
- [ ] Findings fixed or rebutted in each report's `## Outcome`
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
