---
owner: main
date: 2026-09-28
---
# Plan — time-sortable handoff and change-entry names   (2026-09-28)

## Why

The Resume Protocol (`context/AGENTS.md:9`) reads the newest handoff and the three newest change entries. It finds "newest" by file name, and names are `YYYY-MM-DD-<slug>.md`.

When one day has several files, they sort by slug, not by time:
- This repo has 8 days with several handoffs, up to 15 on 2026-09-17.
- On 2026-09-28, `2026-09-28-research-okf.md` sorts last but is the oldest of three (commits 13:25, 13:29, 13:35).
- In a fresh clone every mtime is the clone time, so `ls -t` cannot settle the tie either.

**Measured.** In the 2026-09-28 cold-start probe (`.agents/scratch/agy-probe/`, 20 runs), agy spent about 8 of 17 calls on resume run 1 working out which handoff was newest. It opened the wrong one first and recovered through `git log`.

**The operator's request, 2026-09-28:**
- "go ahead with the naming change";
- "add this to the context so we continiue doing this in all following projects";
- "if ./init-project is run on a active project it should warn for this";
- then "go all" on scope items 1–8.

The scope items, as put to the operator in chat and approved with "go all":
1. Only `handoffs/` and `changes/` get timed names.
2. Local time (`date +%F-%H%M`).
3. Existing files are never renamed.
4. One rule sentence: a date-only name is older than every timed name of its day.
5. `/init-project` update mode warns about date-only names and never renames.
6. The rule goes in `context/AGENTS.md` § Project layout, the handoff skill, `agents/changelog.md`, and every text that shows the old pattern.
7. Full verification: size test, scratch install, sandbox start on both CLIs, behavioural runs, reviewer, `/review --docs`, change entry, handoff.
8. No release in this plan; `/preflight v1.7.0` is asked about after Done.

## Scope

- **In:**
  - handoff and change-entry files are named `YYYY-MM-DD-HHMM-<slug>.md`, from local time (`date +%F-%H%M`);
  - one rule sentence in the global context, including how names without a time order against timed ones;
  - the handoff skill's save step and artifact table;
  - the `changelog` agent's naming;
  - a warning from `/init-project` in update mode when date-only names exist.
- **Out:**
  - renaming existing files, which are immutable or append-only and linked from other records;
  - `plans/`, `reviews/`, `testing/` and `retros/` names, which are found by path or by the active plan;
  - a release: after Done the operator is asked about `/preflight v1.7.0` (scope item 8).
- **Follow-ups noted, not built:**
  - timed names for `reviews/` and `testing/`, if the probe data ever shows a same-day ambiguity cost there.

## Surfaces touched

- `context/AGENTS.md:65` (§ Project layout): one sentence. Rendered context today is 20,329 of 22,000 B; the sentence is about 190 B.
- `skills/handoff/SKILL.md`:
  - `:43`, the artifact table row `changes/YYYY-MM-DD-*.md`;
  - `:145`, `# Save as: $DOCS_ROOT/handoffs/$(date +%F)-<topic>.md`;
  - `:62`, the Step 3.1 date line stays `date +%F`, because it covers frontmatter and headings.
- `agents/changelog.md`:
  - `:3`, the description, which tells dispatchers to pass "today's date from date +%F" (the agent has no Bash, `:4`, so it cannot run `date` itself);
  - `:22-24`, the date rules;
  - `:30`, the lane tree comment;
  - `:39`, the entry heading;
  - `:89`, Records.
- `skills/init-project/SKILL.md`:
  - `:22`, the update-mode paragraph;
  - `:106`, the Step 7 summary.
- `skills/preflight/SKILL.md:86`: "Dispatch the `changelog` agent with today's date for the release entry".
- No Go code and no template file. The `handoffs/` and `changes/` template directories are empty.
- Describing text found by `git grep -nE '(changes|handoffs)/[^ \`]*(date|YYYY|<)'` over `context agents skills commands rules README.md AGENTS.md`: exactly the four hits above (`agents/changelog.md:39,89`, `skills/handoff/SKILL.md:43,145`). The re-run of that grep after Phase 1 must show only timed forms.

## Interfaces

- **File name:** `<docs-root>/handoffs/<stamp>-<slug>.md` and `<docs-root>/changes/<stamp>-<slug>.md`.
  - `<stamp>` is the output of `date +%F-%H%M` (local time), run by the writer when it writes the file.
  - The frontmatter `date:` stays `date +%F`.
- **Context sentence**, appended to the paragraph at `context/AGENTS.md:65` after "…append-only;":

  > `handoffs/` and `changes/` files are named `$(date +%F-%H%M)-<slug>.md`, so name order is time order; a name without the time is older than every timed name of its day.

- **Handoff skill:**
  - `:145` becomes `# Save as: $DOCS_ROOT/handoffs/$(date +%F-%H%M)-<topic>.md`;
  - the table row at `:43` becomes `.agents/docs/changes/YYYY-MM-DD-HHMM-*.md (older entries: YYYY-MM-DD-*.md)`.
- **Changelog agent** (no Bash; D3): the dispatcher passes the stamp, and the agent uses it.
  - `:3`: "passing today's date from date +%F in the prompt" becomes "passing the stamp from `date +%F-%H%M` in the prompt";
  - `:22`'s rule 1 becomes: "Use the stamp given in the dispatch prompt (`YYYY-MM-DD-HHMM`): its first ten characters are the date, and the whole stamp starts the file name. A dispatch that gives only a date gets a date-only file name, flagged in your reply.";
  - rule 2 ("If absent, run `date +%F` yourself", which a Bash-less agent cannot do) is deleted, and rule 3 (derive from a sibling artifact, else `DATE-UNKNOWN`) becomes rule 2 (D4);
  - `:30` becomes `# YYYY-MM-DD-HHMM-short-description.md — …`;
  - `:39` becomes `changes/YYYY-MM-DD-HHMM-<desc>.md`;
  - `:89` becomes `<docs-root>/changes/<stamp>-<slug>.md`.
- **Preflight** `:86`: "with today's date" becomes "with the stamp from `date +%F-%H%M`".
- **`/init-project` update mode**, a sentence added to `:22`:

  > In an update, count the files in `<docs-root>/handoffs/` and `<docs-root>/changes/` whose names start with a date not followed by a four-digit time. A directory that does not exist counts zero; for one that exists, run `ls <dir>`, and when that command fails, Step 7's summary says "<dir> could not be listed; date-only names not counted". Otherwise count its output with `grep -E '^[0-9]{4}-[0-9]{2}-[0-9]{2}-' | grep -vcE '^[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{4}-'`. When the total is above zero, Step 7's summary carries this warning, and nothing is renamed:
  >
  > "N handoff/change files have date-only names. New files are named YYYY-MM-DD-HHMM-<slug>.md; the old names stay, and a date-only name counts as older than every timed name of its day."

  Step 7 at `:106` gains "and the date-only-names warning when Step 1 counted any".

## Failure modes

- **An agent writes a date-only name anyway.** Same-day ties come back for that file only. The context rule and both writers state the format, and the behavioural runs measure adherence (n of N).
- **Machines in different timezones.** Timed names from two machines interleave by local clock. The operator accepted local time (scope #2); a later D-row can switch to UTC.
- **A slug that begins with four digits** (for example `2026-09-28-2024-report`) reads as timed. Accepted: the warning is advisory and renames nothing.
- **The switch-over day** (today here, and the first day in any project): an old-style name sorts after that day's timed names in both byte order and the operator's locale. An agent that ignores the rule sentence picks the date-only file as newest, exactly today's behaviour and no worse.
- **`/init-project` cannot list a directory.** A missing directory holds no date-only names and counts zero. An existing one whose `ls` fails is named in the summary as not counted. Neither fails the run.

## Verification

- **Tiers:**
  - the Go suite, including `TestRenderedContextFitsAntigravityBlock`;
  - a scratch-home install with no drift;
  - a sandbox start of each CLI (guide § 4);
  - behavioural runs (`AGENTS.md` § Testing) of the changed agent, the changed skills and the context rule, each reported n of N;
  - an e2e cold resume on agy.
- **E2E criteria:**
  1. In a fresh clone of this repo after the feature lands, a cold agy session asked for status and next step names the newest handoff (this feature's own timed handoff) as its source, even though today's three date-only handoffs sort after it by name. Run 2 times.
  2. `/init-project` run on a scratch clone of this repo reports the date-only-names warning with a count above zero and renames nothing (`git status` shows no rename). Run on a scratch project whose `changes/` holds only timed names, it reports no warning.

## Phases

- [x] Cold-reader check — design ready and implementation ready (`reviews/2026-09-28-coldread-time-sortable-names-r1.md`, `-r2.md`, `-r3.md`: r3 design ready, implementation ready)

### Phase 1 — Rule text

- [x] Context sentence at `context/AGENTS.md:65` — verify: `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -v`, AI — expect: PASS, with the printed size ≤ 22,000 B.
- [x] Handoff skill `:43,:145` — verify: `grep -n 'date +%F-%H%M' skills/handoff/SKILL.md`, AI — expect: the save line.
- [x] Changelog agent `:3,:22-24,:30,:39,:89` — verify: `grep -nE 'HHMM|%H%M|<stamp>' agents/changelog.md`, AI — expect: five lines; `tools:` unchanged.
- [x] Preflight `:86` — verify: `grep -n 'date +%F-%H%M' skills/preflight/SKILL.md`, AI — expect: the dispatch line.
- [x] `/init-project` `:22`, `:106` — verify: `grep -n 'date-only' skills/init-project/SKILL.md`, AI — expect: both lines.
- [x] Describing-text sweep — verify: the Surfaces grep, AI — expect: no date-only pattern left for `handoffs/` or `changes/`. Positive control: the same grep at HEAD prints the four known hits.
- [x] Checkpoint — expect:
  - gofmt, vet and `go test -race ./...` green;
  - `sh-agent install --home .agents/scratch/tsn` then `status`: drift none;
  - `sh-agent list` parses.

  Evidence: `verifier` log `testing/2026-09-28-verify-time-sortable-names-phase-1.md` PASS: race suite ok, context 20,498 of 22,000 B, scratch install drift none, sweep clean with the positive control firing at HEAD; 2026-09-28.

### Phase 2 — Behavioural runs

- [ ] **Changelog agent**, 3 runs, plus a 3-run must-fail control.
  - Setup: `claude -p --agents <json from agents/changelog.md>` in a scratch clone of a small repo, prompt on stdin. `--allowedTools` names exactly the agent's own tools (Read, Grep, Glob, Edit, Write, plus Agent for the dispatch), with no Bash, the same as production.
  - The dispatch says what the new description tells a dispatcher to say, "Stamp: <output of `date +%F-%H%M`>", and names no file-name format.
  - Control: the same dispatch with the definition from HEAD (`git show HEAD:agents/changelog.md`).
  - Verify: `ls changes/` in each clone. AI.
  - Expect: a name matching `^[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{4}-` in n of 3 for the new definition, and the control's names recorded (a control that also writes timed names means the dispatch alone carries the behaviour, which is reported, not hidden). The tool-call order is read from the subagent transcript.
- [ ] **Handoff skill**, 3 runs. The new `SKILL.md` is copied into the scratch clone at `.claude/skills/handoff-probe/SKILL.md`, and the copy's frontmatter line `name: handoff` is changed to `name: handoff-probe`. Nothing else changes, so the probe cannot collide with the installed `handoff`. It is invoked as `/handoff-probe`, from the real home. Verify: `ls <docs-root>/handoffs/`, AI. Expect: a timed name in n of 3.
- [ ] **Context rule**, 3 runs: `claude -p --append-system-prompt "$(cat <rendered new context>)"`, asked to "write a handoff note for this session into the trail" in a scratch clone, naming no format. Expect: a timed name in n of 3.
- [ ] **Sandbox start of each CLI** (guide § 4).
  - Setup: a sandbox home under `.agents/scratch/`, `sh-agent install --home` into it, credential files copied with `chmod 600` and deleted by name after.
  - Each CLI is asked "how are handoff files named, and which is newest when a date-only and a timed name share a day?"
  - Expect: the timed form and "the timed one" on both CLIs; `cli-*.log` has no parse-error line.
- [ ] **E2E 2 (`/init-project` warning).** Sandbox home as above, on a scratch clone of this repo (must fire) and on a scratch project with only timed names (must stay silent). Expect: the warning with a count, then none; `git status` shows no rename.
- [ ] Checkpoint — expect: every run above at or above 2 of 3 (the init-project and sandbox runs: 1 of 1 each); failures are recorded, not hidden. Evidence: `testing/2026-09-28-verify-time-sortable-names-phase-2.md`, <date>.

### Phase 3 — Done

- [ ] `reviewer` (diff only), then `/review --docs` over the range — expect: findings fixed or rebutted in the report's `## Outcome`.
- [ ] Commit and push the feature, with this feature's change entry and handoff written under timed names.
- [ ] **E2E 1 (agy cold resume)**, 2 runs, on a fresh clone of the pushed HEAD, from the Phase 2 sandbox home with the new harness installed (D2). Question: "State the status and the single next step." Verify: the answer's cited handoff, AI. Expect: this feature's timed handoff named as the source in n of 2; an honest failure otherwise.
- [ ] Checkpoint — expect: tree clean and pushed; handoff written. Evidence: `git status`, `git log @{u}..HEAD` empty, <date>.

## Risks

Assumptions:
- Claude Code accepts a `--agents` JSON built from `agents/changelog.md`. Checked by the first Phase 2 run.
- Copied credentials authenticate a sandbox-home CLI, as guide § 4 says. Checked by the sandbox task's first command.

## Decisions

- D1 — 2026-09-28 — Timed names for `handoffs/` and `changes/` only; local time; no renames; `/init-project` warns and never renames; no release in this plan — by: the operator, "go all" on scope items 1–8.
- D2 — 2026-09-28 — E2E 1 can only see the new context sentence once the context is installed where agy reads it. The real home runs v1.6.0 and R6 forbids installing into it from development. So E2E 1 runs from the Phase 2 sandbox home, with the new harness installed, on a fresh clone of the pushed HEAD, not from the real home — by: implementer.
- D3 — 2026-09-28 — The changelog agent keeps its tools; its dispatcher passes the stamp from `date +%F-%H%M` (the description, preflight's dispatch line), and the agent names the file from it. Adding Bash to a records agent widens what it can do in order to read a clock — by: cold read round 1 (`reviews/2026-09-28-coldread-time-sortable-names-r1.md`, finding 1).
- D4 — 2026-09-28 — The changelog agent's rule 2 ("run `date +%F` yourself") is deleted in the same edit, because the agent has no Bash; rule 3 takes its place — by: cold read round 3 (`reviews/2026-09-28-coldread-time-sortable-names-r3.md`); a bug record follows at Done.
- Open: whether E2E 1 should also be repeated on the real home after the operator next runs `sh-agent update` — default: listed in the handoff's Next Steps as an operator smoke check, not run here.
- Rejected:
  - an `INDEX.md` per trail folder: measured at 240 locate calls against 472 writes (`.agents/scratch/trail-search/`, 2026-09-28);
  - `_` as the separator so that timed names sort after date-only ones: the locale collation `ls` uses ignores punctuation, so it would not hold;
  - renaming today's date-only files: they are immutable records, and other records link them by name;
  - adding Bash to the changelog agent: see D3.
