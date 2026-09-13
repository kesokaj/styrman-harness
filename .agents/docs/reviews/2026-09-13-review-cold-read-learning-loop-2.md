---
owner: reviewer
date: 2026-09-13
---
# Review — cold-read-learning-loop-2   (2026-09-13)

Range reviewed: cold read of `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` against the code it names (`cmd/sh-agent/main.go`, `internal/harness/{repo.go,apply.go,state.go,inventory.go,plan.go}`, `skills/{retro,distill,feedback,handoff,audit,init-project,research}/SKILL.md`). No diff exists yet; this is the pre-code cold-reader check the plan itself gates on (D11).

## Verdict: REQUEST CHANGES (design needs revision before Phase 1 starts)

The plan is unusually well specified (concrete signatures, failure modes, a verification matrix) and gets the hard factual claim right — `go run` really does skip VCS stamping while `go build` does not, verified directly (`go run ./cmd/sh-agent version` → `(devel)`; a built binary → `b64ffe886729+dirty`), so D2's premise holds. The one real defect is that the plan's own pseudocode for D1 contradicts its own interface signature for `installedCopy`, so two implementers would wire the install/sync branch differently. Everything else is a well-covered design; the remaining items are non-blocking defaults an implementer can pick alone.

## Findings

- **[high / confidence: high]** Plan `Interfaces` section, lines 48–49 vs `cmd/sh-agent/main.go:111–112,324–333` — `installedCopy(exe, home string) bool` is typed to take the running binary's **path** (compared via `filepath.EvalSymlinks`), but the sample call site writes `installedCopy(self, home)` where `self` is the pre-existing `selfExecutable()`'s return value, which is `[]byte` (file **contents**, read via `os.ReadFile` after resolving symlinks) — not a path. The plan never says whether `selfExecutable()` is restructured to also expose the resolved path, whether `os.Executable()` is called a second time just for the check, or whether the `install|sync` case is reordered so the path check happens before the byte-read. Two implementers would produce different call orders and possibly different fallback behavior when `go` is missing (the fallback needs the **bytes** again, from a function that up to now only had one call site). This must be settled before Phase 1's first task starts.

## Gaps (must be decided before implementation)

1. **`installedCopy` wiring** (see finding above): decide the exact signature boundary — does `selfExecutable()` return `(path string, data []byte, err error)`, or does `run()` call `os.Executable()` once and pass the string to both `installedCopy` and a separate byte-read helper? Cite: plan lines 48–49; `cmd/sh-agent/main.go:111-120,324-333`.
2. **Cleanup on `Reinstall` failure** (D2, plan line 30/50): "the build directory is removed afterwards" does not say whether removal happens when `go build` succeeds but the built binary's `install` fails, or only on the full-success path. A `defer os.RemoveAll` reads as the obvious choice but the plan doesn't say so, and the existing `Reinstall` in `internal/harness/repo.go:92-103` has no such cleanup today to model it on.
3. **D6 scope-classification mechanics** (plan line 34): "the same shape is already recorded for another project in `~/.config/sh-agent/feedback.jsonl` or the index" has no matching rule (same `rule_ref`? same `category`+`signal`? same target file?). This is prose read by a model rather than code, so it is lower stakes, but the plan should still name the field(s) compared or two runs of `/retro` on the same evidence could classify the same finding differently.

## Readiness — questions a first-pass implementer would have to ask

- Q1 (**BLOCKING**): "`installedCopy(self, home)` — is `self` the path or the bytes, and where in `run()` does the check sit relative to `selfExecutable()`?" Two implementers build genuinely different control flow without an answer (see finding above).
- Q2 (non-blocking): "Does `Reinstall`'s scratch build dir get removed via `defer` (cleaned up even on error) or only after a successful run?" A reasonable implementer defaults to `defer` and the plan doesn't forbid it.
- Q3 (non-blocking): "What exact fields decide D6's 'same shape... for another project' cross-project match?" A reasonable implementer picks `rule_ref` (already a field in `feedback.jsonl` per `skills/feedback/SKILL.md`) as the match key and documents the choice in the retro record's "Scope rule that fired" line, which the plan already requires.
- Q4 (non-blocking): "Exact line placement of the new `binary` status line relative to the existing `retro` line?" Plan says "after `installed`"; either position before or after the existing `retro` line reads fine and doesn't change behavior.
- Q5 (non-blocking): "Should `git rev-parse --short HEAD`'s variable abbreviation length ever mismatch `buildVersion()`'s 12-char slice and cause a `rev[:len(commit)]` panic?" `buildVersion()` returns a 12-hex-char string (or `(devel)`/`unknown`, both handled explicitly per D3), and `git rev-parse --short` defaults well under 12 chars for a repo this size, so this is currently safe; worth a bounds check but not a design fork.

## Not examined

- `internal/harness/plan.go`'s full `RenderContext`/`Split`/`Parse`/`Lint` implementations (read only the parts needed to confirm frontmatter pass-through to both vendors).
- The `skills/init-project/assets/agents-tree` full tree and `README.md`/`.agents/docs/deploy/local.md` current text (Phase 4 targets), not read in this pass.
- `agents/*.md` (`bugs`, `changelog`, `documentation`, `verifier`, `reviewer`, `security`) — not read; Phase 4's Definition-of-Done dispatch list assumes their current behavior needs no change, unverified here.
- Whether Claude Code and Antigravity actually change model-invocation behavior when `disable-model-invocation` is removed from a deployed `SKILL.md` (confirmed only that the harness copies the frontmatter verbatim to both vendor paths — `internal/harness/plan.go:118-127` — not that either CLI's runtime behaves as the plan assumes; the plan already treats this as a risk to verify in the sandbox, which is the right place for it).
- `TestInstallStatusResyncUninstall` and other existing tests' exact assertions, to confirm they don't hard-code the current `go run`-based `Reinstall` output text (a phrase match there could make Phase 1's `go build` switch break an existing test silently) — spot-checked test names only, not bodies.

## Criteria

The dispatch's task is a cold-read of a plan document, not a diff, so the six-criteria review lens is not the natural fit here; the findings above substitute. For completeness: no dead code introduced (nothing shipped yet); flow tracing is the substance of the one blocking finding; standards conformance and directory layout are addressed by the plan's own "Surfaces touched" and skill-line-count constraints; simplicity is respected (stdlib `os/exec`, no new dependency, `go build` over hand-rolled version tooling); decoupling looks sound (retro's project-vs-harness split keeps the two write paths — project docs root vs. harness queue — in separate code/prose paths per D5–D8).
