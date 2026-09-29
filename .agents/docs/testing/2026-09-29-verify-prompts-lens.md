---
owner: verifier
date: 2026-09-29
---
# Verification — prompts-lens   (2026-09-29)

Range reviewed: working tree (uncommitted) — .agents/docs/testing/skill-matrix.md, skills/VENDORED.md, skills/review-standards/SKILL.md, skills/review/SKILL.md, skills/review-standards/references/prompt-audit/LICENSE.txt (new, staged), skills/review-standards/references/prompt-audit/prompt-audit.md (new, staged). No Go source in the diff; the change adds a vendored prompt-audit reference file under `review-standards` and a `--prompts` lens/flag to the `review` skill, plus a matching `skill-matrix.md` row and `VENDORED.md` provenance row.

Guide used: `.agents/docs/testing/guide.md` (already present; not written this run) — § 1 Go checks, § 2 Harness self-check, § 3 Scratch-home install and drift check.

## Status: PASS

Scope: complete for the tiers this dispatch named (§1, §2, install/status/drift, file placement, checksum, provenance-line match, suppression scan). The behavioural run of the new `--prompts` lens (three `claude -p --agents` runs) is run separately by the main agent in `.agents/scratch/prompts-probe/` and is explicitly out of scope for this dispatch — see "Not exercised here".

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean)
- `go vet ./...` → exit 0, no output
- `go test -v -count=1 -race ./...` → exit 0. Output saved to `.agents/scratch/verify-prompts-home2/go-test-1.txt`.
  counts: pass 171, fail 0, skip 0 at 6b1b486+dirty for `go test -race -v -count=1 ./...`
  `ok github.com/kesokaj/styrman-harness 1.097s`, `ok .../cmd/sh-agent 13.306s`, `ok .../internal/harness 3.133s`
  vs previous: `.agents/docs/testing/2026-09-29-verify-v1.11.0-gate-r5.md` (counts: pass 171, fail 0, skip 0 at 814c257), pass 171→171, skip 0→0 — no change, no explanation needed.
- `go run ./cmd/sh-agent list` → exit 0. `review` entry printed with its updated one-line description truncated in the table; no parse error; 27 skills / 7 agents counts unchanged.
- `go run ./cmd/sh-agent render | wc -l` → 90 (must stay under 100 per guide § 2) → PASS. (No `context/` or `rules/*.md` file changed, so `TestRenderedContextFitsAntigravityBlock` is unaffected — already covered by the `go test ./...` run above, which passed.)
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-prompts-home` → exit 0. Rebuilt from this checkout (`6b1b486+dirty`); "262 written, 0 unchanged, 0 removed, 0 kept — 262 files managed", including both new lines:
  `add ~/.claude/skills/review-standards/references/prompt-audit/LICENSE.txt`
  `add ~/.claude/skills/review-standards/references/prompt-audit/prompt-audit.md`
  and the Antigravity mirrors under `~/.gemini/config/skills/review-standards/references/prompt-audit/`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-prompts-home` → exit 0. `drift none — 261 files match the repo` (261 tracked files vs. 262 "managed" in the install log — the 262nd is the installed `sh-agent` binary itself under `~/.local/bin`, not a repo-tracked file, so this is not a discrepancy).
- `ls .agents/scratch/verify-prompts-home/.claude/skills/review-standards/references/prompt-audit/` → `LICENSE.txt`, `prompt-audit.md` — both present.
- `ls .agents/scratch/verify-prompts-home/.gemini/config/skills/review-standards/references/prompt-audit/` → `LICENSE.txt`, `prompt-audit.md` — both present.
- `sha256sum skills/review-standards/references/prompt-audit/prompt-audit.md` → `c33200da181e10f1f31bd4625059cd30bc6922940aac49017d8ab0c2721878f2` — matches the `c33200da…1878f2` digest claimed in `skills/VENDORED.md` line 13.
- `git diff --unified=0 -- skills/review/SKILL.md skills/review-standards/SKILL.md skills/VENDORED.md .agents/docs/testing/skill-matrix.md | grep -iE '@ts-ignore|@ts-expect-error|eslint-disable|nolint|lint:ignore|t\.Skip|go:build ignore|\.only\(|\.skip\('` → no match (exit 1) — no suppressions or skips added by this change.

## Root cause
N/A — no failures found.

## Coverage of changed code
- `skills/review-standards/references/prompt-audit/prompt-audit.md` (new, vendored verbatim) — not Go/TS source; no unit test applies. Coverage is the install-tier file-placement check above (present under both `.claude/` and `.gemini/config/` scratch homes) plus the checksum match against `skills/VENDORED.md`'s claimed digest, both of which passed.
- `skills/review-standards/references/prompt-audit/LICENSE.txt` (new, vendored) — same install-tier placement check; present in both scratch-home mirrors.
- `skills/review-standards/SKILL.md` — covered by `go run ./cmd/sh-agent list` (parses without error) and the install/drift tier (file still deployed correctly with the new reference-file link added).
- `skills/review/SKILL.md` — the new `--prompts` flag, argument-hint entry, example line, Step 1 flag list, Step 3 dispatch-table row, and Step 5 `--fix` handling are prose/skill-definition changes with no unit-test surface in this repo (`go test ./...` does not parse skill prose). The only test this repo has for skill *behaviour* is the sandbox behavioural-run tier (§ Testing in `AGENTS.md`) and the new `review-prompts` row added to `.agents/docs/testing/skill-matrix.md` in this same diff — that behavioural run is explicitly being executed by the main agent separately in `.agents/scratch/prompts-probe/`, per this dispatch's instructions, and was NOT run here. Until that run completes and is recorded, the new `--prompts` lens itself is UNVERIFIED behaviourally; this log only confirms the static/install/deploy tiers pass.
- `.agents/docs/testing/skill-matrix.md` — new `review-prompts` row is a test-matrix definition, not itself executable by this dispatch; its own "runs" happen when the sandbox matrix is executed (not part of this dispatch).
- `skills/VENDORED.md` — new provenance row checked against the actual file's sha256 above (match).

## Not exercised here
- The behavioural run of the new `--prompts` lens (three `claude -p --agents` runs against a seeded fixture, per `AGENTS.md` § Testing and the new `review-prompts` row in `skill-matrix.md`) — this is being run by the main agent in `.agents/scratch/prompts-probe/`, explicitly excluded from this dispatch. Its result is not reflected in this log's Status.
- `govulncheck ./...` — not named in this dispatch's instructions (guide.md §§1–3 only); not run.
- § 3a (plain-copy install) and § 3b (`report --dry-run`) of the guide — not applicable: this change touches no `IsCheckout`/rebuild logic and no `report.go`/`main.go` report wiring.
- `/review --docs` over this range — not part of a verifier dispatch; out of scope here.

## Environment
Linux 6.12.107+deb13-amd64 (`uname -r` 6.12.107+deb13-amd64). Go toolchain as resolved by `go run`/`go test` in PATH (per project `AGENTS.md`, the Linux box runs a user-level go1.27.1 SDK ahead of `/usr/bin/go`). Repo HEAD `6b1b486`, working tree dirty (the change set under review). No Docker/Colima involved (no Dockerfile changed).
