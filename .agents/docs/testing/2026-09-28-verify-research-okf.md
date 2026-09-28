---
owner: verifier
date: 2026-09-28
---
# Verification — research-okf   (2026-09-28)

Range reviewed: working tree vs HEAD (914dc90) — one new untracked file, `.agents/docs/plans/2026-09-28-research-okf.md`; no source/test/context/skill/agent file changed. Confirmed with `git status --short --untracked-files=all` at pre-flight: only that path shown.

## Status: PASS

Scope: docs-only change (a new plan file under `.agents/docs/plans/`). Ran Go checks (gofmt/vet/test -race, plus TestRenderedContextFitsAntigravityBlock -v per dispatch), scratch-home install + status drift check, `list`, and file-shape checks on the new docs file (frontmatter, numstat non-binary, no control bytes, with the grep validated to fire on a fixture first). complete — every item the dispatch named ran; `testing/guide.md` already existed and needed no update (no procedure changed).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0 (go version go1.27.1 linux/amd64)
- `go test -race ./...` → exit 0
  ```
  ok  	github.com/kesokaj/styrman-harness	1.101s
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	6.551s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	2.916s
  ```
  counts: pass 3 packages (no per-test -v breakdown captured in this run; see next line for the named test), skip 0, fail 0 at 914dc90+dirty for `go test -race ./...`; vs previous: no earlier verify log in testing/ carries a matching counts line for this exact command — vs previous: none.
- `go test ./internal/harness/... -run TestRenderedContextFitsAntigravityBlock -v` → exit 0
  ```
  === RUN   TestRenderedContextFitsAntigravityBlock
      size_test.go:48: Antigravity rules block is 20329 bytes (context 20329, rules 0), budget 22000
  --- PASS: TestRenderedContextFitsAntigravityBlock (0.00s)
  PASS
  ```
  Byte figure: rendered context 20,329 bytes, budget 22,000 — well under Antigravity's 24,000-byte global rules block (ADR-035). Unaffected by this change since no context/ or rules/ file was touched.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-okf` → exit 0, "258 written, 0 unchanged, 0 removed, 0 kept — 258 files managed" (fresh scratch home, so full initial write, not drift)
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-okf` → exit 0
  ```
  repo       /home/kesokaj/workspace/styrman-harness (main @ 914dc90, uncommitted changes)
  installed  2026-09-28T13:22:10+02:00 from commit 914dc90+dirty (vendor: all, 258 files, manifest ~/.config/sh-agent/state.json)
  binary     914dc90d59d2+dirty
  repo has   27 skills, 7 agents, 0 commands, 0 rules
  drift      none — 257 files match the repo
  ```
  "drift none" as required. "uncommitted changes" / "+dirty" is expected and correct — it reflects the untracked plan file under verification, which is not part of any vendor manifest (docs are not deployed by install) and so does not itself register as drift.
- `go run ./cmd/sh-agent list` → exit 0; printed 27 skills and 7 agents, table well-formed, matches `AGENTS.md` § Installed counts exactly.
- `head -5 .agents/docs/plans/2026-09-28-research-okf.md` → starts with `---\nowner: main\ndate: 2026-09-28\n---` — frontmatter present.
- `git diff --no-index --numstat /dev/null .agents/docs/plans/2026-09-28-research-okf.md` → `255  0  /dev/null => .../research-okf.md` (exit 1, expected for `--no-index` with a difference) — numeric columns, confirms the file is treated as TEXT, not binary.
- Control-byte grep, fixture check first: `printf 'line one\x01line two\n' > .agents/scratch/verify-okf-fixture/ctrlbyte.txt; grep -nP '[\x00-\x08\x0B-\x1F\x7F]' ctrlbyte.txt` → exit 0, matched line 1 — confirms the grep can fire.
- `grep -nP '[\x00-\x08\x0B-\x1F\x7F]' .agents/docs/plans/2026-09-28-research-okf.md` → exit 1, no output — no control bytes other than newline in the new file.

## Root cause
N/A — no failures.

## Coverage of changed code
- `.agents/docs/plans/2026-09-28-research-okf.md` is a plan/docs file, not source or test code; no unit test applies. Verified instead by shape (frontmatter present, non-binary per `git diff --no-index --numstat`, no stray control bytes). No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) — the change set contains no source or test file, so this check is vacuously satisfied (nothing to scan).
- All existing Go source/tests unaffected by this docs-only change; `go test -race ./...` covers the full module and is green.

## Not examined
- No CLI start (no `agy`/Claude Code session launch) — out of scope since nothing under `context/`, `agents/`, or `skills/` changed.
- No real-home smoke test — install/status ran only against the scratch home `.agents/scratch/verify-okf`, per the guide's rule to never touch the real `$HOME` from a dev session.
- No behavioural run of the new plan's content (it is a research/planning document, not executable spec or code to run).
- Did not run `govulncheck ./...` (listed in `AGENTS.md` § Setup commands but not named in the dispatch's numbered checklist) or § 3a/3b/3c of `testing/guide.md` (plain-copy install, `report --dry-run`, verb×place-kind matrix) — none apply, since no `cmd/sh-agent` or `internal/harness` rebuild/report/verb logic changed.
- Did not capture a per-test `-v` breakdown for `go test -race ./...` (ran bare per the guide); only the named `TestRenderedContextFitsAntigravityBlock` was run with `-v` as the dispatch specified.

## Environment
OS: Linux (Debian trixie, kernel 6.12.107+deb13-amd64). Go: go1.27.1 linux/amd64 (user-level SDK per `AGENTS.md` note, ahead of `/usr/bin/go`). Docker/Colima: not applicable, no Dockerfile changed. Scratch home used: `.agents/scratch/verify-okf` (left in place under gitignored scratch; not committed).
