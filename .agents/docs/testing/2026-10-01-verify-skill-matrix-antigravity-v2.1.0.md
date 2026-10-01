---
owner: main
date: 2026-10-01
---
# Verification — the harness skills changed since v2.0.0, one run each on Antigravity (2026-10-01)

Range reviewed: v2.0.0..159585a — the skill text under test is that of f62513b, unchanged for these six skills in every later commit (`git diff f62513b..159585a --stat -- skills` names only `skills/google-catalog/`). The project's release rule (`AGENTS.md` § Release) asks for one Antigravity run of each skill the range changed; the `google-catalog` router and its entries are covered by the catalog rows (`testing/2026-10-01-verify-google-skills-catalog-phase-4.md`).

## Status: PASS — 6 of 6

A first batch ran against a sandbox home that still held an older install (the install from an exported tree was refused: "checkout … is not the installed one; pass --repo"), so its six passes were discarded and the rows were run again after `go run ./cmd/sh-agent install --home <sandbox> --vendor agy` from the checkout at f62513b, with `diff -rq` of the six skill directories against the installed copies empty.

| skill changed | kit row | result | seconds |
|---|---|---|---|
| `bash` (`references/makefile-patterns.md`) | `bash` (matrix kit) | pass=1, denied 0 | 168 |
| `cloud-run-deploy` (description shortened) | `cloud-run` | pass=1, denied 0 | 43 |
| `data-gcp` | `data` | pass=1, denied 0 | 23 |
| `review` (and `scripts/restatement-grep.sh`) | `review` (extended kit) | pass=1, denied 0 | 159 |
| `feature` | `feature` | pass=1, denied 0 | 643 |
| `project-memory` | `project-memory` | pass=1, denied 0 | 76 |

Each row is the kit's own pass condition from `testing/skill-matrix.md` (an artifact on disk or a named skill in the answer), in a fresh clone of the demo service, the kit preamble sourced first; 0 "kit: FAIL" lines. The token copied into the sandbox home was deleted by name after the batch; a `find` over the sandbox finds none.

## Not checked

No agent changed in the range, so no agent row is owed. The rows prove the skill loads and its procedure runs on Antigravity; they do not exercise the specific sentences that changed (the Makefile's `--project`, the Cloud SQL role note, issue text as data, the symlinked `MEMORY.md` refusal).
