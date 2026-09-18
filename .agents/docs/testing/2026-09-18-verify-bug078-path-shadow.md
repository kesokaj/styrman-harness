---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-078 path shadow fix   (2026-09-18)

Range reviewed: 4cdc0b5..52574a0

## Status: PASS

Scope: complete — every command in `AGENTS.md` § Setup commands and the guide's §1/§3 tiers ran, the four named tests ran individually with `-count=1 -v`, the whole suite ran with `-count=1` and again with `-count=1 -v` for counts, the live (un-fakeable) install/status shadow check ran against a real scratch home, and `testing/guide.md` exists (no update needed — steps unchanged).

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 ./...` → exit 0.
  ```
  ok  	github.com/kesokaj/styrman-harness	1.098s
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	3.830s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	2.823s
  ```
- `go test -race -count=1 -v ./...` (for counts) → exit 0, `PASS` in all three packages.
  counts: pass 108 top-level + 103 subtests (211 `--- PASS:` lines total), fail 0, skip 0, at 52574a0 for `go test -race -count=1 -v ./...`.
  vs previous: `.agents/docs/testing/2026-09-18-verify-preflight-v1.2.1-fixes.md`, pass 104→108 top-level (191→211 total), skip 0→0. Rise fully explained: this range adds exactly 4 new top-level tests (`TestPathShadows`, `TestResolveExecutable`, `TestInstallAndStatusNameAPathShadow`, `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`) and 16 new subtests (10 under `TestPathShadows`, 6 under `TestResolveExecutable`) — matches 104+4=108 and 87+16=103 exactly. No pre-existing test lost.
- `go test -race -count=1 -v ./internal/harness/... -run 'TestPathShadows|TestResolveExecutable'` → exit 0.
  ```
  --- PASS: TestPathShadows (0.00s)
      --- PASS: TestPathShadows/installed_dir_first_—_nothing_shadows_it
      --- PASS: TestPathShadows/an_earlier_entry_wins
      --- PASS: TestPathShadows/every_earlier_entry_is_named,_in_PATH_order
      --- PASS: TestPathShadows/a_link_to_the_installed_binary_is_not_a_shadow
      --- PASS: TestPathShadows/a_repeated_PATH_entry_is_reported_once
      --- PASS: TestPathShadows/a_directory_that_holds_no_sh-agent_is_skipped
      --- PASS: TestPathShadows/the_installed_dir_absent_from_PATH:_every_candidate_wins
      --- PASS: TestPathShadows/an_unresolvable_installed_binary_still_reports_shadows
      --- PASS: TestPathShadows/an_empty_PATH_entry_means_the_working_directory
      --- PASS: TestPathShadows/a_trailing_slash_on_the_installed_dir_still_matches
  --- PASS: TestResolveExecutable (0.00s)
      --- PASS: TestResolveExecutable/an_executable_file
      --- PASS: TestResolveExecutable/a_symlink_resolves_to_its_target,_so_a_link_is_not_a_second_binary
      --- PASS: TestResolveExecutable/a_file_with_no_execute_bit
      --- PASS: TestResolveExecutable/a_dangling_symlink
      --- PASS: TestResolveExecutable/a_directory
      --- PASS: TestResolveExecutable/a_missing_path
  PASS
  ok  	github.com/kesokaj/styrman-harness/internal/harness	1.015s
  ```
- `go test -race -count=1 -v ./cmd/sh-agent/... -run 'TestInstallAndStatusNameAPathShadow|TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice'` → exit 0.
  ```
  --- PASS: TestInstallAndStatusNameAPathShadow (0.09s)
  --- PASS: TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice (0.06s)
  PASS
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	1.172s
  ```
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug078` → exit 0. Live, un-fakeable check (the scratch home's bin dir is off my real PATH; a real `sh-agent` at `/home/kesokaj/.local/bin/sh-agent` is on it):
  ```
  building sh-agent from /home/kesokaj/workspace/styrman-harness
  260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
  note: .agents/scratch/verify-bug078/.local/bin is not on PATH, so `sh-agent` runs /home/kesokaj/.local/bin/sh-agent instead; put .agents/scratch/verify-bug078/.local/bin ahead of it
  settings: .agents/scratch/verify-bug078/.claude/settings.json created with "attribution": {"commit": "", "pr": ""} — Claude Code's commit and PR trailers are off
  ```
  Matches the expected shape: the `note:` line names the real binary that answers instead, and there is NO `warning: … comes before` line (the bin dir is off PATH, so `pathHint` owns the message and `shadowHint` short-circuits via `binOnPath`).
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug078` → exit 0.
  ```
  repo       /home/kesokaj/workspace/styrman-harness (main @ 52574a0, uncommitted changes)
  installed  2026-09-18T12:42:24+02:00 from commit 52574a0+dirty (vendor: all, 260 files, manifest ~/.config/sh-agent/state.json)
  binary     52574a0ed219+dirty
  retro      none yet — run /retro at the end of a session with corrections
  repo has   27 skills, 7 agents, 0 commands, 0 rules
  claude     ~/.claude  (AGENTS.md, CLAUDE.md, skills/, agents/, commands/, rules/)
  agy        ~/.gemini/config  (rules/AGENTS.md, skills/, agents/, workflows/, rules/)
  drift      none — 259 files match the repo
  ```
  No `shadowed` line, as expected (the scratch home's bin dir is off PATH, so the `binOnPath(home)` guard in `status` suppresses the shadow check entirely). "uncommitted changes" refers to untracked review files from a concurrent dispatch (`.agents/docs/reviews/2026-09-18-review-*.md`), not anything this run created.
- `go run ./cmd/sh-agent list` → exit 0, printed 27 skills / 7 agents tables, no parse errors.
- `govulncheck ./...` → exit 0, `No vulnerabilities found.` (go1.27.1, consistent with AGENTS.md's expectation).

## Root cause  (FAIL only)

N/A — all commands passed.

## Coverage of changed code

Every new/changed non-test function in the diff is named by a test:
- `internal/harness/pathguard.go: PathShadows` → `TestPathShadows` (10 subtests).
- `internal/harness/pathguard.go: ResolveExecutable` → `TestResolveExecutable` (6 subtests).
- `cmd/sh-agent/main.go: binOnPath` → exercised by `TestInstallAndStatusNameAPathShadow` and `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`, and by the live install/status check above.
- `cmd/sh-agent/main.go: shadows` → same two tests (feeds both `shadowHint` and the `status` shadow loop).
- `cmd/sh-agent/main.go: shadowHint` (new) → `TestInstallAndStatusNameAPathShadow` (asserts the `warning: … comes before …` line on install when the bin dir IS on PATH) and `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice` (asserts it does NOT fire when off PATH — confirmed live above too).
- `cmd/sh-agent/main.go: pathHint` (modified — now names what answers) → `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice` and the live install check above (`note: … is not on PATH, so \`sh-agent\` runs … instead`).
- `cmd/sh-agent/main.go: status`'s new `shadowed` line → `TestInstallAndStatusNameAPathShadow` (`wantStatus := "shadowed   " + other + " comes first on PATH, ..."`) and the live status check above (confirmed absent when the bin dir is off PATH).
- `README.md`, `.agents/docs/deploy/remote.md` — doc-only prose changes (verify-before-move / resolution-after wording); no executable code, nothing to cover.

Suppressions added in the diff (`git diff 4cdc0b5..52574a0 -- cmd/sh-agent/main_test.go internal/harness/pathguard_test.go` grepped for `nolint|lint:ignore|ts-ignore|ts-expect-error|eslint-disable|t.Skip|SkipNow|go:build ignore|.skip(|.only(`): none found.

No changed file lacks a covering test.

## Not examined

- `README.md` / `.agents/docs/deploy/remote.md` prose changes were read for shape but not run through `/review --docs` — that is the reviewer's job, not verifier's; noting it is out of this run's scope.
- Guide tiers not triggered by this range: §2 harness self-check (`render`) and §4 sandbox CLI starts — this change touches no `context/`, `agents/`, or `skills/` file, so those tiers do not apply; §3a plain-copy install and §3b report loop — this change does not touch `IsCheckout`/rebuild logic or `report.go`, so those tiers do not apply either.
- Did not run `sh-agent report` (forbidden by the dispatch) and did not touch the real `$HOME` (forbidden); no git command altered the working tree.
- Did not check whether the pre-existing `.agents/docs/reviews/2026-09-18-review-4cdc0b5-52574a0*.md` untracked files (seen in `git status --short`, not created by this run) were already reviewed/resolved — out of scope for verification.

## Environment

Linux 6.12.107+deb13-amd64, go1.27.1 (user-level SDK, per AGENTS.md), govulncheck at `/home/kesokaj/.go/bin/govulncheck`. Real `sh-agent` present at `/home/kesokaj/.local/bin/sh-agent` on the operator's actual PATH — this is what the live shadow-check note line correctly names.
