# Roadmap — styrman-harness

## Now

- Operator deploys the Go harness by hand (`./setup.sh`) and confirms `harness doctor`
  passes on the real machine.
- Confirm the proportionality thresholds feel right in daily use; tune
  `reviewLines` / `reviewFiles` in each project's `.agents/harness.json`.
- Read the first weeks of `dispatch-ledger.jsonl` and decide whether any
  accountability row is still not worth its price.

## Next

- Provenance for the four undocumented vendored skills (`gcloud-auth-verification`,
  `gke-basics`, `skill-repair`, `terraform`) — needs the operator, who knows how they
  were obtained.
- Decide `disableBundledSkills` per skill where the harness duplicates a Claude Code
  built-in.
- Wire `SessionEnd` and `PostToolUseFailure` if a use emerges; both are currently
  unused, and neither should be added without one.
- Confirm whether the Antigravity **IDE** loads a plugin from
  `~/.gemini/config/plugins/` without an explicit `agy plugin install`. Untested — no
  IDE on the current machine.

## Later

- Reconsider the Workflow tool as an opt-in accelerator for `/preflight` on the Claude
  side, once the shell-driven dispatcher has a track record.
- Revisit Agent Plugins if a later version of the spec covers hooks and agents.
