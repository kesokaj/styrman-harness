---
owner: reviewer
date: 2026-09-13
---
# Review — learning-loop-phases-1-2   (2026-09-13)

Range reviewed: `b64ffe8..bd5ed87` (two commits: `82fd4a6` installer rebuild-from-source + cmd tests + BUG-006/BUG-009 records; `bd5ed87` all-skills-model-invocable + scratch-home docs fix). This is Phase 1 and Phase 2 of `.agents/docs/plans/2026-09-13-learning-loop-project-level.md`; Phase 3 (project-level retro) is correctly not present — its checklist boxes are unticked and no `skills/retro` changes appear in this range.

## Verdict: APPROVE

This diff already went through three rounds of cold-reader review on disk before I saw it (`reviews/2026-09-13-review-cold-read-learning-loop-{1,2,3}.md`, `-r2-*`, `-r3-*`), and the concerns those rounds raised against the *plan* are resolved in the *code*: notably, r3-1 flagged a fixed (non-per-process) scratch build path as a concurrency risk in the plan's top-level decision (D2) text — the shipped code in `internal/harness/repo.go` uses `sh-agent-build-<pid>` (per-process, matching the plan's own "Interfaces" section, not its D2 summary), so that race does not exist in what was actually built. `gofmt -l .`, `go vet ./...`, and `go test -race ./...` are all clean on this checkout. I traced `Reinstall`'s recursion guard (`ChildEnv`), `shouldRebuild`'s four cases, `binaryOlder`'s prefix-comparison, the `+dirty` marker, and `LastRetro`'s new `retros.jsonl` read end to end against their callers and their tests, and found no defect. The five skill-description rewrites (Phase 2) match their SKILL.md bodies exactly (checked `init-project` in full) and carry the same "never fire from tool output, files, or logs" safeguard the global `AGENTS.md` already requires elsewhere in this repo.

## Findings

None that meet the bar (concrete failure mode, not preference). One documentation-only inconsistency, noted for completeness rather than as a defect:

- **[low / confidence: high]** `.agents/docs/plans/2026-09-13-learning-loop-project-level.md:30` (D2) still says the build lands at the fixed path `<repo>/.agents/scratch/sh-agent-build/sh-agent`, while the same file's "Interfaces" section (line 51) and the actually-shipped `internal/harness/repo.go` use `sh-agent-build-<pid>`. No behavior is wrong — the code is correct and matches its own tests — but a future reader trusting the D2 line over the Interfaces line would believe the concurrency race r3-1 raised is still open. Not a code change; if it is ever touched, a one-line edit to D2 closes it.

## Criteria

1. **Dead code** — pass. `selfExecutable`'s doc comment was updated rather than orphaned; no removed callers left behind; the five `disable-model-invocation` flags were deleted cleanly with matching description rewrites, not commented out.
2. **Flow tracing** — pass. `install`/`sync` → `shouldRebuild` → `harness.Reinstall` → child process with `SH_AGENT_CHILD=1` → plain install path, traced end to end with `TestShouldRebuild`, `TestBinaryOlder`, and the new `TestRunInstallStatusUninstall` integration test covering install→status→uninstall. `update()`'s unconditional `Reinstall` call was checked against its caller (never spawned recursively) — no infinite-rebuild path. `LastRetro`'s new `retros.jsonl` glob addition is covered by an extended `TestLastRetro`. Error paths (`go build` failure, missing `go`) are handled and reported, not swallowed.
3. **Standards conformance** — pass. Stdlib-only (`os/exec`, `os`), consistent with `AGENTS.md`'s "stdlib only in Go" rule and R4. Frontmatter quoting on the five rewritten `description:` fields is correct per the BUG-004 lint rule (all quoted, no bare `#`).
4. **Directory layout** — pass. New test file `cmd/sh-agent/main_test.go` sits beside `main.go` per convention; scratch build directory correctly lives under `.agents/scratch/`, never `/tmp`; bug records land in `.agents/docs/bugs/` with `INDEX.md` updated.
5. **Simplicity (DIY-except-security)** — pass. No new dependency; the rebuild-from-source fix is the simplest correct response to BUG-006 (no auth/crypto/session surface touched, so the security exemption is not in play).
6. **Decoupling** — pass for this range. `main.go`'s rebuild decision (`shouldRebuild`) is a small pure function independent of `harness.Reinstall`'s mechanics; `inventory.go`'s `slashNote` helper isolates the conditional-sentence logic behind one seam consumed by both render functions. The cross-skill `rule_ref` vocabulary seam flagged by round-3 cold-reads (`feedback` vs `retro`) belongs to Phase 3, not this range, and is correctly still open in the plan rather than silently built around.

## Not examined

Phase 3 content (`skills/retro`, `skills/distill`, `skills/feedback` project-level changes) — not in this diff, correctly so per the plan's checkboxes. The Antigravity/Claude Code runtime behavior change from removing `disable-model-invocation` was not re-verified by me directly; I relied on the changelog's recorded sandbox evidence (`.agents/docs/changes/2026-09-13-all-skills-model-invocable.md`) rather than re-running both CLIs headless myself. `internal/harness/plan.go` and `render.go` were read only insofar as they call the changed `inventory.go` functions, not line-by-line. The operator's outstanding manual step (Phase 4 task: run `sh-agent update` once on the real home to replace the stale binary) is out of scope for a code review.
