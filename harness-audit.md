# Harness Audit (Phase 2)

Date: **2026-09-03**  
Audited Repository: **styrman-harness** (`/Users/spgo/Workspace/styrman-harness`)  
Framework: Model-Agnostic Rated Harness Improvement Specification

---

## Findings Summary

| ID | Area | Finding | Impact | Effort | Confidence | Freshness Grade |
|---|---|---|---|---|---|---|
| **AUD-01** | Architecture / Gates | `CheckFreshness` reports STALE when HEAD touches only non-source files (`scope: ` empty) | Quality / Latency (false gate failures) | S | High | [A3] |
| **AUD-02** | Skills Layout | Non-standard directory layouts in `anti-ui-slop` (`reference/`) and `codebase-design` (root `.md` files) | Portability / Standards Conformance | S | High | [A3] |
| **AUD-03** | Skills Triggers | 13 of 28 skills lack negative triggers (`Do NOT use for...`) | Quality / Cost (false skill activations) | M | High | [C2] |
| **AUD-04** | Architecture / Vendor | Vendor-specific packaging artifacts (`agents/openai.yaml`, checksums, manifests) in skills | Context Noise / Cleanliness | S | High | [A1] |
| **AUD-05** | Adapters / Subagents | Model tier mapping in `derive.go` omits `sonnet` (defaults implicitly to `inherit`) | Ambiguity / Predictability | S | High | [A3] |
| **AUD-06** | Token Economics | Fixed context cold-start loads ~7.1k tokens of trail history | Tokens / Cost (~7k tokens per resume) | M | Medium | [A2] |
| **AUD-07** | Subagent Prompts | Redundant bash boilerplate across subagent system prompts for docs root resolution | Context tokens (~300 tokens per subagent) | S | High | [A2] |
| **AUD-08** | Tooling / Subagents | Subagent `reviewer` and `security` have `Bash` access without restriction | Security / Agency Boundary | S | High | [A3] |

---

## Section-by-Section Checklists & Findings

### A. Architecture, Simplicity, Portability
- **[x] Structural Components Justified:**
  - The Go binary (`harness`) consolidates 12 shell/Python scripts into one binary, dropping hook latency from ~22.8ms to ~2.6ms [ADR-001].
  - Gates inform rather than block, except at commit boundary [ADR-002, R4].
- **[!] Finding AUD-01 (High Confidence): `CheckFreshness` fails on doc-only commits.**
  - **Evidence:** `internal/trail/check.go:58-70`. When `p.Scope()` is empty (because HEAD committed today touched only non-source files such as documentation), `covered` is false, and the check fails with:
    `"no changelog entry covers the current change set (HEAD %s, newest entry %s) — an entry must declare scope: "`
  - **Impact:** Blocks releases or reports false failure on `/preflight` when docs or prompts are updated without source code edits.
  - **Effort:** Small (add check for whether HEAD or dirty tree actually has source changes; if none, report OK).
  - **Freshness:** [A3].
- **[!] Finding AUD-04 (High Confidence): Extraneous vendor packaging files.**
  - **Evidence:** `shared/skills/anti-ui-slop/agents/openai.yaml`, `codebase-design/agents/openai.yaml`, `grill-me/agents/openai.yaml`, `grilling/agents/openai.yaml`, plus `anti-ui-slop/MANIFEST.json`, `CHECKSUMS.sha256`, `MODIFICATIONS.md`, `NOTICE`, `LICENSE`.
  - **Impact:** Clutters repository, leaks non-target vendor schemas into mirror directories (`~/.claude/skills`, `~/.gemini/config/skills`).
  - **Effort:** Small.
  - **Freshness:** [A1].

### B. Context Files (AGENTS.md and Friends)
- **[x] Single Source of Truth:**
  - `shared/AGENTS.md` is the canonical constitution. `claude/CLAUDE.md` and `antigravity/GEMINI.md` are exact generated copies [R1].
  - Root `AGENTS.md` provides project-specific overrides without duplicating the constitution.
- **[x] Map, Not Manual:**
  - `shared/AGENTS.md` is concise (78 lines), points to `.agents/docs/`, and avoids micro-managing syntax or style.
  - Adheres to operator constraints: Go, TypeScript, macOS/Linux, GCP only, Colima/Docker.

### C. Skills (Agent Skills Open Spec)
- **[!] Finding AUD-02 (High Confidence): Spec layout non-conformance in vendored skills.**
  - **Evidence:**
    - `shared/skills/anti-ui-slop`: Has `reference/` (singular) containing 6 markdown playbooks, while having `references/` containing `uizze-reference-policy.md`. The Agent Skills spec explicitly defines `references/` (plural).
    - `shared/skills/codebase-design`: Contains `DESIGN-IT-TWICE.md` and `DEEPENING.md` directly in the skill root folder. Spec requires non-`SKILL.md` documents to live in `references/`.
  - **Impact:** Skill discovery engines and validators fail or fail to locate sub-documents; violates progressive disclosure principles.
  - **Effort:** Small (move files to `references/` and update references in `SKILL.md`).
  - **Freshness:** [A3].
- **[!] Finding AUD-03 (High Confidence): Missing negative triggers ("Do NOT use for...") across 13 skills.**
  - **Evidence:** 13 skills (`audit`, `bash`, `distill`, `find-skills`, `gcloud-auth-verification`, `gke-basics`, `grill-me`, `grilling`, `init-project`, `research`, `retro`, `skill-repair`, `terraform`) have descriptions that specify what they are for, but lack explicit "Do NOT use for..." boundary statements.
  - **Impact:** Multi-model evaluations show negative examples increase routing accuracy from 73% to 85% [C2]. Without negative examples, skills like `bash` and `find-skills` trigger indiscriminately.
  - **Effort:** Medium (update descriptions across the 13 skills).
  - **Freshness:** [C2].

### D. Tools and MCP
- **[x] Tool Primitives Over Thin Wrappers:**
  - No bloated tool sets; standard shell (`run_command`/`Bash`) and file system primitives (`view_file`, `replace_file_content`, `write_to_file`) are used [C4].
- **[x] Zero MCP Turn 0 Overhead:**
  - MCP servers are omitted by default; zero token tax paid at turn 0 for MCP schema definitions.
- **[x] Strict Tool Mapping in `internal/install/derive.go`:**
  - Tool names are strictly validated and mapped to prevent subagent dispatch hangs on Antigravity.

### E. Token Economics & Caching
- **[x] Cache TTL Configuration:**
  - `claude/settings.json` specifies `"promptCacheTtl": "1h"` and `"subagentPromptCacheTtl": "1h"` to maximize KV-cache reuse.
- **[!] Finding AUD-06 (Medium Confidence): Trail resume context size.**
  - **Evidence:** Running the resume protocol cold start ingests ~7.1k tokens (`REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`, `DECISIONS.md`, 3 changes).
  - **Impact:** Increases cost and latency on cold starts.
  - **Mitigation:** The instruction in `shared/AGENTS.md` already advises `grep '^## ADR'` for DECISIONS.md rather than full-file ingestion. Can further optimize by ensuring change summaries are concise (<30 lines).
  - **Freshness:** [A2].

### F. Performance, Speed, Latency
- **[x] Fast Path Execution:**
  - Go hooks execute in ~2.6ms.
  - `harness check` takes ~0.9s (bounded by git subprocess calls).
- **[x] Proportionality Gate:**
  - Thresholds (`reviewLines: 50`, `reviewFiles: 5`) prevent launching expensive subagents on trivial 1-2 line edits [ADR-003].

### G. Long-Horizon and State
- **[x] Deterministic Artifact Trail:**
  - All state persists in `.agents/docs/` (`changes/`, `reviews/`, `testing/`, `plans/`, `handoffs/`).
  - Scopes identify artifacts by content digest (`changeset:<digest>`) spanning commit boundaries [ADR-007].
  - Changelog exception for dirty trees prevents task re-dispatch churn [ADR-009].

### H. Verification and Evals
- **[x] Comprehensive Automated Tests:**
  - Unit and integration tests in `internal/` cover gates, hookio, install, trail, and dispatch.
  - Dispatch contract tested via deterministic stub CLI in CI without requiring real vendor API credits [ADR-008].
  - `go test -race ./...` and `go vet ./...` pass 100%.

### I. Security and Containment
- **[x] Credential Staging Gate:**
  - `internal/gate/secret.go` checks git index and worktree directly rather than parsing shell arguments, defeating 15 bypass vectors [ADR-005].
- **[x] Formatter Isolation:**
  - PostToolUse formatter refuses binaries inside working repository or `node_modules` [ADR-006].
- **[!] Finding AUD-08 (High Confidence): Subagent tool allowlist precision.**
  - **Evidence:** `reviewer.md` and `security.md` declare `tools: Read, Grep, Glob, Bash, Write`.
  - While `reviewer` and `security` need Bash to run `git diff` or inspect commit history, their instructions say: *"You never modify source, config, or tests"*.
  - **Recommendation:** Document that the microVM/gate protects the environment and ensure prompt instructions explicitly remind agents that `Bash` is for read-only inspection.

### J. Observability & Maintenance Loop
- **[x] Ledger & Diagnostics:**
  - Dispatches logged to `.agents/docs/testing/dispatch-ledger.jsonl`.
  - `harness doctor` provides a 10-point check suite ensuring deployed files match the repository exactly.
  - `harness check` reports trail freshness.
