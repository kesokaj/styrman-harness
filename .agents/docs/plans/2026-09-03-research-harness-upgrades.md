# Research — harness upgrades and optimisations   (2026-09-03)

## Question

How should styrman-harness be upgraded and optimised, given: no symlinks, overwrite-on-copy,
a purpose-built installer tool, install-time drift checking, and a fix for chronically stuck
subagent dispatch on both Claude Code and Antigravity.

Companion recon: `.agents/docs/plans/2026-09-03-harness-recon.md`.

## Verified environment

| Fact | Value |
|---|---|
| Claude Code | 2.1.258 |
| agy (Antigravity CLI) | 1.1.24 |
| Antigravity IDE | not installed on this machine |
| python3 / jq | 3.13.5 / jq-1.7 |

## Diagnosis: why dispatch gets stuck

### Claude Code

1. **`trail-check.sh` blocks the Stop hook on every turn during active development.**
   `dod_reason()` fires whenever source changed today and any of `changes/`, `reviews/`,
   `testing/` lacks a dated artifact. `stop_hook_active` caps it at one block per turn, but
   that is still one forced "dispatch reviewer + verifier + changelog now" round-trip at the
   end of *every* turn mid-feature. The DoD is a commit-boundary contract; the gate is
   day-grained and turn-triggered. That mismatch is the felt stall.

2. **`subagent-check.sh` matches artifacts by filename glob, so a dispatch prompt that names
   a path breaks it.** Required: `plans/$today-recon-*.md`. Observed today: the explorer was
   told by its dispatch prompt to write `2026-09-03-harness-recon.md`, which does not match.
   Under the repo's `settings.json` that agent would have been blocked and forced to loop.
   Verified live.

3. **The only loop guard is `stop_hook_active`, which is not in the documented SubagentStop
   payload.** Both scripts depend on it (`trail-check.sh:163`, `subagent-check.sh:80`). If the
   field is absent on SubagentStop, `subagent-check.sh` blocks unconditionally and there is no
   documented failsafe — it keeps the marker on block specifically expecting the flag next time.
   This is the hang vector.

4. No agent declares `tools: Agent`, and `tools: *` does not imply it — so no dispatcher exists.

### Antigravity

1. **The harness plugin is not registered.** `agy plugin list` → "No imported plugins";
   `agy agents` → empty. Files were copied into `~/.gemini/config/plugins/harness/` but `agy`
   requires `agy plugin install <path>`; the documented global path is
   `~/.gemini/antigravity-cli/plugins/<name>/`, absent here. `invoke_subagent reviewer` has
   nothing to resolve.
2. **Installed `hooks.json` is missing the whole `harness-gates` block** (PreInvocation
   resume-context, Stop trail-check). Only `harness-format` is deployed.
3. Antigravity docs, Known Issues: *"Specifying an unmapped or misspelled tool name in the
   `tools` list may cause the subagent process to hang during execution."* `derive-agy-agent.py`
   already guards this correctly and its header comment names the same risk — but `TOOL_MAP`
   is hardcoded and never validated against the installed CLI.

## Orchestration options compared

| Approach | Deterministic | Cross-vendor | Invocable from skill | Notes |
|---|---|---|---|---|
| **Dispatcher subagent** (`tools: Agent(...)`) | No — a model decides | No agy equivalent | Yes | Constructible: nesting allowed 3 layers, 20 concurrent. Adds a layer that can itself stall; violates "determinism never in model memory" |
| **Workflow tool** (JS orchestration) | Yes | Claude-only | Yes, when a skill's instructions say to call it | GA. Up to 1000 agents / 16 concurrent. Expensive. No agy counterpart |
| **Agent teams** | No | Claude-only | No | Experimental, `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`, interactive only |
| **Headless shell driver** (`claude -p --agent X`, `agy -p --agent X`) | Yes | Yes | Yes | Both CLIs support `--agent`, `-p`, `--output-format json`. Hard timeouts, real exit codes. No model in the dispatch decision |

Recommendation: **headless shell driver** as the backbone, Workflow tool as an opt-in
accelerator on the Claude side only.

## Packaging options (Claude side)

| Option | Install | Namespacing | Gains | Cost |
|---|---|---|---|---|
| **A. Plain files** (status quo) | copy to `~/.claude/{skills,agents}` | none — `/audit` | simplest | no versioning, no inventory, no eval |
| **B. skills-dir plugin** | copy to `~/.claude/skills/<name>/`, auto-loads as `<name>@skills-dir` | likely `harness:audit` | `claude plugin details` (component + token-cost inventory), enable/disable, `claude plugin eval` | namespacing must be validated |
| **C. Marketplace plugin** | `/plugin marketplace add`, `/plugin install` | `harness:audit` | real versioning, remote install on a fresh machine | most machinery; plugin cached, no `../` refs |

Option B is the best fit for "no symlinks, just overwrite" — it is still a file copy.
Validate namespacing with `claude plugin init` on a throwaway name before committing.

Agent Plugins 1.0 remains disqualified: hooks, agents and rules are explicitly out of scope
in v1, and Anthropic is not a maintainer.

## Recommendations

1. **One installer tool** — consolidate `setup.sh`, `claude/setup.sh`, `antigravity/setup.sh`,
   `merge-settings.py`, `derive-agy-agent.py` into `harness.py` with subcommands
   `install` / `doctor` / `derive`. Python 3 stdlib: JSON is the dominant data type, python3 is
   already a dependency, and bash would need jq for the same work.
2. **Overwrite, don't merge.** Deleting `merge-settings.py` also deletes its hook-merge bug and
   the unbounded `.bak.*` accumulation. Preserve user settings via `~/.claude/settings.local.json`
   rather than by merging into the deployed file.
3. **`agy plugin install` for the Antigravity side**, driven by the installer — copying is not
   registration.
4. **`doctor` folded into `install`**, run at the end of every install and wired to `SessionStart`:
   deployed-vs-repo hash comparison, every hook `command` exists and is executable, `agy plugin list`
   contains `harness`, `TOOL_MAP` names still valid.
5. **`dod-dispatch`** — compute owed rows, run each owner headless with a timeout, in parallel.
6. **Retune `trail-check.sh` to commit boundaries**, not turn boundaries, and make the artifact
   match content-based (frontmatter `owner:` + `date:`) rather than filename-glob-based.
7. **Add `maxTurns` to all six agents** and make `subagent-check.sh` fail open after one block
   when `stop_hook_active` is absent.

## Optimisations backlog

- `sandbox.*` settings — would subsume much of `destructive-guard.py`
- `promptCacheTtl` / `subagentPromptCacheTtl`
- `effort: max` on `security`, `effort` on the rest; `isolation: worktree` for `verifier`
- `skillOverrides` to retire the `grill-me` stub without deleting it
- `disableBundledSkills` — decide per skill where the harness duplicates a built-in
- Hook events not used: `SessionEnd`, `InstructionsLoaded`, `PostToolUseFailure`
- `VENDORED.md` under-reports: `gcloud-auth-verification`, `gke-basics`, `skill-repair`,
  `terraform` carry third-party license/author frontmatter and are unlisted
- Dangling reference: `antigravity/setup.sh:243` → `SUB_AGENTS.md` does not exist
- No CI, no tests, no enforced shellcheck; repo does not dogfood its own `.agents/` scaffold

## Is the project simple enough, with enough context?

**Contract: yes.** 74-line constitution (~2.3k tokens), 28 skill descriptions (~2.8k tokens).
~5k tokens always-on, well inside guidance. Do not optimise here.

**Machinery: no.** 5 install scripts + 2 Python helpers + 6 gate scripts, three byte-identical
copies of the constitution, zero tests. Recommendation 1 addresses this.

## Open questions

- Does `SubagentStop` actually deliver `stop_hook_active`? Undocumented; the harness's only
  loop guard depends on it. Needs an empirical test.
- Does the Antigravity IDE read `~/.gemini/config/plugins/` even though `agy` does not?
  Untestable here — IDE not installed.
- Does a skills-dir plugin namespace its skills? Must be validated before choosing option B.
