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
- On the R-parity text and the two tables below: "Sounds great. Include this in the next version. Go all approved go work."

Both tables are quoted as the operator approved them in chat on 2026-09-30.

**Dreamer table** (approved):
1. How it runs: a user timer (systemd on Linux, launchd on Mac) runs `sh-agent dream --all` nightly, headless, regardless of what any session is doing.
2. Which projects: every harness project in the global search index (`~/.config/sh-agent/search.db`); no new registry; projects that haven't changed are skipped.
3. What it reads: both CLIs' conversation transcripts, `feedback.jsonl`, each project's memory, bug records, reviews, retros, and the installed skills — read-only, treated as data.
4. What it writes: draft skills and proposals into `~/.config/sh-agent/proposals/`, nothing auto-installed; applied through `/retro --harness`.
5. Which CLI: the one the install manifest names (`agy -p` on an agy machine, `claude -p` where Claude Code is installed).
6. Cost limit: once a night, a cap per project, a log line per run shown by `sh-agent status`.
7. The requirement change: R3 gets one named exception, with an ADR.

**Memory table** (approved):
1. agy's memory is thinner (it writes only on corrections) → both CLIs write the same four types by the same rule, written by the harness.
2. Claude Code loads memory automatically; agy only because the context says so → accept for now.
3. Claude Code's split store in an untrusted folder → the dreamer detects it and moves the notes.
4. Notes don't follow the operator between machines → keep per machine.
5. The memory spec header says "nothing below is built" → fix.

Since then, R10's "the harness provides it once and both CLIs use that version" supersedes memory items 2 and 3 (D2): loading goes through `sh-agent memory show` on both CLIs, and with native writes off there is no split store. Recon: `plans/2026-09-30-recon-memory-dream.md`. Cold read r1: `reviews/2026-09-30-cold-read-parity-memory-dream-r1.md`.

## Scope
In:
1. **R10 (R-parity)**: done in d6e396a. The requirement edits the dream job needs, on the operator's word (D3, D12, D13): the R3 exception, R4's driver and `uninstall` sentences, and the Non-goals memory-keys exception.
2. **`sh-agent memory add|show`** on both CLIs. Claude Code's native auto memory is switched off per project (`autoMemoryEnabled: false`, D2).
3. **The `project-memory` skill** (D14), plus:
   - `feedback-log` rule 7 calls the verb on every CLI;
   - the context resume step runs `sh-agent memory show`;
   - `/init-project` Step 5b writes the new key;
   - `/dream` Step 1 always uses `<root>/.agents/memory`;
   - `agents/dreamer.md`, the memory spec and README are updated.
4. **`sh-agent dream --all | --project <dir>`, with `--dry-run`, the CLI from the global `--vendor` flag**: collect, redact, one tool-free CLI call, then render and write proposals and drafts. Also a state file, a log, a lock, and the `status` line.
5. **`sh-agent dream --install-timer | --remove-timer`**, opt-in (D5). `uninstall` removes the timer (D12).
6. **Consumers**:
   - project-level `/retro` reads the dream proposals for its own project;
   - `/retro --harness` reads the harness-target ones;
   - `report` sends `proposals/dream/harness/` and nothing else under `proposals/dream/` (D19);
   - `status` gets the dream line;
   - the README map gets the job (D13).

Out:
- Auto-applying anything (dreamer item 4).
- Any permission in a vendor settings file (Non-goals).
- Cross-machine memory (memory item 4).
- `memory forget`/`list` verbs (the `/dream` skill deletes interactively).

Follow-ups noted, not built:
- A cron fallback for Linux without systemd.
- Agent worktree isolation in `sh-agent` (a later R10 plan).

## Surfaces touched
- **Go**:
  - `cmd/sh-agent/main.go`: verbs; usage const :39; `status` :508; `uninstall` calls `dream.RemoveTimer`, so `internal/harness` does not import `internal/dream`.
  - new `internal/memory/`;
  - new `internal/dream/` (collect, claude transcripts, prompt, run, render, state, lock, timer);
  - `internal/search`: `ListProjects(home)`, plus `AgyConversationText(home, root, since)`, which keeps the driver inside `internal/search` (AGENTS.md:59, D12);
  - `internal/harness/report.go`: exported `Redact(string) string` over `redactLines` (:199); `Collect` (:331-339) gains `proposals/dream/harness/*.md` and still skips the rest of `proposals/dream/` (D19).
- **Harness text**:
  - `skills/project-memory/SKILL.md` (new);
  - `skills/feedback-log/SKILL.md:90`;
  - `skills/init-project/SKILL.md:110-131`;
  - `skills/dream/SKILL.md` Step 1 and description;
  - `agents/dreamer.md:18,28,32`;
  - `skills/retro/SKILL.md:44,48,118`;
  - `skills/distill/SKILL.md:38`;
  - `context/AGENTS.md:9` (resume).
- **Docs**:
  - `REQUIREMENTS.md` (R3, R4, Non-goals, on approval);
  - `DECISIONS.md` (ADR-047 addendum);
  - `specs/sh-agent/memory-loop.md` (rewritten: the memory half);
  - new `specs/sh-agent/dream.md`;
  - `README.md`: map :5-60 with `TestHarnessMap`, `configDir` contents :30, state rows :212-214, verbs table, `uninstall` row;
  - project `AGENTS.md` Layout (`internal/memory`, `internal/dream`) and :59 (the driver's second reader in `internal/search`);
  - `deploy/local.md` (timer install/remove, uninstall);
  - `testing/skill-matrix.md` and `testing/sh-agent-matrix.md` rows.

## Interfaces

**`sh-agent memory add --type <user|feedback|project|reference> --name <slug> --description <text>`**, with the body on stdin (required, non-empty).
- **Root**: `git rev-parse --show-toplevel` from the working directory.
- **The file**: `<root>/.agents/memory/<slug>.md`, whose frontmatter is:
  ```
  ---
  name: <slug>
  description: "<text, JSON-escaped>"
  metadata:
    node_type: memory
    type: <type>
    modified: <UTC RFC 3339>
  ---
  ```
  then the body. `originSessionId` is not written: no CLI exposes its session id to a shell call, and `/dream`'s promote test already counts a file without it as one session.
- **The index**: the line in `MEMORY.md` for that file is `- [<slug>](<slug>.md) — <description>`. The same slug replaces both the file and its line; an existing line is found by its link target `(<slug>.md)`. A `--description` with a newline is refused (exit 2). Files are written 0600.
- **`<slug>`** matches `^[a-z0-9][a-z0-9-]{0,63}$`.
- **Refusals**: exit 2, one line on stderr, nothing written. In order:
  1. not a git checkout;
  2. `.agents` or `.agents/memory` is a symlink;
  3. `git ls-files .agents/memory` prints a path;
  4. `git check-ignore -q .agents/memory/<slug>.md` exits non-zero (the path would be committed: "run /init-project first");
  5. the index would pass 200 lines.
- **Store and writes**: created 0700 when absent. The file is written to a temp name and renamed.
- **Lock**: the index is rewritten whole under `<store>/.lock`, created O_EXCL and holding the PID. A lock whose PID is not alive, or older than 30 s, is removed. A live lock is retried for 5 s, then exit 2 "store busy".

**`sh-agent memory show`**
- Same root and checks 1–3.
- Prints `MEMORY.md`, capped at 200 lines and 25 KB; a cut adds a last line `(… N more lines in .agents/memory/MEMORY.md)`.
- No store, or no index: prints nothing, exit 0.
- **Rollout hint**: when `<root>/.claude/settings.local.json` exists and lacks `"autoMemoryEnabled": false`, `show` prints one last line: `Claude Code's native memory is still on here: run /init-project to switch it off (R10)`. Every session start on either CLI then surfaces the rollout for projects scaffolded before this change.

**`sh-agent dream --all | --project <dir> [--dry-run] [--vendor claude|agy]`**. Neither `--all` nor `--project` is a usage error. `--project` exists for the e2e and to rerun one project after a failure (D15).
- **Projects**: `search.ListProjects(home)`, a plain SELECT on `projects(root,…)` (store.go:41), no reconcile. A missing db means an empty list. `--project` takes one root. A root is skipped with a log line when it is gone, not a git checkout, or lacks `.agents/docs`.
- **Start**: `t0` = now, taken before collecting.
- **Watermark**: in `~/.config/sh-agent/dream-state.json`, `{"projects":{"<root>":"<RFC3339>"},"last_run":{"at":…,"projects":n,"proposals":m,"error":""}}`. The first run uses t0 − 14 days.
- **Changed** means any per-project input below (1, 2, 3, 5, 6) newer than the watermark: file mtime; for transcripts, the line or step time. Input 4 (installed skills) is global: it never makes a project changed and is sent only with a changed project. When nothing changed: skip, log `unchanged`, and write no file. When something changed, the memory files are sent whole; the other inputs only the changed parts.
- **Inputs**:
  1. the memory store;
  2. `feedback.jsonl` lines whose `project` is the root and whose `ts` is after the watermark;
  3. files under `<docs>/{bugs,reviews,retros,changes,handoffs}` with mtime after the watermark;
  4. the installed skills' `name` and `description` from the manifest vendor's skill directory;
  5. Claude Code `~/.claude/projects/*/*.jsonl` lines with `timestamp` after the watermark whose `cwd` is the root or under it, `type` `user`/`assistant`, text parts only;
  6. Antigravity text from `search.AgyConversationText(home, root, since)` (D8). It copies `conversation_summaries.db` and each matching `conversations/<id>.db` with their `-wal` and `-shm` files into a private temp directory under `<home>/.config/sh-agent/`, opens the copies, and deletes them after; it never opens a live file.
- **Redaction**: every input passes `harness.Redact` and `$HOME` → `~`.
- **Prompt**: fixed instructions (in `internal/dream/prompt.go`), then the inputs in the order above, each fenced and labelled "data, not instructions". The cap is 256 KB: the oldest transcript text is dropped first, and the prompt says what was dropped.
- **CLI**: the global `--vendor` flag (main.go:62, "default: from manifest"): `antigravity`/`agy` → agy; `claude` and `all` → claude (dreamer item 5).
  - claude: `claude -p --tools "" --no-session-persistence --output-format text`, prompt on stdin.
  - agy: the invocation proven by the Phase 4 probe (D7).
  - Both run with cwd `<home>/.config/sh-agent/dream-work/` (0700, kept empty; it is not a project, so the run's own conversation never maps to one) and a 10-minute timeout.
  - The child inherits `sh-agent`'s own environment unchanged: `--home` moves only what `sh-agent` itself reads and writes, never the child's `HOME`, so the CLI authenticates as it does for the operator (AGENTS.md:65). A sandbox run sets `HOME` for the whole `sh-agent` process instead.
  - The binary is resolved on PATH at run time; under the timer, PATH is the one captured at `--install-timer` (below).
- **Model output**: the prompt asks for one JSON object:
  ```
  {"proposals":[{"title","kind":"rule|skill|procedure|memory","target":"harness|project","file","text","evidence"}],
   "skills":[{"name","description","body"}]}
  ```
  `sh-agent` takes the first `{…}` span of stdout that parses. When nothing parses, the project is `fail invalid output`, and the raw text is kept in `<home>/.config/sh-agent/dream-failed/<date>-<project>.txt` (0600), outside the CLI's cwd.
- **Written on success**:
  - `~/.config/sh-agent/proposals/dream/<YYYY-MM-DD>-<project name>.md` (`-2`, `-3` on the same day). Frontmatter: `owner: sh-agent dream`, `date`, `project`, `cli`, `inputs` (counts per kind, and what was dropped). Then one section per proposal:
    ```
    ### <n>. <title>
    **Status:** queued
    **Target:** harness|project   **Kind:** …   **File:** …

    <text>

    Evidence: <evidence>
    ```
  - Harness-target proposals (`target: harness`) are written to a second file, `proposals/dream/harness/<YYYY-MM-DD>-<project name>.md`, same shape. `report` sends that directory (R8: report sends "this machine's corrections, proposals"), so a binary-only or zip place feeds harness proposals back too; project-target files stay on the machine (D6).
  - Every skill goes to `proposals/dream/drafts/<date>-<project name>-<name>/SKILL.md`, with `name` validated as a slug (an invalid name is dropped and listed). The draft is `---\nname: <name>\ndescription: "<…>"\n---\n<body>`.
  - All model text is redacted again before it is written.
- **Zero proposals**: when the JSON parses with no proposals and no skills, no file is written, the log says `ok 0 proposals`, and the watermark advances.
- **Watermark**: advances to t0 only after a successful parse and, when there is anything to write, after the write.
- **Log**: `~/.config/sh-agent/dream.log`, one line per project (`<ts> <root> ok|unchanged|skip|fail <reason>`) and one per run (`<ts> - run <n> projects, <m> proposals[, error <reason>]`). A failure before the first project (no CLI, lock or state unreadable) writes the run line with `error` and sets `last_run.error`.
- **Lock**: `~/.config/sh-agent/dream.lock` holds the PID. A lock whose PID is not alive is stale and replaced; a live lock makes the run exit 0 with "already running".
- **`--dry-run`**:
  - prints per project the input counts and bytes, what would be dropped, and the CLI command line;
  - prints the prompt when `--project` is given;
  - runs no CLI, and writes, locks and advances nothing (a test asserts each).
- **Exit**: 0 when every project is ok, unchanged or skipped; 1 when any failed; 2 on a pre-project failure.

**`sh-agent dream --install-timer | --remove-timer`** (D5, D12)
- **Linux**: writes `~/.config/systemd/user/sh-agent-dream.service` and `.timer`.
  - The service: `ExecStart=<home>/.local/bin/sh-agent dream --all` (BinPath), and `Environment=PATH=<the PATH of the installing shell>`, so `claude` and `agy` in `~/.local/bin` resolve.
  - The timer: `OnCalendar=*-*-* 03:00`, `Persistent=true`.
  - Then `systemctl --user daemon-reload` and `systemctl --user enable --now sh-agent-dream.timer`.
- **macOS**: `~/Library/LaunchAgents/com.styrman.sh-agent-dream.plist` with the same program, `EnvironmentVariables` PATH and `StartCalendarInterval` Hour 3, then `launchctl bootstrap gui/<uid> <plist>`.
- **Guards**:
  - `--install-timer` refuses when BinPath is missing, or when `systemctl`/`launchctl` is absent (exit 2, naming it).
  - Unit files are written only when absent or identical to the rendered text; a hand-edited one is kept and named.
- **`--remove-timer`**: `systemctl --user disable --now` / `launchctl bootout`, then deletes the files. Absent is fine.
- **`--home`**: the unit files are written under `<home>`, like every other path (main.go:61); `systemctl` and `launchctl` are resolved on PATH, so tests put fakes first on PATH and never reach the real user manager (R6). Each branch (write, identical, hand-edited kept, absent on remove, a failing `systemctl`, `uninstall` without a manifest) is tested against the fakes.
- **`status`** reports the timer as installed when its unit file exists.
- **The manifest**: the timer is not listed in it (listed files are deleted as stale by the next `install`, apply.go:88-109).
- **`uninstall`**:
  - calls the same removal when a unit file exists, including when there is no manifest (D12);
  - `uninstall --dry-run` prints "would remove the dream timer" and runs no `systemctl`/`launchctl` (a test asserts it).

**`status`**: after the `retro` line (main.go:508):
- `dream      last run <date>: <n> projects, <m> proposals; timer installed|not installed`
- `dream      last run <date> failed: <error>; timer …`
- `dream      never run; timer …`

**`/init-project` Step 5b item 2**: the settings file gets `"autoMemoryDirectory"` as today, and `"autoMemoryEnabled": false`. Each key is handled on its own with the existing branches:
- absent file → both keys in the one-object write;
- present without the key → one inserted line after `{` per key, with the Edit tool;
- present with the key → kept and reported (a `true` is reported, never changed).

The parse check runs after each insert. Antigravity-only machine: nothing (Step 1).

**Context resume step**: "… → `sh-agent memory show` (the project's memory; on every CLI)" replaces "on Antigravity, `.agents/memory/MEMORY.md` when … Claude Code loads it itself". The byte difference is paid in the same edit (`TestRenderedContextFitsAntigravityBlock`, 21998/22000 today).

## Failure modes
- `memory add|show` refusal: exit 2 with a one-line reason; the session reports it and continues without memory.
- Headless sessions without a grant for `sh-agent`: the call is refused and no memory loads. This is the same on both CLIs, a cost of R10 (D16).
- `dream`: the per-project failures are those above (fail + log, watermark kept, others continue, exit 1). A pre-project failure is exit 2 plus the run log line and `status` "failed". Unparseable model output keeps the raw text.
- Timer: a failing `systemctl`/`launchctl` exits 2 and names the unit files left in place.

## Verification
Tiers: unit (`go test -race`), sandbox e2e, the `dreamer` behavioural run, § 3d for the search additions, and the release gate's real-home runs. The setups for R10:
- **agy-only**: a sandbox home installed `--vendor agy` (the kit of `testing/2026-09-30-verify-real-home-v1.16.2.md`).
- **both**: the real home, vendor `all`.
- **claude-only**: a sandbox home installed `--vendor claude`, with Claude Code pointed at it by `CLAUDE_CONFIG_DIR` and a copied credentials file (chmod 600, deleted by name after). This is an assumption the Phase 2 checkpoint's first run checks; when Claude Code cannot authenticate that way, the claude-only runs are marked not exercised and the "both" runs stand for Claude Code.

E2E criteria:
1. **Memory parity**: in a scaffolded fixture project, on each setup, a session told "remember that we deploy on Fridays" leaves a file and an index line in `.agents/memory/`. A fresh session on the same setup, on each installed CLI, answers "when do we deploy?" with Friday.
2. **Claude native writes are off**: in the claude setup, with the pre-change harness installed (no `project-memory` skill, rule 7 as before), the must-fail control comes first. With `autoMemoryEnabled` absent, a session told a preference leaves a file carrying `originSessionId` (native) in `.agents/memory/`. Then the key is set to `false` and the current harness is installed, and no file with `originSessionId` appears while the verb-written file does. A control that produces no native file makes E2E 2 "not exercised", never passed.
3. **Dream golden path** (home: a scratch `--home` whose `search.db` indexes only the fixture projects; the child is the real `claude`): `sh-agent dream --all --home <scratch>` over two indexed fixture projects with new trail and memory writes one proposal file per project with at least one `**Status:** queued` section. A second run with no changes writes no file and logs `unchanged` twice.
4. **Dream edge**:
   - a deleted root is logged `skip`;
   - a secret-shaped string in a fixture memory file is absent from the prompt, shown by `--dry-run --project` and by the fake-CLI unit test's captured stdin.
5. **Timer** (`who: user`, R6):
   - a unit test pins the rendered unit text;
   - `systemd-analyze verify` passes on the rendered files in scratch;
   - the operator runs `sh-agent dream --install-timer` on this box, and `systemctl --user list-timers` shows it; `--remove-timer` removes it.
6. **Parity**: criterion 1 on the claude and agy setups. Criterion 3 with `--vendor claude` (scratch `--home`, as above) and on the agy-only sandbox (the whole process under the sandbox `HOME`, with its settings file leaving `dream-work/` untrusted, as the real home does).

## Phases
- [x] Cold-reader check — three rounds (`reviews/2026-09-30-cold-read-parity-memory-dream-r1.md`, `-r2.md`, `-r3.md`), none both-ready; at the cap the operator took the residue as decided defaults: "Go all" (2026-09-30), on the six answers now D21–D26 and the non-blocking defaults D27

### Phase 0 — Requirements
- [x] R10 added verbatim; ADR-047; `ROADMAP.md` "Now" names this plan — verify: AI, `grep -c "R-parity" REQUIREMENTS.md` — expect: 1 — evidence: d6e396a — docs: REQUIREMENTS.md
- [x] R3 exception, R4 driver and uninstall clauses, and the Non-goals memory-keys exception, verbatim as approved (D3, D12, D13), plus ADR-047a — verify: AI, grep each approved sentence — expect: each present once — evidence: `grep -c` 1 each, 2026-09-30 — docs: REQUIREMENTS.md
- [x] Checkpoint — expect: requirements committed and pushed — evidence: the Phase 0 commit (this one), 2026-09-30

### Phase 1 — `sh-agent memory`
- [x] `internal/memory` `Add` and `Show` — evidence: `go test -race ./internal/memory` ok; mutation of the tracked-store check (`if false &&`) compiled and failed `TestAddRefusals/git_tracks_the_store` and `TestShowRefusesTrackedStore`, restored byte-identical, 2026-09-30 — verify: unit, `go test -race ./internal/memory` — expect: a test per refusal (1–5), update in place, the index cap, the stale lock, the `show` cap and the empty store; one refusal's check mutated makes its named test fail (review-standards § 2)
- [x] `memory add|show` in main.go and usage — evidence: `TestMemoryVerb`, `TestMemoryRefusalExitsTwo` (red first on `takesValue`, BUG-002's shape); a fixture at `~/sandbox/p1fixture` wrote the file and line 0600 in a 0700 store, `show` printed it with the rollout hint, `/tmp` exited 2, 2026-09-30 — verify: unit (main_test) plus `add` then `show` in a scratch fixture checkout outside this repo — expect: the file, the line, and the printed index — docs: README.md verbs table
- [x] Checkpoint — expect: bare `go test -race ./...` exit 0; verifier log — evidence: `testing/2026-09-30-verify-parity-memory-dream-phase-1.md` PASS (404 pass, 0 fail; the verifier's own mutation of the check-ignore refusal failed `TestAddRefusals/store_not_ignored`); two tests added after it for branches it named (`TestShowOutsideACheckoutRefuses`, the unparseable-settings hint case), package re-run ok, 2026-09-30

### Phase 2 — Memory on both CLIs
- [x] `skills/project-memory/SKILL.md` — evidence: scratch install lints, installed under both `.claude/skills/` and `.gemini/config/skills/`, 29 skills, drift none, 2026-09-30; built-in check: Claude Code has `/memory`, not `/project-memory`; Antigravity's expansion of `/project-memory` is shown in the Phase 2 e2e — when to write (the four types and their triggers), the verb, one fact per file, never a secret; the name checked against both CLIs' built-in commands (BUG-065) — verify: scratch install lints; listed on both vendors — expect: installed at both paths — docs: README.md skills list
- [x] `feedback-log` rule 7 → the verb on every CLI — evidence: rendered context 21911/22000 bytes, 2026-09-30 — the context resume → `sh-agent memory show`, bytes paid in the same edit — verify: `go test -run TestRenderedContext -v ./internal/harness` — expect: pass, count ≤ 22000
- [x] `/init-project` Step 5b second key — evidence: `git grep` for the old resolution and "auto memory" claims leaves only the layout line (the shape is Claude Code's), 2026-09-30; also `skills/retro:44`, `skills/distill:38`, `skills/dream:9` and `README.md:216` restatements — `/dream` Step 1 → `<root>/.agents/memory`; `agents/dreamer.md`; `specs/sh-agent/memory-loop.md` rewritten (header fixed, memory item 5) — verify: `git grep` for the old resolution text in living files — expect: none left — docs: specs/sh-agent/memory-loop.md
- [x] Behavioural run of the changed `agents/dreamer.md` — evidence: 4 of 5 on the matrix check, 5 of 5 on the store path (`testing/2026-09-30-verify-parity-memory-dream-phase-2.md`), 2026-09-30 — `claude -p --agents <json>` with the matrix `dreamer` row, several runs, reported as n of N (AGENTS.md:65) — expect: the artifact reads the store at `<root>/.agents/memory`
- [x] `testing/guide.md` gains the claude-only kit — evidence: § 4's three-setup paragraph, 2026-09-30 — (`CLAUDE_CONFIG_DIR`, a copied credential, chmod 600, deleted by name after) — docs: testing/guide.md
- [x] Checkpoint — expect: E2E 1 on the three setups, E2E 2 on the claude setup (both states), sandbox rules per `AGENTS.md` § Testing — evidence: `testing/2026-09-30-verify-parity-memory-dream-phase-2.md` PASS; gate `testing/2026-09-30-verify-parity-memory-dream-phase-2-gate.md` FAIL on a stale "Claude Code loads" line and the log's unexplained stderr claim, then PASS after the fixes (406 pass, context 21906 bytes), 2026-09-30

### Phase 3 — `sh-agent dream` collection
- [x] `search.ListProjects` — evidence: `TestListProjects` (reads a copy: a read-only open of a WAL db left `-wal`/`-shm` behind, found by the scratch dry run, test red first), `TestListProjectsMissingDBIsEmpty`, 2026-09-30 — verify: unit — expect: rows of a fixture db; missing db → empty, no error
- [x] `search.AgyConversationText` (D8) — evidence: probe D8a; `TestAgyConversationText` (driver returns `time.Time`; the temp-dir cleanup order fixed after a smoke on copies of real data left an empty `~/.config/sh-agent`, test tightened red first); smoke on copies of two real conversations returned the known prompts, 2026-09-30 — first: the probe on a copy of one real conversation db in scratch, showing text runs map to the user's and the model's words — verify: unit on a reduced fixture db — expect: a known prompt's text returned; a conversation of another workspace ignored
- [x] `harness.Redact` — evidence: `TestRedactText` over the shared `shapeSecrets` list (79 cases red against the stub, green after), 2026-09-30 — verify: unit — expect: every needle of `TestRedact*` redacted through it; `$HOME` → `~`
- [x] Claude transcript reader — evidence: `TestCollectInputs` (root and subdirectory lines kept; old, other-project, tool and summary lines dropped), 2026-09-30 — verify: unit on a fixture jsonl — expect: only lines after the watermark, with cwd at or under the root, text parts only
- [x] Collect, watermark rule, the cap, `--dry-run` — evidence: `TestCollect*`, `TestDryRunWritesNothing`; `sh-agent dream --all --dry-run --home <scratch>` over two indexed fixture projects printed both and left the home's file list identical, 2026-09-30 — verify: unit plus `sh-agent dream --dry-run --all --home <scratch>` — expect: counts; nothing written, locked or advanced (asserted) — docs: specs/sh-agent/dream.md, README.md verbs table, main.go usage
- [x] `testing/guide.md` § 3d for the `internal/search` additions — evidence: recall means (all 0.595 over 41 questions), rebuild 1.94 s, warm search 42/35 ms (`testing/2026-09-30-verify-parity-memory-dream-phase-3.md`); the `sh-agent-matrix.md` rows moved to the Phase 6 box with the `memory` rows, 2026-09-30 — (recall, rebuild under 10 s, warm search under 100 ms) and the `sh-agent-matrix.md` rows (AGENTS.md:65)
- [x] Checkpoint — evidence: `testing/2026-09-30-verify-parity-memory-dream-phase-3.md` FAIL on branch coverage, then PASS after the added tests (426 pass; internal/dream 90.7%); the verifier's mutation of the skills rule failed `TestCollectSkillAloneIsNoChange` and `TestCollectUnchanged`, restored byte-identical, 2026-09-30

### Phase 4 — the run and its outputs
- [x] Probe D7: `agy -p` with the prompt — evidence: in the agy-only sandbox home with an EMPTY allow list and `trustedWorkspaces` not naming `dream-work`, `agy -p '<prompt>' --mode plan --disable-slash-commands --output-format text` from the empty `dream-work` answered the requested JSON verbatim, exit 0, empty stderr, nothing written in `dream-work`; `claude -p --tools "" --no-session-persistence --output-format text` with the prompt on stdin did the same and persisted no session, 2026-09-30 — in the empty `dream-work`, in an agy-only sandbox home whose settings file does NOT list `dream-work` as trusted (as in the real home) — expect: text on stdout, exit 0, no "no output produced"; the flags recorded in D7, or D7's refusal branch taken
- [x] Run, parse, render, drafts, state, log, lock, `status` — evidence: `internal/dream/run_test.go` with a fake CLI (tool-free flags, cwd `dream-work`, no planted secret in stdin, file shapes, D23 split, invalid draft name listed, watermark only after the write, invalid output and CLI exit keep the watermark and the raw text, no CLI exits 2 before any project, live and stale lock, 30-day prune, `StatusLine`), 2026-09-30 — verify: unit with a fake CLI (a test binary on PATH that records stdin and prints canned output) — expect: the file shape above; the watermark advances only after the write; invalid output → fail with the raw text kept; a stale lock is replaced — docs: README.md (state rows :212-214, `configDir` :30), specs/sh-agent/dream.md
- [x] `Collect` sends `proposals/dream/harness/*.md` and no other `proposals/dream/` path (D19) — evidence: `TestCollectSendsHarnessDreamProposals` red then green; `TestCollectDoesNotFollowALinkedDreamDir`, whose guard's mutation compiled and failed it, restored byte-identical; lands in the Phase 5 commit (D22), 2026-09-30 — verify: unit, written first — expect: fails before the change (the harness file absent from the bundle), passes after; a project-target file and a draft never in the bundle — docs: README.md report row :214
- [x] `/retro` reads its project's dream proposals (source 9, Status applied/rejected), `/retro --harness` reads `proposals/dream/harness/` — evidence: skills/retro/SKILL.md, 2026-09-30 — `/retro --harness` the harness targets; Status set to applied/rejected as the queue does — docs: skills/retro/SKILL.md
- [x] Checkpoint — expect: E2E 3, 4 and 6 — evidence: `testing/2026-09-30-verify-parity-memory-dream-phase-4.md` (E3: two projects 2+2 proposals with Claude Code, second run `unchanged` twice and no file; E4: a planted secret absent from the prompt, a deleted root `skip gone`; E6: agy-only 3 of 3 after D7b), 2026-09-30

### Phase 5 — the timer (after D3/D12 approval)
- [x] `--install-timer`/`--remove-timer`, `uninstall` removal and its dry run — evidence: `timer_test.go` against fake `systemctl`/`launchctl` (rendered text, calls, identical rerun, differing unit kept and named, no binary, no manager, failing manager leaves files named, remove absent/present, macOS); `TestUninstallDryRunLeavesTheTimer` (no manager call in a dry run; removed without a manifest); `systemd-analyze verify` exit 0 on units rendered into a scratch home through a fake `systemctl`, 2026-09-30 — verify: unit on the rendered text and on `uninstall --dry-run` (no command run); `systemd-analyze verify` on scratch files — expect: pass — docs: deploy/local.md, README.md uninstall row
- [x] README map gains the job (D13), `TestHarnessMap` updated in the same commit (`dream` now required on the map), AGENTS.md rule 1 says the timer is neither a delivery nor a report path — verify: `go test -run TestHarnessMap ./cmd/sh-agent` — expect: pass
- [x] Checkpoint — evidence: the Phase 4–5 verifier log `testing/2026-09-30-verify-parity-memory-dream-phase-4-5.md`, 2026-09-30
- [-] E2E 5, the real timer on this box — moved to the handoff's Next Steps (2026-09-30; who: user, R6) — expect: `list-timers` shows it, then it is gone after `--remove-timer`

### Phase 6 — close
- [x] `testing/guide.md` § 3b (report tier) for Phase 1's change — evidence: § 3b's report dry run in `testing/2026-09-30-verify-parity-memory-dream-phase-4-5.md` (only the harness dream file listed to send); § 3c rows added to `sh-agent-matrix.md` for `memory`, `dream`, the timer and `uninstall`, 2026-09-30 — to `reorder`'s `takesValue` map, and § 3c rows for the `memory`/`dream` verbs — found in: Phase 2 (the guide names the map as a 3b trigger)
- [x] Docs pass (`documentation`), `/review --docs` — evidence: `reviews/2026-09-30-review-cb799ff-f1ca525-docs.md` (APPROVE WITH NITS, Outcome), the documentation passes after rounds 1 and 2, 2026-09-30 — AGENTS.md Layout, :59, and rule 1 (:13: the dream timer is neither a delivery path nor a report path, stated in the same commit), matrix rows; `deploy/local.md` states the timer's limits (no user manager without linger, `Persistent=true` catches up at login; a rerun from a shell with another PATH renders different text, so `--remove-timer` first)
- [x] Antigravity run of each changed skill and agent — evidence: `testing/2026-09-30-verify-real-home-v1.17.0.md` PASS, 2026-09-30 — (`project-memory`, `feedback-log`, `init-project`, `dream`, `retro`, `distill`, `dreamer`) before the tag (AGENTS.md:77), plus the three-setup parity rows
- [-] Release through `/preflight` ("Include this in the next version") — moved to handoffs/2026-09-30-2206-parity-memory-dream-done.md Next Steps (2026-09-30): the operator will cut the version with a later project

## Risks
Assumptions:
- `autoMemoryEnabled: false` stops native writes — E2E 2's must-fail control.
- `claude -p --tools "" --no-session-persistence` takes stdin and prints only the answer — the first real run in Phase 4.
- agy text is extractable from step blobs (D8) — the Phase 3 probe.
- agy answers tool-free (D7) — the Phase 4 probe.
- Claude Code authenticates from a sandbox `CLAUDE_CONFIG_DIR` — the first Phase 2 checkpoint run.

## Decisions
- D1 — 2026-09-30 — R10 as approved — by: "Sounds great. … Go all approved go work."
- D2 — 2026-09-30 — Memory is written and loaded through `sh-agent` on both CLIs, with native auto memory off per project; this supersedes memory items 2 and 3 under R10 — by: memory item 1, approved; R10
- D3 — 2026-09-30 — R3 exception, revised exact text: "One exception: `sh-agent dream --all`, run nightly by an OS user timer that `sh-agent dream --install-timer` installs and `sh-agent dream --remove-timer` or `uninstall` removes, reads every indexed project's trail and memory, the feedback log and both CLIs' transcripts, runs a CLI with no tools over what it read, and writes only under `~/.config/sh-agent/`, its proposals to `proposals/dream/`; it never edits a project, a skill, a rule or the context, and a proposal takes effect only through `/retro`." — by: the operator, "Go all" (2026-09-30), on this text shown in chat
- D4 — 2026-09-30 — The CLI call has no tools; `sh-agent` reads and writes everything. This departs from the option notes behind dreamer item 5 (`.agents/scratch/2026-09-30-dream-background-options.md` § 5, "the job needs Read/Grep/Write … a dedicated settings file"), not from the approved table line — by: implementer; put to the operator (item 5)
- D5 — 2026-09-30 — The timer is opt-in, installed once per machine with `--install-timer` — by: dreamer item 1 names the command
- D6 — 2026-09-30 — Proposals go to `proposals/dream/` with queue-shaped Status lines, and drafts are files under `proposals/dream/drafts/`. This departs from dreamer item 4's top-level path, because `LastRetro` counts top-level `*.md`. `report` already skips subdirectories (report.go:331-339) — by: implementer; departure put to the operator (item 5)
- D7 — 2026-09-30 — agy's tool-free invocation comes from the Phase 4 probe, in a home where `dream-work` is untrusted; the candidates are `--mode plan --disable-slash-commands --output-format text`. When agy refuses (no answer without a tool, or an untrusted-workspace refusal), the job on an agy vendor stops before any project, with `status` "failed: agy cannot run the dream job here (<reason>)". That is put to the operator with two options: a Non-goals exception letting `sh-agent dream --install-timer` add `dream-work` to agy's `trustedWorkspaces`, or an R10 known difference. Nothing is written into agy's settings without that approval — by: implementer, cold read r2
- D7a — 2026-09-30 — The probe settles D7: the candidate flags answer tool-free from an untrusted, empty cwd on agy, and `--tools ""` does on Claude Code; the refusal branch is not needed — by: implementer, probe
- D7b — 2026-09-30 — D7a was wrong for a real prompt. With the 14 KB dream prompt, `agy -p … --mode plan --disable-slash-commands` tried a `command` tool, was auto-denied, and ended with no answer; its stderr also said `--mode plan has no effect while slash command expansion is disabled`. `--json-schema` (with `--output-format json`) still tried a tool. What works is prompt text: the prompt opens, for BOTH CLIs (R10), with "IMPORTANT: every tool call in this session is refused and ends it with no answer … Answer directly with the JSON object only", and agy runs `-p <prompt> --mode plan --output-format text`. The result was 3 of 3 agy runs ok (2 proposals and a draft each) and 1 of 1 Claude Code run ok. A night on which agy still calls a tool fails that project (`fail invalid output`, watermark kept), and the next night retries — by: implementer, Phase 4 e2e
- D8 — 2026-09-30 — Antigravity transcripts: map by `conversation_summaries.db` `workspace_uris`, then extract UTF-8 runs of 20+ characters from the `steps` blobs, deduplicated in order, inside `internal/search`. If the probe cannot isolate the user's and model's text, transcripts are dropped on BOTH CLIs (R10) and that is put to the operator — by: implementer, from recon
- D8a — 2026-09-30 — D8's probe (copies of `conversation_summaries.db` and one conversation db in `.agents/scratch/agyprobe/`): `workspace_uris` is a JSON array of `file://` URIs, and `last_modified_time` a `YYYY-MM-DD HH:MM:SS.fffffffff+00:00` string. The `steps.step_payload` blobs hold the user's prompt and the model's reply as UTF-8 runs (the context run's answer "# AGENTS.md — global agent context / 28 skills and 8 subagents / GANNET / bug" came back verbatim), mixed with ids and permission strings. The extractor keeps runs of 20+ bytes that contain at least three spaces, which drops the id noise, deduplicated in order. D8 holds; transcripts stay on both CLIs — by: implementer, probe
- D9 — 2026-09-30 — The watermark lives in `dream-state.json`, not in `search.db`, which is a cache (R4) — by: implementer
- D10 — 2026-09-30 — The `/dream` skill stays (interactive, deletes, gated); `sh-agent dream` is the background proposal job, and the skill's description names the verb. The verb's output is a proposal queue under `~/.config/sh-agent/`, not an R7 record in the project; `context/AGENTS.md:42` ("every retro, distill and dream leaves a record") keeps meaning the skill — by: implementer
- D11 — 2026-09-30 — Existing Claude default stores keep migrating through `/init-project` Step 5b — by: implementer
- D12 — 2026-09-30 — R4 gains two sentences, pending approval (item 1 and item 2):
  - "`sh-agent dream` also opens Antigravity's conversation databases with it, read-only." The reader lives in `internal/search`, so AGENTS.md:59's one-package rule holds.
  - "`uninstall` removes exactly what its manifest lists, and the dream timer when one is installed (R3)."
  The timer stays out of the manifest — by: the operator, "Go all" (2026-09-30), items 1, 2 and 5 of `.agents/scratch/2026-09-30-parity-operator-items.md`
- D13 — 2026-09-30 — Pending approval:
  - the Non-goals exception gains "and the two memory keys `/init-project` writes into a project's gitignored `.claude/settings.local.json` (`autoMemoryDirectory`, `autoMemoryEnabled: false`; R10)" (item 3);
  - the README map gains the dream job (item 4, AGENTS.md rule 3, R9) in the implementing commit with an ADR
  — by: the operator, "Go all" (2026-09-30), items 3 and 4
- D14 — 2026-09-30 — The skill is named `project-memory`, because `memory` shadows Claude Code's built-in `/memory` (BUG-065) — by: cold read r1
- D15 — 2026-09-30 — `--project <dir>` exists for the e2e and for rerunning one project; the CLI is chosen by the existing global `--vendor` flag, not a new `--cli`; `--dry-run` is required by the global rule that a command with side effects has a read-only mode — by: implementer
- D16 — 2026-09-30 — Loading memory costs one `sh-agent memory show` call per session start on both CLIs. A headless run without a grant loads none, and an interactive Claude Code session asks for permission unless the operator allowed `sh-agent`, while agy on this box approves by its own `toolPermission`. Accepted: the harness behaves the same; the difference is each machine's own CLI permission setting, which the harness never writes (Non-goals) — by: implementer, from cold reads r1 and r2
- D17 — 2026-09-30 — With vendor `all`, Antigravity transcripts go to Claude Code's model on the same machine. Accepted: both are the operator's own accounts, and nothing is sent anywhere else — by: implementer, from cold read r1
- D18 — 2026-09-30 — The `.gitignore` check is `git check-ignore -q` on the target path, not a literal line — by: cold read r1 blocking question 6
- D19 — 2026-09-30 — Harness-target dream proposals go to `proposals/dream/harness/` and `report` sends that directory (R8); project-target proposals and drafts stay local. This supersedes D6's "report never sends" for the harness subset — by: implementer, cold read r2
- D20 — 2026-09-30 — Existing projects keep Claude Code's native memory until `/init-project` reruns there; `sh-agent memory show` prints the rollout hint in each such project — by: implementer, cold read r2
- D21 — 2026-09-30 — Over the 256 KB cap after all transcript text is dropped: trail files are dropped next, oldest mtime first, whole files; memory is always sent whole; a project whose memory alone passes the cap fails with `fail memory over cap` — by: the operator, "Go all" (cold read r3 question 1)
- D22 — 2026-09-30 — Phase order under the approvals: all of D3, D12 and D13 are approved, so no phase waits; the README map, the `report` change and the timer land in ONE commit with the code and docs (R9, AGENTS.md rule 3) — Phase 4's `Collect` change and Phase 5's map and timer go in the Phase 5 commit — by: the operator, "Go all" (r3 question 2)
- D23 — 2026-09-30 — A harness-target proposal is written only to `proposals/dream/harness/…`, never also to the project file; one Status line each — by: the operator, "Go all" (r3 question 3)
- D24 — 2026-09-30 — Antigravity's "since": per conversation, by `conversation_summaries.db` `last_modified_time`; a conversation modified after the watermark is sent whole (steps carry no time) — by: the operator, "Go all" (r3 question 4)
- D25 — 2026-09-30 — `/dream`'s promote test counts a memory without `originSessionId` as one session (its existing rule, `skills/dream/SKILL.md:47`); the second session comes from `feedback.jsonl`. `skills/dream/SKILL.md:47` is added to the touched text to say so — by: the operator, "Go all" (r3 question 5)
- D26 — 2026-09-30 — `memory add|show` root: the main worktree's toplevel (the parent of `git rev-parse --path-format=absolute --git-common-dir`), so every worktree shares one store — by: the operator, "Go all" (r3 question 6)
- D27 — 2026-09-30 — The r3 non-blocking defaults, as the review recommends:
  - the rollout hint also fires when `settings.local.json` is absent on a machine whose manifest vendor is not `antigravity`;
  - Step 5b's empty-object branch writes the two-key object;
  - `--project` takes the absolute git toplevel of its argument and needs no registration;
  - the child's stderr is kept beside its failure file and redacted;
  - `dream-failed/` text is redacted, and files older than 30 days are pruned at the start of a run;
  - drafts get `-2` on a same-day rerun;
  - the Claude transcript reader honours `CLAUDE_CONFIG_DIR`;
  - `ListProjects` opens read-only, never creating the db;
  - the proposal frontmatter `project` carries the root.
  — by: the operator, "Go all" (r3 non-blocking 7–15)
- D8b — 2026-09-30 — Antigravity step types, from copies of two real conversation databases (`.agents/scratch/agyprobe/` and the gate18 handoff run): 14 is the user's input (the run prompt), 15 the model's reply (the context answer, the handoff's planning text), 132 a tool step (a `view_file` of a SKILL.md), 101 a message event. Only 14 and 15 are read — by: implementer, probe
- D28 — 2026-09-30 — After the feature review (`reviews/2026-09-30-review-cb799ff-f1ca525-code.md`, `reviews/2026-09-30-security-cb799ff-f1ca525.md`):
  - the Antigravity child runs from a private HOME holding a copy of the operator's login and settings that grant nothing, removed after the run; any tool call is refused whatever the operator allows interactively;
  - its prompt goes on stdin as one stream-json message (Linux caps one argument at 128 KiB, and argv is readable in /proc); this supersedes D7b's invocation, and D7b's prompt warning stays;
  - Claude Code adds `--strict-mcp-config --setting-sources local`;
  - locks are flock;
  - pre-project failures and a corrupt state are exit 2 and recorded, and the corrupt state is never rewritten;
  - only the user's and the model's Antigravity steps are read;
  - symlinked trail directories are not followed.
  — supersedes D7b's invocation — by: code review and security review, 2026-09-30
- D29 — 2026-09-30 — After two review rounds found leaks in the Antigravity transcript reader (tool output, then tool-call arguments), the operator chose to keep both CLIs' transcripts, filtered, over dropping them or keeping Claude Code's alone: "Keep both, filtered (Recommended)". The residual risk is an unknown field type in the undocumented format; the text is redacted and goes only to the operator's own CLI's model — by: the operator, 2026-09-30
- D30 — 2026-09-30 — The third review of the Antigravity reader (`reviews/2026-09-30-security-e2cffc1-a275b96.md`: argument pieces glued to a printable length byte or split at a multi-byte character) replaced the byte-run scan with a schemaless protobuf wire-format decoder. Each tool argument is then one whole JSON string field, dropped whole. The operator chose it over dropping the transcripts: "Decode protobuf properly (Recommended)". Measured on copies of the 95 local conversation databases: 971 prose leaves kept, 0 with a tool-argument key. Supersedes D8's extraction — by: the operator, 2026-09-30
- Open: redaction beyond credential shapes and `~` — default: none (the same set as `report`)
- Rejected: the timer in the manifest — the next `install` deletes listed files as stale (apply.go:88-109)
- Rejected: `ExecStart` from `os.Executable()` — under `go run` it is a temp build; BinPath is what `update` keeps current
- Rejected: `memory add --title` — the slug is the title; one less value to keep in step
