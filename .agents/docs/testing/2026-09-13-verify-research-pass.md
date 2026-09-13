---
owner: main
date: 2026-09-13
---
# Verification — research pass over the harness at b64ffe8   (2026-09-13)

Range reviewed: b64ffe8 (HEAD, clean tree; nothing changed by this pass)

## Status: PASS for the repo and the sandbox; one defect found in the real home's installed binary

Scope: unit tests, installer behaviour in a scratch home under `.agents/scratch/`, one headless start of each CLI from the sandbox home with HEAD installed. Not run: the full skill matrix (`skill-matrix.md`), the Mac, interactive starts.

## Commands

- `gofmt -l .` → nothing; `go vet ./...` → clean; `go test -race -count=1 ./...` → `ok internal/harness 1.168s`, `cmd/sh-agent [no test files]`; `go test -cover ./internal/harness/` → 78.3%.
- `go run ./cmd/sh-agent install --home .agents/scratch/home-test` → `256 written`; second run → `0 written, 256 unchanged`; `status` → `drift none — 255 files match the repo`, `vendor: all`.
- `install --vendor claude` → `127 removed`, `.gemini` gone, `status` reads the vendor from the manifest; `install --vendor all` → `127 written`.
- Tampered `~/.claude/skills/golang/SKILL.md` and removed `~/.gemini/config/agents/bugs/agent.md` → `status` exit 1 with `differs` and `missing`; `diff` showed the tampered line; `install` repaired both.
- Foreign `~/.claude/CLAUDE.md` → `keep` line, never overwritten; `status` note printed.
- `update --check` → `up to date with upstream (b64ffe8)`, exit 0.
- `uninstall --dry-run` → `would remove 255 files`; `uninstall` → `removed 255 files and the manifest`; `find` left nothing.
- Stale-binary repro: copied `~/.local/bin/sh-agent` (built 2026-09-03 21:05, `strings` shows no `target vendor`, `help` shows no `--vendor`) into the scratch home's `.local/bin` and ran *that* `install --home <scratch>` → `0 written, 256 unchanged`; the binary's md5 unchanged. An installed binary re-installs itself; only `go run ./cmd/sh-agent install` or `sh-agent update` rebuilds. Real home: `sh-agent version` → `(devel)`; `state.json` has no `vendor` key; content matches the repo (`drift none`), only the binary is old.
- Sandbox (`HOME=/home/kesokaj/sandbox/home`, both CLIs logged in since 2026-09-03): `install --home` → `24 written, 232 unchanged`, `drift none`. Claude Code 2.1.270 `-p` context questions, no tools: first heading, 28 skills / 7 subagents, the piped-verification hard rule and DoD step 5 quoted verbatim from the 2026-09-10/11 retro commits, `bug` named for a failing test. Antigravity 1.1.25 `-p --mode accept-edits`: identical answers; newest `cli-20260913_204140.log` → 0 `Failed to parse`.
- Antigravity real-home logs: newest is 2026-09-03, so the 2026-09-08..11 retro changes had not been started on this box under Antigravity until the sandbox run above.

## Findings

1. Installed binary stale and self-perpetuating (fix planned: `plans/2026-09-13-learning-loop-project-level.md` Phase 1).
2. `sh-agent version` prints `(devel)` for a `go run` build; the staleness was invisible.
3. `AGENTS.md`, `README.md` say `--home /tmp/...`; the deployed context forbids `/tmp` (same shape as BUG-003).
4. Trail open since 2026-09-04: commits 830f1fb, 67616cb, 9f457a1, 50ac771, b64ffe8 have no changes entry, no verification log, no handoff; `ROADMAP.md` still says no tuning before 2026-09-17.
5. `cmd/sh-agent` has no tests (`reorder`, the BUG-002 site, untested).

## Coverage of changed code

Nothing changed in this pass.

## Environment

Linux (Debian 13, amd64), Go 1.24.4, Claude Code 2.1.270, Antigravity CLI 1.1.25.
