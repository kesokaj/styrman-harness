---
owner: verifier
date: 2026-09-29
---
# Verification — prompts-lens-r2   (2026-09-29)

Range reviewed: working tree (uncommitted) — .agents/docs/testing/skill-matrix.md, skills/VENDORED.md, skills/review-standards/SKILL.md, skills/review/SKILL.md, skills/review-standards/references/prompt-audit/LICENSE.txt (new, staged), skills/review-standards/references/prompt-audit/prompt-audit.md (new, staged). Round 2 re-verification after three review findings were fixed with prose edits: skills/review/SKILL.md (`--prompts` lens row, one clause under Step 5), skills/review-standards/SKILL.md (reference bullet moved under its own heading), .agents/docs/testing/skill-matrix.md (row edit). No Go source changed.

## Status: PASS

Scope: complete — every command the dispatch named ran, in order, per testing/guide.md §§ 1–2 plus the install tier; testing/guide.md exists (unchanged this round).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -race ./...` → exit 0
  - `ok  	github.com/kesokaj/styrman-harness	1.102s`
  - `ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	13.381s`
  - `ok  	github.com/kesokaj/styrman-harness/internal/harness	3.148s`
  - counts: not printed (run without `-v`; dispatch named the bare command, not `-v -count=1`, so package-level PASS is the granularity available). counts: not applicable in the pass/fail/skip line format — no per-test counts to compare against a previous log.
  - vs previous: `.agents/docs/testing/2026-09-29-verify-prompts-lens.md` — same status (round 1 also PASS, all three packages ok) at commit 6b1b486 in both rounds (working tree dirty in both; the round-1 log predates the three prose fixes).
- `go run ./cmd/sh-agent list` → exit 0; skills and agents tables render (27 skills incl. `review`, `review-standards`; 7 agents), no parse errors observed in output
- `go run ./cmd/sh-agent render | wc -l` → exit 0, 90 lines
- `rm -rf .agents/scratch/verify-prompts-home-r2 && go run ./cmd/sh-agent install --home .agents/scratch/verify-prompts-home-r2` → exit 0; `262 written, 0 unchanged, 0 removed, 0 kept — 262 files managed`; includes the two new prompt-audit reference files under both `.claude/skills/review-standards/references/prompt-audit/` and `~/.gemini/config/skills/review-standards/references/prompt-audit/` in the write log
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-prompts-home-r2` → exit 0; `drift      none — 261 files match the repo` (repo shown `uncommitted changes`, binary `6b1b486d2db2+dirty`, as expected on a dirty working tree)
- `grep -c "sub-agent roster check" .agents/scratch/verify-prompts-home-r2/.claude/skills/review/SKILL.md` → `1`, exit 0
- `grep -c "sub-agent roster check" .agents/scratch/verify-prompts-home-r2/.gemini/config/skills/review/SKILL.md` → `1`, exit 0
- `sha256sum skills/review-standards/references/prompt-audit/prompt-audit.md` → `c33200da181e10f1f31bd4625059cd30bc6922940aac49017d8ab0c2721878f2` — matches expected prefix `c33200da`
- `git diff -- skills/review/SKILL.md skills/review-standards/SKILL.md .agents/docs/testing/skill-matrix.md skills/VENDORED.md | grep -nE '^\+.*(nolint|lint:ignore|ts-ignore|ts-expect-error|eslint-disable|t\.Skip|SkipNow|go:build ignore|\.skip\(|\.only\()'` → exit 1 (no matches) — no suppressions, skips, or `.only(` markers added in the round-2 diff

## Root cause
N/A — no failures this round.

## Coverage of changed code
No Go source changed this round; all five changed/added files are prose or reference-doc content (skill markdown, a testing matrix row, a vendored license+reference pair, a vendoring ledger row). Coverage is judged by the harness's own doc-verification path (install + status + drift + content greps), all of which passed above. `LICENSE.txt` and `prompt-audit.md` are inert reference material (not executed), verified by content hash and by their presence in both installed trees.

## Not examined
- The behavioural run of the `--prompts` lens itself: out of scope for this dispatch — it is the main agent's run in `.agents/scratch/prompts-probe/`, referenced from `.agents/docs/testing/2026-09-29-prompts-lens-behavioural.md`, not re-run or re-read here.
- `govulncheck ./...` (named in AGENTS.md Setup Commands but not in this dispatch's command list) — not run this round.
- The two new review reports under `.agents/docs/reviews/2026-09-29-review-worktree-{code,docs}.md` — not read or re-verified; only the resulting SKILL.md/skill-matrix.md/VENDORED.md edits were checked.
- `/review --docs` re-run over this round's diff — not performed by this verifier; it is a separate gate per the project's own Definition of Done.

## Environment
OS: Linux 6.12.107+deb13-amd64 (Debian). Go: go1.27.1 linux/amd64 (user-level SDK ahead of /usr/bin/go, per AGENTS.md). Git HEAD: 6b1b486 on `main`, working tree dirty (the change set under review). No Docker/Colima involved (no Dockerfile changed).
