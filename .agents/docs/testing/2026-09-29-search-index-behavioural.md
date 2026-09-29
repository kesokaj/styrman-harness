---
owner: main
date: 2026-09-29
---
# Behavioural runs — `sh-agent search` in the context, the agents and the skills   (2026-09-29)

Plan: `plans/2026-09-29-search-index.md`, Phase 3. Rows `search-first` and `search-agents` in `testing/skill-matrix.md`.

## Status: PASS on Claude Code (round 2); Antigravity partial — the agents that completed searched, two did not complete (vendor limit)

## Kit

`.agents/scratch/bkit/` (gitignored, ADR-010): `fixture/` is `notebox`, 65 files. It has a Go store, `DECISIONS.md` in which ADR-002 ("sessions in process memory") is superseded by ADR-003 ("sessions in SQLite"), BUG-001 (fixed) and BUG-002 (open) naming `internal/store/store.go`, `docs/usage.md` still saying sessions are lost on restart, 40 filler `changes/` and 15 filler `reviews/` entries, and a `.gitignore` holding `.agents/scratch/*`.

- `run-ctx.sh <with|control> <n>`: a fresh clone and `claude -p`. The `with` arm adds `--append-system-prompt` with the rendered working-tree context. The real home's installed context (v1.12.1, without the entry) is loaded in both arms.
- `run-agent.sh <agent> <n>`: a fresh clone plus a commit adding `Store.Close`, then `claude -p --agents <json built from the working-tree agents/<agent>.md>`, dispatching the agent with a fixed task. For `documentation`, `bugs` and `changelog` the task carries `sh-agent search --refs internal/store/store.go` output (D26).
- `run-agy.sh <case>`: the same cases on `agy -p`, with `HOME` set to a sandbox home installed from the working tree. The token is copied, and a run-specific `settings.json` trusts only the clone and grants `command(git)`, `command(sh-agent)`, `command(ls)`, `command(cat)`, `command(grep)`, `command(date)`, `command(go)`, `command(head)`, `command(find)`, `command(wc)`, `read_file(*)` and `write_file(*)`.
- Every script proves its fixture facts before the CLI starts: the ADR lines, the `Close` commit, the stale doc line, the new binary first on `PATH`, and BUG-001 in the search output. `--allowedTools` names each tool used.
- Tool order comes from the Claude transcripts, with the subagent's own transcript under `<session>/subagents/`. A call counts only when `sh-agent search` is invoked, not when it appears inside a heredoc. Antigravity is judged by its conversation databases, counting the command's own output lines (`index current:`, `referenced by N`), because the definitions themselves carry the command text.

## Round 1 — first wording, clones under `.agents/scratch/` (void for location, kept as evidence of the wording)

The clones sat inside this checkout, so Claude Code loaded this repository's `AGENTS.md` as an ancestor file. All three `changelog` runs wrote their entry into this repository's own trail. The stray untracked file was deleted, and re-applying the edit script to a clean `git archive HEAD` showed no tracked file touched (BUG-106).

| Case | Ran `sh-agent search` |
|---|---|
| Context, with the entry (without "Run it before Grep") | 0 of 3 (Grep, then Read) |
| Context, control | 0 of 3 |
| `security` | 3 of 3 |
| `reviewer` | 1 of 3; one run passed `store.go` and got nothing (D28) |
| `explorer` | 0 of 3 |
| `dreamer` | 0 of 3 (it used the Grep tool for its rule check) |
| `documentation` / `bugs` | 3 of 3 fixed `docs/usage.md` / 3 of 3 cited BUG-001 in BUG-003 |
| `changelog` | void (wrote into this repository) |

## Round 2 — final wording (D29, D31), clones under `~/sandbox/search/cruns/`

| Case | Result |
|---|---|
| Context, with the entry | 2 of 3 ran `sh-agent search session storage` as the first call; 1 of 3 used Grep first. All 3 answered ADR-003, superseding ADR-002 |
| Context, control | 0 of 3 (Grep first); all 3 answered correctly |
| `reviewer` | 3 of 3 ran `sh-agent search --refs … --kind bug` as call 2 (call 1: `git diff`) |
| `security` | 3 of 3 ran it (calls 2 to 4) |
| `explorer` | 3 of 3 ran it as its first call |
| `dreamer` | 3 of 3 ran it before writing its rule row |
| `documentation` | 3 of 3: `docs/usage.md` no longer says sessions are lost on restart |
| `bugs` | 3 of 3: the new BUG-003 names BUG-001 |
| `changelog` | 3 of 3: the entry, written in the clone, cites BUG-001 or ADR-003 |

After round 2, `git status` of this repository showed only the intended edits and the BUG-106 record.

## Antigravity (agy 1.2.13) — round 1, first wording

No run's stderr has an "auto-denied" line.

| Case | Result |
|---|---|
| Context question | answered correctly; no `sh-agent search` output in the conversation |
| `reviewer` / `security` | report written; the subagent's conversation holds `sh-agent search` output (3 lines each) |
| `explorer` | brief written; no search output |
| `dreamer` | no artifact: print mode ended ("I have dispatched the dreamer subagent and am waiting") before the subagent returned |
| `documentation` | `docs/usage.md` updated |
| `bugs` / `changelog` | BUG-003 names BUG-001 / the entry cites them |

## Antigravity — round 2, final wording (context, `explorer`, `dreamer`, `reviewer`, `security`)

No run's stderr has an "auto-denied" line.

| Case | Result |
|---|---|
| Context question | answered correctly (ADR-003, still in force); no `sh-agent search` output, 0 of 1 |
| `explorer` | brief written; the subagent's conversation holds 9 lines of search output (round 1: none) |
| `reviewer` | report written; 3 lines of search output |
| `security` | void: no findings file, no search output; print mode ended before the subagent returned (round 1 of the same agent completed and searched) |
| `dreamer` | void again: no artifact in either round |

Summed over both rounds on Antigravity: among the subagent runs that completed, `reviewer` 2 of 2, `security` 1 of 1 and `explorer` 1 of 2 (the one on the final wording) searched. The context question was followed 0 of 2 times on agy. `agy -p` ending before a dispatched subagent returns is the vendor behaviour the context's headless rule is about. It is recorded here, not fixed: the harness cannot make `agy -p` wait.

## What was not exercised

- The fixture is small (65 files). Grep answers the context question in one call there, so the context result says whether the entry is followed, not whether it saves tokens. Phase 5's larger-project task measures that.
- Interactive sessions; only headless runs.
- The skills' new lines (`review-standards`, `dream`, `handoff`, `bug`, `feature`, `retro`, `research`, `init-project`) were not run behaviourally. Their install was checked (`install` and `status` drift none), and `/review --docs` reads them at Done.
