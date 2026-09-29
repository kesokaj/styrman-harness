---
owner: main
date: 2026-09-29
status: in progress — cold-reader check next; ROADMAP "Now" names this plan
---
# Plan — one memory store for both CLIs, and the `dreamer` agent   (2026-09-29)

## Why

The operator, 2026-09-29: "Cant we build dreaming and memory as a subagents that triggers when something in the project/conversation highlights to generate a skill/rule or procedure? … I want them to work the same for claude and agy." Research `plans/2026-09-29-research-memory-subagents.md` compared four shapes; the operator chose B ("Option B", 2026-09-29): one `dreamer` subagent for the consolidation judgement, capture kept in-thread, and a vendor-neutral store by pointing Claude Code's `autoMemoryDirectory` at a harness-owned directory. Spec: `specs/sh-agent/memory-loop.md`.

## Scope

- **In:**
  - The store at `~/.config/sh-agent/memory/<slug>/`, set up by `/init-project` (new and update paths): the `autoMemoryDirectory` key in `.claude/settings.local.json`, the `.gitignore` line, the directory, and the move of an existing default store.
  - The Resume Protocol clause for Antigravity (`context/AGENTS.md`), the Antigravity capture rule in `feedback-log`, and `/dream`'s split-store check.
  - `agents/dreamer.md` per the spec; `/retro` Step 1 source 8 and `/distill` Step 1 dispatching it; two `skill-matrix.md` rows.
- **Out:**
  - Sending memory through `sh-agent report` (D1); any Go change; any hook or timer; a `memory` capture subagent (research, D-row there); auto-applying anything the `dreamer` proposes.
  - Moving stores of projects other than the one `/init-project` runs in.
- **Follow-ups noted, not built:** a `sh-agent status` line naming the project's store (a Go change, out of a file drop); a `--since` argument on the dreamer dispatch beyond the newest-record date.

## Surfaces touched

- `skills/init-project/SKILL.md` — a new Step 5b between the symlinks (Step 5, lines 79–92) and README/.gitignore (Step 6, lines 94–104); Step 6's `.gitignore` list gains one line; Step 7's summary names the store.
- `context/AGENTS.md` § Resume, then checkpoint (line 9): one clause; headroom 357 bytes (`TestRenderedContextFitsAntigravityBlock`, 21,643 of 22,000 on 2026-09-29).
- `skills/feedback-log/SKILL.md` § Behavior Rules (line 80): one rule for Antigravity.
- `skills/dream/SKILL.md` Step 1: the split-store check.
- `agents/dreamer.md` — new; `internal/harness/derive.go`'s `toolMap` already maps every tool it names (Read, Grep, Glob, Bash, Write), so no Go change.
- `skills/retro/SKILL.md` Step 1 (after source 7) and Step 4 (the record cites the artifact); `skills/distill/SKILL.md` Step 1 (dispatch first).
- `.agents/docs/testing/skill-matrix.md` — rows `init-memory` and `dreamer`.
- `README.md` — nothing: no line on the map moves (`sh-agent report` ships the same files; the store is machine-local like `feedback.jsonl`), so project rule 3 does not fire; `/review --docs` confirms. `specs/sh-agent/delivery.md` — nothing.

## Interfaces

- **Store path:** `~/.config/sh-agent/memory/<slug>/`, `<slug>` = the absolute repository root with every character other than a letter, a digit or a hyphen replaced by a hyphen (`git rev-parse --show-toplevel | sed 's#[^A-Za-z0-9-]#-#g'`) — Claude Code's `<project>` name under `~/.claude/projects/`, read off this box: `/home/kesokaj/workspace/styrman-harness` → `-home-kesokaj-workspace-styrman-harness`, and `…/styrman-harness/.agents/scratch/dream-probe/runs/1/clone` → `…-styrman-harness--agents-scratch-dream-probe-runs-1-clone` (the `.` became `-`, the hyphens stayed); D9. Layout as the spec says.
- **`/init-project` Step 5b — the memory store** (new and update paths, after the symlinks):
  1. `slug=$(git rev-parse --show-toplevel | tr / -)`; `store=~/.config/sh-agent/memory/$slug`; `mkdir -p "$store"`.
  2. `.claude/settings.local.json`: absent → write `{"autoMemoryDirectory": "~/.config/sh-agent/memory/<slug>"}`; present without the key → add the key with the Edit tool as one new line after the opening `{` (`"autoMemoryDirectory": "~/.config/sh-agent/memory/<slug>",`), touching no other line — never `jq` or a re-serialisation, which reformat the file — then `python3 -c 'import json,sys; json.load(open(sys.argv[1]))' .claude/settings.local.json` to prove it still parses; say what was added; present with the key → keep it, report its value, and use that value as `store` for step 3.
  3. Migration: `old=~/.claude/projects/$slug/memory`; when `old` holds files and `store` holds none, `mv "$old"/* "$store"/` and report each file; when both hold files, move nothing and report both paths (the operator merges); when `old` is absent, print `no default store at <old>` and, as a hint, every `~/.claude/projects/*/memory` directory whose name ends with the repository's last path component; when `old` exists and is empty, print `default store empty at <old>`; nothing is moved in either case.
  4. Step 6's `.gitignore` list gains `.claude/settings.local.json`.
  5. Step 7's summary names the store path and what moved.
  Headless refusal of writes under `.claude/` (the existing Step 5 note) applies: report the file's exact content for the user to create.
- **Resume clause** (`context/AGENTS.md` § Resume, appended to the reading order): `→ on Antigravity, MEMORY.md in the project's memory store (~/.config/sh-agent/memory/<slug>/) when it exists; Claude Code loads it itself.` (about 140 bytes).
- **`feedback-log` rule** (§ Behavior Rules): on Antigravity, an entry that is a correction or a stated preference is also written as one memory file in the store — Claude Code's shape (frontmatter `name`, `description`, `metadata.type` `feedback` or `user`, `originSessionId`, `modified` from `date -u +%FT%TZ`), body with `**Why:**` and `**How to apply:**` — plus its `MEMORY.md` line; when the store does not exist, nothing (the project has not run `/init-project` since the store existed). On Claude Code nothing changes: its auto memory writes the file itself.
- **`/dream` Step 1 addition:** after resolving `<store>` from settings, also `ls ~/.claude/projects/<slug>/memory/`; when both hold memory files, print both paths with their counts, say "split store — merge before running", and stop with nothing written.
- **`agents/dreamer.md`:** frontmatter `name: dreamer`, one-line double-quoted `description` (what it returns, when to dispatch, never edits), `tools: Read, Grep, Glob, Bash, Write`, `model: sonnet`, `effort: medium`, `maxTurns: 40`. Body per the spec: write early; what it reads; the candidate tests; the artifact template; Records (only its artifact and drafts under `.agents/scratch/dreamer/`); Output (at most ten lines). Dispatch contract: `Store: <path>. Docs root: <docs-root>. Since: <date|all>. Promoted: <memory files, or none>.` plus "Write your artifact with the Write tool; report only means no source edits." and a tool-call budget.
- **`/retro` Step 1 source 8:** "**Dreamer candidates** — after source 7, dispatch `dreamer` with `Store:`, `Docs root:`, `Since:` and `Promoted: <the memory files source 7 flagged promote, or none>`, a 30-call budget and the artifact sentence. The dreamer does not re-derive `promote`: its rule candidates are the promoted memories (with a target file and the rule text), bug rows whose prevention rule is written nowhere, and trail rows; source 7's `promote` findings are replaced by the dreamer's rows for the same memory files, one row per memory, so no memory appears twice; skill and procedure candidates are noted for `/distill`." Step 4's record cites the artifact path under `## Findings`. (D8)
- **`/distill` Step 1:** first dispatch `dreamer` the same way (`Promoted: none` when `/dream` did not run); its skill and procedure candidates are the first entries of the candidate list, with their evidence; the session transcript scan follows as today and adds only patterns the artifact does not already name (same name or same steps), so a pattern appears once. (D8)
- **Matrix rows:** `init-memory` (base: demo HEAD scaffolded by an older harness, plus a seeded default store `~/.claude/projects/<slug of the clone>/memory/` with two files and an index; prompt: `/init-project` update path; check: `.claude/settings.local.json` holds the key with the clone's slug, the store holds the two files and the index, the old directory is empty, `.gitignore` has the line); `dreamer` (base: demo HEAD plus a seeded store — two `feedback` memories on one fact from two sessions, one `project` memory — and a trail with two `bugs/INDEX.md` rows for one recurring defect whose prevention rule is written nowhere and a release procedure of five steps recorded twice (a handoff and a `changes/` entry); dispatch text as the contract with `Promoted: <the two feedback files>` (what `/retro` passes after `/dream`); check: the artifact exists with `Status: COMPLETE`, a rule candidate naming both promoted files with a target file and rule text, a rule candidate naming the bug rows, a skill or procedure candidate naming the release steps, "Left alone" naming the project memory, and `git status --short` shows only the artifact and `.agents/scratch/`); `init-memory-absent` (the `init-memory` base without the seeded default store; check: the output contains `no default store at`, `.claude/settings.local.json` holds the key, the store directory exists and is empty).

## Failure modes

| Situation | Result |
|---|---|
| `/init-project` cannot write under `.claude/` (headless guard) | the settings file's exact content is reported for the user; the store is still created and the migration still runs |
| The setting is ignored (untrusted workspace) | Claude writes to the default store; `/dream` reports a split store and stops; the `dreamer` reads the settings store and lists the other as not examined |
| `~/.config/sh-agent/` denied to an Antigravity session | the resume clause reads nothing and says so; `feedback-log` logs the entry to chat as it does for the log itself |
| Old and new stores both hold files at migration | nothing moved, both paths reported |
| No default store at the computed path (a slug mismatch, or a project Claude Code never wrote memory for) | `no default store at <old>` plus the hint list; nothing moved; the store is still created and the setting written |
| `dreamer` cut at its turn limit | partial artifact, `IN PROGRESS`, named as such by the caller |
| `dreamer` finds nothing | an artifact with an empty Candidates table and a full "Left alone" |

## Verification

Tiers: unit (the Go suite; the size test for the context clause) / install (scratch home, drift none, `dreamer` under both vendors' agent directories, `sh-agent list` naming it) / behavioural on Claude Code / real-home smoke and an Antigravity run (`who: user`, D-row below) / docs review.

- **Behavioural, Phase 1:** a scratch clone of a small repository scaffolded by an older harness (the matrix's `init-update` base), a seeded default store under the real home's `~/.claude/projects/<slug of the clone>/memory/` (per-project, harmless, removed after); `claude -p` running `/init-project` with the update-path answers, 3 runs; expect the `init-memory` row's check in 3 of 3; then the `init-memory-absent` variant, 2 runs, expecting its check in 2 of 2. Then in one of those clones, a `claude -p` "remember that the demo deploys on Fridays" with `.claude/settings.local.json` in place; expect a new memory file in the new store and none in the default (the positive control for the setting being honoured; if Claude Code refuses the project setting headless as untrusted, the run says so and the split-store path is what gets exercised instead — either result is recorded). `/dream --dry-run` on a fixture with files in both stores: expect the split-store stop line, nothing written, 3 of 3.
- **Behavioural, Phase 2:** `claude -p --agents <json built from agents/dreamer.md>`, the dispatch contract on stdin, the seeded store and trail of the `dreamer` matrix row, 3 runs; expect the row's check in 3 of 3, the artifact carrying the definition's own section headings (only the new definition writes them), executed and refused calls counted separately. Then `/retro --dry-run` with the working-tree retro skill linked, 2 runs: expect the artifact dispatched and its rule candidates as findings with `rule_ref null`.
- **Docs:** `/review --docs` over the range; the living docs include the spec, `init-project`, `feedback-log`, `retro`, `distill`, `dream`, `context/`.

E2E criteria: (1) in this checkout, the `/init-project` update path writes `.claude/settings.local.json`, moves the two real memories into `~/.config/sh-agent/memory/-home-kesokaj-workspace-styrman-harness/`, and `/dream --dry-run` prints that store line with two files; (2) `/retro --dry-run` here dispatches the `dreamer`, whose artifact names BUG-096's prevention rule as a rule candidate (it is written nowhere loaded) and the read-neighbours memory as its evidence.

## Phases

- [ ] Cold-reader check — design ready and implementation ready (`reviews/` paths), or the operator's same-turn override quoted with its date
### Phase 1 — the store
- [ ] `skills/init-project/SKILL.md` Step 5b, Step 6 line, Step 7 summary — verify: behavioural tier (`init-memory` 3 runs, `init-memory-absent` 2 runs), AI — expect: 3 of 3 and 2 of 2 — docs: `skills/init-project/SKILL.md` (it is `/init-project`'s documentation) and `specs/sh-agent/memory-loop.md`
- [ ] `context/AGENTS.md` resume clause — verify: `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1` bare, AI — expect: PASS under 22,000
- [ ] `skills/feedback-log/SKILL.md` Antigravity capture rule — verify: read against the store layout in the spec, AI — expect: the rule names every frontmatter field the spec lists
- [ ] `skills/dream/SKILL.md` split-store check — verify: behavioural tier, AI — expect: the stop line in 3 of 3, nothing written
- [ ] The setting honoured headless — verify: the "remember" run above, AI — expect: the memory file in the new store, or the refusal recorded
- [ ] Checkpoint — expect: the store set up by `/init-project` and proven on the fixture; the resume clause installed under both vendors — evidence: `testing/<date>-verify-dreamer-memory-loop-phase-1.md`, <date>
### Phase 2 — the dreamer
- [ ] `agents/dreamer.md` per Interfaces — verify: install tier, AI — expect: `sh-agent install` succeeds (every tool mapped), `status` drift none, the agent under both vendors' agent directories, `sh-agent list` naming it
- [ ] Behavioural runs, N=3 — verify: behavioural tier, AI — expect: the `dreamer` row's check in 3 of 3
- [ ] `skills/retro/SKILL.md` source 8 and record citation; `skills/distill/SKILL.md` Step 1 — verify: two `/retro --dry-run` probes, AI — expect: the artifact dispatched and its rule candidates as findings — docs: the two skill files are their own documentation; the spec's Consumers line
- [ ] `skill-matrix.md` rows `init-memory` and `dreamer` — verify: the checks run on the fixture clones, AI — expect: PASS
- [ ] E2E in this checkout (criteria 1 and 2) — verify: e2e, who: both (the migration moves the operator's real memory files; the operator says go) — expect: as stated
- [ ] Real-home install and smoke, plus the first Antigravity run of the resume clause and `/dream` — verify: real-home tier, user — expect: three checks PASS on each CLI; an agy session's first lines name the store's `MEMORY.md`
- [ ] `/review --docs` over the range — verify: docs tier, AI — expect: no finding, or each fixed or rebutted
- [ ] Checkpoint — expect: Definition of Done closed, `changes/` entry, commit, push — evidence: verify log path, <date>

## Risks

Assumptions:
- Claude Code honours `autoMemoryDirectory` from `.claude/settings.local.json` in a headless run of a scratch clone (trust) — checked by the "remember" run, Phase 1; the split-store path covers the refusal.
- Claude Code creates the store directory when absent — not relied on: `/init-project` runs `mkdir -p`.
- The Antigravity CLI can read `~/.config/sh-agent/memory/` from a project session — the same directory `feedback-log` already writes to on agy (one approval per session per the README § harness state); checked in the operator's smoke, Phase 2.
- `~/.config/sh-agent/` is 0700 and owned by the operator (the installer tightens it) — the store inherits that; checked by `stat -c %a` in the Phase 1 log.

## Decisions

- D1 (2026-09-29): memory stays out of `sh-agent report` — by: implementer, from `internal/harness/report.go:263-288` (an allowlist: `feedback.jsonl`, `retros.jsonl`, `proposals/`). Why: memory holds facts about the operator; the `dreamer` artifact in `retros/` is what carries learning into the project, and the report already ships `retros.jsonl`. The research record's sentence "memory rides the feedback loop for free" was wrong and is corrected there. Rejected: extending the allowlist — a Go change and a redaction question for a personal store.
- D2 (2026-09-29): the store's per-project name is Claude Code's slug, so one rule names both the default and the harness store — by: implementer. Rejected: a hash or a project id — nothing on agy would map it back.
- D3 (2026-09-29): capture stays in-thread; on Antigravity the `feedback-log` skill writes the memory file — by: the research record's recommendation B, chosen by the operator. Rejected: a `memory` subagent — blind to the conversation.
- D4 (2026-09-29): the `dreamer` writes no `retros.jsonl` line; the retro or distill record that consumed it cites it — by: implementer. Rejected: a fourth record kind — nothing reads it.
- D5 (2026-09-29): `.claude/settings.local.json`, gitignored, carries the key — by: implementer. Why: the path is per machine (`~/`), and the docs read `autoMemoryDirectory` from the local scope. Rejected: `.claude/settings.json` (committed, and the harness's own `.claude/*` are symlinks to `.agents/` where a settings file would be shared).
- D6 (2026-09-29): the trust rule is handled by detection, not prevention — `/dream` stops on a split store — by: implementer, from the memory docs § Storage location ("honors it under the same workspace trust rule as hooks").
- D7 (2026-09-29): the sandbox CLI-start tier is substituted as in the dream plan's D10 — Claude Code's start is the behavioural run, Antigravity's is the operator's smoke — by: implementer.
- D8 (2026-09-29): one judgement, not two — the `dreamer` consumes `/dream`'s verdict table (the `Promoted:` line of its dispatch) and never re-derives `promote`; `/retro` replaces source 7's `promote` findings with the dreamer's rows for the same memory files; `/distill` takes the dreamer's skill and procedure candidates first and adds from the transcript only what the artifact does not name — by: cold read round 1, blocking gap 1 (`reviews/2026-09-29-review-dreamer-memory-loop-coldread-r1.md`), resolved on the operator's "Go all on the cold-read gaps and keep building" (2026-09-29). Rejected: dropping `/dream`'s `promote` flag — `/dream` runs standalone too, where no dreamer is dispatched.
- D9 (2026-09-29): the slug rule is "every character other than a letter, a digit or a hyphen becomes a hyphen", read off this box's `~/.claude/projects/` names (a `.` became `-`, hyphens stayed); `/init-project`'s migration prints "no default store at <old>" when the directory is absent, so a mismatch is visible, and lists any `~/.claude/projects/*/memory` directory whose name ends with the repository's last path component as a hint — by: cold read round 1, gap 2.
- D10 (2026-09-29): the settings edit is one inserted line with the Edit tool plus a parse check; "byte for byte" was unachievable with `jq` — by: cold read round 1, gap 3.
- D11 (2026-09-29): the dreamer's rule-candidate sources are the `Promoted:` files (given), recurring bug rows with no written rule, and a memory a trail row restates — never the two-memory test, which is `/dream`'s; the matrix row passes the seeded pair on the `Promoted:` line as `/retro` would — by: cold read round 2, G4 (`reviews/2026-09-29-review-dreamer-memory-loop-coldread-r2.md`), on the operator's "Go all on the cold-read gaps" (2026-09-29). The spec's Candidate tests line now says so.
- D12 (2026-09-29): an absent default store at migration prints `no default store at <old>` with the hint list, an empty one prints `default store empty at <old>`; both in Interfaces, the Failure modes and the `init-memory-absent` runs — by: cold read round 2, G5, on the same words.
- Open: none.
