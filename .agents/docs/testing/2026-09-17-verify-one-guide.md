---
owner: verifier
date: 2026-09-17
---
# Verification — one testing guide (guide.md rename + agent/context updates)

Range reviewed: working tree (uncommitted) — staged rename `.agents/docs/testing/sh-agent/guide.md` -> `.agents/docs/testing/guide.md`; modified `agents/documentation.md`, `agents/verifier.md`, `context/AGENTS.md`, `skills/handoff/SKILL.md`, `skills/init-project/SKILL.md`, `skills/init-project/assets/AGENTS.md`, `skills/review/SKILL.md`; untracked `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md` (cited per dispatch, not re-run: see that file for repro-5..11 evidence on BUG-063). Also present in the tree but not part of this dispatch's named scope: modified `.agents/docs/bugs/BUG-062-...md` and `BUG-063-...md` (bug records updated to `status: fixed`, presumably by the main agent in this same session; not touched by this verifier run).

Guide used: `.agents/docs/testing/guide.md` (confirmed present, dated 2026-09-17, per ADR-032/BUG-062).

## Status: PASS

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-oneguide/go-test-1.txt`; counts: pass 124, fail 0, skip 1 at e9fa15f+dirty for `go test -race -v -count=1 ./...` (skip is `TestFindRepo`, working directory is inside a harness repo — same pre-existing/expected skip noted in prior logs). vs previous: `.agents/docs/testing/2026-09-17-verify-bug-062.md`, pass 124→124, skip 1→1 (unchanged).
- `go run ./cmd/sh-agent install --home .agents/scratch/oneguide-home` → exit 0; "260 written, 0 unchanged, 0 removed, 0 kept".
- `go run ./cmd/sh-agent status --home .agents/scratch/oneguide-home` → exit 0; `binary   e9fa15fd91a3+dirty` (no warning); `drift      none — 259 files match the repo`.
- Confirmed installed `.claude/agents/verifier.md` and `.gemini/config/agents/verifier/agent.md` reference `$docs_root/testing/guide.md` throughout (steps 2, 4, Scope line, Records line) and contain no `testing/<app>` residue.
- Confirmed installed `.claude/agents/documentation.md` and `.gemini/config/agents/documentation/agent.md` reference `testing/guide.md` (layout diagram and "Shared seam" line), no `testing/<app>`.
- Confirmed installed context copies `.claude/AGENTS.md` and `.gemini/config/rules/AGENTS.md` say `docs/testing/guide.md` (one per project, written by the verifier whenever it is absent, maintained by both), no `testing/<app>`.
- `go run ./cmd/sh-agent list` → exit 0; skills and agents sections print counts (27 skills, 7 subagents) with no parse errors.

## Root cause
N/A — Status is PASS.

## Coverage of changed code
This is a docs/prompt-text change set (agent markdown, context markdown, skill markdown, a file rename) with no `.go` files touched — `git diff --stat` against HEAD shows no `.go` paths. Go tests are unaffected by this change and pass as a baseline/regression check, not as targeted coverage. No suppressions (`t.Skip`, `//nolint`, `.only(`, etc.) were added by this change set — the only skip present (`TestFindRepo`) is pre-existing and unchanged in count from the previous log.

## Not examined
- Sandbox start of each CLI (guide's tier 4: one Claude Code start, one Antigravity start, asked a context question quoting new wording) — not run in this dispatch; budget-scoped to the items the caller listed.
- Real-home smoke checks (guide's tier 5) — not run; the caller's scope was the scratch-home install/status/list checks plus the installed-file greps, not a real-home install.
- `/review --docs` over this range — not run by this verifier dispatch (review is a separate agent/step per AGENTS.md Definition of Done step 2).
- Full contents of `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md` beyond what was cited by the dispatch instruction — read the file's diff excerpt via BUG-063's own record but did not re-verify its repro-5..11 claims independently.
- The concurrent edits to `BUG-062`/`BUG-063` status fields — out of this dispatch's named scope, not evaluated for correctness.

## Environment
Linux 6.12.107+deb13-amd64, go version go1.27.1 linux/amd64 (user-level SDK ahead of system go, per AGENTS.md). No Docker/Colima involved (no Dockerfile in this change set).
