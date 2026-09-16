---
owner: reviewer
date: 2026-09-16
---
# Review — bug-g2-deploy-template (2026-09-16)

Range reviewed: working tree, two named paths (uncommitted diff of exactly these files)
- `skills/init-project/SKILL.md` blob (working tree): `569a554da73cc8fe1db5f1acf789be1b5908ba87`
- `skills/init-project/assets/agents-tree/docs/deploy/local.md` blob (working tree): `52ca0d24563d74173da3e05884d70e0c2b9e6d84`

## Verdict: APPROVE WITH NITS

Nothing breaks Markdown rendering and the SKILL.md text and template are mutually consistent on the two intentional changes (Build/Run/Test blocks filled verbatim from setup answers; Verify left comment-only, matching `remote.md`/`gcp.md`, until a real endpoint exists — a genuine improvement over the old hardcoded `PORT`/`8080`/`healthz` example, which no other file in the repo referenced so nothing else went stale). The one real risk is the Configure table's single row, which folds an instruction ("one row per setting the code reads; none yet when the code does not exist") into the same cell as the `<VARIABLE>`/`<default>`/`<what it sets>` placeholders, with no clear directive to delete the row outright when there is no config yet — unlike `remote.md`/`gcp.md`, which keep instructions in an HTML comment above an empty table and never inside a cell.

## Findings

- **[medium / confidence: medium]** `skills/init-project/assets/agents-tree/docs/deploy/local.md:27` — the Configure table's only row is `| `<VARIABLE>` | `<default>` | `<what it sets>` — one row per setting the code reads; none yet when the code does not exist |`. The instructional clause ("none yet when the code does not exist") lives inside the third cell next to the placeholder itself, rather than in a comment above the table the way `remote.md` and `gcp.md` do it for their empty tables. An agent following SKILL.md's own new instruction ("a `<placeholder>` the answers do not fill … is deleted or left as a placeholder, never replaced with an invented value") could reasonably read this as "leave the row as-is" rather than "delete the row when there is no config" — the row's own text never says "delete this row." Result: a freshly scaffolded project's committed `local.md` could carry a literal `<VARIABLE>` / `<default>` row with a trailing meta-instruction sentence baked into the Purpose column, which is exactly the "placeholder an agent could still mistake for a real value" failure mode in scope here. Fix: move the "one row per setting … none yet" guidance into the existing HTML comment above the table (as `remote.md`/`gcp.md` do for their config-shaped tables) and leave the row itself as a plain `<VARIABLE>` / `<default>` / `<what it sets>` placeholder, or state explicitly that the row is deleted (not filled) when no config exists yet.

## Criteria

1. Dead code — n/a (docs-only diff, no code).
2. Flow tracing — SKILL.md's Step 4 prose and the template's Build/Run/Test/Verify sections agree with each other end to end; no other file in the repo references the old `PORT`/`8080`/`healthz` example or `local.md`'s internal structure in a way this diff breaks (checked `README.md`, `context/`, `agents/verifier.md`, `agents/documentation.md`, this repo's own filled `.agents/docs/deploy/local.md`).
3. Standards conformance — matches the global describing-text rule (SKILL.md text was updated alongside the template it describes); one instance (the Configure row) mixes instruction and placeholder in a way the sibling templates avoid — see finding above.
4. Directory layout — unchanged, files stay where they were.
5. Simplicity — no new dependency or abstraction; the change is a text/placeholder tightening.
6. Decoupling — n/a, single-owner doc/skill pair.

## Not examined

Did not review any other uncommitted file in the working tree (scoped out per instructions). Did not check git history/blame for prior versions of these two files beyond the diff shown. Did not run `/init-project` end to end against a scratch home to observe actual agent behavior on the Configure row — the finding above is a textual-ambiguity risk assessment, not an observed failure.

## Outcome

- Finding 1 (medium) — fixed: the instruction moved into an HTML comment above the Configure table ("No code yet, or no settings: delete the placeholder row and keep the empty table"), the row is now the three placeholders only, and `skills/init-project/SKILL.md` Step 4 says the Configure row is deleted when the code reads no setting yet.
- Round 1 closed 2026-09-16; a live `/init-project` run on both CLIs at the fixed tree is the verification owed (plan Phase 3, sandbox rows after the Claude lane).
