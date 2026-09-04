# Handoff — 2026-09-04 — argolis-devbox migration instructions written from this checkout; the harness itself unchanged

Session ID: 09809bdb-a200-4633-a977-fde22b3ae51f
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean, HEAD `6893e22` (unchanged this session; only this handoff is added)

## Context

The operator used this checkout as the working seat to prepare the migration of `argolis-devbox` (their single-user GCP devbox VM) from the old `devshell_prompt` harness to styrman-harness, the same way styrman-crew was migrated on 2026-09-03/04. Their instruction, verbatim: "DO NOT DO ANY WORK IN THE CURRENT PROJECT styrman-harness, only write instructions how to replace the current harness in argolis-devbox to this harness in styrman-harness. Put the instructions in argolis-devbox and I will start them later." Nothing in this repo was changed.

## Completed

- Recon of argolis-devbox and four instruction files written into `/home/kesokaj/workspace/argolis-devbox/.agents/docs/`: `plans/2026-09-04-recon-vm-harness.md`, `plans/2026-09-04-harness-migration.md` (two-scope checklist), `plans/2026-09-04-vm-harness-replacement.md` (scope B plan draft: six operator questions, decisions D1–D9, four phases, e2e criteria), `handoffs/2026-09-04-harness-migration.md` (with start prompts A–E). Cold-read by `explorer` (one wrong line citation fixed). Commits `7e9d071` and `7c73752` on argolis-devbox `main`, pushed at the operator's word ("push it"); the devbox VM auto-updates and restarts its two services within five minutes of a push.
- Memory note for this checkout: `~/.claude/projects/-home-kesokaj-workspace-styrman-harness/memory/project-harness-migrations.md`.

## Open / In Progress

- **This repo, unchanged from `2026-09-03-harness-verified.md`:** the Mac still needs its first install (`git clone … && go run ./cmd/sh-agent install`), `sh-agent status`, one interactive start of each CLI, and `sh-agent update --check` after a push from the other machine (plan section 12, `TODO.md` high priority). No tuning before 2026-09-17 (ADR-009), then `/retro`.
- **argolis-devbox:** scope A and scope B not started. They run from `/home/kesokaj/workspace/argolis-devbox` with start prompt A from its handoff, never from here. Phase 4 verification needs the operator in the devbox terminal (no `.env` on this box).
- **Downstream mirrors now pin this repo's commit:** styrman-crew `image/harness/COMMIT` = `6893e22` (v1.20.0); argolis-devbox will pin the HEAD at its Phase 1. A change to `context/`, `skills/`, or `agents/` here reaches neither VM until its mirror is refreshed (crew: `image/harness/VENDORED.md`; devbox: `deploy/harness/VENDORED.md` once it exists). Not on `ROADMAP.md` yet; noted here, not edited inline.

## Key Decisions

- The devbox gets the crew model: a mirror of this repo at a recorded commit inside the devbox repo, `sh-agent` built on the VM, installed keyed on the mirror commit, old artifacts moved never deleted, no settings files. Rationale and the rejected clone-on-VM alternative (private repo; the VM's token is scoped to argolis-devbox) are in the devbox plan.
- Pushed to argolis-devbox only after the operator said so, because every push restarts the VM's services.

## Corrections

- None in this repo this session (no feedback-log entry written).

## Agent notes (unverified)

- The devbox handoff's Environment State was corrected after the push in `7c73752`; its plan still says "re-read `git rev-parse HEAD` at the time" for the pin, which is right — do not assume `6893e22` if this repo moves first.
- Two identical migrations in two days suggest a `/distill` candidate: "migrate a project to styrman-harness" (recon → two-scope checklist → plan with D1–D9 → handoff with prompts; mirror + `sh-agent install` + move-never-delete). Third repeat is the trigger per the global context.

## Environment State

- Branch: main, clean apart from this handoff; nothing built or installed; `sh-agent status` in the real home unchanged (no drift at `6893e22`).
- Running processes: none. The devbox VM is applying `7e9d071`/`7c73752` on its own timer.

## Artifacts to Read

- `/home/kesokaj/workspace/argolis-devbox/.agents/docs/handoffs/2026-09-04-harness-migration.md` — the devbox handoff with the start prompts (start there for devbox work)
- `/home/kesokaj/workspace/argolis-devbox/.agents/docs/plans/2026-09-04-vm-harness-replacement.md` — the devbox scope B plan
- `/home/kesokaj/workspace/styrman-harness/.agents/docs/handoffs/2026-09-03-harness-verified.md` — the previous handoff here; its Mac item is still the only open item for this repo
- `/home/kesokaj/workspace/styrman-crew/image/harness/VENDORED.md` — how a downstream mirror is refreshed

## Suggested Skills

- None for this repo. For the devbox session: `feature` (scope B), `bash` (the boot scripts).

## Next Steps

1. This repo: the Mac, as in the previous handoff; nothing else until 2026-09-17.
2. argolis-devbox: open a session in its checkout and paste start prompt A.
3. After the next change to `context/`, `skills/`, or `agents/` here: refresh the crew mirror (and the devbox mirror once Phase 1 has landed) and push each project.
