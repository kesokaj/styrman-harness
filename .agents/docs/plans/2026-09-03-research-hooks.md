---
owner: main
date: 2026-09-03
---
# Research — hooks for the harness   (2026-09-03)

## Question

The operator may want hooks after all: (1) a deploy gate that verifies everything is documented and up to date, (2) a warning when the wrong model is in use, (3) a loop protector. What do Claude Code and the Antigravity CLI actually offer, what would each hook look like on both, what does it cost, and what is worth building?

Constraints that shape the answer: zero scripts in the repo (any hook is a subcommand of the `sh-agent` binary), no CI, both vendors served alike, and the previous harness's hooks were reverted because a wrong block cost eight forced continuations and a degraded agent (old ADR-002, R4). Whatever is built must fail open and block only at a boundary the operator would block at by hand.

## Summary

Both CLIs have a `PreToolUse` hook that can deny, ask, or allow a tool call from a JSON payload on stdin, and a `Stop` hook that can refuse a stop. That is enough for all three wishes. Claude Code additionally exposes the model to hooks (`SessionStart` sometimes, `PostModelSwitch` always) and has an `if` filter so a hook only spawns for matching commands; Antigravity puts `modelName` in every payload and has `PreInvocation` for injecting a message before each model call. Neither has built-in repeated-call detection.

Recommendation, in order of value: build the **loop protector** and the **model warning** (small, deterministic, session-wide), and build the **deploy gate as "a fresh `/preflight` SHIP report for HEAD must exist"** rather than as a hook that judges documentation itself — a hook can check that the check happened; only the `preflight` skill and the `documentation` agent can check content. Do **not** build a Stop-hook that blocks until docs are written: that is the design that was reverted. All three are subcommands of `sh-agent` (`sh-agent hook <name>`), so the repo stays at zero scripts. The one real cost is that Claude Code hooks live in `~/.claude/settings.json`, so the installer would have to add and remove a `hooks` block there — the operator has said no settings changes before; this is the decision to make.

## What each vendor offers

| | Claude Code | Antigravity CLI (agy 1.1.24) |
|---|---|---|
| Config location | `~/.claude/settings.json` (user), `.claude/settings.json` / `.claude/settings.local.json` (project), plugin `hooks/hooks.json`, subagent frontmatter `hooks:` (scoped to that agent), skill frontmatter `hooks:` (registered when the skill is invoked, kept for the session) | one `hooks.json` in the customization root: `~/.gemini/config/hooks.json` (global) or `.agents/hooks.json` (workspace); named hooks, merged by name; `"enabled": false` per hook |
| Events | `SessionStart`, `SessionEnd`, `UserPromptSubmit`, `PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PostToolBatch`, `PermissionRequest`, `PermissionDenied`, `Stop`, `StopFailure`, `SubagentStart`, `SubagentStop`, `PreModelSwitch`, `PostModelSwitch`, `PreCompact`, `PostCompact`, `Notification`, `FileChanged`, `CwdChanged`, and more | `PreToolUse`, `PostToolUse`, `PreInvocation` (before each model call), `PostInvocation` (after tool calls finish), `Stop` |
| Matcher | tool name, `A\|B`, regex; plus `if: "Bash(git push *)"` in permission-rule syntax so the hook only spawns for matching arguments | tool name regex on `PreToolUse`/`PostToolUse` (`run_command`, `browser_.*`, `*`) |
| Input (stdin JSON) | `session_id`, `transcript_path`, `cwd`, `permission_mode`, `hook_event_name`; `PreToolUse`: `tool_name`, `tool_input`, `tool_use_id`, `effort.level`; `Stop`: `stop_hook_active`, `last_assistant_message`; `SessionStart`: `reason`, `model` (not always); `PostModelSwitch`: the model | common: `conversationId`, `workspacePaths`, `transcriptPath`, `artifactDirectoryPath`, **`modelName`**; `PreToolUse`: `toolCall{name,args{CommandLine…}}`, `stepIdx`; `PostToolUse`: `error`; `PreInvocation`: `invocationNum`, `initialNumSteps`; `Stop`: `executionNum`, `terminationReason` (`model_stop`, `max_steps_exceeded`, `error`), `fullyIdle` |
| Output | exit 0 = continue (stdout becomes context for `SessionStart`/`UserPromptSubmit`); exit 2 + stderr = block; or JSON: `hookSpecificOutput.permissionDecision` allow/deny/ask/defer, `permissionDecisionReason`, `additionalContext`, `updatedInput`; `decision: "block"` + `reason` on `Stop`/`PostToolUse`; `systemMessage` | JSON only: `PreToolUse` → `decision` allow/deny/ask/force_ask, `reason`, `overwrite`; `PostToolUse` → `{}`; `PreInvocation`/`PostInvocation` → `injectSteps[{ephemeralMessage\|userMessage\|toolCall}]`, `terminationBehavior` force_continue/terminate; `Stop` → `decision: "continue"` + `reason` |
| Hook types | `command`, `http`, `mcp_tool`, `prompt` (a model judges), `agent` | `command` only |
| Timeout | per hook; defaults 10 min for commands, 30 s for prompt-submit/model-switch | per hook, default 30 s; hooks run synchronously |
| Working dir / command | runs with the session cwd | runs in the directory holding `hooks.json`, via `sh -c`, `~` expanded — so the command must be `~/.local/bin/sh-agent …` or on PATH |
| Built-in loop limits | subagent `maxTurns`; `/goal` evaluator; `stop_hook_active` so a Stop hook is not re-invoked forever | `max_steps_exceeded` termination; 60 s default subagent interaction timeout |
| Model control | `model`, `enforceAvailableModels`, `fallbackModel` in settings; `PreModelSwitch` hook can deny a switch; status line JSON has `model.id`/`display_name` | `--model` flag; `modelName` in every hook payload (`auto` when not pinned) |

Sources: https://code.claude.com/docs/en/hooks-guide.md, https://code.claude.com/docs/en/hooks.md, https://code.claude.com/docs/en/sub-agents.md, https://code.claude.com/docs/en/skills.md, https://code.claude.com/docs/en/settings-reference.md, https://code.claude.com/docs/en/statusline.md, https://code.claude.com/docs/en/goal.md; the "Lifecycle Hooks (`hooks.json`)" guide embedded in the agy binary and `agy changelog`.

## The three hooks, designed for both vendors

All three are `sh-agent hook <name>` subcommands reading stdin JSON and printing the vendor's JSON. Vendor is detected from the payload (`hook_event_name` vs `conversationId`). Every hook fails open: any error, unknown payload, or timeout means "allow, say nothing".

### 1. Deploy gate — "everything documented and up to date"

**What a hook can and cannot decide.** Whether the README, spec, and changelog *match the code* is a judgement; a hook can only check facts on disk. The facts worth checking before a deploy: the tree is clean, HEAD is what will be deployed, and a release gate has run on exactly this HEAD. `/preflight` already does the judgement part (docs freshness via the `documentation` agent, trail freshness, tests, review) and writes `reviews/<date>-preflight-*.md` with `Range reviewed: <base>..<head>` and `VERDICT: SHIP`. So:

- **Trigger:** Claude `PreToolUse`, `matcher: Bash`, `if: Bash(gcloud run deploy *)` plus `Bash(gcloud run services update *)`, `Bash(docker push *)`, `Bash(terraform apply *)`, `Bash(kubectl apply *)`; Antigravity `PreToolUse` on `run_command`, the same regexes applied to `args.CommandLine` inside the hook.
- **Check:** `git status --porcelain` empty; newest `.agents/docs/reviews/*-preflight-*.md` has a `Range reviewed:` line ending in the current `HEAD` short SHA and a `VERDICT: SHIP` line. A `docker push` to `pkg.dev` also requires the image tag to equal the HEAD short SHA when the tag is visible in the command.
- **Result:** pass → allow silently. Fail → `ask` (Claude `permissionDecision: ask`, agy `decision: ask`) with the reason ("no preflight report for HEAD a1b2c3d — run /preflight, or approve to deploy anyway"). Ask, not deny: the operator stays in control, and a deliberate hot-fix deploy costs one keypress instead of a fight.
- **Not gated:** `git push`. Git is backup; a gate on push blocks backups.
- **Cost:** one git call and one directory scan per deploy command; the hook does not spawn for other commands on Claude (the `if` filter), and costs a few milliseconds per command on agy.

Alternative considered and rejected: a `Stop` hook that blocks until a changelog entry exists. It fires every turn, it is wrong whenever the turn was not the end of a task, and it is the exact mechanism that was reverted. A softer variant — `Stop` hook that only *warns* ("source changed since HEAD, no changes/ entry today") via `systemMessage`, never blocks — is cheap and could be added later if the prose Definition of Done is being ignored.

### 2. Wrong-model warning

- **Claude:** `SessionStart` (payload `model` when present) and `PostModelSwitch` (always). If the model id is not in the allowed list, print one line to stdout (exit 0 → it lands in context) and set `systemMessage` so the operator sees it: "Model is claude-haiku-4-5; expected claude-fable-5-1 or claude-opus-5". Optionally `PreModelSwitch` with `permissionDecision: ask` to confirm a switch to a model outside the list. No blocking.
- **Antigravity:** `PreInvocation` fires before every model call and carries `modelName`. On `invocationNum == 1` only, if `modelName` is not allowed (or is `auto` when a pin is required), inject `{"ephemeralMessage": "…"}`; nothing on later invocations.
- **Allowed list:** given as arguments in the hook config itself (`sh-agent hook model --allow claude-fable-5-1,claude-opus-5` / `--allow pro`), so there is no extra config file; the installer templates the lists from one place in the repo (`context/` or a small `hooks/` directory the installer reads).
- **Cheaper alternative on Claude:** no hook; pin `model` in `~/.claude/settings.json` (already done: `claude-fable-5-1[1m]`) and read the status line. The hook only adds value if sessions get started with `--model` or switched with `/model`.

### 3. Loop protector

- **Signal:** the same tool call (tool name + argument hash) repeated three times in a row **with no file-changing call in between**. The "no edit in between" rule removes the common false positive: running `go test ./...` three times is fine when an `Edit` sits between the runs. `PostToolUse` records the result hash as well; identical input and identical output twice in a row is a loop by definition.
- **Claude:** `PreToolUse` on every tool (`matcher: ""`) reads `session_id`, `tool_name`, `tool_input`; `PostToolUse` appends the result hash. State is one small file per session under the CLI's own temp dir (`$TMPDIR/sh-agent/<session_id>`), so nothing persists. On a loop: `permissionDecision: ask`, reason "same call for the 3rd time with no change in between — loop?". `Stop` is left alone (no Stop hook), so no risk of the forced-continuation spiral.
- **Antigravity:** `PreToolUse` on `*` with `toolCall.name` + `args` hash keyed by `conversationId`; `PostToolUse` has no result body, so only the input rule applies. On a loop: `decision: ask` with the reason. `PostInvocation` could add `terminationBehavior: terminate` when `invocationNum` passes a hard ceiling with no edits since, as a second line of defence.
- **Subagents:** already capped by `maxTurns` (20–40) in `agents/*.md`; the hook applies to them too on Claude.
- **Cost:** one hash and one small file append per tool call, a few milliseconds. No model calls.

## What it would take

| Piece | Where | Size |
|---|---|---|
| `sh-agent hook deploy-gate\|model\|loop` | `internal/hook/` + `cmd/sh-agent` | ~350 lines Go incl. vendor payload detection, plus tests against recorded payloads |
| Claude config | `hooks` block in `~/.claude/settings.json`; installer merges by our command prefix (`~/.local/bin/sh-agent hook`) and removes only those entries on uninstall | ~80 lines Go; the one settings write the harness has avoided so far |
| Antigravity config | named entries `sh-agent-deploy`, `sh-agent-model`, `sh-agent-loop` in `~/.gemini/config/hooks.json`; installer merges by name | ~40 lines Go |
| Repo source of truth | `hooks/` directory with one JSON per hook in a vendor-neutral shape (event, matcher/if, command args) that the installer renders into both formats — so a new hook is a file, not code, once the subcommand exists | small |
| Docs | README, `context/AGENTS.md` (the gate's "run /preflight first"), ADR | small |

Option B for Claude that avoids the global settings write: `/init-project` writes the hooks into the project's `.claude/settings.json`. Per project, versioned with the project, but every project must be initialised for the hooks to exist, and the model warning would not fire outside a project.

Option C: skill-scoped hooks. `cloud-run-deploy` and `preflight` could carry the deploy gate in their SKILL.md `hooks:` frontmatter, active for the session once the skill is invoked. No settings write at all, but the gate exists only after the skill was loaded — an agent that runs `gcloud run deploy` without loading the skill is not gated. Fine as a belt when the braces exist, not on its own.

## Risks

- A hook that misparses a payload and denies is worse than no hook. Hence ask-not-deny everywhere, fail-open, and payload fixtures in tests.
- Claude's `if` syntax and hook payload fields have changed across versions; the `model` field on `SessionStart` is documented as "not always present". Pin what the code expects and test against fixtures recorded from the installed version.
- Antigravity hooks run synchronously with a 30 s default timeout; a slow git repo could stall a `run_command`. Keep the deploy gate to `git status` and a directory scan.
- The loop protector's state file must be per session and cleaned by the OS temp dir; never write it into the repo.
- None of this is verified live on either CLI yet; the live check of the base harness comes first.

## Open questions for the operator

1. Build all three, or only the loop protector and the model warning?
2. Is a merged `hooks` block in `~/.claude/settings.json` acceptable (installer-managed, removed on uninstall), or per-project via `/init-project` (Option B)?
3. Deploy gate: `ask` (recommended) or hard `deny`?
4. Allowed models: Claude `claude-fable-5-1`, `claude-opus-5`? Antigravity `pro` only, or is `auto` fine?
5. Loop threshold: 3 identical calls with no edit in between (recommended), or stricter?
