---
owner: main
date: 2026-09-21
---
# Assessment — the "Agent context audit" prompt, and whether to create `rules/`   (2026-09-21)

Status: PROPOSAL, nothing adopted. No plan is active; ROADMAP "Now" stays empty until the operator decides. A first draft of this file was checked by three independent refuters (measurements, fit, rules lenses; workflow `wf_cc624239-e08`, 2026-09-21); every correction they made is applied below and listed at the end.

## Why

The operator, 2026-09-21, own words: "I will give you a prompt, take a look at it and say what you think. Could this help us or are we already optimized? Write down a plan file so we keep your thoughts even if we dont do this." and "Should we create rules? We dont have any."

The prompt is `agent-context-audit-prompt.md` (untracked at the repo root when this was written; Open question 1 says where it should live). It is a read-only audit of every instruction surface against a five-layer loading model (L0 enforcement, L1 always-on, L2 scoped rules, L3 skills, L4 references): defects first, a scorecard (A1–D2), a pruning review, a goldfish test, a target layout, and a proposed adherence experiment. Every statement is to be labelled MEASURED, DOCUMENTED or INFERRED. This file follows that convention.

## Short answer

We are not optimized, and the prompt's "defects first" step found a live defect within the hour:

**Antigravity CLI loads only the first 24,000 bytes of our global context.** The deployed file is 27,114 bytes. The last 3,114 bytes — most of the Environment & Stack bullets and the Installed inventory — are absent from every Antigravity session, on every place, in every release since v1.2.0 (2026-09-17). Nothing in the tests or smoke checks could see it. Severity medium: Antigravity-only, and nearly every cut bullet has a second home in a skill.

Beyond the defect, the always-on context has drifted to 3.2× the size ADR-009 cut it to on 2026-09-03, by paragraph density (67 → 74 lines, 7,349 → 23,194 chars), and the retro loop that grew it has no delete step. The prompt's method — Step 3 defects, Step 4 pruning classes with counts, Step 6 one-row-per-move layout, a goldfish on the always-on files — is the right tool for the consolidation Feature already queued. Its L0 layer, its Cursor/Codex rows, its goldfish over every skill and its adherence experiment do not apply here. `rules/` comes later, as an output of that pruning, after one probe that the split even helps.

## Measured today (2026-09-21, this box, `main` at dbc26dd)

Commands: `wc -l -m -c`, `head -c`, `tail -c +N`, `git show <rev>:context/AGENTS.md | wc -m -c`, `git tag --contains`, an awk pass over each `SKILL.md` frontmatter, `strings` over the Antigravity conversation store, `agy -p` and `claude -p` from the real home with a no-tools question. CLI versions by `agy --version` (1.2.7) and `claude --version` (2.1.278).

| Surface | Lines | Chars | Bytes | Load | Label |
|---|---|---|---|---|---|
| `context/AGENTS.md` (source) | 74 | 23,194 | 23,324 | always, both vendors, via the deployed file | MEASURED |
| `context/dev-setup.md` (source, appended at install) | 13 | 2,935 | 2,947 | always | MEASURED |
| Generated comment + inventory the installer adds (`render.go`) | 7 | — | 134 + 708 | always | MEASURED (deployed − sources = 843) |
| `~/.claude/AGENTS.md` (deployed) | 96 | 26,972 | 27,114 | always, Claude Code, whole file | MEASURED |
| `~/.gemini/config/rules/AGENTS.md` (identical bytes) | 96 | 26,972 | 27,114 | always, Antigravity, **first 24,000 bytes** | MEASURED (see Defect) |
| Project `AGENTS.md` (this repo) | 77 | 11,706 | 11,743 | always, in this checkout | MEASURED |
| `MEMORY.md` (auto memory) | 1 | 164 | 166 | always, Claude Code | MEASURED |
| Always-on total on Claude Code in this repo | 174 | ≈38,842 | ≈39,023 | ≈9.7k tokens at 4 bytes/token | INFERRED (bytes/4) |
| Longest lines in `context/AGENTS.md` | lines 57, 16, 45 | 1,884 / 1,767 / 1,679 | — | — | MEASURED |
| Directive words (never/always/must/only/every) in `context/AGENTS.md` | — | 98 | — | — | MEASURED; discrete-instruction count (the prompt's A3) NOT done |
| Skill descriptions: name+description total, 27 skills | — | 12,312 | — | every turn (each CLI injects its listing) | MEASURED |
| Skill descriptions over 250 chars / over 1,024 / over 1,536 | 25 / 0 / 0 | longest `cloud-run-deploy` 989 | — | — | MEASURED |
| Largest `SKILL.md` | `terraform` 393 lines | — | — | on demand | MEASURED |
| Agent descriptions | 391–551 chars each | — | — | every turn | MEASURED |
| `rules/` | empty (`.gitkeep`) | — | — | — | MEASURED |

Growth of `context/AGENTS.md` (43 commits): ADR-009's cut is 87a862d (2026-09-03, "Trim the context…"): 16,053 chars / 136 lines → 7,349 / 67. Then 9,124 (09-10) → 12,632 (09-14, 71 lines) → 17,806 (09-17) → 21,046 (09-17, 5bc56fb) → 23,194 (09-18, 74 lines). 23,194 / 7,349 = 3.16×. MEASURED. ADR-009 counted tokens too ("about 3.5k fewer per session") and set "under 70 lines"; no size has been checked since, by any unit. DOCUMENTED (`.agents/docs/architecture/DECISIONS.md` ADR-009; `.agents/docs/testing/2026-09-03-verify-context-trim.md`).

The growth engine: the retro skill's harness-target table (`skills/retro/SKILL.md` lines 52–59) routes every "rule gap or restatement" to `context/AGENTS.md` and every stack finding to `context/dev-setup.md`, with no delete or merge action and no `rules/` row; the `bugs` agent graduates recurring rules INTO the context. Nineteen Graduated rows in `bugs/INDEX.md` name `context/AGENTS.md` (BUG-029, 031, 044, 045, 057, 058, 061, 066, 068–072, 078–082, 084); BUG-083 went to this repo's `AGENTS.md` § Testing. DOCUMENTED.

## Defect — Antigravity CLI truncates the deployed context at 24,000 bytes

- **Observation (MEASURED, twice):** the CLI's conversation store (`~/.gemini/antigravity-cli/conversations/*.db`, read with `strings`) holds the prompt it built: the file wrapped in `<RULE[/home/kesokaj/.gemini/config/rules/AGENTS.md]>` inside `<user_rules>`, cut mid-sentence, then the CLI's own marker, then the closing tag. 2026-09-21 (`48ae34c2-*`, the probe for this file): `<truncated 3114 bytes>`, 27,114 − 3,114 = 24,000. 2026-09-18 22:20 (`3c034e5c-*`, a v1.3.0-era session): `<truncated 2553 bytes>`, 26,553 − 2,553 = 24,000. The marker counts bytes (the tail is 3,102 chars). Asked with no tools for its last heading, `agy -p` answered `## Environment & Stack` and reported the 3,114-byte truncation itself; it quoted the `date +%F` and `rm -r` rules correctly, so the file loads up to the cut. `claude -p` with the same question named `## Installed` and "Skills (27)" (DOCUMENTED: Claude Code loads a file up to 4 MiB in full).
- **What is cut (MEASURED, `tail -c +24001`):** the second half of the GCP-defaults bullet (after "Once given, record them in" — the ask-first half loads), then Backend Go, Frontend TanStack, Machines/Colima, root steps with `sudo -n`, test locally first, build amd64, Artifact Registry naming, Cloud Run first, gcloud auth and Secret Manager, and the whole `## Installed` inventory.
- **The cap:** 24,000 bytes, observed on two sessions three days apart (MEASURED). Whether it is per rules file or per `<user_rules>` block is UNKNOWN: only one rules file exists on this box, and one file cannot tell the two apart. The vendor page fetched today says "Rules files are limited to 12,000 characters each" (DOCUMENTED, https://antigravity.google/docs/rules-workflows/); the CLI's observed number is different, so the page's per-file wording cannot be transferred to the CLI.
- **Since when (MEASURED):** the crossing commit is 5bc56fb itself (2026-09-17 23:09): 134 + 21,158 + 1 + 2,947 + 708 = 24,948 bytes; the commit before it rendered 22,609. `git tag --contains 5bc56fb` → v1.2.0, v1.2.1, v1.2.2, v1.2.3, v1.3.0, v1.4.0. Every release binary since v1.2.0 installs a context Antigravity truncates; the Mac (Antigravity-only) and the devbox lose 2,553 bytes at v1.3.0 and 3,114 at v1.4.0.
- **Why nothing caught it (MEASURED):** no test pins a deployed file's size against a vendor cap (`TestHarnessMap` pins verb names). The four 2026-09-17 real-home smoke checks installed commits that predate 5bc56fb. The one Antigravity context check after it (d278f5a, 2026-09-18) asked for a Hard-rules sentence, which sits at bytes 13,875–22,609 and loads. The 2026-09-03 procedure's other questions — skill count, pruned names, the grilling command — DO sit in the truncated tail (bytes 26,406+), and the CLI answered them from its own skill listing. A tail question is therefore not enough; the question has to be one only the file can answer.
- **Severity: medium**, Antigravity-only. Every cut stack bullet except "never suggest another cloud provider", "remote boxes only on request", the Cloud Run-first line and the names-only inventory has a skill home (amd64/buildx: `docker`, `cloud-run-deploy`, `preflight`; Secret Manager: four skills; `gcp.md` defaults: five; slog and layout: `golang`; Colima: `docker`, `cloud-run-deploy`; `sudo -n`: `docker`) — but a skill loads on description match, a rule in the context loads always, so the loss is of certainty, not of the words. Antigravity use on this box since the crossing: two sessions, one of them the probe. Confidence: high on the cut and its position; medium on what "per file" means.
- **Recorded as BUG-085** (`.agents/docs/bugs/`), status open; the fix is Phase 1 below.

## Where the prompt fits us, and where it does not

Fits — each would change something here:

- **Step 3 "defects first: anything silently truncated".** Found the defect. The harness had the "silence from your own instrumentation" rule and still had no probe that reads the END of what it installs. Its lasting form here is small: a size test that renders the real repo (the `TestHarnessMap` pattern) and a canary question.
- **Step 4 pruning classes with counts** (KEEP / DERIVABLE / VAGUE / COMMON-KNOWLEDGE / MOVE-TO-SCOPED / MOVE-TO-SKILL / DUPLICATE / STALE). Nothing in `retro`, `bugs`, `review` or `preflight` measures or prunes the context; the loop only adds. A pass over the 74 paragraphs of `context/AGENTS.md` and the 13 bullets of `dev-setup.md`, with a count per class and the projected size if only KEEP remained, is the consolidation Feature's first deliverable. The prompt's A3 (discrete instructions, target ≤100) is the adherence lever its own R7 evidence supports; that count has not been done.
- **Step 6 target layout, one row per move** (content | current layer | proposed layer | reason | check fixed | size change). The operator decides row by row, which is how this operator wants to decide.
- **A goldfish on the three always-on files** (Step 5's G1–G3 on `context/AGENTS.md`, `context/dev-setup.md`, this repo's `AGENTS.md`), after the prune, to check that what remains is sufficient.
- **A check with a unit.** ADR-009 set a line ceiling and a token count and neither was ever re-checked. Bytes are what the vendor cap and the cost measure; instruction count is what adherence measures. Pin the first in a test; count the second in the pruning review.

Does not fit — strike or run report-only:

- **L0 / E1 / MOVE-TO-ENFORCEMENT.** REQUIREMENTS R3, ADR-001, ADR-007 and the non-goals: no hooks, gates, CI, or permission edits — researched and rejected 2026-09-03. E1 would list every MUST/NEVER as lacking a backstop; that is the design, on record. (Claude Code's memory page says a PreToolUse hook is the only way to block an action regardless of the model; the operator knows and chose otherwise.)
- **R4 Cursor and R5 Codex rows, and the 32 KiB and 8,000-char ceilings drawn from them.** Two vendors only. The tightest limit across OUR harnesses is Antigravity's 24,000 bytes observed (12,000 chars documented) and Claude Code's 200-line target. K3's 8,000-char listing ceiling is Codex's; Claude Code documents a 1,536-char cut per skill and no listing budget (DOCUMENTED, https://code.claude.com/docs/en/skills). The prompt's "~1% of the context window, least-used dropped first" is UNVERIFIABLE from that page.
- **A written loading model per vendor (Step 0 as a spec file).** README's "Adding things by hand" table already states the load mode per surface and the per-vendor scoping key. The one missing fact is the cap: one bullet in REQUIREMENTS § Known limits, not a spec.
- **Step 5 over all 27 skills.** Cold-reader checks exist for plans (feature skill Step 3) and handoffs; agent changes get behavioural N-of-N runs (this repo's `AGENTS.md` § Testing, BUG-083). Twenty-seven skills is the prompt's own >15-run territory, and most are vendored or stable.
- **Step 7 adherence experiment (canary, 3 tasks × 2 layouts × N sessions).** The behavioural-run recipe already exists, and the permanent tail probe is one smoke-check question, not an experiment. Early-versus-late adherence across two layouts is the measure→rule→stale loop the operator closed on 2026-09-18 ("no more retro after this"). Struck unless the operator reopens it.
- **D1 "flag AI-generated context files not human-verified".** Every file here is agent-written and operator-reviewed by trail; uninformative for this repo.

Reference limits checked today (the prompt's Step 1, our two vendors only):

- R1 Claude Code memory — CONFIRMED: target under 200 lines per file; a file up to 4 MiB loads in full; imports load at launch and save no context; `AGENTS.md` read natively when no `CLAUDE.md` is present, otherwise via import; `MEMORY.md` first 200 lines or 25 KB. `.claude/rules/` and `~/.claude/rules/`: "Rules without a `paths` field are loaded unconditionally… Path-scoped rules trigger when Claude reads files matching the pattern, not on every tool use"; `~/.claude/rules/` applies to every project; symlinked rules pointing outside the working directory load only after external-import approval, and then only the unscoped ones. (https://code.claude.com/docs/en/memory)
- R2 Claude Code skills — CONFIRMED 1,536-char cut and the 5,000/25,000-token compaction budget; no listing-budget percentage on the page. `paths` is also a skill frontmatter key. (https://code.claude.com/docs/en/skills)
- R3 Antigravity — CONFIRMED 12,000 chars per rules file (documented); global rules documented at `~/.gemini/GEMINI.md` or, for the CLI, `~/.gemini/antigravity-cli/rules/`; workspace rules `.agents/rules/`; four activation modes (Manual, Always On, Model Decision, Glob). The harness deploys to `~/.gemini/config/rules/`, which CLI 1.2.7 reads: its prompt names the file as `<RULE[/home/kesokaj/.gemini/config/rules/AGENTS.md]>` (MEASURED); no `~/.gemini/antigravity-cli/rules/` exists on this box; whether `config/rules` stays the CLI's path in later versions is UNVERIFIED. (https://antigravity.google/docs/rules-workflows/)
- R4, R5: n/a. R6 (Agent Skills spec), R7, R8: not fetched.

## The `rules/` question

`rules/<name>.md` is built, tested (`source_test.go`, `harness_test.go`, `pathguard_test.go`) and documented in README's deployment table: each file lands as `~/.claude/rules/<name>.md` (always-on, or scoped with `paths:`) and `~/.gemini/config/rules/<name>.md` (always-on, or `trigger: glob` + `globs:`). The directory has been empty since 2026-09-03. A rule file with no frontmatter is a valid always-on rule on both CLIs (the deployed `AGENTS.md` has none). MEASURED and DOCUMENTED. No rule file has ever been loaded live on either CLI from this harness; the smoke checks cover none.

What a rule file would and would not do here:

- **It does not fix the defect by itself.** Moving `context/dev-setup.md` to `rules/dev-setup.md` leaves the main file at 24,167 bytes, still over. The inventory is generated by `render.go`, asserted inside `.claude/AGENTS.md` by `TestPlanFromFSSource` and `TestInventoryAndRender`, and read by the smoke check's skill-count question: moving or dropping it is a code change, not a file drop. And with one file observed, the cap may be per `<user_rules>` block, in which case a split moves the cut to whichever file sorts last and fixes nothing. One probe settles it: a second small rule file in the real home with a canary line, `agy -p` asked for the canary and for `AGENTS.md`'s last heading.
- **Its content already has a home.** The tail's stack bullets are in the skills (grep above): their pruning class is DUPLICATE, not MOVE-TO-SCOPED. Deleting the duplicates alone brings the deployed file under 24,000 bytes with no rule file and no code change.
- **Path scoping is Claude Code's deterministic lever, with a caveat.** A `paths:` rule loads when a matching file is READ (DOCUMENTED above), never when the operator says a word. "Ship means /preflight" and "ask for the GCP defaults" must stay always-on or in a skill description. On Antigravity, Model Decision mode is description-triggered, but README exposes only glob for Antigravity and CLI 1.2.7's support for it is UNVERIFIED.
- **Rules do not change the growth engine.** The retro target table and the `bugs` graduation rule route every new rule into the context; a `rules/` row and a delete/merge action in that table are what change the next retro's output, and they belong in the same commit as the first rule file.

Recommendation: later. After the pruning review shows whether any MOVE-TO-SCOPED residue exists, with a smoke check per rule file on both CLIs before any content depends on it, and with the retro table changed in the same commit.

**Probe result (2026-09-21, after the refute pass, operator-approved):** a 45-byte second rules file beside `AGENTS.md` in the real home never appeared in the CLI's prompt (no `<RULE[...]>` block, canary UNKNOWN, string absent from the conversation store), and the AGENTS.md cut grew from `<truncated 3114 bytes>` to `<truncated 3456 bytes>`. MEASURED. The cap is on the whole rules block, not per file (INFERRED from the two numbers): a split buys nothing, pruning is the only fix, and `rules/` on Antigravity is limited to whatever room the context leaves. Recorded in BUG-085.

## Recommendation

One Feature, the consolidation queued since 2026-09-18, run with the prompt's method inside it. Phases, each a checkpoint:

0. **Probe and record.** BUG-085 record (done). REQUIREMENTS § Known limits: one bullet — Antigravity CLI loads 24,000 bytes of a rules file, documented 12,000 chars, per-file or per-block unknown (operator approval, since REQUIREMENTS is immutable without it). The two-file cap probe in the real home, its result written to the bug record and the bullet.
1. **Under the cap, pinned.** A test that renders the real repo's Antigravity rules file and fails above the observed cap minus a margin (Open question 2) — it fails today. Step 4 pruning classes with counts over `context/AGENTS.md` and `context/dev-setup.md`, written to `.agents/docs/plans/`; the operator approves each DUPLICATE, DERIVABLE, COMMON-KNOWLEDGE and incident-narration deletion by number. A canary sentence as the last line of the deployed context, and the real-home smoke check asks for it on both CLIs. Changelog, the BUG-085 status closed, a smoke run on both CLIs.
2. **Target layout.** Step 6 rows for what remains: KEEP stays; MOVE-TO-SKILL goes to the skill that already owns the topic; any MOVE-TO-SCOPED residue becomes the first `rules/` file(s), with the retro table row, the `bugs` graduation rule and a smoke check per rule file in the same commit. ADR for the layout (rule 3 of this repo's `AGENTS.md` does not fire — `rules/` is already on the map — but the structural decision deserves the record).
3. **Sufficiency.** Goldfish G1–G3 on the three always-on files, fresh `reviewer` dispatches given the file alone; the A3 instruction count; fixes go where a goldfish would look.
4. **Release.** `/preflight` a MINOR, since binary-only places get a changed context and the size test; the Mac and devbox update.

Struck from the prompt for this repo: L0/E1, R4/R5, Step 0 as a spec file, Step 5 over skills, Step 7, the scorecard as a deliverable, D1. Not recommended: running the prompt as written in one session.

## Open questions for the operator

1. Where the prompt file lives: `.agents/docs/plans/2026-09-21-context-audit-prompt.md` beside this file (recommended: committed, referenced), or `skills/context-audit/SKILL.md` with `disable-model-invocation: true` if it is to be re-run.
2. The margin for the size pin: 22,000 bytes (8% headroom, recommended) or the bare 24,000. Not 12,000 chars — `context/AGENTS.md` alone is 23,194 chars, so that pin turns the Bug into halving the file, which is the Feature.
3. Whether the consolidation Feature queued since 2026-09-18 is this Feature. Recommended: yes.
4. Whether Phase 0's probe runs before the Feature is scoped (ten minutes in the real home) or as its first task. Recommended: before, so the scope confirmation knows whether a split is even on the table.

## Refute pass (2026-09-21)

Three fresh readers, given this file's first draft, the prompt and the repo, told to refute. Verdicts: adopt-in-part ×3; rules later ×3. Corrections applied: "per rules file" was labelled MEASURED on one file's evidence; the crossing is 5bc56fb itself and every tag since v1.2.0, not "this box today"; the 2026-09-17 smoke checks predate the crossing and the procedure's tail questions are answered by the CLI's own listing; the split leaves 24,167 bytes and the inventory is code; a 12,000-char pin contradicts the fix; severity high → medium; ADR-009 ratio 3.5× → 3.16× and 67 → 74 lines; BUG-083 cited to the wrong file and BUG-029/031 omitted; rule 3 of this repo's `AGENTS.md` does not fire on a `context/` → `rules/` move; the Step 0 spec and Step 7 experiment were oversold; the "1,884-char paragraph" reason for rules named a paragraph the split would not touch; the Claude path-rule trigger was unlabelled (now DOCUMENTED: on read). Journal: `~/.claude/projects/-home-kesokaj-workspace-styrman-harness/22858f28-*/subagents/workflows/wf_cc624239-e08/journal.jsonl`.

## Sources

- Vendor pages fetched 2026-09-21 with WebFetch: Antigravity rules and workflows; Claude Code memory; Claude Code skills (URLs above).
- Repo: `DECISIONS.md` ADR-001, ADR-007, ADR-009; `REQUIREMENTS.md` R3, R5, § Known limits; `bugs/INDEX.md` Graduated table; `README.md` map and deployment table; `internal/harness/plan.go`, `render.go`, `inventory.go`; `skills/retro/SKILL.md` target table; `agents/bugs.md`; `.agents/docs/testing/2026-09-03-verify-real-home.md`, `sh-agent-matrix.md` § Real machines.
- Runs: `agy -p` and `claude -p` from the real home, 2026-09-21, no tools, four context questions; `strings` over `~/.gemini/antigravity-cli/conversations/{48ae34c2,3c034e5c}-*.db` (dump in `.agents/scratch/agy-conv-strings.txt`).
