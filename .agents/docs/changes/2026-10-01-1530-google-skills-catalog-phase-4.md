---
owner: main
date: 2026-10-01
---
# 2026-10-01-1530 — Google skills catalog, Phase 4: both CLIs, recall, timing; the search no-index marker

Plan `plans/2026-10-01-google-skills-catalog.md` Phase 4. Log: `testing/2026-10-01-verify-google-skills-catalog-phase-4.md`.

- **Both CLIs, from each run's own record**: no catalog entry is listed as a skill (Claude Code's init event; Antigravity by two agreeing listings, self-report); six routing tasks read the expected entry, 6 of 6 on each CLI; a deploy request loads the harness's `preflight` and reads no entry (D23: the criterion named `cloud-run-deploy`, the harness's routing says `preflight`). Rows `catalog-no-leak`, `catalog-route`, `catalog-precedence` in `testing/skill-matrix.md`.
- **`sh-agent search` leaves a directory holding `.sh-agent-noindex` out of its index** (`internal/search/reconcile.go`, `TestNoIndexMarkerSkipsTheDirectory`, seen red; the marker counts committed, staged or not yet added). With the catalog indexed, recall on the fixture fell from 0.633 to 0.596 and catalog files led the hits for Google product words; the catalog now carries the marker (written by the refresh script, outside CHECKSUMS) and recall is 0.633 again over 1,421 files (D22). Docs: `specs/sh-agent/search.md`, README.
- **Timing and size**: `go test ./...` 5.7 → 12.0 s and the first install into an empty home 2.3 → 13.2 s (both over the plan's expectations; D24: accepted, the operator's one-line choice); install again 0.41 → 0.50 s; 268 → 1,838 installed files; binary 13.9 → 20.0 MB; zip 2,575 entries of 5,000.
- The interactive check (does a CLI prompt before reading under the installed skill directory) is the operator's and moved to the handoff.

Verification: `testing/2026-10-01-verify-google-skills-catalog-phase-4-go.md` PASS — 357 top-level tests with -race, three must-fail controls (the marker's effect, its place outside CHECKSUMS, the script writing it), recall 0.633, 0 catalog hits for a catalog word with a positive control. Agents dispatched: verifier ×1; Phase 3's security reviewers still in flight (three of five returned).
