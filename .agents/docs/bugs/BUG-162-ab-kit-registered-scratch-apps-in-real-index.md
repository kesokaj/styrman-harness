---
owner: main
date: 2026-10-07
id: BUG-162
status: fixed
area: testing (headless A/B run kits)
---
# BUG-162: the A/B kit's harness runs registered scratch apps in the real global search index

## Symptom
On 2026-10-07 the main agent ran headless `claude -p` A/B runs (plan `2026-10-07-lean-process`) from a kit it wrote in the session scratchpad. Two runs scaffolded a project, and their agents ran `sh-agent search`/`index`:
- `harness-full`, starting 08:07
- `lean-full`, starting 09:15

Both registered their scratch app in the operator's real global search index, `~/.config/sh-agent/search.db`, which was modified at 08:09. A `strings` read of a copy showed `.../scratchpad/abtest/runs/harness-full/app`. The nightly `sh-agent dream --all` reads every indexed project.

## Repro
Run a harness-loaded `claude -p` session in a scaffolded project. Isolate only `CLAUDE_CONFIG_DIR`, leave `HOME` real, and put no wrapper on `PATH`. The agent's `sh-agent search` registers the project globally.

## Root cause
The kit isolated the CLI's config (`CLAUDE_CONFIG_DIR`) but not the harness binary's state. It did not source `.agents/docs/testing/kit-preamble.sh`, whose `PATH` wrapper pins `--home` on every `sh-agent` call (BUG-107). The kit was written in the session scratchpad as an experiment, so the project rule "every run kit sources `testing/kit-preamble.sh` first" (AGENTS.md § Testing) was never read against it.

## Fix
- Each scratch app's `.agents/docs` was renamed (`docs-archived`), then `sh-agent index --global` was run. The search spec's "no longer has .agents/docs" row drops such a registration.
- Verified: `sh-agent search --global crm booking` returns only `styrman-harness` records. Before the fix, both apps' plans and reviews matched.
- The kit's credential symlinks were deleted by name, and a `find` over the scratch trees for credential names returned nothing.
- The kit is not committed (plan D14): rebuilding it on the preamble is a follow-up.

## Prevention
Any headless run of a harness-loaded CLI is a run kit, wherever its script lives: it sources `kit-preamble.sh` first, or puts an `sh-agent` wrapper pinning `--home <sandbox>` first on `PATH`. Isolating `CLAUDE_CONFIG_DIR` alone isolates the CLI, not the harness.

## Related
BUG-107 (the wrapper this kit skipped), BUG-106 (fixtures under a non-checkout parent), BUG-160 (credentials in scratch homes). Plan `.agents/docs/plans/2026-10-07-lean-process.md` D14.
