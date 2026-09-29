---
owner: main
date: 2026-09-29
---
# Behavioural runs — memory loop, Phase 2 (the `dreamer`)   (2026-09-29)

Range reviewed: working tree over d721b95 — `agents/dreamer.md` (new), `skills/retro/SKILL.md` source 8 and record line, `skills/distill/SKILL.md` Step 1, the three `skill-matrix.md` rows. Plan `plans/2026-09-29-dreamer-memory-loop.md`.

## Status: PASS with one soft column — the skill candidate 4 of 6; everything else at its full count

## Kits

`.agents/scratch/memory-probe/`: `dreamer-fixture/` (a trail with two `bugs/INDEX.md` rows for one recurring defect whose prevention rule is written nowhere, and a five-step release procedure recorded in a handoff and a `changes/` entry), `dreamer-store/` (two `feedback` memories on one fact from two sessions, one open `project` memory, an index), `agents-json.py` (the `--agents` JSON built from `agents/dreamer.md`), `run-dreamer.sh`, `run-retro2.sh` (the `dream-probe` fixture with the working-tree `retro`, `dream` and `feedback-log` linked, the dreamer via `--agents`). All runs `claude -p --setting-sources project,local` (plan D15), real home for credentials, prompt on stdin, read-only tools plus Write for the artifact. Seeds asserted before the runs (two distinct session ids, no rule naming the fact in the global context or the clone, `second occurrence` in the index, `same five steps` in the handoff).

## The `dreamer` on the seeded store and trail (d1–d3 first definition; d4–d6 after the `Since`-bounded reading of bug rows and the same-day suffix)

Dispatch: `Store: <store>. Docs root: .agents/docs. Since: all. Promoted: reviews-in-swedish.md, swedish-summaries.md.`

| Run | Artifact `Status: COMPLETE` | Rule row for the promoted pair (both files named) | Rule row for the recurring bug rows | Skill or procedure row for the release steps | `v2-migration` under Left alone | Tree clean apart from the artifact and `.agents/scratch/` |
|---|---|---|---|---|---|---|
| d1 | yes | yes | yes | yes | yes | yes |
| d2 | yes | yes | yes | yes | yes | yes |
| d3 | yes | yes | yes | **no** | yes | yes |
| d4 | yes | yes | yes | **no** | yes | yes |
| d5 | yes | yes | yes | yes | yes | yes |
| d6 | yes | yes | yes | yes (draft `SKILL.md` under `.agents/scratch/dreamer/`) | yes | yes |

The skill candidate is the soft column: 4 of 6. The two misses read the two release records but did not count them as one procedure done twice. The definition was not reworded a third time for it (the second rewording bounded the bug rows and named `changes/` entries as procedure evidence); the rate is reported as it is, and the matrix row's check keeps the column, so a future run that misses it fails the row. Cost about $0.18 per run (the main agent's part; the subagent's tokens are inside).

## `/retro --dry-run` dispatching the dreamer (t1–t2 first wording; t3–t4 after the dry-run clause)

| Run | `## Memory` with the store line | Memory pair as findings with `rule_ref null` | `dreamer` dispatched (Agent tool, transcript) | Dreamer artifact written and cited | Tree clean |
|---|---|---|---|---|---|
| t1 | yes | yes | **no** — "it always writes an artifact file, and a dry run writes nothing" | no | yes |
| t2 | yes | yes | **no**, same reason | no | yes |
| t3 | yes | yes | yes (one dispatch, `subagent_type: dreamer`) | yes, `retros/2026-09-29-dreamer.md` | yes |
| t4 | yes | yes | yes | yes | yes |

The first wording let the run read `--dry-run`'s "nothing written" as covering the dreamer's artifact; source 8 now says the dry run still dispatches, an agent's artifact being trail, not a change. 2 of 2 after the clause.

## E2E in this checkout

1. **Criterion 1, the store:** Step 5b run by hand here (this session is the interactive one): `.claude/settings.local.json` written and parsing (the whole `.claude/` is gitignored in this repo; the Step 6 line added anyway), the store created at `~/.config/sh-agent/memory/-home-kesokaj-workspace-styrman-harness/`, the two real memories and the index moved, the default path empty; `/dream`'s first line by hand: `store: <that path> (2 memory files; default path ~/.claude/projects/<slug>/memory: 0)`, both files indexed. Met.
2. **Criterion 2, the dreamer on the real trail:** first run (`retros/2026-09-29-dreamer.md`) COMPLETE with zero candidates — it read a sample of the bug index and never opened BUG-096; the definition then gained the `Since`-bounded reading of every index row. Second run (`retros/2026-09-29-dreamer-2.md`) COMPLETE with one rule candidate: "read the neighbours before inserting rule text", evidence the read-neighbours memory plus BUG-095 and BUG-096, the grep over every loaded rule cited as empty, a target named. Met on the rerun; both artifacts kept.

## Not exercised

- `/distill` dispatching the dreamer: the Step 1 text is the same dispatch as retro's source 8; no distill probe was run.
- Antigravity: no run of the dreamer (derived `agent.md` present with every tool mapped), the resume clause or the capture rule; the operator's real-home smoke.
- The matrix rows' check commands were not run verbatim (their `<slug>` and `<store>` placeholders); the equivalent conditions were checked by the kit scripts above.
- The interactive settings write by `/init-project` (Claude Code refuses it headless); this checkout's file was written by hand from Step 5b's text.

## Cleanup

The stores these runs created under the real home (`~/.config/sh-agent/memory/*probe-runs*`) are removed after this log, by that name pattern only.

## After the security round (`reviews/2026-09-29-security-memory-loop.md`)

One more seeded `/init-project` run (i4) against the fixed Step 5b: the headless `.claude/` guard refused the settings write as before, and this time the run printed `migration skipped: the setting is not in place`, left both memory files and the index in the default store, reported the one-line settings content, and created the new store empty with mode 700 (`mkdir -m 700`, since a separate `chmod` was refused — the step now says so). Before the fix, runs i1–i3 moved the files despite the refused write; this run shows the guard. The empty-slug stop was checked by hand: outside a repository the slug command prints an empty string, which the step now refuses before creating anything.

## After the store moved inside the project (D17) and the round-2 security fixes — every probe rerun against the final text

All kits now pass `--setting-sources project,local` (the init kit had lacked it: runs i5 and a3 loaded the installed `init-project` and said so themselves — "the installed copy took precedence over the project's link … a probe meant to test Step 5b needs `--setting-sources project,local`" — and declined to mix the copies; recorded, not counted).

| Probe | Runs | Result |
|---|---|---|
| `/init-project`, seeded default store (i6) | 1 | store `.agents/memory/` created, mode 700; both `.gitignore` lines; the settings write refused headless, content reported, `migration skipped`, the default store untouched — 1 of 1 |
| `/init-project`, no default store, a tracked `{}` settings file (a4) | 1 | Step 5b's trust check first: content reported, nothing written into the tracked file, nothing moved, store created, both `.gitignore` lines — the tracked-file case of security finding 1, 1 of 1 (the matrix row now seeds an untracked `{}` so the `no memory files at` branch is what it tests) |
| `/dream` split store (s13–s15: settings store at the clone's `.agents/memory/`, the default path seeded) | 3 | first line with both counts, `split store — merge before running`, no table, both stores unchanged — 3 of 3 |
| `/dream` with a foreign `--settings` value (run 9, the kit not yet repointed: value `<run>/mem`) | 1 | refused as neither of the two exact paths, fell back to the default path, wrote nothing — the round-2 root bound working, 1 of 1 |
| `/dream` regression, store at the clone's `.agents/memory/` (run 10) | 1 | store line names it (6 memory files, default path 0), verdicts and the `promote` flag as before, no split, store unchanged — 1 of 1 |
| `dreamer` (d7–d9, store at the clone's `.agents/memory/`, `Promoted:` the pair) | 3 | artifact COMPLETE, the pair, the recurring bug rows and the project memory left alone 3 of 3; the skill candidate 2 of 3 (6 of 9 over all rounds); no refusal — the exact-path rule accepted the store |
| `/retro --dry-run` dispatching the dreamer (t5, t6) | 2 | dispatched, artifact cited, the pair as findings — 2 of 2 |
| E2E criterion 2 on this checkout's real store (`retros/2026-09-29-dreamer-3.md`) | 1 | COMPLETE, two rule candidates: read the neighbours (memory + BUG-095 + BUG-096) and "a project-supplied setting naming a filesystem destination is untrusted when its file is tracked" (BUG-101 + BUG-102) — met |

E2E criterion 1 after the move: this checkout's store is `.agents/memory/` (three files, mode 700, gitignored), `.claude/settings.local.json` holds its absolute path, the default path is empty, and the dream first line by hand reads `store: <root>/.agents/memory (2 memory files; default path …: 0)`.
