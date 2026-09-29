---
owner: main
date: 2026-09-29
---
# Search index, Phase 0: R4, baselines   (2026-09-29-1747)

Plan: `plans/2026-09-29-search-index.md`.

- **R4 changed** on the operator's approval of the exact text: one dependency, `modernc.org/sqlite`, for the search index, not vendored. ADR-045 supersedes ADR-002 ("stdlib only") and ADR-029 (offline first build) in part. Project `AGENTS.md`, `README.md` and `deploy/local.md` say the first build of a version fetches one module. Committed separately, before this entry.
- **Recall question set** `internal/search/testdata/recall.tsv`: 41 questions from this repository's history, 10 of them drift questions with hand-checked answer sets.
- **Recall test** `internal/search/recall_test.go`: runs only with `SH_AGENT_RECALL=1`. It has a grep arm now; Phase 4 adds the index arm. `TestRecallQuestionsAreWellFormed` always runs and checks the counts and that every answer path exists.
- **Baselines:** `testing/2026-09-29-token-baseline.md` (subagents spend 59% of tool output on file reads) and `testing/2026-09-29-recall-baseline.md` (one grep pass: mean recall 0.379, drift 0.353).
- **Plan:** three cold-read rounds (`reviews/2026-09-29-coldread-search-index-r1.md` to `-r3.md`); the cap residue resolved by the operator's choice, D19.

Agents dispatched: `explorer` twice (recon, consumer audit), `reviewer` three times (cold read). Diff read by the main agent: one test file and its data, under the reviewer threshold for Phase 0 alone; the Phase 1 review covers it with the rest of the source.

Verification: `gofmt -l cmd internal` empty; `go vet ./...` exit 0; `go test -race ./...` exit 0.
