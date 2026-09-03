# Roadmap — styrman-harness

## Now

Nothing in flight. The harness is built, installed on the operator's machine, and
passing its own `doctor` at 10/10.

## Next

Nothing planned. The next change should come from using it — a gate that fires when it
should not, a row that costs more than it returns, a vendor that moves underneath it.
`.agents/docs/testing/dispatch-ledger.jsonl` records what every dispatch costs, so
"this row is not worth its price" can be settled with numbers rather than impressions.

## Later

Two things worth revisiting only if circumstances change:

- **Workflow-tool orchestration** for `/preflight` on the Claude Code side, if the
  shell-driven dispatcher ever proves too limited. It is Claude-only, which is why it
  was not chosen; that trade only shifts if Antigravity gains an equivalent.
- **Agent Plugins**, if a later version of the spec covers hooks and agents. Version 1
  packages Skills and MCP only — most of this harness falls outside it.
