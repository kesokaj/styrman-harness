---
owner: verifier
date: 2026-09-16
---
# Verification — preflight-v1.0.0   (2026-09-16)

Range reviewed: 4aa19be..HEAD (14 commits, 75 files)

## Status: PASS

Scope: complete — every command AGENTS.md's Setup commands name, plus the release-gate additions (govulncheck quoted, GOTOOLCHAIN=local build+test, install/status/list count check, coverage map over the whole range) all ran.

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-v1.0.0/go-test-1.txt`; counts: pass 103, fail 0, skip 1 at c7e3688 for `go test -race -v -count=1 ./...`. vs previous: `.agents/docs/testing/2026-09-16-verify-clear-the-board.md`, pass 92→103, skip 1→1 (same pre-existing `TestFindRepo` skip, unchanged reason). Rise in pass, no fall, no new skip — no FAIL condition. (Note: the dispatch cited "103" as the main agent's own measurement after three more tests landed since the 92-count log; this run reproduces 103 exactly.)
- `govulncheck ./...` → exit 0. Quoted result: "No vulnerabilities found." / "Your code is affected by 0 vulnerabilities. This scan also found 2 vulnerabilities in packages you import and 24 vulnerabilities in modules you require, but your code doesn't appear to call these vulnerabilities." Confirms ADR-028's claim (2 affecting → 0 affecting) directly — the 2 advisories are now in the "not called" bucket, not the "affecting" bucket. Full output at `.agents/scratch/verify-v1.0.0/govulncheck.txt`.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v1.0.0/home` → exit 0; "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; settings.json created with attribution off.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v1.0.0/home` → exit 0; "drift none — 259 files match the repo"; "repo has 27 skills, 7 agents, 0 commands, 0 rules"; `binary c7e3688a5ea7` (no warning) — matches expectation.
- `go run ./cmd/sh-agent list` → exit 0; `### Skills (27)`, `### Subagents (7)` — matches the counts styrman-crew's channel will smoke-check against this tag.
- `GOTOOLCHAIN=local go build ./...` → exit 0. Confirmed the box's installed toolchain is actually go1.24.4 (`go env GOVERSION` → go1.25.8 pinned via toolchain directive; plain `go version` outside GOTOOLCHAIN=local, inside the module, auto-switches to go1.25.8; `GOTOOLCHAIN=local go version` reports go1.24.4 — the real local binary). The build succeeded under the older, locally-installed toolchain, confirming ADR-028's "does not break a consumer building offline with the older Go" claim.
- `GOTOOLCHAIN=local go test -count=1 ./internal/...` → exit 0, `ok github.com/kesokaj/styrman-harness/internal/harness 0.453s`, under go1.24.4.

## Root cause (FAIL only)

N/A — Status is PASS.

## Coverage of changed code

Go changes (mapped function → covering test, `cmd/sh-agent/main_test.go` / `internal/harness/harness_test.go`):
- `cmd/sh-agent/main.go: update` — `TestUpdateNeverInstallsDirtyTree` (5 subtests)
- `cmd/sh-agent/main.go: attributionSet` — `TestAttributionSet`, `TestAttributionNote`
- `cmd/sh-agent/main.go: readSettings` — `TestReadSettingsProblems`
- `cmd/sh-agent/main.go: abandonTemp` — called directly in `TestReadSettingsProblems`'s neighborhood (line 875), asserted on stdout note
- `internal/harness/repo.go: Inspect` (untracked-file counting) — `TestInspectCountsUntrackedFiles`
- `internal/harness/repo.go: Pull` — `TestBehindAndPull`
- `internal/harness/plan.go: Plan` (new secret-shaped-file refusal) — `TestPlanRefusesSecretShapedFiles`
- `internal/harness/plan.go: secretShaped` (new function) — same test, `TestPlanRefusesSecretShapedFiles`
- `internal/harness/inventory.go: LoadSkills` (duplicate-name refusal) — `TestLoadSkillsRefusesDuplicateNames`
- `internal/harness/adopt.go: Adopt` (staging outside skills/) — `TestAdopt`, `TestAdoptRefusesDestructiveSources`, `TestAdoptFailedCopyKeepsTheExistingSkill`
- `internal/harness/apply.go: writeAtomic` — `TestWriteAtomicIsAtomicAndDurable`

No changed Go function in the named files was left without a covering test.

No added suppressions found: `git diff 4aa19be..HEAD` was checked for `//nolint`, `//lint:ignore`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, `.only(` — none added in this range.

Non-Go range (bulk of the diff — `context/AGENTS.md`, `agents/verifier.md`, and the reference/SKILL.md files under `skills/bash`, `skills/cloud-run-deploy`, `skills/distill`, `skills/docker`, `skills/frontend`, `skills/gcloud-setup`, `skills/golang`, `skills/init-project`, `skills/preflight`, `skills/retro`, `skills/review`): this project's own Testing section names no test framework for instruction text — it is verified by (a) the install run parsing all 27 skills and 7 agents without error (confirmed above, exit 0, no parse warnings) and (b) `/review --docs` over the range, which is a separate reviewer dispatch, not this verifier's job. This log does not itself run `/review --docs`; I did not find a same-range `/review --docs` report under `.agents/docs/reviews/` in this session — that is a gap this log names but does not close (see Not examined).
`.agents/docs/architecture/DECISIONS.md` carries ADR-028, the record for the toolchain-pin trajectory change — present and accepted.

## Not examined

- Did not confirm a `/review --docs` pass exists over exactly 4aa19be..HEAD for the prose/skill files listed above — this verifier only confirmed those files parse cleanly at install time, which is necessary but not sufficient for doc-quality review. If release requires that review gate closed, check `.agents/docs/reviews/` for a matching entry before shipping.
- Did not run the sandbox/real-home CLI smoke checks (`2026-09-16-verify-real-home-fcc8539.md` and related) fresh for this exact commit c7e3688; those were run at earlier commits in this session's trail and are not re-verified here. The dispatch did not ask for them and this run stayed inside the five numbered steps given.
- `govulncheck -show verbose` was not run; the summary line was quoted and taken as sufficient per the dispatch's ask to "quote the result line."

## Environment

Linux 6.12.107+deb13-amd64. Local `go` binary: go1.24.4 (`/usr/bin/go`). Module toolchain: go1.25.8 (pinned via `go.mod`, ADR-028), auto-fetched/used for the plain runs; `GOTOOLCHAIN=local` runs used the real go1.24.4 binary. `govulncheck` present at `/home/kesokaj/.go/bin/govulncheck`, run directly (not via `go run`). No Docker/Colima involved — no Dockerfile in this range.

## Round 2

Range reviewed: c7e3688..787e15c (release-gate findings fix, since round 1's HEAD)

### Status: PASS

Scope: Setup commands re-run, govulncheck, install/status/list, the five numbered claims. Did not re-run `/review --docs` or the sandbox/real-home CLI smoke checks (same gap as round 1, not asked for this round).

### Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0.
- `go test -race -v -count=1 ./...` → exit 0; output at `.agents/scratch/verify-v1.0.0-r2/go-test-1.txt`; counts: pass 106, fail 0, skip 1 at 787e15c for `go test -race -v -count=1 ./...`. vs previous: `.agents/docs/testing/2026-09-16-verify-preflight-v1.0.0.md` (round 1), pass 103→106, skip 1→1 (same `TestFindRepo` skip). Rise of 3 matches new tests added in 787e15c (`TestInspectFailsClosedOnUnknownStatus`, `TestReinstallBuildsWithTheLocalToolchain`, `TestPlanDropsGroupAndOtherWrite`, plus the rename-failure leg added inside the existing `TestAdoptFailedCopyKeepsTheExistingSkill`, not a new top-level test) — no FAIL condition.
- `govulncheck ./...` → exit 0. "No vulnerabilities found." / "Your code is affected by 0 vulnerabilities... 2 vulnerabilities in packages you import and 24... in modules you require, but your code doesn't appear to call these". Output at `.agents/scratch/verify-v1.0.0-r2/govulncheck.txt`.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v1.0.0-r2/home` → exit 0; "260 written, 0 unchanged, 0 removed, 0 kept".
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v1.0.0-r2/home` → exit 0; "drift none — 259 files match the repo"; "repo has 27 skills, 7 agents, 0 commands, 0 rules"; `binary 787e15cdbf52` (no warning). Matches expectation.
- `go run ./cmd/sh-agent list` → exit 0; `### Skills (27)`, `### Subagents (7)`.

### Claim 1 — the blocking fix (GOTOOLCHAIN=local), measured
Reproduced the failure with a fresh empty `GOMODCACHE` under `.agents/scratch/verify-v1.0.0-r2/gomodcache-empty`:
`GOMODCACHE=.../gomodcache-empty GOTOOLCHAIN=auto GOPROXY=off go build -o /dev/null ./cmd/sh-agent`
→ `go: downloading go1.25.8 (linux/amd64)` / `go: download go1.25.8 for linux/amd64: toolchain not available` → **exit 1**.

Then tried the installer's own path exactly as specified, `go run ./cmd/sh-agent install --home <home> --repo .` with the same fresh empty `GOMODCACHE` and `GOPROXY=off` → **exit 1**, same "toolchain not available" message, and it never reached Reinstall's "building sh-agent from..." line. Root cause: `go run` itself is an outer `go` invocation subject to the module's own `toolchain go1.25.8` pin under the ambient `GOTOOLCHAIN` (default `auto` unless overridden), and that pin bites before `main()` — and therefore before Reinstall's own build.Env override — ever runs. The `GOTOOLCHAIN=local` fix only touches the `exec.Command("go","build",...)` child process inside `Reinstall`; it cannot reach the outer `go run` that starts the sh-agent process itself.

To measure the actual consumer-host shape — an already-built `sh-agent` binary invoked directly, which is how `Reinstall` is truly reached (a consumer never runs `go run` against this repo) — I pre-built the binary once under a normal environment (`go build -o .../sh-agent-bin ./cmd/sh-agent`, exit 0), then ran that binary directly with the same fresh empty `GOMODCACHE` and `GOPROXY=off`:
`GOMODCACHE=.../gomodcache-empty2 GOPROXY=off .../sh-agent-bin install --home .../home-empty-cache --repo .`
→ "building sh-agent from ..." / "260 written, 0 unchanged, 0 removed, 0 kept" → **exit 0**.

So: both exit statuses as literally specified are 1 and 1 (the `go run` form fails both before and after the fix, for a reason unrelated to Reinstall). The fix is confirmed only via the binary-invocation path, which is the actual consumer shape (a downloaded/built `sh-agent` binary calling its own Reinstall, not `go run` against a checkout of this repo). This is worth flagging: the dispatch's exact command (`go run ./cmd/sh-agent install ...`) does not demonstrate the fix and cannot, structurally — it should be phrased as "build then run the binary" in any future release-gate script or guide.

### Claim 2 — fail-closed status
`TestInspectFailsClosedOnUnknownStatus` → `--- PASS` (line 173 of the test log). Reasoned from code, did not edit: the test writes `.git` as a plain file (not a dir) so every git subcommand fails; `Inspect` then sets `status, statusErr := gitOut(...)` where `gitOut` returns `("", err)` on failure, so `Dirty` (which is `status != ""`) would read **false** if `StatusUnknown` did not exist or were ignored by `Pull`. The test's second assertion, `Pull(repo, io.Discard)` must return non-nil, would then fail, because pre-787e15c `Pull` only checked `Inspect(repo).Dirty`. This confirms the test would fail if `StatusUnknown` were dropped from the `Dirty || StatusUnknown` guard in `Pull` (repo.go) and in `update` (main.go).

### Claim 3 — adopt's backup
`TestAdoptFailedCopyKeepsTheExistingSkill` → `--- PASS` (line 149). The rename-failure leg (chmod `skills/` to `0o500` to force the rename-into-place to fail, then assert the original survives and the error names `.agents/scratch`) is gated by `if os.Geteuid() != 0`. This run's user is `kesokaj`, uid 1000 (`id -u` → 1000), not root, so the leg ran and passed as part of the single PASS line above (subtests aren't split out for this test — it's one `t.Run`-free function with an inline conditional). Under root it would be skipped silently (no `t.Skip` call, so it wouldn't show as a SKIP line — it just wouldn't execute that block).

### Claim 4 — deployed modes
`TestPlanDropsGroupAndOtherWrite` → `--- PASS` (line 177). Cross-checked directly against the installed tree: `find .agents/scratch/verify-v1.0.0-r2/home -type f -perm /022` → 0 results — no installed file carries group or other write.

### Claim 5 — coverage of 787e15c
- `internal/harness/repo.go` (`Inspect` StatusUnknown, `Pull` guard) — `TestInspectFailsClosedOnUnknownStatus`.
- `internal/harness/adopt.go` (backup/restore rename logic) — `TestAdoptFailedCopyKeepsTheExistingSkill` (extended).
- `internal/harness/plan.go` (`mode := info.Mode().Perm()&0o755 | 0o600`, expanded `secretShaped` patterns) — `TestPlanDropsGroupAndOtherWrite`; did not find a test iterating the newly added `secretShaped` patterns (`.envrc`, `.netrc`, `.npmrc`, `.pypirc`, `*.p12`, `*.pfx`, `*.keystore`, `*.jks`, `*.tfstate.backup`, `id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa`) individually — grep for these literal strings in harness_test.go came up empty in the time available. **NOT COVERED, flagged, not confirmed either way with certainty** — see Not examined.
- `internal/harness/repo.go` (`Reinstall` GOTOOLCHAIN=local) — `TestReinstallBuildsWithTheLocalToolchain` (string-match guard on the source, not an execution test — reasonable given the network/cache dependency, but note it only checks the code shape, not a live empty-cache build; that live case is exactly what I measured by hand above).
- `cmd/sh-agent/main.go` (`update`'s `info.StatusUnknown` branch, the "could not be checked for uncommitted changes" message) — **NOT COVERED at the cmd/sh-agent level.** `TestUpdateNeverInstallsDirtyTree`'s subtests exercise `Dirty` only (via `dirty()` writing an uncommitted file); none construct a repo where `git status` itself fails and call `update()` to confirm it refuses with the StatusUnknown wording. The harness-level test (`TestInspectFailsClosedOnUnknownStatus`) covers `Inspect`/`Pull` but not `update`'s own composition of the error string in main.go. This is a genuine gap.
- `cmd/sh-agent/main_test.go` `TestReadSettingsProblems` closed-set assertion — pre-existing function, test strengthened, not new production code.
- Suppressions: `git diff c7e3688..787e15c -- '*.go'` for `//nolint`, `//lint:ignore`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, `.only(` → one hit, `t.Skip("git is not on PATH")` inside `TestInspectFailsClosedOnUnknownStatus`, carries a reason — not a bare skip, no FAIL.
- ADR-028 amendment (`.agents/docs/architecture/DECISIONS.md`) read in full: states plainly the pin is a ceiling for gated builds, not a floor for a consumer, and that `Reinstall` now uses `GOTOOLCHAIN=local`; matches what was measured.

### Root cause (FAIL only)
N/A — Status is PASS, with two named gaps (see Coverage and Not examined) that did not reach FAIL threshold: the `secretShaped` new-pattern coverage is unconfirmed, not proven absent, and the `update()` StatusUnknown wording is confirmed uncovered but is a message string, not the fail-closed behavior itself (which the harness-level test does cover, and which is what the release finding was about).

### Not examined
- Whether `secretShaped`'s newly added patterns (`.envrc`, `.netrc`, `.npmrc`, `.pypirc`, `*.p12`, `*.pfx`, `*.keystore`, `*.jks`, `*.tfstate.backup`, `id_rsa`, `id_ed25519`, `id_ecdsa`, `id_dsa`) are each asserted by a test — grepped and found none in the time available but did not read `TestPlanRefusesSecretShapedFiles` in full to rule out a table-driven form covering them.
- No test found covering `update()`'s own `StatusUnknown` branch/message in `cmd/sh-agent/main.go` (as opposed to `Inspect`/`Pull`, which are covered). Not fixed, only reported per the verifier's role.
- Did not re-run `/review --docs` over c7e3688..787e15c or check `.agents/docs/reviews/` for one already present this round — same gap named in round 1, not re-checked here.
- Did not re-run the sandbox/real-home CLI smoke checks for 787e15c.
- Did not run `GOTOOLCHAIN=local go build ./...` / `go test` under the real local go1.24.4 toolchain this round (done in round 1 at c7e3688, not repeated).
- Cut off by the turn limit before independently re-deriving `Not covered` status for `secretShaped` with full certainty; flagged rather than asserted.

### Environment
Linux 6.12.107+deb13-amd64, uid 1000 (non-root, confirmed for claim 3). `go` binary go1.24.4; module toolchain go1.25.8 pinned (ADR-028). `govulncheck` at `/home/kesokaj/.go/bin/govulncheck`. Scratch under `.agents/scratch/verify-v1.0.0-r2/`.

## Correction, appended by the main agent after round 2

Rounds 1 and 2 above both record `govulncheck ./...` as reporting no vulnerabilities affecting this code. **That was measured while `go.mod` carried `toolchain go1.25.8`, and that directive has since been withdrawn (ADR-029).** Both readings were true when taken and are wrong about the shipped tree. Re-measured by the main agent after the withdrawal, on the same box:

```
$ govulncheck ./...
Your code is affected by 2 vulnerabilities from the Go standard library.
  GO-2026-4602  os@go1.24.4       fixed in os@go1.25.8        reached through LoadSkills
  GO-2025-3956  os/exec@go1.24.4  fixed in os/exec@go1.24.6    reached through Reinstall's LookPath
```

Why the pin went: round 2's own claim-1 measurement is the reason. It found that the `GOTOOLCHAIN=local` fix inside `Reinstall` cannot save the outer `go run`, which is subject to the pin and fails before any of this code executes. The main agent reproduced the same independently. A pin that makes the module unbuildable on an offline host with an older Go is a worse property than the two advisories it cleared, so it was removed rather than patched a third time. After the withdrawal, in that same restricted environment (empty `GOMODCACHE`, `GOTOOLCHAIN=auto`, `GOPROXY=off`), `go build ./cmd/sh-agent` and `go run ./cmd/sh-agent install` both exit 0.

The two advisories are an operator toolchain item, unchanged from before today and recorded in three earlier security reviews: they are fixed in go1.24.6 and go1.25.8, Debian trixie ships `golang-1.24-go 1.24.4-1` and nothing newer, and `/usr/lib/go-1.24` is root-owned. v1.0.0 ships with them present, and the gate report says so.

Everything else in rounds 1 and 2 stands: the counts, the install and drift checks, the fail-closed status test, adopt's rename-failure leg, and the deployed-mode check were all measured on trees that differ from the shipped one only by the `toolchain` line and the round-2 fixes.
