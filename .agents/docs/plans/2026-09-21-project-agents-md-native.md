---
owner: main
date: 2026-09-21
---
# Plan — Projects use `AGENTS.md` alone: no project `CLAUDE.md` stub   (2026-09-21)

## Why

The operator, 2026-09-21: "I saw that claude can use AGENTS.md from version 2.1.277. That would be greats so we use the full AGENTS.md for everything. Go and verify that and If it is true fix it and cut a new minor", then "Your recommendation. Approved." Verified today on this box (Claude Code 2.1.278): a probe project with only `AGENTS.md` (no `CLAUDE.md` in or above it) was read natively on the second session ("The probe word is HERON-4412, from …/AGENTS.md"); the first session after the upgrade read nothing, as the memory page documents. The user level is different: with the `~/.claude/CLAUDE.md` stub removed from a sandbox home, the global context did not load (UNKNOWN); the docs and the 2.1.277 changelog entry cover project `AGENTS.md` only. So the project stub goes and the global stub stays.

## Scope

In: `/init-project` no longer writes a project `CLAUDE.md`; its asset is removed; its description, Step 2 loop and bullet say so and name when a stub is still needed by hand (Claude Code before 2.1.277; Bedrock, Vertex, Foundry or a telemetry-disabled session, which fetch no feature flags; the first session after an upgrade reads nothing either way). An update run leaves an existing stub alone (D1). This repo's own `CLAUDE.md` is deleted (D2). `context/AGENTS.md` § Project layout says a project root is `AGENTS.md` alone. ADR-036 supersedes the project half of ADR-005. Release v1.6.0, MINOR (init-project behaviour changes for every user).

Out: the global `~/.claude/CLAUDE.md` stub written by `sh-agent install` (`plan.go:89`, `render.go:17`), verified required; removing stubs from existing projects; Antigravity (reads `AGENTS.md` natively already, ADR-005).

Follow-ups noted, not built: none.

## Surfaces touched

- `skills/init-project/SKILL.md:3` (description names "CLAUDE.md import"), `:22` (update rule mentions the `CLAUDE.md` fold), `:42-44` (the template loop lists `CLAUDE.md`), `:48` (the `CLAUDE.md` bullet); `skills/init-project/assets/CLAUDE.md` (deleted).
- `CLAUDE.md` at this repo's root (deleted); `AGENTS.md` § Layout / Key files if they name it (grep).
- `context/AGENTS.md:65` — "Root: `AGENTS.md` (`CLAUDE.md` is the line `@AGENTS.md`; no `GEMINI.md`)".
- `REQUIREMENTS.md` R2 — the sentence "(`~/.claude/CLAUDE.md` and a project's `CLAUDE.md` are one-line imports)" is the claim this plan reverses; replaced by the operator-approved text in Interfaces (D4), in the same commit as ADR-036 (the immutability rule; BUG-025 was this exact violation).
- `.agents/docs/testing/skill-matrix.md:22` — the `init-deploy` row's prompt says "including the .claude symlinks, CLAUDE.md, .gitignore lines"; `CLAUDE.md` is dropped from it and its check gains `! test -e CLAUDE.md`.
- `.agents/docs/testing/2026-09-21-verify-agents-md-native.md` (new) — the probes of today that back § Why, written before the code lands.
- `.agents/docs/architecture/DECISIONS.md` — ADR-036. README: no project-stub statement (grepped: none). `TestHarnessMap`, `plan.go`: untouched (the global stub is not in scope).

## Interfaces

- `skills/init-project/SKILL.md`: description drops "CLAUDE.md import"; the Step 3 template loop becomes `for f in AGENTS.md REQUIREMENTS.md TODO.md ROADMAP.md`; the `CLAUDE.md` bullet becomes: "`CLAUDE.md` — not written: Claude Code 2.1.277+ reads `AGENTS.md` natively when no `CLAUDE.md` sits in or above the project. An existing `CLAUDE.md` is left exactly as it is (an update never removes it; one that holds content other than `@AGENTS.md` is offered the fold as before). Write the one-line stub `@AGENTS.md` by hand only where native reading is unavailable: Claude Code before 2.1.277, Bedrock, Vertex, Foundry or a telemetry-disabled session; note that the first session after a Claude Code upgrade reads no `AGENTS.md` at all." The update sentence at `:22` keeps the fold exception.
- `context/AGENTS.md` Project layout line: "Root: `AGENTS.md` alone (no `CLAUDE.md`, no `GEMINI.md`; a `CLAUDE.md` stub `@AGENTS.md` only where Claude Code cannot read `AGENTS.md`: before 2.1.277, Bedrock, telemetry off), …" — about 100 bytes longer than today's parenthetical; the size test stays green (20,205 of 22,000 today).
- `REQUIREMENTS.md` R2, verbatim (approved 2026-09-21): "The context follows the agents.md standard and is named `AGENTS.md` at every level. A project has no `CLAUDE.md`: Claude Code 2.1.277+ reads `AGENTS.md` natively, and a one-line `@AGENTS.md` stub is written by hand only where it cannot (before 2.1.277, Bedrock, Vertex, Foundry, or telemetry disabled). `~/.claude/CLAUDE.md` stays a one-line import: Claude reads a user-level `AGENTS.md` only as an ancestor file, and a project with its own `CLAUDE.md`, or one outside the home directory, would then lose the global context. Skills follow the Agent Skills `SKILL.md` standard so that anything from skills.sh or any other source drops into `skills/` unchanged."
- ADR-036: supersedes ADR-005's sentence "a project's `CLAUDE.md` [is] the single line `@AGENTS.md`"; keeps the user-level stub with the verification that it is required; names the conditions above.

## Failure modes

- A project on an older Claude, Bedrock, or telemetry-off reads no instructions: the init-project bullet and ADR name the by-hand stub; nothing in the harness detects it (no hooks, R3).
- A place pulls this repo with a Claude older than 2.1.277: this repo's instructions stop loading there. The devbox is the only other Claude place; the handoff asks the operator to check `claude --version` before pulling.

## Verification

Tiers: unit (`gofmt -l cmd internal && go vet ./... && go test -race ./...`, the size test); integration (scratch install + status; a sandbox start of Claude Code with a question quoting the new Project layout line); behavioural, one run: `/init-project` headless in a fresh scratch repo (`.agents/scratch/init-probe`, `--add-dir <home>/.claude/skills/init-project`, answers in the prompt as the matrix's `init-deploy` row) and `ls CLAUDE.md` afterwards must fail; e2e: from this checkout with its `CLAUDE.md` deleted, `claude -p` with no tools answers a question only this repo's `AGENTS.md` can (the first heading "AGENTS.md — styrman-harness" and the canary GANNET from the global context) — run twice if the first is a first-session-after-upgrade case.

E2E criteria: (1) golden: in this checkout without `CLAUDE.md`, Claude Code names this repo's AGENTS.md first heading; (2) edge: `/init-project` in a fresh repo leaves no `CLAUDE.md`.

## Phases

- [x] Cold-reader check — design ready and implementation ready (`reviews/2026-09-21-coldread-project-agents-md-native-r1.md`, `-r2.md`, `-r3.md`: round 3 both tokens, 2026-09-21)

### Phase 1 — The change
- [x] `skills/init-project/SKILL.md` edits (Interfaces) and `git rm skills/init-project/assets/CLAUDE.md` (2026-09-21: grep shows the description, the update sentence and the new bullet) — verify: `grep -n "CLAUDE.md" skills/init-project/SKILL.md` shows only the new bullet and the update-fold sentence, AI — expect: two hits — docs: `skills/init-project/SKILL.md`
- [x] `git rm CLAUDE.md` at the repo root (2026-09-21; no stale claim in `AGENTS.md`/README); grep `AGENTS.md`/README for a claim that this repo has one — verify: `git ls-files CLAUDE.md` empty, AI — expect: empty; no stale claim
- [x] `context/AGENTS.md` Project layout line (Interfaces) (2026-09-21: size test PASS, 20,328 of 22,000) — verify: size test bare, AI — expect: PASS against the 22,000-byte budget (today 20,205; the new line is a few bytes longer or shorter)
- [x] Testing log `testing/2026-09-21-verify-agents-md-native.md` (written 2026-09-21): the four probes (project scope with/without stub; user scope in a home outside `/home`: no stub → both files load; project `CLAUDE.md` → global context dropped; stub → loads) — verify: file exists, AI — expect: written first
- [x] `REQUIREMENTS.md` R2 replaced by the main agent with the operator-approved Interfaces text (the file's "agents never edit" rule bars subagents and unapproved edits; the operator approved this exact text, D4) and ADR-036 by `changelog` (its lane: `DECISIONS.md` only), read against the diff, one commit — verify: `git diff --stat`, AI — expect: R2 one paragraph replaced; DECISIONS.md appended only — docs: `REQUIREMENTS.md`
- [x] `skill-matrix.md` `init-deploy` row (2026-09-21): drop `CLAUDE.md` from the prompt, add `! test -e CLAUDE.md` to the check — verify: `grep -c CLAUDE.md` on that row, AI — expect: 1 (the check only)
- [x] Integration and behavioural runs (Verification) — verify: AI — expect: scratch install drift none; Claude Code quotes the new line; the init probe has no `CLAUDE.md`; this checkout answers from its `AGENTS.md` (2026-09-21: sandbox home reinstalled, drift none, 257 files; both CLIs quoted the new Root line and GANNET, Antigravity store 0 truncation markers, 0 parse errors; headless `/init-project` in `.agents/scratch/init-probe`: scaffold written, `CLAUDE.md` absent — that run did not make the init commit the matrix row expects, one run, out of this plan's scope; e2e golden: this checkout answered `# AGENTS.md — styrman-harness` from its `AGENTS.md` twice)
- [x] Review row: diff under 50 lines read by the main agent; `/review --docs` over the range (an init-project behaviour changed) — expect: findings fixed or rebutted (2026-09-21: `reviews/2026-09-21-review-v1.5.0-HEAD-docs.md` APPROVE WITH NITS, one punctuation nit fixed)
- [x] Checkpoint — expect: `go test -race ./...` exit 0 — evidence: `verifier` log `testing/2026-09-21-project-agents-md-native.md` PASS (251 tests, size 20,328, install drift none, asset gone, global stub kept, root CLAUDE.md absent), 2026-09-21; the punctuation fix after it re-checked bare by the main agent

### Phase 2 — Release
- [x] `/preflight v1.6.0` — who: user's word "cut a new minor" (2026-09-21) — expect: SHIP, signed tag, five assets (gate `reviews/2026-09-21-preflight-v1.6.0.md` SHIP; the tag and assets follow this release commit, 2026-09-21)

## Risks

Assumptions: the devbox's Claude Code is 2.1.277 or newer — checked by the operator (`claude --version`) before pulling this repo there; named in the handoff. The feature-flag fetch works on this box (it did: the second probe read the file).

## Decisions

D1 — 2026-09-21 — An update run never removes an existing project `CLAUDE.md`; the pre-existing offer to fold a `CLAUDE.md` that holds other content into `AGENTS.md` (on the user's yes, leaving the one-line stub) stays as it is — by: "Your recommendation. Approved."
D2 — 2026-09-21 — This repo's own `CLAUDE.md` is removed; the devbox checks its Claude version first — by: "Your recommendation. Approved."
D3 — 2026-09-21 — The global `~/.claude/CLAUDE.md` stub stays — by: "Approve your recommendations" (2026-09-21), after the measured probes of D5
D4 — 2026-09-21 — REQUIREMENTS R2 edit, exact text in Interfaces — by: "Approve your recommendations" (2026-09-21); ADR-036 and the edit land in one commit
D5 — 2026-09-21 — Measured in a probe home outside `/home` (Claude Code 2.1.278): with no stub, `~/.claude/AGENTS.md` and the project's `AGENTS.md` both load (as ancestor `.claude/AGENTS.md` plus project file); a project with its own `CLAUDE.md` switches Claude to CLAUDE.md-only mode and the global context is dropped; with the stub it loads in every case. A probe home under `/home/kesokaj` is invalid: the real `~/.claude/CLAUDE.md` counts as an ancestor `.claude/CLAUDE.md` there — by: implementer, measured (BUG-087's rule: the arrangement)
Rejected: an `--claude-stub` flag on `/init-project` (a by-hand one-liner is enough; a flag is code); removing stubs on update (breaks older or flag-less Claude sessions silently).
