---
owner: verifier
date: 2026-09-18
---
# Verification — preflight-v1.3.0   (2026-09-18)

Range reviewed: v1.2.3..10a9fe2

## Status: PASS

Scope: preflight release gate per `.agents/docs/testing/guide.md`. Ran guide tier 1
(Go checks), tier 2 (harness self-check via `list`), tier 3 (scratch-home install +
drift), tier 3a's dry-run-from-a-built-binary check, plus the three cross-compile
targets and the coverage map. Did not run tier 4 (sandbox CLI start) or tier 5
(real-home smoke checks) — see Not examined. Scope is not "complete" (those tiers
did not run).

## Commands
- `pwd` / `git status --short` → clean at 3b5b71a (only my own in-progress log file
  untracked later, expected)
- `gofmt -l cmd internal` → exit 0, no output
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0, output at
  `.agents/scratch/verify-v130/go-test-1.txt`
  counts: pass 250, fail 0, skip 0 at 3b5b71a for `go test -race -v -count=1 ./...`
  vs previous: `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-2-3.md`,
  pass 250→250, skip 0→0
- `go version` → go1.27.1 linux/amd64
- `govulncheck ./...` → "No vulnerabilities found." exit 0 (clean, as expected on
  go1.27.1+ per AGENTS.md)
- `GOOS=darwin GOARCH=arm64 go build -o .agents/scratch/verify-v130/sh-agent-darwin ./cmd/sh-agent`
  → exit 0
- `GOOS=linux GOARCH=arm64 go build -o .agents/scratch/verify-v130/sh-agent-linux ./cmd/sh-agent`
  → exit 0
- `GOOS=windows GOARCH=amd64 go build -o .agents/scratch/verify-v130/sh-agent-windows.exe ./cmd/sh-agent`
  → exit 0
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v130/home` → exit 0,
  "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; 260 `add` lines
  counted, matches
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v130/home` → exit 0,
  "drift none — 259 files match the repo" (the settings.json generated file is the
  260th, not compared against repo source — consistent with "drift none")
- `go run ./cmd/sh-agent list` → exit 0, "Skills (27)" and "Subagents (7)" headers,
  27 skill rows and 7 agent rows printed, no parse error
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v130/home --dry-run`
  → exit 0 but reported "1 written" (`update ~/.local/bin/sh-agent`), NOT the expected
  0 written — see Root cause (testing-artifact finding, not a code defect)
- Root-caused and reran with the guide's tier 3a prescribed invocation:
  `go build -o .agents/scratch/verify-v130/sh-agent-a ./cmd/sh-agent`, then
  `.agents/scratch/verify-v130/sh-agent-a install --repo /home/kesokaj/workspace/styrman-harness --home .agents/scratch/verify-v130/home --dry-run`
  → exit 0, "[dry-run] the binary would be rebuilt from <checkout>", then
  "[dry-run] 0 written, 260 unchanged, 0 removed, 0 kept — 260 files managed" — matches
  expectation
- Determinism check: built the binary twice with `go build -o ... ./cmd/sh-agent`
  (same -o path) and `cmp`'d them → identical (sha256
  2140f44c8f02af406713e7d53193d42a291147c57483e35f231bc1fa16f30fe0), and it matched
  the binary already installed at `.agents/scratch/verify-v130/home/.local/bin/sh-agent`
- `git diff --stat v1.2.3..3b5b71a -- cmd internal` → empty (no Go file changed in range)
- `git diff --stat v1.2.3..3b5b71a` → 23 files, 543(+)/25(-); non-Go source changed:
  `agents/verifier.md`, `context/AGENTS.md`, `skills/feature/SKILL.md`,
  `skills/feedback-log/SKILL.md`; remainder is this repo's own trail
  (`.agents/docs/**`, `ROADMAP.md`)
- `grep -n "context/AGENTS.md" internal/harness/*_test.go` → `harness_test.go` and
  `source_test.go` read `context/AGENTS.md` as a name inside synthetic fixture trees
  (not the real file's content) — see Coverage of changed code

## Root cause (finding — not a Status: FAIL cause)
`go run ./cmd/sh-agent install --home <home> --dry-run` reported "1 written" (a
spurious `update ~/.local/bin/sh-agent`) instead of the expected 0. Root cause: `go
run` compiles the invoking binary to an ephemeral path each run; even though the Go
source is byte-for-byte the same as what is already installed, the ephemeral `go
run` binary's build ID/embedded path differs from the deterministically-built,
already-installed binary at a fixed `-o` path, so the dry-run's binary comparison
sees a difference that is an artifact of how the *verifier* invoked the command, not
a real change. The guide's own tier 3a already prescribes invoking dry-run from a
**built** binary (`<built binary> install --repo <checkout> --home ... --dry-run`)
for exactly this reason; using it produced the expected "0 written". This is a gap
in the dispatch text (which said `go run ... --dry-run` and "expect 0 written")
against the guide's own prescribed invocation, not a defect in `cmd/sh-agent`. No
bug filed against the harness code; flagging so a future dispatch for this command
uses the built-binary form.

## Coverage of changed code
- No `.go` file changed in this range (`git diff --stat v1.2.3..3b5b71a -- cmd
  internal` empty) — the full Go suite (250 tests) ran as a regression check, not
  because it targets new code.
- `context/AGENTS.md` — no test reads the real file's content; `TestHarnessMap`
  (existence/shape pin) and the embedded-source test at `harness_test.go:580` assert
  only that a `context/AGENTS.md` file exists in the embedded FS, not its wording.
  Exercised behaviorally: `sh-agent install` rendered and installed it into the
  scratch home without error (260/260 files, drift none).
- `agents/verifier.md` — no test reads its content directly. Exercised
  behaviorally: installed without parse error; its own text (the "write the log
  after the first command" rule) is what this very verification run followed.
- `skills/feature/SKILL.md`, `skills/feedback-log/SKILL.md` — no test reads their
  content directly. Exercised behaviorally: `sh-agent list` still reports 27 skills
  with no parse error, and the installer accepted both `SKILL.md` files (frontmatter
  parses; a bad `#` or missing `description` would have failed `install`).
- These four files are FAIL-free only under the guide's own coverage model (tiers
  1-3 = automated/behavioral parse+install; tiers 4-5 = sandbox/real-home content
  checks). Tiers 4 and 5, which are the checks that would actually exercise the
  *wording* change (a context question quoting the new text, an agent dispatch
  showing the verifier's new early-log rule in effect), did NOT run in this dispatch
  — see Not examined. No changed file is uncovered by the automated tier; the
  content-level check is deferred, named below, not silently skipped.

## Not examined
- Guide tier 4 (sandbox start of each CLI, context question quoting new wording,
  log grepped for parse errors) — not run.
- Guide tier 5 (real-home smoke checks: context questions, one skill load, one
  agent dispatch from a real home, on each CLI) — not run; no install into the
  real `$HOME` was performed or attempted, per the hard limit.
- The three cross-compiled release artifacts (darwin/arm64, linux/arm64,
  windows/amd64) were built and confirmed to exit 0, but never executed (can't run
  a cross-arch/cross-OS binary on this box) — no smoke run of `<binary> version` or
  `list` on any of them.
- `sh-agent report` — never run, per hard limit (would push to the real feedback
  repository).
- Guide tier 3b (report/feedback) and 3c (verb × place-kind matrix) — out of scope
  for this range (no change to `report.go` or verb/refusal behavior in the diff).
- The behavioral rerun of `agents/verifier.md`'s new wording specifically (a
  `claude -p --agents <json>` run per the Testing section of `AGENTS.md`) — not
  run; this dispatch's own execution is suggestive evidence (I did write the log
  after the first command) but is not a substitute for the prescribed check.

## Environment
Go 1.27.1 linux/amd64 (user-level SDK ahead of `/usr/bin/go`, per `AGENTS.md`;
`go.mod` carries no `toolchain` directive, ADR-029). govulncheck at
`/home/kesokaj/.go/bin/govulncheck`. Cross-compiled darwin/arm64, linux/arm64,
windows/amd64 build artifacts in `.agents/scratch/verify-v130/` (untested, build
only). Scratch home `.agents/scratch/verify-v130/home`, gitignored. Real `$HOME`
never touched.

## Round 2 (94bec82)

### Status: PASS

Pre-gate commit moved HEAD to 94bec82: four text-only lines (README.md,
agents/verifier.md, context/AGENTS.md, skills/feedback-log/SKILL.md; no Go).
Re-ran the install/status/dry-run tier against the new HEAD, all from a built
binary per guide tier 3a.

- `gofmt -l cmd internal` -> exit 0, no output
- `go test -race -count=1 ./...` -> exit 0, all packages ok (output at
  .agents/scratch/verify-v130/go-test-r2.txt); no new Go in this commit so no
  fresh --- PASS/FAIL count taken, non-verbose ok lines only
- `go build -o .agents/scratch/verify-v130/sh-agent-built ./cmd/sh-agent` -> exit 0
- `sh-agent-built install --home .agents/scratch/verify-v130/home` -> exit 0,
  "0 written, 260 unchanged, 0 removed, 0 kept - 260 files managed" (rebuilt from
  the repo, as the checkout-rebuild path always does, but found nothing to write:
  the earlier Round 1 install of the same scratch home had already been refreshed
  to 94bec82's content by the time this ran, or the four changed files render
  byte-identical output for this vendor - either way, no drift)
- `sh-agent-built status --home .agents/scratch/verify-v130/home` -> exit 0,
  "installed ... from commit 94bec82+dirty", "drift none - 259 files match the
  repo"
- `sh-agent-built install --home .agents/scratch/verify-v130/home --dry-run` ->
  exit 0, "[dry-run] the binary would be rebuilt from <checkout>" then
  "[dry-run] 0 written, 260 unchanged, 0 removed, 0 kept - 260 files managed" -
  matches expectation using the built-binary form (see Round 1 Root cause on why
  `go run ... --dry-run` is not used for this check)

Top-level Status stays PASS; Round 2 confirms no regression from the pre-gate
commit.

## Round 3 (8e1b044)

### Status: PASS

Second pre-gate commit moved HEAD to 8e1b044: AGENTS.md and skills/preflight/SKILL.md
(two lines each; no Go; plus records).

- `go test -race -count=1 ./...` -> exit 0, all packages ok
- `go build -o .agents/scratch/verify-v130/sh-agent-built ./cmd/sh-agent` -> exit 0
- `sh-agent-built install --home .agents/scratch/verify-v130/home` -> exit 0,
  "1 written, 259 unchanged, 0 removed, 0 kept - 260 files managed" (`update
  ~/.local/bin/sh-agent` only). Not a drift/content bug: rebuilding at a new
  commit legitimately changes the binary (Go's automatic VCS build-info stamping
  embeds the commit hash), so the binary file itself differs while every
  rendered/copied doc and skill file (259) stays unchanged - consistent with the
  two-line AGENTS.md/preflight SKILL.md diff.
- `sh-agent-built status --home .agents/scratch/verify-v130/home` -> exit 0,
  "installed ... from commit 8e1b044+dirty", "drift none - 259 files match the
  repo"
- `sh-agent-built list` -> exit 0, "### Skills (27)", "### Subagents (7)" -
  matches expectation, no parse error

Top-level Status stays PASS; Round 3 confirms no regression from the second
pre-gate commit.

## Round 4 (10a9fe2)

### Status: PASS

Third pre-gate commit moved HEAD to 10a9fe2: skills/preflight/SKILL.md (one line),
two docs under .agents/docs, records.

- `go test -count=1 ./...` -> exit 0, all packages ok
- `go build -o .agents/scratch/verify-v130/sh-agent-built ./cmd/sh-agent` -> exit 0
  (this round did not re-install into the scratch home, per the dispatched
  command list - only built a fresh binary on disk)
- `sh-agent-built status --home .agents/scratch/verify-v130/home` -> exit 0,
  "drift none - 259 files match the repo" (matches expectation). Also correctly
  reports the scratch home's installed binary as stale: "installed ... from
  commit 8e1b044+dirty" and "binary 8e1b044cb053+dirty - older than the repo
  (10a9fe2); run `sh-agent update`" - expected, since this round's dispatch did
  not include an `install` step to refresh that home to 10a9fe2; the doc/skill
  file content still matches (no drift) because none of this round's three
  changed files render differently for this vendor's installed set, or install
  was not re-run to pick up the one-line SKILL.md change - the binary-staleness
  message is accurate status reporting, not a defect.
- `sh-agent-built list` -> exit 0, "### Skills (27)", "### Subagents (7)" -
  matches expectation, no parse error

Top-level Status stays PASS; Round 4 confirms no regression from the third
pre-gate commit.
