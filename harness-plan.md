# Harness Improvement Plan (Phase 3)

Date: **2026-09-03**  
Target Systems: **Claude Code** & **Google Antigravity / Gemini CLI**  
Goal: Cheaper, faster, simpler, more reliable, and strictly conformant to open agent standards.

---

## 1. Ranked Backlog

### Priority 1: Fix `CheckFreshness` False-Positive Failure on Non-Source Commits (Free Win / Bug Fix)
- **Problem:** `internal/trail/check.go` evaluates `p.Scope()`. When HEAD committed today contains only non-source files (documentation, prompt files, markdown outside `sourceGlobs`), `p.Scope()` is `""`. `CheckFreshness` marks `covered = false` and prints:
  `"STALE: no changelog entry covers the current change set (HEAD %s, newest entry %s) — an entry must declare scope: "`
- **Fix:** In `internal/trail/check.go`, if `p.Scope() == ""` because there are no dirty source files and `len(p.HeadSourceToday()) == 0`, report:
  `OK: no source changes in current state (clean tree; HEAD has no source changes)`
  and do not fail.
- **Expected Impact:**
  - Quality / Reliability: Eliminates false STALE reports on `/preflight` and `harness check` when committing docs.
  - Latency: Eliminates wasted model cycles trying to satisfy an impossible empty scope.
- **Effort:** Small.
- **Measurement / Verification:** `go test -race ./internal/trail/...` + verify `go run ./cmd/harness check` outputs `RESULT: FRESH`.
- **Rollback:** Revert git commit.
- **Evidence Grade:** [A3] (R4 & R7 compliance).

---

### Priority 2: Standardize Agent Skills Directory Layouts & Prune Packaging Clutter (Free Win / Standards Conformance)
- **Problem:**
  - `shared/skills/codebase-design`: `DESIGN-IT-TWICE.md` and `DEEPENING.md` reside in the skill root folder rather than under `references/`.
  - `shared/skills/anti-ui-slop`: Has both `reference/` (singular) and `references/` (plural) folders.
  - 4 skills (`anti-ui-slop`, `codebase-design`, `grill-me`, `grilling`) contain `agents/openai.yaml`.
  - `anti-ui-slop` contains dead packaging files: `MANIFEST.json`, `CHECKSUMS.sha256`, `MODIFICATIONS.md`, `NOTICE`, `LICENSE`.
- **Fix:**
  - In `codebase-design`: Move `DESIGN-IT-TWICE.md` and `DEEPENING.md` into `references/`; update `SKILL.md` links.
  - In `anti-ui-slop`: Consolidate `reference/*.md` into `references/*.md`; update `SKILL.md` links.
  - Remove extraneous `agents/openai.yaml` from all 4 skills.
  - Remove extraneous packaging clutter (`MANIFEST.json`, `CHECKSUMS.sha256`, etc.) from `anti-ui-slop` (license attribution is already preserved in `VENDORED.md`).
- **Expected Impact:**
  - Portability: 100% conformance to the Agent Skills open specification (`agentskills.io`).
  - Cleaner repository surface, eliminating dead files mirrored to user configs.
- **Effort:** Small.
- **Measurement / Verification:** `harness doctor` passes with `--home /tmp/<scratch>`; `go test -race ./...` passes.
- **Rollback:** `git checkout shared/skills`.
- **Evidence Grade:** [A3] (Living spec active Aug 2026).

---

### Priority 3: Add Negative Triggers ("Do NOT use for...") to 13 Skills (Quality Win)
- **Problem:** 13 skills (`audit`, `bash`, `distill`, `find-skills`, `gcloud-auth-verification`, `gke-basics`, `grill-me`, `grilling`, `init-project`, `research`, `retro`, `skill-repair`, `terraform`) have positive trigger definitions but lack negative boundary definitions.
- **Fix:** Add concise `Do NOT use for: ...` clauses to each skill's `description:` in its `SKILL.md` frontmatter.
  - Example for `bash`: *"Do NOT use for executing application logic (use golang/frontend) or cloud infrastructure (use terraform)."*
  - Example for `research`: *"Do NOT use for targeted local file edits or quick inspections (use direct file tools)."*
  - Example for `terraform`: *"Do NOT use for direct GCP CLI actions (use gcloud-setup) or Docker management (use docker)."*
- **Expected Impact:**
  - Routing accuracy jumps from ~73% to ~85% (per first-party Anthropic/OpenAI C2 measurements).
  - Fewer false skill activations saves turn 0 context and unnecessary tool loading.
- **Effort:** Medium.
- **Measurement / Verification:** Automated trigger check script confirms all 28 skills have negative boundaries.
- **Rollback:** Revert git commit.
- **Evidence Grade:** [C2].

---

### Priority 4: Explicit Antigravity Model Mapping in `derive.go`
- **Problem:** `internal/install/derive.go` only maps `opus` -> `pro` and `haiku` -> `flash`. Agents with `model: sonnet` (`documentation`, `explorer`, `verifier`) omit the `model:` tag when generated for Antigravity, defaulting to whatever model happens to be active.
- **Fix:** In `internal/install/derive.go`, document and clarify model mappings, or map `sonnet` to explicit Antigravity model tiers (or keep `inherit` with documented rationale).
- **Expected Impact:** Predictable subagent cost and capability execution across vendors.
- **Effort:** Small.
- **Measurement / Verification:** Run `harness derive` and verify generated agent YAML in `antigravity/plugins/harness/agents/`.
- **Rollback:** Revert git commit.
- **Evidence Grade:** [A3].

---

### Priority 5: Cross-Vendor Portable Launcher for `grill-me`
- **Problem:** `shared/skills/grill-me/SKILL.md` contains `Call the Skill tool with "grilling"`. This works only on Claude Code (which has a `Skill` tool) and fails on Antigravity.
- **Fix:** Update `grill-me/SKILL.md` with cross-platform instructions:
  *"Engage the grilling interview by loading and following the instructions in the `grilling` skill (using the Skill tool on Claude Code or reading `shared/skills/grilling/SKILL.md` on Antigravity)."*
- **Expected Impact:** Portability across both target platforms.
- **Effort:** Small.
- **Measurement / Verification:** Inspection and cross-CLI execution.
- **Rollback:** Revert git commit.
- **Evidence Grade:** [B4].

---

## 2. The "Stop Doing" List

1. **Stop Failing Trail Freshness on Doc-Only Commits:**
   - *Assumption encoded:* Every commit implies a code changeset that owes a changelog entry with a source-digest scope.
   - *Why stale:* Committing documentation, prompt updates, or notes changes zero source lines. Demanding a changelog entry with `scope: ` (empty) halts release workflows for non-code changes.
2. **Stop Carrying Third-Party Package Metadata in Skills:**
   - *Assumption encoded:* Upstream package files (`agents/openai.yaml`, `MANIFEST.json`, `CHECKSUMS.sha256`) were copied verbatim during initial vendoring.
   - *Why stale:* The harness runs Claude Code and Antigravity, not OpenAI Assistants or NPM package distribution. These files waste space and leak dead config into `~/.claude` and `~/.gemini`.
3. **Stop Leaving Skill Negative Boundaries Implicit:**
   - *Assumption encoded:* Models can infer when NOT to use a skill from positive descriptions alone.
   - *Why stale:* Evidence across vendors shows routing accuracy suffers (~12% false-activation rate) without explicit negative boundaries.

---

## 3. Core vs Adapter Refactor List

1. **Core (`shared/skills/`):**
   - Normalize directory structure: all skill reference materials strictly under `references/`.
   - All skill triggers strictly follow `Use when / Do NOT use when`.
2. **Adapter (`internal/install/derive.go` & `antigravity/plugins/`):**
   - Regenerate derived Antigravity agents via `harness derive` whenever Claude source definitions change.
   - Maintain strict `toolMap` enforcement to prevent Antigravity subagent hangs.

---

## 4. Explicit Non-Changes (Things that look wrong but measured fine)

1. **`shared/AGENTS.md` Constitution Duplication to `claude/CLAUDE.md` and `antigravity/GEMINI.md`:**
   - *Why leave it:* While some tools support `@AGENTS.md` import, copying verbatim guarantees 100% offline determinism, eliminates symlink fragility across OSes, and is guarded by `harness doctor` against drift [R1].
2. **Go binary implementation vs native scripting:**
   - *Why leave it:* The Go binary (`harness`) executes in 2.6ms and consolidates 12 drifting scripts into one checked artifact [ADR-001].
3. **Proportionality Thresholds (`reviewLines: 50`, `reviewFiles: 5`):**
   - *Why leave it:* Solo operator engineering benefits significantly from not running 3 subagents on small 1-5 line edits [ADR-003].
4. **Developer preferences and stack constraints:**
   - Operator's preference for Go, TypeScript, GCP (Cloud Run, GKE), Colima/Docker, and macOS/Linux must remain strictly intact as instructed.

---

## 5. Goldfish Verification of the Plan

A fresh zero-context session given only `harness-plan.md` and `harness-inventory.md` can execute each numbered backlog item in order, knows the exact verification commands, and has the rollback plan for every step.
