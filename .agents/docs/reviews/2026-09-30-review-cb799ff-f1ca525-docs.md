---
owner: reviewer
date: 2026-09-30
---
# Review — cb799ff-f1ca525-docs   (2026-09-30)

Range reviewed: cb799ff..f1ca525

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Docs lens, read-only. Nothing in the living docs still says Claude Code writes or loads memory natively (every hit of "auto memory", `autoMemory*` and "loads it itself" at f1ca525 is either the R10 "switched off" statement or a description of the old default path). The worst gap: `deploy/local.md:59` (the operator's pre-report read list) still says report sends `proposals/*.md` only, while `internal/harness/report.go` now also sends `proposals/dream/harness/*.md`. Every new command and flag is in the source usage text, the README and the specs; the gaps are a missing README table row for the timer flags, stale restatements of `uninstall`, and skill-matrix rows that describe the removed `/dream` settings mechanism. No finding is blocking.

## Findings
- **[medium / confidence: high]** stale: `.agents/docs/deploy/local.md:59` ("What leaves: … `proposals/*.md`") and `:66` ("Text files (`proposals/*.md`)"), `.agents/docs/deploy/remote.md:66` ("retro proposals and index"), `.agents/docs/specs/sh-agent/delivery.md:54` ("the retro proposals") — `internal/harness/report.go` (diff hunk "harnessDir") now also copies `proposals/dream/harness/<name>.md`. The operator reads this list before every report, so a file class that leaves the machine is missing from it. Only `local.md:94` and README line 191/69 say `report` sends `harness/`. Fix: add `proposals/dream/harness/*.md` to the four sentences (and say drafts, project proposals, state, log and failures never leave).
- **[low / confidence: high]** stale: `.agents/docs/deploy/local.md:130` ("`sh-agent uninstall` removes exactly the files the manifest lists") and `cmd/sh-agent/main.go:10` and `:52` (usage: "remove exactly what install wrote / the files install wrote, nothing else") vs `README.md:184` and `main.go:63`/`:172-181`, which remove the dream timer too, with no manifest left. Fix: add the timer clause to those three.
- **[low / confidence: high]** `README.md:180` and `README.md:181-194` — the command table has no row for `dream --install-timer | --remove-timer` (only the prose bullet at README.md:69 and `deploy/local.md:81-95` name them), and the `status` row does not mention the new `dream  last run …` line printed by `main.go:614`. Fix: add both to the table.
- **[low / confidence: medium]** stale test row: `.agents/docs/testing/skill-matrix.md:31` (`dream` row: "`/dream` accepts only that exact path or the default one"; `--settings` holding `autoMemoryDirectory`) and `:37` (`dream-agy-only`: "Step 1 item 0: no settings file on that machine") — `skills/dream/SKILL.md` no longer reads settings and has no Step 1 item 0 (store is `<root>/.agents/memory` on every machine). The rows were run against the old skill; the plan's own box for the Antigravity run of `dream`, `retro`, `distill` and `init-project` is still unticked (`plans/2026-09-30-parity-memory-dream.md:274`). Fix: rewrite both rows' setup and comments to the new Step 1 and re-run them.
- **[low / confidence: medium]** narrowed: `.agents/docs/specs/sh-agent/delivery.md:77` ("Built and shipped: … version and, since v1.13.0, search and index") lists the verbs and omits `memory` and `dream`. Fix: name them with their specs.
- **[low / confidence: medium]** `skills/init-project/SKILL.md:116` — "since Claude Code would not load the new store" (item 3, the migration guard) describes the pre-R10 mechanism; under R10 the store is read by `sh-agent memory show`, and the guard's real reason is that the setting points leftovers at the store. Fix: reword.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/dream.md` (What it does) and `README.md:191` — a project's first dream reads back 14 days (`internal/dream/dream.go:53`, `firstLook`, pinned by `collect_test.go:352`); both say only "newer than the watermark". Fix: state the 14-day first look.

## Checked and consistent
- `memory add|show`: source usage (`main.go:62,65`), README rows (192-193), `specs/sh-agent/memory-loop.md`, `skills/project-memory`, `feedback-log` rule 7, `context/AGENTS.md:9` (Resume step now `sh-agent memory show`; "Claude Code loads it itself" removed), sh-agent-matrix row.
- `dream --all|--project|--dry-run`: usage `:64`, README 191, `specs/sh-agent/dream.md`, `deploy/local.md`, matrix rows; exit codes (1 failed project, 2 run failed before the first, lock exits 0), 256 KB cap, 30-day `dream-failed/` prune, `--mode plan` for agy and `--tools ""` for Claude all match `internal/dream`.
- `--install-timer|--remove-timer`: usage `:63`, README 69 and 219-220, `deploy/local.md` (03:00, PATH of the installing shell, linger, `Persistent=true`), spec; `uninstall` timer removal in README 184, `deploy/local.md:87`, matrix row; `status` dream line in spec, `deploy/local.md:85`, `main.go:614`.
- AGENTS.md rule 1 and Layout, README map and Dream bullet, REQUIREMENTS R3/R4 text agree. `retro` renumbering (source 9 dream proposals, 10 cross-project) matches `specs/sh-agent/dream.md:39`; no other reference to the old number 9 found.
- Dropped text with a new home: `feedback-log` rule 7's "skip when the directory does not exist / symlink / tracked" moved into the verb's refusals (README 192); `/dream`'s split-store and untrusted-workspace rows moved to `memory-loop.md:16,33`.

## Unverified
- The `sh-agent` on this box's PATH prints a `--help` with no `memory` or `dream` lines: it is an older installed build, not the tree at f1ca525. The usage constant in the source (`main.go:62-65`) does list them; I did not build the tree (read-only).
- Whether `README.md` renders the new mermaid `linkStyle` indices (0-22) on the right edges: not rendered.
- The real user timer on this box (plan E2E 5, `plans/…parity-memory-dream.md:269`) is the operator's step and is unticked; the docs describe it as working.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Code correctness (docs lens only). The full text of `internal/search/dream.go` and the dated `testing/2026-09-30-verify-*` logs (not living docs). `skills/project-memory/SKILL.md` and `.agents/docs/specs/sh-agent/memory-loop.md` were read by grep hits and diff previews, not line by line. `*_test.go` files beside changed code were grepped for old wording only (no `stale test:` hit found).

## Outcome
1. medium, report contents stale in `deploy/local.md`, `remote.md` and `delivery.md`: fixed by `documentation`. Each names `proposals/dream/harness/*.md` and its origin, and the operator's read-before-report list includes it.
2. low, `uninstall` descriptions and the README command table: fixed. `main.go` doc comment and usage, `deploy/local.md` Roll back, a README row for `--install-timer|--remove-timer`, and the `status` row's dream line.
3. low, stale skill-matrix rows (`dream`, `dream-agy-only`): fixed; the store is `<root>/.agents/memory` on every machine.
4. low, `delivery.md:77`, `init-project` Step 5b item 3's rationale and the undocumented 14-day first look: fixed. The built list gains `memory` and `dream`; the migration rationale no longer assumes native memory loads the store; the spec states the first look.
