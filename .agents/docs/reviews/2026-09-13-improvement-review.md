---
owner: main
date: 2026-09-13
---
# Improvement review — the harness after the learning-loop plan   (2026-09-13)

Range reviewed: working tree at 5cdd4d6 (Phases 1–4 of `plans/2026-09-13-learning-loop-project-level.md` committed)

Method: eight independent readers, one lens each (trail, context, agents, workflow skills, knowledge skills, binary, vendors, learning loop), each proposal then handed to a skeptic told to refute it (already covered, wrong evidence, conflicts with `REQUIREMENTS.md` or an ADR, speculative), then a completeness critic. 55 agents; 46 proposals raised, 21 survived, 25 refuted, 4 added by the critic. Duplicates merged below. Every item cites what its reader verified in the files; nothing here is applied yet.

The operator's own words framing this: "the agent storing files under .agent/docs is amazing. It is helping so much when I need to switch models and agents."

## What must not be broken (the readers' unanimous strengths)

The resume order and fixed checkpoint moments in `context/AGENTS.md`; every agent writing a dated, frontmattered artifact and replying in ten lines; the feature skill's cold-reader gate; the manifest-only install contract; fail-closed derivation and frontmatter lint; the exact-diff proposal files of `/retro`.

## Proposals — numbered for a go / no-go per item

Recommendation in brackets: `[fix]` do it, `[yours]` only the operator can decide. Cost S = minutes, M = an hour.

### The trail (what makes switching models and agents cheap)

1. **One owner for `TODO.md` and ROADMAP "Now", and done items are deleted, not ticked.** `[fix]` S. Today `context/AGENTS.md` step 4 tells the main agent to clear finished work, the handoff skill forbids the main agent editing `TODO.md` inline, and the changelog agent may only append — nobody holds the delete verb. styrman-crew's `TODO.md` is 340 lines with 69 ticked items and 0 open, read second at every resume. Change: main agent owns both lines (DoD step 4); drop `TODO.md` from the changelog agent's lane and handoff Step 3 item 3; `documentation` keeps ROADMAP Next/Later only; finished lines are deleted, the `changes/` entry is the record.
2. **Handoff template: frontmatter in, `--global` out, repo-relative paths, committed.** `[fix]` S. Add `owner:`/`date:` (18 of 29 crew handoffs lack it); delete `/handoff --global` and `~/.config/sh-agent/handoffs/` (nothing reads it); paths repo-relative — 27 of 29 crew handoffs embed `/home/kesokaj/…`, which does not resolve on the Mac (`/Users/spgo/…`); keep `Session ID` (feeds `claude --resume` on the same machine); fold "Suggested Skills" into Next Steps; add a step that commits and pushes the handoff — a handoff that is not pushed does not exist on the other machine.
3. **Bug skill: the BUG record ships in the fix commit.** `[fix]` S. Step 5 currently commits, then dispatches `bugs`; move the dispatch into DoD step 4 so the record and INDEX row travel with the fix and the changelog entry can cite the id.
4. **Init-project templates match the agents they scaffold for.** `[fix]` S. The ADR skeleton in `assets/agents-tree/docs/architecture/DECISIONS.md` differs from the shape the changelog agent writes (crew keeps an agent-memory file just to remember which wins); the `AGENTS.md` template names `staticcheck` and fake `internal/config/`, `internal/api/` rows.
5. **Name the downstream mirrors as a deploy target and pin the git root.** `[fix]` S. styrman-crew `image/harness/` and argolis-devbox `deploy/harness/` vendor this repo as .git-less copies and build `sh-agent` from them; both are seven commits stale with the BUG-006-era binary, and nothing in this repo says so. Add a "Downstream mirrors" section to `deploy/remote.md` and a refresh step to the retro apply path. Code: `Inspect` must require `git rev-parse --show-toplevel` to equal the checkout — a nested plain copy today inherits the enclosing repo's commit and `update` could pull that repo (security round 3 residual, `TODO.md`).

### Agents

6. **Drop `memory: project` from bugs, changelog, documentation.** `[fix]` S. A Claude-only, per-machine cache that duplicates the trail and contradicts R7's spirit; crew and homey-troubleshooting have committed `.claude/agent-memory/` files carrying a second, Claude-only trail. Replace each memory paragraph with one sentence: the next bug number comes from `INDEX.md`, conventions are read from the newest entries. Remove the `.claude/agent-memory/` gitignore lines and the derive comment that lists `memory`.
7. **Explorer gets the write-early section; the three agents stop naming `maxTurns`.** `[fix]` S. A killed explorer leaves nothing today. The derived Antigravity file has no `maxTurns` key, so "(`maxTurns` in this file's frontmatter)" is false there — say "(set by the caller)".

### Context and workflow skills

8. **`dev-setup.md`: no unconditional Colima command.** `[fix]` S. `colima status || colima start` cannot run on the Linux box; the `docker` skill already has the Mac/Linux split. Same one-line fix in `cloud-run-deploy` step 3. Drop "(arm64)".
9. **Feature skill: e2e never deploys.** `[fix]` S. Step 5's "Deploy (only when the user asked for GCP)" is the one path that ships to Cloud Run without the security row, the tag, or a rollback; route it through `/preflight` after Step 6.
10. **Research full mode hands into the feature skill at its Step 1.** `[fix]` S. Research writes a plan in its own shape; one plan format means a resume finds unticked boxes and a readiness verdict whichever way the work started.
11. **Feedback skill: when in doubt, log; a correction of order or procedure counts even as a mild question.** `[fix]` S. Two of four retros started from an empty log and reconstructed their input from the transcript; the 2026-09-04 retro's P6 was never applied.
12. **Two duplicates created or exposed today.** `[fix]` S. The rendered context now carries "every skill loads on its own…" twice (context line 71 and the inventory sentence from Phase 1) — shorten line 71 to the list of workflow skills. The Hard-rules bullet "Plan mode or a `grilling` interrogation before large features. Trivial tasks: just act." duplicates the Feature and Trivial bullets — delete it.

### Knowledge skills

13. **golang references: 78 dangling links and sections the stack rejects.** `[fix]` M. `references/golang-ecosystem.md` carries CI, a DI framework, and `stay-updated`; `golang-testing.md` carries testify; 78 `./references/<missing>.md` links; the ecosystem file is past the Read tool's 2000-line default so its tail is never seen. Cut to what the stack uses and strip the links.
14. **docker skill: a Go multi-stage Dockerfile and the amd64 rule.** `[fix]` M. Its example builds the wrong language and its build rule contradicts the context; delete the registry section (Artifact Registry is in `cloud-run-deploy`).
15. **frontend skill: TanStack Start config guidance is stale.** `[yours]` S. The tree and the Nitro-preset line predate the `vite.config.ts` + `tanstackStart()` + `nitro()` setup; the reader could not verify the current plugin import against your projects — confirm what your `web/` uses before this is rewritten.
16. **gcloud-setup: delete the project-discovery rule, merge `gcloud-auth-verification` in.** `[fix]` S. Lines 105–111 tell the agent to discover the project id, which the context forbids (ask for the defaults); the second skill's three rules fold into one block; skill count 27.

### Binary and vendors

17. **`status` describes the installed binary, not the running process.** `[fix]` S. `go run ./cmd/sh-agent status` always warns "no build stamp" about itself; read the installed file's build info (`debug/buildinfo.ReadFile`) so the line is trustworthy from `go run`, from the installed copy, and under `go test`.
18. **Vendor facts are stale; re-run the smoke checks.** `[yours]`. Antigravity auto-updated to 1.2.2 during this session (its binary is dated 2026-09-13 20:41; the day's e2e runs after 21:26 ran on it and passed); `REQUIREMENTS.md` Known limits still says the Mac build is "confirmed at the next `sh-agent update`"; Claude Code was never verified on the Mac. Needs you on the Mac: `sh-agent update`, then the three smoke checks on each CLI.

### Decisions only you can take

19. **R7 wording.** `[yours]` Optional sentence for `REQUIREMENTS.md` R7: "The harness's own state (manifest, feedback log, the harness proposal queue, the cross-project index) lives in `~/.config/sh-agent/` and nowhere else; a project's retro and distill records are the project's own state, in its trail." Nothing built today violates R7 as read (ADR-015); the sentence removes the ambiguity three cold readers tripped on.
20. **Where the queue and index live across machines.** `[yours]` They are per machine under `~/.config/sh-agent/`; a retro on the Mac is invisible to a harness session here until the project record is pulled. The skeptic refused moving them into the harness repo's own trail for now ("only if a retro on the Mac is actually reported unreachable"). Decide after the first Mac retro.

## Refuted (the skeptics' reasons, one line each)

25 proposals were refuted; the full reasons are in the workflow record. The shapes: already covered (resume chain, docs-root rule stated 13 times but the fallback is needed per file), wrong evidence (Antigravity does preload `skills:`; the ROADMAP freeze was already dropped today), conflicts with an ADR (cold reads via explorer instead of reviewer; moving state into the repo), or speculative benefit (capping the bug index, memory-file races that never occurred). One surviving panel item (a test for the rebuild path) was already implemented today and is not listed above.

## Not examined

The Mac itself; the full skill matrix on either CLI; the vendored `anti-ui-slop`, `terraform`, `gke-basics` content beyond their descriptions.

## Outcome (2026-09-13, later the same session)

The operator approved every item ("go all"). Applied in the commit that follows a4146dc: items 1–17 and 19 (the R7 sentence, with ADR-016). Item 18: both CLIs re-checked from the sandbox home on Antigravity 1.2.2 (context, `golang` auto-load, `explorer` dispatch — all pass, log clean); the real home on this box and the Mac stay the operator's (`TODO.md`); the Known-limits sentence in `REQUIREMENTS.md` about the Mac was not edited — proposed wording is in the handoff. Item 20: deferred by decision until the first retro on the Mac. Item 15 was verified against styrman-crew's `web/` before the rewrite: SPA mode, client embedded in the Go binary, no Nitro — the panel's "Nitro plugin" detail was wrong and is not in the skill.
