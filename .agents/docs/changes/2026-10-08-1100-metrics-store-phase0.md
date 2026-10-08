---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 0: requirements, ADR, spec

- Plan `plans/2026-10-08-metrics-store.md`: the trajectory and agent/sandbox metrics store, on both CLIs, with visual HTML reports and feedback through dream and `/retro` proposals. D1–D29 were settled from the operator's answers and the cold read.
- Cold read: `reviews/2026-10-08-review-metrics-store-plan-cold-read.md` (`reviewer`). It found 15 gaps and 7 blocking questions. All are answered as D7–D24, and D28 supersedes D18 (Antigravity tool names come from `steps.metadata`; `step_payload` stays unread).
- `REQUIREMENTS.md`: R3, R4 and R7 edited as approved by the operator ("go all approved", 2026-10-08). The driver now serves the metrics store too. `metrics.db` is a store of record: never a cache, never sent, migrated forward, kept by `uninstall`. Dream ingests first.
- `architecture/DECISIONS.md`: ADR-055.
- `specs/sh-agent/metrics.md` (new, status draft): verbs, flags, the transcript-root rule, `Verb`, the tool vocabulary, `show --json`, the report, metric definitions, the schema, invariants, failure modes.
- Verification: `go test ./internal/harness` exit 0 after the REQUIREMENTS edit. No code changed.
- Agents dispatched: `Explore` (recon), `reviewer` (cold read).
