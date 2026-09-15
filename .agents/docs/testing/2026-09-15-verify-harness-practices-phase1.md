---
owner: verifier
date: 2026-09-15
---
# Verification — harness practices Phase 1   (2026-09-15)

Range reviewed: d1ff612..working tree (live working tree, nothing committed) — `git diff d1ff612 --stat` plus the untracked files `git status --short` lists.

Tracked, modified (9 files, 134 insertions, 26 deletions):
`.agents/docs/bugs/BUG-018-bug-skill-missing-autonomy-clause.md`, `.agents/docs/bugs/INDEX.md`,
`.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`, `.agents/docs/testing/sh-agent/guide.md`,
`AGENTS.md`, `README.md`, `skills/bug/SKILL.md`, `skills/retro/SKILL.md`, `skills/review/SKILL.md`

Untracked: `.agents/docs/bugs/BUG-019` through `BUG-023`, and seven `.agents/docs/reviews/2026-09-15-review-sweep-*.md` files.

## Status: PASS

Scope: Phase 1 of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (boxes 1.1, 1.3, 1.4, 1.5) plus the guide's §1 (Go checks), §2 (harness self-check), §3 (scratch-home install/drift). §3a and §4 out of scope for this run, per dispatch.

## Commands

Guide §1 — Go checks (each bare, own exit status):
- `gofmt -l .` → prints nothing, exit 0
- `go vet ./...` → prints nothing, exit 0
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/cmd/sh-agent 1.247s`; `ok github.com/kesokaj/styrman-harness/internal/harness 1.360s`, exit 0

Guide §2 — harness self-check:
- `go run ./cmd/sh-agent list` → exit 0; header lines read "### Skills (27)" and "### Subagents (7)", table lists 27 skill rows and 7 agent rows, all parsed
- `go run ./cmd/sh-agent render | wc -l` → `93` (matches the expected 93)

Guide §3 — scratch-home install and drift check (new dir under `.agents/scratch/`, never `/tmp` or real `$HOME`):
- `rm -rf .agents/scratch/verify-phase1-home && mkdir -p .agents/scratch/verify-phase1-home`
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase1-home` → "building sh-agent from /home/kesokaj/workspace/styrman-harness"; "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; note about PATH only, exit 0
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase1-home` → `repo ... main @ d1ff612, uncommitted changes`; `binary d1ff61239f77+dirty`; `repo has 27 skills, 7 agents, 0 commands, 0 rules`; `drift none — 259 files match the repo`, exit 0

Guide §3a — not applicable. `git diff d1ff612 --stat -- cmd internal` → no output (no file under `cmd/` or `internal/` changed in this range). Confirmed skipped per guide's rule ("skip only when no `.go` file is in the diff — say so").

## Plan boxes (Phase 1, exact commands)

**1.1** I-1a
- `grep -c 'narrowed: <old file:line>' skills/review/SKILL.md` → `1` (expect: `1`) — matches
- positive control: `git show c192622:skills/review/SKILL.md | grep -c 'narrowed:'` → `0` (expect: `0`) — matches

**1.3** I-2a, then I-6
- `grep -c "whatever autonomous instruction is in effect" skills/bug/SKILL.md` → `1` (expect: `1`) — matches
- `grep '^| BUG-018' .agents/docs/bugs/INDEX.md` → row's status cell reads `fixed`: `| BUG-018 | 2026-09-15 | skills/bug (gates) | fixed | bug skill's missing-repro gate has no "not skipped on an autonomous instruction" clause; feature skill kept it | a trim that moves instructions lists each removed rule's new home or why it was dropped; reviewer checks that list against the diff |` (expect: `fixed`) — matches

**1.4** I-3a
- `grep -c 'dropped: <why>' skills/retro/SKILL.md` → `1` (expect: `1`) — matches

**1.5** I-4a and I-5a-e through `documentation`
- `git grep -nE '27 skills|seven subagents|7 agents parse' -- README.md AGENTS.md .agents/docs/testing/sh-agent/guide.md` → no output, exit 1 (grep's normal "no match" exit) (expect: prints nothing) — matches
  - positive control against `c192622`: prints 5 lines (guide.md:22, AGENTS.md:7, README.md:3, README.md:86, README.md:87) (expect: prints lines) — matches
- `grep -c 'searched at \`skills/\` because' AGENTS.md` → `1` (expect: `1`) — matches

All four Phase 1 verify boxes and both positive controls behave exactly as their "— expect:" clauses specify.

## Root cause

N/A — all checks PASS, nothing failed.

## Coverage of changed code

Not applicable in the code-coverage sense — this range is documentation/prose (skills, plan, bug records, reviews, AGENTS.md, README.md) with no `.go` file changed (confirmed by the empty `git diff d1ff612 --stat -- cmd internal`). §1 Go checks ran regardless per the guide's baseline tier and passed clean.

## Not examined

- **§3a (plain-copy install)** — not run. Confirmed inapplicable: no file under `cmd/` or `internal/` is in this diff range.
- **§4 / plan's Phase 1 checkpoint's sandbox CLI start** — deferred to the Phase 5 checkpoint per dispatch; not run here.
- **Box 1.2** (docs-lens control, scratch reviewer dispatch) — not part of this dispatch's scope (only 1.1, 1.3, 1.4, 1.5 named); not run.
- **Box 1.6** (one-time sweep) — already marked in-progress in the plan itself with its own evidence trail (`reviews/2026-09-15-review-sweep-merged.md`); not re-verified here.
- **The Phase 1 Checkpoint box's own full expect list** (e.g. "a verifier log over the phase range with Status PASS whose Not examined section names the sandbox CLI start deferred to the Phase 5 checkpoint") — this log serves that purpose but the box itself was not ticked (plan-file edits are outside a verifier's remit; that is the main agent's or `changelog`'s job).
- **Reviews of the untracked BUG-019..023 and review-sweep files' content** — read only their existence via `git status`; their internal correctness was not assessed (out of scope for a Phase 1 command-verification run).

## Environment

- OS: Linux 6.12.107+deb13-amd64
- Go: as resolved by `go run` / `go vet` / `go test` in this checkout (no explicit version pin found in this run; toolchain used is whatever `go` on PATH resolved)
- No Docker/Colima involved in this range (no Dockerfile changed).
