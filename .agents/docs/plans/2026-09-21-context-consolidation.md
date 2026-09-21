---
owner: main
date: 2026-09-21
---
# Plan — Context consolidation: under Antigravity's block cap, pinned, canaried, and no regrowth   (2026-09-21)

## Why

BUG-085 (`.agents/docs/bugs/BUG-085-antigravity-global-context-truncated-24000-bytes.md`): the Antigravity CLI builds one `<user_rules>` block from `~/.gemini/config/rules/*.md` and cuts it at 24,000 bytes; the deployed context is 27,114 bytes, so the stack bullets and the inventory have been missing from every Antigravity session since v1.2.0. The cap is on the whole block (a 45-byte second file was dropped and cost the first 342 bytes; probe of 2026-09-21), so only pruning fixes it. The context has grown to 3.16× the size ADR-009 cut it to, and the retro loop that grew it has no delete step. The operator, 2026-09-21: "Go with your recommendations" — rows 1–8 of the findings table and the five scope gaps in `.agents/docs/plans/2026-09-21-context-audit-assessment.md` (the assessment) and the chat of that date; the audit method comes from `.agents/docs/plans/2026-09-21-context-audit-prompt.md` (the prompt).

## Scope

In: a canary line the installer emits as the LAST line of the rendered context; a test that renders the real repo and fails above 22,000 bytes; the pruning of `context/AGENTS.md` and `context/dev-setup.md` under that budget using the prompt's Step 4 classes, every deletion approved by the operator by number; the count of discrete instructions (the prompt's A3); a delete/merge action, a `rules/` row and a size step in the retro skill's target table; a sentence in the `bugs` agent's rules so a graduation never proposes a home the size test refuses; the canary question in the real-home smoke check on both CLIs; the approved bullet in `REQUIREMENTS.md` § Known limits with its ADR; goldfish G1–G3 on the three always-on files; e2e from the real home on both CLIs; one MINOR release.

Out: any `rules/` file (deferred: the block cap makes a split worthless until the context has room); skill descriptions (deferred); hooks or any enforcement (R3, ADR-007); Cursor/Codex rows; the prompt's Step 7 experiment; goldfish over skills; a runtime warning in `sh-agent status` when the rendered context exceeds the cap.

Follow-ups noted, not built: whether a second rules file loads when the block is under the cap (needs a smaller context; test once Phase 1 lands, before any `rules/` file); a `status` line for an over-cap context.

## Surfaces touched

Read this session; the assessment carries the measurements, no recon dispatched (D6).

- `internal/harness/render.go:21-54` — `RenderContext` writes Marker, `context/AGENTS.md`, other `context/*.md` in name order, then `RenderInventoryCompact`; the canary goes after the inventory.
- `internal/harness/harness_test.go:165` `TestInventoryAndRender`, `internal/harness/source_test.go:43` `TestPlanFromFSSource` — assert rendered content; gain the canary-last assertion.
- New `internal/harness/size_test.go` — renders the real checkout (`DirSource(filepath.Join("..", ".."))`, `LoadInventory`, `RenderContext`; the pattern of `cmd/sh-agent/main_test.go:1346` `TestHarnessMap`, which reads `../../README.md`).
- `context/AGENTS.md` (74 lines, 23,324 bytes), `context/dev-setup.md` (13 lines, 2,947 bytes) — pruned.
- `skills/retro/SKILL.md:52-59` target table and Step 2 — rows and a size step.
- `agents/bugs.md` § Rules (line 76 onward, "If three records share a prevention rule…") — one sentence.
- `REQUIREMENTS.md` § Known limits — one bullet (approved text in Interfaces). `.agents/docs/architecture/DECISIONS.md` — the ADR.
- `.agents/docs/testing/2026-09-03-verify-real-home.md:17` context-questions row, `.agents/docs/testing/guide.md` § 5 — the canary question. `README.md:208` context row of the deployment table — says the rendered file ends with the canary. This repo's `AGENTS.md` § Testing — the smoke-check sentence names the canary.
- `.agents/docs/plans/2026-09-21-context-pruning-review.md` — Phase 1's numbered classification (new). `.agents/docs/reviews/` — goldfish reports (Phase 3).
- `ROADMAP.md` Now → this plan; `TODO.md` stays empty.

## Interfaces

- `render.go`: `const Canary = "Canary: this context ends with the word GANNET."`; `RenderContext` appends `"\n" + Canary + "\n"` after `RenderInventoryCompact(inv)`. The word is fixed and quoted by the smoke check; it is not a rule and is never answerable from a CLI's skill or agent listing.
- `size_test.go`: `const antigravityRulesBlockBytes = 24000` (observed, BUG-085) and `const renderedContextBudgetBytes = 22000` (D2); `TestRenderedContextFitsAntigravityBlock` renders the real checkout and fails with `rendered context is N bytes, budget 22000 (Antigravity cuts the rules block at 24000): cut M bytes` when `len(out) > renderedContextBudgetBytes`. It runs in the normal `go test ./...`; from a checkout only (`DirSource` on `../..`); when `context/AGENTS.md` is missing there the test fails, never skips.
- Pruning review file (`plans/2026-09-21-context-pruning-review.md`): frontmatter; a numbered table `# | file:line | class | proposed action | bytes` over every paragraph and bullet of the two files, class ∈ {KEEP, DERIVABLE, VAGUE, COMMON-KNOWLEDGE, MOVE-TO-SCOPED, MOVE-TO-SKILL, MOVE-TO-ENFORCEMENT (always "n/a: R3"), DUPLICATE, CONFLICT, STALE} — a DUPLICATE or MOVE-TO-SKILL row names the skill or file that holds or will hold the text; counts per class; the A3 count of discrete instructions (one checkable directive = one) per paragraph and in total; projected bytes if only KEEP remained and after the proposed actions. The operator answers go or no per number, or go all; only then are actions applied.
- `skills/retro/SKILL.md` target table, two rows added: `| Rule that only retells a fixed bug, or duplicates a skill | delete from AGENTS.md §, or merge into the project skill | delete from context/AGENTS.md; merge into skills/<name>/SKILL.md |` and `| Rule scoped to a file type or path | .agents/rules/<name>.md | rules/<name>.md (shares Antigravity's 24,000-byte block with the context — REQUIREMENTS § Known limits) |`; Step 2 gains item 5: "Before a harness proposal adds to `context/AGENTS.md` or `context/dev-setup.md`, run `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1` with the addition in place; over budget, the proposal names the deletion or merge that pays for it, or targets a skill."
- `agents/bugs.md` § Rules, one bullet: "A recommended permanent home is never `context/AGENTS.md` or `context/dev-setup.md` when `TestRenderedContextFitsAntigravityBlock` would fail with the sentence added; recommend the skill that owns the topic instead, and say so in the record under Related."
- `REQUIREMENTS.md` § Known limits, first bullet, verbatim (approved 2026-09-21): "**Antigravity global rules are one 24,000-byte block.** CLI 1.2.7 builds one block from `~/.gemini/config/rules/*.md` and cuts it at 24,000 bytes with a `<truncated N bytes>` marker; a second file costs the first one room and is dropped without a marker (BUG-085, measured 2026-09-21; the vendor page says 12,000 characters per file). The deployed context is pinned under 22,000 bytes by a test, and `rules/` files share the block."
- Smoke check (`2026-09-03-verify-real-home.md` context row): the question list gains "the canary word (expect GANNET, the last line)"; a FAIL is any answer but GANNET, and on Antigravity the newest `~/.gemini/antigravity-cli/conversations/*.db` must hold no `<truncated` line inside the `<RULE[...AGENTS.md]>` block (`strings <db> | grep -c '^<truncated'` → 0).

## Failure modes

- Size test red: the run prints the rendered bytes, the budget and the bytes to cut; the change does not land. Nothing changes at install time — a binary-only place installs what its tag rendered.
- Canary absent from a CLI's answer: the file is truncated or not loaded; a bug record, never a passed check.
- A pruned rule that a session needed: the Phase 3 goldfish names it as a G1 or G3 gap and it comes back in the shortest form, counted against the budget.
- Approved deletions fall short of the budget: the size test stays red; a second numbered table, never a raised budget (Open row).

## Verification

Tiers: unit (`gofmt -l cmd internal && go vet ./... && go test -race ./...`, the size test, the render test); integration (`go run ./cmd/sh-agent install --home .agents/scratch/x` and `status --home .agents/scratch/x` drift none; the sandbox start of each CLI, guide § 4, log grepped for parse errors); behavioural (`agents/bugs.md`: three `claude -p --agents` runs in a scratch clone, tool order from the transcript, reported n of 3 — this repo's `AGENTS.md` § Testing); e2e (real home, both CLIs, after the box's install has moved to the pushed commit); cloud smoke n/a.

E2E criteria:
- Golden path: from the real home, `agy -p` and `claude -p` asked with no tools "what word does the global context end with?" both answer GANNET, and the Antigravity conversation store for that run has no `<truncated` line in the `<RULE[...AGENTS.md]>` block.
- Edge: with `context/AGENTS.md` padded by 2,100 bytes in a scratch copy of the repo, `TestRenderedContextFitsAntigravityBlock` fails naming the excess (a must-fail mutation: target unique, restored from a saved copy, `git diff` clean after, the test named and the code compiled).

## Phases

- [ ] Cold-reader check — design ready and implementation ready (<reviews/ paths>), or the user's same-turn override quoted with its date

### Phase 0 — Record, canary, docs
- [ ] `REQUIREMENTS.md` § Known limits bullet (Interfaces text) and its ADR by `changelog`, read against the diff, one commit — verify: `git diff --stat` shows the two files, AI — expect: bullet first under Known limits; ADR names BUG-085 and the 22,000 budget — docs: `REQUIREMENTS.md`
- [ ] `Canary` in `render.go`, appended last; `TestInventoryAndRender` asserts the rendered context's last line is `Canary` — verify: unit, `go test -race ./internal/harness -run 'TestInventoryAndRender|TestPlanFromFSSource' -count=1`, AI — expect: red before the render change (assertion fails on the inventory line), green after
- [ ] Smoke-check docs: the canary question in `2026-09-03-verify-real-home.md` context row and guide § 5; README context row; this repo's `AGENTS.md` § Testing — verify: `grep -n GANNET` over the four files, AI — expect: four hits — docs: `README.md`, `.agents/docs/testing/2026-09-03-verify-real-home.md`
- [ ] Checkpoint — expect: `go test -race ./...` exit 0; scratch-home install then `status` drift none; `tail -1 .agents/scratch/x/.claude/AGENTS.md` is the canary — evidence: command and exit status, <date>

### Phase 1 — Pin and prune
- [ ] `size_test.go` with `TestRenderedContextFitsAntigravityBlock`, run bare first — verify: unit, `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1`, AI — expect: FAIL naming about 27,160 bytes and 22,000
- [ ] Pruning review file: every paragraph and bullet of `context/AGENTS.md` and `context/dev-setup.md` classified, A3 count, per-class counts, projected sizes — verify: the file exists with one row per paragraph (row count = paragraph count by `grep -c`), AI — expect: projected size after actions ≤ 22,000
- [ ] Operator approval by number — who: user — expect: go or no per row, or go all
- [ ] Apply approved actions: deletions; merges land in the named skill in the same commit; every restatement of a moved claim found by grep before review (BUG-084 rule); comments and docs that describe moved text updated — verify: unit, the size test bare, AI — expect: PASS printing the rendered bytes
- [ ] Integration: scratch-home install and status; sandbox start of each CLI (guide § 4) — verify: integration, AI — expect: drift none; zero parse errors in both logs
- [ ] Review row: `reviewer` over the phase range (diff only) and `/review --docs` over the same range (the context is an interface) — who: AI — expect: findings fixed or rebutted in `## Outcome`
- [ ] Checkpoint — expect: `go test -race ./...` exit 0 with the size test green — evidence: `verifier` log `testing/<date>-context-consolidation-phase-1.md`, <date>

### Phase 2 — Stop the regrowth
- [ ] `skills/retro/SKILL.md`: two table rows and Step 2 item 5 (Interfaces) — verify: integration, scratch-home install and status, sandbox CLI start, AI — expect: drift none, skill listed, no parse error
- [ ] `agents/bugs.md` § Rules: the one bullet (Interfaces) — verify: behavioural, three `claude -p --agents` runs in a scratch clone where the dispatch describes a third-repeat prevention rule whose obvious home is the context; the record's Related must name the size test, which only the new text asks for; tool order read from the transcript, AI — expect: reported as n of 3, with n stated whatever it is
- [ ] Checkpoint — expect: `go test -race ./...` exit 0; install and status drift none — evidence: command and exit status, the behavioural tally, <date>

### Phase 3 — Sufficiency (goldfish)
- [ ] Three `reviewer` dispatches, each given ONE file (`context/AGENTS.md`, `context/dev-setup.md`, this repo's `AGENTS.md`) and the prompt's G1, G2, G3 questions verbatim, the G3 task being "release v1.4.0 as this repo did it" for the project file and "close a bug record after a fix" for the global one — verify: reports under `reviews/`, AI — expect: per file, G1 errors, G2 findings kept, G3 questions, and the two verdicts SUFFICIENT and MINIMAL
- [ ] Fixes for material G1 errors and blocking G3 questions, in the file a goldfish would look, each re-checked against the size budget — verify: the size test bare, AI — expect: PASS
- [ ] Checkpoint — expect: three reports on disk, `go test -race ./...` exit 0 — evidence: report paths, command and exit status, <date>

### Phase 4 — E2E and release
- [ ] This box's real home at the pushed commit (the hourly timer or `sh-agent update`) — who: user — expect: `sh-agent status` installed from the head commit, drift none
- [ ] Real-home smoke on both CLIs from a fresh clone of the demo service: the three checks plus the canary; on Antigravity, the conversation-store grep — verify: e2e, AI — expect: GANNET on both; `grep -c '^<truncated'` 0
- [ ] Must-fail mutation of the size test (E2E edge) — verify: unit, AI — expect: the padded copy fails naming the excess; `git diff` clean after restore
- [ ] Mac smoke (Antigravity-only place) after its `sh-agent update` — who: user — expect: GANNET
- [ ] BUG-085 status fixed in the record and the index (`bugs`), `changes/` entry, `documentation` for README and the guide — verify: `git diff --stat`, AI — expect: the files named
- [ ] Checkpoint — expect: everything above ticked, `go test -race ./...` exit 0 — evidence: `verifier` log `testing/<date>-context-consolidation-phase-4.md`, <date>
- [ ] Release `/preflight v1.5.0` — who: user — expect: signed tag, binaries published, `sh-agent version` prints v1.5.0 on a binary-only place

## Risks

Assumptions: the Mac's Antigravity CLI has the same 24,000-byte block — checked by the Mac smoke, Phase 4 (operator). About 5,200 bytes of DUPLICATE, DERIVABLE, COMMON-KNOWLEDGE and incident narration exist in the two files (27,114 + the canary − 22,000) — checked by the Phase 1 review's projected size, first task of that phase; if it falls short, the Open row. The canary is not in any skill or agent description — checked by `grep -r GANNET skills agents` returning nothing, Phase 0.

## Decisions

D1 — 2026-09-21 — Scope: findings rows 1–8 and the four open decisions of the assessment; rows 9 record, 10–11 defer, 12–14 strike — by: "Go with your recommendations"
D2 — 2026-09-21 — Budget 22,000 bytes for the rendered context, cap 24,000 observed — by: "Go with your recommendations"
D3 — 2026-09-21 — The canary is emitted last by `render.go`, after the inventory, so it proves the true end of the file — by: "Go with your recommendations" (gap 1)
D4 — 2026-09-21 — Pruning is approved by number from one table; nothing deleted before — by: "Go with your recommendations" (gap 3)
D5 — 2026-09-21 — `agents/bugs.md` gets three behavioural runs; the retro skill gets install, status and a CLI start — by: "Go with your recommendations" (gap 4)
D6 — 2026-09-21 — No explorer recon; the assessment is the brief — by: "Go with your recommendations" (gap 5)
D7 — 2026-09-21 — One tag at the end, v1.5.0 MINOR, from the operator's own turn — by: "Go with your recommendations"
D8 — 2026-09-21 — Phase 0 carries the canary and Phase 1 the red size test, so no checkpoint commits a red suite — by: implementer
D9 — 2026-09-21 — The size test lives in `internal/harness` and reads the checkout at `../..`, like `TestHarnessMap`; a release binary renders the same bytes (`TestEmbeddedPlanEqualsCheckout`) — by: implementer
Rejected: splitting into `rules/` files — the cap is per block (BUG-085 probe). Pinning 12,000 chars — `context/AGENTS.md` alone is 23,194 chars. Dropping the inventory — asserted by tests and read by the smoke check. A size check at install time — a binary-only place cannot act on it.
Open: the approved deletions leave the file over 22,000 — default: a second numbered table of further cuts; the budget is never raised by the implementer.
