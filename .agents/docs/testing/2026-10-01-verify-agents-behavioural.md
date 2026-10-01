---
owner: main
date: 2026-10-01
---
# Behavioural runs — the eight agent definitions after Phase 4, and the Antigravity `skills:` probe   (2026-10-01)

Plan `2026-09-30-audit-fix.md` Phase 4. Run by the main agent with `.agents/scratch/audit-fix/kit-agents.sh` (sources `testing/kit-preamble.sh`: fixture asserts, clone outside every checkout, real `~/.config/sh-agent` fingerprinted, `--home` wrapper on PATH) and `kit-agy-probe.sh`.

## Status: PASS (reported as n of N; two definitions were strengthened after their first three runs and re-run)

## Setup

- Fixture `~/sandbox/audit-fix/proj2`: a scaffolded Go demo (cmd/server, internal/greet with a test, `.agents/docs/` tree with a testing guide, a bug index with a Graduated table, DECISIONS.md, ROADMAP.md with one Next item), committed f48f0a1; a fresh clone per run; one comment changed in `internal/greet/greet.go` before each run so reviewers and the verifier have a diff.
- Sandbox home `sbNew`: `sh-agent install --home … --vendor claude` from the working tree, so the agents under test are the installed derived files; a copy of `~/.claude/.credentials.json` for the batch, deleted by name after. Claude Code 2.1.286 in print mode (`HOME=<sb> CLAUDE_CONFIG_DIR=<sb>/.claude claude -p … --allowedTools … Agent …`), the main agent told to dispatch one named agent and stop.
- Measurement per run from the run's transcripts under `<sb>/.claude/projects/<clone>/` (main and `subagents/`): the Agent tool's `subagent_type`; per agent, one artifact marker only the new definition writes.

## Results

| agent | runs | dispatched | marker | result |
|---|---|---|---|---|
| verifier | 5 | 5/5 | every line under `## Commands` starts with its step number; Status PASS | 5/5 |
| reviewer | 3 | 3/3 | transcript shows `skills/review-standards/SKILL.md` read by path; report has `Range reviewed:` | 3/3 |
| security | 3 | 3/3 | findings file with `Range reviewed:` (the review-standards read is the reviewer's marker, not security's) | 3/3 |
| changelog | 5 | 5/5 | first line under `## Files changed` names the `git log --stat <range>` the dispatch carried; subagent transcripts show 0 `git show` calls | range line 4/5 (c3 wrote `commit f48f0a1`); no git show 5/5 |
| documentation, first definition | 3 | 3/3 | ROADMAP.md untouched though the dispatch asked for a Later section | 1/3 — the lane rule in the body did not hold against a conflicting dispatch |
| documentation, Records block names ROADMAP/TODO/DECISIONS as never-yours | 3 | 3/3 | same | 3/3, reply names the lane 3/3 |
| bugs | 6 | 6/6 | the record file is the first Write or Edit of the run; an INDEX row written | record first 6/6 (the kit's first measure counted the Read of INDEX.md and read 1/3; re-measured by Write/Edit calls only); INDEX row 6/6 |
| explorer | 1 | 1/1 | recon brief under `plans/` | 1/1 |
| dreamer | 1 | 1/1 | artifact under `retros/` | 1/1 |

The literal prefix "from the dispatch's" in the changelog template was copied in 0/5 runs while the range line itself was written in 4/5: the template's first line, not the prose sentence, is what the agent reproduces. Left as is: the behaviour that matters (the range named, no shell command) held 5/5.

## Antigravity `skills:` probe (D5a → D5b)

Sandbox home `sbAgy` (`sh-agent install --home … --vendor agy`, token copied chmod 600 and deleted by name after, a settings file written for the run with `read_file(*)`, `command(ls)`, `command(cat)` and the clone as the only trusted workspace). Two hand-written agents: `probe` with `skills: ["skills/review-standards"]` and `probe-nokey` without it; `agy -p` (1.2.14) in a clone of `proj`, the main agent told to dispatch the agent and relay its answer.

| agent | runs | auto-denied | invalid allow entry | six criteria named verbatim | installed skill path named |
|---|---|---|---|---|---|
| probe (with the key) | 3 (s1, k1, k2) | 0 | 0 | 3/3 | 3/3 |
| probe-nokey (control) | 2 (n1, n2) | 0 | 0 | 2/2 | 2/2 |

No hang, no parse error (the only matching `cli.log` lines are one harmless `cs path` warning per start). The key is accepted but the control shows a subagent reaches a global skill by the path its body names without it, so derive keeps dropping `skills:` and its comment cites this probe (D5b). These five starts are also guide § 4's Antigravity start with the derived agents installed.

## Not examined

Antigravity runs of the eight derived agents themselves (the probe exercised the subagent mechanism and the installed derived files parsed; the agents' behaviour on agy is Phase 9's skill-matrix work). The changelog agent's reply wording. Runs on macOS.
