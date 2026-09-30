---
owner: main
date: 2026-09-30
---
# Phases 4–5 of the parity plan: `sh-agent dream` runs, and its nightly timer     (2026-09-30)

## Summary
The background dreamer is complete (R3's exception, ADR-047a). Under D22, the map, `report`'s change and the timer land in this one commit (R9, AGENTS.md rule 3).
- **The run**: `sh-agent dream --all | --project DIR` sends each changed project's collected, redacted prompt to one CLI with no tools, from an empty `~/.config/sh-agent/dream-work/`. It parses the first decodable `{"proposals","skills"}` and writes queued proposals: project targets to `proposals/dream/<date>-<project>.md`, harness targets to `proposals/dream/harness/`, draft skills to `proposals/dream/drafts/`. The watermark advances only after a successful write. A PID lock, `dream.log`, `dream-state.json`, `dream-failed/` (scrubbed, 30-day prune), exit 1 or 2, and the `status` line complete it.
- **The agy invocation** was wrong as first probed (D7a). With a real 14 KB prompt, agy called a tool, was auto-denied, and ended with no answer, and `--mode plan` did nothing beside `--disable-slash-commands`. The prompt now opens, for both CLIs, with "every tool call in this session is refused and ends it with no answer", and agy runs `--mode plan --output-format text`: 3 of 3 agy runs and 1 of 1 Claude Code run ok (D7b).
- **`report`** sends `proposals/dream/harness/*.md` (R8, D19) and nothing else of the dream output. It refuses a symlinked `proposals/dream` or `harness` directory; the guard's mutation compiled and failed its test.
- **`/retro`** source 9 reads the project's dream proposals and sets their Status; `/retro --harness` reads `proposals/dream/harness/`.
- **The timer**: `--install-timer` writes a systemd user service and timer (Linux) or a launchd agent (macOS), with ExecStart at `~/.local/bin/sh-agent` and the installing shell's PATH, then enables it. A unit on disk with other text is kept and named. `--remove-timer` removes it, and so does `uninstall`, even with no manifest; `uninstall --dry-run` only says so. The timer is not in the manifest.
- **Map**: a new node for the dream job, with edges from sessions and into `~/.config/sh-agent/`; `TestHarnessMap` requires `dream`. AGENTS.md rule 1 says the timer is neither a delivery nor a report path. README and `deploy/local.md` describe the install and its limits (linger, PATH).

## Type
feature (phases 4 and 5 of 6)

## Files changed
- `internal/dream/{run,timer}.go` and their tests
- `internal/dream/{dream,prompt}.go`
- `internal/harness/report.go` and `report_test.go`
- `cmd/sh-agent/main.go` and `memory_test.go`
- `cmd/sh-agent/main_test.go` (`TestHarnessMap`)
- `README.md` (map, bullets, rows)
- `AGENTS.md` (rule 1)
- `.agents/docs/deploy/local.md`
- `.agents/docs/specs/sh-agent/dream.md` (new)
- `skills/retro/SKILL.md`
- the plan (D7b, ticks)
- `testing/…-phase-4.md` and `testing/…-phase-4-5.md`

## Agents dispatched
verifier (`testing/2026-09-30-verify-parity-memory-dream-phase-4-5.md`, PASS: 443 tests; § 3b report dry run lists only the harness file; its mutation of the watermark rule failed `TestRunInvalidOutputKeepsRawAndWatermark`; the main agent's e2e log audited, no claim contradicted).

## Verification
- `go test -race -count=1 ./...` exit 0.
- E2E (`testing/2026-09-30-verify-parity-memory-dream-phase-4.md`):
  - Claude Code over two seeded projects: 2+2 proposals, then `unchanged` twice;
  - a planted secret absent from the prompt;
  - a deleted root `skip gone`;
  - agy-only 3 of 3 after D7b.
- `systemd-analyze verify` exit 0 on units rendered into a scratch home.

Not exercised: the real timer on this box (the plan's `who: user` box), and a real `report` push carrying the harness file.

## Requirements
R3's exception, R4's dream clauses, R8, R10 (ADR-047, ADR-047a). The map change was approved in ADR-047a.
