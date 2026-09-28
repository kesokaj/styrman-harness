---
owner: verifier
date: 2026-09-28
---
# Verification — open-items (2026-09-28)

Range reviewed: HEAD (c86d3dd)..working tree (uncommitted)

## Status: PASS

Scope: complete — every command the dispatch named ran; `.agents/docs/testing/guide.md` exists.

## Commands
- `git status --short --untracked-files=all` → exactly three modified files: `.agents/docs/bugs/INDEX.md`, `AGENTS.md`, `skills/cloud-run-deploy/SKILL.md` — matches the expected change set, nothing else.
- `git diff --numstat` → `1 1 .agents/docs/bugs/INDEX.md`, `2 2 AGENTS.md`, `1 1 skills/cloud-run-deploy/SKILL.md` — all numeric columns, no binary (`-`/`-`) files.
- `gofmt -l cmd internal` → exit 0, no output (no files need formatting).
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0. Output saved to `.agents/scratch/go-test-open-items.txt`.
  counts: pass 251, fail 0, skip 0 at c86d3dd+dirty for `go test -race -v ./...`; vs previous: `.agents/docs/testing/2026-09-28-verify-research-okf.md` records `go test -race ./...` (no `-v`, package-level only, "no per-test -v breakdown captured"), so no directly comparable per-test count exists in `testing/` — vs previous: none (comparable).
  `TestRenderedContextFitsAntigravityBlock` reports "Antigravity rules block is 20329 bytes ... budget 22000" — pass, and consistent with the AGENTS.md shortening in this change set (smaller, not larger).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-items` → exit 0, "258 written, 0 unchanged, 0 removed, 0 kept — 258 files managed". Never touched the real `$HOME`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-items` → exit 0, output includes `drift      none — 257 files match the repo` (the settings.json note above it is expected installer behaviour, unrelated to drift).
- Content check: `.agents/scratch/verify-items/.claude/skills/cloud-run-deploy/SKILL.md` line 37 reads `6. **Service account:** the service account the service will run as exists ...` — confirms the installed copy carries the renumbered item.
- `go run ./cmd/sh-agent list` → exit 0; `cloud-run-deploy | auto | End-to-end release workflow for Google Cloud Run: ...` still printed — skill still parses.
- `git diff -- skills/cloud-run-deploy/SKILL.md` → confirms the change is body-only: the unnumbered pre-flight bullet (`- The service account ...`) became `6. **Service account:** ...`; frontmatter `description:` (line 3) untouched by the diff.
- `git diff -- AGENTS.md .agents/docs/bugs/INDEX.md` → confirms: AGENTS.md rule 6 dropped the "v1.3.0 gate spent three review rounds ..." sentence, replaced with a pointer to `.agents/docs/retros/2026-09-18-retro-2.md`; the toolchain paragraph in Setup commands was shortened. bugs/INDEX.md: BUG-011 row removed from the open table and added to the Graduated table pointing at `context/AGENTS.md` § Hard rules.

## Root cause
N/A — no failures found.

## Coverage of changed code
Docs-only change (SKILL.md renumbering, AGENTS.md prose shortening, bugs INDEX.md table move). No Go source changed, so no `_test.go` coverage question applies. Coverage verified through the project's documented docs-change procedure: install to scratch home + status drift-none + content grep + `sh-agent list` parse check, per `AGENTS.md` § Testing and `.agents/docs/testing/guide.md`.

## Not examined
- No CLI sandbox start (Claude Code / Antigravity process launch). Not needed: the skill-file change is body text only (renumbering an existing bullet into the numbered list), and `git diff` confirms the `description:` frontmatter line is untouched — the guide's sandbox-start tier targets frontmatter/parse-shape changes, not renumbering within an already-valid body.
- No real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`). Not run — the real `$HOME` was never touched per the dispatch's explicit instruction; scratch-home install/status/content checks substitute per the docs-change procedure.
- No behavioural agent run (`claude -p --agents ...`). Not applicable — no `agents/*.md` file changed in this change set.
- `govulncheck ./...` was not run (not requested by the dispatch and no dependency/toolchain change is in scope).

## Environment
Linux (Debian trixie), user-level go1.27.1 SDK ahead of `/usr/bin/go` per `.agents/docs/deploy/local.md`. `go test` run with `-race -v -count` not explicitly added (dispatch specified `go test -race ./...` bare per AGENTS.md Setup commands, run with added `-v` for per-test counts as required by this agent's protocol; no `-run` narrowing used, full package run).
