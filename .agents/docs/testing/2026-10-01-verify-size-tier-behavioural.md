---
owner: main
date: 2026-10-01
---
# Behavioural run — the size-tier paragraph in the global context   (2026-10-01)

Plan `2026-09-30-audit-fix.md` Phase 3 (G11, D3): a project without a `docs/` root runs Trivial and Task only — no verifier log, no changes entry, no handoff — and the session offers `/init-project` once. Run by the main agent with the kit `.agents/scratch/audit-fix/kit-size-tier.sh` (sources `testing/kit-preamble.sh` first: fixture asserts, clone outside every checkout, real `~/.config/sh-agent` fingerprinted, `--home` wrapper on PATH).

## Status: PASS

## Setup

- Fixture: `~/sandbox/audit-fix/proj` (a three-file Go program, no `.agents/`, committed b542088); one fresh clone per run under `~/sandbox/audit-fix/runs/<id>`; the clone's parent is no git checkout (the preamble asserts it).
- Two sandbox homes, installed with `sh-agent install --home <sb> --vendor claude`: `sbNew` from the working tree (the new context, 14,8xx bytes, carrying the paragraph), `sbOld` from a git worktree at 329a99b (the previous context, 21,906 bytes, without it). `grep -c 'Unscaffolded project' <sb>/.claude/AGENTS.md` printed 1 and 0.
- Claude Code 2.1.286 in print mode from each home: `HOME=<sb> CLAUDE_CONFIG_DIR=<sb>/.claude claude -p --permission-mode acceptEdits --output-format text --allowedTools …`, the prompt on stdin: add a `-name` flag to main.go, finish the task the way the rules say, do not ask, commit, no push. A copy of `~/.claude/.credentials.json` sat in each sandbox home for the batch and was deleted by name after.
- Measurement per run, from the run's own transcript under `<sb>/.claude/projects/<clone path>/`: an `Agent` tool call with `"subagent_type": "verifier"`; `.agents/docs` created in the clone; the reply mentions `/init-project`; commits in the clone; `go run . -name kit` prints `hello, kit`.

## Results

| context | runs | verifier dispatched | .agents/docs created | offered /init-project | flag works |
|---|---|---|---|---|---|
| new (sbNew: smoke, n1–n5) | 6 | 0 of 6 | 0 of 6 | 6 of 6 | 6 of 6 |
| previous (sbOld: o1–o5) | 5 | 5 of 5 | 5 of 5 | 0 of 5 | 5 of 5 |

A sample new-context reply: "this repo has no `docs/` folder, so the rules only call for a commit here: no test log, change entry or handoff. If you want those, you can set the project up with `/init-project`." Commits: 5 of 5 batch runs committed (the smoke run could not: a fresh clone has no git identity; the kit sets one since).

## Correction during the run

The kit's first "verifier dispatched" measure grepped the transcript for the word `"verifier"`, which the loaded context itself contains, so the smoke and the first two new-context runs read as 1. The measure was changed to the Agent tool's `"subagent_type": "verifier"` field before the remaining runs and every run was re-measured from its transcript with that pattern; the table above is the re-measurement.

## Not examined

Antigravity (no `agy -p` run: the paragraph is vendor-neutral text and the installed context is the same file on both CLIs; Phase 9's real-home smoke and skill matrix cover agy). The Trivial path (a one-line change) was not run separately.
