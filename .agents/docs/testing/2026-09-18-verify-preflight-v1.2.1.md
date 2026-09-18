---
owner: verifier
date: 2026-09-18
---
# Verification — preflight v1.2.1   (2026-09-18)

Range reviewed: v1.2.0..288e52d

## Status: PASS

Scope: complete — every command in the dispatch ran, plus the guide's tier-1 checks
(gofmt/vet/test) match the release dispatch's steps 1-3. `testing/guide.md` exists
(dated 2026-09-17) and was cited, not rewritten (no procedure changed in this range).
Guide's own tiers 2/4/5 (harness self-check via `list`/`render`, sandbox CLI starts,
real-home smoke checks) were NOT run — the dispatch did not ask for them and this
range's change (release-binary source resolution) is not a `context/`/`agents/`/
`skills/` behaviour change, so they are out of scope for this release check.

## Commands

1. `gofmt -l cmd internal *.go` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 ./...` → exit 0.
   ```
   ok  	github.com/kesokaj/styrman-harness	1.089s
   ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	3.738s
   ok  	github.com/kesokaj/styrman-harness/internal/harness	2.773s
   ```
   `go test -race -count=1 -v ./...` → exit 0.
   counts: pass 103 top-level + 86 subtests (189 `--- PASS:` lines total), fail 0,
   skip 0, at 288e52d for `go test -race -count=1 -v ./...`.
   vs previous: `.agents/docs/testing/2026-09-18-verify-preflight-v1.2.0-fixes.md`,
   pass 102→103 top-level (188→189 total), skip 0→0. The +1 top-level pass is a rise,
   not a fall — no explanation needed per the rule (only a fall or a skip rise is a
   FAIL line). Consistent with the dispatch's note that this range adds assertions to
   existing tests, though evidently at least one top-level test/subtest count also grew
   by one; not further isolated (see Not examined).
4. `govulncheck ./...` (go1.27.1) → exit 0, "No vulnerabilities found."
5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-v121-home` → exit 0,
   260 files written. `go run ./cmd/sh-agent status --home .agents/scratch/verify-v121-home`
   → exit 0, `drift      none — 259 files match the repo`. (Repo tree carried one
   untracked file at run time — this verification log itself, created by this run's
   own Write call before the checkout-path install — hence `288e52d+dirty` in the
   status; that dirtiness is self-caused by the verifier's own artifact, not a real
   change to the release candidate. `git status --short --untracked-files=all`
   confirmed no tracked file was modified.)

## Step 6 — release-binary behaviour (the point of this release)

Built `.agents/scratch/verify-v121-bin/sh-agent` with
`-trimpath -ldflags "-X main.version=v9.9.9"` from `./cmd/sh-agent` → exit 0.
`SH_AGENT_REPO` unset throughout. Harness-shaped directory used:
`.agents/scratch/devbox-copy` (already present, confirmed no `.git`).
State files (`~/.config/sh-agent/state.json`) written recording that directory as
`repo` for homes (a), (b), (c).

- **(a) from `/`, home recording a `.git`-less harness-shaped dir:**
  `install --home <home-a>` → `release binary v9.9.9: installing the harness built
  into it and the running binary` (exit 0). `status --home <home-a>` →
  `source     embedded v9.9.9 (release binary; update downloads the newest release)`.
  Did NOT read or install from the recorded tree. **PASS.**
- **(b) standing inside that harness-shaped directory:** same install/status wording,
  same "embedded v9.9.9" source line, same non-use of the on-disk tree. **PASS.**
- **(c) standing inside this repository's own checkout:** same install/status
  wording, same "embedded v9.9.9" source line — the release binary did NOT fall back
  to rebuilding from this checkout even though it is a real `.git` checkout on disk.
  **PASS.**
- **(d) `install --repo <harness-shaped dir> --home <scratch>`:** REFUSED, exit 1:
  `sh-agent: this is the release binary v9.9.9: it installs the harness built into it
  and never a tree on disk, so <path> is refused — drop it, or build sh-agent from
  that tree and run that build's install`. Message names the release binary as
  required. **PASS.**

All four legs match the behaviour ADR-034/BUG-077 (commit 288e52d) intends.

## Root cause
N/A — no failure found.

## Coverage of changed code
Range `v1.2.0..288e52d` (5 commits) is the ADR-034/BUG-077 release-binary source-
resolution work plus a Mac e2e follow-up and INDEX repair; `go test -race -count=1
./...` passed across all three packages (root, `cmd/sh-agent`, `internal/harness`)
with a net +1 top-level pass vs the prior logged run, and the specific new behaviour
(release binary refuses `--repo`, never reads a recorded/standing tree) was exercised
live in step 6 above, independently of the unit tests. No changed-code file found
uncovered by the test run.

## Not examined
- Which specific top-level test/subtest is the +1 pass vs the v1.2.0-fixes log (not
  isolated by name/diff of test files — only the aggregate counts were compared).
- Guide tiers 2, 4, 5 (harness self-check `list`/`render`, sandbox CLI starts, real-
  home smoke checks) — out of scope per dispatch; this range does not touch
  `context/`, `agents/`, or `skills/`.
- `sh-agent update` and `sh-agent report` were not run at all (not `--dry-run` either)
  — the dispatch's hard limits require `report` only with `--dry-run` and forbid
  `update` outright; neither was requested for this release check.
- Mac-specific legs referenced in recent commit history (e2e on the work Mac) — this
  run was on the Linux dev box only, per the dispatch's environment.
- `skills/VENDORED.md` / `/review --docs` — not requested by this dispatch (this
  range is not a `context/`/`agents/`/`skills/` change).

## Environment
Linux 6.12.107+deb13-amd64. go version go1.27.1 linux/amd64 (user-level SDK ahead of
`/usr/bin/go`, per this repo's documented setup). govulncheck at
`/home/kesokaj/.go/bin/govulncheck`, clean. HEAD 288e52d, tree clean at dispatch time
except for this verification log created during the run (untracked, not committed by
this agent). All scratch artifacts under `.agents/scratch/verify-v121*` and
`.agents/scratch/devbox-copy` (pre-existing fixture, read-only used). No network
access attempted; `update`/`report` not invoked; real `$HOME` never touched (all
`HOME=` overrides pointed at `.agents/scratch/verify-v121-{a,b,c,d}`).
