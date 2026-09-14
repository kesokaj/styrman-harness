---
owner: verifier
date: 2026-09-14
---
# Verification guide — sh-agent

The standing procedure for this repo. Written under the rule in `agents/verifier.md`
("What to run"): the first run under a new project writes this guide; later runs cite
it and record outcomes only. Update this file in the same run that changes the steps.

## 1. Go checks (every change; skip only when no `.go` file is in the diff — say so)

```bash
gofmt -l .                # must print nothing
go vet ./...
go test -race ./...
```

## 2. Harness self-check (every change to `context/`, `agents/`, or `skills/`)

```bash
go run ./cmd/sh-agent list      # sanity: 27 skills, 7 agents parse; counts change only when a skill/agent is added or removed
go run ./cmd/sh-agent render | wc -l   # rendered context/AGENTS.md must stay under 100 lines
```

## 3. Scratch-home install and drift check

Never install into the operator's real `$HOME` from a dev session.

```bash
rm -rf .agents/scratch/<slug>-home
mkdir -p .agents/scratch/<slug>-home
go run ./cmd/sh-agent install --home .agents/scratch/<slug>-home
go run ./cmd/sh-agent status --home .agents/scratch/<slug>-home   # must report "drift none"
```

Confirm every `agents/*.md` and its derived Antigravity `agent.md` under
`.gemini/config/agents/<name>/agent.md` in the scratch home carry the same
`## Records` section (exactly one each) when an agent definition changed.

## 4. Sandbox start of each CLI (change to `context/`, `agents/`, or `skills/`)

One start of Claude Code and one of Antigravity in the sandbox (paths and flags in
`.agents/docs/testing/2026-09-03-verify-sandbox.md`), each asked a context question
that quotes the new wording, log grepped for parse errors. The full skill-matrix
scenarios are `.agents/docs/testing/skill-matrix.md` — run the matrix in full only
before saying "complete" or "all skills"; a targeted change runs only the matrix row(s)
it affects.

## 5. Real-home smoke checks (after any install into a real home)

The three checks in `.agents/docs/testing/2026-09-03-verify-real-home.md` (context
questions, one skill load, one agent dispatch) from that home, on each CLI, before the
install is called done. The Mac procedure is
`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`.

## Full matrix

`.agents/docs/testing/skill-matrix.md` is the reproducible spec of every skill/agent
sandbox scenario. Cite it rather than duplicating rows here; add a row there when a
skill or agent gains new behavior worth an end-to-end check.
