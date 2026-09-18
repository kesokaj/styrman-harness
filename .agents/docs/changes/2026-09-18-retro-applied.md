---
owner: main
date: 2026-09-18
---
# Retro 2026-09-18 applied: seven harness-scope proposals     (2026-09-18)

## Summary

The seven proposals of `retros/2026-09-18-retro.md` (queued in `~/.config/sh-agent/proposals/2026-09-18-retro.md`) are applied on the operator's "apply all": four sentences in `context/AGENTS.md`, two triggers in `feedback-log`, the verifier's early-log rule, and one clause in the feature skill's cold-reader brief. The three bugs whose rules graduated — BUG-078, BUG-079, BUG-080 — move to the Graduated table, and the two events the retro found without numbers get theirs (BUG-081, BUG-082), graduated on arrival because their rules land in this same commit. Every other place picks the change up with `sh-agent update`.

## Type

context, skills, agents (harness rules)

## Files

- `context/AGENTS.md` — § Definition of Done step 2: a fix's test builds the ARRANGEMENT the finding measured (BUG-081). § Hard rules: an installed COMMAND is verified by the shell that will run it (BUG-078, BUG-079); the read-only mode resolves paths as the real run does, a child-cwd path is made absolute first (BUG-080); a must-fail mutation counts only when the mutated code compiled and a named test failed (BUG-082). § Work: the dispatch sentence a `verifier` gets is spelled out ("after the FIRST command, Status IN PROGRESS, update after every one") — the one edit beyond the seven proposals, from the behavioural runs below (BUG-083).
- `skills/feedback-log/SKILL.md` — triggers 6 Redirection (signal `redirection`, new table row) and 7 A missed filing (signal `restatement`).
- `agents/verifier.md` — the log is written after the FIRST command with `IN PROGRESS`, updated after EVERY command; a log still in progress at cutoff is a failed dispatch.
- `skills/feature/SKILL.md` — the cold-reader brief asks whether each test fixture in a verify clause produces the condition it claims.
- `.agents/docs/bugs/INDEX.md`, `BUG-078…080` (Related lines), `BUG-081`, `BUG-082` — graduation; `BUG-083` — recorded, open.

## Verification

- `go run ./cmd/sh-agent install --home .agents/scratch/retro-home` → 260 written; `status` → drift none, 259 files match; `list` ok; `gofmt -l cmd internal` empty; `go test ./...` ok (the Go code is untouched; `TestHarnessMap` and the inventory tests read the changed files).
- Context change, behavioural: `claude -p --append-system-prompt "<the new Hard-rules bullet>"` from a scratch directory quoted the new sentence verbatim and named BUG-078 and BUG-079 — the text given, not the installed file, so it proves the wording reads as intended, not the install.
- Verifier change, behavioural: `claude -p --agents <json built from the working-tree agents/verifier.md>` in a scratch clone (`.agents/scratch/verifier-run`). First run: every subagent Bash and Write call hit the permission gate, no agent ran, no log — not evidence either way. Three runs with `--allowedTools` for the verifier's tools and a dispatch that did NOT repeat the early-write sentence: 1 of 3 wrote the log with `IN PROGRESS` after the first command and updated it after the second; 2 of 3 ran every command and wrote once at the end (tool-call order read from each subagent transcript). The definition alone does not bind — BUG-083, the class of BUG-062/063/064; the dispatch text remains the enforcement and now says "after the FIRST command, IN PROGRESS, then after every command". The proposal stays applied: the wording is what the dispatch repeats.
- Not exercised: Antigravity start of the changed files (the operator's Mac; the previous two context changes started clean there and nothing here changes frontmatter); the feature and feedback-log skills in a live run.

## Definition of Done

- Review: `reviewer` with the docs brief → `reviews/2026-09-18-review-worktree-retro-docs.md`, REQUEST CHANGES, three findings; § Outcome: finding 2 (signal vocabulary) fixed; finding 3 (two sentences citing no BUG id) fixed by recording BUG-081 and BUG-082 and citing them; finding 1 (reviewer, security and explorer still say "one finding worth recording") left — the proposal was the verifier's, whose passing commands read as "no finding"; for a reviewer a finding is the unit; whether those three want the per-step rule is the next retro's question, in the handoff.
- Agents dispatched: `reviewer` (docs lens), `bugs` ×3 (graduation; BUG-081/082; BUG-083). No `verifier` dispatch: no Go source changed; the checks above ran bare in the main thread and their status was read.
- Status set to `applied 2026-09-18 (<commit>)` in the queue file and the record after the commit.
