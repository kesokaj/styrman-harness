---
owner: main
date: 2026-10-01
---
# Verification — Google skills catalog, Phase 4: both CLIs in the sandbox, recall, timing (2026-10-01)

Range reviewed: bfd25d3 plus the working tree (the search no-index marker), installed into the sandbox home `~/sandbox/audit-fix/sbCat` (`go run ./cmd/sh-agent install --home …`, vendor all, 35 skills)

## Status: PASS — no leak on both CLIs, 6 of 6 routed tasks on each, precedence 1 of 1 on each; two timing expectations missed and put to the operator; the interactive check is the operator's

Kit: `.agents/scratch/audit-fix/kit-catalog.sh <cli> <row>` (not committed, ADR-010): `testing/kit-preamble.sh` sourced (0 `kit: FAIL` lines), a fresh clone of the demo service per row outside every checkout, credentials copied for the batch and deleted by name after it (`ls`: no such file, both). Claude Code: `claude -p --output-format stream-json --verbose --allowedTools Read Grep Glob Skill "Bash(grep:*)" "Bash(ls:*)"`; Antigravity 1.2.14: `agy -p … --output-format stream-json` with `read_file(*)` and the read commands granted. Each row is judged from the run's own record, never from the reply: Claude Code's stream-json events (the `init` event's `skills` and `slash_commands`; `tool_use` inputs), Antigravity's conversation database for that clone (`strings <db>`).

| row | Claude Code | Antigravity |
|---|---|---|
| `catalog-no-leak` | PASS: the init event lists 53 skills and 89 slash commands; all 35 harness skills present, `google-catalog` present, no catalog entry name | PASS by self-report (its init event lists no skills): two independent runs agree — 35 names, every harness skill, `google-catalog`, no catalog entry name |
| `catalog-route` 1 BigQuery slot cost | PASS: `Skill` google-catalog; reads `bigquery-optimization`, `bigquery-slot-cost-optimizer` | PASS: the same two entries |
| 2 Spanner interleaved schema | PASS: `spanner-basics` | PASS: `spanner-basics` |
| 3 tuning a Gemini model | PASS: `agent-platform-tuning` | PASS: `agent-platform-tuning` |
| 4 agents-cli scaffold | PASS: `google-agents-cli-scaffold` (and `-workflow`) | PASS: the same |
| 5 SecOps case triage | PASS: `secops-triage` | PASS: `secops-triage` |
| 6 PAM entitlement | PASS: `iam-helper-for-privileged-access-management` | PASS: the same |
| `catalog-precedence` | PASS: `Skill` preflight (a deploy request means `/preflight`), no catalog read | PASS: the bodies of `preflight` and `cloud-run-deploy` are in the conversation, `google-catalog`'s is not, no catalog read |
| denials | 0 | 0 |

Route: 6 of 6 on each CLI. The kit's first verdicts were wrong in two places and were re-judged from the saved records, not re-run: its no-leak check compared against the sandbox's skill directory, which holds a `synced` directory Claude Code makes for itself (the repository's 35 are all listed); its precedence check wanted `cloud-run-deploy` by name where the harness's own routing sends a deploy request to `preflight` (D23).

## Search recall (`SH_AGENT_RECALL=1 go test ./internal/search -run TestRecall -v`)

| tree | all (41) | words (26) | drift (10) | indexed files |
|---|---|---|---|---|
| 754a7ce, before the catalog | 0.633 | 0.731 | 0.293 | — |
| bfd25d3, catalog indexed | 0.596 | 0.692 | 0.243 | 2,205 |
| working tree, catalog behind `.sh-agent-noindex` | 0.633 | 0.731 | 0.293 | 1,421 |

The catalog cost 0.037 of recall and put its entries first for Google product words; under the plan's 0.60 line it is excluded: `sh-agent search` now leaves a directory holding a `.sh-agent-noindex` file out of its index (`internal/search/reconcile.go`, `TestNoIndexMarkerSkipsTheDirectory`, seen red first; D22). `sh-agent search bigquery reservation slots` in this checkout: 10 catalog hits before, none after.

## Timing and size (warm build cache; second of two runs)

| measure | 754a7ce, before | with the catalog | ratio | expectation |
|---|---|---|---|---|
| `go test ./...` | 5.7 s | 12.0 s | 2.1 | under 2: missed |
| `go test -race ./...` (verifier logs) | — | about 62 s, `cmd/sh-agent` alone | — | — |
| first `sh-agent install` into an empty home | 2.3 s | 13.2 s | 5.8 | under 3: missed |
| `sh-agent install` again, nothing changed | 0.41 s | 0.50 s | 1.2 | — |
| installed files | 268 | 1,838 | | |
| each vendor's skills directory | 1.3 MB | 9.6 MB | | |
| repository zip | 4.8 MB, 1,486 entries | 7.1 MB, 2,575 entries | | under 4,000 entries: met (cap 5,000) |
| linux/amd64 binary (`go build -trimpath`) | 13.9 MB | 20.0 MB | | |

The two misses are one-time or developer-side costs: the first install after an update writes 1,570 more files, each atomically; the suite's extra time is the tests that copy the whole checkout. Put to the operator as a one-line choice (D24), recommendation: accept.

## What was not exercised

An interactive session on either CLI (whether it prompts before reading under the installed skill directory: the operator's box); the Mac; a seventh family; Antigravity's no-leak from a record rather than its own listing.
