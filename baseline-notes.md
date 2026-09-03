# Baseline Notes (Phase 1)

Date: **2026-09-03**  
Target Environment: **macOS (Apple Silicon), Go 1.24 toolchain, Claude Code & Google Antigravity**

---

## 1. Turn 0 Overhead Analysis

1. **Context Files (`shared/AGENTS.md`, `AGENTS.md`):**
   - Fixed constitution is 2,593 tokens; repo override is 455 tokens. Total = 3,048 tokens.
   - For Claude Code (200k window), this represents 1.5% of context.
   - For Gemini / Antigravity (1M–2M window), this represents <0.3% of context.
   - The constitution is tight (78 lines) and adheres to the "map, not manual" standard.

2. **Skills Metadata (28 Skills):**
   - 28 skills frontmatter (name + description) sum to **2,773 tokens**.
   - These tokens are loaded on turn 0 in every session by both CLIs to facilitate tool/skill routing.
   - Average metadata per skill: ~99 tokens, perfectly matching the Agent Skills recommendation (~100 tokens).

3. **Resume Protocol Cold Start Tax:**
   - Reading `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`, `DECISIONS.md`, and 3 newest changelog entries adds **~7,117 tokens**.
   - `DECISIONS.md` has grown to 179 lines (2,785 tokens) across 9 ADRs.
   - In `shared/AGENTS.md` line 9, the resume protocol specifies:
     *"`grep '^## ADR'` for older ones — never read the whole file"*, which helps agents avoid ingesting the entire ADR history when running `harness resume`.

---

## 2. Key Findings & Active Bug in Baseline

### Finding 1 (Active Bug): `CheckFreshness` Failure on Non-Source Commits
- **Location:** `internal/trail/check.go:58-70`
- **Symptom:** Running `go run ./cmd/harness check` in the repository immediately after HEAD commit `fe4f7b1` (which added `harness-improvement-prompt.md`) outputs:
  ```
  STALE: no changelog entry covers the current change set (HEAD 2026-09-03, newest entry 2026-09-03) — an entry must declare `scope: `
  RESULT: STALE — the documentation trail does not cover the code.
  ```
- **Root Cause:**
  - `p.Scope()` inspects `p.DirtySource()` and `p.HeadSourceToday()`.
  - When HEAD only touches non-source files (e.g., `.md` not in `sourceGlobs`), `files` is empty, so `p.Scope()` returns `""`.
  - In `check.go`:
    ```go
    scope := p.Scope()
    covered := scope != "" && p.HasArtifact(RowChangelog, scope)
    ```
    `covered` evaluates to `false` because `scope == ""`.
  - The switch statement then hits `default:` and fails, demanding an impossible empty scope string: `"an entry must declare scope: "`.
  - Contrast this with `internal/trail/owed.go`:
    ```go
    if len(rep.Changed) == 0 {
        return rep // Nothing owed!
    }
    ```
    `Owed()` correctly realizes that if no source files changed, nothing is owed. `CheckFreshness()` fails to make this check, creating a false-positive trail failure on doc-only commits.

### Finding 2: Missing Negative Boundaries in Skill Triggers
- Out of 28 skills, **13 skills lack negative triggers ("Do NOT use for...")**:
  `audit`, `bash`, `distill`, `find-skills`, `gcloud-auth-verification`, `gke-basics`, `grill-me`, `grilling`, `init-project`, `research`, `retro`, `skill-repair`, `terraform`.
- In empirical multi-vendor evaluations (OpenAI/Anthropic C2 study), negative trigger examples increased routing accuracy from **73% to 85%**. Without explicit negative boundaries, skills like `bash`, `research`, and `find-skills` frequently false-trigger on routine terminal operations.

### Finding 3: Specification Layout Inconsistencies in Skills
- **`anti-ui-slop`:** Uses non-standard `reference/` instead of `references/` (while also having `references/`!). Also carries 5 packaging and checksum files (`MANIFEST.json`, `CHECKSUMS.sha256`, `MODIFICATIONS.md`, `NOTICE`, `LICENSE`) and an unused `agents/openai.yaml`.
- **`codebase-design`:** Contains `DESIGN-IT-TWICE.md` and `DEEPENING.md` directly at the root of the skill folder rather than under `references/`.
- **`agents/openai.yaml`:** Present in 4 skills (`anti-ui-slop`, `codebase-design`, `grill-me`, `grilling`), completely dead on Claude Code and Antigravity.

### Finding 4: Subagent Model Aliasing in `internal/install/derive.go`
- In `derive.go`, `modelMap` maps:
  ```go
  var modelMap = map[string]string{"opus": "pro", "haiku": "flash"}
  ```
- Subagents with `model: sonnet` (such as `documentation`, `explorer`, `verifier`) are omitted from `antigravity/plugins/harness/agents/*.md`, which Antigravity treats as `model: inherit`.
- While `inherit` is valid, defining explicit mapping or documenting this behavior prevents ambiguity.

---

## 3. Performance & Verification Health

- `go test -race ./...` runs all packages cleanly in ~23.8 seconds.
- `go vet ./...` is clean (0 issues).
- `harness doctor --home /tmp/<scratch>` achieves 10/10 passing checks when installed to a test directory.
- `harness` compiled binary check executes in ~0.90s.
