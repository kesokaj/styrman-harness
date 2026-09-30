---
owner: main
date: 2026-09-30
---
# Plan — R-parity, `sh-agent memory`, background `sh-agent dream`   (2026-09-30)

## Why
The operator, 2026-09-30:
- "The dream skills should be a features that runs autonomously in all projects registered on the machine and checks conversations, assets, feedback, critique, skills, for improvements. This should run in the background as a monitor or subagent that does not care about what is happening. It just collects, generates new skills if needed and such. That was the plan."
- "Dreamer is one and memory is the other. … remember these should work for both Claude or agy and both."
- "In the best of worlds I want the setup and executions to be exactly the same over agy and Claude. … Can we build in that feature in sh-agent?"
- On the proposed R-parity text, the dreamer table (1–7) and the memory table (1–5): "Sounds great. Include this in the next version. Go all approved go work."

The built memory loop (`specs/sh-agent/memory-loop.md`) is Claude-native for writing and loading and has no background job. Recon: `plans/2026-09-30-recon-memory-dream.md`.

## Scope
In:
1. **R-parity** in `REQUIREMENTS.md`, exactly as approved, plus its ADR. **R3's exception sentence** for the dream timer goes in when the operator approves its exact text (D3).
2. **`sh-agent memory add|show`**: both CLIs write and load project memory through the binary. Claude Code's native auto memory is switched off per project (`autoMemoryEnabled: false`), so both CLIs run the same path (R-parity).
3. **A `memory` skill** (when to write, the four types, how to call the verb). `feedback-log` rule 7 calls the verb on both CLIs. The context's resume step runs `sh-agent memory show` on both. `/init-project`, `/dream` Step 1, `agents/dreamer.md`, the memory spec and README are updated to match.
4. **`sh-agent dream [--all | <dir>] [--dry-run] [--cli claude|agy]`**: collects per project, redacts, makes one tool-free headless CLI call, and writes one proposal file per project. It keeps state, a log, and a `dream` line in `status`.
5. **`sh-agent dream --install-timer | --remove-timer`**: a systemd user timer on Linux or a launchd agent on macOS. It is opt-in, installed once per machine by the operator, and removed by `uninstall` (waits on D3).
6. **Consumers**: `/retro --harness` reads `proposals/dream/`; `sh-agent report` never sends `proposals/dream/`; `status` shows the dream line.

Out:
- Auto-installing anything the dreamer drafts (dreamer table item 4, approved: proposals only).
- A permission list written into either CLI's settings (REQUIREMENTS § Non-goals: "never a permission").
- Syncing memory between machines (memory table item 4: per machine).
- A `memory forget`/edit verb: deletes stay with the interactive `/dream` skill.

Follow-ups noted, not built:
- A cron fallback for Linux without systemd.
- Agent worktree isolation provided by `sh-agent` (the parity table example); it gets its own plan under R-parity.

## Surfaces touched
- Go: `cmd/sh-agent/main.go` (verbs, usage const :39, `status` :508); new `internal/memory/` (store read/write); new `internal/dream/` (collect, transcripts, run, state, timer); `internal/search` (exported project list, recon: store.go:41, search.go:225/284); `internal/harness/report.go` (exported redaction wrapper over `redactLines` :128 / `redactDocument` :257; `Collect` excludes `proposals/dream/`); `internal/harness/state.go` (`LastRetro` :137 unchanged — it reads only top-level `proposals/*.md`); `internal/harness/apply.go` (`Uninstall` :270 calls timer removal).
- Harness text: `skills/memory/SKILL.md` (new); `skills/feedback-log/SKILL.md:90`; `skills/init-project/SKILL.md:110-131`; `skills/dream/SKILL.md` Step 1; `agents/dreamer.md:18,28,32`; `skills/retro/SKILL.md:44,118` (reads `proposals/dream/`); `skills/distill:38`; `context/AGENTS.md:9` (resume) and `:42`.
- Docs: `REQUIREMENTS.md`, `DECISIONS.md`, `specs/sh-agent/memory-loop.md` (rewritten: memory and dream as two parts), new `specs/sh-agent/dream.md`, `README.md:214` plus the verbs table, `deploy/local.md` (timer install), `testing/skill-matrix.md`, `testing/sh-agent-matrix.md`.

## Interfaces

**`sh-agent memory add --type <user|feedback|project|reference> --name <slug> --description <text> [--title <text>]`**
- The body is read from stdin (required, non-empty). The command runs in the project: root = `git rev-parse --show-toplevel` from the working directory.
- It writes `<root>/.agents/memory/<slug>.md`: frontmatter `name`, `description` (JSON-quoted), `metadata:` `type`, `modified` (UTC RFC 3339), then the body. It sets that slug's line in `MEMORY.md` to `- [<title or slug>](<slug>.md) — <description>`; the same slug replaces both file and line (update, not duplicate).
- `<slug>` matches `^[a-z0-9][a-z0-9-]{0,63}$`.
- Refusals, each exit 2 with one line, and nothing written:
  - not a git checkout;
  - `.agents` or `.agents/memory` is a symlink;
  - git tracks anything under `.agents/memory`;
  - `.gitignore` lacks `.agents/memory` ("run /init-project first");
  - the index would exceed 200 lines.
- The store is created with mode 0700 when absent. A write goes to a temp file and is renamed; the index is rewritten whole under a lock file `<store>/.lock` (O_EXCL, stale after 30 s).

**`sh-agent memory show`**: same root and trust checks. Prints `MEMORY.md`, capped at 200 lines and 25 KB, with a trailing line naming the cut if one happened. No store, or an empty one, prints nothing and exits 0. A trust refusal exits 2.

**`sh-agent dream [--all | <dir>] [--dry-run] [--cli claude|agy]`**
- **Projects**: the `projects` table of `~/.config/sh-agent/search.db`, read with a plain SELECT, no reconcile (`search.ListProjects(home) ([]Project, error)`), or `<dir>` alone. A root is skipped (logged with the reason) when it is gone, not a git checkout, or lacks `.agents/docs`.
- **Watermark**: `~/.config/sh-agent/dream-state.json`, `{"<root>": {"last":"<RFC3339>"}}`. It advances only after the proposal file is written. The first run uses the last 14 days.
- **Collected since the watermark**, in this order:
  1. every memory file;
  2. `feedback.jsonl` entries whose `project` is the root;
  3. files under `<docs>/{bugs,reviews,retros,changes,handoffs}` with mtime after the watermark;
  4. the installed skills' names and descriptions (read from the vendor skill dir the manifest names);
  5. Claude Code transcripts: `~/.claude/projects/*/*.jsonl` whose `cwd` equals the root, with `user`/`assistant` text parts only;
  6. Antigravity conversations: ids from `~/.gemini/antigravity-cli/conversation_summaries.db` whose `workspace_uris` names the root, with text extracted per D8.
- Nothing changed since the watermark: skip, and log "unchanged".
- **Redaction**: `harness.Redact(string) string` (credential shapes, as `report`) plus `$HOME` → `~`. The prompt is capped at 256 KB; the oldest transcript text is dropped first and the drop is named in the prompt.
- **CLI**: `--cli`, else the manifest vendor (`antigravity` → agy; `claude` and `all` → claude).
  - claude: `claude -p --tools "" --no-session-persistence --output-format text`, prompt on stdin.
  - agy: the invocation the Phase 4 probe proves (D7).
  - Both run with cwd `~/.config/sh-agent/dream-work/` (created 0700, empty), so a run is never a project.
  - Timeout 10 min per project.
- **Output**: `~/.config/sh-agent/proposals/dream/<YYYY-MM-DD>-<project name>.md`. It gets frontmatter `owner: sh-agent dream`, `date`, `project` (root), `cli`, `inputs` (counts per kind), then the model's text, redacted again. The prompt asks for `## Proposals` (numbered: kind rule/skill/procedure/memory, target, text, evidence) and `## Draft skills` (full SKILL.md in fenced blocks). A same-day second file gets `-2`.
- **Log**: `~/.config/sh-agent/dream.log`, one line per project per run (`<ts> <root> ok|skip|fail <reason>`), plus `last-run` in the state file.
- **`--dry-run`**: prints per project what would be collected (counts and bytes) and the CLI command. It runs no CLI and writes nothing.
- **Lock**: `~/.config/sh-agent/dream.lock` (O_EXCL). A second run exits 0 with "already running".

**`sh-agent dream --install-timer | --remove-timer`**
- Linux: writes `~/.config/systemd/user/sh-agent-dream.{service,timer}`:
  - the service's `ExecStart` is the absolute path of the running binary plus `dream --all`;
  - the timer has `OnCalendar=*-*-* 03:00`, `Persistent=true`;
  - then runs `systemctl --user daemon-reload` and `systemctl --user enable --now sh-agent-dream.timer`.
- macOS: writes `~/Library/LaunchAgents/com.styrman.sh-agent-dream.plist` (`StartCalendarInterval` Hour 3), then `launchctl bootstrap gui/<uid> <plist>`.
- `--remove-timer` disables and unloads, then deletes the files; absent files are fine. Without `systemctl` (Linux) or `launchctl` (macOS): exit 2, naming what is missing.
- `uninstall` runs the remove path when a unit file exists.

**`status`**: a line `dream      last run <date> (<n> projects, <m> proposals), timer installed|not installed`, or `dream      never run`.

**`/init-project` Step 5b**: `.claude/settings.local.json` gains `"autoMemoryEnabled": false` beside `autoMemoryDirectory`, with the same trust checks and insertion method. On an Antigravity-only machine nothing changes (no `.claude/`).

**Context resume step**: "… → `sh-agent memory show` (the project's memory, on every CLI)". This replaces the Antigravity-only clause.

## Failure modes
- `memory add` outside a checkout, or with an untrusted store: exit 2 with one line; the session goes on without the memory.
- `memory show` on a refused store: exit 2 with one line naming why. The resume step reports it and continues.
- `dream`:
  - a project whose CLI call fails or times out: logged `fail <reason>`, watermark unchanged, the next project continues, final exit 1;
  - no CLI on PATH: exit 2 before any project;
  - search.db absent (no project ever indexed): exit 0, "no indexed projects: run `sh-agent index` in a project";
  - a transcript file that cannot be parsed: skipped and named in the proposal's `inputs`.
- Timer: a failed `systemctl`/`launchctl` call exits 2 and leaves the written unit files in place, named for the operator.

## Verification
Tiers: unit (Go, `-race`), sandbox e2e on three homes (`--vendor claude`, `--vendor agy`, `all`), and a real-home smoke at the release gate.

E2E criteria:
1. **Memory parity**: in a scaffolded project, on each of the three homes, a session told "remember that we deploy on Fridays" leaves one memory file and one index line in `.agents/memory/`. A fresh session on the same home, and on the other CLI where both are installed, answers "when do we deploy?" from memory.
2. **Claude native writes are off**: after a Claude Code session in a scaffolded project, nothing new exists under `~/.claude/projects/<slug>/memory/`.
3. **Dream golden path**: `sh-agent dream --all` over two indexed projects with new trail and memory writes one proposal file per project under `proposals/dream/`, with `## Proposals`. A second run with no changes writes nothing and logs `unchanged`.
4. **Dream edge**: a project whose root was deleted is skipped with a log line; a secret-shaped string placed in a memory file does not appear in the proposal.
5. **Timer**: on Linux, `--install-timer` makes `systemctl --user list-timers` show `sh-agent-dream.timer`; `--remove-timer` removes it; `uninstall` removes it when present.
6. **Parity**: criteria 1 and 3 hold with `--cli claude` and with `--cli agy`.

## Phases
- [ ] Cold-reader check — design ready and implementation ready (<reviews/ paths>)

### Phase 0 — Requirements
- [ ] R-parity added to `REQUIREMENTS.md` verbatim as approved; ADR-047; `ROADMAP.md` "Now" names this plan — verify: AI, `grep -c "R-parity" REQUIREMENTS.md` — expect: 1 — docs: REQUIREMENTS.md
- [ ] R3 exception sentence added verbatim once the operator approves D3's text; ADR-047 extended — verify: AI, grep — expect: the approved sentence present — docs: REQUIREMENTS.md
- [ ] Checkpoint — expect: requirements committed — evidence: commit, date

### Phase 1 — `sh-agent memory`
- [ ] `internal/memory`: `Add(root, Entry, body)` and `Show(root) (string, error)` with the trust checks, lock and atomic write — verify: unit, `go test -race ./internal/memory` — expect: tests for each refusal, update-in-place, index cap, `show` cap and empty store pass; each fails without its check (mutation of one refusal)
- [ ] `memory add|show` wired in main.go and the usage text — verify: unit + a scratch checkout run of `add` then `show` — expect: the file, the index line, and the printed index — docs: README.md verbs table
- [ ] Checkpoint — expect: full `go test -race ./...` green, verifier log — evidence: <log>, date

### Phase 2 — Memory on both CLIs (text)
- [ ] `skills/memory/SKILL.md`: when to write (the four types, Claude Code's own triggers restated), `sh-agent memory add` usage, never a secret, one fact per file — verify: `sh-agent install --home <scratch>` lint passes; skill listed on both vendors — expect: installed at both paths — docs: README.md skills list
- [ ] `feedback-log` rule 7 calls `sh-agent memory add` on every CLI; `context/AGENTS.md` resume → `sh-agent memory show`, bytes paid by cuts in the same edit — verify: `TestRenderedContextFitsAntigravityBlock` — expect: pass, byte count ≤ 22000 read from `-v`
- [ ] `/init-project` Step 5b writes `autoMemoryEnabled: false`; `/dream` Step 1: the store is `<root>/.agents/memory` on every machine; `agents/dreamer.md` and the memory spec match — verify: `/review --docs` grep for `autoMemoryDirectory` resolution in living text — expect: no copy saying the old thing — docs: specs/sh-agent/memory-loop.md
- [ ] Checkpoint — expect: E2E criteria 1 and 2 on the three sandbox homes (sandbox rules: `AGENTS.md` § Testing, BUG-115) — evidence: <log>, date

### Phase 3 — `sh-agent dream` collection
- [ ] `search.ListProjects(home)` — verify: unit — expect: lists the rows of a fixture db; missing db → empty, no error
- [ ] `harness.Redact` exported wrapper plus the `$HOME` rewrite — verify: unit — expect: the report tests' needles redact through it
- [ ] Claude transcript reader (cwd match, text parts only) — verify: unit on a fixture jsonl — expect: only the matching session's user/assistant text
- [ ] Antigravity reader per D8 — first task: the probe of D8 — verify: unit on a fixture made from a real conversation db copied to scratch, fields reduced — expect: the text of a known prompt is extracted; an unmapped conversation is ignored
- [ ] Collect, watermark, 256 KB cap, `--dry-run` — verify: unit + `sh-agent dream --dry-run --all --home <scratch>` — expect: per-project counts; nothing written (asserted) — docs: specs/sh-agent/dream.md
- [ ] Checkpoint — evidence: <log>, date

### Phase 4 — the headless run and its outputs
- [ ] Probe (D7): `agy -p` in an empty cwd with the dream prompt and the flags under test returns text without a tool call and without "no output produced"; record the flags — verify: sandbox run on an agy-only home — expect: text on stdout, exit 0
- [ ] Run, write proposal, state, log, lock, `status` line — verify: unit with a fake CLI binary on PATH (a Go test helper) — expect: the proposal file shape; watermark advances only after the write; a failing CLI leaves it — docs: README.md, specs/sh-agent/dream.md
- [ ] `report` never sends `proposals/dream/`; `/retro --harness` reads it — verify: unit on `Collect` — expect: the path absent from a report bundle — docs: skills/retro/SKILL.md
- [ ] Checkpoint — expect: E2E criteria 3, 4 and 6 on the sandbox homes — evidence: <log>, date

### Phase 5 — the timer (waits on D3)
- [ ] `--install-timer`/`--remove-timer`, `uninstall` removal — verify: unit on the unit-file text; Linux run in a scratch home (`XDG_CONFIG_HOME`/`--home`) with `systemctl --user` on this box — expect: E2E criterion 5 — docs: deploy/local.md, README.md
- [ ] Checkpoint — evidence: <log>, date

### Phase 6 — close
- [ ] Docs pass (`documentation`), `/review --docs`, skill-matrix and sh-agent-matrix rows — verify: reviewer docs lens — expect: no stale text
- [ ] Release through `/preflight` (the operator: "Include this in the next version")

## Risks
Assumptions:
- `autoMemoryEnabled: false` in `.claude/settings.local.json` stops Claude Code's native writes — checked by E2E criterion 2, the first task of Phase 2's checkpoint.
- `claude -p --tools ""` accepts the prompt on stdin and returns plain text — checked by the Phase 4 unit fake and one real run in Phase 4.
- agy conversation text is extractable (D8) — the first task of the Antigravity reader.
- agy can answer tool-free (D7) — the Phase 4 probe.

## Decisions
- D1 — 2026-09-30 — R-parity text as approved — by: "Sounds great. … Go all approved go work."
- D2 — 2026-09-30 — Memory is written and loaded through `sh-agent` on both CLIs; Claude Code's native auto memory is switched off per project, because R-parity says "a native feature is used only where it gives the same result on both" — by: memory table item 1 ("written by the harness on both CLIs"), approved "Go all approved"
- D3 — 2026-09-30 — R3 exception, exact text shown to the operator 2026-09-30: "One exception: `sh-agent dream --all`, run by an OS user timer that `sh-agent dream --install-timer` installs and `--remove-timer` removes, reads every indexed project's trail, memory, the feedback log and both CLIs' transcripts, and writes proposals and draft skills to `~/.config/sh-agent/proposals/` only; it never edits a project, a skill, a rule or the context, and a proposal takes effect only through `/retro --harness`." — by: awaiting the operator's yes on this exact text
- D4 — 2026-09-30 — The dream job's CLI call has no tools; `sh-agent` does all reading and writing. This closes the injection path and needs no permission in any vendor setting (REQUIREMENTS § Non-goals) — by: implementer, from recon risks 1 and 3
- D5 — 2026-09-30 — The timer is opt-in (`--install-timer`, once per machine), never part of `install` — by: implementer; dreamer table item 1 names the command
- D6 — 2026-09-30 — Proposals go to `proposals/dream/`, a subdirectory, so `LastRetro` (top-level `*.md` only, state.go:137) does not count a dream as a retro, and `report` excludes the subdirectory (transcript-derived content never leaves the machine) — by: implementer, from recon open question 7
- D7 — 2026-09-30 — agy invocation for a tool-free answer is decided by the Phase 4 probe; candidates `--mode plan --disable-slash-commands --output-format text` in an empty cwd. If no invocation answers without a tool call, dream uses claude for text and says so in `status`; on an Antigravity-only machine the job then stops with that reason, an R-parity known difference to be put to the operator — by: implementer
- D8 — 2026-09-30 — Antigravity transcripts: map by `conversation_summaries.db` `workspace_uris`, then extract UTF-8 text runs of 20+ characters from the `steps` blobs (protobuf, no schema), deduplicated in order. If the probe shows this cannot isolate the user and model text, transcripts are left out on BOTH CLIs (R-parity) and the job reads memory, feedback and trail only; this is put to the operator — by: implementer, from recon
- D9 — 2026-09-30 — The watermark lives in `~/.config/sh-agent/dream-state.json`, not in `search.db`, which is a cache (R4) — by: implementer, from recon
- D10 — 2026-09-30 — The `/dream` skill stays: it is the interactive consolidation with deletes and a gate. The `sh-agent dream` verb is the background proposal job. The names are kept because the operator calls both "dream"; the skill's description names the verb — by: implementer
- D11 — 2026-09-30 — Existing Claude default stores keep migrating through `/init-project` Step 5b as today; no new migration — by: implementer
- Open: redaction beyond credential shapes and the home path (emails, names) — default: none beyond `report`'s deny-list and `~`, since nothing leaves the machine except to the CLI the operator already uses for that project
- Rejected: auto-applying memory merges from the background job — the operator approved proposals only (dreamer table item 4)
- Rejected: `memory list`/`forget` verbs — `show` and the `/dream` skill cover them
