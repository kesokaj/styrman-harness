# Recon — harness full survey (2026-09-03)

## Question
Survey the styrman-harness repo (a dotfiles-style installer that deploys a shared agent
contract into both Claude Code `~/.claude` and Google Antigravity `~/.gemini` + plugins):
install mechanism, what's genuinely shared vs duplicated between vendors, hook/gate wiring,
skills inventory, agent definitions, contract consistency, and gaps/risks.

## Summary
The repo has one real source of truth (`shared/AGENTS.md`) that is copied byte-for-byte into
`claude/CLAUDE.md` and `antigravity/GEMINI.md` — verified identical by diff. Claude Code's six
subagents (`claude/agents/*.md`) are the source; Antigravity's are **derived**, not duplicated,
via `shared/scripts/derive-agy-agent.py`, and the derived bodies diff byte-identical to the
Claude originals (frontmatter differs by design). The install is copy/rsync-based (never
symlinked), backs up clobbered files with a timestamp suffix, and merges `settings.json` via a
custom deep-merge (`shared/scripts/merge-settings.py`) — but that merge has a real bug for the
`hooks` object: because Claude hook entries are dicts, not strings, a changed `timeout` or added
step on an existing matcher is treated as a brand-new list item and appended rather than
replacing the old one, so a stale hook can coexist with its replacement after a merge-redeploy.
This machine's own `~/.claude` is live proof the "deployed verbatim" claim is aspirational, not
enforced: `~/.claude/CLAUDE.md`/`~/.gemini/GEMINI.md` are stale (missing the "Release Ownership"
section and other later edits present in-repo), and `~/.claude/settings.json` is missing the
`PreToolUse` (destructive-guard) and `SubagentStart`/`SubagentStop` (subagent-check) hook blocks
entirely — `~/.agents/bin/` on this machine has only `resume-context.sh` and `trail-check.sh`,
not `destructive-guard.py` or `subagent-check.sh` — meaning `claude/setup.sh` has not been rerun
since those gates were added, and nothing detects or warns about this drift. Antigravity also
gets materially weaker guardrails by design: `destructive-guard.py` and `subagent-check.sh` are
Claude-only (never deployed or wired for agy), so Antigravity has no destructive-command
confirmation gate and no per-subagent artifact enforcement, only the coarser Stop-hook
`trail-check.sh --agy`. `antigravity/setup.sh:243` also logs a pointer to a `SUB_AGENTS.md` file
that does not exist anywhere in the repo. Skills are clean by the requested thresholds (no
SKILL.md over 500 lines; largest is `terraform` at 393) but `VENDORED.md` under-reports
provenance: `gcloud-auth-verification`, `gke-basics`, `skill-repair`, and `terraform` all carry
third-party `license`/`author`/`publisher` frontmatter (Google, Anton Babenko) yet are absent
from the vendored-skills manifest, which lists only 5 skills.

## Relevant files
- `/home/ostling-simon/projects/styrman-harness/setup.sh` — top-level dispatcher; parses
  `--target`/`--dry-run`, defaults to both vendors, execs `<target>/setup.sh` in sequence.
- `/home/ostling-simon/projects/styrman-harness/claude/setup.sh` — Claude Code installer.
- `/home/ostling-simon/projects/styrman-harness/antigravity/setup.sh` — Antigravity installer.
- `/home/ostling-simon/projects/styrman-harness/shared/scripts/merge-settings.py` — deep-merge
  used to fold repo `settings.json` into an existing user `settings.json` without clobbering it.
- `/home/ostling-simon/projects/styrman-harness/shared/scripts/derive-agy-agent.py` — turns a
  Claude agent `.md` into an Antigravity plugin agent `.md` (frontmatter remap + skill path
  rewrite); the actual mechanism behind "two vendor targets."
- `/home/ostling-simon/projects/styrman-harness/shared/AGENTS.md`,
  `claude/CLAUDE.md`, `antigravity/GEMINI.md` — the constitution; verified identical (see Flows).
- `/home/ostling-simon/projects/styrman-harness/claude/settings.json` — Claude hook wiring +
  permissions; the only place all 6 hook events are declared.
- `/home/ostling-simon/projects/styrman-harness/antigravity/plugins/harness/hooks.json` — agy's
  hook wiring (2 events only: `PostToolUse` format hook, `PreInvocation`/`Stop` gates).
- `/home/ostling-simon/projects/styrman-harness/antigravity/plugins/harness/plugin.json` —
  minimal plugin manifest (name/version/description only).
- `/home/ostling-simon/projects/styrman-harness/shared/scripts/{trail-check.sh,resume-context.sh,
  subagent-check.sh,feedback-capture.sh,format-hook.sh,destructive-guard.py}` — the six gate/hook
  scripts; see Flows for exact wiring per script.
- `/home/ostling-simon/projects/styrman-harness/claude/agents/*.md` (6 files) — source-of-truth
  subagent definitions (changelog, documentation, explorer, reviewer, security, verifier).
- `/home/ostling-simon/projects/styrman-harness/antigravity/plugins/harness/agents/*.md` (6
  files) — derived copies; bodies are byte-identical to `claude/agents/*.md` (confirmed by diff).
- `/home/ostling-simon/projects/styrman-harness/shared/skills/` (28 skill dirs) — shared skill
  library deployed verbatim (via `rsync -av --delete`) to both `~/.claude/skills/` and
  `~/.gemini/config/skills/`.
- `/home/ostling-simon/projects/styrman-harness/shared/skills/VENDORED.md` — provenance manifest,
  incomplete (see Risks).
- `/home/ostling-simon/projects/styrman-harness/shared/templates/` — scaffolding for
  `/init-project` (`AGENTS.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`, `.agents/docs/*`
  skeleton) — not used by this repo on itself (see Gaps).
- `/home/ostling-simon/.claude/CLAUDE.md`, `/home/ostling-simon/.gemini/GEMINI.md`,
  `/home/ostling-simon/.claude/settings.json`, `/home/ostling-simon/.agents/bin/` — this
  machine's actual deployed state, used as live evidence of install/drift behavior.

## Flows

### 1. Install/deploy mechanism
- `setup.sh:48-73` — no-arg run deploys both `antigravity` and `claude` in that order; each
  vendor script must be executable or the whole run aborts (`setup.sh:64-67`). `--dry-run`
  propagates to both.
- **Copy, never symlink.** Every deploy step uses `cp`, `cp -a`, or `rsync -av[--delete]` — no
  `ln -s` anywhere in either setup script. Consequence: local edits to files under `~/.claude` or
  `~/.gemini` are NOT reflected back into the repo (except the one exception below), and a
  re-run of setup silently overwrites hand-edited deployed files (after a timestamped backup).
- `claude/setup.sh` steps, in order:
  1. `:69-77` refresh in-repo `claude/CLAUDE.md` FROM `shared/AGENTS.md` (`cp`) — i.e. the
     in-repo Claude copy is regenerated from source every run, not hand-maintained.
  2. `:82-87` backup + `cp` `claude/CLAUDE.md` → `~/.claude/CLAUDE.md`.
  3. `:92-101` `rsync -av --delete shared/skills/ → ~/.claude/skills/` — full mirror, so a skill
     removed from the repo is removed on disk too (destructive but source-controlled).
  4. `:104-128` retires 5 named legacy slash-command files (`audit.md` etc.) from
     `~/.claude/commands/` with a backup+`rm`, then overlays `claude/commands/` if non-empty
     (repo currently ships none).
  5. `:135-156` settings.json: if none exists, plain `cp`; if one exists, backup then
     `python3 shared/scripts/merge-settings.py TARGET SOURCE OUT` → `mv` over target (see Risks
     for the hooks-merge bug).
  6. `:164-181` subagents: retires `testing.md` (an agent since replaced by `verifier`) with
     backup+`rm`, then `rsync -av` (no `--delete`) `claude/agents/` → `~/.claude/agents/` — user-
     added agents are preserved; a repo agent renamed/removed otherwise lingers until explicitly
     retired here.
  7. `:188-210` copies `claude/scripts/statusline.sh` and `shared/scripts/{format-hook.sh,
     feedback-capture.sh}` into `~/.claude/`, `chmod +x` each, with backups.
  8. `:217-231` copies 4 gate scripts (`resume-context.sh`, `trail-check.sh`,
     `subagent-check.sh`, `destructive-guard.py`) into `~/.agents/bin/`, `chmod +x`, with
     backups; warns (not fails) if `jq` is missing, since the gates are jq-dependent and silently
     inert without it.
- `antigravity/setup.sh` steps, in order:
  1. `:82-90` same refresh pattern: `antigravity/GEMINI.md` regenerated from `shared/AGENTS.md`.
  2. `:97-100` backup + `cp` → `~/.gemini/GEMINI.md`.
  3. `:107-114` `rsync -av --delete shared/skills/ → ~/.gemini/config/skills/`.
  4. `:122-127` retires the entire legacy `~/.gemini/config/global_workflows/` tree (backup +
     `rm -rf`) — workflows migrated into skills.
  5. `:134-160` **derives** plugin agents: for each `claude/agents/*.md`, runs
     `derive-agy-agent.py` into a temp dir; only swaps the temp dir into
     `antigravity/plugins/harness/agents/` once every file derived successfully (atomic:
     half-derived output never ships) — the deriver exits non-zero on any unmapped `tools:`
     entry (`derive-agy-agent.py:68-70`), which is itself a fail-fast guard against a broken agy
     agent shipping silently.
  6. `:170-185` `rsync -av --delete antigravity/plugins/harness/ → ~/.gemini/config/plugins/
     harness/` (harness-owned, full mirror); any *other* plugin dir syncs without `--delete`.
  7. `:192-198` copies `shared/scripts/format-hook.sh` → `~/.gemini/format-hook.sh`.
  8. `:205-219` copies only `resume-context.sh` and `trail-check.sh` (NOT `subagent-check.sh` or
     `destructive-guard.py`) into `~/.agents/bin/` — same warn-on-missing-`jq` behavior.
  9. `:225-230` removes a stale `~/.gemini/settings.json` IF it still carries the harness's old
     `_comment` marker string — a narrowly-scoped, non-destructive cleanup of a legacy artifact.
  10. `:245-280` optional post-deploy verification, only if the `agy` CLI is present: validates
      the plugin parses (`agy plugin validate`) and cross-checks skill count against
      `agy -p "/skills"` output — this is the only place either setup script verifies its own
      work landed.
- **Idempotency:** re-running either script is safe in the sense that it always backs up before
  overwriting (`backup()` in both scripts, timestamped `.bak.YYYYMMDD-HHMMSS`), and rsync mirrors
  converge. It is NOT idempotent in the sense of leaving the target unchanged on a no-op re-run —
  every re-run creates a fresh backup file even when content is unchanged, so repeated runs
  accumulate `.bak.*` files indefinitely (no pruning logic found anywhere in the repo).
- **What happens to user-local edits under `~/.claude`:** anything under `skills/`, `agents/`
  (rsync w/ delete or without), `CLAUDE.md`, `statusline.sh`, `format-hook.sh`,
  `feedback-capture.sh`, and the 4 gate scripts is unconditionally overwritten (after backup) on
  every run — there is no diff/prompt/merge for these paths. Only `settings.json` gets a real
  merge; only `claude/commands/` and non-`harness` `~/.gemini/config/plugins/*` subtrees are
  additive (no `--delete`).

### 2. Two vendor targets — shared vs duplicated
- **Constitution**: `shared/AGENTS.md` is the sole source; `claude/CLAUDE.md` and
  `antigravity/GEMINI.md` are *regenerated* copies produced by each `setup.sh`'s step 1
  (`cp shared/AGENTS.md → claude/CLAUDE.md` / `→ antigravity/GEMINI.md`). Verified via `diff`:
  all three files are byte-identical right now.
- **Agents**: `claude/agents/*.md` is the source; `antigravity/plugins/harness/agents/*.md` is
  generated by `derive-agy-agent.py`. Verified via diff (stripping frontmatter) that all 6 agent
  **bodies** are byte-identical between vendors. Frontmatter differs by design:
  `derive-agy-agent.py:59-77` keeps `name`/`description`, adds `subagent: true`/
  `mainAgent: false`, maps `model:` (`opus→pro`, `haiku→flash`, else omitted — confirmed
  `documentation`/`explorer`/`verifier`, all `model: sonnet` in Claude, carry no `model:` line in
  the derived agy file), maps `tools:` via a fixed table (`Read→view_file`, `Grep→grep_search`,
  `Glob→find_by_name`, `Bash→run_command`, `Write→write_to_file`, `Edit→replace_file_content`;
  an unmapped tool is a hard error, never a silent guess), rewrites `skills:` entries to
  `~/.gemini/config/skills/<name>` paths, and rewrites any `~/.claude/skills/` path in the body
  text to `~/.gemini/config/skills/`.
- **Skills**: fully shared, not duplicated or derived — `shared/skills/` is `rsync`'d verbatim
  (identical bytes) to both `~/.claude/skills/` and `~/.gemini/config/skills/`. No vendor-specific
  skill variant mechanism exists.
- **Hooks**: NOT shared — `claude/settings.json` and
  `antigravity/plugins/harness/hooks.json` are two independently hand-written files with
  different event names, matcher syntax, and JSON shapes (Claude Code hook events vs. agy's
  `harness-format`/`harness-gates` grouping and `PreInvocation`/`Stop`/`PostToolUse` naming). Only
  `format-hook.sh`, `resume-context.sh`, and `trail-check.sh` are wired on both sides; Claude
  additionally wires `destructive-guard.py` (`PreToolUse`), `subagent-check.sh`
  (`SubagentStart`/`SubagentStop`), and `feedback-capture.sh` (`UserPromptSubmit`) — none of
  which exist in `hooks.json` or are deployed for Antigravity at all.

### 3. Hook/gate wiring (Claude Code, `claude/settings.json`)
| Hook event | Matcher | Command | What it enforces | Exit-code semantics |
|---|---|---|---|---|
| `PreToolUse` (`:12-22`) | `Bash` | `$HOME/.agents/bin/destructive-guard.py` | Blocks/asks on destructive shell commands (`rm`, `git clean`, `git reset --hard`, etc.) outside an agent scratch root; emits `permissionDecision: allow/ask` via `hookSpecificOutput` JSON (`destructive-guard.py:222-225`); never emits `deny`; any internal error degrades to `ask` (`:359-360`), never silently allows. Always exits 0. |
| `SessionStart` (`:24-34`) | — | `$HOME/.agents/bin/resume-context.sh` | Prints ≤6 lines of Resume-Protocol context (handoff age, ADR index, changes/ index, trail-vs-HEAD warning, missing root files) only inside a harness-opted-in repo (`docs/changes/` exists); prints nothing otherwise. Always exits 0; stdout is injected into model context. |
| `UserPromptSubmit` (`:35-45`) | — | `~/.claude/feedback-capture.sh` | Pattern-matches user corrections ("no,", "that's wrong", "instead", etc.) and appends a JSON line to `~/.agents/feedback/log.jsonl`; explicitly skips any prompt starting with `<` (system-injected reminders/notifications) to avoid noise. Silent (must never print — stdout is injected into context); always exits 0. False positives are accepted by design (`/retro` filters with judgment); false negatives are the real risk the loose regex is tuned against. |
| `PostToolUse` (`:46-58`) | `Edit\|Write` | `~/.claude/format-hook.sh` | Runs `goimports`/`gofmt` on `*.go`, `prettier --no-install` on JS/TS/CSS/JSON/MD/YAML/HTML. No-ops silently if the tool isn't present. Always exits 0. |
| `SubagentStart` (`:59-70`) | `reviewer\|security\|verifier\|changelog\|explorer` | `subagent-check.sh start` | Stamps a per-`agent_id` marker file used by the paired `stop` check. Always exits 0. |
| `SubagentStop` (`:71-82`) | same 5 agents | `subagent-check.sh stop` | Blocks (`decision: block`) once if the named subagent stopped without writing its class-specific dated artifact (`reviews/`, `testing/`, `changes/`, `plans/`) newer than its start marker (or within the last hour if resumed with no marker). Uses `stop_hook_active` to guarantee it never blocks the same stop twice. Always exits 0; blocking is via JSON only. |
| `Stop` (`:83-93`) | — | `trail-check.sh` (no `--agy`) | Session-level DoD gate: if source changed today (dirty or committed-today) but `changes/`, `reviews/-review-`, or `testing/-verify-` artifacts dated today are missing, blocks once with a message naming the missing rows. Opt-in via `docs/changes/` existing; always exits 0. |

- Both `trail-check.sh` and `subagent-check.sh` are careful about false-positive avoidance: they
  scope "source path" narrowly (`trail-check.sh:56-79` — code/IaC/CI/deploy-manifest/dependency-
  manifest extensions and basenames only, explicitly excluding `*.md` and the docs tree itself),
  and both use `stop_hook_active`/marker files to guarantee at most one block per stop, so neither
  can loop forever. Both fail open (`return 0`/empty JSON) when `jq` is absent, when not in a git
  repo, or when the project hasn't opted in (no `docs/changes/` dir) — false negatives (missing
  a real gap) are more likely than false positives (spurious blocking) by explicit design choice
  documented in the scripts' own header comments.
- **destructive-guard.py false-positive/negative surface**: any compound command (`;`, `&&`,
  `|`, `` ` ``, `$`, redirection) containing a destructive verb falls to `ask` rather than being
  reasoned about (`destructive-guard.py:148-150`) — conservative, so false positives (extra
  prompts) are the failure mode, not false negatives. The script ships a 60+ case built-in
  `--selftest` table (`:228-344`) exercising symlink escapes, `git -C` worktree redirection,
  env-var wrapped commands, and glob-before-slash cases.
- **Scripts referenced by `settings.json` that exist**: all 6 do exist in `shared/scripts/` (or
  `claude/scripts/` for `statusline.sh`) and are deployed by `claude/setup.sh`. No dangling
  references found in `claude/settings.json`.
- **Scripts that exist but are never wired for Antigravity**: `destructive-guard.py` and
  `subagent-check.sh` — present in `shared/scripts/`, deployed to `~/.agents/bin/` only by
  `claude/setup.sh`, absent from `antigravity/setup.sh`'s gate-deploy loop
  (`antigravity/setup.sh:210` loops over exactly `resume-context.sh trail-check.sh`) and absent
  from `hooks.json`. This is a real asymmetry, not an oversight noted anywhere as intentional
  beyond the settings.json comment "Claude-only" — Antigravity has no equivalent guardrail for
  either concern.
- **Dangling reference**: `antigravity/setup.sh:243` logs `"...firing on agy is unverified, see
  SUB_AGENTS.md..."` — no `SUB_AGENTS.md` exists anywhere in the repo (confirmed via
  repo-wide `find`). Either a stale reference to a file that was renamed/removed, or a doc that
  was never created.
- **hooks.json (Antigravity) events**: `harness-format.PostToolUse` (format-hook.sh on
  `write_to_file|replace_file_content|edit_file|create_file`), `harness-gates.PreInvocation`
  (resume-context.sh --agy), `harness-gates.Stop` (trail-check.sh --agy). No equivalent to
  Claude's `SubagentStart`/`SubagentStop`/`PreToolUse`/`UserPromptSubmit` hooks exists in the agy
  schema per this file — consistent with agy lacking subagent lifecycle events (per
  `subagent-check.sh:17` comment: "Antigravity has no subagent events").

### 4. Skills inventory (`shared/skills/`, 28 skill directories)
All SKILL.md files are well under the 500-line flag threshold; largest is `terraform/SKILL.md`
at 393 lines, then `mermaid-syntax` (281), `init-project` (220), `handoff` (190). Full listing
(dir / SKILL.md lines / `name:` / description excerpt) captured in the survey; notable groups:

- **Workflow/slash-command skills** (`disable-model-invocation: true`, i.e. manual-only, so
  absence of "Use when" trigger phrasing in their description is by design, not an omission):
  `audit`, `distill`, `research`, `retro`, `init-project`, `handoff`, `grill-me`.
- **All other skills carry explicit trigger phrasing** in their `description:` ("Use when...",
  "Activate for...", "WHEN:", trigger words) — no skill was found with a bare, trigger-less
  description outside the manual-only group above.
- **Vendored / third-party, per `shared/skills/VENDORED.md`** (5 entries, each with source repo,
  pinned commit/version, license, vendored date): `grill-me`, `grilling`, `codebase-design`
  (all from `github.com/mattpocock/skills`, MIT), `find-skills` (`github.com/vercel-labs/skills`,
  MIT), `anti-ui-slop` (uizze.com, MIT, includes its own `LICENSE`/`NOTICE`/`MANIFEST.json`/
  `CHECKSUMS.sha256`/`MODIFICATIONS.md`).
- **Third-party-looking but UNDOCUMENTED in VENDORED.md**: `gcloud-auth-verification`
  (`license: Apache-2.0`, `metadata.publisher: google`), `gke-basics` (`license: Apache-2.0`,
  `metadata.author: Google Cloud`), `skill-repair` (`license: Apache-2.0`,
  `metadata.publisher: google`), `terraform` (`license: Apache-2.0`,
  `metadata.author: Anton Babenko`). These four carry the same provenance signals
  (license/author/publisher frontmatter) as the documented ones but are absent from the manifest
  that VENDORED.md's own header claims covers "skills vendored verbatim from the open skills.sh
  ecosystem."
- **anti-ui-slop** is the largest and most self-contained vendored skill (80K total): ships its
  own `CHECKSUMS.sha256`, Apache-2.0 `LICENSE` (for the bundled "designStack" v4.1.1, per
  `MANIFEST.json`), `NOTICE`, and `MODIFICATIONS.md` disclosing it's a UIZZE-maintained package
  with product-facing terminology.
- **No duplication of a Claude Code built-in was found.** `review-standards/SKILL.md`
  explicitly states it is "Not a review workflow — the reviewer agent and (on Claude Code) the
  built-in `/code-review` do the reviewing," and `preflight/SKILL.md` explicitly excludes
  "reviewing a diff (use the reviewer agent or /code-review)" from its own scope — the repo
  is deliberately avoiding this collision, consistent with the constitution's "the vendor
  built-in `/code-review` stays untouched."
- `golang/` (340K) and `gke-basics/` (188K) are the largest skill directories overall, but that
  size lives entirely in `references/` subdirectories loaded on demand — `SKILL.md` itself stays
  small (67 and 55 lines respectively), matching the intended "small entry point + deep
  references" skill pattern.

### 5. Agent definitions (`claude/agents/`, 6 files)
| Agent | model | tools | Matches DoD table role? |
|---|---|---|---|
| `changelog.md` | haiku | Read, Grep, Glob, Edit, Write | Yes — owns `changes/` + `DECISIONS.md` (append-only text edits); no Bash needed since it never runs commands, only writes docs. |
| `documentation.md` | sonnet | Read, Grep, Glob, Edit, Write | Yes — README/spec/ROADMAP updates; same no-Bash rationale. |
| `explorer.md` | sonnet | Read, Grep, Glob, Bash, Write | Yes — this very agent; Bash is explicitly scoped to read-only inspection (`git log`, `grep`, `ls`, `wc`) per its own prompt; no Edit granted, consistent with "never edits source." |
| `reviewer.md` | opus | Read, Grep, Glob, Bash, Write | Yes — Bash needed for `git diff`/`gh pr diff`; no Edit granted, matches "never modify source." Loads the `review-standards` skill (`skills: review-standards` frontmatter). |
| `security.md` | opus | Read, Grep, Glob, Bash, Write | Yes — Bash needed for `govulncheck`/`npm audit`; no Edit, matches "never edits source." |
| `verifier.md` | sonnet | Read, Grep, Glob, Bash, Write | Yes — Bash needed to run tests/builds; no Edit, matches "you do not fix." |
- Tool grants are consistently minimal-per-role: only the three read-only agents that must run
  shell commands (`explorer`, `reviewer`, `security`, `verifier`) get `Bash`; only the two that
  produce prose records (`changelog`, `documentation`) get `Edit` (append/update existing docs)
  instead of `Bash`. None of the six grants `Edit` alongside `Bash`, which would let a
  review/verify/recon agent both inspect and modify source — this lines up with the constitution's
  "never self-review" and "the only file it creates is its own report" language in each agent's
  own description.
- `reviewer.md` and `security.md` both use `opus`; `explorer.md`/`documentation.md`/
  `verifier.md` use `sonnet`; `changelog.md` uses `haiku` — a deliberate cost/capability tier by
  judgment-weight of the task.
- `changelog.md` additionally carries `memory: project` (per-project agent memory) — not present
  on the other 5 agents in the frontmatter snippets read.

### 6. Contract consistency
- `diff claude/CLAUDE.md antigravity/GEMINI.md` → **empty** (identical).
- `diff claude/CLAUDE.md shared/AGENTS.md` → **empty** (identical).
- So within the repo, the "deployed verbatim" claim holds for the current commit: all three
  files are byte-for-byte the same text.
- **But the actual deployed files on this machine are NOT verbatim with the repo**:
  `diff ~/.claude/CLAUDE.md claude/CLAUDE.md` shows real divergence — the deployed copy is
  missing the entire "Release Ownership" section, the DoD table's `/code-review`-sourced-review
  exception clause, the "Rows fire at commit boundaries..." sentence, the scratch-root/
  `destructive-guard` bullet under Environment & Stack, and the expanded `/preflight` description
  (ship/deploy/release authorization language) — all present in the current `claude/CLAUDE.md`.
  `~/.gemini/GEMINI.md` diffs identically against `antigravity/GEMINI.md` (same missing
  sections) and is byte-identical to `~/.claude/CLAUDE.md` (both stale, both from the same older
  deploy). This is corroborated by `~/.agents/bin/` on this machine containing only
  `resume-context.sh` and `trail-check.sh` (no `destructive-guard.py`, no `subagent-check.sh`),
  and by `~/.claude/settings.json`'s `hooks` object having only `SessionStart`, `PostToolUse`,
  `UserPromptSubmit`, and `Stop` — missing `PreToolUse` and `SubagentStart`/`SubagentStop`
  entirely, which the current `claude/settings.json` defines. All of this is one consistent
  story: this machine's harness was installed before those gates/sections were added and has not
  been redeployed since.

## Conventions
- Shell scripts: `#!/bin/bash`, no `set -e` in the gate/hook scripts (deliberate — a hook must
  never abort mid-way and must always reach its terminal `exit 0`), heavy header comments
  explaining wiring/exit-code contract/false-positive stance up front, small `g()`/`read_stdin()`/
  `mtime_of()` helper functions repeated near-verbatim across `trail-check.sh`,
  `resume-context.sh`, and `subagent-check.sh` (duplicated, not shared via a sourced lib — see
  Risks).
- The two `setup.sh` per-vendor scripts *do* use `set -euo pipefail` (fail-fast is correct for an
  installer, wrong for a hook), and share a `run()`/`backup()`/`log()` helper pattern duplicated
  between them rather than factored into a shared lib.
- Python scripts: stdlib-only, argv-based CLIs, exit via `sys.exit(msg)` on user error.
  `destructive-guard.py` ships a self-contained `--selftest` fixture/table instead of a separate
  test file (no `tests/` directory anywhere in the repo).
- Docs: everything under `.agents/docs/<class>/YYYY-MM-DD-<slug>.md` naming, append-only for
  `changes/` and `DECISIONS.md`, enforced mechanically by `trail-check.sh`'s `--check` mode
  (`:299-311`) via `git diff --name-status` on those paths.
- Frontmatter style: every skill/agent file begins with YAML frontmatter (`name`, `description`,
  optionally `tools`, `model`, `skills`, `license`, `metadata`) — consistent SKILL.md/agent.md
  convention across the whole repo.

## Risks
- **merge-settings.py hooks-merge bug** (`shared/scripts/merge-settings.py:20-33` list-merge
  branch, applied to `claude/settings.json:11-93`'s `hooks.*` arrays): hook entries are objects,
  not strings, so the "list of strings" fast path never applies; the fallback
  (`if item not in merged: merged.append(item)`) compares whole nested dicts by equality. Any
  edit to an existing hook's `timeout`, `matcher`, or command string in a future repo update
  will NOT replace the old entry on redeploy — it will sit alongside it as a second, distinct
  list item, so the hook fires twice (or the old, wrong-timeout version keeps firing forever
  since nothing removes it). The script's own docstring only warns about *removed* keys not
  being removed (`:12-13`) — it does not document this changed-value-appends-a-duplicate case.
- **Stale-install has no detection**: nothing in `trail-check.sh`, `resume-context.sh`, or
  either `setup.sh` compares the deployed `~/.claude/CLAUDE.md`/`settings.json` against the
  repo's current version and warns of drift. This machine is a live example of the failure mode
  the "Deployed verbatim" tagline promises to prevent going undetected indefinitely.
- **`.bak.*` files accumulate without bound**: `backup()` in both setup scripts always writes a
  new timestamped backup on every re-run, with no pruning or cap — repeated re-deploys (e.g. in
  a dotfiles-refresh cron, if one existed) will silently grow `~/.claude/`, `~/.gemini/`, and
  `~/.agents/bin/` indefinitely.
- **Antigravity guardrail asymmetry**: `destructive-guard.py` (destructive-Bash-command gate)
  and `subagent-check.sh` (per-subagent artifact enforcement) exist only for Claude Code. An
  agy session gets neither a destructive-command confirmation gate nor per-subagent DoD
  enforcement — only the session-level `trail-check.sh --agy` Stop hook, which is explicitly a
  "coarse, day-grained backstop" per the constitution's own DoD section, not a substitute for the
  finer-grained checks Claude Code gets.
- **Dangling doc reference**: `antigravity/setup.sh:243` references `SUB_AGENTS.md`, which does
  not exist in the repo — either dead documentation debt or a broken pointer for anyone trying to
  verify agy hook firing.
- **VENDORED.md under-reports provenance**: 4 skills (`gcloud-auth-verification`, `gke-basics`,
  `skill-repair`, `terraform`) carry third-party license/author/publisher frontmatter but are not
  listed in the manifest that exists specifically to track this. If any of these skills is ever
  edited in place (which `VENDORED.md`'s own instructions forbid for *listed* vendored skills —
  "never edit vendored content in place"), there is no manifest entry to catch the violation for
  these four, and no pinned source/commit to re-sync against.
- **No CI, no tests, no shellcheck wiring**: confirmed no `.github/workflows`, `.gitlab-ci.yml`,
  or any CI config anywhere in the repo; no `tests/` directory (the only "test"-named path is the
  templates skeleton's placeholder `docs/testing/.gitkeep`). `shellcheck -S warning` was run by
  hand against every shell script in the repo and returned zero findings — the scripts are
  currently clean, but nothing enforces that on future changes, and the one Python guard's only
  test coverage is its own `--selftest` fixture, never run automatically.
- **Repo does not dogfood its own scaffold**: there is no `AGENTS.md`, `REQUIREMENTS.md`,
  `ROADMAP.md`, `TODO.md`, or `.agents/docs/` at the repo root (prior to this recon brief, which
  created `.agents/docs/plans/` for the first time) — the harness that mandates this structure
  for every other project does not itself carry one. (`README.md` was deliberately removed per
  commit `a63a208`, per the dispatch note, which is a separate, apparently intentional decision.)
- **Helper-function duplication across gate scripts**: `g()`, `read_stdin()`, `mtime_of()`, the
  marker-dir-ownership pattern, and the docs-root resolution logic are each re-implemented nearly
  verbatim in `trail-check.sh`, `resume-context.sh`, and `subagent-check.sh` rather than sourced
  from one shared file. A bug fix or portability fix (e.g. the GNU/BSD `stat` fallback) made in
  one script is not guaranteed to be applied to the others; drift between the three is currently
  small but structurally easy to introduce.
- **Platform assumptions**: `format-hook.sh:12` explicitly handles both macOS
  (`${HOMEBREW_PREFIX:-/opt/homebrew}`) and Linux (`~/.go/bin`, `~/go/bin`, `/usr/local/go/bin`)
  Go toolchain locations — this one script is careful; `mtime_of()`'s GNU-first/BSD-fallback
  `stat` pattern in `trail-check.sh`/`subagent-check.sh` is similarly portable. No obviously
  macOS-only or Linux-only hardcoded path was found elsewhere in the scripts surveyed.

## Open questions
- Was `~/.claude`/`~/.gemini` on this machine deployed by an older commit of this same repo (via
  `git checkout` to an old ref + `setup.sh`), or hand-edited independently of the repo, or
  deployed by a completely different mechanism? The file contents match an earlier version of
  `claude/CLAUDE.md`'s prose closely enough to suggest an older repo commit, but no commit
  history exists to confirm (`git log` shows only 2 commits: "first commit" and "Remove
  readme" — the repo's own git history does not span the CLAUDE.md prose evolution visible in the
  diff, meaning the stale deploy predates this repo's current commit history, e.g. from a
  previous incarnation of the harness before it was reduced to 2 commits).
  **Recommend asking the user directly** whether `setup.sh` has been run since this repo's
  current state, since recon cannot distinguish "never redeployed" from "deployed then repo
  history was rewritten/squashed."
- Is the `merge-settings.py` hooks-duplication risk purely theoretical (repo hooks have never
  changed a `timeout`/matcher on an existing entry post-initial-deploy) or has it already
  happened on some user's machine? Not verifiable from the repo alone — would need to inspect a
  real `~/.claude/settings.json` history, which this recon only had one example of (and that one
  predates the current hooks entirely, so it doesn't test the merge path either).
- Does `agy plugin validate` / `agy -p "/skills"` (the only self-verification `antigravity/
  setup.sh` performs, `:245-280`) actually run in CI or only ad hoc on a machine with `agy`
  installed? No CI was found, so presumably ad hoc only — unconfirmed from the repo.
  `subagent-check.sh`'s own header states plainly "Antigravity has no subagent events" — is that
  a confirmed API fact or an assumption; the file offers no citation.
  `antigravity/setup.sh`'s own log line calls agy Stop-hook firing "unverified" — is there a
  newer confirmation elsewhere not captured by this survey (e.g. in the removed README, per
  commit `a63a208`)?
- Whether the four undocumented-provenance skills (`gcloud-auth-verification`, `gke-basics`,
  `skill-repair`, `terraform`) were vendored via the same `npx skills add` flow implied by
  `VENDORED.md`'s "Update via" column for the other five, or hand-authored/copied by a different
  route — cannot be determined from file contents alone; no pinned commit/source URL is present
  in any of the four files themselves.
