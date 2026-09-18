---
owner: main
date: 2026-09-18
---
# Handoff — 2026-09-18 — retro applied; sh-agent done; nothing queued

Session ID: 28ca24dd-a2d8-4278-85aa-70e24bbed445
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 1358bcc (this handoff's commit follows it), clean and pushed. Tags `v1.0.0` … `v1.2.3`, all published.

## Context

The session's Definition of Done — "work until the delivery mechanism sh-agent is flawless" — was met and handed off in `2026-09-18-sh-agent-works-everywhere-done.md`. After it, `/retro` produced seven harness-scope proposals and the operator said "apply all"; this handoff closes that application. No plan is active; `ROADMAP.md` "Now" and `TODO.md` are empty.

## Completed

- **Retro applied** (1358bcc; `changes/2026-09-18-retro-applied.md`): four sentences in `context/AGENTS.md` (installed COMMAND verified by the shell; dry-run resolves paths as the real run, child-cwd paths absolute; a fix's test builds the ARRANGEMENT the finding measured; a mutation counts only when compiled) plus one § Work sentence spelling out the verifier's dispatch phrasing; `feedback-log` triggers 6 Redirection (`redirection`) and 7 A missed filing (`restatement`); `agents/verifier.md` log after the FIRST command; `skills/feature` cold-read fixture clause.
- **Bugs:** BUG-078/079/080 graduated; BUG-081 (fix passed a fresh-fixture test, failed the next round) and BUG-082 (mutation verdict on a build error) recorded and graduated on arrival; BUG-083 recorded, open (below).
- **Verified:** `install --home .agents/scratch/retro-home` 260 managed, `status` drift none, `list` ok, `go test ./...` ok; `claude -p --append-system-prompt` quoted the new Hard-rules sentence with BUG-078/079; three `claude -p --agents` runs of the working-tree verifier in `.agents/scratch/verifier-run`.
- **Docs review:** `reviews/2026-09-18-review-worktree-retro-docs.md`, three findings; two fixed, one left (§ Outcome).
- Statuses: the seven proposals are `applied 2026-09-18 (1358bcc)` in `~/.config/sh-agent/proposals/2026-09-18-retro.md` and `retros/2026-09-18-retro.md`.

## Open / In Progress

- **BUG-083, open:** the verifier's rewritten early-log rule held in 1 run of 3 when the dispatch did not repeat it. The definition stays as written; the enforcement is the dispatch sentence (`context/AGENTS.md` § Work, now spelled out). Nothing further is owed unless a dispatch that repeats it is seen failing.
- **Docs-review finding left:** `agents/reviewer.md:20`, `agents/security.md:18`, `agents/explorer.md:17` still say "as soon as you have one finding worth recording" — the next `/retro` decides whether they get the per-step form the verifier got.
- **Other places pick this up with `sh-agent update`.** This box's real home was last installed from `main` at 20ca4b4 (before 1358bcc); the operator runs `sh-agent update` here and on the Mac and devbox when convenient.
- Carried from the previous handoff, unchanged: the README browser chain on the Mac with `V=v1.2.3`; `linux/arm64` never executed; the darwin swapped-signature edge; the kept feedback clone on this box (`~/.config/sh-agent/feedback`) to delete by hand; the plan Follow-ups for next session (writability probe before any write first; `yamlQuote` line breaks; marker regex pinning; swallowed `filepath.Abs` error; settings editor move).
- Requirements edits awaiting approval: none. `close-map-gaps`: ADR-035, still waiting on another place's proposal.

## Key Decisions

- Finding 1 of the docs review left, not fixed: the proposal was verifier-scoped; a reviewer's unit is the finding (review § Outcome).
- BUG-081/082 recorded and graduated in one commit because their rules landed in the same commit (the "every failure gets a number" rule, docs-review finding 3).
- One edit beyond the seven proposals: the § Work dispatch sentence for the verifier, from the 1-of-3 result; the changes entry names it.

## Corrections

- none logged: `~/.config/sh-agent/feedback.jsonl` has no entry for this session. The operator's "Are we done with the sh-agent now? it seems like you are fixing skills" was a status question, answered; not logged as a correction.

## Agent notes (unverified)

- The behavioural-run recipe that worked: build the agent JSON from the working-tree `agents/<name>.md` (frontmatter → description/tools/model, body → prompt), `git clone` the repo into `.agents/scratch/<name>-run`, `claude -p --allowedTools "<the agent's tools as patterns>" --agents "$(cat …json)" < /dev/null`, then read the tool-call order from `~/.claude/projects/…<scratch dir>…/<session>/subagents/*.jsonl`. Without `--allowedTools` every subagent call hits the permission gate and the run proves nothing — worth a line in the project `AGENTS.md` § Testing if it recurs.
- `.agents/scratch/` is large (releases, downloads, matrix fixtures, mutation harnesses, verify homes, three verifier-run logs) — gitignored, deletable.

## Environment State

- `main` at 1358bcc + this handoff, clean, pushed. Six tags, all published; repository `PRIVATE`.
- Suite: `gofmt`, `go vet`, `go test ./...` clean on go1.27.1 at 1358bcc (`-race` and `govulncheck` last run at the v1.2.3 gate; no Go source changed since).
- `ROADMAP.md` "Now" empty; `TODO.md` empty.

## Artifacts to Read

- `.agents/docs/changes/2026-09-18-retro-applied.md` — what changed, how it was verified, the 1-of-3 result.
- `.agents/docs/retros/2026-09-18-retro.md` — the seven proposals and their statuses.
- `.agents/docs/bugs/BUG-083-verifier-early-write-rule-not-followed-every-run.md` — the open one.
- `.agents/docs/handoffs/2026-09-18-sh-agent-works-everywhere-done.md` — the sh-agent closeout this handoff follows.

## Next Steps

1. **Operator:** `sh-agent update` on this box, the Mac and devbox to pick up 1358bcc (the context and skills changed; nothing in the binary).
2. **Next session:** the writability probe before any write (plan Follow-ups, `bug` skill), then the other Follow-ups as Tasks.
3. `/retro` when it next runs: the three "one finding worth recording" agents (docs-review finding 1); BUG-083's dispatch-vs-definition evidence.
4. `close-map-gaps` when the other place's proposal arrives (ADR-035).
