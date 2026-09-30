---
owner: main
date: 2026-09-30
---
# Behavioural runs — the verifier and another agent's artifact under the docs root   (2026-09-30)

Retro `retros/2026-09-30-retro.md`, Proposal 6: `agents/verifier.md` gains "The tree under test is source, tests and config: … treat a new or changed path under `<docs-root>` or `.agents/scratch/` — another agent's artifact, written while you ran — as seen, listed in the log by name, never as the tree moving under you."

## Status: PASS on Claude Code — new definition 9 of 9 classify the foreign file as seen; control 1 of 8 calls the tree "not frozen", 0 of 9 new

## Kit

`.agents/scratch/vkit/` (gitignored): `fixture/` is `greetbox`, a two-commit Go module with `AGENTS.md` (Setup commands: gofmt, vet, `go test -race`), `.agents/docs/testing/guide.md` and an empty `reviews/`. `run.sh <new|control> <n>` clones it to `~/sandbox/vkit/runs/<arm>-<n>` (outside every checkout; `git -C ~/sandbox rev-parse --git-dir` fails), proves the fixture facts (two commits, the guide, the command line in `AGENTS.md`, `go test` green, the pinned `sh-agent` wrapper first on `PATH`, the foreign file absent), then runs `claude -p --agents <json from the definition>` with the prompt on stdin and `--allowedTools` naming the verifier's five tools plus Agent. The `new` arm builds its JSON from the working-tree `agents/verifier.md`; `control` from `git show HEAD:agents/verifier.md`. A background writer waits for the verifier's log to appear (its early write) and 2–4 s later writes `.agents/docs/reviews/<date>-security-probe.md`, another agent's artifact by shape. The dispatch text never states the changed rule; it asks for the project's commands, the early log, and the git status at start and end.

Three rounds, six or five runs each, differing in the dispatch's last sentence and in what context the subagent could read:

| Round | Git-status instruction | `--setting-sources` | Why |
|---|---|---|---|
| 1 | "record git status at the start and again at the end, and say what differs" | `project,local` | first shape |
| 2 (strict) | "confirm at the end that the tree is as you found it plus your log, and record what it shows" | `project,local` | the wording under which today's verifiers flagged a violation |
| 3 (prod) | as round 2 | default (the real home's context, with the frozen-tree rule, loads) | in rounds 1–2 neither arm could read the rule the control was expected to cite |

Attribution: every run's subagent transcript under `~/.claude/projects/-home-kesokaj-sandbox-vkit-runs-<arm>-<n>/*/subagents/` was grepped for the new sentence ("source, tests and config"): present in 9 of 9 `new` transcripts, absent in 8 of 8 `control` transcripts — the `--agents` definition ran in every case, over the real home's installed `verifier` in round 3. Tool order was the same in every run: Bash, Read, Read, then the commands. No `.err` file holds an auto-denied line. `~/.config/sh-agent/` lists the same five entries before and after (`feedback.jsonl proposals retros.jsonl search.db state.json`); the wrapper pinned every `sh-agent` call to `.agents/scratch/vkit/home`, which stayed empty (BUG-107).

## Results

| Round | Arm | Foreign file listed by name | Classified as another agent's artifact, seen | Called the tree "not frozen" / a violation | Status of the log |
|---|---|---|---|---|---|
| 1 | new (3) | 3 of 3 | 3 of 3 ("another agent's artifact, seen and not touched. No source, test or config path changed") | 0 of 3 | FAIL on coverage (Shout untested), as the fixture intends |
| 1 | control (2) | 1 of 2 (control-2's end status ran before the file landed) | 1 of 2 ("appeared during my run … presumably another agent") | 0 of 2 | FAIL on coverage |
| 2 | new (3) | 3 of 3 | 3 of 3 | 0 of 3 | FAIL on coverage |
| 2 | control (3) | 3 of 3 | 3 of 3 ("not mine … another agent likely wrote it") | 0 of 3 | FAIL on coverage |
| 3 | new (3) | 3 of 3 | 3 of 3 ("seen and not touched"; "no source, test or config path changed, so the tree was not moving under the commands") | 0 of 3 | FAIL on coverage |
| 3 | control (3) | 3 of 3 | 2 of 3 neutral | **1 of 3** (control-2: "another agent wrote it during this run, so the tree was not frozen") | FAIL on coverage |

Summed: the new definition names the file and classes it as seen in 9 of 9 runs and never calls the tree not frozen (0 of 9); the old definition reaches the wrong classification 1 of 8 times, and only when it can read the global frozen-tree rule (round 3). The two production logs of today (`testing/2026-09-30-verify-redaction-rules.md`, `-verify-v1.13.1-gate.md`) were the same shape as that one control run.

## What was not exercised

- A source, test or config file changed during the run (the other half of the new sentence: "IS the tree moving"): not staged; the fixture writer only produces the docs-root case.
- Antigravity: no run; the next release owes one run of the changed agent there (`AGENTS.md` § Release).
- The subagent's log Status when the foreign file is the only change and the code is fully covered (every run FAILed on the fixture's untested `Shout`, by design, so PASS-vs-FAIL was not the measured variable).
- Runs with the log written late (after the writer's 150 s fallback): none occurred; every log appeared within the first minute.
