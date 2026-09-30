---
owner: reviewer
date: 2026-09-30
---
# Review — cold-read-parity-memory-dream-r1   (2026-09-30)

Range reviewed: plan file `.agents/docs/plans/2026-09-30-parity-memory-dream.md` (as committed in d6e396a), with the recon brief, `REQUIREMENTS.md` (R1–R10, Non-goals), project `AGENTS.md`, and the code and text it names: `internal/harness/{report,state,apply}.go`, `cmd/sh-agent/main.go` FlagSet, `skills/init-project` Step 5b, `skills/feedback-log` rule 7, `specs/sh-agent/memory-loop.md`, `skills/retro` `--harness`, `README.md` map and state rows, `context/AGENTS.md:42`, BUG-065, the systemd user manager on this box, and the approved options file `.agents/scratch/2026-09-30-dream-background-options.md`.

## Status: COMPLETE

## Verdict: REQUEST CHANGES — design needs revision / implementation not ready

The worst problems are three requirement conflicts and one undeliverable feature. (1) The Antigravity reader has to open SQLite databases from `internal/dream`. `REQUIREMENTS.md:19` (R4) limits the one driver to the search index, and `AGENTS.md:59` names `internal/search` as the only package that imports it. (2) The timer's file removal in `uninstall` contradicts R4's "`uninstall` removes exactly what its manifest lists". Putting the timer in the manifest instead gets it deleted as stale by the next `install` (`apply.go:88-109`). (3) The timer cannot be verified in a scratch home without touching the real one (R6). (4) As designed, the nightly job fails on this box every night: the systemd user manager's PATH is `/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin`, and `claude` and `agy` live in `~/.local/bin`. Separately, the new skill name `memory` collides with Claude Code's built-in `/memory`, which reintroduces BUG-065. Two of the E2E criteria cannot fail as written.

## Findings

### Gaps (question 1)

- **[high / confidence: high]** `plan:65,137,170` vs `REQUIREMENTS.md:19`, `AGENTS.md:59`. D8 reads `~/.gemini/antigravity-cli/conversation_summaries.db` and the `conversations/<uuid>.db` step blobs. Both are SQLite, and the Surfaces list puts the reader in the new `internal/dream/`. R4 says `modernc.org/sqlite` is "used only by the search index" and that any other use "is a requirements change". AGENTS.md:59 says "used only by `internal/search`". The recon flagged this as "needs a decision" (recon:33), and the plan never makes it. Fix: either put an exported read-only function in `internal/search` (and argue that it is still "the search index's" driver), or propose an R4 edit for operator approval. Either choice needs a D-row.
- **[high / confidence: high]** `plan:85` (`uninstall` runs the remove path) vs `REQUIREMENTS.md:19` ("`uninstall` removes exactly what its manifest lists") and `apply.go:268-269` ("removes every path in the manifest and nothing else"). The plan does not say whether `--install-timer` records its unit files in the manifest. Both options fail:
  - Not recorded: `uninstall` deletes files the manifest does not list, which breaks R4.
  - Recorded: the next `install`/`update` treats them as stale and deletes them (`apply.go:88-109`), leaving an enabled timer with no unit file. `status` also reports them as `stale` drift (`apply.go:204`).

  The approved options file (item 1) says "uninstall must remove it (the manifest lists it)". There are two more holes. `Uninstall` returns early when there is no manifest (`apply.go:275-277`), so a timer installed on a machine whose manifest is gone is never removed. The plan also never says that `uninstall --dry-run` must not run `systemctl`/`launchctl` (global hard rule on read-only modes).
- **[high / confidence: high]** `plan:148`, E2E 5 (`plan:111`) vs R6 (`REQUIREMENTS.md:27`). The fixture "a scratch home (`XDG_CONFIG_HOME`/`--home`) with `systemctl --user` on this box" does not produce the condition. The running user manager resolves units only through its own `UnitPath` (read on this box: `~/.config/systemd/user`, `/etc/systemd/user`, …; no scratch path). So `systemctl --user enable --now sh-agent-dream.timer` fails with "Unit file … does not exist" for a unit written under `<scratch>/.config/systemd/user`. Setting `XDG_CONFIG_HOME` on the client does not change the manager's search path. The only way to get the "list-timers shows it" state is to write into the real `~/.config/systemd/user` and have `enable` create `timers.target.wants/` there. That touches the real home and arms a real 03:00 `dream --all` against the real `search.db`. Fix: verify with `systemd-analyze --user verify <scratch unit>` plus a unit test on the text, and move the list-timers check to the operator's real-home release smoke.
- **[high / confidence: high]** `plan:80` (ExecStart) and `plan:98` ("no CLI on PATH: exit 2"). The systemd user manager's environment on this box is `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin`, but `claude` and `agy` are `/home/kesokaj/.local/bin/{claude,agy}` (`command -v`). Every timer run would exit 2 before the first project. The same applies on macOS, where launchd's default PATH is `/usr/bin:/bin:/usr/sbin:/sbin`. It also fails silently: `status` shows `never run`, because `last-run` and the log lines are only written per project. Fix: the unit carries `Environment=PATH=…` captured at `--install-timer` time, or the absolute CLI path resolved then. A pre-project failure also needs a log line and a `status` hint.
- **[medium / confidence: high]** `plan:80` — "the absolute path of the running binary". Under `go run` this is a temp build under the Go cache, and in a checkout it may be `./sh-agent` in the repository. The unit then points at a file that goes away or never moves forward with `update`. Fix: name `BinPath` (`~/.local/bin/sh-agent`, `plan.go:11-17`).
- **[medium / confidence: high]** `plan:76` — the dream lock is O_EXCL with no staleness rule. A run killed by timeout, logout or reboot leaves `dream.lock` behind, and every later nightly run exits 0 with "already running", indefinitely and silently. The memory lock has a 30 s stale rule (`plan:52`); the dream lock needs one too, for example a PID liveness check as in `report.go:500` `alive`.
- **[medium / confidence: high]** E2E 2 (`plan:108`) cannot fail. Step 5b keeps `autoMemoryDirectory` ("beside `autoMemoryDirectory`", `plan:89`), so Claude Code's native writes go to `<root>/.agents/memory`, never to `~/.claude/projects/<slug>/memory/`. The check "nothing new exists under `~/.claude/projects/<slug>/memory/`" passes whether or not `autoMemoryEnabled: false` works. Fix: assert that no file appears in `.agents/memory` other than those written by the verb (for example, no `originSessionId` in any new file), and run the must-fail state first, with the key absent, to show a native write does happen in that sandbox. This is the plan's first listed assumption (`plan:157`), so today it is unchecked.
- **[medium / confidence: medium]** E2E 4 (`plan:110`): "a secret-shaped string placed in a memory file does not appear in the proposal". A model rarely quotes a random token from a memory file, so the check passes with redaction switched off. Fix: assert on the prompt itself, via the fake CLI's captured stdin (Phase 4's test helper) or `--dry-run` output.
- **[medium / confidence: high]** `plan:144` and Surfaces `plan:36` ("`Collect` excludes `proposals/dream/`"). `Collect` already reads only the top-level `.md` entries of `proposals/` (`report.go:331-339`), and a directory named `dream` is neither `.md` nor read recursively. The Phase 4 unit test therefore passes before any change, so it cannot fail. Say that it is a regression pin for existing behaviour, not a code change. Separately, `redactLines` is at `report.go:199`, not `:128`, which is `needleIn`.
- **[high / confidence: high]** `plan:20,128` — the new skill `skills/memory/` takes the name of Claude Code's built-in `/memory` command. This is BUG-065's class exactly (`bugs/INDEX.md:51`: "a skill/agent/command name never shadows a built-in slash command of either CLI"). Rename it, for example to `memory-log` or `project-memory`, and check agy's built-in list too.
- **[medium / confidence: high]** `plan:117-121` vs `REQUIREMENTS.md:31` (R7) and `context/AGENTS.md:42` ("every retro, distill and dream leaves a record in `.agents/docs/retros/`"). `sh-agent dream` writes its per-project output only to `~/.config/sh-agent/proposals/dream/`, and D3 says it "never edits a project". R7 says a project's dream records live in the project's trail. The plan lists `context/AGENTS.md:42` as touched but does not say whether the verb's output is a "dream record" under R7 (D10 keeps the two names on purpose). Stale restatements to fix as well: `README.md:212` ("a dream leaves its record in the project's `retros/` and no line here"), `README.md:30` (the map's `configDir` contents omit `proposals/dream/`, `dream-state.json`, `dream.log`, `dream.lock` and `dream-work/`), and `README.md:214` (the report list).
- **[medium / confidence: medium]** `plan:23` vs `AGENTS.md:15` (rule 3: adding a line to the map needs the operator, an ADR and one commit) and `REQUIREMENTS.md:47` (R9: the procedures and the map change together). The job adds a producer into `configDir` and an edge to `retroHarness` (`README.md:15,30,51`), yet the plan's docs list names `README.md:214` and the verbs table but not the map. Decide whether dream is on the map. If it is, the operator's word and `TestHarnessMap` belong in the same commit.
- **[medium / confidence: high]** `plan:23`, D6 — the consumer is undefined in three places:
  - (a) `/retro --harness` lists proposals by their `**Status:**` line (`skills/retro/SKILL.md:48`), but dream files carry frontmatter plus `## Proposals` and no Status line. Nothing says how one is marked applied or rejected, so each file is re-listed forever.
  - (b) A proposal whose target is a project (kinds `rule`/`memory`) has no consumer. `/retro --harness` runs in the harness checkout and applies harness targets, and project-level `/retro` never reads `proposals/dream/`.
  - (c) On binary-only and zip places, `/retro --harness` never runs ("same machine only", `README.md:51`), and D6 keeps dream files out of `report`. Their output is never seen.
- **[medium / confidence: high]** D3 (`plan:165`) — the R3 exception text says the job writes "to `~/.config/sh-agent/proposals/` only". The design also writes `dream-state.json`, `dream.log`, `dream.lock` and `dream-work/` (`plan:58,71,74,76`). The operator would be approving text that the implementation breaks on its first run. Reword it before asking, for example "writes only under `~/.config/sh-agent/`, and its proposals and draft skills only to `proposals/`".
- **[medium / confidence: high]** Scope items 1 and D1–D6 cite "dreamer table item N" and "memory table item N". The dreamer table exists only in the gitignored `.agents/scratch/2026-09-30-dream-background-options.md`; the memory table was not found anywhere. The approved decisions are therefore not in the trail. The plan also departs from approved items without saying so:
  - item 4 said `proposals/<date>-dream-<project>.md` (top level, "the existing queue") and `proposals/drafts/<name>/SKILL.md`; the plan uses `proposals/dream/` with inline drafts (D6);
  - item 5 said Read/Grep/Write tools with a dedicated agy settings file; the plan goes tool-free (D4);
  - item 1 said "the manifest lists it"; the plan is silent (see the R4 finding above).

  Copy both tables into the plan, and mark each departure as one that needs the operator's word.
- **[medium / confidence: medium]** `plan:89` — `/init-project` Step 5b. Every project already scaffolded falls in the branch "present with the key → keep the file as it is" (`skills/init-project/SKILL.md:115`), so "same insertion method" does not say how `autoMemoryEnabled` gets into a file that already holds `autoMemoryDirectory`. Also, Non-goals (`REQUIREMENTS.md`, the "Editing `settings.json` or any vendor setting … never another key" bullet) is not reconciled with a second vendor-setting key written by a skill. The plan's Out list cites Non-goals only for permissions.
- **[medium / confidence: medium]** `plan:58-66` — the watermark. The plan does not say whether "advances" means run start or write time. With write time, anything changed during a 10-minute run is lost. Memory is collected whole ("every memory file"), so the plan also needs to say what makes a project "unchanged" (E2E 3 depends on it). Claude transcripts also need a rule: file mtime or per-line `timestamp`.
- **[low / confidence: medium]** Risks do not name a cost of R-parity. With native memory off, loading memory needs a `Bash`/`run_command` call to `sh-agent memory show` at every session start, and writing needs one per memory. Both ask for permission because Non-goals bar writing a permission. In headless runs without a grant, the call is refused and memory silently never loads. Either list this as a known difference or accept it in a D-row.
- **[low / confidence: medium]** Cross-vendor disclosure. With manifest vendor `all`, Antigravity transcripts are sent to Claude (recon risk 1). No D-row decides it; the Open row covers only redaction scope.
- **[low / confidence: medium]** `plan:44`. The frontmatter drops `node_type: memory` and `originSessionId`, which `skills/feedback-log/SKILL.md:90` and `memory-loop.md:12` name. `/dream` and `dreamer` read that shape, so say that the spec rewrite changes it.
- **[low / confidence: high]** Docs clauses missing:
  - project `AGENTS.md` Layout (`:36-46`) for `internal/memory` and `internal/dream`;
  - `AGENTS.md:59`, if the driver question moves;
  - the `uninstall` row in README and `deploy/local.md` for the timer removal;
  - `README.md:212,214,30` (see above).
- **[low / confidence: high]** Phase 0 box 1 (`plan:118`) is already done in d6e396a (R10 at `REQUIREMENTS.md:49`, ADR-047 at `DECISIONS.md:501`) but is not ticked. The context budget is 21998 of 22000 bytes (`go test -run RenderedContext -v`), so the Phase 2 resume edit has 2 bytes of headroom; the plan does say it pays in the same edit.
- **[low / confidence: high]** E2E 6 (`plan:112`) applies `--cli`, a `dream` flag, to criterion 1, which is about memory sessions. Reword it as "criterion 1 on the claude and agy homes; criterion 3 with `--cli claude` and `--cli agy`".

### Scope (question 2)
Items the plan builds that `## Why` did not name, with the D-row status of each:
- `sh-agent dream <dir>` (a second, single-project form of the operation): no D-row. It is useful for tests, but the plan should say so or drop it.
- `memory add --title`: an optional value the slug already covers. No D-row, and the memory table that might hold it is not in the trail.
- `--cli claude|agy`: in the approved options item 5 ("overridable with `--cli`"), but that record is scratch only, so it needs a D-row.
- `status` line and `dream.log`: approved options item 6, same record caveat.
- `--dry-run`: required by the global read-only rule. Not a gap.

### Readiness (question 3)
BLOCKING (two careful implementers would build different things):
1. Where does the SQLite read of Antigravity's databases live, given R4 and AGENTS.md:59: an `internal/search` export, or a requirements change?
2. Is the timer in the manifest? How does `uninstall` remove it without breaking R4, and how does `install` avoid deleting it as stale?
3. How does the unit find `claude`/`agy` (PATH or an absolute path captured at install), and which `sh-agent` path goes into `ExecStart`?
4. How is a dream proposal marked consumed, who applies project-target proposals, and what happens to the output on places without a harness checkout?
5. What is the watermark's value (run start or write time), and what counts as "changed" when memory is collected whole?
6. Is "`.gitignore` lacks `.agents/memory`" a literal-line check or `git check-ignore`? A `.agents/` line gives opposite results.
7. How does Step 5b add `autoMemoryEnabled` to a settings file that already holds `autoMemoryDirectory`, and how is Non-goals' "never another key" reconciled?
8. What is the skill named, given the collision with the built-in `/memory` (BUG-065)?
9. Is the dream job on the README map (AGENTS.md rule 3, R9), and is its output an R7 "dream record"?
10. What is the stale rule for `dream.lock`?
11. How is the timer verified without the real systemd manager and the real home (R6)?

NON-BLOCKING (a default is stated, or a reasonable implementer can decide):
- Does `memory add` wait on a held `.lock` or fail at once?
- Claude transcripts from a session started in a subdirectory of the root: exact `cwd` match per `plan:64` misses them; accept that.
- Frontmatter fields `node_type` and `originSessionId`.
- Whether agy needs a workspace-trust entry for `dream-work/`. If so, that is a vendor-setting edit; the D7 fallback covers it.
- Whether to use redaction marker wording other than "[redacted by sh-agent report]".
- Cross-vendor transcript policy (the Open row defaults to none).

`design needs revision` — `implementation not ready`

## Criteria
1. Dead code: not applicable to a plan. One planned change is a no-op: the `Collect` exclusion already holds (`report.go:331-339`).
2. Flow tracing: **fail**. Timer → manifest → `install` stale removal → `uninstall` is untraced. Unit → PATH → CLI fails. Proposal → consumer is undefined for project targets and non-checkout places.
3. Standards conformance: **fail**. R4 (driver outside `internal/search`, uninstall beyond the manifest), R6 (timer verification), BUG-065 (skill name), and R7 and restatements left standing (`README.md:30,212,214`, `context/AGENTS.md:42`). R10 is satisfied in intent. AGENTS.md § Testing (`:65-73`) is cited for the sandbox runs (BUG-115 exists); E2E 3's real `claude -p` under `--home` must keep the real `HOME` in the child's environment (AGENTS.md:65, "a CLI started from a scratch HOME cannot authenticate"), which the plan does not say.
4. Directory layout: pass. `internal/memory` and `internal/dream` fit, but AGENTS.md Layout has no clause for them.
5. Simplicity (DIY-except-security): pass with notes. The Antigravity protobuf string extraction (D8) is not crypto. The `<dir>` form and `--title` are unrequested surface. No new dependency, unless question 1 resolves to one.
6. Decoupling: pass with a note. The exported `harness.Redact` and `search.ListProjects` seams are clean. The timer's reach into `Uninstall` couples `internal/harness` to `internal/dream`; decide which package owns timer removal.

## Not examined
- `agents/dreamer.md` and `skills/dream/SKILL.md`: not read line by line; the plan's cited lines were taken from the recon.
- `skills/distill:38`.
- `internal/search/{search,store}.go` bodies: `tryOpen`'s behaviour on a schema mismatch, relevant to `ListProjects` and "missing db → empty".
- `main.go` status at :508.
- `testing/skill-matrix.md` and `testing/sh-agent-matrix.md`.
- agy's built-in slash-command list.
- Whether `autoMemoryEnabled` is a real Claude Code setting key: not checked against the Claude docs, and still unverified in the plan too.
- macOS launchd behaviour: not tested, since this is a Linux box.
