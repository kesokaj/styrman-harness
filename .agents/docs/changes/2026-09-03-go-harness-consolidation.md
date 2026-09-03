---
owner: changelog
date: 2026-09-03
---

# Consolidate harness tooling into one Go binary; proportional DoD, informing gates

## Summary
Replaced twelve bash/Python scripts (setup.sh x3, merge-settings.py,
derive-agy-agent.py, trail-check.sh, subagent-check.sh, resume-context.sh,
destructive-guard.py, format-hook.sh, feedback-capture.sh, statusline.sh) with a
single Go binary (`cmd/harness` + `internal/{trail,gate,hookio,install,dispatch}`).
`setup.sh` now only bootstraps `go run ./cmd/harness install`.

Alongside the consolidation: Definition-of-Done rows became proportional
(resumability rows always fire; accountability rows fire on thresholds in
`.agents/harness.json`), verification became a script that only escalates to a
model to explain a failure (measured 2s / $0.0000 on this repo), gates now inform
rather than block except at one commit boundary (an always-blocking SubagentStop
hook previously cost 8 forced continuations, $0.236, 37.8s for a one-word
subagent output), artifacts are matched by YAML frontmatter (`owner`/`date`)
instead of filename glob, and `destructive-guard` was deleted in favor of a
credential-staging gate (operator runs in an isolated microVM). Also added:
`harness doctor` (drift detection), `harness check` (trail freshness), a
per-dispatch cost ledger, CI, and the root `AGENTS.md`/`REQUIREMENTS.md`/
`ROADMAP.md`/`TODO.md` scaffold. Agent frontmatter gained `maxTurns` and
`effort`; the verifier agent gained `isolation: worktree`.

The architectural reasoning for these changes is recorded in ADR-001 through
ADR-004 (`.agents/docs/architecture/DECISIONS.md`) — this entry records that the
implementation landed, not the rationale.

Deployment is explicitly out of scope: the operator deploys `~/.agents/bin/harness`,
`~/.claude`, and `~/.gemini` by hand, later.

## Type
task

## Files changed
- `cmd/harness/main.go` — new binary entrypoint
- `internal/trail/*.go` — trail freshness, owed-rows computation, git/source helpers
- `internal/gate/*.go` — trail gate, subagent-stop gate, secret/credential gate, resume-context injection
- `internal/hookio/hookio.go` — hook I/O plumbing shared across gates
- `internal/install/{install,derive,doctor}.go` — installer, Claude→Antigravity agent derivation, drift doctor
- `internal/dispatch/{cli,dispatch,verify}.go` — dispatch cost ledger and script-first verification
- `setup.sh` — reduced to a 3-line bootstrap
- `AGENTS.md`, `REQUIREMENTS.md` (R1–R8, non-goals), `ROADMAP.md`, `TODO.md` — root scaffold
- `.agents/docs/architecture/DECISIONS.md` — ADR-001..004 (already recorded, not duplicated here)
- deleted: `merge-settings.py`, `derive-agy-agent.py`, `trail-check.sh`, `subagent-check.sh`, `resume-context.sh`, `destructive-guard.py`, `format-hook.sh`, `feedback-capture.sh`, `statusline.sh`, two prior `setup.sh` variants

## Verification
`go test -race ./...` and `go vet ./...` per `AGENTS.md`'s working rules; hook
startup measured at ~2.6ms (Go) vs ~22.8ms (Python); verification measured at
2s / $0.0000 on this repo. Per repo convention, a `verifier` subagent owns the
full verification log under `.agents/docs/testing/` — not duplicated here.

## Requirements
R1 (one contract, two vendors), R2 (one deployable artifact), R3 (deploy
overwrites, drift is loud), R4 (gates inform, block only at a boundary), R5
(proportional DoD rows), R6 (spend a model only where judgement is needed), R7
(artifacts identified by frontmatter), R8 (isolation covers host, gate covers
remote).
