---
owner: main
date: 2026-10-08
---
# Verification — metrics store, sandbox runs on both CLIs   (2026-10-08)

Range reviewed: 2c0f215..c407eb8 plus the uncommitted Phase 7 tree (dream metrics block, `/retro` source 11, the no-store refusal, `sh-agent` second-level verbs). That tree is committed with this log.

## Status: FAIL (one matrix-row check: Claude Code `retro`, 0 of 2; every metrics check passes on both CLIs)

## Method

The scratch kit `.agents/scratch/kits/metrics-sbx.sh <claude|agy>` is gitignored. It is modelled on `testing/agy-rel/run.sh` and `testing/ab-kit/run.sh`. Each run:
- sources `kit-preamble.sh` first, for the `--home`-pinned wrapper built fresh from this tree and the exit checks;
- uses a fresh home under `~/sandbox/metrics-sbx/homes/<cli>`, with `sh-agent install --vendor <cli>`;
- clones the agy-rel demo fixture;
- seeds three identical `-race` corrections in the home's `feedback.jsonl`;
- links the CLI's credential (never copies it) and removes it by name before the exit check;
- starts the CLI from `env -i` with the credential broker off.

The sequence: a first session ("list the Go files, run go vet once"), then `sh-agent metrics ingest`, then `/retro`, then a second `ingest`, then `show --json`, `show` and `report`. The check script `metrics-sbx-check.sh` re-reads the saved outputs and the kit's own store; RUN and CHECK are separate (BUG-140).
- Claude Code: `claude -p` with `--model claude-sonnet-5-5`, `--permission-mode acceptEdits`, the matrix's `Bash(<cmd>:*)` allowlist plus `python3`, and `--add-dir <home>/.config/sh-agent`.
- Antigravity: `agy -p --mode accept-edits`, with settings granting `command(*)` (the retro skill runs compound shell, as agy-rel's `retro` row does).

Preamble exit lines on every run: `no processes left by the run`, `real config unchanged`, `no credentials left in KIT_HOME`.

## Results

| Check (per run) | Claude Code | Antigravity |
|---|---|---|
| ingest after a real first session adds a session | 2 of 2 (`1 sessions, 4 steps`) | 1 of 1 (`1 sessions, 36 steps`) |
| `show --json`: this CLI ≥ 2 sessions and > 0 tool calls, the other CLI 0 | 2 of 2 (12 and 7 tool calls) | 1 of 1 (64 tool calls; notes name agy's missing cost, denials, interrupts, compactions) |
| the retro session called `sh-agent metrics` (source 11), read from the kit's own store | 2 of 2 | 1 of 1 |
| the retro record mentions the metrics | 2 of 2 ("Metrics show 1 session and 0 corrections, so no regression to cite") | 1 of 1 ("tool error rate 0.091 (1/11 calls) … no metric regression against the prior 30-day window") |
| `report` rendered, no `<script>` | 2 of 2 | 1 of 1 |
| no CLI or ingest call failed; agy: no auto-denied call | 2 of 2 | 1 of 1 |
| **skill-matrix `retro` row check** (`proposals/*retro*` in the home) | **0 of 2** | 1 of 1 |

## The failing row

**Observed.** Both Claude Code retros proposed the `-race` rule as **project** scope, under Step 2's `target-in-project` rule. Their reasoning: "the correction concerns this project's documented command, with no hit in another project", because the fixture's `AGENTS.md` documents `go test ./...`. They wrote the record and a `retros.jsonl` line, and queued no harness file.

**What this does not establish.**
- The change under test added source 11 to Step 1 and nothing to Step 2's scope rules, but no control run at the pre-feature commit was made. That the change is uninvolved is **not proven**.
- The earlier passing runs of this row were on Antigravity (2026-10-07) and on Opus. This kit ran Claude with Sonnet, which is also untested as a cause.

**Open:** a control run of the Claude row at 2c0f215, with the same model.

## Also exercised

- **Scratch-home install** (`.agents/scratch/metrics-install-home`): `status` reports `drift none — 1789 files match the repo`, and the `binary` line has no warning. The installed binary's `metrics show` with no store printed `no metrics store at …`, exit 1.
- **Normaliser.** Run 1 stored `sh-agent metrics --json --project --since`, losing `show`. The fix (keep `sh-agent`'s second-level verbs) is under `TestVerb`, and run 2 stored `sh-agent metrics show --json --project --since`.

## Not examined

- The Claude row on Opus, and the control at the pre-feature commit.
- A plain-copy place and a release-binary place (sh-agent-matrix rows: not run).
- The real home: the operator's first real ingest and the § 5 real-home smoke.
