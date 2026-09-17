---
owner: main
date: 2026-09-17
---
# Verification — `sh-agent report` end to end, from dave to GitHub   (2026-09-17)

Range: the binary deployed into the real home by the operator's hourly `sh-agent-update.timer` at 09:03:48 (`sh-agent version` → `4a610e93cf0b`, "3 written, 257 unchanged" in `journalctl --user -u sh-agent-update`), against `git@github.com:kesokaj/styrman-harness-feedback.git`, created private and EMPTY minutes before.

## Status: PASS — all three criteria of the plan's Verification section

| # | Criterion | Command | Observed |
|---|---|---|---|
| 0 | What would leave (read before sending) | `go run ./cmd/sh-agent report --dry-run` | only `state.json` present under `~/.config/sh-agent/` (feedback.jsonl, retros.jsonl, proposals/ absent and named); no redactions; `would clone … into ~/.config/sh-agent/feedback`; exit 0 |
| 1 | Golden path | `sh-agent report` from `~` | cloned the empty repository, `reports/dave/2026-09-17-090741/{report.json,state.json}`, commit `fce67aa` pushed; `git rev-parse --abbrev-ref @{upstream}` → `origin/main`; exit 0 |
| 2 | Second run the same day | `sh-agent report` | a second directory `2026-09-17-090754`, commit `71f749e`, plain push; the first directory untouched; exit 0 |
| 3 | `--dry-run` against the real checkout | `sh-agent report --dry-run` | `git status --porcelain --untracked-files=all` plus `rev-parse HEAD` identical before and after (`71f749e…`, no status lines); exit 0 |
| 4 | A fresh clone shows them | `git clone … .agents/scratch/feedback-fresh` | both run directories, both commits; `gh repo view --json` → `PRIVATE`, default branch `main` |

The shipped `report.json` of run 1: `host dave`, `linux/amd64`, `vendor all`, `files 260`, `harness_commit 4a610e9`, `installed_at 2026-09-17T09:03:48+02:00`, `drift 0`, `skills 27`, `agents 7`, `go go1.27.1`, `redactions []`.

## Not exercised

- **Redaction on real data.** This home had no `feedback.jsonl`, `retros.jsonl` or proposals, so no redaction fired end to end; redaction is covered by the unit tests only (`internal/harness/report_test.go`, with the must-fail mutant recorded in the plan).
- A rejected push, a failed first clone, and a same-second collision against the real remote — covered against local bare repositories in `TestSendReport` and `TestResolveFeedbackFailures`, not here.
- Any machine but this one; macOS; a place with no harness checkout.
- The run went through the harness checkout found from `state.json`, so the no-checkout `null` path was not exercised end to end (unit: `TestReportWithoutCheckout`, `TestReportRecord`).
