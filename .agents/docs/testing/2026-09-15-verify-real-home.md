---
owner: main
date: 2026-09-15
---
# Verification — the harness installed into this box's real home, smoke checks on both CLIs   (2026-09-15)

Range reviewed: the harness at a33978d, installed into `/home/kesokaj`; the smoke checks ran from that install. It was then reinstalled at 8202d6a, which changed only a test and trail files.

## Status: PASS for both CLIs

Scope: task 7.5 of `plans/2026-09-15-harness-practices-g01-g26.md` (D48):
- `sh-agent install` and `sh-agent status` into the operator's real home;
- then the three smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI, with the real `HOME`, each from a fresh clone of a sandbox demo.

This is a smoke set, not the matrix; the matrix is `testing/2026-09-15-verify-harness-practices-e2e.md`. The runner was `.agents/scratch/realhome-smoke.sh`, with its outputs in `.agents/scratch/realhome-smoke/` and `.agents/scratch/realhome/`.

## Install and status

| step | command | result |
|---|---|---|
| before | `~/.local/bin/sh-agent status` | exit 1: installed on 2026-09-14 from 8046e4d, `binary 8046e4d5f24b — older than the repo (a33978d)`, harness files differing from the repo |
| install at a33978d, clean checkout | `go run ./cmd/sh-agent install` | exit 0: `building sh-agent from /home/kesokaj/workspace/styrman-harness`, `45 written, 215 unchanged, 0 removed, 0 kept — 260 files managed` |
| status | `~/.local/bin/sh-agent status` | exit 0: installed from a33978d, `binary a33978d7ab9e` with no warning, `drift none — 259 files match the repo` |
| reinstall at 8202d6a, after the smoke checks | `go run ./cmd/sh-agent install` | exit 0: `1 written, 259 unchanged` (the binary) |
| status | `~/.local/bin/sh-agent status` | exit 0: installed from `8202d6a+dirty`, `binary 8202d6a5105b+dirty` with no warning, `drift none — 259 files match the repo` |

The `+dirty` comes from the uncommitted Done verifier log, a trail file that is not installed.

`~/.claude/CLAUDE.md` is still the generated import: the generated header and the line `@AGENTS.md`.

## Smoke checks

Each check ran with the real `HOME`, in a fresh clone of a sandbox demo.

For the checks, Antigravity's `~/.gemini/antigravity-cli/settings.json` was backed up and widened with `command(*)`, `read_file(*)`, `write_file(*)`, `edit_file(*)` and `list_dir(*)`, then restored on exit. `cmp` against the backup reports the file identical, and its allow list is back to its six narrow entries.

**Check 1 — context, no tools:** the first heading, the Installed counts, and Definition of Done step 2's reviewer sentence.
- **Claude Code: PASS (4 s).** `# AGENTS.md — global agent context`; 27 skills, 7 subagents; step 2 quoted with "or any fix under the `bug` skill", matching `~/.claude/AGENTS.md`.
- **Antigravity: PASS (7 s).** The same heading, 27 and 7, and step 2 quoted verbatim.

**Check 2 — skill auto-load:** review `internal/greet/greet.go` against the Go standards, change nothing, and name the rules and their skills.
- **Claude Code: PASS (237 s).** Applied `golang` and `review-standards` rule by rule, with two low findings. No code changed; it wrote and committed only its own review report (the trail rule, BUG-012).
- **Antigravity: PASS (35 s).** Named its rules from `golang` and `review-standards`, including `references/go-review.md`. No findings, a clean clone, no commit.

**Check 3 — agent dispatch:** the `verifier` on an uncommitted comment edit; Antigravity dispatches `explorer` first.
- **Claude Code: PASS (59 s).** The `verifier` ran in the clone: `Status: PASS`, `Scope: complete`, pass 13, log `testing/2026-09-15-verify-greet-comment.md`. The tree holds only its log and the planted edit.
- **Antigravity: PASS (135 s).** `explorer` mapped the three packages into `plans/2026-09-15-recon-packages.md`. The `verifier` reported `Status: PASS`, log `testing/2026-09-15-verify-uncommitted.md`. The CLI waited for its background agent ("root agent idle; waiting for 1 background task(s)") and returned.

`grep -c "Failed to parse"` prints `0` on the newest Antigravity log, `cli-20260915_233425.log`, which this run created.

## Observations

- **Claude Code check 2** could not read the skills' `references/` files from the clone: this headless run's permissions do not cover `~/.claude/skills/*/references/`. It named the rules it therefore left unchecked. An interactive session prompts instead.
- **Claude Code checks 2 and 3** printed that the claude.ai Homey MCP connector needs authorizing. That is unrelated to the harness.
- **Antigravity check 3** waited for its background agent, bounded by `--print-timeout`, unlike the print-mode turn in BUG-031.

## Not verified here

- The Mac.
- The full matrix from the real home; the matrix ran in the sandbox home.
- Interactive prompts on either CLI; every check was headless.
- A session started before this install; each check started a new CLI process.

## Environment

- Linux 6.12.107+deb13-amd64 x86_64, Debian GNU/Linux 13.
- Claude Code 2.1.272 and Antigravity CLI 1.2.3 from `~/.local/bin`; go1.24.4 linux/amd64.
- Harness a33978d installed for the checks, then `8202d6a+dirty` after the reinstall.
