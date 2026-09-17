---
owner: verifier
date: 2026-09-17
---
# Verification — retro-context   (2026-09-17)

Range reviewed: 24c524e..working tree (uncommitted edits to `context/AGENTS.md`, `agents/changelog.md`; untracked `.agents/docs/retros/2026-09-17-retro.md`)

## Status: PASS

Scope: all four steps named in the dispatch ran and passed. Sandbox CLI headless starts and real-home smoke checks were not run (see Not examined).

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-retro-context/go-test-1.txt`; counts: pass 124, fail 0, skip 1 at 24c524e+dirty for `go test -race -v -count=1 ./...`. The single skip is `TestFindRepo` ("working directory happens to be inside a harness repo"), a pre-existing environment-dependent skip unrelated to this change. `TestNoTmpInHarnessProse` ran and passed (line 133).
  vs previous: `.agents/docs/testing/2026-09-17-verify-delivery-rework-done.md` (counts: pass 118, fail 0, skip ? at 24c524e for the same command) → pass 118→124, skip unchanged (both runs show the same TestFindRepo skip; the prior log's grep only counted top-level `^--- PASS` lines, undercounting subtests — this run counts all `--- PASS:`/`--- FAIL:`/`--- SKIP:` lines including indented subtests, which explains the rise). No fall in pass, no unexplained rise in skip.
- `go run ./cmd/sh-agent install --home .agents/scratch/retro-home` → exit 0; "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed"; settings.json created with attribution off as expected.
- `go run ./cmd/sh-agent status --home .agents/scratch/retro-home` → exit 0; `binary   24c524eac707+dirty` (no warning); `drift    none — 259 files match the repo` (259 vs 260 written is the settings.json file, which is excluded from drift comparison by design, not an anomaly).
- Content grep of the installed copies:
  - `BUG-058`, `tool-call budget`, `BUG-061`, `must not drift` — each found (grep -F) in `context/AGENTS.md`, `.agents/scratch/retro-home/.claude/AGENTS.md`, and `.agents/scratch/retro-home/.gemini/config/rules/AGENTS.md`.
  - Six literal `A` characters (`AAAAAA`, grep -F) — not found in any of the three files (expected: no such artifact should be present; absence is the pass condition).
  - `Every factual claim traces to the dispatch or the diff` — found in `agents/changelog.md`, `.agents/scratch/retro-home/.claude/agents/changelog.md`, and `.agents/scratch/retro-home/.gemini/config/agents/changelog/agent.md`.
- `go run ./cmd/sh-agent list` → exit 0; output saved to `.agents/scratch/verify-retro-context/list.txt`; lists 27 skills and all 7 subagents including `changelog`, parses cleanly with no errors.

## Root cause
N/A (PASS; no failures found).

## Coverage of changed code
- `context/AGENTS.md` — prose-only global context change; covered by `TestNoTmpInHarnessProse` (ran, passed) and by the install/content-grep steps above, which confirm the new sentences (BUG-058, tool-call budget, BUG-061, must not drift) deploy verbatim to both vendor copies.
- `agents/changelog.md` — prose-only subagent instruction change; covered by the install + content-grep step confirming "Every factual claim traces to the dispatch or the diff" deploys to both the Claude Code copy and the derived Antigravity `agent.md`.
- `.agents/docs/retros/2026-09-17-retro.md` (untracked) — a docs/trail artifact, not source; not applicable to test coverage.
- No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) were added by this change; the diff is prose-only in `context/AGENTS.md` and `agents/changelog.md`.

## Not examined
- One headless start of each CLI (Claude Code, Antigravity) in the sandbox home, per `.agents/docs/testing/2026-09-03-verify-sandbox.md` — not run because the logged-in sandbox home no longer exists on this box, as recorded in the Phase-3 plan row.
- The real-home smoke checks (context questions, one skill load, one agent dispatch) from `.agents/docs/testing/2026-09-03-verify-real-home.md` — not run; this verification only used the scratch home `.agents/scratch/retro-home`, never the operator's real `$HOME`, per the project rule against installing there during development.
- `/review --docs` over this range — not part of this dispatch; testing only.

## Environment
- OS: Linux 6.12.107+deb13-amd64.
- Go: user-level SDK ahead of `/usr/bin/go` per project AGENTS.md (no `toolchain` directive in go.mod; `GOTOOLCHAIN=auto` not exercised to pin).
- No Docker/Colima involved (no Dockerfile changed).
- Working tree: dirty (uncommitted `context/AGENTS.md`, `agents/changelog.md`; untracked `.agents/docs/retros/2026-09-17-retro.md`), on top of commit 24c524e.
