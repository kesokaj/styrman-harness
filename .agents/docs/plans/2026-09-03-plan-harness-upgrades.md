# Implementation plan — harness redesign   (2026-09-03, rev 2)

Research: `.agents/docs/plans/2026-09-03-research-harness-upgrades.md`
Recon: `.agents/docs/plans/2026-09-03-harness-recon.md`

Rev 2 replaces rev 1 after operator pushback: no Python, no sandbox, and a demand to question
the design rather than patch it.

Verification owner: **AI**, **User**, or **Both**.

## Measurements this rev is built on

| Measurement | Value | Source |
|---|---|---|
| Cost of one always-blocking SubagentStop hook | **$0.236 / 37.8 s / 8 forced continuations** for a subagent whose entire output was the word "pong" | `claude -p --output-format json`, T1 rig |
| Loop cap on a blocking hook | 9 hook firings, then clean exit | reproduced twice |
| Hook process startup: Go vs Python | **2.6 ms vs 22.8 ms** (~9×) | 30-run median, this machine |
| Hook process startup: bash / jq | 3.4 ms / 4.6 ms | same |
| Always-on context | constitution 74 lines ≈2.3 k tok; 28 skill descriptions ≈2.8 k tok | `wc`, frontmatter extraction |
| Per-run cost accounting available | `total_cost_usd`, `usage`, `modelUsage`, `subagent_stats` | `-p --output-format json` |

## The three design faults

**F1 — Gates coerce where they should inform.** `trail-check.sh` and `subagent-check.sh` both
work by *blocking*. A block is free when right and costs up to 8 model round-trips when wrong.
The blast radius is measured above. Blocking is justified only where correctness genuinely
depends on stopping the agent; everywhere else a gate should inject context or report.

**F2 — Every DoD row is an agent, including the rows that need no model.** The verifier's job is
"run the project's tests and report PASS/FAIL" — that is a script. Running it as a subagent pays
model prices to report an exit code. The changelog re-reads a diff the dispatcher already has.
Only reviewer, security, documentation and explorer genuinely need judgement.

**F3 — Subagents return prose, so context grows with every dispatch.** The explorer's
write-to-disk-return-summary pattern is right, but nothing caps the summary. Today's two
subagents each returned ~2 000 words into main context. This is the "context gets huge" symptom
directly.

Underneath all three: **nothing measures what dispatch costs**, so "0 benefit" has been a feeling
rather than a number. `total_cost_usd` per run makes it a number.

## Decisions

**D1 — One Go binary, not Python, not twelve scripts.** `harness` with subcommands replaces
`setup.sh` ×3, `merge-settings.py`, `derive-agy-agent.py`, `trail-check.sh`, `subagent-check.sh`,
`resume-context.sh`, `destructive-guard.py`, `format-hook.sh`, `feedback-capture.sh`,
`statusline.sh`. Rationale: operator's own language; 2.6 ms hook startup against Python's 22.8 ms
on a hook that fires on every Bash call; one deployable artifact; stdlib JSON and
`os/exec` + `context.WithTimeout` are exactly the primitives the dispatcher needs. Cross-compiled
ARM64 + AMD64, which the constitution already requires for images.

**D2 — `harness dispatch` is a router, not an agent.** It computes owed rows from the diff,
applies proportionality, runs script-rows as scripts, runs model-rows in parallel with hard
timeouts, caps their returns, and logs cost per row. No model decides whether to dispatch.

**D3 — Tier the rows by what they actually need.**

| Row | Needs a model? | Mechanism |
|---|---|---|
| verifier | No — unless it fails | script; escalate to a model only to explain a failure |
| changelog | Barely | haiku/flash, **seeded with the diff** so it never re-reads |
| reviewer | Yes | `/code-review` in-process for ordinary diffs; subagent for large ones |
| security | Yes, but rarely triggered | subagent, only when its trigger surface is touched |
| documentation | Yes, rarely | subagent |
| explorer | Yes | subagent — already correct |

**D4 — Cap subagent returns.** Every DoD agent's final message: verdict line, artifact path, at
most three findings. The detail already lives on disk; sending it twice is the context leak.

**D5 — No sandbox; drop `destructive-guard`.** The operator runs both CLIs inside an isolated
microVM, so host-filesystem guarding is friction without value. **Keep secret scanning** — a
microVM protects the host, not a remote: a secret pushed to GitHub escapes the VM. Keep
`promptCacheTtl` / `subagentPromptCacheTtl`; those are token economics, which the operator
explicitly cares about.

**D6 — Plain files, overwrite.** Unchanged from rev 1 and confirmed by T3: a plugin would
namespace all 28 skills and all 6 agents.

## Open design question for the operator

The trail rows split by purpose. **Resumability** rows — `changes/`, `handoffs/`,
`DECISIONS.md` — are cheap and high value; they are what makes a zero-context agent able to
resume. **Accountability** rows — reviewer, security, verifier reports — are expensive and
largely re-derivable from the diff. That trio exists to make work auditable by someone who did
not do it. For a solo operator, most of the value is caught more cheaply by `/code-review`.

Recommend making the accountability rows **proportional** (fire on diff size / trigger surface)
rather than mandatory on every source change. This is a change to `REQUIREMENTS.md`-level
intent, so it needs an explicit decision. *Owner: User.*

---

## P1 — The binary

**1. `harness` Go module**, subcommands: `install`, `doctor`, `derive`, `hook <event>`,
`dispatch`, `resume`, `statusline`. Cross-compiled arm64+amd64. *Verify: AI (`go test -race`,
`go vet`, both arch builds).*

**2. Delete the twelve superseded scripts** once each has a passing test against the binary.
*Verify: AI.*

**3. `install`** — overwrite everything it manages, no merge, no `.bak`. Installs the binary to
`~/.agents/bin/harness`, deploys skills/agents/CLAUDE.md/settings.json, runs
`agy plugin install`. *Verify: AI (dry-run against a temp HOME).*

**4. `doctor`**, run at the end of `install` and on `SessionStart` — content-hash drift,
every hook `command` present and executable, `import_manifest.json` contains `harness`,
`TOOL_MAP` valid for the installed agy. *Verify: AI.*

## P2 — Dispatch

**5. `harness dispatch`** — owed rows from the diff, proportionality filter, script-rows as
scripts, model-rows in parallel with `context.WithTimeout`, capped returns, one compact table
back to main context. *Verify: Both.*

**6. Cost ledger** — append `total_cost_usd`, `usage`, duration and verdict per row to
`.agents/docs/testing/dispatch-ledger.jsonl`. Makes "0 benefit" measurable. *Verify: AI.*

**7. Verifier becomes a script** that runs the project's declared test commands; a model is
invoked only to explain a failure. *Verify: Both.*

**8. `trail-check` → inform, not block.** Move from Stop-hook blocking to
`UserPromptSubmit` context injection plus a commit-boundary check. Keep blocking only where
correctness depends on it. *Verify: Both — User judges the cadence.*

**9. Artifact matching by frontmatter `owner:`/`date:`**, not filename glob. Today's live
failure: `2026-09-03-harness-recon.md` vs the required `-recon-` infix. *Verify: AI.*

**10. `maxTurns` on all six agents; `effort: max` on security; `isolation: worktree` on
verifier; return-cap wording in every agent body.** *Verify: AI.*

## P3 — Hygiene

**11.** `VENDORED.md` provenance for `gcloud-auth-verification`, `gke-basics`, `skill-repair`,
`terraform`. *Verify: User.*
**12.** Drop the dangling `SUB_AGENTS.md` reference. *Verify: AI.*
**13.** Retire the `grill-me` stub via `skillOverrides`. *Verify: AI.*
**14.** Root `AGENTS.md` / `REQUIREMENTS.md` / `ROADMAP.md` / `TODO.md` — the harness currently
fails its own Resume Protocol. *Verify: AI.*
**15.** CI: `go test -race`, `go vet`, shellcheck on what shell remains, `doctor` smoke test.
*Verify: AI.*

## Out of scope

Sandbox settings (D5). Plugin/marketplace packaging (T3). Agent Plugins 1.0 (excludes hooks and
agents). Workflow tool as the dispatch backbone — Claude-only, no agy counterpart.

## Still unverifiable here

Antigravity **IDE** plugin discovery — only `agy` 1.1.24 is installed on this machine.
*Owner: User.*
