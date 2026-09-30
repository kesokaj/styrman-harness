---
owner: main
date: 2026-09-30
---
# Handoff — 2026-09-30 — R-parity, `sh-agent memory`, background `sh-agent dream`: built and gated; not released

Session ID: 01HVg249S3giFbkjjc1n4Wc5        (same machine: `claude --resume 01HVg249S3giFbkjjc1n4Wc5` or `agy --conversation <id>`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The operator asked for dreaming and memory to work identically on Claude Code, Antigravity and a machine with both (R10), and for a background dreamer over every project. Plan `.agents/docs/plans/2026-09-30-parity-memory-dream.md` is fully built. Then: "You dont have do cut a new version. I can do that later. I have a bigger project that will most certianly cut a version."

## Completed

- **Requirements**: R10 (R-parity), R3's dream exception, R4's dream clauses, and the Non-goals memory keys, on the operator's word (ADR-047, ADR-047a).
- **Built**:
  - `sh-agent memory add|show` on both CLIs, with Claude Code's native auto memory off per project;
  - the `project-memory` skill;
  - `sh-agent dream`: collection, redaction, a tool-free CLI call, proposals and drafts;
  - the Antigravity reader, which decodes protobuf (D30);
  - the agy child's private HOME, which grants nothing (D28);
  - flock locks;
  - the nightly timer (`--install-timer`), which `uninstall` removes.
- **Records**:
  - change entries per phase and review round;
  - bugs BUG-116..123;
  - verify logs per phase, plus the v1.17.0 gate logs `testing/2026-09-30-verify-v1.17.0-gate.md` and `-gate-r2.md` (PASS, 462 tests at 032ece1) and `testing/2026-09-30-verify-real-home-v1.17.0.md` (PASS: both CLIs from the real home, every changed skill and the dreamer on Antigravity).
- **Real home**: installed from a275b96. The binary lacks 032ece1's decoder; the next `sh-agent update` brings it.

## Open / In Progress

- **Not released.** The v1.17.0 gate stopped short on the operator's word. The last code re-review (`reviews/2026-09-30-review-a275b96-head-code.md`, Status IN PROGRESS, verdict pending) and the security re-check of the decoder commit 032ece1 were stopped unfinished, and no gate report or tag exists. The next `/preflight` over v1.16.2..HEAD owes a code and security review of a275b96..HEAD; the earlier ranges are covered, each with an Outcome.
- Requirements edits awaiting approval: `REQUIREMENTS.md:31` (R7) still calls the store "auto memory store". The proposed replacement is "memory store (written by `sh-agent memory` on both CLIs)", awaiting the operator's yes.

## Key Decisions

- Plan D1–D30. The operator's own words are on D1, D2, D3, D12, D13, D29 and D30. D29: keep both CLIs' transcripts, filtered. D30: decode protobuf properly.

## Corrections

- "The dream skills should be a features that runs autonomously in all projects registered on the machine… That was the plan." (feedback log 2026-09-30T18:36:22+02:00).

## Agent notes (unverified)

- BUG-122 (a commit ran despite scan hits) recurs a global rule, a retro graduation candidate. BUG-123's prevention could widen AGENTS.md § Testing's `--home` wrapper bullet to name skills that write `~/.config/sh-agent` with file tools.

## Environment State

- Branch: main, clean after this commit, pushed. The newest tag is still v1.16.2.
- Running processes: none. Every agent was stopped or finished.
- Scratch (deletable): `.agents/scratch/{p2-*,p3home,p4home,p5*,gate1*-agy,sh-agent-p*,agyprobe,…}`, `~/sandbox/{p2,p4,p3v,gate17,gate18,gate19,agyonly,p1fixture}/`. Every copied token is deleted.

## Artifacts to Read

- `.agents/docs/plans/2026-09-30-parity-memory-dream.md`: the design, D1–D30.
- `.agents/docs/specs/sh-agent/dream.md` and `specs/sh-agent/memory-loop.md`: how it works.
- `.agents/docs/testing/2026-09-30-verify-real-home-v1.17.0.md`: the last real-home runs.

## Next Steps

1. When the operator cuts the next version: `/preflight`. It owes the a275b96..HEAD code and security reviews named above.
2. Operator: `sh-agent dream --install-timer` on the machines that should dream nightly (the plan's `who: user` step; `systemctl --user list-timers` shows it). Run `/init-project` once in existing projects to switch Claude Code's native memory off; `sh-agent memory show` reminds you in each one.
3. Operator, one line: approve the R7 wording above (recommended: yes), or leave it.
