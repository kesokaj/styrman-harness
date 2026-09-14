---
owner: main
date: 2026-09-14
---
# Handoff — 2026-09-14 — BUG-015 fixed; the mirror-refresh prompt

Session ID: session_01CtMYM9AGp33edxBQGURNaB        (same machine: `claude --resume`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean at 2c8990e plus this handoff's commit; pushed

## Context

The operator asked for one prompt to paste into fresh sessions in styrman-crew and argolis-devbox so each refreshes its plain copy of this repo. Reading both consumers and the harness diff first surfaced BUG-015: since 82fd4a6, `install` rebuilt from any tree with `go` on PATH, which would have failed every boot in styrman-crew's read-only mirror. The fix landed first; the prompt pins it as the floor. Newest wins: this file supersedes `2026-09-14-session-close.md`.

## Completed

- BUG-015 fixed in 2c8990e: `shouldRebuild` is gated on `harness.IsCheckout` (a `.git` entry in the repo); a plain copy installs the running binary and prints `plain copy, not its own git checkout: the running binary is installed as is`. Failing test first (`TestPlainCopyNeverRebuilds`), describing text moved in `AGENTS.md`, `README.md`, `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`; ADR-020; `.agents/docs/changes/2026-09-14-bug-015-plain-copy-never-rebuilds.md`.
- Verified: `.agents/docs/testing/2026-09-14-verify-bug-015.md` (gofmt/vet/race clean; a real build against a read-only git-less mirror; the checkout path still rebuilds). Reviewed: `.agents/docs/reviews/2026-09-14-review-bug-015.md` (APPROVE, one cosmetic fixed) and `2026-09-14-review-bug-015-docs.md` (one medium fixed). Testing guide gained § 3a.
- The prompt: `.agents/docs/plans/2026-09-14-mirror-refresh.md` — one text for both projects, with per-project sections, checked adversarially against each consumer repo before it was handed over.
- The empty gitignored `.claude/agent-memory/` directories were removed from this checkout so the mirrors' `diff -rq` stays at its documented three lines.

## Open / In Progress

Nothing queued in this repo; `TODO.md` and `ROADMAP.md` empty.
- Requirements edits awaiting approval: none.
Operator-only, unchanged from the session-close handoff: the Mac (update, status, three smoke checks per CLI; Claude Code there unverified), Go ≥ 1.24.6 on both machines, improvement-review item 20 after the first Mac retro, the full skill matrix.

## Key Decisions

- ADR-020: the rebuild belongs to a git checkout; a plain copy never rebuilds. Rejected: consumers setting the internal `SH_AGENT_CHILD` marker or stripping `go` from PATH (pushes a harness defect into every consumer).
- The mirror refresh stays a task in each consumer, run by that project's own `VENDORED.md` recipe; the harness only supplies the prompt and the floor commit.

## Corrections

- none this task (no feedback-log entry).

## Agent notes (unverified)

- The crew boot failure is predicted from code and from a local read-only-mirror run, not observed in a styrman-crew container; the crew session confirms it in the boot log after the refresh.
- The devbox reader saw argolis-devbox HEAD move during its read (another session pulled at 11:19); the prompt tells the session to pull and re-read 032 first.

## Environment State

- Branch: main, clean after this commit; origin at HEAD after the push. Real home updated to this commit (`sh-agent update`, no drift).
- Running processes: none. Sandbox unchanged.

## Artifacts to Read

- `.agents/docs/plans/2026-09-14-mirror-refresh.md` — the prompt and why each line is there
- `.agents/docs/bugs/BUG-015-install-rebuild-plain-copy.md` — the consumer call sites and the prevention rule
- `.agents/docs/handoffs/2026-09-14-session-close.md` — the map of the two-day session this task followed

## Next Steps

1. Operator: paste the prompt from `.agents/docs/plans/2026-09-14-mirror-refresh.md` into a fresh session in styrman-crew and one in argolis-devbox, with all three repos side by side.
2. After both refreshes: the two consumer changes entries name the new COMMIT (2c8990e or later); nothing else is owed here.
