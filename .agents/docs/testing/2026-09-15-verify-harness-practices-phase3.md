---
owner: verifier
date: 2026-09-15
---
# Verification — Phase 3, harness practices G01-G26

Range reviewed: 32c8427..working tree (nothing from Phase 3 committed yet)

`git diff 32c8427 --stat`:
```
 .agents/docs/plans/2026-09-15-harness-practices-g01-g26.md | 4 ++--
 context/AGENTS.md                                          | 2 +-
 skills/golang/SKILL.md                                     | 2 +-
 skills/golang/references/golang-ecosystem.md               | 2 +-
 skills/preflight/SKILL.md                                  | 2 +-
 skills/review-standards/SKILL.md                           | 5 +++++
 skills/review-standards/references/go-review.md            | 1 +
 skills/review-standards/references/typescript-review.md    | 3 ++-
 skills/review/SKILL.md                                     | 4 ++--
 9 files changed, 16 insertions(+), 9 deletions(-)
```
No untracked files (`git status --short --untracked-files=all` shows only the 9 modified files above).

## Status: PASS

Scope: guide §1-3 run; §3a checked and confirmed not applicable; boxes 3.1 and 3.2 verify commands run exactly as written, including 3.1's positive control over c192622. §4 (sandbox CLI starts) deferred to Phase 5 per instruction — not run here.

## Commands

### §1 Go checks
- `.go` files in range: **none** (`git diff 32c8427 --stat -- '*.go'` printed nothing) — Go checks run anyway for completeness.
- `gofmt -l .` → no output, exit=0
- `go vet ./...` → no output, exit=0
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/cmd/sh-agent 1.219s`, `ok github.com/kesokaj/styrman-harness/internal/harness 1.366s`, exit=0

### §2 Harness self-check
- `go run ./cmd/sh-agent list` → exit=0, printed 27 skills / 7 agents, all rows parsed (no parse errors)
- `go run ./cmd/sh-agent render | wc -l` → `93` — matches expect exactly

### §3a Plain-copy install
- `git diff 32c8427 --stat -- cmd internal` → empty. No change to `cmd/` or `internal/` in this range, so §3a does not apply. Skipped.

### §3 Scratch-home install and drift check
- Confirmed `.agents/scratch/verify-phase3-home` did not exist before the run (`ls` → "No such file or directory").
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase3-home` → "building sh-agent from ...", "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed", exit=0
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase3-home` → exit=0, key line: `drift none — 259 files match the repo` — matches "drift none" expect. (Note: install reported 260 files managed vs status's 259 matched; this 1-file gap is the manifest/state-file bookkeeping itself, not content drift — the drift line explicitly reads "none".)
- No agent definition changed in this range (only `context/AGENTS.md` and skill files changed; `agents/*.md` untouched), so the Records-section cross-check in the guide's §3 does not apply, per the dispatch instruction.

### Box 3.1 (I-12a-f, I-13a-b) — verify: integration, grep -cF
Current tree, each expected `1`:
- "that section's fixed lines only" in `skills/review/SKILL.md` → 1
- `in flight: <who>` in `skills/review/SKILL.md` → 1
- "unless the user picks one more at the cap" in `skills/review/SKILL.md` → 1
- "or a fix introduced that round's new defect" in `skills/review/SKILL.md` → 1
- "fix the rest yourself" in `skills/review/SKILL.md` → 1
- "which the main agent fixes" in `context/AGENTS.md` → 1
- "never a command that exited non-zero" in `skills/review/SKILL.md` → 1
- `a reason as \`/review\` Step 4 defines it` in `skills/preflight/SKILL.md` → 1

Positive control, same 8 greps run against `c192622` via `git show c192622:<path> | grep -cF ...`, each expected `0`:
→ 0, 0, 0, 0, 0, 0, 0, 0 (all eight)

All results match the box's expect exactly.

### Box 3.2 (I-14a-c, I-15a-f) — verify: integration, grep -cF
- "Leftover debug output" in `skills/review-standards/SKILL.md` → 1
- "TOCTOU race" in `skills/review-standards/SKILL.md` → 1
- "generated counterpart" in `skills/review-standards/SKILL.md` → 1
- "assertion changed in an existing test" in `skills/review-standards/SKILL.md` → 1
- "describing-text rule" in `skills/review-standards/SKILL.md` → 1
- "lint:ignore" in `skills/review-standards/references/go-review.md` → 1
- "describe.only" in `skills/review-standards/references/typescript-review.md` → 1
- "without a reason" in `skills/review-standards/references/typescript-review.md` → 1
- "without a linked issue" in `skills/review-standards/references/typescript-review.md` → 0
- "demonstrated false positive" in `skills/golang/SKILL.md` → 1
- "A Go project that runs golangci-lint" in `skills/golang/references/golang-ecosystem.md` → 1
- `sed -n 323p skills/golang/references/golang-ecosystem.md | grep -cF nolintlint` → 1 (line 323 read: "...small set (`govet`, `staticcheck`, `errcheck`, `gosec`, `unused`, `misspell`, `nolintlint`)...")

All results match the box's expect exactly: every grep `1` except "without a linked issue" (`0`).

## Coverage of changed code

No `.go` files in range — this is a docs/skills-only change. Every changed file (`context/AGENTS.md`, `skills/golang/SKILL.md`, `skills/golang/references/golang-ecosystem.md`, `skills/preflight/SKILL.md`, `skills/review-standards/SKILL.md`, `skills/review-standards/references/go-review.md`, `skills/review-standards/references/typescript-review.md`, `skills/review/SKILL.md`) is covered by the box 3.1/3.2 `grep -cF` checks above, which target the exact added phrases. `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` changed too (box ticks) — not a "living doc" needing grep coverage under this guide.

## Not examined

- **§4 (sandbox CLI starts)**: explicitly deferred to the plan's Phase 5 checkpoint per dispatch instruction; not run in this pass.
- **§5 (real-home smoke checks)**: not applicable — no install into a real home was performed (scratch home only).
- **`/review --docs`** over this range: not run by the verifier; this repo's testing guide notes a docs-touching change also gets a documentation-agent review, which is outside verifier scope.
- Box 3.1/3.2 text was checked only via the literal grep commands the plan specifies (exact strings, exact files); broader prose quality/consistency of the surrounding sections was not read line-by-line.

## Environment

Linux 6.12.107+deb13-amd64, repo at commit 32c8427 with the 9 uncommitted modifications listed above. Go toolchain per repo `go.mod` (module `github.com/kesokaj/styrman-harness`). No Docker/Colima involved (no Dockerfile in range).
