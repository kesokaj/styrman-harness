---
owner: changelog
date: 2026-09-03
scope: changeset:b18aab80d649
---

# Harness Evaluation Suite & Setup Tooling

Added an evaluation engine, CLI, tasks, and setup tooling to evaluate `styrman-harness` across tokens, tools, turns, time, and quality ratings, including A/B baseline comparison (harnessed vs vanilla).

## Added
- `setup-eval.sh`: Bootstrap and setup script for installing the eval CLI, enabling required GCP APIs (Vertex AI, Cloud Storage, Cloud Run), provisioning GCS bucket for run bundles, and deploying the Web UI.
- `cmd/eval/`: Evaluation CLI entrypoint (`eval-harness`) supporting `list`, `run`, `report`, `compare`, and `doctor`.
- `internal/eval/`:
  - `types.go`: Evaluation models (`Task`, `EvalRunResult`, `ComparisonSummary`).
  - `config.go`: Configuration manager supporting GCP parameters and local overrides.
  - `tasks.go`: Task registry with Greenfield Go CLI, Greenfield Go REST API, and Brownfield Trail Resume.
  - `runner.go`: Isolated sandbox workspace orchestrator supporting Claude Code, Antigravity (`agy`), and mock runners.
  - `telemetry.go`: Extraction parser for tokens, tools, turns, and duration from JSON envelopes and transcripts.
  - `grading.go`: Deterministic build/test/git/trail verification and Vertex AI LLM-as-a-judge scorer.
  - `storage.go`: Local JSONL ledger appending and optional GCS bundle sync.
  - `eval_test.go`: Unit and integration test coverage.
- `web/Dockerfile`: Container packaging for deploying the visualizer to Cloud Run.
- `.gitignore`: Added ignore rules for `bin/`, local eval ledger, and eval config.
