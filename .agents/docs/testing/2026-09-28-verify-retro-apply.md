---
owner: main
date: 2026-09-28
---
# Verification — retro 2026-09-28 proposals applied   (2026-09-28)

Range reviewed: 236285e..worktree. The eight retro proposals edit four files: `context/AGENTS.md` (1, 6, 8), `skills/research/SKILL.md` (2), `skills/feature/SKILL.md` (4, 5) and `agents/verifier.md` (3, 7).

## Status: PASS (one control did not fail; see below)

## Commands

- `gofmt -l cmd internal`, `go vet ./...` and `go test -race ./...` are clean.
- `TestRenderedContextFitsAntigravityBlock` passes. `sh-agent render` measures **21,443** of 22,000 bytes.
- Scratch install into `.agents/scratch/retro-0928`: drift none (257 files). Both rendered contexts carry the new DoD clause (grep 1 and 1) and end with GANNET.

## Behavioural runs

**Setup.** A fresh demo Go repository per run, under `.agents/scratch/retro-probe/runs/`, outside the harness. `check_demo` asserted every fact before each CLI start: its own repo root, `/healthz`, the Setup commands, `testing/guide.md`, and passing tests. `claude -p` got the prompt on stdin and the tools named in `--allowedTools`. The script is `.agents/scratch/retro-probe/run.sh`.

### `agents/verifier.md` (proposals 3 and 7), as `verifier-probe`

The definitions were built from the working tree ("new") and from HEAD ("old"). The dispatch did not mention combining commands.

| Definition | Log written, Status PASS | Subagent Bash calls (refused) | Shared calls with `; echo …$?`: attempted / executed |
|---|---|---|---|
| new | 2 of 2 | 13 (3), 11 (5) | 5 / 3, 3 / 0 |
| old (control) | 2 of 2 | 4 (1), 9 (4) | 3 / 2, 3 / 0 |

*Corrected after the gate review (`reviews/2026-09-28-review-v1.10.0-d3cc2b7-code.md`).* The first version of this table counted attempted shared calls as if they ran, and counted `verifier-old-1` as 0. The permission gate of `--allowedTools` refused 13 subagent Bash calls in total ("Contains simple_expansion", "Output redirection … needs approval", "Contains command_substitution"; the count includes any tool error). Those refusals include every attempt at the "redirected to its own file" form, and all shared calls in `verifier-new-2` and `verifier-old-2`. Under the project's rule, a refused run proves nothing for the refused calls. The verifier text was therefore narrowed: shared calls are only for output read inline, and a long suite's redirect is its own call.

**The control did not fail.** The old definition's verifier already combined commands with per-command status lines on this small repo. The change makes the definition agree with that behaviour, and with the dispatches that asked for it all session. It does not change measured behaviour here. The turn-limit failures it answers happened on this repository's full gate (three runs cut off at 30 turns), which the probe does not reproduce. Proposal 7 (foreground long suites) is not exercised, because the demo suite runs in seconds.

### `context/AGENTS.md` sentence 6 ("is it live?")

`--append-system-prompt` got the rendered context, with the sentence (new) and without it (control). Prompt: "Is the /healthz change live on Cloud Run yet? … do not run any deploy."

| Context | States it is not live | Offers `/preflight` itself | Puts the step on the operator |
|---|---|---|---|
| new | 2 of 2 | **2 of 2** ("Should I run `/preflight` to ship it?") | 0 of 2 |
| control | 2 of 2 | 0 of 2 | 2 of 2 ("To get it live you would ask to 'deploy'"; run 2 reads "say 'ship' or 'deploy' and I'll run `/preflight`", which is close to an offer, so the contrast is weaker than 2 vs 0 suggests) |

## Not exercised

- The shared-call and redirect forms under an `--allowedTools` gate that permits them: the probe's gate refused them, so their turn savings are not measured.

- Context sentences 1 (low findings that add friction) and 8 (stop a finished agent with background work): both need a multi-agent session to reproduce.
- The two `feature` skill edits (4 and 5), which are prose in a skill that loads only for features. The install and parse check cover them.
- The `research` skill bullet (2), for the same reason.
- Antigravity for every probe.
