---
owner: verifier
date: 2026-09-16
---
# Verification — clear-the-board (2026-09-16)

Range reviewed: e08fa65..HEAD (6008f32 Go changes/tests, 3c46356 instruction text/docs)

## Status: PASS

Scope: complete — all commands in the dispatch ran; every listed new/changed function mapped to a covering test that PASSed; the instruction-text half checked against the prose test plus a manual grep for `/tmp` mentions.

## Commands

- `gofmt -l cmd internal` (the AGENTS.md Setup command as of this range, not `gofmt -l .`) → exit 0, prints nothing.
- `go vet ./...` → exit 0.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-board/go-test-1.txt`; counts: pass 92, fail 0, skip 1 at 3c46356 for `go test -race -v -count=1 ./...` (subtests included, e.g. `TestUpdateNeverInstallsDirtyTree`'s 5 subtests and `TestManagedDirsRefuseSymlinks`'s 6). vs previous: `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md`, pass 79→92, skip 1→1 (same pre-existing `TestFindRepo` skip, unchanged reason). Rise in pass, no fall, no new skip — no FAIL condition.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-board/home` → exit 0; `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`; settings.json created with attribution off.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-board/home` → exit 0; `drift none — 259 files match the repo`; `repo has 27 skills, 7 agents, 0 commands, 0 rules`; repo line shows `main @ 3c46356, clean`.
- `go run ./cmd/sh-agent list` → exit 0; `### Skills (27)`, `### Subagents (7)` — matches expected counts.

## Coverage of changed code

Functions the dispatch named, each mapped to its covering test (all PASSed in the run above):

- `update` (cmd/sh-agent/main.go, dirty-checkout guard, BUG-046) — `TestUpdateNeverInstallsDirtyTree` (5 subtests: clean+up-to-date installs, dirty+up-to-date refuses, dirty+behind refuses without pulling, clean+behind pulls then installs, `--check` on dirty reports and installs nothing). Re-checked myself: read the test body: it asserts both the refusal error text and that nothing was installed/pulled on the dirty paths, and that the clean paths still install. I did not re-run it in the must-fail state (guard removed) myself — taking that from the main agent's record per the dispatch.
- `writeAtomic` (internal/harness/apply.go, `Sync` before rename) — `TestWriteAtomicIsAtomicAndDurable`. The test itself notes the added `Sync` is not directly observable without a crash; it asserts the surrounding contract (exact bytes, exact mode, no leftover temp file, safe overwrite). The `Sync`-specific crash-durability property is NOT covered by any test — acknowledged in the test's own comment. Re-checked myself: read the test and confirmed it does not and cannot exercise the Sync's actual effect.
- `Adopt` (internal/harness/adopt.go, stage-then-rename) — `TestAdoptFailedCopyKeepsTheExistingSkill` (copy fails mid-walk on a symlink; asserts the original SKILL.md is untouched and the `.adopt-staging` directory is cleaned up) plus the pre-existing `TestAdopt` and `TestAdoptRefusesDestructiveSources`. I re-checked this one myself in the must-fail state per the dispatch note (the main agent validated with staging removed); I did not redo that removal experiment myself, taking it from that record, but confirmed by reading the test that it forces the copy-failure branch (symlink in source) that the staging logic exists for.
- `readSettings` (cmd/sh-agent/main.go) — `TestReadSettingsProblems`, covering absent file, directory instead of file, and (non-root only) an unreadable regular file, and confirms the problem string never leaks file contents. Re-checked myself: read the test; all three legs run under the test's own euid guard.
- `abandonTemp` (cmd/sh-agent/main.go) — `TestAbandonTempReportsRemoveFailure`, forcing the "settings write failed AND temp removal failed" double-failure path via a read-only parent directory (non-root only); asserts both the original cause and "left behind" plus the temp path appear in output. Re-checked myself.
- `RenderContext` symlink refusal — `TestRenderContextRefusesSymlinks`, two subtests (context dir itself a symlink; a file inside context a symlink), each asserting refusal and that the error names "symlink". Re-checked myself.
- `TestNoTmpInHarnessProse` extended with `vendoredSkills(t, repo)` — parses `skills/VENDORED.md`'s table cells into a skip-set so the walk covers every non-vendored skill directory (terraform is NOT vendored per VENDORED.md but IS exempted by the harness's own listing — re-checked: `vendoredSkills` reads the table directly, so any project skill added later is picked up automatically; the assertion `len(out) == 0 → t.Fatal` guards against the parser silently matching nothing). Re-checked myself by reading the parser and confirming it stats each candidate directory before including it, and by independently grepping `/tmp` across the changed prose files (see below) — all hits fall inside an already-allowed "never /tmp" clause.

Not directly named in the dispatch but exercised by the same range: `harness.Inspect(...).Dirty` (used by both `update`'s guard and `TestUpdateNeverInstallsDirtyTree`) has no separate unit test beyond its use inside these subtests — acceptable, it is pre-existing plumbing, not new in this range.

## Instruction-text half of the range

`context/AGENTS.md`, `agents/documentation.md`, `agents/verifier.md`, and the touched `skills/**` reference files (bash, cloud-run-deploy, distill, docker, frontend/accessibility, gcloud-setup, golang references, init-project, review) are prose-only diffs. Covered by:
- Step 2 (`install` + `status`) proves they parse and deploy without drift.
- `TestNoTmpInHarnessProse`, which now walks every non-vendored skill's own references (not just SKILL.md), passed and its walk includes all the changed reference files above.
- Manual grep for `/tmp` across every changed prose file: all hits are inside an already-allowed "never `/tmp`" clause (context/AGENTS.md:53, agents/verifier.md:99, agents/documentation.md:60, skills/golang/references/golang-ecosystem.md:799).

What NO test covers for this half: content correctness of the prose (whether the BUG-019/020/021/022/031 closures and the plan-file text are factually accurate, well-formed markdown beyond frontmatter, or internally consistent with REQUIREMENTS.md) — that is `/review --docs` territory per this project's AGENTS.md Testing section, not something `go test` or the installer asserts. No evidence in this range's commits that `/review --docs` was run over 6008f32/3c46356; the dispatch did not ask me to run it and I did not.

## Not examined

- Did not myself reproduce the must-fail states for `TestUpdateNeverInstallsDirtyTree` (guard removed), the staging-removed state for `Adopt`, or the pre-fix prose-test tree — taken from the main agent's record as the dispatch permitted.
- Did not run `/review --docs` over the instruction-text commits (out of scope for this dispatch; noted above as a gap the record should track).
- Did not exercise the sandbox CLI-start smoke checks (`.agents/docs/testing/2026-09-03-verify-sandbox.md`) or the real-home smoke checks (`2026-09-03-verify-real-home.md`) — not named in this dispatch, and no source in `internal/harness/derive.go` or the agent-translation path changed in this range.
- Did not independently verify the `Sync()`-before-rename durability guarantee under an actual crash/power-loss — untestable without fault injection; the test file itself says so.

## Environment

Linux 6.12.107+deb13-amd64, go version per repo's go.mod toolchain (ran via `go test`/`go vet`/`go run` from PATH, no version pin issues encountered), non-root user (euid-gated test legs all ran, not skipped), git 3c46356 clean working tree at start and throughout (verify-only run, no source changes made).
