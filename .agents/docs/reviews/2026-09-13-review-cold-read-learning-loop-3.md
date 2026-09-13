---
owner: reviewer
date: 2026-09-13
---
# Review — cold-read-learning-loop-3   (2026-09-13)

Range reviewed: plan file `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` (cold-reader check, not a diff). Code surfaces read: `internal/harness/{repo.go,apply.go,state.go,inventory.go,plan.go}`, `cmd/sh-agent/main.go`, `skills/{retro,distill,feedback,handoff,audit,research,init-project}/SKILL.md`, `context/AGENTS.md`, `.agents/docs/bugs/INDEX.md`, `skills/init-project/assets/agents-tree/`.

## Verdict: design needs revision / implementation not ready

The Phase 1 (binary rebuild) design is solid and mostly traces cleanly against the existing code (`Reinstall` already exists in `repo.go` and just needs its `go run`→`go build` swap; the `install|sync` fallback-to-self-copy-on-missing-`go` path is explicitly decided at plan line 79). Phase 3 (project-level retro) is the weak part: the scope-classification rule (D6) has no operational definition, and the harness-scope "apply" step (D8) assumes a project-scoped session can reach into a separate repo checkout without saying how. Both are decisions two implementers would make differently, and both sit on the critical path of the headline feature ("if I run /retro in my project ... I want that saved on a project level").

## Gaps (must be decided before implementation)

- **D6's "same shape" test is undefined** (plan line 34: "the same shape is already recorded for another project in `~/.config/sh-agent/feedback.jsonl` or the index ... then it is `harness`"). There is no definition of "shape" — exact `category`+`rule_ref` match? free-text similarity of `user_quote`? same target file guessed from the finding? `skills/feedback/SKILL.md` entries carry `category`, `signal`, `rule_ref`, `user_quote` (lines ~14-31) — a plausible match key exists but the plan never picks one. Two implementers will classify the same finding differently (one project-scope, one harness-scope), which changes which files get edited and whether the harness-checkout apply step (D8) fires at all.
- **D8's harness-scope apply step doesn't say how a project session reaches the harness checkout** (plan line 36: "harness-scope → in the harness checkout, then `sh-agent install`, sandbox start of each CLI, changes entry, commit, push, handoff"). A `/retro` run happening inside an arbitrary project has its cwd (and, in Claude Code, plausibly its filesystem access) rooted in that project, not in the harness repo. `harness.FindRepo` (`internal/harness/repo.go:17-33`) walks cwd-and-parents, then `$SH_AGENT_REPO`, then the last-installed repo from `state.json` — none of these help a session whose sandbox can't `cd` there. The plan's own failure-mode row ("harness checkout not at hand: status `queued`", line 83) and D7's "additionally written to the existing queue file... unchanged" (line 35) together leave it unclear whether "queued" is ever anything but the permanent outcome for a project-level session — i.e., is D8's apply step reachable from a project context at all, or does every harness-scope proposal from a project always end up merely queued for a later harness-repo session? The plan needs to say which.
- **Minor, non-blocking but worth fixing before code:** the `status` older-binary comparison (plan line 52) compares `buildVersion()` "trimmed to the short-hash length" against the repo's `git rev-parse --short HEAD`. Git's short-hash length is not a fixed constant (it grows with repo size/collisions), so "trimmed to the short-hash length" needs `len(commit)` used as the slice bound at comparison time (which the prose implies but doesn't spell out) — an implementer who hardcodes 7 will occasionally false-positive/negative the warning.

## Readiness — questions a first-pass implementer would have to ask

1. **[BLOCKING]** What exact signal decides `Scope: project` vs `Scope: harness` for a finding (D6)? Without an operational rule, the retro's central output (`Scope:` tags, per D5's template) is implementer-defined, not plan-defined.
2. **[BLOCKING]** When a harness-scope proposal is produced from inside a project session, is the D8 "apply in the harness checkout, sandbox both CLIs, commit, push, handoff" sequence something the *same* session is expected to attempt (relying on `$SH_AGENT_REPO` or a discoverable checkout), or is it always deferred to a later session that opens the harness repo directly, with "queued" as the only real status a project-level retro can produce? This determines whether Phase 3's Definition-of-Done step is buildable as automation or is prose-only.
3. **[NON-BLOCKING, default reasonable]** `install|sync --dry-run` on the installed copy still triggers a real `go build` into `.agents/scratch/sh-agent-build/` before running the built binary with `--dry-run` passed through (plan line 49). The home directory is untouched either way, so "dry run" is preserved for its documented meaning; an implementer can reasonably decide this without asking.
4. **[NON-BLOCKING, default stated]** `go` missing during `install|sync` on the installed copy: plan line 79 explicitly states the fallback (print warning, self-copy as today) — decided, not a question.
5. **[NON-BLOCKING]** Ordering of the new `status` `binary` line relative to the existing `retro` line — plan says "after `installed`" (line 31/52); current code already emits `installed` then `retro` (`cmd/sh-agent/main.go:199-215`), so "after installed" is ambiguous by one line but inconsequential to behavior; an implementer picks either and it's fine.
6. **[NON-BLOCKING]** Concurrent `sh-agent install` runs racing on the same `<repo>/.agents/scratch/sh-agent-build/sh-agent` build path (two terminals) — not addressed, but rare enough that a reasonable implementer ignores it or lets the later build win.

## Not examined

Phase 4/5 docs-and-trail tasks (bug records, ADRs, changes entries) were not individually verified against `.agents/docs/bugs/INDEX.md` format beyond confirming BUG-002 exists and matches the reorder claim. `context/AGENTS.md` Learning-loop paragraph rewrite and the 80-line budget were not measured against a draft since none exists yet. The `retro`/`distill` SKILL.md 200-line cap and the proposed `references/record.md` split were not stress-tested against the full D5 template length. Antigravity-side behavior of `disable-model-invocation` removal (Phase 2) was not verified beyond reading `internal/harness/inventory.go`.

## Criteria

1. Dead code — n/a (plan, not a diff); the flag-supported-but-currently-unused `SlashOnly` path after Phase 2 (zero skills will carry the flag) is intentionally kept for future third-party skills per the plan — not dead code, a documented forward hook.
2. Flow tracing — mostly traces cleanly (`Reinstall`, `Install`, `main.go install/sync/update` call graph checked hop-to-hop); the one break is D8's cross-repo apply step (see gaps).
3. Standards conformance — n/a, no diff to check against language skills yet.
4. Directory layout — proposed `docs/retros/` and `assets/agents-tree/docs/retros/.gitkeep` additions match the existing `.agents/docs/` convention; no issue.
5. Simplicity — the rebuild-via-`go build` design is simpler than the alternative (embedding version info at compile time some other way) and stays stdlib-only; fine.
6. Decoupling — Phase 3's scope-classification logic (D6) is the one seam that isn't decoupled: it requires a project-level skill invocation to read and reason about global, cross-project state (`feedback.jsonl`, `retros.jsonl`) with an undefined comparison rule, which is exactly the kind of interface two implementers would diverge on.
