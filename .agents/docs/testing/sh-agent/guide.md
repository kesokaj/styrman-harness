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
gofmt -l cmd internal     # must print nothing; NOT `gofmt -l .` — it walks .agents/scratch,
                          # where a stray Go file fails a gate that has nothing to do with it (BUG-048)
go vet ./...
go test -race ./...
```

## 2. Harness self-check (every change to `context/`, `agents/`, or `skills/`)

```bash
go run ./cmd/sh-agent list      # sanity: every skill and agent parses; its section headers print the counts
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

## 3a. Plain-copy install (change touching checkout-vs-mirror rebuild logic)

When a change affects `IsCheckout` / rebuild-decision logic in `cmd/sh-agent` or
`internal/harness/repo.go`, the scratch-home tier above only exercises the checkout
path. Also build the real binary and run it against a read-only, `.git`-less copy:

```bash
go build -o .agents/scratch/<slug>/sh-agent ./cmd/sh-agent
mkdir -p .agents/scratch/<slug>/mirror
for d in context skills agents commands rules; do
  mkdir -p ".agents/scratch/<slug>/mirror/$d"
  rsync -a --exclude .git "$d/" ".agents/scratch/<slug>/mirror/$d/"
done
chmod -R a-w .agents/scratch/<slug>/mirror
.agents/scratch/<slug>/sh-agent install --repo .agents/scratch/<slug>/mirror \
  --home .agents/scratch/<slug>/home2   # must say "plain copy" and not rebuild
ls -a .agents/scratch/<slug>/mirror     # must show no .agents/ created
.agents/scratch/<slug>/sh-agent status --repo .agents/scratch/<slug>/mirror \
  --home .agents/scratch/<slug>/home2   # first line "plain copy"; drift none
chmod -R u+w .agents/scratch/<slug>/mirror   # restore before cleanup
```

Then confirm the checkout path still rebuilds with
`<built binary> install --repo <this checkout> --home .agents/scratch/<slug>/home3 --dry-run`
→ `[dry-run] the binary would be rebuilt from <checkout>`.

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
