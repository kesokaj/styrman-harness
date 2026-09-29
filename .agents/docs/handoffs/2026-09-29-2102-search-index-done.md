---
owner: main
date: 2026-09-29
---
# Handoff — 2026-09-29 — `sh-agent search` built; nothing queued in the trail

Session ID: f538b90a-5aa7-44d1-8311-a160b5bb0167
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main`, clean after this commit, pushed. Not released: v1.12.1 is the newest tag.

## Context

The operator asked for research on an embedded SQLite index for the agent. The first verdict was "build nothing" at this repository's scale; the operator then chose to build it for understanding and memory. The plan was `plans/2026-09-29-search-index.md`, now closed: a per-project index and a global search index, context and agent instructions, a recall test that gates vectors, and sandboxed efficiency runs.

## Completed

- **Research:** `plans/2026-09-29-research-sqlite-index.md`. **Recon and consumer audit:** `plans/2026-09-29-recon-project-index.md`, `plans/2026-09-29-audit-search-consumers.md`.
- **R4** replaced on the operator's approval of the exact text: one dependency, `modernc.org/sqlite` v1.46.1, not vendored. ADR-045 and ADR-045a.
- **`sh-agent search`** (words, `--refs`, `--open`, `--kind`, `--limit`, `--global`) and **`sh-agent index`** (`--rebuild`, `--global`), in `internal/search/`. Spec: `specs/sh-agent/search.md`.
- **Instructions:** the context entry "Search first" (21,988 of 22,000 bytes), plus lines in 7 agents and 8 skills.
- **Measured:**
  - Recall: 0.677 against grep's 0.379 on 41 questions.
  - Vectors: not adopted (+6.4 points, under the 10-point bar).
  - Behavioural: every changed agent searched in 3 of 3 runs on Claude Code; the context entry drove a first-call search in 2 of 3.
  - Efficiency on sample apps: no measurable difference. The main agent searched in 0 of 28 real sessions.
- **Reviews:** code APPROVE twice; docs lens with Outcome; security three rounds with Outcome (the cap reached; the operator chose at each gate). The final verifier PASSED at 4c9a327.
- **Bugs:** BUG-106 (run clones inside the checkout) and BUG-107 (test runs registered in the real home's global search index), both graduated to the project `AGENTS.md` § Testing.

## Open / In Progress

- Nothing in the trail: the plan is closed, and `ROADMAP.md` and `TODO.md` are empty.
- Requirements edits awaiting approval: none.

## Key Decisions

See the plan's `## Decisions`, D1–D37. The ones a next session must know:

- **D21:** the driver is pinned at v1.46.1 because later releases raise the `go` directive past 1.24.4. **D25:** `govulncheck` shows one unreachable `x/sys` advisory whose fix needs Go 1.25.
- **D26:** `documentation`, `bugs` and `changelog` have no shell, so the main agent puts `sh-agent search --refs` output in their dispatch (the operator's choice).
- **D30:** vectors closed (below the gate; the "pure Go" model needs ONNX Runtime).
- **D34:** registration stays automatic for checkouts with `.agents/docs` (the operator's choice; the crafted-clone residual is accepted). **D35:** a file edited and restored while `search` waits for the lock can stay indexed with the edit until `index --rebuild` (the operator's choice at the security cap).
- **D36:** the context entry does not make the main agent search in real work (0 of 28); `security` and `reviewer`, which have it as a fixed step, made all 33 calls.

## Corrections

- "on the audit list, why not all places that generate lots of md files?" — feedback log 2026-09-29T21:02:42+02:00. The consumer audit had covered only existing grep steps, and the writer agents were added (D13).

## Agent notes (unverified)

- The main agent likely ignores the context entry because Claude Code's own system prompt steers it to the built-in Grep tool; subagents with an ordered step comply. A step in the resume procedure, or a skill, might change that. Untested.
- `agy -p` ends before dispatched subagents return often enough that Antigravity app-lane comparisons are unreliable. Any future agy efficiency run needs a way to wait, or has to measure single-session tasks.
- The rendered context has 12 bytes of room left, so any future context line needs a cut.

## Environment State

- Branch `main`, clean, pushed. No running processes. Tree-local scratch: `.agents/scratch/{bkit,eff-kit,vecprobe,venv,…}` (gitignored).
- Outside the checkout: `~/sandbox/search/` holds the run clones and sandbox homes. The copied agy tokens and settings were deleted.
- The real home is still at v1.12.1: no `search` verb, no new context. The real `~/.config/sh-agent/` has no `search.db`.

## Artifacts to Read

- `.agents/docs/plans/2026-09-29-search-index.md` — the closed plan, D1–D37, and the one follow-up.
- `.agents/docs/specs/sh-agent/search.md` — commands, storage, freshness, failure modes, measurements.
- `.agents/docs/testing/2026-09-29-search-efficiency.md` — the efficiency runs and what they show.
- `.agents/docs/testing/2026-09-29-search-recall.md` — the recall test and the vector gate.
- `.agents/docs/testing/2026-09-29-search-index-behavioural.md` — which agents search, n of N.
- `.agents/docs/reviews/2026-09-29-security-search-index.md` — three rounds with the Outcome.
- `.agents/docs/changes/2026-09-29-2102-search-index-phase-5-done.md` — the last change entry, and the earlier ones per phase.

## Next Steps

1. **Operator: decide main-agent adoption** (plan § Follow-ups, D36). Options: a line in the context's resume procedure (needs a cut, 12 bytes free), a skill that triggers on "where did we decide" questions, or leave search to the subagents. My recommendation is to try the resume-procedure step first and measure it with the same `search-first` probe.
2. **Operator: release when ready.** `/preflight v1.13.0` (a MINOR release: new verbs, nothing breaks), from your own turn. The gate needs the Antigravity runs of the changed agents and skills that `AGENTS.md` § Release requires; the ones in this session's behavioural log cover five agents on agy.
3. **Operator, after installing** (`sh-agent update` then `hash -r`): the three smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI, plus one `sh-agent search --refs ADR-040` in this checkout. The first build of the new version fetches `modernc.org/sqlite` (`deploy/local.md`).
